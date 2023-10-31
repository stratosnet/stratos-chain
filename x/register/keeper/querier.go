package keeper

import (
	"fmt"

	db "github.com/cometbft/cometbft-db"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	pagiquery "github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/register/types"
)

func (k Keeper) getNodeDeposit(ctx sdk.Context, bondStatus stakingtypes.BondStatus, nodeAddress stratos.SdsAddress, tokens sdkmath.Int) (
	unbondingDeposit, unbondedDeposit, bondedDeposit sdkmath.Int, err error) {

	unbondingDeposit = sdkmath.NewInt(0)
	unbondedDeposit = sdkmath.NewInt(0)
	bondedDeposit = sdkmath.NewInt(0)

	switch bondStatus {
	case stakingtypes.Unbonding:
		unbondingDeposit = k.GetUnbondingNodeBalance(ctx, nodeAddress)
	case stakingtypes.Unbonded:
		unbondedDeposit = tokens
	case stakingtypes.Bonded:
		bondedDeposit = tokens
	default:
		err := fmt.Sprintf("Invalid status of node %s, expected Bonded, Unbonded, or Unbonding, got %s",
			nodeAddress.String(), bondStatus)
		return sdkmath.Int{}, sdkmath.Int{}, sdkmath.Int{}, errors.Wrap(sdkerrors.ErrPanic, err)
	}
	return unbondingDeposit, unbondedDeposit, bondedDeposit, nil
}

func GetIterator(prefixStore storetypes.KVStore, start []byte, reverse bool) db.Iterator {
	if reverse {
		var end []byte
		if start != nil {
			itr := prefixStore.Iterator(start, nil)
			defer itr.Close()
			if itr.Valid() {
				itr.Next()
				end = itr.Key()
			}
		}
		return prefixStore.ReverseIterator(nil, end)
	}
	return prefixStore.Iterator(start, nil)
}

func FilteredPaginate(cdc codec.Codec,
	prefixStore storetypes.KVStore,
	queryOwnerAddr sdk.AccAddress,
	pageRequest *pagiquery.PageRequest,
	onResult func(key []byte, value []byte, accumulate bool) (bool, error),
) (*pagiquery.PageResponse, error) {

	// if the PageRequest is nil, use default PageRequest
	if pageRequest == nil {
		pageRequest = &pagiquery.PageRequest{}
	}

	offset := pageRequest.Offset
	key := pageRequest.Key
	limit := pageRequest.Limit
	countTotal := pageRequest.CountTotal
	reverse := pageRequest.Reverse

	if offset > 0 && key != nil {
		return nil, fmt.Errorf("invalid request, either offset or key is expected, got both")
	}

	if limit == 0 {
		limit = types.QueryDefaultLimit

		// count total results when the limit is zero/not supplied
		countTotal = pageRequest.CountTotal
	}

	if len(key) != 0 {
		iterator := GetIterator(prefixStore, key, reverse)
		defer iterator.Close()

		var numHits uint64
		var nextKey []byte
		var ownerAddr sdk.AccAddress

		for ; iterator.Valid(); iterator.Next() {
			if numHits == limit {
				nextKey = iterator.Key()
				break
			}

			if iterator.Error() != nil {
				return nil, iterator.Error()
			}

			if prefixStore.Has(types.MetaNodeKey) {
				metaNode, err := types.UnmarshalMetaNode(cdc, iterator.Value())
				if err != nil {
					continue
				}

				ownerAddr, err = sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
				if err != nil {
					continue
				}
			} else {
				resourceNode, err := types.UnmarshalResourceNode(cdc, iterator.Value())
				if err != nil {
					continue
				}

				ownerAddr, err = sdk.AccAddressFromBech32(resourceNode.GetOwnerAddress())
				if err != nil {
					continue
				}
			}

			if queryOwnerAddr.String() != ownerAddr.String() {
				continue
			}

			hit, err := onResult(iterator.Key(), iterator.Value(), true)
			if err != nil {
				return nil, err
			}

			if hit {
				numHits++
			}
		}

		return &pagiquery.PageResponse{
			NextKey: nextKey,
		}, nil
	}

	iterator := GetIterator(prefixStore, nil, reverse)
	defer iterator.Close()

	end := offset + limit

	var numHits uint64
	var nextKey []byte
	var ownerAddr sdk.AccAddress

	for ; iterator.Valid(); iterator.Next() {
		if iterator.Error() != nil {
			return nil, iterator.Error()
		}

		if prefixStore.Has(types.MetaNodeKey) {
			metaNode, err := types.UnmarshalMetaNode(cdc, iterator.Value())
			if err != nil {
				continue
			}

			ownerAddr, err = sdk.AccAddressFromBech32(metaNode.GetOwnerAddress())
			if err != nil {
				continue
			}
		} else {
			resourceNode, err := types.UnmarshalResourceNode(cdc, iterator.Value())
			if err != nil {
				continue
			}

			ownerAddr, err = sdk.AccAddressFromBech32(resourceNode.GetOwnerAddress())
			if err != nil {
				continue
			}
		}

		if queryOwnerAddr.String() != ownerAddr.String() {
			continue
		}
		accumulate := numHits >= offset && numHits < end
		hit, err := onResult(iterator.Key(), iterator.Value(), accumulate)
		if err != nil {
			return nil, err
		}

		if hit {
			numHits++
		}

		if numHits == end+1 {
			nextKey = iterator.Key()

			if !countTotal {
				break
			}
		}
	}

	res := &pagiquery.PageResponse{NextKey: nextKey}
	if countTotal {
		res.Total = numHits
	}

	return res, nil
}

