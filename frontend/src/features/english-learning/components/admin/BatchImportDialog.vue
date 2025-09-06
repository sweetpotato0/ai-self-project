<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="批量导入歌曲"
    width="800px"
    :before-close="handleClose"
    class="batch-import-dialog"
  >
    <div class="import-content">
      <!-- 导入方式选择 -->
      <div class="import-methods">
        <h3 class="section-title">选择导入方式</h3>
        <el-radio-group v-model="importMethod" class="method-group">
          <el-radio-button label="json">JSON 格式</el-radio-button>
          <el-radio-button label="csv">CSV 文件</el-radio-button>
          <el-radio-button label="template">使用模板</el-radio-button>
        </el-radio-group>
      </div>

      <!-- JSON 导入 -->
      <div v-if="importMethod === 'json'" class="import-section">
        <div class="section-header">
          <h4>JSON 数据导入</h4>
          <el-button 
            size="small" 
            type="primary" 
            text 
            @click="showJsonExample"
          >
            查看示例格式
          </el-button>
        </div>
        <el-input
          v-model="jsonData"
          type="textarea"
          :rows="12"
          placeholder="请粘贴 JSON 格式的歌曲数据..."
          class="json-input"
        />
        <div class="json-actions">
          <el-button @click="validateJson">验证格式</el-button>
          <el-button type="primary" @click="parseJsonData">解析数据</el-button>
        </div>
      </div>

      <!-- CSV 文件导入 -->
      <div v-if="importMethod === 'csv'" class="import-section">
        <div class="section-header">
          <h4>CSV 文件导入</h4>
          <el-button 
            size="small" 
            type="primary" 
            text 
            @click="downloadCsvTemplate"
          >
            下载模板
          </el-button>
        </div>
        <el-upload
          ref="csvUploadRef"
          :auto-upload="false"
          :show-file-list="true"
          :on-change="handleCsvChange"
          :before-remove="handleCsvRemove"
          accept=".csv"
          drag
          class="csv-upload"
        >
          <el-icon class="el-icon--upload"><Upload /></el-icon>
          <div class="el-upload__text">
            将 CSV 文件拖拽到此处，或 <em>点击选择文件</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              只能上传 CSV 格式文件，且不超过 5MB
            </div>
          </template>
        </el-upload>
      </div>

      <!-- 模板导入 -->
      <div v-if="importMethod === 'template'" class="import-section">
        <div class="section-header">
          <h4>使用导入模板</h4>
          <el-text type="info" size="small">
            快速创建多首歌曲的基本信息，后续可逐一完善
          </el-text>
        </div>
        <div class="template-form">
          <el-form ref="templateFormRef" :model="templateForm" label-width="120px">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="默认分类">
                  <el-select v-model="templateForm.category_id" placeholder="选择分类" style="width: 100%">
                    <el-option
                      v-for="category in categories"
                      :key="category.id"
                      :label="category.name_cn"
                      :value="category.id"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="默认难度">
                  <el-select v-model="templateForm.difficulty" placeholder="选择难度" style="width: 100%">
                    <el-option label="入门 (1星)" :value="1" />
                    <el-option label="初级 (2星)" :value="2" />
                    <el-option label="中级 (3星)" :value="3" />
                    <el-option label="高级 (4星)" :value="4" />
                    <el-option label="专家 (5星)" :value="5" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="歌曲数量">
              <el-input-number v-model="templateForm.count" :min="1" :max="20" />
            </el-form-item>
          </el-form>
          <el-button type="primary" @click="generateTemplate">生成模板</el-button>
        </div>
      </div>

      <!-- 数据预览 -->
      <div v-if="previewData.length > 0" class="preview-section">
        <div class="section-header">
          <h4>数据预览</h4>
          <el-text type="success" size="small">
            共解析出 {{ previewData.length }} 条有效数据
          </el-text>
        </div>
        <div class="preview-table-container">
          <el-table 
            :data="previewData" 
            max-height="300"
            stripe
            border
            size="small"
          >
            <el-table-column prop="title" label="标题" min-width="150" show-overflow-tooltip />
            <el-table-column prop="title_cn" label="中文标题" min-width="120" show-overflow-tooltip />
            <el-table-column prop="category_name" label="分类" width="100" />
            <el-table-column prop="difficulty" label="难度" width="80">
              <template #default="scope">
                <el-rate 
                  :model-value="scope.row.difficulty" 
                  disabled 
                  show-score 
                  text-color="#99A9BF"
                  score-template="{value}星"
                />
              </template>
            </el-table-column>
            <el-table-column prop="audio_url" label="音频" width="80">
              <template #default="scope">
                <el-tag :type="scope.row.audio_url ? 'success' : 'danger'" size="small">
                  {{ scope.row.audio_url ? '有' : '无' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="scope">
                <el-button 
                  type="danger" 
                  size="small" 
                  text 
                  @click="removePreviewItem(scope.$index)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <!-- 导入选项 -->
        <div class="import-options">
          <h4>导入选项</h4>
          <el-checkbox v-model="importOptions.skipExisting">跳过已存在的歌曲（根据标题判断）</el-checkbox>
          <el-checkbox v-model="importOptions.publishAfterImport">导入后直接发布</el-checkbox>
        </div>
      </div>

      <!-- 进度显示 -->
      <div v-if="importing" class="import-progress">
        <el-progress 
          :percentage="importProgress" 
          :status="importStatus"
          :show-text="true"
        />
        <div class="progress-text">
          {{ importProgressText }}
        </div>
      </div>

      <!-- 导入结果 -->
      <div v-if="importResult" class="import-result">
        <el-result
          :icon="importResult.success ? 'success' : 'error'"
          :title="importResult.title"
          :sub-title="importResult.message"
        >
          <template #extra>
            <div class="result-stats">
              <el-statistic title="成功导入" :value="importResult.successCount" />
              <el-statistic title="失败数量" :value="importResult.failureCount" />
            </div>
            <div v-if="importResult.errors.length > 0" class="error-list">
              <h4>错误详情：</h4>
              <el-scrollbar max-height="200px">
                <ul class="error-items">
                  <li v-for="(error, index) in importResult.errors" :key="index">
                    {{ error }}
                  </li>
                </ul>
              </el-scrollbar>
            </div>
          </template>
        </el-result>
      </div>
    </div>

    <!-- JSON 示例对话框 -->
    <el-dialog
      v-model="jsonExampleVisible"
      title="JSON 格式示例"
      width="600px"
      append-to-body
    >
      <el-input
        :model-value="jsonExampleData"
        type="textarea"
        :rows="20"
        readonly
        class="json-example"
      />
      <template #footer>
        <el-button @click="copyJsonExample">复制示例</el-button>
        <el-button type="primary" @click="jsonExampleVisible = false">确定</el-button>
      </template>
    </el-dialog>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button 
          v-if="previewData.length > 0 && !importing && !importResult"
          type="primary" 
          :loading="importing"
          @click="startImport"
        >
          开始导入 ({{ previewData.length }} 条)
        </el-button>
        <el-button 
          v-if="importResult"
          type="primary"
          @click="handleImportComplete"
        >
          完成
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { useEnglishLearningStore } from '../../stores/englishLearningStore'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  categories: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const store = useEnglishLearningStore()

// 响应式数据
const importMethod = ref('json')
const jsonData = ref('')
const previewData = ref([])
const importing = ref(false)
const importProgress = ref(0)
const importStatus = ref('')
const importProgressText = ref('')
const importResult = ref(null)
const jsonExampleVisible = ref(false)
const csvUploadRef = ref()

// 模板表单
const templateForm = ref({
  category_id: null,
  difficulty: 1,
  count: 5
})
const templateFormRef = ref()

// 导入选项
const importOptions = ref({
  skipExisting: true,
  publishAfterImport: false
})

// JSON 示例数据
const jsonExampleData = `[
  {
    "title": "Twinkle Twinkle Little Star",
    "title_cn": "小星星",
    "description": "经典英语儿歌，适合初学者",
    "lyrics": "Twinkle, twinkle, little star\\nHow I wonder what you are\\nUp above the world so high\\nLike a diamond in the sky",
    "lyrics_cn": "一闪一闪小星星\\n我想知道你是什么\\n高高挂在天空中\\n就像天空中的钻石",
    "audio_url": "https://example.com/audio/twinkle-star.mp3",
    "video_url": "https://example.com/video/twinkle-star.mp4",
    "cover_image": "https://example.com/images/twinkle-star.jpg",
    "category_id": 1,
    "difficulty": 1,
    "age_range": "3-6",
    "tags": ["儿歌", "启蒙", "经典"],
    "is_published": true,
    "duration": 120,
    "sort": 0
  },
  {
    "title": "Old MacDonald Had a Farm",
    "title_cn": "老麦当劳有个农场",
    "description": "有趣的农场动物歌曲",
    "lyrics": "Old MacDonald had a farm, E-I-E-I-O\\nAnd on his farm he had a cow, E-I-E-I-O",
    "audio_url": "https://example.com/audio/old-macdonald.mp3",
    "category_id": 1,
    "difficulty": 2,
    "age_range": "3-6",
    "tags": ["动物", "农场", "互动"],
    "is_published": true,
    "duration": 180,
    "sort": 1
  }
]`

// 方法
const handleClose = async () => {
  if (previewData.value.length > 0 || importing.value) {
    try {
      await ElMessageBox.confirm(
        '关闭对话框将丢失当前数据，确定要关闭吗？',
        '确认关闭',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
    } catch {
      return
    }
  }
  resetDialog()
  emit('update:modelValue', false)
}

const resetDialog = () => {
  importMethod.value = 'json'
  jsonData.value = ''
  previewData.value = []
  importing.value = false
  importProgress.value = 0
  importResult.value = null
  templateForm.value = {
    category_id: null,
    difficulty: 1,
    count: 5
  }
}

const showJsonExample = () => {
  jsonExampleVisible.value = true
}

const copyJsonExample = async () => {
  try {
    await navigator.clipboard.writeText(jsonExampleData)
    ElMessage.success('示例已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

const validateJson = () => {
  try {
    const data = JSON.parse(jsonData.value)
    if (Array.isArray(data)) {
      ElMessage.success(`JSON 格式正确，包含 ${data.length} 条数据`)
    } else {
      ElMessage.warning('JSON 应该是一个数组格式')
    }
  } catch (error) {
    ElMessage.error('JSON 格式错误：' + error.message)
  }
}

const parseJsonData = () => {
  try {
    const data = JSON.parse(jsonData.value)
    if (!Array.isArray(data)) {
      ElMessage.error('数据应该是数组格式')
      return
    }
    
    const parsed = data.map(item => ({
      ...item,
      category_name: getCategoryName(item.category_id),
      tags: Array.isArray(item.tags) ? item.tags : (item.tags || '').split(',').filter(t => t.trim())
    }))
    
    previewData.value = parsed
    ElMessage.success(`成功解析 ${parsed.length} 条数据`)
  } catch (error) {
    ElMessage.error('解析失败：' + error.message)
  }
}

const getCategoryName = (categoryId) => {
  const category = props.categories.find(c => c.id === categoryId)
  return category ? category.name_cn : '未知分类'
}

const handleCsvChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const csvData = e.target.result
      const parsed = parseCsvData(csvData)
      previewData.value = parsed
      ElMessage.success(`成功解析 ${parsed.length} 条数据`)
    } catch (error) {
      ElMessage.error('CSV 解析失败：' + error.message)
    }
  }
  reader.readAsText(file.raw)
}

const handleCsvRemove = () => {
  previewData.value = []
}

const parseCsvData = (csvData) => {
  const lines = csvData.split('\n').filter(line => line.trim())
  if (lines.length < 2) throw new Error('CSV 文件格式不正确')
  
  const headers = lines[0].split(',').map(h => h.trim().replace(/"/g, ''))
  const data = []
  
  for (let i = 1; i < lines.length; i++) {
    const values = lines[i].split(',').map(v => v.trim().replace(/"/g, ''))
    const item = {}
    headers.forEach((header, index) => {
      item[header] = values[index] || ''
    })
    
    // 转换数据类型
    if (item.category_id) item.category_id = parseInt(item.category_id)
    if (item.difficulty) item.difficulty = parseInt(item.difficulty)
    if (item.duration) item.duration = parseInt(item.duration)
    if (item.sort) item.sort = parseInt(item.sort)
    if (item.is_published) item.is_published = item.is_published === 'true'
    if (item.tags) item.tags = item.tags.split(',').map(t => t.trim()).filter(t => t)
    
    item.category_name = getCategoryName(item.category_id)
    data.push(item)
  }
  
  return data
}

const downloadCsvTemplate = () => {
  const csvContent = `title,title_cn,description,lyrics,lyrics_cn,audio_url,video_url,cover_image,category_id,difficulty,age_range,tags,is_published,duration,sort
"Twinkle Twinkle Little Star","小星星","经典英语儿歌","Twinkle, twinkle, little star...","一闪一闪小星星...","https://example.com/audio.mp3","","https://example.com/cover.jpg",1,1,"3-6","儿歌,启蒙",true,120,0`

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', 'songs_template.csv')
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  ElMessage.success('模板已下载')
}

const generateTemplate = () => {
  const templates = []
  const { category_id, difficulty, count } = templateForm.value
  
  for (let i = 1; i <= count; i++) {
    templates.push({
      title: `Song Title ${i}`,
      title_cn: `歌曲标题 ${i}`,
      description: '请填写歌曲描述',
      lyrics: '请填写英文歌词',
      lyrics_cn: '请填写中文歌词',
      audio_url: '',
      video_url: '',
      cover_image: '',
      category_id,
      difficulty,
      age_range: '3-6',
      tags: [],
      is_published: false,
      duration: null,
      sort: i - 1,
      category_name: getCategoryName(category_id)
    })
  }
  
  previewData.value = templates
  ElMessage.success(`生成 ${count} 条模板数据`)
}

const removePreviewItem = (index) => {
  previewData.value.splice(index, 1)
}

const startImport = async () => {
  importing.value = true
  importProgress.value = 0
  importStatus.value = ''
  importProgressText.value = '开始导入...'
  
  const results = {
    success: true,
    successCount: 0,
    failureCount: 0,
    errors: []
  }
  
  try {
    for (let i = 0; i < previewData.value.length; i++) {
      const item = previewData.value[i]
      importProgressText.value = `正在导入: ${item.title}...`
      
      try {
        // 检查是否跳过已存在的歌曲
        if (importOptions.value.skipExisting) {
          const existing = await checkSongExists(item.title)
          if (existing) {
            results.errors.push(`歌曲 "${item.title}" 已存在，跳过`)
            continue
          }
        }
        
        // 处理标签格式
        const songData = {
          ...item,
          tags: Array.isArray(item.tags) ? item.tags.join(',') : item.tags,
          is_published: importOptions.value.publishAfterImport || item.is_published
        }
        delete songData.category_name
        
        await store.createSong(songData)
        results.successCount++
        
      } catch (error) {
        results.failureCount++
        results.errors.push(`歌曲 "${item.title}" 导入失败: ${error.message}`)
      }
      
      importProgress.value = Math.round(((i + 1) / previewData.value.length) * 100)
    }
    
    if (results.failureCount === 0) {
      importStatus.value = 'success'
      importProgressText.value = '导入完成！'
    } else {
      importStatus.value = 'warning'
      importProgressText.value = '导入完成，但有部分失败'
    }
    
  } catch (error) {
    results.success = false
    results.errors.push('导入过程中发生错误: ' + error.message)
    importStatus.value = 'exception'
    importProgressText.value = '导入失败'
  }
  
  importResult.value = {
    ...results,
    title: results.success ? '导入完成' : '导入失败',
    message: `成功 ${results.successCount} 条，失败 ${results.failureCount} 条`
  }
  
  importing.value = false
}

const checkSongExists = async (title) => {
  try {
    const response = await store.searchSongs({ query: title, limit: 1 })
    return response.songs && response.songs.some(song => song.title === title)
  } catch {
    return false
  }
}

const handleImportComplete = () => {
  emit('success')
  resetDialog()
  emit('update:modelValue', false)
}
</script>

<style scoped>
.batch-import-dialog :deep(.el-dialog__body) {
  padding: 24px;
  max-height: 70vh;
  overflow-y: auto;
}

.import-content {
  font-size: 14px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 16px 0;
}

.import-methods {
  margin-bottom: 24px;
}

.method-group {
  width: 100%;
}

.import-section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h4 {
  font-size: 14px;
  font-weight: 600;
  margin: 0;
  color: #374151;
}

.json-input {
  margin-bottom: 16px;
}

.json-actions {
  display: flex;
  gap: 12px;
}

.csv-upload {
  width: 100%;
}

.template-form {
  background: #f8fafc;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.preview-section {
  margin-bottom: 24px;
}

.preview-table-container {
  margin: 16px 0;
}

.import-options {
  margin-top: 16px;
  padding: 16px;
  background: #f0f9ff;
  border-radius: 8px;
}

.import-options h4 {
  font-size: 14px;
  margin: 0 0 12px 0;
  color: #1f2937;
}

.import-progress {
  margin: 20px 0;
  text-align: center;
}

.progress-text {
  margin-top: 8px;
  color: #6b7280;
  font-size: 13px;
}

.import-result {
  text-align: center;
}

.result-stats {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin: 20px 0;
}

.error-list {
  text-align: left;
  max-width: 600px;
  margin: 20px auto 0;
}

.error-list h4 {
  font-size: 14px;
  margin: 0 0 12px 0;
  color: #dc2626;
}

.error-items {
  list-style: none;
  padding: 0;
  margin: 0;
}

.error-items li {
  padding: 4px 0;
  border-bottom: 1px solid #f3f4f6;
  color: #6b7280;
  font-size: 12px;
}

.json-example {
  font-family: 'Courier New', monospace;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .batch-import-dialog :deep(.el-dialog) {
    width: 95vw !important;
    margin: 20px auto;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .result-stats {
    flex-direction: column;
    gap: 16px;
  }
}
</style>