package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/mempool"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	stratos "github.com/stratosnet/stratos-chain/types"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"
	tmrpccoretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// RawTxToEthTx returns a evm MsgEthereum transaction from raw tx bytes.
func RawTxToEthTx(clientCtx client.Context, txBz tmtypes.Tx) ([]*evmtypes.MsgEthereumTx, error) {
	tx, err := clientCtx.TxConfig.TxDecoder()(txBz)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	ethTxs := make([]*evmtypes.MsgEthereumTx, len(tx.GetMsgs()))
	for i, msg := range tx.GetMsgs() {
		ethTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return nil, fmt.Errorf("invalid message type %T, expected %T", msg, &evmtypes.MsgEthereumTx{})
		}
		ethTxs[i] = ethTx
	}
	return ethTxs, nil
}

// EthHeaderFromTendermint is an util function that returns an Ethereum Header
// from a tendermint Header.
func EthHeaderFromTendermint(header tmtypes.Header, bloom ethtypes.Bloom, baseFee *big.Int) *ethtypes.Header {
	txHash := ethtypes.EmptyRootHash
	if len(header.DataHash) == 0 {
		txHash = common.BytesToHash(header.DataHash)
	}

	return &ethtypes.Header{
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		UncleHash:   ethtypes.EmptyUncleHash,
		Coinbase:    common.BytesToAddress(header.ProposerAddress),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      txHash,
		ReceiptHash: ethtypes.EmptyRootHash,
		Bloom:       bloom,
		Difficulty:  big.NewInt(0),
		Number:      big.NewInt(header.Height),
		GasLimit:    0,
		GasUsed:     0,
		Time:        uint64(header.Time.UTC().Unix()),
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       ethtypes.BlockNonce{},
		BaseFee:     baseFee,
	}
}

// BlockMaxGasFromConsensusParams returns the gas limit for the current block from the chain consensus params.
func BlockMaxGasFromConsensusParams(clientCtx client.Context, blockHeight int64) (int64, error) {
	resConsParams, err := tmrpccore.ConsensusParams(nil, &blockHeight)
	if err != nil {
		return int64(^uint32(0)), err
	}

	gasLimit := resConsParams.ConsensusParams.Block.MaxGas
	if gasLimit == -1 {
		// Sets gas limit to max uint32 to not error with javascript dev tooling
		// This -1 value indicating no block gas limit is set to max uint64 with geth hexutils
		// which errors certain javascript dev tooling which only supports up to 53 bits
		gasLimit = int64(^uint32(0))
	}

	return gasLimit, nil
}

// FormatBlock creates an ethereum block from a tendermint header and ethereum-formatted
// transactions.
func FormatBlock(
	header tmtypes.Header, size int, gasLimit int64,
	gasUsed *big.Int, transactions []interface{}, bloom ethtypes.Bloom,
	validatorAddr common.Address, baseFee *big.Int,
) map[string]interface{} {
	var transactionsRoot common.Hash
	if len(transactions) == 0 {
		transactionsRoot = ethtypes.EmptyRootHash
	} else {
		transactionsRoot = common.BytesToHash(header.DataHash)
	}

	result := map[string]interface{}{
		"number":           hexutil.Uint64(header.Height),
		"hash":             hexutil.Bytes(header.Hash()),
		"parentHash":       common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		"nonce":            ethtypes.BlockNonce{},   // PoW specific
		"sha3Uncles":       ethtypes.EmptyUncleHash, // No uncles in Tendermint
		"logsBloom":        bloom,
		"stateRoot":        hexutil.Bytes(header.AppHash),
		"miner":            validatorAddr,
		"mixHash":          common.Hash{},
		"difficulty":       (*hexutil.Big)(big.NewInt(0)),
		"extraData":        "0x",
		"size":             hexutil.Uint64(size),
		"gasLimit":         hexutil.Uint64(gasLimit), // Static gas limit
		"gasUsed":          (*hexutil.Big)(gasUsed),
		"timestamp":        hexutil.Uint64(header.Time.Unix()),
		"transactionsRoot": transactionsRoot,
		"receiptsRoot":     ethtypes.EmptyRootHash,

		"uncles":          []common.Hash{},
		"transactions":    transactions,
		"totalDifficulty": (*hexutil.Big)(big.NewInt(0)),
	}

	if baseFee != nil {
		result["baseFeePerGas"] = (*hexutil.Big)(baseFee)
	}

	return result
}

