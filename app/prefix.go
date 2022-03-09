package app

import (
	stratos "github.com/stratosnet/stratos-chain/types"
)

func SetConfig() {
	config := stratos.GetConfig()

	config.Seal()
}
