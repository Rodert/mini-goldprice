#!/bin/bash

# 金价展示页面 - 快速启动脚本

echo "=========================================="
echo "  金价展示页面 - 独立版本"
echo "  兼容旧版浏览器 (Android 4.4+, Chrome 30+)"
echo "=========================================="
echo ""

# 检查Python是否安装
if command -v python3 &> /dev/null; then
    echo "✅ 检测到 Python 3"
    echo "🚀 启动服务器..."
    echo ""
    echo "访问地址: http://localhost:8000/index.html"
    echo "按 Ctrl+C 停止服务器"
    echo ""
    python3 -m http.server 8000
elif command -v python &> /dev/null; then
    echo "✅ 检测到 Python 2"
    echo "🚀 启动服务器..."
    echo ""
    echo "访问地址: http://localhost:8000/index.html"
    echo "按 Ctrl+C 停止服务器"
    echo ""
    python -m SimpleHTTPServer 8000
elif command -v php &> /dev/null; then
    echo "✅ 检测到 PHP"
    echo "🚀 启动服务器..."
    echo ""
    echo "访问地址: http://localhost:8000/index.html"
    echo "按 Ctrl+C 停止服务器"
    echo ""
    php -S localhost:8000
else
    echo "❌ 未检测到 Python 或 PHP"
    echo ""
    echo "请安装以下任一工具："
    echo "  - Python 3: python3 -m http.server 8000"
    echo "  - Python 2: python -m SimpleHTTPServer 8000"
    echo "  - PHP: php -S localhost:8000"
    echo "  - Node.js: npx http-server -p 8000"
    echo ""
    echo "或者直接在浏览器中打开 index.html 文件"
    exit 1
fi

