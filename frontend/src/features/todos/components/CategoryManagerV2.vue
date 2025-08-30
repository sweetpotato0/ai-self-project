<template>
  <div class="category-manager-v2">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <h3 class="title">分类管理</h3>
        <el-badge :value="totalCategoriesCount" class="count-badge">
          <span class="count-text">共 {{ totalCategoriesCount }} 个分类</span>
        </el-badge>
      </div>
      <div class="toolbar-right">
        <el-button-group size="small">
          <el-button :type="viewMode === 'tree' ? 'primary' : ''" @click="setViewMode('tree')">
            <el-icon><Grid /></el-icon>
            树状图
          </el-button>
          <el-button :type="viewMode === 'list' ? 'primary' : ''" @click="setViewMode('list')">
            <el-icon><List /></el-icon>
            列表
          </el-button>
        </el-button-group>
        <el-divider direction="vertical" />
        <el-button type="primary" @click="openAddDialog()">
          <el-icon><Plus /></el-icon>
          新建分类
        </el-button>
      </div>
    </div>

    <!-- 搜索过滤栏 -->
    <div class="filter-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索分类名称..."
        clearable
        size="small"
        style="width: 300px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select v-model="filterParent" placeholder="筛选父级" size="small" clearable style="width: 200px">
        <el-option label="顶级分类" :value="0" />
        <el-option 
          v-for="category in parentOptions" 
          :key="category.id"
          :label="category.name"
          :value="category.id"
        />
      </el-select>
      <span v-if="filteredCategories.length < totalCategoriesCount" class="filter-result">
        筛选结果：{{ filteredCategories.length }} 个分类
      </span>
    </div>

    <!-- 主要内容区域 -->
    <div class="content-area" :class="{ 'tree-mode': viewMode === 'tree' }">
      <!-- 空状态 -->
      <div v-if="filteredCategories.length === 0 && !searchKeyword" class="empty-state">
        <div class="empty-icon">
          <el-icon size="64"><FolderOpened /></el-icon>
        </div>
        <h4>暂无分类</h4>
        <p>创建您的第一个分类来开始分组管理任务</p>
        <el-button type="primary" @click="openAddDialog()">
          <el-icon><Plus /></el-icon>
          创建分类
        </el-button>
      </div>

      <!-- 搜索无结果 -->
      <div v-else-if="filteredCategories.length === 0" class="no-results">
        <el-icon size="48"><Search /></el-icon>
        <h4>未找到匹配的分类</h4>
        <p>试试其他关键词或清空搜索条件</p>
      </div>

      <!-- 树状视图 -->
      <div v-else-if="viewMode === 'tree'" class="tree-view">
        <category-tree-node
          v-for="category in topLevelFilteredCategories"
          :key="category.id"
          :category="category"
          :categories="categories"
          :search-keyword="searchKeyword"
          @edit="openEditDialog"
          @delete="handleDelete"
          @add-child="openAddDialog"
          :get-todo-count="getCategoryTodoCount"
        />
      </div>

      <!-- 列表视图 -->
      <div v-else class="list-view">
        <div class="list-header">
          <div class="col-name">名称</div>
          <div class="col-parent">父级分类</div>
          <div class="col-count">任务数量</div>
          <div class="col-updated">更新时间</div>
          <div class="col-actions">操作</div>
        </div>
        <div class="list-body">
          <category-list-item
            v-for="category in filteredCategories"
            :key="category.id"
            :category="category"
            :categories="categories"
            :search-keyword="searchKeyword"
            @edit="openEditDialog"
            @delete="handleDelete"
            @add-child="openAddDialog"
            :get-todo-count="getCategoryTodoCount"
          />
        </div>
      </div>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingCategory ? '编辑分类' : '新建分类'"
      width="480px"
      :close-on-click-modal="false"
    >
      <el-form 
        ref="formRef"
        :model="form" 
        :rules="rules"
        label-width="90px"
        label-position="left"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input 
            v-model="form.name" 
            placeholder="请输入分类名称"
            maxlength="20"
            show-word-limit
            clearable
          />
        </el-form-item>

        <el-form-item label="父级分类">
          <el-cascader
            v-model="parentPath"
            :options="parentCascaderOptions"
            :props="cascaderProps"
            placeholder="选择父级分类（可选）"
            clearable
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="分类图标">
          <category-icon-picker v-model="form.icon" />
        </el-form-item>
        
        <el-form-item label="分类颜色">
          <div class="color-picker-container">
            <el-color-picker 
              v-model="form.color"
              :predefine="predefinedColors"
              show-alpha
            />
            <span class="color-preview" :style="{ backgroundColor: form.color }"></span>
          </div>
        </el-form-item>
        
        <el-form-item label="分类描述">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入分类描述（可选）"
            maxlength="200"
            show-word-limit
            :rows="3"
          />
        </el-form-item>

        <el-form-item label="排序权重">
          <el-input-number
            v-model="form.sort_order"
            :min="0"
            :max="999"
            placeholder="数字越小排序越前"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSave" :loading="saving">
            {{ editingCategory ? '保存修改' : '创建分类' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, Edit, Delete, Search, Grid, List, 
  FolderOpened, ArrowRight, More 
} from '@element-plus/icons-vue'
import { useTodoStore } from '@/stores/todo'
import CategoryTreeNode from './CategoryTreeNode.vue'
import CategoryListItem from './CategoryListItem.vue'
import CategoryIconPicker from './CategoryIconPicker.vue'

const todoStore = useTodoStore()

// 响应式状态
const viewMode = ref('tree') // 'tree' | 'list'
const searchKeyword = ref('')
const filterParent = ref(null)
const dialogVisible = ref(false)
const editingCategory = ref(null)
const saving = ref(false)
const parentPath = ref([])

// 表单数据
const form = ref({
  name: '',
  icon: 'folder',
  color: '#409eff',
  description: '',
  parent_id: null,
  sort_order: 0
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 1, max: 20, message: '名称长度应为 1-20 个字符', trigger: 'blur' }
  ]
}

