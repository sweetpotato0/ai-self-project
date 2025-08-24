# 容器死锁问题解决方案

## 🚨 问题描述

在 `backend/internal/container/container.go` 中存在死锁问题：

```go
func (c *Container) GetUserHandler() *handler.UserHandler {
	return c.getSingleton("user_handler", func() interface{} {
		return handler.NewUserHandler(c.GetUserService()) // 这里调用了GetUserService
	}).(*handler.UserHandler)
}

func (c *Container) GetUserService() service.UserServiceInterface {
	return c.getSingleton("user_service", func() interface{} {
		return service.NewUserService()
	}).(service.UserServiceInterface)
}
```

**死锁场景：**
1. 线程A调用 `GetUserHandler()`，获取写锁 `c.mu.Lock()`
2. 在线程A的 `getSingleton` 工厂函数中调用 `c.GetUserService()`
3. `GetUserService()` 又调用 `getSingleton`，尝试再次获取写锁 `c.mu.Lock()`
4. 由于同一个线程已经持有锁，造成死锁

## ✅ 解决方案

### 方案1：预初始化所有服务（推荐）

在容器创建时直接初始化所有服务，避免运行时死锁：

```go
func NewContainer(cfg config.ConfigInterface, db *gorm.DB, redisClient redis.RedisClient) ContainerInterface {
	container := &Container{
		services: make(map[string]interface{}),
		config:   cfg,
		db:       db,
		redis:    redisClient,
	}

	// 在创建时直接初始化所有服务，避免运行时死锁
	container.initializeAllServices()

	return container
}

func (c *Container) initializeAllServices() {
	// 创建所有服务实例
	userService := service.NewUserService()
	todoService := service.NewTodoService()
	// ... 其他服务

	// 创建所有处理器实例
	userHandler := handler.NewUserHandler(userService)
	todoHandler := handler.NewTodoHandler(todoService)
	// ... 其他处理器

	// 注册所有服务和处理器
	c.services["user_service"] = userService
	c.services["user_handler"] = userHandler
	// ... 其他注册
}

// 直接从已初始化的服务中获取
func (c *Container) GetUserService() *service.UserService {
	return c.services["user_service"].(*service.UserService)
}

func (c *Container) GetUserHandler() *handler.UserHandler {
	return c.services["user_handler"].(*handler.UserHandler)
}
```

### 方案2：双重检查锁定模式

使用双重检查锁定模式避免死锁：

```go
func (c *Container) getSingleton(name string, factory func() interface{}) interface{} {
	// 第一次检查（读锁）
	c.mu.RLock()
	if service, exists := c.services[name]; exists {
		c.mu.RUnlock()
		return service
	}
	c.mu.RUnlock()

	// 第二次检查（写锁）
	c.mu.Lock()
	defer c.mu.Unlock()

	// 双重检查，防止在获取写锁期间其他goroutine已经创建了服务
	if service, exists := c.services[name]; exists {
		return service
	}

	service := factory()
	c.services[name] = service

	return service
}
```

### 方案3：分离锁

为不同的服务类型使用不同的锁：

```go
type Container struct {
	mu           sync.RWMutex
	serviceMu    sync.RWMutex
	handlerMu    sync.RWMutex
	services     map[string]interface{}
	handlers     map[string]interface{}
	// ...
}
```

## 🎯 推荐方案

**推荐使用方案1（预初始化）**，原因如下：

1. **简单可靠** - 避免了复杂的锁机制
2. **性能更好** - 启动时一次性初始化，运行时无锁竞争
3. **易于理解** - 代码逻辑清晰，易于维护
4. **避免死锁** - 完全消除了死锁的可能性

## 🔧 实施步骤

1. **修改容器结构** - 移除 `getSingleton` 方法
2. **添加预初始化** - 在 `NewContainer` 中调用 `initializeAllServices`
3. **简化获取方法** - 直接从 `services` map 中获取
4. **更新接口** - 返回具体类型而不是接口（避免接口实现问题）

## 📊 性能对比

| 方案 | 启动时间 | 运行时性能 | 内存使用 | 复杂度 |
|------|----------|------------|----------|--------|
| 原始方案 | 快 | 慢（锁竞争） | 低 | 高（死锁风险） |
| 预初始化 | 稍慢 | 快（无锁） | 稍高 | 低 |
| 双重检查 | 快 | 中等 | 低 | 中等 |
| 分离锁 | 快 | 中等 | 低 | 高 |

## 🚀 总结

通过预初始化所有服务的方式，我们完全解决了死锁问题，同时提高了运行时性能。这是一个简单、可靠、高效的解决方案。
