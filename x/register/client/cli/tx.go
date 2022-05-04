package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	registerTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	registerTxCmd.AddCommand(
		CreateResourceNodeCmd(),
		RemoveResourceNodeCmd(),
		UpdateResourceNodeCmd(),
		UpdateResourceNodeStakeCmd(),

		CreateIndexingNodeCmd(),
		RemoveIndexingNodeCmd(),
		UpdateIndexingNodeCmd(),
		UpdateIndexingNodeStakeCmd(),
		IndexingNodeRegistrationVoteCmd(),
	)

	return registerTxCmd
}

// CreateResourceNodeCmd will create a file upload tx and sign it with the given key.
func CreateResourceNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-resource-node [flags]",
		Short: "create a new resource node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildCreateResourceNodeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsNetworkAddress)
	cmd.Flags().AddFlagSet(FsNodeType)
	cmd.Flags().AddFlagSet(FsDescription)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagNodeType)
	return cmd
}

// CreateIndexingNodeCmd will create a file upload tx and sign it with the given key.
func CreateIndexingNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-indexing-node [flags]",
		Short: "create a new indexing node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildCreateIndexingNodeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsNetworkAddress)
	cmd.Flags().AddFlagSet(FsDescription)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)

	return cmd
}

func RemoveResourceNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-resource-node [resource_node_address] [owner_address]",
		Args:  cobra.ExactArgs(2),
		Short: "remove resource node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			resourceNodeAddr, err := stratos.SdsAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			ownerAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgRemoveResourceNode(resourceNodeAddr, ownerAddr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func RemoveIndexingNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-indexing-node [indexing_node_address] [owner_address]",
		Args:  cobra.ExactArgs(2),
		Short: "remove indexing node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			resourceNodeAddr, err := stratos.SdsAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			ownerAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgRemoveIndexingNode(resourceNodeAddr, ownerAddr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func UpdateResourceNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-resource-node [flags]",
		Short: "update resource node info",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildUpdateResourceNodeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FsNetworkAddress)
	cmd.Flags().AddFlagSet(FsDescription)
	cmd.Flags().AddFlagSet(FsNodeType)
	cmd.Flags().AddFlagSet(FsNetworkAddress)

	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagMoniker)
	_ = cmd.MarkFlagRequired(FlagNodeType)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func UpdateIndexingNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-indexing-node [flags]",
		Short: "update indexing node info",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildUpdateIndexingNodeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FsNetworkAddress)
	cmd.Flags().AddFlagSet(FsDescription)
	cmd.Flags().AddFlagSet(FsNetworkAddress)

	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagMoniker)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// UpdateResourceNodeStakeCmd will add/subtract resource node's stake.
func UpdateResourceNodeStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-resource-node-stake [flags]",
		Short: "update resource node's stake",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildUpdateResourceNodeStakeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
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

// UpdateIndexingNodeStakeCmd will add/subtract indexing node's stake.
func UpdateIndexingNodeStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-indexing-node-stake [flags]",
		Short: "update indexing node's stake",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildUpdateIndexingNodeStakeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
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

// IndexingNodeRegistrationVoteCmd Indexing node registration need to be approved by 2/3 of existing indexing nodes
func IndexingNodeRegistrationVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexing_node_reg_vote",
		Short: "vote for the registration of a new indexing node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildIndexingNodeRegistrationVoteMsg(clientCtx, txf, cmd.Flags())

			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
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

// makes a new CreateResourceNodeMsg.
func newBuildCreateResourceNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateResourceNode, error) {
	fAmount, _ := fs.GetString(FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	networkAddrstr := viper.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrstr)
	if err != nil {
		return txf, nil, err
	}
	ownerAddr := clientCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)
	nodeTypeRef := viper.GetInt(FlagNodeType)

	pubKey, er := stratos.GetPubKeyFromBech32(stratos.Bech32PubKeyTypeSdsP2PPub, pkStr)
	if er != nil {
		return txf, nil, err
	}

	desc := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)

	// validate nodeTypeRef
	newNodeType := types.NodeType(nodeTypeRef)
	if t := newNodeType.Type(); t == "UNKNOWN" {
		return txf, nil, types.ErrNodeType
	}
	msg, er := types.NewMsgCreateResourceNode(networkAddr, pubKey, amount, ownerAddr, &desc, &newNodeType)
	if er != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

