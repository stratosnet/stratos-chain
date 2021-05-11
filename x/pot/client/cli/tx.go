package cli

import (
	"bufio"
	"github.com/cosmos/cosmos-sdk/client/flags"
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
	)
	return potTxCmd
}

// VolumeReportCmd will report nodes volume and sign it with the given key.
func VolumeReportCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Create and sign a volume report",
		//Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := createVolumeReportMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})

			//reporter, err := sdk.AccAddressFromBech32(args[0])
			//if err != nil {
			//	return err
			//}
			//
			//reportReference := args[2]
			//value, e := strconv.ParseInt(args[1], 10, 64)
			//if e != nil {
			//	return err
			//}
			//epoch := sdk.NewInt(value)
			//
			//var nodesVolume = make([]types.SingleNodeVolume, 0)
			////er := types.ModuleCdc.UnmarshalJSON([]byte(args[3]), &nodesVolume)
			//er := cliCtx.Codec.UnmarshalJSON([]byte(args[3]), &nodesVolume)
			//if er != nil {
			//	return er
			//}
			//
			//// build and sign the transaction, then broadcast to Tendermint
			//msg := types.NewMsgVolumeReport(
			//	nodesVolume,
			//	reporter,
			//	epoch,
			//	reportReference,
			//)
			//return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	//cmd.Flags().AddFlagSet(FsReporter)
	cmd.Flags().AddFlagSet(FsEpoch)
	cmd.Flags().AddFlagSet(FsReportReference)
	cmd.Flags().AddFlagSet(FsNodesVolume)

	//_ = cmd.MarkFlagRequired(FlagReporter)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagEpoch)
	_ = cmd.MarkFlagRequired(FlagReportReference)
	_ = cmd.MarkFlagRequired(FlagNodesVolume)

	//cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func createVolumeReportMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	reporter := cliCtx.GetFromAddress()
	//reporter, e := sdk.AccAddressFromBech32(reporterStr)
	//if e != nil {
	//	return txBldr, nil, e
	//}

	reportReference := viper.GetString(FlagReportReference)

	value, er := strconv.ParseInt(viper.GetString(FlagEpoch), 10, 64)
	if er != nil {
		return txBldr, nil, er
	}
	epoch := sdk.NewInt(value)

	var nodesVolume = make([]types.SingleNodeVolume, 0)
	err := cliCtx.Codec.UnmarshalJSON([]byte(viper.GetString(FlagNodesVolume)), &nodesVolume)
	if err != nil {
		return txBldr, nil, err
	}

	msg := types.NewMsgVolumeReport(
		nodesVolume,
		reporter,
		epoch,
		reportReference,
	)
	return txBldr, msg, nil
}
