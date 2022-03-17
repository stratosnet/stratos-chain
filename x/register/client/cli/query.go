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
	stratos "github.com/stratosnet/stratos-chain/types"
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

// GetCmdQueryResourceNode implements the query resource nodes by network address command.
func GetCmdQueryResourceNode(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-resource-nodes [flags]", // []byte
		Short: "Query resource node by network address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query resource node by network address`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			// query resource node by network address
			queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
			if queryFlagNetworkAddr == "" {
				return sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
			}
			resp, err := GetResNodesByNetworkAddr(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)
		},
	}
	cmd.Flags().String(FlagNetworkAddress, "", "(optional) The network address of the node")
	return cmd
}

// GetResNodesByNetworkAddr queries all resource nodes by multiple network IDs (sep: ";")
func GetResNodesByNetworkAddr(cliCtx context.CLIContext, queryRoute string) (res string, err error) {
	queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
	queryByFlagNetworkAddrList := strings.Split(queryFlagNetworkAddr, ";")
	for _, v := range queryByFlagNetworkAddrList {
		resp, _, err := QueryResourceNode(cliCtx, queryRoute, v)
		if err != nil {
			return "null", err
		}
		res += string(resp) + ";"
	}
	return res[:len(res)-1], nil
}

// QueryResourceNode queries resource node by network addr
func QueryResourceNode(cliCtx context.CLIContext, queryRoute, networkAddr string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryResourceNodeByNetworkAddr)
	sdsAddress, err := stratos.SdsAddressFromBech32(networkAddr)
	if err != nil {
		return []byte{}, 0, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
	}

	params := types.NewQueryNodesParams(1, 1, sdsAddress, "", nil)
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return []byte{}, 0, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
	}
	return cliCtx.QueryWithData(route, bz)
}

// GetCmdQueryIndexingNodeList implements the query all indexing nodes by network id command.
func GetCmdQueryIndexingNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-indexing-nodes [flags]", // []byte
		Short: "Query all indexing nodes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all indexing nodes`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			// query all indexing nodes by network address
			queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
			if queryFlagNetworkAddr == "" {
				return sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
			}
			resp, err := GetIndNodesByNetworkAddr(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)

		},
	}
	cmd.Flags().String(FlagNetworkAddress, "", "(optional) The network address of the node")

	return cmd
}

// GetIndNodesByNetworkAddr queries all indexing nodes by multiple network addrs (sep: ";")
func GetIndNodesByNetworkAddr(cliCtx context.CLIContext, queryRoute string) (res string, err error) {
	queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
	queryByFlagNetworkAddrList := strings.Split(queryFlagNetworkAddr, ";")
	for _, v := range queryByFlagNetworkAddrList {
		resp, _, err := QueryIndexingNodes(cliCtx, queryRoute, v)
		if err != nil {
			return "null", err
		}
		res += string(resp) + ";"
	}
	return res[:len(res)-1], nil
}

// QueryIndexingNodes queries all indexing nodes
func QueryIndexingNodes(cliCtx context.CLIContext, queryRoute, networkAddr string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodeByNetworkAddr)
	sdsAddress, err := stratos.SdsAddressFromBech32(networkAddr)
	if err != nil {
		return []byte{}, 0, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
	}

	params := types.NewQueryNodesParams(1, 1, sdsAddress, "", nil)
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return []byte{}, 0, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, "Missing network address")
	}
	return cliCtx.QueryWithData(route, bz)
}
