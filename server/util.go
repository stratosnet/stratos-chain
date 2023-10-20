package server

import (
	"path/filepath"

	tmcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	sdkservertypes "github.com/cosmos/cosmos-sdk/server/types"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/version"

	servercfg "github.com/stratosnet/stratos-chain/server/config"
)

// AddCommands add server commands
func AddCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator sdkservertypes.AppCreator,
	appExport sdkservertypes.AppExporter, addStartFlags sdkservertypes.ModuleInitFlags,
) {
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		sdkserver.ShowNodeIDCmd(),
		sdkserver.ShowValidatorCmd(),
		sdkserver.ShowAddressCmd(),
		sdkserver.VersionCmd(),
		tmcmd.ResetAllCmd,
		tmcmd.ResetStateCmd,
	)

	startCmd := StartCmd(appCreator, defaultNodeHome)
	addStartFlags(startCmd)

	rootCmd.AddCommand(
		startCmd,
		tendermintCmd,
		sdkserver.ExportCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
	)
}

// DefaultBaseAppOptions returns the default baseapp options provided by the Cosmos SDK
func DefaultBaseAppOptions(appOpts sdkservertypes.AppOptions) []func(*baseapp.BaseApp) {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(sdkserver.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	pruningOpts, err := sdkserver.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		// fallback to genesis chain-id
		appGenesis, err := tmtypes.GenesisDocFromFile(filepath.Join(homeDir, "config", "genesis.json"))
		if err != nil {
			panic(err)
		}

		chainID = appGenesis.ChainID
	}

	snapshotStore, err := sdkserver.GetSnapshotStore(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(
		cast.ToUint64(appOpts.Get(sdkserver.FlagStateSyncSnapshotInterval)),
		cast.ToUint32(appOpts.Get(sdkserver.FlagStateSyncSnapshotKeepRecent)),
	)

	return []func(*baseapp.BaseApp){
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(checkMinGasPrices(appOpts)),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(sdkserver.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(sdkserver.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(sdkserver.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(sdkserver.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(sdkserver.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(sdkserver.FlagIAVLCacheSize))),
		baseapp.SetIAVLDisableFastNode(cast.ToBool(appOpts.Get(sdkserver.FlagDisableIAVLFastNode))),
		baseapp.SetMempool(
			mempool.NewSenderNonceMempool(
				mempool.SenderNonceMaxTxOpt(cast.ToInt(appOpts.Get(sdkserver.FlagMempoolMaxTxs))),
			),
		),
		baseapp.SetIAVLLazyLoading(cast.ToBool(appOpts.Get(sdkserver.FlagIAVLLazyLoading))),
		baseapp.SetChainID(chainID),
	}
}

func checkMinGasPrices(appOpts sdkservertypes.AppOptions) string {
	minGasPricesInputStr := cast.ToString(appOpts.Get(sdkserver.FlagMinGasPrices))
	minGasPricesInput, err := sdk.ParseCoinNormalized(minGasPricesInputStr)
	if err != nil {
		panic(err)
	}

	minimalMinGasPricesStr := servercfg.GetMinimalMinGasPricesCoinStr()
	minimalMinGasPrices, err := sdk.ParseCoinNormalized(minimalMinGasPricesStr)
	if err != nil {
		panic(err)
	}

	if minGasPricesInput.IsLT(minimalMinGasPrices) {
		return minimalMinGasPricesStr
	}

	return minGasPricesInput.String()
}

//func ConnectTmWS(tmRPCAddr, tmEndpoint string, logger tmlog.Logger) *rpcclient.WSClient {
//	tmWsClient, err := rpcclient.NewWS(tmRPCAddr, tmEndpoint,
//		rpcclient.MaxReconnectAttempts(256),
//		rpcclient.ReadWait(120*time.Second),
//		rpcclient.WriteWait(120*time.Second),
//		rpcclient.PingPeriod(50*time.Second),
//		rpcclient.OnReconnect(func() {
//			logger.Debug("EVM RPC reconnects to Tendermint WS", "address", tmRPCAddr+tmEndpoint)
//		}),
//	)
//
//	if err != nil {
//		logger.Error(
//			"Tendermint WS client could not be created",
//			"address", tmRPCAddr+tmEndpoint,
//			"error", err,
//		)
//	} else if err := tmWsClient.OnStart(); err != nil {
//		logger.Error(
//			"Tendermint WS client could not start",
//			"address", tmRPCAddr+tmEndpoint,
//			"error", err,
//		)
//	}
//
//	return tmWsClient
//}
//
//func MountGRPCWebServices(
//	router *mux.Router,
//	grpcWeb *grpcweb.WrappedGrpcServer,
//	grpcResources []string,
//	logger tmlog.Logger,
//) {
//	for _, res := range grpcResources {
//
//		logger.Info("[GRPC Web] HTTP POST mounted", "resource", res)
//
//		s := router.Methods("POST").Subrouter()
//		s.HandleFunc(res, func(resp http.ResponseWriter, req *http.Request) {
//			if grpcWeb.IsGrpcWebSocketRequest(req) {
//				grpcWeb.HandleGrpcWebsocketRequest(resp, req)
//				return
//			}
//
//			if grpcWeb.IsGrpcWebRequest(req) {
//				grpcWeb.HandleGrpcWebRequest(resp, req)
//				return
//			}
//		})
//	}
//}
