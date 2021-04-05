.PHONY: build all compile_proto test gofmt

protofiles := $(shell find ./proto -name '*.proto')
gofiles := $(shell find . -name '*.go')
protogofiles := intf/jointrpc/jointrpc.pb.go intf/jointrpc/jointrpc_grpc.pb.go

protopyfiles := python/jointrpc/pb/jointrpc_pb2.py \
	python/jointrpc/pb/jointrpc_pb2_grpc.py \
	python/jointrpc/pb/jointrpc_grpc.py

build: compile_proto bin/jointrpc

all: test build

intf/jointrpc/%.pb.go intf/jointrpc/%_grpc.pb.go: proto/%.proto
	protoc -I proto/ --go_out=. --go-grpc_out=. $<

python/jointrpc/pb/%_pb2.py python/jointrpc/pb/%_pb2_grpc.py python/jointrpc/pb/%_grpc.py: proto/%.proto
	@python -m grpc_tools.protoc -I proto/ \
			--python_out=python/jointrpc/pb/ \
			--grpc_python_out=python/jointrpc/pb/ \
			--grpclib_python_out=python/jointrpc/pb/ $<

	@for f in $(protopyfiles); do sed -ie 's/import jointrpc_pb2/from jointrpc.pb import jointrpc_pb2/g' $$f; done
	@for f in $(shell find python/jointrpc -name '*.pye'); do rm $$f; done

compile_proto: $(protogofiles) $(protopyfiles)

bin/jointrpc: ${gofiles}
	go build -o $@ jointrpc.go

build_arch: ${gofiles}
	GOOS=linux GOARCH=amd64 go build -o build/arch/jointrpc.linux-amd64 jointrpc.go
	GOOS=linux GOARCH=arm go build -o build/arch/jointrpc.linux-arm jointrpc.go
	GOOS=linux GOARCH=arm64 go build -o build/arch/jointrpc.linux-arm64 jointrpc.go
	GOOS=android GOARCH=arm64 go build -o build/arch/jointrpc.android-arm64 jointrpc.go
	GOOS=darwin GOARCH=amd64 go build -o build/arch/jointrpc.darwin-amd64 jointrpc.go
	GOOS=darwin GOARCH=arm64 go build -o build/arch/jointrpc.darwin-arm64 jointrpc.go

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
	go fmt dispatch/*.go
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
