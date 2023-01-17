# Eth

---

## eth_protocolVersion

Returns the current Ethereum protocol version.

**Parameters**

None

**Returns**

`String` - The current Ethereum protocol version

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_protocolVersion","params":[],"id":67}'
// Result
{
    "id":67,
    "jsonrpc": "2.0",
    "result": "54"
}
~~~

---

## eth_syncing
Returns an object with data about the sync status or `false`.

**Parameters**

None

**Returns**

`Object|Boolean`, An object with sync status data or `FALSE`, when not syncing:

* `startingBlock`: `QUANTITY` - The block at which the import started (will only be reset, after the sync reached his head)  
* `currentBlock`: `QUANTITY` - The current block, same as eth_blockNumber  
* `highestBlock`: `QUANTITY` - The estimated highest block  

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": {
        startingBlock: '0x384',
        currentBlock: '0x386',
        highestBlock: '0x454'
    }
}
// Or when not syncing
{
    "id":1,
    "jsonrpc": "2.0",
    "result": false
}
~~~

---

## eth_coinbase
Returns the client coinbase address.

**Parameters**

None

**Returns**

`DATA`, 20 bytes - the current coinbase address.

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_coinbase","params":[],"id":64}'
// Result
{
    "id":64,
    "jsonrpc": "2.0",
    "result": "0x407d73d8a49eeb85d32cf465507dd71d507100c1"
}
~~~

---

## eth_gasPrice
Returns the current price per gas in wei.

**Parameters**

None

**Returns**

`QUANTITY` - integer of the current gas price in wei.

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":73}'
// Result
{
    "id":73,
    "jsonrpc": "2.0",
    "result": "0x1dfd14000" // 8049999872 Wei
}
~~~

---

## eth_accounts
Returns a list of addresses owned by client.

**Parameters**

None

**Returns**

`Array of DATA`, 20 Bytes - addresses owned by the client.

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": ["0x407d73d8a49eeb85d32cf465507dd71d507100c1"]
}
~~~

---

## eth_blockNumber
Returns the number of most recent block.

**Parameters**

None

**Returns**

`QUANTITY` - integer of the current block number the client is on.

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}'
// Result
{
    "id":83,
    "jsonrpc": "2.0",
    "result": "0x4b7" // 1207
}
~~~

---

## eth_getBalance
Returns the balance of the account of given address.

**Parameters**

`DATA`, 20 Bytes - address to check for balance.   
`QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter)   
~~~
params: ["0x407d73d8a49eeb85d32cf465507dd71d507100c1", "latest"]
~~~

**Returns**

`QUANTITY` - integer of the current balance in wei.

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1", "latest"],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0x0234c8a3397aab58" // 158972490234375000
}
~~~

---

## eth_getStorageAt
Returns the value from a storage position at a given address.

**Parameters**

`DATA`, 20 Bytes - address of the storage.   
`QUANTITY` - integer of the position in the storage.   
`QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter)

**Returns**

`DATA` - the value at this storage position.

**Example**  

Calculating the correct position depends on the storage to retrieve. Consider the following contract deployed at `0x295a70b2de5e3953354a6a8344e616ed314d7251` by address `0x391694e7e0b0cce554cb130d723a9d27458f9298`.
~~~
contract Storage {
    uint pos0;
    mapping(address => uint) pos1;
    function Storage() {
        pos0 = 1234;
        pos1[msg.sender] = 5678;
    }
}
~~~

Retrieving the value of pos0 is straight forward:
~~~
curl -X POST --data '{"jsonrpc":"2.0", "method": "eth_getStorageAt", "params": ["0x295a70b2de5e3953354a6a8344e616ed314d7251", "0x0", "latest"], "id": 1}' localhost:8545
{"jsonrpc":"2.0","id":1,"result":"0x00000000000000000000000000000000000000000000000000000000000004d2"}
~~~

Retrieving an element of the map is harder. The position of an element in the map is calculated with:
~~~
keccack(LeftPad32(key, 0), LeftPad32(map position, 0))
~~~

