<template>
  <div class="hover-category-selector">
    <!-- 触发区域 -->
    <div 
      class="selector-trigger"
      @mouseenter="showSelector"
      @mouseleave="hideSelector"
      @click="toggleSelector"
    >
      <div class="selected-display">
        <span v-if="selectedPath.length > 0" class="selected-text">
          {{ selectedPath.map(item => item.name).join(' > ') }}
        </span>
        <span v-else class="placeholder">{{ placeholder }}</span>
      </div>
      <el-icon class="dropdown-icon" :class="{ 'is-reverse': visible }">
        <ArrowDown />
      </el-icon>
    </div>

    <!-- 悬停选择面板 -->
    <div 
      v-if="visible"
      class="selector-panel"
      :class="{ 'panel-left': panelOnLeft }"
      @mouseenter="keepShow"
      @mouseleave="hideSelector"
    >
      <div class="breadcrumb" v-if="selectedPath.length > 0">
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

      <div class="categories-container">
        <div 
          v-for="(level, levelIndex) in levels" 
          :key="levelIndex" 
          class="category-level"
        >
          <div class="level-header">{{ getLevelTitle(levelIndex) }}</div>
          <div class="category-list">
            <div
              v-for="category in level"
              :key="category.id"
              class="category-option"
              :class="{ 
                'selected': isSelected(category, levelIndex),
                'has-children': hasChildren(category)
              }"
              @click="selectCategory(category, levelIndex)"
              @mouseenter="previewChildren(category, levelIndex)"
            >
              <div 
                class="category-color" 
                :style="{ backgroundColor: category.color || '#409eff' }"
              ></div>
              <span class="category-name">{{ category.name }}</span>
              <el-icon v-if="hasChildren(category)" class="arrow-icon">
                <ArrowRight />
              </el-icon>
            </div>
            <div v-if="level.length === 0" class="empty-level">
              暂无分类
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ArrowDown, ArrowRight } from '@element-plus/icons-vue'

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
  excludeIds: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const visible = ref(false)
const selectedCategories = ref([])
const levels = ref([])
const panelOnLeft = ref(false)

let hoverTimer = null

const selectedPath = computed(() => {
  return selectedCategories.value.filter(Boolean)
})

const initializeLevels = () => {
  if (!props.categories || props.categories.length === 0) {
    levels.value = [[]]
    selectedCategories.value = []
    return
  }
  
  const rootCategories = props.categories.filter(cat => 
    cat && !cat.parent_id && !props.excludeIds.includes(cat.id)
  )
  levels.value = [rootCategories]
  selectedCategories.value = []
}

