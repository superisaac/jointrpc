.PHONY: build all compile_proto test gofmt

protofiles := $(shell find ./proto -name '*.proto')
gofiles := $(shell find . -name '*.go')

tmp/protoc.ts: ${protofiles}
	mkdir -p tmp

compile_proto: tmp/protoc.ts ./compileproto.sh
	./compileproto.sh $<

bin/rpctube: ${gofiles}
	go build -o $@ rpctube.go

build: compile_proto bin/rpctube

test:
	go test -v github.com/superisaac/rpctube/jsonrpc
	go test -v github.com/superisaac/rpctube/jsonrpc/schema
	go test -v github.com/superisaac/rpctube/tube
	go test -v github.com/superisaac/rpctube/server

all: test build
	echo ${protofiles}

clean:
	rm -rf bin/rpctube
	rm tmp/protoc.ts

gofmt:
	go fmt client/*.go
	go fmt client/example/*.go
	go fmt server/*.go
	go fmt tube/*.go
	go fmt tube/handler/*.go
	go fmt jsonrpc/*.go
	go fmt jsonrpc/schema/*.go

