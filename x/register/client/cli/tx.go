package cli

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	stratos "github.com/stratosnet/stratos-chain/types"
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
		RemoveResourceNodeCmd(cdc),
		UpdateResourceNodeCmd(cdc),
		UpdateResourceNodeStakeCmd(cdc),

		CreateIndexingNodeCmd(cdc),
		RemoveIndexingNodeCmd(cdc),
		UpdateIndexingNodeCmd(cdc),
		UpdateIndexingNodeStakeCmd(cdc),
		IndexingNodeRegistrationVoteCmd(cdc),
	)...)

	return registerTxCmd
}

// CreateResourceNodeCmd will create a file upload tx and sign it with the given key.
func CreateResourceNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-resource-node [flags]",
		Short: "create a new resource node",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			if !viper.IsSet(FlagNetworkID) {
				return errors.New("required flag(s) \"network-id\" not set")
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
	cmd.Flags().AddFlagSet(FsNetworkID)
	cmd.Flags().AddFlagSet(FsNodeType)
	cmd.Flags().AddFlagSet(FsDescription)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	//_ = cmd.MarkFlagRequired(FlagNetworkAddr)
	_ = cmd.MarkFlagRequired(FlagNodeType)
	return cmd
}

// UpdateResourceNodeStakeCmd will add/subtract resource node's stake.
func UpdateResourceNodeStakeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-resource-node-stake [flags]",
		Short: "update resource node's stake",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr, msg, err := buildUpdateResourceNodeStakeMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsIncrStake)
	cmd.Flags().AddFlagSet(FsStakeDelta)
	cmd.Flags().AddFlagSet(FsNetworkAddress)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagStakeDelta)
	_ = cmd.MarkFlagRequired(FlagIncrStake)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	return cmd
}

// CreateIndexingNodeCmd will create a file upload tx and sign it with the given key.
func CreateIndexingNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-indexing-node [flags]",
		Short: "create a new indexing node",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			if !viper.IsSet(FlagNetworkID) {
				return errors.New("required flag(s) \"network-id\" not set")
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
	cmd.Flags().AddFlagSet(FsNetworkID)
	cmd.Flags().AddFlagSet(FsDescription)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)

	return cmd
}

// UpdateIndexingNodeStakeCmd will add/subtract indexing node's stake.
func UpdateIndexingNodeStakeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-indexing-node-stake [flags]",
		Short: "update indexing node's stake",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr, msg, err := buildUpdateIndexingNodeStakeMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsIncrStake)
	cmd.Flags().AddFlagSet(FsStakeDelta)
	cmd.Flags().AddFlagSet(FsNetworkAddress)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagStakeDelta)
	_ = cmd.MarkFlagRequired(FlagIncrStake)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	return cmd
}

// makes a new CreateResourceNodeMsg.
func buildCreateResourceNodeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	amountStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoin(amountStr)
	if err != nil {
		return txBldr, nil, err
	}

	networkID := viper.GetString(FlagNetworkID)
	networkAddr, err := stratos.SdsAddressFromBech32(networkID)
	if err != nil {
		return txBldr, nil, err
	}
	ownerAddr := cliCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)
	nodeTypeRef := viper.GetInt(FlagNodeType)

	pubKey, er := stratos.GetPubKeyFromBech32(stratos.Bech32PubKeyTypeSdsP2PPub, pkStr)
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

	// validate nodeTypeRef
	if t := types.NodeType(nodeTypeRef).Type(); t == "UNKNOWN" {
		return txBldr, nil, types.ErrNodeType
	}
	msg := types.NewMsgCreateResourceNode(networkAddr, pubKey, amount, ownerAddr, desc, types.NodeType(nodeTypeRef))
	return txBldr, msg, nil
}

// makes a new UpdateResourceNodeStakeMsg.
func buildUpdateResourceNodeStakeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	stakeDeltaStr := viper.GetString(FlagStakeDelta)
	stakeDelta, err := sdk.ParseCoin(stakeDeltaStr)
	if err != nil {
		return txBldr, nil, err
	}

	incrStakeStr := viper.GetString(FlagIncrStake)
	incrStake, err := strconv.ParseBool(incrStakeStr)
	if err != nil {
		return txBldr, nil, err
	}

	networkAddrStr := viper.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrStr)
	if err != nil {
		return txBldr, nil, err
	}

	ownerAddr := cliCtx.GetFromAddress()

	msg := types.NewMsgUpdateResourceNodeStake(networkAddr, ownerAddr, stakeDelta, incrStake)
	return txBldr, msg, nil
}

// makes a new MsgCreateIndexingNode.
func buildCreateIndexingNodeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	amountStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoin(amountStr)
	if err != nil {
		return txBldr, nil, err
	}

	networkID := viper.GetString(FlagNetworkID)
	networkAddr, err := stratos.SdsAddressFromBech32(networkID)
	if err != nil {
		return txBldr, nil, err
	}
	ownerAddr := cliCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)

	pubKey, err := stratos.GetPubKeyFromBech32(stratos.Bech32PubKeyTypeSdsP2PPub, pkStr)

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
	msg := types.NewMsgCreateIndexingNode(networkAddr, pubKey, amount, ownerAddr, desc)
	return txBldr, msg, nil
}