This means to retrieve the storage on pos1["0x391694e7e0b0cce554cb130d723a9d27458f9298"] we need to calculate the position with:
~~~
keccak(
    decodeHex(
        "000000000000000000000000391694e7e0b0cce554cb130d723a9d27458f9298" +
        "0000000000000000000000000000000000000000000000000000000000000001"
    )
)
~~~

The geth console which comes with the web3 library can be used to make the calculation:
~~~
> var key = "000000000000000000000000391694e7e0b0cce554cb130d723a9d27458f9298" + "0000000000000000000000000000000000000000000000000000000000000001"
undefined
> web3.sha3(key, {"encoding": "hex"})
"0x6661e9d6d8b923d5bbaab1b96e1dd51ff6ea2a93520fdc9eb75d059238b8c5e9"
~~~

Now to fetch the storage:
~~~
curl -X POST --data '{"jsonrpc":"2.0", "method": "eth_getStorageAt", "params": ["0x295a70b2de5e3953354a6a8344e616ed314d7251", "0x6661e9d6d8b923d5bbaab1b96e1dd51ff6ea2a93520fdc9eb75d059238b8c5e9", "latest"], "id": 1}' localhost:8545
{"jsonrpc":"2.0","id":1,"result":"0x000000000000000000000000000000000000000000000000000000000000162e"}
~~~

---

### eth_getTransactionCount
Returns the number of transactions sent from an address.

**Parameters**

`DATA`, 20 Bytes - address.   
`QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter)
~~~
params: [
"0x407d73d8a49eeb85d32cf465507dd71d507100c1",
"latest", // state at the latest block
]
~~~

**Returns**

`QUANTITY` - integer of the number of transactions send from this address.

Example
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1","latest"],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0x1" // 1
}
~~~

---

## eth_getBlockTransactionCountByHash
Returns the number of transactions in a block from a block matching the given block hash.

**Parameters**

`DATA`, 32 Bytes - hash of a block
~~~
params: ["0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"]
~~~

**Returns**

`QUANTITY` - integer of the number of transactions in this block.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByHash","params":["0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0xb" // 11
}
~~~

---

## eth_getBlockTransactionCountByNumber
Returns the number of transactions in a block matching the given block number.

**Parameters**

`QUANTITY|TAG` - integer of a block number, or the string `"earliest"`, `"latest"` or `"pending"`, as in the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter).
~~~
params: [
    "0xe8", // 232
]
~~~

**Returns**

`QUANTITY` - integer of the number of transactions in this block.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByNumber","params":["0xe8"],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0xa" // 10
}
~~~

---

## eth_getCode
Returns code at a given address.

**Parameters**

`DATA`, 20 Bytes - address   
`QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter)
~~~
params: [
    "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
    "0x2", // 2
]
~~~

**Returns**

`DATA` - the code from the given address.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getCode","params":["0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b", "0x2"],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0x600160008035811a818181146012578301005b601b6001356025565b8060005260206000f25b600060078202905091905056"
}
~~~

---

## eth_getProof
Returns the account and storage values of the specified account including the Merkle-proof. This call can be used to verify that the data you are pulling from is not tampered with.

**Parameters**

* `DATA`, 20 Bytes - address of the account.   
* `ARRAY`, 32 Bytes - array of storage-keys which should be proofed and included. See `eth_getStorageAt`
* `QUANTITY|TAG` - integer block number, or the string `"latest"` or `"earliest"`, see the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter).
~~~
params: [
    "0x7F0d15C7FAae65896648C8273B6d7E43f58Fa842",
    [
        "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
    ],
    "latest"
]
~~~

**Returns**

