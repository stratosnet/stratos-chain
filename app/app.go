package app

import (
	"io"
	"os"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cast"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"

	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkruntime "github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdkservertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/stratosnet/stratos-chain/app/ante"
	"github.com/stratosnet/stratos-chain/app/upgrades"
	"github.com/stratosnet/stratos-chain/runtime"
	srvflags "github.com/stratosnet/stratos-chain/server/flags"
	"github.com/stratosnet/stratos-chain/x/evm"
	evmkeeper "github.com/stratosnet/stratos-chain/x/evm/keeper"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	potkeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	registerkeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	sdskeeper "github.com/stratosnet/stratos-chain/x/sds/keeper"
)

var (
	_ sdkruntime.AppI            = (*StratosApp)(nil)
	_ sdkservertypes.Application = (*StratosApp)(nil)
)

type EVMKeeperApp interface {
	GetEVMKeeper() *evmkeeper.Keeper
}

type StratosApp struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// Cosmos keepers
	accountKeeper    authkeeper.AccountKeeper
	bankKeeper       bankkeeper.Keeper
	capabilityKeeper *capabilitykeeper.Keeper
	stakingKeeper    *stakingkeeper.Keeper
	slashingKeeper   slashingkeeper.Keeper
	mintKeeper       mintkeeper.Keeper
	distrKeeper      distrkeeper.Keeper
	govKeeper        *govkeeper.Keeper
	crisisKeeper     *crisiskeeper.Keeper
	upgradeKeeper    *upgradekeeper.Keeper
	paramsKeeper     paramskeeper.Keeper
	authzKeeper      authzkeeper.Keeper
	evidenceKeeper   evidencekeeper.Keeper
	feeGrantKeeper   feegrantkeeper.Keeper
	consensusKeeper  consensuskeeper.Keeper

	// IBC keepers
	transferKeeper      ibctransferkeeper.Keeper
	ibcKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	ibcFeeKeeper        ibcfeekeeper.Keeper
	icaControllerKeeper icacontrollerkeeper.Keeper
	icaHostKeeper       icahostkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedIBCTransferKeeper   capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper

	// stratos keepers
	registerKeeper registerkeeper.Keeper
	potKeeper      potkeeper.Keeper
	sdsKeeper      sdskeeper.Keeper
	evmKeeper      *evmkeeper.Keeper

	// simulation manager
	sm *module.SimulationManager
}

func NewStratosApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	appOpts sdkservertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *StratosApp {
	var (
		app        = &StratosApp{}
		appBuilder *runtime.AppBuilder

		// merge the AppConfig and other configuration in one config
		appConfig = depinject.Configs(
			StratosAppConfig,
			depinject.Supply(
				// supply the application options
				appOpts,

				// ADVANCED CONFIGURATION

				//
				// AUTH
				//
				// For providing a custom function required in auth to generate custom account types
				// add it below. By default the auth module uses simulation.RandomGenesisAccounts.
				//
				// authtypes.RandomGenesisAccountsFn(simulation.RandomGenesisAccounts),

				// For providing a custom a base account type add it below.
				// By default the auth module uses authtypes.ProtoBaseAccount().
				//
				// func() authtypes.AccountI { return authtypes.ProtoBaseAccount() },

				//
				// MINT
				//

				// For providing a custom inflation function for x/mint add here your
				// custom function that implements the minttypes.InflationCalculationFn
				// interface.
			),
		)
	)

	//reset DefaultPowerReduction to prevent voting power overflow.
	sdk.DefaultPowerReduction = sdkmath.NewInt(1e12)

	if err := depinject.Inject(appConfig,
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.accountKeeper,
		&app.bankKeeper,
		&app.capabilityKeeper,
		&app.stakingKeeper,
		&app.slashingKeeper,
		&app.mintKeeper,
		&app.distrKeeper,
		&app.govKeeper,
		&app.crisisKeeper,
		&app.upgradeKeeper,
		&app.paramsKeeper,
		&app.authzKeeper,
		&app.evidenceKeeper,
		&app.feeGrantKeeper,
		&app.consensusKeeper,

		// Stratos keepers
		&app.registerKeeper,
		&app.potKeeper,
		&app.sdsKeeper,
		&app.evmKeeper,
	); err != nil {
		panic(err)
	}

	// Below we could construct and set an application specific mempool and
	// ABCI 1.0 PrepareProposal and ProcessProposal handlers. These defaults are
	// already set in the SDK's BaseApp, this shows an example of how to override
	// them.
	//
	// Example:
	//
	// app.App = appBuilder.Build(...)
	// nonceMempool := mempool.NewSenderNonceMempool()
	// abciPropHandler := NewDefaultProposalHandler(nonceMempool, app.App.BaseApp)
	//
	// app.App.BaseApp.SetMempool(nonceMempool)
	// app.App.BaseApp.SetPrepareProposal(abciPropHandler.PrepareProposalHandler())
	// app.App.BaseApp.SetProcessProposal(abciPropHandler.ProcessProposalHandler())
	//
	// Alternatively, you can construct BaseApp options, append those to
	// baseAppOptions and pass them to the appBuilder.
	//
	// Example:
	//
	// prepareOpt = func(app *baseapp.BaseApp) {
	// 	abciPropHandler := baseapp.NewDefaultProposalHandler(nonceMempool, app)
	// 	app.SetPrepareProposal(abciPropHandler.PrepareProposalHandler())
	// }
	// baseAppOptions = append(baseAppOptions, prepareOpt)
	app.App = appBuilder.Build(logger, db, traceStore, baseAppOptions...)

	app.evmKeeper.SetMsgServiceRouter(app.App.BaseApp.MsgServiceRouter())

	evmTracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))
	app.evmKeeper.SetTracer(evmTracer)

	// Add the EVM transient store key
	evmTransientKey := sdk.NewTransientStoreKey(evmtypes.TransientKey)
	app.evmKeeper.SetTransientKey(evmTransientKey)
	// Before v012 for legacy read
	app.evmKeeper.SetParamSpace(app.GetSubspace(evmtypes.ModuleName))
	err := app.RegisterStores(evmTransientKey)
	if err != nil {
		panic(err)
	}

	// set up non depinject support modules store keys
	storeKeys := sdk.NewKVStoreKeys(
		ibcexported.StoreKey, ibctransfertypes.StoreKey, ibcfeetypes.StoreKey,
		icahosttypes.StoreKey, icacontrollertypes.StoreKey,
	)
	for _, key := range storeKeys {
		err = app.RegisterStores(key)
		if err != nil {
			panic(err)
		}
	}

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(app.App.BaseApp, appOpts, app.appCodec, logger, app.kvStoreKeys()); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	/****  Module Options ****/

	// set params subspaces
	for _, m := range []string{ibctransfertypes.ModuleName, ibcexported.ModuleName, icahosttypes.SubModuleName, icacontrollertypes.SubModuleName} {
		app.paramsKeeper.Subspace(m)
	}

	// add capability keeper and ScopeToModule for ibc module
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedIBCTransferKeeper := app.capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAControllerKeeper := app.capabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAHostKeeper := app.capabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)

	// Create IBC keeper
	app.ibcKeeper = ibckeeper.NewKeeper(
		app.appCodec, app.GetKey(ibcexported.StoreKey), app.GetSubspace(ibcexported.ModuleName), app.stakingKeeper, app.upgradeKeeper, scopedIBCKeeper,
	)

	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://docs.cosmos.network/main/modules/gov#proposal-messages
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)). // This should be removed. It is still in place to avoid failures of modules that have not yet been upgraded.
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.ibcKeeper.ClientKeeper)).
		AddRoute(evmtypes.RouterKey, evm.NewEVMChangeProposalHandler(app.evmKeeper))

	app.ibcFeeKeeper = ibcfeekeeper.NewKeeper(
		app.appCodec, app.GetKey(ibcfeetypes.StoreKey),
		app.ibcKeeper.ChannelKeeper, // may be replaced with IBC middleware
		app.ibcKeeper.ChannelKeeper,
		&app.ibcKeeper.PortKeeper, app.accountKeeper, app.bankKeeper,
	)

	// Create IBC transfer keeper
	app.transferKeeper = ibctransferkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(ibctransfertypes.StoreKey),
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.ibcKeeper.ChannelKeeper,
		app.ibcKeeper.ChannelKeeper,
		&app.ibcKeeper.PortKeeper,
		app.accountKeeper,
		app.bankKeeper,
		scopedIBCTransferKeeper,
	)

	// Create interchain account keepers
	app.icaHostKeeper = icahostkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(icahosttypes.StoreKey),
		app.GetSubspace(icahosttypes.SubModuleName),
		app.ibcFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		app.ibcKeeper.ChannelKeeper,
		&app.ibcKeeper.PortKeeper,
		app.accountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)
	app.icaControllerKeeper = icacontrollerkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(icacontrollertypes.StoreKey),
		app.GetSubspace(icacontrollertypes.SubModuleName),
		app.ibcFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		app.ibcKeeper.ChannelKeeper,
		&app.ibcKeeper.PortKeeper,
		scopedICAControllerKeeper,
		app.MsgServiceRouter(),
	)

	// Set legacy router for backwards compatibility with gov v1beta1
	app.govKeeper.SetLegacyRouter(govRouter)

	// Create IBC modules with ibcfee middleware
	transferIBCModule := ibcfee.NewIBCMiddleware(ibctransfer.NewIBCModule(app.transferKeeper), app.ibcFeeKeeper)

	// integration point for custom authentication modules
	// see https://medium.com/the-interchain-foundation/ibc-go-v6-changes-to-interchain-accounts-and-how-it-impacts-your-chain-806c185300d7
	// TODO: Once a public ICA auth module is released, replace nil with the module
	var noAuthzModule porttypes.IBCModule
	icaControllerIBCModule := ibcfee.NewIBCMiddleware(
		icacontroller.NewIBCMiddleware(noAuthzModule, app.icaControllerKeeper),
		app.ibcFeeKeeper,
	)

	icaHostIBCModule := ibcfee.NewIBCMiddleware(icahost.NewIBCModule(app.icaHostKeeper), app.ibcFeeKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferIBCModule).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)

	app.ibcKeeper.SetRouter(ibcRouter)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	legacyModules := []module.AppModule{
		ibc.NewAppModule(app.ibcKeeper),
		ibctransfer.NewAppModule(app.transferKeeper),
		ibcfee.NewAppModule(app.ibcFeeKeeper),
		ica.NewAppModule(&app.icaControllerKeeper, &app.icaHostKeeper),
	}

	if err := app.RegisterModules(legacyModules...); err != nil {
		panic(err)
	}

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.accountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}

	// NOTE: Simulation manager, invariants and upgrade handlers must be after all the modules are registered
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()
	app.ModuleManager.RegisterInvariants(app.crisisKeeper)
	app.registerUpgradeHandlers()

	// register additional types
	ibctm.AppModuleBasic{}.RegisterInterfaces(app.interfaceRegistry)
	solomachine.AppModuleBasic{}.RegisterInterfaces(app.interfaceRegistry)

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedIBCTransferKeeper = scopedIBCTransferKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper
	app.ScopedICAControllerKeeper = scopedICAControllerKeeper

	// custom AnteHandler
	maxEthTxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	options := ante.HandlerOptions{
		AccountKeeper:     app.accountKeeper,
		BankKeeper:        app.bankKeeper,
		EvmKeeper:         app.evmKeeper,
		FeegrantKeeper:    app.feeGrantKeeper,
		IBCKeeper:         app.ibcKeeper,
		SignModeHandler:   app.txConfig.SignModeHandler(),
		SigGasConsumer:    ante.DefaultSigVerificationGasConsumer,
		MaxEthTxGasWanted: maxEthTxGasWanted,
		TxFeeChecker:      ante.CheckTxFeeWithValidatorMinGasPrices,
	}
	if err := options.Validate(); err != nil {
		panic(err)
	}
	app.SetAnteHandler(ante.NewAnteHandler(options))

	// In v0.46, the SDK introduces _postHandlers_. PostHandlers are like
	// antehandlers, but are run _after_ the `runMsgs` execution. They are also
	// defined as a chain, and have the same signature as antehandlers.
	//
	// In baseapp, postHandlers are run in the same store branch as `runMsgs`,
	// meaning that both `runMsgs` and `postHandler` state will be committed if
	// both are successful, and both will be reverted if any of the two fails.
	//
	// The SDK exposes a default postHandlers chain, which comprises of only
	// one decorator: the Transaction Tips decorator. However, some chains do
	// not need it by default, so feel free to comment the next line if you do
	// not need tips.
	// To read more about tips:
	// https://docs.cosmos.network/main/core/tips.html
	//
	// Please note that changing any of the anteHandler or postHandler chain is
	// likely to be a state-machine breaking change, which needs a coordinated
	// upgrade.
	postHandler, err := NewPostHandler()
	if err != nil {
		panic(err)
	}
	app.SetPostHandler(postHandler)

	// A custom InitChainer can be set if extra pre-init-genesis logic is required.
	// By default, when using app wiring enabled module, this is not required.
	// For instance, the upgrade module will set automatically the module version map in its init genesis thanks to app wiring.
	// However, when registering a module manually (i.e. that does not support app wiring), the module version map
	// must be set manually as follow. The upgrade module will de-duplicate the module version map.
	app.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		app.upgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())
		return app.App.InitChainer(ctx, req)
	})

	if err = app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}

