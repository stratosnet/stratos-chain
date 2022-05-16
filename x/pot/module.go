package pot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/stratosnet/stratos-chain/x/pot/client/cli"
	"github.com/stratosnet/stratos-chain/x/pot/client/rest"
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the pot module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the pot module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the pot module's types for the given codec.
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

// ValidateGenesis performs genesis state validation for the register module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return types.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the pot module.
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the pot module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the pot module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the pot module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

//____________________________________________________________________________

// AppModule implements an application module for the pot module.
type AppModule struct {
	AppModuleBasic
	keeper     keeper.Keeper
	bankKeeper types.BankKeeper
	//supplyKeeper   supply.Keeper
	accountKeeper  types.AccountKeeper
	stakingKeeper  types.StakingKeeper
	registerKeeper types.RegisterKeeper
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, &genesisState)

	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(sdk.Context, codec.JSONCodec) json.RawMessage {
	panic("implement me")
}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

func (am AppModule) RegisterServices(module.Configurator) {
	panic("implement me")
}

func (am AppModule) ConsensusVersion() uint64 {
	panic("implement me")
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper, stakingKeeper types.StakingKeeper, registerKeeper types.RegisterKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		bankKeeper:     bankKeeper,
		//supplyKeeper:   supplyKeeper,
		accountKeeper:  accountKeeper,
		stakingKeeper:  stakingKeeper,
		registerKeeper: registerKeeper,
	}
}

// Name returns the pot module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the pot module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

//// Route returns the message routing key for the pot module.
//func (AppModule) Route() string {
//	return types.RouterKey
//}

// NewHandler returns an sdk.Handler for the pot module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the pot module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the register module sdk.Querier.
func (am AppModule) NewQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

//// InitGenesis performs genesis initialization for the pot module. It returns
//// no validator updates.
//func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
//	var genesisState types.GenesisState
//	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
//	InitGenesis(ctx, am.keeper, genesisState)
//	return []abci.ValidatorUpdate{}
//}
//
//// ExportGenesis returns the exported genesis state as raw bytes for the pot
//// module.
//func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
//	gs := ExportGenesis(ctx, am.keeper)
//	return types.ModuleCdc.MustMarshalJSON(gs)
//}

// BeginBlock returns the begin blocker for the pot module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the pot module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
