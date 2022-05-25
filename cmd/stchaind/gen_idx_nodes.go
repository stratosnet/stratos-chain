package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	defaultDemon      = "ustos"
	flagGenIdxNodeDir = "gen-idx-node-dir"
)

// AddGenesisIndexingNodeCmd returns add-genesis-indexing-node cobra Command.
func AddGenesisIndexingNodeCmd(
	genBalancesIterator genutiltypes.GenesisBalancesIterator,
	defaultNodeHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-genesis-indexing-node",
		Short: "Add a genesis indexing node to genesis.json",
		Long: `Add a genesis indexing node to genesis.json. If a node name is given,
the address will be looked up in the local Keybase.
`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			genIdxNodesDir := viper.GetString(flagGenIdxNodeDir)
			if genIdxNodesDir == "" {
				genIdxNodesDir = filepath.Join(config.RootDir, "config", "genidxnodes")
			}

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			appIdxNodes, err := getIndexingNodeInfoFromFile(clientCtx.Codec, genIdxNodesDir, *genDoc, genBalancesIterator)
			if err != nil {
				return fmt.Errorf("failed to get indexing node from file: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			registerGenState := registertypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			for _, appIdxNode := range appIdxNodes {
				registerGenState.IndexingNodes.IndexingNodes = append(registerGenState.IndexingNodes.IndexingNodes, &appIdxNode)
			}

			registerGenStateBz, err := clientCtx.Codec.MarshalJSON(&registerGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal register genesis state: %w", err)
			}

			appState[registertypes.ModuleName] = registerGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagGenIdxNodeDir, "", "directory of genesis indexing nodes info")
	return cmd
}

func getIndexingNodeInfoFromFile(cdc codec.Codec, genIdxNodesDir string, genDoc tmtypes.GenesisDoc, genBalanceIterator genutiltypes.GenesisBalancesIterator,
) (appGenIdxNodes []registertypes.IndexingNode, err error) {
	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genIdxNodesDir)
	if err != nil {
		return appGenIdxNodes, err
	}

	var appState map[string]json.RawMessage

	if err := json.Unmarshal(genDoc.AppState, &appState); err != nil {
		return appGenIdxNodes, err
	}

	balanceMap := make(map[string]exported.GenesisBalance)

	genBalanceIterator.IterateGenesisBalances(cdc, appState,
		func(balance exported.GenesisBalance) (stop bool) {
			balanceMap[balance.GetAddress().String()] = balance
			return false
		},
	)

	for _, fo := range fos {
		filename := filepath.Join(genIdxNodesDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the node info
		var jsonRawIdxNode []byte
		if jsonRawIdxNode, err = ioutil.ReadFile(filename); err != nil {
			return appGenIdxNodes, err
		}

		var genIdxNode registertypes.GenesisIndexingNode
		if err = json.Unmarshal(jsonRawIdxNode, &genIdxNode); err != nil {
			return appGenIdxNodes, err
		}

		indexingNode := genIdxNode.ToIndexingNode()
		appGenIdxNodes = append(appGenIdxNodes, indexingNode)

		ownerAddrStr := indexingNode.GetOwnerAddress()
		ownerBalance, ok := balanceMap[ownerAddrStr]
		if !ok {
			return appGenIdxNodes, fmt.Errorf(
				"account %v not in genesis.json: %+v", ownerAddrStr, balanceMap)
		}

		if ownerBalance.GetCoins().AmountOf(defaultDemon).LT(indexingNode.Tokens) {
			return appGenIdxNodes, fmt.Errorf(
				"insufficient fund for delegation %v: %v < %v",
				ownerBalance.GetAddress(), ownerBalance.GetCoins().AmountOf(defaultDemon), indexingNode.Tokens,
			)
		}

		fmt.Println("Add indexing node: " + indexingNode.GetNetworkAddr() + " success.")
	}

	return appGenIdxNodes, nil
}
