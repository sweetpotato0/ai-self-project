<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="歌曲预览"
    width="900px"
    class="song-preview-dialog"
  >
    <div v-if="song" class="preview-content">
      <!-- 歌曲头部信息 -->
      <div class="song-header">
        <div class="song-cover">
          <img
            :src="song.cover_image || defaultCover"
            :alt="song.title"
            class="cover-image"
          />
        </div>
        <div class="song-info">
          <h2 class="song-title">{{ song.title }}</h2>
          <h3 v-if="song.title_cn" class="song-title-cn">{{ song.title_cn }}</h3>
          <div class="song-meta">
            <el-tag
              v-if="song.category"
              :color="song.category.color"
              effect="light"
              size="large"
            >
              {{ song.category.name_cn }}
            </el-tag>
            <div class="difficulty">
              <span class="label">难度：</span>
              <div class="stars">
                <el-icon
                  v-for="i in 5"
                  :key="i"
                  :class="['star', { active: i <= song.difficulty }]"
                >
                  <Star />
                </el-icon>
              </div>
            </div>
            <div v-if="song.age_range" class="age-range">
              <span class="label">适合年龄：</span>
              <el-tag type="info" size="small">{{ song.age_range }}</el-tag>
            </div>
            <!-- <div v-if="song.duration" class="duration">
              <span class="label">时长：</span>
              <span>{{ formatDuration(song.duration) }}</span>
            </div> -->
          </div>
          <div class="song-stats">
            <div class="stat-item">
              <el-icon><View /></el-icon>
              <span>{{ formatNumber(song.view_count || 0) }} 播放</span>
            </div>
            <div class="stat-item">
              <el-icon><Star /></el-icon>
              <span>{{ formatNumber(song.like_count || 0) }} 收藏</span>
            </div>
            <div class="stat-item">
              <el-tag :type="song.is_published ? 'success' : 'info'">
                {{ song.is_published ? '已发布' : '草稿' }}
              </el-tag>
            </div>
          </div>
          <div v-if="song.description" class="song-description">
            <p>{{ song.description }}</p>
          </div>
        </div>
      </div>

      <!-- 媒体播放器 -->
      <div class="media-section">
        <el-tabs v-model="activeMediaTab" class="media-tabs">
          <el-tab-pane label="音频播放" name="audio" :disabled="!song.audio_url">
            <div v-if="song.audio_url" class="audio-player-container">
              <audio
                ref="audioPlayer"
                :src="song.audio_url"
                controls
                class="audio-player"
                @loadedmetadata="handleAudioLoaded"
                @error="handleMediaError"
              />
              <div class="player-info">
                <el-text type="info" size="small">
                  <el-icon><Headset /></el-icon>
                  点击播放按钮开始播放音频
                </el-text>
              </div>
            </div>
            <div v-else class="no-media">
              <el-empty description="暂无音频文件" />
            </div>
          </el-tab-pane>

          <el-tab-pane
            label="视频播放"
            name="video"
            :disabled="!song.video_url"
          >
            <div v-if="song.video_url" class="video-player-container">
              <!-- YouTube播放器 -->
              <div v-if="isYouTubeUrl(song.video_url)" class="youtube-player">
                <iframe
                  :src="getYouTubeEmbedUrl(song.video_url)"
                  frameborder="0"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                  allowfullscreen
                  class="youtube-iframe"
                ></iframe>
              </div>

              <!-- 标准视频播放器 -->
              <div v-else>
                <video
                  ref="videoPlayer"
                  :src="song.video_url"
                  controls
                  class="video-player"
                  @loadedmetadata="handleVideoLoaded"
                  @error="handleMediaError"
                />
              </div>

              <div class="player-info">
                <el-text type="info" size="small">
                  <el-icon><VideoPlay /></el-icon>
                  {{ isYouTubeUrl(song.video_url) ? 'YouTube视频播放' : '点击播放按钮开始播放视频' }}
                </el-text>
              </div>
            </div>
            <div v-else class="no-media">
              <el-empty description="暂无视频文件" />
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 歌词内容 -->
      <div class="lyrics-section">
        <h3 class="section-title">歌词内容</h3>
        <div class="lyrics-container">
          <div class="lyrics-panel">
            <h4 class="lyrics-title">
              <el-icon><Document /></el-icon>
              英文歌词
            </h4>
            <div class="lyrics-content english-lyrics">
              <pre>{{ song.lyrics || '暂无英文歌词' }}</pre>
            </div>
          </div>
          <div v-if="song.lyrics_cn" class="lyrics-panel">
            <h4 class="lyrics-title">
              <el-icon><Document /></el-icon>
              中文翻译
            </h4>
            <div class="lyrics-content chinese-lyrics">
              <pre>{{ song.lyrics_cn }}</pre>
            </div>
          </div>
        </div>
      </div>

      <!-- 标签信息 -->
      <div v-if="song.tags && song.tags.length > 0" class="tags-section">
        <h3 class="section-title">标签</h3>
        <div class="tags-list">
          <el-tag
            v-for="tag in songTags"
            :key="tag"
            size="small"
            class="tag-item"
          >
            {{ tag }}
          </el-tag>
        </div>
      </div>

      <!-- 其他信息 -->
      <div class="additional-info">
        <el-row :gutter="24">
          <el-col :span="12">
            <div class="info-item">
              <span class="label">创建时间：</span>
              <span>{{ formatDate(song.created_at) }}</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="info-item">
              <span class="label">更新时间：</span>
              <span>{{ formatDate(song.updated_at) }}</span>
            </div>
          </el-col>
        </el-row>
        <el-row :gutter="24">
          <el-col :span="12">
            <div class="info-item">
              <span class="label">排序值：</span>
              <span>{{ song.sort || 0 }}</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="info-item">
              <span class="label">歌曲ID：</span>
              <span>{{ song.id }}</span>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <div v-else class="no-song">
      <el-empty description="没有选择歌曲" />
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="$emit('update:modelValue', false)">关闭</el-button>
        <el-button v-if="song" type="primary" @click="handleEdit">
          编辑歌曲
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import {
  Document,
  Headset,
  Star,
  VideoPlay,
  View
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onUnmounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  song: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'edit'])

