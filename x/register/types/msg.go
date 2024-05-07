package types

import (
	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stratos "github.com/stratosnet/stratos-chain/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgCreateResourceNode{}
	_ sdk.Msg = &MsgRemoveResourceNode{}
	_ sdk.Msg = &MsgUpdateResourceNode{}
	_ sdk.Msg = &MsgUpdateResourceNodeDeposit{}
	_ sdk.Msg = &MsgUpdateEffectiveDeposit{}
	_ sdk.Msg = &MsgCreateMetaNode{}
	_ sdk.Msg = &MsgRemoveMetaNode{}
	_ sdk.Msg = &MsgUpdateMetaNode{}
	_ sdk.Msg = &MsgUpdateMetaNodeDeposit{}
	_ sdk.Msg = &MsgMetaNodeRegistrationVote{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// message type and route constants
const (
	TypeMsgCreateResourceNode        = "create_resource_node"
	TypeMsgRemoveResourceNode        = "remove_resource_node"
	TypeMsgUpdateResourceNode        = "update_resource_node"
	TypeMsgUpdateResourceNodeDeposit = "update_resource_node_deposit"
	TypeMsgUpdateEffectiveDeposit    = "update_effective_deposit"
	TypeMsgCreateMetaNode            = "create_meta_node"
	TypeMsgRemoveMetaNode            = "remove_meta_node"
	TypeMsgUpdateMetaNode            = "update_meta_node"
	TypeMsgUpdateMetaNodeDeposit     = "update_meta_node_deposit"
	TypeMsgMetaNodeRegistrationVote  = "meta_node_registration_vote"
	TypeMsgUpdateParams              = "update_params"
)

// NewMsgCreateResourceNode NewMsg<Action> creates a new Msg<Action> instance
func NewMsgCreateResourceNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	value sdk.Coin, ownerAddr sdk.AccAddress, beneficiaryAddr sdk.AccAddress, description Description, nodeType uint32,
) (*MsgCreateResourceNode, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	} else {
		return nil, ErrEmptyPubKey
	}

	return &MsgCreateResourceNode{
		NetworkAddress:     networkAddr.String(),
		Pubkey:             pkAny,
		Value:              value,
		OwnerAddress:       ownerAddr.String(),
		BeneficiaryAddress: beneficiaryAddr.String(),
		Description:        description,
		NodeType:           nodeType,
	}, nil
}

func (msg MsgCreateResourceNode) Route() string { return RouterKey }

func (msg MsgCreateResourceNode) Type() string { return TypeMsgCreateResourceNode }

// ValidateBasic validity check for the CreateResourceNode
func (msg MsgCreateResourceNode) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.GetNetworkAddress())
	if err != nil {
		return ErrInvalidNetworkAddr
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}
	pkAny := msg.GetPubkey().GetCachedValue().(cryptotypes.PubKey)
	sdsAddr := sdk.AccAddress(pkAny.Address())
	if !netAddr.Equals(sdsAddr) {
		return ErrInvalidNetworkAddr
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		return ErrInvalidOwnerAddr
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}
	if !msg.GetValue().IsPositive() {
		return ErrValueNegative
	}

	if msg.GetDescription().Moniker == "" {
		return ErrEmptyMoniker
	}

	if msg.GetDescription() == (Description{}) {
		return ErrEmptyDescription
	}

	nodeType := NodeType(msg.GetNodeType())
	if nodeType.Type() == "UNKNOWN" {
		return ErrInvalidNodeType
	}

	return nil
}

func (msg MsgCreateResourceNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateResourceNode) GetSigners() []sdk.AccAddress {
	// Owner pays the tx fees
	addr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateResourceNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pk)
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgCreateMetaNode creates a new Msg<Action> instance
func NewMsgCreateMetaNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	value sdk.Coin, ownerAddr sdk.AccAddress, beneficiaryAddr sdk.AccAddress, description Description,
) (*MsgCreateMetaNode, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	} else {
		return nil, ErrEmptyPubKey
	}

	return &MsgCreateMetaNode{
		NetworkAddress:     networkAddr.String(),
		Pubkey:             pkAny,
		Value:              value,
		OwnerAddress:       ownerAddr.String(),
		BeneficiaryAddress: beneficiaryAddr.String(),
		Description:        description,
	}, nil
}

