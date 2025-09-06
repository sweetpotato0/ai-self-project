<template>
  <div class="video-player-page" :class="{ 'fullscreen-mode': isFullscreen }">
    <!-- 全屏模式 -->
    <div v-if="isFullscreen" class="fullscreen-player">
      <!-- 全屏播放器 -->
      <div class="fullscreen-video-container">
        <div class="fullscreen-video-wrapper">
          <iframe
            v-if="youtubeEmbedUrl"
            :src="youtubeEmbedUrl"
            class="fullscreen-video"
            frameborder="0"
            allowfullscreen
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          ></iframe>

          <video
            v-else
            ref="videoPlayer"
            class="fullscreen-video"
            :src="currentEpisode.video_url"
            :poster="currentEpisode.thumbnail || currentSeries.cover_image"
            controls
            preload="metadata"
            @loadedmetadata="handleVideoLoaded"
            @timeupdate="handleTimeUpdate"
            @ended="handleVideoEnd"
            @play="handleVideoPlay"
            @pause="handleVideoPause"
          ></video>
        </div>
      </div>

      <!-- 全屏控制按钮 -->
      <div class="fullscreen-exit-btn">
        <el-button
          :icon="Aim"
          @click="toggleFullscreen"
          type="danger"
          size="large"
          circle
          class="exit-fullscreen-btn"
        >
        </el-button>
        <div class="exit-hint">退出全屏</div>
      </div>

      <!-- 全屏剧集列表 -->
      <div v-if="showSidebar" class="fullscreen-episodes">
        <div class="fullscreen-episodes-header">
          <h3>{{ currentSeries.title }}</h3>
          <div class="fullscreen-sidebar-controls">
            <span>共 {{ episodes.length }} 集</span>
            <el-button
              :icon="ArrowRight"
              @click="toggleSidebar"
              type="info"
              size="small"
              text
              class="fullscreen-hide-btn"
            />
          </div>
        </div>
        <div class="fullscreen-episodes-list">
          <div
            v-for="episode in episodes"
            :key="episode.id"
            class="fullscreen-episode-item"
            :class="{ 'active': episode.id === currentEpisode?.id }"
            @click="selectEpisode(episode)"
          >
            <img
              :src="episode.thumbnail || currentSeries.cover_image"
              :alt="episode.title"
              @error="handleImageError"
            />
            <div class="episode-info">
              <div class="episode-number">第 {{ episode.episode_num || episode.episode_number }} 集</div>
              <div class="episode-title">{{ episode.title }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 悬浮显示按钮（全屏模式下当侧边栏隐藏时） -->
      <div v-if="!showSidebar" class="fullscreen-floating-sidebar-btn">
        <el-button
          :icon="Menu"
          @click="toggleSidebar"
          type="primary"
          size="large"
          circle
          class="fullscreen-show-sidebar-btn"
        />
      </div>
    </div>

    <!-- 普通模式 -->
    <div v-else-if="currentSeries && currentEpisode" class="player-page-content">
      <!-- 顶部导航栏 -->
      <div class="top-navigation">
        <el-button
          :icon="ArrowLeft"
          @click="goBack"
          type=""
          size="small"
          text
          class="back-btn"
        >
          返回列表
        </el-button>
        <div class="page-title">
          <h1>{{ currentSeries.title }}</h1>
          <span class="episode-indicator">第 {{ currentEpisode.episode_num || currentEpisode.episode_number }} 集 · {{ currentEpisode.title }}</span>
        </div>
      </div>

      <div class="player-container" :class="{ 'sidebar-hidden': !showSidebar }">
        <!-- 主播放区域 -->
        <div class="main-player-area">
          <div class="video-player">
            <div class="player-wrapper" ref="playerContainer">
            <!-- YouTube 嵌入播放器 -->
            <iframe
              v-if="youtubeEmbedUrl"
              :src="youtubeEmbedUrl"
              class="video-element"
              frameborder="0"
              allowfullscreen
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            ></iframe>

            <!-- 常规视频播放器（备用） -->
            <video
              v-else
              ref="videoPlayer"
              class="video-element"
              :src="currentEpisode.video_url"
              :poster="currentEpisode.thumbnail || currentSeries.cover_image"
              controls
              preload="metadata"
              @loadedmetadata="handleVideoLoaded"
              @timeupdate="handleTimeUpdate"
              @ended="handleVideoEnd"
              @play="handleVideoPlay"
              @pause="handleVideoPause"
            >
              <!-- 字幕轨道 -->
              <track
                v-if="currentEpisode.subtitle_url"
                kind="subtitles"
                :src="currentEpisode.subtitle_url"
                srclang="en"
                label="English"
                default
              />
            </video>

            <!-- 自定义字幕显示 -->
            <div v-if="showSubtitles && currentSubtitle" class="custom-subtitle">
              {{ currentSubtitle }}
            </div>

            <!-- 全屏按钮 -->
            <div class="video-fullscreen-btn">
              <el-button
                :icon="FullScreen"
                @click="toggleFullscreen"
                type="info"
                size="small"
                circle
                class="fullscreen-toggle"
              />
            </div>
          </div>
          </div>

          <!-- 学习进度 - 移到播放器正下方 -->
          <div class="learning-progress-card">
            <div class="progress-header">
              <h3>学习进度</h3>
              <span class="progress-text">{{ Math.round(learningProgress) }}%</span>
            </div>
            <el-progress
              :percentage="learningProgress"
              :stroke-width="8"
              :color="getProgressColor(learningProgress)"
            />
            <div class="episode-description" v-if="currentEpisode.description">
              <p>{{ currentEpisode.description }}</p>
            </div>
          </div>

          <!-- 播放控制工具栏 (仅用于非YouTube视频) -->
          <div v-if="!youtubeEmbedUrl" class="player-controls">
            <div class="control-row">
              <div class="control-group">
                <el-button-group>
                  <el-button
                    :icon="!isPlaying ? VideoPlay : VideoPause"
                    @click="togglePlay"
                    size="large"
                  />
                  <el-button
                    :icon="Mute"
                    @click="toggleMute"
                    size="large"
                  />
                </el-button-group>
              </div>

              <div class="control-group">
                <el-button
                  :type="showSubtitles ? 'primary' : ''"
                  @click="toggleSubtitles"
                  size="small"
                >
                  字幕
                </el-button>
                <el-button
                  :type="repeatMode ? 'primary' : ''"
                  :icon="RefreshRight"
                  @click="toggleRepeat"
                  size="small"
                >
                  重复
                </el-button>
                <el-dropdown trigger="click">
                  <el-button size="small">
                    速度 {{ playbackRate }}x
                    <el-icon><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="setPlaybackRate(0.5)">0.5x</el-dropdown-item>
                      <el-dropdown-item @click="setPlaybackRate(0.75)">0.75x</el-dropdown-item>
                      <el-dropdown-item @click="setPlaybackRate(1)">1x</el-dropdown-item>
                      <el-dropdown-item @click="setPlaybackRate(1.25)">1.25x</el-dropdown-item>
                      <el-dropdown-item @click="setPlaybackRate(1.5)">1.5x</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>

            <!-- 进度条 -->
            <div class="progress-row">
              <div class="time-info">{{ formatTime(currentTime) }}</div>
              <el-slider
                v-model="progressPercentage"
                :show-tooltip="false"
                @input="handleProgressChange"
                class="progress-slider"
              />
              <div class="time-info">{{ formatTime(duration) }}</div>
            </div>
          </div>
        </div>

        <!-- 侧边栏：剧集列表 -->
        <div v-if="showSidebar" class="episodes-sidebar">
          <div class="sidebar-header">
            <h3>剧集列表</h3>
            <div class="sidebar-controls">
              <span class="episode-count">共 {{ episodes.length }} 集</span>
              <el-button
                :icon="ArrowRight"
                @click="toggleSidebar"
                type="info"
                size="small"
                text
                class="sidebar-hide-btn"
              />
            </div>
          </div>
          <div class="episodes-list">
            <div
              v-for="episode in episodes"
              :key="episode.id"
              class="episode-item"
              :class="{ 'active': episode.id === currentEpisode?.id }"
              @click="selectEpisode(episode)"
            >
              <div class="episode-thumbnail">
                <img
                  :src="episode.thumbnail || currentSeries.cover_image"
                  :alt="episode.title"
                  @error="handleImageError"
                />
                <div class="episode-overlay">
                  <el-icon v-if="episode.id === currentEpisode?.id" class="play-icon">
                    <VideoPlay />
                  </el-icon>
                  <span class="episode-duration">{{ formatDuration(episode.duration) }}</span>
                </div>
              </div>

              <div class="episode-details">
                <div class="episode-number">第 {{ episode.episode_num || episode.episode_number }} 集</div>
                <div class="episode-title">{{ episode.title }}</div>
                <div class="episode-progress" v-if="episode.progress > 0">
                  <el-progress
                    :percentage="episode.progress"
                    :show-text="false"
                    :stroke-width="3"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 悬浮显示按钮（当侧边栏隐藏时） -->
        <div v-if="!showSidebar" class="floating-sidebar-btn">
          <el-button
            :icon="Menu"
            @click="toggleSidebar"
            type="primary"
            size="large"
            circle
            class="show-sidebar-btn"
          />
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading-state">
      <el-loading-directive v-loading="true" class="loading-container">
        <div class="loading-text">加载视频中...</div>
      </el-loading-directive>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useEnglishVideosStore } from '@/features/english-videos/stores/englishVideosStore'
