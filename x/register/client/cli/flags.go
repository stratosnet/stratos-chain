package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagPubKey     = "pubkey"
	FlagAmount     = "amount"
	FlagStakeDelta = "stake-delta"
	FlagIncrStake  = "incr-stake"
	FlagNodeType   = "node-type"

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
	FsPk         = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount     = flag.NewFlagSet("", flag.ContinueOnError)
	FsStakeDelta = flag.NewFlagSet("", flag.ContinueOnError)
	FsIncrStake  = flag.NewFlagSet("", flag.ContinueOnError)
	//FsNetworkAddr             = flag.NewFlagSet("", flag.ContinueOnError)
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
	FsStakeDelta.String(FlagStakeDelta, "", "Stake change of coins to be made (always positive like 100000ustos)")
	FsIncrStake.String(FlagIncrStake, "", "Boolean indicator of increase/decrease of stake delta, true for increase and false for decrease")
	//FsNetworkAddr.String(FlagNetworkAddress, "", "The network address of the node")
	//FsNetworkAddr.String(FlagNetworkAddr, "", "The network address of the node")
	FsNodeType.Int(FlagNodeType, 0, "The value of node_type is determined by the three node "+
		"types (storage=4/database=2/computation=1) and their arbitrary combinations.")

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
