package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stratosnet/stratos-chain/x/register"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
)

// Keeper of the pot store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       params.Subspace
	feeCollectorName string // name of the FeeCollector ModuleAccount
	BankKeeper       bank.Keeper
	SupplyKeeper     supply.Keeper
	AccountKeeper    auth.AccountKeeper
	StakingKeeper    staking.Keeper
	RegisterKeeper   register.Keeper
}

// NewKeeper creates a pot keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, feeCollectorName string,
	bankKeeper bank.Keeper, supplyKeeper supply.Keeper, accountKeeper auth.AccountKeeper, stakingKeeper staking.Keeper,
	registerKeeper register.Keeper,
) Keeper {
	keeper := Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName: feeCollectorName,
		BankKeeper:       bankKeeper,
		SupplyKeeper:     supplyKeeper,
		AccountKeeper:    accountKeeper,
		StakingKeeper:    stakingKeeper,
		RegisterKeeper:   registerKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetVolumeReport(ctx sdk.Context, epoch sdk.Int) (res types.VolumeReportRecord) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VolumeReportStoreKey(epoch))
	if bz == nil {
		return types.VolumeReportRecord{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return res
}

func (k Keeper) SetVolumeReport(ctx sdk.Context, epoch sdk.Int, reportRecord types.VolumeReportRecord) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.VolumeReportStoreKey(epoch)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(reportRecord)
	store.Set(storeKey, bz)
}

func (k Keeper) DeleteVolumeReport(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

func (k Keeper) IsSPNode(ctx sdk.Context, addr sdk.AccAddress) (found bool) {
	_, found = k.RegisterKeeper.GetIndexingNode(ctx, addr)
	return found
}

func (k Keeper) getNodeOwnerMap(ctx sdk.Context) map[string]sdk.AccAddress {
	nodeOwnerMap := make(map[string]sdk.AccAddress)
	nodeOwnerMap = k.RegisterKeeper.GetNodeOwnerMapFromIndexingNodes(ctx, nodeOwnerMap)
	nodeOwnerMap = k.RegisterKeeper.GetNodeOwnerMapFromResourceNodes(ctx, nodeOwnerMap)
	return nodeOwnerMap
}

func (k Keeper) setPotRewardRecordByOwnerHeight(ctx sdk.Context, nodeOwnerMap map[string]sdk.AccAddress, epoch sdk.Int, value []NodeRewardsRecord) {
	potRewardsRecordWithOwnerAddr := make(map[string][]NodeRewardsInfo)
	for _, v := range value {
		if ownerAddr, ok := nodeOwnerMap[v.Record.NodeAddress.String()]; ok {
			if _, ok := potRewardsRecordWithOwnerAddr[ownerAddr.String()]; !ok {
				potRewardsRecordWithOwnerAddr[ownerAddr.String()] = []NodeRewardsInfo{}
			}
			potRewardsRecordWithOwnerAddr[ownerAddr.String()] = append(potRewardsRecordWithOwnerAddr[ownerAddr.String()], v.Record)
		}
	}

	for key, val := range potRewardsRecordWithOwnerAddr {
		k.setPotRewardRecord(ctx, epoch, key, val)
	}
}