`Object` - A account object:
* `balance`: `QUANTITY` - the balance of the account. See `eth_getBalance`
* `codeHash`: `DATA`, 32 Bytes - hash of the code of the account. For a simple Account without code it will return `"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"`
* `nonce`: `QUANTITY`, - nonce of the account. See `eth_getTransactionCount`
* `storageHash`: `DATA`, 32 Bytes - SHA3 of the StorageRoot. All storage will deliver a MerkleProof starting with this rootHash.
* `accountProof`: `ARRAY` - Array of rlp-serialized MerkleTree-Nodes, starting with the stateRoot-Node, following the path of the SHA3 (address) as key.
* `storageProof`: `ARRAY` - Array of storage-entries as requested. Each entry is a object with these properties:
  * `key`: `QUANTITY` - the requested storage key
  * `value`: `QUANTITY` - the storage value
  * `proof`: `ARRAY` - Array of rlp-serialized MerkleTree-Nodes, starting with the storageHash-Node, following the path of the SHA3 (key) as path.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getProof","params":["0x7F0d15C7FAae65896648C8273B6d7E43f58Fa842",["0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"],"latest"],"id":1}'
// Result
{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "accountProof": [
            "0xf90211a...0701bc80",
            "0xf90211a...0d832380",
            "0xf90211a...5fb20c80",
            "0xf90211a...0675b80",
            "0xf90151a0...ca08080"
        ],
        "balance": "0x0",
        "codeHash": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
        "nonce": "0x0",
        "storageHash": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
        "storageProof": [
            {
                "key": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
                "proof": [
                    "0xf90211a...0701bc80",
                    "0xf90211a...0d832380"
                ],
                "value": "0x1"
            }
        ]
    }
}
~~~

---

## eth_sign
The sign method calculates an Ethereum specific signature with: `sign(keccak256("\x19Ethereum Signed Message:\n" + len(message) + message)))`.

By adding a prefix to the message makes the calculated signature recognizable as an Ethereum specific signature. This prevents misuse where a malicious dapp can sign arbitrary data (e.g. transaction) and use the signature to impersonate the victim.

Note: the address to sign with must be unlocked.

**Parameters**

`DATA`, 20 Bytes - address
`DATA`, N Bytes - message to sign

**Returns**

`DATA`: Signature

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sign","params":["0x9b2055d370f73ec7d8a03e965129118dc8f5bf83", "0xdeadbeaf"],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0xa3f20717a250c2b0b729b7e5becbff67fdaef7e0699da4de7ca5895b02a170a12d887fd3b17bfdce3481f10bea41f45ba9f709d39ce8325427b57afcfc994cee1b"
}
~~~

---

## eth_sendTransaction
Creates new message call transaction or a contract creation, if the data field contains code.

**Parameters**

1. `Object` - The transaction object   

* `from`: `DATA`, 20 Bytes - The address the transaction is sent from.
* `to`: `DATA`, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
* `gas`: `QUANTITY` - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.
* `gasPrice`: `QUANTITY` - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas.
* `value`: `QUANTITY` - (optional) Integer of the value sent with this transaction.
* `data`: `DATA` - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters.
* `nonce`: `QUANTITY` - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.
~~~
params: [
    {
        from: "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
        to: "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
        gas: "0x76c0", // 30400
        gasPrice: "0x9184e72a000", // 10000000000000
        value: "0x9184e72a", // 2441406250
        data: "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675",
    },
]
~~~

**Returns**

`DATA`, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.

Use `eth_getTransactionReceipt` to get the contract address, after the transaction was mined, when you created a contract.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{see above}],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331"
}
~~~

---

## eth_sendRawTransaction
Creates new message call transaction or a contract creation for signed transactions.

**Parameters**

`DATA`, The signed transaction data.
~~~
params: [
    "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675",
]
~~~

**Returns**

`DATA`, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.

Use `eth_getTransactionReceipt` to get the contract address, after the transaction was mined, when you created a contract.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":[{see above}],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331"
}
~~~

---

## eth_call
Executes a new message call immediately without creating a transaction on the block chain.

**Parameters**

1. `Object` - The transaction call object

* `from`: `DATA`, 20 Bytes - (optional) The address the transaction is sent from.
* `to`: `DATA`, 20 Bytes - The address the transaction is directed to.
* `gas`: `QUANTITY` - (optional) Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.
* `gasPrice`: `QUANTITY` - (optional) Integer of the gasPrice used for each paid gas
* `value`: `QUANTITY` - (optional) Integer of the value sent with this transaction
* `data`: `DATA` - (optional) Hash of the method signature and encoded parameters. For details see [Ethereum Contract ABI in the Solidity documentation](https://docs.soliditylang.org/en/latest/abi-spec.html)

2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter)

