import api from '@/api/index'

/**
 * 文章模块API接口
 * 包含文章的CRUD操作、分类管理、统计分析等功能
 */
export const articlesApi = {
  /**
   * 获取文章列表
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.limit - 每页数量
   * @param {string} params.status - 状态筛选 (draft|published|archived)
   * @param {string} params.category - 分类筛选
   * @param {string} params.tag - 标签筛选
   * @param {string} params.search - 搜索关键词
   * @param {string} params.sortBy - 排序字段 (created_at|updated_at|view_count|like_count)
   * @param {string} params.sortOrder - 排序方向 (asc|desc)
   * @param {string} params.author - 作者筛选
   * @returns {Promise} 文章列表和分页信息
   */
  getArticles(params = {}) {
    return api.get('/articles', { params })
  },

  /**
   * 获取单个文章详情
   * @param {number|string} articleId - 文章ID
   * @param {boolean} includeContent - 是否包含文章内容
   * @returns {Promise} 文章详细信息
   */
  getArticle(articleId, includeContent = true) {
    return api.get(`/articles/${articleId}`, {
      params: { include_content: includeContent }
    })
  },

  /**
   * 创建新文章
   * @param {Object} articleData - 文章数据
   * @param {string} articleData.title - 标题
   * @param {string} articleData.content - 内容
   * @param {string} articleData.excerpt - 摘要
   * @param {string} articleData.status - 状态 (draft|published)
   * @param {string} articleData.category - 分类
   * @param {Array} articleData.tags - 标签数组
   * @param {string} articleData.coverImage - 封面图片
   * @param {boolean} articleData.allowComments - 是否允许评论
   * @param {string} articleData.publishAt - 定时发布时间
   * @returns {Promise} 创建的文章信息
   */
  createArticle(articleData) {
    return api.post('/articles', {
      title: articleData.title,
      content: articleData.content,
      excerpt: articleData.excerpt || '',
      status: articleData.status || 'draft',
      category_id: articleData.category,
      tags: articleData.tags || [],
      cover_image: articleData.coverImage,
      allow_comments: articleData.allowComments !== false,
      publish_at: articleData.publishAt,
      meta_title: articleData.metaTitle,
      meta_description: articleData.metaDescription,
      meta_keywords: articleData.metaKeywords
    })
  },

  /**
   * 更新文章
   * @param {number|string} articleId - 文章ID
   * @param {Object} articleData - 更新数据
   * @returns {Promise} 更新后的文章信息
   */
  updateArticle(articleId, articleData) {
    return api.put(`/articles/${articleId}`, articleData)
  },

  /**
   * 删除文章
   * @param {number|string} articleId - 文章ID
   * @returns {Promise} 删除结果
   */
  deleteArticle(articleId) {
    return api.delete(`/articles/${articleId}`)
  },

  /**
   * 批量删除文章
   * @param {Array} articleIds - 文章ID数组
   * @returns {Promise} 批量删除结果
   */
  batchDeleteArticles(articleIds) {
    return api.delete('/articles/batch', {
      data: { ids: articleIds }
    })
  },

  /**
   * 更新文章状态
   * @param {number|string} articleId - 文章ID
   * @param {string} status - 新状态 (draft|published|archived)
   * @returns {Promise} 更新结果
   */
  updateArticleStatus(articleId, status) {
    return api.patch(`/articles/${articleId}/status`, {
      status: status
    })
  },

  /**
   * 批量更新文章状态
   * @param {Array} articleIds - 文章ID数组
   * @param {string} status - 新状态
   * @returns {Promise} 批量更新结果
   */
  batchUpdateArticleStatus(articleIds, status) {
    return api.patch('/articles/batch/status', {
      ids: articleIds,
      status: status
    })
  },

  /**
   * 增加文章浏览量
   * @param {number|string} articleId - 文章ID
   * @returns {Promise} 更新结果
   */
  incrementViewCount(articleId) {
    return api.post(`/articles/${articleId}/view`)
  },

  /**
   * 点赞文章
   * @param {number|string} articleId - 文章ID
   * @returns {Promise} 点赞结果
   */
  likeArticle(articleId) {
    return api.post(`/articles/${articleId}/like`)
  },

  /**
   * 取消点赞文章
   * @param {number|string} articleId - 文章ID
   * @returns {Promise} 取消点赞结果
   */
  unlikeArticle(articleId) {
    return api.delete(`/articles/${articleId}/like`)
  },

  /**
   * 收藏文章
   * @param {number|string} articleId - 文章ID
   * @returns {Promise} 收藏结果
   */
  favoriteArticle(articleId) {
    return api.post(`/articles/${articleId}/favorite`)
  },

  /**
   * 取消收藏文章
   * @param {number|string} articleId - 文章ID
   * @returns {Promise} 取消收藏结果
   */
  unfavoriteArticle(articleId) {
    return api.delete(`/articles/${articleId}/favorite`)
  },

  /**
   * 获取文章分类列表
   * @param {Object} params - 查询参数
   * @returns {Promise} 分类列表
   */
  getArticleCategories(params = {}) {
    return api.get('/articles/categories', { params })
  },

  /**
   * 创建文章分类
   * @param {Object} categoryData - 分类数据
   * @param {string} categoryData.name - 分类名称
   * @param {string} categoryData.slug - 分类标识
   * @param {string} categoryData.description - 分类描述
   * @param {string} categoryData.color - 分类颜色
   * @param {number} categoryData.parentId - 父分类ID
   * @returns {Promise} 创建的分类信息
   */
  createArticleCategory(categoryData) {
    return api.post('/articles/categories', categoryData)
  },

  /**
   * 更新文章分类
   * @param {number|string} categoryId - 分类ID
   * @param {Object} categoryData - 更新数据
   * @returns {Promise} 更新结果
   */
  updateArticleCategory(categoryId, categoryData) {
    return api.put(`/articles/categories/${categoryId}`, categoryData)
  },

  /**
   * 删除文章分类
   * @param {number|string} categoryId - 分类ID
   * @returns {Promise} 删除结果
   */
  deleteArticleCategory(categoryId) {
    return api.delete(`/articles/categories/${categoryId}`)
  },

  /**
   * 获取文章标签列表
   * @param {Object} params - 查询参数
   * @returns {Promise} 标签列表
   */
  getArticleTags(params = {}) {
    return api.get('/articles/tags', { params })
  },

  /**
   * 创建文章标签
   * @param {Object} tagData - 标签数据
   * @param {string} tagData.name - 标签名称
   * @param {string} tagData.color - 标签颜色
   * @returns {Promise} 创建的标签信息
   */
  createArticleTag(tagData) {
    return api.post('/articles/tags', tagData)
  },

  /**
   * 获取文章统计信息
   * @param {Object} params - 统计参数
   * @param {string} params.period - 统计周期 (week|month|quarter|year)
   * @param {string} params.startDate - 开始日期
   * @param {string} params.endDate - 结束日期
   * @returns {Promise} 统计数据
   */
  getArticleStats(params = {}) {
    return api.get('/articles/stats', { params })
  },

  /**
   * 获取文章趋势数据
   * @param {Object} params - 趋势参数
   * @param {number} params.days - 天数
   * @returns {Promise} 趋势数据
   */
  getArticleTrends(params = {}) {
    return api.get('/articles/trends', { params })
  },

  /**
   * 获取热门文章
   * @param {Object} params - 查询参数
   * @param {number} params.limit - 数量限制
   * @param {string} params.period - 时间范围 (day|week|month|all)
   * @returns {Promise} 热门文章列表
   */
  getPopularArticles(params = {}) {
    return api.get('/articles/popular', { params })
  },

  /**
   * 获取相关文章
   * @param {number|string} articleId - 文章ID
   * @param {number} limit - 数量限制
   * @returns {Promise} 相关文章列表
   */
  getRelatedArticles(articleId, limit = 5) {
    return api.get(`/articles/${articleId}/related`, {
      params: { limit }
    })
  },

  /**
   * 搜索文章
   * @param {Object} searchParams - 搜索参数
   * @param {string} searchParams.query - 搜索关键词
   * @param {Array} searchParams.categories - 分类筛选
   * @param {Array} searchParams.tags - 标签筛选
   * @param {string} searchParams.dateRange - 日期范围
   * @returns {Promise} 搜索结果
   */
  searchArticles(searchParams) {
    return api.get('/articles/search', { params: searchParams })
  },

  /**
   * 导出文章数据
   * @param {Object} params - 导出参数
   * @param {string} params.format - 导出格式 (csv|excel|pdf|markdown)
   * @param {Array} params.ids - 指定导出的文章ID
   * @param {Object} params.filters - 筛选条件
   * @returns {Promise} 导出文件信息
   */
  exportArticles(params = {}) {
    return api.get('/articles/export', { 
      params,
      responseType: 'blob'
    })
  },

  /**
   * 导入文章数据
   * @param {File} file - 导入文件
   * @param {string} format - 文件格式
   * @returns {Promise} 导入结果
   */
  importArticles(file, format = 'markdown') {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('format', format)
    
    return api.post('/articles/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 上传文章图片
   * @param {File} file - 图片文件
   * @param {string} type - 图片类型 (cover|content)
   * @returns {Promise} 上传结果包含图片URL
   */
  uploadArticleImage(file, type = 'content') {
    const formData = new FormData()
    formData.append('image', file)
    formData.append('type', type)
    
    return api.post('/articles/upload/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 获取文章评论
   * @param {number|string} articleId - 文章ID
   * @param {Object} params - 查询参数
   * @returns {Promise} 评论列表
   */
  getArticleComments(articleId, params = {}) {
    return api.get(`/articles/${articleId}/comments`, { params })
  },

  /**
   * 添加文章评论
   * @param {number|string} articleId - 文章ID
   * @param {Object} commentData - 评论数据
   * @returns {Promise} 评论结果
   */
  addArticleComment(articleId, commentData) {
    return api.post(`/articles/${articleId}/comments`, commentData)
  }
}

export default articlesApi