package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagPubKey    = "pubkey"
	FlagAmount    = "amount"
	FlagNetworkID = "network-id"
	FlagNodeType  = "node-type"

	FlagMoniker         = "moniker"
	FlagIdentity        = "identity"
	FlagWebsite         = "website"
	FlagSecurityContact = "security-contact"
	FlagDetails         = "details"

	FlagNetworkAddress          = "network-address"
	FlagCandidateOwnerAddress   = "candidate-owner-address"
	FlagCandidateNetworkAddress = "candidate-network-address"
	FlagOpinion                 = "opinion"
	FlagVoterNetworkAddress     = "voter-network-address"
)

// common flagsets to add to various functions
var (
	FsPk                      = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount                  = flag.NewFlagSet("", flag.ContinueOnError)
	FsNetworkID               = flag.NewFlagSet("", flag.ContinueOnError)
	FsNodeType                = flag.NewFlagSet("", flag.ContinueOnError)
	FsDescription             = flag.NewFlagSet("", flag.ContinueOnError)
	FsNetworkAddress          = flag.NewFlagSet("", flag.ContinueOnError)
	FsCandidateNetworkAddress = flag.NewFlagSet("", flag.ContinueOnError)
	FsCandidateOwnerAddress   = flag.NewFlagSet("", flag.ContinueOnError)
	FsOpinion                 = flag.NewFlagSet("", flag.ContinueOnError)
	FsVoterNetworkAddress     = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsPk.String(FlagPubKey, "", "The Bech32 encoded PubKey of the node")
	FsAmount.String(FlagAmount, "", "Amount of coins to bond")
	//FsNetworkAddr.String(FlagNetworkAddress, "", "The network address of the node")
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

	FsDescription.String(FlagMoniker, "", "The node's name")
	FsDescription.String(FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	FsDescription.String(FlagWebsite, "", "The node's (optional) website")
	FsDescription.String(FlagSecurityContact, "", "The node's (optional) security contact email")
	FsDescription.String(FlagDetails, "", "The node's (optional) details")

	FsNetworkAddress.String(FlagNetworkAddress, "The address of the PP node", "")
	FsCandidateNetworkAddress.String(FlagCandidateNetworkAddress, "The network address of the candidate PP node", "")
	FsCandidateOwnerAddress.String(FlagCandidateOwnerAddress, "The owner address of the candidate PP node", "")
	FsOpinion.Bool(FlagOpinion, false, "Opinion of the vote for the registration of Indexing node.")
	FsVoterNetworkAddress.String(FlagVoterNetworkAddress, "The address of the PP node that made the vote.", "")
}
