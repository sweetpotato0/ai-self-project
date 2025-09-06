<template>
  <div class="video-player">
    <div class="player-container">
      <!-- 视频播放器 -->
      <div class="video-wrapper">
        <div v-if="episode.video_url" class="video-content">
          <!-- YouTube 嵌入播放器 -->
          <iframe 
            v-if="isYouTube"
            :src="youtubeEmbedUrl"
            title="YouTube video player" 
            frameborder="0" 
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" 
            allowfullscreen
            @load="onVideoLoad"
            ref="videoFrame"
          ></iframe>
          
          <!-- 原生视频播放器 -->
          <video 
            v-else
            ref="videoElement"
            :src="episode.video_url"
            controls
            preload="metadata"
            @loadedmetadata="onVideoLoad"
            @timeupdate="onTimeUpdate"
            @ended="onVideoEnded"
            @pause="onVideoPause"
            @play="onVideoPlay"
          >
            <track 
              v-if="episode.subtitles" 
              kind="subtitles" 
              :src="episode.subtitles" 
              srclang="en" 
              label="English"
            >
          </video>
        </div>
        
        <!-- 缩略图占位符 -->
        <div v-else class="video-placeholder">
          <el-image
            :src="episode.thumbnail"
            fit="cover"
            style="width: 100%; height: 100%;"
          >
            <template #error>
              <div class="image-error">
                <el-icon><VideoPlay /></el-icon>
                <p>视频暂不可用</p>
              </div>
            </template>
          </el-image>
        </div>
      </div>
      
      <!-- 播放控制栏 -->
      <div class="player-controls">
        <div class="progress-info">
          <div class="progress-bar">
            <el-progress 
              :percentage="progressPercentage" 
              :show-text="false"
              :stroke-width="6"
              color="#409eff"
            />
            <span class="time-display">
              {{ formatTime(currentTime) }} / {{ formatTime(episode.duration) }}
            </span>
          </div>
          
          <div class="control-buttons">
            <el-button 
              size="small" 
              :type="isPlaying ? 'warning' : 'primary'"
              :icon="isPlaying ? VideoPause : VideoPlay"
              @click="togglePlay"
              :disabled="!episode.video_url"
            >
              {{ isPlaying ? '暂停' : '播放' }}
            </el-button>
            
            <el-button 
              size="small" 
              :icon="isCompleted ? 'SuccessFilled' : 'CircleCheck'"
              :type="isCompleted ? 'success' : 'info'"
              @click="toggleCompleted"
            >
              {{ isCompleted ? '已完成' : '标记完成' }}
            </el-button>
            
            <el-button 
              size="small" 
              :icon="isFavorited ? 'StarFilled' : 'Star'"
              :type="isFavorited ? 'warning' : 'info'"
              @click="toggleFavorite"
            >
              {{ isFavorited ? '已收藏' : '收藏' }}
            </el-button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 学习笔记 -->
    <div class="learning-notes">
      <el-card>
        <template #header>
          <div class="notes-header">
            <span>学习笔记</span>
            <el-button size="small" @click="saveNotes" :loading="savingNotes">
              保存笔记
            </el-button>
          </div>
        </template>
        
        <el-input
          v-model="userNotes"
          type="textarea"
          :rows="4"
          placeholder="在这里记录你的学习笔记..."
          @blur="autoSaveNotes"
        />
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { VideoPlay, VideoPause, Star, StarFilled, CircleCheck, SuccessFilled } from '@element-plus/icons-vue'
import { useEnglishVideosStore } from '../stores/englishVideosStore'

