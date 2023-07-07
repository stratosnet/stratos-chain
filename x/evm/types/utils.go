package types

import (
	"fmt"
	"math/big"

	"github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/pkg/errors"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmjsonrpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

const maxBitLen = 256

var EmptyCodeHash = crypto.Keccak256(nil)

// DecodeTxResponse decodes an protobuf-encoded byte slice into TxResponse
func DecodeTxResponse(in []byte) (*MsgEthereumTxResponse, error) {
	var txMsgData sdk.TxMsgData
	if err := proto.Unmarshal(in, &txMsgData); err != nil {
		return nil, err
	}

	data := txMsgData.GetData()
	if len(data) == 0 {
		return &MsgEthereumTxResponse{}, nil
	}

	var res MsgEthereumTxResponse

	err := proto.Unmarshal(data[0].GetData(), &res)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal tx response message data")
	}

	return &res, nil
}

// EncodeTransactionLogs encodes TransactionLogs slice into a protobuf-encoded byte slice.
func EncodeTransactionLogs(res *TransactionLogs) ([]byte, error) {
	return proto.Marshal(res)
}

// DecodeTxResponse decodes an protobuf-encoded byte slice into TransactionLogs
func DecodeTransactionLogs(data []byte) (TransactionLogs, error) {
	var logs TransactionLogs
	err := proto.Unmarshal(data, &logs)
	if err != nil {
		return TransactionLogs{}, err
	}
	return logs, nil
}

// UnwrapEthereumMsg extract MsgEthereumTx from wrapping sdk.Tx
func UnwrapEthereumMsg(tx *sdk.Tx, ethHash common.Hash) (*MsgEthereumTx, error) {
	if tx == nil {
		return nil, fmt.Errorf("invalid tx: nil")
	}

	for _, msg := range (*tx).GetMsgs() {
		ethMsg, ok := msg.(*MsgEthereumTx)
		if !ok {
			return nil, fmt.Errorf("invalid tx type: %T", tx)
		}
		if ethMsg.AsTransaction().Hash() == ethHash {
			return ethMsg, nil
		}
	}

	return nil, fmt.Errorf("eth tx not found: %s", ethHash)
}

// BinSearch execute the binary search and hone in on an executable gas limit
func BinSearch(lo, hi uint64, executable func(uint64) (bool, *MsgEthereumTxResponse, error)) (uint64, error) {
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := executable(mid)
		// If the error is not nil(consensus error), it means the provided message
		// call or transaction will never be accepted no matter how much gas it is
		// assigned. Return the error directly, don't struggle any more.
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	return hi, nil
}

// SafeNewIntFromBigInt constructs Int from big.Int, return error if more than 256bits
func SafeNewIntFromBigInt(i *big.Int) (sdk.Int, error) {
	if !IsValidInt256(i) {
		return sdk.NewInt(0), fmt.Errorf("big int out of bound: %s", i)
	}
	return sdk.NewIntFromBigInt(i), nil
}

// IsValidInt256 check the bound of 256 bit number
func IsValidInt256(i *big.Int) bool {
	return i == nil || i.BitLen() <= maxBitLen
}

// GetTmTxByHash return result tx in according of dynamic tx searching
func GetTmTxByHash(hash common.Hash) (*tmrpctypes.ResultTx, error) {
	resTx, err := tmrpccore.Tx(nil, hash.Bytes(), false)
	if err != nil {
		query := fmt.Sprintf("%s.%s='%s'", TypeMsgEthereumTx, AttributeKeyEthereumTxHash, hash.Hex())
		resTxs, err := tmrpccore.TxSearch(new(tmjsonrpctypes.Context), query, false, nil, nil, "")
		if err != nil {
			return nil, err
		}
		if len(resTxs.Txs) == 0 {
			return nil, errors.Errorf("ethereum tx not found for hash %s", hash.Hex())
		}
		return resTxs.Txs[0], nil
	}
	return resTx, nil
}

// BlockMaxGasFromConsensusParams returns the gas limit for the current block from the chain consensus params.
func BlockMaxGasFromConsensusParams(blockHeight *int64) (int64, error) {
	resConsParams, err := tmrpccore.ConsensusParams(nil, blockHeight)
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
