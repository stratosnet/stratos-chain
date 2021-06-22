package pot

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}

func Test(t *testing.T) {
	SetConfig()
	//mApp, k, accountKeeper, bankKeeper, stakingKeeper, registerKeeper := getMockApp(t)
	mApp, _, _ := getMockApp(t)

	resOwnderAcc1 := &auth.BaseAccount{
		Address: resOwner1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", initialStakeRes1)},
	}
	fmt.Println("resOwnderAcc1" + resOwnderAcc1.String())

	idxOwnerAcc1 := &auth.BaseAccount{
		Address: idxOwner1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", initialStakeIdx1)},
	}
	fmt.Println("idxOwnerAcc1" + idxOwnerAcc1.String())

	valOwnerAcc1 := &auth.BaseAccount{
		Address: valAccAddr1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", valInitialStake)},
	}
	fmt.Println("valOwnerAcc1" + valOwnerAcc1.String())

	addrIdxAcc1 := &auth.BaseAccount{
		Address: addrIdx1,
		Coins:   sdk.Coins{sdk.NewCoin("ustos", sdk.ZeroInt())},
	}
	fmt.Println("addrIdxAcc1" + addrIdxAcc1.String())

	accs := []authexported.Account{resOwnderAcc1, idxOwnerAcc1, addrIdxAcc1, valOwnerAcc1}

	mock.SetGenesis(mApp, accs)

	//var nodesVolume = make([]types.SingleNodeVolume, 0)

	volume1 := types.NewSingleNodeVolume(addrRes1, sdk.NewInt(10000000))
	volume2 := types.NewSingleNodeVolume(addrRes2, sdk.NewInt(10000000))
	volume3 := types.NewSingleNodeVolume(addrRes3, sdk.NewInt(10000000))

	nodesVolume := []types.SingleNodeVolume{volume1, volume2, volume3}
	reporter := addrIdx1
	epoch := sdk.NewInt(1)
	reportReference := "ref"

	volumeReportMsg := types.NewMsgVolumeReport(nodesVolume, reporter, epoch, reportReference)
	header := abci.Header{Height: mApp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{volumeReportMsg}, []uint64{0}, []uint64{0}, true, true, privKeyIdx1, idxOwnerPrivKey1)

}

