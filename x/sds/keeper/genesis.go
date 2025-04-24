package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	k.SetParams(ctx, data.Params)

	if len(data.MerkleRoot) > 0 {
		k.SetMerkleRoot(ctx, data.MerkleRoot)
	}
	for _, root := range data.PreviousMerkleRoots {
		k.SetMerkleRootByHeight(ctx, root.Height, root.Root)
	}
	return
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func (k Keeper) ExportGenesis(ctx sdk.Context) (data *types.GenesisState) {
	params := k.GetParams(ctx)

	var roots []types.MerkleRoot
	k.IterateMerkleRoots(ctx, func(height int64, root []byte) (stop bool) {
		roots = append(roots, types.MerkleRoot{
			Height: height,
			Root:   root,
		})
		return false
	})

	return types.NewGenesisState(params, k.GetMerkleRoot(ctx), roots)
}
