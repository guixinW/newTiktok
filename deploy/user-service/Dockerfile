# Stage 1: Builder
# This stage builds the Go application into a static binary.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./
# Unset proxy settings that may be inherited from the host
RUN GOPROXY=direct go mod download

# Copy the entire source code
COPY . .

# Build the application, creating a statically linked binary.
# The output is placed in the /app directory.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/user_service ./cmd/user_service/main.go

# Stage 2: Final
# This stage creates the final, small, and secure image.
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/user_service .

# The application looks for 'config.yaml' in the current directory.
# This file will be provided by a Kubernetes ConfigMap.

# Expose the port the gRPC server listens on
EXPOSE 50051

# The command to run the application
CMD ["/app/user_service"]