import {
  Aim,
  ArrowDown,
  ArrowLeft,
  ArrowRight,
  FullScreen,
  Menu,
  Mute, RefreshRight,
  VideoPause,
  VideoPlay
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const englishVideosStore = useEnglishVideosStore()

// 响应式数据
const videoPlayer = ref(null)
const playerContainer = ref(null)
const currentSeries = ref(null)
const currentEpisode = ref(null)
const episodes = ref([])
const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const learningProgress = ref(0)
const showSubtitles = ref(true)
const currentSubtitle = ref('')
const repeatMode = ref(false)
const playbackRate = ref(1)
const isFullscreen = ref(false)
const showSidebar = ref(true)

// 计算属性
const isAuthenticated = computed(() => authStore.isAuthenticated)

const youtubeEmbedUrl = computed(() => {
  if (!currentEpisode.value?.video_url) return null

  const url = currentEpisode.value.video_url
  const videoId = extractYouTubeVideoId(url)

  if (videoId) {
    return `https://www.youtube.com/embed/${videoId}?autoplay=0&controls=1&rel=0&showinfo=0`
  }

  return null
})

const progressPercentage = computed({
  get: () => duration.value > 0 ? (currentTime.value / duration.value) * 100 : 0,
  set: (val) => {
    if (videoPlayer.value && duration.value > 0) {
      const newTime = (val / 100) * duration.value
      videoPlayer.value.currentTime = newTime
      currentTime.value = newTime
    }
  }
})

// 方法
const extractYouTubeVideoId = (url) => {
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
      return match[1]
    }
  }

  return null
}

