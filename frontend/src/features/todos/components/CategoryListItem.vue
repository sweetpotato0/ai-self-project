<template>
  <div 
    class="category-list-item" 
    :class="{ 'highlighted': isHighlighted }"
  >
    <!-- 名称列 -->
    <div class="col-name">
      <div class="name-content">
        <div class="category-icon" :style="{ backgroundColor: category.color }">
          <el-icon size="14">
            <component :is="getCategoryIcon(category.icon)" />
          </el-icon>
        </div>
        <div class="name-info">
          <div class="category-name">
            <template v-if="highlightedName">
              <span v-html="highlightedName"></span>
            </template>
            <template v-else>
              {{ category.name }}
            </template>
          </div>
          <div v-if="category.description" class="category-description">
            <template v-if="highlightedDescription">
              <span v-html="highlightedDescription"></span>
            </template>
            <template v-else>
              {{ category.description }}
            </template>
          </div>
        </div>
      </div>
    </div>

    <!-- 父级分类列 -->
    <div class="col-parent">
      <template v-if="parentCategory">
        <div class="parent-info">
          <div class="parent-icon" :style="{ backgroundColor: parentCategory.color }"></div>
          <span class="parent-name">{{ parentCategory.name }}</span>
        </div>
      </template>
      <span v-else class="no-parent">顶级分类</span>
    </div>

    <!-- 任务数量列 -->
    <div class="col-count">
      <el-badge v-if="todoCount > 0" :value="todoCount" type="primary">
        <div class="count-display">{{ todoCount }}</div>
      </el-badge>
      <span v-else class="no-tasks">0</span>
    </div>

    <!-- 更新时间列 -->
    <div class="col-updated">
      <el-tooltip :content="formatFullTime(category.updated_at)" placement="top">
        <span class="time-text">{{ formatRelativeTime(category.updated_at) }}</span>
      </el-tooltip>
    </div>

    <!-- 操作列 -->
    <div class="col-actions">
      <el-button-group size="small">
        <el-button 
          text 
          type="primary"
          @click="handleAddChild"
          title="添加子分类"
        >
          <el-icon><Plus /></el-icon>
        </el-button>
        <el-button 
          text 
          type="primary"
          @click="handleEdit"
          title="编辑分类"
        >
          <el-icon><Edit /></el-icon>
        </el-button>
        <el-button 
          text 
          type="danger"
          @click="handleDelete"
          title="删除分类"
        >
          <el-icon><Delete /></el-icon>
        </el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { 
  Plus, Edit, Delete,
  // 确认可用的图标
  Folder, Star, User, Setting, Calendar, Bell, Sunny, Moon
} from '@element-plus/icons-vue'

const props = defineProps({
  category: {
    type: Object,
    required: true
  },
  categories: {
    type: Array,
    default: () => []
  },
  searchKeyword: {
    type: String,
    default: ''
  },
  getTodoCount: {
    type: Function,
    required: true
  }
})

const emit = defineEmits(['edit', 'delete', 'add-child'])

// 计算属性
const todoCount = computed(() => props.getTodoCount(props.category.id))

const parentCategory = computed(() => {
  if (!props.category.parent_id) return null
  return props.categories.find(cat => cat.id === props.category.parent_id)
})

// 搜索高亮
const isHighlighted = computed(() => {
  if (!props.searchKeyword) return false
  const keyword = props.searchKeyword.toLowerCase()
  return props.category.name.toLowerCase().includes(keyword) ||
         (props.category.description && props.category.description.toLowerCase().includes(keyword))
})

const highlightedName = computed(() => {
  if (!props.searchKeyword || !isHighlighted.value) return null
  const keyword = props.searchKeyword
  const regex = new RegExp(`(${keyword})`, 'gi')
  return props.category.name.replace(regex, '<mark>$1</mark>')
})

const highlightedDescription = computed(() => {
  if (!props.searchKeyword || !isHighlighted.value || !props.category.description) return null
  const keyword = props.searchKeyword
  const regex = new RegExp(`(${keyword})`, 'gi')
  return props.category.description.replace(regex, '<mark>$1</mark>')
})

// 图标映射
const iconMap = {
  folder: Folder,
  star: Star,
  user: User,
  setting: Setting,
  calendar: Calendar,
  bell: Bell,
  sunny: Sunny,
  moon: Moon,
  edit: Edit,
  plus: Plus,
  delete: Delete
}

const getCategoryIcon = (iconName) => {
  return iconMap[iconName] || Folder
}

// 时间格式化
const formatRelativeTime = (dateString) => {
  if (!dateString) return '-'
  
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date
  
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return date.toLocaleDateString('zh-CN')
}

const formatFullTime = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString('zh-CN')
}

// 事件处理
const handleEdit = () => {
  emit('edit', props.category)
}

const handleDelete = () => {
  emit('delete', props.category)
}

const handleAddChild = () => {
  emit('add-child', props.category)
}
</script>

<style scoped>
.category-list-item {
  display: grid;
  grid-template-columns: 2fr 1.5fr 100px 120px 120px;
  gap: 16px;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  transition: all 0.2s ease;
}

.category-list-item:hover {
  background-color: #f8fafc;
}

.category-list-item.highlighted {
  background-color: #fff7e6;
  border-color: #ffd591;
}

.category-list-item:last-child {
  border-bottom: none;
}

/* 名称列 */
.col-name {
  min-width: 0;
}

.name-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.category-icon {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.name-info {
  flex: 1;
  min-width: 0;
}

.category-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.4;
}

.category-name :deep(mark) {
  background-color: #fff2cc;
  color: #d48806;
  padding: 0 2px;
  border-radius: 2px;
}

.category-description {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-top: 2px;
  line-height: 1.3;
}

.category-description :deep(mark) {
  background-color: #fff2cc;
  color: #d48806;
  padding: 0 2px;
  border-radius: 2px;
}

/* 父级分类列 */
.col-parent {
  min-width: 0;
}

.parent-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.parent-icon {
  width: 12px;
  height: 12px;
  border-radius: 2px;
  flex-shrink: 0;
}

.parent-name {
  font-size: 13px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-parent {
  font-size: 13px;
  color: #c0c4cc;
  font-style: italic;
}

/* 任务数量列 */
.col-count {
  text-align: center;
}

.count-display {
  width: 32px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 500;
}

.no-tasks {
  color: #c0c4cc;
  font-size: 13px;
}

/* 更新时间列 */
.col-updated {
  text-align: center;
}

.time-text {
  font-size: 12px;
  color: #909399;
  cursor: help;
}

/* 操作列 */
.col-actions {
  display: flex;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.category-list-item:hover .col-actions {
  opacity: 1;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .category-list-item {
    grid-template-columns: 1fr 80px 100px;
  }
  
  .col-parent,
  .col-updated {
    display: none;
  }
  
  .name-content {
    gap: 8px;
  }
  
  .category-icon {
    width: 20px;
    height: 20px;
  }
}
</style>