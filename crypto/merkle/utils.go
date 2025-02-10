package merkle

import (
	"crypto/sha256"
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
