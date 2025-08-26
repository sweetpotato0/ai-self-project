<template>
  <div class="base64-encoder-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-text' })">文本类</el-breadcrumb-item>
          <el-breadcrumb-item>Base64编码</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>Base64编码解码</h1>
      <p>对文本进行Base64编码和解码处理</p>
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
          :rows="8"
          placeholder="请输入要编码或解码的文本..."
          @input="handleInputChange"
        />
        <div class="text-stats">
          字符数: {{ inputText.length }} | 字节数: {{ getByteLength(inputText) }}
        </div>
      </div>

      <div class="operation-section">
        <div class="operation-buttons">
          <el-button @click="encodeToBase64" type="primary" size="large">
            <el-icon><ArrowDown /></el-icon>
            编码 (Base64)
          </el-button>
          <el-button @click="decodeFromBase64" type="success" size="large">
            <el-icon><ArrowUp /></el-icon>
            解码 (文本)
          </el-button>
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
          :rows="8"
          readonly
          placeholder="编码或解码结果将在这里显示..."
        />
        <div class="text-stats" v-if="outputText">
          字符数: {{ outputText.length }} | 字节数: {{ getByteLength(outputText) }}
        </div>
      </div>

      <div class="info-section">
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>Base64 说明</span>
            </div>
          </template>
          <div class="info-content">
            <p><strong>Base64编码</strong>是一种基于64个可打印字符来表示二进制数据的表示方法。</p>
            <ul>
              <li><strong>用途</strong>: 在HTTP环境下传递较长的标识信息、邮件传输、数据存储等</li>
              <li><strong>特点</strong>: 编码后的数据比原始数据略长（约4/3倍）</li>
              <li><strong>字符集</strong>: A-Z, a-z, 0-9, +, / 和填充字符 =</li>
              <li><strong>安全性</strong>: Base64不是加密算法，仅是编码方式，不具备安全性</li>
            </ul>
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
  ArrowDown,
  ArrowUp,
  DocumentCopy,
  InfoFilled
} from '@element-plus/icons-vue'

const router = useRouter()
const inputText = ref('')
const outputText = ref('')

const getByteLength = (str) => {
  return new Blob([str]).size
}

const encodeToBase64 = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要编码的文本')
    return
  }
  
  try {
    const encoded = btoa(unescape(encodeURIComponent(inputText.value)))
    outputText.value = encoded
    ElMessage.success('Base64编码成功')
  } catch (error) {
    ElMessage.error('编码失败: ' + error.message)
  }
}

const decodeFromBase64 = () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要解码的Base64文本')
    return
  }
  
  try {
    const decoded = decodeURIComponent(escape(atob(inputText.value)))
    outputText.value = decoded
    ElMessage.success('Base64解码成功')
  } catch (error) {
    ElMessage.error('解码失败，请检查输入是否为有效的Base64格式')
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
  // 清空之前的输出
  outputText.value = ''
}
</script>

<style scoped>
.base64-encoder-container {
  padding: 20px;
  max-width: 1000px;
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
  font-size: 12px;
  color: #6b7280;
  text-align: right;
}

.operation-section {
  display: flex;
  justify-content: center;
  align-items: center;
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

.info-section {
  margin-top: 20px;
}

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

.info-content {
  color: #374151;
  line-height: 1.6;
}

.info-content ul {
  margin: 16px 0;
  padding-left: 20px;
}

.info-content li {
  margin: 8px 0;
}

@media (max-width: 768px) {
  .base64-encoder-container {
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
}
</style>