package keeper

//import (
//	"testing"
//
//	"github.com/cosmos/cosmos-sdk/codec"
//	"github.com/cosmos/cosmos-sdk/store"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/x/auth"
//	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
//	"github.com/cosmos/cosmos-sdk/x/bank"
//	"github.com/cosmos/cosmos-sdk/x/params"
//	stratos "github.com/stratosnet/stratos-chain/types"
//	"github.com/stratosnet/stratos-chain/x/register/types"
//	"github.com/stretchr/testify/require"
//	abci "github.com/tendermint/tendermint/abci/types"
//	"github.com/tendermint/tendermint/libs/log"
//	dbm "github.com/tendermint/tm-db"
//)
//
//func TestMain(m *testing.M) {
//	config := stratos.GetConfig()
//
//	config.Seal()
//}
//
//func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, auth.AccountKeeper, bank.Keeper, Keeper, params.Keeper) {
//
//	keyParams := sdk.NewKVStoreKey(params.StoreKey)
//	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
//	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
//	keyRegister := sdk.NewKVStoreKey(types.StoreKey)
//
//	db := dbm.NewMemDB()
//	ms := store.NewCommitMultiStore(db)
//
//	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
//	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(keyRegister, sdk.StoreTypeIAVL, db)
//	err := ms.LoadLatestVersion()
//	require.Nil(t, err)
//
//	cdc := MakeTestCodec()
//	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
//	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
//
//	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
//	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), nil)
//
//	keeper := NewKeeper(cdc, keyRegister, pk.Subspace(types.DefaultParamSpace), accountKeeper, bankKeeper)
//	keeper.SetParams(ctx, types.DefaultParams())
//
//	return ctx, accountKeeper, bankKeeper, keeper, pk
//}
//
//// create a codec used only for testing
//func MakeTestCodec() *codec.Codec {
//	var cdc = codec.New()
//
//	// Register AppAccount
//	cdc.RegisterInterface((*authexported.Account)(nil), nil)
//	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/pot/BaseAccount", nil)
//	codec.RegisterCrypto(cdc)
//
//	return cdc
//}
