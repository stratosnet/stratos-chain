package keeper

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	stratos "github.com/stratosnet/stratos-chain/types"
)

func CheckAccAddr(w http.ResponseWriter, r *http.Request, data string) (sdk.AccAddress, bool) {
	AccAddr, err := sdk.AccAddressFromBech32(data)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid NodeAddress.")
		return nil, false
	}
	return AccAddr, true
}

func CheckSdsAddr(w http.ResponseWriter, r *http.Request, data string) (stratos.SdsAddress, bool) {
	sdsAddr, err := stratos.SdsAddressFromBech32(data)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid NodeAddress.")
		return nil, false
	}
	return sdsAddr, true
}

func hasValue(items []stratos.SdsAddress, item stratos.SdsAddress) bool {
	for _, eachItem := range items {
		if eachItem.Equals(item) {
			return true
		}
	}
	return false
}

func hasStringValue(items []string, item string) bool {
	for _, eachItem := range items {
		if len(item) > 0 && eachItem == item {
			return true
		}
	}
	return false
}
