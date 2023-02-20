package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	if !k.ResourceNodeRegEnabled(ctx) {
		return &types.MsgCreateResourceNodeResponse{}, types.ErrResourceNodeRegDisabled
	}

	// check to see if the pubkey or sender has been registered before
	pkAny := msg.GetPubkey()
	cachedPubkey := pkAny.GetCachedValue()
	pk := cachedPubkey.(cryptotypes.PubKey)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgCreateResourceNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgCreateResourceNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	if _, found := k.GetResourceNode(ctx, networkAddr); found {
		ctx.Logger().Error("Resource node already exist")
		return nil, types.ErrResourceNodePubKeyExists
	}
	if msg.Value.Denom != k.BondDenom(ctx) {
		return nil, types.ErrBadDenom
	}

	ozoneLimitChange, err := k.RegisterResourceNode(ctx, networkAddr, pk, ownerAddress, *msg.Description, types.NodeType(msg.NodeType), msg.Value)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrRegisterResourceNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateResourceNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyPubKey, hex.EncodeToString(pk.Bytes())),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.String()),
			sdk.NewAttribute(types.AttributeKeyInitialStake, msg.Value.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgCreateResourceNodeResponse{}, nil
}

func (k msgServer) HandleMsgCreateMetaNode(goCtx context.Context, msg *types.MsgCreateMetaNode) (*types.MsgCreateMetaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check to see if the pubkey or sender has been registered before
	pkAny := msg.GetPubkey()
	cachedPubkey := pkAny.GetCachedValue()
	pk := cachedPubkey.(cryptotypes.PubKey)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgCreateMetaNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	if _, found := k.GetMetaNode(ctx, networkAddr); found {
		ctx.Logger().Error("Meta node already exist")
		return nil, types.ErrMetaNodePubKeyExists
	}
	if msg.Value.Denom != k.BondDenom(ctx) {
		return nil, types.ErrBadDenom
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgCreateMetaNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	ozoneLimitChange, err := k.RegisterMetaNode(ctx, networkAddr, pk, ownerAddress, *msg.Description, msg.Value)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrRegisterMetaNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateMetaNode,
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
	return &types.MsgCreateMetaNodeResponse{}, nil
}

func (k msgServer) HandleMsgRemoveResourceNode(goCtx context.Context, msg *types.MsgRemoveResourceNode) (*types.MsgRemoveResourceNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	p2pAddress, err := stratos.SdsAddressFromBech32(msg.ResourceNodeAddress)
	if err != nil {
		return &types.MsgRemoveResourceNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}
	resourceNode, found := k.GetResourceNode(ctx, p2pAddress)
	if !found {
		return nil, types.ErrNoResourceNodeFound
	}
	if resourceNode.GetStatus() == stakingtypes.Unbonding {
		return nil, types.ErrUnbondingNode
	}

	unbondingStake := k.GetUnbondingNodeBalance(ctx, p2pAddress)
	availableStake := resourceNode.Tokens.Sub(unbondingStake)
	if availableStake.LTE(sdk.ZeroInt()) {
		return nil, types.ErrInsufficientBalance
	}

	ozoneLimitChange, completionTime, err := k.UnbondResourceNode(ctx, resourceNode, availableStake)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnbondResourceNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondingResourceNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyResourceNode, msg.ResourceNodeAddress),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.Neg().String()),
			sdk.NewAttribute(types.AttributeKeyStakeToRemove, sdk.NewCoin(k.BondDenom(ctx), availableStake).String()),
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

func (k msgServer) HandleMsgRemoveMetaNode(goCtx context.Context, msg *types.MsgRemoveMetaNode) (*types.MsgRemoveMetaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	p2pAddress, err := stratos.SdsAddressFromBech32(msg.MetaNodeAddress)
	if err != nil {
		return &types.MsgRemoveMetaNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}
	metaNode, found := k.GetMetaNode(ctx, p2pAddress)
	if !found {
		return nil, types.ErrNoMetaNodeFound
	}

	if metaNode.GetStatus() == stakingtypes.Unbonding {
		return nil, types.ErrUnbondingNode
	}

	unbondingStake := k.GetUnbondingNodeBalance(ctx, p2pAddress)
	availableStake := metaNode.Tokens.Sub(unbondingStake)
	if availableStake.LTE(sdk.ZeroInt()) {
		return nil, types.ErrInsufficientBalance
	}

	ozoneLimitChange, completionTime, err := k.UnbondMetaNode(ctx, metaNode, availableStake)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnbondMetaNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondingMetaNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyMetaNode, msg.MetaNodeAddress),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.Neg().String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingMatureTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})

	return &types.MsgRemoveMetaNodeResponse{}, nil
}

