#!/bin/bash

echo "================================"
echo "用户角色选择问题 - 快速诊断"
echo "================================"
echo ""

BASE_URL="http://localhost:8080/api"

# 1. 登录
echo "1. 登录系统..."
TOKEN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ 登录失败！请检查后端服务是否启动"
  exit 1
fi
echo "✓ 登录成功"
echo ""

# 2. 检查角色列表
echo "2. 检查角色列表..."
roles_response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/roles/all")
role_count=$(echo "$roles_response" | jq '.data | length')

echo "角色数量: $role_count"

if [ "$role_count" -eq 0 ] || [ "$role_count" == "null" ]; then
  echo "❌ 角色数据为空！"
  echo ""
  echo "=== 问题原因 ==="
  echo "数据库中没有角色数据"
  echo ""
  echo "=== 解决方案 ==="
  echo "选择以下任一方案："
  echo ""
  echo "【推荐】方案1: 重新初始化数据库"
  echo "  pkill -f 'go run main.go'"
  echo "  rm -f ./data/gold_admin.db"
  echo "  go run main.go"
  echo ""
  echo "方案2: 手动创建角色（在另一个终端执行）"
  echo "  ./init_roles.sh"
  echo ""
  exit 1
else
  echo "✓ 角色数据正常"
fi
echo ""

# 3. 显示所有角色
echo "3. 可用角色列表..."
echo "$roles_response" | jq '.data[] | "  [\(.id)] \(.name) (\(.code))"' -r
echo ""

# 4. 测试创建用户（带角色）
echo "4. 测试创建用户并分配角色..."
timestamp=$(date +%s)
create_response=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"username\": \"test_role_${timestamp}\",
    \"password\": \"123456\",
    \"real_name\": \"测试用户-${timestamp}\",
    \"status\": 1,
    \"role_ids\": [3]
  }")

create_code=$(echo "$create_response" | jq -r '.code')

if [ "$create_code" == "200" ]; then
  echo "✓ 创建用户成功"
  user_id=$(echo "$create_response" | jq -r '.data.id')
  echo "  用户ID: $user_id"
  echo "  用户名: test_role_${timestamp}"
  echo ""
  
  # 5. 验证角色分配
  echo "5. 验证用户角色分配..."
  user_response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/users/$user_id")
  role_ids=$(echo "$user_response" | jq -r '.data.role_ids[]' 2>/dev/null)
  
  if [ -n "$role_ids" ]; then
    echo "✓ 角色分配成功"
    echo "  已分配角色ID: $role_ids"
  else
    echo "⚠️  用户创建成功，但角色分配可能有问题"
  fi
  
  echo ""
  echo "6. 清理测试数据..."
  delete_response=$(curl -s -X DELETE "$BASE_URL/users/$user_id" \
    -H "Authorization: Bearer $TOKEN")
  
  if [ "$(echo "$delete_response" | jq -r '.code')" == "200" ]; then
    echo "✓ 测试用户已删除"
  fi
else
  echo "❌ 创建用户失败"
  echo "$create_response" | jq
fi

echo ""
echo "================================"
echo "诊断完成！"
echo "================================"
echo ""

if [ "$role_count" -gt 0 ] && [ "$create_code" == "200" ]; then
  echo "✅ 用户角色选择功能正常！"
  echo ""
  echo "如果前端仍然无法选择角色，请检查："
  echo "  1. 浏览器控制台是否有错误"
  echo "  2. Network 标签中 /api/roles/all 请求是否成功"
  echo "  3. 尝试清除浏览器缓存并刷新页面"
else
  echo "❌ 发现问题，请按照上述解决方案操作"
fi






