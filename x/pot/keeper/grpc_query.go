package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pagiquery "github.com/cosmos/cosmos-sdk/types/query"

	"github.com/stratosnet/stratos-chain/x/pot/types"
	registerkeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)

	return &types.QueryParamsResponse{Params: &params}, nil
}

func (q Querier) VolumeReport(c context.Context, req *types.QueryVolumeReportRequest) (*types.QueryVolumeReportResponse, error) {
	if req == nil {
		return &types.QueryVolumeReportResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	epochInt64 := req.GetEpoch()
	epoch := sdkmath.NewInt(epochInt64)

	ctx := sdk.UnwrapSDKContext(c)
	if epoch.LT(sdkmath.ZeroInt()) {
		return &types.QueryVolumeReportResponse{}, status.Error(codes.InvalidArgument, "epoch should be positive value")
	}
	if epoch.Equal(sdkmath.ZeroInt()) {
		epoch = q.GetLastDistributedEpoch(ctx)
	}

	height := ctx.BlockHeight()
	volumeReport := q.GetVolumeReport(ctx, epoch)

	return &types.QueryVolumeReportResponse{
		ReportInfo: &types.ReportInfo{
			Epoch:     epoch.Int64(),
			Reference: volumeReport.ReportReference,
			Reporter:  volumeReport.Reporter,
			TxHash:    volumeReport.TxHash,
		},
		Height: height,
	}, nil
}

func (q Querier) RewardsByEpoch(c context.Context, req *types.QueryRewardsByEpochRequest) (*types.QueryRewardsByEpochResponse, error) {
	if req == nil {
		return &types.QueryRewardsByEpochResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	queryEpoch := sdkmath.NewInt(req.GetEpoch())

	if queryEpoch.LTE(sdkmath.ZeroInt()) {
		return &types.QueryRewardsByEpochResponse{}, status.Error(codes.InvalidArgument, "epoch cannot be equal to or lower than 0")
	}

	walletAddr, err := sdk.AccAddressFromBech32(req.GetWalletAddress())
	if err != nil {
		return &types.QueryRewardsByEpochResponse{}, status.Error(codes.Internal, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	matureEpoch := queryEpoch.Add(sdkmath.NewInt(q.MatureEpoch(ctx)))
	var res []*types.Reward

	store := ctx.KVStore(q.storeKey)
	//RewardStore := prefix.NewStore(store, types.GetIndividualRewardIteratorKey(matureEpoch))
	RewardStore := prefix.NewStore(store, types.GetIndividualRewardKey(walletAddr, matureEpoch))

	rewardsPageRes, err := FilteredPaginate(q.cdc, RewardStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		val, err := UnmarshalIndividualReward(q.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			res = append(res, &val)
		}

		return true, nil
	})
	if err != nil {
		return &types.QueryRewardsByEpochResponse{}, status.Error(codes.Internal, err.Error())
	}
	height := ctx.BlockHeight()

	return &types.QueryRewardsByEpochResponse{Rewards: res, Height: height, Pagination: rewardsPageRes}, nil
}

func UnmarshalIndividualReward(cdc codec.Codec, value []byte) (v types.Reward, err error) {
	err = cdc.UnmarshalLengthPrefixed(value, &v)
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
		limit = types.QueryDefaultLimit

		// count total results when the limit is zero/not supplied
		countTotal = pageRequest.CountTotal
	}

	if len(key) != 0 {
		iterator := GetIterator(prefixStore, key, reverse)
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
			v := iterator.Value()
			reward, err := UnmarshalIndividualReward(cdc, v)
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

	iterator := GetIterator(prefixStore, nil, reverse)
	defer iterator.Close()

	end := offset + limit

	var numHits uint64
	var nextKey []byte

	for ; iterator.Valid(); iterator.Next() {
		if iterator.Error() != nil {
			return nil, iterator.Error()
		}

		v := iterator.Value()
		reward, err := UnmarshalIndividualReward(cdc, v)
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

func (q Querier) RewardsByOwner(c context.Context, req *types.QueryRewardsByOwnerRequest) (*types.QueryRewardsByOwnerResponse, error) {
	if req == nil {
		return &types.QueryRewardsByOwnerResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.GetWalletAddress() == "" {
		return &types.QueryRewardsByOwnerResponse{}, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	height := ctx.BlockHeight()

	walletAddr, err := sdk.AccAddressFromBech32(req.GetWalletAddress())
	if err != nil {
		return &types.QueryRewardsByOwnerResponse{}, err
	}

	immatureTotalReward := q.GetImmatureTotalReward(ctx, walletAddr)
	matureTotalReward := q.GetMatureTotalReward(ctx, walletAddr)
	reward := types.NewRewardInfo(walletAddr, matureTotalReward, immatureTotalReward)
	return &types.QueryRewardsByOwnerResponse{Rewards: &reward, Height: height}, nil

}

func (q Querier) SlashingByOwner(c context.Context, req *types.QuerySlashingByOwnerRequest) (*types.QuerySlashingByOwnerResponse, error) {
	if req == nil {
		return &types.QuerySlashingByOwnerResponse{}, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.GetWalletAddress() == "" {
		return &types.QuerySlashingByOwnerResponse{}, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	height := ctx.BlockHeight()

	walletAddr, err := sdk.AccAddressFromBech32(req.GetWalletAddress())
	if err != nil {
		return &types.QuerySlashingByOwnerResponse{}, err
	}

	slashing := q.registerKeeper.GetSlashing(ctx, walletAddr).String()
	return &types.QuerySlashingByOwnerResponse{Slashing: slashing, Height: height}, nil

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
		limit = types.QueryDefaultLimit

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

func (q Querier) TotalMinedToken(c context.Context, _ *types.QueryTotalMinedTokenRequest) (*types.QueryTotalMinedTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	totalMinedToken := q.GetTotalMinedTokens(ctx)

	return &types.QueryTotalMinedTokenResponse{TotalMinedToken: totalMinedToken}, nil
}

func (q Querier) CirculationSupply(c context.Context, _ *types.QueryCirculationSupplyRequest) (*types.QueryCirculationSupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	circulationSupply := q.GetCirculationSupply(ctx)

	return &types.QueryCirculationSupplyResponse{CirculationSupply: circulationSupply}, nil
}

func (q Querier) TotalRewardByEpoch(c context.Context, req *types.QueryTotalRewardByEpochRequest) (
	*types.QueryTotalRewardByEpochResponse, error) {
	if req == nil {
		return &types.QueryTotalRewardByEpochResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	epochInt64 := req.GetEpoch()
	if sdkmath.NewInt(epochInt64).LTE(sdkmath.ZeroInt()) {
		return &types.QueryTotalRewardByEpochResponse{}, status.Error(codes.InvalidArgument, "epoch should be positive value")
	}
	epoch := sdkmath.NewInt(epochInt64)

	ctx := sdk.UnwrapSDKContext(c)

	var isLegacy = false
	totalReward := q.GetTotalReward(ctx, epoch)
	if totalReward.Size() == 0 {
		totalReward = q.GetLegacyTotalReward(ctx, epoch)
		isLegacy = true
	}
	return &types.QueryTotalRewardByEpochResponse{TotalReward: totalReward, IsLegacy: isLegacy}, nil
}

func (q Querier) Metrics(c context.Context, req *types.QueryMetricsRequest) (
	*types.QueryMetricsResponse, error) {
	if req == nil {
		return &types.QueryMetricsResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	metrics := q.GetMetrics(ctx)
	return &types.QueryMetricsResponse{
		Metrics: metrics,
	}, nil
}
