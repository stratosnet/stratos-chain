package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagReporter        = "reporter-addr"
	FlagEpoch           = "epoch"
	FlagReportReference = "reference"
	FlagNodesVolume     = "nodes-volume"
	FlagAmount          = "amount"
	FlagNodeAddress     = "node-address"
)

var (
	FsReporter        = flag.NewFlagSet("", flag.ContinueOnError)
	FsEpoch           = flag.NewFlagSet("", flag.ContinueOnError)
	FsReportReference = flag.NewFlagSet("", flag.ContinueOnError)
	FsNodesVolume     = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount          = flag.NewFlagSet("", flag.ContinueOnError)
	FsNodeAddress     = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsReporter.String(FlagReporter, "", "the node address of reporter")
	FsEpoch.Int64(FlagEpoch, 0, "the epoch when this PoT message reported.")
	FsReportReference.String(FlagReportReference, "", " the hash used as a reference to this PoT report")
	FsNodesVolume.String(FlagNodesVolume, "", "a string of KEY-VALUE pairs. The KEY is 'node_address' and the VALUE is the proof of traffic of this node")
	FsAmount.String(FlagAmount, "", "Amount of coins to withdraw")
	FsNodeAddress.String(FlagNodeAddress, "", "The address of the node to withdraw")
}
