#!/bin/bash

# Mac 本地 Docker 部署脚本
# 用于在 macOS 上快速启动 infinite-canvas 项目

set -e

echo "🚀 开始部署 infinite-canvas (Mac 本地环境)"
echo "=========================================="

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker Desktop"
    exit 1
fi

# 检查架构
ARCH=$(uname -m)
echo "✅ 检测到系统架构: $ARCH"

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "📝 创建 .env 配置文件..."
    cp .env.example .env
    echo "⚠️  请编辑 .env 文件修改默认密码和配置"
fi

# 创建数据目录
echo "📁 创建数据目录..."
mkdir -p data

# 停止旧容器
if docker ps -a --format '{{.Names}}' | grep -q "^infinite-canvas-local$"; then
    echo "🛑 停止并删除旧容器..."
    docker stop infinite-canvas-local 2>/dev/null || true
    docker rm infinite-canvas-local 2>/dev/null || true
fi

# 构建并启动
echo "🔨 开始构建镜像（首次构建可能需要几分钟）..."
docker compose -f docker-compose.local.yml up -d --build

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 5

# 检查容器状态
if docker ps | grep -q "infinite-canvas-local"; then
    echo ""
    echo "=========================================="
    echo "✅ 部署成功！"
    echo ""
    echo "📍 访问地址: http://localhost:3000"
    echo "👤 默认账号: admin"
    echo "🔑 默认密码: infinite-canvas"
    echo ""
    echo "📊 查看日志: docker logs -f infinite-canvas-local"
    echo "🛑 停止服务: docker compose -f docker-compose.local.yml down"
    echo "=========================================="
else
    echo ""
    echo "❌ 容器启动失败，请查看日志："
    echo "docker logs infinite-canvas-local"
    exit 1
fi
