#!/bin/bash

# 用户和角色管理测试脚本

BASE_URL="http://localhost:8080/api"

echo "================================"
echo "用户和角色管理 - 功能演示"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 1. 登录获取 token
echo -e "${BLUE}=== 步骤1: 登录系统 ===${NC}"
login_response=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo "$login_response" | jq -r '.data.token')
echo -e "${GREEN}✓ 登录成功${NC}"
echo ""

# 2. 查看所有角色
echo -e "${BLUE}=== 步骤2: 查看现有角色 ===${NC}"
curl -s -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/roles/all" | jq '.data[] | {id, name, code}'
echo ""

# 3. 创建新角色
echo -e "${BLUE}=== 步骤3: 创建自定义角色 ===${NC}"
echo "创建角色：价格管理员"
new_role=$(curl -s -X POST "$BASE_URL/roles" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "name": "价格管理员",
      "code": "price_admin",
      "description": "只负责价格管理",
      "sort": 10,
      "status": 1,
      "menu_ids": [1, 10, 11]
    }')

role_id=$(echo "$new_role" | jq -r '.data.id')
echo -e "${GREEN}✓ 角色创建成功，ID: $role_id${NC}"
echo ""

# 4. 查看角色详情
echo -e "${BLUE}=== 步骤4: 查看角色详情 ===${NC}"
curl -s -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/roles/$role_id" | jq
echo ""

# 5. 创建新用户
echo -e "${BLUE}=== 步骤5: 创建新用户 ===${NC}"
echo "创建用户：shop_manager_001"
new_user=$(curl -s -X POST "$BASE_URL/users" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "username": "shop_manager_001",
      "password": "123456",
      "real_name": "店长001",
      "phone": "13800138001",
      "email": "shop001@example.com",
      "status": 1,
      "role_ids": ['$role_id']
    }')

user_id=$(echo "$new_user" | jq -r '.data.id')
echo -e "${GREEN}✓ 用户创建成功，ID: $user_id${NC}"
echo ""

# 6. 查看用户详情
echo -e "${BLUE}=== 步骤6: 查看用户详情（包含角色） ===${NC}"
curl -s -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/users/$user_id" | jq
echo ""

# 7. 修改用户角色
echo -e "${BLUE}=== 步骤7: 修改用户角色 ===${NC}"
echo "将用户角色改为：超级管理员"
curl -s -X PUT "$BASE_URL/users/$user_id" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "role_ids": [1]
    }' | jq '.message'
echo ""

# 8. 再次查看用户详情
echo -e "${BLUE}=== 步骤8: 验证角色变更 ===${NC}"
curl -s -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/users/$user_id" | jq '.data.role_ids'
echo ""

# 9. 测试新用户登录
echo -e "${BLUE}=== 步骤9: 测试新用户登录 ===${NC}"
new_user_login=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"shop_manager_001","password":"123456"}')

echo "用户信息："
echo "$new_user_login" | jq '.data.user_info'
echo ""
echo "菜单权限："
echo "$new_user_login" | jq '.data.menu_list[] | {id, title}'
echo ""

# 10. 获取用户列表
echo -e "${BLUE}=== 步骤10: 获取用户列表 ===${NC}"
curl -s -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/users?page=1&page_size=10" | jq '.data | {total, list: .list[] | {id, username, real_name}}'
echo ""

# 11. 修改用户密码
echo -e "${BLUE}=== 步骤11: 修改用户密码 ===${NC}"
curl -s -X PUT "$BASE_URL/users/$user_id/password" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "password": "new_password_123"
    }' | jq '.message'
echo ""

# 12. 使用新密码登录
echo -e "${BLUE}=== 步骤12: 使用新密码登录 ===${NC}"
new_password_login=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"shop_manager_001","password":"new_password_123"}')

if [ "$(echo "$new_password_login" | jq -r '.code')" == "200" ]; then
    echo -e "${GREEN}✓ 新密码登录成功${NC}"
else
    echo -e "${RED}✗ 新密码登录失败${NC}"
fi
echo ""

# 13. 清理测试数据（可选）
echo -e "${BLUE}=== 步骤13: 清理测试数据 ===${NC}"
read -p "是否删除测试数据？(y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # 删除测试用户
    curl -s -X DELETE "$BASE_URL/users/$user_id" \
        -H "Authorization: Bearer $TOKEN" | jq '.message'
    
    # 删除测试角色
    curl -s -X DELETE "$BASE_URL/roles/$role_id" \
        -H "Authorization: Bearer $TOKEN" | jq '.message'
    
    echo -e "${GREEN}✓ 测试数据已清理${NC}"
fi

echo ""
echo "================================"
echo "演示完成！"
echo "================================"
echo ""
echo "总结："
echo "✓ 创建了自定义角色（价格管理员）"
echo "✓ 创建了新用户（shop_manager_001）"
echo "✓ 为用户分配了角色"
echo "✓ 修改了用户角色"
echo "✓ 修改了用户密码"
echo "✓ 验证了新用户登录和权限"
echo ""





