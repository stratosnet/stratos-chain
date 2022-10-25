package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	stratos "github.com/stratosnet/stratos-chain/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
)

const (
	flagGenMetaNodeDir = "gen-meta-node-dir"
)

// AddGenesisMetaNodeCmd returns add-genesis-meta-node cobra Command.
func AddGenesisMetaNodeCmd(
	genBalancesIterator genutiltypes.GenesisBalancesIterator,
	defaultNodeHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-genesis-meta-node",
		Short: "Add a genesis meta node to genesis.json",
		Long: `Add a genesis meta node to genesis.json. If a node name is given,
the address will be looked up in the local Keybase.
`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			genMetaNodesDir := viper.GetString(flagGenMetaNodeDir)
			if genMetaNodesDir == "" {
				genMetaNodesDir = filepath.Join(config.RootDir, "config", "genmetanodes")
			}

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			appMetaNodes, err := getMetaNodeInfoFromFile(clientCtx.Codec, genMetaNodesDir, *genDoc, genBalancesIterator)
			if err != nil {
				return fmt.Errorf("failed to get meta node from file: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			registerGenState := registertypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			if registerGenState.GetMetaNodes() == nil {
				registerGenState.MetaNodes = registertypes.MetaNodes{}
			}

			for i, _ := range appMetaNodes {
				registerGenState.MetaNodes = append(registerGenState.MetaNodes, appMetaNodes[i])
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
	cmd.Flags().String(flagGenMetaNodeDir, "", "directory of genesis meta nodes info")
	return cmd
}

func getMetaNodeInfoFromFile(cdc codec.Codec, genMetaNodesDir string, genDoc tmtypes.GenesisDoc, genBalanceIterator genutiltypes.GenesisBalancesIterator,
) (appGenMetaNodes []registertypes.MetaNode, err error) {
	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genMetaNodesDir)
	if err != nil {
		return appGenMetaNodes, err
	}

	var appState map[string]json.RawMessage

	if err := json.Unmarshal(genDoc.AppState, &appState); err != nil {
		return appGenMetaNodes, err
	}

	balanceMap := make(map[string]exported.GenesisBalance)

	genBalanceIterator.IterateGenesisBalances(cdc, appState,
		func(balance exported.GenesisBalance) (stop bool) {
			balanceMap[balance.GetAddress().String()] = balance
			return false
		},
	)

	for _, fo := range fos {
		filename := filepath.Join(genMetaNodesDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}
		// get the node info
		var jsonRawMetaNode []byte
		if jsonRawMetaNode, err = ioutil.ReadFile(filename); err != nil {
			return appGenMetaNodes, err
		}

		var genMetaNode registertypes.GenesisMetaNode
		if err = cdc.UnmarshalJSON(jsonRawMetaNode, &genMetaNode); err != nil {
			return appGenMetaNodes, err
		}

		metaNode, err := genMetaNode.ToMetaNode()
		if err != nil {
			return appGenMetaNodes, err
		}

		appGenMetaNodes = append(appGenMetaNodes, metaNode)

		ownerAddrStr := metaNode.GetOwnerAddress()
		ownerBalance, ok := balanceMap[ownerAddrStr]
		if !ok {
			return appGenMetaNodes, fmt.Errorf(
				"account %v not in genesis.json: %+v", ownerAddrStr, balanceMap)
		}

		if ownerBalance.GetCoins().AmountOf(stratos.Wei).LT(metaNode.Tokens) {
			return appGenMetaNodes, fmt.Errorf(
				"insufficient fund for delegation %v: %v < %v",
				ownerBalance.GetAddress(), ownerBalance.GetCoins(), metaNode.Tokens,
			)
		}

		fmt.Println("Add meta node: " + metaNode.GetNetworkAddress() + " success.")
	}

	return appGenMetaNodes, nil
}
