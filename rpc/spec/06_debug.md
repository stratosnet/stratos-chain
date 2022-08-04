# Debug

---

## debug_traceTransaction
The traceTransaction debugging method will attempt to run the transaction in the exact same manner as it was executed on the network. It will replay any transaction that may have been executed prior to this one before it will finally attempt to execute the transaction that corresponds to the given hash.

**Parameters**

Trace Config

**Example**
~~~
// Request
curl -X POST --data 
    '{"jsonrpc":"2.0","method":"debug_traceTransaction","params":
    ["0xddecdb13226339681372b44e01df0fbc0f446fca6f834b2de5ecb1e569022ec8", {"tracer": "{data: [], fault: function(log) {}, step: function(log) { if(log.op.toString() == \"CALL\") this.data.push(log.stack.peek(0)); }, result: function() { return this.data; }}"}],"id":1}'
//Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":[
        {
            "result":[
                "68410", 
                "51470"
            ]
        }
    ]
}
~~~

---

## debug_traceBlockByNumber
The traceBlockByNumber endpoint accepts a block number and will replay the block that is already present in the database.

**Parameters**

Trace Config

**Example**
~~~
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"debug_traceBlockByNumber","params":["0xe", {"tracer": "{data: [], fault: function(log) {}, step: function(log) { if(log.op.toString() == \"CALL\") this.data.push(log.stack.peek(0)); }, result: function() { return this.data; }}"}],"id":1}'
//Result
{
    "jsonrpc":"2.0",
    "id":1,
    "result":[
        {
            "result":[
                "68410", 
                "51470"
            ]
        }
    ]
}
~~~

---