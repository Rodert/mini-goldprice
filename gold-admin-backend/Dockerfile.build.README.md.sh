# 在 mac M1 平台打包 linux x86 amd64 的安装包 

# 构建并提取文件
docker build -f Dockerfile.build -t gold-builder .
docker create --name temp gold-builder
docker cp temp:/app/gold-admin-backend-linux-amd64 ./gold-admin-backend-linux-amd64
docker rm temp


# 验证
file gold-admin-backend-linux-amd64

