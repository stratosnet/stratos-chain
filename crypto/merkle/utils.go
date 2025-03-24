package merkle

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

func CreateSdkCommitment(signer []byte, nonce uint64, data []byte) []byte {
	// signer + nonce + any data
	cb := sha256.Sum256(
		append(
			signer,
			append(
				[]byte(strconv.Itoa(int(nonce))),
				data...,
			)...,
		),
	)
	return cb[:]
}

func GetLeaves(mdata MerkleProofData, roots [][]byte) ([][]byte, error) {
	if len(mdata.GetCommitments()) == 0 {
		return nil, fmt.Errorf("commitments could not be empty")
	}
	if len(roots) != len(mdata.GetCommitments()) {
		return nil, fmt.Errorf("roots do not have not enough commitments for pairing")
	}
	leaves := make([][]byte, 0, len(mdata.GetCommitments())*2)
	for i, c := range mdata.GetCommitments() {
		leaves = append(leaves, roots[i], c)
	}
	return leaves, nil
}
