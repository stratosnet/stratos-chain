package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParamsV011(ctx sdk.Context) (params types.Params) {
	// cover as it does not exist in new blocks
	defer func() {
		_ = recover()
	}()

	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// ----------------------------------------------------------------------------
// Parent Base Fee
// Required by EIP1559 base fee calculation.
// ----------------------------------------------------------------------------

// GetBaseFeeParam get's the base fee from the paramSpace
// return nil if base fee is not enabled
func (k Keeper) GetBaseFeeParamV011(ctx sdk.Context) *big.Int {
	params := k.GetParamsV011(ctx)
	if params.FeeMarketParams.NoBaseFee {
		return nil
	}

	return params.FeeMarketParams.BaseFee.BigInt()
}

// GetBalanceV011 load account's balance of gas token for th old blocks before migrations
func (k *Keeper) GetBalanceV011(ctx sdk.Context, addr common.Address) *big.Int {
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	params := k.GetParamsV011(ctx)
	// fast check if block already not in v011
	if params.EvmDenom == "" {
		return nil
	}
	coin := k.bankKeeper.GetBalance(ctx, cosmosAddr, params.EvmDenom)
	return coin.Amount.BigInt()
}
