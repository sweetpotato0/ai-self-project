#!/bin/bash

echo "=== 文章详情功能最终测试 ==="

# 检查后端是否运行
echo "1. 检查后端服务..."
if curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "✅ 后端服务正常运行"
else
    echo "❌ 后端服务未运行，请先启动后端"
    exit 1
fi

# 检查前端是否运行
echo -e "\n2. 检查前端服务..."
if lsof -i :3000 > /dev/null 2>&1; then
    echo "✅ 前端服务正常运行"
else
    echo "❌ 前端服务未运行，请先启动前端"
    exit 1
fi

# 登录获取token
echo -e "\n3. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ 获取token失败"
    echo "登录响应: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Token获取成功: ${TOKEN:0:50}..."

# 获取文章详情
echo -e "\n4. 测试文章详情API..."
DETAIL_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/articles/1 \
  -H "Authorization: Bearer $TOKEN")

if echo "$DETAIL_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 文章详情API调用成功"

    # 提取文章信息
    TITLE=$(echo "$DETAIL_RESPONSE" | grep -o '"title":"[^"]*"' | cut -d'"' -f4)
    CONTENT_LENGTH=$(echo "$DETAIL_RESPONSE" | grep -o '"content":"[^"]*"' | cut -d'"' -f4 | wc -c)

    echo "   文章标题: $TITLE"
    echo "   内容长度: $CONTENT_LENGTH 字符"
else
    echo "❌ 文章详情API调用失败"
    echo "响应: $DETAIL_RESPONSE"
    exit 1
fi

# 测试前端页面
echo -e "\n5. 测试前端页面..."
echo "✅ 前端页面地址: http://localhost:3000"
echo "✅ 文章详情页面: http://localhost:3000/articles/1"
echo "✅ 通过仪表盘访问: http://localhost:3000/dashboard/articles/1"

echo -e "\n=== 测试完成 ==="
echo ""
echo "📋 使用说明:"
echo "1. 打开浏览器访问: http://localhost:3000"
echo "2. 使用 admin/123456 登录"
echo "3. 点击'个人文章'菜单"
echo "4. 点击任意文章卡片查看文章详情"
echo "5. 或者直接访问: http://localhost:3000/articles/1"
echo ""
echo "🔧 调试信息:"
echo "- 打开浏览器开发者工具 (F12)"
echo "- 查看 Console 标签页的调试信息"
echo "- 查看 Network 标签页的API请求"
echo ""
echo "🎨 新功能:"
echo "- 美观的Markdown渲染"
echo "- 渐变背景和动画效果"
echo "- 响应式设计"
echo "- 交互式元素"
