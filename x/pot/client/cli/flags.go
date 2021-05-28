package cli

import flag "github.com/spf13/pflag"

const (
	FlagAmount      = "amount"
	FlagNodeAddress = "node-address"
)

var (
	FsAmount      = flag.NewFlagSet("", flag.ContinueOnError)
	FsNodeAddress = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsAmount.String(FlagAmount, "", "Amount of coins to withdraw")
	FsNodeAddress.String(FlagNodeAddress, "", "The address of the node to withdraw")
}
