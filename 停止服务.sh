#!/bin/bash

# 停止所有服务

echo "================================"
echo "  停止所有服务"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 停止 Go 后端
echo "停止后端服务..."
pkill -f "go run main.go"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 后端服务已停止${NC}"
else
    echo -e "${RED}✗ 未找到运行中的后端服务${NC}"
fi

# 停止 Node 前端
echo "停止前端服务..."
pkill -f "npm run dev"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 前端服务已停止${NC}"
else
    echo -e "${RED}✗ 未找到运行中的前端服务${NC}"
fi

# 停止可能残留的 node 进程
lsof -ti:9527 | xargs kill -9 2>/dev/null
lsof -ti:8080 | xargs kill -9 2>/dev/null

echo ""
echo "================================"
echo -e "${GREEN}所有服务已停止${NC}"
echo "================================"






