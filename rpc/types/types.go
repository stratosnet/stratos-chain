package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Copied the Account and StorageResult types since they are registered under an
// internal pkg on geth.

// AccountResult struct for account proof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

// StorageResult defines the format for storage proof return
type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        *common.Hash         `json:"blockHash"`
	BlockNumber      *hexutil.Big         `json:"blockNumber"`
	From             common.Address       `json:"from"`
	Gas              hexutil.Uint64       `json:"gas"`
	GasPrice         *hexutil.Big         `json:"gasPrice"`
	GasFeeCap        *hexutil.Big         `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big         `json:"maxPriorityFeePerGas,omitempty"`
	Hash             common.Hash          `json:"hash"`
	Input            hexutil.Bytes        `json:"input"`
	Nonce            hexutil.Uint64       `json:"nonce"`
	To               *common.Address      `json:"to"`
	TransactionIndex *hexutil.Uint64      `json:"transactionIndex"`
	Value            *hexutil.Big         `json:"value"`
	Type             hexutil.Uint64       `json:"type"`
	Accesses         *ethtypes.AccessList `json:"accessList,omitempty"`
	ChainID          *hexutil.Big         `json:"chainId,omitempty"`
	V                *hexutil.Big         `json:"v"`
	R                *hexutil.Big         `json:"r"`
	S                *hexutil.Big         `json:"s"`
}

// StateOverride is the collection of overridden accounts.
type StateOverride map[common.Address]OverrideAccount

// OverrideAccount indicates the overriding fields of account during the execution of
// a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

type FeeHistoryResult struct {
	OldestBlock  *hexutil.Big     `json:"oldestBlock"`
	Reward       [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee      []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64        `json:"gasUsedRatio"`
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes         `json:"raw"`
	Tx  *ethtypes.Transaction `json:"tx"`
}

type OneFeeHistory struct {
	BaseFee      *big.Int   // base fee  for each block
	Reward       []*big.Int // each element of the array will have the tip provided to miners for the percentile given
	GasUsedRatio float64    // the ratio of gas used to the gas limit for each block
}

// NOTE: Forked because tendermint block hash calculated in another way
// default ethereum take rlp from the struct
type Header struct {
	Hash        common.Hash         `json:"hash"             gencodec:"required"`
	ParentHash  common.Hash         `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash         `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address      `json:"miner"            gencodec:"required"`
	Root        common.Hash         `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash         `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash         `json:"receiptsRoot"     gencodec:"required"`
	Bloom       ethtypes.Bloom      `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int            `json:"difficulty"       gencodec:"required"`
	Number      *big.Int            `json:"number"           gencodec:"required"`
	GasLimit    uint64              `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64              `json:"gasUsed"          gencodec:"required"`
	Time        uint64              `json:"timestamp"        gencodec:"required"`
	Extra       []byte              `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash         `json:"mixHash"`
	Nonce       ethtypes.BlockNonce `json:"nonce"`

	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	// TODO: Add support
	// NOTE: Do we need this in real?
	// BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`
}

// Block represents a block returned to RPC clients.
type Block struct {
	Number           hexutil.Uint64      `json:"number"`
	Hash             common.Hash         `json:"hash"`
	ParentHash       common.Hash         `json:"parentHash"`
	Nonce            ethtypes.BlockNonce `json:"nonce"`
	Sha3Uncles       common.Hash         `json:"sha3Uncles"`
	LogsBloom        ethtypes.Bloom      `json:"logsBloom"`
	TransactionsRoot common.Hash         `json:"transactionsRoot"`
	StateRoot        common.Hash         `json:"stateRoot"`
	Miner            common.Address      `json:"miner"`
	MixHash          common.Hash         `json:"mixHash"`
	Difficulty       hexutil.Uint64      `json:"difficulty"`
	TotalDifficulty  hexutil.Uint64      `json:"totalDifficulty"`
	ExtraData        hexutil.Bytes       `json:"extraData"`
	Size             hexutil.Uint64      `json:"size"`
	GasLimit         *hexutil.Big        `json:"gasLimit"`
	GasUsed          *hexutil.Big        `json:"gasUsed"`
	Timestamp        hexutil.Uint64      `json:"timestamp"`
	Uncles           []common.Hash       `json:"uncles"`
	ReceiptsRoot     common.Hash         `json:"receiptsRoot"`
	Transactions     []interface{}       `json:"transactions"`
}

// TransactionReceipt represents a mined transaction returned to RPC clients.
type TransactionReceipt struct {
	// Consensus fields: These fields are defined by the Yellow Paper
	Status            hexutil.Uint64  `json:"status"`
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	LogsBloom         ethtypes.Bloom  `json:"logsBloom"`
	Logs              []*ethtypes.Log `json:"logs"`

	// Implementation fields: These fields are added by geth when processing a transaction.
	// They are stored in the chain database.
	TransactionHash common.Hash     `json:"transactionHash"`
	ContractAddress *common.Address `json:"contractAddress"`
	GasUsed         hexutil.Uint64  `json:"gasUsed"`

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        common.Hash    `json:"blockHash"`
	BlockNumber      hexutil.Big    `json:"blockNumber"`
	TransactionIndex hexutil.Uint64 `json:"transactionIndex"`

	// sender and receiver (contract or EOA) addresses
	From common.Address  `json:"from"`
	To   *common.Address `json:"to"`
}
