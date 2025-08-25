<template>
  <div class="json-tools-container">
    <div class="tool-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-development' })">开发类</el-breadcrumb-item>
          <el-breadcrumb-item>JSON工具</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>JSON 工具箱</h1>
      <p>JSON 格式化、压缩、转换和验证工具</p>
    </div>

    <!-- 功能按钮区 -->
    <div class="action-bar">
      <div class="action-group">
        <el-button 
          type="primary" 
          @click="formatJson"
          :disabled="!jsonInput"
          class="action-btn"
        >
          <el-icon><MagicStick /></el-icon>
          格式化
        </el-button>
        <el-button 
          @click="compressJson"
          :disabled="!jsonInput"
          class="action-btn"
        >
          <el-icon><Fold /></el-icon>
          压缩
        </el-button>
        <el-button 
          @click="convertToYaml"
          :disabled="!isValidJson"
          class="action-btn"
        >
          <el-icon><Expand /></el-icon>
          转YAML
        </el-button>
        <el-button 
          @click="clearAll"
          class="action-btn"
        >
          <el-icon><Delete /></el-icon>
          清空
        </el-button>
      </div>
      <div class="status-group">
        <div class="json-status" :class="{ valid: isValidJson, invalid: !isValidJson && jsonInput }">
          <el-icon><CircleCheckFilled v-if="isValidJson" /><WarningFilled v-else /></el-icon>
          <span>{{ jsonStatus }}</span>
        </div>
      </div>
    </div>

    <!-- 编辑器区域 -->
    <div class="editor-section">
      <div class="editor-panel">
        <div class="panel-header">
          <div class="panel-title">
            <el-icon><Edit /></el-icon>
            JSON 编辑器
          </div>
        </div>
        <div class="simple-editor-container">
          <el-input
            v-model="jsonInput"
            type="textarea"
            :rows="15"
            placeholder="请输入JSON内容..."
            class="json-textarea"
          />
        </div>
      </div>

      <!-- 输出区域 -->
      <div class="output-panel" v-if="outputContent">
        <div class="panel-header">
          <div class="panel-title">
            <el-icon><component :is="outputIcon" /></el-icon>
            {{ outputTitle }}
          </div>
          <div class="panel-actions">
            <el-button 
              type="text" 
              size="small" 
              @click="copyToClipboard(outputContent)"
              class="header-btn"
            >
              <el-icon><DocumentCopy /></el-icon>
              复制
            </el-button>
          </div>
        </div>
        <div class="simple-output-container">
          <el-input
            :model-value="outputContent"
            type="textarea"
            :rows="15"
            readonly
            class="output-textarea"
          />
        </div>
      </div>
    </div>

    <!-- JSON示例 -->
    <div class="examples-section">
      <h3>JSON 示例</h3>
      <div class="examples-grid">
        <div 
          v-for="example in examples" 
          :key="example.name"
          class="example-card"
          @click="useExample(example.json)"
        >
          <div class="example-header">
            <div class="example-title">{{ example.name }}</div>
            <el-icon class="example-icon"><Right /></el-icon>
          </div>
          <div class="example-preview">{{ example.description }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  MagicStick, 
  Fold, 
  Expand, 
  Delete, 
  Edit,
  Document,
  DocumentCopy,
  CircleCheckFilled,
  WarningFilled,
  Right
} from '@element-plus/icons-vue'
import * as yamlParser from 'js-yaml'

const router = useRouter()

// 响应式数据
const jsonInput = ref('')
const outputContent = ref('')
const outputTitle = ref('')
const outputIcon = ref('MagicStick')
const isValidJson = ref(false)
const jsonStatus = ref('请输入JSON')

// JSON示例
const examples = ref([
  {
    name: '用户信息',
    description: '包含用户基本信息的JSON对象',
    json: `{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com",
  "age": 28,
  "address": {
    "province": "北京市",
    "city": "朝阳区",
    "street": "建国路88号"
  },
  "hobbies": ["编程", "阅读", "运动"]
}`
  },
  {
    name: 'API响应',
    description: '典型的API返回数据格式',
    json: `{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "任务一",
        "completed": true,
        "created_at": "2023-08-25T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}`
  },
  {
    name: '配置文件',
    description: '应用程序配置JSON示例',
    json: `{
  "app": {
    "name": "MyApp",
    "version": "1.0.0",
    "debug": false
  },
  "database": {
    "host": "localhost",
    "port": 3306,
    "username": "root",
    "password": "password",
    "database": "myapp"
  }
}`
  }
])

