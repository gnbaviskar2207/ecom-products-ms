.PHONY: tools gen-all buf-gen gov-gen lint test test-race build run docker-build


GOBIN := $(shell go env GOPATH)/bin

tools:
	go install github.com/go-delve/delve/cmd/dlv@v1.24.2
	go install github.com/jmattheis/goverter/cmd/goverter@v1.9.3
	go install golang.org/x/tools/gopls@v0.20.0
	
gen-all: buf-gen gov-gen
	
buf-gen: remove-stale-buf-gen-files buf-generate

gov-gen: remove-stale-goverter-gen-files goverter-gen

buf-generate:
	PATH="$(GOBIN):$$PATH" buf generate

goverter-gen:
	PATH="$(GOBIN):$$PATH" goverter gen ./internal/transform

remove-stale-buf-gen-files:
	find ./gen -mindepth 1 -delete

remove-stale-goverter-gen-files:
	find ./internal/transform/generated -mindepth 1 -delete

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