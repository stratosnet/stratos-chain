package keeper

import (
	"context"

	"github.com/stratosnet/stratos-chain/x/register/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

func (q Querier) ResourceNode(ctx context.Context, request *types.QueryResourceNodeRequest) (*types.QueryResourceNodeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) IndexingNode(ctx context.Context, request *types.QueryIndexingNodeRequest) (*types.QueryIndexingNodeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) StakeByNode(ctx context.Context, request *types.QueryStakeByNodeRequest) (*types.QueryStakeByNodeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) StakeByOwner(ctx context.Context, request *types.QueryStakeByOwnerRequest) (*types.QueryStakeByOwnerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) StakeTotal(ctx context.Context, request *types.QueryTotalStakeRequest) (*types.QueryTotalStakeResponse, error) {
	//TODO implement me
	panic("implement me")
}

var _ types.QueryServer = Querier{}
