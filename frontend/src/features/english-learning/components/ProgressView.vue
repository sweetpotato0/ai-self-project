<template>
  <div class="progress-view">
    <div class="view-header">
      <h3 class="view-title">学习进度</h3>
      <p class="view-subtitle">跟踪您的英语学习进程</p>
    </div>

    <!-- 学习概览统计 -->
    <div class="progress-overview" v-if="stats">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-card class="overview-card">
            <div class="overview-content">
              <div class="overview-icon completed">
                <el-icon><Check /></el-icon>
              </div>
              <div class="overview-info">
                <div class="overview-number">{{ stats.completed_songs || 0 }}</div>
                <div class="overview-label">已完成歌曲</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="overview-card">
            <div class="overview-content">
              <div class="overview-icon study-time">
                <el-icon><Timer /></el-icon>
              </div>
              <div class="overview-info">
                <div class="overview-number">{{ stats.total_study_minutes || 0 }}</div>
                <div class="overview-label">学习分钟</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="overview-card">
            <div class="overview-content">
              <div class="overview-icon streak">
                <el-icon><Promotion /></el-icon>
              </div>
              <div class="overview-info">
                <div class="overview-number">{{ stats.current_streak || 0 }}</div>
                <div class="overview-label">连续天数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="overview-card">
            <div class="overview-content">
              <div class="overview-icon level">
                <el-icon><Trophy /></el-icon>
              </div>
              <div class="overview-info">
                <div class="overview-number">{{ stats.level || 1 }}</div>
                <div class="overview-label">当前等级</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 最喜欢的分类 -->
    <div class="favorite-category" v-if="stats && stats.favorite_category">
      <el-card>
        <template #header>
          <div class="card-header">
            <el-icon><Star /></el-icon>
            <span>最喜欢的分类</span>
          </div>
        </template>
        <div class="favorite-info">
          <div class="category-name">{{ stats.favorite_category }}</div>
          <p class="category-description">您在这个分类上花费了最多的学习时间</p>
        </div>
      </el-card>
    </div>

    <!-- 学习进度列表 -->
    <div class="progress-list" v-loading="loading">
      <el-card>
        <template #header>
          <div class="card-header">
            <el-icon><TrendCharts /></el-icon>
            <span>歌曲学习进度</span>
            <el-button type="text" @click="$emit('refresh')">刷新</el-button>
          </div>
        </template>

        <div class="progress-filters">
          <el-radio-group v-model="progressFilter" @change="filterProgress">
            <el-radio-button value="all">全部</el-radio-button>
            <el-radio-button value="learning">学习中</el-radio-button>
            <el-radio-button value="completed">已完成</el-radio-button>
            <el-radio-button value="favorite">已收藏</el-radio-button>
          </el-radio-group>
        </div>

        <div class="progress-items">
          <div 
            v-for="item in filteredProgress" 
            :key="item.id"
            class="progress-item"
          >
            <!-- 歌曲信息 -->
            <div class="song-info">
              <img 
                :src="item.song?.cover_image || defaultCover" 
                :alt="item.song?.title"
                class="song-thumbnail"
              />
              <div class="song-details">
                <h4 class="song-title">{{ item.song?.title }}</h4>
                <p class="song-category">{{ item.song?.category?.name }}</p>
                <div class="song-meta">
                  <el-tag size="small">{{ getDifficultyText(item.song?.difficulty) }}</el-tag>
                  <el-tag size="small" type="info">{{ item.song?.age_range }}</el-tag>
                </div>
              </div>
            </div>

            <!-- 进度信息 -->
            <div class="progress-info">
              <div class="progress-stats">
                <div class="stat-group">
                  <div class="stat-label">进度</div>
                  <div class="progress-bar-container">
                    <el-progress 
                      :percentage="item.progress || 0"
                      :color="getProgressColor(item.progress || 0)"
                      :stroke-width="8"
                    />
                  </div>
                </div>
                
                <div class="stat-group">
                  <div class="stat-label">播放次数</div>
                  <div class="stat-value">{{ item.play_count || 0 }}次</div>
                </div>
                
                <div class="stat-group">
                  <div class="stat-label">学习时长</div>
                  <div class="stat-value">{{ item.study_time_minutes || 0 }}分钟</div>
                </div>
              </div>

              <!-- 学习状态 -->
              <div class="learning-status">
                <el-tag 
                  :type="item.is_completed ? 'success' : 'warning'"
                  size="large"
                >
                  {{ item.is_completed ? '已完成' : '学习中' }}
                </el-tag>
                
                <el-tag v-if="item.is_liked" type="danger" size="large">
                  <el-icon><Star /></el-icon>
                  已收藏
                </el-tag>
              </div>
            </div>

            <!-- 最后学习时间 -->
            <div class="last-studied" v-if="item.last_studied_at">
              <div class="time-info">
                <el-icon><Clock /></el-icon>
                <span>{{ formatLastStudied(item.last_studied_at) }}</span>
              </div>
            </div>

            <!-- 学习笔记 -->
            <div class="learning-notes" v-if="item.notes">
              <div class="notes-header" @click="toggleNotes(item.song_id)">
                <el-icon><Document /></el-icon>
                <span>学习笔记</span>
                <el-icon class="collapse-icon" :class="{ expanded: expandedNotes.has(item.song_id) }">
                  <ArrowDown />
                </el-icon>
              </div>
              <el-collapse-transition>
                <p class="notes-content" v-if="expandedNotes.has(item.song_id)">{{ item.notes }}</p>
              </el-collapse-transition>
            </div>

            <!-- 操作按钮 -->
            <div class="progress-actions">
              <el-button 
                type="primary" 
                size="small"
                @click="continueLearning(item)"
              >
                <el-icon><VideoPlay /></el-icon>
                继续学习
              </el-button>
              
              <el-button 
                size="small"
                @click="editNotes(item)"
              >
                <el-icon><Edit /></el-icon>
                编辑笔记
              </el-button>
              
              <el-button 
                v-if="!item.is_completed"
                type="success" 
                size="small"
                @click="markCompleted(item)"
              >
                <el-icon><Check /></el-icon>
                标记完成
              </el-button>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="!loading && filteredProgress.length === 0" class="empty-state">
            <el-empty description="暂无学习进度">
              <el-button type="primary" @click="$router.push('/dashboard/english-learning')">
                开始学习
              </el-button>
            </el-empty>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 编辑笔记对话框 -->
    <el-dialog 
      v-model="showNotesDialog" 
      title="编辑学习笔记" 
      width="500px"
    >
      <el-form>
        <el-form-item label="学习笔记">
          <el-input
            v-model="editingNotes"
            type="textarea"
            :rows="4"
            placeholder="记录您的学习心得..."
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showNotesDialog = false">取消</el-button>
        <el-button type="primary" @click="saveNotes">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Check, Timer, Trophy, Star, TrendCharts, 
  Clock, Document, VideoPlay, Edit, Promotion, ArrowDown
} from '@element-plus/icons-vue'

