.PHONY: build all compile_proto test gofmt

protofiles := $(shell find ./proto -name '*.proto')
gofiles := $(shell find . -name '*.go')

build: compile_proto bin/rpctube

all: test build

/tmp/protoc.ts: ${protofiles} ./compileproto.sh
	./compileproto.sh
	echo compile >/tmp/protoc.ts

compile_proto: /tmp/protoc.ts

bin/rpctube: ${gofiles}
	go build -o $@ rpctube.go



test:
	go test -v github.com/superisaac/rpctube/jsonrpc
	go test -v github.com/superisaac/rpctube/jsonrpc/schema
	go test -v github.com/superisaac/rpctube/tube
	go test -v github.com/superisaac/rpctube/server

clean:
	rm -rf bin/rpctube
	rm /tmp/protoc.ts

gofmt:
	go fmt client/*.go
	go fmt client/example/*.go
	go fmt server/*.go
	go fmt tube/*.go
	go fmt tube/handler/*.go
	go fmt jsonrpc/*.go
	go fmt jsonrpc/schema/*.go

