package evm

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/store"
	tmtypes "github.com/cometbft/cometbft/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/rpc/types"
)

type Context struct {
	logger     log.Logger
	ms         storetypes.MultiStore
	blockStore *store.BlockStore
}

func NewContext(logger log.Logger, ms storetypes.MultiStore, blockStore *store.BlockStore) *Context {
	return &Context{
		logger:     logger,
		ms:         ms,
		blockStore: blockStore,
	}
}

func (c *Context) copySdkContext(ms storetypes.MultiStore, header *tmtypes.Header) sdk.Context {
	sdkCtx := sdk.NewContext(ms, tmproto.Header{}, true, c.logger)
	if header != nil {
		return sdkCtx.WithHeaderHash(
			header.Hash(),
		).WithBlockHeader(
			types.FormatTmHeaderToProto(header),
		).WithBlockHeight(
			header.Height,
		).WithProposer(
			sdk.ConsAddress(header.ProposerAddress),
		)
	}
	return sdkCtx
}

func (c *Context) GetSdkContext() sdk.Context {
	return c.copySdkContext(c.ms.CacheMultiStore(), nil)
}

func (c *Context) GetSdkContextWithHeader(header *tmtypes.Header) (sdk.Context, error) {
	if header == nil {
		return c.GetSdkContext(), nil
	}
	if c.blockStore != nil {
		latestHeight := c.blockStore.Height()
		if latestHeight == 0 {
			return sdk.Context{}, fmt.Errorf("block store not loaded")
		}
		if latestHeight == header.Height {
			return c.copySdkContext(c.ms.CacheMultiStore(), header), nil
		}
	}

	cms, err := c.ms.CacheMultiStoreWithVersion(header.Height)
	if err != nil {
		return sdk.Context{}, err
	}
	return c.copySdkContext(cms, header), nil
}
