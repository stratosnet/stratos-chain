## Statos-Chain

`Stratos` is the first decentralized data architecture that provides scalable, reliable, self-balanced storage, database and computation network, and offers a solid foundation for data processing.
`Stratos-Chain` is a Golang implementation of the `Stratos` protocol.

[![Go Report Card](https://goreportcard.com/badge/github.com/stratosnet/stratos-chain)](https://goreportcard.com/badge/github.com/stratosnet/stratos-chain)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Building the source

Prerequisites:
* [Go 1.15+](https://golang.org/doc/install)
* [git](https://github.com/git-guides/install-git)
* [wget](https://phoenixnap.com/kb/wget-command-with-examples)

Platform:
* Linux(tested on Ubuntu18.04)
```bash
git clone https://github.com/stratosnet/stratos-chain.git
cd stratos-chain
make build
```
* Mac OS
```bash
git clone https://github.com/stratosnet/stratos-chain.git
cd stratos-chain
make build-mac
```

Then the binary can be found in ./build folder

## Executables

The `Stratos-Chain` comes with 2 types of executables that can be found in the `./build` directory.

|    Command          | Description        |
| :-----------:     | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|  `stchaincli`   | the client end. It is the command line interface for interacting with `stchaind`. It is the entry point into the Stratos network (main-, test- or private net), capable of running as a full node (default), archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the Stratos network via JSON RPC endpoints. Use `stchaincli --help` and the [stchaincli Index](https://github.com/stratosnet/stratos-chain/wiki/SC-Basic-Transaction-and-Query-Commands) for command line options. |
|   `stchaind`   | the app Daemon (server). Use `stchaind --help` and the [stchaind index](https://github.com/stratosnet/stratos-chain/wiki/SC-Basic-Transaction-and-Query-Commands) for command line options. |


### `stchaincli`
```
Usage:
  stchaincli [command]
```

#### Available Commands
```
Commands:
  status      Query remote node for status
  config      Create or query an application CLI configuration file
  query       Querying subcommands
  tx          Transactions subcommands
  rest-server Start LCD (light-client daemon), a local REST server
  keys        Add or view local private keys
  version     Print the app version
  help        Help about any command
```

### `stchaind`

```
Usage:
  stchaind [command]
```

#### Available Commands
```
Commands:
  init                Initialize private validator, p2p, genesis, and application configuration files
  collect-gentxs      Collect genesis txs and output a genesis.json file
  migrate             Migrate genesis to a specified target version
  gentx               Generate a genesis tx carrying a self delegation
  validate-genesis    validates the genesis file at the default location or at the location passed as an arg
  add-genesis-account Add a genesis account to genesis.json
  faucet              Run a faucet cmd
  debug               Tool for helping with debugging your application
  start               Run the full node
  unsafe-reset-all    Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state
  tendermint          Tendermint subcommands
  export              Export state to JSON
  version             Print the app version
  help                Help about any command
```

## Connect to `Stratos` Network

Going through all the possible command line flags is out of scope here,
but we've enumerated a few common parameter combos to get you up to speed quickly
on how you can run your own `Stratos` instance.

### A Full node on the test Stratos network
For prerequisites and detailed instructions of connecting to test network, please read the [Connect to TestNet](https://github.com/stratosnet/stratos-chain-testnet).

### Full node on the main Stratos network - TBA
Prerequisites and detailed instructions of main network will be added later.

## Contribution

Thank you for considering to help out with the source code! We welcome contributions
from anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to stratos-chain, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting)
   guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary)
   guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "eth, rpc: make trace configs optional"

## License

Copyright 2021 Stratos

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the [License](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
