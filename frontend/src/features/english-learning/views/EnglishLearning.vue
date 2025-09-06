<template>
  <div class="english-learning-container">
    <div class="english-learning-header">
      <div class="header-left">
        <h2 class="page-title">英语学习</h2>
        <p class="page-subtitle">通过歌曲快乐学习英语</p>
      </div>
      <div class="header-right">
        <el-button-group>
          <el-button
            :type="currentView === 'categories' ? 'primary' : 'default'"
            @click="switchTab('categories')"
          >
            <el-icon><List /></el-icon>
            分类
          </el-button>
          <el-button
            :type="currentView === 'songs' ? 'primary' : 'default'"
            @click="switchTab('songs')"
          >
            <el-icon><VideoPlay /></el-icon>
            歌曲
          </el-button>
          <el-button
            :type="currentView === 'progress' ? 'primary' : 'default'"
            @click="switchTab('progress')"
          >
            <el-icon><TrendCharts /></el-icon>
            进度
          </el-button>
          <el-button
            v-if="isAdmin"
            :type="currentView === 'admin' ? 'primary' : 'default'"
            @click="switchTab('admin')"
          >
            <el-icon><Setting /></el-icon>
            管理
          </el-button>
        </el-button-group>
      </div>
    </div>

    <!-- 学习统计卡片 -->
    <div class="stats-cards" v-if="userStats">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon total">
            <el-icon><VideoPlay /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ userStats.total_songs || 0 }}</div>
            <div class="stat-label">总歌曲</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon completed">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ userStats.completed_songs || 0 }}</div>
            <div class="stat-label">已完成</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon study-time">
            <el-icon><Timer /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ userStats.total_study_minutes || 0 }}</div>
            <div class="stat-label">学习分钟</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon level">
            <el-icon><Trophy /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ userStats.level || 1 }}</div>
            <div class="stat-label">等级</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 视图切换内容 -->
    <div class="content-area">
      <!-- 分类视图 -->
      <CategoriesView
        v-if="currentView === 'categories'"
        :categories="categories"
        :loading="loading"
        @category-selected="onCategorySelected"
        @refresh="fetchCategories"
      />

      <!-- 歌曲视图 -->
      <SongsView
        v-else-if="currentView === 'songs'"
        :songs="songs"
        :loading="loading"
        :selected-category="selectedCategory"
        @refresh="fetchSongs"
        @like-song="likeSong"
        @update-progress="updateSongProgress"
        @play-song="playSong"
        @search="handleSongSearch"
        @sort-change="handleSongSort"
        @filter-change="handleSongFilter"
      />

      <!-- 学习进度视图 -->
      <ProgressView
        v-else-if="currentView === 'progress'"
        :progress="userProgress"
        :loading="loading"
        :stats="userStats"
        @refresh="fetchUserProgress"
        @continue-learning="handleContinueLearning"
      />

      <!-- 播放器视图 -->
      <SongPlayerView
        v-else-if="currentView === 'player'"
        :song="currentSong"
        @back="backToSongs"
      />

      <!-- 管理员视图 -->
      <SongManagement
        v-else-if="currentView === 'admin'"
      />
    </div>

    <!-- 推荐歌曲浮动面板 -->
    <div v-if="recommendations.length > 0" class="recommendations-wrapper">
      <!-- 折叠后的箭头按钮 -->
      <el-button
        v-if="!showRecommendations"
        class="expand-btn"
        type="primary"
        :icon="ArrowLeft"
        @click="showRecommendations = true"
        circle
        size="large"
      />

      <!-- 推荐面板 -->
      <el-card v-if="showRecommendations" class="recommendations-panel" shadow="always">
        <template #header>
          <div class="panel-header">
            <div class="header-left">
              <el-icon><Star /></el-icon>
              <span>为您推荐</span>
            </div>
            <el-button
              type="text"
              size="small"
              :icon="ArrowRight"
              @click="showRecommendations = false"
              class="collapse-btn"
            />
          </div>
        </template>
        <div class="recommendations-list">
          <div
            v-for="song in recommendations"
            :key="song.id"
            class="recommendation-item"
            @click="playSong(song)"
          >
            <img :src="song.cover_image || defaultCover" :alt="song.title" class="song-cover" />
            <div class="song-info">
              <div class="song-title">{{ song.title }}</div>
              <div class="song-category">{{ song.category?.name }}</div>
            </div>
            <el-icon class="play-icon"><VideoPlay /></el-icon>
          </div>
        </div>
      </el-card>
    </div>

  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import {
  ArrowLeft, ArrowRight,
  Check,
  List,
  Setting,
  Star,
  Timer,
  TrendCharts,
  Trophy,
  VideoPlay
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import CategoriesView from '../components/CategoriesView.vue'
import ProgressView from '../components/ProgressView.vue'
import SongPlayerView from '../components/SongPlayerView.vue'
import SongsView from '../components/SongsView.vue'
import SongManagement from '../components/admin/SongManagement.vue'
import { useEnglishLearningStore } from '../stores/englishLearningStore'

// 状态管理
const store = useEnglishLearningStore()
const authStore = useAuthStore()

// 路由
const route = useRoute()
const router = useRouter()

// 响应式数据
const currentView = ref(route.query.tab || 'songs')
const loading = ref(false)
// 播放器相关状态
const currentSong = ref(null)
const selectedCategory = ref(null)
const showRecommendations = ref(false)

// 从store获取数据
const categories = computed(() => store.categories)
const songs = computed(() => store.songs)
const userProgress = computed(() => store.userProgress)
const userStats = computed(() => store.userStats)
const recommendations = computed(() => store.recommendations)

const defaultCover = 'https://via.placeholder.com/150x150/667eea/ffffff?text=♪'

// 检查是否为管理员
const isAdmin = computed(() => {
  return authStore.user?.role === 'admin'
})

// 生命周期
onMounted(async () => {
  const promises = [
    fetchCategories(),
    fetchSongs(),
    fetchUserStats(),
    fetchRecommendations()
  ]

  // 如果当前视图是进度页面，也获取用户进度
  if (currentView.value === 'progress') {
    promises.push(fetchUserProgress())
  }

  await Promise.all(promises)

  // 从URL恢复分类选择状态
  const categoryId = route.query.category
  if (categoryId) {
    // 等待分类数据加载完成后再设置选择的分类
    const category = categories.value.find(c => c.id == categoryId)
    if (category) {
      selectedCategory.value = category
    }
  }
})

// 标签切换方法
const switchTab = (tab) => {
  currentView.value = tab
  router.push({
    path: route.path,
    query: { ...route.query, tab }
  })
}

// 监听视图切换
watch(currentView, async (newView) => {
  if (newView === 'progress' && userProgress.value.length === 0) {
    await fetchUserProgress()
  }
})

// 方法
const fetchCategories = async () => {
  loading.value = true
  try {
    await store.fetchCategories({ is_active: true, limit: 20 })
  } catch (error) {
    ElMessage.error('获取分类失败')
  } finally {
    loading.value = false
  }
}

const fetchSongs = async (params = {}) => {
  loading.value = true
  try {
    const queryParams = {
      limit: 20,
      sort: 'sort:asc',
      ...params
    }

    // 只有非管理员才只获取已发布的歌曲
    if (!isAdmin.value) {
      queryParams.is_published = true
    }

    if (selectedCategory.value) {
      queryParams.category_id = selectedCategory.value.id
    }
    await store.fetchSongs(queryParams)
  } catch (error) {
    ElMessage.error('获取歌曲失败')
  } finally {
    loading.value = false
  }
}

const fetchUserProgress = async () => {
  loading.value = true
  try {
    await store.fetchUserProgress()
  } catch (error) {
    ElMessage.error('获取学习进度失败')
  } finally {
    loading.value = false
  }
}

const fetchUserStats = async () => {
  try {
    await store.fetchUserStats()
  } catch (error) {
    console.warn('获取用户统计失败:', error)
  }
}

const fetchRecommendations = async () => {
  try {
    await store.fetchRecommendations({ limit: 5 })
  } catch (error) {
    console.warn('获取推荐失败:', error)
  }
}

const onCategorySelected = (category) => {
  selectedCategory.value = category
  switchTab('songs')
  // 也可以在URL中保存分类信息
  router.push({
    path: route.path,
    query: { ...route.query, tab: 'songs', category: category?.id }
  })
  fetchSongs()
}

// 歌曲播放相关方法已移除，改为独立页面跳转

const likeSong = async (song) => {
  try {
    if (song.is_liked) {
      await store.unlikeSong(song.id)
      ElMessage.success('已取消收藏')
    } else {
      await store.likeSong(song.id)
      ElMessage.success('收藏成功!')
    }
  } catch (error) {
    ElMessage.error('操作失败，请重试')
  }
}

const updateSongProgress = async (songId, progressData) => {
  try {
    await store.updateProgress(songId, progressData)
    ElMessage.success('进度更新成功!')
  } catch (error) {
    ElMessage.error('进度更新失败')
  }
}

// 播放器事件处理
const playSong = (song) => {
  // 直接使用歌曲数据，不需要重新获取
  currentSong.value = song
  currentView.value = 'player'
}

const backToSongs = () => {
  currentView.value = 'songs'
  currentSong.value = null
}

const handleContinueLearning = (song) => {
  // 继续学习某首歌，跳转到播放器视图
  playSong(song)
}

// 歌曲筛选和搜索事件处理
const handleSongSearch = async (searchQuery) => {
  try {
    loading.value = true
    await fetchSongs({ search: searchQuery, page: 1 })
  } catch (error) {
    ElMessage.error('搜索失败')
  }
}

const handleSongSort = async (sortOptions) => {
  try {
    loading.value = true
    const sortParam = `${sortOptions.sort_by}:${sortOptions.sort_order}`
    await fetchSongs({ sort: sortParam, page: 1 })
  } catch (error) {
    ElMessage.error('排序失败')
  }
}

const handleSongFilter = async (filters) => {
  try {
    loading.value = true
    await fetchSongs({ ...filters, page: 1 })
  } catch (error) {
    ElMessage.error('筛选失败')
  }
}
</script>

<style scoped>
.english-learning-container {
  padding: 20px;
  min-height: calc(100vh - 60px);
  background: #f8fafc;
}

.english-learning-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  background: white;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.header-left .page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-left .page-subtitle {
  font-size: 16px;
  color: #6b7280;
  margin: 0;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  transition: transform 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon.completed {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
  color: white;
}

.stat-icon.study-time {
  background: linear-gradient(135deg, #ed8936 0%, #dd6b20 100%);
  color: white;
}

.stat-icon.level {
  background: linear-gradient(135deg, #f6ad55 0%, #ed8936 100%);
  color: white;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  margin-top: 4px;
}

.content-area {
  background: white;
  border-radius: 12px;
  min-height: 500px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.recommendations-wrapper {
  position: fixed;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  z-index: 100;
}

.expand-btn {
  position: fixed;
  right: -20px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 101;
  animation: pulse 2s infinite;
}

.expand-btn:hover {
  right: -15px;
  animation: none;
}

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(64, 158, 255, 0.7); }
  70% { box-shadow: 0 0 0 10px rgba(64, 158, 255, 0); }
  100% { box-shadow: 0 0 0 0 rgba(64, 158, 255, 0); }
}

.recommendations-panel {
  width: 300px;
  max-height: 400px;
  overflow-y: auto;
  transform: translateX(0);
  transition: transform 0.3s ease;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: #667eea;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.collapse-btn {
  padding: 4px;
  min-height: auto;
  color: #9ca3af;
}

.collapse-btn:hover {
  color: #667eea;
  background-color: rgba(102, 126, 234, 0.1);
}

.recommendations-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recommendation-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  background: #f8fafc;
  cursor: pointer;
  transition: all 0.2s ease;
}

.recommendation-item:hover {
  background: #e2e8f0;
  transform: scale(1.02);
}

.song-cover {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  object-fit: cover;
}

.song-info {
  flex: 1;
  min-width: 0;
}

.song-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.song-category {
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

.play-icon {
  color: #667eea;
  font-size: 18px;
}

@media (max-width: 1200px) {
  .recommendations-wrapper {
    position: relative;
    right: auto;
    top: auto;
    transform: none;
    width: 100%;
    margin-top: 20px;
  }

  .expand-btn {
    position: relative;
    right: auto;
    top: auto;
    transform: none;
    display: block;
    margin: 0 auto 10px auto;
  }

  .recommendations-panel {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .english-learning-container {
    padding: 16px;
  }

  .english-learning-header {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .stat-number {
    font-size: 24px;
  }
}
</style>