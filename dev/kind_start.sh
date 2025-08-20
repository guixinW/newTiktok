#!/bin/bash

# 创建Kind集群脚本
# 包含代理配置和数据持久化

set -e

# 获取脚本所在目录并切换到该目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "🚀 开始创建Kind集群..."
echo "📍 工作目录: $SCRIPT_DIR"

# 创建数据目录
echo "📁 创建数据持久化目录..."
mkdir -p "$SCRIPT_DIR/data/mysql" "$SCRIPT_DIR/data/redis" "$SCRIPT_DIR/data/etcd"
echo "✅ 数据目录创建完成: $SCRIPT_DIR/data/mysql, $SCRIPT_DIR/data/redis, $SCRIPT_DIR/data/etcd"

# 检查kind配置文件是否存在
if [[ ! -f "$SCRIPT_DIR/kind-config.yaml" ]]; then
    echo "❌ 配置文件 kind-config.yaml 不存在于脚本目录中"
    exit 1
fi

echo "✅ 找到Kind集群配置文件: $SCRIPT_DIR/kind-config.yaml"

# 检查是否需要代理
read -p "🔗 是否需要配置代理? (y/n): " use_proxy

if [[ $use_proxy == "y" || $use_proxy == "Y" ]]; then
    # 询问代理地址
    read -p "🌐 请输入HTTP代理地址 (默认: http://host.docker.internal:7897): " proxy_url
    proxy_url=${proxy_url:-"http://host.docker.internal:7897"}

    echo "🔧 设置代理环境变量..."
    export HTTP_PROXY="$proxy_url"
    export HTTPS_PROXY="$proxy_url"
    export NO_PROXY="localhost,127.0.0.1,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,.svc,.cluster.local,kubernetes.default.svc"

    echo "✅ 代理配置完成:"
    echo "   HTTP_PROXY=$HTTP_PROXY"
    echo "   HTTPS_PROXY=$HTTPS_PROXY"
    echo "   NO_PROXY=$NO_PROXY"
else
    echo "⏭️  跳过代理配置"
fi

# 检查kind是否已安装
if ! command -v kind &> /dev/null; then
    echo "❌ Kind未安装，请先安装Kind"
    echo "安装命令: curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64 && chmod +x ./kind && sudo mv ./kind /usr/local/bin/kind"
    exit 1
fi

# 检查Docker是否运行
if ! docker info &> /dev/null; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 检查集群是否已存在
if kind get clusters | grep -q "dev-cluster"; then
    echo "⚠️  集群 'dev-cluster' 已存在"
    read -p "是否删除现有集群并重新创建? (y/n): " recreate
    if [[ $recreate == "y" || $recreate == "Y" ]]; then
        echo "🗑️  删除现有集群..."
        kind delete cluster --name dev-cluster
    else
        echo "❌ 取消创建，退出脚本"
        exit 1
    fi
fi

# 创建Kind集群
echo "🎯 创建Kind集群 (这可能需要几分钟时间)..."
if kind create cluster --config "$SCRIPT_DIR/kind-config.yaml"; then
    echo "✅ Kind集群创建成功!"
else
    echo "❌ Kind集群创建失败"
    exit 1
fi

# 验证集群状态
echo "🔍 验证集群状态..."
echo "集群信息:"
kubectl cluster-info --context kind-dev-cluster

echo ""
echo "节点状态:"
kubectl get nodes -o wide

echo ""
echo "🎉 Kind集群创建完成!"
echo ""
echo "📋 集群信息:"
echo "   集群名称: dev-cluster"
echo "   配置文件: $SCRIPT_DIR/kind-config.yaml"
echo "   数据目录: $SCRIPT_DIR/data/mysql, $SCRIPT_DIR/data/redis"
echo "   端口映射: 80:80, 443:443"
echo ""
echo "🔧 常用命令:"
echo "   查看集群: kubectl get nodes"
echo "   删除集群: kind delete cluster --name dev-cluster"
echo "   切换context: kubectl config use-context kind-dev-cluster"
echo ""
echo "✨ 现在可以开始部署应用了!"

# 部署数据库
echo "🚀 正在部署数据库..."
kubectl apply -f ../deploy/database/mysql-pv.yaml
kubectl apply -f ../deploy/database/mysql.yaml
kubectl apply -f ../deploy/database/redis-pv.yaml
kubectl apply -f ../deploy/database/redis.yaml
kubectl apply -f ../deploy/database/etcd-pv.yaml
kubectl apply -f ../deploy/database/etcd.yaml
echo "✅ 数据库部署完成!"

echo "🎉 所有部署完成!"
echo "✨ 集群 'dev-cluster' 已准备就绪!"