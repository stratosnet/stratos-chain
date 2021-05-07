package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"strings"

	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	// "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/register/types"
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
		)...,
	)

	return registerQueryCmd
}

// GetCmdQueryResourceNodeList implements the query all resource nodes by network address command.
func GetCmdQueryResourceNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-resource-node-list [network_address]", // []byte
		Args:  cobra.RangeArgs(1, 1),
		Short: "Query all resource nodes by network address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all resource nodes by network address.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query all resource nodes by network address
			resp, _, err := QueryResourceNodes(cliCtx, queryRoute, args[0])
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(string(resp))
		},
	}
}

// QueryResourceNodes queries all resource nodes by network address
func QueryResourceNodes(cliCtx context.CLIContext, queryRoute, networkAddress string) ([]byte, int64, error) {

	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryResourceNodeList)
	return cliCtx.QueryWithData(route, []byte(networkAddress))
}

// GetCmdQueryIndexingNodeList implements the query all indexing nodes by network address command.
func GetCmdQueryIndexingNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-indexing-node-list [network_address]", // []byte
		Args:  cobra.RangeArgs(1, 1),
		Short: "Query all indexing nodes by network address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all indexing nodes by network address.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query all indexing nodes by network address
			resp, _, err := QueryIndexingNodes(cliCtx, queryRoute, args[0])
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(string(resp))
		},
	}
}

// QueryIndexingNodes queries all indexing nodes by network address
func QueryIndexingNodes(cliCtx context.CLIContext, queryRoute, networkAddress string) ([]byte, int64, error) {

	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodeList)
	return cliCtx.QueryWithData(route, []byte(networkAddress))
}