**Returns**

`DATA` - the return value of executed contract.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_call","params":[{see above}],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0x"
}
~~~

---

## eth_estimateGas
Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain. Note that the estimate may be significantly more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.

**Parameters**

See `eth_call` parameters, expect that all properties are optional. If no gas limit is specified geth uses the block gas limit from the pending block as an upper bound. As a result the returned estimate might not be enough to executed the call/transaction when the amount of gas is higher than the pending block gas limit.

**Returns**

`QUANTITY` - the amount of gas used.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_estimateGas","params":[{see above}],"id":1}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0x5208" // 21000
}
~~~

---

## eth_getBlockByHash
Returns information about a block by hash.

**Parameters**

`DATA`, 32 Bytes - Hash of a block.   
`Boolean` - If true it returns the full transaction objects, if false only the hashes of the transactions.
~~~
params: [
    "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae",
    false,
]
~~~

**Returns**

`Object` - A block object, or `null` when no block was found:

* `number`: `QUANTITY` - the block number. `null` when its pending block.
* `hash`: `DATA`, 32 Bytes - hash of the block. `null` when its pending block.
* `parentHash`: `DATA`, 32 Bytes - hash of the parent block.
* `nonce`: `DATA`, 8 Bytes - hash of the generated proof-of-work. `null` when its pending block.
* `sha3Uncles`: `DATA`, 32 Bytes - SHA3 of the uncles data in the block.
* `logsBloom`: `DATA`, 256 Bytes - the bloom filter for the logs of the block. `null` when its pending block.
* `transactionsRoot`: `DATA`, 32 Bytes - the root of the transaction trie of the block.
* `stateRoot`: `DATA`, 32 Bytes - the root of the final state trie of the block.
* `receiptsRoot`: `DATA`, 32 Bytes - the root of the receipts trie of the block.
* `miner`: `DATA`, 20 Bytes - the address of the beneficiary to whom the mining rewards were given.
* `difficulty`: `QUANTITY` - integer of the difficulty for this block.
* `totalDifficulty`: `QUANTITY` - integer of the total difficulty of the chain until this block.
* `extraData`: `DATA` - the "extra data" field of this block.
* `size`: `QUANTITY` - integer the size of this block in bytes.
* `gasLimit`: `QUANTITY` - the maximum gas allowed in this block.
* `gasUsed`: `QUANTITY` - the total used gas by all transactions in this block.
* `timestamp`: `QUANTITY` - the unix timestamp for when the block was collated.
* `transactions`: `Array` - Array of transaction objects, or 32 Bytes transaction hashes depending on the last given parameter.
* `uncles`: `Array` - Array of uncle hashes.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae", false],"id":1}'
// Result
{
    {
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "difficulty": "0x4ea3f27bc",
        "extraData": "0x476574682f4c5649562f76312e302e302f6c696e75782f676f312e342e32",
        "gasLimit": "0x1388",
        "gasUsed": "0x0",
        "hash": "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae",
        "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "miner": "0xbb7b8287f3f0a933474a79eae42cbca977791171",
        "mixHash": "0x4fffe9ae21f1c9e15207b1f472d5bbdd68c9595d461666602f2be20daf5e7843",
        "nonce": "0x689056015818adbe",
        "number": "0x1b4",
        "parentHash": "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
        "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0x220",
        "stateRoot": "0xddc8b0234c2e0cad087c8b389aa7ef01f7d79b2570bccb77ce48648aa61c904d",
        "timestamp": "0x55ba467c",
        "totalDifficulty": "0x78ed983323d",
        "transactions": [
        ],
        "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
        "uncles": [
        ]
    }
}
~~~

---

## eth_getBlockByNumber
Returns information about a block by block number.

**Parameters**

`QUANTITY|TAG` - integer of a block number, or the string `"earliest"`, `"latest"` or `"pending"`, as in the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter).   
`Boolean` - If `true` it returns the full transaction objects, if `false` only the hashes of the transactions.
~~~
params: [
    "0x1b4", // 436
    true,
]
~~~

