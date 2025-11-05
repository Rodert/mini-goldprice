#!/bin/bash

# 重新初始化数据库脚本

echo "================================"
echo "  重新初始化数据库"
echo "================================"
echo ""

# 停止后端服务
echo "1. 正在停止后端服务..."
pkill -f "go run main.go"
pkill -f "gold-admin"
sleep 2
echo "✓ 后端服务已停止"
echo ""

# 备份旧数据库
echo "2. 备份旧数据库..."
if [ -f "./data/gold_admin.db" ]; then
    timestamp=$(date +%Y%m%d_%H%M%S)
    cp ./data/gold_admin.db ./data/gold_admin.db.backup_$timestamp
    echo "✓ 已备份到: ./data/gold_admin.db.backup_$timestamp"
else
    echo "! 未找到数据库文件，无需备份"
fi
echo ""

# 删除旧数据库
echo "3. 删除旧数据库..."
rm -f ./data/gold_admin.db
echo "✓ 旧数据库已删除"
echo ""

# 启动后端（会自动初始化新数据库）
echo "4. 启动后端服务（会自动初始化新数据库）..."
echo ""
go run main.go

