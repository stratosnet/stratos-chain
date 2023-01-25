package filters

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/tendermint/tendermint/libs/log"
	tmpubsub "github.com/tendermint/tendermint/libs/pubsub"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/stratosnet/stratos-chain/rpc/types"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
)

var (
	// base tm tx event
	txEvents = tmtypes.QueryForEvent(tmtypes.EventTx).String()
	// listen only evm txs
	evmEvents = tmquery.MustParse(fmt.Sprintf("%s='%s' AND %s.%s='%s'", tmtypes.EventTypeKey, tmtypes.EventTx, sdk.EventTypeMessage, sdk.AttributeKeyModule, evmtypes.ModuleName)).String()
	// specify to listen a blocks instead of heads as it contain more details about the block header including current block hash
	headerEvents = tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	// for tendermint subscribe channel capacity to make it buffered
	tmChannelCapacity = 100
)

// EventSystem creates subscriptions, processes events and broadcasts them to the
// subscription which match the subscription criteria using the Tendermint's RPC client.
type EventSystem struct {
	logger     log.Logger
	ctx        context.Context
	tmWSClient *rpcclient.WSClient

	index      filterIndex
	topicChans map[string]chan<- coretypes.ResultEvent
	indexMux   *sync.RWMutex

	// Subscriptions
	txsSub         *Subscription // Subscription for new transaction event
	logsSub        *Subscription // Subscription for new log event
	pendingLogsSub *Subscription // Subscription for pending log event
	chainSub       *Subscription // Subscription for new chain event

	// Unidirectional channels to receive Tendermint ResultEvents
	txsCh         <-chan coretypes.ResultEvent // Channel to receive new pending transactions event
	logsCh        <-chan coretypes.ResultEvent // Channel to receive new log event
	pendingLogsCh <-chan coretypes.ResultEvent // Channel to receive new pending log event
	chainCh       <-chan coretypes.ResultEvent // Channel to receive new chain event

	// Channels
	install   chan *Subscription // install filter for event notification
	uninstall chan *Subscription // remove filter for event notification
	eventBus  *tmtypes.EventBus
}

// NewEventSystem creates a new manager that listens for event on the given mux,
// parses and filters them. It uses the all map to retrieve filter changes. The
// work loop holds its own index that is used to forward events to filters.
//
// The returned manager has a loop that needs to be stopped with the Stop function
// or by stopping the given mux.
func NewEventSystem(logger log.Logger, eventBus *tmtypes.EventBus) *EventSystem {
	index := make(filterIndex)
	for i := filters.UnknownSubscription; i < filters.LastIndexSubscription; i++ {
		index[i] = make(map[rpc.ID]*Subscription)
	}

	es := &EventSystem{
		logger:        logger,
		ctx:           context.Background(),
		index:         index,
		topicChans:    make(map[string]chan<- coretypes.ResultEvent, len(index)),
		indexMux:      new(sync.RWMutex),
		install:       make(chan *Subscription),
		uninstall:     make(chan *Subscription),
		eventBus:      eventBus,
		txsCh:         make(<-chan coretypes.ResultEvent),
		logsCh:        make(<-chan coretypes.ResultEvent),
		pendingLogsCh: make(<-chan coretypes.ResultEvent),
		chainCh:       make(<-chan coretypes.ResultEvent),
	}

	go es.eventLoop()
	return es
}

// WithContext sets a new context to the EventSystem. This is required to set a timeout context when
// a new filter is intantiated.
func (es *EventSystem) WithContext(ctx context.Context) {
	es.ctx = ctx
}

// subscribe performs a new event subscription to a given Tendermint event.
// The subscription creates a unidirectional receive event channel to receive the ResultEvent. By
// default, the subscription timeouts (i.e is canceled) after 5 minutes. This function returns an
// error if the subscription fails (eg: if the identifier is already subscribed) or if the filter
// type is invalid.
func (es *EventSystem) subscribe(sub *Subscription) (*Subscription, context.CancelFunc, error) {
	var (
		err      error
		cancelFn context.CancelFunc
		eventCh  <-chan coretypes.ResultEvent
	)

	es.ctx, cancelFn = context.WithTimeout(context.Background(), deadline)

	switch sub.typ {
	case
		filters.PendingTransactionsSubscription,
		filters.PendingLogsSubscription,
		filters.MinedAndPendingLogsSubscription,
		filters.LogsSubscription,
		filters.BlocksSubscription:

		eventCh, err = es.createEventBusSubscription(string(sub.id), sub.event)
	default:
		err = fmt.Errorf("invalid filter subscription type %d", sub.typ)
	}

	if err != nil {
		sub.err <- err
		return nil, cancelFn, err
	}

	// wrap events in a go routine to prevent blocking
	go func() {
		es.install <- sub
		<-sub.installed
	}()

	sub.eventCh = eventCh
	return sub, cancelFn, nil
}

