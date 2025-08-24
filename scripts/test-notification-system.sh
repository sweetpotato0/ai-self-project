#!/bin/bash

echo "🔔 消息通知系统测试"
echo "=================="

# 检查前端服务是否运行
echo "1. 检查前端服务状态..."
if curl -s http://localhost:5173 > /dev/null; then
    echo "✅ 前端服务运行正常 (http://localhost:5173)"
else
    echo "❌ 前端服务未运行，请先启动: cd frontend-todo && npm run dev"
    exit 1
fi

# 检查后端服务是否运行
echo "2. 检查后端服务状态..."
if curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "✅ 后端服务运行正常 (http://localhost:8080)"
else
    echo "❌ 后端服务未运行，请先启动: cd backend && go run main.go"
    exit 1
fi

echo ""
echo "🔍 检查消息通知API..."
echo "===================="

# 测试获取通知列表API
echo "3. 测试获取通知列表API..."
NOTIFICATION_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/notifications" \
  -H "Authorization: Bearer $(cat ~/.auth_token 2>/dev/null || echo 'test_token')" \
  -H "Content-Type: application/json")

if echo "$NOTIFICATION_RESPONSE" | grep -q "notifications"; then
    echo "✅ 通知列表API正常"
    echo "📋 当前通知数量: $(echo "$NOTIFICATION_RESPONSE" | jq '.notifications | length' 2>/dev/null || echo '未知')"
else
    echo "❌ 通知列表API异常: $NOTIFICATION_RESPONSE"
fi

# 测试获取未读通知数量API
echo "4. 测试获取未读通知数量API..."
UNREAD_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/notifications/unread-count" \
  -H "Authorization: Bearer $(cat ~/.auth_token 2>/dev/null || echo 'test_token')" \
  -H "Content-Type: application/json")

if echo "$UNREAD_RESPONSE" | grep -q "count"; then
    echo "✅ 未读通知数量API正常"
    echo "📊 未读通知数量: $(echo "$UNREAD_RESPONSE" | jq '.count' 2>/dev/null || echo '未知')"
else
    echo "❌ 未读通知数量API异常: $UNREAD_RESPONSE"
fi

echo ""
echo "🔧 消息通知系统问题诊断："
echo "========================"
echo "❓ 可能的问题："
echo "1. 定时任务检查频率太低（当前每小时检查一次）"
echo "2. 没有实时WebSocket通知"
echo "3. 前端没有正确获取通知数据"
echo "4. 用户认证问题"
echo "5. 数据库中没有超时任务"
echo ""

echo "🔧 解决方案："
echo "============"
echo "1. 提高检查频率：修改 scheduler.go 中的检查间隔"
echo "2. 添加实时通知：启动WebSocket服务"
echo "3. 手动触发检查：创建测试任务并设置过期时间"
echo "4. 检查前端通知组件：确保正确显示通知"
echo ""

echo "📝 测试步骤："
echo "============"
echo "1. 登录系统"
echo "2. 创建一个任务，设置截止时间为过去的时间"
echo "3. 等待定时任务检查（或手动触发）"
echo "4. 查看右上角通知图标是否有红点"
echo "5. 点击通知图标查看通知列表"
echo ""

echo "🔧 手动测试命令："
echo "================"
echo "# 创建测试任务（已过期）"
echo "curl -X POST http://localhost:8080/api/v1/todos \\"
echo "  -H 'Authorization: Bearer YOUR_TOKEN' \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"title\":\"测试超时任务\",\"description\":\"这是一个测试任务\",\"due_date\":\"2024-01-01T00:00:00Z\"}'"
echo ""
echo "# 手动触发通知检查"
echo "curl -X POST http://localhost:8080/api/v1/admin/trigger-notifications \\"
echo "  -H 'Authorization: Bearer YOUR_TOKEN'"
echo ""

echo "🌐 访问地址："
echo "============"
echo "前端：http://localhost:5173"
echo "后端API：http://localhost:8080/api/v1/notifications"
echo ""

echo "✨ 建议："
echo "========"
echo "1. 检查后端日志，查看定时任务是否正常运行"
echo "2. 检查数据库中是否有通知记录"
echo "3. 确认前端通知组件是否正确实现"
echo "4. 考虑添加实时WebSocket通知功能"
