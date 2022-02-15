package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagReporter        = "reporter-addr"
	FlagEpoch           = "epoch"
	FlagReportReference = "reference"
	FlagWalletVolumes   = "wallet-volumes"
	FlagAmount          = "amount"
	FlagWalletAddress   = "wallet-address"
	FlagTargetAddress   = "target-address"
)

var (
	FsReporter        = flag.NewFlagSet("", flag.ContinueOnError)
	FsEpoch           = flag.NewFlagSet("", flag.ContinueOnError)
	FsReportReference = flag.NewFlagSet("", flag.ContinueOnError)
	FsWalletVolumes   = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount          = flag.NewFlagSet("", flag.ContinueOnError)
	FsWalletAddress   = flag.NewFlagSet("", flag.ContinueOnError)
	FsTargetAddress   = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsReporter.String(FlagReporter, "", "the node address of reporter")
	FsEpoch.String(FlagEpoch, "", "the epoch when this PoT message reported.")
	FsReportReference.String(FlagReportReference, "", " the hash used as a reference to this PoT report")
	FsWalletVolumes.String(FlagWalletVolumes, "", "a string of KEY-VALUE pairs. The KEY is 'wallet-volumes' and the VALUE is the proof of traffic of this wallet`")
	FsAmount.String(FlagAmount, "", "Amount of coins to withdraw")
	FsWalletAddress.String(FlagWalletAddress, "", "The address of the wallet to withdraw")
	FsTargetAddress.String(FlagTargetAddress, "", "The target account where the money is deposited after withdraw")
}
