# Personal

---

## personal_importRawKey
**Private**: Requires authentication.

Imports the given unencrypted private key (hex encoded string) into the key store, encrypting it with the passphrase.

Returns the address of the new account.

**Parameters**

`privkey`: `string`   
`password`: `string`

**Example**
~~~
// Request
curl -X POST --data 
    '{"jsonrpc":"2.0","method":"personal_importRawKey","params":
    ["c5bd76cd0cd948de17a31261567d219576e992d9066fe1a6bca97496dec634e2c8e06f8949773b300b9f73fabbbc7710d5d6691e96bcf3c9145e15daf6fe07b9", 
    "the key is this"],"id":1}' 
~~~

---

## personal_listAccounts
**Private**: Requires authentication.

Returns a list of addresses for accounts this node manages.

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"personal_listAccounts","params":[],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":[
        "0x3b7252d007059ffc82d16d022da3cbf9992d2f70",
        "0xddd64b4712f7c8f1ace3c145c950339eddaf221d",
        "0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0"
    ]
}
~~~

---

## personal_lockAccount
**Private**: Requires authentication.

Removes the private key with given address from memory. The account can no longer be used to send transactions.

**Parameters**

Account Address

**Example**
~~~
Copy
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"personal_lockAccount","params":["0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":true
}
~~~

---

## personal_newAccount
**Private**: Requires authentication.

Generates a new private key and stores it in the key store directory. The key file is encrypted with the given passphrase. Returns the address of the new account.

**Parameters**

Passphrase

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"personal_newAccount","params":["This is the passphrase"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":"0xf0e4086ad1c6aab5d42161d5baaae2f9ad0571c0"
}
~~~

---

## personal_unlockAccount
**Private**: Requires authentication.

Decrypts the key with the given address from the key store.

Both passphrase and unlock duration are optional when using the JavaScript console. The unencrypted key will be held in memory until the unlock duration expires. If the unlock duration defaults to 300 seconds. An explicit duration of zero seconds unlocks the key until geth exits.

The account can be used with eth_sign and eth_sendTransaction while it is unlocked.

**Parameters**

* Account Address
* Passphrase
* Duration

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"personal_unlockAccount","params":["0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0", "secret passphrase", 30],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":true
}
~~~

---

## personal_sendTransaction
**Private**: Requires authentication.

Validate the given passphrase and submit transaction.

The transaction is the same argument as for `eth_sendTransaction` and contains the `from` address. If the passphrase can be used to decrypt the private key belonging to `tx.from` the transaction is verified, signed and send onto the network.

>The account is not unlocked globally in the node and cannot be used in other RPC calls.

**Parameters**

* Object containing:
  * `from`: `DATA`, 20 Bytes - The address the transaction is sent from.
  * `to`: `DATA`, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
  * `value`: `QUANTITY` - value sent with this transaction
* Passphrase

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"personal_sendTransaction","params":[{"from":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70","to":"0xddd64b4712f7c8f1ace3c145c950339eddaf221d", "value":"0x16345785d8a0000"}, "passphrase"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":"0xd2a31ec1b89615c8d1f4d08fe4e4182efa4a9c0d5758ace6676f485ea60e154c"
}
~~~

---

## personal_sign
**Private**: Requires authentication.

The sign method calculates an Ethereum specific signature with: `sign(keccack256("\x19Ethereum Signed Message:\n" + len(message) + message)))`,

**Parameters**

* Message
* Account Address
* Password

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"personal_sign","params":["0xdeadbeaf", "0x3b7252d007059ffc82d16d022da3cbf9992d2f70", "password"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":"0xf9ff74c86aefeb5f6019d77280bbb44fb695b4d45cfe97e6eed7acd62905f4a85034d5c68ed25a2e7a8eeb9baf1b8401e4f865d92ec48c1763bf649e354d900b1c"
}
~~~

---

## personal_ecRecover
**Private**: Requires authentication.

`ecRecover` returns the address associated with the private key that was used to calculate the signature in `personal_sign`.

**Parameters**

* Message
* Signature returned from `personal_sign`

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"personal_ecRecover","params":["0xdeadbeaf", "0xf9ff74c86aefeb5f6019d77280bbb44fb695b4d45cfe97e6eed7acd62905f4a85034d5c68ed25a2e7a8eeb9baf1b8401e4f865d92ec48c1763bf649e354d900b1c"],"id":1}'
// Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70"
}
~~~

---

## personal_initializeWallet
**Private**: Requires authentication.

Initializes a new wallet at the provided URL, by generating and returning a new private key.

**Parameters**

`url`: `string`

Example
~~~
// Request
curl -X POST --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_initializeWallet", "params": [<url>]}'
~~~

---

## personal_unpair
**Private**: Requires authentication.

Unpair deletes a pairing between wallet and the node.

**Parameters**

* URL
* Pairing password

Example
~~~
// Request
curl -X POST --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_unpair", "params": [<url>, <pin>]}'
~~~

---


