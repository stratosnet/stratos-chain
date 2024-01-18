package keeper

import (
	"context"
	"encoding/json"
	"strconv"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stratosnet/stratos-chain/x/evm/types"
)

var _ types.MsgServer = &Keeper{}

// EthereumTx implements the gRPC MsgServer interface. It receives a transaction which is then
// executed (i.e. applied) against the go-ethereum EVM. The provided SDK Context is set to the Keeper
// so that it can implement and call the StateDB methods without receiving it as a function
// parameter.
func (k *Keeper) EthereumTx(goCtx context.Context, msg *types.MsgEthereumTx) (*types.MsgEthereumTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tx := msg.AsTransaction()
	txIndex := k.GetTxIndexTransient(ctx)

	response, err := k.ApplyTransaction(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply transaction")
	}

	// build events
	eventEthereumTx := &types.EventEthereumTx{
		Amount:  tx.Value().String(),
		EthHash: response.Hash,                            // ethereum transaction hash format
		Index:   strconv.FormatUint(txIndex, 10),          // index of valid ethereum tx
		GasUsed: strconv.FormatUint(response.GasUsed, 10), // eth tx gas used, we can't get it from cosmos tx result when it contains multiple eth tx msgs.
	}
	if len(ctx.TxBytes()) > 0 {
		// tendermint transaction hash format
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		eventEthereumTx.Hash = hash.String()
	}
	if to := tx.To(); to != nil {
		eventEthereumTx.Recipient = to.Hex()
	}
	if response.Failed() {
		eventEthereumTx.EthTxFailed = response.VmError
	}

	eventTxLog := &types.EventTxLog{}
	value, err := json.Marshal(response.Logs)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrJSONMarshal, "failed to encode log")
	}
	eventTxLog.TxLogs = value

	// emit events
	if err := ctx.EventManager().EmitTypedEvents(
		eventEthereumTx,
		eventTxLog,
	); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateParams updates the module parameters
func (k *Keeper) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
