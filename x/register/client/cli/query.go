package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stratosnet/stratos-chain/x/register/keeper"
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
			GetCmdQueryResourceNode(queryRoute, cdc),
			GetCmdQueryIndexingNodeList(queryRoute, cdc),
		)...,
	)

	return registerQueryCmd
}

// GetCmdQueryResourceNodeList implements the query all resource nodes by network id command.
func GetCmdQueryResourceNode(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-resource-node [flags]", // []byte
		Short: "Query resource node by network-id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query resource node by network-id`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			// query resource node by network ID
			queryFlagNetworkID := viper.GetString(FlagNetworkID)
			if queryFlagNetworkID == "" {
				return sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network-id")
			}
			resp, err := GetResNodeByNetworkID(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)
		},
	}
	cmd.Flags().String(FlagNetworkID, "", "(optional) The network id of the node")
	return cmd
}

// GetResNodeByNetworkID queries all resource nodes by multiple network IDs (sep: ";")
func GetResNodeByNetworkID(cliCtx context.CLIContext, queryRoute string) (res string, err error) {
	queryFlagNetworkID := viper.GetString(FlagNetworkID)
	queryByFlagNetworkIDList := strings.Split(queryFlagNetworkID, ";")
	for _, v := range queryByFlagNetworkIDList {
		resp, _, err := QueryResourceNode(cliCtx, queryRoute, v)
		if err != nil {
			return "null", err
		}
		res += string(resp) + ";"
	}
	return res[:len(res)-1], nil
}

// QueryResourceNode queries resource node by network addr
func QueryResourceNode(cliCtx context.CLIContext, queryRoute, networkID string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryResourceNodeByNetworkAddr)
	return cliCtx.QueryWithData(route, []byte(networkID))
}

// GetCmdQueryIndexingNodeList implements the query all indexing nodes by network id command.
func GetCmdQueryIndexingNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-indexing-nodes [flags]", // []byte
		Short: "Query indexing node by network-id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all indexing nodes by network-id`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			// query all indexing nodes by network-id
			queryFlagNetworkID := viper.GetString(FlagNetworkID)
			if queryFlagNetworkID == "" {
				return sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network-id")
			}
			resp, err := GetIndNodesByNetworkID(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)

		},
	}
	cmd.Flags().String(FlagNetworkID, "", "(optional) The network id of the node")

	return cmd
}

// GetIndNodesByNetworkID queries all indexing nodes by multiple network IDs
func GetIndNodesByNetworkID(cliCtx context.CLIContext, queryRoute string) (res string, err error) {
	queryFlagNetworkID := viper.GetString(FlagNetworkID)
	queryByFlagNetworkIDList := strings.Split(queryFlagNetworkID, ";")
	for _, v := range queryByFlagNetworkIDList {
		resp, _, err := QueryIndexingNodes(cliCtx, queryRoute, v)
		if err != nil {
			return "null", err
		}
		res += string(resp) + ";"
	}
	return res[:len(res)-1], nil
}

// QueryIndexingNodes queries one indexing node by network ID
func QueryIndexingNodes(cliCtx context.CLIContext, queryRoute, networkID string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodeByNetworkAddr)
	return cliCtx.QueryWithData(route, []byte(networkID))
}
