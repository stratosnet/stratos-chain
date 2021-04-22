package common

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sds "github.com/stratosnet/stratos-chain/x/sds/types"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// QueryUploadedFile queries the hash of an uploaded file by sender
// validator.
func QueryUploadedFile(cliCtx context.CLIContext, queryRoute, sender string) ([]byte, int64, error) {
	accAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return nil, 0, fmt.Errorf("Invalid sender, please specify a sender in Bech32 format %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryUploadedFile)
	return cliCtx.QueryWithData(route, accAddr)
}
