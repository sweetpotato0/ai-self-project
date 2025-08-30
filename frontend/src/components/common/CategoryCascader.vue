<template>
  <el-cascader
    v-model="selectedPath"
    :options="cascaderOptions"
    :props="cascaderProps"
    :placeholder="placeholder"
    clearable
    filterable
    :show-all-levels="true"
    :separator="separator"
    :size="size"
    :disabled="disabled"
    @change="handleChange"
    class="category-cascader"
  >
    <template #default="{ node, data }">
      <div class="cascader-node">
        <div 
          v-if="data.icon" 
          class="node-icon" 
          :style="{ backgroundColor: data.color || '#409eff' }"
        >
          <el-icon size="14">
            <component :is="getIconComponent(data.icon)" />
          </el-icon>
        </div>
        <span class="node-name">{{ data.name }}</span>
        <span v-if="showCounts && data.todoCount > 0" class="node-count">
          {{ data.todoCount }}
        </span>
      </div>
    </template>
  </el-cascader>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { 
  Folder, Star, User, Setting, Calendar, Bell, 
  Sunny, Moon, Edit, Plus, Delete, More
} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: [Number, null],
    default: null
  },
  categories: {
    type: Array,
    default: () => []
  },
  placeholder: {
    type: String,
    default: '请选择分类'
  },
  size: {
    type: String,
    default: 'default'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  showAllLevels: {
    type: Boolean,
    default: false
  },
  separator: {
    type: String,
    default: ' / '
  },
  showCounts: {
    type: Boolean,
    default: false
  },
  excludeIds: {
    type: Array,
    default: () => []
  },
  getTodoCount: {
    type: Function,
    default: () => 0
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

// 响应式状态
const selectedPath = ref([])

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

const getIconComponent = (iconName) => {
  return iconMap[iconName] || Folder
}

// 级联选择器配置
const cascaderProps = {
  value: 'id',
  label: 'name',
  children: 'children',
  checkStrictly: true, // 可选择任意级别
  emitPath: true // 返回完整路径数组，支持面包屑显示
}

// 构建级联选择器选项
const cascaderOptions = computed(() => {
  if (!props.categories || props.categories.length === 0) {
    return []
  }

  const buildTree = (items, parentId = null) => {
    return items
      .filter(item => {
        // 过滤掉排除的ID
        if (props.excludeIds.includes(item.id)) {
          return false
        }
        return item.parent_id === parentId
      })
      .sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))
      .map(item => {
        const children = buildTree(items, item.id)
        return {
          id: item.id,
          name: item.name,
          icon: item.icon,
          color: item.color,
          todoCount: props.getTodoCount ? props.getTodoCount(item.id) : 0,
          children: children.length > 0 ? children : undefined
        }
      })
  }
  
  return buildTree(props.categories)
})

// 根据选中值构建路径
const initializeSelection = () => {
  if (!props.modelValue || !props.categories || props.categories.length === 0) {
    selectedPath.value = []
    return
  }

  const findPath = (categoryId, categories, path = []) => {
    for (const category of categories) {
      const currentPath = [...path, category.id]
      
      if (category.id === categoryId) {
        return currentPath
      }
      
      if (category.children && category.children.length > 0) {
        const childPath = findPath(categoryId, category.children, currentPath)
        if (childPath) {
          return childPath
        }
      }
    }
    return null
  }

  const path = findPath(props.modelValue, cascaderOptions.value)
  selectedPath.value = path || []
}

// 处理选择变化
const handleChange = (value) => {
  // 当emitPath为true时，value是路径数组，需要提取最后一个值
  const selectedValue = Array.isArray(value) && value.length > 0 ? 
    value[value.length - 1] : null
  
  emit('update:modelValue', selectedValue)
  emit('change', selectedValue)
}

// 监听变化
watch(() => props.modelValue, () => {
  initializeSelection()
}, { immediate: true })

watch(() => props.categories, () => {
  initializeSelection()
}, { deep: true })
</script>

<style scoped>
.category-cascader {
  width: 100%;
}

.cascader-node {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.node-icon {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.node-name {
  flex: 1;
  font-size: 14px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-count {
  font-size: 11px;
  color: #909399;
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 10px;
  flex-shrink: 0;
}

/* 自定义级联选择器样式 */
.category-cascader :deep(.el-cascader__tags) {
  max-width: 100%;
}

.category-cascader :deep(.el-cascader__search-input) {
  font-size: 14px;
}

.category-cascader :deep(.el-cascader-panel__content) {
  min-height: 200px;
}

.category-cascader :deep(.el-cascader-node) {
  padding: 8px 12px;
}

.category-cascader :deep(.el-cascader-node:hover) {
  background-color: #f5f7fa;
}

.category-cascader :deep(.el-cascader-node.is-selectable.in-active-path) {
  color: #409eff;
  background-color: #f0f9ff;
}

.category-cascader :deep(.el-cascader-node.is-selectable.is-active) {
  color: #409eff;
  background-color: #e1f3d8;
  font-weight: 600;
}
</style>