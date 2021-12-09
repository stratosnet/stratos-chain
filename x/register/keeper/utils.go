package keeper

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func CheckAccAddr(w http.ResponseWriter, r *http.Request, data string) (sdk.AccAddress, bool) {
	AccAddr, err := sdk.AccAddressFromBech32(data)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid NodeAddress.")
		return nil, false
	}
	return AccAddr, true
}
