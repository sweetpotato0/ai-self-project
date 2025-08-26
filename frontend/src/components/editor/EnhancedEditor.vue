<template>
  <div class="enhanced-editor">
    <!-- 工具栏 -->
    <div class="editor-toolbar">
      <!-- 文本格式 -->
      <div class="toolbar-section">
        <span class="section-label">文本格式</span>
        <el-button-group>
          <el-button @click="execCommand('bold')" :class="{ active: isActive('bold') }" title="粗体 (Ctrl+B)">
            <el-icon><Star /></el-icon>
          </el-button>
          <el-button @click="execCommand('italic')" :class="{ active: isActive('italic') }" title="斜体 (Ctrl+I)">
            <el-icon><Check /></el-icon>
          </el-button>
          <el-button @click="execCommand('underline')" :class="{ active: isActive('underline') }" title="下划线 (Ctrl+U)">
            <el-icon><Close /></el-icon>
          </el-button>
        </el-button-group>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 列表 -->
      <div class="toolbar-section">
        <span class="section-label">列表</span>
        <el-button-group>
          <el-button @click="execCommand('insertUnorderedList')" :class="{ active: isActive('insertUnorderedList') }" title="无序列表">
            <el-icon><CircleCheck /></el-icon>
          </el-button>
          <el-button @click="execCommand('insertOrderedList')" :class="{ active: isActive('insertOrderedList') }" title="有序列表">
            <el-icon><Histogram /></el-icon>
          </el-button>
        </el-button-group>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 对齐方式 -->
      <div class="toolbar-section">
        <span class="section-label">对齐</span>
        <el-button-group>
          <el-button @click="execCommand('justifyLeft')" :class="{ active: isActive('justifyLeft') }" title="左对齐">
            <el-icon><ArrowLeft /></el-icon>
          </el-button>
          <el-button @click="execCommand('justifyCenter')" :class="{ active: isActive('justifyCenter') }" title="居中对齐">
            <el-icon><ArrowUp /></el-icon>
          </el-button>
          <el-button @click="execCommand('justifyRight')" :class="{ active: isActive('justifyRight') }" title="右对齐">
            <el-icon><ArrowRight /></el-icon>
          </el-button>
        </el-button-group>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 插入元素 -->
      <div class="toolbar-section">
        <span class="section-label">插入</span>
        <el-button-group>
          <el-button @click="insertImage" title="插入图片">
            <el-icon><Picture /></el-icon>
          </el-button>
          <el-button @click="insertLink" title="插入链接">
            <el-icon><Link /></el-icon>
          </el-button>
          <el-button @click="insertTable" title="插入表格">
            <el-icon><Grid /></el-icon>
          </el-button>
          <el-button @click="insertCodeBlock" title="插入代码块">
            <el-icon><Document /></el-icon>
          </el-button>
          <el-button @click="insertQuote" title="插入引用">
            <el-icon><Edit /></el-icon>
          </el-button>
        </el-button-group>
      </div>

      <div class="toolbar-divider"></div>

      <!-- 视图切换 -->
      <div class="toolbar-section">
        <el-button @click="togglePreview" :class="{ active: previewMode }" title="预览模式">
          <el-icon><View /></el-icon>
          {{ previewMode ? '编辑' : '预览' }}
        </el-button>
      </div>
    </div>

    <!-- 编辑器主体 -->
    <div class="editor-main">
      <div
        v-if="!previewMode"
        ref="editorRef"
        class="editor-content"
        contenteditable="true"
        @input="handleInput"
        @paste="handlePaste"
        @drop="handleDrop"
        @dragover="handleDragOver"
        @keydown="handleKeydown"
        v-html="modelValue"
      ></div>

      <div v-else class="editor-preview">
        <div class="preview-content" v-html="renderedContent"></div>
      </div>
    </div>

    <!-- 图片上传对话框 -->
    <el-dialog v-model="imageDialogVisible" title="插入图片" width="500px">
      <div class="image-upload-section">
        <el-upload
          ref="uploadRef"
          :action="uploadUrl"
          :headers="uploadHeaders"
          :before-upload="beforeImageUpload"
          :on-success="handleImageSuccess"
          :on-error="handleImageError"
          :show-file-list="false"
          accept="image/*"
          drag
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将图片拖到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持 jpg/png/gif 格式，文件大小不超过 5MB
            </div>
          </template>
        </el-upload>
      </div>
    </el-dialog>

    <!-- 链接插入对话框 -->
    <el-dialog v-model="linkDialogVisible" title="插入链接" width="400px">
      <el-form :model="linkForm" label-width="80px">
        <el-form-item label="链接地址">
          <el-input v-model="linkForm.url" placeholder="请输入链接地址" />
        </el-form-item>
        <el-form-item label="链接文本">
          <el-input v-model="linkForm.text" placeholder="请输入链接文本" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="linkDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="insertLinkConfirm">确定</el-button>
      </template>
    </el-dialog>

    <!-- 表格插入对话框 -->
    <el-dialog v-model="tableDialogVisible" title="插入表格" width="400px">
      <el-form :model="tableForm" label-width="80px">
        <el-form-item label="行数">
          <el-input-number v-model="tableForm.rows" :min="1" :max="20" />
        </el-form-item>
        <el-form-item label="列数">
          <el-input-number v-model="tableForm.cols" :min="1" :max="10" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tableDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="insertTableConfirm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document, List, Edit, View, UploadFilled, Grid,
  Link, Picture, Setting, Plus, Calendar, TrendCharts,
  CircleCheck, Histogram,
  Star, Check, Close, ArrowRight, ArrowLeft, ArrowUp, ArrowDown
} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '请输入内容...'
  }
})

