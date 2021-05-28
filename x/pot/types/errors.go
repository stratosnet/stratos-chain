package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid               = sdkerrors.Register(ModuleName, 1, "custom error message")
	ErrUnknownAccountAddress = sdkerrors.Register(ModuleName, 2, "account address does not exist")
	ErrOutOfIssuance         = sdkerrors.Register(ModuleName, 3, "mining reward reaches the issuance limit")
	//ErrBadTrafficRewardDistribution = sdkerrors.Register(ModuleName, 3, "traffic pool does not have sufficient coins to distribute")
)
