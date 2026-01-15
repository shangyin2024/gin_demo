.PHONY: help

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
BLUE := \033[36m
RESET := \033[0m

help: ## 显示帮助信息
	@echo "$(BLUE)Gin Demo - 可用命令:$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# ============================================
# 开发环境
# ============================================

.PHONY: dev
dev: ## 启动开发环境（Docker + 数据库迁移）
	@echo "$(BLUE)Starting development environment...$(RESET)"
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 3
	@make migrate-up
	@echo "$(BLUE)Development environment is ready!$(RESET)"

.PHONY: dev-stop
dev-stop: ## 停止开发环境
	@echo "$(BLUE)Stopping development environment...$(RESET)"
	docker-compose down

.PHONY: dev-clean
dev-clean: ## 清理开发环境（包括数据卷）
	@echo "$(BLUE)Cleaning development environment...$(RESET)"
	docker-compose down -v
	@echo "$(BLUE)Development environment cleaned!$(RESET)"

# ============================================
# 代码生成
# ============================================

.PHONY: generate
generate: ## 生成所有代码（sqlc + wire）
	@echo "$(BLUE)Generating code...$(RESET)"
	@make sqlc
	@make wire
	@echo "$(BLUE)Code generation complete!$(RESET)"

.PHONY: sqlc
sqlc: ## 生成 sqlc 代码
	@echo "$(BLUE)Generating sqlc code...$(RESET)"
	sqlc generate
	@echo "$(BLUE)sqlc code generated!$(RESET)"

.PHONY: wire
wire: ## 生成 Wire 依赖注入代码
	@echo "$(BLUE)Generating Wire code...$(RESET)"
	wire ./internal/wire
	@echo "$(BLUE)Wire code generated!$(RESET)"

# ============================================
# 数据库
# ============================================

.PHONY: migrate-up
migrate-up: ## 执行数据库迁移
	@echo "$(BLUE)Running database migrations...$(RESET)"
	sql-migrate up
	@echo "$(BLUE)Migrations applied!$(RESET)"

.PHONY: migrate-down
migrate-down: ## 回滚数据库迁移
	@echo "$(BLUE)Rolling back database migrations...$(RESET)"
	sql-migrate down
	@echo "$(BLUE)Migrations rolled back!$(RESET)"

.PHONY: migrate-status
migrate-status: ## 查看迁移状态
	@sql-migrate status

.PHONY: db-reset
db-reset: ## 重置数据库（危险操作！）
	@echo "$(BLUE)Resetting database...$(RESET)"
	@make migrate-down
	@make migrate-up
	@echo "$(BLUE)Database reset complete!$(RESET)"

# ============================================
# 构建和运行
# ============================================

.PHONY: run
run: ## 运行应用
	@echo "$(BLUE)Starting application...$(RESET)"
	go run main.go

.PHONY: build
build: ## 编译应用
	@echo "$(BLUE)Building application...$(RESET)"
	go build -o bin/gin-demo .
	@echo "$(BLUE)Build complete! Binary: bin/gin-demo$(RESET)"

.PHONY: build-linux
build-linux: ## 编译 Linux 版本
	@echo "$(BLUE)Building for Linux...$(RESET)"
	GOOS=linux GOARCH=amd64 go build -o bin/gin-demo-linux .
	@echo "$(BLUE)Build complete! Binary: bin/gin-demo-linux$(RESET)"

# ============================================
# 测试
# ============================================

.PHONY: test
test: ## 运行测试
	@echo "$(BLUE)Running tests...$(RESET)"
	go test -v -race ./...

.PHONY: test-cover
test-cover: ## 运行测试并生成覆盖率报告
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(BLUE)Coverage report generated: coverage.html$(RESET)"

.PHONY: test-short
test-short: ## 运行短测试（跳过集成测试）
	@echo "$(BLUE)Running short tests...$(RESET)"
	go test -short -v ./...

# ============================================
# 代码质量
# ============================================

.PHONY: lint
lint: ## 代码检查
	@echo "$(BLUE)Running linter...$(RESET)"
	golangci-lint run

.PHONY: lint-fix
lint-fix: ## 自动修复代码问题
	@echo "$(BLUE)Running linter with auto-fix...$(RESET)"
	golangci-lint run --fix

.PHONY: fmt
fmt: ## 格式化代码
	@echo "$(BLUE)Formatting code...$(RESET)"
	go fmt ./...
	@echo "$(BLUE)Code formatted!$(RESET)"

.PHONY: vet
vet: ## 运行 go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	go vet ./...

# ============================================
# 依赖管理
# ============================================

.PHONY: deps
deps: ## 下载依赖
	@echo "$(BLUE)Downloading dependencies...$(RESET)"
	go mod download
	@echo "$(BLUE)Dependencies downloaded!$(RESET)"

.PHONY: tidy
tidy: ## 整理依赖
	@echo "$(BLUE)Tidying dependencies...$(RESET)"
	go mod tidy
	@echo "$(BLUE)Dependencies tidied!$(RESET)"

.PHONY: verify
verify: ## 验证依赖
	@echo "$(BLUE)Verifying dependencies...$(RESET)"
	go mod verify

# ============================================
# Docker
# ============================================

.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@echo "$(BLUE)Building Docker image...$(RESET)"
	docker build -t gin-demo:latest .
	@echo "$(BLUE)Docker image built: gin-demo:latest$(RESET)"

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	@echo "$(BLUE)Running Docker container...$(RESET)"
	docker run -p 8080:8080 --env-file .env gin-demo:latest

# ============================================
# 清理
# ============================================

.PHONY: clean
clean: ## 清理构建产物
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "$(BLUE)Clean complete!$(RESET)"

.PHONY: clean-all
clean-all: clean dev-clean ## 清理所有（包括 Docker）
	@echo "$(BLUE)All clean!$(RESET)"

# ============================================
# 工具安装
# ============================================

.PHONY: tools
tools: ## 安装开发工具
	@echo "$(BLUE)Installing development tools...$(RESET)"
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/rubenv/sql-migrate/...@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(BLUE)Tools installed!$(RESET)"

# ============================================
# 快捷命令组合
# ============================================

.PHONY: init
init: tools deps dev ## 初始化项目（安装工具 + 依赖 + 启动环境）
	@echo "$(BLUE)Project initialized! You can now run 'make run'$(RESET)"

.PHONY: check
check: fmt vet lint test ## 完整检查（格式化 + vet + lint + test）
	@echo "$(BLUE)All checks passed!$(RESET)"

.PHONY: ci
ci: generate check build ## CI 流程（生成代码 + 检查 + 构建）
	@echo "$(BLUE)CI pipeline complete!$(RESET)"

# API 文档
.PHONY: swagger
swagger: ## 生成 Swagger 文档
	@echo "==> 生成 Swagger 文档..."
	@which swag > /dev/null || (echo "❌ swag 未安装，运行: go install github.com/swaggo/swag/cmd/swag@latest" && exit 1)
	swag init
	swag fmt
	@echo "✅ Swagger 文档已生成: docs/swagger.json"

# ============================================
# 性能分析
# ============================================

.PHONY: bench
bench: ## 运行性能基准测试
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	go test -bench=. -benchmem -run=^$$ ./...
	@echo "$(BLUE)Benchmark complete!$(RESET)"

.PHONY: bench-cpu
bench-cpu: ## CPU 性能分析
	@echo "$(BLUE)Running CPU profiling...$(RESET)"
	go test -cpuprofile=cpu.prof -bench=. -run=^$$ ./...
	go tool pprof -http=:8081 cpu.prof
	@echo "$(BLUE)CPU profile: cpu.prof$(RESET)"

.PHONY: bench-mem
bench-mem: ## 内存性能分析
	@echo "$(BLUE)Running memory profiling...$(RESET)"
	go test -memprofile=mem.prof -bench=. -run=^$$ ./...
	go tool pprof -http=:8081 mem.prof
	@echo "$(BLUE)Memory profile: mem.prof$(RESET)"

.PHONY: pprof
pprof: ## 查看实时 pprof（需要应用运行在 debug 模式）
	@echo "$(BLUE)Opening pprof web interface...$(RESET)"
	@echo "Make sure the application is running in debug mode"
	@open http://localhost:8080/debug/pprof/ || xdg-open http://localhost:8080/debug/pprof/ || echo "Open http://localhost:8080/debug/pprof/ in browser"

# ============================================
# 代码分析
# ============================================

.PHONY: complexity
complexity: ## 分析代码复杂度
	@echo "$(BLUE)Analyzing code complexity...$(RESET)"
	@which gocyclo > /dev/null || go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	gocyclo -over 15 .
	@echo "$(BLUE)Complexity analysis complete!$(RESET)"

.PHONY: security
security: ## 安全扫描
	@echo "$(BLUE)Running security scan...$(RESET)"
	@which gosec > /dev/null || go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -fmt=json -out=security-report.json ./...
	@echo "$(BLUE)Security report: security-report.json$(RESET)"

.PHONY: deps-check
deps-check: ## 检查依赖更新
	@echo "$(BLUE)Checking for dependency updates...$(RESET)"
	@which go-mod-outdated > /dev/null || go install github.com/psampaz/go-mod-outdated@latest
	go list -u -m -json all | go-mod-outdated -update -direct
	@echo "$(BLUE)Dependency check complete!$(RESET)"

.PHONY: vuln
vuln: ## 检查漏洞
	@echo "$(BLUE)Scanning for vulnerabilities...$(RESET)"
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...
	@echo "$(BLUE)Vulnerability scan complete!$(RESET)"

# ============================================
# 统计信息
# ============================================

.PHONY: stats
stats: ## 显示项目统计信息
	@echo "$(BLUE)Project Statistics:$(RESET)"
	@echo ""
	@echo "📁 Go Files:"
	@find . -name "*.go" -not -path "./vendor/*" | wc -l | xargs echo "  "
	@echo ""
	@echo "📝 Lines of Code:"
	@find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | tail -1 | awk '{print "   " $$1}'
	@echo ""
	@echo "📦 Packages:"
	@go list ./... | wc -l | xargs echo "  "
	@echo ""
	@echo "🧪 Test Files:"
	@find . -name "*_test.go" -not -path "./vendor/*" | wc -l | xargs echo "  "
	@echo ""
	@echo "📊 Test Coverage:"
	@go test -cover -short ./... 2>&1 | grep "coverage:" | awk '{sum+=$$3; count++} END {if(count>0) printf "   %.1f%%\n", sum/count*100}'

.PHONY: todo
todo: ## 查找代码中的 TODO 和 FIXME
	@echo "$(BLUE)Finding TODOs and FIXMEs...$(RESET)"
	@grep -rn "TODO\|FIXME" --include="*.go" --exclude-dir=vendor . || echo "  No TODOs or FIXMEs found!"

# ============================================
# 数据库工具
# ============================================

.PHONY: db-console
db-console: ## 连接到数据库控制台
	@echo "$(BLUE)Connecting to database...$(RESET)"
	mysql -h localhost -P 3306 -u root -ppassword gin_demo

.PHONY: redis-console
redis-console: ## 连接到 Redis 控制台
	@echo "$(BLUE)Connecting to Redis...$(RESET)"
	redis-cli

# ============================================
# 环境配置
# ============================================

.PHONY: env
env: ## 创建 .env 文件（从 .env.example）
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(BLUE)Created .env file from .env.example$(RESET)"; \
		echo "$(BLUE)Please edit .env with your actual configuration$(RESET)"; \
	else \
		echo "$(BLUE).env file already exists$(RESET)"; \
	fi

