package testutil

import (
	"fmt"

	"github.com/stratosnet/stratos-chain/crypto"
	"github.com/stratosnet/stratos-chain/crypto/bls"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
)

// SignVolumeReport Use ed25519 private key of meta nodes to sign
func SignVolumeReport(volumeReportMsg *pottypes.MsgVolumeReport, privKeys ...[]byte) (*pottypes.MsgVolumeReport, error) {
	if len(privKeys) == 0 {
		return nil, fmt.Errorf("No private keys, failed to sign. ")
	}

	signBytes := volumeReportMsg.GetBLSSignBytes()
	signBytesHash := crypto.Keccak256(signBytes)

	var blsSignatures = make([][]byte, len(privKeys))
	var blsPrivKeys = make([][]byte, len(privKeys))
	var blsPubKeys = make([][]byte, len(privKeys))
	var err error

	for i, privKey := range privKeys {
		blsPrivKeys[i], blsPubKeys[i], err = bls.NewKeyPairFromBytes(privKey)
		if err != nil {
			return nil, err
		}
	}

	for i, blsPrivKey := range blsPrivKeys {
		blsSignatures[i], err = bls.Sign(signBytesHash, blsPrivKey)
		if err != nil {
			return nil, err
		}
	}

	finalBlsSignature, err := bls.AggregateSignatures(blsSignatures...)
	signature := pottypes.NewBLSSignatureInfo(blsPubKeys, finalBlsSignature, signBytesHash)
	volumeReportMsg.BLSSignature = signature
	return volumeReportMsg, nil
}
