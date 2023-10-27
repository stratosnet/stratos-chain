package pool

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	mempl "github.com/cometbft/cometbft/mempool"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/stratosnet/stratos-chain/server/config"
	stratos "github.com/stratosnet/stratos-chain/types"
	"github.com/stratosnet/stratos-chain/x/evm"
	evmkeeper "github.com/stratosnet/stratos-chain/x/evm/keeper"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
)

const (
	// txSlotSize is used to calculate how many data slots a single transaction
	// takes up based on its size. The slots are used as DoS protection, ensuring
	// that validating a new transaction remains a constant operation (in reality
	// O(maxslots), where max slots are 4 currently).
	txSlotSize = 32 * 1024

	// txMaxSize is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions are significantly harder and
	// more expensive to propagate; larger transactions also take more resources
	// to validate whether they fit into the pool or not.
	txMaxSize = 4 * txSlotSize // 128KB
)

var (
	evictionInterval       = time.Minute     // Time interval to check for evictable transactions
	statsReportInterval    = 8 * time.Second // Time interval to report transaction pool stats
	processQueueInterval   = 5 * time.Second // Time interval to process queue transactions in case if their ready and move to pending stage
	processPendingInterval = 5 * time.Second // Time interval of pending tx broadcating to a tm mempool
)

// TxPool contains all currently known transactions. Transactions
// enter the pool when they are received from the local network or submitted
// locally. They exit the pool when they are included in the blockchain.
type TxPool struct {
	config      core.TxPoolConfig
	srvCfg      config.Config
	logger      log.Logger
	clientCtx   client.Context
	evmCtx      *evm.Context
	signerCache types.Signer // should be loaded in dynamic as potentially will not exist during chain initialization
	mempool     mempl.Mempool
	mu          sync.RWMutex
	imu         sync.Mutex // only for initializations

	evmkeeper        *evmkeeper.Keeper // Active keeper to get current state
	pendingNonces    *txNoncer         // Pending state tracking virtual nonces
	maxGasLimitCache uint64            // Current gas limit for transaction caps, only for cache

	beats map[common.Address]time.Time // Last heartbeat from each known account

	pending map[common.Address]*txList // All currently processable transactions
	queue   map[common.Address]*txList // Queued but non-processable transactions
	all     *txLookup                  // All transactions to allow lookups
}

// catchInitPanic used only in methods which are dependent from the store and potentially could be down
// due to not initialized chain
func catchInitPanic(err *error) {
	if errRec := recover(); errRec != nil {
		log.Warn("storage still not initialized, value could not be loaded")
		*err = errors.New("chain not started yet")
	}
}

// NewTxPool creates a new transaction pool to gather, sort and filter inbound
// transactions from the network.
func NewTxPool(config core.TxPoolConfig, srvCfg config.Config, clientCtx client.Context, mempool mempl.Mempool, evmkeeper *evmkeeper.Keeper, evmCtx *evm.Context) (*TxPool, error) {
	pool := &TxPool{
		config:           config,
		srvCfg:           srvCfg,
		clientCtx:        clientCtx,
		evmCtx:           evmCtx,
		mempool:          mempool,
		evmkeeper:        evmkeeper,
		signerCache:      nil,
		maxGasLimitCache: 0,
		beats:            make(map[common.Address]time.Time),
		pending:          make(map[common.Address]*txList),
		queue:            make(map[common.Address]*txList),
		all:              newTxLookup(),
	}
	pool.pendingNonces = newTxNoncer(pool.evmCtx, pool.evmkeeper)

	go pool.eventLoop()

	return pool, nil
}

