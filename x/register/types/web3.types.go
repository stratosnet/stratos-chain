package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	stratos "github.com/stratosnet/stratos-chain/types"
)

func (m *MsgCreateResourceNode) GetWeb3Msg() (*stratos.Web3Msg, error) {
	bTo, err := sdk.AccAddressFromBech32(m.GetNetworkAddress())
	if err != nil {
		return nil, err
	}
	to := common.BytesToAddress(bTo.Bytes())
	value := m.Value.Amount.BigInt()
	return &stratos.Web3Msg{
		To:    &to,
		Value: value,
	}, nil
}

func (m *MsgCreateMetaNode) GetWeb3Msg() (*stratos.Web3Msg, error) {
	bTo, err := sdk.AccAddressFromBech32(m.GetNetworkAddress())
	if err != nil {
		return nil, err
	}
	to := common.BytesToAddress(bTo.Bytes())
	value := m.Value.Amount.BigInt()
	return &stratos.Web3Msg{
		To:    &to,
		Value: value,
	}, nil
}

func (m *MsgRemoveMetaNode) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}

func (m *MsgMetaNodeRegistrationVote) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}

func (m *MsgWithdrawMetaNodeRegistrationDeposit) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}

func (m *MsgUpdateResourceNodeResponse) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}

func (m *MsgUpdateResourceNodeDeposit) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}

func (m *MsgUpdateEffectiveDeposit) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}

func (m *MsgUpdateMetaNode) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}

func (m *MsgUpdateMetaNodeDeposit) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}
