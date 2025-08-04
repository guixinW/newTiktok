#!/bin/bash

# set -e: 如果任何命令执行失败，脚本将立即退出，保证了操作的原子性。
set -e

# ==================== 配置区 ====================

# 你的 kind 集群名称，请确保和 kind-config.yaml 中的 name 一致。
CLUSTER_NAME="dev-cluster"

# 你需要推送到集群的镜像列表。
# 请将你的微服务、基础镜像等都添加到这里。
IMAGES=(
  "user-service:latest"
  "video-service:latest"
)

# ==================== 脚本执行区 ====================

# 1. 检查 kind 集群是否存在且在运行中
echo "🔎 正在检查 kind 集群 '${CLUSTER_NAME}' 是否在运行..."
if ! kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
  echo "❌ 错误: 未找到名为 '${CLUSTER_NAME}' 的 kind 集群。"
  echo "   请先使用 'kind create cluster --name ${CLUSTER_NAME} --config ...' 命令创建集群。"
  exit 1
fi
echo "✅ 集群 '${CLUSTER_NAME}' 已就绪。"
echo ""


# 2. 遍历镜像列表并执行推送命令
echo "🚀 开始将镜像推送到集群 '${CLUSTER_NAME}'..."
echo "--------------------------------------------------------"

for img in "${IMAGES[@]}"; do
  echo "🚢 正在推送镜像: ${img} ..."

  # 检查镜像是否在本地存在
  if ! docker image inspect "${img}" &> /dev/null; then
    echo "⚠️  警告: 镜像 '${img}' 在本地不存在，正在尝试拉取..."
    docker pull "${img}"
  fi

  # 执行 kind load 命令
  kind load docker-image "${img}" --name "${CLUSTER_NAME}"
  echo "✅ 镜像 '${img}' 推送成功。"
  echo "--------------------------------------------------------"
done

echo "🎉 所有镜像已成功推送到集群！"