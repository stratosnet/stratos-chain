package filters

import (
	"context"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	tmpubsub "github.com/cometbft/cometbft/libs/pubsub"
	tmquery "github.com/cometbft/cometbft/libs/pubsub/query"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/stratosnet/stratos-chain/rpc/backend"
	"github.com/stratosnet/stratos-chain/rpc/types"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
)

var (
	// base tm tx event
	txEvents = tmtypes.QueryForEvent(tmtypes.EventTx).String()
	// listen only evm txs
	evmEvents = tmquery.MustCompile(fmt.Sprintf("%s='%s' AND %s.%s='%s'", tmtypes.EventTypeKey, tmtypes.EventTx, sdk.EventTypeMessage, sdk.AttributeKeyModule, evmtypes.ModuleName)).String()
	// specify to listen a blocks instead of heads as it contain more details about the block header including current block hash
	headerEvents = tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	// for tendermint subscribe channel capacity to make it buffered
	tmChannelCapacity = 100
	// consider a filter inactive if it has not been polled for within deadline
	deadline = 5 * time.Minute
	// The maximum number of topic criteria allowed, vm.LOG4 - vm.LOG0
	maxTopics = 4
)

// EventSystem creates subscriptions, processes events and broadcasts them to the
// subscription which match the subscription criteria using the Tendermint's RPC client.
type EventSystem struct {
	logger    log.Logger
	ctx       context.Context
	clientCtx client.Context
	backend   backend.BackendI

	// Unidirectional channels to receive Tendermint ResultEvents
	txsCh   <-chan coretypes.ResultEvent // Channel to receive new pending transactions event
	logsCh  <-chan coretypes.ResultEvent // Channel to receive new log event
	chainCh <-chan coretypes.ResultEvent // Channel to receive new chain event

	// Channels
	install   chan *subscription // install filter for event notification
	uninstall chan *subscription // remove filter for event notification
	eventBus  *tmtypes.EventBus
}

// NewEventSystem creates a new manager that listens for event on the given mux,
// parses and filters them. It uses the all map to retrieve filter changes. The
// work loop holds its own index that is used to forward events to filters.
//
// The returned manager has a loop that needs to be stopped with the Stop function
// or by stopping the given mux.
func NewEventSystem(clientCtx client.Context, logger log.Logger, eventBus *tmtypes.EventBus, b backend.BackendI) *EventSystem {
	es := &EventSystem{
		logger:    logger,
		ctx:       context.Background(),
		clientCtx: clientCtx,
		backend:   b,
		install:   make(chan *subscription),
		uninstall: make(chan *subscription),
		eventBus:  eventBus,
	}

	var (
		err                                                                           error
		cancelPendingTxsSubs, cancelLogsSubs, cancelPendingLogsSubs, cancelHeaderSubs context.CancelFunc
	)

	// Subscribe events
	es.txsCh, cancelPendingTxsSubs, err = es.subscribeTm(
		filters.PendingTransactionsSubscription,
		txEvents,
	)
	if err != nil {
		panic(fmt.Errorf("failed to subscribe pending txs: %w", err))
	}

	es.logsCh, cancelLogsSubs, err = es.subscribeTm(
		filters.LogsSubscription,
		evmEvents,
	)
	if err != nil {
		panic(fmt.Errorf("failed to subscribe logs: %w", err))
	}

	es.chainCh, cancelHeaderSubs, err = es.subscribeTm(
		filters.BlocksSubscription,
		headerEvents,
	)
	if err != nil {
		panic(fmt.Errorf("failed to subscribe headers: %w", err))
	}

	cancel := func() {
		cancelPendingTxsSubs()
		cancelLogsSubs()
		cancelPendingLogsSubs()
		cancelHeaderSubs()
	}

	go es.eventLoop(cancel)
	return es
}

// subscribeTm performs a new event subscription to a given Tendermint event.
// The subscription creates a unidirectional receive event channel to receive the ResultEvent. By
// default, the subscription timeouts (i.e is canceled) after 5 minutes. This function returns an
// error if the subscription fails (eg: if the identifier is already subscribed) or if the filter
// type is invalid.
func (es *EventSystem) subscribeTm(t filters.Type, evt string) (<-chan coretypes.ResultEvent, context.CancelFunc, error) {
	var (
		err      error
		cancelFn context.CancelFunc
		eventCh  <-chan coretypes.ResultEvent
	)

	es.ctx, cancelFn = context.WithTimeout(context.Background(), deadline)

	switch t {
	case
		filters.PendingTransactionsSubscription,
		filters.LogsSubscription,
		filters.BlocksSubscription:

		eventCh, err = es.createEventBusSubscription("web3_"+string(t), evt)
	default:
		err = fmt.Errorf("invalid filter subscription type %d", t)
	}

	if err != nil {
		return nil, cancelFn, err
	}

	return eventCh, cancelFn, nil
}

