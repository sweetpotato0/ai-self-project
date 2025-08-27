<template>
  <div class="category-flat-item">
    <div 
      class="category-item"
      @mouseenter="showChildren"
      @mouseleave="hideChildren"
    >
      <div class="category-info">
        <!-- 层级缩进指示器 -->
        <div class="level-indicator" :style="{ marginLeft: (category.level - 1) * 12 + 'px' }">
          <div 
            class="category-color" 
            :style="{ backgroundColor: category.color || '#409eff' }"
          ></div>
        </div>
        
        <div class="category-details">
          <div class="category-main">
            <span class="category-name">{{ category.name }}</span>
            <span v-if="hasChildren" class="has-children-icon">
              <el-icon><ArrowRight /></el-icon>
            </span>
            <span class="category-count">{{ getTodoCount(category.id) }} 个任务</span>
          </div>
          <div v-if="category.description" class="category-description">
            {{ category.description }}
          </div>
        </div>
      </div>
      
      <div class="category-actions">
        <el-button 
          size="small" 
          type="primary" 
          text 
          @click="$emit('add-child', category)"
          title="添加子分类"
        >
          <el-icon><Plus /></el-icon>
        </el-button>
        <el-button 
          size="small" 
          type="primary" 
          text 
          @click="$emit('edit', category)"
          title="编辑分类"
        >
          <el-icon><Edit /></el-icon>
        </el-button>
        <el-button 
          size="small" 
          type="danger" 
          text 
          @click="$emit('delete', category)"
          title="删除分类"
        >
          <el-icon><Delete /></el-icon>
        </el-button>
      </div>
    </div>
    
    <!-- 悬停显示的子分类面板 -->
    <div 
      v-if="showChildPanel && hasChildren"
      class="children-panel"
      @mouseenter="keepShowChildren"
      @mouseleave="hideChildren"
    >
      <div class="panel-header">
        <span>{{ category.name }} 的子分类</span>
      </div>
      <div class="children-list">
        <div
          v-for="child in children"
          :key="child.id"
          class="child-item"
          @click="handleChildClick(child)"
        >
          <div 
            class="child-color" 
            :style="{ backgroundColor: child.color || '#409eff' }"
          ></div>
          <span class="child-name">{{ child.name }}</span>
          <span class="child-count">{{ getTodoCount(child.id) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Plus, Edit, Delete, ArrowRight } from '@element-plus/icons-vue'

const props = defineProps({
  category: {
    type: Object,
    required: true
  },
  categories: {
    type: Array,
    default: () => []
  },
  getTodoCount: {
    type: Function,
    required: true
  }
})

defineEmits(['edit', 'delete', 'add-child'])

const showChildPanel = ref(false)
let hoverTimer = null

// 获取子分类
const children = computed(() => {
  return props.categories.filter(cat => cat.parent_id === props.category.id)
})

// 是否有子分类
const hasChildren = computed(() => {
  return children.value.length > 0
})

// 显示子分类面板
const showChildren = () => {
  if (!hasChildren.value) return
  
  clearTimeout(hoverTimer)
  hoverTimer = setTimeout(() => {
    showChildPanel.value = true
  }, 300) // 300ms延迟显示
}

// 保持显示子分类面板
const keepShowChildren = () => {
  clearTimeout(hoverTimer)
  showChildPanel.value = true
}

// 隐藏子分类面板
const hideChildren = () => {
  clearTimeout(hoverTimer)
  hoverTimer = setTimeout(() => {
    showChildPanel.value = false
  }, 100) // 100ms延迟隐藏
}

// 处理子分类点击
const handleChildClick = (child) => {
  // 可以在这里处理子分类的点击事件，比如跳转或展开详情
  console.log('Child category clicked:', child)
}
</script>

<style scoped>
.category-flat-item {
  position: relative;
  margin-bottom: 1px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  border-bottom: 1px solid #f0f0f0;
  background: white;
  transition: background-color 0.2s;
  min-height: 48px;
}

.category-item:hover {
  background-color: #f8f9fa;
}

.category-info {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.level-indicator {
  display: flex;
  align-items: center;
  margin-right: 12px;
}

.category-color {
  width: 14px;
  height: 14px;
  border-radius: 3px;
  flex-shrink: 0;
}

.category-details {
  flex: 1;
  min-width: 0;
}

.category-main {
  display: flex;
  align-items: center;
  gap: 8px;
}

.category-name {
  color: #303133;
  font-size: 14px;
  font-weight: 500;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.has-children-icon {
  color: #409eff;
  font-size: 12px;
  display: flex;
  align-items: center;
}

.category-count {
  color: #909399;
  font-size: 11px;
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 10px;
  flex-shrink: 0;
}

.category-description {
  color: #666;
  font-size: 11px;
  line-height: 1.3;
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.category-actions {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s;
}

.category-item:hover .category-actions {
  opacity: 1;
}

.category-actions .el-button {
  padding: 4px;
  margin: 0;
}

/* 悬停面板样式 */
.children-panel {
  position: absolute;
  top: 100%;
  left: 16px;
  right: 16px;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 9999;
  max-height: 200px;
  overflow-y: auto;
  margin-top: 2px;
}

.panel-header {
  padding: 4px 8px;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  font-size: 10px;
  font-weight: 500;
  color: #606266;
}

.children-list {
  max-height: 120px;
  overflow-y: auto;
}

.child-item {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  gap: 6px;
  border-bottom: 1px solid #f9f9f9;
  min-height: 24px;
}

.child-item:hover {
  background-color: #f0f9ff;
}

.child-item:last-child {
  border-bottom: none;
}

.child-color {
  width: 10px;
  height: 10px;
  border-radius: 2px;
  flex-shrink: 0;
}

.child-name {
  flex: 1;
  font-size: 11px;
  color: #606266;
}

.child-count {
  font-size: 9px;
  color: #c0c4cc;
  background: #f5f7fa;
  padding: 1px 3px;
  border-radius: 6px;
}
</style>