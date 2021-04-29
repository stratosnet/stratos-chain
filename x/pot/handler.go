package pot

import (
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
		// this line is used by starport scaffolding # 1
		case types.MsgVolumeReport:
			return handleMsgReportVolume(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handle handleMsgReportVolume.
func handleMsgReportVolume(ctx sdk.Context, k keeper.Keeper, msg types.MsgVolumeReport) (*sdk.Result, error) {
	ctx.Logger().Info("enter handleMsgReportVolume " + msg.ReportReferenceHash)
	k.SetVolumeReportHash(ctx, &msg)
	ctx.Logger().Info("exit handleMsgReportVolume " + msg.ReportReferenceHash)

	for _, singleNodeVolume := range msg.NodesVolume {
		ctx.Logger().Info("enter singleNodeVolume " + singleNodeVolume.NodeAddress.String())
		k.SetSingleNodeVolume(ctx, singleNodeVolume)
		ctx.Logger().Info("exit singleNodeVolume " + singleNodeVolume.Volume.String())
	}

	//ctx.EventManager().EmitEvents(
	//		sdk.Events{
	//			sdk.NewEvent(
	//,			types.EventTypeVolumeReport,
	//			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	//			sdk.NewAttribute(types.AttributeKeyReporter, msg.Reporter.String()),
	//			sdk.NewAttribute(types.AttributeKeyReportReferenceHash, msg.ReportReferenceHash),
	//			sdk.NewAttribute(types.AttributeKeyEpoch, msg.Epoch.String()),
	//			sdk.NewAttribute(types.AttributeKeyNodesVolume, string(types.ModuleCdc.MustMarshalJSON(msg.NodesVolume))),
	//		),
	//		sdk.NewEvent(
	//			sdk.EventTypeMessage,
	//			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	//			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
	//		),
	//	})
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVolumeReport,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyReporter, msg.Reporter.String()),
			sdk.NewAttribute(types.AttributeKeyReportReferenceHash, msg.ReportReferenceHash),
			sdk.NewAttribute(types.AttributeKeyEpoch, msg.Epoch.String()),
			sdk.NewAttribute(types.AttributeKeyNodesVolume, string(types.ModuleCdc.MustMarshalJSON(msg.NodesVolume))),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Reporter.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}
