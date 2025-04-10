package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cometbft/cometbft/blocksync"

	"github.com/stratosnet/stratos-chain/misc"
	"github.com/stratosnet/stratos-chain/rpc/backend"
)

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type jsonrpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func errorJsonrpcResponse(w http.ResponseWriter, r *http.Request, msg string) {
	var request jsonrpcMessage
	body, err := io.ReadAll(r.Body)
	if err == nil {
		_ = json.Unmarshal(body, &request) // Ignore error, we'll just omit ID if parsing fails
	}

	response := jsonrpcMessage{
		Version: "2.0",
		ID:      request.ID, // Use the ID from the parsed request
		Error: &jsonError{
			Code:    -32601,
			Message: msg,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func NewStateSyncHandler(next http.Handler, b backend.TMBackend) http.Handler {
	if b == nil {
		panic("tendermint backend for web3 rpc not set")
	}

	if !b.GetTendermintConfig().StateSync.Enable {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reactor := b.GetStateSyncReactor()
		if reactor == nil {
			next.ServeHTTP(w, r)
			return
		}

		syncer := misc.GetMutableField(reactor, "syncer")
		if !misc.IsNilish(syncer) {
			errorJsonrpcResponse(w, r, "RPC not available: node is in statesync mode")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewBlockSyncHandler(next http.Handler, b backend.TMBackend) http.Handler {
	if b == nil {
		panic("tendermint backend for web3 rpc not set")
	}

	if !b.GetTendermintConfig().BlockSyncMode {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reactor := b.GetBlockSyncReactor()
		if reactor == nil {
			next.ServeHTTP(w, r)
			return
		}

		blockPool, ok := misc.GetMutableField(reactor, "pool").(*blocksync.BlockPool)
		if ok && blockPool.IsRunning() {
			errorJsonrpcResponse(w, r, fmt.Sprintf(
				"RPC not available: node is in blocksync mode (node height: %d < chain height: %d)",
				blockPool.Height(), blockPool.MaxPeerHeight()))
			return
		}
		next.ServeHTTP(w, r)
	})
}
