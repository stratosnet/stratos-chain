package main

import (
	"go/doc/comment"
	"os"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stratosnet/stratos-chain/app"
	stratos "github.com/stratosnet/stratos-chain/types"
)

const EnvPrefix = ""

var (
	// Force to build with go1.19, because sorting algorithm has been rewritten since go1.19
	// The order of sorted results will be different between go1.18 & go1.19 if the values of the compared elements are equal
	doc comment.Doc
)

func main() {
	registerDenoms()
	//TODO: enable when customized cosmos-sdk pushed
	//ledger.InitLedger(ethsecp256k1.MakePubKey, ledger.SignMode_SIGN_MODE_DIRECT)

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
