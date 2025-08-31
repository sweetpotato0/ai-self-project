<template>
  <div class="markdown-editor">
    <Editor
      ref="editorRef"
      :value="modelValue"
      :plugins="plugins"
      :locale="locale"
      @change="handleChange"
      :uploadImages="handleImageUpload"
      class="bytemd-editor"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Editor } from '@bytemd/vue-next'
import gfm from '@bytemd/plugin-gfm'
import highlight from '@bytemd/plugin-highlight'
import math from '@bytemd/plugin-math'
import mermaid from '@bytemd/plugin-mermaid'
import { ElMessage } from 'element-plus'

// 导入必要的 CSS 样式
import 'bytemd/dist/index.css'
import 'github-markdown-css/github-markdown.css'
import 'highlight.js/styles/github.css'
import 'katex/dist/katex.css'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: String,
    default: '500px'
  },
  placeholder: {
    type: String,
    default: '请使用 Markdown 语法编写你的文章...'
  }
})

const emit = defineEmits(['update:modelValue'])

// 编辑器引用
const editorRef = ref()

// 配置插件
const plugins = [
  gfm({
    locale: 'zh-Hans'  // 支持中文
  }),
  highlight(),
  math({
    katexOptions: {
      throwOnError: false  // 数学公式错误时不抛出异常
    }
  }),
  mermaid({
    locale: 'zh-Hans'
  })
]

// 中文语言包
const locale = {
  // 工具栏
  h1: '一级标题',
  h2: '二级标题',
  h3: '三级标题',
  bold: '粗体',
  italic: '斜体',
  quote: '引用',
  link: '链接',
  image: '图片',
  code: '代码',
  codeBlock: '代码块',
  ul: '无序列表',
  ol: '有序列表',
  table: '表格',
  hr: '分隔线',
  
  // 按钮
  writeMode: '编辑模式',
  previewMode: '预览模式',
  splitMode: '分屏模式',
  fullscreen: '全屏',
  exitFullscreen: '退出全屏',
  toc: '目录',
  
  // 上传
  uploadImage: '上传图片',
  uploadFailed: '上传失败',
  
  // 其他
  cheatsheet: '速查表',
  limited: '内容长度超出限制'
}

// 处理内容变化
const handleChange = (value) => {
  emit('update:modelValue', value)
}

// 处理图片上传
const handleImageUpload = async (files) => {
  const uploadResults = []
  
  for (const file of files) {
    try {
      // 验证文件类型和大小
      if (!file.type.startsWith('image/')) {
        ElMessage.error('只能上传图片文件!')
        continue
      }
      
      if (file.size > 5 * 1024 * 1024) { // 5MB
        ElMessage.error('图片大小不能超过 5MB!')
        continue
      }

      // 创建 FormData 并上传
      const formData = new FormData()
      formData.append('file', file)
      
      const response = await fetch('http://localhost:8080/api/v1/upload/image', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: formData
      })

      const result = await response.json()
      
      if (result.code === 200) {
        uploadResults.push({
          url: result.data.url,
          alt: file.name,
          title: file.name
        })
      } else {
        throw new Error(result.message || '上传失败')
      }
    } catch (error) {
      console.error('图片上传失败:', error)
      ElMessage.error(`图片 ${file.name} 上传失败: ${error.message}`)
    }
  }

  return uploadResults
}

// 获取编辑器实例的方法（供外部调用）
const getEditor = () => {
  return editorRef.value
}

// 设置编辑器焦点
const focus = () => {
  if (editorRef.value) {
    editorRef.value.focus()
  }
}

// 插入内容到编辑器
const insertText = (text) => {
  if (editorRef.value) {
    editorRef.value.insertText(text)
  }
}

// 暴露方法给父组件
defineExpose({
  getEditor,
  focus,
  insertText
})

onMounted(() => {
  // 可以在这里进行一些初始化操作
})
</script>

