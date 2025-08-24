#!/bin/bash

echo "=== 测试文章详情功能 ==="

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

# 获取文章列表
echo -e "\n2. 获取文章列表..."
ARTICLES_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/articles \
  -H "Authorization: Bearer $TOKEN")

echo "文章列表响应: $ARTICLES_RESPONSE"

# 提取第一个文章的ID
ARTICLE_ID=$(echo $ARTICLES_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$ARTICLE_ID" ]; then
  echo "没有找到文章，先创建一篇测试文章..."

  # 创建测试文章
  CREATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/articles \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "title": "测试文章详情功能",
      "content": "# 测试文章\n\n这是一个测试文章，用于验证文章详情功能。\n\n## 功能特点\n\n- 支持Markdown格式\n- 可以查看完整内容\n- 支持编辑和删除\n\n## 代码示例\n\n```javascript\nconsole.log(\"Hello, World!\");\n```",
      "summary": "测试文章详情功能的完整实现",
      "status": "published",
      "tags": ["测试", "功能验证"]
    }')

  echo "创建文章响应: $CREATE_RESPONSE"

  # 重新获取文章列表
  ARTICLES_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/articles \
    -H "Authorization: Bearer $TOKEN")

  ARTICLE_ID=$(echo $ARTICLES_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
fi

if [ -z "$ARTICLE_ID" ]; then
  echo "无法获取文章ID"
  exit 1
fi

echo "使用文章ID: $ARTICLE_ID"

# 获取文章详情
echo -e "\n3. 获取文章详情..."
DETAIL_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/articles/$ARTICLE_ID \
  -H "Authorization: Bearer $TOKEN")

echo "文章详情响应: $DETAIL_RESPONSE"

# 测试前端路由
echo -e "\n4. 测试前端路由..."
echo "前端文章详情页面地址: http://localhost:3001/articles/$ARTICLE_ID"
echo "请在前端浏览器中访问上述地址查看文章详情页面"

echo -e "\n=== 测试完成 ==="
echo "后端地址: http://localhost:8080"
echo "前端地址: http://localhost:3001"
echo "测试用户: admin / 123456"
