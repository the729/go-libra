#!/bin/bash

print_help()
{
    echo "Build client binary and connect to testnet."
    echo "\`$0 -r|--release\` to use release build or"
    echo "\`$0\` to use debug build."
}

copy_proto()
{
    echo "Copying .proto from package path $1 to $2"
    cp $libra_rust/$1/*.proto $GOPATH/src/$2

    echo "Patching go_package name: $3"
    sed -i "/package/a option go_package = \"$3\";" $GOPATH/src/$2/*.proto
}

if [[ -z "$GOPATH" ]]; then
    echo "GOPATH not defined. "
    echo "Have you installed go?"
    echo
    exit 1
fi

go_package_base=github.com/the729/go-libra
libra_rust=$1
[[ -z "$libra_rust" ]] && libra_rust=$LIBRAPATH

if [[ ! -z "$libra_rust" ]]; then
    echo "Copying .proto files from Libra Core: ${libra_rust}"
    copy_proto "types/src/proto" "$go_package_base/types/proto" "$go_package_base/generated/types"
    copy_proto "mempool/src/proto/shared" "$go_package_base/mempool/proto/shared" "$go_package_base/generated/mpshared"
    copy_proto "admission_control/admission_control_proto/src/proto" "$go_package_base/admission_control/proto" "$go_package_base/generated/ac"
fi

cd $GOPATH/src/$go_package_base
protoc -I types/proto types/proto/*.proto --go_out=plugins=grpc:$GOPATH/src
protoc -I mempool/proto/shared mempool/proto/shared/*.proto --go_out=plugins=grpc:$GOPATH/src
protoc -I types/proto -I mempool/proto/shared -I admission_control/proto admission_control/proto/*.proto --go_out=plugins=grpc:$GOPATH/src