const initializeVideo = async () => {
  const { seriesId, episodeId } = route.params

  try {
    // 获取系列信息
    const seriesData = await englishVideosStore.getVideoSeries(seriesId)
    currentSeries.value = seriesData

    // 获取剧集列表
    const episodesData = await englishVideosStore.getEpisodes(seriesId, {
      page_size: 100,
      sort_order: 'asc'
    })
    episodes.value = episodesData || []

    // 设置当前集
    let targetEpisode = null

    if (episodeId) {
      // 如果指定了集数，查找对应集数
      targetEpisode = episodes.value.find(ep => (ep.episode_num || ep.episode_number) == episodeId)
    } else {
      // 如果没有指定集数，尝试获取用户的最后播放记录
      if (isAuthenticated.value) {
        const lastWatchedEpisodeNumber = await findLastWatchedEpisodeNumber(seriesId)
        if (lastWatchedEpisodeNumber) {
          targetEpisode = episodes.value.find(ep => (ep.episode_num || ep.episode_number) === lastWatchedEpisodeNumber)
        }
      }

      // 如果没有播放记录，选择第一集
      if (!targetEpisode) {
        targetEpisode = episodes.value.find(ep => (ep.episode_num || ep.episode_number) === 1) || episodes.value[0]
      }
    }

    currentEpisode.value = targetEpisode || episodes.value[0]

    // 如果当前URL没有episodeId，更新URL
    if (!episodeId && currentEpisode.value) {
      router.replace({
        name: 'VideoPlayer',
        params: {
          seriesId: seriesId,
          episodeId: currentEpisode.value.episode_num || currentEpisode.value.episode_number
        }
      })
    }

    // 获取学习进度
    await loadLearningProgress()
  } catch (error) {
    console.error('Failed to initialize video:', error)
    ElMessage.error('加载视频失败')
  }
}

