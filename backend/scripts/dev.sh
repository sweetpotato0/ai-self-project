#!/bin/bash

# 开发环境启动脚本

set -e

echo "🚀 Starting development environment..."

# 检查是否有.env文件
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "⚠️  Please edit .env file with your database and Redis configuration"
    exit 1
fi

# 检查Go版本
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# 下载依赖
echo "📦 Downloading dependencies..."
go mod download

# 等待数据库启动
echo "⏳ Waiting for database to be ready..."
until docker-compose exec postgres pg_isready -U postgres > /dev/null 2>&1; do
    echo "Waiting for PostgreSQL..."
    sleep 2
done

# 运行数据库迁移
echo "🗄️  Running database migrations..."
make migrate

# 启动服务器
echo "🎯 Starting server..."
make serve
