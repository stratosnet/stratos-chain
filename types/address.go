package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// Bech32 conversion constants
const (
	StratosBech32Prefix = "st"

	// PrefixSds is the prefix for sds keys
	PrefixSds = "sds"

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

var _ sdk.Address = SdsAddress{}
var _ yaml.Marshaler = SdsAddress{}

type SdsAddress []byte

// SdsPubKeyFromBech32 returns a SdsPublicKey from a Bech32 string.
func SdsPubKeyFromBech32(pubkeyStr string) (cryptotypes.PubKey, error) {
	bech32Prefix := GetConfig().GetBech32SdsNodeP2PPubPrefix()

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

// SdsAddressFromHex creates an SdsAddress from a hex string.
func SdsAddressFromHex(address string) (addr SdsAddress, err error) {
	bz, err := addressBytesFromHexString(address)
	return SdsAddress(bz), err
}

// AccAddressFromBech32 creates an SdsAddress from a Bech32 string.
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

// Returns boolean for whether two SdsAddress are Equal
func (aa SdsAddress) Equals(aa2 sdk.Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Returns boolean for whether a SdsAddress is empty
func (aa SdsAddress) Empty() bool {
	return aa == nil || len(aa) == 0
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa SdsAddress) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *SdsAddress) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa SdsAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa SdsAddress) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *SdsAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)

	if err != nil {
		return err
	}
	if s == "" {
		*aa = SdsAddress{}
		return nil
	}

	aa2, err := SdsAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *SdsAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	if s == "" {
		*aa = SdsAddress{}
		return nil
	}

	aa2, err := SdsAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
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
