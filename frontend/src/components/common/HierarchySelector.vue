<template>
  <div class="hierarchy-selector">
    <!-- 面包屑显示当前选择路径 -->
    <div v-if="selectedPath.length > 0" class="breadcrumb">
      <span v-for="(item, index) in selectedPath" :key="item.id" class="breadcrumb-item">
        {{ item.name }}
        <el-icon v-if="index < selectedPath.length - 1" class="separator">
          <ArrowRight />
        </el-icon>
      </span>
      <el-button 
        type="primary" 
        text 
        size="small" 
        @click="clearSelection"
        class="clear-btn"
      >
        清空
      </el-button>
    </div>

    <!-- 层级选择区域 -->
    <div class="selector-panels">
      <div 
        v-for="(level, levelIndex) in levels" 
        :key="levelIndex" 
        class="selector-panel"
      >
        <div class="panel-header">{{ getPanelTitle(levelIndex) }}</div>
        <div class="panel-content">
          <div
            v-for="category in level"
            :key="category.id"
            class="category-option"
            :class="{ 
              'selected': isSelected(category, levelIndex),
              'has-children': hasChildren(category)
            }"
            @click="selectCategory(category, levelIndex)"
          >
            <span class="category-name">{{ category.name }}</span>
            <el-icon v-if="hasChildren(category)" class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>
          <div v-if="level.length === 0" class="empty-panel">
            暂无下级分类
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ArrowRight } from '@element-plus/icons-vue'

const props = defineProps({
  categories: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: [Number, null],
    default: null
  },
  placeholder: {
    type: String,
    default: '请选择分类'
  },
  // 排除某些分类ID（编辑时防止循环引用）
  excludeIds: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const selectedCategories = ref([])
const levels = ref([])

// 计算选择路径的面包屑
const selectedPath = computed(() => {
  return selectedCategories.value.filter(Boolean)
})

// 初始化层级数据
const initializeLevels = () => {
  // 获取根分类（过滤排除的ID）
  const rootCategories = props.categories.filter(cat => 
    !cat.parent_id && !props.excludeIds.includes(cat.id)
  )
  levels.value = [rootCategories]
  selectedCategories.value = []
}

// 根据已选值初始化选择状态
const initializeSelection = () => {
  if (!props.modelValue) {
    initializeLevels()
    return
  }

  // 找到选中的分类
  const selectedCategory = props.categories.find(cat => cat.id === props.modelValue)
  if (!selectedCategory) {
    initializeLevels()
    return
  }

  // 构建选择路径
  const path = buildCategoryPath(selectedCategory)
  selectedCategories.value = path
  
  // 构建各个层级的选项
  levels.value = []
  let currentParentId = null
  
  path.forEach((category, index) => {
    // 获取当前层级的所有选项
    const currentLevelOptions = props.categories.filter(cat => {
      const matchesParent = cat.parent_id === currentParentId
      const notExcluded = !props.excludeIds.includes(cat.id)
      return matchesParent && notExcluded
    })
    levels.value[index] = currentLevelOptions
    currentParentId = category.id
  })
  
  // 如果最后选中的分类有子分类，显示下一层
  const lastSelected = path[path.length - 1]
  if (hasChildren(lastSelected)) {
    const childOptions = props.categories.filter(cat => {
      const matchesParent = cat.parent_id === lastSelected.id
      const notExcluded = !props.excludeIds.includes(cat.id)
      return matchesParent && notExcluded
    })
    if (childOptions.length > 0) {
      levels.value.push(childOptions)
    }
  }
}

// 构建分类的完整路径
const buildCategoryPath = (category) => {
  const path = []
  let current = category
  
  while (current) {
    path.unshift(current)
    current = props.categories.find(cat => cat.id === current.parent_id)
  }
  
  return path
}

// 判断分类是否有子分类（排除被禁用的）
const hasChildren = (category) => {
  return props.categories.some(cat => 
    cat.parent_id === category.id && !props.excludeIds.includes(cat.id)
  )
}

// 判断分类是否被选中
const isSelected = (category, levelIndex) => {
  return selectedCategories.value[levelIndex]?.id === category.id
}

// 获取面板标题
const getPanelTitle = (levelIndex) => {
  const titles = ['一级分类', '二级分类', '三级分类', '四级分类', '五级分类']
  return titles[levelIndex] || `${levelIndex + 1}级分类`
}

// 选择分类
const selectCategory = (category, levelIndex) => {
  // 截取到当前层级
  selectedCategories.value = selectedCategories.value.slice(0, levelIndex)
  selectedCategories.value[levelIndex] = category
  
  // 截取层级显示到当前层级+1
  levels.value = levels.value.slice(0, levelIndex + 1)
  
  // 如果有子分类，显示下一层
  if (hasChildren(category)) {
    const childOptions = props.categories.filter(cat => {
      const matchesParent = cat.parent_id === category.id
      const notExcluded = !props.excludeIds.includes(cat.id)
      return matchesParent && notExcluded
    })
    if (childOptions.length > 0) {
      levels.value.push(childOptions)
    }
  }
  
  // 发出选择事件
  const selectedValue = category.id
  emit('update:modelValue', selectedValue)
  emit('change', selectedValue, buildCategoryPath(category))
}

// 清空选择
const clearSelection = () => {
  selectedCategories.value = []
  emit('update:modelValue', null)
  emit('change', null, [])
  initializeLevels()
}

// 监听分类数据变化
watch(() => props.categories, () => {
  initializeSelection()
}, { immediate: true })

// 监听选中值变化
watch(() => props.modelValue, () => {
  initializeSelection()
})

// 监听排除ID变化
watch(() => props.excludeIds, () => {
  initializeSelection()
})
</script>

<style scoped>
.hierarchy-selector {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: white;
}

.breadcrumb {
  padding: 8px 12px;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #606266;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.separator {
  font-size: 10px;
  color: #c0c4cc;
  margin: 0 4px;
}

.clear-btn {
  margin-left: auto;
  font-size: 12px;
}

.selector-panels {
  display: flex;
  min-height: 200px;
  max-height: 300px;
}

.selector-panel {
  flex: 1;
  border-right: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
}

.selector-panel:last-child {
  border-right: none;
}

.panel-header {
  padding: 8px 12px;
  background: #fafafa;
  border-bottom: 1px solid #ebeef5;
  font-size: 12px;
  font-weight: 500;
  color: #606266;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
}

.category-option {
  padding: 8px 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  transition: background-color 0.2s;
  font-size: 13px;
  border-bottom: 1px solid #f5f7fa;
}

.category-option:hover {
  background-color: #f0f9ff;
}

.category-option.selected {
  background-color: #e1f3d8;
  color: #67c23a;
  font-weight: 500;
}

.category-option.has-children {
  color: #409eff;
}

.category-name {
  flex: 1;
}

.arrow-icon {
  font-size: 12px;
  color: #c0c4cc;
}

.category-option.selected .arrow-icon {
  color: #67c23a;
}

.empty-panel {
  padding: 20px;
  text-align: center;
  color: #c0c4cc;
  font-size: 12px;
}
</style>