package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
	"strings"
)

// GetQueryCmd returns the cli query commands for register module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group register queries under a subcommand
	registerQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	registerQueryCmd.AddCommand(
		flags.GetCommands(
			// this line is used by starport scaffolding # 1
			GetCmdQueryResourceNodeList(queryRoute, cdc),
			GetCmdQueryIndexingNodeList(queryRoute, cdc),
			GetCmdQueryNetworkSet(queryRoute, cdc),
		)...,
	)

	return registerQueryCmd
}

// GetCmdQueryResourceNodeList implements the query all resource nodes by network address command.
func GetCmdQueryResourceNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get-resource-node-list", // []byte
		//Args:  cobra.RangeArgs(1, 1),
		Short: "Query all resource nodes by network address.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all resource nodes by network address.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			resp, err := queryByNetworkAddr(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(string(resp[:len(resp)-1]))
		},
	}
	cmd.Flags().AddFlagSet(FsNetworkAddr)
	_ = cmd.MarkFlagRequired(FlagNetworkAddr)

	return cmd
}

// query all resource nodes by multiple network addresses (sep: ";")
func queryByNetworkAddr(cliCtx context.CLIContext, queryRoute string) ([]byte, error) {
	queryByFlagNetworkAddr := viper.GetString(FlagNetworkAddr)
	queryByFlagNetworkAddrList := strings.Split(queryByFlagNetworkAddr, ";")
	var res []byte
	for _, v := range queryByFlagNetworkAddrList {
		resp, _, err := QueryResourceNodes(cliCtx, queryRoute, v)
		if err != nil {
			return nil, err
		}
		res = append(res, resp...)
		res = append(res, ';')
	}
	return res, nil
}

// QueryResourceNodes queries all resource nodes by network address
func QueryResourceNodes(cliCtx context.CLIContext, queryRoute, networkAddress string) ([]byte, int64, error) {

	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryResourceNodeList)
	return cliCtx.QueryWithData(route, []byte(networkAddress))
}

// GetCmdQueryIndexingNodeList implements the query all indexing nodes by network address command.
func GetCmdQueryIndexingNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get-indexing-node-list", // []byte
		//Args:  cobra.RangeArgs(1, 1),
		Short: "Query all indexing nodes by network address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all indexing nodes by network address.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			// query all indexing nodes by network address
			resp, _, err := QueryIndexingNodes(cliCtx, queryRoute, viper.GetString(FlagNetworkAddr))
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(string(resp))
		},
	}
	cmd.Flags().AddFlagSet(FsNetworkAddr)
	_ = cmd.MarkFlagRequired(FlagNetworkAddr)

	return cmd
}

// QueryIndexingNodes queries all indexing nodes by network address
func QueryIndexingNodes(cliCtx context.CLIContext, queryRoute, networkAddress string) ([]byte, int64, error) {

	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodeList)
	return cliCtx.QueryWithData(route, []byte(networkAddress))
}

// GetCmdQueryNetworkSet implements the query all indexing nodes by network address command.
func GetCmdQueryNetworkSet(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "get-network-set",
		//Args:  cobra.RangeArgs(0, 0),
		Short: "Query all network addresses",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all network addresses.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query get-network-set by network address
			resp, _, err := QueryNetworkSet(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(string(resp))
		},
	}
}

// QueryNetworkSet queries all network address
func QueryNetworkSet(cliCtx context.CLIContext, queryRoute string) ([]byte, int64, error) {

	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryNetworkSet)
	return cliCtx.Query(route)
}

// HasMoniker - returns whether or not the node has an moniker(name)
//func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
//	store := ctx.KVStore(k.ownersStoreKey)
//	bz := store.Get([]byte(name))
//	return bz != nil
//}
