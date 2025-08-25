<template>
  <div class="string-generator-container">
    <div class="tool-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-development' })">开发类</el-breadcrumb-item>
          <el-breadcrumb-item>字符串生成</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>字符串生成工具</h1>
      <p>生成随机字符串，支持自定义长度和字符集</p>
    </div>

    <div class="generator-section">
      <div class="config-panel">
        <div class="panel-header">
          <div class="panel-title">
            <el-icon><Setting /></el-icon>
            配置选项
          </div>
        </div>
        
        <div class="config-content">
          <!-- 长度设置 -->
          <div class="config-item">
            <label class="config-label">字符串长度</label>
            <div class="length-controls">
              <el-input-number
                v-model="length"
                :min="1"
                :max="1000"
                :step="1"
                size="large"
                class="length-input"
              />
              <div class="length-presets">
                <el-button 
                  v-for="preset in lengthPresets" 
                  :key="preset"
                  size="small"
                  @click="length = preset"
                  :type="length === preset ? 'primary' : ''"
                >
                  {{ preset }}
                </el-button>
              </div>
            </div>
          </div>

          <!-- 字符集选择 -->
          <div class="config-item">
            <label class="config-label">包含字符类型</label>
            <div class="charset-options">
              <el-checkbox 
                v-model="options.uppercase"
                size="large"
                class="charset-checkbox"
              >
                <span class="checkbox-content">
                  <el-icon><Star /></el-icon>
                  大写字母 (A-Z)
                </span>
              </el-checkbox>
              
              <el-checkbox 
                v-model="options.lowercase"
                size="large"
                class="charset-checkbox"
              >
                <span class="checkbox-content">
                  <el-icon><Star /></el-icon>
                  小写字母 (a-z)
                </span>
              </el-checkbox>
              
              <el-checkbox 
                v-model="options.numbers"
                size="large"
                class="charset-checkbox"
              >
                <span class="checkbox-content">
                  <el-icon><Star /></el-icon>
                  数字 (0-9)
                </span>
              </el-checkbox>
              
              <el-checkbox 
                v-model="options.symbols"
                size="large"
                class="charset-checkbox"
              >
                <span class="checkbox-content">
                  <el-icon><Star /></el-icon>
                  特殊符号 (!@#$%^&*)
                </span>
              </el-checkbox>
              
              <el-checkbox 
                v-model="options.similar"
                size="large"
                class="charset-checkbox"
              >
                <span class="checkbox-content">
                  <el-icon><Warning /></el-icon>
                  包含易混淆字符 (0O1lI)
                </span>
              </el-checkbox>
            </div>
          </div>

          <!-- 生成按钮 -->
          <div class="generate-section">
            <el-button 
              type="primary"
              size="large"
              @click="generateString"
              :disabled="!hasValidOptions"
              class="generate-btn"
            >
              <el-icon><Refresh /></el-icon>
              生成字符串
            </el-button>
            
            <el-button 
              size="large"
              @click="generateMultiple"
              :disabled="!hasValidOptions"
              class="batch-btn"
            >
              <el-icon><CopyDocument /></el-icon>
              批量生成 (5个)
            </el-button>
          </div>
        </div>
      </div>

      <div class="result-panel">
        <div class="panel-header">
          <div class="panel-title">
            <el-icon><Document /></el-icon>
            生成结果
          </div>
          <div class="panel-actions" v-if="generatedStrings.length > 0">
            <el-button 
              type="text"
              size="small"
              @click="clearResults"
              class="clear-btn"
            >
              <el-icon><Delete /></el-icon>
              清空
            </el-button>
          </div>
        </div>
        
        <div class="results-content">
          <div v-if="generatedStrings.length === 0" class="empty-state">
            <el-icon class="empty-icon"><DocumentCopy /></el-icon>
            <p>点击"生成字符串"开始生成</p>
          </div>
          
          <div v-else class="results-list">
            <div 
              v-for="(str, index) in generatedStrings" 
              :key="index"
              class="result-item"
            >
              <div class="result-string">{{ str }}</div>
              <div class="result-actions">
                <el-button 
                  type="text"
                  size="small"
                  @click="copyToClipboard(str)"
                  class="copy-btn"
                >
                  <el-icon><CopyDocument /></el-icon>
                </el-button>
                <el-button 
                  type="text"
                  size="small"
                  @click="removeResult(index)"
                  class="remove-btn"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
          
          <div v-if="generatedStrings.length > 1" class="batch-actions">
            <el-button 
              size="small"
              @click="copyAllResults"
              class="copy-all-btn"
            >
              <el-icon><CopyDocument /></el-icon>
              复制全部
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 使用场景示例 -->
    <div class="examples-section">
      <h3>常用场景</h3>
      <div class="examples-grid">
        <div 
          v-for="example in examples" 
          :key="example.name"
          class="example-card"
          @click="applyExample(example)"
        >
          <div class="example-header">
            <div class="example-icon">
              <el-icon><component :is="example.icon" /></el-icon>
            </div>
            <div class="example-info">
              <div class="example-title">{{ example.name }}</div>
              <div class="example-desc">{{ example.description }}</div>
            </div>
          </div>
          <div class="example-config">
            <span class="config-tag">长度: {{ example.length }}</span>
            <span class="config-tag">类型: {{ example.types.join(', ') }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Setting,
  Star,
  Warning,
  Refresh,
  CopyDocument,
  Document,
  Delete,
  Close,
  Key,
  Lock,
  User,
  Connection
} from '@element-plus/icons-vue'

const router = useRouter()

// 响应式数据
const length = ref(12)
const lengthPresets = ref([6, 8, 12, 16, 20, 32])

const options = ref({
  uppercase: true,
  lowercase: true,
  numbers: true,
  symbols: false,
  similar: false
})

const generatedStrings = ref([])

// 字符集定义
const charsets = {
  uppercase: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
  lowercase: 'abcdefghijklmnopqrstuvwxyz', 
  numbers: '0123456789',
  symbols: '!@#$%^&*()_+-=[]{}|;:,.<>?',
  similar: '0O1lI' // 易混淆字符
}

// 使用场景示例
const examples = ref([
  {
    name: '密码生成',
    description: '包含大小写、数字和符号的强密码',
    icon: 'Lock',
    length: 16,
    types: ['大写', '小写', '数字', '符号'],
    config: {
      uppercase: true,
      lowercase: true,
      numbers: true,
      symbols: true,
      similar: false
    }
  },
  {
    name: 'API密钥',
    description: '用于API认证的随机密钥',
    icon: 'Key',
    length: 32,
    types: ['大写', '小写', '数字'],
    config: {
      uppercase: true,
      lowercase: true,
      numbers: true,
      symbols: false,
      similar: false
    }
  },
  {
    name: '用户名',
    description: '友好的用户名标识符',
    icon: 'User',
    length: 8,
    types: ['小写', '数字'],
    config: {
      uppercase: false,
      lowercase: true,
      numbers: true,
      symbols: false,
      similar: false
    }
  },
  {
    name: '会话ID',
    description: '临时会话标识符',
    icon: 'Connection',
    length: 20,
    types: ['大写', '数字'],
    config: {
      uppercase: true,
      lowercase: false,
      numbers: true,
      symbols: false,
      similar: false
    }
  }
])

// 计算属性
const hasValidOptions = computed(() => {
  return options.value.uppercase || 
         options.value.lowercase || 
         options.value.numbers || 
         options.value.symbols
})

const availableChars = computed(() => {
  let chars = ''
  
  if (options.value.uppercase) chars += charsets.uppercase
  if (options.value.lowercase) chars += charsets.lowercase
  if (options.value.numbers) chars += charsets.numbers
  if (options.value.symbols) chars += charsets.symbols
  
  // 如果不包含易混淆字符，则移除它们
  if (!options.value.similar && chars) {
    const similarChars = charsets.similar
    chars = chars.split('').filter(char => !similarChars.includes(char)).join('')
  }
  
  return chars
})

// 生成随机字符串
const generateRandomString = (len = length.value) => {
  if (!availableChars.value) {
    ElMessage.error('请至少选择一种字符类型')
    return ''
  }
  
  let result = ''
  const charactersLength = availableChars.value.length
  
  for (let i = 0; i < len; i++) {
    result += availableChars.value.charAt(Math.floor(Math.random() * charactersLength))
  }
  
  return result
}

// 生成单个字符串
const generateString = () => {
  const newString = generateRandomString()
  if (newString) {
    generatedStrings.value.unshift(newString)
    // 限制最多保存20个结果
    if (generatedStrings.value.length > 20) {
      generatedStrings.value = generatedStrings.value.slice(0, 20)
    }
    ElMessage.success('字符串生成成功')
  }
}

// 批量生成
const generateMultiple = () => {
  const newStrings = []
  for (let i = 0; i < 5; i++) {
    const str = generateRandomString()
    if (str) newStrings.push(str)
  }
  
  if (newStrings.length > 0) {
    generatedStrings.value.unshift(...newStrings)
    // 限制最多保存20个结果
    if (generatedStrings.value.length > 20) {
      generatedStrings.value = generatedStrings.value.slice(0, 20)
    }
    ElMessage.success(`成功生成 ${newStrings.length} 个字符串`)
  }
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

// 复制所有结果
const copyAllResults = async () => {
  const allText = generatedStrings.value.join('\n')
  try {
    await navigator.clipboard.writeText(allText)
    ElMessage.success('已复制所有结果到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 移除单个结果
const removeResult = (index) => {
  generatedStrings.value.splice(index, 1)
}

// 清空所有结果
const clearResults = () => {
  generatedStrings.value = []
  ElMessage.success('已清空所有结果')
}

// 应用示例配置
const applyExample = (example) => {
  length.value = example.length
  options.value = { ...example.config }
  ElMessage.success(`已应用"${example.name}"配置`)
}
</script>

<style scoped>
.string-generator-container {
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

.generator-section {
  display: grid;
  grid-template-columns: 400px 1fr;
  gap: 32px;
  margin-bottom: 40px;
}

.config-panel,
.result-panel {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.panel-header {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  color: #1f2937;
  font-size: 16px;
}

.panel-actions {
  display: flex;
  gap: 8px;
}

.clear-btn,
.copy-btn,
.remove-btn {
  color: #6b7280;
  transition: all 0.2s ease;
}

.clear-btn:hover {
  color: #ef4444;
}

.copy-btn:hover {
  color: #10b981;
}

.remove-btn:hover {
  color: #ef4444;
}

.config-content {
  padding: 24px;
}

.config-item {
  margin-bottom: 32px;
}

.config-item:last-child {
  margin-bottom: 0;
}

.config-label {
  display: block;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
  font-size: 15px;
}

.length-controls {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.length-input {
  width: 100%;
}

.length-presets {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.charset-options {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.charset-checkbox {
  font-size: 15px;
}

.checkbox-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.generate-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.generate-btn,
.batch-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
}

.generate-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.results-content {
  padding: 24px;
  min-height: 400px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.result-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  transition: all 0.2s ease;
}

.result-item:hover {
  background: #f1f5f9;
  border-color: #409eff;
}

.result-string {
  font-family: 'Fira Code', 'Monaco', 'Menlo', monospace;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  letter-spacing: 0.5px;
  flex: 1;
  margin-right: 16px;
  word-break: break-all;
}

.result-actions {
  display: flex;
  gap: 8px;
}

.batch-actions {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: center;
}

.copy-all-btn {
  background: rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.2);
  color: #409eff;
}

.examples-section {
  margin-top: 40px;
}

.examples-section h3 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 20px;
}

.examples-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
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
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.example-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
}

.example-info {
  flex: 1;
}

.example-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.example-desc {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.4;
}

.example-config {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.config-tag {
  font-size: 12px;
  color: #409eff;
  background: rgba(64, 158, 255, 0.1);
  padding: 4px 8px;
  border-radius: 6px;
  border: 1px solid rgba(64, 158, 255, 0.2);
}

@media (max-width: 1024px) {
  .generator-section {
    grid-template-columns: 1fr;
    gap: 24px;
  }
  
  .examples-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .config-content {
    padding: 20px;
  }
  
  .results-content {
    padding: 20px;
  }
  
  .example-card {
    padding: 16px;
  }
}
</style>