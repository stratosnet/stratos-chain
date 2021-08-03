package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// stateObject represents an Ethereum account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
// Finally, call CommitTrie to write the modified storage trie into a database.
type stateObject struct {
	code Code // contract bytecode, which gets set when code is loaded
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	originStorage Storage // Storage cache of original entries to dedup rewrites
	dirtyStorage  Storage // Storage entries that need to be flushed to disk

	// DB error
	dbErr   error
	stateDB *CommitStateDB
	account *auth.BaseAccount

	keyToOriginStorageIndex map[Hash]int
	keyToDirtyStorageIndex  map[Hash]int

	address sdk.AccAddress

	// cache flags
	//
	// When an object is marked suicided it will be delete from the trie during
	// the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

func newStateObject(db *CommitStateDB, account *auth.BaseAccount) *stateObject {
	//// func newStateObject(db *CommitStateDB, accProto authexported.Account, balance sdk.Int) *stateObject {
	//ethermintAccount, ok := accProto.(*ethermint.EthAccount)
	//if !ok {
	//	panic(fmt.Sprintf("invalid account type for state object: %T", accProto))
	//}
	//
	//// set empty code hash
	//if ethermintAccount.CodeHash == nil {
	//	ethermintAccount.CodeHash = emptyCodeHash
	//}

	return &stateObject{
		stateDB:                 db,
		account:                 account,
		address:                 account.GetAddress(),
		originStorage:           Storage{},
		dirtyStorage:            Storage{},
		keyToOriginStorageIndex: make(map[Hash]int),
		keyToDirtyStorageIndex:  make(map[Hash]int),
	}
}

// stateEntry represents a single key value pair from the StateDB's stateObject mappindg.
// This is to prevent non determinism at genesis initialization or export.
type stateEntry struct {
	// address key of the state object
	address     sdk.AccAddress
	stateObject *stateObject
}