func (k msgServer) HandleMsgMetaNodeRegistrationVote(goCtx context.Context, msg *types.MsgMetaNodeRegistrationVote) (*types.MsgMetaNodeRegistrationVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	candidateNetworkAddress, err := stratos.SdsAddressFromBech32(msg.CandidateNetworkAddress)
	if err != nil {
		return &types.MsgMetaNodeRegistrationVoteResponse{}, sdkerrors.Wrap(types.ErrInvalidCandidateNetworkAddr, err.Error())
	}

	nodeToApprove, found := k.GetMetaNode(ctx, candidateNetworkAddress)
	if !found {
		return nil, types.ErrNoMetaNodeFound
	}
	if nodeToApprove.OwnerAddress != msg.CandidateOwnerAddress {
		return nil, types.ErrInvalidOwnerAddr
	}

	voterNetworkAddress, err := stratos.SdsAddressFromBech32(msg.VoterNetworkAddress)
	if err != nil {
		return &types.MsgMetaNodeRegistrationVoteResponse{}, sdkerrors.Wrap(types.ErrInvalidVoterNetworkAddr, err.Error())
	}
	voter, found := k.GetMetaNode(ctx, voterNetworkAddress)
	if !found {
		return nil, types.ErrInvalidVoterAddr
	}

	candidateOwnerAddress, err := sdk.AccAddressFromBech32(msg.CandidateOwnerAddress)
	if err != nil {
		return &types.MsgMetaNodeRegistrationVoteResponse{}, sdkerrors.Wrap(types.ErrInvalidCandidateOwnerAddr, err.Error())
	}

	if !(voter.Status == stakingtypes.Bonded) || voter.Suspend {
		return nil, types.ErrInvalidVoterStatus
	}

	nodeStatus, err := k.HandleVoteForMetaNodeRegistration(ctx, candidateNetworkAddress, candidateOwnerAddress, types.VoteOpinion(msg.Opinion), voterNetworkAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrVoteMetaNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMetaNodeRegistrationVote,
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

	return &types.MsgMetaNodeRegistrationVoteResponse{}, nil
}

func (k msgServer) HandleMsgUpdateResourceNode(goCtx context.Context, msg *types.MsgUpdateResourceNode) (*types.MsgUpdateResourceNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}
	err = k.UpdateResourceNode(ctx, msg.Description, types.NodeType(msg.NodeType), networkAddr, ownerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUpdateResourceNode, err.Error())
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
		return &types.MsgUpdateResourceNodeStakeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeStakeResponse{}, sdkerrors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	if msg.StakeDelta.Amount.LT(sdk.NewInt(0)) {
		return &types.MsgUpdateResourceNodeStakeResponse{}, errors.New("invalid stake delta")
	}

	ozoneLimitChange, completionTime, _, err := k.UpdateResourceNodeStake(ctx, networkAddr, ownerAddress, *msg.StakeDelta, msg.IncrStake)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUpdateResourceNodeStake, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateResourceNodeStake,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyIncrStakeBool, strconv.FormatBool(msg.IncrStake)),
			sdk.NewAttribute(types.AttributeKeyStakeDelta, msg.StakeDelta.String()),
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

func (k msgServer) HandleMsgUpdateEffectiveStake(goCtx context.Context, msg *types.MsgUpdateEffectiveStake) (*types.MsgUpdateEffectiveStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, reporter := range msg.Reporters {
		reporterSdsAddr, err := stratos.SdsAddressFromBech32(reporter)
		if err != nil {
			return &types.MsgUpdateEffectiveStakeResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
		}
		if !(k.IsSPNode(ctx, reporterSdsAddr)) {
			return &types.MsgUpdateEffectiveStakeResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "MsgUpdateEffectiveStake is not sent by a meta node")
		}
	}

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateEffectiveStakeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	if msg.EffectiveTokens.LTE(sdk.NewInt(0)) {
		return &types.MsgUpdateEffectiveStakeResponse{}, errors.New("effective tokens should be greater than 0")
	}

	_, _, isUnsuspendedDuringUpdate, err := k.UpdateEffectiveStake(ctx, networkAddr, msg.EffectiveTokens)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUpdateResourceNodeStake, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateEffectiveStake,
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyEffectiveStakeAfter, msg.EffectiveTokens.String()),
			sdk.NewAttribute(types.AttributeKeyIsUnsuspended, strconv.FormatBool(isUnsuspendedDuringUpdate)),
		),
	})
	return &types.MsgUpdateEffectiveStakeResponse{}, nil
}

func (k msgServer) HandleMsgUpdateMetaNode(goCtx context.Context, msg *types.MsgUpdateMetaNode) (*types.MsgUpdateMetaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateMetaNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateMetaNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	err = k.UpdateMetaNode(ctx, msg.Description, networkAddr, ownerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUpdateMetaNode, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateMetaNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgUpdateMetaNodeResponse{}, nil
}

func (k msgServer) HandleMsgUpdateMetaNodeStake(goCtx context.Context, msg *types.MsgUpdateMetaNodeStake) (*types.MsgUpdateMetaNodeStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateMetaNodeStakeResponse{}, sdkerrors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateMetaNodeStakeResponse{}, sdkerrors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	if msg.StakeDelta.Amount.LT(sdk.NewInt(0)) {
		return &types.MsgUpdateMetaNodeStakeResponse{}, errors.New("invalid stake delta")
	}

	ozoneLimitChange, completionTime, err := k.UpdateMetaNodeStake(ctx, networkAddr, ownerAddress, *msg.StakeDelta, msg.IncrStake)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUpdateMetaNodeStake, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateMetaNodeStake,
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
	return &types.MsgUpdateMetaNodeStakeResponse{}, nil
}
