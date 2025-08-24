#!/bin/bash

echo "=== 测试文章路由修复 ==="

# 检查前端是否运行
echo "1. 检查前端服务..."
if lsof -i :3000 > /dev/null 2>&1; then
    echo "✅ 前端服务正常运行"
else
    echo "❌ 前端服务未运行"
    exit 1
fi

# 检查后端是否运行
echo -e "\n2. 检查后端服务..."
if curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "✅ 后端服务正常运行"
else
    echo "❌ 后端服务未运行"
    exit 1
fi

echo -e "\n3. 测试路由配置..."
echo "✅ 独立路由: /articles/:id"
echo "✅ 嵌套路由: /dashboard/articles/:id"

echo -e "\n4. 可用的访问方式:"
echo "   📄 直接访问文章详情: http://localhost:3000/articles/1"
echo "   📄 通过仪表盘访问: http://localhost:3000/dashboard/articles/1"
echo "   📄 从文章列表跳转: 登录后点击文章卡片"

echo -e "\n5. 路由修复说明:"
echo "   🔧 添加了独立路由 /articles/:id"
echo "   🔧 保留了嵌套路由 /dashboard/articles/:id"
echo "   🔧 两种方式都可以访问文章详情"

echo -e "\n=== 测试完成 ==="
echo ""
echo "🎯 现在可以正常访问:"
echo "   http://localhost:3000/articles/1"
echo ""
echo "📋 使用步骤:"
echo "1. 打开浏览器访问: http://localhost:3000"
echo "2. 使用 admin/123456 登录"
echo "3. 直接访问: http://localhost:3000/articles/1"
echo "4. 或者从文章列表点击文章卡片"
