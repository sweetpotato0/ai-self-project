<template>
  <div class="song-player-view">
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
              <span class="difficulty-label">难度：</span>
              <div class="difficulty-stars">
                <el-icon
                  v-for="i in 5"
                  :key="i"
                  :class="['star', { active: i <= (song.difficulty || 1) }]"
                >
                  <Star />
                </el-icon>
              </div>
              <span class="difficulty-text">{{ getDifficultyText(song.difficulty) }}</span>
            </div>

            <el-tag v-if="song.age_range" size="large" type="info">
              适合年龄：{{ song.age_range }}
            </el-tag>
          </div>

          <p v-if="song.description" class="song-description">
            {{ song.description }}
          </p>
        </div>
      </div>

      <!-- 主要内容区域：视频和歌词并排 -->
      <div class="main-content" v-if="song.video_url || song.audio_url || song.lyrics">
        <!-- 左侧：播放器区域 -->
        <div class="player-side" v-if="song.video_url || song.audio_url">
          <!-- 音频播放器 -->
          <div class="audio-section" v-if="song.audio_url">
            <h4>音频播放</h4>
            <audio
              ref="audioPlayer"
              :src="song.audio_url"
              controls
              preload="metadata"
              @play="handlePlay"
              @pause="handlePause"
              @ended="handleEnded"
              @timeupdate="handleTimeUpdate"
              @loadedmetadata="handleLoadedMetadata"
            />
          </div>

          <!-- 视频播放器 -->
          <div class="video-section" v-if="song.video_url">
            <h4>视频播放</h4>

            <!-- YouTube播放器 -->
            <div v-if="isYouTubeUrl(song.video_url)" class="youtube-player">
              <iframe
                :src="getYouTubeEmbedUrl(song.video_url)"
                class="youtube-iframe"
                frameborder="0"
                allowfullscreen
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              ></iframe>
            </div>

            <!-- 普通视频播放器 -->
            <video
              v-else
              ref="videoPlayer"
              :src="song.video_url"
              controls
              preload="metadata"
              class="standard-video"
              @play="handlePlay"
              @pause="handlePause"
              @ended="handleEnded"
              @timeupdate="handleTimeUpdate"
              @loadedmetadata="handleLoadedMetadata"
            />
          </div>
        </div>

        <!-- 右侧：歌词和推荐区域 -->
        <div class="right-side" :class="{ 'full-width': !song.video_url && !song.audio_url }">
          <!-- 歌词区域 -->
          <div class="lyrics-section" v-if="song.lyrics">
            <el-tabs v-model="lyricsTab" class="lyrics-tabs">
              <el-tab-pane
                v-if="song.lyrics"
                label="英文歌词"
                name="english"
              >
                <div class="lyrics-content">{{ song.lyrics }}</div>
              </el-tab-pane>

              <el-tab-pane
                v-if="song.lyrics_cn"
                label="中文翻译"
                name="chinese"
              >
                <div class="lyrics-content">{{ song.lyrics_cn }}</div>
              </el-tab-pane>

              <el-tab-pane
                v-if="song.lyrics && song.lyrics_cn"
                label="双语对照"
                name="bilingual"
              >
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

          <!-- 推荐歌曲区域 -->
          <div class="related-songs-section">
            <div class="section-header">
              <h4>相关推荐</h4>
              <el-button
                type="text"
                size="small"
                @click="refreshRecommendations"
                :loading="loadingRecommendations"
              >
                换一批
              </el-button>
            </div>

            <div class="related-songs-list" v-loading="loadingRecommendations">
              <div
                v-for="relatedSong in relatedSongs"
                :key="relatedSong.id"
                class="related-song-item"
                @click="switchToSong(relatedSong)"
              >
                <img
                  :src="relatedSong.cover_image || defaultCover"
                  :alt="relatedSong.title"
                  class="related-song-cover"
                />
                <div class="related-song-info">
                  <div class="related-song-title">{{ relatedSong.title }}</div>
                  <div class="related-song-meta">
                    <span class="category">{{ relatedSong.category?.name }}</span>
                    <div class="difficulty-stars">
                      <el-icon
                        v-for="i in (relatedSong.difficulty || 1)"
                        :key="i"
                        class="star active"
                      >
                        <Star />
                      </el-icon>
                    </div>
                  </div>
                </div>
                <el-button
                  type="primary"
                  :icon="VideoPlay"
                  circle
                  size="small"
                  class="play-btn"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 学习进度控制 -->
      <div class="progress-controls">
        <el-card shadow="never">
          <div class="progress-header">
            <h4>学习进度</h4>
            <div class="progress-info">
              <span class="progress-text">{{ Math.round(learningProgress) }}%</span>
              <span v-if="isPlaying" class="playing-indicator">
                🎵 正在播放
              </span>
              <div v-if="videoDuration > 0 || audioDuration > 0" class="time-info">
                {{ formatTime(currentTime) }} / {{ formatTime(videoDuration || audioDuration) }}
              </div>
            </div>
          </div>

          <el-progress
            :percentage="learningProgress"
            :stroke-width="8"
            :color="getProgressColor(learningProgress)"
          />

          <div class="progress-actions">
            <div class="progress-mode-info">
              <el-text size="small" type="info">
                💡 进度会根据播放时间自动更新
              </el-text>
            </div>
            <el-button-group>
              <el-button
                type="info"
                :icon="RefreshRight"
                @click="restartLearning"
                size="small"
                :disabled="learningProgress === 0"
              >
                重新学习
              </el-button>
              <el-button
                type="success"
                :icon="Plus"
                @click="addProgress(10)"
                size="small"
                :disabled="learningProgress >= 90"
              >
                +10%
              </el-button>
              <el-button
                type="warning"
                :icon="Check"
                @click="markCompleted"
                size="small"
                :disabled="learningProgress >= 100"
              >
                标记完成
              </el-button>
            </el-button-group>
          </div>
        </el-card>
      </div>

      <!-- 学习笔记 -->
      <div class="notes-section">
        <el-card shadow="never">
          <template #header>
            <div class="notes-header">
              <h4>学习笔记 ({{ learningNotesList.length }})</h4>
            </div>
          </template>

          <!-- 笔记列表 (GitHub 样式 - 列表在上面) -->
          <div class="notes-list">
            <div
              v-for="note in learningNotesList"
              :key="note.id"
              class="note-item"
            >
              <div class="note-content">
                <template v-if="editingNoteId === note.id">
                  <!-- 编辑模式 -->
                  <div class="note-edit-container">
                    <div class="note-tabs">
                      <el-tabs v-model="editPreviewMode" class="markdown-tabs">
                        <el-tab-pane label="编辑" name="edit">
                          <el-input
                            v-model="note.content"
                            type="textarea"
                            :rows="6"
                            resize="vertical"
                            class="markdown-input"
                            placeholder="支持 Markdown 格式..."
                            @keydown.esc="cancelNoteEdit"
                            @keydown.ctrl.enter="saveNoteEdit(note, note.content)"
                          />
                        </el-tab-pane>
                        <el-tab-pane label="预览" name="preview">
                          <div class="markdown-preview" v-html="renderMarkdown(note.content)"></div>
                        </el-tab-pane>
                      </el-tabs>
                    </div>
                    <div class="edit-note-actions">
                      <el-button
                        type="primary"
                        size="small"
                        @click="saveNoteEdit(note, note.content)"
                      >
                        保存
                      </el-button>
                      <el-button
                        size="small"
                        @click="cancelNoteEdit"
                      >
                        取消
                      </el-button>
                    </div>
                  </div>
                </template>
                <template v-else>
                  <!-- 显示模式 -->
                  <div class="note-markdown-content" v-html="renderMarkdown(note.content)"></div>
                  <div class="note-meta">
                    <span class="note-time">
                      {{ note.updatedAt !== note.createdAt ? '编辑于' : '创建于' }} {{ formatNoteTime(note.updatedAt) }}
                    </span>
                  </div>
                </template>
              </div>
              <div class="note-actions" v-if="editingNoteId !== note.id">
                <el-button
                  type="primary"
                  :icon="Edit"
                  size="small"
                  text
                  @click="editNote(note.id)"
                >
                  编辑
                </el-button>
                <el-popconfirm
                  title="确定删除这条笔记吗？"
                  @confirm="deleteNote(note.id)"
                >
                  <template #reference>
                    <el-button
                      type="danger"
                      :icon="Delete"
                      size="small"
                      text
                    >
                      删除
                    </el-button>
                  </template>
                </el-popconfirm>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-if="learningNotesList.length === 0" class="empty-notes">
              <p>还没有学习笔记，添加第一条笔记吧！</p>
            </div>
          </div>

          <!-- 添加新笔记 (GitHub 样式 - 输入框在下面) -->
          <div class="add-note-section">
            <div class="add-note-tabs">
              <el-tabs v-model="newNotePreviewMode" class="markdown-tabs">
                <el-tab-pane label="编辑" name="edit">
                  <el-input
                    v-model="newNoteContent"
                    type="textarea"
                    :rows="4"
                    placeholder="使用 Markdown 添加新的学习笔记..."
                    resize="vertical"
                    class="markdown-input"
                  />
                </el-tab-pane>
                <el-tab-pane label="预览" name="preview">
                  <div class="markdown-preview" v-html="renderMarkdown(newNoteContent)"></div>
                </el-tab-pane>
              </el-tabs>
            </div>
            <div class="add-note-actions">
              <span class="markdown-tip">支持 Markdown 语法</span>
              <el-button
                type="primary"
                size="small"
                @click="addNewNote"
                :icon="Plus"
                :disabled="!newNoteContent.trim()"
              >
                添加笔记
              </el-button>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading-state">
      <el-loading-directive v-loading="true" class="loading-container">
        <div class="loading-text">加载歌曲信息中...</div>
      </el-loading-directive>
    </div>
  </div>
