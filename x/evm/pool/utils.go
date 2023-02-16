package pool

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stratosnet/stratos-chain/x/evm/types"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmjsonrpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

func GetTmTxByHash(hash common.Hash) (*tmrpctypes.ResultTx, error) {
	resTx, err := tmrpccore.Tx(nil, hash.Bytes(), false)
	if err != nil {
		query := fmt.Sprintf("%s.%s='%s'", types.TypeMsgEthereumTx, types.AttributeKeyEthereumTxHash, hash.Hex())
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
