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
	if !(k.IsIndexingNode(ctx, msg.Reporter)) {

		ctx.Logger().Info("IsIndexingNode", "false")
		errMsg := fmt.Sprint("message is not sent by a superior peer")
		ctx.Logger().Info(errMsg)

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, errMsg)
	}
	ctx.Logger().Info("IsIndexingNode", "true")
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
