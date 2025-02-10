package merkle

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestRelayerMerkle_CreateRelayAndAccept(t *testing.T) {
	var (
		rlProof MerkleProofData
		err     error
	)
	currentRootHash := NullCommitment[:]
	rlMerkle := NewRelayerMerkleProver()
	{
		// 1. Create prepay new root
		commitment1 := []byte("c1")
		commitment2 := []byte("c2")

		newRootHash := rlMerkle.GetRoot(currentRootHash, [][]byte{
			commitment1,
			commitment2,
		})

		// adding new root to the store
		fmt.Println("newRootHash", common.Bytes2Hex(newRootHash))
	}

	{
		// 2. Relayer prepays proofs for volume report for commitment 1
		commitment1 := []byte("c1")

		rlProof, err = rlMerkle.CreateProofs(
			[][]byte{
				currentRootHash,
			},
			[][]byte{
				commitment1,
			})

		assert.ErrorIs(t, err, nil)

		// broadcast to st volume report
		fmt.Println("newRootHash", common.Bytes2Hex(rlProof.GetRoot()))
		fmt.Println("leafData", rlProof.GetLeaves())
	}

	{
		// 3. Accept volume report for commitment 1
		isValid, err := rlMerkle.VerifyProofs(rlProof.GetRoot(), rlProof.GetProofs(), rlProof.GetLeaves())
		assert.ErrorIs(t, err, nil)
		assert.Equal(t, isValid, true)
		fmt.Println("root hash update", rlProof.GetRoot())
		currentRootHash = rlProof.GetRoot()
	}

	{
		// 4. Relayer prepays proofs for volume report for commitment 2
		prevRootHash := NullCommitment[:]
		commitment2 := []byte("c2")

		rlProof, err = rlMerkle.CreateProofs(
			[][]byte{
				currentRootHash,
				prevRootHash,
			},
			[][]byte{
				NullCommitment[:],
				commitment2,
			})

		assert.ErrorIs(t, err, nil)

		// broadcast to st volume report
		fmt.Println("newRootHash", common.Bytes2Hex(rlProof.GetRoot()))
		fmt.Println("leafData", rlProof.GetLeaves())
	}

	{
		// 5. Accept volume report for commitment 2
		isValid, err := rlMerkle.VerifyProofs(rlProof.GetRoot(), rlProof.GetProofs(), rlProof.GetLeaves())
		assert.ErrorIs(t, err, nil)
		assert.Equal(t, isValid, true)
		fmt.Println("root hash update", rlProof.GetRoot())
		// currentRootHash = rlProof.Root
	}
}
