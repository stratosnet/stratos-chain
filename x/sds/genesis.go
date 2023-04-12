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

	for _, file := range data.GetFiles() {
		keeper.SetFileInfo(ctx, []byte(file.FileHash), file.GetFileInfo())
	}
	return
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)

	var files []types.GenesisFileInfo
	keeper.IterateFileInfo(ctx, func(fileHash string, fileInfo types.FileInfo) (stop bool) {
		files = append(files, types.GenesisFileInfo{FileHash: fileHash, FileInfo: fileInfo})
		return false
	})

	return types.NewGenesisState(&params, files)
}
