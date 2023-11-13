# this script is for generating protobuf files for the new google.golang.org/protobuf API

set -eo pipefail

protoc_install_gopulsar() {
  go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@v1.0.0-beta.2
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
}

protoc_install_gopulsar

echo "Cleaning API directory"
(cd api; find ./ -type f \( -iname \*.pulsar.go -o -iname \*.pb.go -o -iname \*.cosmos_orm.go -o -iname \*.pb.gw.go \) -delete; find . -empty -type d -delete; cd ..)

echo "Generating API"
cd proto

module_dirs=$(find ./stratos/ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $module_dirs; do
  buf generate --template buf.gen.pulsar.yaml --path $dir
done

chmod 755 ../api -R


#
## this script is for generating protobuf files for the new google.golang.org/protobuf API
#echo "Generating API module"
#(cd proto; buf generate --template buf.gen.pulsar.yaml)
#