</template>

<script setup>
import {
  ArrowLeft,
  Check,
  Delete,
  Edit,
  Plus,
  RefreshRight,
  Star
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import MarkdownIt from 'markdown-it'
import hljs from 'markdown-it-highlightjs'
import { onMounted, onUnmounted, ref } from 'vue'
import { useEnglishLearningStore } from '../stores/englishLearningStore'

const props = defineProps({
  song: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['back'])

// 状态管理
const store = useEnglishLearningStore()

// Markdown 配置
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true // 支持 GitHub 风格的换行
}).use(hljs)

// Markdown 渲染函数
const renderMarkdown = (content) => {
  if (!content) return '<p class="empty-preview">没有内容</p>'
  try {
    return md.render(content)
  } catch (error) {
    console.error('Markdown rendering error:', error)
    return `<p class="error-preview">Markdown 渲染错误</p>`
  }
}

// 响应式数据
const isLiked = ref(false)
const learningProgress = ref(0)
const learningNotes = ref('')
const learningNotesList = ref([]) // 新的多条笔记数组
const newNoteContent = ref('') // 新笔记内容
const editingNoteId = ref(null) // 正在编辑的笔记ID
const newNotePreviewMode = ref('edit') // 新笔记预览模式：edit/preview
const editPreviewMode = ref('edit') // 编辑笔记预览模式：edit/preview
const lyricsTab = ref('english')
const playCount = ref(0)
const studyTimeMinutes = ref(0)
const startTime = ref(Date.now())
const relatedSongs = ref([])
const loadingRecommendations = ref(false)

