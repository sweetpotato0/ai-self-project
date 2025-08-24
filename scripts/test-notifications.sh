#!/bin/bash

echo "=== 通知功能测试 ==="

# 1. 检查服务状态
echo "1. 检查服务状态..."
if curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "✅ 后端服务正常"
else
    echo "❌ 后端服务异常"
    exit 1
fi

# 2. 登录获取token
echo "2. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo "✅ 登录成功"
else
    echo "❌ 登录失败"
    exit 1
fi

# 3. 创建测试任务（即将到期）
echo "3. 创建即将到期的测试任务..."
DUE_DATE=$(date -v+1H -u +"%Y-%m-%dT%H:%M:%SZ")
TASK_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"测试即将到期任务\",
    \"description\": \"这是一个测试任务，将在1小时后到期\",
    \"priority_id\": 1,
    \"due_date\": \"$DUE_DATE\"
  }")

if echo "$TASK_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 创建即将到期任务成功"
else
    echo "❌ 创建即将到期任务失败"
    echo "响应: $TASK_RESPONSE"
fi

# 4. 创建测试任务（已逾期）
echo "4. 创建已逾期的测试任务..."
OVERDUE_DATE=$(date -v-1d -u +"%Y-%m-%dT%H:%M:%SZ")
OVERDUE_TASK_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"测试逾期任务\",
    \"description\": \"这是一个测试任务，已经逾期1天\",
    \"priority_id\": 1,
    \"due_date\": \"$OVERDUE_DATE\"
  }")

if echo "$OVERDUE_TASK_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 创建逾期任务成功"
else
    echo "❌ 创建逾期任务失败"
    echo "响应: $OVERDUE_TASK_RESPONSE"
fi

# 5. 手动触发通知检查
echo "5. 手动触发通知检查..."
# 这里可以通过API调用或直接运行通知检查逻辑
echo "ℹ️  通知检查将在后台定时运行"

# 6. 查看通知列表
echo "6. 查看通知列表..."
NOTIFICATIONS_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/notifications?limit=10" \
  -H "Authorization: Bearer $TOKEN")

if echo "$NOTIFICATIONS_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 获取通知列表成功"
    NOTIFICATION_COUNT=$(echo "$NOTIFICATIONS_RESPONSE" | grep -o '"data":\[[^]]*\]' | grep -o '\[.*\]' | jq length 2>/dev/null || echo "0")
    echo "   当前通知数量: $NOTIFICATION_COUNT"
else
    echo "❌ 获取通知列表失败"
    echo "响应: $NOTIFICATIONS_RESPONSE"
fi

# 7. 完成一个任务测试完成通知
echo "7. 完成一个任务测试完成通知..."
# 获取任务列表
TODOS_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer $TOKEN")

if echo "$TODOS_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 获取任务列表成功"
    # 这里可以解析任务ID并完成一个任务
    echo "ℹ️  可以通过前端界面完成任务来测试完成通知"
else
    echo "❌ 获取任务列表失败"
fi

# 8. 发布文章测试发布通知
echo "8. 发布文章测试发布通知..."
ARTICLE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "测试文章",
    "content": "这是一篇测试文章的内容",
    "summary": "测试文章摘要",
    "status": "published"
  }')

if echo "$ARTICLE_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 发布文章成功"
else
    echo "❌ 发布文章失败"
    echo "响应: $ARTICLE_RESPONSE"
fi

echo ""
echo "=== 通知功能测试完成 ==="
echo ""
echo "📋 通知类型说明:"
echo "1. 任务即将到期通知 - 任务到期前24小时"
echo "2. 任务逾期通知 - 任务已逾期"
echo "3. 任务完成通知 - 任务状态改为已完成"
echo "4. 文章发布通知 - 文章状态改为已发布"
echo "5. 系统通知 - 系统维护、更新等"
echo "6. 欢迎通知 - 新用户注册"
echo ""
echo "🔄 定时检查:"
echo "- 每小时自动检查任务到期情况"
echo "- 自动创建相应的通知"
echo ""
echo "📱 前端通知:"
echo "- 右上角通知图标显示未读数量"
echo "- 点击查看通知详情"
echo "- 支持标记已读/未读"
echo ""
echo "🎯 测试建议:"
echo "1. 创建不同到期时间的任务"
echo "2. 完成一些任务"
echo "3. 发布一些文章"
echo "4. 观察通知的产生和显示"
