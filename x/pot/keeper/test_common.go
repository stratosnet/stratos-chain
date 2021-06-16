package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

const (
	AccountAddressPrefix = "st"
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
)

func CreateTestInput(t *testing.T, isCheckTx bool) (
	sdk.Context, auth.AccountKeeper, bank.Keeper, Keeper, staking.Keeper, params.Keeper, supply.Keeper, register.Keeper) {

	SetConfig()

	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyRegister := sdk.NewKVStoreKey(register.StoreKey)
	keyPot := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyRegister, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyPot, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	cdc := MakeTestCodec()
	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())

	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(cdc, keyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace))
	StakingParam := staking.NewParams(staking.DefaultUnbondingTime, staking.DefaultMaxValidators, staking.DefaultMaxEntries, 0, "ustos")
	stakingKeeper.SetParams(ctx, StakingParam)
	registerKeeper := register.NewKeeper(cdc, keyRegister, accountKeeper, bankKeeper, pk.Subspace(register.DefaultParamSpace))
	registerKeeper.SetParams(ctx, register.DefaultParams())

	keeper := NewKeeper(cdc, keyPot, pk.Subspace(types.DefaultParamSpace), auth.FeeCollectorName, bankKeeper, supplyKeeper, accountKeeper, stakingKeeper, registerKeeper)
	keeper.SetParams(ctx, types.DefaultParams())

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)

	initCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultBondDenom, sdk.NewInt(1000000000000)))

	foundationAcc, err := sdk.AccAddressFromBech32("st1qr9set2jaayzjjpm9tw4f3n6f5zfu3hef8wtaw")
	if err != nil {
		panic(err)
	}
	bankKeeper.AddCoins(ctx, foundationAcc, initCoins)

	return ctx, accountKeeper, bankKeeper, keeper, stakingKeeper, pk, supplyKeeper, registerKeeper
}

func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(register.MsgCreateResourceNode{}, "register/MsgCreateResourceNode", nil)
	cdc.RegisterConcrete(register.MsgCreateIndexingNode{}, "register/MsgCreateIndexingNode", nil)
	//cdc.RegisterConcrete(bank.MsgSend{}, "test/staking/Send", nil)
	//cdc.RegisterConcrete(types.MsgCreateValidator{}, "test/staking/CreateValidator", nil)
	//cdc.RegisterConcrete(types.MsgEditValidator{}, "test/staking/EditValidator", nil)
	//cdc.RegisterConcrete(types.MsgUndelegate{}, "test/staking/Undelegate", nil)
	//cdc.RegisterConcrete(types.MsgBeginRedelegate{}, "test/staking/BeginRedelegate", nil)

	// Register AppAccount
	cdc.RegisterInterface((*authexported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/pot/BaseAccount", nil)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}
