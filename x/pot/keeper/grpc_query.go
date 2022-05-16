package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/x/pot/types"
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

	if req.Epoch.LTE(sdk.ZeroInt()) {
		return &types.QueryVolumeReportResponse{}, status.Error(codes.InvalidArgument, "epoch should be positive value")
	}

	epoch, ok := sdk.NewIntFromString(req.Epoch.String())
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
