# Miner

---

## miner_getHashrate
**Private**: Requires authentication.

Get the hashrate in H/s (Hash operations per second).

**Proof-of-Work specific. This endpoint always returns `0`.**

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"miner_setGasPrice","params":[],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":0
}
~~~

---

## miner_setExtra
**Private**: Requires authentication.

Sets the extra data a validator can include when proposing blocks. This is capped at 32 bytes.

**Unsupported. This endpoint always returns an error**

**Parameters**

Data

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"miner_setExtra","params":["data"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":false
}
~~~

---

## miner_setGasPrice
**Private**: Requires authentication.

Sets the minimal gas price used to accept transactions. Any transaction below this limit is excluded from the validator block proposal process.

This method requires a `node` restart after being called because it changes the configuration file.

Make sure your `stchaind` start call is not using the flag `minimum-gas-prices` because this value will be used instead of the one set on the configuration file.

**Parameters**

Hex Gas Price

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"miner_setGasPrice","params":["0x0"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":true
}
~~~

---

## miner_start
**Private**: Requires authentication.

Start the CPU validation process with the given number of threads.

**Unsupported. This endpoint always returns an error**

**Parameters**

Hex Number of threads

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"miner_start","params":["0x1"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":false
}
~~~

---

## miner_stop
**Private**: Requires authentication.

Stop the validation operation.

**Unsupported. This endpoint always performs a no-op.**

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"miner_stop","params":[],"id":1}'
~~~

---

## miner_setGasLimit
**Private**: Requires authentication.

Sets the gas limit the miner will target when mining. Note: on networks where EIP-1559 (opens new window)is activated, this should be set to twice what you want the gas target (i.e. the effective gas used on average per block) to be.

**Unsupported. This endpoint always returns `false`**

**Parameters**

Hex gas limit

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"miner_setGasLimit","params":["0x10000"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":false
}
~~~

---

## miner_setEtherbase
**Private**: Requires authentication.

Sets the etherbase. It changes the wallet where the validator rewards will be deposited.

**Parameters**

Account Address

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"miner_setEtherbase","params":["0x3b7252d007059ffc82d16d022da3cbf9992d2f70"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":true
}
~~~

---
