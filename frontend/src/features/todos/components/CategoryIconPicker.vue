<template>
  <div class="category-icon-picker">
    <div class="selected-icon" @click="showPicker = !showPicker">
      <div class="icon-preview">
        <el-icon size="16">
          <component :is="selectedIconComponent" />
        </el-icon>
      </div>
      <span class="selected-name">{{ getIconName(modelValue) }}</span>
      <el-icon class="dropdown-arrow" :class="{ 'expanded': showPicker }">
        <ArrowDown />
      </el-icon>
    </div>

    <div v-if="showPicker" class="icon-picker-panel">
      <div class="picker-header">
        <span>选择图标</span>
        <el-input
          v-model="searchKeyword"
          placeholder="搜索图标..."
          size="small"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <div class="icon-categories">
        <div
          v-for="category in iconCategories"
          :key="category.name"
          :class="['category-tab', { active: activeCategory === category.name }]"
          @click="activeCategory = category.name"
        >
          {{ category.label }}
        </div>
      </div>

      <div class="icon-grid">
        <div
          v-for="iconItem in filteredIcons"
          :key="iconItem.name"
          :class="['icon-item', { selected: modelValue === iconItem.name }]"
          @click="selectIcon(iconItem.name)"
          :title="iconItem.label"
        >
          <el-icon size="18">
            <component :is="iconItem.component" />
          </el-icon>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { 
  ArrowDown, Search,
  // 确认可用的图标
  Folder, Star, User, Setting, Calendar, Bell,
  Sunny, Moon, Plus, Edit, Delete, More
} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: 'folder'
  }
})

const emit = defineEmits(['update:modelValue'])

// 响应式状态
const showPicker = ref(false)
const searchKeyword = ref('')
const activeCategory = ref('folder')

// 图标数据
const iconCategories = [
  { name: 'folder', label: '文件夹' },
  { name: 'mark', label: '标记' },
  { name: 'time', label: '时间' },
  { name: 'user', label: '用户' },
  { name: 'work', label: '工作' },
  { name: 'life', label: '生活' },
  { name: 'study', label: '学习' },
  { name: 'health', label: '健康' },
  { name: 'entertainment', label: '娱乐' },
  { name: 'transport', label: '交通' },
  { name: 'communication', label: '通讯' },
  { name: 'other', label: '其他' }
]

const allIcons = [
  // 基础分类
  { name: 'folder', label: '文件夹', component: Folder, category: 'folder' },
  { name: 'star', label: '星标', component: Star, category: 'mark' },
  { name: 'user', label: '用户', component: User, category: 'user' },
  { name: 'setting', label: '设置', component: Setting, category: 'work' },
  { name: 'calendar', label: '日历', component: Calendar, category: 'time' },
  { name: 'bell', label: '通知', component: Bell, category: 'communication' },
  { name: 'sunny', label: '太阳', component: Sunny, category: 'health' },
  { name: 'moon', label: '月亮', component: Moon, category: 'health' },
  { name: 'edit', label: '编辑', component: Edit, category: 'study' },
  { name: 'plus', label: '添加', component: Plus, category: 'other' },
  { name: 'delete', label: '删除', component: Delete, category: 'other' },
  { name: 'more', label: '更多', component: More, category: 'other' }
]

// 计算属性
const selectedIconComponent = computed(() => {
  const iconItem = allIcons.find(icon => icon.name === props.modelValue)
  return iconItem ? iconItem.component : Folder
})

const filteredIcons = computed(() => {
  let icons = allIcons

  // 按分类筛选
  if (activeCategory.value) {
    icons = icons.filter(icon => icon.category === activeCategory.value)
  }

  // 按搜索关键词筛选
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    icons = icons.filter(icon => 
      icon.label.toLowerCase().includes(keyword) ||
      icon.name.toLowerCase().includes(keyword)
    )
  }

  return icons
})

// 方法
const getIconName = (iconName) => {
  const iconItem = allIcons.find(icon => icon.name === iconName)
  return iconItem ? iconItem.label : '文件夹'
}

const selectIcon = (iconName) => {
  emit('update:modelValue', iconName)
  showPicker.value = false
}

// 点击外部关闭选择器
const handleClickOutside = (event) => {
  if (!event.target.closest('.category-icon-picker')) {
    showPicker.value = false
  }
}

// 监听点击外部事件
watch(showPicker, (visible) => {
  if (visible) {
    document.addEventListener('click', handleClickOutside)
  } else {
    document.removeEventListener('click', handleClickOutside)
    searchKeyword.value = ''
  }
})
</script>

<style scoped>
.category-icon-picker {
  position: relative;
  width: 100%;
}

.selected-icon {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  transition: border-color 0.2s;
}

.selected-icon:hover {
  border-color: #c0c4cc;
}

.icon-preview {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-radius: 4px;
  color: #606266;
}

.selected-name {
  flex: 1;
  font-size: 14px;
  color: #606266;
}

.dropdown-arrow {
  color: #c0c4cc;
  transition: transform 0.3s;
}

.dropdown-arrow.expanded {
  transform: rotate(180deg);
}

.icon-picker-panel {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  z-index: 9999;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  margin-top: 4px;
  max-height: 400px;
  overflow: hidden;
}

.picker-header {
  padding: 12px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.picker-header span {
  font-weight: 500;
  color: #303133;
  white-space: nowrap;
}

.icon-categories {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  padding: 8px 12px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafbfc;
}

.category-tab {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  color: #606266;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.category-tab:hover {
  background: #e9ecef;
}

.category-tab.active {
  background: #409eff;
  color: white;
}

.icon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
  gap: 4px;
  padding: 12px;
  max-height: 240px;
  overflow-y: auto;
}

.icon-item {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  color: #606266;
}

.icon-item:hover {
  background: #f0f9ff;
  color: #409eff;
}

.icon-item.selected {
  background: #409eff;
  color: white;
}

/* 滚动条样式 */
.icon-grid::-webkit-scrollbar {
  width: 6px;
}

.icon-grid::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.icon-grid::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.icon-grid::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>