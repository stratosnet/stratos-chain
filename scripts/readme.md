# Generate protobuf files

## Install buf

`go install github.com/bufbuild/buf/cmd/buf@v1.28.0`

## Install protobuf (Cosmos's fork of gogo/protobuf)

https://github.com/cosmos/gogoproto/tree/v1.4.11

## Install protoc-gen-go-pulsar

https://github.com/cosmos/cosmos-proto/tree/v1.0.0-beta.3

## Install protoc-gen-go, protoc-gen-grpc-gateway, protoc-gen-swagger

https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.16.0

## Export proto dependency

`$ buf export buf.build/cosmos/cosmos-sdk:v0.47.0 --output ./proto`
