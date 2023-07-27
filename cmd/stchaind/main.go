package main

import (
	"go/doc/comment"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/ledger"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stratosnet/stratos-chain/crypto/ethsecp256k1"

	"github.com/stratosnet/stratos-chain/app"
	stratos "github.com/stratosnet/stratos-chain/types"
)

var (
	// Force to build with go1.19, because sorting algorithm has been rewritten since go1.19
	// The order of sorted results will be different between go1.18 & go1.19 if the values of the compared elements are equal
	doc comment.Doc
)

func main() {
	registerDenoms()
	ledger.InitLedger(ethsecp256k1.MakePubKey, signing.SignMode_SIGN_MODE_DIRECT)

	rootCmd, _ := NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
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
	if err := sdk.RegisterDenom(stratos.Stos, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(stratos.Gwei, sdk.NewDecWithPrec(1, stratos.GweiDenomUnit)); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(stratos.Wei, sdk.NewDecWithPrec(1, stratos.WeiDenomUnit)); err != nil {
		panic(err)
	}
}
