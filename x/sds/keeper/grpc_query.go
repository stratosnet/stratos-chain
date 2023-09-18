package keeper

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/x/sds/types"
)

var _ types.QueryServer = Querier{}

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

func (q Querier) Fileupload(c context.Context, req *types.QueryFileUploadRequest) (*types.QueryFileUploadResponse, error) {
	if req == nil {
		return &types.QueryFileUploadResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GetFileHash() == "" {
		return &types.QueryFileUploadResponse{}, status.Error(codes.InvalidArgument, " Network address cannot be empty")
	}

	_, err := cid.Decode(req.GetFileHash())
	if err != nil {
		return &types.QueryFileUploadResponse{}, fmt.Errorf("invalid file hash %w", err)
	}

	ctx := sdk.UnwrapSDKContext(c)

	fileInfo, found := q.GetFileInfoByFileHash(ctx, []byte(req.GetFileHash()))
	if !found {
		return &types.QueryFileUploadResponse{}, types.ErrNoFileFound
	}

	return &types.QueryFileUploadResponse{FileInfo: &fileInfo}, nil
}

func (q Querier) SimPrepay(c context.Context, request *types.QuerySimPrepayRequest) (*types.QuerySimPrepayResponse, error) {
	if request == nil {
		return &types.QuerySimPrepayResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.GetAmount() == nil {
		return &types.QuerySimPrepayResponse{}, status.Error(codes.InvalidArgument, "Amount cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	noz := q.simulatePurchaseNoz(ctx, request.GetAmount())
	return &types.QuerySimPrepayResponse{Noz: noz}, nil
}

func (q Querier) NozPrice(c context.Context, _ *types.QueryNozPriceRequest) (*types.QueryNozPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	St, Pt, Lt := q.registerKeeper.GetCurrNozPriceParams(ctx)
	nozPrice := q.potKeeper.GetCurrentNozPrice(St, Pt, Lt)
	return &types.QueryNozPriceResponse{Price: nozPrice}, nil
}

func (q Querier) NozSupply(c context.Context, request *types.QueryNozSupplyRequest) (*types.QueryNozSupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	remaining, total := q.potKeeper.NozSupply(ctx)
	return &types.QueryNozSupplyResponse{Remaining: remaining, Total: total}, nil
}

func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)
	return &types.QueryParamsResponse{Params: &params}, nil
}