// GetDepositInfosByResourceNodes Iteration for querying DepositInfos of resource nodes by owner(grpc)
func GetDepositInfosByResourceNodes(
	ctx sdk.Context, k Keeper, resourceNodes types.ResourceNodes,
) ([]*types.DepositInfo, error) {
	resp := make([]*types.DepositInfo, len(resourceNodes))

	for i, resourceNode := range resourceNodes {
		depositInfo, err := GetDepositInfoByResourceNode(ctx, k, resourceNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &depositInfo
	}

	return resp, nil
}

// GetDepositInfosByMetaNodes Iteration for querying DepositInfos of meta nodes by owner(grpc)
func GetDepositInfosByMetaNodes(
	ctx sdk.Context, k Keeper, metaNodes types.MetaNodes,
) ([]*types.DepositInfo, error) {

	resp := make([]*types.DepositInfo, len(metaNodes))

	for i, metaNode := range metaNodes {
		depositInfo, err := GetDepositInfoByMetaNode(ctx, k, metaNode)
		if err != nil {
			return nil, err
		}

		resp[i] = &depositInfo

	}

	return resp, nil
}

func GetDepositInfoByResourceNode(ctx sdk.Context, k Keeper, node types.ResourceNode) (types.DepositInfo, error) {
	networkAddr, _ := stratos.SdsAddressFromBech32(node.GetNetworkAddress())
	depositInfo := types.DepositInfo{}
	unBondingDeposit, unBondedDeposit, bondedDeposit, er := k.getNodeDeposit(
		ctx,
		node.GetStatus(),
		networkAddr,
		node.Tokens,
	)
	if er != nil {
		return depositInfo, er
	}

	if !node.Equal(types.ResourceNode{}) {
		depositInfo = types.NewDepositInfoByResourceNodeAddr(
			k.BondDenom(ctx),
			node,
			unBondingDeposit,
			unBondedDeposit,
			bondedDeposit,
		)
	}
	return depositInfo, nil
}

func GetDepositInfoByMetaNode(ctx sdk.Context, k Keeper, node types.MetaNode) (types.DepositInfo, error) {
	networkAddr, _ := stratos.SdsAddressFromBech32(node.GetNetworkAddress())
	depositInfo := types.DepositInfo{}
	unBondingDeposit, unBondedDeposit, bondedDeposit, er := k.getNodeDeposit(
		ctx,
		node.GetStatus(),
		networkAddr,
		node.Tokens,
	)
	if er != nil {
		return depositInfo, er
	}

	if !node.Equal(types.MetaNode{}) {
		depositInfo = types.NewDepositInfoByMetaNodeAddr(
			k.BondDenom(ctx),
			node,
			unBondingDeposit,
			unBondedDeposit,
			bondedDeposit,
		)
	}
	return depositInfo, nil
}
