import api from './index'

export const statisticsApi = {
  // 获取统一统计数据
  getStatistics(params = {}) {
    return api.get('/statistics', { params })
  },

  // 获取任务统计
  getTodoStats(params = {}) {
    return api.get('/statistics', { params: { type: 'todo', ...params } })
  },

  // 获取文章统计
  getArticleStats(params = {}) {
    return api.get('/statistics', { params: { type: 'article', ...params } })
  },

  // 获取趋势数据
  getTrends(params = {}) {
    return api.get('/statistics/trends', { params })
  },

  // 获取任务完成趋势
  getTodoTrends(params = {}) {
    return api.get('/statistics/trends', { params: { type: 'todo', ...params } })
  },

  // 获取文章浏览量趋势
  getArticleTrends(params = {}) {
    return api.get('/statistics/trends', { params: { type: 'article', ...params } })
  }
}
