package types

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/mempool"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmrpccore "github.com/cometbft/cometbft/rpc/core"
	tmrpccoretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
)

var MempoolCapacity = 100

// HashToUint64 used to convert string or hash string to uint64
func HashToUint64(s string) uint64 {
	h := sha1.New()
	h.Write([]byte(s))

	hash := h.Sum(nil)[:2]
	result := uint64(0)
	for i := 0; i < len(hash); i++ {
		result = result << 8
		result += uint64(hash[i])

	}
	return result
}

// RawTxToEthTx returns an evm MsgEthereum transaction from raw tx bytes.
func RawTxToEthTx(clientCtx client.Context, txBz tmtypes.Tx) ([]*evmtypes.MsgEthereumTx, error) {
	tx, err := clientCtx.TxConfig.TxDecoder()(txBz)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
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

func GetBlockBloom(blockResults *tmrpccoretypes.ResultBlockResults) (ethtypes.Bloom, error) {
	for _, event := range blockResults.EndBlockEvents {
		// <v012 support
		if event.GetType() == evmtypes.EventTypeBlockBloom {
			for _, attr := range event.Attributes {
				if attr.Key == evmtypes.AttributeKeyEthereumBloom {
					return ethtypes.BytesToBloom([]byte(attr.Value)), nil
				}
			}
		} else {
			msg, err := sdk.ParseTypedEvent(event)
			if err != nil {
				// ignore in case of not typed event (in case of not migrated code)
				continue
			}

			bloomEvt, ok := msg.(*evmtypes.EventBlockBloom)
			if !ok {
				continue
			}
			return ethtypes.BytesToBloom([]byte(bloomEvt.Bloom)), nil
		}
	}
	return ethtypes.Bloom{}, fmt.Errorf("block bloom event is not found")
}

// EthHeaderFromTendermint is a util function that returns an Ethereum Header
// from a tendermint Header.
func EthHeaderFromTendermint(header tmtypes.Header) (*Header, error) {
	results, err := tmrpccore.BlockResults(nil, &header.Height)
	if err != nil {
		return nil, err
	}
	gasLimit, err := evmtypes.BlockMaxGasFromConsensusParams(&header.Height)
	if err != nil {
		return nil, err
	}

	bloom, err := GetBlockBloom(results)
	if err != nil {
		return nil, err
	}

	gasUsed := int64(0)
	for _, txResult := range results.TxsResults {
		gasUsed += txResult.GasUsed
	}

	return &Header{
		Hash:        common.BytesToHash(header.Hash()),
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash),
		UncleHash:   common.Hash{},
		Coinbase:    common.BytesToAddress(header.ProposerAddress),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      common.BytesToHash(header.DataHash),
		ReceiptHash: common.Hash{},
		Bloom:       bloom,
		Difficulty:  big.NewInt(1), // NOTE: Maybe move to some constant?
		Number:      big.NewInt(header.Height),
		GasLimit:    uint64(gasLimit),
		GasUsed:     uint64(gasUsed),
		Time:        uint64(header.Time.Unix()),
		Extra:       common.Hex2Bytes(""),
		MixDigest:   common.Hash{},
		Nonce:       ethtypes.BlockNonce{},
		BaseFee:     nil,
		// TODO: Add size somehow for legacy subscription support as for a new Header type after London
		// is not exist but still present in newBlockHeaders call on subscription
	}, nil
}

