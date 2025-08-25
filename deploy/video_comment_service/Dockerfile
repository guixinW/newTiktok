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
# 这是第一层，只有在 go.mod 或 go.sum 文件发生变化时，这一层以及之后的缓存才会失效。
COPY go.mod go.sum ./

# 2. 下载依赖
# 利用上一步的缓存，只有在依赖变更时才需要重新下载。
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    go mod download

# 3. 精细化地复制构建所需的源代码
# 这是最关键的优化：只复制当前服务（video_comment_service）真正依赖的目录。
COPY cmd/video_comment_service/ ./cmd/video_comment_service/
COPY internal/video_comment/ ./internal/video_comment
COPY pkg/ ./pkg

# 4. 编译应用
# CGO_ENABLED=0 创建一个静态链接的二进制文件，不依赖任何系统的 C 库。
# GOOS=linux 指定目标操作系统为 Linux。
# -a 强制重新构建。
# -installsuffix cgo 避免使用标准库的缓存。
RUN export http_proxy=${HTTP_PROXY} && \
    export https_proxy=${HTTPS_PROXY} && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/video_comment_service ./cmd/video_comment_service/main.go


# Stage 2: Final
# 此阶段负责构建最终的运行镜像，它非常小且安全。
FROM alpine:latest
# 建议：为了生产环境的稳定性，可以固定版本，如 FROM alpine:3.20

# 最终镜像的工作目录
WORKDIR /app

# 注意：应用的配置（如数据库地址、密码等）应通过环境变量或挂载的配置文件来管理，
# 而不是直接打包进镜像。你的代码 `config.LoadConfig()` 需要适配这种方式。

# 从 builder 阶段仅复制编译好的二进制文件
COPY --from=builder /app/video_comment_service .

# 暴露 gRPC 服务监听的端口
EXPOSE 50054

# 容器启动时运行的命令
CMD ["/app/video_comment_service"]
