# HWHKit-Go Makefile

.PHONY: help build test clean run docker install lint fmt

# 默认目标
help:
	@echo "HWHKit-Go Development Commands:"
	@echo ""
	@echo "  install     - Install dependencies"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  test-v      - Run tests with verbose output"
	@echo "  bench       - Run benchmarks"
	@echo "  coverage    - Run tests with coverage"
	@echo "  lint        - Run linter"
	@echo "  fmt         - Format code"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker      - Build and run with Docker"
	@echo "  docker-down - Stop Docker containers"
	@echo ""

# 安装依赖
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# 构建应用
build:
	@echo "Building application..."
	go build -o bin/hwhkit-app ./examples/basic/main.go

# 运行应用
run:
	@echo "Running application..."
	go run ./examples/basic/main.go

# 运行测试
test:
	@echo "Running tests..."
	go test ./pkg/...

# 详细测试输出
test-v:
	@echo "Running tests with verbose output..."
	go test -v ./pkg/...

# 运行基准测试
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./pkg/...

# 测试覆盖率
coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./pkg/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 代码检查
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not found, install it with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

# 清理
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache

# Docker 构建和运行
docker:
	@echo "Building and running with Docker..."
	cd docker && docker-compose up --build

# 停止 Docker 容器
docker-down:
	@echo "Stopping Docker containers..."
	cd docker && docker-compose down

# 验证所有模块
verify:
	@echo "Verifying modules..."
	go mod verify

# 安装开发工具
dev-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 生成API文档
docs:
	@echo "Generating API documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g ./examples/basic/main.go; \
	else \
		echo "swag not found, install it with: make dev-tools"; \
	fi

# 完整检查
check: fmt lint test
	@echo "All checks passed!"

# 发布准备
release: clean check build
	@echo "Release preparation completed!"