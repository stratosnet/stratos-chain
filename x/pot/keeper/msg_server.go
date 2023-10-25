package keeper

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cometbft/cometbft/crypto/tmhash"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/crypto"
	"github.com/stratosnet/stratos-chain/crypto/bls"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
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
		return &types.MsgVolumeReportResponse{}, errors.Wrap(types.ErrReporterAddress, err.Error())
	}
	reporterOwner, err := sdk.AccAddressFromBech32(msg.ReporterOwner)
	if err != nil {
		return &types.MsgVolumeReportResponse{}, errors.Wrap(types.ErrReporterOwnerAddr, err.Error())
	}

	if !(k.registerKeeper.OwnMetaNode(ctx, reporterOwner, reporter)) {
		return &types.MsgVolumeReportResponse{}, types.ErrReporterAddressOrOwner
	}

	// ensure epoch increment
	epoch, ok := sdkmath.NewIntFromString(msg.Epoch.String())
	if !ok {
		return &types.MsgVolumeReportResponse{}, types.ErrInvalid
	}
	lastDistributedEpoch := k.GetLastDistributedEpoch(ctx)
	if msg.Epoch.LTE(lastDistributedEpoch) {
		e := errors.Wrapf(types.ErrMatureEpoch, "expected epoch should be greater than %s, got %s",
			lastDistributedEpoch.String(), msg.Epoch.String())
		return &types.MsgVolumeReportResponse{}, e
	}

	blsSignature := msg.GetBLSSignature()

	// verify txDataHash
	signBytes := msg.GetBLSSignBytes()
	txDataHash := crypto.Keccak256(signBytes)
	if !bytes.Equal(txDataHash, blsSignature.GetTxData()) {
		return &types.MsgVolumeReportResponse{}, types.ErrBLSTxDataInvalid
	}

	// verify blsSignature
	verified, err := bls.Verify(blsSignature.GetTxData(), blsSignature.GetSignature(), blsSignature.GetPubKeys()...)
	if err != nil {
		return &types.MsgVolumeReportResponse{}, errors.Wrap(types.ErrBLSVerifyFailed, err.Error())
	}
	if !verified {
		return &types.MsgVolumeReportResponse{}, types.ErrBLSVerifyFailed
	}

	if !k.HasReachedThreshold(ctx, blsSignature.GetPubKeys()) {
		return &types.MsgVolumeReportResponse{}, types.ErrBLSNotReachThreshold
	}

	txBytes := ctx.TxBytes()
	txhash := fmt.Sprintf("%X", tmhash.Sum(txBytes))

	walletVolumes := types.WalletVolumes{Volumes: msg.WalletVolumes}

	err = k.VolumeReport(ctx, walletVolumes, reporter, epoch, msg.ReportReference, txhash)
	if err != nil {
		return nil, errors.Wrap(types.ErrVolumeReport, err.Error())
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventVolumeReport{
		ReportReference: msg.GetReportReference(),
		Epoch:           msg.Epoch,
	})
	if err != nil {
		return nil, errors.Wrap(types.ErrEmitEvent, err.Error())
	}

	return &types.MsgVolumeReportResponse{}, nil
}

func (k msgServer) HandleMsgWithdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	walletAddress, err := sdk.AccAddressFromBech32(msg.WalletAddress)
	if err != nil {
		return &types.MsgWithdrawResponse{}, errors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	targetAddress, err := sdk.AccAddressFromBech32(msg.TargetAddress)
	if err != nil {
		return &types.MsgWithdrawResponse{}, errors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	err = k.Withdraw(ctx, msg.Amount, walletAddress, targetAddress)
	if err != nil {
		return &types.MsgWithdrawResponse{}, errors.Wrap(types.ErrWithdrawFailure, err.Error())
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventWithdraw{
		Amount:        msg.GetAmount(),
		WalletAddress: msg.GetWalletAddress(),
		TargetAddress: msg.GetTargetAddress(),
	})
	if err != nil {
		return nil, errors.Wrap(types.ErrEmitEvent, err.Error())
	}

	return &types.MsgWithdrawResponse{}, nil
}

func (k msgServer) HandleMsgFoundationDeposit(goCtx context.Context, msg *types.MsgFoundationDeposit) (*types.MsgFoundationDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return &types.MsgFoundationDepositResponse{}, errors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	err = k.FoundationDeposit(ctx, msg.Amount, from)
	if err != nil {
		return &types.MsgFoundationDepositResponse{}, errors.Wrap(types.ErrFoundationDepositFailure, err.Error())
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventFoundationDeposit{
		Amount: msg.GetAmount(),
	})
	if err != nil {
		return nil, errors.Wrap(types.ErrEmitEvent, err.Error())
	}

	return &types.MsgFoundationDepositResponse{}, nil
}

func (k msgServer) HandleMsgSlashingResourceNode(goCtx context.Context, msg *types.MsgSlashingResourceNode) (*types.MsgSlashingResourceNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.Reporters) == 0 || len(msg.ReporterOwner) == 0 {
		return &types.MsgSlashingResourceNodeResponse{}, types.ErrReporterAddressOrOwner
	}
	reporterOwners := msg.ReporterOwner
	for idx, reporter := range msg.Reporters {
		reporterSdsAddr, err := stratos.SdsAddressFromBech32(reporter)
		if err != nil {
			return &types.MsgSlashingResourceNodeResponse{}, errors.Wrap(types.ErrReporterAddress, err.Error())
		}
		ownerAddr, err := sdk.AccAddressFromBech32(reporterOwners[idx])
		if err != nil {
			return &types.MsgSlashingResourceNodeResponse{}, errors.Wrap(types.ErrReporterOwnerAddr, err.Error())
		}

		if !(k.registerKeeper.OwnMetaNode(ctx, ownerAddr, reporterSdsAddr)) {
			return &types.MsgSlashingResourceNodeResponse{}, types.ErrReporterAddressOrOwner
		}
	}
	networkAddress, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return &types.MsgSlashingResourceNodeResponse{}, errors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	walletAddress, err := sdk.AccAddressFromBech32(msg.WalletAddress)
	if err != nil {
		return &types.MsgSlashingResourceNodeResponse{}, errors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	nozAmt, ok := sdkmath.NewIntFromString(msg.Slashing.String())
	if !ok {
		return &types.MsgSlashingResourceNodeResponse{}, types.ErrInvalidAmount
	}

	tokenAmt, nodeType, err := k.SlashingResourceNode(ctx, networkAddress, walletAddress, nozAmt, msg.Suspend)
	if err != nil {
		return &types.MsgSlashingResourceNodeResponse{}, errors.Wrap(types.ErrSlashingResourceNodeFailure, err.Error())
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventSlashing{
		WalletAddress:  msg.GetWalletAddress(),
		NetworkAddress: msg.GetNetworkAddress(),
		Amount:         tokenAmt,
		SlashingType:   nodeType.String(),
		Suspend:        msg.GetSuspend(),
	})
	if err != nil {
		return nil, errors.Wrap(types.ErrEmitEvent, err.Error())
	}

	return &types.MsgSlashingResourceNodeResponse{}, nil
}
