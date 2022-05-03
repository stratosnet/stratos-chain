package types

import (
	"strconv"

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
	_ sdk.Msg = &MsgCreateIndexingNode{}
	_ sdk.Msg = &MsgRemoveIndexingNode{}
	_ sdk.Msg = &MsgUpdateIndexingNode{}
	_ sdk.Msg = &MsgUpdateIndexingNodeStake{}
	_ sdk.Msg = &MsgIndexingNodeRegistrationVote{}
)

// message type and route constants
const (
	// TypeMsgCreateResourceNodeTx defines the type string of an CreateResourceNodeTx transaction
	TypeMsgCreateResourceNodeTx        = "create_resource_node"
	TypeMsgRemoveResourceNodeTx        = "remove_resource_node"
	TypeUpdateResourceNodeTx           = "update_resource_node"
	TypeUpdateResourceNodeStakeTx      = "update_resource_node_stake"
	TypeCreateIndexingNodeTx           = "create_indexing_node"
	TypeRemoveIndexingNodeTx           = "remove_indexing_node"
	TypeUpdateIndexingNodeTx           = "update_indexing_node"
	TypeUpdateIndexingNodeStakeTx      = "update_indexing_node_stake"
	TypeIndexingNodeRegistrationVoteTx = "indexing_node_registration_vote"
)

// NewMsgCreateResourceNode NewMsg<Action> creates a new Msg<Action> instance
func NewMsgCreateResourceNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	value sdk.Coin, ownerAddr sdk.AccAddress, description *Description, nodeType *NodeType,
) (*MsgCreateResourceNode, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateResourceNode{
		NetworkAddr:  networkAddr.String(),
		PubKey:       pkAny,
		Value:        value,
		OwnerAddress: ownerAddr.String(),
		Description:  description,
		NodeType:     nodeType.Type(),
	}, nil
}

func (msg MsgCreateResourceNode) Route() string { return RouterKey }

func (msg MsgCreateResourceNode) Type() string { return TypeMsgCreateResourceNodeTx }