// 播放器引用
const audioPlayer = ref(null)
const videoPlayer = ref(null)

// 自动进度跟踪相关
const isPlaying = ref(false)
const videoProgress = ref(0)
const audioDuration = ref(0)
const videoDuration = ref(0)
const currentTime = ref(0)

const defaultCover = 'https://via.placeholder.com/300x300/667eea/ffffff?text=♪'

// 生命周期
onMounted(async () => {
  if (props.song) {
    await initializeSong()
    fetchRelatedSongs()
  }
})

onUnmounted(() => {
  // 保存学习时长
  const sessionTime = Math.floor((Date.now() - startTime.value) / 1000 / 60)
  studyTimeMinutes.value += sessionTime
  saveProgress()
})

// 方法
const initializeSong = async () => {
  // 初始化歌曲数据
  isLiked.value = props.song.is_liked || false

  // 获取用户的学习进度
  try {
    await store.fetchUserProgress()
    const userProgressData = store.userProgress.find(p => p.song_id === props.song.id)
    if (userProgressData) {
      learningProgress.value = userProgressData.progress || 0
      learningNotes.value = userProgressData.notes || ''
      // 解析笔记数据，支持新的数组格式和旧的字符串格式
      try {
        if (userProgressData.notes) {
          const parsedNotes = JSON.parse(userProgressData.notes)
          if (Array.isArray(parsedNotes)) {
            learningNotesList.value = parsedNotes
          } else {
            // 如果是字符串格式，转换为数组格式
            learningNotesList.value = [{
              id: Date.now(),
              content: userProgressData.notes,
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString()
            }]
          }
        } else {
          learningNotesList.value = []
        }
      } catch {
        // 如果解析失败，说明是旧的字符串格式
        learningNotesList.value = userProgressData.notes ? [{
          id: Date.now(),
          content: userProgressData.notes,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }] : []
      }
      playCount.value = userProgressData.play_count || 0
      studyTimeMinutes.value = userProgressData.study_time_minutes || 0
    } else {
      learningProgress.value = 0
      learningNotes.value = ''
      learningNotesList.value = []
    }
  } catch (error) {
    console.warn('Failed to fetch user progress:', error)
    learningProgress.value = 0
    learningNotes.value = ''
    learningNotesList.value = []
  }

  // 如果有歌词，默认显示英文
  if (props.song.lyrics) {
    lyricsTab.value = 'english'
  }
}

