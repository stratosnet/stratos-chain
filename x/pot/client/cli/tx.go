package cli

import (
	"bufio"
	"github.com/spf13/viper"

	//"encoding/hex"
	//"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	//"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	//"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	potTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	potTxCmd.AddCommand(
		VolumeReportCmd(cdc),
		WithdrawCmd(cdc),
	)
	return potTxCmd
}

func WithdrawCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "withdraw POT reward",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := buildWithdrawMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsNodeAddress)

	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagNodeAddress)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// makes a new WithdrawMsg.
func buildWithdrawMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	amountStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoin(amountStr)
	if err != nil {
		return txBldr, nil, err
	}
	nodeAddressStr := viper.GetString(FlagNodeAddress)
	nodeAddress, err := sdk.AccAddressFromBech32(nodeAddressStr)
	if err != nil {
		return txBldr, nil, err
	}
	ownerAddress := cliCtx.GetFromAddress()

	msg := types.NewMsgWithdraw(amount, nodeAddress, ownerAddress)

	return txBldr, msg, nil
}

// VolumeReportCmd will report nodes volume and sign it with the given key.
func VolumeReportCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report [reporter] [epoch] [report_reference] [nodes_volume]",
		Short: "Create and sign a volume report",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			reporter, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			reportReference := args[2]
			value, e := strconv.ParseInt(args[1], 10, 64)
			if e != nil {
				return err
			}
			epoch := sdk.NewInt(value)

			var nodesVolume = make([]types.SingleNodeVolume, 0)
			//er := types.ModuleCdc.UnmarshalJSON([]byte(args[3]), &nodesVolume)
			er := cliCtx.Codec.UnmarshalJSON([]byte(args[3]), &nodesVolume)
			if er != nil {
				return er
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgVolumeReport(
				nodesVolume,
				reporter,
				epoch,
				reportReference,
			)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
