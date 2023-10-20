package keeper

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (k msgServer) HandleMsgCreateResourceNode(goCtx context.Context, msg *types.MsgCreateResourceNode) (
	*types.MsgCreateResourceNodeResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.ResourceNodeRegEnabled(ctx) {
		return &types.MsgCreateResourceNodeResponse{}, types.ErrResourceNodeRegDisabled
	}

	// check to see if the pubKey or sender has been registered before
	pkAny := msg.GetPubkey()
	cachedPubKey := pkAny.GetCachedValue()
	pk := cachedPubKey.(cryptotypes.PubKey)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgCreateResourceNodeResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgCreateResourceNodeResponse{}, errors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	ozoneLimitChange, err := k.RegisterResourceNode(ctx, networkAddr, pk, ownerAddress, msg.Description, types.NodeType(msg.NodeType), msg.GetValue())
	if err != nil {
		return nil, errors.Wrap(types.ErrRegisterResourceNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateResourceNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyPubKey, hex.EncodeToString(pk.Bytes())),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.String()),
			sdk.NewAttribute(types.AttributeKeyInitialDeposit, msg.Value.String()),
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

	// check to see if the pubKey or sender has been registered before
	pkAny := msg.GetPubkey()
	cachedPubKey := pkAny.GetCachedValue()
	pk := cachedPubKey.(cryptotypes.PubKey)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgCreateMetaNodeResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgCreateMetaNodeResponse{}, errors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	ozoneLimitChange, err := k.RegisterMetaNode(ctx, networkAddr, pk, ownerAddress, msg.Description, msg.GetValue())
	if err != nil {
		return nil, errors.Wrap(types.ErrRegisterMetaNode, err.Error())
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

func (k msgServer) HandleMsgRemoveResourceNode(goCtx context.Context, msg *types.MsgRemoveResourceNode) (
	*types.MsgRemoveResourceNodeResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	p2pAddress, err := stratos.SdsAddressFromBech32(msg.GetResourceNodeAddress())
	if err != nil {
		return &types.MsgRemoveResourceNodeResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		return &types.MsgRemoveResourceNodeResponse{}, errors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	depositToRemove, completionTime, err := k.UnbondResourceNode(ctx, p2pAddress, ownerAddress)
	if err != nil {
		return &types.MsgRemoveResourceNodeResponse{}, errors.Wrap(types.ErrUnbondResourceNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondingResourceNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyResourceNode, msg.ResourceNodeAddress),
			sdk.NewAttribute(types.AttributeKeyDepositToRemove, sdk.NewCoin(k.BondDenom(ctx), depositToRemove).String()),
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
		return &types.MsgRemoveMetaNodeResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}
	metaNode, found := k.GetMetaNode(ctx, p2pAddress)
	if !found {
		return &types.MsgRemoveMetaNodeResponse{}, types.ErrNoMetaNodeFound
	}
	if msg.GetOwnerAddress() != metaNode.GetOwnerAddress() {
		return &types.MsgRemoveMetaNodeResponse{}, types.ErrInvalidOwnerAddr
	}

	unBondingDeposit := k.GetUnbondingNodeBalance(ctx, p2pAddress)
	availableDeposit := metaNode.Tokens.Sub(unBondingDeposit)
	if availableDeposit.LTE(sdkmath.ZeroInt()) {
		return &types.MsgRemoveMetaNodeResponse{}, types.ErrInsufficientBalance
	}

	ozoneLimitChange, _, _, completionTime, err := k.UnbondMetaNode(ctx, metaNode, availableDeposit)
	if err != nil {
		return &types.MsgRemoveMetaNodeResponse{}, errors.Wrap(types.ErrUnbondMetaNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondingMetaNode,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyMetaNode, msg.MetaNodeAddress),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.Neg().String()),
			sdk.NewAttribute(types.AttributeKeyDepositToRemove, sdk.NewCoin(k.BondDenom(ctx), availableDeposit).String()),
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

	candidateNetworkAddress, err := stratos.SdsAddressFromBech32(msg.GetCandidateNetworkAddress())
	if err != nil {
		return &types.MsgMetaNodeRegistrationVoteResponse{}, errors.Wrap(types.ErrInvalidCandidateNetworkAddr, err.Error())
	}
	candidateOwnerAddress, err := sdk.AccAddressFromBech32(msg.GetCandidateOwnerAddress())
	if err != nil {
		return &types.MsgMetaNodeRegistrationVoteResponse{}, types.ErrInvalidCandidateOwnerAddr
	}
	voterNetworkAddress, err := stratos.SdsAddressFromBech32(msg.GetVoterNetworkAddress())
	if err != nil {
		return &types.MsgMetaNodeRegistrationVoteResponse{}, errors.Wrap(types.ErrInvalidVoterNetworkAddr, err.Error())
	}
	voterOwnerAddress, err := sdk.AccAddressFromBech32(msg.GetVoterOwnerAddress())
	if err != nil {
		return &types.MsgMetaNodeRegistrationVoteResponse{}, errors.Wrap(types.ErrInvalidVoterOwnerAddr, err.Error())
	}

	nodeStatus, err := k.HandleVoteForMetaNodeRegistration(
		ctx, candidateNetworkAddress, candidateOwnerAddress, types.VoteOpinion(msg.Opinion), voterNetworkAddress, voterOwnerAddress)
	if err != nil {
		return &types.MsgMetaNodeRegistrationVoteResponse{}, errors.Wrap(types.ErrVoteMetaNode, err.Error())
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

func (k msgServer) HandleMsgWithdrawMetaNodeRegistrationDeposit(goCtx context.Context, msg *types.MsgWithdrawMetaNodeRegistrationDeposit) (
	*types.MsgWithdrawMetaNodeRegistrationDepositResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.GetNetworkAddress())
	if err != nil {
		return &types.MsgWithdrawMetaNodeRegistrationDepositResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}
	ownerAddr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		return &types.MsgWithdrawMetaNodeRegistrationDepositResponse{}, types.ErrInvalidOwnerAddr
	}

	completionTime, err := k.WithdrawMetaNodeRegistrationDeposit(ctx, networkAddr, ownerAddr)
	if err != nil {
		return &types.MsgWithdrawMetaNodeRegistrationDepositResponse{}, errors.Wrap(types.ErrUnbondMetaNode, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawMetaNodeRegistrationDeposit,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyUnbondingMatureTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})

	return &types.MsgWithdrawMetaNodeRegistrationDepositResponse{}, nil
}

func (k msgServer) HandleMsgUpdateResourceNode(goCtx context.Context, msg *types.MsgUpdateResourceNode) (*types.MsgUpdateResourceNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}
	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeResponse{}, errors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	err = k.UpdateResourceNode(ctx, msg.Description, types.NodeType(msg.NodeType), networkAddr, ownerAddress)
	if err != nil {
		return nil, errors.Wrap(types.ErrUpdateResourceNode, err.Error())
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

func (k msgServer) HandleMsgUpdateResourceNodeDeposit(goCtx context.Context, msg *types.MsgUpdateResourceNodeDeposit) (
	*types.MsgUpdateResourceNodeDepositResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeDepositResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateResourceNodeDepositResponse{}, errors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter, completionTime, node, err :=
		k.UpdateResourceNodeDeposit(ctx, networkAddr, ownerAddress, msg.GetDepositDelta())
	if err != nil {
		return nil, errors.Wrap(types.ErrUpdateResourceNodeDeposit, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateResourceNodeDeposit,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyDepositDelta, msg.DepositDelta.String()),
			sdk.NewAttribute(types.AttributeKeyCurrentDeposit, sdk.NewCoin(k.BondDenom(ctx), node.Tokens).String()),
			sdk.NewAttribute(types.AttributeKeyAvailableTokenBefore, sdk.NewCoin(k.BondDenom(ctx), availableTokenAmtBefore).String()),
			sdk.NewAttribute(types.AttributeKeyAvailableTokenAfter, sdk.NewCoin(k.BondDenom(ctx), availableTokenAmtAfter).String()),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingMatureTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgUpdateResourceNodeDepositResponse{}, nil
}

func (k msgServer) HandleMsgUpdateEffectiveDeposit(goCtx context.Context, msg *types.MsgUpdateEffectiveDeposit) (*types.MsgUpdateEffectiveDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.Reporters) == 0 || len(msg.ReporterOwner) == 0 {
		return &types.MsgUpdateEffectiveDepositResponse{}, types.ErrReporterAddressOrOwner
	}

	reporterOwners := msg.ReporterOwner
	validReporterCount := 0
	for idx, reporter := range msg.Reporters {
		reporterSdsAddr, err := stratos.SdsAddressFromBech32(reporter)
		if err != nil {
			continue
		}
		ownerAddr, err := sdk.AccAddressFromBech32(reporterOwners[idx])
		if err != nil {
			continue
		}
		if !(k.OwnMetaNode(ctx, ownerAddr, reporterSdsAddr)) {
			continue
		}
		validReporterCount++
	}

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateEffectiveDepositResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	if msg.EffectiveTokens.LTE(sdkmath.NewInt(0)) {
		return &types.MsgUpdateEffectiveDepositResponse{}, types.ErrInvalidEffectiveToken
	}

	_, _, isUnsuspendedDuringUpdate, err := k.UpdateEffectiveDeposit(ctx, networkAddr, msg.EffectiveTokens)
	if err != nil {
		return nil, errors.Wrap(types.ErrUpdateResourceNodeDeposit, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateEffectiveDeposit,
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyEffectiveDepositAfter, msg.EffectiveTokens.String()),
			sdk.NewAttribute(types.AttributeKeyIsUnsuspended, strconv.FormatBool(isUnsuspendedDuringUpdate)),
		),
	})
	return &types.MsgUpdateEffectiveDepositResponse{}, nil
}

func (k msgServer) HandleMsgUpdateMetaNode(goCtx context.Context, msg *types.MsgUpdateMetaNode) (*types.MsgUpdateMetaNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateMetaNodeResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateMetaNodeResponse{}, errors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	err = k.UpdateMetaNode(ctx, msg.Description, networkAddr, ownerAddress)
	if err != nil {
		return nil, errors.Wrap(types.ErrUpdateMetaNode, err.Error())
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

func (k msgServer) HandleMsgUpdateMetaNodeDeposit(goCtx context.Context, msg *types.MsgUpdateMetaNodeDeposit) (
	*types.MsgUpdateMetaNodeDepositResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	networkAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgUpdateMetaNodeDepositResponse{}, errors.Wrap(types.ErrInvalidNetworkAddr, err.Error())
	}

	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return &types.MsgUpdateMetaNodeDepositResponse{}, errors.Wrap(types.ErrInvalidOwnerAddr, err.Error())
	}

	if msg.DepositDelta.Amount.IsNegative() {
		return &types.MsgUpdateMetaNodeDepositResponse{}, errors.Wrap(types.ErrInvalidDepositChange, err.Error())
	}

	ozoneLimitChange, availableTokenAmtBefore, availableTokenAmtAfter, completionTime, node, err := k.UpdateMetaNodeDeposit(ctx, networkAddr, ownerAddress, msg.GetDepositDelta())
	if err != nil {
		return nil, errors.Wrap(types.ErrUpdateMetaNodeDeposit, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateMetaNodeDeposit,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyNetworkAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyDepositDelta, msg.DepositDelta.String()),
			sdk.NewAttribute(types.AttributeKeyCurrentDeposit, sdk.NewCoin(k.BondDenom(ctx), node.Tokens).String()),
			sdk.NewAttribute(types.AttributeKeyAvailableTokenBefore, sdk.NewCoin(k.BondDenom(ctx), availableTokenAmtBefore).String()),
			sdk.NewAttribute(types.AttributeKeyAvailableTokenAfter, sdk.NewCoin(k.BondDenom(ctx), availableTokenAmtAfter).String()),
			sdk.NewAttribute(types.AttributeKeyOZoneLimitChanges, ozoneLimitChange.String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingMatureTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})
	return &types.MsgUpdateMetaNodeDepositResponse{}, nil
}
