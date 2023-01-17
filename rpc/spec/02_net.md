# Net

---

## net_version
Returns the current network id.

**Parameters**

None

**Returns**

`String` - The current network id.

The full list of current network IDs is available at chainlist.org. Sopme common ones are: 1: Ethereum Mainnet 2: Morden testnet (now deprecated) 3: Ropsten testnet 4: Rinkeby testnet 5: Goerli testnet

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}'
// Result
{
"id":67,
"jsonrpc": "2.0",
"result": "3"
}
~~~

---

## net_listening
Returns true if client is actively listening for network connections.

**Parameters**

None

**Returns**

`Boolean` - true when listening, otherwise false.

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"net_listening","params":[],"id":67}'
// Result
{
"id":67,
"jsonrpc":"2.0",
"result":true
}
~~~

---

## net_peerCount
Returns number of peers currently connected to the client.

**Parameters**

None

**Returns**

`QUANTITY` - integer of the number of connected peers.

**Example**

~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":74}'
// Result
{
"id":74,
"jsonrpc": "2.0",
"result": "0x2" // 2
}
~~~

---