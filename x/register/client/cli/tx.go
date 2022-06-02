package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
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

	cmd.Flags().AddFlagSet(flagSetPublicKey())
	cmd.Flags().AddFlagSet(flagSetAmount())
	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetNodeType())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagNodeType)
	_ = cmd.MarkFlagRequired(FlagMoniker)
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
	cmd.Flags().AddFlagSet(flagSetPublicKey())
	cmd.Flags().AddFlagSet(flagSetAmount())
	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}

func RemoveResourceNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "remove-resource-node [flag]",
		//Args:  cobra.ExactArgs(1),
		Short: "remove resource node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildRemoveResourceNodeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(flagSetNetworkAddress())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)

	return cmd
}

func RemoveIndexingNodeCmd() *cobra.Command {
	//cmd := &cobra.Command{
	//	Use:   "remove-indexing-node [indexing_node_address]",
	//	Args:  cobra.ExactArgs(2),
	//	Short: "remove indexing node",
	//	RunE: func(cmd *cobra.Command, args []string) error {
	//		clientCtx, err := client.GetClientTxContext(cmd)
	//		if err != nil {
	//			return err
	//		}
	//
	//		indexingNodeAddr, err := stratos.SdsAddressFromBech32(args[0])
	//		if err != nil {
	//			return err
	//		}
	//		ownerAddr := clientCtx.GetFromAddress()
	//		//ownerAddr, err := sdk.AccAddressFromBech32(args[1])
	//		//if err != nil {
	//		//	return err
	//		//}
	//
	//		msg := types.NewMsgRemoveIndexingNode(indexingNodeAddr, ownerAddr)
	//
	//		return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	//	},
	//}
	//_ = cmd.MarkFlagRequired(flags.FlagFrom)
	//
	//flags.AddTxFlagsToCmd(cmd)
	//
	//return cmd
	cmd := &cobra.Command{
		Use: "remove-indexing-node [flag]",
		//Args:  cobra.ExactArgs(1),
		Short: "remove indeixng node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildRemoveIndexingNodeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(flagSetNetworkAddress())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)

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

	cmd.Flags().AddFlagSet(flagSetNodeType())
	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	//_ = cmd.MarkFlagRequired(FlagMoniker)
	//_ = cmd.MarkFlagRequired(FlagNodeType)
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

	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	//_ = cmd.MarkFlagRequired(FlagMoniker)
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

	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetStakeUpdate())

	flags.AddTxFlagsToCmd(cmd)

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

	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetStakeUpdate())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagStakeDelta)
	_ = cmd.MarkFlagRequired(FlagIncrStake)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	return cmd
}

// IndexingNodeRegistrationVoteCmd Indexing node registration need to be approved by 2/3 of existing indexing nodes
func IndexingNodeRegistrationVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexing_node_reg_vote [flags]",
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

	cmd.Flags().AddFlagSet(flagSetVoting())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagCandidateNetworkAddress)
	_ = cmd.MarkFlagRequired(FlagCandidateOwnerAddress)
	_ = cmd.MarkFlagRequired(FlagOpinion)
	_ = cmd.MarkFlagRequired(FlagVoterNetworkAddress)
	return cmd
}

