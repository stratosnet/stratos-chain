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

		// this line is used by starport scaffolding # 1
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateResourceNode(ctx sdk.Context, msg types.MsgCreateResourceNode, k keeper.Keeper) (*sdk.Result, error) {
	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetResourceNode(ctx, msg.ResourceNodeAddress); found {
		return nil, types.ErrResourceNodeOwnerExists
	}
	if _, found := k.GetResourceNodeByAddr(ctx, sdk.GetConsAddress(msg.PubKey)); found {
		return nil, types.ErrResourceNodePubKeyExists
	}
	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	resourceNode := types.NewResourceNode(msg.ResourceNodeAddress, msg.PubKey, msg.Description)
	k.SetResourceNode(ctx, resourceNode)
	k.SetResourceNodeByAddr(ctx, resourceNode)
	k.SetNewResourceNodeByPowerIndex(ctx, resourceNode)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateResourceNode,
			sdk.NewAttribute(types.AttributeKeyResourceNode, msg.ResourceNodeAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateIndexingNode(ctx sdk.Context, msg types.MsgCreateIndexingNode, k keeper.Keeper) (*sdk.Result, error) {
	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetIndexingNode(ctx, msg.IndexingNodeAddress); found {
		return nil, types.ErrIndexingNodeOwnerExists
	}
	if _, found := k.GetIndexingNodeByAddr(ctx, sdk.GetConsAddress(msg.PubKey)); found {
		return nil, types.ErrIndexingNodePubKeyExists
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	indexingNode := types.NewIndexingNode(msg.IndexingNodeAddress, msg.PubKey, msg.Description)
	k.SetIndexingNode(ctx, indexingNode)
	k.SetIndexingNodeByAddr(ctx, indexingNode)
	k.SetNewIndexingNodeByPowerIndex(ctx, indexingNode)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateIndexingNode,
			sdk.NewAttribute(types.AttributeKeyIndexingNode, msg.IndexingNodeAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
