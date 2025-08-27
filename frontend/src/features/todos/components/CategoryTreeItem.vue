<template>
  <div class="category-tree-item">
    <div 
      class="category-item" 
      :style="{ paddingLeft: level * 16 + 12 + 'px' }"
    >
      <div class="category-info">
        <div 
          class="category-color" 
          :style="{ backgroundColor: category.color || '#409eff' }"
        ></div>
        <div class="category-details">
          <div class="category-main">
            <span class="category-name">{{ category.name }}</span>
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
    
    <!-- 递归渲染子分类 -->
    <template v-if="children.length > 0">
      <category-tree-item
        v-for="child in children"
        :key="child.id"
        :category="child"
        :level="level + 1"
        :get-todo-count="getTodoCount"
        @edit="$emit('edit', $event)"
        @delete="$emit('delete', $event)"
        @add-child="$emit('add-child', $event)"
      />
    </template>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'

const props = defineProps({
  category: {
    type: Object,
    required: true
  },
  level: {
    type: Number,
    default: 0
  },
  getTodoCount: {
    type: Function,
    required: true
  }
})

defineEmits(['edit', 'delete', 'add-child'])

const allCategories = inject('categories', [])

const children = computed(() => {
  return allCategories.value?.filter(cat => cat.parent_id === props.category.id) || []
})
</script>

<style scoped>
.category-tree-item {
  margin-bottom: 1px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid #f0f0f0;
  background: white;
  transition: background-color 0.2s;
  min-height: 42px;
}

.category-item:hover {
  background-color: #f8f9fa;
}

.category-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.category-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
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
  font-size: 13px;
  font-weight: 500;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
</style>