const findLastWatchedEpisodeNumber = async (seriesId) => {
  try {
    const progressData = await englishVideosStore.getUserProgress()
    const seriesProgress = progressData?.find(p => p.series_id == seriesId)

    if (seriesProgress && seriesProgress.last_episode_number) {
      // 返回最后观看的集数
      return seriesProgress.last_episode_number
    }
  } catch (error) {
    console.error('Failed to get last watched episode:', error)
  }
  return null
}

const loadLearningProgress = async () => {
  if (!currentEpisode.value || !isAuthenticated.value) return

  try {
    const progress = await englishVideosStore.getEpisodeProgress(currentEpisode.value.id)
    learningProgress.value = progress?.progress || 0

    // 恢复观看进度
    if (progress?.current_time > 0 && videoPlayer.value) {
      videoPlayer.value.currentTime = progress.current_time
    }
  } catch (error) {
    console.error('Failed to load progress:', error)
  }
}

const selectEpisode = async (episode) => {
  if (episode.id === currentEpisode.value?.id) return

  currentEpisode.value = episode
  learningProgress.value = episode.progress || 0

  // 加载新剧集的观看进度
  await loadLearningProgress()

  // 更新URL
  router.replace({
    name: 'VideoPlayer',
    params: {
      seriesId: route.params.seriesId,
      episodeId: episode.episode_num || episode.episode_number
    }
  })
}

const togglePlay = () => {
  if (videoPlayer.value) {
    if (isPlaying.value) {
      videoPlayer.value.pause()
    } else {
      videoPlayer.value.play()
    }
  }
}

const toggleMute = () => {
  if (videoPlayer.value) {
    videoPlayer.value.muted = !videoPlayer.value.muted
  }
}

const toggleSubtitles = () => {
  showSubtitles.value = !showSubtitles.value
}

const toggleRepeat = () => {
  repeatMode.value = !repeatMode.value
  if (videoPlayer.value) {
    videoPlayer.value.loop = repeatMode.value
  }
}

const setPlaybackRate = (rate) => {
  playbackRate.value = rate
  if (videoPlayer.value) {
    videoPlayer.value.playbackRate = rate
  }
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  
  // 控制body滚动
  if (isFullscreen.value) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}

const toggleSidebar = () => {
  showSidebar.value = !showSidebar.value
}

const handleVideoLoaded = () => {
  if (videoPlayer.value) {
    duration.value = videoPlayer.value.duration
  }
}

const handleTimeUpdate = () => {
  if (videoPlayer.value) {
    currentTime.value = videoPlayer.value.currentTime

    // 更新学习进度
    const progress = (currentTime.value / duration.value) * 100
    if (progress > learningProgress.value) {
      learningProgress.value = Math.min(progress, 100)

      // 每10%保存一次进度
      if (Math.floor(progress / 10) > Math.floor((learningProgress.value - 1) / 10)) {
        saveLearningProgress()
      }
    }
  }
}

const handleVideoPlay = () => {
  isPlaying.value = true
}

const handleVideoPause = () => {
  isPlaying.value = false
}

const handleVideoEnd = () => {
  isPlaying.value = false
  learningProgress.value = 100
  saveLearningProgress()

  // 自动播放下一集
  const currentIndex = episodes.value.findIndex(ep => ep.id === currentEpisode.value.id)
  if (currentIndex < episodes.value.length - 1) {
    selectEpisode(episodes.value[currentIndex + 1])
  }
}

const handleProgressChange = (val) => {
  progressPercentage.value = val
}

