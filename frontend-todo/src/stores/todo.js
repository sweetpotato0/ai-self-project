import { todoApi } from '@/api/todo'
import { ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useTodoStore = defineStore('todo', () => {
  const todos = ref([])
  const loading = ref(false)

  // 获取TODO列表
  const fetchTodos = async () => {
    loading.value = true
    try {
      const response = await todoApi.getTodos()
      todos.value = response.data.todos || []
    } catch (error) {
      ElMessage.error('获取TODO列表失败')
    } finally {
      loading.value = false
    }
  }

  // 创建TODO
  const createTodo = async (todoData) => {
    try {
      const response = await todoApi.createTodo(todoData)
      todos.value.unshift(response.data.todo)
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
      const response = await todoApi.updateTodo(id, todoData)
      const index = todos.value.findIndex(todo => todo.id === id)
      if (index !== -1) {
        todos.value[index] = response.data.todo
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
      await todoApi.deleteTodo(id)
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

  return {
    todos,
    loading,
    fetchTodos,
    createTodo,
    updateTodo,
    deleteTodo,
    toggleTodoStatus,
    getStats
  }
})