const initializeSelection = () => {
  if (!props.modelValue || !props.categories || props.categories.length === 0) {
    initializeLevels()
    return
  }

  const selectedCategory = props.categories.find(cat => cat && cat.id === props.modelValue)
  if (!selectedCategory) {
    initializeLevels()
    return
  }

  const path = buildCategoryPath(selectedCategory)
  selectedCategories.value = path
  
  levels.value = []
  let currentParentId = null
  
  path.forEach((category, index) => {
    const currentLevelOptions = props.categories.filter(cat => {
      const matchesParent = cat.parent_id === currentParentId
      const notExcluded = !props.excludeIds.includes(cat.id)
      return matchesParent && notExcluded
    })
    levels.value[index] = currentLevelOptions
    currentParentId = category.id
  })
  
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

const buildCategoryPath = (category) => {
  const path = []
  let current = category
  
  while (current) {
    path.unshift(current)
    current = props.categories.find(cat => cat.id === current.parent_id)
  }
  
  return path
}

const hasChildren = (category) => {
  if (!props.categories || !category) return false
  return props.categories.some(cat => 
    cat && cat.parent_id === category.id && !props.excludeIds.includes(cat.id)
  )
}

const isSelected = (category, levelIndex) => {
  return selectedCategories.value[levelIndex]?.id === category.id
}

const getLevelTitle = (levelIndex) => {
  const titles = ['一级分类', '二级分类', '三级分类', '四级分类', '五级分类']
  return titles[levelIndex] || `${levelIndex + 1}级分类`
}

const selectCategory = (category, levelIndex) => {
  selectedCategories.value = selectedCategories.value.slice(0, levelIndex)
  selectedCategories.value[levelIndex] = category
  
  levels.value = levels.value.slice(0, levelIndex + 1)
  
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
  
  const selectedValue = category.id
  emit('update:modelValue', selectedValue)
  emit('change', selectedValue, buildCategoryPath(category))
}

const previewChildren = (category, levelIndex) => {
  if (!hasChildren(category)) return
  
  // 预览时暂时显示子级选项
  const childOptions = props.categories.filter(cat => {
    const matchesParent = cat.parent_id === category.id
    const notExcluded = !props.excludeIds.includes(cat.id)
    return matchesParent && notExcluded
  })
  
  if (childOptions.length > 0) {
    levels.value = levels.value.slice(0, levelIndex + 1)
    levels.value.push(childOptions)
  }
}

const clearSelection = () => {
  selectedCategories.value = []
  emit('update:modelValue', null)
  emit('change', null, [])
  initializeLevels()
}

const showSelector = () => {
  clearTimeout(hoverTimer)
  hoverTimer = setTimeout(() => {
    // 检测屏幕空间，决定面板位置
    const triggerElement = document.querySelector('.hover-category-selector .selector-trigger')
    if (triggerElement) {
      const rect = triggerElement.getBoundingClientRect()
      const windowWidth = window.innerWidth
      const panelWidth = 300 // 面板最小宽度
      const spaceOnRight = windowWidth - rect.right
      
      // 如果右边空间不足，显示在左边
      panelOnLeft.value = spaceOnRight < panelWidth
    }
    
    visible.value = true
  }, 200)
}

const keepShow = () => {
  clearTimeout(hoverTimer)
  visible.value = true
}

const hideSelector = () => {
  clearTimeout(hoverTimer)
  hoverTimer = setTimeout(() => {
    visible.value = false
  }, 150)
}

const toggleSelector = () => {
  visible.value = !visible.value
}

watch(() => props.categories, () => {
  initializeSelection()
}, { immediate: true })

watch(() => props.modelValue, () => {
  initializeSelection()
})

watch(() => props.excludeIds, () => {
  initializeSelection()
})
</script>

<style scoped>
.hover-category-selector {
  position: relative;
  width: 100%;
}

.selector-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  transition: border-color 0.2s;
  min-height: 32px;
}

.selector-trigger:hover {
  border-color: #c0c4cc;
}

.selected-display {
  flex: 1;
  overflow: hidden;
}

.selected-text {
  font-size: 13px;
  color: #303133;
}

.placeholder {
  font-size: 13px;
  color: #c0c4cc;
}

.dropdown-icon {
  font-size: 12px;
  color: #c0c4cc;
  transition: transform 0.3s;
}

.dropdown-icon.is-reverse {
  transform: rotateZ(180deg);
}

.selector-panel {
  position: absolute;
  top: 0;
  left: 100%;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 9999;
  max-height: 240px;
  overflow: hidden;
  margin-left: 4px;
  min-width: 300px;
  width: auto;
}

.selector-panel.panel-left {
  left: auto;
  right: 100%;
  margin-left: 0;
  margin-right: 4px;
}

.breadcrumb {
  padding: 6px 8px;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  font-size: 11px;
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

.categories-container {
  display: flex;
  max-height: 200px;
}

.category-level {
  flex: 1;
  border-right: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
  min-width: 100px;
  max-width: 140px;
}

.category-level:last-child {
  border-right: none;
}

.level-header {
  padding: 4px 8px;
  background: #fafafa;
  border-bottom: 1px solid #ebeef5;
  font-size: 10px;
  font-weight: 500;
  color: #606266;
}

.category-list {
  flex: 1;
  overflow-y: auto;
}

.category-option {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  font-size: 11px;
  gap: 4px;
  border-bottom: 1px solid #f9f9f9;
  min-height: 28px;
}

.category-option:hover {
  background-color: #f0f9ff;
}

.category-option.selected {
  background-color: #e1f3d8;
  color: #67c23a;
  font-weight: 500;
}

.category-color {
  width: 8px;
  height: 8px;
  border-radius: 2px;
  flex-shrink: 0;
}

.category-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.arrow-icon {
  font-size: 10px;
  color: #c0c4cc;
}

.category-option.selected .arrow-icon {
  color: #67c23a;
}

.empty-level {
  padding: 12px;
  text-align: center;
  color: #c0c4cc;
  font-size: 10px;
}
</style>