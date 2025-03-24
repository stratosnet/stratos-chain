package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/golang/mock/gomock"

	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stratosnet/stratos-chain/crypto/merkle"
	regkeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	regtestutil "github.com/stratosnet/stratos-chain/x/register/testutil"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

var (
	regAcct = authtypes.NewModuleAddress(regtypes.ModuleName)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	key            *storetypes.KVStoreKey
	registerKeeper regkeeper.Keeper
	accountKeeper  *regtestutil.MockAccountKeeper
	bankKeeper     *regtestutil.MockBankKeeper
	distrKeeper    *regtestutil.MockDistrKeeper
	queryClient    regtypes.QueryClient
	msgServer      regtypes.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(regtypes.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, sdk.NewTransientStoreKey("transient_test"))
	// creating context
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	// registering interfaces
	encCfg := moduletestutil.MakeTestEncodingConfig()
	regtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	authtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	distrtypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	accountKeeper := regtestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := regtestutil.NewMockBankKeeper(ctrl)
	distrKeeper := regtestutil.NewMockDistrKeeper(ctrl)

	registerKeeper := regkeeper.NewKeeper(encCfg.Codec, key, accountKeeper, bankKeeper, distrKeeper, regAcct.String())

	s.ctx = ctx
	s.key = key
	s.registerKeeper = registerKeeper
	s.accountKeeper = accountKeeper
	s.bankKeeper = bankKeeper
	s.distrKeeper = distrKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	regtypes.RegisterQueryServer(queryHelper, regkeeper.Querier{Keeper: registerKeeper})

	s.queryClient = regtypes.NewQueryClient(queryHelper)
	s.msgServer = regkeeper.NewMsgServerImpl(registerKeeper)
}

func (s *KeeperTestSuite) TestParams() {
	ctx, keeper := s.ctx, s.registerKeeper
	require := s.Require()

	expParams := regtypes.DefaultParams()
	keeper.SetParams(ctx, expParams)
	resParams := keeper.GetParams(ctx)
	require.True(expParams.Equal(resParams))
}

// Init used before tests to set initial data
func (s *KeeperTestSuite) Reset(t *testing.T) {
	t.Helper()
	s.registerKeeper.SetParams(s.ctx, regtypes.DefaultParams())
	s.ctx = s.ctx.WithEventManager(sdk.NewEventManager())
}

func (s *KeeperTestSuite) mockAccGetAccount(acc sdk.AccAddress) *gomock.Call {
	return s.accountKeeper.EXPECT().GetAccount(s.ctx, acc)
}

func (s *KeeperTestSuite) TestQuerierMerkleRoot() {
	require := s.Require()

	// explicitly set so in case of changes, this will work
	proover := merkle.NewRelayerMerkleProver()
	s.registerKeeper.SetProover(proover)

	s.T().Run("test querier merkle root call", func(t *testing.T) {
		s.Reset(t)

		ctx, keeper, queryClient := s.ctx, s.registerKeeper, s.queryClient
		signer := sdk.AccAddress([]byte("signer1_______________"))
		data := []byte("test")

		acc := authtypes.NewBaseAccount(signer, nil, 1, 1)
		s.mockAccGetAccount(signer).Return(acc).AnyTimes()

		commitment1 := merkle.CreateSdkCommitment(signer, acc.GetSequence(), data)

		// empty
		req := &regtypes.QueryMerkleRootRequest{}

		resp, err := queryClient.MerkleRoot(ctx, req)
		require.NoError(err)

		require.Equal(merkle.NullCommitment[:], common.Hex2Bytes(resp.Root))
		require.Equal("", resp.Commitment)

		// with some not existing commitment
		req = &regtypes.QueryMerkleRootRequest{
			Commitment: common.Bytes2Hex(commitment1),
		}

		resp, err = queryClient.MerkleRoot(ctx, req)
		require.NoError(err)

		require.Equal(merkle.NullCommitment[:], common.Hex2Bytes(resp.Root))
		require.Equal(commitment1, common.Hex2Bytes(resp.Commitment))

		// after record
		err = keeper.RecordMerkleCommitment(ctx, signer, data)
		require.NoError(err)
		root1 := keeper.GetMerkleRootForCommitment(ctx, commitment1)
		require.Equal(root1, keeper.GetMerkleRoot(ctx))

		req = &regtypes.QueryMerkleRootRequest{
			Commitment: common.Bytes2Hex(commitment1),
		}

		resp, err = queryClient.MerkleRoot(ctx, req)
		require.NoError(err)

		require.Equal(keeper.GetMerkleRoot(ctx), common.Hex2Bytes(resp.Root))
		require.Equal(commitment1, common.Hex2Bytes(resp.Commitment))

		// after nullify
		newRoot := proover.GetRoot(keeper.GetMerkleRoot(ctx), [][]byte{commitment1})
		mdata := merkle.NewMerkleProofBundle(newRoot, [][]byte{commitment1})

		err = keeper.ProcessMerkleProofs(ctx, mdata)
		require.NoError(err)

		req = &regtypes.QueryMerkleRootRequest{
			Commitment: common.Bytes2Hex(commitment1),
		}

		resp, err = queryClient.MerkleRoot(ctx, req)
		require.NoError(err)

		require.Equal(merkle.NullCommitment[:], common.Hex2Bytes(resp.Root))
		require.Equal(commitment1, common.Hex2Bytes(resp.Commitment))
	})
}

