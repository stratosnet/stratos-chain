package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagReporterAddr    = "reporter-addr"
	FlagEpoch           = "epoch"
	FlagReportReference = "reference"
	FlagWalletVolumes   = "wallet-volumes"
	FlagAmount          = "amount"
	FlagWalletAddress   = "wallet-address"
	FlagTargetAddress   = "target-address"
	FlagReporters       = "reporters"
	FlagReporterOwner   = "reporter-owner"
	FlagNetworkAddress  = "network-address"
	FlagSlashing        = "slashing"
	FlagSuspend         = "suspend"
	FlagBLSSignature    = "bls-signature"
)

var (
//FsReporterAddr    = flag.NewFlagSet("", flag.ContinueOnError)
//FsEpoch           = flag.NewFlagSet("", flag.ContinueOnError)
//FsReportReference = flag.NewFlagSet("", flag.ContinueOnError)
//FsWalletVolumes   = flag.NewFlagSet("", flag.ContinueOnError)
//FsAmount          = flag.NewFlagSet("", flag.ContinueOnError)
//FsWalletAddress   = flag.NewFlagSet("", flag.ContinueOnError)
//FsTargetAddress   = flag.NewFlagSet("", flag.ContinueOnError)
//FsReporters       = flag.NewFlagSet("", flag.ContinueOnError)
//FsReportOwner     = flag.NewFlagSet("", flag.ContinueOnError)
//FsNetworkAddress  = flag.NewFlagSet("", flag.ContinueOnError)
//FsSlashing        = flag.NewFlagSet("", flag.ContinueOnError)
//FsSuspend         = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	//FsReporterAddr.String(FlagReporterAddr, "", "the node address of reporter")
	//FsEpoch.String(FlagEpoch, "", "the epoch when this PoT message reported.")
	//FsReportReference.String(FlagReportReference, "", " the hash used as a reference to this PoT report")
	//FsWalletVolumes.String(FlagWalletVolumes, "", "a string of KEY-VALUE pairs. The KEY is 'wallet-volumes' and the VALUE is the proof of traffic of this wallet`")
	//FsAmount.String(FlagAmount, "", "Amount of coins to withdraw")
	//FsWalletAddress.String(FlagWalletAddress, "", "The address of the wallet to withdraw")
	//FsTargetAddress.String(FlagTargetAddress, "", "The target account where the money is deposited after withdraw")
	//FsReporters.String(FlagReporters, "", "the node address list of reporters")
	//FsReportOwner.String(FlagReporterOwner, "", "the node address list of reporters")
	//FsNetworkAddress.String(FlagNetworkAddress, "", "the node address of resource node to slashing")
	//FsSlashing.String(FlagSlashing, "", "the amount of slashing")
	//FsSuspend.Bool(FlagSuspend, false, "if the resource node is suspend")

}

// FlagSetAmount Returns the FlagSet for amount related operations.
func flagSetAmount() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagAmount, "", "Amount of coins to withdraw")
	return fs
}

// flagSetReportVolumes Returns the FlagSet for report volumes.
func flagSetReportVolumes() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagWalletVolumes, "", "a string of KEY-VALUE pairs. The KEY is 'wallet-volumes' and the VALUE is the proof of traffic of this wallet`")
	fs.String(FlagReporterAddr, "", "the node address of reporter")
	fs.String(FlagEpoch, "", "the epoch when this PoT message reported.")
	fs.String(FlagReportReference, "", " the hash used as a reference to this PoT report")
	fs.String(FlagBLSSignature, "", " BLS signature")

	return fs
}

// flagSetEpoch Returns the FlagSet for epoch.
func flagSetEpoch() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagEpoch, "", "the epoch when this PoT message reported.")

	return fs
}

// flagSetReportersAndOwners Returns the FlagSet for reporters and their owners.
func flagSetReportersAndOwners() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagReporters, "", "the node address list of reporters")
	fs.String(FlagReporterOwner, "", "the owner address list of reporter")

	return fs
}

// flagSetWalletAddress Returns the FlagSet for wallet address.
func flagSetWalletAddress() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagWalletAddress, "", "The wallet address to withdraw from")
	return fs
}

// flagSetTargetAddress Returns the FlagSet for target wallet address.
func flagSetTargetAddress() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagTargetAddress, "", "The target wallet address to deposit into after withdrawing")
	return fs
}

// FlagSetNetworkAddress Returns the FlagSet for network address of resource node
func flagSetNetworkAddress() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagNetworkAddress, "", "The address of the PP node")
	return fs
}

// flagSetFsSlashing Returns the FlagSet for slashing amount.
func flagSetSlashing() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagSlashing, "", "the amount of slashing")
	return fs
}

// flagSetSuspend Returns the FlagSet for suspend state of resource node.
func flagSetSuspend() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(FlagSuspend, false, "if the resource node is suspend")
	return fs
}