// 响应式数据
const activeMediaTab = ref('audio')
const audioPlayer = ref()
const videoPlayer = ref()
const defaultCover = 'https://via.placeholder.com/150x150/667eea/ffffff?text=♪'

// 计算属性
const songTags = computed(() => {
  if (!props.song?.tags) return []
  if (Array.isArray(props.song.tags)) return props.song.tags
  return props.song.tags.split(',').filter(tag => tag.trim())
})

// 监听器
watch(() => props.modelValue, (visible) => {
  if (visible) {
    // 设置默认媒体标签
    if (props.song?.audio_url) {
      activeMediaTab.value = 'audio'
    } else if (props.song?.video_url) {
      activeMediaTab.value = 'video'
    }
  } else {
    // 暂停所有媒体播放
    pauseAllMedia()
  }
})

watch(() => activeMediaTab.value, () => {
  // 切换媒体时暂停其他媒体
  pauseAllMedia()
})

// 生命周期
onUnmounted(() => {
  pauseAllMedia()
})

// 方法
const pauseAllMedia = () => {
  if (audioPlayer.value && !audioPlayer.value.paused) {
    audioPlayer.value.pause()
  }
  if (videoPlayer.value && !videoPlayer.value.paused) {
    videoPlayer.value.pause()
  }
}

const handleAudioLoaded = () => {
  console.log('音频加载完成')
}

const handleVideoLoaded = () => {
  console.log('视频加载完成')
}

const handleMediaError = (event) => {
  console.error('媒体文件加载失败:', event)
  ElMessage.error('媒体文件加载失败，请检查文件地址')
}

const handleEdit = () => {
  emit('edit', props.song)
}

