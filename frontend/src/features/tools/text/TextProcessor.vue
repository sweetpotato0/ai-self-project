<template>
  <div class="text-processor-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-text' })">文本类</el-breadcrumb-item>
          <el-breadcrumb-item>文本处理器</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>文本处理器</h1>
      <p>文本格式化、大小写转换、去重等多功能处理工具</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <div class="section-header">
          <h3>输入文本</h3>
          <div class="action-buttons">
            <el-button @click="clearInput" size="small">清空</el-button>
            <el-button @click="pasteFromClipboard" size="small" type="primary">粘贴</el-button>
          </div>
        </div>
        <el-input
          v-model="inputText"
          type="textarea"
          :rows="8"
          placeholder="请输入要处理的文本..."
          @input="updateStats"
        />
        <div class="text-stats">
          <div class="stats-item">字符数: {{ stats.chars }}</div>
          <div class="stats-item">单词数: {{ stats.words }}</div>
          <div class="stats-item">行数: {{ stats.lines }}</div>
          <div class="stats-item">字节数: {{ stats.bytes }}</div>
        </div>
      </div>

      <div class="operations-section">
        <el-card class="operations-card">
          <template #header>
            <div class="card-header">
              <el-icon><Tools /></el-icon>
              <span>处理操作</span>
            </div>
          </template>
          <div class="operations-grid">
            <el-button @click="toUpperCase" type="primary" plain>全部大写</el-button>
            <el-button @click="toLowerCase" type="primary" plain>全部小写</el-button>
            <el-button @click="toTitleCase" type="primary" plain>标题格式</el-button>
            <el-button @click="toCamelCase" type="primary" plain>驼峰命名</el-button>
            <el-button @click="toSnakeCase" type="primary" plain>下划线命名</el-button>
            <el-button @click="toKebabCase" type="primary" plain>短横线命名</el-button>
            <el-button @click="removeDuplicateLines" type="success" plain>去重行</el-button>
            <el-button @click="sortLines" type="success" plain>排序行</el-button>
            <el-button @click="reverseLines" type="success" plain>反转行</el-button>
            <el-button @click="removeEmptyLines" type="warning" plain>删除空行</el-button>
            <el-button @click="trimLines" type="warning" plain>去首尾空格</el-button>
            <el-button @click="addLineNumbers" type="info" plain>添加行号</el-button>
            <el-button @click="removeLineNumbers" type="info" plain>删除行号</el-button>
            <el-button @click="extractEmails" type="danger" plain>提取邮箱</el-button>
            <el-button @click="extractUrls" type="danger" plain>提取URL</el-button>
            <el-button @click="extractNumbers" type="danger" plain>提取数字</el-button>
          </div>
        </el-card>
      </div>

      <div class="output-section">
        <div class="section-header">
          <h3>处理结果</h3>
          <div class="action-buttons">
            <el-button @click="copyResult" size="small" type="primary">
              <el-icon><DocumentCopy /></el-icon>
              复制
            </el-button>
            <el-button @click="clearOutput" size="small">清空</el-button>
          </div>
        </div>
        <el-input
          v-model="outputText"
          type="textarea"
          :rows="8"
          readonly
          placeholder="处理结果将在这里显示..."
        />
        <div class="text-stats" v-if="outputText">
          <div class="stats-item">字符数: {{ outputStats.chars }}</div>
          <div class="stats-item">单词数: {{ outputStats.words }}</div>
          <div class="stats-item">行数: {{ outputStats.lines }}</div>
          <div class="stats-item">字节数: {{ outputStats.bytes }}</div>
        </div>
      </div>

      <div class="examples-section">
        <el-card class="examples-card">
          <template #header>
            <div class="card-header">
              <el-icon><Collection /></el-icon>
              <span>使用示例</span>
            </div>
          </template>
          <div class="examples-list">
            <div 
              v-for="example in examples" 
              :key="example.id"
              class="example-item"
              @click="useExample(example)"
            >
              <div class="example-title">{{ example.title }}</div>
              <div class="example-text">{{ example.text }}</div>
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
import {
  DocumentCopy,
  Tools,
  Collection
} from '@element-plus/icons-vue'

const router = useRouter()
const inputText = ref('')
const outputText = ref('')

const examples = ref([
  {
    id: 1,
    title: '大小写转换示例',
    text: 'Hello World\nthis is a TEST\nJavaScript Programming'
  },
  {
    id: 2,
    title: '去重排序示例',
    text: 'apple\nbanana\napple\ncherry\nbanana\ndate'
  },
  {
    id: 3,
    title: '邮箱提取示例',
    text: '联系我们：support@example.com 或者 admin@test.org\n销售咨询：sales@company.net'
  },
  {
    id: 4,
    title: 'URL提取示例',
    text: '请访问 https://www.example.com 了解更多\n或者查看 http://blog.test.com/article/123'
  }
])

const stats = computed(() => {
  const text = inputText.value
  return {
    chars: text.length,
    words: text.trim() ? text.trim().split(/\s+/).length : 0,
    lines: text.split('\n').length,
    bytes: new Blob([text]).size
  }
})

const outputStats = computed(() => {
  const text = outputText.value
  return {
    chars: text.length,
    words: text.trim() ? text.trim().split(/\s+/).length : 0,
    lines: text.split('\n').length,
    bytes: new Blob([text]).size
  }
})

const updateStats = () => {
  // Stats are computed automatically
}

const toUpperCase = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  outputText.value = inputText.value.toUpperCase()
  ElMessage.success('转换为大写完成')
}

const toLowerCase = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  outputText.value = inputText.value.toLowerCase()
  ElMessage.success('转换为小写完成')
}

