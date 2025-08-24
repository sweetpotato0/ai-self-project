#!/bin/bash

echo "=== 调试文章详情功能 ==="

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

# 获取文章详情
echo -e "\n2. 获取文章详情 (ID: 1)..."
DETAIL_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/articles/1 \
  -H "Authorization: Bearer $TOKEN")

echo "文章详情响应:"
echo "$DETAIL_RESPONSE" | jq '.'

# 检查文章内容
echo -e "\n3. 检查文章内容字段..."
CONTENT=$(echo "$DETAIL_RESPONSE" | jq -r '.data.content // "null"')
echo "文章内容长度: ${#CONTENT}"
echo "文章内容前100字符: ${CONTENT:0:100}"

# 检查前端API调用
echo -e "\n4. 检查前端API配置..."
echo "前端API基础URL: http://localhost:3001"
echo "文章详情页面URL: http://localhost:3001/articles/1"

# 检查浏览器控制台错误
echo -e "\n5. 调试建议:"
echo "- 打开浏览器开发者工具 (F12)"
echo "- 查看 Console 标签页是否有错误信息"
echo "- 查看 Network 标签页检查API请求"
echo "- 检查文章详情页面的Vue组件是否正确加载"

echo -e "\n=== 调试完成 ==="
