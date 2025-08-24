#!/bin/bash

echo "🎯 项目测试脚本管理器"
echo "===================="

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
echo "📋 可用测试脚本："
echo "===================="
echo "1.  test-system.sh                    - 系统基础功能测试"
echo "2.  test-todo-creation.sh             - TODO创建功能测试"
echo "3.  test-calendar-date-fix.sh         - 日历日期修复验证"
echo "4.  test-calendar-form-layout.sh      - 日历表单布局验证"
echo "5.  test-calendar-time-modes.sh       - 日历时间段模式验证"
echo "6.  test-calendar-views.sh            - 日历多视图功能验证"
echo "7.  test-calendar-features.sh         - 日历功能综合测试"
echo "8.  test-calendar-fixes.sh            - 日历问题修复验证"
echo "9.  test-enhanced-editor.sh           - 富文本编辑器测试"
echo "10. test-enhanced-editor-icons.sh     - 编辑器图标修复验证"
echo "11. test-article-detail.sh            - 文章详情功能测试"
echo "12. test-article-detail-final.sh      - 文章详情最终验证"
echo "13. test-article-routing.sh           - 文章路由功能测试"
echo "14. test-article-final.sh             - 文章功能最终验证"
echo "15. test-article-complete.sh          - 文章功能完整测试"
echo "16. test-notifications.sh             - 消息通知功能测试"
echo ""

# 如果没有参数，显示选择菜单
if [ $# -eq 0 ]; then
    echo "请选择要运行的测试："
    echo "0.  运行所有测试"
    echo "输入数字选择特定测试，或输入 'all' 运行所有测试"
    echo ""
    read -p "请输入选择 (0-16 或 all): " choice
else
    choice=$1
fi

# 运行选择的测试
case $choice in
    0|all)
        echo "🚀 运行所有测试..."
        echo ""
        for script in test-*.sh; do
            if [ -f "$script" ] && [ "$script" != "run-all-tests.sh" ]; then
                echo "📋 运行: $script"
                echo "----------------------------------------"
                ./"$script"
                echo ""
                echo "✅ $script 完成"
                echo "========================================"
                echo ""
            fi
        done
        echo "🎉 所有测试完成！"
        ;;
    1)
        echo "📋 运行系统基础功能测试..."
        ./test-system.sh
        ;;
    2)
        echo "📋 运行TODO创建功能测试..."
        ./test-todo-creation.sh
        ;;
    3)
        echo "📋 运行日历日期修复验证..."
        ./test-calendar-date-fix.sh
        ;;
    4)
        echo "📋 运行日历表单布局验证..."
        ./test-calendar-form-layout.sh
        ;;
    5)
        echo "📋 运行日历时间段模式验证..."
        ./test-calendar-time-modes.sh
        ;;
    6)
        echo "📋 运行日历多视图功能验证..."
        ./test-calendar-views.sh
        ;;
    7)
        echo "📋 运行日历功能综合测试..."
        ./test-calendar-features.sh
        ;;
    8)
        echo "📋 运行日历问题修复验证..."
        ./test-calendar-fixes.sh
        ;;
    9)
        echo "📋 运行富文本编辑器测试..."
        ./test-enhanced-editor.sh
        ;;
    10)
        echo "📋 运行编辑器图标修复验证..."
        ./test-enhanced-editor-icons.sh
        ;;
    11)
        echo "📋 运行文章详情功能测试..."
        ./test-article-detail.sh
        ;;
    12)
        echo "📋 运行文章详情最终验证..."
        ./test-article-detail-final.sh
        ;;
    13)
        echo "📋 运行文章路由功能测试..."
        ./test-article-routing.sh
        ;;
    14)
        echo "📋 运行文章功能最终验证..."
        ./test-article-final.sh
        ;;
    15)
        echo "📋 运行文章功能完整测试..."
        ./test-article-complete.sh
        ;;
    16)
        echo "📋 运行消息通知功能测试..."
        ./test-notifications.sh
        ;;
    *)
        echo "❌ 无效选择: $choice"
        echo "请选择 0-16 或 all"
        exit 1
        ;;
esac

echo ""
echo "✨ 测试完成！"
