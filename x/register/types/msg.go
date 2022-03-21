package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/tendermint/tendermint/crypto"
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

type MsgCreateResourceNode struct {
	NetworkAddr  stratos.SdsAddress `json:"network_address" yaml:"network_address"`
	PubKey       crypto.PubKey      `json:"pubkey" yaml:"pubkey"`
	Value        sdk.Coin           `json:"value" yaml:"value"`
	OwnerAddress sdk.AccAddress     `json:"owner_address" yaml:"owner_address"`
	Description  Description        `json:"description" yaml:"description"`
	NodeType     NodeType           `json:"node_type" yaml:"node_type"`
}

// NewMsgCreateResourceNode NewMsg<Action> creates a new Msg<Action> instance
func NewMsgCreateResourceNode(networkAddr stratos.SdsAddress, pubKey crypto.PubKey, value sdk.Coin,
	ownerAddr sdk.AccAddress, description Description, nodeType NodeType,
) MsgCreateResourceNode {
	return MsgCreateResourceNode{
		NetworkAddr:  networkAddr,
		PubKey:       pubKey,
		Value:        value,
		OwnerAddress: ownerAddr,
		Description:  description,
		NodeType:     nodeType,
	}
}

func (msg MsgCreateResourceNode) Route() string {
	return RouterKey
}

func (msg MsgCreateResourceNode) Type() string {
	return "create_resource_node"
}

// ValidateBasic validity check for the CreateResourceNode
func (msg MsgCreateResourceNode) ValidateBasic() error {
	if msg.NetworkAddr.Empty() {
		return ErrEmptyNodeId
	}
	if !msg.NetworkAddr.Equals(stratos.SdsAddress(msg.PubKey.Address())) {
		return ErrInvalidNetworkAddr
	}
	if msg.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	if !msg.Value.IsPositive() {
		return ErrValueNegative
	}

	if msg.Description == (Description{}) {
		return ErrEmptyDescription
	}
	if msg.Description.Moniker == "" {
		return ErrEmptyMoniker
	}
	if msg.NodeType > 7 || msg.NodeType < 1 {
		return ErrInvalidNodeType
	}
	return nil
}

func (msg MsgCreateResourceNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateResourceNode) GetSigners() []sdk.AccAddress {
	// Owner pays the tx fees
	addrs := []sdk.AccAddress{msg.OwnerAddress}
	return addrs
}

type MsgCreateIndexingNode struct {
	NetworkAddr  stratos.SdsAddress `json:"network_addr" yaml:"network_addr"`
	PubKey       crypto.PubKey      `json:"pubkey" yaml:"pubkey"`
	Value        sdk.Coin           `json:"value" yaml:"value"`
	OwnerAddress sdk.AccAddress     `json:"owner_address" yaml:"owner_address"`
	Description  Description        `json:"description" yaml:"description"`
}

// NewMsgCreateIndexingNode NewMsg<Action> creates a new Msg<Action> instance
func NewMsgCreateIndexingNode(networkAddr stratos.SdsAddress, pubKey crypto.PubKey, value sdk.Coin, ownerAddr sdk.AccAddress, description Description,
) MsgCreateIndexingNode {
	return MsgCreateIndexingNode{
		NetworkAddr:  networkAddr,
		PubKey:       pubKey,
		Value:        value,
		OwnerAddress: ownerAddr,
		Description:  description,
	}
}

func (msg MsgCreateIndexingNode) Route() string {
	return RouterKey
}

func (msg MsgCreateIndexingNode) Type() string {
	return "create_indexing_node"
}

func (msg MsgCreateIndexingNode) ValidateBasic() error {
	if msg.NetworkAddr.Empty() {
		return ErrInvalidNetworkAddr
	}
	if !msg.NetworkAddr.Equals(stratos.SdsAddress(msg.PubKey.Address())) {
		return ErrInvalidNetworkAddr
	}
	if msg.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	if !msg.Value.IsPositive() {
		return ErrValueNegative
	}

	if msg.Description == (Description{}) {
		return ErrEmptyDescription
	}
	if msg.Description.Moniker == "" {
		return ErrEmptyMoniker
	}
	return nil
}

func (msg MsgCreateIndexingNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateIndexingNode) GetSigners() []sdk.AccAddress {
	// Owner pays the tx fees
	addrs := []sdk.AccAddress{msg.OwnerAddress}
	return addrs
}

