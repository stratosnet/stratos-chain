package rest

import (
	"bytes"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/bech32"
	"testing"
)

func TestAccAddrPrefix(t *testing.T) {

	hexStr := "c03661732294feb49caf6dc16c7cbb2534986d73"
	acc, err := sdk.AccAddressFromHex(hexStr)

	bech, err := bech32.ConvertAndEncode("st", acc.Bytes())

	if err != nil {
		t.Error(err)
	}

	// accAddr in bech32 is in form of "st1cqmxzuezjnltf890dhqkcl9my56fsmtnunn4z4",
	// where "st" is the hrp (human-reading prefix) and the last 4 digits work as checksum
	fmt.Println(bech)
	hrp, data, err := bech32.DecodeAndConvert(bech)

	if err != nil {
		t.Error(err)
	}
	if hrp != "st" {
		t.Error("Invalid hrp")
	}
	if !bytes.Equal(data, acc.Bytes()) {
		t.Error("Invalid decode")
	}
	fmt.Println(hrp)
	fmt.Println(data)
}