//
//func TestPotVolumeReportMsgs(t *testing.T) {
//	mApp, k, accountKeeper, bankKeeper, stakingKeeper, registerKeeper := getMockApp(t)
//
//	// create validator with 50% commission
//	stakingHandler := staking.NewHandler(stakingKeeper)
//	//createAccount for validator's delegator
//	account := accountKeeper.GetAccount(ctx, valAccAddr1)
//	if account == nil {
//		account = accountKeeper.NewAccountWithAddress(ctx, valAccAddr1)
//		//fmt.Printf("create account: " + account.String() + "\n")
//	}
//
//	_, err := bankKeeper.AddCoins(ctx, valAccAddr1, sdk.NewCoins(valInitialStake))
//	require.NoError(t, err)
//
//
//	commission := staking.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1))
//	msgVal := staking.NewMsgCreateValidator(valOpAddr1, valConsPk1, valInitialStake, staking.Description{"foo_moniker", "", "", "", ""}, commission, sdk.OneInt())
//
//	header := abci.Header{Height: mApp.LastBlockHeight() + 1}
//	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{msgVal}, []uint64{0}, []uint64{0}, true, true, valConsPrivKey1)
//
//
//	res, err := stakingHandler(ctx, msgVal)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//	stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
//
//	//build traffic list
//	var trafficList []types.SingleNodeVolume
//	trafficList = append(trafficList, types.NewSingleNodeVolume(addrRes1, sdk.NewInt(resourceNodeVolume1)))
//	trafficList = append(trafficList, types.NewSingleNodeVolume(addrRes2, sdk.NewInt(resourceNodeVolume2)))
//	trafficList = append(trafficList, types.NewSingleNodeVolume(addrRes3, sdk.NewInt(resourceNodeVolume3)))
//
//	//check prepared data
//	S := registerKeeper.GetInitialGenesisStakeTotal(ctx).ToDec()
//	fmt.Println("S=" + S.String())
//	Pt := k.GetTotalUnissuedPrepay(ctx).ToDec()
//	fmt.Println("Pt=" + Pt.String())
//	Y := k.GetTotalConsumedOzone(trafficList).ToDec()
//	fmt.Println("Y=" + Y.String())
//	Lt := registerKeeper.GetRemainingOzoneLimit(ctx).ToDec()
//	fmt.Println("Lt=" + Lt.String())
//	R := S.Add(Pt).Mul(Y).Quo(Lt.Add(Y))
//	fmt.Println("R=" + R.String())
//
//	fmt.Println("***************************************************************************************")
//
//	//genTokens := sdk.TokensFromConsensusPower(42)
//	//bondTokens := sdk.TokensFromConsensusPower(10)
//	//genCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)
//	//bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, bondTokens)
//	//
//	//acc1 := &auth.BaseAccount{
//	//	Address: addr1,
//	//	Coins:   sdk.Coins{genCoin},
//	//}
//	//acc2 := &auth.BaseAccount{
//	//	Address: addr2,
//	//	Coins:   sdk.Coins{genCoin},
//	//}
//	//accs := []authexported.Account{acc1, acc2}
//	//
//	//mock.SetGenesis(mApp, accs)
//	//mock.CheckBalance(t, mApp, addr1, sdk.Coins{genCoin})
//	//mock.CheckBalance(t, mApp, addr2, sdk.Coins{genCoin})
//	//
//	//// create validator
//	//description := NewDescription("foo_moniker", "", "", "", "")
//	//createValidatorMsg := NewMsgCreateValidator(
//	//	sdk.ValAddress(addr1), priv1.PubKey(), bondCoin, description, commissionRates, sdk.OneInt(),
//	//)
//	//
//	//header := abci.Header{Height: mApp.LastBlockHeight() + 1}
//	//mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{createValidatorMsg}, []uint64{0}, []uint64{0}, true, true, priv1)
//	//mock.CheckBalance(t, mApp, addr1, sdk.Coins{genCoin.Sub(bondCoin)})
//	//
//	//header = abci.Header{Height: mApp.LastBlockHeight() + 1}
//	//mApp.BeginBlock(abci.RequestBeginBlock{Header: header})
//	//
//	//validator := checkValidator(t, mApp, keeper, sdk.ValAddress(addr1), true)
//	//require.Equal(t, sdk.ValAddress(addr1), validator.OperatorAddress)
//	//require.Equal(t, sdk.Bonded, validator.Status)
//	//require.True(sdk.IntEq(t, bondTokens, validator.BondedTokens()))
//	//
//	//header = abci.Header{Height: mApp.LastBlockHeight() + 1}
//	//mApp.BeginBlock(abci.RequestBeginBlock{Header: header})
//	//
//	//// edit the validator
//	//description = NewDescription("bar_moniker", "", "", "", "")
//	//editValidatorMsg := NewMsgEditValidator(sdk.ValAddress(addr1), description, nil, nil)
//	//
//	//header = abci.Header{Height: mApp.LastBlockHeight() + 1}
//	//mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{editValidatorMsg}, []uint64{0}, []uint64{1}, true, true, priv1)
//	//
//	//validator = checkValidator(t, mApp, keeper, sdk.ValAddress(addr1), true)
//	//require.Equal(t, description, validator.Description)
//	//
//	//// delegate
//	//mock.CheckBalance(t, mApp, addr2, sdk.Coins{genCoin})
//	//delegateMsg := NewMsgDelegate(addr2, sdk.ValAddress(addr1), bondCoin)
//	//
//	//header = abci.Header{Height: mApp.LastBlockHeight() + 1}
//	//mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{delegateMsg}, []uint64{1}, []uint64{0}, true, true, priv2)
//	//mock.CheckBalance(t, mApp, addr2, sdk.Coins{genCoin.Sub(bondCoin)})
//	//checkDelegation(t, mApp, keeper, addr2, sdk.ValAddress(addr1), true, bondTokens.ToDec())
//	//
//	//// begin unbonding
//	//beginUnbondingMsg := NewMsgUndelegate(addr2, sdk.ValAddress(addr1), bondCoin)
//	//header = abci.Header{Height: mApp.LastBlockHeight() + 1}
//	//mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{beginUnbondingMsg}, []uint64{1}, []uint64{1}, true, true, priv2)
//	//
//	//// delegation should exist anymore
//	//checkDelegation(t, mApp, keeper, addr2, sdk.ValAddress(addr1), false, sdk.Dec{})
//	//
//	//// balance should be the same because bonding not yet complete
//	//mock.CheckBalance(t, mApp, addr2, sdk.Coins{genCoin.Sub(bondCoin)})
//}

