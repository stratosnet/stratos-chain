package app

import (
	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	StratosBech32Prefix = "st"

	AccountPubKeyPrefix    = StratosBech32Prefix + "pub"
	ValidatorAddressPrefix = StratosBech32Prefix + "valoper"
	ValidatorPubKeyPrefix  = StratosBech32Prefix + "valoperpub"
	ConsNodeAddressPrefix  = StratosBech32Prefix + "valcons"
	ConsNodePubKeyPrefix   = StratosBech32Prefix + "valconspub"
	SdsNodeP2PKeyPrefix    = StratosBech32Prefix + "sdsp2p"
)

func SetConfig() {
	config := stratos.GetConfig()
	config.SetBech32PrefixForAccount(StratosBech32Prefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.SetBech32PrefixForSdsNodeP2P(SdsNodeP2PKeyPrefix)
	config.Seal()
}
