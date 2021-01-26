.PHONY: build all compile_proto test gofmt

protofiles := $(shell find ./proto -name '*.proto')
gofiles := $(shell find . -name '*.go')
protogofiles := intf/jointrpc/jointrpc.pb.go intf/jointrpc/jointrpc_grpc.pb.go

protopyfiles := python/jointrpc_pb2.py \
	python/jointrpc_pb2_grpc.py \
	python/jointrpc_grpc.py

build: compile_proto bin/jointrpc

all: test build

intf/jointrpc/%.pb.go intf/jointrpc/%_grpc.pb.go: proto/%.proto
	protoc -I proto/ --go_out=. --go-grpc_out=. $<

python/%_pb2.py python/%_pb2_grpc.py python/%_grpc.py: proto/%.proto
	python -m grpc_tools.protoc -I proto/ \
			--python_out=python/ \
			--grpc_python_out=python/ \
			--grpclib_python_out=python/ $<

compile_proto: $(protogofiles) $(protopyfiles)

bin/jointrpc: ${gofiles}
	go build -o $@ jointrpc.go

test:
	go test -v github.com/superisaac/jointrpc/datadir
	go test -v github.com/superisaac/jointrpc/jsonrpc
	go test -v github.com/superisaac/jointrpc/jsonrpc/schema
	go test -v github.com/superisaac/jointrpc/rpcrouter
	go test -v github.com/superisaac/jointrpc/client
	go test -v github.com/superisaac/jointrpc/server
	go test -v github.com/superisaac/jointrpc/service/builtin
	go test -v github.com/superisaac/jointrpc/service/neighbor
	go test -v github.com/superisaac/jointrpc/cluster/bridge

clean:
	rm -rf bin/jointrpc

gofmt:
	go fmt datadir/*.go
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
	go fmt service/neighbor/*.go
	go fmt service/vars/*.go
	go fmt command/*.go

install: bin/jointrpc
	install $< /usr/local/bin
