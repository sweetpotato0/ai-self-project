<template>
  <div class="song-player-page">
    <!-- 顶部导航栏 -->
    <div class="player-header">
      <div class="back-button">
        <el-button
          type="primary"
          :icon="ArrowLeft"
          @click="goBack"
          plain
        >
          返回
        </el-button>
      </div>

      <div class="song-title-header" v-if="song">
        <h1>{{ song.title }}</h1>
        <p v-if="song.title_cn" class="title-cn">{{ song.title_cn }}</p>
      </div>

      <div class="player-actions" v-if="song">
        <el-button
          :type="isLiked ? 'danger' : 'default'"
          :icon="Star"
          @click="toggleLike"
        >
          {{ isLiked ? '已收藏' : '收藏' }}
        </el-button>
      </div>
    </div>

    <div v-if="song" class="player-container">
      <!-- 歌曲信息头部 -->
      <div class="song-header">
        <img
          :src="song.cover_image || defaultCover"
          :alt="song.title"
          class="song-cover-large"
        />
        <div class="song-info">
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

    <!-- 加载中状态 -->
    <div v-else class="loading-container">
      <el-skeleton :loading="loading" animated>
        <template #template>
          <el-skeleton-item variant="image" style="width: 240px; height: 240px;" />
          <div style="padding: 14px;">
            <el-skeleton-item variant="h3" style="width: 50%" />
            <div style="display: flex; align-items: center; justify-items: space-between; margin-top: 16px; height: 16px;">
              <el-skeleton-item variant="text" style="margin-right: 16px;" />
              <el-skeleton-item variant="text" style="width: 30%;" />
            </div>
          </div>
        </template>
      </el-skeleton>
    </div>

    <!-- 错误状态 -->
    <div v-if="error" class="error-container">
      <el-result
        icon="warning"
        title="歌曲加载失败"
        :sub-title="error"
      >
        <template #extra>
          <el-button type="primary" @click="loadSong">重新加载</el-button>
          <el-button @click="goBack">返回列表</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup>
import {
  ArrowLeft,
  Document,
  Star, TrendCharts
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useEnglishLearningStore } from '../stores/englishLearningStore'

const route = useRoute()
const router = useRouter()
const store = useEnglishLearningStore()

// 响应式数据
const song = ref(null)
const loading = ref(true)
const error = ref('')

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

// 从路由参数获取歌曲ID
const songId = computed(() => {
  return parseInt(route.params.id)
})

// 页面加载时获取歌曲数据
onMounted(async () => {
  await loadSong()
  studyStartTime.value = Date.now()
})

// 页面卸载时保存学习时长
onUnmounted(() => {
  if (studyStartTime.value) {
    const studyDuration = Math.floor((Date.now() - studyStartTime.value) / 1000 / 60)
    if (studyDuration > 0) {
      studyTimeMinutes.value += studyDuration
      saveProgress()
    }
  }
})

// 监听歌曲变化
watch(() => route.params.id, async (newId) => {
  if (newId) {
    await loadSong()
  }
})

