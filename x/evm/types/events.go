package types

// Evm module events
const (
	EventTypeStratosTx  = TypeMsgStratosTx
	EventTypeEthereumTx = TypeMsgEthereumTx

	AttributeKeyContractAddress = "contract"
	AttributeKeyRecipient       = "recipient"
	AttributeValueCategory      = ModuleName
)
