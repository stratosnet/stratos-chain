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
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	regkeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	regtestutil "github.com/stratosnet/stratos-chain/x/register/testutil"
	regtypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdskeeper "github.com/stratosnet/stratos-chain/x/sds/keeper"
	sdstestutil "github.com/stratosnet/stratos-chain/x/sds/testutil"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

var (
	sdsAcct = authtypes.NewModuleAddress(sdstypes.ModuleName)
	regAcct = authtypes.NewModuleAddress(regtypes.ModuleName)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	key            *storetypes.KVStoreKey
	encCfg         moduletestutil.TestEncodingConfig
	sdsKeeper      sdskeeper.Keeper
	bankKeeper     *sdstestutil.MockBankKeeper
	registerKeeper *sdstestutil.MockRegisterKeeper
	potKeeper      *sdstestutil.MockPotKeeper

	accountKeeper *regtestutil.MockAccountKeeper

	queryClient sdstypes.QueryClient
	msgServer   sdstypes.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(sdstypes.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, sdk.NewTransientStoreKey("transient_test"))
	// creating context
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	// registering interfaces
	encCfg := moduletestutil.MakeTestEncodingConfig()
	sdstypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	regtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	pottypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	authtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	distrtypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	bankKeeper := sdstestutil.NewMockBankKeeper(ctrl)
	registerKeeper := sdstestutil.NewMockRegisterKeeper(ctrl)
	potKeeper := sdstestutil.NewMockPotKeeper(ctrl)

	sdsKeeper := sdskeeper.NewKeeper(encCfg.Codec, key, bankKeeper, registerKeeper, potKeeper, sdsAcct.String())

	s.ctx = ctx
	s.key = key
	s.encCfg = encCfg
	s.sdsKeeper = sdsKeeper
	s.bankKeeper = bankKeeper
	s.registerKeeper = registerKeeper
	s.potKeeper = potKeeper
	s.accountKeeper = regtestutil.NewMockAccountKeeper(ctrl)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	sdstypes.RegisterQueryServer(queryHelper, sdskeeper.Querier{Keeper: sdsKeeper})

	s.queryClient = sdstypes.NewQueryClient(queryHelper)
	s.msgServer = sdskeeper.NewMsgServerImpl(sdsKeeper)
}

func (s *KeeperTestSuite) TestParams() {
	ctx, keeper := s.ctx, s.sdsKeeper
	require := s.Require()

	expParams := sdstypes.DefaultParams()
	keeper.SetParams(ctx, expParams)
	resParams := keeper.GetParams(ctx)
	require.True(expParams.Equal(resParams))
}

func (s *KeeperTestSuite) GetNoMockRegKeeper() regkeeper.Keeper {
	ctrl := gomock.NewController(s.T())
	return regkeeper.NewKeeper(
		s.encCfg.Codec, s.key,
		s.accountKeeper, regtestutil.NewMockBankKeeper(ctrl), regtestutil.NewMockDistrKeeper(ctrl),
		regAcct.String(),
	)
}

// Init used before tests to set initial data
func (s *KeeperTestSuite) Reset(t *testing.T) {
	t.Helper()
	s.ctx = s.ctx.WithEventManager(sdk.NewEventManager())
	s.sdsKeeper.SetParams(s.ctx, sdstypes.DefaultParams())
}

func (s *KeeperTestSuite) mockBankSendCoinsFromAccountToModule(sender sdk.Address, receipientModule string, amount sdk.Coins) *gomock.Call {
	return s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(s.ctx, sender, receipientModule, amount)
}

func (s *KeeperTestSuite) mockBankHasBalance(sender sdk.AccAddress, amount sdk.Coin) *gomock.Call {
	return s.bankKeeper.EXPECT().HasBalance(s.ctx, sender, amount)
}

func (s *KeeperTestSuite) mockRegCalculatePurchaseAmount(amount sdkmath.Int) *gomock.Call {
	return s.registerKeeper.EXPECT().CalculatePurchaseAmount(s.ctx, amount)
}

func (s *KeeperTestSuite) mockRegSetRemainingOzoneLimit(amount sdkmath.Int) *gomock.Call {
	return s.registerKeeper.EXPECT().SetRemainingOzoneLimit(s.ctx, amount)
}

func (s *KeeperTestSuite) mockRegGenerateMerkleProofs(signer sdk.AccAddress, data []byte) *gomock.Call {
	return s.registerKeeper.EXPECT().GenerateMerkleProofs(s.ctx, signer, data)
}

func (s *KeeperTestSuite) mockAccGetAccount(acc sdk.AccAddress) *gomock.Call {
	return s.accountKeeper.EXPECT().GetAccount(s.ctx, acc)
}
