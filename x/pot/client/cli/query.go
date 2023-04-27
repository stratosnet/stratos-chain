package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

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
		GetCmdQueryParams(),
		GetCmdQueryTotalMinedTokens(),
		GetCmdQueryCirculationSupply(),
	)

	return potQueryCmd
}

func GetCmdQueryCirculationSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "circulation-supply",
		Args:  cobra.NoArgs,
		Short: "Query the circulation supply",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the circulation supply.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CirculationSupply(cmd.Context(), &types.QueryCirculationSupplyRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryTotalMinedTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-mined-token",
		Args:  cobra.NoArgs,
		Short: "Query the total mined tokens",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the total mined tokens.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalMinedToken(cmd.Context(), &types.QueryTotalMinedTokenRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current pot parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as pot parameters.

Example:
$ %s query pot params
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

// GetCmdQueryVolumeReport implements the query volume report command.
func GetCmdQueryVolumeReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report [flags]",
		Short: "Query volume report hash by epoch",
		Long: strings.TrimSpace(
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

func checkFlagEpoch(epochStr string) (sdk.Int, error) {
	epochInt64, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return sdk.NewInt(0), fmt.Errorf("invalid epoch: %w", err)
	}
	epoch := sdk.NewInt(epochInt64)
	return epoch, nil
}
