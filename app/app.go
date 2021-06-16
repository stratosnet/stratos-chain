package app

import (
	"encoding/json"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	"github.com/stratosnet/stratos-chain/x/pot"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stratosnet/stratos-chain/x/sds"
	// this line is used by starport scaffolding # 1
)

const (
	appName    = "sds"
	appVersion = "v0.39"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.stratoschaincli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.stratoschaind")
	ModuleBasics    = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		register.AppModuleBasic{},
		pot.AppModuleBasic{},
		//gov.AppModuleBasic{},
		distr.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler,
		),
		sds.AppModuleBasic{},
		// this line is used by starport scaffolding # 2
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName: nil,
		// this line is used by starport scaffolding # 2.1
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
	// module accounts that are allowed to receive tokens
	//allowedReceivingModAcc = map[string]bool{
	//	distr.ModuleName: true,
	//}
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc.Seal()
}

type NewApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	subspaces map[string]params.Subspace

	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	stakingKeeper  staking.Keeper
	supplyKeeper   supply.Keeper
	paramsKeeper   params.Keeper
	govKeeper      gov.Keeper
	upgradeKeeper  upgrade.Keeper
	distrKeeper    distr.Keeper
	registerKeeper register.Keeper
	potKeeper      pot.Keeper
	sdsKeeper      sds.Keeper
	// this line is used by starport scaffolding # 3
	mm *module.Manager

	sm *module.SimulationManager
}

var _ simapp.App = (*NewApp)(nil)

func NewInitApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *NewApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		auth.StoreKey,
		staking.StoreKey,
		supply.StoreKey,
		params.StoreKey,
		gov.StoreKey,
		upgrade.StoreKey,
		register.StoreKey,
		pot.StoreKey,
		sds.StoreKey,
		// this line is used by starport scaffolding # 5
	)

	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	var app = &NewApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[pot.ModuleName] = app.paramsKeeper.Subspace(pot.DefaultParamSpace)
	app.subspaces[register.ModuleName] = app.paramsKeeper.Subspace(register.DefaultParamSpace)
	// this line is used by starport scaffolding # 5.1

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.supplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	// this line is used by starport scaffolding # 5.2

	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
		// this line is used by starport scaffolding # 5.3
		),
	)

	app.upgradeKeeper = upgrade.NewKeeper(
		map[int64]bool{},
		keys[upgrade.StoreKey],
		app.cdc,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper))
	app.govKeeper = gov.NewKeeper(
		app.cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.supplyKeeper,
		&stakingKeeper, govRouter,
	)

	app.registerKeeper = register.NewKeeper(
		app.cdc,
		keys[register.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		app.subspaces[register.ModuleName],
	)

	app.potKeeper = pot.NewKeeper(
		app.cdc,
		keys[pot.StoreKey],
		app.subspaces[pot.ModuleName],
		auth.FeeCollectorName,
		app.bankKeeper,
		app.supplyKeeper,
		app.accountKeeper,
		app.stakingKeeper,
		app.registerKeeper,
	)

	app.sdsKeeper = sds.NewKeeper(
		app.cdc,
		keys[sds.StoreKey],
		app.bankKeeper,
		app.registerKeeper,
		app.potKeeper,
	)

	// this line is used by starport scaffolding # 4

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		register.NewAppModule(app.registerKeeper, app.accountKeeper, app.bankKeeper),
		pot.NewAppModule(app.potKeeper, app.bankKeeper, app.supplyKeeper, app.accountKeeper, app.stakingKeeper, app.registerKeeper),
		sds.NewAppModule(app.sdsKeeper, app.bankKeeper, app.registerKeeper),
		// this line is used by starport scaffolding # 6
	)

	app.mm.SetOrderEndBlockers(
		staking.ModuleName,
		gov.ModuleName,
		// this line is used by starport scaffolding # 6.1
	)

	app.mm.SetOrderInitGenesis(
		// this line is used by starport scaffolding # 6.2
		auth.ModuleName,
		staking.ModuleName,
		bank.ModuleName,
		gov.ModuleName,
		supply.ModuleName,
		register.ModuleName,
		sds.ModuleName,
		pot.ModuleName,
		genutil.ModuleName,
		upgrade.ModuleName,
		// this line is used by starport scaffolding # 7
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.upgradeKeeper.SetUpgradeHandler(appVersion, func(ctx sdk.Context, plan upgrade.Plan) {
		logger.Info("Upgrade Handler working")
	})
	app.SetStoreLoader(bam.StoreLoaderWithUpgrade(&store.StoreUpgrades{
		Renamed: []store.StoreRename{{
			//OldKey: "foo",
			//NewKey: "bar",
		}},
	}),
	)

	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *NewApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState

	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *NewApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *NewApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *NewApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

func (app *NewApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *NewApp) Codec() *codec.Codec {
	return app.cdc
}

func (app *NewApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}
