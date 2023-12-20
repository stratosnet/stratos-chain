package v011

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/register/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, legacySubspace types.ParamsSubspace, cdc codec.Codec) error {
	store := ctx.KVStore(storeKey)

	// migrate params
	if err := migrateParams(ctx, store, cdc, legacySubspace); err != nil {
		return err
	}

	// add beneficiary address
	if err := migrateMetaNodes(store, cdc); err != nil {
		return err
	}

	// add beneficiary address
	if err := migrateResourceNode(store, cdc); err != nil {
		return err
	}

	return nil
}

// migrateParams will set the params to store from legacySubspace
func migrateParams(ctx sdk.Context, store storetypes.KVStore, cdc codec.Codec, legacySubspace types.ParamsSubspace) error {
	var legacyParams types.Params
	legacySubspace.GetParamSet(ctx, &legacyParams)

	if err := legacyParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&legacyParams)
	store.Set(types.ParamsKey, bz)
	return nil
}

func migrateMetaNodes(store sdk.KVStore, cdc codec.Codec) error {
	oldMetaNodeStore := prefix.NewStore(store, MetaNodeKey)
	iterator := oldMetaNodeStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		oldMetaNode := MustUnmarshalMetaNode(cdc, iterator.Value())
		newMetaNode := types.MetaNode{
			NetworkAddress:     oldMetaNode.NetworkAddress,
			Pubkey:             oldMetaNode.Pubkey,
			Suspend:            oldMetaNode.Suspend,
			Status:             oldMetaNode.Status,
			Tokens:             oldMetaNode.Tokens,
			OwnerAddress:       oldMetaNode.OwnerAddress,
			BeneficiaryAddress: oldMetaNode.OwnerAddress,
			Description: types.Description{
				Moniker:         oldMetaNode.Description.Moniker,
				Identity:        oldMetaNode.Description.Identity,
				Website:         oldMetaNode.Description.Website,
				SecurityContact: oldMetaNode.Description.SecurityContact,
				Details:         oldMetaNode.Description.Details,
			},
			CreationTime: oldMetaNode.CreationTime,
		}

		newMetaNodeBz := types.MustMarshalMetaNode(cdc, newMetaNode)
		storeKey := types.GetMetaNodeKey(key)

		oldMetaNodeStore.Delete(iterator.Key())
		store.Set(storeKey, newMetaNodeBz)
	}

	return nil
}

func migrateResourceNode(store sdk.KVStore, cdc codec.Codec) error {
	oldResourceNodeStore := prefix.NewStore(store, ResourceNodeKey)
	iterator := oldResourceNodeStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		oldResourceNode := MustUnmarshalResourceNode(cdc, iterator.Value())
		newResourceNode := types.ResourceNode{
			NetworkAddress:     oldResourceNode.NetworkAddress,
			Pubkey:             oldResourceNode.Pubkey,
			Suspend:            oldResourceNode.Suspend,
			Status:             oldResourceNode.Status,
			Tokens:             oldResourceNode.Tokens,
			OwnerAddress:       oldResourceNode.OwnerAddress,
			BeneficiaryAddress: oldResourceNode.OwnerAddress,
			Description: types.Description{
				Moniker:         oldResourceNode.Description.Moniker,
				Identity:        oldResourceNode.Description.Identity,
				Website:         oldResourceNode.Description.Website,
				SecurityContact: oldResourceNode.Description.SecurityContact,
				Details:         oldResourceNode.Description.Details,
			},
			CreationTime:    oldResourceNode.CreationTime,
			NodeType:        oldResourceNode.NodeType,
			EffectiveTokens: oldResourceNode.EffectiveTokens,
		}

		newResourceNodeBz := types.MustMarshalResourceNode(cdc, newResourceNode)
		storeKey := types.GetResourceNodeKey(key)

		oldResourceNodeStore.Delete(iterator.Key())
		store.Set(storeKey, newResourceNodeBz)
	}

	return nil
}