const goBack = () => {
  emit('back')
}

// YouTube相关方法
const isYouTubeUrl = (url) => {
  if (!url) return false
  const youtubeRegex = /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+/
  return youtubeRegex.test(url)
}

const getYouTubeEmbedUrl = (url) => {
  if (!url) return ''

  let videoId = ''

  // 处理不同格式的YouTube链接
  if (url.includes('youtu.be/')) {
    videoId = url.split('youtu.be/')[1]?.split('?')[0]
  } else if (url.includes('youtube.com/watch')) {
    const urlParams = new URLSearchParams(url.split('?')[1])
    videoId = urlParams.get('v')
  }

  if (videoId) {
    return `https://www.youtube.com/embed/${videoId}?rel=0&modestbranding=1&controls=1`
  }

  return ''
}

// 歌词处理
const getBilingualLyrics = () => {
  if (!props.song?.lyrics || !props.song?.lyrics_cn) return []

  const englishLines = props.song.lyrics.split('\n').filter(line => line.trim())
  const chineseLines = props.song.lyrics_cn.split('\n').filter(line => line.trim())

  const maxLines = Math.max(englishLines.length, chineseLines.length)
  const result = []

  for (let i = 0; i < maxLines; i++) {
    result.push({
      english: englishLines[i] || '',
      chinese: chineseLines[i] || ''
    })
  }

  return result
}

// 播放器事件处理 - 增强的自动进度跟踪
const handlePlay = (event) => {
  isPlaying.value = true
  if (playCount.value === 0) {
    playCount.value = 1
  }
  console.log('播放开始')
}

const handlePause = () => {
  isPlaying.value = false
  console.log('播放暂停')
}

const handleEnded = () => {
  isPlaying.value = false
  // 播放结束时自动设置进度为100%
  learningProgress.value = 100
  saveProgress()
  ElMessage.success('恭喜！您已完成这首歌的学习！')
  console.log('播放结束')
}

const handleTimeUpdate = (event) => {
  const player = event.target
  const duration = player.duration || 0
  const current = player.currentTime || 0

  if (duration > 0) {
    // 更新当前播放时间和总时长
    currentTime.value = current
    if (player.tagName === 'VIDEO') {
      videoDuration.value = duration
      videoProgress.value = (current / duration) * 100
    } else {
      audioDuration.value = duration
    }

    // 自动更新学习进度 - 基于播放时间计算
    const watchedProgress = (current / duration) * 100

    // 调试信息
    console.log(`Progress update: ${current.toFixed(2)}s / ${duration.toFixed(2)}s = ${watchedProgress.toFixed(2)}%`)

    // 只有当自动计算的进度大于当前进度时才更新（防止倒退）
    if (watchedProgress > learningProgress.value) {
      const previousProgress = learningProgress.value
      learningProgress.value = Math.min(watchedProgress, 100)

      console.log(`Learning progress updated from ${previousProgress.toFixed(2)}% to ${learningProgress.value.toFixed(2)}%`)

      // 每2%进度给用户提示
      const progressMilestone = Math.floor(learningProgress.value / 2) * 2
      const lastMilestone = Math.floor(previousProgress / 2) * 2

      if (progressMilestone > lastMilestone && progressMilestone > 0 && progressMilestone < 100) {
        ElMessage({
          message: `学习进度已达到 ${progressMilestone}%`,
          type: 'info',
          duration: 2000
        })

        // 在达到里程碑时保存进度
        saveProgress()
      }
    }
  }
}