// 预定义颜色
const predefinedColors = [
  '#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399',
  '#c45656', '#73767a', '#b88230', '#c0c4cc', '#606266',
  '#8957a1', '#ff7875', '#36cfc9', '#ffc069', '#95de64'
]

// 级联选择器配置
const cascaderProps = {
  value: 'id',
  label: 'name',
  children: 'children',
  checkStrictly: true, // 可选择任意级别
  emitPath: false // 只返回最后一级的值
}

// 计算属性
const categories = computed(() => todoStore.categories)
const totalCategoriesCount = computed(() => categories.value.length)

// 过滤后的分类
const filteredCategories = computed(() => {
  let filtered = categories.value

  // 搜索关键词过滤
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(cat => 
      cat.name.toLowerCase().includes(keyword) ||
      (cat.description && cat.description.toLowerCase().includes(keyword))
    )
  }

  // 父级过滤
  if (filterParent.value !== null) {
    if (filterParent.value === 0) {
      // 只显示顶级分类
      filtered = filtered.filter(cat => !cat.parent_id)
    } else {
      // 显示指定父级的子分类
      filtered = filtered.filter(cat => cat.parent_id === filterParent.value)
    }
  }

  return filtered.sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))
})

// 顶级分类（用于树状视图）
const topLevelFilteredCategories = computed(() => {
  return filteredCategories.value.filter(cat => !cat.parent_id)
})

// 父级选项（用于筛选）
const parentOptions = computed(() => {
  return categories.value
    .filter(cat => !cat.parent_id)
    .sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))
})

// 父级分类级联选项
const parentCascaderOptions = computed(() => {
  const buildTree = (items, parentId = null) => {
    return items
      .filter(item => {
        // 过滤掉当前编辑的分类及其子分类
        if (editingCategory.value && 
           (item.id === editingCategory.value.id || isDescendant(item, editingCategory.value.id))) {
          return false
        }
        return item.parent_id === parentId
      })
      .sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))
      .map(item => ({
        id: item.id,
        name: item.name,
        children: buildTree(items, item.id)
      }))
  }
  
  return buildTree(categories.value)
})

// 工具方法
const isDescendant = (category, ancestorId) => {
  let current = category
  while (current.parent_id) {
    if (current.parent_id === ancestorId) return true
    current = categories.value.find(c => c.id === current.parent_id)
    if (!current) break
  }
  return false
}

const getCategoryTodoCount = (categoryId) => {
  return todoStore.todos.filter(todo => todo.category_id === categoryId).length
}

const setViewMode = (mode) => {
  viewMode.value = mode
}

// 对话框操作
const openAddDialog = (parentCategory = null) => {
  editingCategory.value = null
  form.value = {
    name: '',
    icon: 'folder',
    color: '#409eff',
    description: '',
    parent_id: parentCategory?.id || null,
    sort_order: 0
  }
  
  if (parentCategory) {
    // 构建父级路径
    parentPath.value = [parentCategory.id]
  } else {
    parentPath.value = []
  }
  
  dialogVisible.value = true
}

