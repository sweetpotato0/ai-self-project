#!/bin/bash

# 停止TODO清单管理系统

echo "🛑 正在停止TODO清单管理系统..."

# 停止后端服务
if [ -f .backend.pid ]; then
    BACKEND_PID=$(cat .backend.pid)
    if kill -0 $BACKEND_PID 2>/dev/null; then
        kill $BACKEND_PID
        echo "✅ 后端服务已停止"
    else
        echo "⚠️ 后端服务未运行"
    fi
    rm -f .backend.pid
fi

# 停止前端服务
if [ -f .frontend.pid ]; then
    FRONTEND_PID=$(cat .frontend.pid)
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        kill $FRONTEND_PID
        echo "✅ 前端服务已停止"
    else
        echo "⚠️ 前端服务未运行"
    fi
    rm -f .frontend.pid
fi

# 停止Docker服务
echo "🐳 停止数据库和Redis..."
docker-compose down

echo "✅ 所有服务已停止"
