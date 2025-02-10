package merkle

import "github.com/ethereum/go-ethereum/common"

// NullCommitment used to generate empty commitment for new root link as hash(newRoot, sha256([]byte{})).
// Basically it is sha256([]byte{})
var NullCommitment = [32]byte(common.Hex2Bytes("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"))

var _ MerkleProofData = &MerkleProofBundle{}

// MerkleProofBundle is base for proof data response. It could be changed for upcomming calls by other protobuf
type MerkleProofBundle struct {
	root   []byte
	data   [][]byte
	proofs [][]byte
}

func NewMerkleProofBundle(root []byte, data [][]byte, proofs [][]byte) MerkleProofData {
	return &MerkleProofBundle{
		root:   root,
		data:   data,
		proofs: proofs,
	}
}

func (d *MerkleProofBundle) GetRoot() []byte {
	return d.root
}

func (d *MerkleProofBundle) GetLeaves() [][]byte {
	return d.data
}

func (d *MerkleProofBundle) GetProofs() [][]byte {
	return d.proofs
}

type MerkleProofData interface {
	GetRoot() []byte
	GetLeaves() [][]byte
	GetProofs() [][]byte
}

type MerkleProver interface {
	GetRoot(root []byte, commitments [][]byte) []byte
	CreateProofs(roots [][]byte, commitments [][]byte) (MerkleProofData, error)
	VerifyProofs(rootHash []byte, proofs [][]byte, data [][]byte) (bool, error)
}
