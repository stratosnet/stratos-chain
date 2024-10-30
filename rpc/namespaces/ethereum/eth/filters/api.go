package filters

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/stratosnet/stratos-chain/rpc/backend"
	"github.com/stratosnet/stratos-chain/rpc/types"
)

// filter is a helper struct that holds meta information over the filter type
// and associated subscription in the event system.
type filter struct {
	typ      filters.Type
	deadline *time.Timer // filter is inactive when deadline triggers
	hashes   []common.Hash
	fullTx   bool
	txs      []*types.Transaction
	crit     filters.FilterCriteria
	logs     []*ethtypes.Log
	s        *Subscription // associated subscription in event system
}

// PublicFilterAPI offers support to create and manage filters. This will allow external clients to retrieve various
// information related to the Ethereum protocol such as blocks, transactions and logs.
type PublicFilterAPI struct {
	logger    log.Logger
	clientCtx client.Context
	backend   backend.BackendI
	events    *EventSystem
	filtersMu sync.Mutex
	filters   map[rpc.ID]*filter
}

// NewPublicAPI returns a new PublicFilterAPI instance.
func NewPublicAPI(logger log.Logger, clientCtx client.Context, eventBus *tmtypes.EventBus, b backend.BackendI) *PublicFilterAPI {
	logger = logger.With("api", "filter")
	api := &PublicFilterAPI{
		logger:    logger,
		clientCtx: clientCtx,
		backend:   b,
		filters:   make(map[rpc.ID]*filter),
		events:    NewEventSystem(clientCtx, logger, eventBus, b),
	}

	go api.timeoutLoop()

	return api
}

// timeoutLoop runs every 5 minutes and deletes filters that have not been recently used.
// Tt is started when the api is created.
func (api *PublicFilterAPI) timeoutLoop() {
	var toUninstall []*Subscription
	ticker := time.NewTicker(deadline)
	defer ticker.Stop()

	for {
		<-ticker.C
		api.filtersMu.Lock()
		for id, f := range api.filters {
			select {
			case <-f.deadline.C:
				toUninstall = append(toUninstall, f.s)
				delete(api.filters, id)
			default:
				continue
			}
		}
		api.filtersMu.Unlock()

		// Unsubscribes are processed outside the lock to avoid the following scenario:
		// event loop attempts broadcasting events to still active filters while
		// Unsubscribe is waiting for it to process the uninstall request.
		for _, s := range toUninstall {
			s.Unsubscribe()
		}
		toUninstall = nil
	}
}

// isFilterCapped check the config cap lenght for the filter. Wrapped with mutex
// in order to prevent dirty data read
func (api *PublicFilterAPI) isFilterCapped() bool {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()
	return len(api.filters) >= int(api.backend.RPCFilterCap())
}

// NewPendingTransactionFilter creates a filter that fetches pending transactions
// as transactions enter the pending state.
//
// It is part of the filter package because this filter can be used through the
// `eth_getFilterChanges` polling method that is also used for log filters.
func (api *PublicFilterAPI) NewPendingTransactionFilter(fullTx *bool) rpc.ID {
	api.logger.Debug("NewPendingTransactionFilter", "fullTx", fullTx)

	if api.isFilterCapped() {
		return rpc.ID("error creating pending tx filter: max limit reached")
	}

	pendingTxs := make(chan []*types.Transaction)
	pendingTxSub := api.events.SubscribePendingTxs(pendingTxs)

	api.filtersMu.Lock()
	api.filters[pendingTxSub.ID()] = &filter{typ: filters.PendingTransactionsSubscription, fullTx: fullTx != nil && *fullTx, deadline: time.NewTimer(deadline), hashes: make([]common.Hash, 0), s: pendingTxSub}
	api.filtersMu.Unlock()

	go func() {
		for {
			select {
			case pTx := <-pendingTxs:
				api.filtersMu.Lock()
				if f, found := api.filters[pendingTxSub.ID()]; found {
					f.txs = append(f.txs, pTx...)
				}
				api.filtersMu.Unlock()
			case <-pendingTxSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, pendingTxSub.ID())
				api.filtersMu.Unlock()
				return
			}
		}
	}()

	return pendingTxSub.ID()
}