const props = defineProps({
  progress: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  stats: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['refresh', 'continue-learning', 'update-progress'])

// 响应式数据
const progressFilter = ref('all')
const showNotesDialog = ref(false)
const editingNotes = ref('')
const editingItem = ref(null)
const expandedNotes = ref(new Set()) // 存储展开的笔记ID

const defaultCover = '/images/default-song-cover.jpg'

// 计算属性
const filteredProgress = computed(() => {
  let filtered = props.progress

  switch (progressFilter.value) {
    case 'learning':
      filtered = filtered.filter(item => !item.is_completed)
      break
    case 'completed':
      filtered = filtered.filter(item => item.is_completed)
      break
    case 'favorite':
      filtered = filtered.filter(item => item.is_liked)
      break
  }

  return filtered.sort((a, b) => {
    // 优先显示最近学习的
    if (a.last_studied_at && b.last_studied_at) {
      return new Date(b.last_studied_at) - new Date(a.last_studied_at)
    }
    return b.updated_at - a.updated_at
  })
})

// 方法
const filterProgress = () => {
  // 筛选逻辑在computed中处理
}

const continueLearning = (item) => {
  emit('continue-learning', item.song)
}

const toggleNotes = (songId) => {
  if (expandedNotes.value.has(songId)) {
    expandedNotes.value.delete(songId)
  } else {
    expandedNotes.value.add(songId)
  }
}

const editNotes = (item) => {
  editingItem.value = item
  editingNotes.value = item.notes || ''
  showNotesDialog.value = true
}

const saveNotes = async () => {
  if (!editingItem.value) return

  try {
    await emit('update-progress', editingItem.value.song_id, {
      notes: editingNotes.value
    })
    editingItem.value.notes = editingNotes.value
    showNotesDialog.value = false
    ElMessage.success('笔记保存成功')
  } catch (error) {
    ElMessage.error('笔记保存失败')
  }
}

const markCompleted = async (item) => {
  try {
    await emit('update-progress', item.song_id, {
      is_completed: true,
      progress: 100
    })
    item.is_completed = true
    item.progress = 100
    ElMessage.success('已标记为完成')
  } catch (error) {
    ElMessage.error('操作失败')
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

const formatLastStudied = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date
  
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
    return date.toLocaleDateString()
  }
}
</script>

<style scoped>
.progress-view {
  padding: 24px;
}

.view-header {
  margin-bottom: 24px;
}

.view-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.view-subtitle {
  color: #6b7280;
  margin: 0;
}

.progress-overview {
  margin-bottom: 24px;
}

.overview-card {
  height: 100%;
}

.overview-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.overview-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
}

