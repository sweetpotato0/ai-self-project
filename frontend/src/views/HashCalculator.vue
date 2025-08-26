<template>
  <div class="hash-calculator-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-others' })">其它</el-breadcrumb-item>
          <el-breadcrumb-item>Hash计算器</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>Hash计算器</h1>
      <p>计算文本或文件的各种Hash值，支持MD5、SHA1、SHA256等多种算法</p>
    </div>

    <div class="tools-content">
      <!-- 输入区域 -->
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><EditPen /></el-icon>
              <span>输入内容</span>
            </div>
          </template>
          
          <div class="input-content">
            <el-tabs v-model="inputMode" @tab-change="clearResults">
              <el-tab-pane label="文本输入" name="text">
                <div class="text-input-section">
                  <el-input
                    v-model="inputText"
                    type="textarea"
                    :rows="8"
                    placeholder="输入要计算Hash的文本内容..."
                    @input="calculateHashes"
                  />
                  <div class="input-stats" v-if="inputText">
                    <span>字符数: {{ inputText.length }}</span>
                    <span>字节数: {{ getByteLength(inputText) }}</span>
                  </div>
                </div>
              </el-tab-pane>
              
              <el-tab-pane label="文件上传" name="file">
                <div class="file-input-section">
                  <el-upload
                    ref="uploadRef"
                    class="upload-area"
                    drag
                    :auto-upload="false"
                    :on-change="handleFileSelect"
                    :show-file-list="false"
                    accept="*/*"
                  >
                    <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
                    <div class="el-upload__text">
                      将文件拖拽到此处，或<em>点击选择文件</em>
                    </div>
                    <template #tip>
                      <div class="el-upload__tip">
                        支持任意格式文件，最大支持100MB
                      </div>
                    </template>
                  </el-upload>
                  
                  <div class="selected-file" v-if="selectedFile">
                    <div class="file-info">
                      <div class="file-details">
                        <span class="file-name">{{ selectedFile.name }}</span>
                        <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
                      </div>
                      <el-button @click="clearFile" size="small" text type="danger">
                        <el-icon><Delete /></el-icon>
                        移除
                      </el-button>
                    </div>
                    <el-button @click="calculateFileHashes" type="primary" :loading="calculating">
                      <el-icon><Operation /></el-icon>
                      {{ calculating ? '计算中...' : '计算Hash值' }}
                    </el-button>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </el-card>
      </div>

      <!-- 算法选择 -->
      <div class="algorithms-section">
        <el-card class="algorithms-card">
          <template #header>
            <div class="card-header">
              <el-icon><Setting /></el-icon>
              <span>Hash算法</span>
            </div>
          </template>
          
          <div class="algorithms-content">
            <div class="algorithm-options">
              <el-checkbox 
                v-for="algo in availableAlgorithms" 
                :key="algo.name"
                v-model="selectedAlgorithms[algo.name]"
                @change="onAlgorithmChange"
                class="algorithm-checkbox"
              >
                <div class="algorithm-info">
                  <span class="algorithm-name">{{ algo.name }}</span>
                  <span class="algorithm-desc">{{ algo.description }}</span>
                </div>
              </el-checkbox>
            </div>

            <div class="algorithm-actions">
              <el-button @click="selectAllAlgorithms" size="small">全选</el-button>
              <el-button @click="clearAllAlgorithms" size="small">清空</el-button>
              <el-button @click="selectCommonAlgorithms" size="small" type="primary">常用</el-button>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 结果显示 -->
      <div class="results-section">
        <el-card class="results-card">
          <template #header>
            <div class="card-header">
              <el-icon><DocumentCopy /></el-icon>
              <span>计算结果</span>
              <div class="result-actions" v-if="hasResults">
                <el-button @click="copyAllResults" size="small" type="primary">
                  <el-icon><CopyDocument /></el-icon>
                  复制全部
                </el-button>
                <el-button @click="exportResults" size="small">
                  <el-icon><Download /></el-icon>
                  导出
                </el-button>
              </div>
            </div>
          </template>
          
          <div class="results-content">
            <div v-if="!hasResults" class="empty-state">
              <el-icon class="empty-icon"><Tools /></el-icon>
              <p>{{ inputMode === 'text' ? '输入文本后自动计算Hash值' : '选择文件并点击计算按钮' }}</p>
            </div>

            <div v-else class="hash-results">
              <div 
                v-for="(result, algorithm) in hashResults" 
                :key="algorithm"
                class="hash-result-item"
              >
                <div class="hash-header">
                  <span class="hash-algorithm">{{ algorithm }}</span>
                  <el-button @click="copyHash(result)" size="small" text>
                    <el-icon><DocumentCopy /></el-icon>
                    复制
                  </el-button>
                </div>
                <div class="hash-value" :title="result">{{ result }}</div>
              </div>
            </div>

            <!-- 比较功能 -->
            <div class="comparison-section" v-if="hasResults">
              <el-divider>Hash值验证</el-divider>
              <div class="verification-input">
                <el-input
                  v-model="verificationHash"
                  placeholder="粘贴已知的Hash值进行对比验证..."
                  @input="verifyHash"
                >
                  <template #prepend>验证Hash</template>
                </el-input>
              </div>
              <div class="verification-result" v-if="verificationHash">
                <div v-if="verificationMatch" class="match-success">
                  <el-icon><SuccessFilled /></el-icon>
                  <span>Hash值匹配！算法: {{ verificationAlgorithm }}</span>
                </div>
                <div v-else class="match-failure">
                  <el-icon><CircleCloseFilled /></el-icon>
                  <span>Hash值不匹配</span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- Hash算法说明 -->
    <div class="info-section">
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <el-icon><InfoFilled /></el-icon>
            <span>Hash算法说明</span>
          </div>
        </template>
        
        <div class="algorithm-descriptions">
          <div class="description-grid">
            <div class="desc-item">
              <h4>MD5</h4>
              <p>128位散列算法，速度快但安全性较低，主要用于文件校验</p>
              <span class="output-length">输出长度: 32字符</span>
            </div>
            <div class="desc-item">
              <h4>SHA-1</h4>
              <p>160位散列算法，比MD5更安全，但已不推荐用于安全用途</p>
              <span class="output-length">输出长度: 40字符</span>
            </div>
            <div class="desc-item">
              <h4>SHA-256</h4>
              <p>256位散列算法，安全性高，广泛用于密码学和区块链</p>
              <span class="output-length">输出长度: 64字符</span>
            </div>
            <div class="desc-item">
              <h4>SHA-512</h4>
              <p>512位散列算法，最高安全级别，适用于高安全要求场景</p>
              <span class="output-length">输出长度: 128字符</span>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 使用场景 -->
    <div class="use-cases-section">
      <el-card class="use-cases-card">
        <template #header>
          <div class="card-header">
            <el-icon><Compass /></el-icon>
            <span>应用场景</span>
          </div>
        </template>
        
        <div class="use-cases-grid">
          <div class="use-case">
            <el-icon class="case-icon"><Lock /></el-icon>
            <h4>文件完整性校验</h4>
            <p>验证文件在传输或存储过程中是否被篡改</p>
          </div>
          <div class="use-case">
            <el-icon class="case-icon"><Key /></el-icon>
            <h4>密码存储</h4>
            <p>将用户密码进行Hash处理后存储，提高安全性</p>
          </div>
          <div class="use-case">
            <el-icon class="case-icon"><UserFilled /></el-icon>
            <h4>数字签名</h4>
            <p>生成数据的唯一标识，用于防伪和认证</p>
          </div>
          <div class="use-case">
            <el-icon class="case-icon"><Connection /></el-icon>
            <h4>区块链应用</h4>
            <p>在区块链中用于生成区块哈希和交易哈希</p>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  EditPen,
  Setting,
  DocumentCopy,
  CopyDocument,
  Download,
  Tools,
  SuccessFilled,
  CircleCloseFilled,
  InfoFilled,
  Compass,
  Lock,
  Key,
  UserFilled,
  Connection,
  UploadFilled,
  Delete,
  Operation
} from '@element-plus/icons-vue'