type DataError interface {
	Error() string          // returns the message
	ErrorData() interface{} // returns the error data
}

type dataError struct {
	msg  string
	data string
}

func (d *dataError) Error() string {
	return d.msg
}

func (d *dataError) ErrorData() interface{} {
	return d.data
}

type SDKTxLogs struct {
	Log string `json:"log"`
}

const LogRevertedFlag = "transaction reverted"

func ErrRevertedWith(data []byte) DataError {
	return &dataError{
		msg:  "VM execution error.",
		data: fmt.Sprintf("0x%s", hex.EncodeToString(data)),
	}
}

// GetPendingTxCountByAddress is used to get pending tx count (nonce) for user address
func GetPendingTxCountByAddress(txDecoder sdk.TxDecoder, mem mempool.Mempool, address common.Address) (total uint64) {
	for _, tmTx := range mem.ReapMaxTxs(50) {
		tx, err := txDecoder(tmTx)
		if err != nil {
			continue
		}
		msg := tx.GetMsgs()[0]
		if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			ethTx := ethMsg.AsTransaction()
			var signer ethtypes.Signer
			if ethTx.Protected() {
				signer = ethtypes.LatestSignerForChainID(ethTx.ChainId())
			} else {
				signer = ethtypes.HomesteadSigner{}
			}
			from, _ := ethtypes.Sender(signer, ethTx)
			fmt.Println("signer", from.Hex())
			fmt.Println("checking", address.Hex())
			if bytes.Equal(from.Bytes(), address.Bytes()) {
				total++
			}
		}
	}
	return
}

// GetProof performs an ABCI query with the given key and returns a merkle proof. The desired
// tendermint height to perform the query should be set in the client context. The query will be
// performed at one below this height (at the IAVL version) in order to obtain the correct merkle
// proof. Proof queries at height less than or equal to 2 are not supported.
// Issue: https://github.com/cosmos/cosmos-sdk/issues/6567
func GetProof(height int64, storeKey string, key []byte) ([]byte, *crypto.ProofOps, error) {
	// ABCI queries at height less than or equal to 2 are not supported.
	// Base app does not support queries for height less than or equal to 1.
	// Therefore, a query at height 2 would be equivalent to a query at height 3
	if height <= 2 {
		return nil, nil, fmt.Errorf("proof queries at height <= 2 are not supported")
	}

	// Use the IAVL height if a valid tendermint height is passed in.
	height--

	abciRes, err := tmrpccore.ABCIQuery(nil, fmt.Sprintf("store/%s/key", storeKey), key, height, true)
	if err != nil {
		return nil, nil, err
	}

	return abciRes.Response.Value, abciRes.Response.ProofOps, nil
}

func FormatTmHeaderToProto(header *tmtypes.Header) tmproto.Header {
	return tmproto.Header{
		Version: header.Version,
		ChainID: header.ChainID,
		Height:  header.Height,
		Time:    header.Time,
		LastBlockId: tmproto.BlockID{
			Hash: header.LastBlockID.Hash,
			PartSetHeader: tmproto.PartSetHeader{
				Total: header.LastBlockID.PartSetHeader.Total,
				Hash:  header.LastBlockID.PartSetHeader.Hash,
			},
		},
		LastCommitHash:     header.LastCommitHash,
		DataHash:           header.DataHash,
		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		ConsensusHash:      header.ConsensusHash,
		AppHash:            header.AppHash,
		LastResultsHash:    header.LastResultsHash,
		EvidenceHash:       header.EvidenceHash,
		ProposerAddress:    header.ProposerAddress,
	}
}

