package filters

import (
	"sync"
	"time"

	"github.com/stratosnet/stratos-chain/rpc/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

// Subscription defines a wrapper for the private subscription
type subscription struct {
	id        rpc.ID
	typ       filters.Type
	event     string
	created   time.Time
	logsCrit  filters.FilterCriteria
	logs      chan []*ethtypes.Log
	txs       chan []*types.Transaction
	headers   chan *types.Header
	installed chan struct{} // closed when the filter is installed
	err       chan error
}

// Subscription defines a wrapper for the private subscription
type Subscription struct {
	id        rpc.ID
	f         *subscription
	es        *EventSystem
	unsubOnce sync.Once
}

// ID returns the underlying subscription RPC identifier.
func (s *Subscription) ID() rpc.ID {
	return s.id
}

// Unsubscribe from the current subscription to Tendermint Websocket. It sends an error to the
// subscription error channel if unsubscribe fails.
func (s *Subscription) Unsubscribe() {
	s.unsubOnce.Do(func() {
	uninstallLoop:
		for {
			// write uninstall request and consume logs/hashes. This prevents
			// the eventLoop broadcast method to deadlock when writing to the
			// filter event channel while the subscription loop is waiting for
			// this method to return (and thus not reading these events).
			select {
			case s.es.uninstall <- s.f:
				break uninstallLoop
			case <-s.f.logs:
			case <-s.f.txs:
			case <-s.f.headers:
			}
		}

		// wait for filter to be uninstalled in work loop before returning
		// this ensures that the manager won't use the event channel which
		// will probably be closed by the client asap after this method returns.
		<-s.Err()
	})
}

// Err returns the error channel
func (s *Subscription) Err() <-chan error {
	return s.f.err
}