// makes a new UpdateIndexingNodeStakeMsg.
func buildUpdateIndexingNodeStakeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	stakeDeltaStr := viper.GetString(FlagStakeDelta)
	stakeDelta, err := sdk.ParseCoin(stakeDeltaStr)
	if err != nil {
		return txBldr, nil, err
	}

	incrStakeStr := viper.GetString(FlagIncrStake)
	incrStake, err := strconv.ParseBool(incrStakeStr)
	if err != nil {
		return txBldr, nil, err
	}

	networkAddrStr := viper.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrStr)
	if err != nil {
		return txBldr, nil, err
	}

	ownerAddr := cliCtx.GetFromAddress()

	msg := types.NewMsgUpdateIndexingNodeStake(networkAddr, ownerAddr, stakeDelta, incrStake)
	return txBldr, msg, nil
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

			resourceNodeAddr, err := stratos.SdsAddressFromBech32(args[0])
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

			indexingNodeAddr, err := stratos.SdsAddressFromBech32(args[0])
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

// IndexingNodeRegistrationVoteCmd Indexing node registration need to be approved by 2/3 of existing indexing nodes
func IndexingNodeRegistrationVoteCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexing_node_reg_vote",
		Short: "vote for the registration of a new indexing node",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr, msg, err := buildIndexingNodeRegistrationVoteMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsCandidateNetworkAddress)
	cmd.Flags().AddFlagSet(FsCandidateOwnerAddress)
	cmd.Flags().AddFlagSet(FsOpinion)
	cmd.Flags().AddFlagSet(FsVoterNetworkAddress)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagCandidateNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagCandidateOwnerAddress)
	_ = cmd.MarkFlagRequired(FlagOpinion)
	_ = cmd.MarkFlagRequired(FlagVoterNetworkAddress)
	return cmd
}

func buildIndexingNodeRegistrationVoteMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	candidateNetworkAddrStr := viper.GetString(FlagCandidateNetworkAddress)
	candidateNetworkAddr, err := stratos.SdsAddressFromBech32(candidateNetworkAddrStr)
	if err != nil {
		return txBldr, nil, err
	}
	candidateOwnerAddrStr := viper.GetString(FlagCandidateOwnerAddress)
	candidateOwnerAddr, err := sdk.AccAddressFromBech32(candidateOwnerAddrStr)
	if err != nil {
		return txBldr, nil, err
	}
	opinionVal := viper.GetBool(FlagOpinion)
	opinion := types.VoteOpinionFromBool(opinionVal)
	voterNetworkAddrStr := viper.GetString(FlagVoterNetworkAddress)
	voterNetworkAddr, err := stratos.SdsAddressFromBech32(voterNetworkAddrStr)
	if err != nil {
		return txBldr, nil, err
	}
	voterOwnerAddr := cliCtx.GetFromAddress()

	msg := types.NewMsgIndexingNodeRegistrationVote(candidateNetworkAddr, candidateOwnerAddr, opinion, voterNetworkAddr, voterOwnerAddr)
	return txBldr, msg, nil
}

func UpdateResourceNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-resource-node [flags]",
		Short: "update resource node info",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr, msg, err := buildUpdateResourceNodeMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsNetworkID)
	cmd.Flags().AddFlagSet(FsDescription)
	cmd.Flags().AddFlagSet(FsNodeType)
	cmd.Flags().AddFlagSet(FsNetworkAddress)

	_ = cmd.MarkFlagRequired(FlagNetworkID)
	_ = cmd.MarkFlagRequired(FlagMoniker)
	_ = cmd.MarkFlagRequired(FlagNodeType)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// makes a new MsgUpdateResourceNode.
func buildUpdateResourceNodeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	//networkID := viper.GetString(FlagNetworkID)

	desc := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)

	nodeType := viper.GetInt(FlagNodeType)

	nodeAddrStr := viper.GetString(FlagNetworkAddress)
	nodeAddr, err := stratos.SdsAddressFromBech32(nodeAddrStr)
	if err != nil {
		return txBldr, nil, err
	}

	ownerAddr := cliCtx.GetFromAddress()
	if t := types.NodeType(nodeType).Type(); t == "UNKNOWN" {
		return txBldr, nil, types.ErrNodeType
	}
	msg := types.NewMsgUpdateResourceNode(desc, types.NodeType(nodeType), nodeAddr, ownerAddr)
	return txBldr, msg, nil
}

func UpdateIndexingNodeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-indexing-node [flags]",
		Short: "update indexing node info",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr, msg, err := buildUpdateIndexingNodeMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsNetworkID)
	cmd.Flags().AddFlagSet(FsDescription)
	cmd.Flags().AddFlagSet(FsNetworkAddress)

	_ = cmd.MarkFlagRequired(FlagNetworkID)
	_ = cmd.MarkFlagRequired(FlagMoniker)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// makes a new MsgUpdateIndexingNode.
func buildUpdateIndexingNodeMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	//networkID := viper.GetString(FlagNetworkID)

	desc := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)

	nodeAddrStr := viper.GetString(FlagNetworkAddress)
	nodeAddr, err := stratos.SdsAddressFromBech32(nodeAddrStr)
	if err != nil {
		return txBldr, nil, err
	}

	ownerAddr := cliCtx.GetFromAddress()

	msg := types.NewMsgUpdateIndexingNode(desc, nodeAddr, ownerAddr)
	return txBldr, msg, nil
}
