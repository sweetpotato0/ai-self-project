import { todosApi } from '@/features/todos/api'
import { ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { notificationManager } from '@/features/tools/shared/utils/notificationManager'

export const useTodoStore = defineStore('todo', () => {
  const todos = ref([])
  const categories = ref([])
  const loading = ref(false)

  // 获取TODO列表
  const fetchTodos = async () => {
    loading.value = true
    try {
      const response = await todosApi.getTodos()
      todos.value = response.data.todos || []
      
      // 检查任务到期情况并发送通知
      notificationManager.applyNotificationSettings(todos.value)
    } catch (error) {
      ElMessage.error('获取TODO列表失败')
    } finally {
      loading.value = false
    }
  }

  // 创建TODO
  const createTodo = async (todoData) => {
    try {
      const response = await todosApi.createTodo(todoData)
      const newTodo = response.data.todo
      todos.value.unshift(newTodo)
      
      // 发送新任务通知
      notificationManager.showNewTaskNotification(newTodo)
      
      ElMessage.success('TODO创建成功')
      return true
    } catch (error) {
      ElMessage.error('创建TODO失败')
      return false
    }
  }

  // 更新TODO
  const updateTodo = async (id, todoData) => {
    try {
      const response = await todosApi.updateTodo(id, todoData)
      const updatedTodo = response.data.todo
      const index = todos.value.findIndex(todo => todo.id === id)
      const oldTodo = index !== -1 ? todos.value[index] : null
      
      if (index !== -1) {
        todos.value[index] = updatedTodo
      }
      
      // 如果任务被标记为完成，发送完成通知
      if (oldTodo && !oldTodo.completed && updatedTodo.completed) {
        notificationManager.showCompletionNotification(updatedTodo)
      }
      
      ElMessage.success('TODO更新成功')
      return true
    } catch (error) {
      ElMessage.error('更新TODO失败')
      return false
    }
  }

  // 删除TODO
  const deleteTodo = async (id) => {
    try {
      await todosApi.deleteTodo(id)
      const index = todos.value.findIndex(todo => todo.id === id)
      if (index !== -1) {
        todos.value.splice(index, 1)
      }
      ElMessage.success('TODO删除成功')
      return true
    } catch (error) {
      ElMessage.error('删除TODO失败')
      return false
    }
  }

  // 切换TODO状态
  const toggleTodoStatus = async (id, status) => {
    return await updateTodo(id, { status })
  }

  // 获取统计信息
  const getStats = () => {
    const stats = {
      total: todos.value.length,
      pending: 0,
      inProgress: 0,
      completed: 0,
      cancelled: 0,
      overdue: 0
    }

    const now = new Date()
    todos.value.forEach(todo => {
      switch (todo.status) {
        case 'pending':
          stats.pending++
          break
        case 'in_progress':
          stats.inProgress++
          break
        case 'completed':
          stats.completed++
          break
        case 'cancelled':
          stats.cancelled++
          break
      }

      // 检查是否逾期
      if (todo.due_date && new Date(todo.due_date) < now &&
          !['completed', 'cancelled'].includes(todo.status)) {
        stats.overdue++
      }
    })

    return stats
  }

  // 获取分类列表
  const fetchCategories = async () => {
    try {
      const response = await todosApi.getTodoCategories()
      categories.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch categories:', error)
      ElMessage.error('获取分类列表失败')
    }
  }

  // 创建分类
  const createCategory = async (categoryData) => {
    try {
      const response = await todosApi.createTodoCategory(categoryData)
      const newCategory = response.data.category
      categories.value.push(newCategory)
      ElMessage.success('分类创建成功')
      return true
    } catch (error) {
      ElMessage.error('创建分类失败')
      return false
    }
  }

  // 更新分类
  const updateCategory = async (id, categoryData) => {
    try {
      const response = await todosApi.updateTodoCategory(id, categoryData)
      const updatedCategory = response.data.category
      const index = categories.value.findIndex(cat => cat.id === id)
      if (index !== -1) {
        categories.value[index] = updatedCategory
      }
      ElMessage.success('分类更新成功')
      return true
    } catch (error) {
      ElMessage.error('更新分类失败')
      return false
    }
  }

  // 删除分类
  const deleteCategory = async (id) => {
    try {
      await todosApi.deleteTodoCategory(id)
      const index = categories.value.findIndex(cat => cat.id === id)
      if (index !== -1) {
        categories.value.splice(index, 1)
      }
      ElMessage.success('分类删除成功')
      return true
    } catch (error) {
      ElMessage.error('删除分类失败')
      return false
    }
  }

  return {
    todos,
    categories,
    loading,
    fetchTodos,
    createTodo,
    updateTodo,
    deleteTodo,
    toggleTodoStatus,
    getStats,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory
  }
})