// NewPendingTransactions creates a subscription that is triggered each time a
// transaction enters the transaction pool. If fullTx is true the full tx is
// sent to the client, otherwise the hash is sent.
func (api *PublicFilterAPI) NewPendingTransactions(ctx context.Context, fullTx *bool) (*rpc.Subscription, error) {
	api.logger.Debug("NewPendingTransactions", "fullTx", fullTx)

	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		txs := make(chan []*types.Transaction, 128)
		pendingTxSub := api.events.SubscribePendingTxs(txs)

		defer pendingTxSub.Unsubscribe()

		for {
			select {
			case txs := <-txs:
				// To keep the original behaviour, send a single tx hash in one notification.
				// TODO(rjl493456442) Send a batch of tx hashes in one notification
				for _, tx := range txs {
					if fullTx != nil && *fullTx {
						notifier.Notify(rpcSub.ID, tx)
					} else {
						notifier.Notify(rpcSub.ID, tx.Hash)
					}
				}
			case <-rpcSub.Err():
				return
			}
		}
	}()

	return rpcSub, nil
}

// NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
// It is part of the filter package since polling goes with eth_getFilterChanges.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newblockfilter
func (api *PublicFilterAPI) NewBlockFilter() rpc.ID {
	api.logger.Debug("NewBlockFilter")

	if api.isFilterCapped() {
		return rpc.ID("error creating pending tx filter: max limit reached")
	}

	headers := make(chan *types.Header)
	headersSub := api.events.SubscribeNewHeads(headers)

	api.filtersMu.Lock()
	api.filters[headersSub.ID()] = &filter{typ: filters.BlocksSubscription, deadline: time.NewTimer(deadline), hashes: []common.Hash{}, s: headersSub}
	api.filtersMu.Unlock()

	go func() {
		for {
			select {
			case h := <-headers:
				api.filtersMu.Lock()
				if f, found := api.filters[headersSub.ID()]; found {
					f.hashes = append(f.hashes, h.Hash)
				}
				api.filtersMu.Unlock()
			case <-headersSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, headersSub.ID())
				api.filtersMu.Unlock()
				return
			}
		}
	}()

	return headersSub.ID()
}

// NewHeads send a notification each time a new (header) block is appended to the chain.
func (api *PublicFilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	api.logger.Debug("NewHeads")

	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		headers := make(chan *types.Header)
		headersSub := api.events.SubscribeNewHeads(headers)

		defer headersSub.Unsubscribe()

		for {
			select {
			case h := <-headers:
				notifier.Notify(rpcSub.ID, h)
			case <-rpcSub.Err():
				return
			}
		}
	}()

	return rpcSub, nil
}

// Logs creates a subscription that fires for all new log that match the given filter criteria.
func (api *PublicFilterAPI) Logs(ctx context.Context, crit filters.FilterCriteria) (*rpc.Subscription, error) {
	api.logger.Debug("Logs")

	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	matchedLogs := make(chan []*ethtypes.Log)
	logsSub, err := api.events.SubscribeLogs(crit, matchedLogs)
	if err != nil {
		return nil, err
	}

	go func() {
		defer logsSub.Unsubscribe()

		for {
			select {
			case logs := <-matchedLogs:
				for _, log := range logs {
					log := log
					notifier.Notify(rpcSub.ID, &log)
				}
			case <-rpcSub.Err():
				return
			}
		}
	}()

	return rpcSub, err
}

// NewFilter creates a new filter and returns the filter id. It can be
// used to retrieve logs when the state changes. This method cannot be
// used to fetch logs that are already stored in the state.
//
// Default criteria for the from and to block are "latest".
// Using "latest" as block number will return logs for mined blocks.
// Using "pending" as block number returns logs for not yet mined (pending) blocks.
// In case logs are removed (chain reorg) previously returned logs are returned
// again but with the removed property set to true.
//
// In case "fromBlock" > "toBlock" an error is returned.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newfilter
func (api *PublicFilterAPI) NewFilter(criteria filters.FilterCriteria) (rpc.ID, error) {
	api.logger.Debug("NewFilter")

	if api.isFilterCapped() {
		return "", fmt.Errorf("error creating filter: max limit reached")
	}

	logs := make(chan []*ethtypes.Log)
	logsSub, err := api.events.SubscribeLogs(criteria, logs)
	if err != nil {
		return "", err
	}

	api.filtersMu.Lock()
	api.filters[logsSub.ID()] = &filter{typ: filters.LogsSubscription, crit: criteria, deadline: time.NewTimer(deadline), hashes: []common.Hash{}, s: logsSub}
	api.filtersMu.Unlock()

	go func() {
		for {
			select {
			case l := <-logs:
				api.filtersMu.Lock()
				if f, found := api.filters[logsSub.ID()]; found {
					f.logs = append(f.logs, l...)
				}
				api.filtersMu.Unlock()
			case <-logsSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, logsSub.ID())
				api.filtersMu.Unlock()
				return
			}
		}
	}()

	return logsSub.ID(), nil
}

