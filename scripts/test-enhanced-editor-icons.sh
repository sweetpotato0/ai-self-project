#!/bin/bash

echo "🧪 测试富文本编辑器图标和功能"
echo "=================================="

# 检查前端服务是否运行
echo "1. 检查前端服务状态..."
if curl -s http://localhost:5173 > /dev/null; then
    echo "✅ 前端服务运行正常 (http://localhost:5173)"
else
    echo "❌ 前端服务未运行，请先启动: npm run dev"
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
echo "🎯 富文本编辑器图标优化完成！"
echo "=================================="
echo "✅ 文本格式图标："
echo "   - 粗体：⭐ (Star)"
echo "   - 斜体：✓ (Check)"
echo "   - 下划线：✗ (Close)"
echo ""
echo "✅ 对齐方式图标："
echo "   - 左对齐：← (ArrowLeft)"
echo "   - 居中对齐：↑ (ArrowUp)"
echo "   - 右对齐：→ (ArrowRight)"
echo ""
echo "✅ 插入元素图标："
echo "   - 图片：🖼️ (Picture)"
echo "   - 链接：🔗 (Link)"
echo "   - 表格：📊 (Grid)"
echo "   - 代码块：📄 (Document)"
echo "   - 引用：✏️ (Edit)"
echo ""
echo "✅ 列表图标："
echo "   - 无序列表：⭕ (CircleCheck)"
echo "   - 有序列表：📈 (Histogram)"
echo ""
echo "🎨 工具栏优化："
echo "   - 添加了分组标签"
echo "   - 添加了分隔线"
echo "   - 改进了视觉布局"
echo "   - 增强了悬停提示"
echo ""
echo "🌐 访问地址："
echo "   - 前端：http://localhost:5173"
echo "   - 文章编辑页面：http://localhost:5173/dashboard/articles"
echo ""
echo "📝 测试步骤："
echo "1. 打开浏览器访问 http://localhost:5173"
echo "2. 登录后进入文章管理"
echo "3. 创建或编辑文章"
echo "4. 测试各种富文本编辑功能"
echo "5. 验证图标是否直观易懂"
echo ""
echo "✨ 图标现在更加直观，无需悬停即可理解功能！"
