# 📋 待办事项组件

这个目录包含待办事项功能相关的UI组件。

## 🎯 组件说明

### 设计理念
- **专用性**: 仅用于待办事项功能的组件
- **可复用性**: 在待办事项模块内部复用
- **功能聚合**: 相关功能的组件集中管理

### 已实现的组件
- ✅ **CategoryManager** - 任务分类管理组件
  - 分类列表显示
  - 分类创建/编辑/删除
  - 颜色选择器
  - 任务数量统计

### 待创建的组件
- 待办项组件 (TodoItem)
- 待办表单组件 (TodoForm)
- 待办筛选组件 (TodoFilter)
- 待办统计组件 (TodoStats)

### 使用指南

```javascript
// 在待办事项页面中使用
import CategoryManager from '@/features/todos/components/CategoryManager.vue'

// 组件使用示例
<CategoryManager />
```

---

💡 **注意**: 这些组件专门服务于待办事项功能，不应在其他功能模块中使用。跨模块共享的组件应放在 `src/components/` 目录下。