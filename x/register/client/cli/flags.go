package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagPubKey      = "pubkey"
	FlagAmount      = "amount"
	FlagNetworkAddr = "network-addr"
	FlagNetworkID   = "network-id"
	FlagNodeType    = "node-type"

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
	FsNetworkID         = flag.NewFlagSet("", flag.ContinueOnError)
	FsNodeType          = flag.NewFlagSet("", flag.ContinueOnError)
	FsDescriptionCreate = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsPk.String(FlagPubKey, "", "The Bech32 encoded PubKey of the node")
	FsAmount.String(FlagAmount, "", "Amount of coins to bond")
	FsNetworkAddr.String(FlagNetworkAddr, "", "The network address of the node")
	FsNetworkID.String(FlagNetworkID, "", "The network id of the node")
	FsNodeType.Int(FlagNodeType, 0, `The value of node_type is determined by the three node types('storage', 'database'', and 'computation') and their arbitrary combinations.
Suppose, we define:
	computation 	= 1,
	database 	= 2,
	storage 	= 4,
Then, their combinations:
	computation && database 			= 1 + 2 = 3,
	computation && storage 				= 1 + 4 = 5,
	database && storage 				= 2 + 4 = 6,
	computation && database && storage 		= 1 + 2 + 4 = 7,
As a result, the value of node_type should be one of the following digits:
	1:  "computation",
	2:  "database",
	3:  "computation/database",
	4:  "storage",
	5:  "computation/storage",
	6:  "database/storage",
	7:  "computation/database/storage"`)

	FsDescriptionCreate.String(FlagMoniker, "", "The node's name")
	FsDescriptionCreate.String(FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	FsDescriptionCreate.String(FlagWebsite, "", "The node's (optional) website")
	FsDescriptionCreate.String(FlagSecurityContact, "", "The node's (optional) security contact email")
	FsDescriptionCreate.String(FlagDetails, "", "The node's (optional) details")

}
