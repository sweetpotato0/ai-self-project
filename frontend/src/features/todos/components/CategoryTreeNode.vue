<template>
  <div class="category-tree-node" :class="{ 'highlighted': isHighlighted }">
    <!-- 分类项 -->
    <div 
      class="category-item"
      :style="{ paddingLeft: level * 24 + 'px' }"
      @click="toggleExpanded"
    >
      <!-- 展开/收起图标 -->
      <div class="expand-icon" :class="{ 'expanded': expanded, 'has-children': hasChildren }">
        <el-icon v-if="hasChildren" size="14">
          <ArrowRight />
        </el-icon>
      </div>

      <!-- 分类图标 -->
      <div class="category-icon" :style="{ backgroundColor: category.color }">
        <el-icon size="14">
          <component :is="getCategoryIcon(category.icon)" />
        </el-icon>
      </div>

      <!-- 分类信息 -->
      <div class="category-info">
        <div class="category-main">
          <span class="category-name" :class="{ 'highlighted-text': isHighlighted }">
            <template v-if="highlightedName">
              <span v-html="highlightedName"></span>
            </template>
            <template v-else>
              {{ category.name }}
            </template>
          </span>
          
          <div class="category-badges">
            <el-badge v-if="todoCount > 0" :value="todoCount" class="todo-badge" />
            <span v-if="childrenCount > 0" class="children-count">
              {{ childrenCount }} 个子分类
            </span>
          </div>
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

      <!-- 操作按钮 -->
      <div class="category-actions" @click.stop>
        <el-dropdown trigger="hover" placement="bottom-end">
          <el-button text size="small">
            <el-icon><More /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleAddChild">
                <el-icon><Plus /></el-icon>
                添加子分类
              </el-dropdown-item>
              <el-dropdown-item @click="handleEdit">
                <el-icon><Edit /></el-icon>
                编辑分类
              </el-dropdown-item>
              <el-dropdown-item divided @click="handleDelete" class="delete-item">
                <el-icon><Delete /></el-icon>
                删除分类
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 子分类 -->
    <transition name="expand" appear>
      <div v-show="expanded && hasChildren" class="children-container">
        <category-tree-node
          v-for="child in children"
          :key="child.id"
          :category="child"
          :categories="categories"
          :level="level + 1"
          :search-keyword="searchKeyword"
          :get-todo-count="getTodoCount"
          @edit="$emit('edit', $event)"
          @delete="$emit('delete', $event)"
          @add-child="$emit('add-child', $event)"
        />
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { 
  ArrowRight, Plus, Edit, Delete, More,
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
  level: {
    type: Number,
    default: 0
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

// 响应式状态
const expanded = ref(true) // 默认展开

// 计算属性
const children = computed(() => {
  return props.categories
    .filter(cat => cat.parent_id === props.category.id)
    .sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))
})

const hasChildren = computed(() => children.value.length > 0)
const todoCount = computed(() => props.getTodoCount(props.category.id))
const childrenCount = computed(() => children.value.length)

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
  delete: Delete,
  more: More
}

const getCategoryIcon = (iconName) => {
  return iconMap[iconName] || Folder
}

// 方法
const toggleExpanded = () => {
  if (hasChildren.value) {
    expanded.value = !expanded.value
  }
}

const handleEdit = () => {
  emit('edit', props.category)
}

const handleDelete = () => {
  emit('delete', props.category)
}

const handleAddChild = () => {
  emit('add-child', props.category)
}

// 监听搜索关键词，有搜索时自动展开
watch(() => props.searchKeyword, (newKeyword) => {
  if (newKeyword && hasChildren.value) {
    expanded.value = true
  }
})
</script>

<style scoped>
.category-tree-node {
  margin-bottom: 2px;
}

.category-item {
  display: flex;
  align-items: center;
  height: 44px;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  gap: 8px;
  position: relative;
}

.category-item:hover {
  background-color: #f8fafc;
}

.category-tree-node.highlighted .category-item {
  background-color: #fff7e6;
  border: 1px solid #ffd591;
}

.expand-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.expand-icon.has-children {
  cursor: pointer;
  color: #909399;
}

.expand-icon.has-children:hover {
  background-color: #e9ecef;
  color: #606266;
}

.expand-icon.expanded {
  transform: rotate(90deg);
}

.category-icon {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 500;
  flex-shrink: 0;
}

.category-info {
  flex: 1;
  min-width: 0;
}

.category-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.category-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.category-name.highlighted-text {
  color: #d48806;
}

.category-name :deep(mark) {
  background-color: #fff2cc;
  color: #d48806;
  padding: 0 2px;
  border-radius: 2px;
}

.category-badges {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.todo-badge {
  --el-badge-bg-color: #409eff;
  --el-badge-border-color: #409eff;
}

.children-count {
  font-size: 11px;
  color: #909399;
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 10px;
}

.category-description {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.category-description :deep(mark) {
  background-color: #fff2cc;
  color: #d48806;
  padding: 0 2px;
  border-radius: 2px;
}

.category-actions {
  opacity: 0;
  transition: opacity 0.2s ease;
  flex-shrink: 0;
}

.category-item:hover .category-actions {
  opacity: 1;
}

.children-container {
  margin-left: 12px;
  border-left: 2px solid #f0f0f0;
  padding-left: 8px;
}

.delete-item {
  color: #f56c6c;
}

.delete-item:hover {
  background-color: #fef0f0;
}

/* 动画效果 */
.expand-enter-active,
.expand-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  max-height: 0;
  opacity: 0;
  transform: translateY(-10px);
}

.expand-enter-to,
.expand-leave-from {
  max-height: 1000px;
  opacity: 1;
  transform: translateY(0);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .category-item {
    height: 40px;
    padding: 6px 12px;
  }
  
  .category-badges {
    flex-direction: column;
    gap: 4px;
  }
  
  .children-count {
    display: none;
  }
}
</style>