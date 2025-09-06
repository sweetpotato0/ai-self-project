<template>
  <div class="songs-view">
    <div class="view-header">
      <div class="header-left">
        <h3 class="view-title">
          英语歌曲
          <span v-if="selectedCategory" class="category-badge">
            {{ selectedCategory.name }}
          </span>
        </h3>
        <p class="view-subtitle">通过优美的歌曲学习英语</p>
      </div>
    </div>

    <!-- 搜索和排序栏 -->
    <div class="search-sort-bar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索歌曲..."
        prefix-icon="Search"
        @input="handleSearch"
        class="search-input"
      />
      <el-select
        v-model="sortOption"
        @change="handleSort"
        placeholder="排序"
        class="sort-select"
      >
        <el-option label="默认排序" value="sort:asc" />
        <el-option label="最新发布" value="created_at:desc" />
        <el-option label="最热门" value="view_count:desc" />
        <el-option label="最受喜爱" value="like_count:desc" />
        <el-option label="难度从低到高" value="difficulty:asc" />
      </el-select>
    </div>

    <!-- 筛选栏 -->
    <div class="filters-bar">
      <div class="filter-group">
        <el-text class="filter-label">难度：</el-text>
        <el-radio-group v-model="difficultyFilter" @change="handleFilter">
          <el-radio-button :value="null">全部</el-radio-button>
          <el-radio-button :value="1">入门</el-radio-button>
          <el-radio-button :value="2">初级</el-radio-button>
          <el-radio-button :value="3">中级</el-radio-button>
          <el-radio-button :value="4">高级</el-radio-button>
          <el-radio-button :value="5">专家</el-radio-button>
        </el-radio-group>
      </div>

      <div class="filter-group">
        <el-text class="filter-label">年龄：</el-text>
        <el-select
          v-model="ageRangeFilter"
          @change="handleFilter"
          placeholder="年龄范围"
          style="width: 120px;"
        >
          <el-option label="全部年龄" :value="null" />
          <el-option label="3-6岁" value="3-6" />
          <el-option label="7-12岁" value="7-12" />
          <el-option label="13-18岁" value="13-18" />
          <el-option label="18岁以上" value="18+" />
        </el-select>
      </div>
    </div>

    <!-- 歌曲网格 -->
    <div class="songs-grid" v-loading="loading">
      <el-card
        v-for="song in songs"
        :key="song.id"
        class="song-card"
        shadow="hover"
      >
        <!-- 歌曲封面 -->
        <div class="song-cover" @click="playSong(song)">
          <img
            :src="song.cover_image || defaultCover"
            :alt="song.title"
            class="cover-image"
          />
          <div class="play-overlay">
            <el-button type="primary" :icon="VideoPlay" circle size="large" />
          </div>
        </div>

        <!-- 歌曲信息 -->
        <div class="song-info">
          <h4 class="song-title" @click="playSong(song)">{{ song.title }}</h4>
          <p class="song-title-cn" v-if="song.title_cn">{{ song.title_cn }}</p>

          <!-- 歌曲元数据 -->
          <div class="song-meta">
            <el-tag
              v-if="song.category"
              :color="song.category.color"
              size="small"
              effect="light"
            >
              {{ song.category.name }}
            </el-tag>

            <div class="difficulty-stars">
              <el-icon
                v-for="i in 5"
                :key="i"
                :class="['star', { active: i <= (song.difficulty || 1) }]"
              >
                <Star />
              </el-icon>
            </div>

            <el-tag v-if="song.age_range" size="small" type="info">
              {{ song.age_range }}
            </el-tag>
          </div>

          <!-- 歌曲描述 -->
          <p class="song-description" v-if="song.description">
            {{ song.description }}
          </p>

          <!-- 歌曲统计 -->
          <div class="song-stats">
            <div class="stat-item">
              <el-icon><DataLine /></el-icon>
              <span>{{ formatNumber(song.view_count || 0) }}</span>
            </div>
            <div class="stat-item">
              <el-icon><Star /></el-icon>
              <span>{{ formatNumber(song.like_count || 0) }}</span>
            </div>
            <!-- <div class="stat-item" v-if="song.duration_seconds">
              <el-icon><Timer /></el-icon>
              <span>{{ formatDuration(song.duration_seconds) }}</span>
            </div> -->
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="song-actions">
          <el-button
            :type="song.is_liked ? 'danger' : 'default'"
            :icon="Star"
            @click="toggleLike(song)"
            circle
            size="small"
          />

          <el-button
            type="primary"
            @click="playSong(song)"
            size="small"
          >
            <el-icon><VideoPlay /></el-icon>
            播放
          </el-button>

          <el-dropdown @command="handleSongAction">
            <el-button :icon="MoreFilled" circle size="small" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :command="{ action: 'progress', song }">
                  <el-icon><TrendCharts /></el-icon>
                  查看进度
                </el-dropdown-item>
                <el-dropdown-item :command="{ action: 'vocabulary', song }">
                  <el-icon><Reading /></el-icon>
                  学习词汇
                </el-dropdown-item>
                <el-dropdown-item :command="{ action: 'lyrics', song }">
                  <el-icon><Document /></el-icon>
                  查看歌词
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>

        <!-- 学习进度条 -->
        <div class="progress-bar" v-if="getSongProgress(song.id)">
          <el-progress
            :percentage="getSongProgress(song.id).progress || 0"
            :stroke-width="4"
            :show-text="false"
          />
          <div class="progress-text">
            进度: {{ getSongProgress(song.id).progress || 0 }}%
          </div>
        </div>
      </el-card>

      <!-- 空状态 -->
      <div v-if="!loading && songs.length === 0" class="empty-state">
        <el-empty description="暂无歌曲">
          <el-button type="primary" @click="$emit('refresh')">
            刷新歌曲
          </el-button>
        </el-empty>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination-container" v-if="totalPages > 1">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        :page-count="totalPages"
        layout="prev, pager, next, jumper, sizes, total"
        :page-sizes="[12, 24, 48, 96]"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import {
  DataLine,
  Document, MoreFilled,
  Reading,
  Star,
  Timer, TrendCharts,
  VideoPlay
} from '@element-plus/icons-vue'
import { computed, ref } from 'vue'

