# syntax=docker/dockerfile:1.7

FROM golang:1.24-alpine AS builder

WORKDIR /workspace

COPY ./go.mod ./go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./ .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/products \
    ./cmd/products/main.go



# Final stage with a minimal image
FROM alpine:3.24

WORKDIR /workspace

COPY --from=builder /bin/products /bin/products
COPY --from=builder /workspace/config.yaml /workspace/config.yaml

EXPOSE 50051 9091

ENTRYPOINT ["/bin/products", "-config", "/workspace/config.yaml"]
