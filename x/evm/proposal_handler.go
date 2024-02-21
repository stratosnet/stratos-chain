package evm

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/stratosnet/stratos-chain/x/evm/keeper"
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// NewEVMChangeProposalHandler defines the evm changes proposals
func NewEVMChangeProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		pc, err := keeper.NewProposalCounsil(k, ctx)
		if err != nil {
			return err
		}

		switch c := content.(type) {
		case *types.UpdateImplmentationProposal:
			return pc.UpdateProxyImplementation(c)

		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized proxy proposal content type: %T", c)
		}
	}
}
