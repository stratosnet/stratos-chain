package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"strings"

	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	// "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// GetQueryCmd returns the cli query commands for pot module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group pot queries under a subcommand
	potQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	potQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryVolumeReport(queryRoute, cdc),
		)...,
	)

	return potQueryCmd
}

// GetCmdQueryVolumeReport implements the query volume report command.
func GetCmdQueryVolumeReport(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report [reporter]", // reporter: []byte
		Args:  cobra.RangeArgs(1, 1),
		Short: "Query volume report hash by reporter addr",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query volume report hash by reporter.`),
		),
		//	RunE: func(cmd *cobra.Command, args []string) error {
		//		cliCtx := context.NewCLIContext().WithCodec(cdc)
		//
		//		// query volume report by reporter
		//		resp, _, err := QueryVolumeReport(cliCtx, queryRoute, args[0])
		//		if err != nil {
		//			return err
		//		}
		//		return cliCtx.PrintOutput(string(resp))
		//	},
		//}
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			//txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			// query all indexing nodes by network address
			resp, _, err := QueryVolumeReport(cliCtx, queryRoute, viper.GetString(FlagReporter))
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(string(resp))

		},
	}
	cmd.Flags().AddFlagSet(FsReporter)
	_ = cmd.MarkFlagRequired(FlagReporter)

	return cmd
}

// QueryVolumeReport queries the volume hash by reporter
func QueryVolumeReport(cliCtx context.CLIContext, queryRoute, reporter string) ([]byte, int64, error) {
	accAddr, err := sdk.AccAddressFromBech32(reporter)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid reporter, please specify a reporter in Bech32 format %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryVolumeReportHash)
	return cliCtx.QueryWithData(route, accAddr)
}
