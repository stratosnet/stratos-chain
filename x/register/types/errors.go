package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid                  = sdkerrors.Register(ModuleName, 1, "custom error message")
	ErrBadResourceNodeAddr      = sdkerrors.Register(ModuleName, 2, "resource node address is invalid")
	ErrInsufficientShares       = sdkerrors.Register(ModuleName, 3, "insufficient delegation shares")
	ErrResourceNodeOwnerExists  = sdkerrors.Register(ModuleName, 4, "resource node already exist for this operator address; must use new resource node operator address")
	ErrResourceNodePubKeyExists = sdkerrors.Register(ModuleName, 5, "resource node already exist for this pubkey; must use new resource node pubkey")
	ErrIndexingNodeOwnerExists  = sdkerrors.Register(ModuleName, 4, "indexing node already exist for this operator address; must use new indexing node operator address")
	ErrIndexingNodePubKeyExists = sdkerrors.Register(ModuleName, 5, "indexing node already exist for this pubkey; must use new indexing node pubkey")
)
