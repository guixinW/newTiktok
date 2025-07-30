# --- Build Stage ---
# Use the official Go image as a builder.
# Using alpine for a smaller image size.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

ARG HTTP_PROXY
ARG HTTPS_PROXY

# Copy go.mod and go.sum files to download dependencies first.
# This leverages Docker's layer caching to speed up builds.
COPY go.mod go.sum ./
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    GOPROXY=https://goproxy.cn,direct go mod download
# Copy the rest of the application's source code.
COPY . .

# Build the Go application.
# The output will be a static binary named 'user_service'.
# CGO_ENABLED=0 is important for creating a static binary that can run in a minimal container.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user_service ./cmd/user_service

# --- Final Stage ---
# Use a minimal base image for the final container.
# alpine is a good choice for its small size.
FROM alpine:latest

# Set the working directory.
WORKDIR /app

# Copy the built binary from the builder stage.
COPY --from=builder /app/user_service .

# Expose the port the gRPC server listens on.
# I'll assume the default gRPC port 50051. If it's different, we can change this.
EXPOSE 50051

# The command to run when the container starts.
CMD ["/app/user_service"]