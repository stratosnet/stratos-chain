package common

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sds "github.com/stratosnet/stratos-chain/x/sds/types"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// QueryUploadedFile queries the hash of an uploaded file by sender
func QueryUploadedFile(cliCtx context.CLIContext, queryRoute, fileHashHex string) ([]byte, int64, error) {
	fileHashByteArr, err := hex.DecodeString(fileHashHex)
	if err != nil {
		return nil, 0, fmt.Errorf("Invalid file hash, please specify a hash in hex format %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryUploadedFile)
	return cliCtx.QueryWithData(route, fileHashByteArr)
}

// QueryPrepayBalance queries the prepaid balance by sender in VolumnPool
func QueryPrepayBalance(cliCtx context.CLIContext, queryRoute, sender string) ([]byte, int64, error) {
	accAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return nil, 0, fmt.Errorf("Invalid sender, please specify a sender in Bech32 format %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryPrepay)
	return cliCtx.QueryWithData(route, accAddr)
}