func (es *EventSystem) createEventBusSubscription(subscriber, query string) (out <-chan coretypes.ResultEvent, err error) {
	q, err := tmquery.New(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse query")
	}

	var sub tmtypes.Subscription
	if tmChannelCapacity > 0 {
		sub, err = es.eventBus.Subscribe(es.ctx, subscriber, q, tmChannelCapacity)
	} else {
		sub, err = es.eventBus.SubscribeUnbuffered(es.ctx, subscriber, q)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	outc := make(chan coretypes.ResultEvent, tmChannelCapacity)
	go es.eventsRoutine(sub, subscriber, q, outc)

	return outc, nil
}

func (es *EventSystem) eventsRoutine(
	sub tmtypes.Subscription,
	subscriber string,
	q tmpubsub.Query,
	outc chan<- coretypes.ResultEvent) {
	for {
		select {
		case msg := <-sub.Out():
			result := coretypes.ResultEvent{Query: q.String(), Data: msg.Data(), Events: msg.Events()}
			if cap(outc) == 0 {
				outc <- result
			} else {
				select {
				case outc <- result:
				default:
					// es.logger.Error("wanted to publish ResultEvent, but out channel is full", "result", result, "query", result.Query)
				}
			}
		case <-sub.Cancelled():
			if sub.Err() == tmpubsub.ErrUnsubscribed {
				return
			}

			es.logger.Error("subscription was cancelled, resubscribing...", "err", sub.Err(), "query", q.String())
			sub = es.resubscribe(subscriber, q)
			if sub == nil { // client was stopped
				return
			}
		case <-es.eventBus.Quit():
			return
		}
	}
}

// Try to resubscribe with exponential backoff.
func (es *EventSystem) resubscribe(subscriber string, q tmpubsub.Query) tmtypes.Subscription {
	attempts := 0
	for {
		if !es.eventBus.IsRunning() {
			return nil
		}

		sub, err := es.eventBus.Subscribe(context.Background(), subscriber, q)
		if err == nil {
			return sub
		}

		attempts++
		time.Sleep((10 << uint(attempts)) * time.Millisecond) // 10ms -> 20ms -> 40ms
	}
}

// SubscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel. Default value for the from and to
// block is "latest". If the fromBlock > toBlock an error is returned.
func (es *EventSystem) SubscribeLogs(crit filters.FilterCriteria) (*Subscription, context.CancelFunc, error) {
	var from, to rpc.BlockNumber
	if crit.FromBlock == nil {
		from = rpc.LatestBlockNumber
	} else {
		from = rpc.BlockNumber(crit.FromBlock.Int64())
	}
	if crit.ToBlock == nil {
		to = rpc.LatestBlockNumber
	} else {
		to = rpc.BlockNumber(crit.ToBlock.Int64())
	}

	switch {
	// only interested in pending logs
	case from == rpc.PendingBlockNumber && to == rpc.PendingBlockNumber:
		return es.subscribePendingLogs(crit)

	// only interested in new mined logs, mined logs within a specific block range, or
	// logs from a specific block number to new mined blocks
	case (from == rpc.LatestBlockNumber && to == rpc.LatestBlockNumber),
		(from >= 0 && to >= 0 && to >= from):
		return es.subscribeLogs(crit)

	// interested in mined logs from a specific block number, new logs and pending logs
	case from >= rpc.LatestBlockNumber && (to == rpc.PendingBlockNumber || to == rpc.LatestBlockNumber):
		return es.subscribeMinedPendingLogs(crit)

	default:
		return nil, nil, fmt.Errorf("invalid from and to block combination: from > to (%d > %d)", from, to)
	}
}

// subscribePendingLogs creates a subscription that writes transaction hashes for
// transactions that enter the transaction pool.
func (es *EventSystem) subscribePendingLogs(crit filters.FilterCriteria) (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.PendingLogsSubscription,
		event:     evmEvents,
		logsCrit:  crit,
		created:   time.Now().UTC(),
		logs:      make(chan []*ethtypes.Log),
		installed: make(chan struct{}, 1),
		err:       make(chan error, 1),
	}
	return es.subscribe(sub)
}

// subscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel.
func (es *EventSystem) subscribeLogs(crit filters.FilterCriteria) (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.LogsSubscription,
		event:     evmEvents,
		logsCrit:  crit,
		created:   time.Now().UTC(),
		logs:      make(chan []*ethtypes.Log),
		installed: make(chan struct{}, 1),
		err:       make(chan error, 1),
	}
	return es.subscribe(sub)
}

// subscribeMinedPendingLogs creates a subscription that returned mined and
// pending logs that match the given criteria.
func (es *EventSystem) subscribeMinedPendingLogs(crit filters.FilterCriteria) (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.MinedAndPendingLogsSubscription,
		event:     evmEvents,
		logsCrit:  crit,
		created:   time.Now().UTC(),
		logs:      make(chan []*ethtypes.Log),
		installed: make(chan struct{}, 1),
		err:       make(chan error, 1),
	}
	return es.subscribe(sub)
}

