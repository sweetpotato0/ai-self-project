#!/bin/bash

echo "🔍 测试通知流程"

# 设置token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImlzcyI6Imdpbi13ZWItZnJhbWV3b3JrIiwiZXhwIjoxNzU2MDQ0Nzg0LCJuYmYiOjE3NTU5NTgzODQsImlhdCI6MTc1NTk1ODM4NH0.S16wdafD8aHShXcAoXLN4Of16NNqJrKL6f_sPSBM3zs"

echo "1. 检查后端健康状态..."
curl -s http://localhost:8080/api/v1/health | jq .

echo -e "\n2. 检查TODO任务列表..."
curl -s -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/todos" | jq '.data.todos[] | {id, title, due_date, status}' | head -20

echo -e "\n3. 检查当前通知列表..."
curl -s -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/notifications" | jq .

echo -e "\n4. 检查WebSocket端点..."
curl -s -I "http://localhost:8080/api/v1/ws" | head -5

echo -e "\n5. 等待1分钟让调度器运行..."
sleep 60

echo -e "\n6. 再次检查通知列表..."
curl -s -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/notifications" | jq .

echo -e "\n✅ 测试完成"
