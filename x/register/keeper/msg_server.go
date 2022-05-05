package keeper

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) HandleMsgCreateResourceNode(goCtx context.Context, msg *types.MsgCreateResourceNode) (*types.MsgCreateResourceNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check to see if the pubkey or sender has been registered before
	pk, err := stratos.GetPubKeyFromBech32(stratos.Bech32PubKeyTypeSdsP2PPub, msg.PubKey.String())
	if err != nil {
		return nil, err
	}

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddr)
	if err != nil {
		return &types.MsgCreateResourceNodeResponse{}, err
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgCreateResourceNodeResponse{}, err
	}

	if _, found := k.GetResourceNode(ctx, networkAddr); found {
		ctx.Logger().Error("Resource node already exist")
		return nil, types.ErrResourceNodePubKeyExists
	}
	if msg.Value.Denom != k.BondDenom(ctx) {
		return nil, types.ErrBadDenom
	}
	nodeType, err := strconv.ParseUint(msg.NodeType, 10, 8)
	if err != nil {
		return &types.MsgCreateResourceNodeResponse{}, err
	}
	ozoneLimitChange, err := k.RegisterResourceNode(ctx, networkAddr, pk, ownerAddress, *msg.Description, types.NodeType(nodeType), msg.Value)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateResourceNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, pk.String()),
			sdk.NewAttribute(types.AttributeKeyPubKey, hex.EncodeToString(pk.Bytes())),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.String()),
			sdk.NewAttribute(types.AttributeKeyInitialStake, msg.Value.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgCreateResourceNodeResponse{}, nil
}

func (k msgServer) HandleMsgCreateIndexingNode(goCtx context.Context, msg *types.MsgCreateIndexingNode) (*types.MsgCreateIndexingNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check to see if the pubkey or sender has been registered before
	pk, err := stratos.GetPubKeyFromBech32(stratos.Bech32PubKeyTypeSdsP2PPub, msg.PubKey.String())
	if err != nil {
		return nil, err
	}

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddr)
	if err != nil {
		return &types.MsgCreateIndexingNodeResponse{}, err
	}

	if _, found := k.GetIndexingNode(ctx, networkAddr); found {
		ctx.Logger().Error("Indexing node already exist")
		return nil, types.ErrIndexingNodePubKeyExists
	}
	if msg.Value.Denom != k.BondDenom(ctx) {
		return nil, types.ErrBadDenom
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgCreateIndexingNodeResponse{}, err
	}

	ozoneLimitChange, err := k.RegisterIndexingNode(ctx, networkAddr, pk, ownerAddress, *msg.Description, msg.Value)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateIndexingNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, pk.String()),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgCreateIndexingNodeResponse{}, nil
}

func (k msgServer) HandleMsgRemoveResourceNode(goCtx context.Context, msg *types.MsgRemoveResourceNode) (*types.MsgRemoveResourceNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	p2pAddress, err := stratos.SdsAddressFromBech32(msg.ResourceNodeAddress)
	if err != nil {
		return &types.MsgRemoveResourceNodeResponse{}, err
	}
	resourceNode, found := k.GetResourceNode(ctx, p2pAddress)
	if !found {
		return nil, types.ErrNoResourceNodeFound
	}
	if resourceNode.GetStatus() == stakingtypes.Unbonding {
		return nil, types.ErrUnbondingNode
	}

	ozoneLimitChange, completionTime, err := k.UnbondResourceNode(ctx, resourceNode, resourceNode.Tokens)
	if err != nil {
		return nil, err
	}

	//completionTimeBz := amino.MustMarshalBinaryLengthPrefixed(completionTime)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondingResourceNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyResourceNode, msg.ResourceNodeAddress),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.Neg().String()),
			sdk.NewAttribute(types.AttributeKeyStakeToRemove, resourceNode.Tokens.String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingMatureTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})

	return &types.MsgRemoveResourceNodeResponse{}, nil
}

func (k msgServer) HandleMsgRemoveIndexingNode(goCtx context.Context, msg *types.MsgRemoveIndexingNode) (*types.MsgRemoveIndexingNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	p2pAddress, err := stratos.SdsAddressFromBech32(msg.IndexingNodeAddress)
	if err != nil {
		return &types.MsgRemoveIndexingNodeResponse{}, err
	}
	indexingNode, found := k.GetIndexingNode(ctx, p2pAddress)
	if !found {
		return nil, types.ErrNoIndexingNodeFound
	}

	if indexingNode.GetStatus() == stakingtypes.Unbonding {
		return nil, types.ErrUnbondingNode
	}

	ozoneLimitChange, completionTime, err := k.UnbondIndexingNode(ctx, indexingNode, indexingNode.Tokens)
	if err != nil {
		return nil, err
	}

	//completionTimeBz := amino.MustMarshalBinaryLengthPrefixed(completionTime)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondingIndexingNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyIndexingNode, msg.IndexingNodeAddress),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.Neg().String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingMatureTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})

	return &types.MsgRemoveIndexingNodeResponse{}, nil
}

