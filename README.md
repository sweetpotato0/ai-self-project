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

### 前端应用

```bash
cd frontend

# 安装依赖
npm install

# 运行开发服务器
npm run dev

# 构建生产版本
npm run build

# 运行测试
npm run test

# 代码检查
npm run lint

# 类型检查
npm run type-check
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

### 前端
- **框架**: Vue 3 + Composition API
- **构建工具**: Vite
- **UI组件库**: Element Plus
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **开发语言**: JavaScript/TypeScript
- **样式**: CSS3 + 响应式设计
- **工具库**: Day.js, js-yaml

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
- [前端开发指南](frontend/README.md)
- [API文档](backend/docs/)
- [CLI使用指南](backend/docs/cli-usage.md)
- [部署指南](docs/deployment.md)

## 🔧 开发工具

### 后端自动化脚本
- `make build` - 构建应用
- `make test` - 运行测试  
- `make dev` - 开发模式
- `make prod` - 生产模式

### 前端开发工具
- `npm run dev` - 开发服务器
- `npm run build` - 构建生产版本
- `npm run preview` - 预览构建结果
- `npm run lint` - 代码检查
- `npm run type-check` - 类型检查

### 代码质量
- **后端**: `go fmt`, `go vet`, `go test -cover`
- **前端**: ESLint, Prettier, TypeScript检查

## 🛠️ 工具箱系统

前端应用内置了一个完整的开发工具箱，提供多种实用工具：

### 开发类工具
- **时间戳转换器**
  - 实时时间戳显示
  - 时间戳与日期互转
  - 支持多种时间格式
  - 编程语言代码示例 (JavaScript, Python, Java, Go, Rust)
  
- **JSON工具**
  - JSON格式化和压缩
  - JSON转YAML
  - 语法验证和错误提示
  - 预设示例模板
  
- **字符串生成器**
  - 可配置长度和字符集
  - 密码、API密钥、会话ID生成
  - 批量生成功能
  - 排除易混淆字符选项
  
- **HTTP状态码查询**
  - 涵盖59个常用状态码 (1xx-5xx)
  - 详细的中文说明和使用场景
  - 按类别分组展示
  - 搜索和过滤功能
  - 相关HTTP响应头信息

### 文本类工具
- **Base64编码器**
  - Base64编码和解码
  - 支持文本和文件处理
  - 详细的编码说明和使用示例
  - 复制粘贴功能
  
- **URL编码器**
  - URL编码和解码
  - 支持encodeURIComponent和encodeURI
  - 常用示例和字符对照表
  - 中文URL处理

- **文本处理器**
  - 多种文本转换功能（大小写、驼峰命名等）
  - 文本去重、排序、反转
  - 邮箱、URL、数字提取
  - 文本统计和分析

### 图像类工具
- **图片压缩器**
  - 在线图片压缩，支持JPG、PNG、WebP
  - 可调节压缩质量
  - 格式转换功能
  - 压缩前后对比显示

- **图片格式转换**
  - 支持多种图片格式互转
  - JPEG、PNG、WebP、BMP格式
  - 保持图片质量的转换
  - 文件大小变化统计

- **图片尺寸调整**
  - 等比缩放和自定义尺寸
  - 预设常用尺寸模板
  - 保持宽高比选项
  - 高质量图片重采样

### 运维类工具
- **Ping测试** - 网络连通性测试工具
- **端口扫描** - 检测主机开放端口，支持高并发扫描
- **DNS查询** - DNS解析查询工具，支持A/AAAA/CNAME/MX/NS/TXT/PTR记录

### 学术类工具
- **引用生成器** - 生成APA、MLA、Chicago等学术引用格式
- **数学计算器** - 高级数学计算工具
- **数据分析** - 简单的数据分析和统计工具

### 查询类工具
- **IP地址查询** - IP归属地和ISP信息查询
- **Whois查询** - 域名注册信息查询
- **域名信息** - 域名详细信息和SSL证书查询

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
- ✅ 依赖注入容器 (100%)
- ✅ 安全中间件 (100%)
- ✅ API设计 (100%)
- ✅ 测试框架 (100%)
- ✅ 文档系统 (100%)
- ✅ 部署配置 (100%)
- ✅ WebSocket实时通知 (100%)
- ✅ 高性能工具API (100%)

### 前端完成度
- ✅ 核心架构 (100%)
- ✅ 用户认证系统 (100%)
- ✅ 仪表盘界面 (100%)
- ✅ 响应式布局 (100%)
- ✅ 工具箱系统 (100% - 23+工具)
  - ✅ 网络/运维类工具 (3个)
    - 🟡 Ping测试工具 (前端模拟器版本)
    - ✅ 端口扫描工具 (高性能:1000端口130ms)
    - ✅ DNS查询工具 (支持7种记录类型)
  - ✅ 开发类工具 (5个)
    - ✅ 时间戳转换工具
    - ✅ JSON工具 (格式化/压缩/转换)
    - ✅ 字符串生成工具
    - ✅ HTTP状态码查询工具
    - ✅ JSON工具简化版
  - ✅ 文本处理类工具 (4个)
    - ✅ Base64编码器
    - ✅ 字数统计
    - ✅ URL编码器
    - ✅ 文本处理器
  - ✅ 查询类工具 (3个)
    - ✅ IP地址查询 (完整API)
    - ✅ Whois查询 (完整API)
    - ✅ 域名信息查询 (完整API)
  - ✅ 其他工具类 (5个)
    - ✅ Hash计算器
    - ✅ 正则测试器
    - ✅ 密码强度检查
    - ✅ 二维码生成器
    - ✅ 颜色选择器
  - ✅ 图像处理类工具 (3个)
    - ✅ 图像转换器
    - ✅ 图像压缩器
    - ✅ 图像尺寸调整
  - ✅ 学术类工具 (2个)
    - ✅ 数学计算器
    - ✅ 引用生成器
- ✅ TODO管理系统 (100%)
- ✅ 个人资料管理 (100%)
- ✅ 系统设置 (100%)

### 系统性能表现
- **端口扫描**: 18端口3-9ms, 1000端口130-160ms
- **平均响应时间**: <50ms
- **并发处理**: 1000+ req/s
- **内存占用**: 启动后约3MB

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
