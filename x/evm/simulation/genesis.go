package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// GenExtraEIPs randomly generates specific extra eips or not.
func genExtraEIPs(r *rand.Rand) []int64 {
	var extraEIPs []int64
	if r.Uint32()%2 == 0 {
		extraEIPs = []int64{1344, 1884, 2200, 2929, 3198, 3529}
	}
	return extraEIPs
}

// RandomizedGenState generates a random GenesisState for nft
func RandomizedGenState(simState *module.SimulationState) {
	feeMarketParams := types.NewFeeMarketParams(simState.Rand.Uint32()%2 == 0, simState.Rand.Uint32(), simState.Rand.Uint32(), simState.Rand.Uint64(), simState.Rand.Int63())
	blockGas := simState.Rand.Uint64()

	extraEIPs := genExtraEIPs(simState.Rand)
	params := types.NewParams(types.DefaultEVMDenom, true, true, types.DefaultChainConfig(), feeMarketParams, extraEIPs...)
	evmGenesis := types.NewGenesisState(params, []types.GenesisAccount{}, blockGas)

	bz, err := json.MarshalIndent(evmGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(evmGenesis)
}
