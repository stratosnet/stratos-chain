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

		newRootHash := currentRootHash[:]

		// throw to tm events
		fmt.Println("newRootHash", common.Bytes2Hex(newRootHash))
		fmt.Println("commitment1", common.Bytes2Hex(commitment1))
		fmt.Println("commitment2", common.Bytes2Hex(commitment2))
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
		fmt.Println("commitments", rlProof.GetCommitments())
	}

	{
		// 3. Accept volume report for commitment 1
		leaves, err := GetLeaves(rlProof, [][]byte{
			currentRootHash,
		})
		assert.ErrorIs(t, err, nil)

		isValid, err := rlMerkle.VerifyProofs(rlProof.GetRoot(), leaves)
		assert.ErrorIs(t, err, nil)
		assert.Equal(t, isValid, true)
		fmt.Println("root hash update", rlProof.GetRoot())
		currentRootHash = rlProof.GetRoot()
	}

	{
		// 4. Relayer prepays proofs for volume report for commitment 2
		prevRootHash := NullCommitment[:]
		commitment2 := []byte("c2")
		commitment3 := []byte("c3")

		rlProof, err = rlMerkle.CreateProofs(
			[][]byte{
				currentRootHash,
				prevRootHash,
			},
			[][]byte{
				commitment3,
				commitment2,
			})

		assert.ErrorIs(t, err, nil)

		// broadcast to st volume report
		fmt.Println("newRootHash", common.Bytes2Hex(rlProof.GetRoot()))
		fmt.Println("commitments", rlProof.GetCommitments())
	}

	{
		// 5. Accept volume report for commitment 2
		leaves, err := GetLeaves(rlProof, [][]byte{
			currentRootHash,
			NullCommitment[:],
		})
		assert.ErrorIs(t, err, nil)

		isValid, err := rlMerkle.VerifyProofs(rlProof.GetRoot(), leaves)
		assert.ErrorIs(t, err, nil)
		assert.Equal(t, isValid, true)
		fmt.Println("root hash update", rlProof.GetRoot())
		// currentRootHash = rlProof.Root
	}
}