//func getMockApp(t *testing.T) (*mock.App, Keeper, auth.AccountKeeper, bank.Keeper, staking.Keeper, register.Keeper) {
//	mApp := mock.NewApp()
//
//	RegisterCodec(mApp.Cdc)
//	supply.RegisterCodec(mApp.Cdc)
//	register.RegisterCodec(mApp.Cdc)
//	staking.RegisterCodec(mApp.Cdc)
//
//	keyPot := sdk.NewKVStoreKey(StoreKey)
//	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
//	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
//	keyRegister := sdk.NewKVStoreKey(register.StoreKey)
//
//	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
//	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
//	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
//
//	blacklistedAddrs := make(map[string]bool)
//	blacklistedAddrs[feeCollector.GetAddress().String()] = true
//	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
//	blacklistedAddrs[bondPool.GetAddress().String()] = true
//
//	bankKeeper := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), blacklistedAddrs)
//	maccPerms := map[string][]string{
//		auth.FeeCollectorName:   nil,
//		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
//		staking.BondedPoolName:    {supply.Burner, supply.Staking},
//	}
//	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bankKeeper, maccPerms)
//	stakingKeeper := staking.NewKeeper(mApp.Cdc, keyStaking, supplyKeeper, mApp.ParamsKeeper.Subspace(staking.DefaultParamspace))
//	registerKeeper := register.NewKeeper(mApp.Cdc, keyRegister, mApp.AccountKeeper, bankKeeper, mApp.ParamsKeeper.Subspace(register.DefaultParamSpace))
//
//	keeper := NewKeeper(mApp.Cdc, keyPot, mApp.ParamsKeeper.Subspace(DefaultParamSpace), auth.FeeCollectorName, bankKeeper, supplyKeeper, mApp.AccountKeeper, stakingKeeper, registerKeeper)
//
//	mApp.Router().AddRoute(RouterKey, NewHandler(keeper))
//	mApp.SetEndBlocker(getEndBlocker(keeper))
//	mApp.SetInitChainer(getInitChainer(mApp, keeper, mApp.AccountKeeper, supplyKeeper, []supplyexported.ModuleAccountI{feeCollector, notBondedPool, bondPool}, stakingKeeper, registerKeeper))
//
//	require.NoError(t, mApp.CompleteSetup(keyStaking, keySupply))
//	return mApp, keeper, mApp.AccountKeeper, bankKeeper, stakingKeeper, registerKeeper
//}

