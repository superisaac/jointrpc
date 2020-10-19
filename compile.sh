#!/bin/bash

protofiles=$(find ./proto -name '*.proto')

function exit_if() {
    extcode=$1
    msg=$2
    if [ $extcode -ne 0 ]
    then
        if [ "msg$msg" != "msg" ]; then
            echo $msg >&2
        fi
        exit $extcode
    fi
}

if [ ! -x $GOPATH/bin/protoc-gen-go ]
then
    echo 'No plugin for golang installed, skip the go installation' >&2
    echo 'try go get github.com/golang/protobuf/protoc-gen-go' >&2
else
    echo Compiling go interfaces...
    export GO_PATH=$GOPATH
    export GOBIN=$GOPATH/bin
    export PATH=$GOPATH/bin:$PATH

    # protoc -I proto/ \
    #        --go_out=. --go_opt=paths=source_relative \
    #        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    #        $protofiles
    protoc -I proto/ \
           --go_out=. \
           --go-grpc_out=. \
           $protofiles
    
    exit_if $?
    echo Done
fi