// MsgRemoveResourceNode - struct for removing resource node
type MsgRemoveResourceNode struct {
	ResourceNodeAddress stratos.SdsAddress `json:"resource_node_address" yaml:"resource_node_address"`
	OwnerAddress        sdk.AccAddress     `json:"owner_address" yaml:"owner_address"`
}

// NewMsgRemoveResourceNode creates a new MsgRemoveResourceNode instance.
func NewMsgRemoveResourceNode(resourceNodeAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) MsgRemoveResourceNode {
	return MsgRemoveResourceNode{
		ResourceNodeAddress: resourceNodeAddr,
		OwnerAddress:        ownerAddr,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) Type() string { return "remove_resource_node" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRemoveResourceNode) ValidateBasic() error {
	if msg.ResourceNodeAddress.Empty() {
		return ErrEmptyResourceNodeAddr
	}
	if msg.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	return nil
}

// MsgRemoveIndexingNode - struct for removing indexing node
type MsgRemoveIndexingNode struct {
	IndexingNodeAddress stratos.SdsAddress `json:"indexing_node_address" yaml:"indexing_node_address"`
	OwnerAddress        sdk.AccAddress     `json:"owner_address" yaml:"owner_address"`
}

// NewMsgRemoveIndexingNode creates a new MsgRemoveIndexingNode instance.
func NewMsgRemoveIndexingNode(indexingNodeAddr stratos.SdsAddress, ownerAddr sdk.AccAddress) MsgRemoveIndexingNode {
	return MsgRemoveIndexingNode{
		IndexingNodeAddress: indexingNodeAddr,
		OwnerAddress:        ownerAddr,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) Type() string { return "remove_indexing_node" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRemoveIndexingNode) ValidateBasic() error {
	if msg.IndexingNodeAddress.Empty() {
		return ErrEmptyIndexingNodeAddr
	}
	if msg.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	return nil
}

// MsgUpdateResourceNode struct for updating resource node
type MsgUpdateResourceNode struct {
	Description    Description        `json:"description" yaml:"description"`
	NodeType       NodeType           `json:"node_type" yaml:"node_type"`
	NetworkAddress stratos.SdsAddress `json:"network_address" yaml:"network_address"`
	OwnerAddress   sdk.AccAddress     `json:"owner_address" yaml:"owner_address"`
}

func NewMsgUpdateResourceNode(description Description, nodeType NodeType,
	networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress) MsgUpdateResourceNode {

	return MsgUpdateResourceNode{
		Description:    description,
		NodeType:       nodeType,
		NetworkAddress: networkAddress,
		OwnerAddress:   ownerAddress,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) Type() string { return "update_resource_node" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) ValidateBasic() error {
	if msg.NetworkAddress.Empty() {
		return ErrInvalidNetworkAddr
	}

	if msg.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	if msg.Description.Moniker == "" {
		return ErrEmptyMoniker
	}
	if msg.NodeType > 7 || msg.NodeType < 1 {
		return ErrInvalidNodeType
	}
	return nil
}

// MsgUpdateResourceNodeStake struct for only updating resource node's stake
type MsgUpdateResourceNodeStake struct {
	NetworkAddress stratos.SdsAddress `json:"network_address" yaml:"network_address"`
	OwnerAddress   sdk.AccAddress     `json:"owner_address" yaml:"owner_address"`
	StakeDelta     sdk.Coin           `json:"stake_delta" yaml:"stake_delta"`
	IncrStake      bool               `json:"incr_stake" yaml:"incr_stake"`
}

func NewMsgUpdateResourceNodeStake(networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	stakeDelta sdk.Coin, incrStake bool) MsgUpdateResourceNodeStake {
	return MsgUpdateResourceNodeStake{
		NetworkAddress: networkAddress,
		OwnerAddress:   ownerAddress,
		StakeDelta:     stakeDelta,
		IncrStake:      incrStake,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) Type() string { return "update_resource_node_stake" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateResourceNodeStake) ValidateBasic() error {
	if msg.NetworkAddress.Empty() {
		return ErrInvalidNetworkAddr
	}
	if msg.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	if msg.StakeDelta.Amount.LTE(sdk.ZeroInt()) {
		return ErrInvalidStakeChange
	}
	return nil
}

// MsgUpdateIndexingNode struct for updating indexing node
type MsgUpdateIndexingNode struct {
	Description    Description        `json:"description" yaml:"description"`
	NetworkAddress stratos.SdsAddress `json:"network_address" yaml:"network_address"`
	OwnerAddress   sdk.AccAddress     `json:"owner_address" yaml:"owner_address"`
}

func NewMsgUpdateIndexingNode(description Description, networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
) MsgUpdateIndexingNode {

	return MsgUpdateIndexingNode{
		Description:    description,
		NetworkAddress: networkAddress,
		OwnerAddress:   ownerAddress,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) Type() string { return "update_indexing_node" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNode) ValidateBasic() error {
	if msg.NetworkAddress.Empty() {
		return ErrInvalidNetworkAddr
	}
	if msg.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	if msg.Description.Moniker == "" {
		return ErrEmptyMoniker
	}
	return nil
}

// MsgUpdateIndexingNodeStake struct for updating indexing node's stake
type MsgUpdateIndexingNodeStake struct {
	NetworkAddress stratos.SdsAddress `json:"network_address" yaml:"network_address"`
	OwnerAddress   sdk.AccAddress     `json:"owner_address" yaml:"owner_address"`
	StakeDelta     sdk.Coin           `json:"stake_delta" yaml:"stake_delta"`
	IncrStake      bool               `json:"incr_stake" yaml:"incr_stake"`
}

func NewMsgUpdateIndexingNodeStake(networkAddress stratos.SdsAddress, ownerAddress sdk.AccAddress,
	stakeDelta sdk.Coin, incrStake bool) MsgUpdateIndexingNodeStake {
	return MsgUpdateIndexingNodeStake{
		NetworkAddress: networkAddress,
		OwnerAddress:   ownerAddress,
		StakeDelta:     stakeDelta,
		IncrStake:      incrStake,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) Type() string { return "update_indexing_node" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateIndexingNodeStake) ValidateBasic() error {
	if msg.NetworkAddress.Empty() {
		return ErrInvalidNetworkAddr
	}
	if msg.OwnerAddress.Empty() {
		return ErrEmptyOwnerAddr
	}
	if msg.StakeDelta.Amount.LTE(sdk.ZeroInt()) {
		return ErrInvalidStakeChange
	}
	return nil
}

type MsgIndexingNodeRegistrationVote struct {
	CandidateNetworkAddress stratos.SdsAddress `json:"candidate_network_address" yaml:"candidate_network_address"` // node address of indexing node
	CandidateOwnerAddress   sdk.AccAddress     `json:"candidate_owner_address" yaml:"candidate_owner_address"`     // owner address of indexing node
	Opinion                 VoteOpinion        `json:"opinion" yaml:"opinion"`
	VoterNetworkAddress     stratos.SdsAddress `json:"voter_network_address" yaml:"voter_network_address"` // address of voter (other existed indexing node)
	VoterOwnerAddress       sdk.AccAddress     `json:"voter_owner_address" yaml:"voter_owner_address"`     // address of owner of the voter (other existed indexing node)
}

func NewMsgIndexingNodeRegistrationVote(candidateNetworkAddress stratos.SdsAddress, candidateOwnerAddress sdk.AccAddress, opinion VoteOpinion,
	voterNetworkAddress stratos.SdsAddress, voterOwnerAddress sdk.AccAddress) MsgIndexingNodeRegistrationVote {

	return MsgIndexingNodeRegistrationVote{
		CandidateNetworkAddress: candidateNetworkAddress,
		CandidateOwnerAddress:   candidateOwnerAddress,
		Opinion:                 opinion,
		VoterNetworkAddress:     voterNetworkAddress,
		VoterOwnerAddress:       voterOwnerAddress,
	}
}

func (m MsgIndexingNodeRegistrationVote) Route() string {
	return RouterKey
}

func (m MsgIndexingNodeRegistrationVote) Type() string {
	return "indexing_node_reg_vote"
}

func (m MsgIndexingNodeRegistrationVote) ValidateBasic() error {
	if m.CandidateNetworkAddress.Empty() {
		return ErrEmptyCandidateNetworkAddr
	}
	if m.CandidateOwnerAddress.Empty() {
		return ErrEmptyCandidateOwnerAddr
	}
	if m.VoterNetworkAddress.Empty() {
		return ErrEmptyVoterNetworkAddr
	}
	if m.VoterOwnerAddress.Empty() {
		return ErrEmptyVoterOwnerAddr
	}
	if m.CandidateNetworkAddress.Equals(m.VoterNetworkAddress) {
		return ErrSameAddr
	}
	return nil
}

func (m MsgIndexingNodeRegistrationVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgIndexingNodeRegistrationVote) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress
	addrs = append(addrs, m.VoterOwnerAddress)
	return addrs
}
