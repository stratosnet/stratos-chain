package types

import (
	"fmt"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	CACHE_NOT_DIRTY = 0
	CACHE_DIRTY     = 1
)

type CacheStatus struct {
	Status uint8
}

// NewMetaNode - initialize a new meta node
func NewMetaNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, ownerAddr sdk.AccAddress,
	beneficiaryAddress sdk.AccAddress, description Description, creationTime time.Time) (MetaNode, error) {

	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return MetaNode{}, err
	}
	return MetaNode{
		NetworkAddress:     networkAddr.String(),
		Pubkey:             pkAny,
		Suspend:            true,
		Status:             stakingtypes.Unbonded,
		Tokens:             sdkmath.ZeroInt(),
		OwnerAddress:       ownerAddr.String(),
		BeneficiaryAddress: beneficiaryAddress.String(),
		Description:        description,
		CreationTime:       creationTime,
	}, nil
}

// ConvertToString returns a human-readable string representation of an meta node.
func (v MetaNode) ConvertToString() string {
	pkAny, err := codectypes.NewAnyWithValue(v.GetPubkey())
	if err != nil {
		return ErrUnknownPubKey.Error()
	}
	cachedPubkey := pkAny.GetCachedValue()
	pk := cachedPubkey.(cryptotypes.PubKey)

	pubKey, err := stratos.SdsPubKeyToBech32(pk)
	if err != nil {
		return ErrUnknownPubKey.Error()
	}
	return fmt.Sprintf(`MetaNode:{
		Network Id:			%s
 		Pubkey:				%s
 		Suspend:			%v
 		Status:				%s
 		Tokens:				%s
		Owner Address: 		%s
        Beneficiary Address %s
 		Description:		%s
		CreationTime:		%s
	}`, v.GetNetworkAddress(), pubKey, v.GetSuspend(), v.GetStatus(),
		v.Tokens, v.GetOwnerAddress(), v.GetBeneficiaryAddress(), v.GetDescription(), v.GetCreationTime())
}

// AddToken adds tokens to a meta node
func (v MetaNode) AddToken(amount sdkmath.Int) MetaNode {
	v.Tokens = v.Tokens.Add(amount)
	return v
}

// SubToken removes tokens from a meta node
func (v MetaNode) SubToken(amount sdkmath.Int) MetaNode {
	if amount.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", amount))
	}
	if v.Tokens.LT(amount) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, amount))
	}
	v.Tokens = v.Tokens.Sub(amount)
	return v
}

func (v MetaNode) Validate() error {
	netAddr, err := stratos.SdsAddressFromBech32(v.GetNetworkAddress())
	if err != nil {
		return err
	}

	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}
	pkAny := v.GetPubkey()

	pubkey, ok := pkAny.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return ErrUnknownPubKey
	}

	sdsAddr := stratos.SdsAddress(pubkey.Address())

	if !netAddr.Equals(sdsAddr) {
		return ErrInvalidNetworkAddr
	}
	if len(pkAny.String()) == 0 {
		return ErrEmptyPubKey
	}

	ownerAddr, err := sdk.AccAddressFromBech32(v.GetOwnerAddress())
	if err != nil {
		return err
	}

	if ownerAddr.Empty() {
		return ErrEmptyOwnerAddr
	}

	_, err = sdk.AccAddressFromBech32(v.GetBeneficiaryAddress())
	if err != nil {
		return err
	}

	if v.Tokens.LT(sdkmath.ZeroInt()) {
		return ErrValueNegative
	}
	if v.GetDescription().Moniker == "" {
		return ErrEmptyMoniker
	}
	return nil
}

// IsBonded checks if the node status equals Bonded
func (v MetaNode) IsBonded() bool {
	return v.GetStatus() == stakingtypes.Bonded
}

// IsUnBonded checks if the node status equals Unbonded
func (v MetaNode) IsUnBonded() bool {
	return v.GetStatus() == stakingtypes.Unbonded
}

// IsUnBonding checks if the node status equals Unbonding
func (v MetaNode) IsUnBonding() bool {
	return v.GetStatus() == stakingtypes.Unbonding
}

// IsActivate get active state defined in MetaNode
func (v MetaNode) IsActivate() bool {
	return v.GetStatus() == stakingtypes.Bonded && !v.GetSuspend()
}

// MustMarshalMetaNode returns the metaNode bytes. Panics if fails
func MustMarshalMetaNode(cdc codec.Codec, metaNode MetaNode) []byte {
	return cdc.MustMarshal(&metaNode)
}

// MustUnmarshalMetaNode unmarshal an meta node from a store value. Panics if fails
func MustUnmarshalMetaNode(cdc codec.Codec, value []byte) MetaNode {
	metaNode, err := UnmarshalMetaNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return metaNode
}

// UnmarshalMetaNode unmarshal an meta node from a store value
func UnmarshalMetaNode(cdc codec.Codec, value []byte) (metaNode MetaNode, err error) {
	err = cdc.Unmarshal(value, &metaNode)
	return metaNode, err
}

// MetaNodes is a collection of meta node
type MetaNodes []MetaNode

func NewMetaNodes(metaNodes ...MetaNode) MetaNodes {
	if len(metaNodes) == 0 {
		return MetaNodes{}
	}
	return metaNodes
}

func (v MetaNodes) String() (out string) {
	for _, node := range v {
		out += node.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func (v MetaNodes) Validate() error {
	for _, node := range v {
		if err := node.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type VoteOpinion bool

const (
	Approve            VoteOpinion = true
	Reject             VoteOpinion = false
	VoteOpinionApprove             = "Approve"
	VoteOpinionReject              = "Reject"
)

func VoteOpinionFromBool(b bool) VoteOpinion {
	if b {
		return Approve
	} else {
		return Reject
	}
}

// Equal compares two VoteOpinion instances
func (v VoteOpinion) Equal(v2 VoteOpinion) bool {
	return v == v2
}

// String implements the Stringer interface for VoteOpinion.
func (v VoteOpinion) String() string {
	if v {
		return VoteOpinionApprove
	} else {
		return VoteOpinionReject
	}
}

func NewRegistrationVotePool(nodeAddress stratos.SdsAddress, expireTime time.Time) MetaNodeRegistrationVotePool {
	var approveList = make([]string, 0)
	var rejectList = make([]string, 0)
	return MetaNodeRegistrationVotePool{
		NetworkAddress: nodeAddress.String(),
		ApproveList:    approveList,
		RejectList:     rejectList,
		ExpireTime:     expireTime,
		IsVotePassed:   false,
	}
}

func NewKickMetaNodeVotePool(targetNetworkAddress stratos.SdsAddress, expireTime time.Time) KickMetaNodeVotePool {
	var approveList = make([]string, 0)
	var rejectList = make([]string, 0)
	return KickMetaNodeVotePool{
		TargetNetworkAddress: targetNetworkAddress.String(),
		ApproveList:          approveList,
		RejectList:           rejectList,
		ExpireTime:           expireTime,
		IsVotePassed:         false,
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v MetaNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.Pubkey, &pk)
}
