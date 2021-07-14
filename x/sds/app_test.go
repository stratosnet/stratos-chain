package sds

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/stratosnet/stratos-chain/x/pot"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/stratosnet/stratos-chain/x/sds/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"log"
	"testing"
)

var (
	paramSpecificMinedReward = sdk.NewInt(160000000000)
	paramSpecificEpoch       = sdk.NewInt(100)
)

func TestSdsMsgs(t *testing.T) {

	/********************* initialize mock app *********************/
	SetConfig()
	mApp, _, _, _, _ := getMockApp(t)
	accs := setupAccounts(mApp)
	mock.SetGenesis(mApp, accs)
	mock.CheckBalance(t, mApp, foundationAccAddr, foundationDeposit)

	///********************* create fileUpload msg *********************/
	log.Print("====== Testing MsgFileUpload ======")
	fileHash, _ := hex.DecodeString(testFileHashHex)
	fileUploadMsg := types.NewMsgUpload(fileHash, spNodeAddrIdx1, sdsAccAddr2)
	headerUpload := abci.Header{Height: mApp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, headerUpload, []sdk.Msg{fileUploadMsg}, []uint64{17}, []uint64{0}, true, true, spNodePrivKeyIdx1)
	coin := sdk.NewCoin(DefaultDenom, spNodeInitialStakeIdx1)
	mock.CheckBalance(t, mApp, spNodeAddrIdx1, sdk.Coins{coin})
	///********************* create prepay msg *********************/
	log.Print("====== Testing MsgPrepay ======")
	coinToPrepay := sdk.NewCoin(DefaultDenom, prepayAmt)
	prepayMsg := types.NewMsgPrepay(sdsAccAddr3, sdk.NewCoins(coinToPrepay))
	headerPrepay := abci.Header{Height: mApp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, headerPrepay, []sdk.Msg{prepayMsg}, []uint64{21}, []uint64{0}, true, true, sdsAccPrivKey3)
	newBalanceInt := sdsAccBal3.Sub(prepayAmt)
	newBalanceCoin := sdk.NewCoin(DefaultDenom, newBalanceInt)
	mock.CheckBalance(t, mApp, sdsAccAddr3, sdk.NewCoins(newBalanceCoin))
}

func getMockApp(t *testing.T) (*mock.App, Keeper, bank.Keeper, register.Keeper, pot.Keeper) {
	mApp := mock.NewApp()

	RegisterCodec(mApp.Cdc)
	supply.RegisterCodec(mApp.Cdc)
	staking.RegisterCodec(mApp.Cdc)
	register.RegisterCodec(mApp.Cdc)

	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyRegister := sdk.NewKVStoreKey(register.StoreKey)
	keyPot := sdk.NewKVStoreKey(pot.StoreKey)
	keySds := sdk.NewKVStoreKey(StoreKey)

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	bankKeeper := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     {"fee_collector"},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(mApp.Cdc, keyStaking, supplyKeeper, mApp.ParamsKeeper.Subspace(staking.DefaultParamspace))
	registerKeeper := register.NewKeeper(mApp.Cdc, keyRegister, mApp.AccountKeeper, bankKeeper, mApp.ParamsKeeper.Subspace(register.DefaultParamSpace))
	potKeeper := pot.NewKeeper(mApp.Cdc, keyPot, mApp.ParamsKeeper.Subspace(DefaultParamSpace), auth.FeeCollectorName, bankKeeper, supplyKeeper, mApp.AccountKeeper, stakingKeeper, registerKeeper)
	keeper := NewKeeper(mApp.Cdc, keySds, bankKeeper, registerKeeper, potKeeper)

	mApp.Router().AddRoute(staking.RouterKey, staking.NewHandler(stakingKeeper))
	mApp.Router().AddRoute(RouterKey, NewHandler(keeper))
	mApp.SetEndBlocker(getEndBlocker(keeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper, mApp.AccountKeeper, supplyKeeper,
		[]supplyexported.ModuleAccountI{feeCollector, notBondedPool, bondPool}, stakingKeeper, registerKeeper, potKeeper))

	err := mApp.CompleteSetup(keySds, keyStaking, keySupply, keyRegister, keyPot)
	require.NoError(t, err)

	return mApp, keeper, bankKeeper, registerKeeper, potKeeper
}

// getInitChainer initializes the chainer of the mock app and sets the genesis
// state. It returns an empty ResponseInitChain.
func getInitChainer(mapp *mock.App, keeper Keeper, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper,
	blacklistedAddrs []supplyexported.ModuleAccountI, stakingKeeper staking.Keeper, registerKeeper register.Keeper, potKeeper pot.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		mapp.InitChainer(ctx, req)

		lastResourceNodeTotalStake := initialStakeRes1.Add(initialStakeRes2).Add(initialStakeRes3).Add(initialStakeRes4).Add(initialStakeRes5)
		lastIndexingNodeTotalStake := initialStakeIdx1.Add(initialStakeIdx2).Add(initialStakeIdx3).Add(spNodeInitialStakeIdx1)

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
		lastIndexingNodeStakes = append(lastIndexingNodeStakes, register.LastIndexingNodeStake{Address: spNodeAddrIdx1, Stake: spNodeInitialStakeIdx1})

		resourceNodes := setupAllResourceNodes()
		indexingNodes := setupAllIndexingNodes()

		registerGenesis := register.NewGenesisState(register.DefaultParams(), lastResourceNodeTotalStake, lastResourceNodeStakes, resourceNodes,
			lastIndexingNodeTotalStake, lastIndexingNodeStakes, indexingNodes)

		register.InitGenesis(ctx, registerKeeper, registerGenesis)

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		stakingGenesis := staking.NewGenesisState(staking.NewParams(staking.DefaultUnbondingTime, staking.DefaultMaxValidators, staking.DefaultMaxEntries, 0, DefaultDenom), nil, nil)

		totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000)))
		supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		validators := staking.InitGenesis(ctx, stakingKeeper, accountKeeper, supplyKeeper, stakingGenesis)

		//preset
		registerKeeper.SetRemainingOzoneLimit(ctx, remainingOzoneLimit)
		potKeeper.SetTotalUnissuedPrepay(ctx, totalUnissuedPrepay)

		//pot genesis data load
		pot.InitGenesis(ctx, potKeeper, pot.NewGenesisState(pottypes.DefaultParams(), foundationAccAddr, initialOzonePrice))

		// init bank genesis
		keeper.BankKeeper.SetSendEnabled(ctx, true)

		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

// getEndBlocker returns a staking endblocker.
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		validatorUpdates := keeper.PotKeeper.StakingKeeper.BlockValidatorUpdates(ctx)

		return abci.ResponseEndBlock{
			ValidatorUpdates: validatorUpdates,
		}
	}
	return nil
}