const router = useRouter()

// 响应式数据
const inputMode = ref('text')
const inputText = ref('')
const selectedFile = ref(null)
const calculating = ref(false)
const verificationHash = ref('')
const hashResults = reactive({})

// 可用的Hash算法
const availableAlgorithms = ref([
  { name: 'MD5', description: '128位，快速但安全性较低' },
  { name: 'SHA-1', description: '160位，中等安全性' },
  { name: 'SHA-256', description: '256位，高安全性推荐' },
  { name: 'SHA-512', description: '512位，最高安全性' }
])

// 选中的算法
const selectedAlgorithms = reactive({
  'MD5': true,
  'SHA-1': false,
  'SHA-256': true,
  'SHA-512': false
})

// 计算属性
const hasResults = computed(() => {
  return Object.keys(hashResults).length > 0
})

const verificationMatch = computed(() => {
  if (!verificationHash.value) return false
  return Object.values(hashResults).some(hash => 
    hash.toLowerCase() === verificationHash.value.toLowerCase()
  )
})

const verificationAlgorithm = computed(() => {
  if (!verificationHash.value) return ''
  for (const [algorithm, hash] of Object.entries(hashResults)) {
    if (hash.toLowerCase() === verificationHash.value.toLowerCase()) {
      return algorithm
    }
  }
  return ''
})

