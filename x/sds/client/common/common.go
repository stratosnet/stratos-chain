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
		return nil, 0, fmt.Errorf("invalid file hash, please specify a hash in hex format %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryUploadedFile)
	return cliCtx.QueryWithData(route, fileHashByteArr)
}

// QueryPrepayBalance queries the prepaid balance by sender in VolumnPool
func QueryPrepayBalance(cliCtx context.CLIContext, queryRoute, sender string) ([]byte, int64, error) {
	accAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid sender, please specify a sender in Bech32 format %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryPrepay)
	return cliCtx.QueryWithData(route, accAddr)
}

// QuerySimulatePrepay queries the ongoing price for prepay
func QuerySimulatePrepay(cliCtx context.CLIContext, queryRoute string, amtToPrepay sdk.Int) ([]byte, int64, error) {
	amtByteArray, err := amtToPrepay.MarshalJSON()
	if err != nil {
		return nil, 0, fmt.Errorf("invalid amount, please specify a valid amount to simulate prepay %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QuerySimulatePrepay)
	return cliCtx.QueryWithData(route, amtByteArray)
}

// QueryCurrUozPrice queries the current price for uoz
func QueryCurrUozPrice(cliCtx context.CLIContext, queryRoute string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryCurrUozPrice)
	return cliCtx.QueryWithData(route, nil)
}

// QueryCurrUozPrice queries the current price for uoz
func QueryUozSupply(cliCtx context.CLIContext, queryRoute string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryUozSupply)
	return cliCtx.QueryWithData(route, nil)
}
