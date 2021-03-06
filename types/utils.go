package types

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//ethereum 1 eth = 10^18 wei, stratos 1 stos = 10^9 ustos
//assume 1eth = 1stos, then 1 ustos = 10^9 wei
func UstosToWei(ustosVal sdk.Int) (weiVal sdk.Int, err error) {
	if ustosVal.IsNegative() {
		return weiVal, sdkerrors.Wrap(err, "value is negative")
	}

	return ustosVal.MulRaw(int64(math.Pow10(DenomUnitDiff))), nil
}

func WeiToUstosBigInt(weiValBig *hexutil.Big) (ustosValBig *hexutil.Big, err error) {
	if weiValBig == nil {
		return nil, nil
	}

	ustosVal, err := WeiToUstos(sdk.NewIntFromBigInt(weiValBig.ToInt()))
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(ustosVal.BigInt()), nil
}

func WeiToUstos(weiVal sdk.Int) (ustosVal sdk.Int, err error) {
	if weiVal.IsNegative() {
		return weiVal, sdkerrors.Wrap(err, "value is negative")
	}

	return weiVal.QuoRaw(int64(math.Pow10(DenomUnitDiff))), nil
}
