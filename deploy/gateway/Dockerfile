# Stage 1: Builder
# This stage builds the Go application into a static binary.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Accept proxy settings as build arguments
ARG HTTP_PROXY
ARG HTTPS_PROXY

# Copy go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./

# Use the proxy arguments when running go mod download
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    go mod download

# Copy the entire source code
COPY ../../ .

# Build the application, creating a statically linked binary.
# Ensure the proxy is also available for the build step if needed.
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/gateway ./cmd/gateway/main.go

# Stage 2: Final
# This stage creates the final, small, and secure image.
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/gateway .

# Expose the port the HTTP server listens on
EXPOSE 8081

# The command to run the application
CMD ["/app/gateway"]