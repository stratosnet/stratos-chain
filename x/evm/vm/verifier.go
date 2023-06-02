package vm

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
	return &GenesisContractVerifier{
		applyState:        make(map[uint64][]VerifiedContract),
		verifiedAddresses: map[string]bool{},
	}
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
