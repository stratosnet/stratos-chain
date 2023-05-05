package keeper

import (
	"sort"

	"github.com/stratosnet/stratos-chain/x/evm/types"
)

type ProposalVerifier struct {
	applyState        map[uint64][]*types.ProxyContractInitState
	verifiedAddresses map[string]bool
}

func NewProposalVerifier() *ProposalVerifier {
	return &ProposalVerifier{
		applyState:        make(map[uint64][]*types.ProxyContractInitState),
		verifiedAddresses: map[string]bool{},
	}
}

func (pc *ProposalVerifier) GetStates(height uint64) []*types.ProxyContractInitState {
	return pc.applyState[height]
}

func (pv *ProposalVerifier) ApplyParamsState(params types.Params) {
	keys := make([]string, 0)
	for k, _ := range params.ProxyProposalParams.Contracts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		contract := params.ProxyProposalParams.Contracts[k]
		pv.AddTrustedAddress(contract.Address)

		pv.applyState[contract.Height] = append(pv.applyState[contract.Height], contract)
	}
}

func (pv *ProposalVerifier) AddTrustedAddress(addr string) {
	pv.verifiedAddresses[addr] = true
}

func (pv *ProposalVerifier) IsTrustedAddress(addr string) bool {
	return pv.verifiedAddresses[addr]
}