// 监听播放器加载完成
const handleLoadedMetadata = (event) => {
  const player = event.target
  const duration = player.duration || 0

  if (player.tagName === 'VIDEO') {
    videoDuration.value = duration
  } else {
    audioDuration.value = duration
  }
}

// 进度管理
const addProgress = (amount) => {
  learningProgress.value = Math.min(learningProgress.value + amount, 100)
  saveProgress()
  ElMessage.success(`进度增加${amount}%`)
}

const markCompleted = () => {
  learningProgress.value = 100
  saveProgress()
  ElMessage.success('恭喜！您已完成这首歌的学习！')
}

const restartLearning = () => {
  learningProgress.value = 0
  studyTimeMinutes.value = 0
  playCount.value = 0
  saveProgress()
  ElMessage.success('已重置学习进度，开始重新学习！')
}

const saveProgress = async () => {
  if (!props.song) return

  const progressData = {
    progress: learningProgress.value,
    is_completed: learningProgress.value >= 100,
    play_count: playCount.value,
    study_time_minutes: studyTimeMinutes.value,
    notes: JSON.stringify(learningNotesList.value)
  }

  console.log('Saving progress:', progressData)

  try {
    await store.updateProgress(props.song.id, progressData)
    console.log('Progress saved successfully')
  } catch (error) {
    console.error('Failed to save progress:', error)
  }
}

const saveNotes = async () => {
  if (!props.song) return

  console.log('Saving notes:', learningNotes.value)

  try {
    await saveProgress()
    ElMessage.success('笔记已保存')
  } catch (error) {
    console.error('Failed to save notes:', error)
    ElMessage.error('保存笔记失败，请重试')
  }
}

// 新的多笔记管理功能
const addNewNote = async () => {
  if (!newNoteContent.value.trim()) {
    ElMessage.warning('请输入笔记内容')
    return
  }

  const newNote = {
    id: Date.now(),
    content: newNoteContent.value.trim(),
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }

  learningNotesList.value.push(newNote)
  newNoteContent.value = ''

  try {
    await saveProgress()
    ElMessage.success('笔记已添加')
  } catch (error) {
    console.error('Failed to save new note:', error)
    ElMessage.error('保存笔记失败，请重试')
  }
}

const editNote = (noteId) => {
  editingNoteId.value = noteId
  editPreviewMode.value = 'edit' // 进入编辑模式时默认显示编辑标签页
}

const saveNoteEdit = async (note, newContent) => {
  if (!newContent.trim()) {
    ElMessage.warning('笔记内容不能为空')
    return
  }

  note.content = newContent.trim()
  note.updatedAt = new Date().toISOString()
  editingNoteId.value = null

  try {
    await saveProgress()
    ElMessage.success('笔记已更新')
  } catch (error) {
    console.error('Failed to update note:', error)
    ElMessage.error('更新笔记失败，请重试')
  }
}

const cancelNoteEdit = () => {
  editingNoteId.value = null
}

const deleteNote = async (noteId) => {
  const noteIndex = learningNotesList.value.findIndex(note => note.id === noteId)
  if (noteIndex > -1) {
    learningNotesList.value.splice(noteIndex, 1)

    try {
      await saveProgress()
      ElMessage.success('笔记已删除')
    } catch (error) {
      console.error('Failed to delete note:', error)
      ElMessage.error('删除笔记失败，请重试')
    }
  }
}

const formatNoteTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  return date.toLocaleDateString()
}

