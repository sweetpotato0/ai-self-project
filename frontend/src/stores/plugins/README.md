# Pinia 插件目录

这个目录用于存放 Pinia 状态管理的插件扩展。

## 🔌 插件类型

### 持久化插件
```javascript
// persistPlugin.js
export const persistPlugin = ({ store }) => {
  // 实现状态持久化逻辑
}
```

### 日志插件
```javascript
// loggerPlugin.js
export const loggerPlugin = ({ store }) => {
  // 实现状态变更日志记录
}
```

### 开发工具插件
```javascript
// devtoolsPlugin.js
export const devtoolsPlugin = ({ store }) => {
  // 增强开发工具集成
}
```

## 📦 使用方式

在 `main.js` 中注册插件：

```javascript
import { createPinia } from 'pinia'
import { persistPlugin } from '@/stores/plugins/persistPlugin'

const pinia = createPinia()
pinia.use(persistPlugin)

app.use(pinia)
```

## 🎯 常见插件场景

- **数据持久化**: localStorage/sessionStorage 同步
- **状态日志**: 开发环境下的状态变更追踪
- **权限控制**: 基于角色的状态访问控制
- **数据同步**: 多 tab 页间的状态同步