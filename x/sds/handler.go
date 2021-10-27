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
	// check if reporter addr belongs to an registered sp node
	if _, found := k.RegisterKeeper.GetIndexingNode(ctx, msg.Reporter); found == false {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "Reporter %s isn't an SP node", msg.Reporter.String())
	}
	height := sdk.NewInt(ctx.BlockHeight())
	heightByteArr, _ := height.MarshalJSON()
	var heightReEncoded sdk.Int
	heightReEncoded.UnmarshalJSON(heightByteArr)

	fileInfo := types.NewFileInfo(heightReEncoded, msg.Reporter, msg.Uploader)
	k.SetFileHash(ctx, msg.FileHash, fileInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFileUpload,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyReporter, msg.Reporter.String()),
			sdk.NewAttribute(types.AttributeKeyUploader, msg.Uploader.String()),
			sdk.NewAttribute(types.AttributeKeyFileHash, hex.EncodeToString(msg.FileHash)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// Handle MsgPrepay.
func handleMsgPrepay(ctx sdk.Context, k keeper.Keeper, msg types.MsgPrepay) (*sdk.Result, error) {
	if k.BankKeeper.GetSendEnabled(ctx) == false {
		return nil, nil
	}
	purchased, err := k.Prepay(ctx, msg.Sender, msg.Coins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePrepay,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeKeyCoins, msg.Coins.String()),
			sdk.NewAttribute(types.AttributeKeyPurchasedUoz, purchased.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
