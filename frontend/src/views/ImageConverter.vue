<template>
  <div class="image-converter-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-image' })">图像类</el-breadcrumb-item>
          <el-breadcrumb-item>格式转换</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>图片格式转换</h1>
      <p>支持JPG、PNG、WebP等格式相互转换</p>
    </div>

    <div class="tools-content">
      <div class="upload-section">
        <el-upload
          ref="uploadRef"
          class="upload-container"
          drag
          :auto-upload="false"
          :on-change="handleFileChange"
          :before-upload="beforeUpload"
          accept="image/*"
          :file-list="fileList"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将图片拖拽到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持 JPG、PNG、WebP、BMP、GIF 等格式
            </div>
          </template>
        </el-upload>
      </div>

      <div v-if="originalImage" class="conversion-section">
        <div class="format-selector">
          <h3>选择输出格式</h3>
          <div class="format-options">
            <div 
              v-for="format in formats" 
              :key="format.value"
              class="format-option"
              :class="{ active: selectedFormat === format.value }"
              @click="selectFormat(format.value)"
            >
              <div class="format-icon">{{ format.icon }}</div>
              <div class="format-info">
                <div class="format-name">{{ format.name }}</div>
                <div class="format-desc">{{ format.description }}</div>
              </div>
            </div>
          </div>
        </div>

        <div class="preview-section">
          <div class="image-preview">
            <div class="preview-panel">
              <div class="panel-header">
                <h3>原图预览</h3>
                <div class="image-info">
                  <span>格式: {{ originalFormat }}</span>
                  <span>大小: {{ originalSize }}</span>
                  <span>尺寸: {{ imageDimensions }}</span>
                </div>
              </div>
              <div class="image-container">
                <img :src="originalImage" alt="原图" />
              </div>
            </div>

            <div class="preview-panel" v-if="convertedImage">
              <div class="panel-header">
                <h3>转换后</h3>
                <div class="image-info">
                  <span>格式: {{ selectedFormat.toUpperCase() }}</span>
                  <span>大小: {{ convertedSize }}</span>
                  <span>变化: {{ sizeChange }}</span>
                </div>
              </div>
              <div class="image-container">
                <img :src="convertedImage" alt="转换后" />
              </div>
            </div>
          </div>

          <div class="action-buttons">
            <el-button 
              @click="convertImage" 
              type="primary" 
              size="large"
              :loading="converting"
              :disabled="!selectedFormat"
            >
              <el-icon><Refresh /></el-icon>
              {{ converting ? '转换中...' : '开始转换' }}
            </el-button>
            
            <el-button 
              v-if="convertedImage" 
              @click="downloadImage" 
              type="success" 
              size="large"
            >
              <el-icon><Download /></el-icon>
              下载图片
            </el-button>
            
            <el-button @click="resetConverter" type="default" size="large">
              <el-icon><RefreshLeft /></el-icon>
              重置
            </el-button>
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
  UploadFilled,
  Refresh,
  Download,
  RefreshLeft
} from '@element-plus/icons-vue'

const router = useRouter()

const uploadRef = ref()
const fileList = ref([])
const originalFile = ref(null)
const originalImage = ref('')
const convertedImage = ref('')
const selectedFormat = ref('')
const converting = ref(false)
const originalFileSize = ref(0)
const convertedFileSize = ref(0)
const originalFormat = ref('')
const imageDimensions = ref('')

const formats = ref([
  {
    value: 'jpeg',
    name: 'JPEG',
    icon: '🖼️',
    description: '适合照片，压缩率高'
  },
  {
    value: 'png',
    name: 'PNG',
    icon: '🎨',
    description: '支持透明，无损压缩'
  },
  {
    value: 'webp',
    name: 'WebP',
    icon: '🚀',
    description: '现代格式，高压缩率'
  },
  {
    value: 'bmp',
    name: 'BMP',
    icon: '📷',
    description: '位图格式，无压缩'
  }
])