// 方法
const loadSong = async () => {
  try {
    loading.value = true
    error.value = ''

    const response = await store.fetchSong(songId.value)
    // 提取实际的歌曲数据
    song.value = response.data || response

    // 加载学习进度数据
    loadSongProgress()

  } catch (err) {
    console.error('Failed to load song:', err)
    error.value = '加载歌曲信息失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

const loadSongProgress = () => {
  if (!song.value) return

  // 这里应该从store或API加载真实的进度数据
  // 暂时使用模拟数据
  learningProgress.value = 0
  playCount.value = 0
  studyTimeMinutes.value = 0
  lastStudiedAt.value = new Date()
  learningNotes.value = ''
  isLiked.value = song.value.is_liked || false
}

const goBack = () => {
  router.back()
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
  if (!song.value) return

  const progressData = {
    progress: learningProgress.value,
    is_completed: learningProgress.value >= 100,
    play_count: playCount.value,
    study_time_minutes: studyTimeMinutes.value,
    notes: learningNotes.value
  }

  console.log('Saving progress:', progressData)
  // TODO: 调用API保存进度
}

const saveNotes = () => {
  if (!song.value) return

  console.log('Saving notes:', learningNotes.value)
  ElMessage.success('笔记已保存')
  // TODO: 调用API保存笔记
}

const toggleLike = async () => {
  if (!song.value) return

  try {
    await store.likeSong(song.value.id)
    isLiked.value = !isLiked.value
    ElMessage.success(isLiked.value ? '收藏成功' : '已取消收藏')
  } catch (error) {
    console.error('Toggle like failed:', error)
    ElMessage.error('操作失败，请重试')
  }
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
  if (!song.value?.lyrics || !song.value?.lyrics_cn) return []

  const englishLines = song.value.lyrics.split('\n').filter(line => line.trim())
  const chineseLines = song.value.lyrics_cn.split('\n').filter(line => line.trim())

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
.song-player-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.player-header {
  background: white;
  padding: 12px 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  z-index: 100;
}

.back-button {
  flex-shrink: 0;
}

.song-title-header {
  flex: 1;
  text-align: center;
  margin: 0 20px;
}

.song-title-header h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
}

.title-cn {
  margin: 4px 0 0 0;
  font-size: 16px;
  color: #6b7280;
}

.player-actions {
  flex-shrink: 0;
}

.player-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.song-header {
  display: flex;
  gap: 24px;
  padding: 20px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.song-cover-large {
  width: 40px;
  height: 40px;
  border-radius: 16px;
  object-fit: cover;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.song-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.song-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  /* margin-bottom: 20px; */
  flex-wrap: wrap;
}

.difficulty-display {
  display: flex;
  align-items: center;
  gap: 8px;
}

.star {
  color: #d1d5db;
  font-size: 18px;
}

.star.active {
  color: #fbbf24;
}

.difficulty-text {
  font-size: 16px;
  color: #6b7280;
}

.song-description {
  color: #6b7280;
  line-height: 1.8;
  margin: 0;
  font-size: 16px;
}

/* 主要内容区域 - 左右布局 */
.main-content {
  display: grid;
  grid-template-columns: 3fr 1fr; /* 视频区域占3份，歌词区域占1份 */
  gap: 24px;
  align-items: start;
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.player-side {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.lyrics-side {
  height: 100%;
}

/* 仅音频时的垂直布局 */
.audio-only-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
  background: white;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.audio-section h4,
.video-section h4 {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #374151;
}

.audio-section audio {
  width: 100%;
  border-radius: 12px;
}

.standard-video {
  width: 100%;
  max-height: 500px;
  border-radius: 12px;
}

.youtube-player {
  position: relative;
  width: 100%;
  padding-bottom: 56.25%; /* 16:9 宽高比 */
  height: 0;
  overflow: hidden;
  border-radius: 12px;
  background: #000;
}

.youtube-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 12px;
}

.lyrics-section {
  background: #f9fafb;
  border-radius: 16px;
  padding: 24px;
  height: 550px; /* 调整高度与视频播放器匹配 */
  display: flex;
  flex-direction: column;
}

.lyrics-section h4 {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #374151;
}

.lyrics-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.lyrics-tabs :deep(.el-tabs__content) {
  padding: 20px 0 0 0;
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
  font-size: 14px;
  color: #374151;
  white-space: pre-line;
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: white;
  border-radius: 12px;
  border-left: 4px solid #667eea;
}

.bilingual-lyrics {
  flex: 1;
  overflow-y: auto;
  padding-right: 8px;
}

.bilingual-line {
  margin-bottom: 20px;
  padding: 16px;
  background: white;
  border-radius: 12px;
}

.english-line {
  font-size: 16px;
  color: #1f2937;
  margin-bottom: 8px;
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

.progress-section,
.notes-section {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.progress-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.progress-display {
  flex: 1;
  margin-right: 24px;
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
  gap: 12px;
}

.learning-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: #f9fafb;
  border-radius: 12px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.loading-container,
.error-container {
  max-width: 800px;
  margin: 40px auto;
  padding: 0 24px;
}

@media (max-width: 768px) {
  .player-header {
    padding: 16px;
  }

  .song-title-header h1 {
    font-size: 18px;
  }

  .player-container {
    padding: 20px 16px;
    gap: 24px;
  }

  .song-header {
    flex-direction: column;
    text-align: center;
    padding: 24px;
  }

  .song-cover-large {
    width: 40px;
    height: 40px;
    align-self: center;
  }

  /* 移动端改为垂直布局 */
  .main-content {
    grid-template-columns: 1fr;
    gap: 20px;
    padding: 24px;
  }

  .lyrics-section {
    height: 400px; /* 移动端减少高度 */
  }

  .progress-controls {
    flex-direction: column;
    gap: 20px;
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
}

@media (max-width: 1200px) {
  /* 中等屏幕优化 */
  .main-content {
    gap: 24px;
    grid-template-columns: 4fr 1fr; /* 中等屏幕下也保持视频区域更大 */
  }

  .lyrics-section {
    height: 480px;
  }

  .player-container {
    max-width: 1200px;
  }
}
</style>