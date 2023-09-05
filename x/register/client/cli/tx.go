package cli

import (
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
		UpdateResourceNodeDepositCmd(),

		CreateMetaNodeCmd(),
		RemoveMetaNodeCmd(),
		UpdateMetaNodeCmd(),
		UpdateMetaNodeDepositCmd(),
		MetaNodeRegistrationVoteCmd(),
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

// CreateMetaNodeCmd will create a file upload tx and sign it with the given key.
func CreateMetaNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-meta-node [flags]",
		Short: "create a new meta node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildCreateMetaNodeMsg(clientCtx, txf, cmd.Flags())
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
	cmd.Flags().AddFlagSet(flagSetBeneficiaryAddress())

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

func RemoveMetaNodeCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use: "remove-meta-node [flag]",
		//Args:  cobra.ExactArgs(1),
		Short: "remove meta node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildRemoveMetaNodeMsg(clientCtx, txf, cmd.Flags())
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

func UpdateMetaNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-meta-node [flags]",
		Short: "update meta node info",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildUpdateMetaNodeMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())
	cmd.Flags().AddFlagSet(flagSetBeneficiaryAddress())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	//_ = cmd.MarkFlagRequired(FlagMoniker)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// UpdateResourceNodeDepositCmd will add/subtract resource node's deposit.
func UpdateResourceNodeDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-resource-node-deposit [flags]",
		Short: "update resource node's deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildUpdateResourceNodeDepositMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetDepositUpdate())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagDepositDelta)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	return cmd
}

// UpdateMetaNodeDepositCmd will add/subtract meta node's deposit.
func UpdateMetaNodeDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-meta-node-deposit [flags]",
		Short: "update meta node's deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildUpdateMetaNodeDepositMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetNetworkAddress())
	cmd.Flags().AddFlagSet(flagSetDepositUpdate())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagDepositDelta)
	_ = cmd.MarkFlagRequired(FlagIncrDeposit)
	_ = cmd.MarkFlagRequired(FlagNetworkAddress)
	return cmd
}

// MetaNodeRegistrationVoteCmd Meta node registration need to be approved by 2/3 of existing meta nodes
func MetaNodeRegistrationVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meta-node-reg-vote [flags]",
		Short: "vote for the registration of a new meta node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildMetaNodeRegistrationVoteMsg(clientCtx, txf, cmd.Flags())

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

	nodeTypeVal, err := fs.GetUint32(FlagNodeType)
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

	// validate nodeTypeVal
	nodeType := types.NodeType(nodeTypeVal)
	if t := nodeType.Type(); t == "UNKNOWN" {
		return txf, nil, types.ErrNodeType
	}
	msg, er := types.NewMsgCreateResourceNode(networkAddr, pubKey, amount, ownerAddr, description, nodeTypeVal)
	if er != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

// makes a new MsgCreateMetaNode.
func newBuildCreateMetaNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateMetaNode, error) {
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
	msg, er := types.NewMsgCreateMetaNode(networkAddr, pubKey, amount, ownerAddr, description)
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

	nodeTypeVal, err := fs.GetUint32(FlagNodeType)
	if err != nil {
		return txf, nil, types.ErrInvalidNodeType
	}

	// validate nodeTypeVal
	nodeType := types.NodeType(nodeTypeVal)
	if t := nodeType.Type(); t == "UNKNOWN" {
		return txf, nil, types.ErrNodeType
	}
	msg := types.NewMsgUpdateResourceNode(description, nodeTypeVal, networkAddr, ownerAddr)
	return txf, msg, nil
}

// makes a new MsgUpdateMetaNode.
func newBuildUpdateMetaNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgUpdateMetaNode, error) {
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

	msg := types.NewMsgUpdateMetaNode(description, networkAddr, ownerAddr)
	return txf, msg, nil
}

// newBuildUpdateResourceNodeDepositMsg makes a new MsgUpdateResourceNodeDeposit.
func newBuildUpdateResourceNodeDepositMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (
	tx.Factory, *types.MsgUpdateResourceNodeDeposit, error) {

	depositDeltaStr, err := fs.GetString(FlagDepositDelta)
	if err != nil {
		return txf, nil, err
	}
	depositDelta, err := sdk.ParseCoinNormalized(depositDeltaStr)
	if err != nil {
		return txf, nil, err
	}

	networkAddrStr, _ := fs.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgUpdateResourceNodeDeposit(networkAddr, ownerAddr, depositDelta)
	return txf, msg, nil
}

// newBuildUpdateMetaNodeDepositMsg makes a new MsgUpdateMetaNodeDeposit.
func newBuildUpdateMetaNodeDepositMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgUpdateMetaNodeDeposit, error) {
	depositDeltaStr, err := fs.GetString(FlagDepositDelta)
	if err != nil {
		return txf, nil, err
	}
	depositDelta, err := sdk.ParseCoinNormalized(depositDeltaStr)
	if err != nil {
		return txf, nil, err
	}

	networkAddrStr, _ := fs.GetString(FlagNetworkAddress)
	networkAddr, err := stratos.SdsAddressFromBech32(networkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgUpdateMetaNodeDeposit(networkAddr, ownerAddr, depositDelta)
	return txf, msg, nil
}

func newBuildMetaNodeRegistrationVoteMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgMetaNodeRegistrationVote, error) {
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

	msg := types.NewMsgMetaNodeRegistrationVote(candidateNetworkAddr, candidateOwnerAddr, opinionVal, voterNetworkAddr, voterOwnerAddr)
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

func newBuildRemoveMetaNodeMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgRemoveMetaNode, error) {
	flagNetworkAddrStr, err := fs.GetString(FlagNetworkAddress)
	if err != nil {
		return txf, nil, err
	}
	networkAddr, err := stratos.SdsAddressFromBech32(flagNetworkAddrStr)
	if err != nil {
		return txf, nil, err
	}

	ownerAddr := clientCtx.GetFromAddress()

	msg := types.NewMsgRemoveMetaNode(networkAddr, ownerAddr)

	return txf, msg, nil
}
