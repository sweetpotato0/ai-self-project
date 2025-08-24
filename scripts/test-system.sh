#!/bin/bash

echo "=== 系统测试脚本 ==="

# 测试后端健康检查
echo "1. 测试后端健康检查..."
BACKEND_HEALTH=$(curl -s http://localhost:8080/api/v1/health)
echo "后端响应: $BACKEND_HEALTH"

# 测试用户登录
echo -e "\n2. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "登录成功，Token: ${TOKEN:0:50}..."

# 测试获取文章列表
echo -e "\n3. 测试获取文章列表..."
ARTICLES_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/articles \
  -H "Authorization: Bearer $TOKEN")
echo "文章列表响应: $ARTICLES_RESPONSE"

# 测试获取文章统计
echo -e "\n4. 测试获取文章统计..."
STATS_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/articles/stats \
  -H "Authorization: Bearer $TOKEN")
echo "文章统计响应: $STATS_RESPONSE"

# 测试获取通知
echo -e "\n5. 测试获取通知..."
NOTIFICATIONS_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer $TOKEN")
echo "通知响应: $NOTIFICATIONS_RESPONSE"

echo -e "\n=== 测试完成 ==="
echo "后端地址: http://localhost:8080"
echo "前端地址: http://localhost:3002 (或检查其他端口)"
echo "测试用户: admin / 123456"
