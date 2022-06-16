package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	stratos "github.com/stratosnet/stratos-chain/types"
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
		GetCmdQueryPrepayBalance(),
	)

	return sdsQueryCmd
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
				return sdkerrors.Wrap(types.ErrEmptyFileHash, "Missing file hash")
			}

			// query file by fileHash
			//resp, _, err := common.QueryUploadedFile(clientCtx, queryRoute, args[0])
			//if err != nil {
			//	return err
			//}
			//fi := types.MustUnmarshalFileInfo(cdc, resp)
			//return cliCtx.PrintOutput(fi.String())

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

// GetCmdQueryPrepayBalance implements the query prepay balance command.
func GetCmdQueryPrepayBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepay [acct_addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query balance of prepayment in Volume Pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query balance of prepayment in Volumn Pool.

Example:
$ %s query sds prepay st1yx3kkx9jnqeck59j744nc5qgtv4lt4dc45jcwz
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

			queryAccAddr := strings.TrimSpace(args[0][:])
			if len(queryAccAddr) == 0 {
				return sdkerrors.Wrap(types.ErrEmptySenderAddr, "Missing sender address")
			}
			_, err = stratos.SdsAddressFromBech32(queryAccAddr)
			if err != nil {
				return err
			}

			result, err := queryClient.Prepay(cmd.Context(), &types.QueryPrepayRequest{
				AcctAddr: queryAccAddr,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