// Name returns the name of the App
func (app *StratosApp) Name() string { return app.BaseApp.Name() }

// LegacyAmino returns StratosApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *StratosApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// SimulationManager implements the SimulationApp interface
func (app *StratosApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *StratosApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
	// register swagger API in app.go so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}

// registerUpgrade registers the given upgrade to be supported by the app
func (app *StratosApp) registerUpgrade(upgrade upgrades.Upgrade) {
	app.upgradeKeeper.SetUpgradeHandler(upgrade.Name(), upgrade.Handler())

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgrade.Name() && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// Configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, upgrade.StoreUpgrades()))
	}
}

// AppCodec returns StratosApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *StratosApp) AppCodec() codec.Codec {
	return app.appCodec
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *StratosApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *StratosApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	sk := app.UnsafeFindStoreKey(storeKey)
	kvStoreKey, ok := sk.(*storetypes.KVStoreKey)
	if !ok {
		return nil
	}
	return kvStoreKey
}

// kvStoreKeys returns all the kv store keys registered inside StratosApp
func (app *StratosApp) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}

func (app *StratosApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.accountKeeper
}

func (app *StratosApp) GetBankKeeper() bankkeeper.Keeper {
	return app.bankKeeper
}

func (app *StratosApp) GetStakingKeeper() *stakingkeeper.Keeper {
	return app.stakingKeeper
}

func (app *StratosApp) GetRegisterKeeper() registerkeeper.Keeper {
	return app.registerKeeper
}

func (app *StratosApp) GetPotKeeper() potkeeper.Keeper {
	return app.potKeeper
}

func (app *StratosApp) GetDistrKeeper() distrkeeper.Keeper {
	return app.distrKeeper
}

func (app *StratosApp) GetEVMKeeper() *evmkeeper.Keeper {
	return app.evmKeeper
}

// TxConfig returns StratosApp's txConfig
func (app *StratosApp) TxConfig() client.TxConfig {
	return app.txConfig
}
