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
    sed -i "/^package/a option go_package = \"$3\";" $proj_path/$2/*.proto

    if [[ ! -z "$gopherjs" ]]; then
        echo "Patching gopherjs_package name: $3"
        sed -i "/^package/a option (gopherjs.gopherjs_package) = \"$3\";" $proj_path/$2/*.proto
        sed -i "/^package/a import \"github.com/johanbrandhorst/protobuf/proto/gopherjs.proto\";" $proj_path/$2/*.proto

        sed -i "s/google\/protobuf\/wrappers.proto/github.com\/johanbrandhorst\/protobuf\/ptypes\/wrappers\/wrappers.proto/" $proj_path/$2/*.proto
    fi
}

add_build_constraints()
{
    echo "Adding build constrants to $1"

    sed -i "1 i $2" $1
}

if [[ -z "$GOPATH" ]]; then
    echo "GOPATH not defined. "
    echo "Have you installed go?"
    echo
    exit 1
fi

go_module_name=github.com/the729/go-libra
libra_rust=$LIBRAPATH

if [[ $1 == "gopherjs" ]]; then
    gopherjs=yes
fi

if [[ ! -z "$libra_rust" ]]; then
    echo "Copying .proto files from Libra Core: ${libra_rust}"
    copy_proto "types/src/proto" "types/proto" "$go_module_name/generated/pbtypes"
    copy_proto "mempool/mempool-shared-proto/src/proto" "mempool/proto/shared" "$go_module_name/generated/pbmpshared"
    copy_proto "admission_control/admission-control-proto/src/proto" "admission_control/proto" "$go_module_name/generated/pbac"
fi

cd $proj_path

if [[ ! -z "$gopherjs" ]]; then
    mkdir -p generated/pbtypes
    protoc -I types/proto -I $GOPATH/src types/proto/*.proto \
        --gopherjs_out=plugins=grpc,import_path=pbtypes:$GOPATH/src
    
    add_build_constraints "generated/pbtypes/*.pb.gopherjs.go" "// +build js"

    mkdir -p generated/pbmpshared
    protoc -I mempool/proto/shared -I $GOPATH/src mempool/proto/shared/*.proto \
        --gopherjs_out=plugins=grpc,import_path=pbmpshared:$GOPATH/src

    add_build_constraints "generated/pbmpshared/*.pb.gopherjs.go" "// +build js"

    mkdir -p generated/pbac
    protoc -I types/proto -I mempool/proto/shared -I admission_control/proto -I $GOPATH/src admission_control/proto/*.proto \
        --gopherjs_out=plugins=grpc,import_path=pbac:$GOPATH/src

    add_build_constraints "generated/pbac/*.pb.gopherjs.go" "// +build js"

else
    mkdir -p generated/pbtypes
    protoc -I types/proto types/proto/*.proto --go_out=plugins=grpc,paths=source_relative:generated/pbtypes
    add_build_constraints "generated/pbtypes/*.pb.go" "// +build !js"

    mkdir -p generated/pbmpshared
    protoc -I mempool/proto/shared mempool/proto/shared/*.proto --go_out=plugins=grpc,paths=source_relative:generated/pbmpshared
    add_build_constraints "generated/pbmpshared/*.pb.go" "// +build !js"

    mkdir -p generated/pbac
    protoc -I types/proto -I mempool/proto/shared -I admission_control/proto admission_control/proto/*.proto --go_out=plugins=grpc,paths=source_relative:generated/pbac
    add_build_constraints "generated/pbac/*.pb.go" "// +build !js"
fi
