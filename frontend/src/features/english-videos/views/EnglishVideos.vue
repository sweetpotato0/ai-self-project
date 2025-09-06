<template>
  <div class="english-videos-container">
    <div class="english-videos-header">
      <div class="header-left">
        <!-- <h2 class="page-title">英语视频</h2> -->
        <p class="page-subtitle">通过经典视频系列学习英语</p>
      </div>
      <div class="header-right">
        <el-button-group>
          <el-button
            :type="currentView === 'series' ? 'primary' : 'default'"
            @click="switchView('series')"
          >
            <el-icon><VideoPlay /></el-icon>
            系列
          </el-button>
          <el-button
            :type="currentView === 'progress' ? 'primary' : 'default'"
            @click="switchView('progress')"
          >
            <el-icon><TrendCharts /></el-icon>
            进度
          </el-button>
          <el-button
            v-if="isAdmin"
            :type="currentView === 'admin' ? 'primary' : 'default'"
            @click="switchView('admin')"
          >
            <el-icon><Setting /></el-icon>
            管理
          </el-button>
        </el-button-group>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards" v-if="videoStats">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon total">
            <el-icon><VideoPlay /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ videoStats.total_series || 0 }}</div>
            <div class="stat-label">总系列</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon episodes">
            <el-icon><List /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ videoStats.total_episodes || 0 }}</div>
            <div class="stat-label">总集数</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon views">
            <el-icon><View /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ formatNumber(videoStats.total_views || 0) }}</div>
            <div class="stat-label">总观看</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon completed">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ videoStats.watched_series || 0 }}</div>
            <div class="stat-label">已看完</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 视图切换内容 -->
    <div class="content-area">
      <!-- 系列视图 -->
      <div v-if="currentView === 'series'">
        <!-- 搜索和筛选 -->
        <div class="search-filters">
          <el-row :gutter="20" align="middle">
            <el-col :lg="8" :md="12" :sm="24">
              <el-input
                v-model="searchQuery"
                placeholder="搜索视频系列..."
                :prefix-icon="Search"
                size="large"
                clearable
                @input="handleSearch"
              />
            </el-col>
            <el-col :lg="4" :md="6" :sm="12">
              <el-select
                v-model="selectedDifficulty"
                placeholder="难度筛选"
                size="large"
                clearable
                @change="handleFilter"
              >
                <el-option label="全部难度" value="" />
                <el-option label="入门" :value="1" />
                <el-option label="初级" :value="2" />
                <el-option label="中级" :value="3" />
                <el-option label="高级" :value="4" />
                <el-option label="专家" :value="5" />
              </el-select>
            </el-col>
            <el-col :lg="4" :md="6" :sm="12">
              <el-select
                v-model="selectedAge"
                placeholder="年龄段"
                size="large"
                clearable
                @change="handleFilter"
              >
                <el-option label="全部年龄" value="" />
                <el-option label="3-5岁" value="3-5" />
                <el-option label="6-8岁" value="6-8" />
                <el-option label="9-12岁" value="9-12" />
                <el-option label="青少年" value="13-17" />
                <el-option label="成人" value="18+" />
              </el-select>
            </el-col>
            <el-col :lg="8" :md="24" :sm="24" class="actions">
              <div class="search-actions">
                <el-button-group>
                  <el-button
                    :type="viewMode === 'grid' ? 'primary' : ''"
                    :icon="Grid"
                    @click="viewMode = 'grid'"
                  >
                    网格
                  </el-button>
                  <el-button
                    :type="viewMode === 'list' ? 'primary' : ''"
                    :icon="List"
                    @click="viewMode = 'list'"
                  >
                    列表
                  </el-button>
                </el-button-group>
                <el-button @click="fetchVideoSeries" :loading="loading">
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
              </div>
            </el-col>
          </el-row>
        </div>

        <!-- 视频系列列表 -->
        <div class="video-series-container">
          <div v-if="loading" class="loading-state">
            <el-skeleton :rows="4" animated />
          </div>

          <div v-else-if="filteredSeries.length === 0" class="empty-state">
            <el-empty description="暂无视频系列" />
          </div>

          <div v-else>
            <!-- 网格视图 -->
            <div v-if="viewMode === 'grid'" class="grid-view">
              <el-row :gutter="24">
                <el-col
                  v-for="series in paginatedSeries"
                  :key="series.id"
                  :lg="6" :md="8" :sm="12" :xs="24"
                  class="series-col"
                >
                  <div class="series-card" @click="viewSeries(series)">
                  <div class="series-cover">
                    <img
                      :src="series.cover_image || '/default-video-cover.png'"
                      :alt="series.title"
                      @error="handleImageError"
                    />
                    <div class="series-overlay">
                      <div class="play-btn">
                        <el-icon size="40"><VideoPlay /></el-icon>
                      </div>
                      <div class="episodes-count">
                        {{ series.episode_count || 0 }} 集
                      </div>
                    </div>
                  </div>

                  <div class="series-info">
                    <h3 class="series-title">{{ series.title }}</h3>
                    <p class="series-description">{{ series.description }}</p>

                    <div class="series-meta">
                      <div class="difficulty">
                        <el-icon><Star /></el-icon>
                        <span v-for="i in (series.difficulty || 1)" :key="i" class="star filled">★</span>
                        <span v-for="i in (5 - (series.difficulty || 1))" :key="i + series.difficulty" class="star">★</span>
                      </div>
                      <div class="age-tag">
                        <el-tag size="small" type="info">{{ series.age_range }}</el-tag>
                      </div>
                    </div>

                    <div class="series-stats">
                      <span class="stat-item">
                        <el-icon><View /></el-icon>
                        {{ formatNumber(series.view_count || 0) }}
                      </span>
                      <span class="stat-item">
                        <el-icon><Star /></el-icon>
                        {{ formatNumber(series.like_count || 0) }}
                      </span>
                    </div>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>

          <!-- 列表视图 -->
          <div v-else class="list-view">
            <div
              v-for="series in paginatedSeries"
              :key="series.id"
              class="series-list-item"
              @click="viewSeries(series)"
            >
              <div class="series-thumbnail">
                <img
                  :src="series.cover_image || '/default-video-cover.png'"
                  :alt="series.title"
                  @error="handleImageError"
                />
                <div class="play-overlay">
                  <el-icon><VideoPlay /></el-icon>
                </div>
              </div>

              <div class="series-details">
                <h3 class="series-title">{{ series.title }}</h3>
                <p class="series-description">{{ series.description }}</p>

                <div class="series-info-row">
                  <div class="difficulty">
                    <span v-for="i in (series.difficulty || 1)" :key="i" class="star filled">★</span>
                    <span v-for="i in (5 - (series.difficulty || 1))" :key="i + series.difficulty" class="star">★</span>
                  </div>
                  <div class="age-tag">
                    <el-tag size="small" type="info">{{ series.age_range }}</el-tag>
                  </div>
                  <div class="episodes-count">{{ series.episode_count || 0 }} 集</div>
                </div>

                <div class="series-stats">
                  <span class="stat">{{ formatNumber(series.view_count || 0) }} 次观看</span>
                  <span class="stat">{{ formatNumber(series.like_count || 0) }} 个收藏</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[12, 24, 36, 48]"
              :small="false"
              :disabled="loading"
              :background="true"
              layout="total, sizes, prev, pager, next, jumper"
              :total="filteredSeries.length"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
          </div>
        </div>
      </div>

      <!-- 学习进度视图 -->
      <div v-else-if="currentView === 'progress'" class="progress-view">
        <div class="progress-content">
          <h3>学习进度</h3>
          <p>正在观看的视频系列</p>

          <div v-if="userProgress.length === 0" class="empty-state">
            <el-empty description="暂无学习进度" />
          </div>

          <div v-else class="progress-list">
            <div
              v-for="progress in userProgress"
              :key="progress.series_id"
              class="progress-item"
              @click="continueWatching(progress)"
            >
              <div class="progress-cover">
                <img
                  :src="progress.series?.cover_image || '/default-video-cover.png'"
                  :alt="progress.series?.title"
                />
                <div class="progress-overlay">
                  <div class="progress-percentage">{{ Math.round(progress.watch_progress * 100) }}%</div>
                </div>
              </div>
              <div class="progress-info">
                <h4>{{ progress.series?.title }}</h4>
                <p>第 {{ progress.current_episode }} / {{ progress.series?.episode_count }} 集</p>
                <div class="progress-bar">
                  <el-progress :percentage="Math.round(progress.watch_progress * 100)" :show-text="false" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 管理员视图 -->
      <div v-else-if="currentView === 'admin'" class="admin-view">
        <VideoSeriesManagement />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useEnglishVideosStore } from '@/features/english-videos/stores/englishVideosStore'
