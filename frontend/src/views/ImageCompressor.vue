<template>
  <div class="image-compressor-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-image' })">图像类</el-breadcrumb-item>
          <el-breadcrumb-item>图片压缩</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>图片压缩工具</h1>
      <p>在线图片压缩，支持JPG、PNG等格式，保持画质的同时减小文件大小</p>
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
          accept="image/jpeg,image/png,image/webp"
          :file-list="fileList"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将图片拖拽到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持 JPG、PNG、WebP 格式，单张图片不超过 10MB
            </div>
          </template>
        </el-upload>
      </div>

      <div v-if="originalImage" class="compression-section">
        <div class="settings-panel">
          <el-card class="settings-card">
            <template #header>
              <div class="card-header">
                <el-icon><Setting /></el-icon>
                <span>压缩设置</span>
              </div>
            </template>
            <div class="settings-content">
              <div class="setting-item">
                <label>压缩质量</label>
                <el-slider
                  v-model="quality"
                  :min="10"
                  :max="100"
                  :step="5"
                  show-stops
                  show-tooltip
                  @change="compressImage"
                />
                <div class="quality-labels">
                  <span class="label-left">高压缩</span>
                  <span class="label-center">平衡</span>
                  <span class="label-right">高质量</span>
                </div>
              </div>
              
              <div class="setting-item">
                <label>输出格式</label>
                <el-radio-group v-model="outputFormat" @change="compressImage">
                  <el-radio label="jpeg">JPEG</el-radio>
                  <el-radio label="png">PNG</el-radio>
                  <el-radio label="webp">WebP</el-radio>
                </el-radio-group>
              </div>

              <div class="setting-item">
                <el-button @click="compressImage" type="primary" :loading="compressing">
                  <el-icon><Refresh /></el-icon>
                  重新压缩
                </el-button>
                <el-button @click="resetImage" type="default">
                  <el-icon><RefreshLeft /></el-icon>
                  重置
                </el-button>
              </div>
            </div>
          </el-card>
        </div>

        <div class="comparison-section">
          <div class="image-comparison">
            <div class="image-panel original-panel">
              <div class="panel-header">
                <h3>原图</h3>
                <div class="image-info">
                  <span>{{ originalSize }}</span>
                  <span>{{ originalDimensions }}</span>
                </div>
              </div>
              <div class="image-container">
                <img :src="originalImage" alt="原图" />
              </div>
            </div>

            <div class="image-panel compressed-panel">
              <div class="panel-header">
                <h3>压缩后</h3>
                <div class="image-info">
                  <span>{{ compressedSize }}</span>
                  <span>压缩率: {{ compressionRatio }}</span>
                </div>
              </div>
              <div class="image-container">
                <img v-if="compressedImage" :src="compressedImage" alt="压缩后" />
                <div v-else class="loading-placeholder">
                  <el-icon v-if="compressing" class="is-loading"><Loading /></el-icon>
                  <span v-if="compressing">压缩中...</span>
                  <span v-else>等待压缩</span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="compressedImage" class="download-section">
            <el-button @click="downloadImage" type="success" size="large">
              <el-icon><Download /></el-icon>
              下载压缩图片
            </el-button>
          </div>
        </div>
      </div>

      <div class="info-section">
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>使用说明</span>
            </div>
          </template>
          <div class="info-content">
            <div class="info-item">
              <h4>支持格式</h4>
              <p>支持 JPEG、PNG、WebP 格式的图片压缩和转换</p>
            </div>
            <div class="info-item">
              <h4>质量说明</h4>
              <ul>
                <li><strong>10-30:</strong> 高压缩率，适合网页缩略图</li>
                <li><strong>40-70:</strong> 平衡质量与大小，适合一般使用</li>
                <li><strong>80-100:</strong> 高质量，适合重要图片</li>
              </ul>
            </div>
            <div class="info-item">
              <h4>格式特点</h4>
              <ul>
                <li><strong>JPEG:</strong> 适合照片，不支持透明背景</li>
                <li><strong>PNG:</strong> 支持透明背景，适合图标和截图</li>
                <li><strong>WebP:</strong> 现代格式，压缩率更高，兼容性较新</li>
              </ul>
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
  UploadFilled,
  Setting,
  Refresh,
  RefreshLeft,
  Download,
  Loading,
  InfoFilled
} from '@element-plus/icons-vue'

