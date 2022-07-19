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

func flagSetDescriptionCreate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagMoniker, "", "The node's name")
	fs.String(FlagIdentity, "", "The (optional) identity signature (ex. UPort or Keybase)")
	fs.String(FlagWebsite, "", "The node's (optional) website")
	fs.String(FlagSecurityContact, "", "The node's (optional) security contact email")
	fs.String(FlagDetails, "", "The node's (optional) details")

	return fs
}

// FlagSetAmount Returns the FlagSet for amount related operations.
func flagSetAmount() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagAmount, "", "Amount of coins to bond")
	return fs
}

// FlagSetPublicKey Returns the flagset for Public Key related operations.
func flagSetPublicKey() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagPubKey, "", "The resource node's Protobuf JSON encoded public key")
	return fs
}

// FlagSetNetworkAddress Returns the flagset for network address of resource node
func flagSetNetworkAddress() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagNetworkAddress, "", "The address of the PP node")
	return fs
}

// FlagSetNodeType Returns the flagset for node type of resource node
func flagSetNodeType() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Uint32(FlagNodeType, 0, "The value of node_type is determined by the three node types (storage=4/database=2/computation=1) and their arbitrary combinations.")
	return fs
}

func flagSetStakeUpdate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagStakeDelta, "", "Stake change of coins to be made (always positive like 100000ustos)")
	fs.String(FlagIncrStake, "", "Boolean indicator of increase/decrease of stake delta, true for increase and false for decrease")

	return fs
}

func flagSetVoting() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagCandidateNetworkAddress, "The network address of the candidate PP node", "")
	fs.String(FlagCandidateOwnerAddress, "The owner address of the candidate PP node", "")
	fs.Bool(FlagOpinion, false, "Opinion of the vote for the registration of Meta node.")
	fs.String(FlagVoterNetworkAddress, "The address of the PP node that made the vote.", "")
	return fs
}
