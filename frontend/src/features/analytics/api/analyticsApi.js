import api from '@/api/index'

/**
 * 全局统计分析 API
 * 跨模块的统计数据和趋势分析
 */
export const analyticsApi = {
  /**
   * 获取全局统计数据汇总
   * @param {Object} params - 查询参数
   * @param {string} params.period - 时间周期: day, week, month, year
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @returns {Promise} API响应
   */
  getOverallStats(params = {}) {
    return api.get('/analytics/overview', { params })
  },

  /**
   * 获取全局趋势数据
   * @param {Object} params - 查询参数
   * @param {string} params.metric - 指标类型: tasks, articles, users
   * @param {string} params.period - 时间周期
   * @param {number} params.days - 查询天数
   * @returns {Promise} API响应
   */
  getTrends(params = {}) {
    return api.get('/analytics/trends', { params })
  },

  /**
   * 获取模块间对比统计
   * @param {Object} params - 查询参数
   * @param {Array} params.modules - 模块列表: ['todos', 'articles']
   * @param {string} params.metric - 对比指标
   * @param {string} params.period - 时间周期
   * @returns {Promise} API响应
   */
  getModuleComparison(params = {}) {
    return api.get('/analytics/comparison', { params })
  },

  /**
   * 获取用户活跃度统计
   * @param {Object} params - 查询参数
   * @param {string} params.period - 时间周期
   * @param {string} params.group_by - 分组方式: hour, day, week
   * @returns {Promise} API响应
   */
  getUserActivity(params = {}) {
    return api.get('/analytics/activity', { params })
  },

  /**
   * 获取性能统计数据
   * @param {Object} params - 查询参数
   * @param {string} params.metric - 性能指标: response_time, error_rate
   * @param {number} params.limit - 限制数量
   * @returns {Promise} API响应
   */
  getPerformanceStats(params = {}) {
    return api.get('/analytics/performance', { params })
  },

  /**
   * 导出统计报告
   * @param {Object} params - 导出参数
   * @param {string} params.format - 导出格式: pdf, excel, csv
   * @param {string} params.report_type - 报告类型: summary, detailed
   * @param {string} params.period - 统计周期
   * @returns {Promise} API响应
   */
  exportReport(params = {}) {
    return api.post('/analytics/export', params, {
      responseType: 'blob'
    })
  }
}