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

	root, proofs := merkle.ProofsFromByteSlices(data)

	var mProofs [][]byte

	for _, proof := range proofs {
		proofData, err := proof.ToProto().Marshal()
		if err != nil {
			return nil, err
		}
		mProofs = append(mProofs, proofData)
	}

	rlProof := NewMerkleProofBundle(root, data, mProofs)

	return rlProof, nil
}

func (ms RelayerMerkleProver) VerifyProofs(rootHash []byte, proofs [][]byte, data [][]byte) (bool, error) {
	newRootHash, newProofs := merkle.ProofsFromByteSlices(data)
	if !bytes.Equal(rootHash, newRootHash) {
		return false, fmt.Errorf("roots not equal: %s != %s", common.Bytes2Hex(rootHash), common.Bytes2Hex(newRootHash))
	}

	if len(newProofs) != len(proofs) {
		return false, fmt.Errorf("proofs length not match")
	}

	for i, np := range newProofs {
		npData, err := np.ToProto().Marshal()
		if err != nil {
			return false, err
		}
		if !bytes.Equal(npData, proofs[i]) {
			return false, fmt.Errorf("proof aunt not match")
		}
	}

	return true, nil
}
