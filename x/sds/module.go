package sds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	modulev1 "github.com/stratosnet/stratos-chain/api/stratos/sds/module/v1"
	potkeeper "github.com/stratosnet/stratos-chain/x/pot/keeper"
	registerkeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/sds/client/cli"
	"github.com/stratosnet/stratos-chain/x/sds/keeper"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

const (
	consensusVersion = 1
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
	_ appmodule.AppModule        = AppModule{}
	_ depinject.OnePerModuleType = AppModule{}
)

// AppModuleBasic defines the basic application module used by the sds module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the sds module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the sds module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the register
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the sds module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return types.ValidateGenesis(data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the sds module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root tx command for the sds module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the sds module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

//____________________________________________________________________________

// AppModule implements an application module for the sds module.
type AppModule struct {
	AppModuleBasic
	keeper         keeper.Keeper
	bankKeeper     types.BankKeeper
	registerKeeper types.RegisterKeeper
	potKeeper      types.PotKeeper

	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace types.ParamsSubspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, k keeper.Keeper, bankKeeper types.BankKeeper, registerKeeper types.RegisterKeeper,
	potKeeper types.PotKeeper, legacySubspace types.ParamsSubspace,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         k,
		bankKeeper:     bankKeeper,
		registerKeeper: registerKeeper,
		potKeeper:      potKeeper,
		legacySubspace: legacySubspace,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	querier := keeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), querier)

	//m := keeper.NewMigrator(am.keeper)
	//_ = cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
}

// RegisterInvariants registers the sds module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs genesis initialization for the sds module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the sds
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(&gs)
}

// BeginBlock returns the beginBlocker for the sds module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the endBlocker for the sds module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, req, am.keeper)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 {
	return consensusVersion
}

// IsOnePerModuleType implements depinject.OnePerModuleType
func (am AppModule) IsOnePerModuleType() {
}

// IsAppModule implements appmodule.AppModule
func (am AppModule) IsAppModule() {
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the profiles module.
type AppModuleSimulation struct{}

// GenerateGenesisState implements AppModuleSimulation
func (am AppModule) GenerateGenesisState(input *module.SimulationState) {
}

// RegisterStoreDecoder implements AppModuleSimulation
func (am AppModule) RegisterStoreDecoder(registry sdk.StoreDecoderRegistry) {
}

// WeightedOperations implements AppModuleSimulation
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// App Wiring Setup

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(
			ProvideModule,
		),
	)
}

type ModuleInputs struct {
	depinject.In

	Config         *modulev1.Module
	Cdc            codec.Codec
	Key            *storetypes.KVStoreKey
	BankKeeper     bankkeeper.Keeper
	RegisterKeeper registerkeeper.Keeper
	PotKeeper      potkeeper.Keeper

	// LegacySubspace is used solely for migration of x/params managed parameters
	LegacySubspace types.ParamsSubspace `optional:"true"`
}

type ModuleOutputs struct {
	depinject.Out

	SdsKeeper keeper.Keeper
	Module    appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.BankKeeper,
		in.RegisterKeeper,
		in.PotKeeper,
		authority.String(),
	)

	m := NewAppModule(
		in.Cdc,
		k,
		in.BankKeeper,
		in.RegisterKeeper,
		in.PotKeeper,
		in.LegacySubspace,
	)

	return ModuleOutputs{SdsKeeper: k, Module: m}
}