const toggleLike = async () => {
  if (!props.song) return

  try {
    await store.likeSong(props.song.id)
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

// 格式化时间显示
const formatTime = (seconds) => {
  if (!seconds || isNaN(seconds)) return '0:00'
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = Math.floor(seconds % 60)
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

// 改进的智能推荐系统 - 多策略后备推荐
const fetchRelatedSongs = async () => {
  if (!props.song) return

  loadingRecommendations.value = true
  try {
    let recommendedSongs = []

    // 检查管理员权限 - 这里简化处理，实际项目中应该从用户信息获取
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    const isAdmin = user?.role === 'admin'


    // 策略1: 同分类歌曲推荐
    if (props.song?.category?.id) {
      try {
        const queryParams = {
          category_id: props.song.category.id,
          limit: 10
        }
        // 只有非管理员才只获取已发布的歌曲
        if (!isAdmin) {
          queryParams.is_published = true
        }

        await store.fetchSongs(queryParams)
        const categorySongs = (store.songs || []).filter(song => song.id !== props.song.id)
        recommendedSongs.push(...categorySongs)
      } catch (error) {
        console.warn('Category-based recommendation failed:', error)
      }
    }

    // 策略2: 同难度歌曲推荐（如果推荐不足）
    if (props.song?.difficulty && recommendedSongs.length < 6) {
      try {
        const queryParams = {
          difficulty: props.song.difficulty,
          limit: 8
        }
        if (!isAdmin) {
          queryParams.is_published = true
        }

        await store.fetchSongs(queryParams)
        const difficultySongs = (store.songs || []).filter(
          song => song.id !== props.song.id && !recommendedSongs.find(r => r.id === song.id)
        )
        recommendedSongs.push(...difficultySongs)
      } catch (error) {
        console.warn('Difficulty-based recommendation failed:', error)
      }
    }

    // 策略3: 热门歌曲推荐（如果推荐数量仍不足）
    if (recommendedSongs.length < 4) {
      try {
        const queryParams = {
          sort_by: 'view_count',
          sort_order: 'desc',
          limit: 10
        }
        if (!isAdmin) {
          queryParams.is_published = true
        }

        await store.fetchSongs(queryParams)
        const popularSongs = (store.songs || []).filter(
          song => song.id !== props.song.id && !recommendedSongs.find(r => r.id === song.id)
        )
        recommendedSongs.push(...popularSongs)
      } catch (error) {
        console.warn('Popular-based recommendation failed:', error)
      }
    }

    // 随机排序并限制数量
    const shuffled = [...recommendedSongs].sort(() => 0.5 - Math.random())
    relatedSongs.value = shuffled.slice(0, 4)

  } catch (error) {
    console.error('获取推荐歌曲失败:', error)
    relatedSongs.value = []
  } finally {
    loadingRecommendations.value = false
  }
}

const refreshRecommendations = () => {
  fetchRelatedSongs()
}

const switchToSong = (song) => {
  // 切换到新歌曲
  emit('back') // 先返回到歌曲列表，然后播放新歌曲
  setTimeout(() => {
    // 通过全局事件或其他方式播放新歌曲
    // 这里暂时返回到歌曲列表，用户可以重新选择
  }, 100)
}
</script>

<style scoped>
.song-player-view {
  min-height: 100%;
  background: #f5f5f5;
  padding: 0;
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
  width: 100%;
  margin: 0;
  padding: 20px;
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

.difficulty-label,
.difficulty-text {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.song-description {
  font-size: 16px;
  line-height: 1.6;
  color: #4b5563;
  margin: 0;
  font-size: 16px;
}

/* 主要内容区域 - 左右布局 */
.main-content {
  display: grid;
  grid-template-columns: 2.5fr 1fr; /* 视频区域稍微调小，歌词推荐区域占1份 */
  gap: 24px;
  align-items: start;
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 当没有播放器时，右侧区域占满全宽 */
.main-content:has(.right-side.full-width) {
  grid-template-columns: 1fr;
}

/* 兼容性更好的写法 - 当只有歌词时 */
.right-side.full-width {
  grid-column: 1 / -1;
}

.player-side {
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

.right-side {
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: 900px; /* 增加到1.5倍高度 (600 * 1.5) */
}

/* 当没有播放器时，右侧区域布局调整 */
.right-side.full-width {
  display: grid;
  grid-template-columns: 2fr 1fr; /* 歌词占2份，推荐占1份 */
  gap: 24px;
  height: auto; /* 自适应高度 */
}

.right-side.full-width .lyrics-section {
  height: auto;
  min-height: 500px; /* 设置最小高度 */
}

.right-side.full-width .related-songs-section {
  height: auto;
  min-height: 500px; /* 设置最小高度 */
}

.lyrics-section {
  background: #f9fafb;
  border-radius: 16px;
  padding: 20px;
  height: 525px; /* 增加到1.5倍高度 (350 * 1.5) */
  display: flex;
  flex-direction: column;
  overflow: hidden; /* 防止内容溢出 */
}

.related-songs-section {
  background: #f9fafb;
  border-radius: 16px;
  padding: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #374151;
}

.related-songs-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.related-song-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: white;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 2px solid transparent;
}

.related-song-item:hover {
  background: #e2e8f0;
  border-color: #667eea;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.related-song-cover {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  object-fit: cover;
  flex-shrink: 0;
}

.related-song-info {
  flex: 1;
  min-width: 0;
}

.related-song-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.related-song-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
}

.category {
  color: #6b7280;
  background: #e5e7eb;
  padding: 2px 6px;
  border-radius: 4px;
}

.play-btn {
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.related-song-item:hover .play-btn {
  opacity: 1;
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
  overflow-y: auto; /* 启用垂直滚动 */
  padding: 20px;
  background: white;
  border-radius: 12px;
  border-left: 4px solid #667eea;
  max-height: 420px; /* 增加到1.5倍高度 (280 * 1.5) */
  scrollbar-width: thin; /* Firefox 细滚动条 */
  scrollbar-color: #cbd5e1 #f1f5f9; /* Firefox 滚动条颜色 */
}

/* Webkit 浏览器滚动条样式 */
.lyrics-content::-webkit-scrollbar {
  width: 6px;
}

.lyrics-content::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 3px;
}

.lyrics-content::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

.lyrics-content::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

.bilingual-lyrics {
  flex: 1;
  overflow-y: auto; /* 启用垂直滚动 */
  padding-right: 8px;
  max-height: 420px; /* 增加到1.5倍高度 (280 * 1.5) */
  scrollbar-width: thin; /* Firefox 细滚动条 */
  scrollbar-color: #cbd5e1 #f1f5f9; /* Firefox 滚动条颜色 */
}

/* 双语歌词滚动条样式 */
.bilingual-lyrics::-webkit-scrollbar {
  width: 6px;
}

.bilingual-lyrics::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 3px;
}

.bilingual-lyrics::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

.bilingual-lyrics::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

.bilingual-line {
  margin-bottom: 20px;
  padding: 16px;
  background: white;
  border-radius: 12px;
  border-left: 4px solid #667eea;
}

.english-line {
  font-size: 16px;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 8px;
  line-height: 1.6;
}

.chinese-line {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
}

.progress-controls {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.progress-header h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #374151;
}

.progress-text {
  font-size: 24px;
  font-weight: 700;
  color: #667eea;
}

.progress-actions {
  margin-top: 20px;
  text-align: center;
}

.progress-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.playing-indicator {
  font-size: 12px;
  color: #67c23a;
  animation: pulse 1.5s ease-in-out infinite;
}

.time-info {
  font-size: 14px;
  color: #6b7280;
  font-family: monospace;
}

.progress-mode-info {
  margin-bottom: 12px;
  text-align: center;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.notes-section {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.notes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notes-header h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #374151;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: white;
  border-radius: 16px;
  margin: 20px;
}

.loading-container {
  width: 200px;
  height: 150px;
  align-self: center;
}

.loading-text {
  margin-top: 20px;
  text-align: center;
  color: #6b7280;
}

@media (max-width: 1200px) {
  /* 中等屏幕优化 */
  .main-content {
    gap: 20px;
    grid-template-columns: 2.2fr 1fr; /* 中等屏幕下适当调整比例 */
    padding: 16px;
  }

  .right-side {
    height: 500px;
  }

  .lyrics-section {
    height: 280px;
  }

  .player-container {
    padding: 16px;
  }

  .song-cover-large {
    width: 40px;
    height: 40px;
    align-self: center;
  }
}

@media (max-width: 768px) {
  .player-container {
    padding: 12px;
    gap: 16px;
  }

  .song-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 16px;
  }

  .song-cover-large {
    width: 40px;
    height: 40px;
    align-self: center;
  }
  /* 移动端改为垂直布局 */
  .main-content {
    grid-template-columns: 1fr;
    gap: 16px;
    padding: 16px;
  }

  .right-side {
    height: auto;
    gap: 16px;
  }

  /* 移动端没有播放器时也改为垂直布局 */
  .right-side.full-width {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .right-side.full-width .lyrics-section,
  .right-side.full-width .related-songs-section {
    min-height: 300px; /* 移动端减少最小高度 */
  }

  .lyrics-section {
    height: 450px; /* 移动端1.5倍高度 (300 * 1.5) */
  }

  .lyrics-content {
    max-height: 330px; /* 移动端1.5倍高度 (220 * 1.5) */
  }

  .bilingual-lyrics {
    max-height: 330px; /* 移动端1.5倍高度 (220 * 1.5) */
  }

  .related-songs-section {
    height: 250px; /* 移动端推荐区域高度 */
  }
  .progress-controls {
    padding: 16px;
  }

  .player-header {
    padding: 12px;
  }

  .song-title-header h1 {
    font-size: 20px;
  }
}

/* 多笔记系统样式 - GitHub 风格 */
.add-note-section {
  margin-top: 20px;
  border-top: 1px solid #ebeef5;
  padding-top: 16px;
}

.add-note-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.markdown-tip {
  font-size: 12px;
  color: #909399;
  font-style: italic;
}

.notes-list {
  max-height: 400px;
  overflow-y: auto;
}

.note-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  margin-bottom: 12px;
  background: #fafbfc;
  transition: all 0.2s ease;
}

.note-item:hover {
  border-color: #c6e2ff;
  background: #f5f9ff;
}

.note-content {
  flex: 1;
  min-width: 0;
}

/* Markdown 内容样式 */
.note-markdown-content,
.markdown-preview {
  font-size: 14px;
  line-height: 1.6;
  color: #333;
  margin-bottom: 8px;
  word-wrap: break-word;
}

.note-markdown-content h1,
.note-markdown-content h2,
.note-markdown-content h3,
.markdown-preview h1,
.markdown-preview h2,
.markdown-preview h3 {
  margin: 16px 0 8px 0;
  font-weight: 600;
}

.note-markdown-content h1,
.markdown-preview h1 {
  font-size: 18px;
  border-bottom: 1px solid #ebeef5;
  padding-bottom: 4px;
}

.note-markdown-content h2,
.markdown-preview h2 {
  font-size: 16px;
}

.note-markdown-content h3,
.markdown-preview h3 {
  font-size: 14px;
}

.note-markdown-content p,
.markdown-preview p {
  margin: 8px 0;
}

.note-markdown-content code,
.markdown-preview code {
  background: #f6f8fa;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 12px;
}

.note-markdown-content pre,
.markdown-preview pre {
  background: #f6f8fa;
  border-radius: 6px;
  padding: 12px;
  overflow-x: auto;
  margin: 12px 0;
}

.note-markdown-content pre code,
.markdown-preview pre code {
  background: none;
  padding: 0;
}

.note-markdown-content blockquote,
.markdown-preview blockquote {
  border-left: 4px solid #dfe2e5;
  padding-left: 12px;
  margin: 12px 0;
  color: #666;
}

.note-markdown-content ul,
.note-markdown-content ol,
.markdown-preview ul,
.markdown-preview ol {
  margin: 8px 0;
  padding-left: 24px;
}

.note-markdown-content li,
.markdown-preview li {
  margin: 4px 0;
}

.empty-preview,
.error-preview {
  color: #909399;
  font-style: italic;
  text-align: center;
  padding: 20px;
}

.error-preview {
  color: #f56c6c;
}

/* Markdown 编辑器标签页样式 */
.markdown-tabs {
  margin-bottom: 0;
}

.markdown-tabs .el-tabs__header {
  margin-bottom: 8px;
}

.markdown-input textarea {
  font-family: 'Monaco', 'Consolas', monospace !important;
  font-size: 13px !important;
  line-height: 1.5 !important;
}

.markdown-preview {
  min-height: 120px;
  padding: 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: #fafafa;
}

/* 编辑容器 */
.note-edit-container {
  width: 100%;
}

.note-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.note-time {
  font-size: 12px;
  color: #909399;
}

.note-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.edit-note-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

.empty-notes {
  text-align: center;
  color: #909399;
  padding: 32px 16px;
  font-style: italic;
}

.empty-notes p {
  margin: 0;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .note-item {
    flex-direction: column;
    align-items: stretch;
  }

  .note-actions {
    align-self: flex-end;
    margin-top: 8px;
  }
}
</style>