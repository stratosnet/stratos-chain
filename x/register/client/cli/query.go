package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// GetQueryCmd returns the cli query commands for register module
func GetQueryCmd() *cobra.Command {
	// Group register queries under a subcommand
	registerQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	registerQueryCmd.AddCommand(
		GetCmdQueryResourceNode(),
		GetCmdQueryIndexingNode(),
		//GetCmdQueryIndexingNodeList(),
	)

	return registerQueryCmd
}

// GetCmdQueryResourceNode implements the query resource nodes by network address command.
func GetCmdQueryResourceNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-resource-nodes [flag]",
		Short: "Query a resource node by its network address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual resource node by its network address.
Example:
$ %s query register get-resource-nodes --network-id=%sstsds1np4d8re98lpgrcdqcas8yt85gl3rvj268leg6v
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// query resource node by network address
			queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
			if queryFlagNetworkAddr == "" {
				return sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
			}

			result, err := queryClient.ResourceNode(cmd.Context(), &types.QueryResourceNodeRequest{
				// Leaving status empty on purpose to query all validators.
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}
	cmd.Flags().String(FlagNetworkAddress, "", "(optional) The network address of the node")
	return cmd
}

// GetCmdQueryIndexingNode implements the query indexing nodes by network address command.
func GetCmdQueryIndexingNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-indexing-nodes [flag]",
		Short: "Query an indexing node by its network address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual indexing node by its network address.
Example:
$ %s query register get-indexing-nodes --network-id=%sstsds1faej5w4q6hgnt0ft598dlm408g4p747y4krwca
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// query resource node by network address
			queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
			if queryFlagNetworkAddr == "" {
				return sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
			}

			result, err := queryClient.IndexingNode(cmd.Context(), &types.QueryIndexingNodeRequest{
				// Leaving status empty on purpose to query all validators.
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}
	cmd.Flags().String(FlagNetworkAddress, "", "(optional) The network address of the node")
	return cmd
}

//
//// GetResNodesByNetworkAddr queries all resource nodes by multiple network IDs (sep: ";")
//func GetResNodesByNetworkAddr(cliCtx context.CLIContext, queryRoute string) (res string, err error) {
//	queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
//	queryByFlagNetworkAddrList := strings.Split(queryFlagNetworkAddr, ";")
//	for _, v := range queryByFlagNetworkAddrList {
//		resp, _, err := QueryResourceNode(cliCtx, queryRoute, v)
//		if err != nil {
//			return "null", err
//		}
//		res += string(resp) + ";"
//	}
//	return res[:len(res)-1], nil
//}
//
//// QueryResourceNode queries resource node by network addr
//func QueryResourceNode(cliCtx context.CLIContext, queryRoute, networkAddr string) ([]byte, int64, error) {
//	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryResourceNodeByNetworkAddr)
//	sdsAddress, err := stratos.SdsAddressFromBech32(networkAddr)
//	if err != nil {
//		return []byte{}, 0, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
//	}
//
//	params := types.NewQueryNodesParams(1, 1, sdsAddress, "", nil)
//	bz, err := cliCtx.Codec.MarshalJSON(params)
//	if err != nil {
//		return []byte{}, 0, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
//	}
//	return cliCtx.QueryWithData(route, bz)
//}
//
//// GetCmdQueryIndexingNodeList implements the query all indexing nodes by network id command.
//func GetCmdQueryIndexingNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "get-indexing-nodes [flags]", // []byte
//		Short: "Query all indexing nodes",
//		Long: strings.TrimSpace(
//			fmt.Sprintf(`Query all indexing nodes`),
//		),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			inBuf := bufio.NewReader(cmd.InOrStdin())
//			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
//
//			// query all indexing nodes by network address
//			queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
//			if queryFlagNetworkAddr == "" {
//				return sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
//			}
//			resp, err := GetIndNodesByNetworkAddr(cliCtx, queryRoute)
//			if err != nil {
//				return err
//			}
//			return cliCtx.PrintOutput(resp)
//
//		},
//	}
//	cmd.Flags().String(FlagNetworkAddress, "", "(optional) The network address of the node")
//
//	return cmd
//}
//
//// GetIndNodesByNetworkAddr queries all indexing nodes by multiple network addrs (sep: ";")
//func GetIndNodesByNetworkAddr(cliCtx context.CLIContext, queryRoute string) (res string, err error) {
//	queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
//	queryByFlagNetworkAddrList := strings.Split(queryFlagNetworkAddr, ";")
//	for _, v := range queryByFlagNetworkAddrList {
//		resp, _, err := QueryIndexingNodes(cliCtx, queryRoute, v)
//		if err != nil {
//			return "null", err
//		}
//		res += string(resp) + ";"
//	}
//	return res[:len(res)-1], nil
//}
//
//// QueryIndexingNodes queries all indexing nodes
//func QueryIndexingNodes(cliCtx context.CLIContext, queryRoute, networkAddr string) ([]byte, int64, error) {
//	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodeByNetworkAddr)
//	sdsAddress, err := stratos.SdsAddressFromBech32(networkAddr)
//	if err != nil {
//		return []byte{}, 0, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
//	}
//
//	params := types.NewQueryNodesParams(1, 1, sdsAddress, "", nil)
//	bz, err := cliCtx.Codec.MarshalJSON(params)
//	if err != nil {
//		return []byte{}, 0, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
//	}
//	return cliCtx.QueryWithData(route, bz)
//}
//
//// Route returns the message routing key for the staking module.
//func (am AppModule) Route() sdk.Route {
//	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
//}
//
//// QuerierRoute returns the staking module's querier route name.
//func (AppModule) QuerierRoute() string {
//	return types.QuerierRoute
//}
//
//// LegacyQuerierHandler returns the staking module sdk.Querier.
//func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
//	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
//}
