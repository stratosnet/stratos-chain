package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagPubKey      = "pubkey"
	FlagAmount      = "amount"
	FlagNetworkAddr = "network-addr"

	FlagMoniker         = "moniker"
	FlagIdentity        = "identity"
	FlagWebsite         = "website"
	FlagSecurityContact = "security-contact"
	FlagDetails         = "details"
)

// common flagsets to add to various functions
var (
	FsPk                = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount            = flag.NewFlagSet("", flag.ContinueOnError)
	FsNetworkAddr       = flag.NewFlagSet("", flag.ContinueOnError)
	fsDescriptionCreate = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsPk.String(FlagPubKey, "", "The Bech32 encoded PubKey of the node")
	FsAmount.String(FlagAmount, "", "Amount of coins to bond")
	FsNetworkAddr.String(FlagNetworkAddr, "", "The network address of the node")

	fsDescriptionCreate.String(FlagMoniker, "", "The node's name")
	fsDescriptionCreate.String(FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	fsDescriptionCreate.String(FlagWebsite, "", "The node's (optional) website")
	fsDescriptionCreate.String(FlagSecurityContact, "", "The node's (optional) security contact email")
	fsDescriptionCreate.String(FlagDetails, "", "The node's (optional) details")
}
