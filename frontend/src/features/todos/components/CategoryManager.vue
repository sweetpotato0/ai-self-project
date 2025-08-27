<template>
  <div class="category-manager">
    <div class="category-header">
      <h3>任务分类管理</h3>
      <el-button type="primary" size="small" @click="openCategoryDialog()">
        <el-icon><Plus /></el-icon>
        新增分类
      </el-button>
    </div>

    <div class="category-list">
      <div v-if="topLevelCategories.length === 0" class="empty-state">
        <el-icon><Folder /></el-icon>
        <p>暂无分类，点击上方按钮创建第一个分类</p>
      </div>

      <!-- 只显示顶级分类 -->
      <div v-for="category in topLevelCategories" :key="category.id">
        <category-flat-item 
          :category="category" 
          :categories="categories"
          @edit="openCategoryDialog"
          @delete="deleteCategory"
          @add-child="(parent) => openCategoryDialog(null, parent.id)"
          :get-todo-count="getCategoryTodoCount"
        />
      </div>
    </div>

    <!-- 分类编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingCategory ? '编辑分类' : '新增分类'"
      width="400px"
    >
      <el-form :model="categoryForm" label-width="80px" ref="categoryFormRef">
        <el-form-item label="分类名称" prop="name" required>
          <el-input 
            v-model="categoryForm.name" 
            placeholder="请输入分类名称"
            maxlength="20"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="父分类">
          <hover-category-selector
            v-model="categoryForm.parent_id"
            :categories="categories"
            :exclude-ids="editingCategory ? [editingCategory.id] : []"
            placeholder="请选择父分类（留空为顶级分类）"
          />
        </el-form-item>
        
        <el-form-item label="分类颜色">
          <el-color-picker 
            v-model="categoryForm.color"
            :predefine="predefinedColors"
          />
        </el-form-item>
        
        <el-form-item label="分类描述">
          <el-input
            v-model="categoryForm.description"
            type="textarea"
            placeholder="请输入分类描述（可选）"
            maxlength="100"
            show-word-limit
            :rows="3"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveCategory" :loading="loading">
            {{ editingCategory ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Folder } from '@element-plus/icons-vue'
import { useTodoStore } from '@/stores/todo'
import CategoryFlatItem from './CategoryFlatItem.vue'
import HoverCategorySelector from '@/components/common/HoverCategorySelector.vue'

const todoStore = useTodoStore()
const categoryFormRef = ref()

// 响应式数据
const dialogVisible = ref(false)
const editingCategory = ref(null)
const loading = ref(false)

const categoryForm = ref({
  name: '',
  color: '#409eff',
  description: '',
  parent_id: null
})

const predefinedColors = [
  '#409eff', '#67c23a', '#e6a23c', '#f56c6c', 
  '#909399', '#c45656', '#73767a', '#b88230',
  '#c0c4cc', '#606266', '#303133', '#8957a1'
]

// 计算属性
const categories = computed(() => todoStore.categories)

// 只获取顶级分类（没有parent_id的分类）
const topLevelCategories = computed(() => {
  return categories.value.filter(cat => !cat.parent_id)
})

// 获取可用的父分类（排除当前编辑的分类及其子分类）
const availableParentCategories = computed(() => {
  let available = categories.value.filter(cat => {
    // 如果是编辑模式，排除当前分类及其子分类
    if (editingCategory.value) {
      return cat.id !== editingCategory.value.id && 
             cat.parent_id !== editingCategory.value.id
    }
    return true
  })
  return available
})

// 方法
const getCategoryTodoCount = (categoryId) => {
  return todoStore.todos.filter(todo => todo.category_id === categoryId).length
}

const openCategoryDialog = (category = null, parentId = null) => {
  editingCategory.value = category
  if (category) {
    categoryForm.value = {
      name: category.name,
      color: category.color || '#409eff',
      description: category.description || '',
      parent_id: category.parent_id || null
    }
  } else {
    categoryForm.value = {
      name: '',
      color: '#409eff',
      description: '',
      parent_id: parentId || null
    }
  }
  dialogVisible.value = true
}

const saveCategory = async () => {
  if (!categoryForm.value.name.trim()) {
    ElMessage.error('请输入分类名称')
    return
  }

  loading.value = true
  try {
    const categoryData = {
      name: categoryForm.value.name.trim(),
      color: categoryForm.value.color,
      description: categoryForm.value.description.trim(),
      parent_id: categoryForm.value.parent_id
    }

    let success
    if (editingCategory.value) {
      success = await todoStore.updateCategory(editingCategory.value.id, categoryData)
    } else {
      success = await todoStore.createCategory(categoryData)
    }

    if (success) {
      dialogVisible.value = false
    }
  } finally {
    loading.value = false
  }
}

const deleteCategory = async (category) => {
  const todoCount = getCategoryTodoCount(category.id)
  if (todoCount > 0) {
    ElMessage.error(`该分类下还有 ${todoCount} 个任务，无法删除`)
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除分类"${category.name}"吗？此操作不可撤销。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await todoStore.deleteCategory(category.id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除分类失败:', error)
    }
  }
}

// 生命周期
onMounted(() => {
  todoStore.fetchCategories()
})
</script>

<style scoped>
.category-manager {
  background: white;
  border-radius: 8px;
  padding: 20px;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.category-header h3 {
  margin: 0;
  color: #303133;
}

.category-list {
  min-height: 200px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
}

.empty-state .el-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  margin-bottom: 12px;
  transition: all 0.2s;
}

.category-item:hover {
  border-color: #c0c4cc;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.category-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.category-color {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  flex-shrink: 0;
}

.category-details h4 {
  margin: 0 0 4px 0;
  color: #303133;
  font-size: 16px;
}

.category-details p {
  margin: 0 0 4px 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.4;
}

.category-count {
  color: #909399;
  font-size: 12px;
}

.category-actions {
  display: flex;
  gap: 4px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>