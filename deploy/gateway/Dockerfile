# Dockerfile for gateway service
# IMPORTANT: This file is designed to be run from the project root directory.
# Example command: docker build -f deploy/gateway/Dockerfile .

# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

ARG HTTP_PROXY
ARG HTTPS_PROXY

# 1. Copy dependency files (path is relative to build context, i.e., project root)
COPY go.mod go.sum ./

# 2. Download dependencies
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    go mod download

# 3. Copy only necessary source code for the gateway service
# This replaces the inefficient "COPY ../../ ."
COPY cmd/gateway/ ./cmd/gateway/
# The following line assumes the gateway has its own internal logic. If not, remove it.
COPY internal/gateway/ ./internal/gateway/
# The gateway typically needs all shared packages and protobuf definitions.
COPY pkg/ ./pkg/

# 4. Build the gateway service
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/gateway ./cmd/gateway/main.go


# Stage 2: Final
# This stage is already well-optimized. No changes needed.
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/gateway .

# Expose the port the HTTP gateway listens on
EXPOSE 8081

# The command to run the service
CMD ["/app/gateway"]