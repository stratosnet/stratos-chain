package rpc

import (
	"context"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stratosnet/stratos-chain/server/config"
	"github.com/tendermint/tendermint/libs/log"
)

type Web3Server struct {
	httpURI string
	wsURI   string
	enabled bool
	logger  log.Logger
}

func NewWeb3Server(cfg config.Config, logger log.Logger) *Web3Server {
	return &Web3Server{
		httpURI: cfg.JSONRPC.Address,
		wsURI:   cfg.JSONRPC.WsAddress,
		enabled: cfg.JSONRPC.Enable,
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
	for _, api := range apis {
		if err := server.RegisterName(api.Namespace, api.Service); err != nil {
			web3.logger.Error(
				"failed to register service in JSON RPC namespace",
				"namespace", api.Namespace,
				"service", api.Service,
			)
			return err
		}
	}
	return nil
}

func (web3 *Web3Server) StartHTTP(apis []rpc.API) error {
	rpcSrv := rpc.NewServer()
	handler := node.NewHTTPHandlerStack(rpcSrv, []string{"*"}, []string{"localhost"}) // TODO: Replace cors and vshosts from config
	if err := web3.registerAPIs(rpcSrv, apis); err != nil {
		return err
	}
	return web3.start(web3.httpURI, handler)
}

func (web3 *Web3Server) StartWS(apis []rpc.API) error {
	rpcSrv := rpc.NewServer()
	handler := rpcSrv.WebsocketHandler([]string{}) // TODO: Add config origins
	if err := web3.registerAPIs(rpcSrv, apis); err != nil {
		return err
	}
	return web3.start(web3.wsURI, handler)
}
