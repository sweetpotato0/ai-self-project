# TypeScript 类型定义目录

这个目录用于存放全局的 TypeScript 类型定义文件。

## 📂 目录结构规划

```
types/
├── api.ts          # API 接口类型定义
├── components.ts   # 组件 Props 类型定义
├── stores.ts       # 状态管理类型定义
├── tools.ts        # 工具相关类型定义
├── common.ts       # 通用类型定义
└── index.ts        # 类型导出入口
```

## 🎯 使用场景

- **API 类型**: 定义请求/响应的数据结构
- **组件类型**: 定义组件 Props、Emits 等接口
- **状态类型**: 定义 Store 的 state、actions 类型
- **工具类型**: 定义工具函数的参数和返回值类型

## 💡 命名规范

- 接口名称使用 `PascalCase` (如: `UserInfo`)
- 类型别名使用 `PascalCase` (如: `UserId`)
- 枚举使用 `PascalCase` (如: `UserStatus`)

## 📝 使用示例

```typescript
// api.ts
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: UserInfo
}

// 在组件中使用
import type { LoginRequest } from '@/types/api'
```