const emit = defineEmits(['update:modelValue'])

// 响应式数据
const editorRef = ref()
const previewMode = ref(false)
const imageDialogVisible = ref(false)
const linkDialogVisible = ref(false)
const tableDialogVisible = ref(false)

// 表单数据
const linkForm = ref({
  url: '',
  text: ''
})

const tableForm = ref({
  rows: 3,
  cols: 3
})

// 上传配置
const uploadUrl = 'http://localhost:8080/api/v1/upload/image'
const uploadHeaders = computed(() => ({
  'Authorization': `Bearer ${localStorage.getItem('token')}`
}))

// 计算属性
const renderedContent = computed(() => {
  return props.modelValue
})

// 方法
const execCommand = (command, value = null) => {
  document.execCommand(command, false, value)
  editorRef.value?.focus()
}

const isActive = (command, value = null) => {
  return document.queryCommandState(command)
}

const handleInput = () => {
  if (editorRef.value) {
    emit('update:modelValue', editorRef.value.innerHTML)
  }
}

const handlePaste = (event) => {
  event.preventDefault()
  const text = event.clipboardData.getData('text/plain')
  document.execCommand('insertText', false, text)
}

const handleDrop = (event) => {
  event.preventDefault()
  const files = event.dataTransfer.files
  if (files.length > 0 && files[0].type.startsWith('image/')) {
    uploadImage(files[0])
  }
}

const handleDragOver = (event) => {
  event.preventDefault()
}

const handleKeydown = (event) => {
  // 支持快捷键
  if (event.ctrlKey || event.metaKey) {
    switch (event.key) {
      case 'b':
        event.preventDefault()
        execCommand('bold')
        break
      case 'i':
        event.preventDefault()
        execCommand('italic')
        break
      case 'u':
        event.preventDefault()
        execCommand('underline')
        break
    }
  }
}

const togglePreview = () => {
  previewMode.value = !previewMode.value
}

const insertLink = () => {
  const selection = window.getSelection()
  if (selection.toString()) {
    linkForm.value.text = selection.toString()
  }
  linkDialogVisible.value = true
}

const insertLinkConfirm = () => {
  if (linkForm.value.url && linkForm.value.text) {
    const linkHtml = `<a href="${linkForm.value.url}" target="_blank">${linkForm.value.text}</a>`
    document.execCommand('insertHTML', false, linkHtml)
    linkDialogVisible.value = false
    linkForm.value = { url: '', text: '' }
  } else {
    ElMessage.warning('请输入链接地址和文本')
  }
}

const insertImage = () => {
  imageDialogVisible.value = true
}

const insertTable = () => {
  tableDialogVisible.value = true
}

const insertTableConfirm = () => {
  const { rows, cols } = tableForm.value
  let tableHtml = '<table border="1" style="border-collapse: collapse; width: 100%; margin: 10px 0;">'

  for (let i = 0; i < rows; i++) {
    tableHtml += '<tr>'
    for (let j = 0; j < cols; j++) {
      tableHtml += '<td style="padding: 8px; border: 1px solid #ddd;">&nbsp;</td>'
    }
    tableHtml += '</tr>'
  }
  tableHtml += '</table>'

  document.execCommand('insertHTML', false, tableHtml)
  tableDialogVisible.value = false
  tableForm.value = { rows: 3, cols: 3 }
}

const insertCodeBlock = () => {
  const codeHtml = `<pre style="background: #2d3748; color: #e2e8f0; padding: 20px; border-radius: 8px; overflow-x: auto; margin: 15px 0; font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace; font-size: 13px; line-height: 1.6;"><code>// 在这里输入代码</code></pre>`
  document.execCommand('insertHTML', false, codeHtml)
}

const insertQuote = () => {
  const quoteHtml = `<blockquote style="margin: 15px 0; padding: 15px 20px; background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%); border-left: 4px solid #409eff; border-radius: 0 8px 8px 0; font-style: italic; color: #6c757d;">在这里输入引用内容</blockquote>`
  document.execCommand('insertHTML', false, quoteHtml)
}

const beforeImageUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  return true
}

