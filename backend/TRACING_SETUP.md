# 🔍 Tracing中间件设置说明

## 📋 概述

本项目已经完整集成了OpenTelemetry分布式追踪系统，使用OTLP（OpenTelemetry Protocol）作为标准协议，替代了已弃用的Jaeger导出器。

## ✅ 已完成的集成

### 1. **中间件集成**
- ✅ **Tracing中间件**: 已添加到路由中间件链
- ✅ **Prometheus指标中间件**: 已添加到路由中间件链
- ✅ **OpenTelemetry初始化**: 在服务启动时自动初始化

### 2. **技术栈更新**
- ✅ **OTLP导出器**: 使用`go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp`
- ✅ **Prometheus指标**: 使用`go.opentelemetry.io/otel/exporters/prometheus`
- ✅ **标准化协议**: 遵循OpenTelemetry标准

## 🛠️ 配置说明

### 环境变量配置

```bash
# OpenTelemetry 配置
TELEMETRY_ENABLED=true
OTLP_ENDPOINT=http://localhost:4318/v1/traces
ENABLE_TRACING=true
ENABLE_METRICS=true
TELEMETRY_SERVICE_NAME=ai-self-project-backend
TELEMETRY_SERVICE_VERSION=1.0.0
TELEMETRY_ENVIRONMENT=development
```

### 配置系统集成

项目现在使用统一的配置系统管理OpenTelemetry设置：

```go
// 获取配置
cfg := config.Get()
telemetryCfg := cfg.GetTelemetry()

// 检查是否启用
if telemetryCfg.Enabled {
    // 使用 telemetryCfg.OTLPEndpoint
    // 使用 telemetryCfg.EnableTracing
    // 使用 telemetryCfg.EnableMetrics
}
```

### Kubernetes配置

```yaml
env:
- name: TELEMETRY_ENABLED
  value: "true"
- name: OTLP_ENDPOINT
  value: "http://jaeger-collector:4318/v1/traces"
- name: ENABLE_TRACING
  value: "true"
- name: ENABLE_METRICS
  value: "true"
- name: TELEMETRY_SERVICE_NAME
  value: "ai-self-project-backend"
- name: TELEMETRY_SERVICE_VERSION
  value: "1.0.0"
- name: TELEMETRY_ENVIRONMENT
  value: "production"
```

## 🔧 使用方法

### 1. **自动追踪**
所有HTTP请求都会自动被追踪，包括：
- 请求方法、路径、状态码
- 请求持续时间
- 错误信息
- 用户代理信息

### 2. **手动追踪**
在业务代码中使用追踪：

```go
import (
    "gin-web-framework/internal/middleware"
    "context"
)

// 数据库操作追踪
ctx, span := middleware.DatabaseTracing(ctx, "query", "SELECT * FROM users")
defer span.End()

// Redis操作追踪
ctx, span := middleware.RedisTracing(ctx, "get", "user:123")
defer span.End()

// 业务操作追踪
ctx, span := middleware.BusinessTracing(ctx, "create_user", map[string]interface{}{
    "user_id": 123,
    "email": "user@example.com",
})
defer span.End()
```

### 3. **获取追踪ID**
```go
traceID := middleware.GetTraceID(ctx)
spanID := middleware.GetSpanID(ctx)
```

## 🚀 部署配置

### 1. **本地开发**
```bash
# 启动OTLP Collector (可选)
docker run -p 4318:4318 otel/opentelemetry-collector:latest

# 启动应用
make run
```

### 2. **Kubernetes部署**
```bash
# 部署OTLP Collector
make otlp-setup

# 部署监控系统
make monitoring-setup

# 部署应用
make k8s-deploy
```

### 3. **Helm部署**
```bash
# 使用Helm部署
make helm-deploy
```

## 📊 监控和可视化

### 1. **Jaeger UI**
- 访问: http://localhost:16686
- 查看分布式追踪
- 分析请求链路

### 2. **Prometheus**
- 访问: http://localhost:9090
- 查看应用指标
- 配置告警规则

### 3. **Grafana**
- 访问: http://localhost:3000
- 可视化监控面板
- 自定义仪表板

## 🔍 追踪数据

### HTTP请求追踪
- **请求方法**: GET, POST, PUT, DELETE等
- **请求路径**: /api/v1/users, /api/v1/todos等
- **响应状态**: 200, 404, 500等
- **响应大小**: 字节数
- **请求持续时间**: 毫秒级精度

### 数据库追踪
- **操作类型**: query, transaction等
- **SQL语句**: 实际执行的SQL
- **执行时间**: 数据库操作耗时

### Redis追踪
- **操作类型**: get, set, del等
- **键名**: 操作的Redis键
- **执行时间**: Redis操作耗时

### 业务追踪
- **操作名称**: create_user, update_todo等
- **参数信息**: 业务操作的关键参数
- **执行时间**: 业务逻辑耗时

## 🛡️ 安全考虑

### 1. **敏感信息过滤**
- 自动过滤密码、令牌等敏感字段
- 只记录必要的业务参数
- 支持自定义过滤规则

### 2. **采样策略**
- 支持配置采样率
- 避免过多追踪数据
- 平衡性能和可观测性

### 3. **数据保留**
- 配置数据保留期限
- 自动清理过期数据
- 符合数据保护要求

## 🔧 故障排除

### 1. **追踪不工作**
```bash
# 检查OTLP端点
curl http://localhost:4318/v1/traces

# 检查环境变量
echo $OTLP_ENDPOINT
echo $ENABLE_TRACING
```

### 2. **性能问题**
```bash
# 检查采样率配置
# 检查批处理设置
# 检查内存使用情况
```

### 3. **数据丢失**
```bash
# 检查网络连接
# 检查OTLP Collector状态
# 检查Jaeger存储
```

## 📈 最佳实践

### 1. **追踪设计**
- 为关键业务操作添加追踪
- 使用有意义的操作名称
- 包含必要的上下文信息

### 2. **性能优化**
- 合理配置采样率
- 使用异步追踪
- 避免过度追踪

### 3. **监控告警**
- 设置响应时间告警
- 监控错误率
- 配置服务可用性告警

## 🎯 总结

Tracing中间件现在已经完全集成到项目中，提供了：

- ✅ **完整的分布式追踪**
- ✅ **标准化的OTLP协议**
- ✅ **自动化的追踪收集**
- ✅ **可视化的追踪分析**
- ✅ **生产就绪的配置**

项目现在具备了企业级的可观测性能力，可以轻松监控和分析分布式系统的性能和行为。
