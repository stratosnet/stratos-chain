package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/bech32"
	"gopkg.in/yaml.v2"
)

// Bech32PubKeyType defines a string type alias for a Bech32 public key type.
type Bech32PubKeyType string

// Bech32 conversion constants
const (
	StratosBech32Prefix = "st"

	Bech32PubKeyTypeAccPub    Bech32PubKeyType = "accpub"
	Bech32PubKeyTypeValPub    Bech32PubKeyType = "valpub"
	Bech32PubKeyTypeConsPub   Bech32PubKeyType = "conspub"
	Bech32PubKeyTypeSdsP2PPub Bech32PubKeyType = "sdsp2p"

	AccountPubKeyPrefix     = StratosBech32Prefix + "pub"
	ValidatorAddressPrefix  = StratosBech32Prefix + "valoper"
	ValidatorPubKeyPrefix   = StratosBech32Prefix + "valoperpub"
	ConsNodeAddressPrefix   = StratosBech32Prefix + "valcons"
	ConsNodePubKeyPrefix    = StratosBech32Prefix + "valconspub"
	SdsNodeP2PPubkeyPrefix  = StratosBech32Prefix + "sdspub"
	SdsNodeP2PAddressPrefix = StratosBech32Prefix + "sds"

	CoinType = 606

	HDPath = "m/44'/606'/0'/0/0"
)

// Bech32ifyPubKey returns a Bech32 encoded string containing the appropriate
// prefix based on the key type provided for a given PublicKey.
func Bech32ifyPubKey(pkt Bech32PubKeyType, pubkey crypto.PubKey) (string, error) {
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

	return bech32.ConvertAndEncode(bech32Prefix, pubkey.Bytes())
}

// GetPubKeyFromBech32 returns a PublicKey from a bech32-encoded PublicKey with
// a given key type.
func GetPubKeyFromBech32(pkt Bech32PubKeyType, pubkeyStr string) (crypto.PubKey, error) {
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

	pk, err := tmamino.PubKeyFromBytes(bz)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

type SdsAddress []byte

func (a SdsAddress) Equals(addr sdk.Address) bool {
	if a.Empty() && addr.Empty() {
		return true
	}

	return bytes.Equal(a.Bytes(), addr.Bytes())
}

func (a SdsAddress) Empty() bool {
	if a == nil {
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

func (a SdsAddress) Bytes() []byte {
	return a
}

func (a SdsAddress) String() string {
	if a.Empty() {
		return ""
	}

	bech32PrefixSdsAddr := GetConfig().GetBech32SdsNodeP2PAddrPrefix()

	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixSdsAddr, a.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

func (a SdsAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(a.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", a)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(a))))
	}
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

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func SdsAddressFromBech32(address string) (addr SdsAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return SdsAddress{}, nil
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
