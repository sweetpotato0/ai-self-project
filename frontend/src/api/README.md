# 🌐 全局API配置目录

这个目录用于管理**全局共享的API配置**，不包含具体业务API接口。

## 📋 目录说明

### 🎯 设计理念
- **职责分离**: 全局配置与业务接口分离
- **统一管理**: axios实例、拦截器等基础配置集中管理
- **功能模块化**: 具体业务API在 `features/*/api/` 中管理

### 📁 当前文件

```
api/
├── README.md          # 本说明文件
├── index.js          # axios实例和拦截器配置 ⭐ 核心配置文件
└── category.js       # 通用分类API (跨模块使用)
```

### 🔧 核心文件说明

**`index.js` - axios核心配置**
- ✅ axios实例创建
- ✅ 请求/响应拦截器
- ✅ 错误处理机制
- ✅ token自动添加
- ✅ 统一错误消息

**`category.js` - 通用分类API**
- 跨多个功能模块使用的分类接口
- 文章分类、任务分类等通用分类操作

## 🚀 与功能API的区别

| 维度 | `src/api/` | `features/*/api/` |
|------|------------|-------------------|
| **职责** | 全局配置、跨模块API | 功能专用API |
| **内容** | axios配置、通用接口 | 业务逻辑接口 |
| **使用范围** | 全项目共享 | 功能模块内部 |
| **示例** | 分类管理 | 用户认证、文章CRUD、统计数据 |

## 📝 使用指南

### 1. axios配置使用
```javascript
// 在功能API中使用全局axios实例
import api from '@/api/index'

export const myApi = {
  getData() {
    return api.get('/my-endpoint')
  }
}
```

### 2. 通用API使用
```javascript
// 使用通用分类API
import { categoryApi } from '@/api/category'

const categories = await categoryApi.getCategories('article')
```

### 3. 功能API使用
```javascript
// 使用功能专用API
import { authApi } from '@/features/auth/api'
import { analyticsApi } from '@/features/analytics/api'

const user = await authApi.login({ email, password })
const trends = await analyticsApi.getTrends({ metric: 'tasks' })
```

## ⚠️ 重要原则

### ✅ 应该放在这里的
- axios基础配置
- 跨模块使用的通用API
- 全局错误处理
- 请求/响应拦截器

### ❌ 不应该放在这里的
- 功能特定的业务API
- 单一模块使用的接口
- 复杂的业务逻辑

## 🔄 迁移说明

**已完成的API重构:**
- ✅ `auth.js` → `features/auth/api/`
- ✅ `todo.js` → `features/todos/api/`
- ✅ `article.js` → `features/articles/api/`
- ✅ `notification.js` → `features/notifications/api/`
- ✅ `settings.js` → `features/settings/api/`
- ✅ `statistics.js` → `features/analytics/api/`

**保留的文件:**
- ✅ `index.js` - 核心axios配置
- ✅ `category.js` - 跨模块通用分类API

---

💡 **注意**: 这种架构确保了清晰的职责分离，全局配置统一管理，功能API模块化管理，提高了项目的可维护性和扩展性。