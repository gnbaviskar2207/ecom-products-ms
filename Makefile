.PHONY: tools generate lint test test-race build run docker-build


GOBIN := $(shell go env GOPATH)/bin

tools:
	go install github.com/go-delve/delve/cmd/dlv@v1.24.2
	go install github.com/jmattheis/goverter/cmd/goverter@v1.9.3
	go install golang.org/x/tools/gopls@v0.20.0
	
generate:
	PATH="$(GOBIN):$$PATH" buf generate
	
lint:
	buf lint
	go vet ./...

test:
	go test -cover ./...

test-race:
	go test -race ./...

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/products ./cmd/products

run:
	go run ./cmd/products/main.go --config=./config.yaml

# docker-build:
# 	docker build -f ./Dockerfile -t dropshipping-product-service:local .