// eventLoop starting a main logic of tx pool, orchaestrator of queues
func (pool *TxPool) eventLoop() {
	var (
		prevPending, prevQueued int
		// Start the stats reporting and transaction eviction tickers
		report         = time.NewTicker(statsReportInterval)
		evict          = time.NewTicker(evictionInterval)
		processQueue   = time.NewTicker(processQueueInterval)
		processPending = time.NewTicker(processPendingInterval)
	)
	defer report.Stop()
	defer evict.Stop()
	defer processQueue.Stop()
	defer processPending.Stop()

	for {
		select {
		// Handle stats reporting ticks
		// NOTE: Took some example from go ethereum, however it could be enabled only during for debuging,
		// so performance should be better without, but let's keep at least for now in order to watch how it will
		// work with tm pool
		case <-report.C:
			pool.mu.RLock()
			pending, queued := pool.stats()
			total := pool.all.Count()
			pool.mu.RUnlock()

			if pending != prevPending || queued != prevQueued {
				log.Debug("Transaction pool status report", "executable", pending, "queued", queued, "total", total)
				prevPending, prevQueued = pending, queued
			}

		// Handle inactive account transaction eviction
		case <-evict.C:
			pool.mu.Lock()
			for addr := range pool.queue {
				// Any non-locals old enough should be removed
				if time.Since(pool.beats[addr]) > pool.config.Lifetime {
					list := pool.queue[addr].Flatten()
					for _, tx := range list {
						pool.removeTx(tx.Hash())
					}
				}
			}
			pool.mu.Unlock()

		// Handle queue processing and moving to a next step
		case <-processQueue.C:
			pool.mu.Lock()
			pool.processQueue()
			pool.mu.Unlock()

		// Handle pending queue as a last process step before tm mempool execution
		case <-processPending.C:
			pool.mu.Lock()
			pool.processPending()
			pool.mu.Unlock()
		}
	}
}

// getSigner return signer for existing chain id
func (pool *TxPool) getSigner() (_ types.Signer, err error) {
	defer catchInitPanic(&err)

	if pool.signerCache == nil {
		pool.imu.Lock()
		defer pool.imu.Unlock()

		if pool.signerCache == nil {
			sdkCtx := pool.evmCtx.GetSdkContext()
			params := pool.evmkeeper.GetParams(sdkCtx)
			pool.signerCache = types.LatestSignerForChainID(params.ChainConfig.ChainID.BigInt())
		}
	}
	return pool.signerCache, nil
}

// getMaxGasLimit return current gas limit for transaction caps
func (pool *TxPool) getMaxGasLimit() (_ uint64, err error) {
	defer catchInitPanic(&err)

	if pool.maxGasLimitCache == 0 {
		pool.imu.Lock()
		defer pool.imu.Unlock()

		if pool.maxGasLimitCache == 0 {
			gasLimit, err := evmtypes.BlockMaxGasFromConsensusParams(nil)
			if err != nil {
				return 0, fmt.Errorf("failed to get tx pool current max gas: %w (possible DB not started?)", err)
			}
			pool.maxGasLimitCache = uint64(gasLimit)
		}
	}
	return pool.maxGasLimitCache, nil
}

// Get returns a transaction if it is contained in the pool and nil otherwise.
func (pool *TxPool) Get(hash common.Hash) *types.Transaction {
	return pool.all.Get(hash)
}

// Has returns an indicator whether txpool has a transaction cached with the
// given hash.
func (pool *TxPool) Has(hash common.Hash) bool {
	// checking first locals
	if tx := pool.all.Get(hash); tx != nil {
		return true
	}
	// attempt to db store
	if tx, _ := evmtypes.GetTmTxByHash(hash); tx != nil {
		return true
	}
	return false
}

func (pool *TxPool) MinGasPrice() *big.Int {
	sdkCtx := pool.evmCtx.GetSdkContext()
	params := pool.evmkeeper.GetParams(sdkCtx)
	minGasPrice := pool.srvCfg.GetMinGasPrices()
	amt := minGasPrice.AmountOf(params.EvmDenom).TruncateInt64()
	if amt == 0 {
		return new(big.Int).SetInt64(stratos.DefaultGasPrice)
	}
	return new(big.Int).SetInt64(amt)
}

// Nonce returns the next nonce of an account, with all transactions executable
// by the pool already applied on top.
func (pool *TxPool) Nonce(addr common.Address) uint64 {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.pendingNonces.get(addr)
}

// stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) stats() (int, int) {
	pending := 0
	for _, list := range pool.pending {
		pending += list.Len()
	}
	queued := 0
	for _, list := range pool.queue {
		queued += list.Len()
	}
	return pending, queued
}

// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (pool *TxPool) validateTx(tx *types.Transaction) error {
	log.Trace("Validating tx", "hash", tx.Hash())
	sdkCtx := pool.evmCtx.GetSdkContext()

	currentMaxGas, err := pool.getMaxGasLimit()
	if err != nil {
		return err
	}

	signer, err := pool.getSigner()
	if err != nil {
		return err
	}

	// Accept only legacy transactions until EIP-2718/2930 activates.
	if !tx.Protected() {
		return core.ErrTxTypeNotSupported
	}
	// Reject access list transactions until EIP-2930 activates.
	if tx.Type() == types.AccessListTxType {
		return core.ErrTxTypeNotSupported
	}
	// Reject transactions over defined size to prevent DOS attacks
	if uint64(tx.Size()) > txMaxSize {
		return core.ErrOversizedData
	}
	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return core.ErrNegativeValue
	}
	// Ensure the transaction doesn't exceed the current block limit gas.
	if currentMaxGas < tx.Gas() {
		return core.ErrGasLimit
	}
	// Sanity check for extremely large numbers
	if tx.GasFeeCap().BitLen() > 256 {
		return core.ErrFeeCapVeryHigh
	}
	if tx.GasTipCap().BitLen() > 256 {
		return core.ErrTipVeryHigh
	}
	// Ensure gasFeeCap is greater than or equal to gasTipCap.
	if tx.GasFeeCapIntCmp(tx.GasTipCap()) < 0 {
		return core.ErrTipAboveFeeCap
	}
	// Make sure the transaction is signed properly.
	from, err := types.Sender(signer, tx)
	if err != nil {
		return core.ErrInvalidSender
	}
	// Ensure the transaction adheres to nonce ordering
	if pool.evmkeeper.GetNonce(sdkCtx, from) > tx.Nonce() {
		return core.ErrNonceTooLow
	}
	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if pool.evmkeeper.GetBalance(sdkCtx, from).Cmp(tx.Cost()) < 0 {
		return core.ErrInsufficientFunds
	}
	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := core.IntrinsicGas(tx.Data(), tx.AccessList(), tx.To() == nil, true, true)
	if err != nil {
		return err
	}
	if tx.Gas() < intrGas {
		return core.ErrIntrinsicGas
	}
	return nil
}

// broadcastTx inserts a new transaction into the tendermint mempool and propagate
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) broadcastTx(tx *types.Transaction) error {
	log.Trace("Preparing to broadcast tx", "hash", tx.Hash())
	ethereumTx := &evmtypes.MsgEthereumTx{}
	if err := ethereumTx.FromEthereumTx(tx); err != nil {
		log.Trace("transaction converting failed", "error", err.Error())
		return err
	}

	params := pool.evmkeeper.GetParams(pool.evmCtx.GetSdkContext())
	cosmosTx, err := ethereumTx.BuildTx(pool.clientCtx.TxConfig.NewTxBuilder(), params.EvmDenom)
	if err != nil {
		log.Trace("failed to build cosmos tx", "error", err.Error())
		return err
	}

	// Encode transaction by default Tx encoder
	packet, err := pool.clientCtx.TxConfig.TxEncoder()(cosmosTx)
	if err != nil {
		log.Trace("failed to encode eth tx using default encoder", "error", err.Error())
		return err
	}

	err = pool.mempool.CheckTx(packet, nil, mempl.TxInfo{})
	if err != nil {
		log.Trace("failed to send eth tx packet to mempool", "error", err.Error())
		return err
	}
	return nil
}

// Add validates a transaction and inserts it into the non-executable queue for later
// pending promotion and execution. If the transaction is a replacement for an already
// pending or queued one, it overwrites the previous transaction if its price is higher.
func (pool *TxPool) Add(tx *types.Transaction) (replaced bool, err error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	hash := tx.Hash()

	if pool.Has(hash) {
		log.Trace("Discarding already known transaction", "hash", hash)
		return false, core.ErrAlreadyKnown
	}
	err = pool.validateTx(tx)
	if err != nil {
		return false, err
	}

	signer, _ := pool.getSigner()
	// Try to replace an existing transaction in the pending pool
	from, _ := types.Sender(signer, tx) // already validated
	if list := pool.pending[from]; list != nil && list.Overlaps(tx) {
		log.Trace("Tx overlap", "hash", tx.Hash())
		// Nonce already pending, check if required price bump is met
		inserted, old := list.Add(tx, pool.config.PriceBump)
		if !inserted {
			return false, core.ErrReplaceUnderpriced
		}
		// New transaction is better, replace old one
		if old != nil {
			pool.all.Remove(old.Hash())
		}
		pool.all.Add(tx)
		// Successful promotion, bump the heartbeat
		pool.beats[from] = time.Now()
		return old != nil, nil
	}

	return pool.enqueueTx(tx, true)
}

