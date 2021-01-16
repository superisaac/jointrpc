.PHONY: build all compile_proto test gofmt

protofiles := $(shell find ./proto -name '*.proto')
gofiles := $(shell find . -name '*.go')

build: compile_proto bin/jointrpc

all: test build

./tmp/protoc.ts: ${protofiles} ./compileproto.sh
	mkdir -p tmp
	./compileproto.sh
	echo compile >./tmp/protoc.ts

compile_proto: ./tmp/protoc.ts

bin/jointrpc: ${gofiles}
	go build -o $@ jointrpc.go

test:
	go test -v github.com/superisaac/jointrpc/jsonrpc
	go test -v github.com/superisaac/jointrpc/jsonrpc/schema
	go test -v github.com/superisaac/jointrpc/rpcrouter
	go test -v github.com/superisaac/jointrpc/client
	go test -v github.com/superisaac/jointrpc/server
	go test -v github.com/superisaac/jointrpc/service/builtin
	go test -v github.com/superisaac/jointrpc/service/mirror
	go test -v github.com/superisaac/jointrpc/cluster/bridge

clean:
	rm -rf bin/jointrpc
	rm ./tmp/protoc.ts

gofmt:
	go fmt misc/*.go
	go fmt client/*.go
	go fmt client/example/*.go
	go fmt server/*.go
	go fmt rpcrouter/*.go
	go fmt rpcrouter/handler/*.go
	go fmt jsonrpc/*.go
	go fmt jsonrpc/schema/*.go
	go fmt encoding/*.go
	go fmt cluster/bridge/*.go
	go fmt service/*.go
	go fmt service/builtin/*.go
	go fmt service/mirror/*.go
	go fmt service/vars/*.go
	go fmt command/*.go

install: bin/jointrpc
	install $< /usr/local/bin