// 工具方法
const formatDuration = (seconds) => {
  if (!seconds) return '--'
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

const formatNumber = (num) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

const formatDate = (dateString) => {
  if (!dateString) return '--'
  return new Date(dateString).toLocaleString('zh-CN')
}

// YouTube相关方法
const isYouTubeUrl = (url) => {
  if (!url) return false
  const youtubeRegex = /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+/
  return youtubeRegex.test(url)
}

const getYouTubeEmbedUrl = (url) => {
  if (!url) return ''

  // 提取YouTube视频ID
  let videoId = ''

  if (url.includes('youtu.be/')) {
    // 短链接格式: https://youtu.be/VIDEO_ID
    videoId = url.split('youtu.be/')[1]?.split('?')[0]
  } else if (url.includes('youtube.com/watch')) {
    // 标准链接格式: https://www.youtube.com/watch?v=VIDEO_ID
    const urlParams = new URLSearchParams(url.split('?')[1])
    videoId = urlParams.get('v')
  } else if (url.includes('youtube.com/embed/')) {
    // 已经是embed格式
    return url
  }

  if (videoId) {
    // 移除可能的时间戳参数
    videoId = videoId.split('&')[0].split('#')[0]
    return `https://www.youtube.com/embed/${videoId}?rel=0&modestbranding=1`
  }

  return ''
}
</script>

<style scoped>
.song-preview-dialog :deep(.el-dialog__body) {
  padding: 24px;
  max-height: 70vh;
  overflow-y: auto;
}

.preview-content {
  font-size: 14px;
}

.song-header {
  display: flex;
  gap: 24px;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #ebeef5;
}

.song-cover {
  flex-shrink: 0;
}

.cover-image {
  width: 150px;
  height: 150px;
  border-radius: 12px;
  object-fit: cover;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.song-info {
  flex: 1;
  min-width: 0;
}

.song-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.song-title-cn {
  font-size: 18px;
  color: #6b7280;
  margin: 0 0 16px 0;
}

.song-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: center;
  /* margin-bottom: 16px; */
}

.difficulty {
  display: flex;
  align-items: center;
  gap: 4px;
}

.stars {
  display: flex;
  gap: 2px;
}

.star {
  color: #d1d5db;
  font-size: 14px;
}

.star.active {
  color: #fbbf24;
}

.label {
  font-weight: 500;
  color: #374151;
}

.song-stats {
  display: flex;
  gap: 20px;
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #6b7280;
  font-size: 13px;
}

.song-description {
  color: #4b5563;
  line-height: 1.6;
}

.media-section {
  margin-bottom: 24px;
}

.media-tabs :deep(.el-tabs__content) {
  padding-top: 16px;
}

.audio-player-container,
.video-player-container {
  text-align: center;
}

.audio-player {
  width: 100%;
  max-width: 500px;
  height: 50px;
}

.video-player {
  width: 100%;
  max-width: 600px;
  height: 300px;
  border-radius: 8px;
}

.youtube-player {
  position: relative;
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
  padding-bottom: 56.25%; /* 16:9 宽高比 */
  height: 0;
  overflow: hidden;
  border-radius: 8px;
  background: #000;
}

.youtube-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 8px;
}

.player-info {
  margin-top: 12px;
}

.no-media {
  padding: 40px 0;
}

.lyrics-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 16px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.lyrics-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.lyrics-panel {
  background: #f8fafc;
  border-radius: 8px;
  overflow: hidden;
}

.lyrics-title {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin: 0;
  padding: 12px 16px;
  background: #e5e7eb;
  display: flex;
  align-items: center;
  gap: 6px;
}

.lyrics-content {
  padding: 16px;
  max-height: 300px;
  overflow-y: auto;
}

.lyrics-content pre {
  margin: 0;
  font-family: inherit;
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.8;
}

.english-lyrics {
  color: #1f2937;
}

.chinese-lyrics {
  color: #4b5563;
}

.tags-section {
  margin-bottom: 24px;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  margin: 0;
}

.additional-info {
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
  font-size: 13px;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.info-item .label {
  width: 80px;
  flex-shrink: 0;
}

.no-song {
  padding: 60px 0;
  text-align: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .song-preview-dialog :deep(.el-dialog) {
    width: 95vw !important;
    margin: 20px auto;
  }

  .song-header {
    flex-direction: column;
    gap: 16px;
  }

  .cover-image {
    width: 120px;
    height: 120px;
    margin: 0 auto;
  }

  .song-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .lyrics-container {
    grid-template-columns: 1fr;
  }
}
</style>