package vm

import "github.com/ethereum/go-ethereum/common"

var (
	// Default reserved contract controll addresses
	ConsensusAddress  = common.HexToAddress("0x1000000000000000000000000000000000000000")
	ProxyOwnerAddress = common.HexToAddress("0x1000000000000000000000000000000000000001")
	// Default reserved contracts
	PrepayContractAddress = common.HexToAddress("0x1000000000000000000000000000000000010101")
)

type VerifiedContract interface {
	GetHeight() uint64
	GetAddress() string
	GetBin() string
	GetInit() string
}

type GenesisContractVerifier struct {
	applyState        map[uint64][]VerifiedContract
	verifiedAddresses map[string]bool
}

func NewGenesisContractVerifier() *GenesisContractVerifier {
	gcv := &GenesisContractVerifier{
		applyState:        make(map[uint64][]VerifiedContract),
		verifiedAddresses: map[string]bool{},
	}
	// init trusted addresses
	gcv.initTrustedAddresses()
	return gcv
}

func (gcv *GenesisContractVerifier) initTrustedAddresses() {
	gcv.AddTrustedAddress(PrepayContractAddress.Hex())
}

func (gcv *GenesisContractVerifier) GetContracts(height uint64) []VerifiedContract {
	return gcv.applyState[height]
}

func (gcv *GenesisContractVerifier) AddContract(contract VerifiedContract, trusted bool) {
	gcv.applyState[contract.GetHeight()] = append(gcv.applyState[contract.GetHeight()], contract)
	if trusted {
		gcv.AddTrustedAddress(contract.GetAddress())
	}
}

func (gcv *GenesisContractVerifier) AddTrustedAddress(addr string) {
	gcv.verifiedAddresses[addr] = true
}

func (gcv *GenesisContractVerifier) IsTrustedAddress(addr string) bool {
	return gcv.verifiedAddresses[addr]
}
