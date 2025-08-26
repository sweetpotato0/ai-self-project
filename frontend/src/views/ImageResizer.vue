<template>
  <div class="image-resizer-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-image' })">图像类</el-breadcrumb-item>
          <el-breadcrumb-item>尺寸调整</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>图片尺寸调整</h1>
      <p>调整图片尺寸，支持等比缩放和自定义尺寸</p>
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
              支持各种图片格式，单张图片不超过 10MB
            </div>
          </template>
        </el-upload>
      </div>

      <div v-if="originalImage" class="resize-section">
        <div class="settings-panel">
          <el-card class="settings-card">
            <template #header>
              <div class="card-header">
                <el-icon><Setting /></el-icon>
                <span>尺寸设置</span>
              </div>
            </template>
            <div class="settings-content">
              <div class="original-info">
                <h4>原图信息</h4>
                <div class="info-item">
                  <span>尺寸: {{ originalWidth }} × {{ originalHeight }}</span>
                </div>
                <div class="info-item">
                  <span>宽高比: {{ aspectRatio }}</span>
                </div>
              </div>

              <div class="resize-options">
                <h4>调整方式</h4>
                <el-radio-group v-model="resizeMode" @change="onResizeModeChange">
                  <el-radio label="percentage">按百分比</el-radio>
                  <el-radio label="dimensions">自定义尺寸</el-radio>
                  <el-radio label="preset">预设尺寸</el-radio>
                </el-radio-group>
              </div>

              <div v-if="resizeMode === 'percentage'" class="percentage-settings">
                <label>缩放比例 (%)</label>
                <el-slider
                  v-model="scalePercentage"
                  :min="10"
                  :max="200"
                  :step="5"
                  show-input
                  @change="calculateByPercentage"
                />
              </div>

              <div v-if="resizeMode === 'dimensions'" class="dimension-settings">
                <div class="dimension-inputs">
                  <div class="input-group">
                    <label>宽度 (px)</label>
                    <el-input-number 
                      v-model="newWidth" 
                      :min="1" 
                      :max="4000"
                      @change="onWidthChange"
                    />
                  </div>
                  <div class="input-group">
                    <label>高度 (px)</label>
                    <el-input-number 
                      v-model="newHeight" 
                      :min="1" 
                      :max="4000"
                      @change="onHeightChange"
                    />
                  </div>
                </div>
                <div class="maintain-aspect">
                  <el-checkbox v-model="maintainAspect">保持宽高比</el-checkbox>
                </div>
              </div>

              <div v-if="resizeMode === 'preset'" class="preset-settings">
                <div class="preset-options">
                  <div 
                    v-for="preset in presetSizes" 
                    :key="preset.name"
                    class="preset-item"
                    :class="{ active: selectedPreset === preset.name }"
                    @click="selectPreset(preset)"
                  >
                    <div class="preset-name">{{ preset.name }}</div>
                    <div class="preset-size">{{ preset.width }} × {{ preset.height }}</div>
                  </div>
                </div>
              </div>

              <div class="action-buttons">
                <el-button @click="resizeImage" type="primary" :loading="resizing">
                  <el-icon><Refresh /></el-icon>
                  调整尺寸
                </el-button>
                <el-button @click="resetResizer" type="default">
                  <el-icon><RefreshLeft /></el-icon>
                  重置
                </el-button>
              </div>
            </div>
          </el-card>
        </div>

        <div class="preview-section">
          <div class="image-comparison">
            <div class="image-panel">
              <div class="panel-header">
                <h3>原图</h3>
                <div class="image-info">
                  <span>{{ originalWidth }} × {{ originalHeight }}</span>
                </div>
              </div>
              <div class="image-container">
                <img :src="originalImage" alt="原图" />
              </div>
            </div>

            <div class="image-panel">
              <div class="panel-header">
                <h3>调整后</h3>
                <div class="image-info">
                  <span v-if="newWidth && newHeight">{{ newWidth }} × {{ newHeight }}</span>
                  <span v-else>等待调整</span>
                </div>
              </div>
              <div class="image-container">
                <img v-if="resizedImage" :src="resizedImage" alt="调整后" />
                <div v-else class="loading-placeholder">
                  <el-icon v-if="resizing" class="is-loading"><Loading /></el-icon>
                  <span v-if="resizing">调整中...</span>
                  <span v-else>预览将在这里显示</span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="resizedImage" class="download-section">
            <el-button @click="downloadImage" type="success" size="large">
              <el-icon><Download /></el-icon>
              下载调整后的图片
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
  Setting,
  Refresh,
  RefreshLeft,
  Download,
  Loading
} from '@element-plus/icons-vue'

const router = useRouter()

const uploadRef = ref()
const fileList = ref([])
const originalImage = ref('')
const resizedImage = ref('')
const resizing = ref(false)

const originalWidth = ref(0)
const originalHeight = ref(0)
const newWidth = ref(0)
const newHeight = ref(0)

const resizeMode = ref('percentage')
const scalePercentage = ref(100)
const maintainAspect = ref(true)
const selectedPreset = ref('')

const presetSizes = ref([
  { name: '头像', width: 200, height: 200 },
  { name: '缩略图', width: 150, height: 150 },
  { name: '小图标', width: 64, height: 64 },
  { name: '中等尺寸', width: 400, height: 300 },
  { name: '标准尺寸', width: 800, height: 600 },
  { name: '高清尺寸', width: 1920, height: 1080 },
  { name: '正方形小', width: 300, height: 300 },
  { name: '正方形大', width: 600, height: 600 }
])

