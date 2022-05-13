package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

type singleWalletVolumeStr struct {
	WalletAddress string `json:"wallet_address"`
	Volume        string `json:"volume"`
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

	potTxCmd.AddCommand(
		VolumeReportCmd(),
		WithdrawCmd(),
		FoundationDepositCmd(),
		SlashingResourceNodeCmd(),
	)
	return potTxCmd
}

func WithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "withdraw POT reward",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := buildWithdrawMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsTargetAddress)

	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// makes a new WithdrawMsg.
func buildWithdrawMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgWithdraw, error) {
	amountStr, _ := fs.GetString(FlagAmount)
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return txf, nil, err
	}

	walletAddress := clientCtx.GetFromAddress()

	var targetAddress sdk.AccAddress
	if viper.IsSet(FlagTargetAddress) {
		targetAddressStr := viper.GetString(FlagTargetAddress)
		targetAddress, err = sdk.AccAddressFromBech32(targetAddressStr)
		if err != nil {
			return txf, nil, err
		}
	} else {
		targetAddress = walletAddress
	}

	msg := types.NewMsgWithdraw(amount, walletAddress, targetAddress)

	return txf, msg, nil
}

// VolumeReportCmd will report wallets volume and sign it with the given key.
func VolumeReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report [flags]",
		Short: "Create and sign a volume report",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := createVolumeReportMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)

			//inBuf := bufio.NewReader(cmd.InOrStdin())
			//txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			//cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			//txBldr, msg, err := createVolumeReportMsg(cliCtx, txBldr)
			//if err != nil {
			//	return err
			//}
			//return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsReporterAddr)
	cmd.Flags().AddFlagSet(FsEpoch)
	cmd.Flags().AddFlagSet(FsReportReference)
	cmd.Flags().AddFlagSet(FsWalletVolumes)

	_ = cmd.MarkFlagRequired(FlagReporterAddr)
	_ = cmd.MarkFlagRequired(FlagEpoch)
	_ = cmd.MarkFlagRequired(FlagReportReference)
	_ = cmd.MarkFlagRequired(FlagWalletVolumes)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func createVolumeReportMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgVolumeReport, error) {
	reporterStr, _ := fs.GetString(FlagReporterAddr)
	reporter, err := stratos.SdsAddressFromBech32(reporterStr)
	if err != nil {
		return txf, nil, err
	}

	reportReference := viper.GetString(FlagReportReference)
	value, err := strconv.ParseInt(viper.GetString(FlagEpoch), 10, 64)
	if err != nil {
		return txf, nil, err
	}
	epoch := sdk.NewInt(value)
	var walletVolumesStr = make([]singleWalletVolumeStr, 0)
	err = json.Unmarshal([]byte(viper.GetString(FlagWalletVolumes)), &walletVolumesStr)
	if err != nil {
		return txf, nil, err
	}

	var walletVolumes = make([]*types.SingleWalletVolume, 0)
	for _, n := range walletVolumesStr {
		walletAcc, err := sdk.AccAddressFromBech32(n.WalletAddress)
		if err != nil {
			return txf, nil, err
		}
		volumeInt64, err := strconv.ParseInt(n.Volume, 10, 64)
		if err != nil {
			return txf, nil, err
		}
		volume := sdk.NewInt(volumeInt64)
		walletVolumes = append(walletVolumes, types.NewSingleWalletVolume(walletAcc, volume))
	}

	reporterOwner := clientCtx.GetFromAddress()

	msg := types.NewMsgVolumeReport(
		walletVolumes,
		reporter,
		epoch,
		reportReference,
		reporterOwner,
		types.BLSSignatureInfo{},
	)
	return txf, msg, nil
}

func FoundationDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "foundation-deposit",
		Short: "Deposit to foundation account",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := buildFoundationDepositMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAmount)

	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func buildFoundationDepositMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgFoundationDeposit, error) {
	amountStr, _ := fs.GetString(FlagAmount)
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return txf, nil, err
	}
	from := clientCtx.GetFromAddress()
	msg := types.NewMsgFoundationDeposit(amount, from)
	return txf, msg, nil
}

func SlashingResourceNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slashing",
		Short: "slashing resource node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := buildSlashingResourceNodeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsReporters)
	cmd.Flags().AddFlagSet(FsReportOwner)
	cmd.Flags().AddFlagSet(FsNetworkAddress)
	cmd.Flags().AddFlagSet(FsWalletAddress)
	cmd.Flags().AddFlagSet(FsSlashing)
	cmd.Flags().AddFlagSet(FsSuspend)

	_ = cmd.MarkFlagRequired(FlagReporters)
	_ = cmd.MarkFlagRequired(FlagReporterOwner)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagWalletAddress)
	_ = cmd.MarkFlagRequired(FlagSlashing)
	_ = cmd.MarkFlagRequired(FlagSuspend)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func buildSlashingResourceNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgSlashingResourceNode, error) {
	var reportersStr = make([]string, 0)
	err := json.Unmarshal([]byte(viper.GetString(FlagReporters)), &reportersStr)
	if err != nil {
		return txf, nil, err
	}
	var reporters = make([]stratos.SdsAddress, 0)
	for _, val := range reportersStr {
		reporterAddr, err := stratos.SdsAddressFromBech32(val)
		if err != nil {
			return txf, nil, err
		}
		reporters = append(reporters, reporterAddr)
	}

	var reporterOwnerStr = make([]string, 0)
	err = json.Unmarshal([]byte(viper.GetString(FlagReporterOwner)), &reporterOwnerStr)
	if err != nil {
		return txf, nil, err
	}
	var reporterOwner = make([]sdk.AccAddress, 0)
	for _, val := range reporterOwnerStr {
		reporterOwnerAddr, err := sdk.AccAddressFromBech32(val)
		if err != nil {
			return txf, nil, err
		}
		reporterOwner = append(reporterOwner, reporterOwnerAddr)
	}

	networkAddressStr := viper.GetString(FlagNetworkAddress)
	networkAddress, err := stratos.SdsAddressFromBech32(networkAddressStr)
	if err != nil {
		return txf, nil, err
	}

	walletAddressStr := viper.GetString(FlagWalletAddress)
	walletAddress, err := sdk.AccAddressFromBech32(walletAddressStr)
	if err != nil {
		return txf, nil, err
	}

	slashingVal, err := strconv.ParseInt(viper.GetString(FlagSlashing), 10, 64)
	if err != nil {
		return txf, nil, err
	}
	slashing := sdk.NewInt(slashingVal)

	suspendVal := viper.GetString(FlagSuspend)
	suspend, err := strconv.ParseBool(suspendVal)
	if err != nil {
		return txf, nil, err
	}

	msg := types.NewMsgSlashingResourceNode(reporters, reporterOwner, networkAddress, walletAddress, slashing, suspend)
	return txf, msg, nil
}
