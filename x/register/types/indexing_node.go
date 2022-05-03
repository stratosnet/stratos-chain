package types

import (
	"bytes"
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	goamino "github.com/tendermint/go-amino"
)

// IndexingNodes is a collection of indexing node
//type IndexingNodes []IndexingNode

//func (v IndexingNodes) String() (out string) {
//	for _, node := range v {
//		out += node.String() + "\n"
//	}
//	return strings.TrimSpace(out)
//}

//Sort IndexingNodes sorts IndexingNode array in ascending owner address order
//func (v IndexingNodes) Sort() {
//	sort.Sort(v.GetIndexingNodes())
//}

//// Len Implements sort interface
//func (v IndexingNodes) Len() int {
//	return len(v.GetIndexingNodes())
//}
//
//// Less Implements sort interface
//func (v IndexingNodes) Less(i, j int) bool {
//	return v.GetIndexingNodes()[i].Tokens < (v.GetIndexingNodes()[j].Tokens)
//}
//
//// Swap Implements sort interface
//func (v IndexingNodes) Swap(i, j int) {
//	it := v.GetIndexingNodes()[i]
//	v.GetIndexingNodes()[i] = v.GetIndexingNodes()[j]
//	v.GetIndexingNodes()[j] = it
//}

func (v IndexingNodes) Validate() error {
	for _, node := range v.GetIndexingNodes() {
		if err := node.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// NewIndexingNode - initialize a new indexing node
func NewIndexingNode(networkAddr stratos.SdsAddress, pubKey cryptotypes.PubKey, ownerAddr sdk.AccAddress, description *Description, creationTime time.Time) (IndexingNode, error) {
	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return IndexingNode{}, err
	}
	return IndexingNode{
		NetworkAddr:  networkAddr.String(),
		PubKey:       pkAny,
		Suspend:      true,
		Status:       stakingtypes.Unbonded,
		Tokens:       sdk.ZeroInt(),
		OwnerAddress: ownerAddr.String(),
		Description:  description,
		CreationTime: creationTime,
	}, nil
}

// ConvertToString returns a human-readable string representation of an indexing node.
func (v IndexingNode) ConvertToString() string {
	pkAny, err := codectypes.NewAnyWithValue(v.GetPubKey())
	if err != nil {
		return ErrUnknownPubKey.Error()
	}
	pubKey, err := stratos.GetPubKeyFromBech32(stratos.Bech32PubKeyTypeAccPub, pkAny.String())
	if err != nil {
		return ErrUnknownPubKey.Error()
	}
	return fmt.Sprintf(`IndexingNode:{
		Network Id:			%s
 		Pubkey:				%s
 		Suspend:			%v
 		Status:				%s
 		Tokens:				%s
		Owner Address: 		%s
 		Description:		%s
		CreationTime:		%s
	}`, v.GetNetworkAddr(), pubKey, v.GetSuspend(), v.GetStatus(),
		v.Tokens, v.GetOwnerAddress(), v.GetDescription(), v.GetCreationTime())
}

// AddToken adds tokens to a indexing node
func (v IndexingNode) AddToken(amount sdk.Int) IndexingNode {
	v.Tokens = v.Tokens.Add(amount)
	return v
}

// SubToken removes tokens from a indexing node
func (v IndexingNode) SubToken(amount sdk.Int) IndexingNode {
	if amount.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", amount))
	}
	if v.Tokens.LT(amount) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, amount))
	}
	v.Tokens = v.Tokens.Sub(amount)
	return v
}

func (v IndexingNode) Validate() error {
	netAddr, err := stratos.SdsAddressFromBech32(v.GetNetworkAddr())
	if err != nil {
		return err
	}

	if netAddr.Empty() {
		return ErrEmptyNodeNetworkAddress
	}
	pkAny, err := codectypes.NewAnyWithValue(v.GetPubKey())
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
		return err
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
	return nil
}

// IsBonded checks if the node status equals Bonded
func (v IndexingNode) IsBonded() bool {
	return v.GetStatus() == stakingtypes.Bonded
}

// IsUnBonded checks if the node status equals Unbonded
func (v IndexingNode) IsUnBonded() bool {
	return v.GetStatus() == stakingtypes.Unbonded
}

// IsUnBonding checks if the node status equals Unbonding
func (v IndexingNode) IsUnBonding() bool {
	return v.GetStatus() == stakingtypes.Unbonding
}

// MustMarshalIndexingNode returns the indexingNode bytes. Panics if fails
func MustMarshalIndexingNode(cdc *goamino.Codec, indexingNode IndexingNode) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(indexingNode)
}

// MustUnmarshalIndexingNode unmarshal an indexing node from a store value. Panics if fails
func MustUnmarshalIndexingNode(cdc *goamino.Codec, value []byte) IndexingNode {
	indexingNode, err := UnmarshalIndexingNode(cdc, value)
	if err != nil {
		panic(err)
	}
	return indexingNode
}

// UnmarshalIndexingNode unmarshal an indexing node from a store value
func UnmarshalIndexingNode(cdc *goamino.Codec, value []byte) (indexingNode IndexingNode, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &indexingNode)
	return indexingNode, err
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

type IndexingNodeRegistrationVotePool struct {
	NodeAddress stratos.SdsAddress   `json:"node_address" yaml:"node_address"`
	ApproveList []stratos.SdsAddress `json:"approve_list" yaml:"approve_list"`
	RejectList  []stratos.SdsAddress `json:"reject_list" yaml:"reject_list"`
	ExpireTime  time.Time            `json:"expire_time" yaml:"expire_time"`
}

func NewRegistrationVotePool(nodeAddress stratos.SdsAddress, approveList []stratos.SdsAddress, rejectList []stratos.SdsAddress, expireTime time.Time) IndexingNodeRegistrationVotePool {
	return IndexingNodeRegistrationVotePool{
		NodeAddress: nodeAddress,
		ApproveList: approveList,
		RejectList:  rejectList,
		ExpireTime:  expireTime,
	}
}

func (v1 IndexingNode) Equal(v2 IndexingNode) bool {
	bz1 := goamino.MustMarshalBinaryLengthPrefixed(&v1)
	bz2 := goamino.MustMarshalBinaryLengthPrefixed(&v2)
	return bytes.Equal(bz1, bz2)
}