func getMockApp(t *testing.T) (*mock.App, Keeper, auth.AccountKeeper) {
	mApp := mock.NewApp()

	RegisterCodec(mApp.Cdc)
	supply.RegisterCodec(mApp.Cdc)
	staking.RegisterCodec(mApp.Cdc)
	register.RegisterCodec(mApp.Cdc)

	//keyParams := sdk.NewKVStoreKey(params.StoreKey)
	//tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	//keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyRegister := sdk.NewKVStoreKey(register.StoreKey)
	keyPot := sdk.NewKVStoreKey(StoreKey)

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	//pk := params.NewKeeper(mApp.Cdc, keyParams, tkeyParams)
	//mApp.ParamsKeeper = pk

	bankKeeper := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     {"fee_collector"},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(mApp.Cdc, keyStaking, supplyKeeper, mApp.ParamsKeeper.Subspace(staking.DefaultParamspace))
	registerKeeper := register.NewKeeper(mApp.Cdc, keyRegister, mApp.AccountKeeper, bankKeeper, mApp.ParamsKeeper.Subspace(register.DefaultParamSpace))

	keeper := NewKeeper(mApp.Cdc, keyPot, mApp.ParamsKeeper.Subspace(DefaultParamSpace), auth.FeeCollectorName, bankKeeper, supplyKeeper, mApp.AccountKeeper, stakingKeeper, registerKeeper)

	//mApp.Router().AddRoute(staking.RouterKey, staking.NewHandler(stakingKeeper))
	//mApp.Router().AddRoute(register.RouterKey, register.NewHandler(registerKeeper))
	mApp.Router().AddRoute(RouterKey, NewHandler(keeper))
	mApp.SetEndBlocker(getEndBlocker(keeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper, mApp.AccountKeeper, supplyKeeper,
		[]supplyexported.ModuleAccountI{feeCollector, notBondedPool, bondPool}, stakingKeeper, registerKeeper))

	err := mApp.CompleteSetup(keyStaking, keySupply, keyRegister, keyPot)
	require.NoError(t, err)

	return mApp, keeper, mApp.AccountKeeper
}

// getEndBlocker returns a staking endblocker.
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	//return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	//	validatorUpdates := EndBlocker(ctx, keeper)
	//
	//	return abci.ResponseEndBlock{
	//		ValidatorUpdates: validatorUpdates,
	//	}
	//}
	return nil
}

// getInitChainer initializes the chainer of the mock app and sets the genesis
// state. It returns an empty ResponseInitChain.
func getInitChainer(mapp *mock.App, keeper Keeper, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper,
	blacklistedAddrs []supplyexported.ModuleAccountI, stakingKeeper staking.Keeper, registerKeeper register.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		mapp.InitChainer(ctx, req)

		//registerGenesis := register.DefaultGenesisState()
		lastResourceNodeTotalStake := initialStakeRes1.Add(initialStakeRes2).Add(initialStakeRes3).Add(initialStakeRes4).Add(initialStakeRes5)
		lastIndexingNodeTotalStake := initialStakeIdx1.Add(initialStakeIdx2).Add(initialStakeIdx3)

		var lastResourceNodeStakes []register.LastResourceNodeStake
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes1, Stake: initialStakeRes1})
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes2, Stake: initialStakeRes2})
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes3, Stake: initialStakeRes3})
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes4, Stake: initialStakeRes4})
		lastResourceNodeStakes = append(lastResourceNodeStakes, register.LastResourceNodeStake{Address: addrRes5, Stake: initialStakeRes5})

		var lastIndexingNodeStakes []register.LastIndexingNodeStake
		lastIndexingNodeStakes = append(lastIndexingNodeStakes, register.LastIndexingNodeStake{Address: addrIdx1, Stake: initialStakeIdx1})
		lastIndexingNodeStakes = append(lastIndexingNodeStakes, register.LastIndexingNodeStake{Address: addrIdx2, Stake: initialStakeIdx2})
		lastIndexingNodeStakes = append(lastIndexingNodeStakes, register.LastIndexingNodeStake{Address: addrIdx3, Stake: initialStakeIdx3})

		resourceNodes := setupAllResourceNodes()
		indexingNodes := setupAllIndexingNodes()

		registerGenesis := register.NewGenesisState(register.DefaultParams(), lastResourceNodeTotalStake, lastResourceNodeStakes, resourceNodes,
			lastIndexingNodeTotalStake, lastIndexingNodeStakes, indexingNodes)
		register.InitGenesis(ctx, registerKeeper, registerGenesis)

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		stakingGenesis := staking.NewGenesisState(staking.NewParams(staking.DefaultUnbondingTime, staking.DefaultMaxValidators, staking.DefaultMaxEntries, 0, "ustos"), nil, nil)

		totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000)))
		supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		validators := staking.InitGenesis(ctx, stakingKeeper, accountKeeper, supplyKeeper, stakingGenesis)

		InitGenesis(ctx, keeper, types.DefaultGenesisState())

		return abci.ResponseInitChain{
			Validators: validators,
		}
		//validators := staking.InitGenesis(ctx, stakingKeeper, accountKeeper, supplyKeeper, stakingGenesis)
		//return abci.ResponseInitChain{
		//	Validators: validators,
		//}

	}

	return nil
}

