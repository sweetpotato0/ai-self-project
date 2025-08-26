# 🏪 全局状态管理目录

这个目录是项目的**统一状态管理中心**，基于 Pinia 实现，管理整个应用的全局状态。

## 📋 设计理念

### 🎯 统一管理原则
- **集中化管理**：所有应用状态统一放在此目录下
- **按功能分离**：每个功能模块对应一个独立的 Store 文件
- **职责清晰**：避免状态分散和管理混乱

### 🔧 架构优势
1. **易于维护** - 开发者明确知道状态文件的位置
2. **便于调试** - Pinia DevTools 可以统一监控所有状态
3. **避免冗余** - 不会出现"状态应该放在哪里"的困扰
4. **团队协作** - 统一的状态管理规范

## 📁 目录结构

```
stores/
├── README.md           # 本说明文件
├── auth.js            # 🔐 用户认证状态
├── todo.js            # ✅ 待办事项状态
├── article.js         # 📝 文章管理状态
├── notification.js    # 🔔 通知系统状态
├── settings.js        # ⚙️ 系统设置状态
├── statistics.js      # 📊 数据统计状态
├── category.js        # 🏷️ 分类管理状态
├── modules/           # 🗂️ 全局共享状态模块
│   ├── README.md      # 模块说明
│   ├── app.js         # 应用级状态（主题、语言等）
│   ├── user.js        # 用户全局状态
│   ├── cache.js       # 全局缓存管理
│   └── permission.js  # 权限控制状态
└── plugins/           # 🔧 Pinia 插件
    ├── README.md      # 插件说明
    ├── persistence.js # 持久化插件
    └── logger.js      # 日志插件
```

## 🎯 Store 文件职责

| Store 文件 | 功能描述 | 主要状态 |
|-----------|----------|----------|
| **auth.js** | 用户认证管理 | 用户信息、登录状态、权限 |
| **todo.js** | 待办事项管理 | 任务列表、筛选条件、统计 |
| **article.js** | 文章内容管理 | 文章列表、当前文章、统计 |
| **notification.js** | 通知系统 | 通知列表、未读数、设置 |
| **settings.js** | 系统设置 | 用户偏好、主题、语言 |
| **statistics.js** | 数据统计 | 图表数据、趋势分析 |
| **category.js** | 分类管理 | 分类树、标签系统 |

## 🔄 Store 标准结构

每个 Store 文件都遵循统一的结构模式：

```javascript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useFeatureStore = defineStore('feature', () => {
  // 📦 状态定义
  const state = ref(initialState)
  const loading = ref(false)
  
  // 🔄 计算属性
  const computedValue = computed(() => {
    return processState(state.value)
  })
  
  // ⚡ 异步操作
  const fetchData = async () => {
    loading.value = true
    try {
      // API 调用
      const result = await api.getData()
      state.value = result
    } catch (error) {
      console.error('获取数据失败:', error)
    } finally {
      loading.value = false
    }
  }
  
  // 🔧 同步操作
  const updateState = (newData) => {
    state.value = { ...state.value, ...newData }
  }
  
  // 🧹 重置状态
  const reset = () => {
    state.value = initialState
    loading.value = false
  }
  
  return {
    // 状态
    state,
    loading,
    // 计算属性
    computedValue,
    // 方法
    fetchData,
    updateState,
    reset
  }
})
```

## 📝 使用指南

### 1. 在组件中使用 Store

```javascript
<script setup>
import { useFeatureStore } from '@/stores/feature'

const featureStore = useFeatureStore()

// 访问状态
const data = featureStore.state

// 调用方法
const handleFetch = () => {
  featureStore.fetchData()
}
</script>
```

### 2. Store 之间的通信

```javascript
// 在一个 Store 中使用另一个 Store
import { useAuthStore } from './auth'

export const useFeatureStore = defineStore('feature', () => {
  const authStore = useAuthStore()
  
  const fetchUserData = async () => {
    if (authStore.isLoggedIn) {
      // 执行需要登录的操作
    }
  }
  
  return { fetchUserData }
})
```

### 3. 持久化状态

```javascript
import { defineStore } from 'pinia'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref({})
  
  return { settings }
}, {
  // 启用持久化
  persist: {
    key: 'settings',
    storage: localStorage,
    paths: ['settings'] // 指定需要持久化的状态
  }
})
```

## 🚀 最佳实践

### ✅ 推荐做法
1. **命名规范** - Store 文件使用小写+短横线，Store 名称使用驼峰
2. **状态扁平** - 避免过深的嵌套状态结构
3. **异步处理** - 所有 API 调用都应该有 loading 状态和错误处理
4. **类型安全** - 使用 TypeScript 定义状态类型
5. **文档注释** - 为复杂的状态和方法添加 JSDoc 注释

### ❌ 避免做法
1. **不要** 在 Store 中直接操作 DOM
2. **不要** 在 Store 中使用 Vue 路由
3. **不要** 在 Store 中存储临时的 UI 状态
4. **不要** 创建过于复杂的嵌套状态
5. **不要** 忘记处理异步操作的错误情况

## 🔧 开发工具

### Pinia DevTools
- 安装 Vue DevTools 浏览器扩展
- 在开发模式下可以实时查看和修改 Store 状态
- 支持时间旅行调试和状态导入/导出

### 调试技巧
```javascript
// 在 Store 中添加调试信息
const updateState = (newData) => {
  if (import.meta.env.DEV) {
    console.log('State updated:', { old: state.value, new: newData })
  }
  state.value = { ...state.value, ...newData }
}
```

## 🔄 状态迁移

当需要修改状态结构时：
1. **向后兼容** - 先添加新字段，保留旧字段
2. **渐进迁移** - 逐步迁移组件使用新的状态字段
3. **版本控制** - 为状态结构变更做好版本记录
4. **清理旧代码** - 确认所有组件都迁移后再删除旧字段

---

💡 **提示**: 如果你需要添加新的 Store 文件，请参考现有文件的结构和命名规范，确保项目的一致性和可维护性。