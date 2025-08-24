# Gin CLI 工具使用指南

## 概述

Gin CLI 是一个基于 Cobra 框架的命令行工具，提供了完整的应用程序管理功能。

## 安装

```bash
# 构建CLI工具
make build

# 或者直接构建
go build -o bin/gin-cli cmd/cli/main.go
```

## 基本用法

### 查看帮助
```bash
./bin/gin-cli --help
```

### 查看版本
```bash
./bin/gin-cli version
```

## 命令详解

### 服务器管理

#### 启动服务器
```bash
# 使用默认配置启动
./bin/gin-cli serve

# 指定端口启动
./bin/gin-cli serve --port 3000

# 指定模式启动
./bin/gin-cli serve --mode release

# 组合使用
./bin/gin-cli serve --port 3000 --mode release
```

#### 健康检查
```bash
# 检查默认地址
./bin/gin-cli health

# 检查指定地址
./bin/gin-cli health --url http://localhost:3000

# 设置超时时间
./bin/gin-cli health --timeout 5
```

### 数据库管理

#### 运行迁移
```bash
# 正常迁移
./bin/gin-cli migrate

# 强制迁移（忽略错误）
./bin/gin-cli migrate --force
```

### 用户管理

#### 创建用户
```bash
./bin/gin-cli user create
```

#### 列出用户
```bash
./bin/gin-cli user list
```

## 全局标志

### 配置文件
```bash
# 使用指定配置文件
./bin/gin-cli --config /path/to/config.env serve

# 使用简写形式
./bin/gin-cli -c /path/to/config.env serve
```

### 详细输出
```bash
# 启用详细输出
./bin/gin-cli --verbose serve

# 使用简写形式
./bin/gin-cli -v serve
```

## 高级用法

### 自动补全

#### Bash 补全
```bash
# 生成 bash 补全脚本
./bin/gin-cli completion bash > ~/.local/share/bash-completion/completions/gin-cli

# 或者添加到 ~/.bashrc
echo 'source <(./bin/gin-cli completion bash)' >> ~/.bashrc
```

#### Zsh 补全
```bash
# 生成 zsh 补全脚本
./bin/gin-cli completion zsh > ~/.zsh/completion/_gin-cli

# 添加到 ~/.zshrc
echo 'fpath=(~/.zsh/completion $fpath)' >> ~/.zshrc
echo 'autoload -U compinit' >> ~/.zshrc
echo 'compinit' >> ~/.zshrc
```

#### Fish 补全
```bash
# 生成 fish 补全脚本
./bin/gin-cli completion fish > ~/.config/fish/completions/gin-cli.fish
```

### 别名设置

为了方便使用，可以设置别名：

```bash
# 添加到 ~/.bashrc 或 ~/.zshrc
alias gin='./bin/gin-cli'

# 然后就可以使用
gin serve
gin health
gin user list
```

## 配置文件

CLI 工具支持通过 `--config` 标志指定配置文件：

```bash
# 使用自定义配置文件
./bin/gin-cli --config production.env serve

# 配置文件格式（.env）
SERVER_PORT=8080
GIN_MODE=release
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=password
DB_NAME=gin_web_framework
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=dev_jwt_secret_key_here_must_be_at_least_32_chars_long
JWT_EXPIRE_HOURS=24
```

## 错误处理

CLI 工具提供了详细的错误信息：

```bash
# 数据库连接失败
./bin/gin-cli serve
Error: failed to initialize database: dial tcp localhost:5432: connect: connection refused

# 配置文件不存在
./bin/gin-cli --config nonexistent.env serve
Error: failed to load config: open nonexistent.env: no such file or directory

# 无效的参数
./bin/gin-cli serve --port invalid
Error: invalid port: strconv.Atoi: parsing "invalid": invalid syntax
```

## 最佳实践

### 1. 开发环境
```bash
# 启动开发服务器
./bin/gin-cli serve --mode debug

# 检查服务状态
./bin/gin-cli health
```

### 2. 生产环境
```bash
# 启动生产服务器
./bin/gin-cli serve --mode release --port 8080

# 定期健康检查
./bin/gin-cli health --url https://your-domain.com
```

### 3. 数据库管理
```bash
# 部署前运行迁移
./bin/gin-cli migrate

# 强制迁移（谨慎使用）
./bin/gin-cli migrate --force
```

### 4. 用户管理
```bash
# 创建管理员用户
./bin/gin-cli user create

# 查看所有用户
./bin/gin-cli user list
```

## 故障排除

### 常见问题

1. **权限问题**
   ```bash
   chmod +x bin/gin-cli
   ```

2. **依赖问题**
   ```bash
   go mod tidy
   go build -o bin/gin-cli cmd/cli/main.go
   ```

3. **配置文件问题**
   ```bash
   # 检查配置文件是否存在
   ls -la .env

   # 复制示例配置
   cp env.example .env
   ```

4. **数据库连接问题**
   ```bash
   # 检查数据库是否运行
   docker-compose ps

   # 重启数据库
   docker-compose restart postgres
   ```

## 扩展开发

### 添加新命令

1. 在 `cmd/cli/commands/` 目录下创建新文件
2. 实现命令逻辑
3. 在 `root.go` 中注册新命令

### 示例：添加配置命令

```go
// cmd/cli/commands/config.go
package commands

import (
    "github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "config",
        Short: "Manage configuration",
        RunE: func(cmd *cobra.Command, args []string) error {
            // 实现配置管理逻辑
            return nil
        },
    }
}
```

然后在 `root.go` 中添加：

```go
rootCmd.AddCommand(
    newServeCmd(),
    newMigrateCmd(),
    newUserCmd(),
    newHealthCmd(),
    newVersionCmd(),
    newConfigCmd(), // 添加新命令
)
```

## 总结

Gin CLI 工具提供了完整的应用程序管理功能，包括：

- ✅ 服务器启动和管理
- ✅ 数据库迁移
- ✅ 用户管理
- ✅ 健康检查
- ✅ 配置管理
- ✅ 自动补全支持
- ✅ 详细的帮助文档
- ✅ 错误处理和调试信息

通过合理使用这些命令，可以大大提高开发和运维效率。
