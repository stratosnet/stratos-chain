package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	stratos "github.com/stratosnet/stratos-chain/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	flag "github.com/spf13/pflag"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	sdsTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	sdsTxCmd.AddCommand(
		FileUploadTxCmd(),
		PrepayTxCmd(),
	)
	return sdsTxCmd
}

// FileUploadTxCmd will create a file upload tx and sign it with the given key.
func FileUploadTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload [flags]",
		Short: "Create and sign a file upload tx",
		//Args:  cobra.RangeArgs(0, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildFileuploadMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}
	//cmd = flags.PostCommands(cmd)[0]
	//cmd.Flags().String(flags.FlagFrom, "", "from address")
	//cmd.Flags().String(FlagFileHash, "", "Hash of uploaded file")
	//cmd.Flags().String(FlagReporter, "", "Reporter of file")
	//cmd.Flags().String(FlagUploader, "", "Uploader of file")
	cmd.Flags().AddFlagSet(flagSetFileHash())
	cmd.Flags().AddFlagSet(flagSetReporter())
	cmd.Flags().AddFlagSet(flagSetUploader())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagFileHash)
	_ = cmd.MarkFlagRequired(FlagReporter)
	_ = cmd.MarkFlagRequired(FlagUploader)

	return cmd
}

// PrepayTxCmd will create a prepay tx and sign it with the given key.
func PrepayTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepay [from_address] [coins]",
		Short: "Create and sign a prepay tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromAddr, fromName, _, err := client.GetFromFields(cliCtx.Keyring, args[0], cliCtx.GenerateOnly)
			if err != nil {
				return err
			}

			clientCtx := cliCtx.WithFrom(args[0]).WithFromAddress(fromAddr).WithFromName(fromName)

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildPrepayMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
			//
			//
			//
			//inBuf := bufio.NewReader(cmd.InOrStdin())
			//txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			//cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			//
			//coins, err := sdk.ParseCoins(args[1])
			//if err != nil {
			//	return err
			//}
			//
			//// build and sign the transaction, then broadcast to Tendermint
			//msg := types.NewMsgPrepay(cliCtx.GetFromAddress(), coins)
			//return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	//cmd = flags.PostCommands(cmd)[0]

	return cmd
}

// makes a new newBuildFileuploadMsg
func newBuildFileuploadMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgFileUpload, error) {
	fileHash, err := fs.GetString(FlagFileHash)
	if err != nil {
		return txf, nil, err
	}
	_, err = hex.DecodeString(fileHash)
	if err != nil {
		return txf, nil, err
	}

	flagReporterStr, err := fs.GetString(FlagReporter)
	if err != nil {
		return txf, nil, err
	}
	_, err = stratos.SdsAddressFromBech32(flagReporterStr)
	if err != nil {
		return txf, nil, err
	}

	flagUploaderStr, err := fs.GetString(FlagUploader)
	if err != nil {
		return txf, nil, err
	}
	_, err = sdk.AccAddressFromBech32(flagUploaderStr)
	if err != nil {
		return txf, nil, err
	}

	msg := types.NewMsgUpload(
		fileHash,
		clientCtx.GetFromAddress().String(),
		flagReporterStr,
		flagUploaderStr,
	)

	return txf, msg, nil
}

// makes a new newBuildPrepayMsg
func newBuildPrepayMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgPrepay, error) {
	coin, err := sdk.ParseCoinNormalized(fs.Arg(1))
	if err != nil {
		return txf, nil, err
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := types.NewMsgPrepay(clientCtx.GetFromAddress().String(), sdk.NewCoins(coin))

	return txf, msg, nil
}
