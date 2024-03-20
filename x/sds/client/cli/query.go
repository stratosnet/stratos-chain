package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// GetQueryCmd returns the cli query commands for sds module
func GetQueryCmd() *cobra.Command {
	// Group sds queries under a subcommand
	sdsQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	sdsQueryCmd.AddCommand(
		GetCmdQueryUploadedFile(),
		GetCmdQueryParams(),
	)

	return sdsQueryCmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current sds parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as sds parameters.

Example:
$ %s query sds params
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

// GetCmdQueryUploadedFile implements the query uploaded file command.
func GetCmdQueryUploadedFile() *cobra.Command {
	cmd := &cobra.Command{
		//return &cobra.Command{
		Use:   "upload [file_hash]",
		Args:  cobra.ExactArgs(1),
		Short: "Query uploaded file info by hash",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query uploaded file info by hash.

Example:
$ %s query sds upload c03661732294feb49caf6dc16c7cbb2534986d73
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

			queryFileHash := strings.TrimSpace(args[0][:])
			if len(queryFileHash) == 0 {
				return errors.Wrap(types.ErrEmptyFileHash, "Missing file hash")
			}

			result, err := queryClient.Fileupload(cmd.Context(), &types.QueryFileUploadRequest{
				FileHash: queryFileHash,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "upload")
	return cmd
}
