package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Web3Msg struct {
	From  *common.Address
	To    *common.Address
	Value *big.Int
}

type Web3MsgType interface {
	GetWeb3Msg() (*Web3Msg, error)
}