const originalSize = computed(() => formatFileSize(originalFileSize.value))
const convertedSize = computed(() => formatFileSize(convertedFileSize.value))

const sizeChange = computed(() => {
  if (originalFileSize.value && convertedFileSize.value) {
    const diff = convertedFileSize.value - originalFileSize.value
    const percent = (diff / originalFileSize.value * 100).toFixed(1)
    return diff > 0 ? `+${percent}%` : `${percent}%`
  }
  return ''
})

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const beforeUpload = (file) => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请上传图片文件')
    return false
  }
  
  const isValidSize = file.size / 1024 / 1024 < 10
  if (!isValidSize) {
    ElMessage.error('图片大小不能超过 10MB')
    return false
  }
  
  return true
}

const handleFileChange = (file) => {
  if (!beforeUpload(file.raw)) return
  
  originalFile.value = file.raw
  originalFileSize.value = file.raw.size
  originalFormat.value = file.raw.type.split('/')[1].toUpperCase()
  
  const reader = new FileReader()
  reader.onload = (e) => {
    originalImage.value = e.target.result
    
    const img = new Image()
    img.onload = () => {
      imageDimensions.value = `${img.width} × ${img.height}`
    }
    img.src = e.target.result
  }
  reader.readAsDataURL(file.raw)
}

const selectFormat = (format) => {
  selectedFormat.value = format
  convertedImage.value = ''
}

const convertImage = async () => {
  if (!originalFile.value || !selectedFormat.value) return
  
  converting.value = true
  
  try {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    const img = new Image()
    
    img.onload = () => {
      canvas.width = img.width
      canvas.height = img.height
      ctx.drawImage(img, 0, 0)
      
      const mimeType = `image/${selectedFormat.value}`
      const quality = selectedFormat.value === 'jpeg' ? 0.9 : undefined
      const dataURL = canvas.toDataURL(mimeType, quality)
      convertedImage.value = dataURL
      
      // 计算转换后文件大小
      const base64Data = dataURL.split(',')[1]
      const byteCharacters = atob(base64Data)
      convertedFileSize.value = byteCharacters.length
      
      converting.value = false
      ElMessage.success(`成功转换为 ${selectedFormat.value.toUpperCase()} 格式`)
    }
    
    img.src = originalImage.value
  } catch (error) {
    converting.value = false
    ElMessage.error('转换失败: ' + error.message)
  }
}

const downloadImage = () => {
  if (!convertedImage.value) return
  
  const link = document.createElement('a')
  link.download = `converted_image.${selectedFormat.value}`
  link.href = convertedImage.value
  link.click()
  ElMessage.success('下载开始')
}

const resetConverter = () => {
  originalFile.value = null
  originalImage.value = ''
  convertedImage.value = ''
  selectedFormat.value = ''
  fileList.value = []
  originalFileSize.value = 0
  convertedFileSize.value = 0
  originalFormat.value = ''
  imageDimensions.value = ''
  uploadRef.value?.clearFiles()
}
</script>

<style scoped>
.image-converter-container {
  padding: 20px;
  max-width: 1200px;
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
  gap: 24px;
}

.upload-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.conversion-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.format-selector h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 16px 0;
}

.format-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.format-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.format-option:hover {
  border-color: #409eff;
  background-color: #f8fafc;
}

.format-option.active {
  border-color: #409eff;
  background-color: rgba(64, 158, 255, 0.05);
}

.format-icon {
  font-size: 24px;
}

.format-name {
  font-weight: 600;
  color: #2c3e50;
}

.format-desc {
  font-size: 12px;
  color: #6b7280;
}

.image-preview {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.preview-panel {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.panel-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.image-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 16px;
}

.image-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  background: #f9fafb;
  border-radius: 8px;
  border: 2px dashed #d1d5db;
}

.image-container img {
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
  border-radius: 4px;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

@media (max-width: 768px) {
  .image-converter-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .format-options {
    grid-template-columns: 1fr;
  }
  
  .image-preview {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>