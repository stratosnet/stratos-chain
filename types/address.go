package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
)

// Bech32PubKeyType defines a string type alias for a Bech32 public key type.
type Bech32PubKeyType string

const (
	StratosBech32Prefix = "st"

	// PrefixSds is the prefix for sds keys
	PrefixSds = "sds"

	Bech32PubKeyTypeAccPub    Bech32PubKeyType = "accpub"
	Bech32PubKeyTypeValPub    Bech32PubKeyType = "valpub"
	Bech32PubKeyTypeConsPub   Bech32PubKeyType = "conspub"
	Bech32PubKeyTypeSdsP2PPub Bech32PubKeyType = "sdsp2p"

	// AccountAddressPrefix defines the Bech32 prefix of an account's address (st)
	AccountAddressPrefix = StratosBech32Prefix
	// AccountPubKeyPrefix defines the Bech32 prefix of an account's public key (stpub)
	AccountPubKeyPrefix = StratosBech32Prefix + sdk.PrefixPublic
	// ValidatorAddressPrefix defines the Bech32 prefix of a validator's operator address (stvaloper)
	ValidatorAddressPrefix = StratosBech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// ValidatorPubKeyPrefix defines the Bech32 prefix of a validator's operator public key (stvaloperpub)
	ValidatorPubKeyPrefix = StratosBech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// ConsNodeAddressPrefix defines the Bech32 prefix of a consensus node address (stvalcons)
	ConsNodeAddressPrefix = StratosBech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// ConsNodePubKeyPrefix defines the Bech32 prefix of a consensus node public key (stvalconspub)
	ConsNodePubKeyPrefix = StratosBech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
	// SdsNodeP2PPubkeyPrefix defines the Bech32 prefix of an sds account's public key (stsdspub)
	SdsNodeP2PPubkeyPrefix = StratosBech32Prefix + PrefixSds + sdk.PrefixPublic
	// SdsNodeP2PAddressPrefix defines the Bech32 prefix of an sds account's address (stsds)
	SdsNodeP2PAddressPrefix = StratosBech32Prefix + PrefixSds
)

// GetPubKeyFromBech32 returns a PublicKey from a bech32-encoded PublicKey with
// a given key type.
func GetPubKeyFromBech32(pkt Bech32PubKeyType, pubkeyStr string) (cryptotypes.PubKey, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = GetConfig().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = GetConfig().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = GetConfig().GetBech32ConsensusPubPrefix()

	case Bech32PubKeyTypeSdsP2PPub:
		bech32Prefix = GetConfig().GetBech32SdsNodeP2PPubPrefix()
	}

	bz, err := sdk.GetFromBech32(pubkeyStr, bech32Prefix)
	if err != nil {
		return nil, err
	}

	pk, err := legacy.PubKeyFromBytes(bz)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

var _ sdk.Address = SdsAddress{}
var _ yaml.Marshaler = SdsAddress{}

type SdsAddress []byte

var _ sdk.Address = SdsAddress{}

// SdsAddressFromHex creates an SdsAddress from a hex string.
func SdsAddressFromHex(address string) (addr SdsAddress, err error) {
	bz, err := addressBytesFromHexString(address)
	return bz, err
}

// SdsAddressFromBech32 creates an SdsAddress from a Bech32 string.
func SdsAddressFromBech32(address string) (addr SdsAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return SdsAddress{}, errors.New("empty address string is not allowed")
	}

	bech32PrefixSdsAddr := GetConfig().GetBech32SdsNodeP2PAddrPrefix()

	bz, err := sdk.GetFromBech32(address, bech32PrefixSdsAddr)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return SdsAddress(bz), nil
}

// Equals Returns boolean for whether two SdsAddress are Equal
func (a SdsAddress) Equals(addr sdk.Address) bool {
	if a.Empty() && addr.Empty() {
		return true
	}

	return bytes.Equal(a.Bytes(), addr.Bytes())
}

func (a SdsAddress) Empty() bool {
	if a == nil || len(a) == 0 {
		return true
	}

	aa2 := SdsAddress{}
	return bytes.Equal(a.Bytes(), aa2.Bytes())
}

func (a SdsAddress) Marshal() ([]byte, error) {
	return a, nil
}

func (a SdsAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (a *SdsAddress) Unmarshal(data []byte) error {
	*a = data
	return nil
}

// MarshalYAML marshals to YAML using Bech32.
func (a SdsAddress) MarshalYAML() (interface{}, error) {
	return a.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (a *SdsAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)

	if err != nil {
		return err
	}
	if s == "" {
		*a = SdsAddress{}
		return nil
	}

	aa2, err := SdsAddressFromBech32(s)
	if err != nil {
		return err
	}

	*a = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (a *SdsAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	if s == "" {
		*a = SdsAddress{}
		return nil
	}

	aa2, err := SdsAddressFromBech32(s)
	if err != nil {
		return err
	}

	*a = aa2
	return nil
}

// Bytes returns the raw address bytes.
func (aa SdsAddress) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa SdsAddress) String() string {
	if aa.Empty() {
		return ""
	}

	bech32PrefixSdsAddr := GetConfig().GetBech32SdsNodeP2PAddrPrefix()

	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixSdsAddr, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
func (aa SdsAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}

func addressBytesFromHexString(address string) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("decoding Bech32 address failed: must provide an address")
	}

	return hex.DecodeString(address)
}
