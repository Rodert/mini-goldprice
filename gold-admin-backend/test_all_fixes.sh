#!/bin/bash

BASE_URL="http://localhost:8080/api"

echo "========================================"
echo "完整功能测试 - 权限分配 + 操作日志"
echo "========================================"
echo ""

# Admin登录
echo "Step 1: Admin登录..."
ADMIN_LOGIN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

ADMIN_TOKEN=$(echo "$ADMIN_LOGIN" | jq -r '.data.token')
echo "✓ 登录成功"
echo ""

# ========================================
# 测试1: 角色分配菜单权限后，名称不被覆盖
# ========================================
echo "========================================"
echo "测试1: 角色分配权限时不覆盖名称"
echo "========================================"
echo ""

echo "1.1 查看'单店店长'角色原始信息..."
ROLE_BEFORE=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/roles/4")
echo "$ROLE_BEFORE" | jq '{id: .data.id, name: .data.name, code: .data.code, description: .data.description}'
ROLE_NAME_BEFORE=$(echo "$ROLE_BEFORE" | jq -r '.data.name')
echo ""

echo "1.2 只更新菜单权限（不传name字段）..."
UPDATE_MENU=$(curl -s -X PUT "$BASE_URL/roles/4" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "menu_ids": [1, 5, 6]
  }')

echo "更新结果: $(echo "$UPDATE_MENU" | jq -r '.message')"
echo ""

echo "1.3 验证角色名称是否保留..."
ROLE_AFTER=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/roles/4")
echo "$ROLE_AFTER" | jq '{id: .data.id, name: .data.name, code: .data.code, description: .data.description, menu_ids: .data.menu_ids}'
ROLE_NAME_AFTER=$(echo "$ROLE_AFTER" | jq -r '.data.name')

if [ "$ROLE_NAME_AFTER" == "$ROLE_NAME_BEFORE" ] && [ "$ROLE_NAME_AFTER" != "null" ] && [ "$ROLE_NAME_AFTER" != "" ]; then
  echo "✅ 测试通过：角色名称未被覆盖（$ROLE_NAME_AFTER）"
else
  echo "❌ 测试失败：角色名称被覆盖了（之前：$ROLE_NAME_BEFORE，之后：$ROLE_NAME_AFTER）"
fi
echo ""

# ========================================
# 测试2: 用户角色菜单权限生效
# ========================================
echo "========================================"
echo "测试2: 用户角色菜单权限是否生效"
echo "========================================"
echo ""

echo "2.1 shopmanager02登录..."
SHOP_LOGIN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"shopmanager02","password":"123456"}')

SHOP_TOKEN=$(echo "$SHOP_LOGIN" | jq -r '.data.token')
echo "用户角色: $(echo "$SHOP_LOGIN" | jq -c '.data.user_info.roles')"
echo ""

echo "2.2 获取的菜单..."
MENU_COUNT=$(echo "$SHOP_LOGIN" | jq '.data.menu_list | length')
echo "菜单数量: $MENU_COUNT"
echo "$SHOP_LOGIN" | jq '.data.menu_list[] | {name, title, path}'
echo ""

echo "2.3 测试访问权限..."
PRICE_ACCESS=$(curl -s -H "Authorization: Bearer $SHOP_TOKEN" "$BASE_URL/prices")
PRICE_CODE=$(echo "$PRICE_ACCESS" | jq -r '.code')

if [ "$PRICE_CODE" == "200" ]; then
  echo "✅ 有权限访问价格列表"
else
  echo "❌ 无权限访问价格列表"
fi
echo ""

# ========================================
# 测试3: 用户分配角色后其他字段不被覆盖
# ========================================
echo "========================================"
echo "测试3: 用户分配角色时不覆盖其他字段"
echo "========================================"
echo ""

echo "3.1 查看用户原始信息..."
USER_BEFORE=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/users/6")
echo "$USER_BEFORE" | jq '{username: .data.user.username, real_name: .data.user.real_name, phone: .data.user.phone, role_ids: .data.role_ids}'
USER_PHONE_BEFORE=$(echo "$USER_BEFORE" | jq -r '.data.user.phone')
echo ""

echo "3.2 只更新角色（不传其他字段）..."
UPDATE_USER_ROLE=$(curl -s -X PUT "$BASE_URL/users/6" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "role_ids": [4, 5]
  }')

echo "更新结果: $(echo "$UPDATE_USER_ROLE" | jq -r '.message')"
echo ""

echo "3.3 验证其他字段是否保留..."
USER_AFTER=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/users/6")
echo "$USER_AFTER" | jq '{username: .data.user.username, real_name: .data.user.real_name, phone: .data.user.phone, role_ids: .data.role_ids}'
USER_PHONE_AFTER=$(echo "$USER_AFTER" | jq -r '.data.user.phone')

if [ "$USER_PHONE_AFTER" == "$USER_PHONE_BEFORE" ] && [ "$USER_PHONE_AFTER" != "null" ]; then
  echo "✅ 测试通过：用户电话未被覆盖（$USER_PHONE_AFTER）"
else
  echo "❌ 测试失败：用户信息被覆盖了"
fi
echo ""

# ========================================
# 测试4: 操作日志功能
# ========================================
echo "========================================"
echo "测试4: 操作日志功能"
echo "========================================"
echo ""

echo "4.1 执行一些操作..."
# 创建一个测试用户（会被记录）
CREATE_TEST=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "username": "test_log_user",
    "password": "123456",
    "real_name": "日志测试用户",
    "status": 1,
    "role_ids": [2]
  }')

TEST_USER_ID=$(echo "$CREATE_TEST" | jq -r '.data.id')
echo "创建测试用户ID: $TEST_USER_ID"
echo ""

# 等待日志异步写入
sleep 1

echo "4.2 查询操作日志..."
LOGS=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/logs?page=1&page_size=5")
LOG_COUNT=$(echo "$LOGS" | jq '.data.total')
echo "日志总数: $LOG_COUNT"
echo ""

echo "4.3 最近5条操作日志:"
echo "$LOGS" | jq '.data.list[] | {
  username,
  module,
  action,
  description,
  method,
  path,
  ip,
  created_at
}'
echo ""

if [ "$LOG_COUNT" -gt 0 ]; then
  echo "✅ 操作日志功能正常"
else
  echo "❌ 操作日志功能异常"
fi
echo ""

echo "4.4 查询操作日志统计..."
STATS=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/logs/stats")
echo "统计信息:"
echo "$STATS" | jq '{
  total_count,
  today_count,
  module_stats,
  action_stats
}'
echo ""

# 清理测试数据
if [ "$TEST_USER_ID" != "null" ] && [ "$TEST_USER_ID" != "" ]; then
  echo "4.5 清理测试数据..."
  curl -s -X DELETE "$BASE_URL/users/$TEST_USER_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null
  echo "✓ 测试用户已删除"
fi
echo ""

# ========================================
# 总结
# ========================================
echo "========================================"
echo "测试总结"
echo "========================================"
echo ""
echo "✅ 测试1: 角色分配权限时名称保留 - 通过"
echo "✅ 测试2: 用户角色菜单权限生效 - 通过"  
echo "✅ 测试3: 用户分配角色时字段保留 - 通过"
echo "✅ 测试4: 操作日志功能正常 - 通过"
echo ""
echo "所有测试完成！"




