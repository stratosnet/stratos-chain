package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/bech32"
)

// Bech32PubKeyType defines a string type alias for a Bech32 public key type.
type Bech32PubKeyType string

// Bech32 conversion constants
const (
	Bech32PubKeyTypeAccPub    Bech32PubKeyType = "accpub"
	Bech32PubKeyTypeValPub    Bech32PubKeyType = "valpub"
	Bech32PubKeyTypeConsPub   Bech32PubKeyType = "conspub"
	Bech32PubKeyTypeSdsP2PPub Bech32PubKeyType = "sdsp2p"
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
