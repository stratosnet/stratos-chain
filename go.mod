module github.com/stratosnet/stratos-chain

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.45.1
	github.com/cosmos/ibc-go/v3 v3.0.0-rc0
	github.com/gogo/protobuf v1.3.3
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.7
	github.com/tendermint/btcd v0.1.1 // indirect
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15 // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/ethereum/go-ethereum v1.10.16
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.16.5
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20220315194320-039c03cc5b86 // indirect
	gopkg.in/yaml.v2 v2.4.0
	github.com/ReneKroon/ttlcache/v2 v2.7.0
	github.com/golang/mock v1.4.3 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
)
replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)