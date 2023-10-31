package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

type singleWalletVolumeStr struct {
	WalletAddress string `json:"wallet_address"`
	Volume        string `json:"volume"`
}

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
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

			msg, err := buildWithdrawMsg(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetAmount())
	cmd.Flags().AddFlagSet(flagSetTargetAddress())
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// makes a new WithdrawMsg.
func buildWithdrawMsg(clientCtx client.Context, fs *flag.FlagSet) (*types.MsgWithdraw, error) {
	amountStr, err := fs.GetString(FlagAmount)
	if err != nil {
		return nil, err
	}
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return nil, err
	}

	walletAddress := clientCtx.GetFromAddress()

	var targetAddress sdk.AccAddress
	flagTargetAddress := fs.Lookup(FlagTargetAddress)
	if flagTargetAddress == nil {
		targetAddress = walletAddress
	} else {
		targetAddressStr, _ := fs.GetString(FlagTargetAddress)
		targetAddress, err = sdk.AccAddressFromBech32(targetAddressStr)
		if err != nil {
			return nil, err
		}
	}

	msg := types.NewMsgWithdraw(amount, walletAddress, targetAddress)

	return msg, nil
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

			msg, err := createVolumeReportMsg(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(flagSetReportVolumes())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagReporterAddr)
	_ = cmd.MarkFlagRequired(FlagEpoch)
	_ = cmd.MarkFlagRequired(FlagReportReference)
	_ = cmd.MarkFlagRequired(FlagWalletVolumes)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func createVolumeReportMsg(clientCtx client.Context, fs *flag.FlagSet) (*types.MsgVolumeReport, error) {
	reporterStr, err := fs.GetString(FlagReporterAddr)
	if err != nil {
		return nil, err
	}
	reporter, err := stratos.SdsAddressFromBech32(reporterStr)
	if err != nil {
		return nil, err
	}

	reportReference, err := fs.GetString(FlagReportReference)
	if err != nil {
		return nil, err
	}
	//flagEpochInt64, err := fs.GetInt64(FlagEpoch)
	flagEpochStr, err := fs.GetString(FlagEpoch)
	if err != nil {
		return nil, err
	}
	value, err := strconv.ParseInt(flagEpochStr, 10, 64)
	if err != nil {
		return nil, err
	}
	epoch := sdkmath.NewInt(value)

	flagWalletVolumes, err := fs.GetString(FlagWalletVolumes)
	if err != nil {
		return nil, err
	}
	walletVolumesStr := make([]singleWalletVolumeStr, 0)
	err = json.Unmarshal([]byte(flagWalletVolumes), &walletVolumesStr)
	if err != nil {
		return nil, err
	}

	var walletVolumes = make([]types.SingleWalletVolume, 0)
	for _, n := range walletVolumesStr {
		walletAcc, err := sdk.AccAddressFromBech32(n.WalletAddress)
		if err != nil {
			return nil, err
		}
		volumeInt64, err := strconv.ParseInt(n.Volume, 10, 64)
		if err != nil {
			return nil, err
		}
		volume := sdkmath.NewInt(volumeInt64)
		walletVolumes = append(walletVolumes, types.NewSingleWalletVolume(walletAcc, volume))
	}

	reporterOwner := clientCtx.GetFromAddress()

	blsSigture, err := fs.GetString(FlagBLSSignature)
	if err != nil {
		return nil, err
	}

	var sig types.BaseBLSSignatureInfo
	err = json.Unmarshal([]byte(blsSigture), &sig)
	if err != nil {
		return nil, err
	}

	// TODO: change pubkey
	pubKeys := make([][]byte, len(sig.PubKeys))
	for i, v := range sig.PubKeys {
		pubKeys[i] = []byte(v)
	}

	signature := types.NewBLSSignatureInfo(pubKeys, []byte(sig.Signature), []byte(sig.TxData))

	msg := types.NewMsgVolumeReport(
		walletVolumes,
		reporter,
		epoch,
		reportReference,
		reporterOwner,
		signature,
	)
	return msg, nil
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

			msg, err := buildFoundationDepositMsg(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(flagSetAmount())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func buildFoundationDepositMsg(clientCtx client.Context, fs *flag.FlagSet) (*types.MsgFoundationDeposit, error) {
	amountStr, err := fs.GetString(FlagAmount)
	if err != nil {
		return nil, err
	}
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return nil, err
	}
	from := clientCtx.GetFromAddress()
	msg := types.NewMsgFoundationDeposit(amount, from)
	return msg, nil
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

			msg, err := buildSlashingResourceNodeMsg(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(flagSetReportersAndOwners())
	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetWalletAddress())
	cmd.Flags().AddFlagSet(flagSetSlashing())
	cmd.Flags().AddFlagSet(flagSetSuspend())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagReporters)
	_ = cmd.MarkFlagRequired(FlagReporterOwner)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagWalletAddress)
	_ = cmd.MarkFlagRequired(FlagSlashing)
	_ = cmd.MarkFlagRequired(FlagSuspend)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func buildSlashingResourceNodeMsg(clientCtx client.Context, fs *flag.FlagSet) (*types.MsgSlashingResourceNode, error) {
	var reportersStr = make([]string, 0)
	flagReportersStr, err := fs.GetString(FlagReporters)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(flagReportersStr), &reportersStr)
	if err != nil {
		return nil, err
	}

	var reporters = make([]stratos.SdsAddress, 0)
	for _, val := range reportersStr {
		reporterAddr, err := stratos.SdsAddressFromBech32(val)
		if err != nil {
			return nil, err
		}
		reporters = append(reporters, reporterAddr)
	}

	var reporterOwnerStr = make([]string, 0)
	flagReporterOwnerStr, err := fs.GetString(FlagReporterOwner)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(flagReporterOwnerStr), &reporterOwnerStr)
	if err != nil {
		return nil, err
	}

	var reporterOwner = make([]sdk.AccAddress, 0)
	for _, val := range reporterOwnerStr {
		reporterOwnerAddr, err := sdk.AccAddressFromBech32(val)
		if err != nil {
			return nil, err
		}
		reporterOwner = append(reporterOwner, reporterOwnerAddr)
	}

	flagNetworkAddressStr, err := fs.GetString(FlagNetworkAddress)
	if err != nil {
		return nil, err
	}
	networkAddress, err := stratos.SdsAddressFromBech32(flagNetworkAddressStr)
	if err != nil {
		return nil, err
	}

	flagWalletAddressStr, err := fs.GetString(FlagWalletAddress)
	if err != nil {
		return nil, err
	}
	walletAddress, err := sdk.AccAddressFromBech32(flagWalletAddressStr)
	if err != nil {
		return nil, err
	}

	flagSlashingStr, err := fs.GetString(FlagSlashing)
	if err != nil {
		return nil, err
	}
	slashingVal, err := strconv.ParseInt(flagSlashingStr, 10, 64)
	if err != nil {
		return nil, err
	}
	slashing := sdkmath.NewInt(slashingVal)

	suspend, err := fs.GetBool(FlagSuspend)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgSlashingResourceNode(reporters, reporterOwner, networkAddress, walletAddress, slashing, suspend)
	return msg, nil
}