<style scoped>
.markdown-editor {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

:deep(.bytemd) {
  height: v-bind(height);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
}

:deep(.bytemd-toolbar) {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-bottom: 1px solid #e4e7ed;
  padding: 8px 16px;
}

:deep(.bytemd-toolbar-icon) {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  transition: all 0.3s ease;
}

:deep(.bytemd-toolbar-icon):hover {
  background-color: rgba(64, 158, 255, 0.1);
  color: #409eff;
  transform: translateY(-1px);
}

:deep(.bytemd-toolbar-icon.bytemd-toolbar-icon-active) {
  background-color: #409eff;
  color: white;
}

:deep(.bytemd-body) {
  background: #fafbfc;
}

:deep(.CodeMirror) {
  height: 100%;
  font-size: 14px;
  line-height: 1.6;
  font-family: 'JetBrains Mono', 'Fira Code', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

:deep(.CodeMirror-focused .CodeMirror-cursor) {
  border-left: 2px solid #409eff;
}

:deep(.CodeMirror-selected) {
  background-color: rgba(64, 158, 255, 0.1);
}

:deep(.bytemd-preview) {
  padding: 20px;
  line-height: 1.8;
  color: #2c3e50;
}

:deep(.bytemd-preview h1),
:deep(.bytemd-preview h2),
:deep(.bytemd-preview h3),
:deep(.bytemd-preview h4),
:deep(.bytemd-preview h5),
:deep(.bytemd-preview h6) {
  color: #2c3e50;
  margin: 24px 0 16px;
  font-weight: 600;
  border-bottom: 1px solid #eaecef;
  padding-bottom: 8px;
}

:deep(.bytemd-preview h1) { font-size: 28px; }
:deep(.bytemd-preview h2) { font-size: 24px; }
:deep(.bytemd-preview h3) { font-size: 20px; }
:deep(.bytemd-preview h4) { font-size: 18px; }
:deep(.bytemd-preview h5) { font-size: 16px; }
:deep(.bytemd-preview h6) { font-size: 14px; }

:deep(.bytemd-preview p) {
  margin: 16px 0;
  text-align: justify;
}

:deep(.bytemd-preview blockquote) {
  margin: 16px 0;
  padding: 16px 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-left: 4px solid #409eff;
  border-radius: 0 8px 8px 0;
  color: #6c757d;
}

:deep(.bytemd-preview code) {
  background: #f8f9fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
  color: #e83e8c;
  font-family: 'JetBrains Mono', 'Fira Code', 'Monaco', 'Menlo', monospace;
}

:deep(.bytemd-preview pre) {
  background: #f6f8fa;
  border: none;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 16px 0;
  font-size: 13px;
  line-height: 1.6;
  font-family: 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
}

:deep(.bytemd-preview pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
  border-radius: 0;
  font-size: inherit;
}

/* 代码高亮样式增强 */
:deep(.hljs) {
  display: block;
  overflow-x: auto;
  padding: 0;
  color: #24292e;
  background: transparent;
}

:deep(.hljs-comment),
:deep(.hljs-quote) {
  color: #6a737d;
  font-style: italic;
}

:deep(.hljs-keyword),
:deep(.hljs-selector-tag),
:deep(.hljs-subst) {
  color: #d73a49;
  font-weight: bold;
}

:deep(.hljs-number),
:deep(.hljs-literal),
:deep(.hljs-variable),
:deep(.hljs-template-variable),
:deep(.hljs-tag .hljs-attr) {
  color: #005cc5;
}

:deep(.hljs-string),
:deep(.hljs-doctag) {
  color: #032f62;
}

:deep(.hljs-title),
:deep(.hljs-section),
:deep(.hljs-selector-id) {
  color: #6f42c1;
  font-weight: bold;
}

:deep(.hljs-class .hljs-title) {
  color: #6f42c1;
}

:deep(.hljs-tag),
:deep(.hljs-name),
:deep(.hljs-attribute) {
  color: #22863a;
  font-weight: normal;
}

:deep(.hljs-regexp),
:deep(.hljs-link) {
  color: #032f62;
}

:deep(.hljs-symbol),
:deep(.hljs-bullet) {
  color: #e36209;
}

:deep(.hljs-built_in),
:deep(.hljs-builtin-name) {
  color: #005cc5;
}

:deep(.hljs-meta) {
  color: #6a737d;
}

:deep(.hljs-deletion) {
  background: #ffeef0;
}

:deep(.hljs-addition) {
  background: #f0fff4;
}

:deep(.hljs-emphasis) {
  font-style: italic;
}

:deep(.hljs-strong) {
  font-weight: bold;
}

:deep(.bytemd-preview table) {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

:deep(.bytemd-preview th),
:deep(.bytemd-preview td) {
  border: none;
  border-right: 1px solid #e2e8f0;
  border-bottom: 1px solid #e2e8f0;
  padding: 12px 16px;
  text-align: left;
}

:deep(.bytemd-preview th:last-child),
:deep(.bytemd-preview td:last-child) {
  border-right: none;
}

:deep(.bytemd-preview tr:last-child td) {
  border-bottom: none;
}

:deep(.bytemd-preview th) {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  font-weight: 600;
  color: #2c3e50;
}

:deep(.bytemd-preview tr:nth-child(even)) {
  background-color: #fafbfc;
}

:deep(.bytemd-preview tr:hover) {
  background-color: #f0f2f5;
}

:deep(.bytemd-preview a) {
  color: #409eff;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-bottom 0.3s ease;
}

:deep(.bytemd-preview a:hover) {
  border-bottom: 1px solid #409eff;
}

:deep(.bytemd-preview img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  margin: 16px 0;
  display: block;
  margin-left: auto;
  margin-right: auto;
}

:deep(.bytemd-preview ul),
:deep(.bytemd-preview ol) {
  padding-left: 24px;
  margin: 16px 0;
}

:deep(.bytemd-preview li) {
  margin: 8px 0;
}

/* 数学公式样式 */
:deep(.katex) {
  font-size: 1.1em;
}

:deep(.katex-display) {
  margin: 20px 0;
  text-align: center;
}

/* Mermaid 图表样式 */
:deep(.mermaid) {
  text-align: center;
  margin: 20px 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  :deep(.bytemd) {
    height: 400px;
  }
  
  :deep(.bytemd-toolbar) {
    padding: 6px 12px;
  }
  
  :deep(.bytemd-toolbar-icon) {
    width: 28px;
    height: 28px;
  }
  
  :deep(.bytemd-preview) {
    padding: 16px;
  }
}

/* 深色主题支持 */
@media (prefers-color-scheme: dark) {
  :deep(.bytemd) {
    --bytemd-color-border: #484c56;
    --bytemd-color-bg: #1e1e1e;
    --bytemd-color-text: #d4d4d4;
  }
}
</style>