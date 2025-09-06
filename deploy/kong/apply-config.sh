#!/bin/bash

# 脚本：通过 kubectl apply 部署或更新 Kong 网关。
#
# !!! 注意：请务必在项目的根目录下运行此脚本 !!!

set -e

# --- 配置 ---
KONG_NAMESPACE="kong"
# 定义所有文件的相对路径 (相对于项目根目录)
CONFIG_FILE="./deploy/kong/kong.yaml"
PROTO_FILE="./deploy/kong/user.proto"
DEPLOYMENT_FILE="./deploy/kong/kong-deployment.yaml"

CONFIG_MAP_NAME="kong-declarative-config"

# --- 脚本主体 ---

echo "▶️ 开始 Kong 部署/更新流程..."

# 1. 检查前提条件：Kind 集群是否在运行
echo "🔎 正在检查 Kind 集群状态..."
if ! kind get clusters | grep -q "dev-cluster"; then
    echo "❌ 错误：未找到正在运行的 Kind 集群。"
    exit 1
fi
echo "✅ Kind 集群正在运行。"

# 2. 确保 Kong 的命名空间存在
echo "🔎 正在确保命名空间 '$KONG_NAMESPACE' 存在..."
kubectl create namespace "$KONG_NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# 3. 从 .proto 和 kong.yaml 创建或更新 ConfigMap
echo "🔄 正在更新 ConfigMap '$CONFIG_MAP_NAME' বৈদেশিক..."
if [ ! -f "$CONFIG_FILE" ] || [ ! -f "$PROTO_FILE" ]; then
    echo "❌ 错误：配置文件 $CONFIG_FILE 或 $PROTO_FILE 未找到。"
    exit 1
fi
kubectl create configmap "$CONFIG_MAP_NAME" \
  --from-file=user.proto="$PROTO_FILE" \
  --from-file=kong.yaml="$CONFIG_FILE" \
  -n "$KONG_NAMESPACE" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ ConfigMap 已是最新状态。"

# 4. 应用固化的部署文件
echo "🚀 正在通过 kubectl apply 部署 Kong..."
if [ ! -f "$DEPLOYMENT_FILE" ]; then
    echo "❌ 错误：部署文件 '$DEPLOYMENT_FILE' 未找到。"
    exit 1
fi
kubectl apply -f "$DEPLOYMENT_FILE"
echo "✅ Kong 部署清单已应用。正在触发滚动更新以加载最新配置..."
kubectl rollout restart deployment/kong-kong -n "$KONG_NAMESPACE"

# 5. 等待 Deployment 完成其滚动更新
kubectl rollout status deployment/kong-kong -n "$KONG_NAMESPACE" --timeout=5m

echo "🎉 部署完成，Kong 已就绪！"