**Returns**

See `eth_getBlockByHash`

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}'
~~~

Result see `eth_getBlockByHash`

---

## eth_getTransactionByHash
Returns the information about a transaction requested by transaction hash.

**Parameters**

`DATA`, 32 Bytes - hash of a transaction
~~~
params: ["0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"]
~~~

**Returns**

`Object` - A transaction object, or `null` when no transaction was found:

* `blockHash`: `DATA`, 32 Bytes - hash of the block where this transaction was in. `null` when it's pending.
* `blockNumber`: `QUANTITY` - block number where this transaction was in. `null` when it's pending.
* `from`: `DATA`, 20 Bytes - address of the sender.
* `gas`: `QUANTITY` - gas provided by the sender.
* `gasPrice`: `QUANTITY` - gas price provided by the sender in Wei.
* `hash`: `DATA`, 32 Bytes - hash of the transaction.
* `input`: `DATA` - the data send along with the transaction.
* `nonce`: `QUANTITY` - the number of transactions made by the sender prior to this one.
* `to`: `DATA`, 20 Bytes - address of the receiver. `null` when it's a contract creation transaction.
* `transactionIndex`: `QUANTITY` - integer of the transactions index position in the block. `null` when it's pending.
* `value`: `QUANTITY` - value transferred in Wei.
* `v`: `QUANTITY` - ECDSA recovery id
* `r`: `QUANTITY` - ECDSA signature r
* `s`: `QUANTITY` - ECDSA signature s

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":{
        "blockHash":"0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
        "blockNumber":"0x5daf3b", // 6139707
        "from":"0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
        "gas":"0xc350", // 50000
        "gasPrice":"0x4a817c800", // 20000000000
        "hash":"0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b",
        "input":"0x68656c6c6f21",
        "nonce":"0x15", // 21
        "to":"0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
        "transactionIndex":"0x41", // 65
        "value":"0xf3dbb76162000", // 4290000000000000
        "v":"0x25", // 37
        "r":"0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea",
        "s":"0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c"
    }
}
~~~

---

## eth_getTransactionByBlockHashAndIndex
Returns information about a transaction by block hash and transaction index position.

**Parameters**

`DATA`, 32 Bytes - hash of a block.   
`QUANTITY` - integer of the transaction index position.
~~~
params: [
    "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
    "0x0", // 0
]
~~~

**Returns**

See `eth_getTransactionByHash`

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByBlockHashAndIndex","params":["0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b", "0x0"],"id":1}'
~~~

Result see `eth_getTransactionByHash`

---

## eth_getTransactionByBlockNumberAndIndex
Returns information about a transaction by block number and transaction index position.

**Parameters**