import {
  Check,
  Grid, List,
  Refresh,
  Search,
  Setting,
  Star,
  TrendCharts,
  VideoPlay,
  View
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'
import VideoSeriesManagement from '@/features/english-videos/components/admin/VideoSeriesManagement.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const englishVideosStore = useEnglishVideosStore()

// 响应式数据
const searchQuery = ref('')
const selectedDifficulty = ref('')
const selectedAge = ref('')
const viewMode = ref('grid')
const currentPage = ref(1)
const pageSize = ref(12)
const searchTimeout = ref(null)
const currentView = ref(route.query.tab || 'series')
const videoStats = ref(null)
const userProgress = ref([])

// 使用Store中的响应式数据
const { videoSeries, loading } = storeToRefs(englishVideosStore)

// 计算属性
const isAuthenticated = computed(() => authStore.isAuthenticated)

const isAdmin = computed(() => {
  console.log(authStore.user)
  return authStore.user?.role === 'admin'
})

const filteredSeries = computed(() => {
  if (!Array.isArray(videoSeries.value)) {
    return []
  }
  let result = [...videoSeries.value]

  // 搜索过滤
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(series =>
      series.title.toLowerCase().includes(query) ||
      series.description.toLowerCase().includes(query)
    )
  }

  // 难度过滤
  if (selectedDifficulty.value) {
    result = result.filter(series => series.difficulty === selectedDifficulty.value)
  }

  // 年龄过滤
  if (selectedAge.value) {
    result = result.filter(series => series.age_range === selectedAge.value)
  }

  return result
})