const props = defineProps({
  songs: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  selectedCategory: {
    type: Object,
    default: null
  },
  userProgress: {
    type: Array,
    default: () => []
  },
  total: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 20
  },
  totalPages: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits([
  'song-selected', 'refresh', 'like-song', 'update-progress',
  'filter-change', 'search', 'sort-change', 'page-change', 'size-change', 'play-song'
])

// 响应式数据
const searchQuery = ref('')
const sortOption = ref('sort:asc')
const difficultyFilter = ref(null)
const ageRangeFilter = ref(null)
const searchTimeout = ref(null)

const defaultCover = '/images/default-song-cover.jpg'

// 计算属性
const progressMap = computed(() => {
  const map = {}
  props.userProgress.forEach(progress => {
    map[progress.song_id] = progress
  })
  return map
})

// 方法
const playSong = (song) => {
  // 触发播放事件，让父组件处理
  emit('play-song', song)
}

const toggleLike = (song) => {
  emit('like-song', song)
}

const getSongProgress = (songId) => {
  return progressMap.value[songId]
}

const handleSearch = () => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }

  searchTimeout.value = setTimeout(() => {
    emit('search', searchQuery.value)
  }, 500)
}

const handleSort = () => {
  const [sortBy, sortOrder] = sortOption.value.split(':')
  emit('sort-change', { sort_by: sortBy, sort_order: sortOrder })
}

const handleFilter = () => {
  const filters = {}
  if (difficultyFilter.value !== null) {
    filters.difficulty = difficultyFilter.value
  }
  if (ageRangeFilter.value !== null) {
    filters.age_range = ageRangeFilter.value
  }
  emit('filter-change', filters)
}

const handleSongAction = ({ action, song }) => {
  switch (action) {
    case 'progress':
      // 显示进度详情
      break
    case 'vocabulary':
      // 显示词汇学习
      break
    case 'lyrics':
      // 显示歌词
      break
  }
}

const handlePageChange = (page) => {
  emit('page-change', page)
}

const handleSizeChange = (size) => {
  emit('size-change', size)
}

const formatNumber = (num) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

const formatDuration = (seconds) => {
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}
</script>

<style scoped>
.songs-view {
  padding: 24px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.view-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.category-badge {
  font-size: 14px;
  font-weight: 500;
  padding: 4px 12px;
  background: #667eea;
  color: white;
  border-radius: 20px;
}

.view-subtitle {
  color: #6b7280;
  margin: 0;
}

.filters-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  margin-bottom: 24px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-label {
  font-weight: 500;
  color: #374151;
  white-space: nowrap;
}

.search-sort-bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.search-input {
  width: 300px;
}

.sort-select {
  width: 160px;
}

.songs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.song-card {
  position: relative;
  transition: all 0.3s ease;
}

.song-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
}

.song-cover {
  position: relative;
  width: 100%;
  height: 200px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  margin-bottom: 16px;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.song-card:hover .cover-image {
  transform: scale(1.05);
}

.play-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.song-cover:hover .play-overlay {
  opacity: 1;
}

.song-info {
  margin-bottom: 16px;
}

.song-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
  cursor: pointer;
  transition: color 0.2s ease;
}

.song-title:hover {
  color: #667eea;
}

.song-title-cn {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 8px 0;
}

.song-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  /* margin-bottom: 8px; */
  flex-wrap: wrap;
}

.difficulty-stars {
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

.song-description {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.song-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #6b7280;
}

.song-actions {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: flex-end;
}

.progress-bar {
  margin-top: 12px;
}

.progress-text {
  font-size: 12px;
  color: #6b7280;
  text-align: center;
  margin-top: 4px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.empty-state {
  grid-column: 1 / -1;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

@media (max-width: 768px) {
  .songs-view {
    padding: 16px;
  }

  .view-header {
    flex-direction: column;
    gap: 16px;
  }

  .header-right {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .filters-bar {
    flex-direction: column;
    gap: 16px;
  }

  .songs-grid {
    grid-template-columns: 1fr;
  }

  .song-actions {
    flex-wrap: wrap;
    gap: 8px;
  }
}
</style>