const router = useRouter()

const uploadRef = ref()
const fileList = ref([])
const originalFile = ref(null)
const originalImage = ref('')
const compressedImage = ref('')
const quality = ref(80)
const outputFormat = ref('jpeg')
const compressing = ref(false)

const originalFileSize = ref(0)
const compressedFileSize = ref(0)
const imageDimensions = ref({ width: 0, height: 0 })

const originalSize = computed(() => {
  return formatFileSize(originalFileSize.value)
})

const compressedSize = computed(() => {
  return formatFileSize(compressedFileSize.value)
})

const originalDimensions = computed(() => {
  if (imageDimensions.value.width && imageDimensions.value.height) {
    return `${imageDimensions.value.width} × ${imageDimensions.value.height}`
  }
  return ''
})

const compressionRatio = computed(() => {
  if (originalFileSize.value && compressedFileSize.value) {
    const ratio = ((originalFileSize.value - compressedFileSize.value) / originalFileSize.value * 100).toFixed(1)
    return `${ratio}%`
  }
  return '0%'
})

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const beforeUpload = (file) => {
  const isValidType = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  if (!isValidType) {
    ElMessage.error('只能上传 JPG、PNG、WebP 格式的图片')
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
  if (!beforeUpload(file.raw)) {
    return
  }
  
  originalFile.value = file.raw
  originalFileSize.value = file.raw.size
  
  const reader = new FileReader()
  reader.onload = (e) => {
    originalImage.value = e.target.result
    
    // 获取图片尺寸
    const img = new Image()
    img.onload = () => {
      imageDimensions.value = {
        width: img.width,
        height: img.height
      }
      // 自动压缩
      compressImage()
    }
    img.src = e.target.result
  }
  reader.readAsDataURL(file.raw)
}

const compressImage = async () => {
  if (!originalFile.value) return
  
  compressing.value = true
  
  try {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    const img = new Image()
    
    img.onload = () => {
      canvas.width = img.width
      canvas.height = img.height
      ctx.drawImage(img, 0, 0)
      
      const mimeType = `image/${outputFormat.value}`
      const dataURL = canvas.toDataURL(mimeType, quality.value / 100)
      compressedImage.value = dataURL
      
      // 计算压缩后文件大小
      const base64Data = dataURL.split(',')[1]
      const byteCharacters = atob(base64Data)
      compressedFileSize.value = byteCharacters.length
      
      compressing.value = false
      ElMessage.success('图片压缩完成')
    }
    
    img.src = originalImage.value
  } catch (error) {
    compressing.value = false
    ElMessage.error('压缩失败: ' + error.message)
  }
}

const downloadImage = () => {
  if (!compressedImage.value) return
  
  const link = document.createElement('a')
  link.download = `compressed_image.${outputFormat.value}`
  link.href = compressedImage.value
  link.click()
  ElMessage.success('下载开始')
}

const resetImage = () => {
  originalFile.value = null
  originalImage.value = ''
  compressedImage.value = ''
  fileList.value = []
  originalFileSize.value = 0
  compressedFileSize.value = 0
  imageDimensions.value = { width: 0, height: 0 }
  quality.value = 80
  outputFormat.value = 'jpeg'
  uploadRef.value?.clearFiles()
}
</script>

<style scoped>
.image-compressor-container {
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
  gap: 24px;
}

.upload-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.upload-container {
  width: 100%;
}

.compression-section {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 24px;
}

.settings-card {
  border: 1px solid #e5e7eb;
  height: fit-content;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #2c3e50;
}

.settings-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.setting-item label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #374151;
}

.quality-labels {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #6b7280;
  margin-top: 8px;
}

.comparison-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.image-comparison {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.image-panel {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.panel-header {
  margin-bottom: 16px;
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

.loading-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #6b7280;
}

.loading-placeholder .el-icon {
  font-size: 32px;
}

.download-section {
  display: flex;
  justify-content: center;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.info-section {
  margin-top: 20px;
}

.info-card {
  border: 1px solid #e5e7eb;
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

.info-item ul {
  margin: 8px 0 0 20px;
  padding: 0;
}

.info-item li {
  margin: 4px 0;
}

@media (max-width: 768px) {
  .image-compressor-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .compression-section {
    grid-template-columns: 1fr;
  }
  
  .image-comparison {
    grid-template-columns: 1fr;
  }
}
</style>