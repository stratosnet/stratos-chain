package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

	err = ctx.EventManager().EmitTypedEvents(
		&types.EventFileUpload{
			Sender:   msg.GetFrom(),
			Reporter: msg.GetReporter(),
			Uploader: msg.GetUploader(),
			FileHash: msg.GetFileHash(),
		},
		&types.EventMessage{
			Module: types.ModuleName,
			Sender: msg.GetFrom(),
		},
	)
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

	err = ctx.EventManager().EmitTypedEvents(
		&types.EventPrePay{
			Sender:       msg.GetSender(),
			Beneficiary:  msg.GetBeneficiary(),
			Amount:       msg.GetAmount().String(),
			PurchasedNoz: purchased.String(),
		},
		&types.EventMessage{
			Module: types.ModuleName,
			Sender: msg.GetSender(),
		},
	)
	if err != nil {
		return nil, errors.Wrap(types.ErrEmitEvent, err.Error())
	}

	return &types.MsgPrepayResponse{}, nil
}

// UpdateParams updates the module parameters
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventMessage{
		Module: types.ModuleName,
		Sender: msg.Authority,
		Action: sdk.MsgTypeURL(msg),
	})
	if err != nil {
		return nil, errors.Wrap(types.ErrEmitEvent, err.Error())
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
