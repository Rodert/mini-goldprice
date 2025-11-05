.PHONY: run build clean test init-admin

# 运行开发服务器
run:
	go run main.go

# 安装依赖
deps:
	go mod tidy
	go mod download

# 编译
build:
	go build -ldflags="-w -s" -o gold-admin main.go

# 清理
clean:
	rm -f gold-admin
	rm -rf data/
	rm -rf logs/

# 测试
test:
	go test -v ./...

# 创建默认管理员账号（用户名: admin 密码: admin123）
init-admin:
	@echo "创建默认管理员账号..."
	@go run tools/create_admin.go

# 帮助
help:
	@echo "可用命令:"
	@echo "  make run          - 运行开发服务器"
	@echo "  make deps         - 安装依赖"
	@echo "  make build        - 编译项目"
	@echo "  make clean        - 清理生成文件"
	@echo "  make test         - 运行测试"
	@echo "  make init-admin   - 创建默认管理员"





