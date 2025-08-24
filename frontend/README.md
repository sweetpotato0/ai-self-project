# TODO清单管理系统 - 前端

基于Vue3 + Element Plus的现代化TODO清单管理系统前端。

## 技术栈

- **Vue 3** - 渐进式JavaScript框架
- **Vue Router 4** - 官方路由管理器
- **Pinia** - 状态管理库
- **Element Plus** - Vue 3组件库
- **Axios** - HTTP客户端
- **Vite** - 构建工具

## 功能特性

- ✅ 用户认证（登录/注册）
- ✅ TODO CRUD操作
- ✅ 优先级管理
- ✅ 状态管理
- ✅ 响应式设计
- ✅ 现代化UI

## 快速开始

### 安装依赖

```bash
cd frontend
npm install
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
frontend/
├── src/
│   ├── api/          # API接口
│   ├── components/   # 组件
│   ├── router/       # 路由配置
│   ├── stores/       # 状态管理
│   ├── utils/        # 工具函数
│   ├── views/        # 页面组件
│   ├── App.vue       # 根组件
│   └── main.js       # 入口文件
├── index.html        # HTML模板
├── package.json      # 依赖配置
├── vite.config.js    # Vite配置
└── README.md         # 说明文档
```

## API接口

### 认证相关
- `POST /api/v1/users/register` - 用户注册
- `POST /api/v1/users/login` - 用户登录
- `GET /api/v1/users/profile` - 获取用户信息

### TODO相关
- `GET /api/v1/todos` - 获取TODO列表
- `POST /api/v1/todos` - 创建TODO
- `PUT /api/v1/todos/:id` - 更新TODO
- `DELETE /api/v1/todos/:id` - 删除TODO

## 开发说明

### 状态管理

使用Pinia进行状态管理：

- `auth` - 认证状态
- `todo` - TODO数据状态

### 路由守卫

- 未登录用户自动跳转到登录页
- 已登录用户访问登录页自动跳转到TODO页面

### 样式规范

- 使用Element Plus组件库
- 响应式设计，支持移动端
- 统一的颜色和间距规范
