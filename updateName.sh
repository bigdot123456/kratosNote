#!/bin/sh
cd api
rm -f api*.go
go generate
cd ../internal/di
go generate
cd ../..
rm -f ./main
go build ./cmd/main.go
./main --conf ./configs