// subscribe installs the subscription in the event broadcast loop.
func (es *EventSystem) subscribe(sub *subscription) *Subscription {
	es.install <- sub
	<-sub.installed
	return &Subscription{id: sub.id, f: sub, es: es}
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
					es.logger.Error("wanted to publish ResultEvent, but out channel is full", "result", result, "query", result.Query)
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

		var (
			err error
			sub tmtypes.Subscription
		)
		if tmChannelCapacity > 0 {
			sub, err = es.eventBus.Subscribe(es.ctx, subscriber, q, tmChannelCapacity)
		} else {
			sub, err = es.eventBus.SubscribeUnbuffered(es.ctx, subscriber, q)
		}
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
func (es *EventSystem) SubscribeLogs(crit filters.FilterCriteria, logs chan []*ethtypes.Log) (*Subscription, error) {
	if len(crit.Topics) > maxTopics {
		return nil, errExceedMaxTopics
	}

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

	// Pending logs are not supported anymore.
	if from == rpc.PendingBlockNumber || to == rpc.PendingBlockNumber {
		return nil, errPendingLogsUnsupported
	}

	// only interested in new mined logs
	if from == rpc.LatestBlockNumber && to == rpc.LatestBlockNumber {
		return es.subscribeLogs(crit, logs), nil
	}

	// only interested in mined logs within a specific block range
	if from >= 0 && to >= 0 && to >= from {
		return es.subscribeLogs(crit, logs), nil
	}

	// interested in logs from a specific block number to new mined blocks
	if from >= 0 && to == rpc.LatestBlockNumber {
		return es.subscribeLogs(crit, logs), nil
	}

	return nil, errInvalidBlockRange
}

// subscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel.
func (es *EventSystem) subscribeLogs(crit filters.FilterCriteria, logs chan []*ethtypes.Log) *Subscription {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.LogsSubscription,
		event:     evmEvents,
		logsCrit:  crit,
		created:   time.Now().UTC(),
		logs:      logs,
		txs:       make(chan []*types.Transaction),
		headers:   make(chan *types.Header),
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return es.subscribe(sub)
}

// SubscribePendingTxs creates a subscription that writes transactions for
// transactions that enter the transaction pool.
func (es *EventSystem) SubscribePendingTxs(txs chan []*types.Transaction) *Subscription {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.PendingTransactionsSubscription,
		event:     evmEvents,
		created:   time.Now().UTC(),
		logs:      make(chan []*ethtypes.Log),
		txs:       txs,
		headers:   make(chan *types.Header),
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return es.subscribe(sub)
}

// SubscribeNewHeads subscribes to new block headers events.
func (es *EventSystem) SubscribeNewHeads(headers chan *types.Header) *Subscription {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.BlocksSubscription,
		event:     headerEvents,
		created:   time.Now().UTC(),
		logs:      make(chan []*ethtypes.Log),
		txs:       make(chan []*types.Transaction),
		headers:   headers,
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return es.subscribe(sub)
}

type filterIndex map[filters.Type]map[rpc.ID]*subscription

func (es *EventSystem) handleLogs(fls filterIndex, ev coretypes.ResultEvent) {
	_, isMsgEthereumTx := ev.Events[fmt.Sprintf("%s.%s", evmtypes.EventTypeEthereumTx, evmtypes.AttributeKeyEthereumTxHash)]
	if !isMsgEthereumTx {
		// >v011
		_, isMsgEthereumTx = ev.Events[fmt.Sprintf("%s.%s", "stratos.evm.v1.EventEthereumTx", "eth_hash")]
		if !isMsgEthereumTx {
			return
		}
	}

	data, _ := ev.Data.(tmtypes.EventDataTx)

	txResponse, err := evmtypes.DecodeTxResponse(data.TxResult.Result.Data)
	if err != nil {
		return
	}

	for _, f := range fls[filters.LogsSubscription] {
		matchedLogs := FilterLogs(evmtypes.LogsToEthereum(txResponse.Logs), f.logsCrit.FromBlock, f.logsCrit.ToBlock, f.logsCrit.Addresses, f.logsCrit.Topics)
		if len(matchedLogs) > 0 {
			f.logs <- matchedLogs
		}
	}
}

func (es *EventSystem) handleTxsEvent(fls filterIndex, ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataTx)
	height := uint64(data.Height)
	header := es.backend.CurrentHeader()
	if header == nil {
		es.logger.Error("hader not found during handleTxsEvent", "height", height)
		return
	}
	// NOTE: We do not know index at this time, left right now as 0
	index := uint64(0)
	for _, f := range fls[filters.PendingTransactionsSubscription] {
		tx, err := types.TmTxToEthTx(es.clientCtx.TxConfig.TxDecoder(), tmtypes.Tx(data.Tx), &header.Hash, &height, &index)
		if err != nil {
			continue
		}
		f.txs <- []*types.Transaction{tx}
	}
}

func (es *EventSystem) handleChainEvent(fls filterIndex, ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataNewBlockHeader)
	for _, f := range fls[filters.BlocksSubscription] {
		header, err := types.EthHeaderFromTendermint(data.Header)
		if err != nil {
			continue
		}
		// override dynamicly miner address
		sdkCtx, _, err := es.backend.GetEVMContext().GetSdkContextWithHeader(&data.Header)
		if err != nil {
			continue
		}
		validator, err := es.backend.GetEVMKeeper().GetCoinbaseAddress(sdkCtx)
		if err != nil {
			continue
		}
		header.Coinbase = validator
		f.headers <- header
	}
}

// eventLoop (un)installs filters and processes mux events.
func (es *EventSystem) eventLoop(cancel func()) {
	defer cancel()

	index := make(filterIndex)
	for i := filters.UnknownSubscription; i < filters.LastIndexSubscription; i++ {
		index[i] = make(map[rpc.ID]*subscription)
	}

	for {
		select {
		case txEvent := <-es.txsCh:
			es.handleTxsEvent(index, txEvent)
		case headerEv := <-es.chainCh:
			es.handleChainEvent(index, headerEv)
		case logsEv := <-es.logsCh:
			es.handleLogs(index, logsEv)

		case f := <-es.install:
			index[f.typ][f.id] = f
			close(f.installed)

		case f := <-es.uninstall:
			delete(index[f.typ], f.id)
			close(f.err)
		}
	}
}
