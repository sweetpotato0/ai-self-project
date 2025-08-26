/**
 * 待办事项模块API导出文件
 * 统一管理待办事项相关的所有API接口
 */

import todosApi from './todosApi'

// 导出待办事项API
export { todosApi }

// 默认导出所有API
export default {
  todosApi
}

/**
 * API使用示例：
 * 
 * 在Store中使用：
 * import { todosApi } from '@/features/todos/api'
 * 
 * export const useTodosStore = defineStore('todos', () => {
 *   const todos = ref([])
 *   const loading = ref(false)
 *   
 *   const fetchTodos = async (params = {}) => {
 *     loading.value = true
 *     try {
 *       const response = await todosApi.getTodos(params)
 *       todos.value = response.data.todos
 *       return response
 *     } catch (error) {
 *       ElMessage.error('获取待办事项失败')
 *       throw error
 *     } finally {
 *       loading.value = false
 *     }
 *   }
 *   
 *   const createTodo = async (todoData) => {
 *     try {
 *       const response = await todosApi.createTodo(todoData)
 *       todos.value.unshift(response.data.todo)
 *       ElMessage.success('创建成功')
 *       return response
 *     } catch (error) {
 *       ElMessage.error('创建失败')
 *       throw error
 *     }
 *   }
 *   
 *   return {
 *     todos,
 *     loading,
 *     fetchTodos,
 *     createTodo
 *   }
 * })
 * 
 * 在组件中使用：
 * import { todosApi } from '@/features/todos/api'
 * import { ElMessage } from 'element-plus'
 * 
 * const handleStatusUpdate = async (todoId, newStatus) => {
 *   try {
 *     await todosApi.updateTodoStatus(todoId, newStatus)
 *     ElMessage.success('状态更新成功')
 *     // 刷新数据
 *     await fetchTodos()
 *   } catch (error) {
 *     ElMessage.error('状态更新失败')
 *   }
 * }
 */