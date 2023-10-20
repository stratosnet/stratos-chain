# Generate protobuf files

## Install protobuf (Cosmos's fork of gogo/protobuf)

https://github.com/cosmos/gogoproto/tree/v1.4.10

## Install protoc-gen-go-pulsar

https://github.com/cosmos/cosmos-proto/tree/v1.0.0-beta.2

## Install protoc-gen-go, protoc-gen-grpc-gateway, protoc-gen-swagger

https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.16.0

## Export proto dependency

`$ buf export buf.build/cosmos/cosmos-sdk:v0.47.0 --output ./proto`
