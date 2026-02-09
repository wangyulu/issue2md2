.PHONY: test build install clean fmt vet lint

# 运行所有测试
test:
	go test ./...

# 构建二进制文件
build:
	go build -o bin/issue2md ./cmd/issue2md

# 安装到 $GOPATH/bin
install:
	go install ./cmd/issue2md

# 清理构建产物
clean:
	rm -rf bin/

# 格式化代码
fmt:
	go fmt ./...

# 运行 go vet
vet:
	go vet ./...

# 运行静态检查
lint:
	golangci-lint run

# 显示帮助
help:
	@echo "Available targets:"
	@echo "  test    - Run all tests"
	@echo "  build   - Build binary"
	@echo "  clean   - Clean build artifacts"
	@echo "  install - Install to GOPATH/bin"
	@echo "  fmt     - Format code"
	@echo "  vet     - Run go vet"
	@echo "  lint    - Run static check"
