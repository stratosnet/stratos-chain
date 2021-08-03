package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/stratosnet/stratos-chain/x/evm/types"
)

// GetQueryCmd defines evm module queries through the cli
func GetQueryCmd(moduleName string, cdc *codec.Codec) *cobra.Command {
	evmQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the evm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	evmQueryCmd.AddCommand(flags.GetCommands(
	//GetCmdGetStorageAt(moduleName, cdc),
	//GetCmdGetCode(moduleName, cdc),
	)...)
	return evmQueryCmd
}