// makes a new CreateResourceNodeMsg.
func newBuildCreateResourceNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateResourceNode, error) {
	flagAmountStr, err := fs.GetString(FlagAmount)
	if err != nil {
		return txf, nil, err
	}
	amount, err := sdk.ParseCoinNormalized(flagAmountStr)
	if err != nil {
		return txf, nil, err
	}

	flagNetworkAddrStr, err := fs.GetString(FlagNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	networkAddr, err := stratos.SdsAddressFromBech32(flagNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	pkStr, err := fs.GetString(FlagPubKey)
	if err != nil {
		return txf, nil, err
	}

	pubKey, err := stratos.SdsPubKeyFromBech32(pkStr)
	if err != nil {
		return txf, nil, err
	}

	nodeTypeRef, err := fs.GetInt(FlagNodeType)
	if err != nil {
		return txf, nil, err
	}

	moniker, _ := fs.GetString(FlagMoniker)
	identity, _ := fs.GetString(FlagIdentity)
	website, _ := fs.GetString(FlagWebsite)
	security, _ := fs.GetString(FlagSecurityContact)
	details, _ := fs.GetString(FlagDetails)
	description := types.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)

	// validate nodeTypeRef
	newNodeType := types.NodeType(nodeTypeRef)
	if t := newNodeType.Type(); t == "UNKNOWN" {
		return txf, nil, types.ErrNodeType
	}
	msg, er := types.NewMsgCreateResourceNode(networkAddr, pubKey, amount, ownerAddr, &description, strconv.Itoa(nodeTypeRef))
	if er != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

// makes a new MsgCreateIndexingNode.
func newBuildCreateIndexingNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateIndexingNode, error) {
	flagAmountStr, err := fs.GetString(FlagAmount)
	if err != nil {
		return txf, nil, err
	}
	amount, err := sdk.ParseCoinNormalized(flagAmountStr)
	if err != nil {
		return txf, nil, err
	}

	flagNetworkAddrStr, err := fs.GetString(FlagNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	networkAddr, err := stratos.SdsAddressFromBech32(flagNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	pkStr, err := fs.GetString(FlagPubKey)
	if err != nil {
		return txf, nil, err
	}
	pubKey, er := stratos.SdsPubKeyFromBech32(pkStr)
	if er != nil {
		return txf, nil, err
	}

	moniker, _ := fs.GetString(FlagMoniker)
	identity, _ := fs.GetString(FlagIdentity)
	website, _ := fs.GetString(FlagWebsite)
	security, _ := fs.GetString(FlagSecurityContact)
	details, _ := fs.GetString(FlagDetails)
	description := types.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)
	msg, er := types.NewMsgCreateIndexingNode(networkAddr, pubKey, amount, ownerAddr, &description)
	if er != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

// makes a new MsgUpdateResourceNode.
func newBuildUpdateResourceNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgUpdateResourceNode, error) {
	flagNetworkAddrStr, err := fs.GetString(FlagNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	networkAddr, err := stratos.SdsAddressFromBech32(flagNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	moniker, _ := fs.GetString(FlagMoniker)
	identity, _ := fs.GetString(FlagIdentity)
	website, _ := fs.GetString(FlagWebsite)
	security, _ := fs.GetString(FlagSecurityContact)
	details, _ := fs.GetString(FlagDetails)
	description := types.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)

	nodeTypeRef, err := fs.GetInt(FlagNodeType)
	if err != nil {
		return txf, nil, types.ErrInvalidNodeType
	}

	if nodeTypeRef > 7 || nodeTypeRef < 0 {
		return txf, nil, types.ErrInvalidNodeType
	}
	nodeTypeStr := strconv.Itoa(nodeTypeRef)
	msg := types.NewMsgUpdateResourceNode(description, nodeTypeStr, networkAddr, ownerAddr)
	return txf, msg, nil
}

// makes a new MsgUpdateIndexingNode.
func newBuildUpdateIndexingNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgUpdateIndexingNode, error) {
	flagNetworkAddrStr, err := fs.GetString(FlagNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	networkAddr, err := stratos.SdsAddressFromBech32(flagNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	moniker, _ := fs.GetString(FlagMoniker)
	identity, _ := fs.GetString(FlagIdentity)
	website, _ := fs.GetString(FlagWebsite)
	security, _ := fs.GetString(FlagSecurityContact)
	details, _ := fs.GetString(FlagDetails)
	description := types.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)

	msg := types.NewMsgUpdateIndexingNode(description, networkAddr, ownerAddr)
	return txf, msg, nil
}

// newBuildUpdateResourceNodeStakeMsg makes a new UpdateResourceNodeStakeMsg.
func newBuildUpdateResourceNodeStakeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgUpdateResourceNodeStake, error) {
	stakeDeltaStr, err := fs.GetString(FlagStakeDelta)
	if err != nil {
		return txf, nil, err
	}
	stakeDelta, err := sdk.ParseCoinNormalized(stakeDeltaStr)
	if err != nil {
		return txf, nil, err
	}

	incrStakeStr, err := fs.GetString(FlagIncrStake)
	if err != nil {
		return txf, nil, err
	}
	incrStake, err := strconv.ParseBool(incrStakeStr)
	if err != nil {
		return txf, nil, err
	}

	networkAddrStr, _ := fs.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgUpdateResourceNodeStake(networkAddr, ownerAddr, &stakeDelta, incrStake)
	return txf, msg, nil
}

// newBuildUpdateIndexingNodeStakeMsg makes a new UpdateIndexingNodeStakeMsg.
func newBuildUpdateIndexingNodeStakeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgUpdateIndexingNodeStake, error) {
	stakeDeltaStr, err := fs.GetString(FlagStakeDelta)
	if err != nil {
		return txf, nil, err
	}
	stakeDelta, err := sdk.ParseCoinNormalized(stakeDeltaStr)
	if err != nil {
		return txf, nil, err
	}

	incrStakeStr, err := fs.GetString(FlagIncrStake)
	if err != nil {
		return txf, nil, err
	}
	incrStake, err := strconv.ParseBool(incrStakeStr)
	if err != nil {
		return txf, nil, err
	}

	networkAddrStr, _ := fs.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgUpdateIndexingNodeStake(networkAddr, ownerAddr, &stakeDelta, incrStake)
	return txf, msg, nil
}

func newBuildIndexingNodeRegistrationVoteMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgIndexingNodeRegistrationVote, error) {
	candidateNetworkAddrStr, err := fs.GetString(FlagCandidateNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	candidateNetworkAddr, err := stratos.SdsAddressFromBech32(candidateNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	candidateOwnerAddrStr, err := fs.GetString(FlagCandidateOwnerAddress)
	if err != nil {
		return txf, nil, err
	}
	candidateOwnerAddr, err := sdk.AccAddressFromBech32(candidateOwnerAddrStr)
	if err != nil {
		return txf, nil, err
	}

	opinionVal, err := fs.GetBool(FlagOpinion)
	if err != nil {
		return txf, nil, err
	}
	//opinion := types.VoteOpinionFromBool(opinionVal)
	voterNetworkAddrStr, err := fs.GetString(FlagVoterNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	voterNetworkAddr, err := stratos.SdsAddressFromBech32(voterNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	voterOwnerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgIndexingNodeRegistrationVote(candidateNetworkAddr, candidateOwnerAddr, opinionVal, voterNetworkAddr, voterOwnerAddr)
	return txf, msg, nil
}

func newBuildRemoveResourceNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgRemoveResourceNode, error) {
	flagNetworkAddrStr, err := fs.GetString(FlagNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	networkAddr, err := stratos.SdsAddressFromBech32(flagNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgRemoveResourceNode(networkAddr, ownerAddr)

	return txf, msg, nil
}

func newBuildRemoveIndexingNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgRemoveIndexingNode, error) {
	flagNetworkAddrStr, err := fs.GetString(FlagNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	networkAddr, err := stratos.SdsAddressFromBech32(flagNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgRemoveIndexingNode(networkAddr, ownerAddr)

	return txf, msg, nil
}
