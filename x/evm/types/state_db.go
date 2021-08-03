package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"sync"
)

// CommitStateDB implements the Geth state.StateDB interface. Instead of using
// a trie and database for querying and persistence, the Keeper uses KVStores
// and an AccountKeeper to facilitate state transitions.
//
// TODO: This implementation is subject to change in regards to its statefull
// manner. In otherwords, how this relates to the keeper in this module.
type CommitStateDB struct {
	// TODO: We need to store the context as part of the structure itself opposed
	// to being passed as a parameter (as it should be) in order to implement the
	// StateDB interface. Perhaps there is a better way.
	ctx sdk.Context

	storeKey      sdk.StoreKey
	paramSpace    params.Subspace
	accountKeeper AccountKeeper

	// array that hold 'live' objects, which will get modified while processing a
	// state transition
	stateObjects         []stateEntry
	addressToObjectIndex map[sdk.AccAddress]int // map from address to the index of the state objects slice
	stateObjectsDirty    map[sdk.AccAddress]struct{}

	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash Hash
	txIndex      int
	logSize      uint

	// TODO: Determine if we actually need this as we do not need preimages in
	// the SDK, but it seems to be used elsewhere in Geth.
	preimages           []preimageEntry
	hashToPreimageIndex map[Hash]int // map from hash to the index of the preimages slice

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memo-ized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionID int

	// Per-transaction access list
	accessList *accessList

	// mutex for state deep copying
	lock sync.Mutex
}

// NewCommitStateDB returns a reference to a newly initialized CommitStateDB
// which implements Geth's state.StateDB interface.
//
// CONTRACT: Stores used for state must be cache-wrapped as the ordering of the
// key/value space matters in determining the merkle root.
func NewCommitStateDB(
	ctx sdk.Context, storeKey sdk.StoreKey, paramSpace params.Subspace, ak AccountKeeper,
) *CommitStateDB {
	return &CommitStateDB{
		ctx:                  ctx,
		storeKey:             storeKey,
		paramSpace:           paramSpace,
		accountKeeper:        ak,
		stateObjects:         []stateEntry{},
		addressToObjectIndex: make(map[ethcmn.Address]int),
		stateObjectsDirty:    make(map[ethcmn.Address]struct{}),
		preimages:            []preimageEntry{},
		hashToPreimageIndex:  make(map[Hash]int),
		journal:              newJournal(),
		accessList:           newAccessList(),
	}
}
