package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagPubKey       = "pubkey"
	FlagAmount       = "amount"
	FlagDepositDelta = "deposit-delta"
	FlagIncrDeposit  = "incr-deposit"
	FlagNodeType     = "node-type"

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
	FlagBeneficiaryAddress      = "beneficiary-address"
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

func flagSetDepositUpdate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagDepositDelta, "", "Deposit change of coins to be made (always positive like 100000wei)")
	fs.String(FlagIncrDeposit, "", "Boolean indicator of increase/decrease of deposit delta, true for increase and false for decrease")

	return fs
}

func flagSetVoting() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagCandidateNetworkAddress, "", "The network address of the candidate PP node")
	fs.String(FlagCandidateOwnerAddress, "", "The owner address of the candidate PP node")
	fs.Bool(FlagOpinion, false, "Opinion of the vote for the registration of Meta node.")
	fs.String(FlagVoterNetworkAddress, "", "The address of the PP node that made the vote.")
	return fs
}

func flagSetBeneficiaryAddress() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagBeneficiaryAddress, "", "The beneficiary address of the meta node")
	return fs
}
