package merkle

import "github.com/ethereum/go-ethereum/common"

// NullCommitment used to generate empty commitment for new root link as hash(newRoot, sha256([]byte{})).
// Basically it is sha256([]byte{})
var NullCommitment = [32]byte(common.Hex2Bytes("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"))

var _ MerkleProofData = &MerkleProofBundle{}

// MerkleProofBundle is base for proof data response. It could be changed for upcomming calls by other protobuf
type MerkleProofBundle struct {
	root        []byte
	commitments [][]byte
}

func NewMerkleProofBundle(root []byte, commitments [][]byte) MerkleProofData {
	return &MerkleProofBundle{
		root:        root,
		commitments: commitments,
	}
}

func (d *MerkleProofBundle) GetRoot() []byte {
	return d.root
}

func (d *MerkleProofBundle) GetCommitments() [][]byte {
	return d.commitments
}

type MerkleProofData interface {
	GetRoot() []byte
	GetCommitments() [][]byte
}

type MerkleProver interface {
	GetRoot(root []byte, commitments [][]byte) []byte
	CreateProofs(roots [][]byte, commitments [][]byte) (MerkleProofData, error)
	VerifyProofs(rootHash []byte, data [][]byte) (bool, error)
}
