package rpc

import (
	"context"
	"net/http"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stratosnet/stratos-chain/rpc/backend"
	"github.com/stratosnet/stratos-chain/server/config"
)

type Web3Server struct {
	modules []string
	httpURI string
	wsURI   string
	enabled bool
	backend backend.BackendI
	logger  log.Logger
}

func NewWeb3Server(cfg config.Config, backend backend.BackendI, logger log.Logger) *Web3Server {
	return &Web3Server{
		modules: cfg.JSONRPC.API,
		httpURI: cfg.JSONRPC.Address,
		wsURI:   cfg.JSONRPC.WsAddress,
		enabled: cfg.JSONRPC.Enable,
		backend: backend,
		logger:  logger,
	}
}

func (web3 *Web3Server) start(uri string, handler http.Handler) error {
	if !web3.enabled {
		web3.logger.Info("Web3 api disabled, skipping")
		return nil
	}

	var (
		err error
	)
	channel := make(chan error)
	timeout := make(chan error)

	srv := &http.Server{
		Addr:         uri,
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	srv.SetKeepAlivesEnabled(true)

	//Timeout Go routine
	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		timeout <- nil
	}()

	go func(ch chan error) {
		web3.logger.Info("starting Web3 RPC server on " + uri)
		err := srv.ListenAndServe()
		if err != nil {
			web3.logger.Error("server error, details: %s", err)
		}
		srv.Shutdown(context.TODO())
		ch <- err
	}(channel)

	select {
	case err = <-channel:
	case err = <-timeout:
	}

	return err
}

func (web3 *Web3Server) registerAPIs(server *rpc.Server, apis []rpc.API) error {
	// Generate the allow list based on the allowed modules
	allowList := make(map[string]bool)
	for _, module := range web3.modules {
		allowList[module] = true
	}

	for _, api := range apis {
		if allowList[api.Namespace] || len(allowList) == 0 {
			if err := server.RegisterName(api.Namespace, api.Service); err != nil {
				web3.logger.Error(
					"failed to register service in JSON RPC namespace",
					"namespace", api.Namespace,
					"service", api.Service,
				)
				return err
			}
		}
	}
	return nil
}

func (web3 *Web3Server) StartHTTP(apis []rpc.API) error {
	rpcSrv := rpc.NewServer()
	handler := node.NewHTTPHandlerStack(rpcSrv, []string{}, []string{"*"}, []byte{})
	handler = NewBlockSyncHandler(handler, web3.backend)
	handler = NewStateSyncHandler(handler, web3.backend)
	if err := web3.registerAPIs(rpcSrv, apis); err != nil {
		return err
	}
	return web3.start(web3.httpURI, handler)
}

func (web3 *Web3Server) StartWS(apis []rpc.API) error {
	rpcSrv := rpc.NewServer()
	handler := rpcSrv.WebsocketHandler([]string{})
	handler = NewBlockSyncHandler(handler, web3.backend)
	handler = NewStateSyncHandler(handler, web3.backend)
	if err := web3.registerAPIs(rpcSrv, apis); err != nil {
		return err
	}
	return web3.start(web3.wsURI, handler)
}