// makes a new MsgCreateIndexingNode.
func newBuildCreateIndexingNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateResourceNode, error) {
	fAmount, _ := fs.GetString(FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	networkAddrstr := viper.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrstr)
	if err != nil {
		return txf, nil, err
	}
	ownerAddr := clientCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)

	pubKey, er := stratos.GetPubKeyFromBech32(stratos.Bech32PubKeyTypeSdsP2PPub, pkStr)
	if er != nil {
		return txf, nil, err
	}

	desc := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)
	msg, er := types.NewMsgCreateIndexingNode(networkAddr, pubKey, amount, ownerAddr, &desc)
	if er != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

// makes a new MsgUpdateResourceNode.
func newBuildUpdateResourceNodeMsg(clientCtx client.Context, txf tx.Factory, _ *flag.FlagSet) (tx.Factory, *types.MsgUpdateResourceNode, error) {
	networkAddrstr := viper.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrstr)
	if err != nil {
		return txf, nil, err
	}
	ownerAddr := clientCtx.GetFromAddress()
	nodeTypeRef := viper.GetInt(FlagNodeType)

	desc := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)

	newNodeType := types.NodeType(nodeTypeRef)
	if t := newNodeType.Type(); t == "UNKNOWN" {
		return txf, nil, types.ErrNodeType
	}
	msg := types.NewMsgUpdateResourceNode(desc, newNodeType, networkAddr, ownerAddr)
	return txf, msg, nil
}

// makes a new MsgUpdateIndexingNode.
func newBuildUpdateIndexingNodeMsg(clientCtx client.Context, txf tx.Factory, _ *flag.FlagSet) (tx.Factory, *types.MsgUpdateIndexingNode, error) {
	desc := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)

	networkAddrstr := viper.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrstr)
	if err != nil {
		return txf, nil, err
	}
	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgUpdateIndexingNode(desc, networkAddr, ownerAddr)
	return txf, msg, nil
}

// newBuildUpdateResourceNodeStakeMsg makes a new UpdateResourceNodeStakeMsg.
func newBuildUpdateResourceNodeStakeMsg(clientCtx client.Context, txf tx.Factory, _ *flag.FlagSet) (tx.Factory, *types.MsgUpdateResourceNodeStake, error) {
	stakeDeltaStr := viper.GetString(FlagStakeDelta)
	stakeDelta, err := sdk.ParseCoinNormalized(stakeDeltaStr)
	if err != nil {
		return txf, nil, err
	}

	incrStakeStr := viper.GetString(FlagIncrStake)
	incrStake, err := strconv.ParseBool(incrStakeStr)
	if err != nil {
		return txf, nil, err
	}

	networkAddrStr := viper.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgUpdateResourceNodeStake(networkAddr, ownerAddr, &stakeDelta, incrStake)
	return txf, msg, nil
}

// newBuildUpdateIndexingNodeStakeMsg makes a new UpdateIndexingNodeStakeMsg.
func newBuildUpdateIndexingNodeStakeMsg(clientCtx client.Context, txf tx.Factory, _ *flag.FlagSet) (tx.Factory, *types.MsgUpdateIndexingNodeStake, error) {
	stakeDeltaStr := viper.GetString(FlagStakeDelta)
	stakeDelta, err := sdk.ParseCoinNormalized(stakeDeltaStr)
	if err != nil {
		return txf, nil, err
	}

	incrStakeStr := viper.GetString(FlagIncrStake)
	incrStake, err := strconv.ParseBool(incrStakeStr)
	if err != nil {
		return txf, nil, err
	}

	networkAddrStr := viper.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgUpdateIndexingNodeStake(networkAddr, ownerAddr, &stakeDelta, incrStake)
	return txf, msg, nil
}

func newBuildIndexingNodeRegistrationVoteMsg(clientCtx client.Context, txf tx.Factory, _ *flag.FlagSet) (tx.Factory, *types.MsgIndexingNodeRegistrationVote, error) {
	candidateNetworkAddrStr := viper.GetString(FlagCandidateNetworkAddress)
	candidateNetworkAddr, err := stratos.SdsAddressFromBech32(candidateNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}
	candidateOwnerAddrStr := viper.GetString(FlagCandidateOwnerAddress)
	candidateOwnerAddr, err := sdk.AccAddressFromBech32(candidateOwnerAddrStr)
	if err != nil {
		return txf, nil, err
	}
	opinionVal := viper.GetBool(FlagOpinion)
	//opinion := types.VoteOpinionFromBool(opinionVal)
	voterNetworkAddrStr := viper.GetString(FlagVoterNetworkAddress)
	voterNetworkAddr, err := stratos.SdsAddressFromBech32(voterNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}
	voterOwnerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgIndexingNodeRegistrationVote(candidateNetworkAddr, candidateOwnerAddr, opinionVal, voterNetworkAddr, voterOwnerAddr)
	return txf, msg, nil
}
