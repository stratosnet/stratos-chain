package pot

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stratosnet/stratos-chain/x/pot/keeper"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgVolumeReport:
			ctx.Logger().With("pot", "enter NewHandler")
			return handleMsgReportVolume(ctx, k, msg)
		case types.MsgWithdraw:
			return handleMsgWithdraw(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handle handleMsgReportVolume.
func handleMsgReportVolume(ctx sdk.Context, k keeper.Keeper, msg types.MsgVolumeReport) (*sdk.Result, error) {
	//ctx.Logger().Info("enter handleMsgReportVolume start", "true")
	//ctx.Logger().Info("ctx in pot:" + string(types.ModuleCdc.MustMarshalJSON(ctx)))
	//ctx.Logger().Info("Reporter in pot:" + string(types.ModuleCdc.MustMarshalJSON(msg.Reporter)))
	if !(k.IsSPNode(ctx, msg.Reporter)) {

		ctx.Logger().Info("Sender Info:", "IsSPNode", "false")
		errMsg := fmt.Sprint("message is not sent by a superior peer")
		ctx.Logger().Info(errMsg)

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, errMsg)
	}
	ctx.Logger().Info("Sender Info: ", "IsSPNode", "true")
	k.SetVolumeReport(ctx, msg.Reporter, msg.ReportReference)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVolumeReport,
			sdk.NewAttribute(types.AttributeKeyReportReference, hex.EncodeToString([]byte(msg.ReportReference))),
			sdk.NewAttribute(types.AttributeKeyEpoch, msg.Epoch.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdraw(ctx sdk.Context, k keeper.Keeper, msg types.MsgWithdraw) (*sdk.Result, error) {
	err := k.Withdraw(ctx, msg.Amount, msg.NodeAddress, msg.OwnerAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdraw,
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNodeAddress, msg.NodeAddress.String()),
			sdk.NewAttribute(types.AttributeKeyOwnerAddress, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
