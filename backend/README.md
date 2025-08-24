# Gin Web Framework

一个基于Gin的Web后端项目框架，支持PostgreSQL、MySQL、SQLite数据库和Redis缓存。

## 功能特性

- 🚀 基于Gin框架的高性能Web服务
- 🗄️ 支持PostgreSQL、MySQL和SQLite数据库
- 🔄 Redis缓存支持
- 🔐 JWT认证中间件
- 📝 结构化日志记录
- 🛡️ CORS跨域支持
- 📊 统一响应格式
- 🔧 配置管理
- 🏗️ 模块化架构
- ✅ TODO清单管理模块
  - 任务CRUD操作
  - 优先级管理
  - 分类管理
  - 状态跟踪
  - 截止时间管理
  - 通知系统

## 项目结构

```
backend/
├── cmd/                  # 命令行工具
│   └── cli/             # 主CLI工具（基于Cobra）
│       ├── main.go
│       └── commands/    # 命令定义
├── config/              # 配置管理
│   └── config.go
├── internal/            # 内部包
│   ├── database/        # 数据库连接
│   ├── redis/          # Redis连接
│   ├── router/         # 路由管理
│   ├── middleware/     # 中间件
│   ├── handler/        # HTTP处理器
│   ├── service/        # 业务逻辑层
│   └── models/         # 数据模型
├── pkg/                # 公共包
│   ├── jwt/            # JWT工具
│   ├── auth/           # 认证工具
│   ├── logger/         # 日志管理
│   ├── response/       # 响应格式
│   └── validator/      # 验证工具
├── scripts/            # 自动化脚本
├── docs/               # 文档
├── main.go            # 主程序入口
├── go.mod             # Go模块文件
├── env.example        # 环境配置示例
├── Makefile           # 构建工具
└── README.md         # 项目说明
```

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd gin-web-framework
```

### 2. 使用开发脚本（推荐）

```bash
# 一键启动开发环境
./scripts/dev.sh
```

### 3. 手动设置

#### 安装依赖
```bash
make deps
```

#### 配置环境
```bash
make dev-setup  # 复制 env.example 到 .env
# 编辑 .env 文件配置数据库和Redis
```

#### 启动服务
```bash
# 使用Docker Compose（推荐）
make docker-run

# 或者本地运行
make migrate  # 运行数据库迁移
make serve    # 启动服务器
```

### 4. 验证服务

```bash
make health  # 健康检查
./scripts/api-test.sh  # API测试
```

## API接口

### 健康检查
- `GET /api/v1/health` - 健康检查

### 用户管理
- `POST /api/v1/users/register` - 用户注册
- `POST /api/v1/users/login` - 用户登录
- `GET /api/v1/users/profile` - 获取用户资料 (需要认证)
- `PUT /api/v1/users/profile` - 更新用户资料 (需要认证)

### 产品管理
- `GET /api/v1/products` - 获取产品列表
- `GET /api/v1/products/:id` - 获取单个产品
- `POST /api/v1/products` - 创建产品 (需要认证)
- `PUT /api/v1/products/:id` - 更新产品 (需要认证)
- `DELETE /api/v1/products/:id` - 删除产品 (需要认证)

### TODO管理
- `GET /api/v1/todos` - 获取TODO列表 (需要认证)
- `POST /api/v1/todos` - 创建TODO (需要认证)
- `PUT /api/v1/todos/:id` - 更新TODO (需要认证)
- `DELETE /api/v1/todos/:id` - 删除TODO (需要认证)

## 数据库支持

项目支持多种数据库，可以根据需要选择：

### PostgreSQL（推荐用于生产环境）
```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=password
DB_NAME=gin_web_framework
DB_SSLMODE=disable
```

### MySQL
```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_NAME=gin_web_framework
```

### SQLite（推荐用于开发环境）
```env
DB_DRIVER=sqlite
DB_PATH=./data/app.db
```

### 快速切换数据库

项目提供了预配置的环境文件，可以快速切换数据库：

```bash
# 使用PostgreSQL
make serve-postgres
make migrate-postgres

# 使用MySQL
make serve-mysql
make migrate-mysql

# 使用SQLite（推荐开发使用）
make serve-sqlite
make migrate-sqlite
```

## 可用命令

### Make命令
```bash
make help          # 显示所有可用命令
make build         # 构建所有工具
make serve         # 启动服务器
make migrate       # 运行数据库迁移
make user-create   # 创建用户（交互式）
make user-list     # 列出所有用户
make health        # 健康检查
make test          # 运行测试
make docker-run    # Docker Compose启动
```

### 开发脚本
```bash
./scripts/dev.sh        # 开发环境启动
./scripts/test.sh       # 运行测试套件
./scripts/api-test.sh   # API端点测试
```

### 命令行工具
```bash
# 编译后可用的工具
./bin/gin-cli       # 主CLI工具（基于Cobra框架）
```

### CLI工具使用
```bash
# 查看所有命令
./bin/gin-cli --help

# 启动服务器
./bin/gin-cli serve

# 数据库迁移
./bin/gin-cli migrate

# 用户管理
./bin/gin-cli user create
./bin/gin-cli user list

# 健康检查
./bin/gin-cli health

# 查看版本
./bin/gin-cli version
```

详细使用说明请参考：[CLI使用指南](docs/cli-usage.md)

## 开发指南

### 前端开发

前端项目位于 `../frontend-todo/` 目录，请参考前端的README文档。

```bash
# 进入前端目录
cd ../frontend-todo

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

前端访问地址：http://localhost:3000

### 添加新的模型

1. 在 `internal/models/` 目录下创建新的模型文件
2. 在 `internal/database/database.go` 的 `AutoMigrate()` 函数中添加模型

### 添加新的API接口

1. 在 `internal/service/` 目录下实现业务逻辑
2. 在 `internal/handler/` 目录下添加处理器函数
3. 在 `internal/router/router.go` 中添加路由

### 添加新的中间件

在 `internal/middleware/` 目录下创建中间件函数，并在路由中使用。

## 部署

### Docker部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
```

构建和运行：

```bash
docker build -t gin-web-framework .
docker run -p 8080:8080 gin-web-framework
```

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License
