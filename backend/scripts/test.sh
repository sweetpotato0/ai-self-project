#!/bin/bash

# 测试脚本

set -e

echo "🧪 Running tests..."

# 设置测试环境变量
export GIN_MODE=test
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_PASSWORD=password
export DB_NAME=gin_web_framework_test
export DB_SSLMODE=disable

# 创建测试数据库
echo "📝 Setting up test database..."
docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS gin_web_framework_test;"
docker-compose exec postgres psql -U postgres -c "CREATE DATABASE gin_web_framework_test;"

# 运行测试
echo "🔍 Running unit tests..."
go test -v ./...

# 运行基准测试
echo "⚡ Running benchmark tests..."
go test -bench=. ./...

# 生成测试覆盖率报告
echo "📊 Generating coverage report..."
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "✅ Tests completed! Coverage report saved to coverage.html"