// ValidateBasic validity check for the CreateResourceNode
func (msg MsgCreateResourceNode) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.GetNetworkAddr())
	if err != nil {
		return err
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	pkAny, err := codectypes.NewAnyWithValue(msg.GetPubKey())
	if err != nil {
		return err
	}
	sdsAddr, err := stratos.SdsAddressFromBech32(pkAny.String())
	if err != nil {
		return err
	}
	if !netAddr.Equals(sdsAddr) {
		return ErrInvalidNetworkAddr
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		return err
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

	nodeTypeNum, err := strconv.Atoi(msg.GetNodeType())
	if err != nil {
		return ErrInvalidNodeType
	}
	if nodeTypeNum > 7 || nodeTypeNum < 1 {
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

// NewMsgCreateIndexingNode NewMsg<Action> creates a new Msg<Action> instance
func NewMsgCreateIndexingNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	value sdk.Coin, ownerAddr sdk.AccAddress, description *Description,
) (*MsgCreateResourceNode, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateResourceNode{
		NetworkAddr:  networkAddr.String(),
		PubKey:       pkAny,
		Value:        value,
		OwnerAddress: ownerAddr.String(),
		Description:  description,
	}, nil
}

func (msg MsgCreateIndexingNode) Route() string { return RouterKey }

func (msg MsgCreateIndexingNode) Type() string { return TypeCreateIndexingNodeTx }

func (msg MsgCreateIndexingNode) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.GetNetworkAddr())
	if err != nil {
		return err
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	pkAny, err := codectypes.NewAnyWithValue(msg.GetPubKey())
	if err != nil {
		return err
	}
	sdsAddr, err := stratos.SdsAddressFromBech32(pkAny.String())
	if err != nil {
		return err
	}
	if !netAddr.Equals(sdsAddr) {
		return ErrInvalidNetworkAddr
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		return err
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

func (msg MsgCreateIndexingNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateIndexingNode) GetSigners() []sdk.AccAddress {
	// Owner pays the tx fees
	addr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}

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
		return err
	}
	if sdsAddress.Empty() {
		return ErrEmptyResourceNodeAddr
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return err
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}
	return nil
}

// NewMsgRemoveIndexingNode creates a new MsgRemoveIndexingNode instance.
func NewMsgRemoveIndexingNode(indexingNodeAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) *MsgRemoveIndexingNode {
	return &MsgRemoveIndexingNode{
		IndexingNodeAddress: indexingNodeAddr.String(),
		OwnerAddress:        ownerAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) Type() string { return TypeRemoveIndexingNodeTx }

// GetSigners implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) ValidateBasic() error {
	sdsAddress, err := stratos.SdsAddressFromBech32(msg.IndexingNodeAddress)
	if err != nil {
		return err
	}
	if sdsAddress.Empty() {
		return ErrEmptyIndexingNodeAddr
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return err
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}
	return nil
}

func NewMsgUpdateResourceNode(description Description, nodeType NodeType,
	networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress) *MsgUpdateResourceNode {

	return &MsgUpdateResourceNode{
		Description:    description,
		NodeType:       nodeType.Type(),
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
		return err
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return err
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	if msg.Description.Moniker == "" {
		return ErrEmptyMoniker
	}

	nodeTypeNum, err := strconv.Atoi(msg.NodeType)
	if err != nil {
		return ErrInvalidNodeType
	}
	if nodeTypeNum > 7 || nodeTypeNum < 1 {
		return ErrInvalidNodeType
	}
	return nil
}

func NewMsgUpdateResourceNodeStake(networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	stakeDelta *sdk.Coin, incrStake bool) *MsgUpdateResourceNodeStake {
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
		return err
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return err
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	if msg.StakeDelta.Amount.LTE(sdk.ZeroInt()) {
		return ErrInvalidStakeChange
	}
	return nil
}

func NewMsgUpdateIndexingNode(description Description, networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
) *MsgUpdateIndexingNode {

	return &MsgUpdateIndexingNode{
		Description:    description,
		NetworkAddress: networkAddress.String(),
		OwnerAddress:   ownerAddress.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) Type() string { return TypeUpdateIndexingNodeTx }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return err
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return err
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	if msg.Description.Moniker == "" {
		return ErrEmptyMoniker
	}

	return nil
}

func NewMsgUpdateIndexingNodeStake(networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	stakeDelta *sdk.Coin, incrStake bool) *MsgUpdateIndexingNodeStake {
	return &MsgUpdateIndexingNodeStake{
		NetworkAddress: networkAddress.String(),
		OwnerAddress:   ownerAddress.String(),
		StakeDelta:     stakeDelta,
		IncrStake:      incrStake,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) Type() string { return TypeUpdateIndexingNodeStakeTx }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) ValidateBasic() error {
	netAddr, err := stratos.SdsAddressFromBech32(msg.NetworkAddress)
	if err != nil {
		return err
	}
	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return err
	}
	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	if msg.StakeDelta.Amount.LTE(sdk.ZeroInt()) {
		return ErrInvalidStakeChange
	}
	return nil
}

func NewMsgIndexingNodeRegistrationVote(candidateNetworkAddress stratos.SdsAddress, candidateOwnerAddress sdk.AccAddress, opinion bool,
	voterNetworkAddress stratos.SdsAddress, voterOwnerAddress sdk.AccAddress) *MsgIndexingNodeRegistrationVote {

	return &MsgIndexingNodeRegistrationVote{
		CandidateNetworkAddress: candidateNetworkAddress.String(),
		CandidateOwnerAddress:   candidateOwnerAddress.String(),
		Opinion:                 opinion,
		VoterNetworkAddress:     voterNetworkAddress.String(),
		VoterOwnerAddress:       voterOwnerAddress.String(),
	}
}

func (mmsg MsgIndexingNodeRegistrationVote) Route() string { return RouterKey }

func (msg MsgIndexingNodeRegistrationVote) Type() string { return TypeIndexingNodeRegistrationVoteTx }

func (msg MsgIndexingNodeRegistrationVote) ValidateBasic() error {
	if msg.CandidateNetworkAddress.Empty() {
		return ErrEmptyCandidateNetworkAddr
	}
	if msg.CandidateOwnerAddress.Empty() {
		return ErrEmptyCandidateOwnerAddr
	}
	if msg.VoterNetworkAddress.Empty() {
		return ErrEmptyVoterNetworkAddr
	}
	if msg.VoterOwnerAddress.Empty() {
		return ErrEmptyVoterOwnerAddr
	}
	if msg.CandidateNetworkAddress.Equals(msg.VoterNetworkAddress) {
		return ErrSameAddr
	}
	return nil
}

func (msg MsgIndexingNodeRegistrationVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgIndexingNodeRegistrationVote) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.VoterOwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}