const openEditDialog = (category) => {
  editingCategory.value = category
  form.value = {
    name: category.name,
    icon: category.icon || 'folder',
    color: category.color || '#409eff',
    description: category.description || '',
    parent_id: category.parent_id || null,
    sort_order: category.sort_order || 0
  }
  
  // 构建父级路径
  if (category.parent_id) {
    const buildParentPath = (catId) => {
      const cat = categories.value.find(c => c.id === catId)
      if (!cat) return []
      if (cat.parent_id) {
        return [...buildParentPath(cat.parent_id), cat.id]
      }
      return [cat.id]
    }
    parentPath.value = buildParentPath(category.parent_id)
  } else {
    parentPath.value = []
  }
  
  dialogVisible.value = true
}

// 监听级联选择器变化
const handleParentChange = () => {
  form.value.parent_id = parentPath.value.length > 0 ? 
    parentPath.value[parentPath.value.length - 1] : null
}

// 保存分类
const formRef = ref()

const handleSave = async () => {
  try {
    if (!formRef.value) return
    
    await formRef.value.validate()
    
    saving.value = true
    
    // 获取父级ID（级联选择器的最后一级）
    form.value.parent_id = parentPath.value.length > 0 ? 
      parentPath.value[parentPath.value.length - 1] : null
    
    const categoryData = { ...form.value }
    
    let success
    if (editingCategory.value) {
      success = await todoStore.updateCategory(editingCategory.value.id, categoryData)
    } else {
      success = await todoStore.createCategory(categoryData)
    }
    
    if (success) {
      dialogVisible.value = false
      ElMessage.success(editingCategory.value ? '分类更新成功' : '分类创建成功')
    }
  } catch (error) {
    console.error('保存分类失败:', error)
  } finally {
    saving.value = false
  }
}

// 删除分类
const handleDelete = async (category) => {
  const todoCount = getCategoryTodoCount(category.id)
  const childCount = categories.value.filter(c => c.parent_id === category.id).length
  
  let message = `确定要删除分类"${category.name}"吗？`
  if (todoCount > 0) {
    message += `\n该分类下有 ${todoCount} 个任务。`
  }
  if (childCount > 0) {
    message += `\n该分类下有 ${childCount} 个子分类。`
  }
  message += '\n此操作不可撤销。'
  
  try {
    await ElMessageBox.confirm(message, '确认删除', {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
      dangerouslyUseHTMLString: true
    })
    
    const success = await todoStore.deleteCategory(category.id)
    if (success) {
      ElMessage.success('分类删除成功')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除分类失败:', error)
    }
  }
}

// 生命周期
onMounted(async () => {
  await todoStore.fetchCategories()
})
</script>

<style scoped>
.category-manager-v2 {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

/* 工具栏样式 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.count-badge {
  --el-badge-bg-color: transparent;
}

.count-text {
  font-size: 13px;
  color: #909399;
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 12px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 过滤栏样式 */
.filter-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: #fafbfc;
  border-radius: 8px;
}

.filter-result {
  font-size: 13px;
  color: #409eff;
  margin-left: auto;
}

/* 内容区域样式 */
.content-area {
  flex: 1;
  min-height: 400px;
  overflow: hidden;
}

/* 空状态样式 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 400px;
  color: #909399;
}

.empty-icon {
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-state h4 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #606266;
}

.empty-state p {
  margin: 0 0 24px 0;
  font-size: 14px;
}

.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: #909399;
}

.no-results h4 {
  margin: 16px 0 8px 0;
  color: #606266;
}

.no-results p {
  margin: 0;
  font-size: 14px;
}

/* 树状视图样式 */
.tree-view {
  height: 100%;
  overflow-y: auto;
  padding: 8px 0;
}

/* 列表视图样式 */
.list-view {
  height: 100%;
  display: flex;
  flex-direction: column;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
}

.list-header {
  display: grid;
  grid-template-columns: 2fr 1.5fr 100px 120px 120px;
  gap: 16px;
  align-items: center;
  padding: 12px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
  font-weight: 600;
  font-size: 13px;
  color: #606266;
}

.list-body {
  flex: 1;
  overflow-y: auto;
}

/* 对话框样式 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.color-picker-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.color-preview {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .category-manager-v2 {
    padding: 16px;
  }
  
  .toolbar {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .toolbar-right {
    justify-content: space-between;
  }
  
  .filter-bar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .list-header {
    grid-template-columns: 1fr 80px 100px;
  }
  
  .col-parent,
  .col-updated {
    display: none;
  }
}
</style>