const props = defineProps({
  episode: {
    type: Object,
    required: true
  },
  seriesId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['progress-update', 'episode-complete'])

const englishVideosStore = useEnglishVideosStore()

// 响应式数据
const videoElement = ref(null)
const videoFrame = ref(null)
const currentTime = ref(0)
const isPlaying = ref(false)
const isCompleted = ref(false)
const isFavorited = ref(false)
const userNotes = ref('')
const savingNotes = ref(false)
const progressTimer = ref(null)
const watchStartTime = ref(0)
const totalWatchTime = ref(0)

// 计算属性
const progressPercentage = computed(() => {
  if (!props.episode.duration) return 0
  return Math.round((currentTime.value / props.episode.duration) * 100)
})

const isYouTube = computed(() => {
  const url = props.episode.video_url || ''
  return url.includes('youtube.com') || url.includes('youtu.be')
})

const youtubeEmbedUrl = computed(() => {
  if (!isYouTube.value) return ''
  
  const url = props.episode.video_url
  let videoId = ''
  
  if (url.includes('youtube.com/watch?v=')) {
    videoId = url.split('v=')[1].split('&')[0]
  } else if (url.includes('youtu.be/')) {
    videoId = url.split('youtu.be/')[1].split('?')[0]
  } else {
    // 假设直接是视频ID
    videoId = url
  }
  
  return `https://www.youtube.com/embed/${videoId}?enablejsapi=1&origin=${window.location.origin}`
})

// 方法
const loadEpisodeProgress = async () => {
  try {
    const progress = await englishVideosStore.getEpisodeProgress(props.episode.id)
    if (progress) {
      currentTime.value = progress.current_time || 0
      isCompleted.value = progress.is_completed || false
      userNotes.value = progress.notes || ''
      totalWatchTime.value = progress.watch_time_minutes || 0
    }
  } catch (error) {
    console.error('Failed to load episode progress:', error)
  }
}

const updateProgress = async (progressData) => {
  try {
    await englishVideosStore.updateEpisodeProgress(props.episode.id, {
      progress: progressPercentage.value,
      current_time: Math.floor(currentTime.value),
      watch_time_minutes: Math.floor(totalWatchTime.value),
      is_completed: isCompleted.value,
      notes: userNotes.value,
      ...progressData
    })
    
    emit('progress-update', {
      episodeId: props.episode.id,
      progress: progressPercentage.value,
      isCompleted: isCompleted.value
    })
  } catch (error) {
    console.error('Failed to update progress:', error)
    ElMessage.error('保存进度失败')
  }
}

const onVideoLoad = () => {
  if (videoElement.value && currentTime.value > 0) {
    videoElement.value.currentTime = currentTime.value
  }
}

const onTimeUpdate = () => {
  if (videoElement.value) {
    currentTime.value = videoElement.value.currentTime
    
    // 每30秒自动保存进度
    if (Math.floor(currentTime.value) % 30 === 0) {
      updateProgress()
    }
  }
}

const onVideoPlay = () => {
  isPlaying.value = true
  watchStartTime.value = Date.now()
  
  // 启动观看时长计时器
  progressTimer.value = setInterval(() => {
    totalWatchTime.value += 1/60 // 每秒增加1/60分钟
  }, 1000)
}

const onVideoPause = () => {
  isPlaying.value = false
  
  // 计算本次观看时长并累加
  if (watchStartTime.value > 0) {
    const sessionTime = (Date.now() - watchStartTime.value) / 1000 / 60 // 转换为分钟
    totalWatchTime.value += sessionTime
    watchStartTime.value = 0
  }
  
  // 清除计时器
  if (progressTimer.value) {
    clearInterval(progressTimer.value)
    progressTimer.value = null
  }
  
  // 保存进度
  updateProgress()
}

const onVideoEnded = () => {
  isPlaying.value = false
  isCompleted.value = true
  
  // 清除计时器
  if (progressTimer.value) {
    clearInterval(progressTimer.value)
    progressTimer.value = null
  }
  
  // 保存完成状态
  updateProgress({ is_completed: true })
  
  // 通知父组件视频已完成
  emit('episode-complete', props.episode.id)
  
  ElMessage.success('恭喜完成本集学习！')
}

const togglePlay = () => {
  if (!videoElement.value) return
  
  if (isPlaying.value) {
    videoElement.value.pause()
  } else {
    videoElement.value.play()
  }
}

const toggleCompleted = () => {
  isCompleted.value = !isCompleted.value
  updateProgress({ is_completed: isCompleted.value })
  
  if (isCompleted.value) {
    emit('episode-complete', props.episode.id)
    ElMessage.success('已标记为完成')
  } else {
    ElMessage.info('已取消完成标记')
  }
}

const toggleFavorite = async () => {
  try {
    if (props.seriesId) {
      await englishVideosStore.toggleSeriesLike(props.seriesId)
      isFavorited.value = !isFavorited.value
      ElMessage.success(isFavorited.value ? '已添加到收藏' : '已取消收藏')
    }
  } catch (error) {
    ElMessage.error('收藏操作失败')
  }
}

const saveNotes = async () => {
  if (savingNotes.value) return
  
  savingNotes.value = true
  try {
    await updateProgress({ notes: userNotes.value })
    ElMessage.success('笔记保存成功')
  } catch (error) {
    ElMessage.error('笔记保存失败')
  } finally {
    savingNotes.value = false
  }
}

const autoSaveNotes = () => {
  // 自动保存笔记
  updateProgress({ notes: userNotes.value })
}

const formatTime = (seconds) => {
  if (!seconds) return '00:00'
  
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  
  if (hours > 0) {
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  } else {
    return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
}

// 生命周期钩子
onMounted(() => {
  loadEpisodeProgress()
})

onUnmounted(() => {
  // 清理定时器
  if (progressTimer.value) {
    clearInterval(progressTimer.value)
  }
  
  // 保存最终进度
  updateProgress()
})

// 监听集数变化
watch(() => props.episode.id, () => {
  // 清理状态
  currentTime.value = 0
  isPlaying.value = false
  isCompleted.value = false
  userNotes.value = ''
  totalWatchTime.value = 0
  
  // 重新加载进度
  loadEpisodeProgress()
})
</script>

<style scoped>
.video-player {
  max-width: 1200px;
  margin: 0 auto;
}

.player-container {
  margin-bottom: 20px;
}

.video-wrapper {
  position: relative;
  width: 100%;
  height: 0;
  padding-bottom: 56.25%; /* 16:9 宽高比 */
  background: #000;
  border-radius: 8px;
  overflow: hidden;
}

.video-content {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.video-content iframe,
.video-content video {
  width: 100%;
  height: 100%;
}

.video-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #909399;
}

.image-error .el-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.player-controls {
  margin-top: 15px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
}

.progress-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.progress-bar {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 15px;
}

.progress-bar .el-progress {
  flex: 1;
}

.time-display {
  font-size: 14px;
  color: #606266;
  white-space: nowrap;
}

.control-buttons {
  display: flex;
  gap: 10px;
}

.learning-notes {
  margin-top: 20px;
}

.notes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

@media (max-width: 768px) {
  .progress-info {
    flex-direction: column;
    gap: 15px;
  }
  
  .progress-bar {
    width: 100%;
  }
  
  .control-buttons {
    width: 100%;
    justify-content: center;
    flex-wrap: wrap;
  }
}
</style>