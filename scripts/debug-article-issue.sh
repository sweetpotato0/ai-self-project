#!/bin/bash

echo "=== 调试文章详情页面问题 ==="

# 检查后端
echo "1. 检查后端服务..."
if curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "✅ 后端服务正常"
else
    echo "❌ 后端服务异常"
    exit 1
fi

# 检查前端
echo -e "\n2. 检查前端服务..."
if lsof -i :3000 > /dev/null 2>&1; then
    echo "✅ 前端服务正常"
else
    echo "❌ 前端服务异常"
    exit 1
fi

# 测试登录
echo -e "\n3. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}')

echo "登录响应: $LOGIN_RESPONSE"

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ 获取token失败"
    exit 1
fi

echo "✅ Token获取成功: ${TOKEN:0:50}..."

# 测试文章详情API
echo -e "\n4. 测试文章详情API..."
DETAIL_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/articles/1 \
  -H "Authorization: Bearer $TOKEN")

echo "文章详情API响应: $DETAIL_RESPONSE"

if echo "$DETAIL_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 文章详情API正常"

    # 提取文章信息
    TITLE=$(echo "$DETAIL_RESPONSE" | grep -o '"title":"[^"]*"' | cut -d'"' -f4)
    CONTENT=$(echo "$DETAIL_RESPONSE" | grep -o '"content":"[^"]*"' | cut -d'"' -f4)

    echo "   文章标题: $TITLE"
    echo "   内容长度: ${#CONTENT} 字符"
    echo "   内容预览: ${CONTENT:0:100}..."
else
    echo "❌ 文章详情API异常"
    echo "错误响应: $DETAIL_RESPONSE"
fi

echo -e "\n5. 问题诊断..."
echo "可能的问题:"
echo "1. 用户未登录或token过期"
echo "2. 前端API调用失败"
echo "3. 路由配置问题"
echo "4. 组件渲染问题"

echo -e "\n6. 解决方案..."
echo "1. 确保用户已登录: http://localhost:3000"
echo "2. 检查浏览器控制台错误信息"
echo "3. 检查Network标签页的API请求"
echo "4. 确认文章ID存在: 当前测试ID为1"

echo -e "\n=== 调试完成 ==="
