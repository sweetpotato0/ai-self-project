#!/bin/bash

echo "🔍 完整流程测试"

# 设置token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImlzcyI6Imdpbi13ZWItZnJhbWV3b3JrIiwiZXhwIjoxNzU2MDQ0Nzg0LCJuYmYiOjE3NTU5NTgzODQsImlhdCI6MTc1NTk1ODM4NH0.S16wdafD8aHShXcAoXLN4Of16NNqJrKL6f_sPSBM3zs"

echo "1. 检查后端健康状态..."
curl -s http://localhost:8080/api/v1/health | jq .

echo -e "\n2. 检查TODO任务列表（确认有超时任务）..."
curl -s -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/todos" | jq '.data.todos[] | {id, title, due_date, status}' | head -10

echo -e "\n3. 检查当前通知列表..."
curl -s -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/notifications" | jq .

echo -e "\n4. 检查WebSocket端点..."
curl -s -I "http://localhost:8080/api/v1/ws" | head -5

echo -e "\n5. 检查前端是否运行..."
curl -s http://localhost:3001 | head -5

echo -e "\n6. 等待2分钟让调度器运行..."
echo "   在此期间，后端应该自动检查超时任务并创建通知"
sleep 120

echo -e "\n7. 再次检查通知列表..."
curl -s -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/notifications" | jq .

echo -e "\n8. 如果通知列表仍然为空，说明通知检查有问题"
echo -e "\n9. 建议手动触发通知检查或检查后端日志"

echo -e "\n✅ 测试完成"
echo -e "\n📝 下一步："
echo "   - 如果通知列表为空，检查后端日志"
echo "   - 测试WebSocket连接（使用浏览器打开 http://localhost:3001）"
echo "   - 检查前端控制台是否有WebSocket连接错误"
