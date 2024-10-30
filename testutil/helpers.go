package testutil

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/simapp"
	simappparams "cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/std"
	simstestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratosapp "github.com/stratosnet/stratos-chain/app"
	stratos "github.com/stratosnet/stratos-chain/types"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

// MakeTestEncodingConfig creates an EncodingConfig for testing. This function
// should be used only in tests or when creating a new app instance (NewApp*()).
// App user shouldn't create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeTestEncodingConfig() simappparams.EncodingConfig {
	encodingConfig := simappparams.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	stratosapp.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	stratosapp.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// SetupWithGenesisNodeSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisNodeSet(t *testing.T,
	valSet *tmtypes.ValidatorSet,
	metaNodes []registertypes.MetaNode,
	resourceNodes []registertypes.ResourceNode,
	genAccs []authtypes.GenesisAccount,
	chainId string,
	freshStart bool,
	balances ...banktypes.Balance) *stratosapp.StratosApp {

	app, genesisState := setup(chainId, true, 5)
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	validatorBondAmt := sdkmath.NewInt(1000000)

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress: sdk.ValAddress(val.Address).String(),
			ConsensusPubkey: pkAny,
			Jailed:          false,
			Status:          stakingtypes.Bonded,
			Tokens:          validatorBondAmt,
			DelegatorShares: sdkmath.LegacyOneDec(),
			Description:     stakingtypes.Description{},
			UnbondingHeight: int64(0),
			UnbondingTime:   time.Unix(0, 0).UTC(),
			// 50% commission
			Commission: stakingtypes.NewCommission(
				sdkmath.LegacyNewDecWithPrec(5, 1),
				sdkmath.LegacyNewDecWithPrec(5, 1),
				sdkmath.LegacyZeroDec()),
			MinSelfDelegation: sdkmath.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdkmath.LegacyOneDec()))
	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(
		stakingtypes.NewParams(
			stakingtypes.DefaultUnbondingTime,
			stakingtypes.DefaultMaxValidators,
			stakingtypes.DefaultMaxEntries,
			stakingtypes.DefaultHistoricalEntries,
			stratos.Wei,
			stakingtypes.DefaultMinCommissionRate),
		validators,
		delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(stratos.Wei, validatorBondAmt)},
	})

	// add bonded amount of resource nodes to module account
	resNodeBondedAmt := sdkmath.ZeroInt()
	for _, resNode := range resourceNodes {
		resNodeBondedAmt = resNodeBondedAmt.Add(resNode.Tokens)
	}
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(registertypes.ResourceNodeBondedPool).String(),
		Coins:   sdk.Coins{sdk.NewCoin(stratos.Wei, resNodeBondedAmt)},
	})

	// add bonded amount of meta nodes to module account
	metaNodeBondedAmt := sdkmath.ZeroInt()
	for _, metaNode := range metaNodes {
		metaNodeBondedAmt = metaNodeBondedAmt.Add(metaNode.Tokens)
	}
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(registertypes.MetaNodeBondedPool).String(),
		Coins:   sdk.Coins{sdk.NewCoin(stratos.Wei, metaNodeBondedAmt)},
	})

	initRemainingOzoneLimit := sdkmath.ZeroInt()
	if !freshStart {
		initRemainingOzoneLimit = resNodeBondedAmt.
			Add(metaNodeBondedAmt).
			ToLegacyDec().
			Quo(registertypes.DefaultDepositNozRate).
			TruncateInt()
	}

	// update total supply
	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	registerGenesis := registertypes.NewGenesisState(
		registertypes.DefaultParams(),
		resourceNodes,
		metaNodes,
		initRemainingOzoneLimit,
		make([]registertypes.Slashing, 0),
		registertypes.DefaultDepositNozRate,
		make([]registertypes.MetaNodeRegistrationVotePool, 0),
		make([]registertypes.UnbondingNode, 0),
		make([]registertypes.KickMetaNodeVotePool, 0),
	)
	genesisState[registertypes.ModuleName] = app.AppCodec().MustMarshalJSON(registerGenesis)

	potGenesis := pottypes.DefaultGenesisState()
	potGenesis.Params.MatureEpoch = 1
	genesisState[pottypes.ModuleName] = app.AppCodec().MustMarshalJSON(potGenesis)

	sdsGenesis := sdstypes.DefaultGenesisState()
	genesisState[sdstypes.ModuleName] = app.AppCodec().MustMarshalJSON(sdsGenesis)

	evmGenesis := evmtypes.DefaultGenesisState()
	genesisState[evmtypes.ModuleName] = app.AppCodec().MustMarshalJSON(evmGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         chainId,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		ChainID:            chainId,
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	return app
}

func setup(chainID string, withGenesis bool, invCheckPeriod uint) (*stratosapp.StratosApp, simapp.GenesisState) {
	db := dbm.NewMemDB()

	appOptions := make(simstestutil.AppOptionsMap)
	appOptions[flags.FlagHome] = stratosapp.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = invCheckPeriod

	app := stratosapp.NewStratosApp(log.NewNopLogger(), db, nil, true, appOptions, baseapp.SetChainID(chainID))
	if withGenesis {
		//return app, stratosapp.ModuleBasics.DefaultGenesis(app.AppCodec())
		return app, stratosapp.NewDefaultGenesisState(app.AppCodec())
	}
	return app, simapp.GenesisState{}
}

// CheckBalance checks the balance of an account.
func CheckBalance(t *testing.T, app *stratosapp.StratosApp, addr sdk.AccAddress, balances sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	require.True(t, balances.IsEqual(app.GetBankKeeper().GetAllBalances(ctxCheck, addr)))
}