const saveLearningProgress = async () => {
  if (!currentEpisode.value || !isAuthenticated.value) return

  try {
    await englishVideosStore.updateEpisodeProgress(currentEpisode.value.id, {
      progress: Math.floor(learningProgress.value),
      current_time: Math.floor(currentTime.value),
      watch_time_minutes: Math.floor(currentTime.value / 60),
      is_completed: learningProgress.value >= 80
    })
  } catch (error) {
    console.error('Failed to save progress:', error)
  }
}

const goBack = () => {
  router.push({ name: 'EnglishVideos' })
}

const handleImageError = (e) => {
  e.target.src = '/default-video-cover.png'
}

const formatTime = (seconds) => {
  if (!seconds || isNaN(seconds)) return '0:00'
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = Math.floor(seconds % 60)
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

const formatDuration = (seconds) => {
  if (!seconds) return '未知'
  const minutes = Math.floor(seconds / 60)
  return `${minutes} 分钟`
}

const getProgressColor = (percentage) => {
  if (percentage >= 80) return '#67c23a'
  if (percentage >= 50) return '#e6a23c'
  return '#409eff'
}

// 监听路由变化
watch(() => route.params, initializeVideo, { immediate: true })

// 生命周期
onMounted(() => {
  // 添加键盘快捷键
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = '' // 恢复body滚动
  saveLearningProgress() // 离开时保存进度
})

const handleKeydown = (e) => {
  // 处理全屏模式下的ESC键
  if (e.key === 'Escape' && isFullscreen.value) {
    e.preventDefault()
    toggleFullscreen()
    return
  }

  // 只在非YouTube视频时处理键盘快捷键
  if (!videoPlayer.value || youtubeEmbedUrl.value) return

  switch (e.key) {
    case ' ':
      e.preventDefault()
      togglePlay()
      break
    case 'ArrowLeft':
      e.preventDefault()
      videoPlayer.value.currentTime = Math.max(0, currentTime.value - 10)
      break
    case 'ArrowRight':
      e.preventDefault()
      videoPlayer.value.currentTime = Math.min(duration.value, currentTime.value + 10)
      break
    case 'm':
    case 'M':
      toggleMute()
      break
    case 'f':
    case 'F':
      e.preventDefault()
      toggleFullscreen()
      break
  }
}
</script>

<style scoped>
.video-player-page {
  min-height: 100vh;
  height: 100vh;
  background: #f9f9f9;
  color: #212121;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.video-player-page.fullscreen-mode {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9999;
  background: #f8f9fa;
  margin: 0;
  padding: 0;
}

/* 全屏模式样式 */
.fullscreen-player {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  display: grid;
  grid-template-columns: 1fr 400px;
  grid-template-rows: 1fr;
  height: 100vh;
  width: 100vw;
  gap: 12px;
  padding: 12px;
  box-sizing: border-box;
  background: #000;
}

.fullscreen-video-container {
  grid-column: 1;
  grid-row: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e0e0e0;
  padding: 20px;
}

.fullscreen-video-wrapper {
  position: relative;
  width: 80%; /* 全屏模式下合适的大小 */
  max-width: 80%;
  padding-bottom: 45%; /* 80% * 56.25% = 16:9 aspect ratio */
  height: 0;
  background: #000;
  border-radius: 8px;
  overflow: hidden;
  margin: 0 auto; /* 居中显示 */
}

.fullscreen-video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.fullscreen-exit-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 10000;
  display: flex;
  align-items: center;
  gap: 8px;
  animation: fadeInSlide 0.5s ease-out;
}

.exit-fullscreen-btn {
  background: rgba(245, 108, 108, 0.9) !important;
  border: 2px solid rgba(255, 255, 255, 0.8) !important;
  color: white !important;
  backdrop-filter: blur(8px);
  box-shadow: 0 4px 16px rgba(245, 108, 108, 0.4);
  transition: all 0.3s ease;
}

.exit-fullscreen-btn:hover {
  background: rgba(245, 108, 108, 1) !important;
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(245, 108, 108, 0.6);
}

