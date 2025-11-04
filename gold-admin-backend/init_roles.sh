#!/bin/bash

echo "================================"
echo "角色数据初始化脚本"
echo "================================"
echo ""

BASE_URL="http://localhost:8080/api"

# 检查服务是否运行
echo "检查服务状态..."
if ! curl -s -f http://localhost:8080/health > /dev/null; then
  echo "❌ 后端服务未启动！"
  echo ""
  echo "请先启动服务："
  echo "  cd /Users/xuanxuanzi/home/s/javapub/mini-goldprice/gold-admin-backend"
  echo "  go run main.go"
  exit 1
fi
echo "✓ 服务正常运行"
echo ""

# 登录获取 token
echo "登录系统..."
TOKEN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ 登录失败！"
  exit 1
fi
echo "✓ 登录成功"
echo ""

# 检查是否已有角色
echo "检查现有角色..."
existing_count=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/roles/all" | jq '.data | length')

if [ "$existing_count" -gt 0 ]; then
  echo "⚠️  已存在 $existing_count 个角色"
  echo ""
  read -p "是否继续创建？可能会有重复 (y/N): " confirm
  if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
    echo "操作已取消"
    exit 0
  fi
fi
echo ""

# 创建角色
echo "开始创建角色..."
echo ""

# 定义角色数组
declare -a roles=(
  '{"name":"超级管理员","code":"super_admin","description":"拥有所有权限","sort":1,"status":1}'
  '{"name":"总部店长","code":"head_manager","description":"管理所有店铺","sort":2,"status":1}'
  '{"name":"单店店长","code":"shop_manager","description":"管理单个店铺","sort":3,"status":1}'
  '{"name":"店员","code":"shop_staff","description":"处理日常业务（只读）","sort":4,"status":1}'
  '{"name":"财务","code":"finance","description":"查看数据、导出报表","sort":5,"status":1}'
)

declare -a role_names=(
  "超级管理员"
  "总部店长"
  "单店店长"
  "店员"
  "财务"
)

success_count=0
fail_count=0

for i in "${!roles[@]}"; do
  role="${roles[$i]}"
  name="${role_names[$i]}"
  
  echo -n "创建角色: $name ... "
  
  response=$(curl -s -X POST "$BASE_URL/roles" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$role")
  
  code=$(echo "$response" | jq -r '.code')
  message=$(echo "$response" | jq -r '.message')
  
  if [ "$code" == "200" ]; then
    echo "✓ 成功"
    ((success_count++))
  else
    echo "✗ 失败 ($message)"
    ((fail_count++))
  fi
done

echo ""
echo "================================"
echo "初始化完成！"
echo "================================"
echo "成功: $success_count"
echo "失败: $fail_count"
echo ""

# 显示所有角色
echo "当前角色列表："
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/roles/all" | jq '.data[] | "  [\(.id)] \(.name) - \(.description)"' -r

echo ""
echo "✅ 角色初始化完成！现在可以创建用户并分配角色了。"