func setupAllResourceNodes() []register.ResourceNode {
	resourceNode1 := register.NewResourceNode("sds://resourceNode1", pubKeyRes1, resOwner1, register.NewDescription("sds://resourceNode1", "", "", "", ""), "4")
	resourceNode2 := register.NewResourceNode("sds://resourceNode2", pubKeyRes2, resOwner2, register.NewDescription("sds://resourceNode2", "", "", "", ""), "4")
	resourceNode3 := register.NewResourceNode("sds://resourceNode3", pubKeyRes3, resOwner3, register.NewDescription("sds://resourceNode3", "", "", "", ""), "4")
	resourceNode4 := register.NewResourceNode("sds://resourceNode4", pubKeyRes4, resOwner4, register.NewDescription("sds://resourceNode4", "", "", "", ""), "4")
	resourceNode5 := register.NewResourceNode("sds://resourceNode5", pubKeyRes5, resOwner5, register.NewDescription("sds://resourceNode5", "", "", "", ""), "4")

	resourceNode1.AddToken(initialStakeRes1)
	resourceNode2.AddToken(initialStakeRes2)
	resourceNode3.AddToken(initialStakeRes3)
	resourceNode4.AddToken(initialStakeRes4)
	resourceNode5.AddToken(initialStakeRes5)

	var resourceNodes []register.ResourceNode
	resourceNodes = append(resourceNodes, resourceNode1)
	resourceNodes = append(resourceNodes, resourceNode2)
	resourceNodes = append(resourceNodes, resourceNode3)
	resourceNodes = append(resourceNodes, resourceNode4)
	resourceNodes = append(resourceNodes, resourceNode5)
	return resourceNodes
}

func setupAllIndexingNodes() []register.IndexingNode {
	var indexingNodes []register.IndexingNode
	indexingNode1 := register.NewIndexingNode("sds://indexingNode1", pubKeyIdx1, idxOwner1, register.NewDescription("sds://indexingNode1", "", "", "", ""))
	indexingNode2 := register.NewIndexingNode("sds://indexingNode2", pubKeyIdx2, idxOwner2, register.NewDescription("sds://indexingNode2", "", "", "", ""))
	indexingNode3 := register.NewIndexingNode("sds://indexingNode3", pubKeyIdx3, idxOwner3, register.NewDescription("sds://indexingNode3", "", "", "", ""))

	indexingNode1.AddToken(initialStakeIdx1)
	indexingNode2.AddToken(initialStakeIdx2)
	indexingNode3.AddToken(initialStakeIdx3)

	indexingNode1.Status = sdk.Bonded
	indexingNode2.Status = sdk.Bonded
	indexingNode3.Status = sdk.Bonded

	indexingNodes = append(indexingNodes, indexingNode1)
	indexingNodes = append(indexingNodes, indexingNode2)
	indexingNodes = append(indexingNodes, indexingNode3)

	return indexingNodes

}
