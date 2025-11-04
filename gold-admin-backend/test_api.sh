#!/bin/bash

# API 测试脚本
BASE_URL="http://localhost:8080/api"

echo "================================"
echo "后台管理系统 API 测试"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local url=$3
    local data=$4
    local token=$5
    
    echo -e "${YELLOW}测试: $name${NC}"
    
    if [ -z "$token" ]; then
        if [ "$method" = "GET" ]; then
            response=$(curl -s -w "\n%{http_code}" "$url")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
        fi
    else
        if [ "$method" = "GET" ]; then
            response=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $token" "$url")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "$data" "$url")
        fi
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✓ 成功 (HTTP $http_code)${NC}"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo -e "${RED}✗ 失败 (HTTP $http_code)${NC}"
        echo "$body"
    fi
    echo ""
}

# 1. 健康检查
echo "=== 1. 健康检查 ==="
test_api "Ping" "GET" "http://localhost:8080/ping"

# 2. 登录
echo "=== 2. 登录测试 ==="
login_response=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}')

echo "$login_response" | jq '.'

# 提取 token
TOKEN=$(echo "$login_response" | jq -r '.data.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取 token！${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 登录成功，获取到 Token${NC}"
echo ""

# 3. 获取用户信息
echo "=== 3. 获取当前用户信息 ==="
test_api "获取用户信息" "GET" "$BASE_URL/user/info" "" "$TOKEN"

# 4. 首页看板统计
echo "=== 4. 首页看板统计 ==="
test_api "统计数据" "GET" "$BASE_URL/dashboard/stats" "" "$TOKEN"
test_api "最近动态" "GET" "$BASE_URL/dashboard/activities" "" "$TOKEN"
test_api "预约趋势" "GET" "$BASE_URL/dashboard/trend" "" "$TOKEN"

# 5. 用户管理
echo "=== 5. 用户管理 ==="
test_api "用户列表" "GET" "$BASE_URL/users?page=1&page_size=10" "" "$TOKEN"

# 6. 角色管理
echo "=== 6. 角色管理 ==="
test_api "角色列表" "GET" "$BASE_URL/roles?page=1&page_size=10" "" "$TOKEN"
test_api "所有角色" "GET" "$BASE_URL/roles/all" "" "$TOKEN"

# 7. 菜单管理
echo "=== 7. 菜单管理 ==="
test_api "菜单树" "GET" "$BASE_URL/menus" "" "$TOKEN"

# 8. 价格管理
echo "=== 8. 价格管理 ==="
test_api "价格列表" "GET" "$BASE_URL/prices" "" "$TOKEN"

# 9. 店铺管理
echo "=== 9. 店铺管理 ==="
test_api "店铺列表" "GET" "$BASE_URL/shops" "" "$TOKEN"
test_api "所有店铺" "GET" "$BASE_URL/shops/all" "" "$TOKEN"

# 10. 预约管理
echo "=== 10. 预约管理 ==="
test_api "预约列表" "GET" "$BASE_URL/appointments" "" "$TOKEN"
test_api "预约统计" "GET" "$BASE_URL/appointments/stats" "" "$TOKEN"

# 11. 创建测试数据
echo "=== 11. 创建测试数据 ==="

# 创建预约
appointment_data='{
    "metal_type": "gold_9999",
    "service_type": "store",
    "appointment_time": "2025-11-05T10:00:00Z",
    "name": "测试用户",
    "phone": "13800138000",
    "note": "测试预约"
}'
test_api "创建预约" "POST" "$BASE_URL/appointments" "$appointment_data" "$TOKEN"

# 12. 登出
echo "=== 12. 登出 ==="
test_api "登出" "POST" "$BASE_URL/logout" "" "$TOKEN"

echo "================================"
echo "测试完成！"
echo "================================"
