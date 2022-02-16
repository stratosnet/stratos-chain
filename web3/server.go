package web3

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/stratosnet/stratos-chain/web3/api"
	"github.com/stratosnet/stratos-chain/web3/config"
	"github.com/stratosnet/stratos-chain/web3/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

// Server holds an RPC server that is served over HTTP
type Server struct {
	logger    log.Logger
	apiConfig *config.APIConfig
	services  map[string]types.Web3Service
}

func NewServer(apiConfig *config.APIConfig) *Server {
	server := &Server{
		logger:    log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("web3"),
		apiConfig: apiConfig,
		services:  make(map[string]types.Web3Service, 0),
	}
	server.defaultRegisterForAll()

	return server
}

// defaultRegisterForAll regs services from the namespaces
func (s *Server) defaultRegisterForAll() {
	s.RegisterService("web3", api.NewAPI())
}

func (s *Server) RegisterService(name string, service types.Web3Service) {
	s.services[name] = service
}

// ServiceList represents a cuurent list of registered servises
func (s *Server) ServiceList() map[string]types.Web3Service {
	return s.services
}

func RegisterApis(src *rpc.Server, apis map[string]types.Web3Service) error {
	for name, svc := range apis {
		err := src.RegisterName(name, svc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) start(rpcInfo interface{}, apis map[string]types.Web3Service) error {
	var (
		err               error
		uri               string
		enabled           bool
		availableAPINames []string
		availableAPIs     = make(map[string]types.Web3Service, 0)
		name              string
		handler           http.Handler
	)
	rpcSrv := rpc.NewServer()
	channel := make(chan error)
	timeout := make(chan error)
	keepAlive := true

	switch rpcCfg := rpcInfo.(type) {
	case *config.HTTPConfig:
		name = "HTTP"
		uri = fmt.Sprintf("%s:%d", rpcCfg.Addr, rpcCfg.Port)
		enabled = rpcCfg.Enabled
		availableAPINames = rpcCfg.API
		keepAlive = rpcCfg.KeepAlive
		handler = node.NewHTTPHandlerStack(rpcSrv, s.apiConfig.HTTPConfig.CORSDomain, s.apiConfig.HTTPConfig.VHosts)
	case *config.WSConfig:
		name = "WS"
		uri = fmt.Sprintf("%s:%d", rpcCfg.Addr, rpcCfg.Port)
		enabled = rpcCfg.Enabled
		availableAPINames = rpcCfg.API
		handler = rpcSrv.WebsocketHandler(s.apiConfig.WSConfig.Origins)
	default:
		s.logger.Info("Config for Web3 RPC not properly configured, skipping")
		return nil
	}

	srv := &http.Server{
		Addr:         uri,
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	srv.SetKeepAlivesEnabled(keepAlive)

	if !enabled {
		s.logger.Info("Web3 " + name + " RPC server not enabled, skipping")
		return nil
	}

	for _, apiName := range availableAPINames {
		apiService, ok := apis[apiName]
		if ok {
			availableAPIs[apiName] = apiService
		}
	}

	if err := RegisterApis(rpcSrv, availableAPIs); err != nil {
		return err
	}

	//Timeout Go routine
	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		timeout <- nil
	}()

	go func(ch chan error) {
		s.logger.Info("starting Web3 " + name + " RPC server on " + uri)
		err := srv.ListenAndServe()
		if err != nil {
			s.logger.Error("server error, details: %s", err)
		}
		srv.Shutdown(nil)
		rpcSrv.Stop()
		ch <- err
	}(channel)

	select {
	case err = <-channel:
	case err = <-timeout:
	}

	return err
}

func (s *Server) StartHTTP(apis map[string]types.Web3Service) error {
	if s.apiConfig == nil {
		s.logger.Info("Config for Web3 API not set, skipping")
		return nil
	}
	return s.start(s.apiConfig.HTTPConfig, apis)
}

func (s *Server) StartWS(apis map[string]types.Web3Service) error {
	if s.apiConfig == nil {
		s.logger.Info("Config for Web3 API not set, skipping")
		return nil
	}
	return s.start(s.apiConfig.WSConfig, apis)
}
