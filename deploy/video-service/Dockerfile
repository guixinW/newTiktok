# Dockerfile for video-service
# IMPORTANT: This file is designed to be run from the project root directory.
# Example command: docker build -f deploy/video-service/Dockerfile .

# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

ARG HTTP_PROXY
ARG HTTPS_PROXY

# 1. Copy dependency files for caching
COPY go.mod go.sum ./

# 2. Download dependencies
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    go mod download

# 3. Copy only necessary source code for video_service
# This is the key optimization to isolate build caches between services.
COPY cmd/video_service/ ./cmd/video_service/
COPY internal/video/ ./internal/video/
COPY pkg/ ./pkg/

# 4. Build the video_service application
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/video_service ./cmd/video_service/main.go


# Stage 2: Final
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/video_service .

# Expose the port for the video service gRPC server
EXPOSE 50052

# The command to run the application
CMD ["/app/video_service"]