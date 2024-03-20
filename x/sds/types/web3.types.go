package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	stratos "github.com/stratosnet/stratos-chain/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
)

func (m *MsgFileUpload) GetWeb3Msg() (*stratos.Web3Msg, error) {
	return nil, nil // SKIP
}

func (m *MsgPrepay) GetWeb3Msg() (*stratos.Web3Msg, error) {
	bFrom, err := sdk.AccAddressFromBech32(m.GetSender())
	if err != nil {
		return nil, err
	}
	from := common.BytesToAddress(bFrom.Bytes())
	to := common.BytesToAddress(authtypes.NewModuleAddress(registertypes.TotalUnissuedPrepay).Bytes())
	value := m.Amount.AmountOf(stratos.Wei).BigInt()
	return &stratos.Web3Msg{
		From:  &from,
		To:    &to,
		Value: value,
	}, nil
}
