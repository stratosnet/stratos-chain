package simulation

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// RandomizedGenState generates a random GenesisState for nft
func RandomizedGenState(simState *module.SimulationState) {
	feeMarketParams := types.NewFeeMarketParams(simState.Rand.Uint32()%2 == 0, simState.Rand.Uint32(), simState.Rand.Uint32(), simState.Rand.Uint64(), simState.Rand.Int63())
	blockGas := simState.Rand.Uint64()

	params := types.NewParams(types.DefaultEVMDenom, true, true, types.DefaultChainConfig(), feeMarketParams)
	if simState.Rand.Uint32()%2 == 0 {
		params = types.NewParams(types.DefaultEVMDenom, true, true, types.DefaultChainConfig(), feeMarketParams, 1344, 1884, 2200, 2929, 3198, 3529)
	}
	evmGenesis := types.NewGenesisState(params, []types.GenesisAccount{}, blockGas)

	bz, err := json.MarshalIndent(evmGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(evmGenesis)
}
