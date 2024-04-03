package main

import (
	"crypto/tls"
	"os"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/ledger"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/stratos-chain/crypto/ethsecp256k1"

	"github.com/stratosnet/stratos-chain/app"
	stratos "github.com/stratosnet/stratos-chain/types"
)

const EnvPrefix = ""

var (
	// Force to build with go1.20
	noUsage tls.CertificateVerificationError
)

func main() {
	registerDenoms()

	ledger.InitLedger(ethsecp256k1.MakePubKey, ledger.SignMode_SIGN_MODE_DIRECT)
	ledger.SetAppName("stratos")

	rootCmd, _ := NewRootCmd()
	if err := svrcmd.Execute(rootCmd, EnvPrefix, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

// RegisterDenoms registers the base and display denominations to the SDK.
func registerDenoms() {
	if err := sdk.RegisterDenom(stratos.Stos, sdkmath.LegacyOneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(stratos.Gwei, sdkmath.LegacyNewDecWithPrec(1, stratos.GweiDenomUnit)); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(stratos.Wei, sdkmath.LegacyNewDecWithPrec(1, stratos.WeiDenomUnit)); err != nil {
		panic(err)
	}
}
