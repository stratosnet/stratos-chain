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
	nozAmt sdk.Int, suspend bool) (tokenAmt sdk.Int, nodeType registertypes.NodeType, err error) {

	node, ok := k.registerKeeper.GetResourceNode(ctx, p2pAddr)
	if !ok {
		return sdk.ZeroInt(), registertypes.NodeType(0), registertypes.ErrNoResourceNodeFound
	}
	toBeSuspended := node.Suspend == false && suspend == true
	node.Suspend = suspend

	//slashing amt is equivalent to reward traffic calculation
	trafficList := []types.SingleWalletVolume{{
		WalletAddress: node.OwnerAddress,
		Volume:        &nozAmt,
	}}
	totalConsumedNoz := k.GetTotalConsumedNoz(trafficList).ToDec()
	slashTokenAmt := k.GetTrafficReward(ctx, totalConsumedNoz)

	oldSlashing := k.registerKeeper.GetSlashing(ctx, walletAddr)

	// only slashing the reward token for now.
	newSlashing := oldSlashing.Add(slashTokenAmt.TruncateInt())

	// directly change oz limit while node being suspended
	if toBeSuspended {
		effectiveStakeChange := sdk.ZeroInt().Sub(node.EffectiveTokens)
		node.EffectiveTokens = sdk.ZeroInt()
		k.registerKeeper.DecreaseOzoneLimitBySubtractStake(ctx, effectiveStakeChange.Abs())
	}

	k.registerKeeper.SetResourceNode(ctx, node)
	k.registerKeeper.SetSlashing(ctx, walletAddr, newSlashing)

	return slashTokenAmt.TruncateInt(), registertypes.NodeType(node.NodeType), nil
}
