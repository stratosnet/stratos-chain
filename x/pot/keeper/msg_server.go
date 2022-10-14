package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
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

func (k msgServer) HandleMsgVolumeReport(goCtx context.Context, msg *types.MsgVolumeReport) (*types.MsgVolumeReportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	reporter, err := stratos.SdsAddressFromBech32(msg.Reporter)
	if err != nil {
		return &types.MsgVolumeReportResponse{}, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	if !(k.IsSPNode(ctx, reporter)) {
		return &types.MsgVolumeReportResponse{}, sdkerrors.Wrap(types.ErrInvalidAddress, "Volume report is not sent by a superior peer")
	}

	// ensure epoch increment
	epoch, ok := sdk.NewIntFromString(msg.Epoch.String())
	if !ok {
		return &types.MsgVolumeReportResponse{}, types.ErrInvalid
	}
	lastEpoch := k.GetLastReportedEpoch(ctx)
	if msg.Epoch.LTE(lastEpoch) {
		e := sdkerrors.Wrapf(types.ErrMatureEpoch, "expected epoch should be greater than %s, got %s",
			lastEpoch.String(), msg.Epoch.String())
		return &types.MsgVolumeReportResponse{}, e
	}

	// TODO: verify BLS signature

	txBytes := ctx.TxBytes()
	txhash := fmt.Sprintf("%X", tmhash.Sum(txBytes))

	err = k.VolumeReport(ctx, msg.WalletVolumes, reporter, epoch, msg.ReportReference, txhash)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrVolumeReport, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVolumeReport,
			sdk.NewAttribute(types.AttributeKeyReportReference, hex.EncodeToString([]byte(msg.ReportReference))),
			sdk.NewAttribute(types.AttributeKeyEpoch, msg.Epoch.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ReporterOwner),
		),
	})

	return &types.MsgVolumeReportResponse{}, nil
}

func (k msgServer) HandleMsgWithdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	walletAddress, err := sdk.AccAddressFromBech32(msg.WalletAddress)
	if err != nil {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	targetAddress, err := sdk.AccAddressFromBech32(msg.TargetAddress)
	if err != nil {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	err = k.Withdraw(ctx, msg.Amount, walletAddress, targetAddress)
	if err != nil {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrap(types.ErrWithdrawFailure, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdraw,
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.WalletAddress),
		),
	})
	return &types.MsgWithdrawResponse{}, nil
}

func (k msgServer) HandleMsgFoundationDeposit(goCtx context.Context, msg *types.MsgFoundationDeposit) (*types.MsgFoundationDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return &types.MsgFoundationDepositResponse{}, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	err = k.FoundationDeposit(ctx, msg.Amount, from)
	if err != nil {
		return &types.MsgFoundationDepositResponse{}, sdkerrors.Wrap(types.ErrFoundationDepositFailure, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFoundationDeposit,
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})
	return &types.MsgFoundationDepositResponse{}, nil
}

func (k msgServer) HandleMsgSlashingResourceNode(goCtx context.Context, msg *types.MsgSlashingResourceNode) (*types.MsgSlashingResourceNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, reporter := range msg.Reporters {
		reporterSdsAddr, err := stratos.SdsAddressFromBech32(reporter)
		if err != nil {
			return &types.MsgSlashingResourceNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
		}
		if !(k.IsSPNode(ctx, reporterSdsAddr)) {
			return &types.MsgSlashingResourceNodeResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Slashing msg is not sent by a meta node")
		}
	}
	networkAddress, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgSlashingResourceNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	walletAddress, err := sdk.AccAddressFromBech32(msg.WalletAddress)
	if err != nil {
		return &types.MsgSlashingResourceNodeResponse{}, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	slashing, ok := sdk.NewIntFromString(msg.Slashing.String())
	if !ok {
		return &types.MsgSlashingResourceNodeResponse{}, types.ErrInvalidAmount
	}
	amt, nodeType, err := k.SlashingResourceNode(ctx, networkAddress, walletAddress, slashing, msg.Suspend)
	if err != nil {
		return &types.MsgSlashingResourceNodeResponse{}, sdkerrors.Wrap(types.ErrSlashingResourceNodeFailure, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSlashing,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyNodeP2PAddress, msg.NetworkAddress),
			sdk.NewAttribute(types.AttributeKeyAmount, amt.String()),
			sdk.NewAttribute(types.AttributeKeySlashingNodeType, nodeType.String()),
			sdk.NewAttribute(types.AttributeKeyNodeSuspended, strconv.FormatBool(msg.Suspend)),
		),
	})
	return &types.MsgSlashingResourceNodeResponse{}, nil
}
