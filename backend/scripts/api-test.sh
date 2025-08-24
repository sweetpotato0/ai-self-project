#!/bin/bash

# API测试脚本

set -e

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

echo "🧪 Testing API endpoints..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试健康检查
test_health() {
    echo -e "${YELLOW}Testing health endpoint...${NC}"
    response=$(curl -s -w "HTTP_STATUS:%{http_code}" "$BASE_URL/health")
    status_code=$(echo $response | grep -o "HTTP_STATUS:[0-9]*" | cut -d: -f2)

    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Health check passed${NC}"
    else
        echo -e "${RED}❌ Health check failed with status $status_code${NC}"
        exit 1
    fi
}

# 测试用户注册
test_register() {
    echo -e "${YELLOW}Testing user registration...${NC}"

    username="testuser_$(date +%s)"
    email="test_$(date +%s)@example.com"
    password="testpass123"

    response=$(curl -s -w "HTTP_STATUS:%{http_code}" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$username\",\"email\":\"$email\",\"password\":\"$password\"}" \
        "$BASE_URL/users/register")

    status_code=$(echo $response | grep -o "HTTP_STATUS:[0-9]*" | cut -d: -f2)

    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ User registration passed${NC}"
        USER_USERNAME=$username
        USER_PASSWORD=$password
    else
        echo -e "${RED}❌ User registration failed with status $status_code${NC}"
        echo $response
        exit 1
    fi
}

# 测试用户登录
test_login() {
    echo -e "${YELLOW}Testing user login...${NC}"

    response=$(curl -s -w "HTTP_STATUS:%{http_code}" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$USER_USERNAME\",\"password\":\"$USER_PASSWORD\"}" \
        "$BASE_URL/users/login")

    status_code=$(echo $response | grep -o "HTTP_STATUS:[0-9]*" | cut -d: -f2)
    body=$(echo $response | sed 's/HTTP_STATUS:[0-9]*$//')

    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ User login passed${NC}"
        TOKEN=$(echo $body | jq -r '.data.token')
        if [ "$TOKEN" = "null" ]; then
            echo -e "${RED}❌ No token received${NC}"
            exit 1
        fi
    else
        echo -e "${RED}❌ User login failed with status $status_code${NC}"
        echo $body
        exit 1
    fi
}

# 测试获取用户资料
test_profile() {
    echo -e "${YELLOW}Testing get profile...${NC}"

    response=$(curl -s -w "HTTP_STATUS:%{http_code}" \
        -H "Authorization: Bearer $TOKEN" \
        "$BASE_URL/users/profile")

    status_code=$(echo $response | grep -o "HTTP_STATUS:[0-9]*" | cut -d: -f2)

    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Get profile passed${NC}"
    else
        echo -e "${RED}❌ Get profile failed with status $status_code${NC}"
        echo $response
        exit 1
    fi
}

# 测试产品列表
test_products() {
    echo -e "${YELLOW}Testing get products...${NC}"

    response=$(curl -s -w "HTTP_STATUS:%{http_code}" \
        "$BASE_URL/products")

    status_code=$(echo $response | grep -o "HTTP_STATUS:[0-9]*" | cut -d: -f2)

    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Get products passed${NC}"
    else
        echo -e "${RED}❌ Get products failed with status $status_code${NC}"
        echo $response
        exit 1
    fi
}

# 检查服务是否运行
check_service() {
    if ! curl -s "$BASE_URL/health" > /dev/null; then
        echo -e "${RED}❌ Service is not running. Please start the server first.${NC}"
        echo "Run: make serve"
        exit 1
    fi
}

# 检查jq是否安装
check_dependencies() {
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}❌ jq is required but not installed.${NC}"
        echo "Install with: brew install jq (on macOS) or apt-get install jq (on Ubuntu)"
        exit 1
    fi
}

# 主函数
main() {
    echo "🚀 Starting API tests..."

    check_dependencies
    check_service

    test_health
    test_register
    test_login
    test_profile
    test_products

    echo -e "${GREEN}🎉 All tests passed!${NC}"
}

main