.overview-icon.completed {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.overview-icon.study-time {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.overview-icon.streak {
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
}

.overview-icon.level {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.overview-number {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
}

.overview-label {
  font-size: 14px;
  color: #6b7280;
  margin-top: 4px;
}

.favorite-category {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.favorite-info {
  text-align: center;
}

.category-name {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.category-description {
  color: #6b7280;
  margin: 0;
}

.progress-filters {
  margin-bottom: 20px;
}

.progress-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.progress-item {
  padding: 20px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #f9fafb;
  transition: all 0.2s ease;
}

.progress-item:hover {
  background: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.song-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.song-thumbnail {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  object-fit: cover;
}

.song-details {
  flex: 1;
}

.song-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
}

.song-category {
  color: #6b7280;
  margin: 0 0 8px 0;
  font-size: 14px;
}

.song-meta {
  display: flex;
  gap: 8px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.progress-stats {
  display: flex;
  gap: 32px;
}

.stat-group {
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 8px;
}

.progress-bar-container {
  width: 120px;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.learning-status {
  display: flex;
  gap: 8px;
}

.last-studied {
  margin-bottom: 16px;
}

.time-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #6b7280;
}

.learning-notes {
  margin-bottom: 16px;
}

.notes-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
  cursor: pointer;
  user-select: none;
  transition: color 0.2s ease;
}

.notes-header:hover {
  color: #667eea;
}

.collapse-icon {
  margin-left: auto;
  transition: transform 0.3s ease;
}

.collapse-icon.expanded {
  transform: rotate(180deg);
}

.notes-content {
  color: #6b7280;
  line-height: 1.5;
  margin: 0;
  padding: 12px;
  background: white;
  border-radius: 8px;
}

.progress-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

@media (max-width: 768px) {
  .progress-view {
    padding: 16px;
  }

  .progress-overview .el-col {
    margin-bottom: 16px;
  }

  .progress-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .progress-stats {
    flex-direction: column;
    gap: 16px;
    width: 100%;
  }

  .progress-actions {
    justify-content: center;
  }
}
</style>