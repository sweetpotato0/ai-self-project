<template>
  <div class="episode-management">
    <div class="management-header">
      <h2>剧集管理 - {{ seriesInfo?.title_cn || '系列详情' }}</h2>
      <div class="header-actions">
        <el-button @click="$emit('back')">返回系列列表</el-button>
        <el-button type="primary" :icon="Plus" @click="showCreateDialog">
          新建剧集
        </el-button>
        <el-button type="success" :icon="Upload" @click="showBatchImportDialog">
          批量导入
        </el-button>
      </div>
    </div>

    <!-- 系列信息 -->
    <el-card class="series-info" v-if="seriesInfo">
      <div class="info-row">
        <span><strong>系列标题:</strong> {{ seriesInfo.title }} ({{ seriesInfo.title_cn }})</span>
        <span><strong>难度:</strong> {{ seriesInfo.difficulty }}级</span>
        <span><strong>年龄段:</strong> {{ seriesInfo.age_range }}</span>
        <span><strong>剧集数量:</strong> {{ episodesList.length }}</span>
      </div>
    </el-card>

    <!-- 剧集列表 -->
    <el-card class="episode-list">
      <template #header>
        <span>剧集列表 ({{ episodesList.length }})</span>
      </template>

      <el-table :data="episodesList" style="width: 100%" v-loading="loading">
        <el-table-column prop="episode_num" label="集数" width="80" sortable />
        <el-table-column label="缩略图" width="120">
          <template #default="{ row }">
            <el-image
              :src="row.thumbnail"
              fit="cover"
              style="width: 80px; height: 45px; border-radius: 4px;"
            >
              <template #error>
                <div class="image-placeholder">
                  <el-icon><VideoPlay /></el-icon>
                </div>
              </template>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="英文标题" />
        <el-table-column prop="title_cn" label="中文标题" />
        <el-table-column prop="duration" label="时长" width="80">
          <template #default="{ row }">
            {{ formatDuration(row.duration) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_published ? 'success' : 'info'">
              {{ row.is_published ? '已发布' : '未发布' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="editEpisode(row)">编辑</el-button>
            <el-button size="small" type="info" @click="previewEpisode(row)">预览</el-button>
            <el-button size="small" type="danger" @click="deleteEpisode(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && episodesList.length === 0" description="暂无剧集数据" />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      :title="isEditing ? '编辑剧集' : '新建剧集'"
      v-model="showDialog"
      width="700px"
    >
      <el-form
        :model="formData"
        :rules="formRules"
        ref="formRef"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="集数" prop="episode_num">
              <el-input-number v-model="formData.episode_num" :min="1" :max="9999" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="时长(秒)" prop="duration">
              <el-input-number v-model="formData.duration" :min="1" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="英文标题" prop="title">
          <el-input v-model="formData.title" placeholder="英文标题" />
        </el-form-item>

        <el-form-item label="中文标题" prop="title_cn">
          <el-input v-model="formData.title_cn" placeholder="中文标题" />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="剧集描述"
          />
        </el-form-item>

        <el-form-item label="视频URL" prop="video_url">
          <el-input v-model="formData.video_url" placeholder="视频文件URL" />
        </el-form-item>

        <el-form-item label="缩略图URL" prop="thumbnail">
          <el-input v-model="formData.thumbnail" placeholder="缩略图URL" />
        </el-form-item>

        <el-form-item label="字幕文件" prop="subtitles">
          <el-input
            v-model="formData.subtitles"
            type="textarea"
            :rows="2"
            placeholder="字幕内容或字幕文件URL"
          />
        </el-form-item>

        <el-form-item label="文字稿" prop="transcript">
          <el-input
            v-model="formData.transcript"
            type="textarea"
            :rows="3"
            placeholder="视频文字稿"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="formData.sort" :min="1" :max="1000" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="发布状态">
              <el-switch v-model="formData.is_published" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEditing ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 批量导入对话框 -->
    <el-dialog
      title="批量导入剧集"
      v-model="showBatchDialog"
      width="800px"
    >
      <div class="batch-import-content">
        <p class="import-tip">
          请按照以下格式输入剧集信息，每行一个剧集，字段用逗号分隔：<br>
          <code>集数,英文标题,中文标题,视频URL,缩略图URL,时长(秒),描述</code>
        </p>
        <el-input
          v-model="batchImportText"
          type="textarea"
          :rows="10"
          placeholder="示例：&#10;1,Peppa Pig S1E01,小猪佩奇第一季第1集,dQw4w9WgXcQ,https://img.youtube.com/vi/dQw4w9WgXcQ/0.jpg,300,佩奇和乔治在泥坑里玩耍&#10;2,Peppa Pig S1E02,小猪佩奇第一季第2集,xyz123,https://img.youtube.com/vi/xyz123/0.jpg,295,佩奇学会了骑自行车"
        />
        <div class="import-preview" v-if="parsedEpisodes.length > 0">
          <h4>预览 ({{parsedEpisodes.length}} 个剧集)：</h4>
          <el-table :data="parsedEpisodes.slice(0, 5)" size="small" style="width: 100%">
            <el-table-column prop="episode_number" label="集数" width="60" />
            <el-table-column prop="title" label="英文标题" width="150" show-overflow-tooltip />
            <el-table-column prop="title_cn" label="中文标题" width="150" show-overflow-tooltip />
            <el-table-column prop="video_url" label="视频URL" width="100" show-overflow-tooltip />
            <el-table-column prop="duration" label="时长" width="80" />
          </el-table>
          <p v-if="parsedEpisodes.length > 5" class="preview-note">
            仅显示前5条，共 {{parsedEpisodes.length}} 条记录
          </p>
        </div>
      </div>

      <template #footer>
        <el-button @click="showBatchDialog = false">取消</el-button>
        <el-button type="primary" @click="handleBatchImport" :loading="batchImporting" :disabled="parsedEpisodes.length === 0">
          导入 ({{parsedEpisodes.length}}) 个剧集
        </el-button>
      </template>
    </el-dialog>

    <!-- 预览对话框 -->
    <el-dialog
      title="剧集预览"
      v-model="showPreviewDialog"
      width="900px"
    >
      <div class="episode-preview" v-if="previewEpisodeData">
        <div class="preview-header">
          <h3>{{ previewEpisodeData.title_cn }} ({{ previewEpisodeData.title }})</h3>
          <p><strong>第{{ previewEpisodeData.episode_num }}集</strong></p>
        </div>

        <div class="preview-content">
          <div class="video-player" v-if="previewEpisodeData.video_url">
            <!-- YouTube 嵌入播放器 -->
            <iframe
              v-if="getYouTubeEmbedUrl(previewEpisodeData.video_url)"
              :src="getYouTubeEmbedUrl(previewEpisodeData.video_url)"
              class="preview-video"
              frameborder="0"
              allowfullscreen
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            ></iframe>

            <!-- 常规视频播放器 -->
            <video
              v-else
              :src="previewEpisodeData.video_url"
              :poster="previewEpisodeData.thumbnail"
              controls
              preload="metadata"
              class="preview-video"
            >
              您的浏览器不支持视频播放
            </video>
          </div>

          <div v-else class="video-placeholder">
            <div class="video-error">
              <el-icon><VideoPlay /></el-icon>
              <p>无视频URL</p>
            </div>
          </div>

          <div class="episode-description" v-if="previewEpisodeData.description">
            <h4>描述</h4>
            <p>{{ previewEpisodeData.description }}</p>
          </div>

          <div class="episode-transcript" v-if="previewEpisodeData.transcript">
            <h4>文字稿</h4>
            <p class="transcript-content">{{ previewEpisodeData.transcript }}</p>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Upload, VideoPlay } from '@element-plus/icons-vue'
import { useEnglishVideosStore } from '../../stores/englishVideosStore'

const props = defineProps({
  seriesId: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['back'])

const englishVideosStore = useEnglishVideosStore()

// 响应式数据
const seriesInfo = ref(null)
const episodesList = ref([])
const loading = ref(false)
const showDialog = ref(false)
const showBatchDialog = ref(false)
const showPreviewDialog = ref(false)
const isEditing = ref(false)
const submitting = ref(false)
const batchImporting = ref(false)
const formRef = ref()
const batchImportText = ref('')
const previewEpisodeData = ref(null)

// 表单数据
const formData = reactive({
  id: null,
  title: '',
  title_cn: '',
  description: '',
  video_url: '',
  thumbnail: '',
  duration: 0,
  episode_num: 1,
  subtitles: '',
  transcript: '',
  sort: 1,
  is_published: true
})

// 表单验证规则
const formRules = {
  title: [
    { required: true, message: '请输入英文标题', trigger: 'blur' }
  ],
  title_cn: [
    { required: true, message: '请输入中文标题', trigger: 'blur' }
  ],
  video_url: [
    { required: true, message: '请输入视频URL', trigger: 'blur' }
  ],
  episode_num: [
    { required: true, message: '请输入集数', trigger: 'blur' },
    { type: 'number', min: 1, message: '集数必须大于0', trigger: 'blur' }
  ],
  duration: [
    { required: true, message: '请输入时长', trigger: 'blur' },
    { type: 'number', min: 1, message: '时长必须大于0秒', trigger: 'blur' }
  ]
}

// 计算属性 - 解析批量导入数据
const parsedEpisodes = computed(() => {
  if (!batchImportText.value.trim()) return []

  const lines = batchImportText.value.trim().split('\n')
  const episodes = []

  lines.forEach((line, index) => {
    const parts = line.split(',').map(part => part.trim())
    if (parts.length >= 7) {
      episodes.push({
        episode_number: parseInt(parts[0]) || (index + 1),
        title: parts[1],
        title_cn: parts[2],
        video_url: parts[3],
        thumbnail_url: parts[4],
        duration: parseInt(parts[5]) || 0,
        description: parts[6],
        is_published: true,
        sort_order: parseInt(parts[0]) || (index + 1)
      })
    }
  })

  return episodes
})

// 方法
const loadSeriesInfo = async () => {
  try {
    const response = await englishVideosStore.getVideoSeries(props.seriesId)
    seriesInfo.value = response
  } catch (error) {
    ElMessage.error('获取系列信息失败: ' + error.message)
  }
}

const loadEpisodeList = async () => {
  loading.value = true
  try {
    const response = await englishVideosStore.getEpisodes(props.seriesId)
    episodesList.value = response || []
  } catch (error) {
    ElMessage.error('获取剧集列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  isEditing.value = false
  resetForm()

  // 设置默认集数为当前最大集数+1
  const maxEpisodeNum = Math.max(...episodesList.value.map(ep => ep.episode_num), 0)
  formData.episode_num = maxEpisodeNum + 1

  showDialog.value = true
}

const editEpisode = (episode) => {
  isEditing.value = true
  Object.assign(formData, episode)
  showDialog.value = true
}

const resetForm = () => {
  Object.assign(formData, {
    id: null,
    title: '',
    title_cn: '',
    description: '',
    video_url: '',
    thumbnail: '',
    duration: 0,
    episode_num: 1,
    subtitles: '',
    transcript: '',
    sort: 1,
    is_published: true
  })
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true

  try {
    if (isEditing.value) {
      await englishVideosStore.updateEpisode(formData.id, formData)
      ElMessage.success('更新剧集成功')
    } else {
      await englishVideosStore.createEpisode(props.seriesId, formData)
      ElMessage.success('创建剧集成功')
    }

    showDialog.value = false
    loadEpisodeList()
  } catch (error) {
    ElMessage.error(`${isEditing.value ? '更新' : '创建'}剧集失败: ${error.message}`)
  } finally {
    submitting.value = false
  }
}

const deleteEpisode = async (episode) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除剧集"${episode.title_cn}"吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await englishVideosStore.deleteEpisode(episode.id)
    ElMessage.success('删除剧集成功')
    loadEpisodeList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除剧集失败: ' + error.message)
    }
  }
}

const previewEpisode = (episode) => {
  previewEpisodeData.value = episode
  showPreviewDialog.value = true
}

const showBatchImportDialog = () => {
  batchImportText.value = ''
  showBatchDialog.value = true
}

const handleBatchImport = async () => {
  if (parsedEpisodes.value.length === 0) {
    ElMessage.warning('请输入有效的剧集数据')
    return
  }

  batchImporting.value = true

  try {
    await englishVideosStore.batchImportEpisodes(props.seriesId, parsedEpisodes.value)
    ElMessage.success(`成功导入 ${parsedEpisodes.value.length} 个剧集`)
    showBatchDialog.value = false
    loadEpisodeList()
  } catch (error) {
    ElMessage.error('批量导入失败: ' + error.message)
  } finally {
    batchImporting.value = false
  }
}

const formatDuration = (seconds) => {
  if (!seconds) return '00:00'

  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  if (hours > 0) {
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  } else {
    return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
}

const getYouTubeEmbedUrl = (url) => {
  if (!url) return null

  // 支持的YouTube URL格式:
  // https://www.youtube.com/watch?v=VIDEO_ID
  // https://youtu.be/VIDEO_ID
  // https://www.youtube.com/embed/VIDEO_ID

  const patterns = [
    /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([^&\n?#]+)/,
    /youtube\.com\/watch\?.*v=([^&\n?#]+)/
  ]

  for (const pattern of patterns) {
    const match = url.match(pattern)
    if (match && match[1]) {
      return `https://www.youtube.com/embed/${match[1]}?autoplay=0&controls=1&rel=0&showinfo=0`
    }
  }

  return null
}

// 生命周期
onMounted(() => {
  loadSeriesInfo()
  loadEpisodeList()
})
</script>

<style scoped>
.episode-management {
  padding: 20px;
}

.management-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.series-info {
  margin-bottom: 20px;
}

.info-row {
  display: flex;
  gap: 30px;
  align-items: center;
}

.episode-list {
  margin-top: 20px;
}

.image-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 45px;
  background-color: #f5f7fa;
  border-radius: 4px;
  color: #909399;
}

.batch-import-content {
  margin: 20px 0;
}

.import-tip {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 15px;
  color: #606266;
  line-height: 1.6;
}

.import-tip code {
  background: #e6a23c;
  color: white;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
}

.import-preview {
  margin-top: 15px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 6px;
}

.import-preview h4 {
  margin: 0 0 10px 0;
  color: #409eff;
}

.preview-note {
  margin: 8px 0 0 0;
  font-size: 12px;
  color: #909399;
}

.video-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #999;
}

.video-error .el-icon {
  font-size: 40px;
  margin-bottom: 10px;
}

.episode-preview {
  padding: 20px 0;
}

.preview-header {
  text-align: center;
  margin-bottom: 20px;
}

.preview-header h3 {
  margin: 0 0 10px 0;
  color: #303133;
}

.preview-content {
  margin-top: 20px;
}

.video-placeholder {
  margin-bottom: 20px;
  width: 100%;
  height: 300px;
  border-radius: 8px;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-player {
  margin-bottom: 20px;
  position: relative;
  width: 100%;
  padding-bottom: 56.25%; /* 16:9 aspect ratio */
  height: 0;
  background: #000;
  border-radius: 8px;
  overflow: hidden;
}

.preview-video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.episode-description, .episode-transcript {
  margin-bottom: 20px;
}

.episode-description h4, .episode-transcript h4 {
  margin: 0 0 10px 0;
  color: #606266;
  font-size: 16px;
}

.transcript-content {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 6px;
  line-height: 1.6;
  white-space: pre-wrap;
}
</style>