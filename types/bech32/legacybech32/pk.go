// Deprecated: The module provides legacy bech32 functions which will be removed in a future
// release.
package legacybech32

import (
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/stratosnet/stratos-chain/types"
)

// TODO: when removing this package remove:
// + sdk:config.GetBech32AccountPubPrefix (and other related functions)
// + Bech32PrefixAccAddr and other related constants

// Deprecated: Bech32PubKeyType defines a string type alias for a Bech32 public key type.
type Bech32PubKeyType string

// Bech32 conversion constants
const (
	Bech32PubKeyTypeAccPub    Bech32PubKeyType = "accpub"
	Bech32PubKeyTypeValPub    Bech32PubKeyType = "valpub"
	Bech32PubKeyTypeConsPub   Bech32PubKeyType = "conspub"
	Bech32PubKeyTypeSdsP2PPub Bech32PubKeyType = "sdsp2p"
)

// Deprecated: MarshalPubKey returns a Bech32 encoded string containing the appropriate
// prefix based on the key type provided for a given PublicKey.
func MarshalPubKey(pkt Bech32PubKeyType, pubkey cryptotypes.PubKey) (string, error) {
	bech32Prefix := getPrefix(pkt)
	return bech32.ConvertAndEncode(bech32Prefix, legacy.Cdc.MustMarshal(pubkey))
}

// Deprecated: MustMarshalPubKey calls MarshalPubKey and panics on error.
func MustMarshalPubKey(pkt Bech32PubKeyType, pubkey cryptotypes.PubKey) string {
	res, err := MarshalPubKey(pkt, pubkey)
	if err != nil {
		panic(err)
	}

	return res
}

func getPrefix(pkt Bech32PubKeyType) string {
	cfg := types.GetConfig()
	switch pkt {
	case Bech32PubKeyTypeAccPub:
		return cfg.GetBech32AccountPubPrefix()
	case Bech32PubKeyTypeValPub:
		return cfg.GetBech32ValidatorPubPrefix()
	case Bech32PubKeyTypeConsPub:
		return cfg.GetBech32ConsensusPubPrefix()
	case Bech32PubKeyTypeSdsP2PPub:
		return cfg.GetBech32SdsNodeP2PPubPrefix()
	}

	return ""
}

// Deprecated: UnmarshalPubKey returns a PublicKey from a bech32-encoded PublicKey with
// a given key type.
func UnmarshalPubKey(pkt Bech32PubKeyType, pubkeyStr string) (cryptotypes.PubKey, error) {
	bech32Prefix := getPrefix(pkt)

	bz, err := sdk.GetFromBech32(pubkeyStr, bech32Prefix)
	if err != nil {
		return nil, err
	}
	return legacy.PubKeyFromBytes(bz)
}