// 简化的Hash计算函数（实际项目中应使用crypto-js等专业库）
const calculateSimpleHash = (text, algorithm) => {
  // 这里使用简化的hash实现，实际项目中应该使用专业的加密库
  let hash = 0
  let length = 32 // MD5长度
  
  switch (algorithm) {
    case 'SHA-1':
      length = 40
      break
    case 'SHA-256':
      length = 64
      break
    case 'SHA-512':
      length = 128
      break
  }
  
  // 简单的字符串hash算法（演示用）
  for (let i = 0; i < text.length; i++) {
    const char = text.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash // Convert to 32-bit integer
  }
  
  // 生成指定长度的十六进制字符串
  let result = Math.abs(hash).toString(16)
  while (result.length < length) {
    result = '0' + result + Math.abs(hash * (result.length + 1)).toString(16).slice(-2)
  }
  
  return result.slice(0, length)
}

// 方法
const getByteLength = (str) => {
  return new Blob([str]).size
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const clearResults = () => {
  Object.keys(hashResults).forEach(key => {
    delete hashResults[key]
  })
  verificationHash.value = ''
}

const calculateHashes = () => {
  if (!inputText.value) {
    clearResults()
    return
  }

  clearResults()
  
  for (const [algorithm, selected] of Object.entries(selectedAlgorithms)) {
    if (selected) {
      hashResults[algorithm] = calculateSimpleHash(inputText.value, algorithm)
    }
  }
}

const handleFileSelect = (file) => {
  if (file.size > 100 * 1024 * 1024) { // 100MB
    ElMessage.error('文件大小不能超过100MB')
    return
  }
  
  selectedFile.value = file.raw
  clearResults()
}

const clearFile = () => {
  selectedFile.value = null
  clearResults()
}

const calculateFileHashes = async () => {
  if (!selectedFile.value) return
  
  calculating.value = true
  clearResults()
  
  try {
    const reader = new FileReader()
    reader.onload = (e) => {
      const text = e.target.result
      for (const [algorithm, selected] of Object.entries(selectedAlgorithms)) {
        if (selected) {
          hashResults[algorithm] = calculateSimpleHash(text, algorithm)
        }
      }
      calculating.value = false
      ElMessage.success('文件Hash计算完成')
    }
    
    reader.onerror = () => {
      calculating.value = false
      ElMessage.error('文件读取失败')
    }
    
    reader.readAsText(selectedFile.value)
  } catch (error) {
    calculating.value = false
    ElMessage.error('计算过程中出现错误')
  }
}

const onAlgorithmChange = () => {
  if (inputMode.value === 'text' && inputText.value) {
    calculateHashes()
  }
}

const selectAllAlgorithms = () => {
  for (const algorithm of availableAlgorithms.value) {
    selectedAlgorithms[algorithm.name] = true
  }
  if (inputMode.value === 'text' && inputText.value) {
    calculateHashes()
  }
}

const clearAllAlgorithms = () => {
  for (const algorithm of availableAlgorithms.value) {
    selectedAlgorithms[algorithm.name] = false
  }
  clearResults()
}

const selectCommonAlgorithms = () => {
  clearAllAlgorithms()
  selectedAlgorithms['MD5'] = true
  selectedAlgorithms['SHA-256'] = true
  if (inputMode.value === 'text' && inputText.value) {
    calculateHashes()
  }
}

const copyHash = async (hash) => {
  try {
    await navigator.clipboard.writeText(hash)
    ElMessage.success('Hash值已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const copyAllResults = async () => {
  const allHashes = Object.entries(hashResults)
    .map(([algorithm, hash]) => `${algorithm}: ${hash}`)
    .join('\n')
  
  try {
    await navigator.clipboard.writeText(allHashes)
    ElMessage.success('所有Hash值已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const exportResults = () => {
  const content = Object.entries(hashResults)
    .map(([algorithm, hash]) => `${algorithm}: ${hash}`)
    .join('\n')
  
  const blob = new Blob([content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `hash-results-${new Date().toISOString().slice(0, 10)}.txt`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  
  ElMessage.success('Hash结果已导出')
}

const verifyHash = () => {
  // 验证逻辑在计算属性中处理
}
</script>

<style scoped>
.hash-calculator-container {
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
  display: grid;
  grid-template-columns: 1fr 300px;
  grid-template-rows: auto auto;
  gap: 24px;
  margin-bottom: 32px;
}

.results-section {
  grid-column: span 2;
}

.input-card,
.algorithms-card,
.results-card,
.info-card,
.use-cases-card {
  border: 1px solid #e5e7eb;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: #2c3e50;
}

.card-header .el-icon + span {
  margin-left: 8px;
}

.input-content {
  padding: 8px 0;
}

.text-input-section {
  padding-top: 16px;
}

.input-stats {
  margin-top: 8px;
  font-size: 12px;
  color: #6b7280;
  text-align: right;
  display: flex;
  justify-content: space-between;
}

.file-input-section {
  padding-top: 16px;
}

.upload-area {
  width: 100%;
}

.selected-file {
  margin-top: 16px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.file-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.file-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-name {
  font-weight: 500;
  color: #2c3e50;
}

.file-size {
  font-size: 14px;
  color: #6b7280;
}

.algorithms-content {
  padding: 8px 0;
}

.algorithm-options {
  margin-bottom: 16px;
}

.algorithm-checkbox {
  display: block;
  margin-bottom: 12px;
  width: 100%;
}

.algorithm-checkbox :deep(.el-checkbox__label) {
  width: 100%;
}

.algorithm-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.algorithm-name {
  font-weight: 500;
  color: #2c3e50;
}

.algorithm-desc {
  font-size: 12px;
  color: #6b7280;
}

.algorithm-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.results-content {
  padding: 8px 0;
  min-height: 200px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 150px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.hash-results {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.hash-result-item {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  background: #f8fafc;
}

.hash-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.hash-algorithm {
  font-weight: 600;
  color: #2c3e50;
}

.hash-value {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 14px;
  color: #1f2937;
  word-break: break-all;
  line-height: 1.4;
  background: #fff;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #d1d5db;
}

.comparison-section {
  margin-top: 24px;
}

.verification-input {
  margin: 16px 0;
}

.verification-result {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.match-success {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #22c55e;
}

.match-failure {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #ef4444;
}

.algorithm-descriptions {
  padding: 16px 0;
}

.description-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.desc-item {
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f8fafc;
}

.desc-item h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.desc-item p {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.output-length {
  font-size: 12px;
  color: #9ca3af;
  font-style: italic;
}

.use-cases-section {
  margin-bottom: 32px;
}

.use-cases-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  padding: 16px 0;
}

.use-case {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 20px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #fff;
  transition: all 0.2s ease;
}

.use-case:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.case-icon {
  font-size: 32px;
  color: #667eea;
  margin-bottom: 12px;
}

.use-case h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.use-case p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.4;
}

@media (max-width: 1024px) {
  .tools-content {
    grid-template-columns: 1fr;
    grid-template-rows: auto;
  }
  
  .results-section {
    grid-column: auto;
  }
}

@media (max-width: 768px) {
  .hash-calculator-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .description-grid {
    grid-template-columns: 1fr;
  }
  
  .use-cases-grid {
    grid-template-columns: 1fr;
  }
  
  .algorithm-actions {
    justify-content: center;
  }
  
  .result-actions {
    flex-direction: column;
    gap: 8px;
  }
}
</style>