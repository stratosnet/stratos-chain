package app

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// RandomGenesisAccounts is used by the auth module to create random genesis accounts in simulation when a genesis.json is not specified.
// In contrast, the default auth module's RandomGenesisAccounts implementation creates only base accounts and vestings accounts.
func RandomGenesisAccounts(simState *module.SimulationState) authtypes.GenesisAccounts {
	emptyCodeHash := crypto.Keccak256(nil)
	genesisAccs := make(authtypes.GenesisAccounts, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		bacc := authtypes.NewBaseAccountWithAddress(acc.Address)

		ethacc := &stratos.EthAccount{
			BaseAccount: bacc,
			CodeHash:    common.BytesToHash(emptyCodeHash).String(),
		}
		genesisAccs[i] = ethacc
	}

	return genesisAccs
}
