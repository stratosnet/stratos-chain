package types

import (
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
	_ sdk.Msg = &MsgUpdateResourceNodeStake{}
	_ sdk.Msg = &MsgUpdateEffectiveStake{}
	_ sdk.Msg = &MsgCreateMetaNode{}
	_ sdk.Msg = &MsgRemoveMetaNode{}
	_ sdk.Msg = &MsgUpdateMetaNode{}
	_ sdk.Msg = &MsgUpdateMetaNodeStake{}
	_ sdk.Msg = &MsgMetaNodeRegistrationVote{}
)

// message type and route constants
const (
	TypeMsgCreateResourceNodeTx    = "create_resource_node"
	TypeMsgRemoveResourceNodeTx    = "remove_resource_node"
	TypeUpdateResourceNodeTx       = "update_resource_node"
	TypeUpdateResourceNodeStakeTx  = "update_resource_node_stake"
	TypeUpdateEffectiveStakeTx     = "update_effective_stake"
	TypeCreateMetaNodeTx           = "create_meta_node"
	TypeRemoveMetaNodeTx           = "remove_meta_node"
	TypeUpdateMetaNodeTx           = "update_meta_node"
	TypeUpdateMetaNodeStakeTx      = "update_meta_node_stake"
	TypeMetaNodeRegistrationVoteTx = "meta_node_registration_vote"
)

// NewMsgCreateResourceNode NewMsg<Action> creates a new Msg<Action> instance
func NewMsgCreateResourceNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	value sdk.Coin, ownerAddr sdk.AccAddress, description *Description, nodeType uint32,
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
		NetworkAddress: networkAddr.String(),
		Pubkey:         pkAny,
		Value:          value,
		OwnerAddress:   ownerAddr.String(),
		Description:    description,
		NodeType:       nodeType,
	}, nil
}

func (msg MsgCreateResourceNode) Route() string { return RouterKey }

func (msg MsgCreateResourceNode) Type() string { return TypeMsgCreateResourceNodeTx }

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

	if *msg.GetDescription() == (Description{}) {
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

// NewMsgCreateMetaNode creates a new Msg<Action> instance
func NewMsgCreateMetaNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	value sdk.Coin, ownerAddr sdk.AccAddress, description *Description,
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
		NetworkAddress: networkAddr.String(),
		Pubkey:         pkAny,
		Value:          value,
		OwnerAddress:   ownerAddr.String(),
		Description:    description,
	}, nil
}

func (msg MsgCreateMetaNode) Route() string { return RouterKey }

