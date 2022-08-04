# JSON-RPC Methods

## Endpoints

| Method                                               | Namespace | Implemented | Public | Notes              |
|:-----------------------------------------------------|:----------|:-----------:|:------:|:-------------------|
| [web3_clientVersion](01_web3.md)                     | Web3      |      ✔      |   ✔    |                    |
| [web3_sha3](01_web3.md)                              | Web3      |      ✔      |   ✔    |                    |
| [net_version](02_net.md)                             | Net       |      ✔      |   ✔    |                    |
| [net_listening](02_net.md)                           | Net       |      ✔      |   ✔    |                    |
| [net_peerCount](02_net.md)                           | Net       |      ✔      |   ✔    |                    |
| [eth_protocolVersion](03_eth.md)                     | Eth       |      ✔      |   ✔    |                    |
| [eth_syncing](03_eth.md)                             | Eth       |      ✔      |   ✔    |                    |
| [eth_coinbase](03_eth.md)                            | Eth       |      ✔      |        |                    |
| eth_mining                                           | Eth       |             |        |                    |
| [eth_gasPrice](03_eth.md)                            | Eth       |      ✔      |   ✔    |                    |
| [eth_accounts](03_eth.md)                            | Eth       |      ✔      |   ✔    |                    |
| [eth_blockNumber](03_eth.md)                         | Eth       |      ✔      |   ✔    |                    |
| [eth_getBalance](03_eth.md)                          | Eth       |      ✔      |   ✔    |                    |
| [eth_getStorageAt](03_eth.md)                        | Eth       |      ✔      |   ✔    |                    |
| [eth_getTransactionCount](03_eth.md)                 | Eth       |      ✔      |   ✔    |                    |
| [eth_getBlockTransactionCountByHash](03_eth.md)      | Eth       |      ✔      |   ✔    |                    |
| [eth_getBlockTransactionCountByNumber](03_eth.md)    | Eth       |      ✔      |   ✔    |                    |
| [eth_getCode](03_eth.md)                             | Eth       |      ✔      |   ✔    |                    |
| [eth_getProof](03_eth.md)                            | Eth       |      ✔      |        |                    |
| [eth_sign](03_eth.md)                                | Eth       |      ✔      |   ✔    |                    |
| eth_signTransaction                                  | Eth       |             |        |                    |
| [eth_sendTransaction](03_eth.md)                     | Eth       |      ✔      |   ✔    |                    |
| [eth_sendRawTransaction](03_eth.md)                  | Eth       |      ✔      |   ✔    |                    |
| [eth_call](03_eth.md)                                | Eth       |      ✔      |   ✔    |                    |
| [eth_estimateGas](03_eth.md)                         | Eth       |      ✔      |   ✔    |                    |
| [eth_getBlockByHash](03_eth.md)                      | Eth       |      ✔      |   ✔    |                    |
| [eth_getBlockByNumber](03_eth.md)                    | Eth       |      ✔      |   ✔    |                    |
| [eth_getTransactionByHash](03_eth.md)                | Eth       |      ✔      |   ✔    |                    |
| [eth_getTransactionByBlockHashAndIndex](03_eth.md)   | Eth       |      ✔      |   ✔    |                    |
| [eth_getTransactionByBlockNumberAndIndex](03_eth.md) | Eth       |      ✔      |        |                    |
| [eth_getTransactionReceipt](03_eth.md)               | Eth       |      ✔      |   ✔    |                    |
| eth_getCompilers                                     | Eth       |             |        |                    |
| eth_compileSolidity                                  | Eth       |             |        |                    |
| eth_compileLLL                                       | Eth       |             |        |                    |
| eth_compileSerpent                                   | Eth       |             |        |                    |
| [eth_newFilter](03_eth.md)                           | Eth       |      ✔      |   ✔    |                    |
| [eth_newBlockFilter](03_eth.md)                      | Eth       |      ✔      |   ✔    |                    |
| [eth_newPendingTransactionFilter](03_eth.md)         | Eth       |      ✔      |   ✔    |                    |
| [eth_uninstallFilter](03_eth.md)                     | Eth       |      ✔      |   ✔    |                    |
| [eth_getFilterChanges](03_eth.md)                    | Eth       |      ✔      |   ✔    |                    |
| [eth_getFilterLogs](03_eth.md)                       | Eth       |      ✔      |   ✔    |                    |
| [eth_getLogs](03_eth.md)                             | Eth       |      ✔      |   ✔    |                    |
| eth_hashrate                                         | Eth       |     N/A     |        | PoW-only           |
| eth_getUncleCountByBlockHash                         | Eth       |     N/A     |        | PoW-only           |
| eth_getUncleCountByBlockNumber                       | Eth       |     N/A     |        | PoW-only           |
| eth_getUncleByBlockHashAndIndex                      | Eth       |     N/A     |        | PoW-only           |
| eth_getUncleByBlockNumberAndIndex                    | Eth       |     N/A     |        | PoW-only           |
| eth_getWork                                          | Eth       |     N/A     |        | PoW-only           |           
| eth_submitWork                                       | Eth       |     N/A     |        | PoW-only           |
| eth_submitHashrate                                   | Eth       |     N/A     |        | PoW-only           |
| [eth_subscribe](04_websocket.md)                     | Websocket |      ✔      |        |                    |
| [eth_unsubscribe](04_websocket.md)                   | Websocket |      ✔      |        |                    |
| [personal_importRawKey](05_personal.md)              | Personal  |      ✔      |   ❌    |                    |
| [personal_listAccounts](05_personal.md)              | Personal  |      ✔      |   ❌    |                    | 
| [personal_lockAccount](05_personal.md)               | Personal  |      ✔      |   ❌    |                    |
| [personal_newAccount](05_personal.md)                | Personal  |      ✔      |   ❌    |                    |
| [personal_unlockAccount](05_personal.md)             | Personal  |      ✔      |   ❌    |                    |
| [personal_sendTransaction](05_personal.md)           | Personal  |      ✔      |   ❌    |                    |
| [personal_sign](05_personal.md)                      | Personal  |      ✔      |   ❌    |                    |
| [personal_ecRecover](05_personal.md)                 | Personal  |      ✔      |   ❌    |                    |
| [personal_initializeWallet](05_personal.md)          | Personal  |      ✔      |   ❌    |                    |
| [personal_unpair](05_personal.md)                    | Personal  |      ✔      |   ❌    |                    |
| db_putString                                         | DB        |             |        |                    |
| db_getString                                         | DB        |             |        |                    |
| db_putHex                                            | DB        |             |        |                    |
| db_getHex                                            | DB        |             |        |                    |
| shh_post                                             | SSH       |             |        |                    |
| shh_version                                          | SSH       |             |        |                    |
| shh_newIdentity                                      | SSH       |             |        |                    |
| shh_hasIdentity                                      | SSH       |             |        |                    |
| shh_newGroup                                         | SSH       |             |        |                    |
| shh_addToGroup                                       | SSH       |             |        |                    |
| shh_newFilter                                        | SSH       |             |        |                    |
| shh_uninstallFilter                                  | SSH       |             |        |                    |
| shh_getFilterChanges                                 | SSH       |             |        |                    |
| shh_getMessages                                      | SSH       |             |        |                    |
| admin_addPeer                                        | Admin     |             |        |                    |
| admin_datadir                                        | Admin     |             |        |                    | 
| admin_nodeInfo                                       | Admin     |             |        |                    |
| admin_peers                                          | Admin     |             |        |                    |
| admin_startRPC                                       | Admin     |             |        |                    |
| admin_startWS                                        | Admin     |             |        |                    |
| admin_stopRPC                                        | Admin     |             |        |                    |
| admin_stopWS                                         | Admin     |             |        |                    |
| clique_getSnapshot                                   | Clique    |             |        |                    |
| clique_getSnapshotAtHash                             | Clique    |             |        |                    |
| clique_getSigners                                    | Clique    |             |        |                    |
| clique_proposals                                     | Clique    |             |        |                    |
| clique_propose                                       | Clique    |             |        |                    |
| clique_discard                                       | Clique    |             |        |                    |
| clique_status                                        | Clique    |             |        |                    |
| debug_backtraceAt                                    | Debug     |             |        |                    |
| debug_blockProfile                                   | Debug     |      ✔      |        |                    |
| debug_cpuProfile                                     | Debug     |      ✔      |        |                    |
| debug_dumpBlock                                      | Debug     |             |        |                    |
| debug_gcStats                                        | Debug     |      ✔      |        |                    |
| debug_getBlockRlp                                    | Debug     |             |        |                    |
| debug_goTrace                                        | Debug     |      ✔      |        |                    |
| debug_freeOSMemory                                   | Debug     |      ✔      |        |                    |
| debug_memStats                                       | Debug     |      ✔      |        |                    |
| debug_mutexProfile                                   | Debug     |      ✔      |        |                    |
| debug_seedHash                                       | Debug     |             |        |                    |
| debug_setHead                                        | Debug     |             |        |                    |
| debug_setBlockProfileRate                            | Debug     |      ✔      |        |                    |
| debug_setGCPercent                                   | Debug     |      ✔      |        |                    |
| debug_setMutexProfileFraction                        | Debug     |      ✔      |        |                    |
| debug_stacks                                         | Debug     |      ✔      |        |                    |
| debug_startCPUProfile                                | Debug     |      ✔      |        |                    |
| debug_startGoTrace                                   | Debug     |      ✔      |        |                    |
| debug_stopCPUProfile                                 | Debug     |      ✔      |        |                    |
| debug_stopGoTrace                                    | Debug     |      ✔      |        |                    |
| debug_traceBlock                                     | Debug     |      ✔      |        |                    |
| [debug_traceBlockByNumber](06_debug.md)              | Debug     |      ✔      |        |                    |
| debug_traceBlockByHash                               | Debug     |      ✔      |        |                    |
| debug_traceBlockFromFile                             | Debug     |             |        |                    |
| debug_standardTraceBlockToFile                       | Debug     |             |        |                    |
| debug_standardTraceBadBlockToFile                    | Debug     |             |        |                    |
| [debug_traceTransaction](06_debug.md)                | Debug     |      ✔      |        |                    |
| debug_verbosity                                      | Debug     |             |        |                    |
| debug_vmodule                                        | Debug     |             |        |                    |
| debug_writeBlockProfile                              | Debug     |      ✔      |        |                    |
| debug_writeMemProfile                                | Debug     |      ✔      |        |                    |
| debug_writeMutexProfile                              | Debug     |      ✔      |        |                    |
| les_serverInfo                                       | Les       |             |        |                    |
| les_clientInfo                                       | Les       |             |        |                    |
| les_priorityClientInfo                               | Les       |             |        |                    |
| les_addBalance                                       | Les       |             |        |                    |
| les_setClientParams                                  | Les       |             |        |                    |
| les_setDefaultParams                                 | Les       |             |        |                    |
| les_latestCheckpoint                                 | Les       |             |        |                    |
| les_getCheckpoint                                    | Les       |             |        |                    |
| les_getCheckpointContractAddress                     | Les       |             |        |                    |
| [miner_getHashrate](07_miner.md)                     | Miner     |      ✔      |   ❌    | No-op              |
| [miner_setExtra](07_miner.md)                        | Miner     |      ✔      |   ❌    | No-op              |         
| [miner_setGasPrice](07_miner.md)                     | Miner     |      ✔      |   ❌    | Needs node restart | 
| [miner_start](07_miner.md)                           | Miner     |      ✔      |   ❌    | No-op              |
| [miner_stop](07_miner.md)                            | Miner     |      ✔      |   ❌    | No-op              |
| [miner_setGasLimit](07_miner.md)                     | Miner     |      ✔      |   ❌    | No-op              |
| [miner_setEtherbase](07_miner.md)                    | Miner     |      ✔      |   ❌    |                    |
| [txpool_content](08_txpool.md)                       | TxPool    |      ✔      |        |                    |
| [txpool_inspect](08_txpool.md)                       | TxPool    |      ✔      |        |                    |
| [txpool_status](08_txpool.md)                        | TxPool    |      ✔      |        |                    |