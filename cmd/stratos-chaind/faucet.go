package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
)

const (
	flagFrom = "from" // optional
	flagTo   = "to"
	flagAmt  = "amt" // denom fixed as stos

	defaultNodeURI        = "tcp://127.0.0.1:26657"
	defaultKeyringBackend = "test"
	defaultHome           = "build/node/stratos-chaincli"
	defaultDenom          = "stos"
	defaultChainId        = "test-chain"
)

// global to load command line args
var faucetArgs = FaucetArgs{}

// struct to hold the command-line args
type FaucetArgs struct {
	from  sdk.AccAddress
	to    sdk.AccAddress
	coins sdk.Coins
}

// AddFaucetCmd returns faucet cobra Command
func AddFaucetCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "faucet",
		Short: "Run a faucet cmd",
		Args:  cobra.RangeArgs(0, 5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			if viper.IsSet(flagFrom) {
				fromAddr := viper.GetString(flagFrom)
				fromAddrBytes, err := sdk.AccAddressFromBech32(fromAddr)
				if err != nil {
					return fmt.Errorf("failed to parse bech32 address fro FROM Address: %w", err)
				}
				faucetArgs.from = fromAddrBytes
			}

			if faucetArgs.from == nil {
				genesis := ctx.Config.GenesisFile()
				faucetArgs.from, err = getFirstAccAddressFromGenesis(cdc, genesis)
				if err != nil {
					return fmt.Errorf("failed to parse genesis: %w", err)
				}
				fmt.Printf("No sender account specified, using account 0 for faucet\n")
			}

			var toTransferAmt int
			if toTransferAmt = viper.GetInt(flagAmt); toTransferAmt <= 0 {
				return fmt.Errorf("Invalid amount in faucet")
			}
			coin := sdk.Coin{Amount: sdk.NewInt(int64(toTransferAmt)), Denom: defaultDenom}
			faucetArgs.coins = sdk.Coins{coin}

			toAddr := viper.GetString(flagTo)
			toAddrBytes, err := sdk.AccAddressFromBech32(toAddr)
			if err != nil {
				return fmt.Errorf("failed to parse bech32 address for To Address: %w", err)
			}
			faucetArgs.to = toAddrBytes

			ctx.Logger.Info("Starting faucet...")

			// start threads
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			viper.Set(flags.FlagBroadcastMode, "async")
			if !viper.IsSet(flags.FlagChainID) {
				viper.Set(flags.FlagChainID, defaultChainId)
			}
			viper.Set(flags.FlagSkipConfirmation, true)
			if !viper.IsSet(flags.FlagKeyringBackend) {
				viper.Set(flags.FlagKeyringBackend, defaultKeyringBackend)
			}
			if !viper.IsSet(flags.FlagNode) {
				viper.Set(flags.FlagNode, defaultNodeURI)
			}
			if !viper.IsSet(flags.FlagHome) {
				viper.Set(flags.FlagHome, defaultHome)
			}
			viper.Set(flags.FlagTrustNode, true)
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, faucetArgs.from.String()).WithCodec(cdc)
			doFaucet(cliCtx, txBldr.WithChainID(viper.GetString(flags.FlagChainID)), faucetArgs.to, faucetArgs.from, coin) // send coin to temp account

			// print stats
			fmt.Println("####################################################################")
			fmt.Println("################        Terminating faucet        ###############")
			fmt.Println("####################################################################")
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagTo, "", "to address")
	cmd.Flags().String(flagAmt, "", "amt to transfer in faucet")
	cmd.Flags().String(flagFrom, "", "from address")
	cmd.Flags().String(flags.FlagChainID, "", "chain id")

	return cmd
}

func getFirstAccAddressFromGenesis(cdc *codec.Codec, genesisFilePath string) (accAddr sdk.AccAddress, err error) {
	var genDoc *tmtypes.GenesisDoc
	if genDoc, err = tmtypes.GenesisDocFromFile(strings.ReplaceAll(genesisFilePath, "cli", "d")); err != nil {
		return nil, fmt.Errorf("error loading genesis doc from %s: %s", genesisFilePath, err.Error())
	}
	var genState map[string]json.RawMessage
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genState); err != nil {
		return nil, fmt.Errorf("error unmarshalling genesis doc %s: %s", genesisFilePath, err.Error())
	}
	var addresses []sdk.AccAddress
	auth.GenesisAccountIterator{}.IterateGenesisAccounts(
		cdc, genState, func(acc exported.Account) (stop bool) {
			addresses = append(addresses, acc.GetAddress())
			return false
		},
	)
	if len(addresses) > 0 {
		return addresses[0], nil
	}
	return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "No account initiated in genesis")
}

func doFaucet(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, to sdk.AccAddress, from sdk.AccAddress, coin sdk.Coin) {
	//// build and sign the transaction, then broadcast to Tendermint
	msg := bank.NewMsgSend(from, to, sdk.Coins{coin})
	fmt.Printf("From: %s, To: %s, Coin: %s\n", msg.FromAddress.String(), msg.ToAddress.String(), msg.Amount.String())
	err := utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
	if err != nil {
		fmt.Println(err)
	}
}