func (msg MsgCreateMetaNode) Type() string { return TypeCreateMetaNodeTx }

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
	if !msg.GetValue().IsPositive() {
		return ErrValueNegative
	}

	if msg.GetDescription().Moniker == "" {
		return ErrEmptyMoniker
	}

	if *msg.GetDescription() == (Description{}) {
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
func (msg MsgRemoveResourceNode) Type() string { return TypeMsgRemoveResourceNodeTx }

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
func (msg MsgRemoveMetaNode) Type() string { return TypeRemoveMetaNodeTx }

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

func NewMsgUpdateResourceNode(description Description, nodeType uint32,
	networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress) *MsgUpdateResourceNode {

	return &MsgUpdateResourceNode{
		Description:    description,
		NodeType:       nodeType,
		NetworkAddress: networkAddress.String(),
		OwnerAddress:   ownerAddress.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) Type() string { return TypeUpdateResourceNodeTx }

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

	//if msg.Description.Moniker == "" {
	//	return ErrEmptyMoniker
	//}

	nodeType := NodeType(msg.NodeType)
	if nodeType.Type() == "UNKNOWN" {
		return ErrInvalidNodeType
	}
	return nil
}

func NewMsgUpdateResourceNodeStake(networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	stakeDelta sdk.Coin, incrStake bool) *MsgUpdateResourceNodeStake {
	return &MsgUpdateResourceNodeStake{
		NetworkAddress: networkAddress.String(),
		OwnerAddress:   ownerAddress.String(),
		StakeDelta:     stakeDelta,
		IncrStake:      incrStake,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) Type() string { return TypeUpdateResourceNodeStakeTx }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) ValidateBasic() error {
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

	if msg.StakeDelta.Amount.LTE(sdk.ZeroInt()) {
		return ErrInvalidStakeChange
	}
	return nil
}

func NewMsgUpdateMetaNode(description Description, networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
) *MsgUpdateMetaNode {

	return &MsgUpdateMetaNode{
		Description:    description,
		NetworkAddress: networkAddress.String(),
		OwnerAddress:   ownerAddress.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) Type() string { return TypeUpdateMetaNodeTx }

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

func NewMsgUpdateMetaNodeStake(networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	stakeDelta sdk.Coin, incrStake bool) *MsgUpdateMetaNodeStake {
	return &MsgUpdateMetaNodeStake{
		NetworkAddress: networkAddress.String(),
		OwnerAddress:   ownerAddress.String(),
		StakeDelta:     stakeDelta,
		IncrStake:      incrStake,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeStake) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeStake) Type() string { return TypeUpdateMetaNodeStakeTx }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeStake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateMetaNodeStake) ValidateBasic() error {
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

	if msg.StakeDelta.Amount.LTE(sdk.ZeroInt()) {
		return ErrInvalidStakeChange
	}
	return nil
}

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

func (mmsg MsgMetaNodeRegistrationVote) Route() string { return RouterKey }

func (msg MsgMetaNodeRegistrationVote) Type() string { return TypeMetaNodeRegistrationVoteTx }

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

func NewMsgUpdateEffectiveStake(reporters []stratos.SdsAddress, reporterOwner []sdk.AccAddress,
	networkAddress stratos.SdsAddress, newEffectiveStake sdk.Int) *MsgUpdateEffectiveStake {

	reporterStrSlice := make([]string, 0)
	for _, reporter := range reporters {
		reporterStrSlice = append(reporterStrSlice, reporter.String())
	}

	reporterOwnerStrSlice := make([]string, 0)
	for _, reporterOwner := range reporterOwner {
		reporterOwnerStrSlice = append(reporterOwnerStrSlice, reporterOwner.String())
	}
	return &MsgUpdateEffectiveStake{
		Reporters:       reporterStrSlice,
		ReporterOwner:   reporterOwnerStrSlice,
		NetworkAddress:  networkAddress.String(),
		EffectiveTokens: newEffectiveStake,
	}
}

func (m MsgUpdateEffectiveStake) Route() string {
	return RouterKey
}

func (m MsgUpdateEffectiveStake) Type() string {
	return "update_effective_stake"
}

func (m MsgUpdateEffectiveStake) ValidateBasic() error {
	if len(m.NetworkAddress) == 0 {
		return ErrInvalidNetworkAddr
	}
	if len(m.Reporters) == 0 {
		return ErrReporterAddress
	}
	if len(m.ReporterOwner) == 0 || len(m.Reporters) != len(m.ReporterOwner) {
		return ErrInvalidOwnerAddr
	}
	for _, r := range m.Reporters {
		if len(r) == 0 {
			return ErrReporterAddress
		}
	}

	for _, owner := range m.ReporterOwner {
		_, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			return ErrInvalidOwnerAddr
		}
	}

	if m.EffectiveTokens.LT(sdk.ZeroInt()) {
		return ErrInvalidAmount
	}
	return nil
}

func (m MsgUpdateEffectiveStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUpdateEffectiveStake) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	for _, owner := range m.ReporterOwner {
		reporterOwner, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			panic(err)
		}
		addrs = append(addrs, reporterOwner)
	}
	if len(addrs) == 0 {
		panic("no valid signer for MsgUpdateEffectiveStake")
	}
	return addrs
}