const paginatedSeries = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredSeries.value.slice(start, start + pageSize.value)
})

// 方法
const fetchVideoSeries = async () => {
  try {
    const params = {
      search: searchQuery.value,
      difficulty: selectedDifficulty.value,
      age_range: selectedAge.value,
      sort_order: 'asc'
    }
    
    // 过滤空值
    const filteredParams = Object.fromEntries(
      Object.entries(params).filter(([_, value]) => value !== '' && value !== null && value !== undefined)
    )

    await englishVideosStore.fetchVideoSeries(filteredParams)
  } catch (error) {
    console.error('Failed to fetch video series:', error)
    ElMessage.error('获取视频系列失败')
  }
}

const handleSearch = () => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }

  searchTimeout.value = setTimeout(() => {
    currentPage.value = 1
  }, 300)
}

const handleFilter = () => {
  currentPage.value = 1
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

const viewSeries = (series) => {
  router.push({
    name: 'VideoSeries',
    params: {
      seriesId: series.id
    }
  })
}

const handleImageError = (e) => {
  e.target.src = '/default-video-cover.png'
}

const formatNumber = (num) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

// 新增方法
const switchView = (view) => {
  currentView.value = view
  router.push({
    path: route.path,
    query: { ...route.query, tab: view }
  })
}

const fetchVideoStats = async () => {
  try {
    const data = await englishVideosStore.getVideoStats()
    videoStats.value = data
  } catch (error) {
    console.error('Failed to fetch video stats:', error)
  }
}

const fetchUserProgress = async () => {
  if (!isAuthenticated.value) return

  try {
    const data = await englishVideosStore.getUserProgress()
    userProgress.value = data || []
  } catch (error) {
    console.error('Failed to fetch user progress:', error)
  }
}

const continueWatching = (progress) => {
  router.push({
    name: 'VideoPlayer',
    params: {
      seriesId: progress.series_id,
      episodeId: progress.current_episode || 1
    }
  })
}


// 生命周期
onMounted(async () => {
  await Promise.all([
    fetchVideoSeries(),
    fetchVideoStats(),
    fetchUserProgress()
  ])
})
</script>

<style scoped>
.english-videos-container {
  padding: 20px;
  min-height: calc(100vh - 60px);
  background: #f8fafc;
}

.english-videos-header {
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

.stat-icon.episodes {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
  color: white;
}

.stat-icon.views {
  background: linear-gradient(135deg, #ed8936 0%, #dd6b20 100%);
  color: white;
}

.stat-icon.completed {
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
  padding: 24px;
}

.search-filters {
  margin-bottom: 24px;
  padding: 20px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.actions {
  text-align: right;
}

.search-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
}

/* 进度视图样式 */
.progress-view {
  padding: 24px;
}

.progress-content h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.progress-content p {
  color: #6b7280;
  margin: 0 0 24px 0;
}

.progress-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.progress-item {
  display: flex;
  gap: 16px;
  background: #f8fafc;
  padding: 16px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid #e2e8f0;
}

.progress-item:hover {
  background: #f1f5f9;
  border-color: #667eea;
  transform: translateY(-2px);
}

.progress-cover {
  position: relative;
  width: 80px;
  height: 60px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
}

.progress-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.progress-overlay {
  position: absolute;
  top: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 4px 8px;
  font-size: 12px;
  border-radius: 0 8px 0 8px;
}

.progress-info {
  flex: 1;
  min-width: 0;
}

.progress-info h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.progress-info p {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 12px 0;
}

.progress-bar {
  margin-top: 8px;
}

/* 管理员视图样式 */
.admin-view {
  padding: 24px;
}

.admin-content h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.admin-content p {
  color: #6b7280;
  margin: 0 0 24px 0;
}

.admin-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.admin-table {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.video-series-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  padding: 32px;
}

.loading-state, .empty-state {
  padding: 80px 20px;
  text-align: center;
}

/* 网格视图样式 */
.grid-view .series-col {
  margin-bottom: 32px;
}

.series-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  height: 100%;
}

.series-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
}

.series-cover {
  position: relative;
  width: 100%;
  padding-bottom: 56.25%; /* 16:9 aspect ratio */
  overflow: hidden;
}

.series-cover img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.series-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0) 0%, rgba(0,0,0,0.3) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.series-card:hover .series-overlay {
  opacity: 1;
}