// EthBlockFromTendermint returns a JSON-RPC compatible Ethereum blockfrom a given Tendermint block.
func EthBlockFromTendermint(txDecoder sdk.TxDecoder, block *tmtypes.Block, fullTx bool) (*Block, error) {
	header, err := EthHeaderFromTendermint(block.Header)
	if err != nil {
		return nil, err
	}

	transactions := make([]interface{}, 0, len(block.Txs))
	for i, tmTx := range block.Txs {
		if !fullTx {
			transactions = append(transactions, fmt.Sprintf("0x%s", common.Bytes2Hex(tmTx.Hash())))
			continue
		}
		blockHash := common.BytesToHash(block.Hash())
		blockHeight := uint64(block.Height)
		txIndex := uint64(i)
		tx, err := TmTxToEthTx(txDecoder, tmTx, &blockHash, &blockHeight, &txIndex)
		if err != nil {
			// NOTE: Add debug?
			continue
		}
		transactions = append(transactions, tx)
	}

	return &Block{
		Number:           hexutil.Uint64(header.Number.Uint64()),
		Hash:             header.Hash,
		ParentHash:       header.ParentHash,
		Nonce:            ethtypes.BlockNonce{}, // PoW specific
		Sha3Uncles:       common.Hash{},         // No uncles in Tendermint
		LogsBloom:        header.Bloom,
		TransactionsRoot: header.TxHash,
		StateRoot:        header.Root,
		Miner:            header.Coinbase,
		MixHash:          common.Hash{},
		Difficulty:       hexutil.Uint64(header.Difficulty.Uint64()),
		TotalDifficulty:  hexutil.Uint64(header.Difficulty.Uint64()),
		ExtraData:        common.Hex2Bytes(""),
		Size:             hexutil.Uint64(block.Size()),
		GasLimit:         (*hexutil.Big)(new(big.Int).SetUint64(header.GasLimit)),
		GasUsed:          (*hexutil.Big)(new(big.Int).SetUint64(header.GasUsed)),
		Timestamp:        hexutil.Uint64(header.Time),
		Transactions:     transactions,
		Uncles:           make([]common.Hash, 0),
		Withdrawals:      make([]common.Hash, 0),
		ReceiptsRoot:     common.Hash{},
		BaseFee:          nil,
	}, nil
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
	for _, tmTx := range mem.ReapMaxTxs(MempoolCapacity) {
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

func GetPendingTx(txDecoder sdk.TxDecoder, mem mempool.Mempool, hash common.Hash) (*Transaction, error) {
	for _, uTx := range mem.ReapMaxTxs(MempoolCapacity) {
		if bytes.Equal(uTx.Hash(), hash.Bytes()) {
			return TmTxToEthTx(txDecoder, uTx, nil, nil, nil)
		}
	}
	return nil, nil
}

func GetPendingTxs(txDecoder sdk.TxDecoder, mem mempool.Mempool) []*Transaction {
	txs := make([]*Transaction, 0, MempoolCapacity)
	for _, uTx := range mem.ReapMaxTxs(MempoolCapacity) {
		if tx, err := TmTxToEthTx(txDecoder, uTx, nil, nil, nil); err == nil {
			txs = append(txs, tx)
		}
	}
	sort.Sort(TxByNonce(txs))
	return txs
}

func GetPendingTxsLen(mem mempool.Mempool) int {
	return len(mem.ReapMaxTxs(MempoolCapacity))
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

// GetTxHash get hash depends on what type, unfortunately system support two tx hash algo
// in order to have opportunity for off chain tx hash computation
func GetTxHash(txDecoder sdk.TxDecoder, tmTx tmtypes.Tx) (common.Hash, error) {
	tx, err := txDecoder(tmTx)
	if err != nil {
		return common.Hash{}, err
	}
	if len(tx.GetMsgs()) == 0 {
		return common.Hash{}, fmt.Errorf("tx contain empty msgs")
	}
	msg := tx.GetMsgs()[0]
	if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
		return ethMsg.AsTransaction().Hash(), nil
	} else {
		return common.BytesToHash(tmTx.Hash()), nil
	}
}

// TmTxToEthTx convert ethereum and rest transaction on ethereum based structure
func TmTxToEthTx(
	txDecoder sdk.TxDecoder,
	tmTx tmtypes.Tx,
	blockHash *common.Hash,
	blockNumber, index *uint64,
) (*Transaction, error) {
	tx, err := txDecoder(tmTx)
	if err != nil {
		return nil, err
	}

	// the `msgIndex` is inferred from tx events, should be within the bound.
	// always taking first into account
	msg := tx.GetMsgs()[0]

	if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
		tx := ethMsg.AsTransaction()
		return NewRPCTransaction(tx, blockHash, blockNumber, index)
	} else {
		addr := msg.GetSigners()[0]
		from := common.BytesToAddress(addr.Bytes())

		nonce := uint64(0)
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
				nonce = sig.Sequence
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

		var bNumber *big.Int
		if blockNumber != nil {
			bNumber = new(big.Int).SetUint64(*blockNumber)
		}

		to := new(common.Address)
		value := new(big.Int).SetInt64(0)

		switch sdkTx := msg.(type) {
		case stratos.Web3MsgType:
			web3Msg, err := sdkTx.GetWeb3Msg()
			if err != nil {
				return nil, err
			}
			// Just in case not implemented
			if web3Msg != nil {
				if web3Msg.From != nil {
					from = *web3Msg.From
				}
				if web3Msg.To != nil {
					to = web3Msg.To
				}
				if web3Msg.Value != nil {
					value = web3Msg.Value
				}
			}

			// NOTE: As we could not add GetWeb3Msg method, it will be handled directly here
			// cosmos.bank.v1beta1.MsgSend
		case *banktypes.MsgSend:
			bFrom, err := sdk.AccAddressFromBech32(sdkTx.FromAddress)
			if err != nil {
				return nil, err
			}
			bTo, err := sdk.AccAddressFromBech32(sdkTx.ToAddress)
			if err != nil {
				return nil, err
			}
			from = common.BytesToAddress(bFrom.Bytes())
			*to = common.BytesToAddress(bTo.Bytes())
			value = sdkTx.Amount.AmountOf(stratos.Wei).BigInt()
		}

		return &Transaction{
			BlockHash:        blockHash,
			BlockNumber:      (*hexutil.Big)(bNumber),
			Type:             hexutil.Uint64(HashToUint64(sdk.MsgTypeURL(msg))),
			From:             from,
			Gas:              hexutil.Uint64(gas),
			GasPrice:         (*hexutil.Big)(gasPrice),
			Hash:             common.BytesToHash(tmTx.Hash()),
			Input:            make(hexutil.Bytes, 0),
			Nonce:            hexutil.Uint64(nonce),
			To:               to,
			TransactionIndex: (*hexutil.Uint64)(index),
			Value:            (*hexutil.Big)(value),
			V:                (*hexutil.Big)(v),
			R:                (*hexutil.Big)(r),
			S:                (*hexutil.Big)(s),
		}, nil
	}
}

// NewRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewRPCTransaction(
	tx *ethtypes.Transaction, blockHash *common.Hash, blockNumber, index *uint64,
) (*Transaction, error) {
	// Determine the signer. For replay-protected transactions, use the most permissive
	// signer, because we assume that signers are backwards-compatible with old
	// transactions. For non-protected transactions, the homestead signer is used
	// because the return value of ChainId is zero for those transactions.
	var signer ethtypes.Signer
	if tx.Protected() {
		signer = ethtypes.LatestSignerForChainID(tx.ChainId())
	} else {
		signer = ethtypes.HomesteadSigner{}
	}
	from, _ := ethtypes.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()
	result := &Transaction{
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
	if blockHash != nil {
		result.BlockHash = blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(*blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(index)
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

// FindEthereumTxEvent returns the msg index of the eth tx in cosmos tx, and the tx evt,
// returns -1 and nil if not found.
func FindEthereumTxEvent(events []abci.Event, txHash string) (int, *evmtypes.EventEthereumTx) {
	msgIndex := -1
	for _, event := range events {
		var evtEthTx *evmtypes.EventEthereumTx
		// <v012 support
		if event.GetType() == evmtypes.EventTypeEthereumTx {
			msgIndex++

			value := FindAttribute(event.Attributes, evmtypes.AttributeKeyEthereumTxHash)
			if value != txHash {
				continue
			}

			// found, convert attributes to map for later lookup
			attrs := make(map[string]string, len(event.Attributes))
			for _, attr := range event.Attributes {
				attrs[attr.Key] = attr.Value
			}

			evtEthTx = &evmtypes.EventEthereumTx{
				EthHash:     value,
				EthTxFailed: attrs[evmtypes.AttributeKeyEthereumTxFailed],
			}
		} else {
			msg, err := sdk.ParseTypedEvent(event)
			if err != nil {
				// ignore in case of not typed event (in case of not migrated code)
				continue
			}

			var ok bool
			evtEthTx, ok = msg.(*evmtypes.EventEthereumTx)
			if !ok {
				continue
			}

			msgIndex++

			if evtEthTx.EthHash != txHash {
				continue
			}
		}

		return msgIndex, evtEthTx
	}
	// not found
	return -1, &evmtypes.EventEthereumTx{}
}

// FindAttribute find event attribute with specified key, if not found returns nil.
func FindAttribute(attrs []abci.EventAttribute, key string) string {
	for _, attr := range attrs {
		if attr.GetKey() != key {
			continue
		}
		return attr.Value
	}
	return ""
}

var errPrune = fmt.Errorf("version mismatch on immutable IAVL tree; version does not exist. Version has either been pruned, or is for a future block height")

func IsPruneError(err error) bool {
	return strings.Contains(err.Error(), errPrune.Error())
}

var ErrNotArchiveNode = fmt.Errorf("the data is available only in the archive node")
