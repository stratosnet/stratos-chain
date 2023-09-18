package v0_12_0

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.Codec) error {
	store := ctx.KVStore(storeKey)

	if err := migrateMetaNodes(store, cdc); err != nil {
		return err
	}
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
