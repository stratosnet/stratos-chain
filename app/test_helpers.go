package app

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

	stratostestutil "github.com/stratosnet/stratos-chain/testutil/stratos"
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
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// SetupWithGenesisNodeSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisNodeSet(t *testing.T,
	freshStart bool,
	valSet *tmtypes.ValidatorSet,
	metaNodes []registertypes.MetaNode,
	resourceNodes []registertypes.ResourceNode,
	genAccs []authtypes.GenesisAccount,
	chainId string,
	balances ...banktypes.Balance) *StratosApp {

	app, genesisState := setup(true, 5)
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	validatorBondedAmt := sdkmath.ZeroInt()
	bondAmt := sdkmath.NewInt(1000000)

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdkmath.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
			MinSelfDelegation: sdkmath.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdkmath.LegacyOneDec()))
		validatorBondedAmt = validatorBondedAmt.Add(bondAmt)
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

	initRemainingOzoneLimit := sdkmath.ZeroInt()
	if !freshStart {
		// add bonded amount to bonded pool module account
		balances = append(balances, banktypes.Balance{
			Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
			Coins:   sdk.Coins{sdk.NewCoin(stratos.Wei, validatorBondedAmt)},
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

		initRemainingOzoneLimit = resNodeBondedAmt.ToLegacyDec().
			Quo(registertypes.DefaultDepositNozRate).
			TruncateInt()
	}

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin(stratos.Wei, validatorBondedAmt))...)
	}

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	//registerGenesis := registertypes.DefaultGenesisState()
	registerGenesis := registertypes.NewGenesisState(
		registertypes.DefaultParams(),
		resourceNodes,
		metaNodes,
		initRemainingOzoneLimit,
		make([]registertypes.Slashing, 0),
		registertypes.DefaultDepositNozRate,
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
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: stratostestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
			ChainId:         chainId,
		},
	)

	// commit genesis changes
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
		ChainID:            chainId,
	}})

	return app
}

func setup(withGenesis bool, invCheckPeriod uint) (*StratosApp, simapp.GenesisState) {
	db := dbm.NewMemDB()

	appOptions := make(simstestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = invCheckPeriod

	app := NewStratosApp(log.NewNopLogger(), db, nil, true, appOptions)
	if withGenesis {
		return app, app.DefaultGenesis()
	}
	return app, simapp.GenesisState{}
}

// CheckBalance checks the balance of an account.
func CheckBalance(t *testing.T, app *StratosApp, addr sdk.AccAddress, balances sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	require.True(t, balances.IsEqual(app.GetBankKeeper().GetAllBalances(ctxCheck, addr)))
}