.exit-hint {
  color: white;
  background: rgba(0, 0, 0, 0.8);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  white-space: nowrap;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

@keyframes fadeInSlide {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.fullscreen-episodes {
  grid-column: 2;
  grid-row: 1;
  background: white;
  border-radius: 8px;
  padding: 12px;
  overflow-y: auto;
  border: 1px solid #e0e0e0;
}

.fullscreen-episodes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e0e0e0;
}

.fullscreen-episodes-header h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
  font-weight: 600;
}

.fullscreen-sidebar-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.fullscreen-sidebar-controls span {
  font-size: 12px;
  color: #666;
}

.fullscreen-hide-btn {
  color: #666 !important;
  transition: all 0.2s ease;
}

.fullscreen-hide-btn:hover {
  color: #409eff !important;
  transform: translateX(2px);
}

.fullscreen-episodes-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.fullscreen-episode-item {
  display: flex;
  gap: 12px;
  padding: 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: transparent;
}

.fullscreen-episode-item:hover {
  background: #e9ecef;
}

.fullscreen-episode-item.active {
  background: rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.3);
}

.fullscreen-episode-item img {
  width: 120px;
  height: 68px;
  border-radius: 4px;
  object-fit: cover;
  flex-shrink: 0;
}

.episode-info {
  flex: 1;
  min-width: 0;
}

.episode-info .episode-number {
  font-size: 12px;
  color: #606060;
  margin-bottom: 4px;
}

.episode-info .episode-title {
  font-size: 14px;
  font-weight: 400;
  color: #030303;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

/* 全屏模式悬浮侧边栏按钮 */
.fullscreen-floating-sidebar-btn {
  position: absolute;
  top: 50%;
  right: 20px;
  transform: translateY(-50%);
  z-index: 10000;
  animation: float 2s ease-in-out infinite;
}

.fullscreen-show-sidebar-btn {
  background: linear-gradient(135deg, #409eff, #67c23a) !important;
  border: none !important;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
  transition: all 0.3s ease;
}

.fullscreen-show-sidebar-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.6);
}

/* 页面内容容器 */
.player-page-content {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 顶部导航栏 */
.top-navigation {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 16px 20px;
  background: white;
  border-bottom: 1px solid #e0e0e0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
}

.back-btn {
  flex-shrink: 0;
}

.page-title {
  flex: 1;
  min-width: 0;
}

.page-title h1 {
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
  color: #212121;
}

.episode-indicator {
  font-size: 14px;
  color: #606060;
}

.player-container {
  display: grid;
  grid-template-columns: 1fr 300px; /* 缩小侧边栏宽度 */
  flex: 1;
  overflow: hidden;
  transition: grid-template-columns 0.3s ease;
  position: relative;
}

.player-container.sidebar-hidden {
  grid-template-columns: 1fr 0px;
}


.main-player-area {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  overflow-y: auto;
  height: 100%;
}

.video-player {
  background: #000;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
  width: 75%; /* 和播放器宽度一致 */
  margin: 0 auto; /* 居中显示 */
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
}

.fullscreen-controls {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
}

.fullscreen-btn {
  background: rgba(0, 0, 0, 0.7);
  color: white !important;
  border: none;
}

.fullscreen-btn:hover {
  background: rgba(0, 0, 0, 0.9);
}

/* 返回控制 */
.back-control {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 10;
}

.back-btn {
  background: rgba(0, 0, 0, 0.7);
  color: white !important;
  border: none;
}

.back-btn:hover {
  background: rgba(0, 0, 0, 0.9);
}

/* 视频全屏按钮 */
.video-fullscreen-btn {
  position: absolute;
  bottom: 16px;
  right: 16px;
  z-index: 10;
}

.fullscreen-toggle {
  background: rgba(0, 0, 0, 0.8) !important;
  border: none !important;
  color: white !important;
  transition: all 0.3s ease;
}

.fullscreen-toggle:hover {
  background: rgba(64, 158, 255, 0.8) !important;
  transform: scale(1.1);
}

/* 侧边栏控制 */
.sidebar-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sidebar-hide-btn {
  color: #ccc !important;
  transition: all 0.2s ease;
}

.sidebar-hide-btn:hover {
  color: #409eff !important;
  transform: translateX(2px);
}

/* 悬浮显示按钮 */
.floating-sidebar-btn {
  position: fixed;
  top: 50%;
  right: 20px;
  transform: translateY(-50%);
  z-index: 1000;
  animation: float 2s ease-in-out infinite;
}

.show-sidebar-btn {
  background: linear-gradient(135deg, #409eff, #67c23a) !important;
  border: none !important;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
  transition: all 0.3s ease;
}

.show-sidebar-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.6);
}

