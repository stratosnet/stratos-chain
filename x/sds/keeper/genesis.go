package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/sds/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	k.SetParams(ctx, data.Params)

	for _, file := range data.GetFiles() {
		k.SetFileInfo(ctx, []byte(file.FileHash), file.GetFileInfo())
	}
	return
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func (k Keeper) ExportGenesis(ctx sdk.Context) (data *types.GenesisState) {
	params := k.GetParams(ctx)

	var files []types.GenesisFileInfo
	k.IterateFileInfo(ctx, func(fileHash string, fileInfo types.FileInfo) (stop bool) {
		files = append(files, types.GenesisFileInfo{FileHash: fileHash, FileInfo: fileInfo})
		return false
	})

	return types.NewGenesisState(params, files)
}
