#!/bin/bash
#
# Usage:
#   ./gen_proto.sh
#       - build all .proto into .pb.go
#   LIBRAPATH=/path/to/rust/libra ./gen_proto.sh
#       - copy .proto from rust LibraCore project, patch them, then build
#

proj_path="$( cd "$(dirname "$0")/.." ; pwd -P )"
echo $proj_path

copy_proto()
{
    echo "Copying .proto from package path $1 to $2"
    cp $libra_rust/$1/*.proto $proj_path/$2

    echo "Patching go_package name: $3"
    sed -i "/package/a option go_package = \"$3\";" $proj_path/$2/*.proto
}

if [[ -z "$GOPATH" ]]; then
    echo "GOPATH not defined. "
    echo "Have you installed go?"
    echo
    exit 1
fi

go_module_name=github.com/the729/go-libra
libra_rust=$1
[[ -z "$libra_rust" ]] && libra_rust=$LIBRAPATH

if [[ ! -z "$libra_rust" ]]; then
    echo "Copying .proto files from Libra Core: ${libra_rust}"
    copy_proto "types/src/proto" "types/proto" "$go_module_name/generated/pbtypes"
    copy_proto "mempool/src/proto/shared" "mempool/proto/shared" "$go_module_name/generated/pbmpshared"
    copy_proto "admission_control/admission_control_proto/src/proto" "admission_control/proto" "$go_module_name/generated/pbac"
fi

cd $proj_path
mkdir -p generated/pbtypes
protoc -I types/proto types/proto/*.proto --go_out=plugins=grpc,paths=source_relative:generated/pbtypes

mkdir -p generated/pbmpshared
protoc -I mempool/proto/shared mempool/proto/shared/*.proto --go_out=plugins=grpc,paths=source_relative:generated/pbmpshared

mkdir -p generated/pbac
protoc -I types/proto -I mempool/proto/shared -I admission_control/proto admission_control/proto/*.proto --go_out=plugins=grpc,paths=source_relative:generated/pbac
