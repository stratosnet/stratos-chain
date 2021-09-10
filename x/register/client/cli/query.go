package cli

import (
	"bufio"
	"errors"
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
			//GetCmdQueryNetworkSet(queryRoute, cdc),
		)...,
	)

	return registerQueryCmd
}

// GetCmdQueryResourceNodeList implements the query all resource nodes by network id command.
func GetCmdQueryResourceNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get-resource-nodes [flags]", // []byte
		//Args:  cobra.RangeArgs(1, 1),
		Short: "Query all resource nodes by network id.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all resource nodes by network id.`),
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
			queryFlagNetworkID := viper.GetString(FlagNetworkID)
			if queryFlagNetworkID == "" {
				return errors.New("at least one of the flags 'network-id' and 'moniker' must be set")
			}
			resp, err := GetResNodesByNetworkID(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsNetworkID)
	cmd.Flags().AddFlagSet(FsDescription)
	//_ = cmd.MarkFlagRequired(FlagNetworkID)

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

// GetResNodesByNetworkID queries all resource nodes by multiple network IDs (sep: ";")
func GetResNodesByNetworkID(cliCtx context.CLIContext, queryRoute string) (res string, err error) {
	queryFlagNetworkID := viper.GetString(FlagNetworkID)
	queryByFlagNetworkIDList := strings.Split(queryFlagNetworkID, ";")
	for _, v := range queryByFlagNetworkIDList {
		resp, _, err := QueryResourceNodes(cliCtx, queryRoute, v)
		if err != nil {
			return "null", err
		}
		res += string(resp) + ";"
	}
	return res[:len(res)-1], nil
}

// QueryResourceNodes queries all resource nodes by network id
func QueryResourceNodes(cliCtx context.CLIContext, queryRoute, networkID string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryResourceNodeList)
	return cliCtx.QueryWithData(route, []byte(networkID))
}

// GetCmdQueryIndexingNodeList implements the query all indexing nodes by network id command.
func GetCmdQueryIndexingNodeList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get-indexing-nodes [flags]", // []byte
		//Args:  cobra.RangeArgs(1, 1),
		Short: "Query all indexing nodes by network id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all indexing nodes by network id.`),
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
			queryFlagNetworkID := viper.GetString(FlagNetworkID)
			if queryFlagNetworkID == "" {
				return errors.New("at least one of the flags 'network-id' and 'moniker' must be set")
			}
			resp, err := GetIndNodesByNetworkID(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)

		},
	}
	cmd.Flags().AddFlagSet(FsNetworkID)
	cmd.Flags().AddFlagSet(FsDescription)
	//_ = cmd.MarkFlagRequired(FlagNetworkID)

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
func QueryIndNodesByMoniker(cliCtx context.CLIContext, queryRoute, networkID string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodeByMoniker)
	return cliCtx.QueryWithData(route, []byte(networkID))
}

// GetIndNodesByNetworkID queries all indexing nodes by multiple network IDs (sep: ";")
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

// QueryIndexingNodes queries all resource nodes by network is
func QueryIndexingNodes(cliCtx context.CLIContext, queryRoute, networkID string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryIndexingNodeList)
	return cliCtx.QueryWithData(route, []byte(networkID))
}

//// GetCmdQueryNetworkSet implements the query all indexing nodes by network id command.
//func GetCmdQueryNetworkSet(queryRoute string, cdc *codec.Codec) *cobra.Command {
//	return &cobra.Command{
//		Use: "get-network-ids",
//		//Args:  cobra.RangeArgs(0, 0),
//		Short: "Query all network ids",
//		Long: strings.TrimSpace(
//			fmt.Sprintf(`Query all network ids.`),
//		),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			cliCtx := context.NewCLIContext().WithCodec(cdc)
//
//			// query get-network-set by network id
//			resp, _, err := QueryNetworkSet(cliCtx, queryRoute)
//			if err != nil {
//				return err
//			}
//			return cliCtx.PrintOutput(string(resp))
//		},
//	}
//}

//// QueryNetworkSet queries all network id
//func QueryNetworkSet(cliCtx context.CLIContext, queryRoute string) ([]byte, int64, error) {
//	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryNetworkSet)
//	return cliCtx.Query(route)
//}
