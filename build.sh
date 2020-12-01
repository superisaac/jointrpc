#!/bin/bash

rm -rf bin/*
exec go build -o bin rpctube.go

