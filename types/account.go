package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v2"
)

var _ exported.Account = (*StAccount)(nil)
var _ exported.GenesisAccount = (*StAccount)(nil)

func init() {
	authtypes.RegisterAccountTypeCodec(&StAccount{}, StAccountName)
}

// ----------------------------------------------------------------------------
// Main Stratos account
// ----------------------------------------------------------------------------

// StAccount implements the auth.Account interface and embeds an
// auth.BaseAccount type. It is compatible with the auth.AccountKeeper.
type StAccount struct {
	*authtypes.BaseAccount `json:"base_account" yaml:"base_account"`
	CodeHash               []byte `json:"code_hash" yaml:"code_hash"`
}

// ProtoAccount defines the prototype function for BaseAccount used for an
// AccountKeeper.
func ProtoAccount() exported.Account {
	return &StAccount{
		BaseAccount: &auth.BaseAccount{},
		CodeHash:    ethcrypto.Keccak256(nil),
	}
}

// EthAddress returns the account address ethereum format.
func (acc StAccount) EthAddress() ethcmn.Address {
	return ethcmn.BytesToAddress(acc.Address.Bytes())
}

// TODO: remove on SDK v0.40

// Balance returns the balance of an account.
func (acc StAccount) Balance(denom string) sdk.Int {
	return acc.GetCoins().AmountOf(denom)
}

// SetBalance sets an account's balance of the given coin denomination.
//
// CONTRACT: assumes the denomination is valid.
func (acc *StAccount) SetBalance(denom string, amt sdk.Int) {
	coins := acc.GetCoins()
	diff := amt.Sub(coins.AmountOf(denom))
	switch {
	case diff.IsPositive():
		// Increase coins to amount
		coins = coins.Add(sdk.NewCoin(denom, diff))
	case diff.IsNegative():
		// Decrease coins to amount
		coins = coins.Sub(sdk.NewCoins(sdk.NewCoin(denom, diff.Neg())))
	default:
		return
	}

	if err := acc.SetCoins(coins); err != nil {
		panic(fmt.Errorf("could not set %s coins for address %s: %w", denom, acc.EthAddress().String(), err))
	}
}

type stratosAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	EthAddress    string         `json:"eth_address" yaml:"eth_address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	CodeHash      string         `json:"code_hash" yaml:"code_hash"`
}

// MarshalYAML returns the YAML representation of an account.
func (acc StAccount) MarshalYAML() (interface{}, error) {
	alias := stratosAccountPretty{
		Address:       acc.Address,
		EthAddress:    acc.EthAddress().String(),
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      ethcmn.Bytes2Hex(acc.CodeHash),
	}

	var err error

	if acc.PubKey != nil {
		alias.PubKey, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}
	}

	bz, err := yaml.Marshal(alias)
	if err != nil {
		return nil, err
	}

	return string(bz), err
}

// MarshalJSON returns the JSON representation of an StAccount.
func (acc StAccount) MarshalJSON() ([]byte, error) {
	var ethAddress = ""

	if acc.BaseAccount != nil && acc.Address != nil {
		ethAddress = acc.EthAddress().String()
	}

	alias := stratosAccountPretty{
		Address:       acc.Address,
		EthAddress:    ethAddress,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      ethcmn.Bytes2Hex(acc.CodeHash),
	}

	var err error

	if acc.PubKey != nil {
		alias.PubKey, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into an StAccount.
func (acc *StAccount) UnmarshalJSON(bz []byte) error {
	var (
		alias stratosAccountPretty
		err   error
	)

	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	switch {
	case !alias.Address.Empty() && alias.EthAddress != "":
		// Both addresses provided. Verify correctness
		ethAddress := ethcmn.HexToAddress(alias.EthAddress)
		ethAddressFromAccAddress := ethcmn.BytesToAddress(alias.Address.Bytes())

		if !bytes.Equal(ethAddress.Bytes(), alias.Address.Bytes()) {
			err = sdkerrors.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"expected %s, got %s",
				ethAddressFromAccAddress.String(), ethAddress.String(),
			)
		}

	case !alias.Address.Empty() && alias.EthAddress == "":
		// unmarshal sdk.AccAddress only. Do nothing here
	case alias.Address.Empty() && alias.EthAddress != "":
		// retrieve sdk.AccAddress from ethereum address
		ethAddress := ethcmn.HexToAddress(alias.EthAddress)
		alias.Address = sdk.AccAddress(ethAddress.Bytes())
	case alias.Address.Empty() && alias.EthAddress == "":
		err = sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"account must contain address in Ethereum Hex or Cosmos Bech32 format",
		)
	}

	if err != nil {
		return err
	}

	acc.BaseAccount = &authtypes.BaseAccount{
		Coins:         alias.Coins,
		Address:       alias.Address,
		AccountNumber: alias.AccountNumber,
		Sequence:      alias.Sequence,
	}
	acc.CodeHash = ethcmn.Hex2Bytes(alias.CodeHash)

	if alias.PubKey != "" {
		acc.BaseAccount.PubKey, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, alias.PubKey)
		if err != nil {
			return err
		}
	}
	return nil
}

// String implements the fmt.Stringer interface
func (acc StAccount) String() string {
	out, _ := yaml.Marshal(acc)
	return string(out)
}