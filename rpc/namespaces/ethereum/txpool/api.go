package txpool

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/stratosnet/stratos-chain/rpc/backend"
	"github.com/stratosnet/stratos-chain/rpc/types"
)

// PublicAPI offers and API for the transaction pool. It only operates on data that is non-confidential.
type PublicAPI struct {
	logger    log.Logger
	clientCtx client.Context
	backend   backend.BackendI
}

// NewPublicAPI creates a new tx pool service that gives information about the transaction pool.
func NewPublicAPI(logger log.Logger, clientCtx client.Context, backend backend.BackendI) *PublicAPI {
	return &PublicAPI{
		logger:    logger.With("module", "txpool"),
		clientCtx: clientCtx,
		backend:   backend,
	}
}

func (api *PublicAPI) getPendingList() map[common.Address]map[string]*types.Transaction {
	pTxs := types.GetPendingTxs(api.clientCtx.TxConfig.TxDecoder(), api.backend.GetMempool())
	pending := make(map[common.Address]map[string]*types.Transaction)
	for _, pTx := range pTxs {
		addrMap, ok := pending[pTx.From]
		if !ok {
			addrMap = make(map[string]*types.Transaction)
			pending[pTx.From] = addrMap
		}
		addrMap[fmt.Sprintf("%d", pTx.Nonce)] = pTx
	}
	return pending
}

// Content returns the transactions contained within the transaction pool
func (api *PublicAPI) Content() map[string]map[common.Address]map[string]*types.Transaction {
	api.logger.Debug("txpool_content")
	content := map[string]map[common.Address]map[string]*types.Transaction{
		"pending": api.getPendingList(),
		"queued":  make(map[common.Address]map[string]*types.Transaction),
	}
	return content
}

// Inspect returns the content of the transaction pool and flattens it into an
func (api *PublicAPI) Inspect() map[string]map[string]map[string]string {
	api.logger.Debug("txpool_inspect")
	content := map[string]map[string]map[string]string{
		"pending": make(map[string]map[string]string),
		"queued":  make(map[string]map[string]string),
	}
	pending := api.getPendingList()
	// Define a formatter to flatten a transaction into a string
	var format = func(tx *types.Transaction) string {
		if to := tx.To; to != nil {
			return fmt.Sprintf("%s: %d wei + %d gas × %d wei", tx.To.Hex(), tx.Value.ToInt(), tx.Gas, tx.GasPrice.ToInt())
		}
		return fmt.Sprintf("contract creation: %d wei + %d gas × %d wei", tx.Value.ToInt(), tx.Gas, tx.GasPrice.ToInt())
	}
	// Flatten the pending transactions
	for account, txs := range pending {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce)] = format(tx)
		}
		content["pending"][account.Hex()] = dump
	}

	return content
}

// Status returns the number of pending and queued transaction in the pool.
func (api *PublicAPI) Status() map[string]hexutil.Uint {
	api.logger.Debug("txpool_status")
	return map[string]hexutil.Uint{
		"pending": hexutil.Uint(types.GetPendingTxsLen(api.backend.GetMempool())),
		"queued":  hexutil.Uint(0),
	}
}
