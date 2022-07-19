package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagFileHash = "file-hash"
	FlagReporter = "reporter"
	FlagUploader = "uploader"
)

func flagSetFileHash() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagFileHash, "", "The hash of uploaded file")
	return fs
}

func flagSetReporter() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagReporter, "", "The reporter address of meta node that reported the file")
	return fs
}

func flagSetUploader() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagUploader, "", "The owner address of resource node that uploaded the file")
	return fs
}
