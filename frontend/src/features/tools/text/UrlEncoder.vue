<template>
  <div class="url-encoder-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-text' })">文本类</el-breadcrumb-item>
          <el-breadcrumb-item>URL编码</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>URL编码解码</h1>
      <p>对URL和文本进行编码和解码处理</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <div class="section-header">
          <h3>输入内容</h3>
          <div class="action-buttons">
            <el-button @click="clearInput" size="small">清空</el-button>
            <el-button @click="pasteFromClipboard" size="small" type="primary">粘贴</el-button>
          </div>
        </div>
        <el-input
          v-model="inputText"
          type="textarea"
          :rows="6"
          placeholder="请输入要编码或解码的URL或文本..."
          @input="handleInputChange"
        />
        <div class="text-stats">
          字符数: {{ inputText.length }}
        </div>
      </div>

      <div class="operation-section">
        <div class="operation-buttons">
          <el-button @click="encodeURL" type="primary" size="large">
            <el-icon><ArrowDown /></el-icon>
            URL编码
          </el-button>
          <el-button @click="decodeURL" type="success" size="large">
            <el-icon><ArrowUp /></el-icon>
            URL解码
          </el-button>
        </div>
        <div class="encoding-options">
          <el-radio-group v-model="encodingType">
            <el-radio label="component">encodeURIComponent (推荐)</el-radio>
            <el-radio label="uri">encodeURI</el-radio>
          </el-radio-group>
        </div>
      </div>

      <div class="output-section">
        <div class="section-header">
          <h3>输出结果</h3>
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
          :rows="6"
          readonly
          placeholder="编码或解码结果将在这里显示..."
        />
        <div class="text-stats" v-if="outputText">
          字符数: {{ outputText.length }}
        </div>
      </div>

      <div class="examples-section">
        <el-card class="examples-card">
          <template #header>
            <div class="card-header">
              <el-icon><Collection /></el-icon>
              <span>常用示例</span>
            </div>
          </template>
          <div class="examples-grid">
            <div 
              v-for="example in examples" 
              :key="example.id"
              class="example-item"
              @click="useExample(example)"
            >
              <div class="example-title">{{ example.title }}</div>
              <div class="example-original">原文: {{ example.original }}</div>
              <div class="example-encoded">编码: {{ example.encoded }}</div>
            </div>
          </div>
        </el-card>
      </div>

      <div class="info-section">
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>URL编码说明</span>
            </div>
          </template>
          <div class="info-content">
            <div class="info-item">
              <h4>encodeURIComponent</h4>
              <p>编码所有特殊字符，包括 : / ? # [ ] @，适用于URL参数值</p>
            </div>
            <div class="info-item">
              <h4>encodeURI</h4>
              <p>保留URL结构字符 : / ? # [ ] @，适用于完整URL编码</p>
            </div>
            <div class="common-chars">
              <h4>常见字符编码对照</h4>
              <div class="chars-grid">
                <div class="char-item">空格 → %20</div>
                <div class="char-item">中文 → %E4%B8%AD%E6%96%87</div>
                <div class="char-item">+ → %2B</div>
                <div class="char-item">& → %26</div>
                <div class="char-item">= → %3D</div>
                <div class="char-item"># → %23</div>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowDown,
  ArrowUp,
  DocumentCopy,
  InfoFilled,
  Collection
} from '@element-plus/icons-vue'

const router = useRouter()
const inputText = ref('')
const outputText = ref('')
const encodingType = ref('component')

const examples = ref([
  {
    id: 1,
    title: '中文URL',
    original: 'https://example.com/搜索?q=测试',
    encoded: 'https://example.com/%E6%90%9C%E7%B4%A2?q=%E6%B5%8B%E8%AF%95'
  },
  {
    id: 2,
    title: '包含空格的参数',
    original: 'name=John Doe&age=25',
    encoded: 'name=John%20Doe&age=25'
  },
  {
    id: 3,
    title: '特殊字符',
    original: 'search=hello+world&filter=A&B',
    encoded: 'search=hello%2Bworld&filter=A%26B'
  },
  {
    id: 4,
    title: '邮箱地址',
    original: 'email=user@example.com',
    encoded: 'email=user%40example.com'
  }
])

const encodeURL = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要编码的内容')
    return
  }
  
  try {
    let encoded
    if (encodingType.value === 'component') {
      encoded = encodeURIComponent(inputText.value)
    } else {
      encoded = encodeURI(inputText.value)
    }
    outputText.value = encoded
    ElMessage.success('URL编码成功')
  } catch (error) {
    ElMessage.error('编码失败: ' + error.message)
  }
}

const decodeURL = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要解码的URL')
    return
  }
  
  try {
    const decoded = decodeURIComponent(inputText.value)
    outputText.value = decoded
    ElMessage.success('URL解码成功')
  } catch (error) {
    ElMessage.error('解码失败，请检查输入是否为有效的URL编码格式')
  }
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

const handleInputChange = () => {
  outputText.value = ''
}

const useExample = (example) => {
  inputText.value = example.original
  outputText.value = ''
  ElMessage.success(`已加载示例: ${example.title}`)
}
</script>

<style scoped>
.url-encoder-container {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.tools-header {
  margin-bottom: 30px;
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
  font-size: 12px;
  color: #6b7280;
  text-align: right;
}

.operation-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 20px 0;
}

.operation-buttons {
  display: flex;
  gap: 20px;
}

.operation-buttons .el-button {
  padding: 12px 24px;
  font-size: 16px;
  font-weight: 500;
}

.encoding-options {
  display: flex;
  justify-content: center;
}

.examples-section,
.info-section {
  margin-top: 20px;
}

.examples-card,
.info-card {
  border: 1px solid #e5e7eb;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #2c3e50;
}

.examples-grid {
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

.example-original,
.example-encoded {
  font-size: 12px;
  color: #6b7280;
  margin: 4px 0;
  word-break: break-all;
}

.example-original {
  color: #059669;
}

.example-encoded {
  color: #dc2626;
}

.info-content {
  color: #374151;
  line-height: 1.6;
}

.info-item {
  margin-bottom: 16px;
}

.info-item h4 {
  color: #2c3e50;
  margin: 0 0 8px 0;
  font-size: 16px;
}

.common-chars h4 {
  color: #2c3e50;
  margin: 16px 0 12px 0;
}

.chars-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 8px;
}

.char-item {
  padding: 8px 12px;
  background: #f3f4f6;
  border-radius: 6px;
  font-family: monospace;
  font-size: 12px;
}

@media (max-width: 768px) {
  .url-encoder-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .operation-buttons {
    flex-direction: column;
    gap: 12px;
  }
  
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .action-buttons {
    align-self: flex-end;
  }
  
  .examples-grid {
    grid-template-columns: 1fr;
  }
  
  .chars-grid {
    grid-template-columns: 1fr;
  }
}
</style>