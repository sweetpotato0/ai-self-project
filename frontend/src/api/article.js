import api from './index'

export const articleApi = {
  // 创建文章
  createArticle(data) {
    return api.post('/articles', data)
  },

  // 获取用户文章列表
  getArticles(params) {
    return api.get('/articles', { params })
  },

  // 获取文章统计信息
  getArticleStats() {
    return api.get('/articles/stats')
  },

  // 根据ID获取文章
  getArticleById(id) {
    return api.get(`/articles/${id}`)
  },

  // 更新文章
  updateArticle(id, data) {
    return api.put(`/articles/${id}`, data)
  },

  // 删除文章
  deleteArticle(id) {
    return api.delete(`/articles/${id}`)
  }
}