func (msg MsgCreateMetaNode) Route() string { return RouterKey }

func (msg MsgCreateMetaNode) Type() string { return TypeMsgCreateMetaNode }

func (msg MsgCreateMetaNode) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.GetNetworkAddress())
	if err != nil {
		return ErrInvalidNetworkAddr
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	pkAny := msg.GetPubkey().GetCachedValue().(cryptotypes.PubKey)
	sdsAddr := sdk.AccAddress(pkAny.Address())
	if !netAddr.Equals(sdsAddr) {
		return ErrInvalidNetworkAddr
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		return ErrInvalidOwnerAddr
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	beneficiaryAddress, err := sdk.AccAddressFromBech32(msg.GetBeneficiaryAddress())
	if err != nil {
		return ErrInvalidBeneficiaryAddr
	}
	if beneficiaryAddress.Empty() {
		return ErrInvalidBeneficiaryAddr
	}

	if !msg.GetValue().IsPositive() {
		return ErrValueNegative
	}

	if msg.GetDescription().Moniker == "" {
		return ErrEmptyMoniker
	}

	if msg.GetDescription() == (Description{}) {
		return ErrEmptyDescription
	}

	return nil
}

func (msg MsgCreateMetaNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateMetaNode) GetSigners() []sdk.AccAddress {
	// Owner pays the tx fees
	addr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}

}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateMetaNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pk)
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgRemoveResourceNode creates a new MsgRemoveResourceNode instance.
func NewMsgRemoveResourceNode(resourceNodeAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) *MsgRemoveResourceNode {
	return &MsgRemoveResourceNode{
		ResourceNodeAddress: resourceNodeAddr.String(),
		OwnerAddress:        ownerAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) Type() string { return TypeMsgRemoveResourceNode }

// GetSigners implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) ValidateBasic() error {
	sdsAddress, err := stratos.SdsAddressFromBech32(msg.ResourceNodeAddress)
	if err != nil {
		return ErrInvalidNetworkAddr
	}
	if sdsAddress.Empty() {
		return ErrEmptyResourceNodeAddr
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return ErrInvalidOwnerAddr
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgRemoveMetaNode creates a new MsgRemoveMetaNode instance.
func NewMsgRemoveMetaNode(metaNodeAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) *MsgRemoveMetaNode {
	return &MsgRemoveMetaNode{
		MetaNodeAddress: metaNodeAddr.String(),
		OwnerAddress:    ownerAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgRemoveMetaNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgRemoveMetaNode) Type() string { return TypeMsgRemoveMetaNode }

// GetSigners implements the sdk.Msg interface.
func (msg MsgRemoveMetaNode) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgRemoveMetaNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRemoveMetaNode) ValidateBasic() error {
	sdsAddress, err := stratos.SdsAddressFromBech32(msg.MetaNodeAddress)
	if err != nil {
		return ErrInvalidNetworkAddr
	}
	if sdsAddress.Empty() {
		return ErrEmptyMetaNodeAddr
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return ErrInvalidOwnerAddr
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateResourceNode(description Description, networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	beneficiaryAddr sdk.AccAddress) *MsgUpdateResourceNode {

	return &MsgUpdateResourceNode{
		Description:        description,
		NetworkAddress:     networkAddress.String(),
		OwnerAddress:       ownerAddress.String(),
		BeneficiaryAddress: beneficiaryAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) Type() string { return TypeMsgUpdateResourceNode }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return ErrInvalidNetworkAddr
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return ErrInvalidOwnerAddr
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateResourceNodeDeposit(networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	depositDelta sdk.Coin) *MsgUpdateResourceNodeDeposit {
	return &MsgUpdateResourceNodeDeposit{
		NetworkAddress: networkAddress.String(),
		OwnerAddress:   ownerAddress.String(),
		DepositDelta:   depositDelta,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeDeposit) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeDeposit) Type() string { return TypeMsgUpdateResourceNodeDeposit }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeDeposit) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeDeposit) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return ErrInvalidNetworkAddr
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return ErrInvalidOwnerAddr
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	if msg.DepositDelta.Amount.LTE(sdkmath.ZeroInt()) {
		return ErrInvalidDepositChange
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateMetaNode(description Description, networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	beneficiaryAddress sdk.AccAddress) *MsgUpdateMetaNode {

	return &MsgUpdateMetaNode{
		Description:        description,
		NetworkAddress:     networkAddress.String(),
		OwnerAddress:       ownerAddress.String(),
		BeneficiaryAddress: beneficiaryAddress.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) Type() string { return TypeMsgUpdateMetaNode }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return ErrInvalidNetworkAddr
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return ErrInvalidOwnerAddr
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	if msg.Description.Moniker == "" {
		return ErrEmptyMoniker
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateMetaNodeDeposit(networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	depositDelta sdk.Coin) *MsgUpdateMetaNodeDeposit {
	return &MsgUpdateMetaNodeDeposit{
		NetworkAddress: networkAddress.String(),
		OwnerAddress:   ownerAddress.String(),
		DepositDelta:   depositDelta,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeDeposit) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeDeposit) Type() string { return TypeMsgUpdateMetaNodeDeposit }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeDeposit) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeDeposit) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return ErrInvalidNetworkAddr
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return ErrInvalidOwnerAddr
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	if msg.DepositDelta.Amount.LTE(sdkmath.ZeroInt()) {
		return ErrInvalidDepositChange
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgMetaNodeRegistrationVote(candidateNetworkAddress stratos.SdsAddress, candidateOwnerAddress sdk.AccAddress, opinion bool,
	voterNetworkAddress stratos.SdsAddress, voterOwnerAddress sdk.AccAddress) *MsgMetaNodeRegistrationVote {

	return &MsgMetaNodeRegistrationVote{
		CandidateNetworkAddress: candidateNetworkAddress.String(),
		CandidateOwnerAddress:   candidateOwnerAddress.String(),
		Opinion:                 opinion,
		VoterNetworkAddress:     voterNetworkAddress.String(),
		VoterOwnerAddress:       voterOwnerAddress.String(),
	}
}

func (msg MsgMetaNodeRegistrationVote) Route() string { return RouterKey }

func (msg MsgMetaNodeRegistrationVote) Type() string { return TypeMsgMetaNodeRegistrationVote }

func (msg MsgMetaNodeRegistrationVote) ValidateBasic() error {
	candidateNetworkAddress, err := stratos.SdsAddressFromBech32(msg.CandidateNetworkAddress)
	if err != nil {
		return ErrInvalidCandidateNetworkAddr
	}
	if candidateNetworkAddress.Empty() {
		return ErrEmptyCandidateNetworkAddr
	}

	voterNetworkAddr, err := stratos.SdsAddressFromBech32(msg.VoterNetworkAddress)
	if err != nil {
		return ErrInvalidVoterNetworkAddr
	}
	if voterNetworkAddr.Empty() {
		return ErrEmptyVoterNetworkAddr
	}

	candidateOwnerAddr, err := sdk.AccAddressFromBech32(msg.CandidateOwnerAddress)
	if err != nil {
		return ErrInvalidCandidateOwnerAddr
	}
	if candidateOwnerAddr.Empty() {
		return ErrEmptyCandidateOwnerAddr
	}

	voterOwnerAddr, err := sdk.AccAddressFromBech32(msg.VoterOwnerAddress)
	if err != nil {
		return ErrInvalidVoterOwnerAddr
	}
	if voterOwnerAddr.Empty() {
		return ErrEmptyVoterOwnerAddr
	}

	if candidateNetworkAddress.Equals(voterNetworkAddr) {
		return ErrSameAddr
	}
	return nil
}

func (msg MsgMetaNodeRegistrationVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgMetaNodeRegistrationVote) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.VoterOwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgKickMetaNodeVote(targetNetworkAddress stratos.SdsAddress, opinion bool, voterNetworkAddress stratos.SdsAddress,
	voterOwnerAddress sdk.AccAddress) *MsgKickMetaNodeVote {

	return &MsgKickMetaNodeVote{
		TargetNetworkAddress: targetNetworkAddress.String(),
		Opinion:              opinion,
		VoterNetworkAddress:  voterNetworkAddress.String(),
		VoterOwnerAddress:    voterOwnerAddress.String(),
	}
}

func (msg MsgKickMetaNodeVote) ValidateBasic() error {
	targetNetworkAddress, err := stratos.SdsAddressFromBech32(msg.GetTargetNetworkAddress())
	if err != nil {
		return ErrInvalidTargetNetworkAddr
	}
	if targetNetworkAddress.Empty() {
		return ErrInvalidTargetNetworkAddr
	}

	voterNetworkAddress, err := stratos.SdsAddressFromBech32(msg.GetVoterNetworkAddress())
	if err != nil {
		return ErrInvalidVoterNetworkAddr
	}
	if voterNetworkAddress.Empty() {
		return ErrInvalidVoterNetworkAddr
	}

	voterOwnerAddress, err := sdk.AccAddressFromBech32(msg.GetVoterOwnerAddress())
	if err != nil {
		return ErrInvalidVoterOwnerAddr
	}
	if voterOwnerAddress.Empty() {
		return ErrInvalidVoterOwnerAddr
	}

	return nil
}

func (msg MsgKickMetaNodeVote) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.GetVoterOwnerAddress())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateEffectiveDeposit(reporters []stratos.SdsAddress, reporterOwner []sdk.AccAddress,
	networkAddress stratos.SdsAddress, newEffectiveDeposit sdkmath.Int) *MsgUpdateEffectiveDeposit {

	reporterStrSlice := make([]string, 0)
	for _, reporter := range reporters {
		reporterStrSlice = append(reporterStrSlice, reporter.String())
	}

	reporterOwnerStrSlice := make([]string, 0)
	for _, reporterOwner := range reporterOwner {
		reporterOwnerStrSlice = append(reporterOwnerStrSlice, reporterOwner.String())
	}
	return &MsgUpdateEffectiveDeposit{
		Reporters:       reporterStrSlice,
		ReporterOwner:   reporterOwnerStrSlice,
		NetworkAddress:  networkAddress.String(),
		EffectiveTokens: newEffectiveDeposit,
	}
}

func (msg MsgUpdateEffectiveDeposit) Route() string {
	return RouterKey
}

func (msg MsgUpdateEffectiveDeposit) Type() string {
	return TypeMsgUpdateEffectiveDeposit
}

func (msg MsgUpdateEffectiveDeposit) ValidateBasic() error {
	if len(msg.NetworkAddress) == 0 {
		return ErrInvalidNetworkAddr
	}
	if len(msg.Reporters) == 0 {
		return ErrReporterAddress
	}
	if len(msg.ReporterOwner) == 0 || len(msg.Reporters) != len(msg.ReporterOwner) {
		return ErrInvalidOwnerAddr
	}
	for _, r := range msg.Reporters {
		if len(r) == 0 {
			return ErrReporterAddress
		}
	}

	for _, owner := range msg.ReporterOwner {
		_, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			return ErrInvalidOwnerAddr
		}
	}

	if msg.EffectiveTokens.LT(sdkmath.ZeroInt()) {
		return ErrInvalidAmount
	}
	return nil
}

func (msg MsgUpdateEffectiveDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUpdateEffectiveDeposit) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	for _, owner := range msg.ReporterOwner {
		reporterOwner, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			panic(err)
		}
		addrs = append(addrs, reporterOwner)
	}
	if len(addrs) == 0 {
		panic("no valid signer for MsgUpdateEffectiveDeposit")
	}
	return addrs
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateParams(params Params, authority string) *MsgUpdateParams {
	return &MsgUpdateParams{
		Params:    params,
		Authority: authority,
	}
}

// Route implements legacytx.LegacyMsg
func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

// Type implements legacytx.LegacyMsg
func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return err
	}

	return msg.Params.Validate()
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority := sdk.MustAccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

// GetSigners implements legacytx.LegacyMsg
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
