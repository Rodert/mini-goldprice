#!/bin/bash

# 沪金汇管理系统 - 一键启动脚本

echo "================================"
echo "  沪金汇管理系统 - 启动脚本"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        echo -e "${YELLOW}警告: 端口 $port 已被占用${NC}"
        return 1
    else
        echo -e "${GREEN}✓ 端口 $port 可用${NC}"
        return 0
    fi
}

# 启动后端
start_backend() {
    echo ""
    echo "=== 启动后端服务 ==="
    
    cd gold-admin-backend
    
    # 检查依赖
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}✗ 找不到 go.mod 文件${NC}"
        exit 1
    fi
    
    echo "正在启动 Go 后端服务..."
    echo "端口: 8080"
    echo ""
    
    # 后台启动
    nohup go run main.go > ../logs/backend.log 2>&1 &
    BACKEND_PID=$!
    
    echo -e "${GREEN}✓ 后端服务已启动 (PID: $BACKEND_PID)${NC}"
    echo "日志文件: logs/backend.log"
    
    cd ..
    
    # 等待服务启动
    echo "等待后端服务启动..."
    sleep 5
    
    # 测试后端
    if curl -s http://localhost:8080/ping > /dev/null; then
        echo -e "${GREEN}✓ 后端服务运行正常${NC}"
    else
        echo -e "${RED}✗ 后端服务启动失败${NC}"
        echo "请查看日志: tail -f logs/backend.log"
    fi
}

# 启动前端
start_frontend() {
    echo ""
    echo "=== 启动前端服务 ==="
    
    cd gold-admin-frontend
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo "首次启动，正在安装依赖..."
        npm install
    fi
    
    echo "正在启动 Vue 前端服务..."
    echo "端口: 9527"
    echo ""
    
    # 后台启动
    nohup npm run dev > ../logs/frontend.log 2>&1 &
    FRONTEND_PID=$!
    
    echo -e "${GREEN}✓ 前端服务已启动 (PID: $FRONTEND_PID)${NC}"
    echo "日志文件: logs/frontend.log"
    
    cd ..
}

# 主函数
main() {
    # 创建日志目录
    mkdir -p logs
    
    # 检查端口
    echo "检查端口..."
    check_port 8080
    check_port 9527
    
    echo ""
    read -p "是否继续启动服务？(y/n) " -n 1 -r
    echo
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "已取消启动"
        exit 0
    fi
    
    # 启动后端
    start_backend
    
    # 启动前端
    start_frontend
    
    echo ""
    echo "================================"
    echo -e "${GREEN}✓ 所有服务已启动${NC}"
    echo "================================"
    echo ""
    echo "访问地址:"
    echo "  前端管理界面: http://localhost:9527"
    echo "  后端 API:     http://localhost:8080"
    echo ""
    echo "默认登录账号:"
    echo "  用户名: admin"
    echo "  密码:   admin123"
    echo ""
    echo "查看日志:"
    echo "  后端: tail -f logs/backend.log"
    echo "  前端: tail -f logs/frontend.log"
    echo ""
    echo "停止服务:"
    echo "  ./停止服务.sh"
    echo ""
    
    # 等待一会儿，让前端服务启动
    echo "等待前端服务启动..."
    sleep 10
    
    # 自动打开浏览器
    echo "正在打开浏览器..."
    if command -v open &> /dev/null; then
        open http://localhost:9527
    elif command -v xdg-open &> /dev/null; then
        xdg-open http://localhost:9527
    else
        echo "请手动打开浏览器访问: http://localhost:9527"
    fi
}

# 执行主函数
main


