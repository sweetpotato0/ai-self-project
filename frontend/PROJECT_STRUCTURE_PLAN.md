# 项目结构重构方案

## 🎯 重构目标

1. **功能模块化**：按功能领域组织代码
2. **组件复用性**：提取公共组件，减少重复代码
3. **可维护性**：清晰的目录结构和命名规范
4. **扩展性**：便于添加新功能和模块

## 📁 新的项目结构

```
src/
├── App.vue
├── main.js
├── assets/                    # 静态资源
│   ├── images/
│   ├── icons/
│   └── fonts/
├── components/                # 通用组件
│   ├── common/               # 基础通用组件
│   │   ├── LoadingSpinner.vue
│   │   ├── ErrorAlert.vue
│   │   ├── ExportButton.vue
│   │   ├── CopyButton.vue
│   │   └── ConfirmDialog.vue
│   ├── ui/                   # UI组件
│   │   ├── Card.vue
│   │   ├── InputGroup.vue
│   │   ├── ResultPanel.vue
│   │   └── StatusTag.vue
│   ├── charts/               # 图表组件
│   │   ├── BarChart.vue
│   │   ├── LineChart.vue
│   │   ├── PieChart.vue
│   │   └── ScatterChart.vue
│   ├── tools/                # 工具组件
│   │   ├── ToolCard.vue
│   │   ├── ToolHeader.vue
│   │   ├── ToolNavigation.vue
│   │   └── HistoryManager.vue
│   └── forms/                # 表单组件
│       ├── FormField.vue
│       ├── FileUpload.vue
│       └── CodeEditor.vue
├── features/                  # 功能模块
│   ├── auth/                 # 认证模块
│   │   ├── components/
│   │   │   └── LoginForm.vue
│   │   ├── stores/
│   │   │   └── authStore.js
│   │   ├── api/
│   │   │   └── authApi.js
│   │   └── views/
│   │       └── Login.vue
│   ├── dashboard/            # 仪表盘模块
│   │   ├── components/
│   │   │   ├── Sidebar.vue
│   │   │   ├── TopBar.vue
│   │   │   └── StatsCard.vue
│   │   ├── stores/
│   │   │   └── dashboardStore.js
│   │   └── views/
│   │       ├── Dashboard.vue
│   │       └── DashboardHome.vue
│   ├── tools/                # 工具模块
│   │   ├── shared/           # 工具共享组件和逻辑
│   │   │   ├── components/
│   │   │   │   ├── ToolsGrid.vue
│   │   │   │   ├── CategoryCard.vue
│   │   │   │   └── ToolBreadcrumb.vue
│   │   │   ├── composables/
│   │   │   │   ├── useHistory.js
│   │   │   │   ├── useExport.js
│   │   │   │   └── useValidation.js
│   │   │   └── utils/
│   │   │       ├── formatters.js
│   │   │       └── validators.js
│   │   ├── development/      # 开发工具
│   │   │   ├── components/
│   │   │   │   └── CodeFormatter.vue
│   │   │   ├── views/
│   │   │   │   ├── ToolsDevelopment.vue
│   │   │   │   ├── TimestampConverter.vue
│   │   │   │   ├── JsonTools.vue
│   │   │   │   ├── StringGenerator.vue
│   │   │   │   └── HttpStatusCodes.vue
│   │   │   └── utils/
│   │   │       └── codeUtils.js
│   │   ├── text/             # 文本工具
│   │   │   ├── components/
│   │   │   │   └── TextTransformer.vue
│   │   │   ├── views/
│   │   │   │   ├── ToolsText.vue
│   │   │   │   ├── Base64Encoder.vue
│   │   │   │   ├── UrlEncoder.vue
│   │   │   │   └── TextProcessor.vue
│   │   │   └── utils/
│   │   │       └── textUtils.js
│   │   ├── image/            # 图像工具
│   │   │   ├── components/
│   │   │   │   ├── ImagePreview.vue
│   │   │   │   └── ImageUploader.vue
│   │   │   ├── views/
│   │   │   │   ├── ToolsImage.vue
│   │   │   │   ├── ImageCompressor.vue
│   │   │   │   ├── ImageConverter.vue
│   │   │   │   └── ImageResizer.vue
│   │   │   └── utils/
│   │   │       └── imageUtils.js
│   │   ├── network/          # 网络工具
│   │   │   ├── components/
│   │   │   │   └── NetworkTester.vue
│   │   │   ├── views/
│   │   │   │   ├── ToolsOperations.vue
│   │   │   │   ├── PingTool.vue
│   │   │   │   ├── PortScanner.vue
│   │   │   │   └── DNSLookup.vue
│   │   │   └── utils/
│   │   │       └── networkUtils.js
│   │   ├── query/            # 查询工具
│   │   │   ├── components/
│   │   │   │   └── QueryResult.vue
│   │   │   ├── views/
│   │   │   │   ├── ToolsQuery.vue
│   │   │   │   ├── IPLookup.vue
│   │   │   │   ├── WhoisLookup.vue
│   │   │   │   └── DomainInfo.vue
│   │   │   └── utils/
│   │   │       └── queryUtils.js
│   │   ├── academic/         # 学术工具
│   │   │   ├── components/
│   │   │   │   ├── Citation.vue
│   │   │   │   └── DataChart.vue
│   │   │   ├── views/
│   │   │   │   ├── ToolsAcademic.vue
│   │   │   │   ├── CitationGenerator.vue
│   │   │   │   ├── MathCalculator.vue
│   │   │   │   └── DataAnalyzer.vue
│   │   │   └── utils/
│   │   │       ├── mathUtils.js
│   │   │       └── citationUtils.js
│   │   ├── document/         # 文档工具
│   │   │   ├── components/
│   │   │   │   ├── PDFViewer.vue
│   │   │   │   └── MarkdownEditor.vue
│   │   │   ├── views/
│   │   │   │   ├── ToolsDocument.vue
│   │   │   │   ├── PDFTools.vue
│   │   │   │   └── WordCloudGenerator.vue
│   │   │   └── utils/
│   │   │       └── documentUtils.js
│   │   ├── others/           # 其他工具
│   │   │   ├── views/
│   │   │   │   ├── ToolsOthers.vue
│   │   │   │   ├── QRCodeGenerator.vue
│   │   │   │   ├── ColorPicker.vue
│   │   │   │   ├── PasswordStrengthChecker.vue
│   │   │   │   ├── RegexTester.vue
│   │   │   │   └── HashCalculator.vue
│   │   │   └── utils/
│   │   │       └── miscUtils.js
│   │   └── views/
│   │       └── Tools.vue     # 工具主页
│   ├── articles/             # 文章模块
│   │   ├── components/
│   │   │   ├── ArticleCard.vue
│   │   │   ├── ArticleDialog.vue
│   │   │   └── EnhancedEditor.vue
│   │   ├── stores/
│   │   │   └── articleStore.js
│   │   ├── api/
│   │   │   └── articleApi.js
│   │   └── views/
│   │       ├── Articles.vue
│   │       └── ArticleDetail.vue
│   ├── todos/                # 待办事项模块
│   │   ├── components/
│   │   │   ├── TodoItem.vue
│   │   │   └── TodoForm.vue
│   │   ├── stores/
│   │   │   └── todoStore.js
│   │   ├── api/
│   │   │   └── todoApi.js
│   │   └── views/
│   │       └── Todos.vue
│   ├── calendar/             # 日历模块
│   │   ├── components/
│   │   │   └── CalendarView.vue
│   │   └── views/
│   │       └── Calendar.vue
│   ├── analytics/            # 分析模块
│   │   ├── components/
│   │   │   └── AnalyticsChart.vue
│   │   └── views/
│   │       └── Analytics.vue
│   ├── notifications/        # 通知模块
│   │   ├── components/
│   │   │   ├── NotificationPanel.vue
│   │   │   └── WebSocketNotification.vue
│   │   ├── stores/
│   │   │   └── notificationStore.js
│   │   ├── api/
│   │   │   └── notificationApi.js
│   │   └── views/
│   │       └── Notifications.vue
│   ├── profile/              # 个人资料模块
│   │   ├── components/
│   │   │   └── ProfileForm.vue
│   │   └── views/
│   │       └── Profile.vue
│   └── settings/             # 设置模块
│       ├── components/
│       │   └── SettingsPanel.vue
│       ├── stores/
│       │   └── settingsStore.js
│       ├── api/
│       │   └── settingsApi.js
│       └── views/
│           └── Settings.vue
├── composables/              # 组合式函数
│   ├── useAuth.js
│   ├── useNotification.js
│   ├── useTheme.js
│   ├── useLocalStorage.js
│   └── useWebSocket.js
├── router/                   # 路由配置
│   ├── index.js
│   ├── guards.js
│   └── modules/
│       ├── authRoutes.js
│       ├── toolsRoutes.js
│       └── dashboardRoutes.js
├── stores/                   # 全局状态管理
│   ├── index.js
│   ├── modules/
│   │   ├── app.js
│   │   └── user.js
│   └── plugins/
├── utils/                    # 通用工具函数
│   ├── index.js
│   ├── api.js
│   ├── dateTime.js
│   ├── format.js
│   ├── validation.js
│   ├── constants.js
│   ├── localeManager.js
│   └── notificationManager.js
├── styles/                   # 样式文件
│   ├── main.css
│   ├── variables.css
│   ├── themes.css
│   ├── components.css
│   └── utilities.css
└── types/                    # TypeScript类型定义
    ├── api.ts
    ├── components.ts
    └── tools.ts
```

