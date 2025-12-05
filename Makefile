.PHONY: build test run clean deps

# 默认目标
all: deps build test

# 构建项目
build:
	go build -o bin/ai-doc-library cmd/main.go

# 运行测试
test:
	go test -v -cover ./...

# 运行项目
run:
	go run cmd/main.go

# 清理构建产物
clean:
	rm -rf bin/

# 安装依赖
deps:
	go mod tidy
	go mod download

# 生成覆盖率报告
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 格式化代码
fmt:
	go fmt ./...

# 检查代码规范
lint:
	golangci-lint run

# 初始化数据库
init-db:
	psql -h localhost -U postgres -d ai_doc_library -f scripts/init.sql

# 运行前端开发服务器
dev-frontend:
	cd web && python3 -m http.server 8080

# 构建Docker镜像
docker-build:
	docker build -t ai-doc-library:latest .

# 运行Docker容器
docker-run:
	docker run -p 8080:8080 ai-doc-library:latest