`QUANTITY|TAG` - a block number, or the string `"earliest"`, `"latest"` or `"pending"`, as in the [default block parameter](https://ethereum.org/en/developers/docs/apis/json-rpc/#default-block-parameter).   
`QUANTITY` - the transaction index position.
~~~
params: [
    "0x29c", // 668
    "0x0", // 0
]
~~~

**Returns**

See `eth_getTransactionByHash`

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params":["0x29c", "0x0"],"id":1}'
~~~

Result see `eth_getTransactionByHash`

---

## eth_getTransactionReceipt
Returns the receipt of a transaction by transaction hash.

**Note** That the receipt is not available for pending transactions.

**Parameters**

`DATA`, 32 Bytes - hash of a transaction
~~~
params: ["0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"]
~~~

**Returns** 

`Object` - A transaction receipt object, or `null` when no receipt was found:

* `transactionHash` : `DATA`, 32 Bytes - hash of the transaction.
* `transactionIndex`: `QUANTITY` - integer of the transactions index position in the block.
* `blockHash`: `DATA`, 32 Bytes - hash of the block where this transaction was in.
* `blockNumber`: `QUANTITY` - block number where this transaction was in.
* `from`: `DATA`, 20 Bytes - address of the sender.
* `to`: `DATA`, 20 Bytes - address of the receiver. `null` when it's a contract creation transaction.
* `cumulativeGasUsed` : `QUANTITY` - The total amount of gas used when this transaction was executed in the block.
* `gasUsed` : `QUANTITY` - The amount of gas used by this specific transaction alone.
* `contractAddress` : `DATA`, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise `null`.
* `logs`: `Array` - Array of log objects, which this transaction generated.
* `logsBloom`: `DATA`, 256 Bytes - Bloom filter for light clients to quickly retrieve related logs. It also returns either :
* `root` : `DATA` 32 bytes of post-transaction stateroot (pre Byzantium)
* `status`: `QUANTITY` either `1` (success) or `0` (failure)

**Example**
~~~
// Request
curl -X POST --data 
'{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":
["0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"],"id":1}'
// Result
{
    "id":1,
    "jsonrpc":"2.0",
    "result": {
        transactionHash: '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238',
        transactionIndex:  '0x1', // 1
        blockNumber: '0xb', // 11
        blockHash: '0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b',
        cumulativeGasUsed: '0x33bc', // 13244
        gasUsed: '0x4dc', // 1244
        contractAddress: '0xb60e8dd61c5d32be8058bb8eb970870f07233155', // or null, if none was created
        logs: [{
            // logs as returned by getFilterLogs, etc.
        }, ...],
        logsBloom: "0x00...0", // 256 byte bloom filter
        status: '0x1'
    }
}
~~~

---

## eth_newFilter
Creates a filter object, based on filter options, to notify when the state changes (logs). To check if the state has changed, call `eth_getFilterChanges`.

**A note on specifying topic filters**: Topics are order-dependent. A transaction with a log with topics [A, B] will be matched by the following topic filters:

* `[]` "anything"
* `[A]` "A in first position (and anything after)"
* `[null, B]` "anything in first position AND B in second position (and anything after)"
* `[A, B]` "A in first position AND B in second position (and anything after)"
* `[[A, B], [A, B]]` "(A OR B) in first position AND (A OR B) in second position (and anything after)"

**Parameters**

`Object` - The filter options:

* `fromBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
* `toBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
* `address`: `DATA|Array`, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.
* `topics`: `Array of DATA`, - (optional) Array of 32 Bytes `DATA` topics. Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
~~~
params: [
    {
        fromBlock: "0x1",
        toBlock: "0x2",
        address: "0x8888f1f195afa192cfee860698584c030f4c9db1",
        topics: [
            "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
            null,
            [
                "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
                "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc",
            ],
        ],
    },
]
~~~

**Returns**

`QUANTITY` - A filter id.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_newFilter","params":[{"topics":["0x12341234"]}],"id":73}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": "0x1" // 1
}
~~~

---

## eth_newBlockFilter
Creates a filter in the node, to notify when a new block arrives. To check if the state has changed, call eth_getFilterChanges.

**Parameters**

None

**Returns**

`QUANTITY` - A filter id.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_newBlockFilter","params":[],"id":73}'
// Result
{
    "id":1,
    "jsonrpc":  "2.0",
    "result": "0x1" // 1
}
~~~

---

## eth_newPendingTransactionFilter
Creates a filter in the node, to notify when new pending transactions arrive. To check if the state has changed, call `eth_getFilterChanges`.

**Parameters**

None

**Returns**

`QUANTITY` - A filter id.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_newPendingTransactionFilter","params":[],"id":73}'
// Result
{
    "id":1,
    "jsonrpc":  "2.0",
    "result": "0x1" // 1
}
~~~

---

## eth_uninstallFilter
Uninstalls a filter with given id. Should always be called when watch is no longer needed. Additionally, Filters timeout when they aren't requested with `eth_getFilterChanges` for a period of time.

**Parameters**

`QUANTITY` - The filter id.
~~~
params: [
    "0xb", // 11
]
~~~

**Returns**

`Boolean` - `true` if the filter was successfully uninstalled, otherwise `false`.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_uninstallFilter","params":["0xb"],"id":73}'
// Result
{
    "id":1,
    "jsonrpc": "2.0",
    "result": true
}
~~~

