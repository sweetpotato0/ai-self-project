<template>
  <div class="result-container" v-if="visible">
    <div class="result-header">
      <h3>{{ title }}</h3>
      <div class="result-actions">
        <el-button 
          v-if="copyable" 
          type="primary" 
          size="small" 
          @click="copyResult"
          :icon="Copy"
        >
          复制
        </el-button>
        <el-button 
          v-if="downloadable" 
          type="success" 
          size="small" 
          @click="downloadResult"
          :icon="Download"
        >
          下载
        </el-button>
        <el-button 
          v-if="clearable" 
          type="default" 
          size="small" 
          @click="clearResult"
          :icon="Delete"
        >
          清空
        </el-button>
      </div>
    </div>
    
    <div class="result-content" :class="{ 'scrollable': scrollable }">
      <slot>
        <pre v-if="type === 'text'">{{ content }}</pre>
        <div v-else-if="type === 'html'" v-html="content"></div>
        <img v-else-if="type === 'image'" :src="content" :alt="title" />
        <div v-else>{{ content }}</div>
      </slot>
    </div>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import { Copy, Download, Delete } from '@element-plus/icons-vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: true
  },
  title: {
    type: String,
    default: '结果'
  },
  content: {
    type: [String, Object],
    default: ''
  },
  type: {
    type: String,
    default: 'text',
    validator: (value) => ['text', 'html', 'image', 'json'].includes(value)
  },
  copyable: {
    type: Boolean,
    default: true
  },
  downloadable: {
    type: Boolean,
    default: false
  },
  clearable: {
    type: Boolean,
    default: false
  },
  scrollable: {
    type: Boolean,
    default: true
  },
  filename: {
    type: String,
    default: 'result.txt'
  }
})

const emit = defineEmits(['copy', 'download', 'clear'])

const copyResult = () => {
  const textToCopy = props.type === 'json' ? JSON.stringify(props.content, null, 2) : String(props.content)
  navigator.clipboard.writeText(textToCopy).then(() => {
    ElMessage.success('复制成功！')
    emit('copy', textToCopy)
  }).catch(() => {
    ElMessage.error('复制失败！')
  })
}

const downloadResult = () => {
  const content = props.type === 'json' ? JSON.stringify(props.content, null, 2) : String(props.content)
  const blob = new Blob([content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = props.filename
  a.click()
  URL.revokeObjectURL(url)
  ElMessage.success('下载成功！')
  emit('download', content)
}

const clearResult = () => {
  emit('clear')
}
</script>

<style scoped>
.result-container {
  margin-top: 20px;
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
}

.result-header {
  background: #f5f7fa;
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #dcdfe6;
}

.result-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.result-actions {
  display: flex;
  gap: 8px;
}

.result-content {
  padding: 16px;
  background: white;
  min-height: 100px;
}

.result-content.scrollable {
  max-height: 400px;
  overflow-y: auto;
}

.result-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.result-content img {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .result-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .result-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>