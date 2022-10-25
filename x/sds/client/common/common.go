package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
)

// QuerySimulatePrepay queries the ongoing price for prepay
func QuerySimulatePrepay(clientCtx client.Context, amtToPrepay sdk.Int) ([]byte, int64, error) {
	amtByteArray, err := amtToPrepay.MarshalJSON()
	if err != nil {
		return nil, 0, fmt.Errorf("invalid amount, please specify a valid amount to simulate prepay %w", err)
	}
	route := fmt.Sprintf("custom/%s/%s", sdstypes.QuerierRoute, sdstypes.QuerySimulatePrepay)
	return clientCtx.QueryWithData(route, amtByteArray)
}

// QueryCurrNozPrice queries the current price for noz
func QueryCurrNozPrice(clientCtx client.Context) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", sdstypes.QuerierRoute, sdstypes.QueryCurrNozPrice)
	return clientCtx.QueryWithData(route, nil)
}

// QueryNozSupply QueryCurrNozPrice queries the current price for noz
func QueryNozSupply(clientCtx client.Context) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", sdstypes.QuerierRoute, sdstypes.QueryNozSupply)
	return clientCtx.QueryWithData(route, nil)
}
