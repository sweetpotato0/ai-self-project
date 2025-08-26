<template>
  <div class="json-formatter-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-development' })">开发类</el-breadcrumb-item>
          <el-breadcrumb-item>JSON格式化</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>JSON格式化工具</h1>
      <p>格式化、压缩、验证JSON数据</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Edit /></el-icon>
              <span>JSON输入</span>
            </div>
          </template>
          <div class="input-content">
            <div class="toolbar">
              <el-button-group>
                <el-button @click="formatJson" type="primary" :disabled="!inputJson.trim()">
                  <el-icon><MagicStick /></el-icon>
                  格式化
                </el-button>
                <el-button @click="compressJson" :disabled="!inputJson.trim()">
                  <el-icon><Compress /></el-icon>
                  压缩
                </el-button>
                <el-button @click="validateJson">
                  <el-icon><CircleCheck /></el-icon>
                  验证
                </el-button>
              </el-button-group>
              
              <el-button-group>
                <el-button @click="clearInput">
                  <el-icon><Delete /></el-icon>
                  清空
                </el-button>
                <el-button @click="pasteFromClipboard">
                  <el-icon><DocumentCopy /></el-icon>
                  粘贴
                </el-button>
                <el-button @click="loadSample">
                  <el-icon><Document /></el-icon>
                  示例
                </el-button>
              </el-button-group>
            </div>
            
            <el-input
              v-model="inputJson"
              type="textarea"
              :rows="15"
              placeholder="请输入JSON数据..."
              class="json-input"
              @input="onInputChange"
            />
            
            <div class="input-status" v-if="validationResult">
              <el-alert
                :title="validationResult.message"
                :type="validationResult.type"
                :closable="false"
                show-icon
              />
            </div>
          </div>
        </el-card>
      </div>

      <div class="output-section" v-if="outputJson">
        <el-card class="output-card">
          <template #header>
            <div class="card-header">
              <el-icon><View /></el-icon>
              <span>格式化结果</span>
              <div class="header-actions">
                <el-button @click="copyOutput" size="small" text>
                  <el-icon><DocumentCopy /></el-icon>
                  复制
                </el-button>
              </div>
            </div>
          </template>
          <div class="output-content">
            <el-input
              v-model="outputJson"
              type="textarea"
              :rows="15"
              readonly
              class="json-output"
            />
          </div>
        </el-card>
      </div>

      <div class="analysis-section" v-if="jsonAnalysis">
        <el-card class="analysis-card">
          <template #header>
            <div class="card-header">
              <el-icon><DataAnalysis /></el-icon>
              <span>JSON分析</span>
            </div>
          </template>
          <div class="analysis-content">
            <div class="analysis-grid">
              <div class="analysis-item">
                <div class="analysis-label">大小</div>
                <div class="analysis-value">{{ formatBytes(jsonAnalysis.size) }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">字符数</div>
                <div class="analysis-value">{{ jsonAnalysis.characters.toLocaleString() }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">行数</div>
                <div class="analysis-value">{{ jsonAnalysis.lines.toLocaleString() }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">深度</div>
                <div class="analysis-value">{{ jsonAnalysis.depth }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">对象数</div>
                <div class="analysis-value">{{ jsonAnalysis.objects }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">数组数</div>
                <div class="analysis-value">{{ jsonAnalysis.arrays }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">字符串数</div>
                <div class="analysis-value">{{ jsonAnalysis.strings }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">数字数</div>
                <div class="analysis-value">{{ jsonAnalysis.numbers }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">布尔值数</div>
                <div class="analysis-value">{{ jsonAnalysis.booleans }}</div>
              </div>
              <div class="analysis-item">
                <div class="analysis-label">null值数</div>
                <div class="analysis-value">{{ jsonAnalysis.nulls }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <div class="tree-section" v-if="jsonTree.length > 0">
        <el-card class="tree-card">
          <template #header>
            <div class="card-header">
              <el-icon><Share /></el-icon>
              <span>JSON结构树</span>
            </div>
          </template>
          <div class="tree-content">
            <el-tree
              :data="jsonTree"
              :props="{ children: 'children', label: 'label' }"
              default-expand-all
              show-checkbox
              node-key="id"
              class="json-tree"
            />
          </div>
        </el-card>
      </div>

      <div class="tools-section">
        <el-card class="tools-card">
          <template #header>
            <div class="card-header">
              <el-icon><Tools /></el-icon>
              <span>快捷工具</span>
            </div>
          </template>
          <div class="tools-content-inner">
            <div class="tools-grid">
              <el-button @click="escapeJson" :disabled="!outputJson">
                <el-icon><Key /></el-icon>
                转义JSON
              </el-button>
              <el-button @click="unescapeJson" :disabled="!inputJson.trim()">
                <el-icon><Unlock /></el-icon>
                反转义
              </el-button>
              <el-button @click="sortKeys" :disabled="!outputJson">
                <el-icon><Sort /></el-icon>
                排序键名
              </el-button>
              <el-button @click="downloadJson" :disabled="!outputJson">
                <el-icon><Download /></el-icon>
                下载JSON
              </el-button>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@/composables/useClipboard'
import {
  Edit,
  View,
  MagicStick,
  Compress,
  CircleCheck,
  Delete,
  DocumentCopy,
  Document,
  DataAnalysis,
  Share,
  Tools,
  Key,
  Unlock,
  Sort,
  Download
} from '@element-plus/icons-vue'

const router = useRouter()
const { copy } = useClipboard()

const inputJson = ref('')
const outputJson = ref('')
const validationResult = ref(null)

const jsonAnalysis = computed(() => {
  if (!outputJson.value) return null
  
  try {
    const parsed = JSON.parse(outputJson.value)
    const analysis = analyzeJson(parsed)
    analysis.size = new Blob([outputJson.value]).size
    analysis.characters = outputJson.value.length
    analysis.lines = outputJson.value.split('\n').length
    return analysis
  } catch {
    return null
  }
})

const jsonTree = computed(() => {
  if (!outputJson.value) return []
  
  try {
    const parsed = JSON.parse(outputJson.value)
    return buildTree(parsed, 'root')
  } catch {
    return []
  }
})

const analyzeJson = (obj, depth = 0) => {
  let stats = {
    depth: depth,
    objects: 0,
    arrays: 0,
    strings: 0,
    numbers: 0,
    booleans: 0,
    nulls: 0
  }
  
  if (obj === null) {
    stats.nulls++
  } else if (typeof obj === 'object') {
    if (Array.isArray(obj)) {
      stats.arrays++
      obj.forEach(item => {
        const childStats = analyzeJson(item, depth + 1)
        mergeStats(stats, childStats)
      })
    } else {
      stats.objects++
      Object.values(obj).forEach(value => {
        const childStats = analyzeJson(value, depth + 1)
        mergeStats(stats, childStats)
      })
    }
  } else if (typeof obj === 'string') {
    stats.strings++
  } else if (typeof obj === 'number') {
    stats.numbers++
  } else if (typeof obj === 'boolean') {
    stats.booleans++
  }
  
  return stats
}

const mergeStats = (target, source) => {
  target.depth = Math.max(target.depth, source.depth)
  target.objects += source.objects
  target.arrays += source.arrays
  target.strings += source.strings
  target.numbers += source.numbers
  target.booleans += source.booleans
  target.nulls += source.nulls
}

const buildTree = (obj, key, path = '') => {
  const currentPath = path ? `${path}.${key}` : key
  
  if (obj === null) {
    return [{
      id: currentPath,
      label: `${key}: null`,
      children: []
    }]
  }
  
  if (typeof obj !== 'object') {
    const type = typeof obj
    const value = type === 'string' ? `"${obj}"` : String(obj)
    return [{
      id: currentPath,
      label: `${key}: ${value} (${type})`,
      children: []
    }]
  }
  
  if (Array.isArray(obj)) {
    return [{
      id: currentPath,
      label: `${key}: Array[${obj.length}]`,
      children: obj.flatMap((item, index) => 
        buildTree(item, `[${index}]`, currentPath)
      )
    }]
  }
  
  return [{
    id: currentPath,
    label: `${key}: Object{${Object.keys(obj).length}}`,
    children: Object.entries(obj).flatMap(([childKey, childValue]) =>
      buildTree(childValue, childKey, currentPath)
    )
  }]
}

const onInputChange = () => {
  validationResult.value = null
}

const validateJson = () => {
  if (!inputJson.value.trim()) {
    validationResult.value = {
      type: 'warning',
      message: '请输入JSON数据'
    }
    return false
  }
  
  try {
    JSON.parse(inputJson.value)
    validationResult.value = {
      type: 'success',
      message: 'JSON格式正确'
    }
    return true
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `JSON格式错误: ${error.message}`
    }
    return false
  }
}

const formatJson = () => {
  if (!validateJson()) return
  
  try {
    const parsed = JSON.parse(inputJson.value)
    outputJson.value = JSON.stringify(parsed, null, 2)
    ElMessage.success('JSON格式化成功')
  } catch (error) {
    ElMessage.error('格式化失败: ' + error.message)
  }
}

const compressJson = () => {
  if (!validateJson()) return
  
  try {
    const parsed = JSON.parse(inputJson.value)
    outputJson.value = JSON.stringify(parsed)
    ElMessage.success('JSON压缩成功')
  } catch (error) {
    ElMessage.error('压缩失败: ' + error.message)
  }
}

const escapeJson = () => {
  outputJson.value = JSON.stringify(outputJson.value)
  ElMessage.success('JSON转义成功')
}

const unescapeJson = () => {
  try {
    const unescaped = JSON.parse(inputJson.value)
    if (typeof unescaped === 'string') {
      outputJson.value = unescaped
      ElMessage.success('JSON反转义成功')
    } else {
      ElMessage.error('输入的不是转义后的JSON字符串')
    }
  } catch (error) {
    ElMessage.error('反转义失败: ' + error.message)
  }
}

const sortKeys = () => {
  try {
    const parsed = JSON.parse(outputJson.value)
    const sorted = sortObjectKeys(parsed)
    outputJson.value = JSON.stringify(sorted, null, 2)
    ElMessage.success('键名排序成功')
  } catch (error) {
    ElMessage.error('排序失败: ' + error.message)
  }
}

const sortObjectKeys = (obj) => {
  if (obj === null || typeof obj !== 'object') {
    return obj
  }
  
  if (Array.isArray(obj)) {
    return obj.map(item => sortObjectKeys(item))
  }
  
  const sortedKeys = Object.keys(obj).sort()
  const sortedObj = {}
  
  sortedKeys.forEach(key => {
    sortedObj[key] = sortObjectKeys(obj[key])
  })
  
  return sortedObj
}

const clearInput = () => {
  inputJson.value = ''
  outputJson.value = ''
  validationResult.value = null
}

const pasteFromClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    inputJson.value = text
    validationResult.value = null
    ElMessage.success('已粘贴到输入框')
  } catch (error) {
    ElMessage.error('粘贴失败: ' + error.message)
  }
}

const copyOutput = async () => {
  if (!outputJson.value) return
  await copy(outputJson.value)
}

const loadSample = () => {
  inputJson.value = JSON.stringify({
    "name": "张三",
    "age": 30,
    "city": "北京",
    "skills": ["JavaScript", "Vue.js", "Node.js"],
    "address": {
      "street": "朝阳路100号",
      "zipCode": "100000",
      "coordinates": {
        "lat": 39.9042,
        "lng": 116.4074
      }
    },
    "isEmployed": true,
    "spouse": null,
    "projects": [
      {
        "name": "项目A",
        "status": "completed",
        "technologies": ["Vue", "Element Plus"]
      },
      {
        "name": "项目B",
        "status": "in-progress",
        "technologies": ["React", "TypeScript"]
      }
    ]
  }, null, 2)
  
  validationResult.value = null
  ElMessage.success('示例数据已加载')
}

const downloadJson = () => {
  if (!outputJson.value) return
  
  const blob = new Blob([outputJson.value], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `formatted_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  
  ElMessage.success('JSON文件已下载')
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<style scoped>
.json-formatter-container {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.tools-header {
  margin-bottom: 30px;
  text-align: left;
}

.tools-header h1,
.tools-header p {
  text-align: center;
}

.breadcrumb {
  margin-bottom: 16px;
}

.breadcrumb .el-breadcrumb-item {
  cursor: pointer;
}

.tools-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.tools-header p {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

.tools-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.header-actions {
  margin-left: auto;
}

.input-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.json-input,
.json-output {
  font-family: 'Courier New', Monaco, monospace;
  font-size: 14px;
}

.json-input :deep(.el-textarea__inner) {
  font-family: 'Courier New', Monaco, monospace;
}

.json-output :deep(.el-textarea__inner) {
  font-family: 'Courier New', Monaco, monospace;
  background: #f8fafc;
}

.input-status {
  margin-top: 8px;
}

.analysis-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
}

.analysis-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.analysis-label {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 4px;
}

.analysis-value {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
}

.tree-content {
  max-height: 400px;
  overflow-y: auto;
}

.json-tree {
  font-family: 'Courier New', Monaco, monospace;
  font-size: 13px;
}

.tools-content-inner {
  padding: 16px 0;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}

@media (max-width: 768px) {
  .json-formatter-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .toolbar {
    flex-direction: column;
  }
  
  .analysis-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .tools-grid {
    grid-template-columns: 1fr;
  }
}
</style>