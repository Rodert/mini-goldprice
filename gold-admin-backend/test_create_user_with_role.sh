#!/bin/bash

BASE_URL="http://localhost:8080/api"

echo "1. 登录获取 token..."
TOKEN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')

echo "Token: ${TOKEN:0:20}..."
echo ""

echo "2. 获取所有角色列表..."
ROLES=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/roles/all")
echo "$ROLES" | jq '.data[] | {id, name, code}'
echo ""

echo "3. 创建用户并分配'单店店长'角色（ID=4）..."
CREATE_RESULT=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "shopmanager02",
    "password": "123456",
    "real_name": "李店长",
    "status": 1,
    "role_ids": [4]
  }')

echo "$CREATE_RESULT" | jq
USER_ID=$(echo "$CREATE_RESULT" | jq -r '.data.id')
echo ""

echo "4. 查询新建用户的详细信息（验证角色是否分配成功）..."
curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/users/$USER_ID" | jq
echo ""

echo "5. 检查数据库中的user_roles表..."
sqlite3 ./data/gold_admin.db "SELECT * FROM user_roles WHERE user_id = $USER_ID;"
echo ""

echo "测试完成！"
