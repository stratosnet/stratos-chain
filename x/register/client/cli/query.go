package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

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
		GetCmdQueryMetaNode(),
		GetCmdQueryResourceNodesCnt(),
		GetCmdQueryMetaNodesCnt(),
		GetCmdQueryParams(),
	)

	return registerQueryCmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current register parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as register parameters.

Example:
$ %s query register params
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

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryResourceNodesCnt implements the query total number of bonded resource nodes.
func GetCmdQueryResourceNodesCnt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-resource-count",
		Short: "Query the total number of bonded resource nodes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the total number of bonded resource nodes.
Example:
$ %s query register get-resource-count
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

			result, err := queryClient.BondedResourceNodeCount(cmd.Context(), &types.QueryBondedResourceNodeCountRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryMetaNodesCnt implements the query total number of bonded meta nodes.
func GetCmdQueryMetaNodesCnt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-meta-count",
		Short: "Query the total number of bonded meta nodes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the total number of bonded meta nodes.
Example:
$ %s query register get-meta-count
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

			result, err := queryClient.BondedMetaNodeCount(cmd.Context(), &types.QueryBondedMetaNodeCountRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryResourceNode implements the query resource nodes by network address command.
func GetCmdQueryResourceNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-resource-node",
		Short: "Query a resource node by its network address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual resource node by its network address.
Example:
$ %s query register get-resource-node --network-address=stsds1np4d8re98lpgrcdqcas8yt85gl3rvj268leg6v
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
			if len(queryFlagNetworkAddr) == 0 {
				return errors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
			}

			result, err := queryClient.ResourceNode(cmd.Context(), &types.QueryResourceNodeRequest{
				NetworkAddr: queryFlagNetworkAddr,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result.GetNode())
		},
	}

	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	//cmd.Flags().String(FlagNetworkAddress, "", "The network address of node")
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryMetaNode implements the query meta nodes by network address command.
func GetCmdQueryMetaNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-meta-node",
		Short: "Query an meta node by its network address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual meta node by its network address.
Example:
$ %s query register get-meta-node --network-address=stsds1faej5w4q6hgnt0ft598dlm408g4p747y4krwca
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
			if len(queryFlagNetworkAddr) == 0 {
				return errors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
			}

			result, err := queryClient.MetaNode(cmd.Context(), &types.QueryMetaNodeRequest{
				NetworkAddr: queryFlagNetworkAddr,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result.GetNode())
		},
	}

	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