---

## eth_getFilterChanges
Polling method for a filter, which returns an array of logs which occurred since last poll.

**Parameters**

`QUANTITY` - the filter id.
~~~
params: [
    "0x16", // 22
]
~~~

Returns

`Array` - Array of log objects, or an empty array if nothing has changed since last poll.

* For filters created with `eth_newBlockFilter` the return are block hashes (`DATA`, 32 Bytes), e.g. `["0x3454645634534..."]`.
* For filters created with `eth_newPendingTransactionFilter` the return are transaction hashes (`DATA`, 32 Bytes), e.g. `["0x6345343454645..."]`.
* For filters created with `eth_newFilter` logs are objects with following params:
   * `removed`: `TAG` - `true` when the log was removed, due to a chain reorganization. `false` if it's a valid log.
   * `logIndex`: `QUANTITY` - integer of the log index position in the block. `null` when its pending log.
   * `transactionIndex`: `QUANTITY` - integer of the transactions index position log was created from. `null` when its pending log.
   * `transactionHash`: `DATA`, 32 Bytes - hash of the transactions this log was created from. `null` when its pending log.
   * `blockHash`: `DATA`, 32 Bytes - hash of the block where this log was in. `null` when it's pending. `null` when its pending log.
   * `blockNumber`: `QUANTITY` - the block number where this log was in. `null` when it's pending. `null` when its pending log.
   * `address`: `DATA`, 20 Bytes - address from which this log originated.
   * `data`: `DATA` - contains one or more 32 Bytes non-indexed arguments of the log.
   * `topics`: `Array of DATA` - Array of 0 to 4 32 Bytes `DATA` of indexed log arguments. (In solidity: The first topic is the hash of the signature of the event (e.g. `Deposit(address,bytes32,uint256)`), except you declared the event with the `anonymous` specifier.)

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getFilterChanges","params":["0x16"],"id":73}'
// Result
{
    "id":1,
    "jsonrpc":"2.0",
    "result": [{
        "logIndex": "0x1", // 1
        "blockNumber":"0x1b4", // 436
        "blockHash": "0x8216c5785ac562ff41e2dcfdf5785ac562ff41e2dcfdf829c5a142f1fccd7d",
        "transactionHash":  "0xdf829c5a142f1fccd7d8216c5785ac562ff41e2dcfdf5785ac562ff41e2dcf",
        "transactionIndex": "0x0", // 0
        "address": "0x16c5785ac562ff41e2dcfdf829c5a142f1fccd7d",
        "data":"0x0000000000000000000000000000000000000000000000000000000000000000",
        "topics": ["0x59ebeb90bc63057b6515673c3ecf9438e5058bca0f92585014eced636878c9a5"]
    },{
        ...
    }]
}
~~~

---

## eth_getFilterLogs
Returns an array of all logs matching filter with given id.

**Parameters**

`QUANTITY` - The filter id.
~~~
params: [
    "0x16", // 22
]
~~~

**Returns**

See `eth_getFilterChanges`

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getFilterLogs","params":["0x16"],"id":74}'
~~~

Result see eth_getFilterChanges

---

## eth_getLogs
Returns an array of all logs matching a given filter object.

**Parameters**

`Object` - The filter options:

* `fromBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
* `toBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
* `address`: `DATA|Array`, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.
* `topics`: `Array of DATA`, - (optional) Array of 32 Bytes `DATA` topics. Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
* `blockhash`: `DATA`, 32 Bytes - (optional, future) With the addition of EIP-234, `blockHash` will be a new filter option which restricts the logs returned to the single block with the 32-byte hash `blockHash`. Using `blockHash` is equivalent to `fromBlock` = `toBlock` = the block number with hash `blockHash`. If `blockHash` is present in the filter criteria, then neither `fromBlock` nor `toBlock` are allowed.

~~~
params: [
    {
        topics: [
            "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
        ],
    },
]
~~~

**Returns**

See `eth_getFilterChanges`

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"topics":["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"]}],"id":74}'
~~~

Result see eth_getFilterChanges

---

