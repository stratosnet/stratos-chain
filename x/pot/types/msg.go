package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	VolumeReportMsgType      = "volume_report"
	WithdrawMsgType          = "withdraw"
	FoundationDepositMsgType = "foundation_deposit"
)

// verify interface at compile time
var (
	_ sdk.Msg = &MsgVolumeReport{}
	_ sdk.Msg = &MsgWithdraw{}
	_ sdk.Msg = &MsgFoundationDeposit{}
)

type MsgVolumeReport struct {
	WalletVolumes   []SingleWalletVolume `json:"wallet_volumes" yaml:"wallet_volumes"`     // volume report
	Reporter        sdk.AccAddress       `json:"reporter" yaml:"reporter"`                 // node p2p address of the reporter
	Epoch           sdk.Int              `json:"epoch" yaml:"epoch"`                       // volume report epoch
	ReportReference string               `json:"report_reference" yaml:"report_reference"` // volume report reference
	ReporterOwner   sdk.AccAddress       `json:"reporter_owner" yaml:"reporter_owner"`     // owner address of the reporter
}

// NewMsgVolumeReport creates a new Msg<Action> instance
func NewMsgVolumeReport(
	walletVolumes []SingleWalletVolume,
	reporter sdk.AccAddress,
	epoch sdk.Int,
	reportReference string,
	reporterOwner sdk.AccAddress,
) MsgVolumeReport {
	return MsgVolumeReport{
		WalletVolumes:   walletVolumes,
		Reporter:        reporter,
		Epoch:           epoch,
		ReportReference: reportReference,
		ReporterOwner:   reporterOwner,
	}
}

type QueryVolumeReportRecord struct {
	Reporter        sdk.AccAddress
	ReportReference string
	TxHash          string
	walletVolumes   []SingleWalletVolume
}

func NewQueryVolumeReportRecord(reporter sdk.AccAddress, reportReference string, txHash string, walletVolumes []SingleWalletVolume) QueryVolumeReportRecord {
	return QueryVolumeReportRecord{
		Reporter:        reporter,
		ReportReference: reportReference,
		TxHash:          txHash,
		walletVolumes:   walletVolumes,
	}
}

type VolumeReportRecord struct {
	Reporter        sdk.AccAddress
	ReportReference string
	TxHash          string
}

func NewReportRecord(reporter sdk.AccAddress, reportReference string, txHash string) VolumeReportRecord {
	return VolumeReportRecord{
		Reporter:        reporter,
		ReportReference: reportReference,
		TxHash:          txHash,
	}
}

// Route Implement
func (msg MsgVolumeReport) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgVolumeReport) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	addrs = append(addrs, msg.ReporterOwner)
	return addrs
}

// Type Implement
func (msg MsgVolumeReport) Type() string { return VolumeReportMsgType }

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgVolumeReport) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgVolumeReport) ValidateBasic() error {
	if msg.Reporter.Empty() {
		return ErrEmptyReporterAddr
	}
	if !(len(msg.WalletVolumes) > 0) {
		return ErrEmptyWalletVolumes
	}

	if !(msg.Epoch.IsPositive()) {
		return ErrEpochNotPositive
	}

	if !(len(msg.ReportReference) > 0) {
		return ErrEmptyReportReference
	}
	if msg.ReporterOwner.Empty() {
		return ErrEmptyReporterOwnerAddr
	}

	for _, item := range msg.WalletVolumes {
		if item.Volume.IsNegative() {
			return ErrNegativeVolume
		}
		if item.WalletAddress.Empty() {
			return ErrMissingWalletAddress
		}
	}
	return nil
}

type MsgWithdraw struct {
	Amount        sdk.Coin       `json:"amount" yaml:"amount"`
	WalletAddress sdk.AccAddress `json:"wallet_address" yaml:"wallet_address"`
	TargetAddress sdk.AccAddress `json:"target_address" yaml:"target_address"`
}

func NewMsgWithdraw(amount sdk.Coin, walletAddress sdk.AccAddress, targetAddress sdk.AccAddress) MsgWithdraw {
	return MsgWithdraw{
		Amount:        amount,
		WalletAddress: walletAddress,
		TargetAddress: targetAddress,
	}
}

// Route Implement
func (msg MsgWithdraw) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.WalletAddress}
}

// Type Implement
func (msg MsgWithdraw) Type() string { return WithdrawMsgType }

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgWithdraw) ValidateBasic() error {
	if !(msg.Amount.IsPositive()) {
		return ErrWithdrawAmountNotPositive
	}
	if msg.WalletAddress.Empty() {
		return ErrMissingWalletAddress
	}
	if msg.TargetAddress.Empty() {
		return ErrMissingTargetAddress
	}
	return nil
}

type MsgFoundationDeposit struct {
	Amount sdk.Coin       `json:"amount" yaml:"amount"`
	From   sdk.AccAddress `json:"from" yaml:"from"`
}

func NewMsgFoundationDeposit(amount sdk.Coin, from sdk.AccAddress) MsgFoundationDeposit {
	return MsgFoundationDeposit{
		Amount: amount,
		From:   from,
	}
}

// Route Implement
func (msg MsgFoundationDeposit) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgFoundationDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// Type Implement
func (msg MsgFoundationDeposit) Type() string { return FoundationDepositMsgType }

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgFoundationDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgFoundationDeposit) ValidateBasic() error {
	if !(msg.Amount.IsPositive()) {
		return ErrWithdrawAmountNotPositive
	}
	if msg.From.Empty() {
		return ErrEmptyFromAddr
	}
	return nil
}
