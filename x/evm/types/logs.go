package types

import (
	"errors"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// TransactionLogs define the logs generated from a transaction execution
// with a given hash. It it used for import/export data as transactions are not persisted
// on blockchain state after an upgrade.
type TransactionLogs struct {
	Hash string          `json:"hash"`
	Logs []*ethtypes.Log `json:"logs"`
}

// NewTransactionLogs creates a new NewTransactionLogs instance.
func NewTransactionLogs(hash ethcmn.Hash, logs []*ethtypes.Log) TransactionLogs { // nolint: interfacer
	return TransactionLogs{
		Hash: hash.String(),
		Logs: logs,
	}
}

// MarshalLogs encodes an array of logs using amino
func MarshalLogs(logs []*ethtypes.Log) ([]byte, error) {
	return ModuleCdc.MarshalBinaryLengthPrefixed(logs)
}

// UnmarshalLogs decodes an amino-encoded byte array into an array of logs
func UnmarshalLogs(in []byte) ([]*ethtypes.Log, error) {
	logs := []*ethtypes.Log{}
	err := ModuleCdc.UnmarshalBinaryLengthPrefixed(in, &logs)
	return logs, err
}

// Validate performs a basic validation of a GenesisAccount fields.
func (tx TransactionLogs) Validate() error {
	if stratos.IsEmptyHash(tx.Hash) {
		return fmt.Errorf("hash cannot be the empty %s", tx.Hash)
	}

	for i, log := range tx.Logs {
		if err := ValidateLog(log); err != nil {
			return fmt.Errorf("invalid log %d: %w", i, err)
		}
		if log.TxHash.String() != tx.Hash {
			return fmt.Errorf("log tx hash mismatch (%s ≠ %s)", log.TxHash.String(), tx.Hash)
		}
	}
	return nil
}

// ValidateLog performs a basic validation of an ethereum Log fields.
func ValidateLog(log *ethtypes.Log) error {
	if log == nil {
		return errors.New("log cannot be nil")
	}
	if stratos.IsZeroAddress(log.Address.String()) {
		return fmt.Errorf("log address cannot be empty %s", log.Address.String())
	}
	if stratos.IsEmptyHash(log.BlockHash.String()) {
		return fmt.Errorf("block hash cannot be the empty %s", log.BlockHash.String())
	}
	if log.BlockNumber == 0 {
		return errors.New("block number cannot be zero")
	}
	if stratos.IsEmptyHash(log.TxHash.String()) {
		return fmt.Errorf("tx hash cannot be the empty %s", log.TxHash.String())
	}
	return nil
}