## 🔄 重构计划

### 第一阶段：基础结构调整
1. 创建新的目录结构
2. 移动现有文件到对应的功能模块
3. 更新所有导入路径

### 第二阶段：组件提取和复用
1. 提取公共工具组件
2. 创建通用UI组件
3. 抽象公共逻辑为组合式函数

### 第三阶段：功能模块化
1. 将每个功能封装为独立模块
2. 每个模块包含自己的组件、状态、API和视图
3. 模块间通过明确的接口通信

### 第四阶段：代码优化
1. 统一代码风格和命名规范
2. 添加类型定义
3. 优化构建配置

## 🎨 命名规范

- **文件夹**：kebab-case (如：`user-profile`)
- **组件文件**：PascalCase (如：`UserProfile.vue`)
- **工具函数**：camelCase (如：`formatDate.js`)
- **常量**：UPPER_SNAKE_CASE (如：`API_BASE_URL`)
- **组合式函数**：以use开头 (如：`useAuth.js`)

## 🚀 重构优势

1. **更好的可维护性**：相关代码集中在一起
2. **更高的复用性**：公共组件和逻辑可以复用
3. **更清晰的依赖关系**：模块边界明确
4. **更好的开发体验**：功能独立，便于并行开发
5. **更容易测试**：模块化结构便于单元测试