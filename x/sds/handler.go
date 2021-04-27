package sds

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/sds/keeper"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgFileUpload:
			return handleMsgFileUpload(ctx, k, msg)
		case types.MsgPrepay:
			return handleMsgPrepay(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handle MsgFileUpload.
func handleMsgFileUpload(ctx sdk.Context, k keeper.Keeper, msg types.MsgFileUpload) (*sdk.Result, error) {
	k.SetFileHash(ctx, msg.Sender, msg.FileHash)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFileUpload,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeKeyFileHash, hex.EncodeToString(msg.FileHash)),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// Handle MsgPrepay.
func handleMsgPrepay(ctx sdk.Context, k keeper.Keeper, msg types.MsgPrepay) (*sdk.Result, error) {
	if k.BankKeeper.GetSendEnabled(ctx) == false {
		return nil, nil
	}
	err := k.Prepay(ctx, msg.Sender, msg.Coins)
	if err == nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePrepay,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
				sdk.NewAttribute(types.AttributeKeyCoins, msg.Coins.String()),
			),
		)
		return &sdk.Result{Events: ctx.EventManager().Events()}, nil
	}
	return nil, nil
}
