package keeper

import (
	"context"
	"fmt"
	"strconv"

	//"github.com/cosmos/cosmos-sdk/client"
	//"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pagiquery "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stratosnet/stratos-chain/x/pot/types"
	registerkeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (q Querier) VolumeReport(c context.Context, req *types.QueryVolumeReportRequest) (*types.QueryVolumeReportResponse, error) {
	if req == nil {
		return &types.QueryVolumeReportResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	epochInt64, err := strconv.ParseInt(req.Epoch, 10, 64)
	if err != nil {
		return &types.QueryVolumeReportResponse{}, status.Error(codes.InvalidArgument, "invalid epoch")
	}

	if sdk.NewInt(epochInt64).LTE(sdk.ZeroInt()) {
		return &types.QueryVolumeReportResponse{}, status.Error(codes.InvalidArgument, "epoch should be positive value")
	}

	epoch, ok := sdk.NewIntFromString(req.Epoch)
	if !ok {
		return &types.QueryVolumeReportResponse{}, status.Error(codes.InvalidArgument, "invalid epoch")
	}
	ctx := sdk.UnwrapSDKContext(c)
	height := ctx.BlockHeight()
	volumeReport := q.GetVolumeReport(ctx, epoch)

	return &types.QueryVolumeReportResponse{
		ReportInfo: &types.ReportInfo{
			Epoch:     epoch.String(),
			Reference: volumeReport.ReportReference,
		},
		Height: height,
	}, nil
}

func (q Querier) PotRewardsByEpoch(c context.Context, req *types.QueryPotRewardsByEpochRequest) (*types.QueryPotRewardsByEpochResponse, error) {
	if req == nil {
		return &types.QueryPotRewardsByEpochResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	queryEpochStr := req.GetEpoch()

	if queryEpochStr == "" {
		return &types.QueryPotRewardsByEpochResponse{}, status.Error(codes.InvalidArgument, "epoch cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	queryEpochInt64, err := strconv.ParseInt(queryEpochStr, 10, 64)
	if err != nil {
		return nil, err
	}
	queryEpoch := sdk.NewInt(queryEpochInt64)
	matureEpoch := queryEpoch.Add(sdk.NewInt(q.MatureEpoch(ctx)))
	var res []*types.Reward

	store := ctx.KVStore(q.storeKey)
	RewardStore := prefix.NewStore(store, types.GetIndividualRewardIteratorKey(matureEpoch))

	//if req.Pagination.Limit == 0 {
	//	req.Pagination.Limit = registertypes.QueryDefaultLimit
	//
	//	// count total results when the limit is zero/not supplied
	//	req.Pagination.CountTotal = true
	//}

	rewardsPageRes, err := FilteredPaginate(q.cdc, RewardStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		val, err := UnmarshalIndividualReward(q.cdc, value)
		if err != nil {
			return true, err
		}

		if accumulate {
			res = append(res, &val)
		}

		return true, nil
	})
	if err != nil {
		return &types.QueryPotRewardsByEpochResponse{}, status.Error(codes.Internal, err.Error())
	}

	//	height := ctx.BlockHeight()
	//	return &types.QueryPotRewardsByEpochResponse{Rewards: res, Height: height}, nil
	//	q.IteratorIndividualReward(ctx, matureEpoch, func(walletAddress sdk.AccAddress, individualReward types.Reward) (stop bool) {
	//		if !((individualReward.RewardFromMiningPool.Empty() || individualReward.RewardFromMiningPool.IsZero()) &&
	//			(individualReward.RewardFromTrafficPool.Empty() || individualReward.RewardFromTrafficPool.IsZero())) {
	//			res = append(res, &individualReward)
	//		}
	//		return false
	//	})
	//
	//	offset := req.Pagination.Offset
	//	page := 1
	//	if offset != 0 {
	//		page =
	//	}
	//	start, end := client.Paginate(len(res), params.Page, params.Limit, QueryDefaultLimit)
	//	if start < 0 || end < 0 {
	//		return &types.QueryPotRewardsByEpochResponse{}, nil
	//	} else {
	//		res = res[start:end]
	//	}
	height := ctx.BlockHeight()
	return &types.QueryPotRewardsByEpochResponse{Rewards: res, Height: height, Pagination: rewardsPageRes}, nil
	//}
}

func UnmarshalIndividualReward(cdc codec.BinaryCodec, value []byte) (v types.Reward, err error) {
	err = cdc.Unmarshal(value, &v)
	return v, err
}

func FilteredPaginate(cdc codec.Codec,
	prefixStore storetypes.KVStore,
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
		limit = registertypes.QueryDefaultLimit

		// count total results when the limit is zero/not supplied
		countTotal = pageRequest.CountTotal
	}

	if len(key) != 0 {
		iterator := registerkeeper.GetIterator(prefixStore, key, reverse)
		defer iterator.Close()

		var numHits uint64
		var nextKey []byte

		for ; iterator.Valid(); iterator.Next() {
			if numHits == limit {
				nextKey = iterator.Key()
				break
			}

			if iterator.Error() != nil {
				return nil, iterator.Error()
			}
			reward, err := UnmarshalIndividualReward(cdc, iterator.Value())
			if err != nil {
				continue
			}
			if (reward.RewardFromMiningPool.Empty() || reward.RewardFromMiningPool.IsZero()) &&
				(reward.RewardFromTrafficPool.Empty() || reward.RewardFromTrafficPool.IsZero()) {
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

	iterator := registerkeeper.GetIterator(prefixStore, nil, reverse)
	defer iterator.Close()

	end := offset + limit

	var numHits uint64
	var nextKey []byte

	for ; iterator.Valid(); iterator.Next() {
		if iterator.Error() != nil {
			return nil, iterator.Error()
		}

		reward, err := UnmarshalIndividualReward(cdc, iterator.Value())
		if err != nil {
			continue
		}
		if (reward.RewardFromMiningPool.Empty() || reward.RewardFromMiningPool.IsZero()) &&
			(reward.RewardFromTrafficPool.Empty() || reward.RewardFromTrafficPool.IsZero()) {
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

func (q Querier) PotRewardsByOwner(c context.Context, req *types.QueryPotRewardsByOwnerRequest) (*types.QueryPotRewardsByOwnerResponse, error) {
	if req == nil {
		return &types.QueryPotRewardsByOwnerResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.GetWalletAddress() == "" {
		return &types.QueryPotRewardsByOwnerResponse{}, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	height := ctx.BlockHeight()

	walletAddr, err := sdk.AccAddressFromBech32(req.GetWalletAddress())
	if err != nil {
		return &types.QueryPotRewardsByOwnerResponse{}, err
	}

	immatureTotalReward := q.GetImmatureTotalReward(ctx, walletAddr)
	matureTotalReward := q.GetMatureTotalReward(ctx, walletAddr)
	reward := types.NewPotRewardInfo(walletAddr, matureTotalReward, immatureTotalReward)
	return &types.QueryPotRewardsByOwnerResponse{Rewards: &reward, Height: height}, nil

}

func (q Querier) PotSlashingByOwner(c context.Context, req *types.QueryPotSlashingByOwnerRequest) (*types.QueryPotSlashingByOwnerResponse, error) {
	if req == nil {
		return &types.QueryPotSlashingByOwnerResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.GetWalletAddress() == "" {
		return &types.QueryPotSlashingByOwnerResponse{}, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	height := ctx.BlockHeight()

	walletAddr, err := sdk.AccAddressFromBech32(req.GetWalletAddress())
	if err != nil {
		return &types.QueryPotSlashingByOwnerResponse{}, err
	}

	slashing := q.RegisterKeeper.GetSlashing(ctx, walletAddr).String()
	return &types.QueryPotSlashingByOwnerResponse{Slashing: slashing, Height: height}, nil

}

func Paginate(
	prefixStore storetypes.KVStore,
	pageRequest *pagiquery.PageRequest,
	onResult func(key []byte, value []byte) error,
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
		limit = QueryDefaultLimit

		// count total results when the limit is zero/not supplied
		countTotal = true
	}

	if len(key) != 0 {
		iterator := registerkeeper.GetIterator(prefixStore, key, reverse)
		defer iterator.Close()

		var count uint64
		var nextKey []byte

		for ; iterator.Valid(); iterator.Next() {

			if count == limit {
				nextKey = iterator.Key()
				break
			}
			if iterator.Error() != nil {
				return nil, iterator.Error()
			}
			err := onResult(iterator.Key(), iterator.Value())
			if err != nil {
				return nil, err
			}

			count++
		}

		return &pagiquery.PageResponse{
			NextKey: nextKey,
		}, nil
	}

	iterator := registerkeeper.GetIterator(prefixStore, nil, reverse)
	defer iterator.Close()

	end := offset + limit

	var count uint64
	var nextKey []byte

	for ; iterator.Valid(); iterator.Next() {
		count++

		if count <= offset {
			continue
		}
		if count <= end {
			err := onResult(iterator.Key(), iterator.Value())
			if err != nil {
				return nil, err
			}
		} else if count == end+1 {
			nextKey = iterator.Key()

			if !countTotal {
				break
			}
		}
		if iterator.Error() != nil {
			return nil, iterator.Error()
		}
	}

	res := &pagiquery.PageResponse{NextKey: nextKey}
	if countTotal {
		res.Total = count
	}

	return res, nil
}
