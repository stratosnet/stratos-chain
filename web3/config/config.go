package config

type APIConfig struct {
	HTTPConfig *HTTPConfig `toml:"http"`
	WSConfig   *WSConfig   `toml:"ws"`
}

func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		HTTPConfig: &HTTPConfig{
			Addr:       "127.0.0.1",
			Port:       18545,
			KeepAlive:  true,
			API:        []string{"web3"},
			CORSDomain: make([]string, 0),
			VHosts:     []string{"127.0.0.1"},
		},
		WSConfig: &WSConfig{
			Addr:    "127.0.0.1",
			Port:    18546,
			API:     []string{"web3"},
			Origins: make([]string, 0),
		},
	}
}

type HTTPConfig struct {
	Enabled bool   `toml:"enabled" desc:"Enable the HTTP-RPC server"`
	Addr    string `toml:"addr" desc:"HTTP-RPC server listening interface (default: \"localhost\")"`
	Port    int    `toml:"port" desc:"HTTP-RPC server listening port (default: 8545)"`

	KeepAlive  bool     `toml:"keepalive" desc:"API's keep alive. Note: only very resource-constrained environments or servers in the process of shutting down should disable them (default: true)"`
	API        []string `toml:"api" desc:"API's offered over the HTTP-RPC interface"`
	CORSDomain []string `toml:"corsdomain" desc:"Comma separated list of domains from which to accept cross origin requests (browser enforced)"`
	VHosts     []string `toml:"vhosts" desc:"Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default: \"localhost\")"`
}

type WSConfig struct {
	Enabled bool   `toml:"enabled" desc:"Enable the WS-RPC server"`
	Addr    string `toml:"addr" desc:"WS-RPC server listening interface (default: \"localhost\")"`
	Port    int    `toml:"port" desc:"WS-RPC server listening port (default: 8546)"`

	API     []string `toml:"api" desc:"API's offered over the WS-RPC interface"`
	Origins []string `toml:"origins" desc:"Origins from which to accept websockets requests"`
}
