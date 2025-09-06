<template>
  <div class="categories-view">
    <div class="view-header">
      <h3 class="view-title">学习分类</h3>
      <p class="view-subtitle">选择您感兴趣的学习分类</p>
    </div>

    <div class="categories-grid" v-loading="loading">
      <el-card 
        v-for="category in categories" 
        :key="category.id"
        class="category-card"
        shadow="hover"
        @click="selectCategory(category)"
      >
        <div class="category-content">
          <!-- 分类图标 -->
          <div class="category-icon" :style="{ backgroundColor: category.color }">
            <el-icon v-if="category.icon" :size="32">
              <component :is="getIconComponent(category.icon)" />
            </el-icon>
            <el-icon v-else :size="32">
              <Folder />
            </el-icon>
          </div>

          <!-- 分类信息 -->
          <div class="category-info">
            <h4 class="category-name">{{ category.name }}</h4>
            <p class="category-name-cn">{{ category.name_cn }}</p>
            <p class="category-description">{{ category.description }}</p>
          </div>

          <!-- 分类统计 -->
          <div class="category-stats">
            <div class="stat-item">
              <span class="stat-number">{{ category.song_count || 0 }}</span>
              <span class="stat-label">首歌曲</span>
            </div>
          </div>

          <!-- 活跃状态 -->
          <div v-if="category.is_active" class="active-badge">
            <el-icon><Check /></el-icon>
          </div>
        </div>

        <!-- 悬停效果 -->
        <div class="category-overlay">
          <el-button type="primary" size="large" round>
            <el-icon><VideoPlay /></el-icon>
            开始学习
          </el-button>
        </div>
      </el-card>

      <!-- 空状态 -->
      <div v-if="!loading && categories.length === 0" class="empty-state">
        <el-empty description="暂无学习分类">
          <el-button type="primary" @click="$emit('refresh')">
            刷新分类
          </el-button>
        </el-empty>
      </div>
    </div>

    <!-- 分类筛选和排序 -->
    <div class="categories-toolbar">
      <div class="toolbar-left">
        <el-radio-group v-model="filterActive" @change="handleFilterChange">
          <el-radio-button :value="null">全部</el-radio-button>
          <el-radio-button :value="true">活跃</el-radio-button>
          <el-radio-button :value="false">非活跃</el-radio-button>
        </el-radio-group>
      </div>

      <div class="toolbar-right">
        <el-select v-model="sortBy" @change="handleSortChange" placeholder="排序方式">
          <el-option label="默认排序" value="sort" />
          <el-option label="名称排序" value="name" />
          <el-option label="创建时间" value="created_at" />
          <el-option label="歌曲数量" value="song_count" />
        </el-select>

        <el-button 
          :icon="sortOrder === 'asc' ? 'ArrowUp' : 'ArrowDown'"
          @click="toggleSortOrder"
          circle
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { 
  Folder, Check, VideoPlay, ArrowUp, ArrowDown,
  Headset, Reading, Star, Trophy
} from '@element-plus/icons-vue'

const props = defineProps({
  categories: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['category-selected', 'refresh', 'filter-change'])

// 响应式数据
const filterActive = ref(null)
const sortBy = ref('sort')
const sortOrder = ref('asc')

// 图标映射
const iconComponents = {
  'Music': Headset,
  'Headset': Headset,
  'Reading': Reading,
  'Star': Star,
  'Trophy': Trophy,
  'Folder': Folder
}

// 计算属性
const filteredCategories = computed(() => {
  let filtered = props.categories

  if (filterActive.value !== null) {
    filtered = filtered.filter(category => category.is_active === filterActive.value)
  }

  // 排序
  filtered = [...filtered].sort((a, b) => {
    let aVal = a[sortBy.value]
    let bVal = b[sortBy.value]

    if (typeof aVal === 'string') {
      aVal = aVal.toLowerCase()
      bVal = bVal.toLowerCase()
    }

    if (sortOrder.value === 'asc') {
      return aVal > bVal ? 1 : -1
    } else {
      return aVal < bVal ? 1 : -1
    }
  })

  return filtered
})

// 方法
const selectCategory = (category) => {
  emit('category-selected', category)
}

const getIconComponent = (iconName) => {
  return iconComponents[iconName] || Folder
}

const handleFilterChange = () => {
  emit('filter-change', { 
    is_active: filterActive.value,
    sort_by: sortBy.value,
    sort_order: sortOrder.value
  })
}

const handleSortChange = () => {
  emit('filter-change', { 
    is_active: filterActive.value,
    sort_by: sortBy.value,
    sort_order: sortOrder.value
  })
}

const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  handleSortChange()
}
</script>

<style scoped>
.categories-view {
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

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.category-card {
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
}

.category-card:hover .category-overlay {
  opacity: 1;
}

.category-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
  z-index: 1;
}

.category-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin: 0 auto;
}

.category-info {
  text-align: center;
}

.category-name {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
}

.category-name-cn {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 8px 0;
}

.category-description {
  font-size: 14px;
  color: #9ca3af;
  margin: 0;
  line-height: 1.5;
}

.category-stats {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.stat-item {
  text-align: center;
}

.stat-number {
  display: block;
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

.active-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 24px;
  height: 24px;
  background: #10b981;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
}

.category-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 2;
}

.categories-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-top: 1px solid #e5e7eb;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.empty-state {
  grid-column: 1 / -1;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

@media (max-width: 768px) {
  .categories-view {
    padding: 16px;
  }

  .categories-grid {
    grid-template-columns: 1fr;
  }

  .categories-toolbar {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .toolbar-left,
  .toolbar-right {
    justify-content: center;
  }
}
</style>