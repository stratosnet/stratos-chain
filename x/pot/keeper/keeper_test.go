package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
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
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stratosnet/stratos-chain/crypto/merkle"
	sttypes "github.com/stratosnet/stratos-chain/types"
	potkeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	pottestutil "github.com/stratosnet/stratos-chain/x/pot/testutil"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	regkeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	regtestutil "github.com/stratosnet/stratos-chain/x/register/testutil"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
)

var (
	potAcct  = authtypes.NewModuleAddress(pottypes.ModuleName)
	regAcct  = authtypes.NewModuleAddress(regtypes.ModuleName)
	prepAcct = authtypes.NewModuleAddress(regtypes.TotalUnissuedPrepay)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	key            *storetypes.KVStoreKey
	encCfg         moduletestutil.TestEncodingConfig
	potKeeper      potkeeper.Keeper
	accountKeeper  *pottestutil.MockAccountKeeper
	bankKeeper     *pottestutil.MockBankKeeper
	distrKeeper    *pottestutil.MockDistrKeeper
	registerKeeper *pottestutil.MockRegisterKeeper
	stakingKeeper  *pottestutil.MockStakingKeeper
	queryClient    pottypes.QueryClient
	msgServer      pottypes.MsgServer
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
	pottypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	authtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	distrtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	stakingtypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	accountKeeper := pottestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := pottestutil.NewMockBankKeeper(ctrl)
	distrKeeper := pottestutil.NewMockDistrKeeper(ctrl)
	registerKeeper := pottestutil.NewMockRegisterKeeper(ctrl)
	stakingKeeper := pottestutil.NewMockStakingKeeper(ctrl)

	potKeeper := potkeeper.NewKeeper(encCfg.Codec, key, accountKeeper, bankKeeper, distrKeeper, registerKeeper, stakingKeeper, potAcct.String())

	s.ctx = ctx
	s.encCfg = encCfg
	s.key = key
	s.potKeeper = potKeeper
	s.accountKeeper = accountKeeper
	s.bankKeeper = bankKeeper
	s.distrKeeper = distrKeeper
	s.registerKeeper = registerKeeper
	s.stakingKeeper = stakingKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	pottypes.RegisterQueryServer(queryHelper, potkeeper.Querier{Keeper: potKeeper})

	s.queryClient = pottypes.NewQueryClient(queryHelper)
	s.msgServer = potkeeper.NewMsgServerImpl(potKeeper)
}

func (s *KeeperTestSuite) GetNoMockRegKeeper() regkeeper.Keeper {
	ctrl := gomock.NewController(s.T())
	return regkeeper.NewKeeper(
		s.encCfg.Codec, s.key,
		s.accountKeeper, regtestutil.NewMockBankKeeper(ctrl), s.distrKeeper,
		regAcct.String(),
	)
}

func (s *KeeperTestSuite) TestParams() {
	ctx, keeper := s.ctx, s.potKeeper
	require := s.Require()

	expParams := pottypes.DefaultParams()
	keeper.SetParams(ctx, expParams)
	resParams := keeper.GetParams(ctx)
	require.True(expParams.Equal(resParams))
}

// Init used before tests to set initial data
func (s *KeeperTestSuite) Reset(t *testing.T) {
	t.Helper()
	s.ctx = s.ctx.WithEventManager(sdk.NewEventManager())
	s.potKeeper.SetParams(s.ctx, pottypes.DefaultParams())
	s.potKeeper.InitVariable(s.ctx)
}

func (s *KeeperTestSuite) mockRegOwnMetaNode(ownerAddr sdk.AccAddress, p2pAddr sttypes.SdsAddress) *gomock.Call {
	return s.registerKeeper.EXPECT().OwnMetaNode(s.ctx, ownerAddr, p2pAddr)
}

func (s *KeeperTestSuite) mockRegGetBondedMetaNodeCnt() *gomock.Call {
	return s.registerKeeper.EXPECT().GetBondedMetaNodeCnt(s.ctx)
}

func (s *KeeperTestSuite) mockRegGetRemainingOzoneLimit() *gomock.Call {
	return s.registerKeeper.EXPECT().GetRemainingOzoneLimit(s.ctx)
}

func (s *KeeperTestSuite) mockRegSetRemainingOzoneLimit(amount sdkmath.Int) *gomock.Call {
	return s.registerKeeper.EXPECT().SetRemainingOzoneLimit(s.ctx, amount)
}

func (s *KeeperTestSuite) mockRegGetDepositNozRate() *gomock.Call {
	return s.registerKeeper.EXPECT().GetDepositNozRate(s.ctx)
}

func (s *KeeperTestSuite) mockRegGetEffectiveTotalDeposit() *gomock.Call {
	return s.registerKeeper.EXPECT().GetEffectiveTotalDeposit(s.ctx)
}

func (s *KeeperTestSuite) mockRegGetTotalUnissuedPrepay() *gomock.Call {
	return s.registerKeeper.EXPECT().GetTotalUnissuedPrepay(s.ctx)
}

func (s *KeeperTestSuite) mockRegGetMetaNodeIterator() *gomock.Call {
	return s.registerKeeper.EXPECT().GetMetaNodeIterator(s.ctx)
}

func (s *KeeperTestSuite) mockBankSendCoinsFromModuleToModule(msender string, mreceipient string, amount sdk.Coins) *gomock.Call {
	return s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(s.ctx, msender, mreceipient, amount)
}

func (s *KeeperTestSuite) mockAccGetModuleAddress(name string) *gomock.Call {
	return s.accountKeeper.EXPECT().GetModuleAddress(name)
}

func (s *KeeperTestSuite) mockDistrFundCommunityPool(amount sdk.Coins, sender sdk.AccAddress) *gomock.Call {
	return s.distrKeeper.EXPECT().FundCommunityPool(s.ctx, amount, sender)
}

func (s *KeeperTestSuite) mockRegProcessMerkleProofs(mdata merkle.MerkleProofData) *gomock.Call {
	return s.registerKeeper.EXPECT().ProcessMerkleProofs(s.ctx, mdata)
}
