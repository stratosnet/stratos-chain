package main

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stratosnet/stratos-chain/x/register"
	registerexported "github.com/stratosnet/stratos-chain/x/register/exported"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	defaultDemon      = "ustos"
	flagGenIdxNodeDir = "gen-idx-node-dir"
)

// GenesisAccountsIterator defines the expected iterating genesis accounts object (noalias)
type GenesisAccountsIterator interface {
	IterateGenesisAccounts(
		cdc *codec.Codec,
		appGenesis map[string]json.RawMessage,
		iterateFn func(authexported.Account) (stop bool),
	)
}

func getIndexingNodeInfoFromFile(cdc *codec.Codec, genIdxNodesDir string, genDoc tmtypes.GenesisDoc, genAccIterator GenesisAccountsIterator,
) (appGenIdxNodes registerexported.GenesisIndexingNodes, err error) {
	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genIdxNodesDir)
	if err != nil {
		return appGenIdxNodes, err
	}

	// prepare a map of all accounts in genesis state to then validate
	// against the validators addresses
	var appState map[string]json.RawMessage
	if err := cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
		return appGenIdxNodes, err
	}

	addrMap := make(map[string]authexported.Account)
	genAccIterator.IterateGenesisAccounts(cdc, appState,
		func(acc authexported.Account) (stop bool) {
			addrMap[acc.GetAddress().String()] = acc
			return false
		},
	)

	for _, fo := range fos {
		filename := filepath.Join(genIdxNodesDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawIdxNode []byte
		if jsonRawIdxNode, err = ioutil.ReadFile(filename); err != nil {
			return appGenIdxNodes, err
		}

		var genIdxNode registerexported.GenesisIndexingNode
		if err = cdc.UnmarshalJSON(jsonRawIdxNode, &genIdxNode); err != nil {
			return appGenIdxNodes, err
		}
		appGenIdxNodes = append(appGenIdxNodes, genIdxNode)

		if err := genIdxNode.Validate(); err != nil {
			return appGenIdxNodes, fmt.Errorf("failed to validate new genesis account: %w", err)
		}

		ownerAddrStr := genIdxNode.GetOwnerAddr().String()
		ownerAccount, ownerOk := addrMap[ownerAddrStr]
		if !ownerOk {
			return appGenIdxNodes, fmt.Errorf(
				"account %v not in genesis.json: %+v", ownerAccount, addrMap)
		}

		if ownerAccount.GetCoins().AmountOf(defaultDemon).LT(genIdxNode.GetTokens()) {
			return appGenIdxNodes, fmt.Errorf(
				"insufficient fund for delegation %v: %v < %v",
				ownerAccount.GetAddress(), ownerAccount.GetCoins().AmountOf(defaultDemon), genIdxNode.GetTokens(),
			)
		}
	}

	return appGenIdxNodes, nil
}

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisIndexingNodeCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string, genAccIterator GenesisAccountsIterator,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-genesis-indexing-node [node_name]",
		Short: "Add a genesis indexing node to genesis.json",
		Long: `Add a genesis indexing node to genesis.json. If a node name is given,
the address will be looked up in the local Keybase.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			genIdxNodesDir := viper.GetString(flagGenIdxNodeDir)
			if genIdxNodesDir == "" {
				genIdxNodesDir = filepath.Join(config.RootDir, "config", "genidxnodes")
			}

			appIdxNodes, err := getIndexingNodeInfoFromFile(cdc, genIdxNodesDir, *genDoc, genAccIterator)

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			registerGenState := register.GetGenesisStateFromAppState(cdc, appState)

			for _, appIdxNode := range appIdxNodes {
				registerGenState.IndexingNodes = append(registerGenState.IndexingNodes, appIdxNode)
				registerGenState.LastIndexingNodeStakes = append(registerGenState.LastIndexingNodeStakes,
					register.LastIndexingNodeStake{Address: appIdxNode.GetNetworkAddr(), Stake: appIdxNode.GetTokens()})
			}

			registerGenStateBz, err := cdc.MarshalJSON(registerGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal register genesis state: %w", err)
			}

			appState[register.ModuleName] = registerGenStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")

	return cmd
}
