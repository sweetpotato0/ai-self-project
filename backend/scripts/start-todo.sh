#!/bin/bash

# TODO清单管理系统启动脚本
# 包含后端和前端的一键启动

set -e

echo "🚀 启动TODO清单管理系统..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "❌ Node.js未安装，请先安装Node.js"
    exit 1
fi

# 检查npm是否安装
if ! command -v npm &> /dev/null; then
    echo "❌ npm未安装，请先安装npm"
    exit 1
fi

echo "📦 启动数据库和Redis..."
# 启动数据库和Redis
docker-compose up -d postgres redis

echo "⏳ 等待数据库启动..."
sleep 5

echo "🔧 配置环境..."
# 复制环境配置
if [ ! -f .env ]; then
    cp env.example .env
    echo "✅ 已创建.env配置文件"
fi

echo "🏗️ 构建后端..."
# 构建后端
make build

echo "🗄️ 运行数据库迁移..."
# 运行数据库迁移
./bin/gin-cli migrate

echo "📦 安装前端依赖..."
# 安装前端依赖
cd ../frontend-todo
if [ ! -d node_modules ]; then
    npm install
fi
cd ../backend

echo "🌐 启动后端服务..."
# 启动后端服务（后台运行）
./bin/gin-cli serve &
BACKEND_PID=$!

echo "⏳ 等待后端服务启动..."
sleep 3

echo "🎨 启动前端服务..."
# 启动前端服务（后台运行）
cd ../frontend-todo
npm run dev &
FRONTEND_PID=$!
cd ../backend

echo "✅ 服务启动完成！"
echo ""
echo "📊 服务地址："
echo "   后端API: http://localhost:8080"
echo "   前端页面: http://localhost:3000"
echo "   健康检查: http://localhost:8080/api/v1/health"
echo ""
echo "🔧 管理命令："
echo "   停止服务: ./scripts/stop-todo.sh"
echo "   查看日志: docker-compose logs -f"
echo "   CLI工具: ./bin/gin-cli --help"
echo ""

# 保存进程ID
echo $BACKEND_PID > .backend.pid
echo $FRONTEND_PID > .frontend.pid

# 等待用户中断
echo "按 Ctrl+C 停止服务"
trap 'echo "🛑 正在停止服务..."; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; rm -f .backend.pid .frontend.pid; echo "✅ 服务已停止"; exit 0' INT

wait
