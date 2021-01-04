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
	go test -v github.com/superisaac/jointrpc/joint
	go test -v github.com/superisaac/jointrpc/server

clean:
	rm -rf bin/jointrpc
	rm ./tmp/protoc.ts

gofmt:
	go fmt misc/*.go
	go fmt client/*.go
	go fmt client/example/*.go
	go fmt server/*.go
	go fmt joint/*.go
	go fmt joint/handler/*.go
	go fmt jsonrpc/*.go
	go fmt jsonrpc/schema/*.go
	go fmt encoding/*.go
	go fmt cluster/bridge/*.go
	go fmt cluster/mirror/*.go

install: bin/jointrpc
	install $< /usr/local/bin
