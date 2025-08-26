# Pinia 状态模块目录

这个目录用于存放全局共享的 Pinia 状态模块。

## 🏗️ 模块结构规划

```
modules/
├── app.js          # 应用全局状态（主题、语言、布局等）
├── user.js         # 用户全局状态（跨模块的用户信息）
├── cache.js        # 全局缓存管理
├── permission.js   # 权限状态管理
└── notification.js # 全局通知状态
```

## 🎯 使用原则

### 全局状态 vs 功能状态
- **全局状态**: 跨多个功能模块使用的状态放在这里
- **功能状态**: 功能内部的状态放在对应的 `features/*/stores/`

### 示例：全局用户状态

```javascript
// modules/user.js
import { defineStore } from 'pinia'

export const useGlobalUserStore = defineStore('globalUser', {
  state: () => ({
    profile: null,
    permissions: [],
    preferences: {}
  }),
  
  getters: {
    hasPermission: (state) => (permission) => {
      return state.permissions.includes(permission)
    }
  },
  
  actions: {
    updateProfile(profile) {
      this.profile = profile
    }
  }
})
```

## 📋 模块职责划分

- **app.js**: 应用级配置（主题、语言、侧边栏状态等）
- **user.js**: 跨模块的用户信息（权限、偏好设置等）
- **cache.js**: 应用级缓存（常用数据、配置缓存等）
- **permission.js**: 权限控制（路由权限、功能权限等）