func (s *KeeperTestSuite) TestGenerateMerkleProofs() {
	require := s.Require()

	s.T().Run("test base relayer merkle proof generation", func(t *testing.T) {
		s.Reset(t)

		ctx, keeper := s.ctx, s.registerKeeper

		signer := sdk.AccAddress([]byte("signer1_______________"))
		data := []byte("test")

		oldRoot := keeper.GetMerkleRoot(ctx)

		require.Equal(oldRoot, merkle.NullCommitment[:])

		acc := authtypes.NewBaseAccount(signer, nil, 1, 1)

		s.mockAccGetAccount(signer).Return(acc).AnyTimes()

		// explicitly set so in case of changes, this will work
		s.registerKeeper.SetProover(merkle.NewRelayerMerkleProver())

		expCommitment := merkle.CreateSdkCommitment(signer, acc.GetSequence(), data)

		err := keeper.RecordMerkleCommitment(ctx, signer, data)
		require.NoError(err)

		expRoot := keeper.GetMerkleRoot(ctx)

		require.Equal(oldRoot, expRoot)

		require.Equal(1, len(ctx.EventManager().ABCIEvents()))

		for _, evt := range ctx.EventManager().ABCIEvents() {
			msg, _ := sdk.ParseTypedEvent(evt)
			switch revt := msg.(type) {
			case *regtypes.EventMerkleDataUpdated:
				require.Equal(expRoot, revt.Root)
				require.Equal(expCommitment, revt.Commitment)
			}
		}
	})
}

func (s *KeeperTestSuite) TestAckMerkleLeaves() {
	require := s.Require()

	// explicitly set so in case of changes, this will work
	proover := merkle.NewRelayerMerkleProver()
	s.registerKeeper.SetProover(proover)

	s.T().Run("nullify one leaf", func(t *testing.T) {
		s.Reset(t)

		ctx, keeper := s.ctx, s.registerKeeper

		commitment1 := []byte("c1")

		root := keeper.GetMerkleRoot(ctx)

		rlProof, _ := proover.CreateProofs(
			[][]byte{
				root,
			},
			[][]byte{
				commitment1,
			})

		err := keeper.CreateMerkleCommitment(ctx, commitment1, root)
		require.NoError(err)

		_, err = keeper.NullifyMerkleCommitments(ctx, rlProof.GetCommitments())
		require.NoError(err)
		require.Equal(1, len(ctx.EventManager().ABCIEvents()))

		res := keeper.GetMerkleRootForCommitment(ctx, commitment1)
		require.Equal(merkle.NullCommitment[:], res)

		for _, evt := range ctx.EventManager().ABCIEvents() {
			msg, _ := sdk.ParseTypedEvent(evt)
			switch revt := msg.(type) {
			case *regtypes.EventMerkleDataUpdated:
				require.Equal(root, revt.Root)
				require.Equal(commitment1, revt.Commitment)
				require.Equal(regtypes.EventMerkleDataUpdated_NULLIFY, revt.ActionType)
			}
		}
	})

	s.T().Run("nullify same commitment and fail", func(t *testing.T) {
		s.Reset(t)

		ctx, keeper := s.ctx, s.registerKeeper

		commitment1 := merkle.NullCommitment[:]

		root := keeper.GetMerkleRoot(ctx)

		rlProof, _ := proover.CreateProofs(
			[][]byte{
				root,
			},
			[][]byte{
				commitment1,
			})

		_, err := keeper.NullifyMerkleCommitments(ctx, rlProof.GetCommitments())
		require.Error(err)

		res := ctx.KVStore(s.key).Get(regtypes.GetMerkleCommitmentKey(commitment1))
		require.Equal(0, len(res))
		require.Equal(0, len(ctx.EventManager().ABCIEvents()))
	})

	s.T().Run("nullify multiple leaves", func(t *testing.T) {
		s.Reset(t)

		ctx, keeper := s.ctx, s.registerKeeper

		commitment1 := []byte("c4")
		commitment2 := merkle.NullCommitment[:]
		commitment3 := []byte("c5")

		root := keeper.GetMerkleRoot(ctx)

		rlProof, _ := proover.CreateProofs(
			[][]byte{
				root,
				root,
				root,
			},
			[][]byte{
				commitment1,
				commitment2,
				commitment3,
			})

		keeper.CreateMerkleCommitment(ctx, commitment1, root)
		keeper.CreateMerkleCommitment(ctx, commitment3, root)

		_, err := keeper.NullifyMerkleCommitments(ctx, rlProof.GetCommitments())

		require.NoError(err)

		res1 := keeper.GetMerkleRootForCommitment(ctx, commitment1)
		require.Equal(merkle.NullCommitment[:], res1)
		res2 := keeper.GetMerkleRootForCommitment(ctx, commitment2)
		require.Equal(merkle.NullCommitment[:], res2)
		res3 := keeper.GetMerkleRootForCommitment(ctx, commitment3)
		require.Equal(merkle.NullCommitment[:], res3)

		require.Equal(2, len(ctx.EventManager().ABCIEvents()))

		evt1 := ctx.EventManager().ABCIEvents()[0]
		msg1, _ := sdk.ParseTypedEvent(evt1)
		revt1 := msg1.(*regtypes.EventMerkleDataUpdated)
		require.Equal(root, revt1.Root)
		require.Equal(commitment1, revt1.Commitment)
		require.Equal(regtypes.EventMerkleDataUpdated_NULLIFY, revt1.ActionType)

		evt3 := ctx.EventManager().ABCIEvents()[1]
		msg3, _ := sdk.ParseTypedEvent(evt3)
		revt3 := msg3.(*regtypes.EventMerkleDataUpdated)
		require.Equal(root, revt3.Root)
		require.Equal(commitment3, revt3.Commitment)
		require.Equal(regtypes.EventMerkleDataUpdated_NULLIFY, revt3.ActionType)
	})
}