// enqueueTx inserts a new transaction into the non-executable transaction queue.
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) enqueueTx(tx *types.Transaction, addAll bool) (bool, error) {
	log.Trace("Tx enqeue", "hash", tx.Hash())

	signer, _ := pool.getSigner()
	from, _ := types.Sender(signer, tx) // already validated
	if pool.queue[from] == nil {
		pool.queue[from] = newTxList(false)
	}
	inserted, old := pool.queue[from].Add(tx, pool.config.PriceBump)
	if !inserted {
		// An older transaction was better, discard this
		return false, core.ErrReplaceUnderpriced
	}
	if addAll {
		pool.all.Add(tx)
	}
	// If we never record the heartbeat, do it right now.
	if _, exist := pool.beats[from]; !exist {
		pool.beats[from] = time.Now()
	}

	// Set the potentially new pending nonce and notify any subsystems of the new tx
	pool.pendingNonces.set(from, tx.Nonce()+1)

	return old != nil, nil
}

// removeTx removes a single transaction from the queue, moving all subsequent
// transactions back to the future queue.
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) removeTx(hash common.Hash) {
	log.Error("Tx remove", "hash", hash)
	// Fetch the transaction we wish to delete
	tx := pool.Get(hash)
	if tx == nil {
		log.Error("Tx not found during removal", "hash", hash)
		return
	}
	signer, _ := pool.getSigner()
	addr, _ := types.Sender(signer, tx) // already validated during insertion

	// Remove it from the list of known transactions
	pool.all.Remove(hash)
	// Remove the transaction from the pending lists and reset the account nonce
	if pending := pool.pending[addr]; pending != nil {
		if removed, invalids := pending.Remove(tx); removed {
			// If no more pending transactions are left, remove the list
			if pending.Empty() {
				delete(pool.pending, addr)
			}
			// Postpone any invalidated transactions
			for _, tx := range invalids {
				// Internal shuffle shouldn't touch the lookup set.
				pool.enqueueTx(tx, false)
			}
			// Update the account nonce if needed
			pool.pendingNonces.setIfLower(addr, tx.Nonce())
			return
		}
	}
	// Transaction is in the future queue
	if future := pool.queue[addr]; future != nil {
		future.Remove(tx)
		log.Error("Removing tx from future list", "hash", hash, "addr", addr, "future", future)
		if future.Empty() {
			log.Error("Futures empty, cleaning", "hash", hash, "addr", addr)
			delete(pool.queue, addr)
			delete(pool.beats, addr)
		}
	}
}

func (pool *TxPool) processQueue() error {
	promoteAddrs := make([]common.Address, 0, len(pool.queue))
	for addr := range pool.queue {
		promoteAddrs = append(promoteAddrs, addr)
	}
	pool.promoteExecutables(promoteAddrs)
	pool.demoteUnexecutables()

	// Update all accounts to the latest known pending nonce
	nonces := make(map[common.Address]uint64, len(pool.pending))
	for addr, list := range pool.pending {
		highestPending := list.LastElement()
		nonces[addr] = highestPending.Nonce() + 1
	}
	pool.pendingNonces.setAll(nonces)
	return nil
}

func (pool *TxPool) processPending() error {
	runnerAddrs := make([]common.Address, 0, len(pool.pending))
	for addr := range pool.pending {
		runnerAddrs = append(runnerAddrs, addr)
	}

	for _, addr := range runnerAddrs {
		list := pool.pending[addr]
		if list == nil {
			continue // Just in case someone calls with a non existing account
		}
		readies := list.Ready(pool.pendingNonces.get(addr))
		for _, tx := range readies {
			err := pool.broadcastTx(tx)
			if err != nil {
				pool.logger.Error("Broadcast failed", "error", err)
			}
			pool.all.Remove(tx.Hash())
		}
		log.Trace("Processed pending transactions", "count", len(readies))

		if list.Empty() {
			delete(pool.pending, addr)
		}
	}

	return nil
}

