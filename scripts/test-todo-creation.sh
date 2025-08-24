#!/bin/bash

echo "=== 测试TODO创建功能 ==="

# 登录获取token
echo "1. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "获取token失败"
  exit 1
fi

echo "Token: ${TOKEN:0:50}..."

# 测试创建TODO（不包含日期）
echo -e "\n2. 测试创建TODO（不包含日期）..."
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "测试任务1",
    "description": "这是一个测试任务，不包含日期字段",
    "priority_id": 2,
    "category_id": 1,
    "estimated_hours": 4
  }'

echo -e "\n\n3. 测试创建TODO（包含日期）..."
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "测试任务2",
    "description": "这是一个包含日期的测试任务",
    "priority_id": 3,
    "category_id": 1,
    "start_date": "2024-08-23T10:00:00Z",
    "due_date": "2024-08-23T18:00:00Z",
    "estimated_hours": 8
  }'

echo -e "\n\n4. 获取TODO列表..."
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer $TOKEN"

echo -e "\n\n=== 测试完成 ==="
