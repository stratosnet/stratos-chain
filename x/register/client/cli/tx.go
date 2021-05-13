package cli

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	registerTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	registerTxCmd.AddCommand(flags.PostCommands(
		CreateResourceNodeCmd(cdc),
		CreateIndexingNodeCmd(cdc),
		RemoveResourceNodeCmd(cdc),
		RemoveIndexingNodeCmd(cdc),
	)...)

	return registerTxCmd
}

// CreateResourceNodeCmd will create a file upload tx and sign it with the given key.
func CreateResourceNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-resource-node",
		Short: "create a new resource node",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			if !viper.IsSet(FlagNetworkAddr) {
				return errors.New("required flag(s) \"network-addr\" not set")
			}

			if !viper.IsSet(FlagMoniker) {
				return errors.New("required flag(s) \"moniker\" not set")
			}
			txBldr, msg, err := buildCreateResourceNodeMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsNetworkAddr)
	cmd.Flags().AddFlagSet(FsNodeType)
	cmd.Flags().AddFlagSet(FsDescriptionCreate)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	//_ = cmd.MarkFlagRequired(FlagNetworkAddr)
	_ = cmd.MarkFlagRequired(FlagNodeType)

	return cmd
}

// CreateIndexingNodeCmd will create a file upload tx and sign it with the given key.
func CreateIndexingNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-indexing-node",
		Short: "create a new indexing node",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			if !viper.IsSet(FlagNetworkAddr) {
				return errors.New("required flag(s) \"network-adr\" not set")
			}
			if !viper.IsSet(FlagMoniker) {
				return errors.New("required flag(s) \"moniker\" not set")
			}
			txBldr, msg, err := buildCreateIndexingNodeMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsNetworkAddr)
	cmd.Flags().AddFlagSet(FsDescriptionCreate)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	//_ = cmd.MarkFlagRequired(FlagNetworkAddr)

	return cmd
}

// makes a new CreateResourceNodeMsg.
func buildCreateResourceNodeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	amountStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoin(amountStr)
	if err != nil {
		return txBldr, nil, err
	}

	networkAddr := viper.GetString(FlagNetworkAddr)
	ownerAddr := cliCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)
	nodeTypeRef := viper.GetInt(FlagNodeType)

	pk, er := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pkStr)
	if er != nil {
		return txBldr, nil, err
	}

	desc := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)
	if !ValueInSlice(nodeTypeRef, types.NodeTypes) {
		return txBldr, nil, types.ErrNodeType
	}
	msg := types.NewMsgCreateResourceNode(networkAddr, pk, amount, ownerAddr, desc, fmt.Sprintf("%d: %s", nodeTypeRef, types.NodeTypesMap[nodeTypeRef]))

	return txBldr, msg, nil
}

func ValueInSlice(v int, list []int) bool {
	for _, b := range list {
		if b == v {
			return true
		}
	}
	return false
}

func RemoveResourceNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-resource-node [resource_node_address] [owner_address]",
		Args:  cobra.ExactArgs(2),
		Short: "remove resource node",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[1]).WithCodec(cdc)

			resourceNodeAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			ownerAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgRemoveResourceNode(resourceNodeAddr, ownerAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func RemoveIndexingNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-indexing-node [indexing_node_address] [owner_address]",
		Args:  cobra.ExactArgs(2),
		Short: "remove indexing node",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[1]).WithCodec(cdc)

			indexingNodeAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			ownerAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgRemoveIndexingNode(indexingNodeAddr, ownerAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// makes a new MsgCreateIndexingNode.
func buildCreateIndexingNodeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	amountStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoin(amountStr)
	if err != nil {
		return txBldr, nil, err
	}

	networkAddr := viper.GetString(FlagNetworkAddr)
	ownerAddr := cliCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)

	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	desc := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)
	msg := types.NewMsgCreateIndexingNode(networkAddr, pk, amount, ownerAddr, desc)
	return txBldr, msg, nil
}
