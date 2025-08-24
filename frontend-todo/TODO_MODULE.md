# TODO清单管理系统模块

## 功能特性

- ✅ 用户认证（登录/注册）
- ✅ TODO CRUD操作
- ✅ 优先级管理（5个等级）
- ✅ 分类管理（工作/学习/生活等）
- ✅ 状态管理（待处理/进行中/已完成/已取消）
- ✅ 截止时间管理
- ✅ 现代化前端界面

## 技术栈

### 后端
- Gin + GORM + PostgreSQL/MySQL + Redis + JWT

### 前端
- Vue3 + Element Plus + Pinia + Vite

## 快速启动

```bash
# 一键启动
make todo-start

# 停止服务
make todo-stop
```

## 访问地址

- 前端: http://localhost:3000
- 后端: http://localhost:8080

## API接口

- `GET /api/v1/todos` - 获取TODO列表
- `POST /api/v1/todos` - 创建TODO
- `PUT /api/v1/todos/:id` - 更新TODO
- `DELETE /api/v1/todos/:id` - 删除TODO

## 项目结构

```
├── internal/
│   ├── models/todo.go          # TODO数据模型
│   ├── service/todo_service.go # 业务逻辑
│   ├── handler/todo_handler.go # API处理器
│   └── database/seed.go        # 种子数据
├── frontend/                   # Vue3前端
│   ├── src/views/Todos.vue     # 主页面
│   ├── src/stores/todo.js      # 状态管理
│   └── src/api/todo.js         # API接口
└── scripts/
    ├── start-todo.sh           # 启动脚本
    └── stop-todo.sh            # 停止脚本
```
