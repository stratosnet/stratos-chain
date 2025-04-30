package keeper

import (
	"context"
	"encoding/base64"

	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cosmos/gogoproto/proto"
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

func (q Querier) SimPrepay(c context.Context, request *types.QuerySimPrepayRequest) (*types.QuerySimPrepayResponse, error) {
	if request == nil {
		return &types.QuerySimPrepayResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.GetAmount() == "" {
		return &types.QuerySimPrepayResponse{}, status.Error(codes.InvalidArgument, "Amount cannot be empty")
	}
	amount, err := sdk.ParseCoinNormalized(request.GetAmount())
	if err != nil {
		return &types.QuerySimPrepayResponse{}, status.Error(codes.InvalidArgument, "Amount invalid")
	}
	ctx := sdk.UnwrapSDKContext(c)

	noz := q.simulatePurchaseNoz(ctx, sdk.NewCoins(amount))
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

func (q Querier) MerkleRoot(c context.Context, request *types.QueryMerkleRootRequest) (*types.QueryMerkleRootResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if request.Height < 1 {
		return &types.QueryMerkleRootResponse{
			Root: base64.StdEncoding.EncodeToString(q.GetMerkleRoot(ctx)),
		}, nil
	}
	return &types.QueryMerkleRootResponse{
		Root: base64.StdEncoding.EncodeToString(q.GetMerkleRootByHeight(ctx, request.Height)),
	}, nil
}

func (q Querier) VerifyUpload(c context.Context, request *types.QueryVerifyUploadRequest) (*types.QueryVerifyUploadResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if len(request.Proof) == 0 {
		return &types.QueryVerifyUploadResponse{Verified: false, Error: "proof is missing (files before v13 have no proof)"}, nil
	}
	if len(request.FileHash) == 0 {
		return &types.QueryVerifyUploadResponse{}, status.Error(codes.InvalidArgument, "filehash is missing")
	}

	proofBytes, err := base64.StdEncoding.DecodeString(request.Proof)
	if err != nil {
		return &types.QueryVerifyUploadResponse{}, status.Error(codes.InvalidArgument, "proof is not base64 encoded")
	}
	protoProof := &crypto.Proof{}
	err = proto.Unmarshal(proofBytes, protoProof)
	if err != nil {
		return &types.QueryVerifyUploadResponse{}, status.Error(codes.InvalidArgument, "proof is not a valid crypto proof protobuf object")
	}
	merkleProof, err := merkle.ProofFromProto(protoProof)
	if err != nil {
		return &types.QueryVerifyUploadResponse{}, status.Error(codes.InvalidArgument, "proof is not a valid merkle proof protobuf object")
	}

	merkleRoot := q.GetMerkleRootByHeight(ctx, request.Height)
	err = merkleProof.Verify(merkleRoot, []byte(request.FileHash))
	if err != nil {
		return &types.QueryVerifyUploadResponse{
			Verified: false,
			Error:    err.Error(),
		}, nil
	}
	return &types.QueryVerifyUploadResponse{Verified: true}, nil
}
