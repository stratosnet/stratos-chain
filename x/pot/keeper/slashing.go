package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
)

/*
This function only record slashing amount.

Deduct slashing amount when:
1, calculate upcoming mature reward, deduct from mature_total & upcoming mature reward.
2, unstaking meta node.
3, unstaking resource node.
*/
func (k Keeper) SlashingResourceNode(ctx sdk.Context, p2pAddr stratos.SdsAddress, walletAddr sdk.AccAddress,
	nozAmt sdk.Int, suspend bool, effectiveStakeAfter sdk.Int) (tokenAmt sdk.Int, nodeType registertypes.NodeType, unsuspended bool, err error) {

	node, ok := k.RegisterKeeper.GetResourceNode(ctx, p2pAddr)
	if !ok {
		return sdk.ZeroInt(), registertypes.NodeType(0), false, registertypes.ErrNoResourceNodeFound
	}

	toBeUnsuspended := node.Suspend == true && suspend == false

	// before calc ozone limit change, get unbonding stake and calc effective stake to trigger ozLimit change
	unbondingStake := k.RegisterKeeper.GetUnbondingNodeBalance(ctx, p2pAddr)
	// no effective stake after subtracting unbonding stake
	if node.Tokens.LTE(unbondingStake) {
		return sdk.ZeroInt(), registertypes.NodeType(0), toBeUnsuspended, registertypes.ErrInsufficientBalance
	}
	availableStake := node.Tokens.Sub(unbondingStake)
	if availableStake.LT(effectiveStakeAfter) {
		return sdk.ZeroInt(), registertypes.NodeType(0), toBeUnsuspended, registertypes.ErrInsufficientBalance
	}

	effectiveStakeBefore := node.EffectiveTokens
	effectiveStakeChange := effectiveStakeAfter.Sub(effectiveStakeBefore)

	node.Suspend = suspend

	//slashing amt is equivalent to reward traffic calculation
	trafficList := []*types.SingleWalletVolume{{
		WalletAddress: node.OwnerAddress,
		Volume:        &nozAmt,
	}}
	totalConsumedNoz := k.GetTotalConsumedNoz(trafficList).ToDec()
	slashTokenAmt := k.GetTrafficReward(ctx, totalConsumedNoz)

	oldSlashing := k.RegisterKeeper.GetSlashing(ctx, walletAddr)

	// only slashing the reward token for now.
	newSlashing := oldSlashing.Add(slashTokenAmt.TruncateInt())

	k.RegisterKeeper.SetResourceNode(ctx, node)
	k.RegisterKeeper.SetSlashing(ctx, walletAddr, newSlashing)

	if effectiveStakeChange.IsNegative() {
		k.RegisterKeeper.DecreaseOzoneLimitBySubtractStake(ctx, effectiveStakeChange)
	}
	if effectiveStakeChange.IsPositive() {
		k.RegisterKeeper.IncreaseOzoneLimitByAddStake(ctx, effectiveStakeChange)
	}

	return slashTokenAmt.TruncateInt(), registertypes.NodeType(node.NodeType), toBeUnsuspended, nil
}