const handleImageSuccess = (response) => {
  if (response.code === 200) {
    const imageUrl = response.data.url
    const imageHtml = `<img src="${imageUrl}" alt="图片" style="max-width: 100%; height: auto; margin: 10px 0; border-radius: 4px;" />`
    document.execCommand('insertHTML', false, imageHtml)
    imageDialogVisible.value = false
    ElMessage.success('图片上传成功')
  } else {
    ElMessage.error('图片上传失败')
  }
}

const handleImageError = () => {
  ElMessage.error('图片上传失败')
}

const uploadImage = async (file) => {
  const formData = new FormData()
  formData.append('file', file)

  try {
    const response = await fetch(uploadUrl, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: formData
    })

    const result = await response.json()
    if (result.code === 200) {
      const imageUrl = result.data.url
      const imageHtml = `<img src="${imageUrl}" alt="图片" style="max-width: 100%; height: auto; margin: 10px 0; border-radius: 4px;" />`
      document.execCommand('insertHTML', false, imageHtml)
      ElMessage.success('图片上传成功')
    } else {
      ElMessage.error('图片上传失败')
    }
  } catch (error) {
    console.error('Upload error:', error)
    ElMessage.error('图片上传失败')
  }
}

// 监听modelValue变化
watch(() => props.modelValue, (newValue) => {
  if (editorRef.value && !previewMode.value) {
    editorRef.value.innerHTML = newValue
  }
}, { immediate: true })
</script>

<style scoped>
.enhanced-editor {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-bottom: 1px solid #dcdfe6;
  flex-wrap: wrap;
}

.toolbar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.section-label {
  font-size: 11px;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 2px;
}

.toolbar-divider {
  width: 1px;
  height: 40px;
  background: linear-gradient(to bottom, transparent, #ddd, transparent);
  margin: 0 5px;
}

.toolbar-section .el-button {
  padding: 8px 12px;
  font-size: 14px;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.toolbar-section .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.toolbar-section .el-button.active {
  background: linear-gradient(135deg, #409eff 0%, #36a3f7 100%);
  color: white;
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.editor-main {
  min-height: 400px;
  background: white;
}

.editor-content {
  min-height: 400px;
  padding: 25px;
  outline: none;
  line-height: 1.8;
  font-size: 14px;
  color: #333;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.editor-content:empty:before {
  content: attr(placeholder);
  color: #c0c4cc;
  font-style: italic;
}

.editor-content:focus {
  background: #fafbfc;
}

.editor-preview {
  padding: 25px;
  background: #fafbfc;
  min-height: 400px;
  border-left: 4px solid #409eff;
}

.preview-content {
  line-height: 1.8;
  font-size: 14px;
  color: #333;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.preview-content h1,
.preview-content h2,
.preview-content h3,
.preview-content h4,
.preview-content h5,
.preview-content h6 {
  margin: 20px 0 12px 0;
  font-weight: 600;
  color: #2c3e50;
  border-bottom: 2px solid #ecf0f1;
  padding-bottom: 8px;
}

.preview-content h1 { font-size: 28px; }
.preview-content h2 { font-size: 24px; }
.preview-content h3 { font-size: 20px; }
.preview-content h4 { font-size: 18px; }
.preview-content h5 { font-size: 16px; }
.preview-content h6 { font-size: 14px; }

.preview-content p {
  margin: 12px 0;
  text-align: justify;
}

.preview-content ul,
.preview-content ol {
  margin: 12px 0;
  padding-left: 25px;
}

.preview-content li {
  margin: 6px 0;
}

.preview-content img {
  max-width: 100%;
  height: auto;
  margin: 15px 0;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.preview-content img:hover {
  transform: scale(1.02);
}

.preview-content table {
  width: 100%;
  border-collapse: collapse;
  margin: 15px 0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.preview-content th,
.preview-content td {
  padding: 12px;
  border: 1px solid #e4e7ed;
  text-align: left;
}

.preview-content th {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  font-weight: 600;
  color: #2c3e50;
}

.preview-content tr:nth-child(even) {
  background-color: #fafbfc;
}

.preview-content tr:hover {
  background-color: #f0f2f5;
}

.preview-content a {
  color: #409eff;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-bottom 0.3s ease;
}

.preview-content a:hover {
  border-bottom: 1px solid #409eff;
}

.preview-content blockquote {
  margin: 15px 0;
  padding: 15px 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-left: 4px solid #409eff;
  border-radius: 0 8px 8px 0;
  font-style: italic;
  color: #6c757d;
}

.preview-content code {
  background: #f8f9fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  color: #e83e8c;
}

.preview-content pre {
  background: #2d3748;
  color: #e2e8f0;
  padding: 20px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 15px 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
}

.image-upload-section {
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .editor-toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }

  .toolbar-section {
    justify-content: center;
  }

  .editor-content,
  .editor-preview {
    padding: 15px;
  }
}
</style>
