package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewFileInfo constructor
func NewFileInfo(height *sdk.Int, reporter, uploader string) FileInfo {
	return FileInfo{
		Height:   height,
		Reporter: reporter,
		Uploader: uploader,
	}
}

// MustMarshalFileInfo returns the fileInfo's bytes. Panics if fails
func MustMarshalFileInfo(cdc codec.Codec, file FileInfo) []byte {
	return cdc.MustMarshalLengthPrefixed(&file)
}

// MustUnmarshalFileInfo unmarshal a file's info from a store value. Panics if fails
func MustUnmarshalFileInfo(cdc codec.Codec, value []byte) FileInfo {
	file, err := UnmarshalFileInfo(cdc, value)
	if err != nil {
		panic(err)
	}
	return file
}

// UnmarshalFileInfo unmarshal a file's info from a store value
func UnmarshalFileInfo(cdc codec.Codec, value []byte) (fi FileInfo, err error) {
	err = cdc.UnmarshalLengthPrefixed(value, &fi)
	return fi, err
}
