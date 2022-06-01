package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
)

/*
	This function only record slashing amount.

	Deduct slashing amount when:
	1, calculate upcoming mature reward, deduct from mature_total & upcoming mature reward.
	2, unstaking indexing node.
	3, unstaking resource node.
*/
func (k Keeper) SlashingResourceNode(ctx sdk.Context, p2pAddr stratos.SdsAddress, walletAddr sdk.AccAddress,
	ozAmt sdk.Int, suspend bool) (amt sdk.Int, nodeType registertypes.NodeType, err error) {

	node, ok := k.RegisterKeeper.GetResourceNode(ctx, p2pAddr)
	if !ok {
		return sdk.ZeroInt(), registertypes.NodeType(0), registertypes.ErrNoResourceNodeFound
	}

	node.Suspend = suspend

	//slashing amt is equivalent to reward traffic calculation
	trafficList := []*types.SingleWalletVolume{{
		WalletAddress: node.OwnerAddress,
		Volume:        &ozAmt,
	}}
	totalConsumedUoz := k.GetTotalConsumedUoz(trafficList).ToDec()
	slash := k.GetTrafficReward(ctx, totalConsumedUoz)

	oldSlashing := k.RegisterKeeper.GetSlashing(ctx, walletAddr)

	// only slashing the reward token for now.
	newSlashing := oldSlashing.Add(slash.TruncateInt())

	k.RegisterKeeper.SetResourceNode(ctx, node)
	k.RegisterKeeper.SetSlashing(ctx, walletAddr, newSlashing)
	resourceNodeType, err := strconv.Atoi(node.NodeType)
	if err != nil {
		return sdk.ZeroInt(), registertypes.NodeType(0), registertypes.ErrNodeType
	}

	return slash.TruncateInt(), registertypes.NodeType(resourceNodeType), nil
}