@keyframes float {
  0%, 100% {
    transform: translateY(-50%) scale(1);
  }
  50% {
    transform: translateY(-50%) scale(1.05);
  }
}

/* 侧边栏动画 */
.episodes-sidebar {
  transform: translateX(0);
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.player-container.sidebar-hidden .episodes-sidebar {
  transform: translateX(100%);
  opacity: 0;
}

.player-wrapper {
  position: relative;
  width: 100%; /* 填满主播放区域 */
  padding-bottom: 56.25%; /* 16:9 aspect ratio */
  height: 0;
  max-width: none;
  background: #000; /* 黑色背景 */
}


.video-element {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.custom-subtitle {
  position: absolute;
  bottom: 60px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 16px;
  text-align: center;
  max-width: 80%;
}

.player-controls {
  background: rgba(0, 0, 0, 0.9);
  padding: 16px;
}

.control-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.control-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.time-info {
  font-size: 14px;
  color: #ccc;
  min-width: 50px;
}

.progress-slider {
  flex: 1;
}

.video-info {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.episode-header {
  margin-bottom: 12px;
}

.episode-title {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: #212121;
  line-height: 1.3;
}

.episode-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 14px;
  color: #606060;
}

.episode-description {
  margin-bottom: 20px;
  line-height: 1.6;
  color: #424242;
  font-size: 15px;
}

.learning-progress-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin: 20px auto;
  width: 75%; /* 和播放器宽度保持一致 */
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.progress-header h3 {
  margin: 0;
  font-size: 16px;
  color: #212121;
  font-weight: 600;
}

.progress-text {
  font-weight: 600;
  color: #1976d2;
}

.episodes-sidebar {
  background: white;
  border-radius: 12px;
  padding: 16px;
  overflow-y: auto;
  margin: 0 16px 16px 0; /* 顶部对齐，去掉上边距 */
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  height: fit-content; /* 适应内容高度 */
  max-height: calc(100vh - 200px); /* 限制最大高度 */
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e0e0e0;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 18px;
  color: #212121;
  font-weight: 600;
}

.episode-count {
  font-size: 12px;
  color: #606060;
}

.episodes-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.episode-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: transparent;
}

.episode-item:hover {
  background: #f5f5f5;
}

.episode-item.active {
  background: #e3f2fd;
  border: 1px solid #1976d2;
}

.episode-thumbnail {
  position: relative;
  width: 70px; /* 缩小缩略图 */
  height: 39px; /* 保持16:9比例 */
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
}

.episode-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.episode-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, transparent 0%, rgba(0,0,0,0.6) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.play-icon {
  color: #409eff;
  font-size: 16px;
}

.episode-duration {
  position: absolute;
  bottom: 2px;
  right: 4px;
  font-size: 10px;
  color: white;
  background: rgba(0, 0, 0, 0.8);
  padding: 1px 3px;
  border-radius: 2px;
}

.episode-details {
  flex: 1;
  min-width: 0;
}

.episode-number {
  font-size: 12px;
  color: #606060;
  margin-bottom: 2px;
}

.episode-title {
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 4px;
  color: #212121;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.episode-progress {
  margin-top: 4px;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f9f9f9;
}

.loading-container {
  padding: 40px;
}

.loading-text {
  color: #212121;
  text-align: center;
  margin-top: 16px;
  font-size: 16px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .player-container {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr auto;
  }

  .episodes-sidebar {
    max-height: 300px;
    order: 3;
  }

  .episodes-list {
    flex-direction: row;
    overflow-x: auto;
    padding-bottom: 8px;
  }

  .episode-item {
    flex-direction: column;
    min-width: 120px;
    text-align: center;
  }
}

@media (max-width: 768px) {
  .main-player-area {
    padding: 12px;
  }

  .episode-title {
    font-size: 20px;
  }

  .episode-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>