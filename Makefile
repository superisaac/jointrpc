protofiles := $(shell find ./proto -name '*.proto')
gofiles := $(shell find . -name '*.go')
protogofiles := intf/jointrpc/jointrpc.pb.go intf/jointrpc/jointrpc_grpc.pb.go
goarchs := linux-amd64 linux-arm linux-arm64 android-arm64 darwin-amd64 darwin-arm64

buildarchdirs := $(foreach a,$(goarchs),build/arch/jointrpc-$a)

protopyfiles := python/jointrpc/pb/jointrpc_pb2.py \
	python/jointrpc/pb/jointrpc_pb2_grpc.py \
	python/jointrpc/pb/jointrpc_grpc.py

# goflag := -gcflags=-G=3

build: compile_proto bin/jointrpc bin/jointrpc-server

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
	go build $(goflag) -o $@ jointrpc.go

bin/jointrpc-server: ${gofiles}
	go build $(goflag) -o $@ jointrpc_server.go

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
	go test -v github.com/superisaac/jointrpc/playbook

clean:
	rm -rf bin/jointrpc bin/jointrpc-server build dist

gofmt:
	go fmt datadir/*.go
	go fmt misc/*.go
	go fmt client/*.go
	go fmt client/example/*.go
	go fmt server/*.go
	go fmt rpcrouter/*.go
	go fmt dispatch/*.go
	go fmt playbook/*.go
	go fmt jsonrpc/*.go
	go fmt jsonrpc/schema/*.go
	go fmt encoding/*.go
	go fmt playbook/*.go
	go fmt cluster/bridge/*.go
	go fmt service/*.go
	go fmt service/builtin/*.go
	go fmt service/neighbor/*.go
	go fmt service/vars/*.go
	go fmt cmd/*.go
	go fmt *.go

install: bin/jointrpc bin/jointrpc-server
	install $^ /usr/local/bin


# cross build distributions of multiple targets
dist:
	@for arch in $(goarchs); do \
		$(MAKE) dist/jointrpc-$$arch.tar.gz; \
	done

dist/jointrpc-%.tar.gz: build/arch/jointrpc-%
	@mkdir -p dist
	tar czvf $@ $<

build/arch/jointrpc-%: ${gofiles}
	GOOS=$(shell echo $@|cut -d- -f 2) \
	GOARCH=$(shell echo $@|cut -d- -f 3) \
	go build $(goflag) -o $@/jointrpc jointrpc.go

	GOOS=$(shell echo $@|cut -d- -f 2) \
	GOARCH=$(shell echo $@|cut -d- -f 3) \
	go build $(goflag) -o $@/jointrpc-server jointrpc_server.go


.PHONY: build all compile_proto test gofmt dist $(goarchs)
.SECONDARY: $(buildarchdirs)