// GetBlockCumulativeGas returns the cumulative gas used on a block up to a given
// transaction index. The returned gas used includes the gas from both the SDK and
// EVM module transactions.
func GetBlockCumulativeGas(blockResults *tmrpccoretypes.ResultBlockResults, idx int) uint64 {
	var gasUsed uint64

	for i := 0; i < idx && i < len(blockResults.TxsResults); i++ {
		tx := blockResults.TxsResults[i]
		gasUsed += uint64(tx.GasUsed)
	}
	return gasUsed
}

func GetPendingTx(txDecoder sdk.TxDecoder, mem mempool.Mempool, hash common.Hash, chainID *big.Int) (*RPCTransaction, error) {
	for _, uTx := range mem.ReapMaxTxs(50) {
		if bytes.Equal(uTx.Hash(), hash.Bytes()) {
			return TmTxToEthTx(txDecoder, uTx, nil, nil, nil)
		}
	}
	return nil, nil
}

func GetNonEVMSignatures(sig []byte) (v, r, s *big.Int) {
	var tmpV byte
	if len(sig) == 65 {
		tmpV = sig[len(sig)-1:][0]
	} else {
		// in case of 64 length
		tmpV = byte(int(sig[0]) % 2)
	}

	v = new(big.Int).SetBytes([]byte{tmpV + 27})
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	return
}

// TmTxToEthTx convert ethereum and rest transaction on ethereum based structure
func TmTxToEthTx(
	txDecoder sdk.TxDecoder,
	tmTx tmtypes.Tx,
	blockHash *common.Hash,
	blockNumber, index *uint64,
) (*RPCTransaction, error) {
	tx, err := txDecoder(tmTx)
	if err != nil {
		return nil, err
	}

	// the `msgIndex` is inferred from tx events, should be within the bound.
	// always taking first into account
	msg := tx.GetMsgs()[0]

	if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
		tx := ethMsg.AsTransaction()
		return NewRPCTransaction(tx, *blockHash, *blockNumber, *index)
	} else {
		addr := msg.GetSigners()[0]
		from := common.BytesToAddress(addr.Bytes())

		v := new(big.Int).SetInt64(0)
		r := new(big.Int).SetInt64(0)
		s := new(big.Int).SetInt64(0)
		sigTx, ok := tx.(authsigning.SigVerifiableTx)
		if ok {
			sigs, _ := sigTx.GetSignaturesV2()
			if len(sigs) > 0 {
				sig := sigs[0]
				sigProto := signing.SignatureDataToProto(sig.Data)
				v, r, s = GetNonEVMSignatures(sigProto.GetSingle().GetSignature())
			}
		} else {
			fmt.Printf("failed to get signature for %s\n", tmTx.Hash())
		}

		gas := uint64(0)
		gasPrice := new(big.Int).SetInt64(stratos.DefaultGasPrice)
		if feeTx, ok := tx.(sdk.FeeTx); ok {
			gas = feeTx.GetGas()
			gasPrice = new(big.Int).Div(
				feeTx.GetFee().AmountOf("wei").BigInt(), // TODO: mv somehow wei from config
				new(big.Int).SetUint64(gas),
			)
		}

		return &RPCTransaction{
			BlockHash:        blockHash,
			BlockNumber:      (*hexutil.Big)(new(big.Int).SetUint64(*blockNumber)),
			Type:             hexutil.Uint64(0),
			From:             from,
			Gas:              hexutil.Uint64(gas),
			GasPrice:         (*hexutil.Big)(gasPrice),
			Hash:             common.BytesToHash(tmTx.Hash()),
			Input:            make(hexutil.Bytes, 0),
			Nonce:            hexutil.Uint64(0),
			To:               new(common.Address),
			TransactionIndex: (*hexutil.Uint64)(index),
			Value:            (*hexutil.Big)(new(big.Int).SetInt64(0)), // NOTE: How to get value in generic way?
			V:                (*hexutil.Big)(v),
			R:                (*hexutil.Big)(r),
			S:                (*hexutil.Big)(s),
		}, nil
	}
}

