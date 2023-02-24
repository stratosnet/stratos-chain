package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	VolumeReportMsgType      = "volume_report"
	WithdrawMsgType          = "withdraw"
	LegacyWithdrawMsgType    = "legacy_withdraw"
	FoundationDepositMsgType = "foundation_deposit"
)

// verify interface at compile time
var (
	_ sdk.Msg = &MsgVolumeReport{}
	_ sdk.Msg = &MsgWithdraw{}
	_ sdk.Msg = &MsgLegacyWithdraw{}
	_ sdk.Msg = &MsgFoundationDeposit{}
	_ sdk.Msg = &MsgSlashingResourceNode{}
)

// NewMsgVolumeReport creates a new MsgVolumeReport instance
func NewMsgVolumeReport(
	walletVolumes []*SingleWalletVolume,
	reporter stratos.SdsAddress,
	epoch sdk.Int,
	reportReference string,
	reporterOwner sdk.AccAddress,
	blsSignature BLSSignatureInfo,
) *MsgVolumeReport {
	return &MsgVolumeReport{
		WalletVolumes:   walletVolumes,
		Reporter:        reporter.String(),
		Epoch:           &epoch,
		ReportReference: reportReference,
		ReporterOwner:   reporterOwner.String(),
		BLSSignature:    &blsSignature,
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

// Route Implement
func (msg MsgVolumeReport) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgVolumeReport) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	reporterOwner, err := sdk.AccAddressFromBech32(msg.ReporterOwner)
	if err != nil {
		panic(err)
	}
	addrs = append(addrs, reporterOwner)
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
	if len(msg.Reporter) == 0 {
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
	if len(msg.ReporterOwner) == 0 {
		return ErrEmptyReporterOwnerAddr
	}

	_, err := sdk.AccAddressFromBech32(msg.GetReporterOwner())
	if err != nil {
		return ErrReporterAddress
	}

	for _, item := range msg.WalletVolumes {
		if item.Volume.IsNegative() {
			return ErrNegativeVolume
		}
		if len(item.WalletAddress) == 0 {
			return ErrMissingWalletAddress
		}
	}

	if len(msg.BLSSignature.Signature) == 0 {
		return ErrBLSSignatureInvalid
	}
	if len(msg.BLSSignature.TxData) == 0 {
		return ErrBLSTxDataInvalid
	}
	for _, pubKey := range msg.BLSSignature.PubKeys {
		if len(pubKey) == 0 {
			return ErrBLSPubkeysInvalid
		}
	}

	return nil
}

func NewMsgWithdraw(amount sdk.Coins, walletAddress sdk.AccAddress, targetAddress sdk.AccAddress) *MsgWithdraw {
	return &MsgWithdraw{
		Amount:        amount,
		WalletAddress: walletAddress.String(),
		TargetAddress: targetAddress.String(),
	}
}

// Route Implement
func (msg MsgWithdraw) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	walletAddress, err := sdk.AccAddressFromBech32(msg.WalletAddress)
	if err != nil {
		panic(err)
	}
	addrs = append(addrs, walletAddress)
	return addrs
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
	if !(msg.Amount.IsValid()) {
		return ErrWithdrawAmountInvalid
	}
	if len(msg.WalletAddress) == 0 {
		return ErrMissingWalletAddress
	}
	if len(msg.TargetAddress) == 0 {
		return ErrMissingTargetAddress
	}
	_, err := sdk.AccAddressFromBech32(msg.GetWalletAddress())
	if err != nil {
		return ErrInvalidAddress
	}
	return nil
}

func NewMsgLegacyWithdraw(amount sdk.Coins, from sdk.AccAddress, targetAddress sdk.AccAddress) *MsgLegacyWithdraw {
	return &MsgLegacyWithdraw{
		Amount:        amount,
		From:          from.String(),
		TargetAddress: targetAddress.String(),
	}
}

// Route Implement
func (msg MsgLegacyWithdraw) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgLegacyWithdraw) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	addrs = append(addrs, from)
	return addrs
}

// Type Implement
func (msg MsgLegacyWithdraw) Type() string { return LegacyWithdrawMsgType }

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgLegacyWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgLegacyWithdraw) ValidateBasic() error {
	if !(msg.Amount.IsValid()) {
		return ErrWithdrawAmountInvalid
	}
	if len(msg.From) == 0 {
		return ErrEmptyFromAddr
	}
	if len(msg.TargetAddress) == 0 {
		return ErrMissingTargetAddress
	}
	_, err := sdk.AccAddressFromBech32(msg.GetFrom())
	if err != nil {
		return ErrInvalidAddress
	}
	return nil
}

func NewMsgFoundationDeposit(amount sdk.Coins, from sdk.AccAddress) *MsgFoundationDeposit {
	return &MsgFoundationDeposit{
		Amount: amount,
		From:   from.String(),
	}
}

// Route Implement
func (msg MsgFoundationDeposit) Route() string { return RouterKey }

// GetSigners Implement
func (msg MsgFoundationDeposit) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	addrs = append(addrs, from)
	return addrs
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
	if !(msg.Amount.IsValid()) {
		return ErrFoundationDepositAmountInvalid
	}
	if len(msg.From) == 0 {
		return ErrEmptyFromAddr
	}
	_, err := sdk.AccAddressFromBech32(msg.GetFrom())
	if err != nil {
		return ErrInvalidAddress
	}
	return nil
}

func NewMsgSlashingResourceNode(reporters []stratos.SdsAddress, reporterOwner []sdk.AccAddress,
	networkAddress stratos.SdsAddress, walletAddress sdk.AccAddress, slashing sdk.Int, suspend bool) *MsgSlashingResourceNode {

	reporterStrSlice := make([]string, 0)
	for _, reporter := range reporters {
		reporterStrSlice = append(reporterStrSlice, reporter.String())
	}

	reporterOwnerStrSlice := make([]string, 0)
	for _, reporterOwner := range reporterOwner {
		reporterOwnerStrSlice = append(reporterOwnerStrSlice, reporterOwner.String())
	}
	return &MsgSlashingResourceNode{
		Reporters:      reporterStrSlice,
		ReporterOwner:  reporterOwnerStrSlice,
		NetworkAddress: networkAddress.String(),
		WalletAddress:  walletAddress.String(),
		Slashing:       &slashing,
		Suspend:        suspend,
	}
}

func (m MsgSlashingResourceNode) Route() string {
	return RouterKey
}

func (m MsgSlashingResourceNode) Type() string {
	return "slashing_resource_node"
}

func (m MsgSlashingResourceNode) ValidateBasic() error {
	if len(m.NetworkAddress) == 0 {
		return ErrMissingTargetAddress
	}
	if len(m.WalletAddress) == 0 {
		return ErrMissingWalletAddress
	}
	for _, r := range m.Reporters {
		if len(r) == 0 {
			return ErrReporterAddress
		}
	}

	for _, owner := range m.ReporterOwner {
		_, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			return ErrInvalidAddress
		}
	}

	if m.Slashing.LT(sdk.ZeroInt()) {
		return ErrInvalidAmount
	}
	return nil
}

func (m MsgSlashingResourceNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgSlashingResourceNode) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	for _, owner := range m.ReporterOwner {
		reporterOwner, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			panic(err)
		}
		addrs = append(addrs, reporterOwner)
	}
	return addrs
}
