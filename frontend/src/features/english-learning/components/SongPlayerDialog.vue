<template>
  <el-dialog
    v-model="dialogVisible"
    :title="song?.title || '歌曲播放器'"
    width="1100px"
    :before-close="handleClose"
    class="song-player-dialog"
  >
    <div v-if="song" class="player-container">
      <!-- 歌曲信息头部 -->
      <div class="song-header">
        <img
          :src="song.cover_image || defaultCover"
          :alt="song.title"
          class="song-cover-large"
        />
        <div class="song-info">
          <h2 class="song-title">{{ song.title }}</h2>
          <p class="song-title-cn" v-if="song.title_cn">{{ song.title_cn }}</p>

          <div class="song-meta">
            <el-tag
              v-if="song.category"
              :color="song.category.color"
              size="large"
              effect="light"
            >
              {{ song.category.name }}
            </el-tag>

            <div class="difficulty-display">
              <el-icon
                v-for="i in 5"
                :key="i"
                :class="['star', { active: i <= (song.difficulty || 1) }]"
              >
                <Star />
              </el-icon>
              <span class="difficulty-text">{{ getDifficultyText(song.difficulty) }}</span>
            </div>

            <el-tag v-if="song.age_range" size="large" type="info">
              适合年龄: {{ song.age_range }}
            </el-tag>
          </div>

          <p class="song-description" v-if="song.description">
            {{ song.description }}
          </p>
        </div>
      </div>

      <!-- 主要内容区域：视频和歌词并排 -->
      <div class="main-content" v-if="song.video_url || song.lyrics">
        <!-- 左侧：播放器区域 -->
        <div class="player-side">
          <!-- 音频播放器 -->
          <div class="audio-section" v-if="song.audio_url">
            <h4>音频播放</h4>
            <audio
              ref="audioPlayer"
              :src="song.audio_url"
              controls
              preload="metadata"
              @timeupdate="handleTimeUpdate"
              @ended="handleAudioEnded"
            >
              您的浏览器不支持音频播放。
            </audio>
          </div>

          <!-- 视频播放器 -->
          <div class="video-section" v-if="song.video_url">
            <h4>视频播放</h4>

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
            <video
              v-else
              ref="videoPlayer"
              :src="song.video_url"
              controls
              preload="metadata"
              @timeupdate="handleTimeUpdate"
              @ended="handleVideoEnded"
              class="standard-video"
            >
              您的浏览器不支持视频播放。
            </video>
          </div>
        </div>

        <!-- 右侧：歌词区域 -->
        <div class="lyrics-side" v-if="song.lyrics">
          <div class="lyrics-section">
            <h4>歌词</h4>
            <el-tabs v-model="lyricsTab" class="lyrics-tabs">
              <el-tab-pane label="英文" name="english">
                <div class="lyrics-content" v-html="formatLyrics(song.lyrics)"></div>
              </el-tab-pane>
              <el-tab-pane
                v-if="song.lyrics_cn"
                label="中文"
                name="chinese"
              >
                <div class="lyrics-content" v-html="formatLyrics(song.lyrics_cn)"></div>
              </el-tab-pane>
              <el-tab-pane label="双语" name="bilingual" v-if="song.lyrics_cn">
                <div class="bilingual-lyrics">
                  <div
                    v-for="(line, index) in getBilingualLyrics()"
                    :key="index"
                    class="bilingual-line"
                  >
                    <div class="english-line">{{ line.english }}</div>
                    <div class="chinese-line">{{ line.chinese }}</div>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </div>

      <!-- 仅有音频时的简单布局 -->
      <div class="audio-only-section" v-else-if="song.audio_url">
        <div class="audio-section">
          <h4>音频播放</h4>
          <audio
            ref="audioPlayer"
            :src="song.audio_url"
            controls
            preload="metadata"
            @timeupdate="handleTimeUpdate"
            @ended="handleAudioEnded"
          >
            您的浏览器不支持音频播放。
          </audio>
        </div>

        <!-- 歌词显示 -->
        <div class="lyrics-section" v-if="song.lyrics">
          <h4>歌词</h4>
          <el-tabs v-model="lyricsTab" class="lyrics-tabs">
            <el-tab-pane label="英文歌词" name="english">
              <div class="lyrics-content" v-html="formatLyrics(song.lyrics)"></div>
            </el-tab-pane>
            <el-tab-pane
              v-if="song.lyrics_cn"
              label="中文翻译"
              name="chinese"
            >
              <div class="lyrics-content" v-html="formatLyrics(song.lyrics_cn)"></div>
            </el-tab-pane>
            <el-tab-pane label="双语对照" name="bilingual" v-if="song.lyrics_cn">
              <div class="bilingual-lyrics">
                <div
                  v-for="(line, index) in getBilingualLyrics()"
                  :key="index"
                  class="bilingual-line"
                >
                  <div class="english-line">{{ line.english }}</div>
                  <div class="chinese-line">{{ line.chinese }}</div>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>

      <!-- 学习进度区域 -->
      <div class="progress-section">
        <el-card>
          <template #header>
            <div class="section-header">
              <el-icon><TrendCharts /></el-icon>
              <span>学习进度</span>
            </div>
          </template>

          <div class="progress-controls">
            <div class="progress-display">
              <el-progress
                :percentage="learningProgress"
                :stroke-width="12"
                :color="getProgressColor(learningProgress)"
                class="progress-bar"
              />
              <div class="progress-text">完成度: {{ learningProgress }}%</div>
            </div>

            <div class="progress-buttons">
              <el-button
                type="success"
                :disabled="learningProgress >= 100"
                @click="updateProgress(25)"
              >
                +25% 进度
              </el-button>
              <el-button
                type="success"
                :disabled="learningProgress >= 100"
                @click="updateProgress(50)"
              >
                +50% 进度
              </el-button>
              <el-button
                type="warning"
                :disabled="learningProgress >= 100"
                @click="markCompleted"
              >
                标记完成
              </el-button>
            </div>
          </div>

          <!-- 学习统计 -->
          <div class="learning-stats">
            <div class="stat-item">
              <div class="stat-label">播放次数</div>
              <div class="stat-value">{{ playCount }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">学习时长</div>
              <div class="stat-value">{{ studyTimeMinutes }}分钟</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">最后学习</div>
              <div class="stat-value">{{ formatLastStudied(lastStudiedAt) }}</div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 学习笔记区域 -->
      <div class="notes-section">
        <el-card>
          <template #header>
            <div class="section-header">
              <el-icon><Document /></el-icon>
              <span>学习笔记</span>
            </div>
          </template>

          <el-input
            v-model="learningNotes"
            type="textarea"
            :rows="4"
            placeholder="记录您的学习心得、生词、语法要点..."
            @blur="saveNotes"
          />
        </el-card>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <div class="footer-left">
          <el-button
            :type="isLiked ? 'danger' : 'default'"
            :icon="Star"
            @click="toggleLike"
          >
            {{ isLiked ? '已收藏' : '收藏' }}
          </el-button>
        </div>

        <div class="footer-right">
          <el-button @click="handleClose">关闭</el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import {
  Document,
  Star, TrendCharts
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'

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

const emit = defineEmits([
  'update:modelValue',
  'progress-updated',
  'song-liked',
  'notes-saved'
])

// 响应式数据
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const lyricsTab = ref('english')
const learningProgress = ref(0)
const playCount = ref(0)
const studyTimeMinutes = ref(0)
const lastStudiedAt = ref(null)
const learningNotes = ref('')
const isLiked = ref(false)
const studyStartTime = ref(null)

const audioPlayer = ref(null)
const videoPlayer = ref(null)

const defaultCover = '/images/default-song-cover.jpg'

// 监听器
watch(() => props.song, (newSong) => {
  if (newSong) {
    loadSongProgress()
  }
}, { immediate: true })

watch(dialogVisible, (visible) => {
  if (visible) {
    studyStartTime.value = Date.now()
  } else {
    if (studyStartTime.value) {
      const studyDuration = Math.floor((Date.now() - studyStartTime.value) / 1000 / 60)
      if (studyDuration > 0) {
        studyTimeMinutes.value += studyDuration
        saveProgress()
      }
    }
  }
})

// 方法
const loadSongProgress = () => {
  if (!props.song) return

  // 这里应该从store或API加载真实的进度数据
  // 暂时使用模拟数据
  learningProgress.value = 0
  playCount.value = 0
  studyTimeMinutes.value = 0
  lastStudiedAt.value = new Date()
  learningNotes.value = ''
  isLiked.value = props.song.is_liked || false
}

const handleTimeUpdate = (event) => {
  // 可以根据播放进度自动更新学习进度
  const media = event.target
  const progress = (media.currentTime / media.duration) * 100

  // 如果播放了一定比例，自动增加进度
  if (progress > 50 && learningProgress.value < 25) {
    learningProgress.value = 25
  }
}

const handleAudioEnded = () => {
  playCount.value++
  if (learningProgress.value < 50) {
    learningProgress.value = 50
  }
  saveProgress()
}

const handleVideoEnded = () => {
  playCount.value++
  if (learningProgress.value < 75) {
    learningProgress.value = 75
  }
  saveProgress()
}

const updateProgress = (increment) => {
  learningProgress.value = Math.min(learningProgress.value + increment, 100)
  saveProgress()
}

const markCompleted = () => {
  learningProgress.value = 100
  saveProgress()
  ElMessage.success('恭喜！您已完成这首歌的学习！')
}

const saveProgress = () => {
  if (!props.song) return

  const progressData = {
    progress: learningProgress.value,
    is_completed: learningProgress.value >= 100,
    play_count: playCount.value,
    study_time_minutes: studyTimeMinutes.value,
    notes: learningNotes.value
  }

  emit('progress-updated', progressData)
}

const saveNotes = () => {
  if (!props.song) return

  const progressData = {
    notes: learningNotes.value
  }

  emit('notes-saved', progressData)
  ElMessage.success('笔记已保存')
}

const toggleLike = () => {
  if (!props.song) return

  isLiked.value = !isLiked.value
  emit('song-liked', props.song.id)
}

const handleClose = () => {
  // 暂停播放
  if (audioPlayer.value) {
    audioPlayer.value.pause()
  }
  if (videoPlayer.value) {
    videoPlayer.value.pause()
  }

  dialogVisible.value = false
}

const getDifficultyText = (difficulty) => {
  const texts = ['', '入门', '初级', '中级', '高级', '专家']
  return texts[difficulty] || '未知'
}

const getProgressColor = (percentage) => {
  if (percentage >= 80) return '#67c23a'
  if (percentage >= 50) return '#e6a23c'
  return '#f56c6c'
}

const formatLyrics = (lyrics) => {
  if (!lyrics) return ''
  return lyrics.replace(/\n/g, '<br>')
}

const getBilingualLyrics = () => {
  if (!props.song?.lyrics || !props.song?.lyrics_cn) return []

  const englishLines = props.song.lyrics.split('\n').filter(line => line.trim())
  const chineseLines = props.song.lyrics_cn.split('\n').filter(line => line.trim())

  const maxLength = Math.max(englishLines.length, chineseLines.length)
  const result = []

  for (let i = 0; i < maxLength; i++) {
    result.push({
      english: englishLines[i] || '',
      chinese: chineseLines[i] || ''
    })
  }

  return result
}

const formatLastStudied = (date) => {
  if (!date) return '从未'

  const now = new Date()
  const diff = now - new Date(date)
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return new Date(date).toLocaleDateString()
  }
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
.song-player-dialog :deep(.el-dialog) {
  max-height: 90vh;
  overflow-y: auto;
}

.player-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.song-header {
  display: flex;
  gap: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.song-cover-large {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  object-fit: cover;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.song-info {
  flex: 1;
}

.song-title {
  font-size: 24px;
  font-weight: 700;
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
  align-items: center;
  gap: 16px;
  /* margin-bottom: 16px; */
  flex-wrap: wrap;
}

.difficulty-display {
  display: flex;
  align-items: center;
  gap: 8px;
}

.star {
  color: #d1d5db;
  font-size: 16px;
}

.star.active {
  color: #fbbf24;
}

.difficulty-text {
  font-size: 14px;
  color: #6b7280;
}

.song-description {
  color: #6b7280;
  line-height: 1.6;
  margin: 0;
}

/* 主要内容区域 - 左右布局 */
.main-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  align-items: start;
}

.player-side {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.lyrics-side {
  height: 100%;
}

/* 仅音频时的垂直布局 */
.audio-only-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.audio-section h4,
.video-section h4 {
  margin: 0 0 12px 0;
  font-size: 16px;
  font-weight: 600;
  color: #374151;
}

.audio-section audio {
  width: 100%;
  border-radius: 8px;
}

.standard-video {
  width: 100%;
  max-height: 400px;
  border-radius: 8px;
}

.youtube-player {
  position: relative;
  width: 100%;
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

.lyrics-section {
  background: #f9fafb;
  border-radius: 12px;
  padding: 20px;
  height: 450px; /* 固定高度与视频播放器匹配 */
  display: flex;
  flex-direction: column;
}

.lyrics-section h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #374151;
}

.lyrics-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.lyrics-tabs :deep(.el-tabs__content) {
  padding: 16px 0 0 0;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.lyrics-tabs :deep(.el-tab-pane) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.lyrics-content {
  line-height: 2;
  font-size: 16px;
  color: #374151;
  white-space: pre-line;
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: white;
  border-radius: 8px;
  border-left: 4px solid #667eea;
}

.bilingual-lyrics {
  flex: 1;
  overflow-y: auto;
  padding-right: 8px;
}

.bilingual-line {
  margin-bottom: 16px;
  padding: 12px;
  background: white;
  border-radius: 8px;
}

.english-line {
  font-size: 16px;
  color: #1f2937;
  margin-bottom: 4px;
  font-weight: 500;
}

.chinese-line {
  font-size: 14px;
  color: #6b7280;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.progress-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.progress-display {
  flex: 1;
  margin-right: 20px;
}

.progress-bar {
  margin-bottom: 8px;
}

.progress-text {
  text-align: center;
  font-size: 14px;
  color: #6b7280;
}

.progress-buttons {
  display: flex;
  gap: 8px;
}

.learning-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

@media (max-width: 768px) {
  .song-header {
    flex-direction: column;
    text-align: center;
  }

  /* 移动端改为垂直布局 */
  .main-content {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .lyrics-section {
    height: 350px; /* 移动端减少高度 */
  }

  .progress-controls {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .progress-buttons {
    justify-content: center;
    flex-wrap: wrap;
  }

  .learning-stats {
    grid-template-columns: 1fr;
  }

  .song-meta {
    justify-content: center;
  }

  .song-player-dialog :deep(.el-dialog) {
    width: 95vw !important;
    margin: 20px auto;
  }
}

@media (max-width: 1200px) {
  /* 中等屏幕优化 */
  .main-content {
    gap: 20px;
  }

  .lyrics-section {
    height: 400px;
  }
}
</style>