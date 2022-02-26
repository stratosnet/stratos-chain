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
)

var (
	FsReporterAddr    = flag.NewFlagSet("", flag.ContinueOnError)
	FsEpoch           = flag.NewFlagSet("", flag.ContinueOnError)
	FsReportReference = flag.NewFlagSet("", flag.ContinueOnError)
	FsWalletVolumes   = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount          = flag.NewFlagSet("", flag.ContinueOnError)
	FsWalletAddress   = flag.NewFlagSet("", flag.ContinueOnError)
	FsTargetAddress   = flag.NewFlagSet("", flag.ContinueOnError)
	FsReporters       = flag.NewFlagSet("", flag.ContinueOnError)
	FsReportOwner     = flag.NewFlagSet("", flag.ContinueOnError)
	FsNetworkAddress  = flag.NewFlagSet("", flag.ContinueOnError)
	FsSlashing        = flag.NewFlagSet("", flag.ContinueOnError)
	FsSuspend         = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsReporterAddr.String(FlagReporterAddr, "", "the node address of reporter")
	FsEpoch.String(FlagEpoch, "", "the epoch when this PoT message reported.")
	FsReportReference.String(FlagReportReference, "", " the hash used as a reference to this PoT report")
	FsWalletVolumes.String(FlagWalletVolumes, "", "a string of KEY-VALUE pairs. The KEY is 'wallet-volumes' and the VALUE is the proof of traffic of this wallet`")
	FsAmount.String(FlagAmount, "", "Amount of coins to withdraw")
	FsWalletAddress.String(FlagWalletAddress, "", "The address of the wallet to withdraw")
	FsTargetAddress.String(FlagTargetAddress, "", "The target account where the money is deposited after withdraw")
	FsReporters.String(FlagReporters, "", "the node address list of reporters")
	FsReportOwner.String(FlagReporterOwner, "", "the node address list of reporters")
	FsNetworkAddress.String(FlagNetworkAddress, "", "the node address of resource node to slashing")
	FsSlashing.String(FlagSlashing, "", "the amount of slashing")
	FsSuspend.String(FlagSuspend, "", "if the resource node is suspend")

}
