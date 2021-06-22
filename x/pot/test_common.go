package pot

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	AccountAddressPrefix = "st"

	stos2ustos = 1000000000
	oz2uoz     = 1000000000

	resourceNodeVolume1 = 500 * oz2uoz
	resourceNodeVolume2 = 300 * oz2uoz
	resourceNodeVolume3 = 200 * oz2uoz
	totalVolume         = resourceNodeVolume1 + resourceNodeVolume2 + resourceNodeVolume3
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"

	foundationAcc     = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	foundationDeposit = sdk.NewCoins(sdk.NewCoin("ustos", sdk.NewInt(40000000*stos2ustos)))

	resOwnerPrivKey1 = ed25519.GenPrivKey()
	resOwnerPrivKey2 = ed25519.GenPrivKey()
	resOwnerPrivKey3 = ed25519.GenPrivKey()
	resOwnerPrivKey4 = ed25519.GenPrivKey()
	resOwnerPrivKey5 = ed25519.GenPrivKey()
	idxOwnerPrivKey1 = ed25519.GenPrivKey()
	idxOwnerPrivKey2 = ed25519.GenPrivKey()
	idxOwnerPrivKey3 = ed25519.GenPrivKey()

	resOwner1 = sdk.AccAddress(resOwnerPrivKey1.PubKey().Address())
	resOwner2 = sdk.AccAddress(resOwnerPrivKey2.PubKey().Address())
	resOwner3 = sdk.AccAddress(resOwnerPrivKey3.PubKey().Address())
	resOwner4 = sdk.AccAddress(resOwnerPrivKey4.PubKey().Address())
	resOwner5 = sdk.AccAddress(resOwnerPrivKey5.PubKey().Address())
	idxOwner1 = sdk.AccAddress(idxOwnerPrivKey1.PubKey().Address())
	idxOwner2 = sdk.AccAddress(idxOwnerPrivKey2.PubKey().Address())
	idxOwner3 = sdk.AccAddress(idxOwnerPrivKey3.PubKey().Address())

	pubKeyRes1       = ed25519.GenPrivKey().PubKey()
	addrRes1         = sdk.AccAddress(pubKeyRes1.Address())
	initialStakeRes1 = sdk.NewInt(3 * stos2ustos)

	pubKeyRes2       = ed25519.GenPrivKey().PubKey()
	addrRes2         = sdk.AccAddress(pubKeyRes2.Address())
	initialStakeRes2 = sdk.NewInt(3 * stos2ustos)

	pubKeyRes3       = ed25519.GenPrivKey().PubKey()
	addrRes3         = sdk.AccAddress(pubKeyRes3.Address())
	initialStakeRes3 = sdk.NewInt(3 * stos2ustos)

	pubKeyRes4       = ed25519.GenPrivKey().PubKey()
	addrRes4         = sdk.AccAddress(pubKeyRes4.Address())
	initialStakeRes4 = sdk.NewInt(3 * stos2ustos)

	pubKeyRes5       = ed25519.GenPrivKey().PubKey()
	addrRes5         = sdk.AccAddress(pubKeyRes5.Address())
	initialStakeRes5 = sdk.NewInt(3 * stos2ustos)

	privKeyIdx1      = ed25519.GenPrivKey()
	pubKeyIdx1       = privKeyIdx1.PubKey()
	addrIdx1         = sdk.AccAddress(pubKeyIdx1.Address())
	initialStakeIdx1 = sdk.NewInt(5 * stos2ustos)

	pubKeyIdx2       = ed25519.GenPrivKey().PubKey()
	addrIdx2         = sdk.AccAddress(pubKeyIdx2.Address())
	initialStakeIdx2 = sdk.NewInt(5 * stos2ustos)

	pubKeyIdx3       = ed25519.GenPrivKey().PubKey()
	addrIdx3         = sdk.AccAddress(pubKeyIdx3.Address())
	initialStakeIdx3 = sdk.NewInt(5 * stos2ustos)

	valOpPk1        = ed25519.GenPrivKey().PubKey()
	valOpAddr1      = sdk.ValAddress(valOpPk1.Address())
	valAccAddr1     = sdk.AccAddress(valOpPk1.Address())
	valConsPrivKey1 = ed25519.GenPrivKey()
	valConsPk1      = valConsPrivKey1.PubKey()
	valInitialStake = sdk.NewInt(15 * stos2ustos)

	totalUnissuedPrePay = sdk.NewInt(5000 * stos2ustos)
	remainingOzoneLimit = sdk.NewInt(5000 * oz2uoz)

	epoch1    = sdk.NewInt(1)
	epoch2017 = epoch1.Add(sdk.NewInt(2016))
	epoch4033 = epoch2017.Add(sdk.NewInt(2016))
)
