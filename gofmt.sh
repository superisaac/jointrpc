#!/bin/bash

go fmt client/*.go
go fmt client/example/*.go
go fmt server/*.go
go fmt tube/*.go
go fmt tube/handler/*.go
go fmt jsonrpc/*.go
go fmt jsonrpc/schema/*.go



