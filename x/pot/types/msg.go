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
	NodesVolume     []SingleNodeVolume `json:"nodes_volume" yaml:"nodes_volume"`         // volume report
	Reporter        sdk.AccAddress     `json:"reporter" yaml:"reporter"`                 // node address of the reporter
	Epoch           sdk.Int            `json:"report_epoch" yaml:"report_epoch"`         // volume report epoch
	ReportReference string             `json:"report_reference" yaml:"report_reference"` // volume report reference
	ReporterOwner   sdk.AccAddress     `json:"reporter_owner" yaml:"reporter_owner"`     // owner address of the reporter
}

// NewMsgVolumeReport creates a new Msg<Action> instance
func NewMsgVolumeReport(
	nodesVolume []SingleNodeVolume,
	reporter sdk.AccAddress,
	epoch sdk.Int,
	reportReference string,
	reporterOwner sdk.AccAddress,
) MsgVolumeReport {
	return MsgVolumeReport{
		NodesVolume:     nodesVolume,
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
	NodesVolume     []SingleNodeVolume
}

func NewQueryVolumeReportRecord(reporter sdk.AccAddress, reportReference string, txHash string, nodesVolume []SingleNodeVolume) QueryVolumeReportRecord {
	return QueryVolumeReportRecord{
		Reporter:        reporter,
		ReportReference: reportReference,
		TxHash:          txHash,
		NodesVolume:     nodesVolume,
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
	if !(len(msg.NodesVolume) > 0) {
		return ErrEmptyNodesVolume
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

	for _, item := range msg.NodesVolume {
		if item.Volume.IsNegative() {
			return ErrNegativeVolume
		}
		if item.NodeAddress.Empty() {
			return ErrMissingNodeAddress
		}
	}
	return nil
}

type MsgWithdraw struct {
	Amount       sdk.Coin       `json:"amount" yaml:"amount"`
	NodeAddress  sdk.AccAddress `json:"node_address" yaml:"node_address"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"`
}

func NewMsgWithdraw(amount sdk.Coin, nodeAddress sdk.AccAddress, ownerAddress sdk.AccAddress) MsgWithdraw {
	return MsgWithdraw{
		Amount:       amount,
		NodeAddress:  nodeAddress,
		OwnerAddress: ownerAddress,
	}
}

// Route Implement
func (msg MsgWithdraw) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
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
	if msg.NodeAddress.Empty() {
		return ErrMissingNodeAddress
	}
	if msg.OwnerAddress.Empty() {
		return ErrMissingOwnerAddress
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