.play-btn {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
}

.episodes-count {
  position: absolute;
  bottom: 12px;
  right: 12px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.series-info {
  padding: 20px;
}

.series-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.series-description {
  font-size: 14px;
  color: #7f8c8d;
  margin: 0 0 16px 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.series-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.difficulty {
  display: flex;
  align-items: center;
  gap: 4px;
}

.star {
  color: #e0e0e0;
  font-size: 14px;
}

.star.filled {
  color: #f39c12;
}

.series-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #ecf0f1;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #7f8c8d;
}

/* 列表视图样式 */
.series-list-item {
  display: flex;
  gap: 20px;
  padding: 20px;
  border-bottom: 1px solid #ecf0f1;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.series-list-item:hover {
  background: #f8f9fa;
}

.series-list-item:last-child {
  border-bottom: none;
}

.series-thumbnail {
  position: relative;
  width: 160px;
  height: 90px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
}

.series-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.play-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.series-list-item:hover .play-overlay {
  opacity: 1;
}

.series-details {
  flex: 1;
  min-width: 0;
}

.series-info-row {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 12px 0;
}

.episodes-count {
  font-size: 14px;
  color: #666;
}

.pagination-container {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    padding: 20px 0;
  }

  .main-content {
    padding: 0 16px 20px;
  }

  .search-filters {
    padding: 16px;
  }

  .video-series-container {
    padding: 20px;
  }

  .actions {
    text-align: left;
    margin-top: 16px;
  }

  .series-list-item {
    flex-direction: column;
    gap: 12px;
  }

  .series-thumbnail {
    width: 100%;
    height: 200px;
  }
}
</style>