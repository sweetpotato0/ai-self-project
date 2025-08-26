import api from '@/api/index'

/**
 * 待办事项模块API接口
 * 包含任务的CRUD操作、分类管理、统计等功能
 */
export const todosApi = {
  /**
   * 获取待办事项列表
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.limit - 每页数量
   * @param {string} params.status - 状态筛选 (pending|in_progress|completed|cancelled)
   * @param {string} params.category - 分类筛选
   * @param {string} params.priority - 优先级筛选 (low|medium|high|urgent)
   * @param {string} params.search - 搜索关键词
   * @param {string} params.sortBy - 排序字段
   * @param {string} params.sortOrder - 排序方向 (asc|desc)
   * @returns {Promise} 待办事项列表和分页信息
   */
  getTodos(params = {}) {
    return api.get('/todos', { params })
  },

  /**
   * 获取单个待办事项详情
   * @param {number|string} todoId - 待办事项ID
   * @returns {Promise} 待办事项详细信息
   */
  getTodo(todoId) {
    return api.get(`/todos/${todoId}`)
  },

  /**
   * 创建新的待办事项
   * @param {Object} todoData - 待办事项数据
   * @param {string} todoData.title - 标题
   * @param {string} todoData.description - 描述
   * @param {string} todoData.priority - 优先级
   * @param {string} todoData.category - 分类
   * @param {string} todoData.dueDate - 截止日期
   * @param {Array} todoData.tags - 标签数组
   * @param {number} todoData.estimatedHours - 预估工时
   * @returns {Promise} 创建的待办事项信息
   */
  createTodo(todoData) {
    return api.post('/todos', {
      title: todoData.title,
      description: todoData.description,
      priority: todoData.priority || 'medium',
      category: todoData.category,
      due_date: todoData.dueDate,
      tags: todoData.tags || [],
      estimated_hours: todoData.estimatedHours || 0,
      status: 'pending'
    })
  },

  /**
   * 更新待办事项
   * @param {number|string} todoId - 待办事项ID
   * @param {Object} todoData - 更新数据
   * @returns {Promise} 更新后的待办事项信息
   */
  updateTodo(todoId, todoData) {
    return api.put(`/todos/${todoId}`, todoData)
  },

  /**
   * 删除待办事项
   * @param {number|string} todoId - 待办事项ID
   * @returns {Promise} 删除结果
   */
  deleteTodo(todoId) {
    return api.delete(`/todos/${todoId}`)
  },

  /**
   * 批量删除待办事项
   * @param {Array} todoIds - 待办事项ID数组
   * @returns {Promise} 批量删除结果
   */
  batchDeleteTodos(todoIds) {
    return api.delete('/todos/batch', {
      data: { ids: todoIds }
    })
  },

  /**
   * 更新待办事项状态
   * @param {number|string} todoId - 待办事项ID
   * @param {string} status - 新状态
   * @returns {Promise} 更新结果
   */
  updateTodoStatus(todoId, status) {
    return api.patch(`/todos/${todoId}/status`, {
      status: status
    })
  },

  /**
   * 批量更新待办事项状态
   * @param {Array} todoIds - 待办事项ID数组
   * @param {string} status - 新状态
   * @returns {Promise} 批量更新结果
   */
  batchUpdateTodoStatus(todoIds, status) {
    return api.patch('/todos/batch/status', {
      ids: todoIds,
      status: status
    })
  },

  /**
   * 为待办事项添加标签
   * @param {number|string} todoId - 待办事项ID
   * @param {Array} tags - 标签数组
   * @returns {Promise} 添加结果
   */
  addTodoTags(todoId, tags) {
    return api.post(`/todos/${todoId}/tags`, {
      tags: tags
    })
  },

  /**
   * 移除待办事项标签
   * @param {number|string} todoId - 待办事项ID
   * @param {Array} tags - 要移除的标签数组
   * @returns {Promise} 移除结果
   */
  removeTodoTags(todoId, tags) {
    return api.delete(`/todos/${todoId}/tags`, {
      data: { tags: tags }
    })
  },

  /**
   * 获取待办事项分类列表
   * @returns {Promise} 分类列表
   */
  getTodoCategories() {
    return api.get('/todos/categories')
  },

  /**
   * 创建新分类
   * @param {Object} categoryData - 分类数据
   * @param {string} categoryData.name - 分类名称
   * @param {string} categoryData.color - 分类颜色
   * @param {string} categoryData.description - 分类描述
   * @returns {Promise} 创建的分类信息
   */
  createTodoCategory(categoryData) {
    return api.post('/todos/categories', categoryData)
  },

  /**
   * 更新分类信息
   * @param {number|string} categoryId - 分类ID
   * @param {Object} categoryData - 更新数据
   * @returns {Promise} 更新结果
   */
  updateTodoCategory(categoryId, categoryData) {
    return api.put(`/todos/categories/${categoryId}`, categoryData)
  },

  /**
   * 删除分类
   * @param {number|string} categoryId - 分类ID
   * @returns {Promise} 删除结果
   */
  deleteTodoCategory(categoryId) {
    return api.delete(`/todos/categories/${categoryId}`)
  },

  /**
   * 获取所有标签
   * @returns {Promise} 标签列表
   */
  getTodoTags() {
    return api.get('/todos/tags')
  },

  /**
   * 获取待办事项统计信息
   * @param {Object} params - 统计参数
   * @param {string} params.period - 统计周期 (week|month|quarter|year)
   * @param {string} params.startDate - 开始日期
   * @param {string} params.endDate - 结束日期
   * @returns {Promise} 统计数据
   */
  getTodoStats(params = {}) {
    return api.get('/todos/stats', { params })
  },

  /**
   * 获取待办事项趋势数据
   * @param {Object} params - 趋势参数
   * @param {number} params.days - 天数
   * @returns {Promise} 趋势数据
   */
  getTodoTrends(params = {}) {
    return api.get('/todos/trends', { params })
  },

  /**
   * 导出待办事项数据
   * @param {Object} params - 导出参数
   * @param {string} params.format - 导出格式 (csv|excel|pdf)
   * @param {Array} params.ids - 指定导出的ID数组
   * @param {Object} params.filters - 筛选条件
   * @returns {Promise} 导出文件信息
   */
  exportTodos(params = {}) {
    return api.get('/todos/export', { 
      params,
      responseType: 'blob'
    })
  },

  /**
   * 导入待办事项数据
   * @param {File} file - 导入文件
   * @param {string} format - 文件格式
   * @returns {Promise} 导入结果
   */
  importTodos(file, format = 'csv') {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('format', format)
    
    return api.post('/todos/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 获取今日待办事项
   * @returns {Promise} 今日待办事项列表
   */
  getTodayTodos() {
    return api.get('/todos/today')
  },

  /**
   * 获取过期的待办事项
   * @returns {Promise} 过期待办事项列表
   */
  getOverdueTodos() {
    return api.get('/todos/overdue')
  },

  /**
   * 获取即将到期的待办事项
   * @param {number} days - 未来几天内 (默认7天)
   * @returns {Promise} 即将到期的待办事项列表
   */
  getUpcomingTodos(days = 7) {
    return api.get('/todos/upcoming', {
      params: { days }
    })
  }
}

export default todosApi