package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MsgType = "volume_report"

// verify interface at compile time
var (
	_ sdk.Msg = &MsgVolumeReport{}
	_ sdk.Msg = &MsgWithdraw{}
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

type ReportRecord struct {
	Reporter        sdk.AccAddress
	ReportReference string
	TxHash          string
	NodesVolume     []SingleNodeVolume
}

func NewReportRecord(reporter sdk.AccAddress, reportReference string, txHash string, nodesVolume []SingleNodeVolume) ReportRecord {
	return ReportRecord{
		Reporter:        reporter,
		ReportReference: reportReference,
		TxHash:          txHash,
		NodesVolume:     nodesVolume,
	}
}

// Route Implement
func (msg MsgVolumeReport) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgVolumeReport) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	//addrs = append(addrs, msg.Reporter)
	addrs = append(addrs, msg.ReporterOwner)
	return addrs
}

// Type Implement
func (msg MsgVolumeReport) Type() string { return MsgType }

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
func (msg MsgWithdraw) Type() string { return "withdraw" }

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
