package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/sds/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

	ctx := sdk.UnwrapSDKContext(c)

	_, err := hex.DecodeString(req.GetFileHash())
	if err != nil {
		return &types.QueryFileUploadResponse{}, fmt.Errorf("invalid file hash, please specify a hash in hex format %w", err)
	}
	fileInfoBytes, err := q.GetFileInfoBytesByFileHash(ctx, []byte(req.GetFileHash()))
	if err != nil {
		return &types.QueryFileUploadResponse{}, err
	}
	fileInfo, err := types.UnmarshalFileInfo(q.cdc, fileInfoBytes)
	if err != nil {
		return &types.QueryFileUploadResponse{}, err
	}

	return &types.QueryFileUploadResponse{FileInfo: &fileInfo}, nil
}

func (q Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)
	return &types.QueryParamsResponse{Params: &params}, nil
}

var _ types.QueryServer = Querier{}
