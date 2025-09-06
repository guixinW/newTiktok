#!/bin/bash

# ==============================================================================
#  一键启动/停止 Kubernetes 服务端口转发的脚本
#
#  功能:
#    - 启动多个服务的端口转发，并在后台运行。
#    - 停止由本脚本启动的所有端口转发进程。
#
#  用法:
#    ./port_forward.sh start   # 启动所有端口转发
#    ./port_forward.sh stop    # 停止所有端口转发
#
# ==============================================================================

# --- 配置 ---
# 定义需要转发的服务和端口
# 格式: "service_name namespace local_port:service_port [local_port2...]"
SERVICES=(
  "mysql default 3306:3306"
  "redis default 6379:6379"
  "etcd default 2379:2379"
  "kafka-service default 19092:19092"
  "user-service default 50051:50051"
  "video-service default 50052:50052"
  "gateway-service default 8081:8081"
  "kong-kong-proxy kong 8000:80" # Forward Kong proxy to localhost:8000
)

# 用于存储后台进程ID (PID) 的文件
PID_FILE="port_forward_pids.txt"


# --- 函数定义 ---

# 启动所有端口转发
start_forwards() {
  # 检查 kubectl 是否可用
  if ! command -v kubectl &> /dev/null; then
    echo "错误: kubectl 命令未找到。请确保 kubectl 已安装并位于您的 PATH 中。"
    exit 1
  fi

  # 如果 PID 文件已存在，先执行停止操作，以防有残留进程
  if [ -f "$PID_FILE" ]; then
    echo "发现已存在的 PID 文件，正在尝试清理旧的转发进程..."
    stop_forwards
  fi

  echo "正在启动端口转发..."

  # 遍历服务列表并启动转发
  for S in "${SERVICES[@]}"; do
    # 将字符串分割为服务名、命名空间和端口列表
    # parts[0]=服务名, parts[1]=命名空间, parts[2...]=端口映射
    read -r -a parts <<< "$S"
    local SERVICE_NAME="service/${parts[0]}"
    local NAMESPACE="${parts[1]}"
    local PORTS=("${parts[@]:2}") # 获取从第三个元素开始的所有端口映射

    # 根据命名空间构建命令
    local KUBE_CMD="kubectl port-forward -n ${NAMESPACE} ${SERVICE_NAME} ${PORTS[@]}"

    # 在后台执行端口转发，并将输出重定向到 /dev/null
    # 使用 & 将命令置于后台
    eval "$KUBE_CMD" > /dev/null 2>&1 &

    # 获取刚刚启动的后台进程的 PID
    local PID=$!

    # 将 PID 写入文件，并打印信息
    echo "$PID" >> "$PID_FILE"
    echo "  > 转发 '${SERVICE_NAME}' (Namespace: $NAMESPACE, PID: $PID) 已在后台启动。"
  done

  echo -e "\n所有端口转发已在后台启动。您现在可以开始测试了。"
  echo "测试完成后，请运行 './port_forward.sh stop' 来清理所有进程。"
}

# 停止所有端口转发
stop_forwards() {
  if [ -f "$PID_FILE" ]; then
    echo "正在停止端口转发..."
    # 读取 PID 文件中的每一个 PID 并尝试停止它
    while read -r PID; do
      # 检查进程是否存在，然后停止它
      if ps -p "$PID" > /dev/null; then
        echo "  > 正在停止进程 PID: $PID"
        kill "$PID"
      else
        echo "  > 进程 PID: $PID 已不存在。"
      fi
    done < "$PID_FILE"

    # 清理 PID 文件
    rm "$PID_FILE"
    echo "所有转发进程已清理。"
  else
    echo "未找到 PID 文件。似乎没有正在运行的转发进程。"
  fi
}


# --- 主逻辑 ---

# 根据传入的第一个参数决定执行哪个函数
case "$1" in
  start)
    start_forwards
    ;;
  stop)
    stop_forwards
    ;;
  *)
    # 如果参数不正确，打印用法说明
    echo "用法: $0 {start|stop}"
    exit 1
    ;;
esac