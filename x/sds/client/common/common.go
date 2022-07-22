package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

// QueryUploadedFile queries the hash of an uploaded file by sender
//func QueryUploadedFile(clientCtx client.Context, queryRoute, fileHashHex string) ([]byte, int64, error) {
//	fileHashByteArr, err := hex.DecodeString(fileHashHex)
//	if err != nil {
//		return nil, 0, fmt.Errorf("invalid file hash, please specify a hash in hex format %w", err)
//	}
//	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryUploadedFile)
//	return clientCtx.QueryWithData(route, fileHashByteArr)
//}

// QueryPrepayBalance queries the prepaid balance by sender in VolumnPool
//func QueryPrepayBalance(clientCtx client.Context, queryRoute, sender string) ([]byte, int64, error) {
//	accAddr, err := sdk.AccAddressFromBech32(sender)
//	if err != nil {
//		return nil, 0, fmt.Errorf("invalid sender, please specify a sender in Bech32 format %w", err)
//	}
//	route := fmt.Sprintf("custom/%s/%s", queryRoute, sds.QueryPrepay)
//	return clientCtx.QueryWithData(route, accAddr)
//}

// QuerySimulatePrepay queries the ongoing price for prepay
func QuerySimulatePrepay(clientCtx client.Context, amtToPrepay sdk.Int) ([]byte, int64, error) {
	amtByteArray, err := amtToPrepay.MarshalJSON()
	if err != nil {
		return nil, 0, fmt.Errorf("invalid amount, please specify a valid amount to simulate prepay %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", sdstypes.QuerierRoute, sdstypes.QuerySimulatePrepay)
	return clientCtx.QueryWithData(route, amtByteArray)
}

// QueryCurrUozPrice queries the current price for uoz
func QueryCurrUozPrice(clientCtx client.Context) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", sdstypes.QuerierRoute, sdstypes.QueryCurrUozPrice)
	return clientCtx.QueryWithData(route, nil)
}

// QueryUozSupply QueryCurrUozPrice queries the current price for uoz
func QueryUozSupply(clientCtx client.Context) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", sdstypes.QuerierRoute, sdstypes.QueryUozSupply)
	return clientCtx.QueryWithData(route, nil)
}
