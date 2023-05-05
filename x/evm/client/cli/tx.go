package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	rpctypes "github.com/stratosnet/stratos-chain/rpc/types"
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(NewRawTxCmd())
	return cmd
}

// NewRawTxCmd command build cosmos transaction from raw ethereum transaction
func NewRawTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw [tx-hex]",
		Short: "Build cosmos transaction from raw ethereum transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := hexutil.Decode(args[0])
			if err != nil {
				return errors.Wrap(err, "failed to decode ethereum tx hex bytes")
			}

			msg := &types.MsgEthereumTx{}
			if err := msg.UnmarshalBinary(data); err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rsp, err := rpctypes.NewQueryClient(clientCtx).Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			tx, err := msg.BuildTx(clientCtx.TxConfig.NewTxBuilder(), rsp.Params.EvmDenom)
			if err != nil {
				return err
			}

			if clientCtx.GenerateOnly {
				json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
				if err != nil {
					return err
				}

				return clientCtx.PrintString(fmt.Sprintf("%s\n", json))
			}

			if !clientCtx.SkipConfirm {
				out, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
				if err != nil {
					return err
				}

				_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

				buf := bufio.NewReader(os.Stdin)
				ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

				if err != nil || !ok {
					_, _ = fmt.Fprintf(os.Stderr, "%s\n", "canceled transaction")
					return err
				}
			}

			txBytes, err := clientCtx.TxConfig.TxEncoder()(tx)
			if err != nil {
				return err
			}

			// broadcast to a Tendermint node
			res, err := clientCtx.BroadcastTx(txBytes)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

const (
	FlagProxyAddress          = "proxy-address"
	FlagImplementationAddress = "implementation-address"
	FlagData                  = "data"
	FlagValue                 = "value"
)

// NewEVMProxyImplmentationUpgrade implements a command handler for submitting a software upgrade with implementation upgrade for existing gensis proxies
func NewEVMProxyImplmentationUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evm-proxy-upgrade (--proxy-address [address]) (--implmentation-address [address]) (--data [data]) (--value [value]) [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Submit an implemntation upgrade for genesis proxy",
		Long:  "Initial proxy implementation upgrade for defined genesis proxy addresses",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()

			proxyAddrStr, err := cmd.Flags().GetString(FlagProxyAddress)
			if err != nil {
				return err
			}
			if !common.IsHexAddress(proxyAddrStr) {
				return fmt.Errorf("%s is not a valid Ethereum address", proxyAddrStr)
			}
			proxyAddr := common.HexToAddress(proxyAddrStr)

			implAddrStr, err := cmd.Flags().GetString(FlagImplementationAddress)
			if err != nil {
				return err
			}
			if !common.IsHexAddress(implAddrStr) {
				return fmt.Errorf("%s is not a valid Ethereum address", implAddrStr)
			}
			implAddr := common.HexToAddress(implAddrStr)

			dataStr, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}
			data, err := hexutil.Decode(dataStr)
			if err != nil {
				return err
			}

			fmt.Println("data", data)

			valueStr, err := cmd.Flags().GetString(FlagValue)
			if err != nil {
				return err
			}
			valueCoin, err := sdk.ParseCoinNormalized(valueStr)
			if err != nil {
				return err
			}
			value := valueCoin.Amount

			depositStr, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}
			content := types.NewUpdateImplmentationProposal(proxyAddr, implAddr, data, &value)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagProxyAddress, "", "proxy address of the contract")
	cmd.Flags().String(FlagImplementationAddress, "", "implementation address which should be used for proxy upgrade")
	cmd.Flags().String(FlagData, "0x", "addition smart contract data for proxy execution (optional)")
	cmd.Flags().String(FlagValue, "0wei", "value of tokens should be used in data execution with payable modifier (optional)")
	cmd.Flags().String(govcli.FlagDeposit, "", "deposit of proposal (optional)")
	cmd.MarkFlagRequired(FlagProxyAddress)
	cmd.MarkFlagRequired(FlagImplementationAddress)

	return cmd
}
