package ante

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	evmtypes "github.com/stratosnet/stratos-chain/x/evm/types"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper, EVM Keeper and Fee Market Keeper.
type HandlerOptions struct {
	AccountKeeper          evmtypes.AccountKeeper
	BankKeeper             evmtypes.BankKeeper
	ExtensionOptionChecker authante.ExtensionOptionChecker
	IBCKeeper              *ibckeeper.Keeper
	EvmKeeper              EVMKeeper
	FeegrantKeeper         authante.FeegrantKeeper
	SignModeHandler        authsigning.SignModeHandler
	SigGasConsumer         func(meter sdk.GasMeter, sig signing.SignatureV2, params authtypes.Params) error
	MaxEthTxGasWanted      uint64
	TxFeeChecker           authante.TxFeeChecker
}

func (options HandlerOptions) Validate() error {
	if options.AccountKeeper == nil {
		return errors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return errors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return errors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.EvmKeeper == nil {
		return errors.Wrap(sdkerrors.ErrLogic, "evm keeper is required for AnteHandler")
	}
	return nil
}

func newEthAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewEthSetUpContextDecorator(options.EvmKeeper), // outermost AnteDecorator. SetUpContext must be called first
		NewEthMempoolFeeDecorator(options.EvmKeeper),   // Check eth effective gas price against minimal-gas-prices
		NewEthValidateBasicDecorator(options.EvmKeeper),
		NewEthTxPayloadVerificationDecorator(),
		NewEthTxCosmosMsgVerifierDecorator(options.EvmKeeper), // NEW cosmos msg verifier, loop protector
		NewEthSigVerificationDecorator(options.EvmKeeper),
		NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper),
		NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxEthTxGasWanted),
		NewCanTransferDecorator(options.EvmKeeper),
		NewEthIncrementSenderSequenceDecorator(options.AccountKeeper), // innermost AnteDecorator.
	)
}

func newCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		RejectMessagesDecorator{}, // reject MsgEthereumTxs
		authante.NewSetUpContextDecorator(),
		authante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		authante.NewValidateBasicDecorator(),
		authante.NewTxTimeoutHeightDecorator(),
		authante.NewValidateMemoDecorator(options.AccountKeeper),
		authante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		authante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewSetPubKeyDecorator(options.AccountKeeper),
		authante.NewValidateSigCountDecorator(options.AccountKeeper),
		authante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		authante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		authante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	)
}
