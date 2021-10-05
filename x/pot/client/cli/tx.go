package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

type singleNodeVolumeStr struct {
	NodeAddress string `json:"node_address"`
	Volume      string `json:"node_volume"`
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	potTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	potTxCmd.AddCommand(flags.PostCommands(
		VolumeReportCmd(cdc),
		WithdrawCmd(cdc),
		FoundationDepositCmd(cdc),
	)...)
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

	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagNodeAddress)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

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
		Use:   "report [flags]",
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
		},
	}
	cmd.Flags().AddFlagSet(FsReporter)
	cmd.Flags().AddFlagSet(FsEpoch)
	cmd.Flags().AddFlagSet(FsReportReference)
	cmd.Flags().AddFlagSet(FsNodesVolume)

	_ = cmd.MarkFlagRequired(FlagReporter)
	_ = cmd.MarkFlagRequired(FlagEpoch)
	_ = cmd.MarkFlagRequired(FlagReportReference)
	_ = cmd.MarkFlagRequired(FlagNodesVolume)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func createVolumeReportMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	reporterStr := viper.GetString(FlagReporter)
	reporter, err := sdk.AccAddressFromBech32(reporterStr)
	if err != nil {
		return txBldr, nil, err
	}

	reportReference := viper.GetString(FlagReportReference)
	value, err := strconv.ParseInt(viper.GetString(FlagEpoch), 10, 64)
	if err != nil {
		return txBldr, nil, err
	}

	epoch := sdk.NewInt(value)
	lastEpochByte, _, err := cliCtx.QueryStore(types.LastMaturedEpochKey, types.StoreKey)
	if err != nil {
		return txBldr, nil, err
	}
	var lastEpoch sdk.Int
	if lastEpochByte == nil {
		lastEpoch = sdk.ZeroInt()
	} else {
		cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(lastEpochByte, &lastEpoch)
	}
	if epoch.LTE(lastEpoch) {
		e := errors.Wrapf(types.ErrMatureEpoch, "expected epoch should be greater than %s, got %s",
			lastEpoch.String(), epoch.String())
		return txBldr, nil, e
	}

	var nodesVolumeStr = make([]singleNodeVolumeStr, 0)
	err = cliCtx.Codec.UnmarshalJSON([]byte(viper.GetString(FlagNodesVolume)), &nodesVolumeStr)
	if err != nil {
		return txBldr, nil, err
	}

	var nodesVolume = make([]types.SingleNodeVolume, 0)
	for _, n := range nodesVolumeStr {
		nodeAcc, err := sdk.AccAddressFromBech32(n.NodeAddress)
		if err != nil {
			return txBldr, nil, err
		}
		volumeInt64, err := strconv.ParseInt(n.Volume, 10, 64)
		if err != nil {
			return txBldr, nil, err
		}
		nodeVolume := sdk.NewInt(volumeInt64)
		nodesVolume = append(nodesVolume, types.NewSingleNodeVolume(nodeAcc, nodeVolume))
	}

	reporterOwner := cliCtx.GetFromAddress()

	msg := types.NewMsgVolumeReport(
		nodesVolume,
		reporter,
		epoch,
		reportReference,
		reporterOwner,
	)
	return txBldr, msg, nil
}

func FoundationDepositCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "foundation-deposit",
		Short: "Deposit to foundation account",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := buildFoundationDepositMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAmount)

	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func buildFoundationDepositMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	amountStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoin(amountStr)
	if err != nil {
		return txBldr, nil, err
	}
	from := cliCtx.GetFromAddress()
	msg := types.NewMsgFoundationDeposit(amount, from)
	return txBldr, msg, nil
}
