package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
)

type NodeType uint32

const (
	STORAGE     NodeType = 4
	DATABASE    NodeType = 2
	COMPUTATION NodeType = 1
)

func (n NodeType) Type() string {
	switch n {
	case 7:
		return "storage/database/computation"
	case 6:
		return "database/storage"
	case 5:
		return "computation/storage"
	case 4:
		return "storage"
	case 3:
		return "computation/database"
	case 2:
		return "database"
	case 1:
		return "computation"
	}
	return "UNKNOWN"
}

func (n NodeType) String() string {
	return n.Type()
}

// ResourceNodes is a collection of resource node
type ResourceNodes []ResourceNode

func NewResourceNodes(resourceNodes ...ResourceNode) ResourceNodes {
	if len(resourceNodes) == 0 {
		return ResourceNodes{}
	}
	return resourceNodes
}

func (v ResourceNodes) String() (out string) {
	for _, node := range v {
		out += node.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func (v ResourceNodes) Validate() error {
	for _, node := range v {
		if err := node.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// NewResourceNode - initialize a new resource node
func NewResourceNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, ownerAddr sdk.AccAddress,
	description *Description, nodeType NodeType, creationTime time.Time) (ResourceNode, error) {
	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return ResourceNode{}, err
	}
	return ResourceNode{
		NetworkAddress:  networkAddr.String(),
		Pubkey:          pkAny,
		Suspend:         true,
		Status:          stakingtypes.Unbonded,
		Tokens:          sdk.ZeroInt(),
		OwnerAddress:    ownerAddr.String(),
		Description:     description,
		NodeType:        uint32(nodeType),
		CreationTime:    creationTime,
		EffectiveTokens: sdk.ZeroInt(),
	}, nil
}

// ConvertToString returns a human-readable string representation of a resource node.
func (v ResourceNode) ConvertToString() string {
	pkAny, err := codectypes.NewAnyWithValue(v.GetPubkey())
	if err != nil {
		return ErrUnknownPubKey.Error()
	}
	pubKey, err := stratos.SdsPubKeyFromBech32(pkAny.String())
	if err != nil {
		return ErrUnknownPubKey.Error()
	}
	return fmt.Sprintf(`ResourceNode:{
		Network Id:	        %s
		Pubkey:				%s
		Suspend:			%v
		Status:				%s
		Tokens:				%s
		Owner Address: 		%s
		NodeType:           %v
		Description:		%s
		CreationTime:		%s
	}`, v.GetNetworkAddress(), pubKey, v.GetSuspend(), v.GetStatus(), v.Tokens,
		v.GetOwnerAddress(), v.NodeType, v.GetDescription(), v.GetCreationTime())
}

// AddToken adds tokens to a resource node
func (v ResourceNode) AddToken(amount sdk.Int) ResourceNode {
	v.Tokens = v.Tokens.Add(amount)
	return v
}

// SubToken removes tokens from a resource node
func (v ResourceNode) SubToken(amount sdk.Int) ResourceNode {
	if amount.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", amount))
	}
	if v.Tokens.LT(amount) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, amount))
	}
	v.Tokens = v.Tokens.Sub(amount)
	return v
}

func (v ResourceNode) Validate() error {
	netAddr, err := stratos.SdsAddressFromBech32(v.GetNetworkAddress())
	if err != nil {
		return err
	}

	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}
	pkAny, err := codectypes.NewAnyWithValue(v.GetPubkey())
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
	if len(pkAny.String()) == 0 {
		return ErrEmptyPubKey
	}

	ownerAddr, err := sdk.AccAddressFromBech32(v.GetOwnerAddress())
	if err != nil {
		panic(err)
	}

	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	if v.Tokens.LT(sdk.ZeroInt()) {
		return ErrValueNegative
	}
	if v.GetDescription().Moniker == "" {
		return ErrEmptyMoniker
	}

	nodeType := NodeType(v.GetNodeType())
	if nodeType.Type() == "UNKNOWN" {
		return ErrInvalidNodeType
	}
	return nil
}

// IsBonded checks if the node status equals Bonded
func (v ResourceNode) IsBonded() bool {
	return v.GetStatus() == stakingtypes.Bonded
}

// IsUnBonded checks if the node status equals Unbonded
func (v ResourceNode) IsUnBonded() bool {
	return v.GetStatus() == stakingtypes.Unbonded
}

// IsUnBonding checks if the node status equals Unbonding
func (v ResourceNode) IsUnBonding() bool {
	return v.GetStatus() == stakingtypes.Unbonding
}

// MustMarshalResourceNode returns the resourceNode bytes. Panics if fails
func MustMarshalResourceNode(cdc codec.Codec, resourceNode ResourceNode) []byte {
	return cdc.MustMarshal(&resourceNode)
}

// MustUnmarshalResourceNode unmarshal a resourceNode from a store value. Panics if fails
func MustUnmarshalResourceNode(cdc codec.Codec, value []byte) ResourceNode {
	resourceNode, err := UnmarshalResourceNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return resourceNode
}

// UnmarshalResourceNode unmarshal a resourceNode from a store value
func UnmarshalResourceNode(cdc codec.Codec, value []byte) (v ResourceNode, err error) {
	err = cdc.Unmarshal(value, &v)
	return v, err
}

func (s *Staking) GetOwnerAddr() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(s.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return addr
}
func (s *Staking) GetShares() sdk.Dec { return s.Value }

// Stakings is a collection of Staking
type Stakings []Staking

func (ss Stakings) String() (out string) {
	for _, staking := range ss {
		out += staking.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v ResourceNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.Pubkey, &pk)
}
