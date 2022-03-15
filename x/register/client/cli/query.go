package cli

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
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
			GetCmdQueryResourceNodeList(queryRoute, cdc),
			GetCmdQueryIndexingNodeList(queryRoute, cdc),
		)...,
	)

	return registerQueryCmd
}

// GetCmdQueryResourceNodeList implements the query all resource nodes by network id command.
func GetCmdQueryResourceNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-resource-nodes [flags]", // []byte
		Short: "Query all resource nodes by network id or moniker",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all resource nodes by network id or moniker`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			// query all resource nodes by moniker
			queryFlagMoniker := viper.GetString(FlagMoniker)
			if queryFlagMoniker != "" {
				resp, err := GetResNodesByMoniker(cliCtx, queryRoute, queryFlagMoniker)
				if err != nil {
					return err
				}
				return cliCtx.PrintOutput(resp)
			}

			// query all resource nodes by network id
			queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
			if queryFlagNetworkAddr == "" {
				return errors.New("at least one of the flags 'network-addr' and 'moniker' must be set")
			}
			resp, err := GetResNodesByNetworkAddr(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)
		},
	}
	cmd.Flags().String(FlagNetworkAddress, "", "(optional) The network address of the node")
	cmd.Flags().String(FlagMoniker, "", "(optional) The name of the node")

	return cmd
}

func GetResNodesByMoniker(cliCtx context.CLIContext, queryRoute string, queryFlagMoniker string) (res string, err error) {
	queryByFlagMonikerList := strings.Split(queryFlagMoniker, ";")
	for _, v := range queryByFlagMonikerList {
		resp, _, err := QueryResNodesByMoniker(cliCtx, queryRoute, v)
		if err != nil {
			return "null", err
		}
		res += string(resp) + ";"
	}
	return res[:len(res)-1], nil
}

func QueryResNodesByMoniker(cliCtx context.CLIContext, queryRoute, moniker string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryResourceNodeByMoniker)
	return cliCtx.QueryWithData(route, []byte(moniker))
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
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryResourceNodesByNetworkAddr)
	return cliCtx.QueryWithData(route, []byte(networkAddr))
}

// GetCmdQueryIndexingNodeList implements the query all indexing nodes by network id command.
func GetCmdQueryIndexingNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-indexing-nodes [flags]", // []byte
		Short: "Query all indexing nodes by network id or moniker",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all indexing nodes by network id or moniker`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			// query all resource nodes by moniker
			queryFlagMoniker := viper.GetString(FlagMoniker)
			if queryFlagMoniker != "" {
				resp, err := GetIndByMoniker(cliCtx, queryRoute, queryFlagMoniker)
				if err != nil {
					return err
				}
				return cliCtx.PrintOutput(resp)
			}

			// query all indexing nodes by network id
			queryFlagNetworkAddr := viper.GetString(FlagNetworkAddress)
			if queryFlagNetworkAddr == "" {
				return errors.New("at least one of the flags 'network-addr' and 'moniker' must be set")
			}
			resp, err := GetIndNodesByNetworkAddr(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)

		},
	}
	cmd.Flags().String(FlagNetworkAddress, "", "(optional) The network address of the node")
	cmd.Flags().String(FlagMoniker, "", "(optional) The name of the node")

	return cmd
}

func GetIndByMoniker(cliCtx context.CLIContext, queryRoute string, queryFlagMoniker string) (res string, err error) {
	queryByFlagMonikerList := strings.Split(queryFlagMoniker, ";")
	for _, v := range queryByFlagMonikerList {
		resp, _, err := QueryIndNodesByMoniker(cliCtx, queryRoute, v)
		if err != nil {
			return "null", err
		}
		res += string(resp) + ";"
	}
	return res[:len(res)-1], nil
}

// QueryIndNodesByMoniker queries all indexing nodes by network ID
func QueryIndNodesByMoniker(cliCtx context.CLIContext, queryRoute, networkAddr string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodeByMoniker)
	return cliCtx.QueryWithData(route, []byte(networkAddr))
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

// QueryIndexingNodes queries all resource nodes by network is
func QueryIndexingNodes(cliCtx context.CLIContext, queryRoute, networkAddr string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodesByNetworkAddr)
	return cliCtx.QueryWithData(route, []byte(networkAddr))
}
