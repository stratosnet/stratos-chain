package merkle

import (
	"bytes"
	"fmt"

	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/ethereum/go-ethereum/common"
)

var _ MerkleProver = &RelayerMerkleProver{}

type RelayerMerkleProver struct{}

func NewRelayerMerkleProver() MerkleProver {
	return &RelayerMerkleProver{}
}

func (ms RelayerMerkleProver) GetRoot(root []byte, commitments [][]byte) []byte {
	var data [][]byte

	for _, c := range commitments {
		data = append(data, root)
		data = append(data, c)
	}

	root, _ = merkle.ProofsFromByteSlices(data)
	return root
}

func (ms RelayerMerkleProver) CreateProofs(roots [][]byte, commitments [][]byte) (MerkleProofData, error) {
	if len(roots) != len(commitments) {
		return nil, fmt.Errorf("roots do not have not enough commitments for pairing")
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("roots/commitments could not be empty")
	}

	var data [][]byte

	for i, c := range commitments {
		data = append(data, roots[i])
		data = append(data, c)
	}

	root, _ := merkle.ProofsFromByteSlices(data)

	rlProof := NewMerkleProofBundle(root, commitments)

	return rlProof, nil
}

func (ms RelayerMerkleProver) VerifyProofs(rootHash []byte, data [][]byte) (bool, error) {
	newRootHash, _ := merkle.ProofsFromByteSlices(data)
	if !bytes.Equal(rootHash, newRootHash) {
		return false, fmt.Errorf("roots not equal: %s != %s", common.Bytes2Hex(rootHash), common.Bytes2Hex(newRootHash))
	}

	return true, nil
}

func (ms RelayerMerkleProver) CreateFileUploadProofs(rootHash []byte, fileHashes []string) ([]byte, []*merkle.Proof) {
	var data [][]byte

	for _, file := range fileHashes {
		data = append(data, rootHash)
		data = append(data, []byte(file))
	}
	root, proofs := merkle.ProofsFromByteSlices(data)
	return root, proofs
}