// NewTransactionFromData returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewRPCTransaction(
	tx *ethtypes.Transaction, blockHash common.Hash, blockNumber, index uint64,
) (*RPCTransaction, error) {
	// Determine the signer. For replay-protected transactions, use the most permissive
	// signer, because we assume that signers are backwards-compatible with old
	// transactions. For non-protected transactions, the homestead signer signer is used
	// because the return value of ChainId is zero for those transactions.
	var signer ethtypes.Signer
	if tx.Protected() {
		signer = ethtypes.LatestSignerForChainID(tx.ChainId())
	} else {
		signer = ethtypes.HomesteadSigner{}
	}
	from, _ := ethtypes.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()
	fmt.Printf("\n v: %s, r: %s, s: %s \n", v, r, s)
	result := &RPCTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	switch tx.Type() {
	case ethtypes.AccessListTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case ethtypes.DynamicFeeTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		result.GasPrice = (*hexutil.Big)(tx.GasPrice())
	}
	return result, nil
}

// BaseFeeFromEvents parses the feemarket basefee from cosmos events
func BaseFeeFromEvents(events []abci.Event) *big.Int {
	for _, event := range events {
		if event.Type != evmtypes.EventTypeFeeMarket {
			continue
		}

		for _, attr := range event.Attributes {
			if bytes.Equal(attr.Key, []byte(evmtypes.AttributeKeyBaseFee)) {
				result, success := new(big.Int).SetString(string(attr.Value), 10)
				if success {
					return result
				}

				return nil
			}
		}
	}
	return nil
}

// FindTxAttributes returns the msg index of the eth tx in cosmos tx, and the attributes,
// returns -1 and nil if not found.
func FindTxAttributes(events []abci.Event, txHash string) (int, map[string]string) {
	msgIndex := -1
	for _, event := range events {
		if event.Type != evmtypes.EventTypeEthereumTx {
			continue
		}

		msgIndex++

		value := FindAttribute(event.Attributes, []byte(evmtypes.AttributeKeyEthereumTxHash))
		if !bytes.Equal(value, []byte(txHash)) {
			continue
		}

		// found, convert attributes to map for later lookup
		attrs := make(map[string]string, len(event.Attributes))
		for _, attr := range event.Attributes {
			attrs[string(attr.Key)] = string(attr.Value)
		}
		return msgIndex, attrs
	}
	// not found
	return -1, nil
}

// FindTxAttributesByIndex search the msg in tx events by txIndex
// returns the msgIndex, returns -1 if not found.
func FindTxAttributesByIndex(events []abci.Event, txIndex uint64) int {
	strIndex := []byte(strconv.FormatUint(txIndex, 10))
	txIndexKey := []byte(evmtypes.AttributeKeyTxIndex)
	msgIndex := -1
	for _, event := range events {
		if event.Type != evmtypes.EventTypeEthereumTx {
			continue
		}

		msgIndex++

		value := FindAttribute(event.Attributes, txIndexKey)
		if !bytes.Equal(value, strIndex) {
			continue
		}

		// found, convert attributes to map for later lookup
		return msgIndex
	}
	// not found
	return -1
}

// FindAttribute find event attribute with specified key, if not found returns nil.
func FindAttribute(attrs []abci.EventAttribute, key []byte) []byte {
	for _, attr := range attrs {
		if !bytes.Equal(attr.Key, key) {
			continue
		}
		return attr.Value
	}
	return nil
}

// GetUint64Attribute parses the uint64 value from event attributes
func GetUint64Attribute(attrs map[string]string, key string) (uint64, error) {
	value, found := attrs[key]
	if !found {
		return 0, fmt.Errorf("tx index attribute not found: %s", key)
	}
	var result int64
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	if result < 0 {
		return 0, fmt.Errorf("negative tx index: %d", result)
	}
	return uint64(result), nil
}

// AccumulativeGasUsedOfMsg accumulate the gas used by msgs before `msgIndex`.
func AccumulativeGasUsedOfMsg(events []abci.Event, msgIndex int) (gasUsed uint64) {
	for _, event := range events {
		if event.Type != evmtypes.EventTypeEthereumTx {
			continue
		}

		if msgIndex < 0 {
			break
		}
		msgIndex--

		value := FindAttribute(event.Attributes, []byte(evmtypes.AttributeKeyTxGasUsed))
		var result int64
		result, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			continue
		}
		gasUsed += uint64(result)
	}
	return
}
