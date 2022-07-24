## Stratos-Chain

`Stratos` is the first decentralized data architecture that provides scalable, reliable, self-balanced storage, database and computation network, and offers a solid foundation for data processing.
`Stratos-Chain` is a Golang implementation of the `Stratos` protocol.

[![Go Report Card](https://goreportcard.com/badge/github.com/stratosnet/stratos-chain)](https://goreportcard.com/badge/github.com/stratosnet/stratos-chain)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

--- ---

## Recommended hardware/software

- <b>Recommended Hardware</b>

        * CPU           i5 (4 cores)
        * RAM           16GB
        * Hard disk     2TB


- <b>Software(tested version)</b>

        * Ubuntu 18.04+
        * Go 1.16+ linux/amd64


- <b>Platform</b>

    * Linux(tested on Ubuntu18.04 and 20.04)
    * Mac OS
    * Windows

      It is possible to build and run the software on Windows. However, we did not test it on Windows completely.
      It may give you unexpected results, or it may require additional setup.

      An alternative option is to install a separate virtual Linux system using [VirtualBox](https://www.virtualbox.org/wiki/Downloads) or [VMware Workstation](https://www.vmware.com/ca/products/workstation-player/workstation-player-evaluation.html)

--- ---

## Connect to `Stratos` Network

### A Full node on the test Stratos network
For prerequisites and detailed instructions of connecting to `Tropos Incentive Testnet` network, please refer to [Connecting to Tropos Incentive Testnet](https://github.com/stratosnet/sds/wiki/Tropos-Incentive-Testnet).

### Full node on the main Stratos network - TBA
Prerequisites and detailed instructions of main network will be added later.

--- ---

## Stratos Explorer

* https://explorer-tropos.thestratos.org/

---

## References

<details>
    <summary><b><code>Stratos-chain</code> document List</b></summary>

<br>

* [Tropos Incentive Testnet](https://github.com/stratosnet/sds/wiki/Tropos-Incentive-Testnet)
 
* ['stchaind' Commands(part1)](https://github.com/stratosnet/stratos-chain/wiki/Stratos-Chain-%60stchaind%60-Commands(part1))

* [stchaind' Commands(part2)](https://github.com/stratosnet/stratos-chain/wiki/Stratos-Chain-%60stchaind%60-Commands(part2))

* [gRPC Queries](https://github.com/stratosnet/stratos-chain/wiki/Stratos-Chain-gRPC-Queries)

* [REST APIs](https://github.com/stratosnet/stratos-chain/wiki/Stratos-Chain-REST-APIs)

* [How to become a validator](https://github.com/stratosnet/stratos-chain/wiki/How-to-Become-a-Validator)

</details>

--- ---

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
* Pull requests need to be based on and opened against the `main` branch.
* Commit messages should be prefixed with the package(s) they modify.
  * E.g. "eth, rpc: make trace configs optional"

--- ---

## License

Copyright 2022 Stratos

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the [License](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
