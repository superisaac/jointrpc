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

function compile_go() {
    if [ ! -x $GOPATH/bin/protoc-gen-go ]
    then
        echo 'No plugin for golang installed, skip the go installation' >&2
        echo 'try go get github.com/golang/protobuf/protoc-gen-go' >&2
    else
        echo Compiling go interfaces...
        export GO_PATH=$GOPATH
        export GOBIN=$GOPATH/bin
        export PATH=$GOPATH/bin:$PATH

        protoc -I proto/ \
               --go_out=. \
               --go-grpc_out=. \
               $protofiles
        
        exit_if $?
        echo Done
    fi
}

function compile_python() {
    echo Compiling python interfaces...
    python -m grpc_tools.protoc -I proto/ \
           --python_out=python/ \
           --grpc_python_out=python/ \
           $protofiles
    exit_if $?

    if [ yes`which protoc-gen-grpclib_python` != yes ]; then
        python -m grpc_tools.protoc -I proto/ \
               --grpclib_python_out=python/ \
               $protofiles
        
        exit_if $?
    else
        echo 'No plugin for grpclib installed, skip the go installation' >&2
    fi

    for dir in $(find ./python -type d)
    do
        touch $dir/__init__.py
    done
    echo Done
}

compile_go
compile_python
