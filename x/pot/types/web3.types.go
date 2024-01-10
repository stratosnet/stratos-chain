package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	stratos "github.com/stratosnet/stratos-chain/types"
)

func (m *MsgVolumeReport) GetWeb3Msg() (*stratos.Web3Msg, error) {
	bFrom, err := sdk.AccAddressFromBech32(m.GetReporter())
	if err != nil {
		return nil, err
	}
	from := common.BytesToAddress(bFrom.Bytes())
	return &stratos.Web3Msg{
		From: &from,
	}, nil
}

func (m *MsgWithdraw) GetWeb3Msg() (*stratos.Web3Msg, error) {
	bFrom, err := sdk.AccAddressFromBech32(m.GetWalletAddress())
	if err != nil {
		return nil, err
	}
	bTo, err := sdk.AccAddressFromBech32(m.GetTargetAddress())
	if err != nil {
		return nil, err
	}
	from := common.BytesToAddress(bFrom.Bytes())
	to := common.BytesToAddress(bTo.Bytes())
	value := m.Amount.AmountOf(stratos.Wei).BigInt()
	return &stratos.Web3Msg{
		From:  &from,
		To:    &to,
		Value: value,
	}, nil
}

func (m *MsgFoundationDeposit) GetWeb3Msg() (*stratos.Web3Msg, error) {
	bFrom, err := sdk.AccAddressFromBech32(m.GetFrom())
	if err != nil {
		return nil, err
	}
	from := common.BytesToAddress(bFrom.Bytes())
	to := common.BytesToAddress(authtypes.NewModuleAddress(FoundationAccount).Bytes())
	value := m.Amount.AmountOf(stratos.Wei).BigInt()
	return &stratos.Web3Msg{
		From:  &from,
		To:    &to,
		Value: value,
	}, nil
}

func (m *MsgSlashingResourceNode) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}
