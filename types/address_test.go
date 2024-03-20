package types

import (
	"encoding/base64"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/stratosnet/stratos-chain/app"

	"github.com/stretchr/testify/require"
)

func initCodec() codec.Codec {
	encodingConfig := app.MakeTestEncodingConfig()
	cdc := encodingConfig.Codec
	return cdc
}

// go test -v address_test.go address.go config.go coin.go -run TestSdsPubKeyToBech32
func TestSdsPubKeyToBech32(t *testing.T) {
	tests := []struct {
		name           string
		pubkeyJson     string
		expectedBech32 string
		wantErr        bool
	}{
		{"test1", "{\"@type\":\"/cosmos.crypto.ed25519.PubKey\",\"key\":\"2OAeLO0+KrBkSxuFKU1ofJqGb4RtA8GpD8XCZlMYw2A=\"}",
			"stsdspub1mrsput8d8c4tqeztrwzjjntg0jdgvmuyd5pur2g0chpxv5cccdsqvayhan", false},
	}
	cdc := initCodec()
	cfg := GetConfig()
	cfg.Seal()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var pubKey cryptotypes.PubKey
			err := cdc.UnmarshalInterfaceJSON([]byte(tt.pubkeyJson), &pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf(err.Error())
			}

			bech32Str, err := SdsPubKeyToBech32(pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf(err.Error())
			}

			require.Equal(t, bech32Str, tt.expectedBech32)
		})
	}
}

// go test -v address_test.go address.go config.go coin.go -run TestSdsPubKeyFromBech32
func TestSdsPubKeyFromBech32(t *testing.T) {
	tests := []struct {
		name           string
		bech32PubKey   string
		expectedBase64 string
		wantErr        bool
	}{
		{"test1", "stsdspub1mrsput8d8c4tqeztrwzjjntg0jdgvmuyd5pur2g0chpxv5cccdsqvayhan",
			"2OAeLO0+KrBkSxuFKU1ofJqGb4RtA8GpD8XCZlMYw2A=", false},
	}

	cfg := GetConfig()
	cfg.Seal()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubKey, err := SdsPubKeyFromBech32(tt.bech32PubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf(err.Error())
			}
			base64Str := base64.StdEncoding.EncodeToString(pubKey.Bytes())
			require.Equal(t, base64Str, tt.expectedBase64)
		})
	}

}

func TestSdsAddress_Unmarshal(t *testing.T) {

	tests := []struct {
		name    string
		aa      string
		args    string
		wantErr bool
	}{
		{"test1", "stsds14c3em44vlh276cujnr2ez802uyjyeqrrsu9fuh", "stsds14c3em44vlh276cujnr2ez802uyjyeqrrsu9fuh", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := SdsAddressFromBech32(tt.aa)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			aa := &SdsAddress{}
			bz, err := addr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = aa.UnmarshalJSON(bz)
			if !aa.Equals(addr) || (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			addr, err := SdsAddressFromBech32(tt.aa)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			aa := &SdsAddress{}
			bz, err := addr.MarshalYAML()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = aa.UnmarshalYAML([]byte(bz.(string)))
			if !aa.Equals(addr) || (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
