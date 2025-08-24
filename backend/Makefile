# Makefile for Gin Web Framework Backend

# 变量定义
BINARY_NAME=server
BUILD_DIR=bin
MAIN_FILE=cmd/cli/main.go

# 默认目标
.PHONY: all
all: clean build

# 构建应用
.PHONY: build
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行应用
.PHONY: run
run:
	@echo "Running application..."
	@go run $(MAIN_FILE) serve

# 测试
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 性能测试
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. ./...

# 代码格式化
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
.PHONY: lint
lint:
	@echo "Running linter..."
	@go vet ./...

# 清理
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# 依赖管理
.PHONY: deps
deps:
	@echo "Updating dependencies..."
	@go mod tidy
	@go mod download

# 安装
.PHONY: install
install: deps build
	@echo "Installation completed"

# 开发模式
.PHONY: dev
dev:
	@echo "Starting development server..."
	@go run $(MAIN_FILE) serve --mode debug

# 生产模式
.PHONY: prod
prod:
	@echo "Starting production server..."
	@go run $(MAIN_FILE) serve --mode release

# 安全扫描
.PHONY: security-scan
security-scan:
	@echo "Running security scan..."
	@if command -v trivy > /dev/null; then \
		trivy fs .; \
	else \
		echo "Trivy not installed. Install with: brew install trivy"; \
	fi

# 构建Docker镜像
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t ai-self-project-backend:latest .

# 推送Docker镜像
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image..."
	@docker tag ai-self-project-backend:latest $(REGISTRY)/ai-self-project-backend:$(TAG)
	@docker push $(REGISTRY)/ai-self-project-backend:$(TAG)

# 部署到Kubernetes
.PHONY: k8s-deploy
k8s-deploy:
	@echo "Deploying to Kubernetes..."
	@kubectl apply -f k8s/

# 部署到Kubernetes (Helm)
.PHONY: helm-deploy
helm-deploy:
	@echo "Deploying with Helm..."
	@helm upgrade --install ai-self-project-backend ./helm \
		--namespace default \
		--create-namespace \
		--set image.tag=$(TAG)

# 部署到Istio
.PHONY: istio-deploy
istio-deploy:
	@echo "Deploying to Istio..."
	@kubectl apply -f k8s/istio.yaml

# 查看Kubernetes状态
.PHONY: k8s-status
k8s-status:
	@echo "Kubernetes deployment status:"
	@kubectl get pods -l app=ai-self-project-backend
	@kubectl get services -l app=ai-self-project-backend
	@kubectl get ingress -l app=ai-self-project-backend

# 查看日志
.PHONY: k8s-logs
k8s-logs:
	@echo "Fetching logs..."
	@kubectl logs -l app=ai-self-project-backend --tail=100 -f

# 端口转发
.PHONY: k8s-port-forward
k8s-port-forward:
	@echo "Port forwarding to service..."
	@kubectl port-forward svc/ai-self-project-backend-service 8080:80

# 清理Kubernetes资源
.PHONY: k8s-cleanup
k8s-cleanup:
	@echo "Cleaning up Kubernetes resources..."
	@kubectl delete -f k8s/ --ignore-not-found=true
	@helm uninstall ai-self-project-backend --ignore-not-found=true

# 运行集成测试
.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	@go test -tags=integration ./... -v

# 生成API文档
.PHONY: docs
docs:
	@echo "Generating API documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/cli/main.go; \
	else \
		echo "Swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 验证Kubernetes配置
.PHONY: k8s-validate
k8s-validate:
	@echo "Validating Kubernetes configurations..."
	@kubectl apply --dry-run=client -f k8s/
	@helm template ai-self-project-backend ./helm | kubectl apply --dry-run=client -f -

# 监控设置
.PHONY: monitoring-setup
monitoring-setup:
	@echo "Setting up monitoring..."
	@kubectl apply -f k8s/monitoring.yaml
	@kubectl apply -f k8s/otlp-collector.yaml

# OTLP设置
.PHONY: otlp-setup
otlp-setup:
	@echo "Setting up OTLP collector..."
	@kubectl apply -f k8s/otlp-collector.yaml

# 完整部署流程
.PHONY: deploy
deploy: docker-build k8s-deploy k8s-status
	@echo "Deployment completed!"

# 完整Helm部署流程
.PHONY: helm-deploy-full
helm-deploy-full: docker-build helm-deploy k8s-status
	@echo "Helm deployment completed!"

# 帮助
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  benchmark     - Run performance benchmarks"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Update dependencies"
	@echo "  install       - Install application"
	@echo "  dev           - Start development server"
	@echo "  prod          - Start production server"
	@echo "  security-scan - Run security scan"
	@echo "  docker-build  - Build Docker image"
	@echo "  k8s-deploy    - Deploy to Kubernetes"
	@echo "  helm-deploy   - Deploy with Helm"
	@echo "  istio-deploy  - Deploy to Istio"
	@echo "  k8s-status    - Show Kubernetes status"
	@echo "  k8s-logs      - Show Kubernetes logs"
	@echo "  k8s-cleanup   - Clean up Kubernetes resources"
	@echo "  deploy        - Full deployment pipeline"
	@echo "  monitoring-setup - Setup monitoring and OTLP"
	@echo "  otlp-setup    - Setup OTLP collector"
	@echo "  help          - Show this help message"
