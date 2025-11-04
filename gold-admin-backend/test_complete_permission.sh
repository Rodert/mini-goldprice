#!/bin/bash

BASE_URL="http://localhost:8080/api"

echo "================================"
echo "完整权限测试"
echo "================================"
echo ""

# Admin登录
ADMIN_TOKEN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')

echo "测试1: 查看单店店长角色的菜单配置"
echo "===================================="
ROLE_INFO=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/roles/4")
echo "$ROLE_INFO" | jq '{id: .data.id, name: .data.name, menu_count: (.data.menu_ids | length)}'
echo "分配的菜单ID: $(echo "$ROLE_INFO" | jq -c '.data.menu_ids')"
echo ""

echo "测试2: shopmanager02登录并查看菜单"
echo "===================================="
SHOP_LOGIN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"shopmanager02","password":"123456"}')

echo "2.1 登录结果:"
echo "$SHOP_LOGIN" | jq '{username: .data.user_info.username, roles: .data.user_info.roles}'

echo ""
echo "2.2 获得的菜单:"
echo "$SHOP_LOGIN" | jq '.data.menu_list[] | {name, title, path}'

SHOP_TOKEN=$(echo "$SHOP_LOGIN" | jq -r '.data.token')
echo ""

echo "测试3: 修改用户信息时是否保留角色"
echo "===================================="
echo "3.1 修改shopmanager02的真实姓名..."
UPDATE_RESULT=$(curl -s -X PUT "$BASE_URL/users/6" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "real_name": "李店长-已更新",
    "role_ids": [4]
  }')

echo "更新结果: $(echo "$UPDATE_RESULT" | jq -r '.message')"

echo ""
echo "3.2 重新查询用户信息验证角色是否保留..."
USER_DETAIL=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/users/6")
echo "$USER_DETAIL" | jq '{username: .data.user.username, real_name: .data.user.real_name, role_ids: .data.role_ids}'
echo ""

echo "测试4: 创建新用户并分配多个角色"
echo "===================================="
MULTI_ROLE_USER=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "username": "test_multi_role",
    "password": "123456",
    "real_name": "多角色测试",
    "status": 1,
    "role_ids": [3, 5]
  }')

MULTI_USER_ID=$(echo "$MULTI_ROLE_USER" | jq -r '.data.id')
echo "创建用户ID: $MULTI_USER_ID"

echo ""
echo "4.1 验证多角色用户的角色..."
MULTI_USER_DETAIL=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/users/$MULTI_USER_ID")
echo "$MULTI_USER_DETAIL" | jq '{username: .data.user.username, role_ids: .data.role_ids}'

echo ""
echo "4.2 多角色用户登录查看菜单..."
MULTI_LOGIN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"test_multi_role","password":"123456"}')

echo "角色: $(echo "$MULTI_LOGIN" | jq -c '.data.user_info.roles')"
echo "菜单数量: $(echo "$MULTI_LOGIN" | jq '.data.menu_list | length')"
echo ""

echo "测试5: 修改角色时验证数据完整性"
echo "===================================="
echo "5.1 修改用户角色（从[4]改为[3,4,5]）..."
UPDATE_ROLE=$(curl -s -X PUT "$BASE_URL/users/6" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "real_name": "李店长-已更新",
    "phone": "13800138000",
    "role_ids": [3, 4, 5]
  }')

echo "更新结果: $(echo "$UPDATE_ROLE" | jq -r '.message')"

echo ""
echo "5.2 验证所有字段是否保留..."
VERIFY_RESULT=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/users/6")
echo "$VERIFY_RESULT" | jq '{
  username: .data.user.username,
  real_name: .data.user.real_name,
  phone: .data.user.phone,
  role_ids: .data.role_ids,
  status: .data.user.status
}'

echo ""
echo "5.3 用户重新登录验证菜单是否更新..."
NEW_LOGIN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"shopmanager02","password":"123456"}')

echo "新角色: $(echo "$NEW_LOGIN" | jq -c '.data.user_info.roles')"
echo "新菜单数量: $(echo "$NEW_LOGIN" | jq '.data.menu_list | length')"
echo ""

echo "================================"
echo "测试完成！"
echo "================================"
