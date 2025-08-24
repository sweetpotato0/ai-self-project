# AI Self Project - Enterprise Web Framework

一个完整的、生产就绪的企业级Web应用框架，采用现代化的架构设计和最佳实践。

## 🏗️ 项目架构

```
ai-self-project/
├── backend/                 # Go后端服务
│   ├── cmd/                # 命令行工具
│   ├── config/             # 配置文件
│   ├── internal/           # 内部包
│   ├── pkg/                # 公共包
│   ├── scripts/            # 自动化脚本
│   └── docs/               # 文档
├── frontend/               # 前端应用
├── docs/                   # 项目文档
└── scripts/                # 项目级脚本
```

## 🚀 快速开始

### 后端服务

```bash
cd backend

# 安装依赖
go mod tidy

# 运行开发服务器
make dev

# 或者直接运行
go run cmd/cli/main.go serve
```

### 可用命令

```bash
# 查看所有可用命令
make help

# 构建应用
make build

# 运行测试
make test

# 性能测试
make benchmark

# 代码格式化
make fmt

# 代码检查
make lint
```

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **数据库**: PostgreSQL / MySQL / SQLite
- **缓存**: Redis
- **认证**: JWT
- **文档**: Swagger/OpenAPI
- **测试**: testify + httptest
- **容器**: Docker

### 架构特性
- ✅ 依赖注入容器
- ✅ 分层架构设计
- ✅ 中间件系统
- ✅ 安全防护
- ✅ 性能监控
- ✅ 错误处理
- ✅ 日志系统
- ✅ 配置管理

## 📚 文档

- [后端架构文档](backend/README.md)
- [API文档](backend/docs/)
- [CLI使用指南](backend/docs/cli-usage.md)
- [部署指南](docs/deployment.md)

## 🔧 开发工具

### 自动化脚本
- `make build` - 构建应用
- `make test` - 运行测试
- `make dev` - 开发模式
- `make prod` - 生产模式

### 代码质量
- 代码格式化: `go fmt`
- 代码检查: `go vet`
- 测试覆盖率: `go test -cover`
- 性能测试: `go test -bench`

## 🐳 Docker支持

```bash
# 构建镜像
docker build -t ai-self-project .

# 运行容器
docker-compose up -d
```

## 📊 项目状态

### 后端完成度
- ✅ 核心架构 (100%)
- ✅ 安全中间件 (100%)
- ✅ API设计 (100%)
- ✅ 测试框架 (100%)
- ✅ 文档系统 (100%)
- ✅ 部署配置 (100%)

### 前端状态
- 🔄 待更新

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 支持

如果你遇到任何问题或有建议，请：

1. 查看 [文档](docs/)
2. 搜索 [Issues](../../issues)
3. 创建新的 [Issue](../../issues/new)

---

**这是一个学习项目，展示了企业级Web应用开发的最佳实践。**
