package sds

import (
	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/sds/keeper"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock, _ keeper.Keeper) {
	// 	TODO: fill out if your application requires beginBlock, if not you can delete this function
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock, keeper keeper.Keeper) []abci.ValidatorUpdate {
	newFiles := keeper.ClearNewFiles(ctx)
	if len(newFiles) == 0 {
		return []abci.ValidatorUpdate{}
	}

	previousRoot := keeper.GetMerkleRoot(ctx)
	prover := keeper.GetMerkleProver()
	newRoot, proofs := prover.CreateFileUploadProofs(previousRoot, newFiles)
	keeper.SetMerkleRoot(ctx, newRoot)
	keeper.SetMerkleRootByHeight(ctx, ctx.BlockHeight(), newRoot)

	if len(newFiles)*2 != len(proofs) {
		keeper.Logger(ctx).Error("Invalid number of merkle proofs for file upload", "expected", len(newFiles)*2, "actual", len(proofs))
		return []abci.ValidatorUpdate{}
	}

	var protoProofs []*types.NewFileMerkleProof
	for i, fileHash := range newFiles {
		proof := proofs[2*i+1]
		protoProofs = append(protoProofs, &types.NewFileMerkleProof{
			FileHash: fileHash,
			Total:    proof.Total,
			Index:    proof.Index,
			LeafHash: proof.LeafHash,
			Aunts:    proof.Aunts,
		})
	}

	err := ctx.EventManager().EmitTypedEvents(
		&types.EventNewFilesUploaded{
			Height: ctx.BlockHeight(),
			Proofs: protoProofs,
		})
	if err != nil {
		keeper.Logger(ctx).Error(err.Error())
		return []abci.ValidatorUpdate{}
	}
	return []abci.ValidatorUpdate{}
}