// SubscribeNewHeads subscribes to new block headers events.
func (es EventSystem) SubscribeNewHeads() (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.BlocksSubscription,
		event:     headerEvents,
		created:   time.Now().UTC(),
		headers:   make(chan *types.Header),
		installed: make(chan struct{}, 1),
		err:       make(chan error, 1),
	}
	return es.subscribe(sub)
}

// SubscribePendingTxs subscribes to new pending transactions events from the mempool.
func (es EventSystem) SubscribePendingTxs() (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.PendingTransactionsSubscription,
		event:     txEvents,
		created:   time.Now().UTC(),
		hashes:    make(chan []common.Hash),
		installed: make(chan struct{}, 1),
		err:       make(chan error, 1),
	}
	return es.subscribe(sub)
}

type filterIndex map[filters.Type]map[rpc.ID]*Subscription

func (es *EventSystem) handleLogs(ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataTx)
	// logReceipt := onetypes.GetTxEthLogs(&data.TxResult.Result, data.Index)
	resultData, err := evmtypes.DecodeTransactionLogs(data.TxResult.Result.Data)
	if err != nil {
		return
	}

	if len(resultData.Logs) == 0 {
		return
	}

	for _, f := range es.index[filters.LogsSubscription] {
		matchedLogs := FilterLogs(evmtypes.LogsToEthereum(resultData.Logs), f.logsCrit.FromBlock, f.logsCrit.ToBlock, f.logsCrit.Addresses, f.logsCrit.Topics)
		if len(matchedLogs) > 0 {
			f.logs <- matchedLogs
		}
	}
}

func (es *EventSystem) handleTxsEvent(ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataTx)
	for _, f := range es.index[filters.PendingTransactionsSubscription] {
		// NOTE: In previous version, data.Tx return types.Tx, but right now just bytes,
		// so we need manually covert in order to get sha256 hash fot tx payload
		f.hashes <- []common.Hash{common.BytesToHash(tmtypes.Tx(data.Tx).Hash())}
	}
}

func (es *EventSystem) handleChainEvent(ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataNewBlockHeader)
	for _, f := range es.index[filters.BlocksSubscription] {
		header, err := types.EthHeaderFromTendermint(data.Header)
		if err != nil {
			continue
		}
		f.headers <- header
	}
}

// eventLoop (un)installs filters and processes mux events.
func (es *EventSystem) eventLoop() {
	var (
		err                                                                           error
		cancelPendingTxsSubs, cancelLogsSubs, cancelPendingLogsSubs, cancelHeaderSubs context.CancelFunc
	)

	// Subscribe events
	es.txsSub, cancelPendingTxsSubs, err = es.SubscribePendingTxs()
	if err != nil {
		panic(fmt.Errorf("failed to subscribe pending txs: %w", err))
	}

	defer cancelPendingTxsSubs()

	es.logsSub, cancelLogsSubs, err = es.SubscribeLogs(filters.FilterCriteria{})
	if err != nil {
		panic(fmt.Errorf("failed to subscribe logs: %w", err))
	}

	defer cancelLogsSubs()

	es.pendingLogsSub, cancelPendingLogsSubs, err = es.subscribePendingLogs(filters.FilterCriteria{})
	if err != nil {
		panic(fmt.Errorf("failed to subscribe pending logs: %w", err))
	}

	defer cancelPendingLogsSubs()

	es.chainSub, cancelHeaderSubs, err = es.SubscribeNewHeads()
	if err != nil {
		panic(fmt.Errorf("failed to subscribe headers: %w", err))
	}

	defer cancelHeaderSubs()

	// Ensure all subscriptions get cleaned up
	defer func() {
		es.txsSub.Unsubscribe(es)
		es.logsSub.Unsubscribe(es)
		es.pendingLogsSub.Unsubscribe(es)
		es.chainSub.Unsubscribe(es)
	}()

	for {
		select {
		case txEvent := <-es.txsSub.eventCh:
			es.handleTxsEvent(txEvent)
		case headerEv := <-es.chainSub.eventCh:
			es.handleChainEvent(headerEv)
		case logsEv := <-es.logsSub.eventCh:
			es.handleLogs(logsEv)
		case logsEv := <-es.pendingLogsSub.eventCh:
			es.handleLogs(logsEv)

		case f := <-es.install:
			if f.typ == filters.MinedAndPendingLogsSubscription {
				// the type are logs and pending logs subscriptions
				es.index[filters.LogsSubscription][f.id] = f
				es.index[filters.PendingLogsSubscription][f.id] = f
			} else {
				es.index[f.typ][f.id] = f
			}
			close(f.installed)

		case f := <-es.uninstall:
			if f.typ == filters.MinedAndPendingLogsSubscription {
				// the type are logs and pending logs subscriptions
				delete(es.index[filters.LogsSubscription], f.id)
				delete(es.index[filters.PendingLogsSubscription], f.id)
			} else {
				delete(es.index[f.typ], f.id)
			}
			close(f.err)
			// System stopped
		case <-es.txsSub.Err():
			return
		case <-es.logsSub.Err():
			return
		case <-es.pendingLogsSub.Err():
			return
		case <-es.chainSub.Err():
			return
		}
	}
}