// 验证JSON
const validateJson = (text) => {
  if (!text.trim()) {
    isValidJson.value = false
    jsonStatus.value = '请输入JSON'
    return false
  }
  
  try {
    JSON.parse(text)
    isValidJson.value = true
    jsonStatus.value = 'JSON格式正确'
    return true
  } catch (error) {
    isValidJson.value = false
    jsonStatus.value = `JSON格式错误: ${error.message}`
    return false
  }
}

// 格式化JSON
const formatJson = () => {
  if (!jsonInput.value.trim()) {
    ElMessage.warning('请先输入JSON内容')
    return
  }
  
  try {
    const parsed = JSON.parse(jsonInput.value)
    const formatted = JSON.stringify(parsed, null, 2)
    
    outputContent.value = formatted
    outputTitle.value = '格式化结果'
    outputIcon.value = 'MagicStick'
    
    ElMessage.success('JSON格式化完成')
  } catch (error) {
    ElMessage.error(`格式化失败: ${error.message}`)
  }
}

// 压缩JSON
const compressJson = () => {
  if (!jsonInput.value.trim()) {
    ElMessage.warning('请先输入JSON内容')
    return
  }
  
  try {
    const parsed = JSON.parse(jsonInput.value)
    const compressed = JSON.stringify(parsed)
    
    outputContent.value = compressed
    outputTitle.value = '压缩结果'
    outputIcon.value = 'Fold'
    
    ElMessage.success('JSON压缩完成')
  } catch (error) {
    ElMessage.error(`压缩失败: ${error.message}`)
  }
}

// 转换为YAML
const convertToYaml = () => {
  if (!isValidJson.value) {
    ElMessage.warning('请先输入有效的JSON内容')
    return
  }
  
  try {
    const parsed = JSON.parse(jsonInput.value)
    const yamlContent = yamlParser.dump(parsed, { 
      indent: 2,
      lineWidth: 80,
      noRefs: true 
    })
    
    outputContent.value = yamlContent
    outputTitle.value = 'YAML结果'
    outputIcon.value = 'Document'
    
    ElMessage.success('转换为YAML完成')
  } catch (error) {
    ElMessage.error(`转换失败: ${error.message}`)
  }
}

// 清空所有内容
const clearAll = () => {
  jsonInput.value = ''
  outputContent.value = ''
  outputTitle.value = ''
  
  ElMessage.success('已清空所有内容')
}

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 使用示例
const useExample = (exampleJson) => {
  jsonInput.value = exampleJson
  ElMessage.success('已载入示例JSON')
}

// 监听输入变化
watch(jsonInput, (newValue) => {
  validateJson(newValue)
})
</script>

<style scoped>
.json-tools-container {
  padding: 20px 0;
  max-width: 1400px;
  margin: 0 auto;
}

.tool-header {
  margin-bottom: 30px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.breadcrumb .el-breadcrumb-item {
  cursor: pointer;
}

.tool-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.tool-header p {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.action-group {
  display: flex;
  gap: 12px;
}

.action-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.action-btn:disabled {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.5);
  transform: none;
}

.status-group {
  display: flex;
  align-items: center;
}

.json-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  font-weight: 500;
  transition: all 0.3s ease;
}

.json-status.valid {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.json-status.invalid {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.editor-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  margin-bottom: 40px;
}

.editor-panel,
.output-panel {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.panel-header {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  padding: 16px 20px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #1f2937;
}

.panel-actions {
  display: flex;
  gap: 8px;
}

.header-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #6b7280;
  transition: all 0.2s ease;
}

.header-btn:hover {
  color: #409eff;
  transform: translateY(-1px);
}

.simple-editor-container,
.simple-output-container {
  padding: 0;
}

.json-textarea,
.output-textarea {
  font-family: 'Fira Code', 'Monaco', 'Menlo', monospace !important;
  font-size: 14px;
  line-height: 1.6;
}

.json-textarea :deep(.el-textarea__inner) {
  background-color: #1e1e1e;
  color: #d4d4d4;
  border: 1px solid #3e3e3e;
  border-radius: 8px;
}

.output-textarea :deep(.el-textarea__inner) {
  background-color: #0f172a;
  color: #e2e8f0;
  border: 1px solid #334155;
  border-radius: 8px;
}

.examples-section {
  margin-top: 40px;
}

.examples-section h3 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 16px;
}

.examples-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.example-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.example-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.15);
  border-color: #409eff;
}

.example-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.example-title {
  font-weight: 600;
  color: #1f2937;
}

.example-icon {
  color: #9ca3af;
  transition: all 0.3s ease;
}

.example-card:hover .example-icon {
  color: #409eff;
  transform: translateX(4px);
}

.example-preview {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
}

@media (max-width: 1024px) {
  .editor-section {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .action-bar {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .action-group {
    flex-wrap: wrap;
  }
  
  .examples-grid {
    grid-template-columns: 1fr;
  }
}
</style>