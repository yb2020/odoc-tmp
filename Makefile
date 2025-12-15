.PHONY: all build clean test test-unit test-integration test-benchmark test-coverage lint fmt vet

# 默认目标
all: lint test build

# 构建应用
build:
	@echo "Building application..."
	go build -o bin/server ./cmd/server

# 清理构建产物
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf coverage/

# 运行所有测试
test: test-unit test-integration test-benchmark

# 运行单元测试
test-unit:
	@echo "Running unit tests..."
	go test -v ./test/unit/simple_test.go

# 运行集成测试
test-integration:
	@echo "Running integration tests..."
	go test -v ./test/integration/...

# 运行基准测试
test-benchmark:
	@echo "Running benchmark tests..."
	go test -v -bench=. ./test/benchmark/...

# 生成测试覆盖率报告
test-coverage:
	@echo "Generating test coverage report..."
	mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report generated at coverage/coverage.html"

# 代码格式化
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 代码静态分析
vet:
	@echo "Vetting code..."
	go vet ./...

# 代码质量检查
lint: fmt vet
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping lint"; \
	fi