const aspectRatio = computed(() => {
  if (originalWidth.value && originalHeight.value) {
    const ratio = originalWidth.value / originalHeight.value
    return ratio.toFixed(3)
  }
  return '0'
})

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
  
  const reader = new FileReader()
  reader.onload = (e) => {
    originalImage.value = e.target.result
    
    const img = new Image()
    img.onload = () => {
      originalWidth.value = img.width
      originalHeight.value = img.height
      newWidth.value = img.width
      newHeight.value = img.height
      scalePercentage.value = 100
    }
    img.src = e.target.result
  }
  reader.readAsDataURL(file.raw)
}

const onResizeModeChange = () => {
  selectedPreset.value = ''
  if (resizeMode.value === 'percentage') {
    calculateByPercentage()
  } else if (resizeMode.value === 'dimensions') {
    newWidth.value = originalWidth.value
    newHeight.value = originalHeight.value
  }
}

const calculateByPercentage = () => {
  newWidth.value = Math.round(originalWidth.value * scalePercentage.value / 100)
  newHeight.value = Math.round(originalHeight.value * scalePercentage.value / 100)
}

const onWidthChange = () => {
  if (maintainAspect.value && originalWidth.value && originalHeight.value) {
    const ratio = originalHeight.value / originalWidth.value
    newHeight.value = Math.round(newWidth.value * ratio)
  }
}

const onHeightChange = () => {
  if (maintainAspect.value && originalWidth.value && originalHeight.value) {
    const ratio = originalWidth.value / originalHeight.value
    newWidth.value = Math.round(newHeight.value * ratio)
  }
}

const selectPreset = (preset) => {
  selectedPreset.value = preset.name
  if (maintainAspect.value) {
    // 根据原图宽高比调整预设尺寸
    const originalRatio = originalWidth.value / originalHeight.value
    const presetRatio = preset.width / preset.height
    
    if (originalRatio > presetRatio) {
      // 原图更宽，以宽度为准
      newWidth.value = preset.width
      newHeight.value = Math.round(preset.width / originalRatio)
    } else {
      // 原图更高，以高度为准
      newHeight.value = preset.height
      newWidth.value = Math.round(preset.height * originalRatio)
    }
  } else {
    newWidth.value = preset.width
    newHeight.value = preset.height
  }
}

const resizeImage = async () => {
  if (!originalImage.value || !newWidth.value || !newHeight.value) return
  
  resizing.value = true
  
  try {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    const img = new Image()
    
    img.onload = () => {
      canvas.width = newWidth.value
      canvas.height = newHeight.value
      
      // 高质量缩放
      ctx.imageSmoothingEnabled = true
      ctx.imageSmoothingQuality = 'high'
      ctx.drawImage(img, 0, 0, newWidth.value, newHeight.value)
      
      resizedImage.value = canvas.toDataURL('image/png')
      resizing.value = false
      ElMessage.success('图片尺寸调整完成')
    }
    
    img.src = originalImage.value
  } catch (error) {
    resizing.value = false
    ElMessage.error('调整失败: ' + error.message)
  }
}

const downloadImage = () => {
  if (!resizedImage.value) return
  
  const link = document.createElement('a')
  link.download = `resized_${newWidth.value}x${newHeight.value}.png`
  link.href = resizedImage.value
  link.click()
  ElMessage.success('下载开始')
}

const resetResizer = () => {
  originalImage.value = ''
  resizedImage.value = ''
  fileList.value = []
  originalWidth.value = 0
  originalHeight.value = 0
  newWidth.value = 0
  newHeight.value = 0
  resizeMode.value = 'percentage'
  scalePercentage.value = 100
  maintainAspect.value = true
  selectedPreset.value = ''
  uploadRef.value?.clearFiles()
}
</script>

<style scoped>
.image-resizer-container {
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

.upload-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.resize-section {
  display: grid;
  grid-template-columns: 350px 1fr;
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

.original-info h4,
.resize-options h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 12px 0;
}

.info-item {
  font-size: 14px;
  color: #6b7280;
  margin: 4px 0;
}

.percentage-settings label,
.dimension-inputs label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.dimension-inputs {
  display: flex;
  gap: 16px;
}

.input-group {
  flex: 1;
}

.maintain-aspect {
  margin-top: 12px;
}

.preset-options {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.preset-item {
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: center;
}

.preset-item:hover {
  border-color: #409eff;
  background-color: #f8fafc;
}

.preset-item.active {
  border-color: #409eff;
  background-color: rgba(64, 158, 255, 0.05);
}

.preset-name {
  font-weight: 500;
  color: #2c3e50;
  font-size: 13px;
}

.preset-size {
  font-size: 11px;
  color: #6b7280;
  margin-top: 4px;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.image-comparison {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.image-panel {
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
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 16px;
}

.image-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 250px;
  background: #f9fafb;
  border-radius: 8px;
  border: 2px dashed #d1d5db;
}

.image-container img {
  max-width: 100%;
  max-height: 400px;
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

.download-section {
  display: flex;
  justify-content: center;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

@media (max-width: 768px) {
  .image-resizer-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .resize-section {
    grid-template-columns: 1fr;
  }
  
  .dimension-inputs {
    flex-direction: column;
  }
  
  .preset-options {
    grid-template-columns: 1fr;
  }
  
  .image-comparison {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>