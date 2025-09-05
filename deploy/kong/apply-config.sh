#!/bin/bash

# 脚本：部署或更新在 Kind 集群中运行的 DB-less Kong 网关。
#
# 这个脚本是幂等的，可以安全地重复运行。
# 它会确保 Kong 的配置与本地文件同步，并且服务处于运行状态。
#
# 用法：在 'deploy/kong' 目录下运行 ./apply-config.sh

set -e

# --- 配置 ---
KONG_NAMESPACE="kong"
VALUES_FILE="./values.yaml"
CONFIG_FILE="./kong.yaml"
KONG_DEPLOYMENT="kong-kong"
CONFIG_MAP_NAME="kong-declarative-config"
HELM_RELEASE_NAME="kong"

# --- 脚本主体 ---

echo "▶️ 开始 Kong 部署/更新流程..."

# 1. 检查前提条件：Kind 集群是否在运行
echo "🔎 正在检查 Kind 集群状态..."
if ! kind get clusters | grep -q "dev-cluster"; then
    echo "❌ 错误：未找到正在运行的 Kind 集群。"
    echo "请先启动您的 Kind 集群 (例如，使用 'kind create cluster')。"
    exit 1
fi
echo "✅ Kind 集群正在运行。"

# 2. 确保 Kong 的命名空间存在
echo "🔎 正在确保命名空间 '$KONG_NAMESPACE' 存在..."
kubectl create namespace "$KONG_NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# 3. 从 kong.yaml 创建或更新声明式配置的 ConfigMap
echo "🔄 正在从 '$CONFIG_FILE' 更新 ConfigMap '$CONFIG_MAP_NAME'..."
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ 错误：在当前目录下未找到 Kong 配置文件 '$CONFIG_FILE'。"
    exit 1
fi
# 从 kong.yaml 文件创建或更新 ConfigMap
kubectl create configmap "$CONFIG_MAP_NAME" \
  --from-file=kong.yaml="$CONFIG_FILE" \

  -n "$KONG_NAMESPACE" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ ConfigMap 已是最新状态。"

# 4. 使用 Helm 部署或更新 Kong
echo "🚀 正在通过 Helm 部署/更新 Kong..."
if [ ! -f "$VALUES_FILE" ]; then
    echo "❌ 错误：在当前目录下未找到 Helm 值文件 '$VALUES_FILE'。"
    exit 1
fi
helm upgrade --install "$HELM_RELEASE_NAME" kong/kong \
  --values "$VALUES_FILE" \
  -n "$KONG_NAMESPACE" \
  --create-namespace

echo "✅ Helm 操作完成。正在等待 Deployment 生效..."

# 5. 等待 Deployment 完成其滚动更新
kubectl rollout status deployment/"$KONG_DEPLOYMENT" -n "$KONG_NAMESPACE" --timeout=5m

echo "🎉 部署完成，Kong 已就绪！"
echo "💡 您现在可以修改 'kong.yaml' 并重新运行此脚本来应用更改。"