#!/bin/bash

# set -e: 如果任何命令失败，脚本将立即退出
set -e

# --- 核心修正 ---
# 获取脚本文件自身所在的目录的绝对路径
# 这样，无论你从哪个目录执行这个脚本，它总能找到正确的YAML文件
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# --- 介绍性输出 ---
echo "🚀 Starting database services deployment..."
echo "Executing from script directory: ${SCRIPT_DIR}"
echo "-------------------------------------------"

# --- 第1步: 创建持久化存储 (PVs and PVCs) ---
echo "STEP 1: Applying PersistentVolumes and Claims for MySQL..."
kubectl apply -f "${SCRIPT_DIR}/mysql-pv.yaml"

echo "STEP 2: Applying PersistentVolumes and Claims for Redis..."
kubectl apply -f "${SCRIPT_DIR}/redis-pv.yaml"

echo "✅ Persistent storage configured."
echo "-------------------------------------------"

# --- 第2步: 部署应用 (StatefulSets) ---
echo "STEP 3: Applying MySQL StatefulSet..."
kubectl apply -f "${SCRIPT_DIR}/mysql.yaml"

echo "STEP 4: Applying Redis StatefulSet..."
kubectl apply -f "${SCRIPT_DIR}/redis.yaml"

echo "✅ Database applications deployment initiated."
echo "-------------------------------------------"

# --- 第3步: 监控启动状态 ---
echo "⏳ Waiting for PersistentVolumeClaims to be bound..."
sleep 5
kubectl get pvc

echo ""
echo "👀 Monitoring Pod startup status (press Ctrl+C to exit)..."
kubectl get pods -w