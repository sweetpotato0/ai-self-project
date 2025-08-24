import api from './index'

export const todoApi = {
  // 获取TODO列表
  getTodos() {
    return api.get('/todos')
  },

  // 创建TODO
  createTodo(data) {
    return api.post('/todos', data)
  },

  // 更新TODO
  updateTodo(id, data) {
    return api.put(`/todos/${id}`, data)
  },

  // 删除TODO
  deleteTodo(id) {
    return api.delete(`/todos/${id}`)
  }
}
