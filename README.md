## Statos-Chain


`Stratos` is the first decentralized data architecture that provides scalable, reliable, self-balanced storage, database and computation network, and offers a solid foundation for data processing.
`Stratos-Chain` is a Golang implementation of the `Stratos` protocol.

[![Go Report Card](https://goreportcard.com/badge/github.com/stratosnet/stratos-chain)](https://goreportcard.com/badge/github.com/stratosnet/stratos-chain)

Automated builds are available for stable releases and the unstable master branch. Binary
archives are published at https://geth.ethereum.org/downloads/.

## Building the source

Prerequisites
* [Go 1.15+](https://golang.org/doc/install)
* [git](https://github.com/git-guides/install-git)
* [wget](https://phoenixnap.com/kb/wget-command-with-examples)

Platform
* Linux(tested on Ubuntu18.04)
* Mac OS

For details about building from the source code, please read the [Installation Instructions](https://github.com/stratosnet/stratos-chain-testnet/blob/main/README.md).

## Executables

The `Stratos-Chain` comes with several wrappers/executables that can be found in the `cmd` directory.


|    Command          | Description        |
| :-----------:     | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|  `stchaincli`   | the client end. It is the command line interface for interacting with `stchaind`. It is the entry point into the Stratos network (main-, test- or private net), capable of running as a full node (default), archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the Stratos network via JSON RPC endpoints. Use `stchaincli --help` and the [stchaincli Index](https://github.com/stratosnet/stratos-chain/wiki/SC-Basic-Transaction-and-Query-Commands) for command line options. |
|   `stchaind`   | the app Daemon (server). Use `stchaind --help` and the [stchaind Index](https://github.com/stratosnet/stratos-chain/wiki/SC-Basic-Transaction-and-Query-Commands) for command line options. |


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

#### Available Flags(options)

```
Flags:
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
  -h, --help              help for stchaincli
      --home string       directory for config and data (default "/home/node0/.stchaincli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             Print out full stack trace on errors
```

#### Getting detailed command Help Info
```
stchaincli [command] --help      More information about a 'stchaincli' command or sub-command 
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

#### Available Flags(options)
```
Flags:
  -h, --help                    help for stchaind
      --home string             directory for config and data (default "/home/hong/.stchaind")
      --inv-check-period uint   Assert registered invariants every N blocks
      --log_level string        Log level (default "main:info,state:info,*:error")
      --trace                   print out full stack trace on errors
```


#### Getting detailed command Help Info
```
stchaind [command] --help      More information about a 'stchaind' command or sub-command 
```
## Connect to `Stratos` Network

Going through all the possible command line flags is out of scope here,
but we've enumerated a few common parameter combos to get you up to speed quickly
on how you can run your own `stchaind` instance.

### A Full node on the test Stratos network
For prerequisites and detailed instructions of connecting to test network, please read the [Instructions to connect to TestNet](https://github.com/stratosnet/stratos-chain-testnet).

### Full node on the main Stratos network - TBA
For prerequisites and detailed instructions of connecting to test network, please read the [Instructions to connect to MainNet](https://github.com/stratosnet/stratos-chain-testnet).

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

The stratos-chain library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The stratos-chain binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
{"mode":"full","isActive":false}