// GetLogs returns logs matching the given argument that are stored within the state.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getlogs
func (api *PublicFilterAPI) GetLogs(ctx context.Context, crit filters.FilterCriteria) ([]*ethtypes.Log, error) {
	api.logger.Debug("GetLogs", "crit", crit)

	if len(crit.Topics) > maxTopics {
		return nil, errExceedMaxTopics
	}

	var filter *Filter
	if crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = NewBlockFilter(api.logger, api.backend, crit)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		if begin > 0 && end > 0 && begin > end {
			return nil, errInvalidBlockRange
		}
		// Construct the range filter
		filter = NewRangeFilter(api.logger, api.backend, begin, end, crit.Addresses, crit.Topics)
	}

	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx, int(api.backend.RPCLogsCap()), int64(api.backend.RPCBlockRangeCap()))
	if err != nil {
		return nil, err
	}

	return returnLogs(logs), err
}

// UninstallFilter removes the filter with the given filter id.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_uninstallfilter
func (api *PublicFilterAPI) UninstallFilter(id rpc.ID) bool {
	api.logger.Debug("UninstallFilter", "id", id)

	api.filtersMu.Lock()
	f, found := api.filters[id]
	if found {
		delete(api.filters, id)
	}
	api.filtersMu.Unlock()

	if found {
		f.s.Unsubscribe()
	}
	return found
}

// GetFilterLogs returns the logs for the filter with the given id.
// If the filter could not be found an empty array of logs is returned.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterlogs
func (api *PublicFilterAPI) GetFilterLogs(ctx context.Context, id rpc.ID) ([]*ethtypes.Log, error) {
	api.logger.Debug("GetFilterLogs", "id", id)

	api.filtersMu.Lock()
	f, found := api.filters[id]
	api.filtersMu.Unlock()

	if !found || f.typ != filters.LogsSubscription {
		return nil, errFilterNotFound
	}

	var filter *Filter
	if f.crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = NewBlockFilter(api.logger, api.backend, f.crit)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if f.crit.FromBlock != nil {
			begin = f.crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if f.crit.ToBlock != nil {
			end = f.crit.ToBlock.Int64()
		}
		// Construct the range filter
		filter = NewRangeFilter(api.logger, api.backend, begin, end, f.crit.Addresses, f.crit.Topics)
	}
	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx, int(api.backend.RPCLogsCap()), int64(api.backend.RPCBlockRangeCap()))
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), nil
}

// GetFilterChanges returns the logs for the filter with the given id since
// last time it was called. This can be used for polling.
//
// For pending transaction and block filters the result is []common.Hash.
// (pending)Log filters return []Log.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterchanges
func (api *PublicFilterAPI) GetFilterChanges(id rpc.ID) (interface{}, error) {
	api.logger.Debug("GetFilterChanges", "id", id)

	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	if f, found := api.filters[id]; found {
		if !f.deadline.Stop() {
			// timer expired but filter is not yet removed in timeout loop
			// receive timer value and reset timer
			<-f.deadline.C
		}
		f.deadline.Reset(deadline)

		switch f.typ {
		case filters.BlocksSubscription:
			hashes := f.hashes
			f.hashes = nil
			return returnHashes(hashes), nil
		case filters.PendingTransactionsSubscription:
			if f.fullTx {
				txs := make([]*types.Transaction, 0, len(f.txs))
				txs = append(txs, f.txs...)
				f.txs = nil
				return txs, nil
			} else {
				hashes := make([]common.Hash, 0, len(f.txs))
				for _, tx := range f.txs {
					hashes = append(hashes, tx.Hash)
				}
				f.txs = nil
				return hashes, nil
			}
		case filters.LogsSubscription:
			logs := make([]*ethtypes.Log, len(f.logs))
			copy(logs, f.logs)
			f.logs = nil
			return returnLogs(logs), nil
		}
	}

	return []interface{}{}, errFilterNotFound
}

// Syncing provides information when this nodes starts synchronising with the OneLedger network and when it's finished.
func (api *PublicFilterAPI) Syncing() (*rpc.Subscription, error) {
	api.logger.Debug("Syncing")

	return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
}
