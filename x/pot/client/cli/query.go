package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// GetQueryCmd returns the cli query commands for pot module
func GetQueryCmd() *cobra.Command {
	// Group pot queries under a subcommand
	potQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	potQueryCmd.AddCommand(
		GetCmdQueryVolumeReport(),
	)

	return potQueryCmd
}

// GetCmdQueryVolumeReport implements the query volume report command.
func GetCmdQueryVolumeReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report [flags]", // reporter: []byte
		Short: "Query volume report hash by epoch",
		Long: strings.TrimSpace(
			//fmt.Sprintf(`Query volume report hash by reporter.`),
			fmt.Sprintf(`Query volume report hash by epoch.`),
		),

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			epochStr := viper.GetString(FlagEpoch)
			epoch, err := checkFlagEpoch(epochStr)
			if err != nil {
				return err
			}

			result, err := queryClient.VolumeReport(cmd.Context(), &types.QueryVolumeReportRequest{
				Epoch: epoch.Int64(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}
	cmd.Flags().AddFlagSet(flagSetEpoch())
	_ = cmd.MarkFlagRequired(FlagEpoch)

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

//func QueryVolumeReport(cliCtx context.CLIContext, queryRoute string, epoch sdk.Int) (types.ReportInfo, int64, error) {
//	route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryVolumeReport)
//	resp, height, err := cliCtx.QueryWithData(route, []byte(epoch.String()))
//	if err != nil {
//		return types.ReportInfo{}, height, err
//	}
//	reportRes := types.NewReportInfo(epoch, string(resp))
//	return reportRes, height, nil
//}

func checkFlagEpoch(epochStr string) (sdk.Int, error) {
	epochInt64, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return sdk.NewInt(0), fmt.Errorf("invalid epoch: %w", err)
	}
	epoch := sdk.NewInt(epochInt64)
	return epoch, nil
}
