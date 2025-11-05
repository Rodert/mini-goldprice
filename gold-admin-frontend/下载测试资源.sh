#!/bin/bash

# 金价展示页面 - 下载测试资源脚本
# 用于快速下载占位图片进行功能测试

echo "======================================"
echo "  金价展示页面 - 测试资源下载工具"
echo "======================================"
echo ""

# 创建资源目录
ASSETS_DIR="public/assets"
mkdir -p $ASSETS_DIR

echo "📁 资源目录: $ASSETS_DIR"
echo ""

# 检查 curl 是否可用
if ! command -v curl &> /dev/null; then
    echo "❌ 错误: 未找到 curl 命令"
    echo "请安装 curl 后重试"
    exit 1
fi

echo "⬇️  开始下载测试图片..."
echo ""

# 下载产品图片（使用占位图片服务）
download_image() {
    local num=$1
    local color=$2
    local file="$ASSETS_DIR/product${num}.jpg"
    
    echo "  下载 product${num}.jpg (${color})..."
    
    if [ -f "$file" ]; then
        echo "    ⚠️  文件已存在，跳过"
    else
        curl -s "https://via.placeholder.com/400x400/${color}/000000?text=Product+${num}" -o "$file"
        if [ $? -eq 0 ]; then
            echo "    ✅ 下载成功"
        else
            echo "    ❌ 下载失败"
        fi
    fi
}

# 下载5张占位图片
download_image 1 "FFD700"  # 金色
download_image 2 "C0C0C0"  # 银色
download_image 3 "CD7F32"  # 铜色
download_image 4 "FFD700"  # 金色
download_image 5 "C0C0C0"  # 银色

echo ""
echo "======================================"
echo "✅ 测试图片下载完成！"
echo "======================================"
echo ""
echo "📝 注意事项:"
echo "  1. 视频文件 jewelry-video.mp4 需要手动添加"
echo "  2. 可以从以下网站下载免费视频:"
echo "     - Pexels: https://www.pexels.com/zh-cn/search/videos/jewelry/"
echo "     - Pixabay: https://pixabay.com/zh/videos/search/jewelry/"
echo ""
echo "🚀 下一步:"
echo "  1. 启动后端: cd gold-admin-backend && go run main.go"
echo "  2. 启动前端: cd gold-admin-frontend && npm run serve"
echo "  3. 访问页面: http://localhost:8080/display"
echo ""
echo "📖 详细说明请查看: 金价展示页面使用指南.md"
echo ""

