# Stage 1: Builder
# 此阶段负责编译 Go 应用，生成一个静态链接的二进制文件。
FROM golang:1.23-alpine AS builder

# 容器内的工作目录
WORKDIR /app

# 接收代理设置作为构建参数，以便在需要时使用
ARG HTTP_PROXY
ARG HTTPS_PROXY

# --- 缓存优化步骤 ---

# 1. 仅复制依赖管理文件
COPY go.mod go.sum ./

# 2. 下载依赖
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    go mod download

# 3. 精细化地复制构建所需的源代码
COPY cmd/user-service/ ./cmd/user-service/
COPY internal/user/ ./internal/user
COPY internal/common/ ./internal/common/

# 4. 编译应用，并压缩二进制文件
# 新增 -ldflags="-s -w" 来剥离调试信息，减小二进制文件体积
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o /app/user-service ./cmd/user-service/main.go


# Stage 2: Final
# 此阶段负责构建最终的运行镜像，它非常小且安全。
# 使用 Google 的 distroless 镜像作为基础，它比 alpine 更小、更安全
FROM gcr.io/distroless/static-debian11

# 最终镜像的工作目录
WORKDIR /app

# 从 builder 阶段仅复制编译好的二进制文件
COPY --from=builder /app/user-service .

# 暴露 gRPC 服务监听的端口
EXPOSE 50051

# 容器启动时运行的命令
CMD ["/app/user-service"]