const toTitleCase = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  outputText.value = inputText.value.replace(/\w\S*/g, (txt) => 
    txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase()
  )
  ElMessage.success('转换为标题格式完成')
}

const toCamelCase = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  outputText.value = inputText.value
    .replace(/(?:^\w|[A-Z]|\b\w)/g, (word, index) => 
      index === 0 ? word.toLowerCase() : word.toUpperCase()
    )
    .replace(/\s+/g, '')
  ElMessage.success('转换为驼峰命名完成')
}

const toSnakeCase = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  outputText.value = inputText.value
    .replace(/\W+/g, ' ')
    .split(/ |\B(?=[A-Z])/)
    .map(word => word.toLowerCase())
    .join('_')
  ElMessage.success('转换为下划线命名完成')
}

const toKebabCase = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  outputText.value = inputText.value
    .replace(/\W+/g, ' ')
    .split(/ |\B(?=[A-Z])/)
    .map(word => word.toLowerCase())
    .join('-')
  ElMessage.success('转换为短横线命名完成')
}

const removeDuplicateLines = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const lines = inputText.value.split('\n')
  const uniqueLines = [...new Set(lines)]
  outputText.value = uniqueLines.join('\n')
  ElMessage.success(`去重完成，从 ${lines.length} 行减少到 ${uniqueLines.length} 行`)
}

const sortLines = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const lines = inputText.value.split('\n')
  const sortedLines = lines.sort()
  outputText.value = sortedLines.join('\n')
  ElMessage.success('行排序完成')
}

const reverseLines = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const lines = inputText.value.split('\n')
  const reversedLines = lines.reverse()
  outputText.value = reversedLines.join('\n')
  ElMessage.success('行反转完成')
}

const removeEmptyLines = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const lines = inputText.value.split('\n')
  const nonEmptyLines = lines.filter(line => line.trim() !== '')
  outputText.value = nonEmptyLines.join('\n')
  ElMessage.success(`删除空行完成，从 ${lines.length} 行减少到 ${nonEmptyLines.length} 行`)
}

const trimLines = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const lines = inputText.value.split('\n')
  const trimmedLines = lines.map(line => line.trim())
  outputText.value = trimmedLines.join('\n')
  ElMessage.success('去除首尾空格完成')
}

const addLineNumbers = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const lines = inputText.value.split('\n')
  const numberedLines = lines.map((line, index) => `${index + 1}. ${line}`)
  outputText.value = numberedLines.join('\n')
  ElMessage.success('添加行号完成')
}

const removeLineNumbers = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const lines = inputText.value.split('\n')
  const cleanLines = lines.map(line => line.replace(/^\d+\.\s*/, ''))
  outputText.value = cleanLines.join('\n')
  ElMessage.success('删除行号完成')
}

const extractEmails = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const emailRegex = /\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b/g
  const emails = inputText.value.match(emailRegex) || []
  outputText.value = [...new Set(emails)].join('\n')
  ElMessage.success(`提取到 ${emails.length} 个邮箱地址`)
}

const extractUrls = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const urlRegex = /https?:\/\/[^\s]+/g
  const urls = inputText.value.match(urlRegex) || []
  outputText.value = [...new Set(urls)].join('\n')
  ElMessage.success(`提取到 ${urls.length} 个URL`)
}

const extractNumbers = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要处理的文本')
    return
  }
  const numberRegex = /\d+\.?\d*/g
  const numbers = inputText.value.match(numberRegex) || []
  outputText.value = [...new Set(numbers)].join('\n')
  ElMessage.success(`提取到 ${numbers.length} 个数字`)
}

const clearInput = () => {
  inputText.value = ''
}

const clearOutput = () => {
  outputText.value = ''
}

const copyResult = async () => {
  if (!outputText.value) {
    ElMessage.warning('没有可复制的内容')
    return
  }
  
  try {
    await navigator.clipboard.writeText(outputText.value)
    ElMessage.success('复制成功')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const pasteFromClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    inputText.value = text
    ElMessage.success('粘贴成功')
  } catch (error) {
    ElMessage.error('粘贴失败')
  }
}

const useExample = (example) => {
  inputText.value = example.text
  outputText.value = ''
  ElMessage.success(`已加载示例: ${example.title}`)
}
</script>

<style scoped>
.text-processor-container {
  padding: 20px;
  max-width: 1200px;
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
  gap: 20px;
}

.input-section,
.output-section {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.text-stats {
  margin-top: 8px;
  display: flex;
  gap: 16px;
  justify-content: flex-end;
  flex-wrap: wrap;
}

.stats-item {
  font-size: 12px;
  color: #6b7280;
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 4px;
}

.operations-section {
  margin: 20px 0;
}

.operations-card {
  border: 1px solid #e5e7eb;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #2c3e50;
}

.operations-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}

.examples-section {
  margin-top: 20px;
}

.examples-card {
  border: 1px solid #e5e7eb;
}

.examples-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.example-item {
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.example-item:hover {
  border-color: #409eff;
  background-color: #f8fafc;
}

.example-title {
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.example-text {
  font-size: 12px;
  color: #6b7280;
  font-family: monospace;
  white-space: pre-line;
  line-height: 1.4;
}

@media (max-width: 768px) {
  .text-processor-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .action-buttons {
    align-self: flex-end;
  }
  
  .operations-grid {
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    gap: 8px;
  }
  
  .examples-list {
    grid-template-columns: 1fr;
  }
  
  .text-stats {
    justify-content: flex-start;
    gap: 8px;
  }
}
</style>