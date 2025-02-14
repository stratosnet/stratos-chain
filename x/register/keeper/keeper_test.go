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

		s.mockAccGetAccount(signer).Return(acc)

		// explicitly set so in case of changes, this will work
		s.registerKeeper.SetProover(merkle.NewRelayerMerkleProver())

		expRoot := common.Hex2Bytes("f74e8557b482dc4d487178364a6c9bb05a21f41410f6d880c9136a1045b0a93d")
		expCommitment := common.Hex2Bytes("20b0225e68a64f45ace8916cde4d7410ad50b7e2ac53c2c42c46c6afffc6425c")

		err := keeper.GenerateMerkleProofs(ctx, signer, data)
		require.NoError(err)

		newRoot := keeper.GetMerkleRoot(ctx)

		require.NotEqual(oldRoot, newRoot)
		require.Equal(expRoot, newRoot)

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

	s.T().Run("ack one leaf", func(t *testing.T) {
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

		err := keeper.AckMerkleLeaves(ctx, rlProof.GetLeaves())

		require.NoError(err)

		res := ctx.KVStore(s.key).Get(regtypes.GetMerkleCommitmentKey(commitment1))

		require.Equal(1, int(res[0]))

		require.Equal(1, len(ctx.EventManager().ABCIEvents()))

		for _, evt := range ctx.EventManager().ABCIEvents() {
			msg, _ := sdk.ParseTypedEvent(evt)
			switch revt := msg.(type) {
			case *regtypes.EventCommitmentAcknowledged:
				require.Equal(root, revt.Root)
				require.Equal(commitment1, revt.Commitment)
			}
		}
	})

	s.T().Run("odd leaves", func(t *testing.T) {
		s.Reset(t)

		ctx, keeper := s.ctx, s.registerKeeper

		commitment1 := []byte("c1")
		commitment2 := []byte("c2")

		root := keeper.GetMerkleRoot(ctx)

		err := keeper.AckMerkleLeaves(ctx, [][]byte{
			root, commitment1, commitment2,
		})

		require.Error(err)
		require.ErrorContains(err, "even")
		require.Equal(0, len(ctx.EventManager().ABCIEvents()))
	})

	s.T().Run("null commitment skip", func(t *testing.T) {
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

		err := keeper.AckMerkleLeaves(ctx, rlProof.GetLeaves())

		require.NoError(err)

		res := ctx.KVStore(s.key).Get(regtypes.GetMerkleCommitmentKey(commitment1))
		require.Equal(0, len(res))
		require.Equal(0, len(ctx.EventManager().ABCIEvents()))
	})

	s.T().Run("ack multiple leaves", func(t *testing.T) {
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

		err := keeper.AckMerkleLeaves(ctx, rlProof.GetLeaves())

		require.NoError(err)

		res1 := ctx.KVStore(s.key).Get(regtypes.GetMerkleCommitmentKey(commitment1))
		require.Equal(1, int(res1[0]))
		res2 := ctx.KVStore(s.key).Get(regtypes.GetMerkleCommitmentKey(commitment2))
		require.Equal(0, len(res2))
		res3 := ctx.KVStore(s.key).Get(regtypes.GetMerkleCommitmentKey(commitment3))
		require.Equal(1, int(res3[0]))

		require.Equal(2, len(ctx.EventManager().ABCIEvents()))

		evt1 := ctx.EventManager().ABCIEvents()[0]
		msg1, _ := sdk.ParseTypedEvent(evt1)
		revt1 := msg1.(*regtypes.EventCommitmentAcknowledged)
		require.Equal(root, revt1.Root)
		require.Equal(commitment1, revt1.Commitment)

		evt3 := ctx.EventManager().ABCIEvents()[1]
		msg3, _ := sdk.ParseTypedEvent(evt3)
		revt3 := msg3.(*regtypes.EventCommitmentAcknowledged)
		require.Equal(root, revt3.Root)
		require.Equal(commitment3, revt3.Commitment)
	})
}
