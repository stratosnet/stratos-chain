package sds

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/sds/keeper"
	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, *data.Params)

	for _, file := range data.GetFileUploads() {
		keeper.SetFileHash(ctx, []byte(file.FileHash), *file.FileInfo)
	}
	return
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)

	var fileUpload []*types.FileUpload
	keeper.IterateFileUpload(ctx, func(fileHash string, fileInfo types.FileInfo) (stop bool) {
		fileUpload = append(fileUpload, &types.FileUpload{FileHash: fileHash, FileInfo: &fileInfo})
		return false
	})

	return types.NewGenesisState(&params, fileUpload)
}