// promoteExecutables moves transactions that have become processable from the
// future queue to the set of pending transactions. During this process, all
// invalidated transactions (low nonce, low balance) are deleted.
func (pool *TxPool) promoteExecutables(accounts []common.Address) []*types.Transaction {
	// Track the promoted transactions to broadcast them at once
	var promoted []*types.Transaction

	sdkCtx := pool.evmCtx.GetSdkContext()
	currentMaxGas, _ := pool.getMaxGasLimit()

	// Iterate over all accounts and promote any executable transactions
	for _, addr := range accounts {
		list := pool.queue[addr]
		if list == nil {
			continue // Just in case someone calls with a non existing account
		}
		// Drop all transactions that are deemed too old (low nonce)
		forwards := list.Forward(pool.evmkeeper.GetNonce(sdkCtx, addr))
		for _, tx := range forwards {
			hash := tx.Hash()
			pool.all.Remove(hash)
		}
		log.Trace("Removed old queued transactions", "count", len(forwards))
		// Drop all transactions that are too costly (low balance or out of gas)
		drops, _ := list.Filter(pool.evmkeeper.GetBalance(sdkCtx, addr), currentMaxGas)
		for _, tx := range drops {
			hash := tx.Hash()
			pool.all.Remove(hash)
		}
		log.Trace("Removed unpayable queued transactions", "count", len(drops))

		// Gather all executable transactions and promote them
		readies := list.Ready(pool.pendingNonces.get(addr))
		for _, tx := range readies {
			if pool.promoteTx(addr, tx) {
				promoted = append(promoted, tx)
			}
		}
		log.Trace("Promoted queued transactions", "count", len(promoted))

		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.queue, addr)
			delete(pool.beats, addr)
		}
	}
	return promoted
}

// demoteUnexecutables removes invalid and processed transactions from the pools
// executable/pending queue and any subsequent transactions that become unexecutable
// are moved back into the future queue.
//
// Note: transactions are not marked as removed in the priced list because re-heaping
// is always explicitly triggered by SetBaseFee and it would be unnecessary and wasteful
// to trigger a re-heap is this function
func (pool *TxPool) demoteUnexecutables() {
	sdkCtx := pool.evmCtx.GetSdkContext()
	currentMaxGas, _ := pool.getMaxGasLimit()

	// Iterate over all accounts and demote any non-executable transactions
	for addr, list := range pool.pending {
		nonce := pool.evmkeeper.GetNonce(sdkCtx, addr)

		// Drop all transactions that are deemed too old (low nonce)
		olds := list.Forward(nonce)
		for _, tx := range olds {
			hash := tx.Hash()
			pool.all.Remove(hash)
			log.Trace("Removed old pending transaction", "hash", hash)
		}
		// Drop all transactions that are too costly (low balance or out of gas), and queue any invalids back for later
		drops, invalids := list.Filter(pool.evmkeeper.GetBalance(sdkCtx, addr), currentMaxGas)
		for _, tx := range drops {
			hash := tx.Hash()
			log.Trace("Removed unpayable pending transaction", "hash", hash)
			pool.all.Remove(hash)
		}

		for _, tx := range invalids {
			hash := tx.Hash()
			log.Trace("Demoting pending transaction", "hash", hash)

			// Internal shuffle shouldn't touch the lookup set.
			pool.enqueueTx(tx, false)
		}
		// If there's a gap in front, alert (should never happen) and postpone all transactions
		if list.Len() > 0 && list.txs.Get(nonce) == nil {
			gapped := list.Cap(0)
			for _, tx := range gapped {
				hash := tx.Hash()
				log.Trace("Demoting invalidated transaction", "hash", hash)

				// Internal shuffle shouldn't touch the lookup set.
				pool.enqueueTx(tx, false)
			}
		}
		// Delete the entire pending entry if it became empty.
		if list.Empty() {
			delete(pool.pending, addr)
		}
	}
}

// promoteTx adds a transaction to the pending (processable) list of transactions
// and returns whether it was inserted or an older was better.
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) promoteTx(addr common.Address, tx *types.Transaction) bool {
	log.Trace("Preparing to promote tx", "hash", tx.Hash())
	// Try to insert the transaction into the pending queue
	if pool.pending[addr] == nil {
		pool.pending[addr] = newTxList(true)
	}
	list := pool.pending[addr]

	inserted, old := list.Add(tx, pool.config.PriceBump)
	if !inserted {
		// An older transaction was better, discard this
		pool.all.Remove(tx.Hash())
		return false
	}
	// Otherwise discard any previous transaction and mark this
	if old != nil {
		pool.all.Remove(old.Hash())
	}

	// Successful promotion, bump the heartbeat
	pool.beats[addr] = time.Now()
	return true
}
