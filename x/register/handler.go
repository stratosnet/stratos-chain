package register

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/register/keeper"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateResourceNode:
			return handleMsgCreateResourceNode(ctx, msg, k)
		case types.MsgCreateIndexingNode:
			return handleMsgCreateIndexingNode(ctx, msg, k)
		case types.MsgRemoveResourceNode:
			return handleMsgRemoveResourceNode(ctx, msg, k)
		case types.MsgRemoveIndexingNode:
			return handleMsgRemoveIndexingNode(ctx, msg, k)

		// this line is used by starport scaffolding # 1
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateResourceNode(ctx sdk.Context, msg types.MsgCreateResourceNode, k keeper.Keeper) (*sdk.Result, error) {
	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetResourceNode(ctx, sdk.AccAddress(msg.PubKey.Address())); found {
		ctx.Logger().Error("Resource node already exist")
		return nil, ErrResourceNodePubKeyExists
	}
	if msg.Value.Denom != k.BondDenom(ctx) {
		return nil, ErrBadDenom
	}
	//if _, err := msg.Description.EnsureLength(); err != nil {
	//	return nil, err
	//}

	resourceNode := types.NewResourceNode(msg.NetworkAddress, msg.PubKey, msg.OwnerAddress, msg.Description)
	err := k.AddResourceNodeTokens(ctx, resourceNode, msg.Value)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateResourceNode,
			sdk.NewAttribute(types.AttributeKeyResourceNode, msg.NetworkAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNodeAddress, sdk.AccAddress(msg.PubKey.Address()).String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateIndexingNode(ctx sdk.Context, msg types.MsgCreateIndexingNode, k keeper.Keeper) (*sdk.Result, error) {
	ctx.Logger().Info(fmt.Sprintf("in handleMsgCreateIndexingNode, indexingNodeAddress = %s", sdk.AccAddress(msg.PubKey.Address()).String()))
	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetIndexingNode(ctx, sdk.AccAddress(msg.PubKey.Address())); found {
		return nil, ErrIndexingNodePubKeyExists
	}
	if msg.Value.Denom != k.BondDenom(ctx) {
		return nil, ErrBadDenom
	}
	//if _, err := msg.Description.EnsureLength(); err != nil {
	//	return nil, err
	//}

	indexingNode := types.NewIndexingNode(msg.NetworkAddress, msg.PubKey, msg.OwnerAddress, msg.Description)
	err := k.AddIndexingNodeTokens(ctx, indexingNode, msg.Value)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateIndexingNode,
			sdk.NewAttribute(types.AttributeKeyIndexingNode, msg.NetworkAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNodeAddress, sdk.AccAddress(msg.PubKey.Address()).String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveResourceNode(ctx sdk.Context, msg types.MsgRemoveResourceNode, k keeper.Keeper) (*sdk.Result, error) {
	resourceNode, found := k.GetResourceNode(ctx, msg.ResourceNodeAddress)
	if !found {
		return nil, ErrNoResourceNodeFound
	}
	err := k.SubtractResourceNodeTokens(ctx, resourceNode, resourceNode.GetTokens())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveResourceNode,
			sdk.NewAttribute(types.AttributeKeyResourceNode, msg.ResourceNodeAddress.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveIndexingNode(ctx sdk.Context, msg types.MsgRemoveIndexingNode, k keeper.Keeper) (*sdk.Result, error) {
	ctx.Logger().Info("in handleMsgRemoveIndexingNode, indexingNodeAddress = " + msg.IndexingNodeAddress.String())
	ctx.Logger().Info("in handleMsgRemoveIndexingNode, ownerAddress = " + msg.OwnerAddress.String())
	indexingNode, found := k.GetIndexingNode(ctx, msg.IndexingNodeAddress)
	if !found {
		return nil, ErrNoIndexingNodeFound
	}
	err := k.SubtractIndexingNodeTokens(ctx, indexingNode, indexingNode.GetTokens())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveIndexingNode,
			sdk.NewAttribute(types.AttributeKeyIndexingNode, msg.IndexingNodeAddress.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
