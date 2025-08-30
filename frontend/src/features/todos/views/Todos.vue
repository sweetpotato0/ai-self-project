<template>
  <div class="todos-container">
    <div class="todos-header">
      <h2>TODO清单管理</h2>
      <el-button type="primary" @click="showDialog = true">
        <el-icon><Plus /></el-icon>
        新建任务
      </el-button>
    </div>

    <el-tabs v-model="activeTab" class="todos-tabs">
      <el-tab-pane label="任务列表" name="todos">
        <div class="todos-content">
      <!-- 筛选和搜索 -->
      <div class="filter-section">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-select v-model="filters.status" placeholder="状态筛选" clearable>
              <el-option label="全部" value="" />
              <el-option label="待处理" value="pending" />
              <el-option label="进行中" value="in_progress" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
          </el-col>
          <el-col :span="6">
            <el-select v-model="filters.priority" placeholder="优先级筛选" clearable>
              <el-option label="全部" value="" />
              <el-option label="低" :value="1" />
              <el-option label="中" :value="2" />
              <el-option label="高" :value="3" />
              <el-option label="紧急" :value="4" />
              <el-option label="立即" :value="5" />
            </el-select>
          </el-col>
          <el-col :span="6">
            <category-cascader
              v-model="filters.category"
              :categories="todoStore.categories"
              placeholder="分类筛选"
              :show-counts="true"
              :get-todo-count="getCategoryTodoCount"
              size="default"
            />
          </el-col>
          <el-col :span="6">
            <el-input
              v-model="filters.search"
              placeholder="搜索任务..."
              clearable
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-col>
        </el-row>
      </div>

      <!-- 任务列表 -->
      <el-table
        :data="filteredTodos"
        v-loading="todoStore.loading"
        class="todos-table"
      >
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <div class="todo-title">
              <span :class="{ 'completed': row.status === 'completed' }">
                {{ row.title }}
              </span>
              <el-tag
                v-if="row.status === 'completed'"
                type="success"
                size="small"
              >
                已完成
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />

        <el-table-column prop="priority_id" label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority_id)" size="small">
              {{ getPriorityText(row.priority_id) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="category_id" label="分类" width="100">
          <template #default="{ row }">
            <el-tag v-if="getCategoryName(row.category_id)" type="info" size="small">
              {{ getCategoryName(row.category_id) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="时间信息" width="200">
          <template #default="{ row }">
            <div class="time-info">
              <div v-if="row.start_date" class="time-item">
                <span class="time-label">开始:</span>
                <span>{{ formatDateTimeLocal(row.start_date) }}</span>
              </div>
              <div v-if="row.due_date" class="time-item">
                <span class="time-label">截止:</span>
                <span :class="{ 'overdue': isOverdue(row.due_date) }">
                  {{ formatDateTimeLocal(row.due_date) }}
                </span>
              </div>
              <div v-if="row.completed_at" class="time-item">
                <span class="time-label">完成:</span>
                <span>{{ formatDateTimeLocal(row.completed_at) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="工时" width="120">
          <template #default="{ row }">
            <div class="hours-info">
              <div v-if="row.estimated_hours" class="hours-item">
                <span class="hours-label">预估:</span>
                <span>{{ row.estimated_hours }}h</span>
              </div>
              <div v-if="row.actual_hours" class="hours-item">
                <span class="hours-label">实际:</span>
                <span>{{ row.actual_hours }}h</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="editTodo(row)">编辑</el-button>
            <el-button
              v-if="row.status !== 'completed'"
              size="small"
              type="success"
              @click="completeTodo(row)"
            >
              完成
            </el-button>
            <el-button size="small" type="danger" @click="deleteTodo(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane label="分类管理" name="categories">
        <CategoryManagerV2 />
      </el-tab-pane>
    </el-tabs>

    <!-- 新建/编辑任务对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="editingTodo ? '编辑任务' : '新建任务'"
      width="600px"
    >
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入任务标题" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入任务描述"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
                    <el-form-item label="优先级" prop="priority_id">
          <el-select v-model="form.priority_id" placeholder="选择优先级">
                <el-option label="低" :value="1" />
                <el-option label="中" :value="2" />
                <el-option label="高" :value="3" />
                <el-option label="紧急" :value="4" />
                <el-option label="立即" :value="5" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分类">
              <category-cascader
                v-model="form.category_id"
                :categories="todoStore.categories"
                placeholder="选择分类"
                :show-counts="true"
                :get-todo-count="getCategoryTodoCount"
                size="default"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="全天任务">
              <el-switch v-model="form.isAllDay" />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 全天任务：只需要一个日期 -->
        <el-row :gutter="20" v-if="form.isAllDay">
          <el-col :span="12">
            <el-form-item label="任务日期">
              <el-date-picker
                v-model="form.start_date"
                type="date"
                placeholder="选择任务日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 非全天任务：开始时间，结束时间 -->
        <template v-if="!form.isAllDay">
          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item label="开始时间">
                <div class="datetime-group">
                  <el-date-picker
                    v-model="form.start_date"
                    type="date"
                    placeholder="开始日期"
                    format="YYYY-MM-DD"
                    value-format="YYYY-MM-DD"
                    style="width: 200px"
                  />
                  <el-time-picker
                    v-model="form.start_time"
                    placeholder="开始时间"
                    format="HH:mm"
                    value-format="HH:mm"
                    style="width: 120px"
                  />
                </div>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item label="结束时间">
                <div class="datetime-group">
                  <el-date-picker
                    v-model="form.due_date"
                    type="date"
                    placeholder="结束日期"
                    format="YYYY-MM-DD"
                    value-format="YYYY-MM-DD"
                    style="width: 200px"
                  />
                  <el-time-picker
                    v-model="form.due_time"
                    placeholder="结束时间"
                    format="HH:mm"
                    value-format="HH:mm"
                    style="width: 120px"
                  />
                </div>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item label="时间预设">
                <div class="time-presets">
                  <el-button
                    v-for="preset in timePresets"
                    :key="preset.label"
                    size="small"
                    @click="setTimePreset(preset)"
                    :type="isTimePresetActive(preset) ? 'primary' : 'default'"
                  >
                    {{ preset.label }}
                  </el-button>
                </div>
              </el-form-item>
            </el-col>
          </el-row>
        </template>

        <el-row :gutter="20">
          <el-col :span="12">
                        <el-form-item label="预估工时">
              <el-input-number
                v-model="form.estimated_hours"
                :min="0"
                :max="24"
                :precision="1"
                placeholder="预估工时（小时）"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
                        <el-form-item label="实际工时">
              <el-input-number
                v-model="form.actual_hours"
                :min="0"
                :max="24"
                :precision="1"
                placeholder="实际工时（小时）"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveTodo" :loading="saving">
          {{ editingTodo ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useTodoStore } from '@/stores/todo'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { formatDateTime, isOverdue } from '@/features/tools/shared/utils/dateTime'
import CategoryManagerV2 from '../components/CategoryManagerV2.vue'
import CategoryCascader from '@/components/common/CategoryCascader.vue'

const authStore = useAuthStore()
const todoStore = useTodoStore()

// 响应式数据
const activeTab = ref('todos')
const showDialog = ref(false)
const editingTodo = ref(null)
const saving = ref(false)
const formRef = ref()

const filters = reactive({
  status: '',
  priority: '',
  category: '',
  search: ''
})

// 时间预设
const timePresets = [
  { label: '上午', start: '09:00', end: '12:00' },
  { label: '下午', start: '14:00', end: '18:00' },
  { label: '晚上', start: '19:00', end: '22:00' },
  { label: '全天', start: '00:00', end: '23:59' },
  { label: '1小时', start: '09:00', end: '10:00' },
  { label: '2小时', start: '09:00', end: '11:00' },
  { label: '4小时', start: '09:00', end: '13:00' }
]

const form = reactive({
  title: '',
  description: '',
  priority_id: 2,
  category_id: null,
  isAllDay: false,
  start_date: '',
  due_date: '',
  start_time: '09:00',
  due_time: '10:00',
  estimated_hours: 0,
  actual_hours: 0
})

const formRules = {
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' }
  ],
  priority_id: [
    { required: true, message: '请选择优先级', trigger: 'change' }
  ],
  due_date: [
    { required: true, message: '请选择截止时间', trigger: 'change' }
  ]
}

// 计算属性
const filteredTodos = computed(() => {
  let todos = todoStore.todos || []

  if (filters.status) {
    todos = todos.filter(todo => todo.status === filters.status)
  }

    if (filters.priority) {
    todos = todos.filter(todo => todo.priority_id === filters.priority)
  }

  if (filters.category) {
    todos = todos.filter(todo => todo.category_id === filters.category)
  }

  if (filters.search) {
    const search = filters.search.toLowerCase()
    todos = todos.filter(todo =>
      todo.title.toLowerCase().includes(search) ||
      todo.description?.toLowerCase().includes(search)
    )
  }

  return todos
})

// 方法
const getPriorityType = (priority) => {
  const typeMap = { 1: 'success', 2: 'warning', 3: 'danger', 4: 'danger', 5: 'danger' }
  return typeMap[priority] || 'info'
}

const getPriorityText = (priority) => {
  const textMap = { 1: '低', 2: '中', 3: '高', 4: '紧急', 5: '立即' }
  return textMap[priority] || '未知'
}

const getStatusType = (status) => {
  const typeMap = {
    'pending': 'info',
    'in_progress': 'warning',
    'completed': 'success',
    'cancelled': 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    'pending': '待处理',
    'in_progress': '进行中',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return textMap[status] || '未知'
}

const formatDateTimeLocal = formatDateTime

// 获取分类任务数量
const getCategoryTodoCount = (categoryId) => {
  return todoStore.todos.filter(todo => todo.category_id === categoryId).length
}

// 获取分类名称
const getCategoryName = (categoryId) => {
  if (!categoryId) return null
  const category = todoStore.categories.find(cat => cat.id === categoryId)
  return category ? category.name : null
}

// 设置时间预设
const setTimePreset = (preset) => {
  form.start_time = preset.start
  form.due_time = preset.end
}

// 检查时间预设是否激活
const isTimePresetActive = (preset) => {
  return form.start_time === preset.start && form.due_time === preset.end
}

const editTodo = (todo) => {
  editingTodo.value = todo
  form.title = todo.title || ''
  form.description = todo.description || ''
  form.priority_id = todo.priority_id || 2
  form.category_id = todo.category_id || null
      // 处理时间字段
    if (todo.start_date && todo.due_date) {
      console.log('原始时间数据:', {
        start_date: todo.start_date,
        due_date: todo.due_date
      })
      
      // 将UTC时间转换为本地时间进行编辑
      const startDate = new Date(todo.start_date)
      const endDate = new Date(todo.due_date)
      
      // 获取本地时间的各个组件
      const startHour = startDate.getHours()
      const startMinute = startDate.getMinutes()
      const endHour = endDate.getHours()
      const endMinute = endDate.getMinutes()
      
      console.log('本地时间解析后:', {
        startTime: `${startHour.toString().padStart(2, '0')}:${startMinute.toString().padStart(2, '0')}`,
        endTime: `${endHour.toString().padStart(2, '0')}:${endMinute.toString().padStart(2, '0')}`
      })
      
      // 检查是否为全天任务（UTC时间的00:00-23:59在本地时间可能会有时区偏移）
      // 但我们检查时间差是否接近24小时来判断是否为全天任务
      const timeDiffMs = endDate.getTime() - startDate.getTime()
      const timeDiffHours = timeDiffMs / (1000 * 60 * 60)
      
      if (Math.abs(timeDiffHours - 24) < 0.1) { // 允许小的误差
        // 全天任务
        form.isAllDay = true
        const year = startDate.getFullYear()
        const month = (startDate.getMonth() + 1).toString().padStart(2, '0')
        const day = startDate.getDate().toString().padStart(2, '0')
        form.start_date = `${year}-${month}-${day}`
        form.due_date = ''
        form.start_time = '09:00'
        form.due_time = '10:00'
      } else {
        // 非全天任务
        form.isAllDay = false
        const startYear = startDate.getFullYear()
        const startMonth = (startDate.getMonth() + 1).toString().padStart(2, '0')
        const startDay = startDate.getDate().toString().padStart(2, '0')
        form.start_date = `${startYear}-${startMonth}-${startDay}`
        
        const endYear = endDate.getFullYear()
        const endMonth = (endDate.getMonth() + 1).toString().padStart(2, '0')
        const endDay = endDate.getDate().toString().padStart(2, '0')
        form.due_date = `${endYear}-${endMonth}-${endDay}`
        
        form.start_time = `${startHour.toString().padStart(2, '0')}:${startMinute.toString().padStart(2, '0')}`
        form.due_time = `${endHour.toString().padStart(2, '0')}:${endMinute.toString().padStart(2, '0')}`
      }
    } else {
      form.start_date = ''
      form.due_date = ''
      form.isAllDay = false
      form.start_time = '09:00'
      form.due_time = '10:00'
    }
  form.estimated_hours = todo.estimated_hours || 0
  form.actual_hours = todo.actual_hours || 0
  showDialog.value = true
}

const completeTodo = async (todo) => {
  try {
    await ElMessageBox.confirm(
      '确定要将此任务标记为已完成吗？',
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await todoStore.updateTodo(todo.id, {
      ...todo,
      status: 'completed',
      completedAt: new Date().toISOString()
    })

    ElMessage.success('任务已完成')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const deleteTodo = async (todo) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除此任务吗？此操作不可恢复。',
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await todoStore.deleteTodo(todo.id)
    ElMessage.success('任务已删除')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const saveTodo = async () => {
  try {
    await formRef.value.validate()
    saving.value = true

    const todoData = {
      title: form.title,
      description: form.description,
      priority_id: form.priority_id,
      category_id: form.category_id,
      estimated_hours: form.estimated_hours,
      actual_hours: form.actual_hours
    }

    // 只有当日期不为空时才添加到请求中
    // 处理时间格式
    if (form.isAllDay) {
      // 全天任务
      if (form.start_date) {
        const [year, month, day] = form.start_date.split('-').map(Number)
        const startDateTime = new Date(year, month - 1, day, 0, 0, 0)
        const endDateTime = new Date(year, month - 1, day, 23, 59, 59)
        todoData.start_date = startDateTime.toISOString()
        todoData.due_date = endDateTime.toISOString()
      }
    } else {
      // 非全天任务
      if (form.start_date) {
        const [year, month, day] = form.start_date.split('-').map(Number)
        const [startHour, startMinute] = form.start_time.split(':').map(Number)
        const startDateTime = new Date(year, month - 1, day, startHour, startMinute, 0)
        todoData.start_date = startDateTime.toISOString()
      }

      if (form.due_date) {
        const [year, month, day] = form.due_date.split('-').map(Number)
        const [endHour, endMinute] = form.due_time.split(':').map(Number)
        const endDateTime = new Date(year, month - 1, day, endHour, endMinute, 0)
        todoData.due_date = endDateTime.toISOString()
      } else if (form.start_date) {
        // 如果没有设置结束日期，默认使用开始日期
        const [year, month, day] = form.start_date.split('-').map(Number)
        const [endHour, endMinute] = form.due_time.split(':').map(Number)
        const endDateTime = new Date(year, month - 1, day, endHour, endMinute, 0)
        todoData.due_date = endDateTime.toISOString()
      }
    }

    if (editingTodo.value) {
      await todoStore.updateTodo(editingTodo.value.id, todoData)
      ElMessage.success('任务更新成功')
    } else {
      await todoStore.createTodo(todoData)
      ElMessage.success('任务创建成功')
    }

    showDialog.value = false
    resetForm()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    saving.value = false
  }
}

const resetForm = () => {
  editingTodo.value = null
  form.title = ''
  form.description = ''
  form.priority_id = 2
  form.category_id = null
  form.start_date = ''
  form.due_date = ''
  form.estimated_hours = 0
  form.actual_hours = 0
}

// 生命周期
onMounted(async () => {
  await Promise.all([
    todoStore.fetchTodos(),
    todoStore.fetchCategories()
  ])
})
</script>

<style scoped>
.todos-container {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.todos-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.todos-header h2 {
  margin: 0;
  color: #303133;
}

.todos-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.todos-tabs :deep(.el-tabs__content) {
  flex: 1;
}

.todos-tabs :deep(.el-tab-pane) {
  height: 100%;
}

.todos-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.filter-section {
  margin-bottom: 20px;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.todos-table {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.todo-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.todo-title .completed {
  text-decoration: line-through;
  color: #909399;
}

.time-info {
  font-size: 12px;
}

.time-item {
  margin-bottom: 2px;
}

.time-label {
  color: #909399;
  margin-right: 4px;
}

.overdue {
  color: #f56c6c;
  font-weight: 600;
}

.hours-info {
  font-size: 12px;
}

.hours-item {
  margin-bottom: 2px;
}

.hours-label {
  color: #909399;
  margin-right: 4px;
}
</style>
