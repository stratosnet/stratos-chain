package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Use:   "report [flags]", // reporter: []byte
		Short: "Query volume report hash by epoch",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query volume report hash by epoch.`),
		),

		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			epochStr := viper.GetString(FlagEpoch)
			epoch, err := checkFlagEpoch(epochStr)
			if err != nil {
				return err
			}
			resp, _, err := QueryVolumeReport(cliCtx, queryRoute, epoch)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(resp)

		},
	}
	cmd.Flags().AddFlagSet(FsEpoch)
	_ = cmd.MarkFlagRequired(FlagEpoch)

	return cmd
}

func QueryVolumeReport(cliCtx context.CLIContext, queryRoute string, epoch sdk.Int) (types.ReportInfo, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryVolumeReportHash)
	resp, height, err := cliCtx.QueryWithData(route, []byte(epoch.String()))
	if err != nil {
		return types.ReportInfo{}, height, err
	}
	reportRes := types.NewReportInfo(epoch, string(resp))
	return reportRes, height, nil
}

func checkFlagEpoch(epochStr string) (sdk.Int, error) {
	epochInt64, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return sdk.NewInt(0), fmt.Errorf("invalid epoch: %w", err)
	}
	epoch := sdk.NewInt(epochInt64)
	return epoch, nil
}