func (k msgServer) HandleMsgIndexingNodeRegistrationVote(goCtx context.Context, msg *types.MsgIndexingNodeRegistrationVote) (*types.MsgIndexingNodeRegistrationVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	candidateNetworkAddress, err := stratos.SdsAddressFromBech32(msg.CandidateNetworkAddress)
	if err != nil {
		return &types.MsgIndexingNodeRegistrationVoteResponse{}, err
	}

	nodeToApprove, found := k.GetIndexingNode(ctx, candidateNetworkAddress)
	if !found {
		return nil, types.ErrNoIndexingNodeFound
	}
	ownerAddress, err := stratos.SdsAddressFromBech32(nodeToApprove.OwnerAddress)
	if err != nil {
		return &types.MsgIndexingNodeRegistrationVoteResponse{}, err
	}
	if !ownerAddress.Equals(candidateNetworkAddress) {
		return nil, types.ErrInvalidOwnerAddr
	}
	voterNetworkAddress, err := stratos.SdsAddressFromBech32(msg.VoterNetworkAddress)
	if err != nil {
		return &types.MsgIndexingNodeRegistrationVoteResponse{}, err
	}
	voter, found := k.GetIndexingNode(ctx, voterNetworkAddress)
	if !found {
		return nil, types.ErrInvalidVoterAddr
	}

	candidateOwnerAddress, err := sdk.AccAddressFromBech32(msg.CandidateOwnerAddress)
	if err != nil {
		return &types.MsgIndexingNodeRegistrationVoteResponse{}, err
	}

	if !(voter.Status == stakingtypes.Bonded) || voter.Suspend {
		return nil, types.ErrInvalidVoterStatus
	}

	nodeStatus, err := k.HandleVoteForIndexingNodeRegistration(ctx, candidateNetworkAddress, candidateOwnerAddress, types.VoteOpinion(msg.Opinion), voterNetworkAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIndexingNodeRegistrationVote,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.VoterOwnerAddress),
			sdk.NewAttribute(types.AttributeKeyVoterNetworkAddress, msg.VoterNetworkAddress),
			sdk.NewAttribute(types.AttributeKeyCandidateNetworkAddress, msg.CandidateNetworkAddress),
			sdk.NewAttribute(types.AttributeKeyCandidateStatus, nodeStatus.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.VoterOwnerAddress),
		),
	})

	return &types.MsgIndexingNodeRegistrationVoteResponse{}, nil
}

func (k msgServer) HandleMsgUpdateResourceNode(goCtx context.Context, msg *types.MsgUpdateResourceNode) (*types.MsgUpdateResourceNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeResponse{}, err
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeResponse{}, err
	}
	nodeType, err := strconv.ParseUint(msg.NodeType, 10, 8)
	if err != nil {
		return &types.MsgUpdateResourceNodeResponse{}, err
	}
	err = k.UpdateResourceNode(ctx, msg.Description, types.NodeType(nodeType), networkAddr, ownerAddress)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateResourceNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgUpdateResourceNodeResponse{}, nil
}

func (k msgServer) HandleMsgUpdateResourceNodeStake(goCtx context.Context, msg *types.MsgUpdateResourceNodeStake) (*types.MsgUpdateResourceNodeStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeStakeResponse{}, err
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeStakeResponse{}, err
	}

	ozoneLimitChange, completionTime, err := k.UpdateResourceNodeStake(ctx, networkAddr, ownerAddress, *msg.StakeDelta, msg.IncrStake)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateResourceNodeStake,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyIncrStakeBool, strconv.FormatBool(msg.IncrStake)),
			sdk.NewAttribute(types.AttributeKeyStakeDelta, msg.StakeDelta.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingMatureTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgUpdateResourceNodeStakeResponse{}, nil
}

func (k msgServer) HandleMsgUpdateIndexingNode(goCtx context.Context, msg *types.MsgUpdateIndexingNode) (*types.MsgUpdateIndexingNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateIndexingNodeResponse{}, err
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateIndexingNodeResponse{}, err
	}

	err = k.UpdateIndexingNode(ctx, msg.Description, networkAddr, ownerAddress)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateIndexingNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgUpdateIndexingNodeResponse{}, nil
}

func (k msgServer) HandleMsgUpdateIndexingNodeStake(goCtx context.Context, msg *types.MsgUpdateIndexingNodeStake) (*types.MsgUpdateIndexingNodeStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateIndexingNodeStakeResponse{}, err
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateIndexingNodeStakeResponse{}, err
	}

	ozoneLimitChange, completionTime, err := k.UpdateIndexingNodeStake(ctx, networkAddr, ownerAddress, *msg.StakeDelta, msg.IncrStake)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateIndexingNodeStake,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyIncrStakeBool, strconv.FormatBool(msg.IncrStake)),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingMatureTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgUpdateIndexingNodeStakeResponse{}, nil
}
