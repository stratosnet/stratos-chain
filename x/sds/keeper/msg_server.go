package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/sds/types"
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

// HandleMsgFileUpload Handles MsgFileUpload.
func (k msgServer) HandleMsgFileUpload(c context.Context, msg *types.MsgFileUpload) (*types.MsgFileUploadResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	reporter, err := stratos.SdsAddressFromBech32(msg.GetReporter())
	if err != nil {
		return &types.MsgFileUploadResponse{}, errors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	reporterOwner, err := sdk.AccAddressFromBech32(msg.GetFrom())
	if err != nil {
		return &types.MsgFileUploadResponse{}, errors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	uploader, err := sdk.AccAddressFromBech32(msg.Uploader)
	if err != nil {
		return &types.MsgFileUploadResponse{}, errors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	err = k.FileUpload(ctx, msg.GetFileHash(), reporter, reporterOwner, uploader)
	if err != nil {
		return &types.MsgFileUploadResponse{}, err
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventFileUpload{
		Sender:   msg.GetFrom(),
		Reporter: msg.GetReporter(),
		Uploader: msg.GetUploader(),
		FileHash: msg.GetFileHash(),
	})
	if err != nil {
		return nil, errors.Wrap(types.ErrEmitEvent, err.Error())
	}

	return &types.MsgFileUploadResponse{}, nil
}

// HandleMsgPrepay Handles MsgPrepay.
func (k msgServer) HandleMsgPrepay(c context.Context, msg *types.MsgPrepay) (*types.MsgPrepayResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.GetSender())
	if err != nil {
		return &types.MsgPrepayResponse{}, errors.Wrap(types.ErrInvalidSenderAddr, err.Error())
	}

	_, err = sdk.AccAddressFromBech32(msg.GetBeneficiary())
	if err != nil {
		return &types.MsgPrepayResponse{}, errors.Wrap(types.ErrInvalidBeneficiaryAddr, err.Error())
	}

	purchased, err := k.Prepay(ctx, sender, msg.GetAmount())
	if err != nil {
		return nil, errors.Wrap(types.ErrPrepayFailure, err.Error())
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventPrePay{
		Sender:       msg.GetSender(),
		Beneficiary:  msg.GetBeneficiary(),
		Amount:       msg.GetAmount(),
		PurchasedNoz: purchased,
	})
	if err != nil {
		return nil, errors.Wrap(types.ErrEmitEvent, err.Error())
	}

	return &types.MsgPrepayResponse{}, nil
}
