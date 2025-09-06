import api from '@/api/index'

/**
 * 英语视频模块API接口
 */
export const englishVideosApi = {
  /**
   * 获取视频系列列表
   * @param {Object} params - 查询参数
   * @returns {Promise} API响应
   */
  getVideoSeries(params = {}) {
    const defaultParams = {
      page: 1,
      page_size: 20,
      ...params
    }
    return api.get('/english-videos/series', { params: defaultParams })
  },

  /**
   * 获取视频系列详情
   * @param {number} seriesId - 系列ID
   * @returns {Promise} API响应
   */
  getVideoSeriesDetail(seriesId) {
    return api.get(`/english-videos/series/${seriesId}`)
  },

  /**
   * 获取系列的剧集列表
   * @param {number} seriesId - 系列ID
   * @param {Object} params - 查询参数
   * @returns {Promise} API响应
   */
  getEpisodes(seriesId, params = {}) {
    const defaultParams = {
      page: 1,
      page_size: 20,
      sort_order: 'asc',
      ...params
    }
    return api.get(`/english-videos/series/${seriesId}/episodes`, { params: defaultParams })
  },

  /**
   * 获取剧集详情
   * @param {number} episodeId - 剧集ID
   * @returns {Promise} API响应
   */
  getEpisodeDetail(episodeId) {
    return api.get(`/english-videos/episodes/${episodeId}`)
  },

  /**
   * 获取用户对剧集的观看进度
   * @param {number} episodeId - 剧集ID
   * @returns {Promise} API响应
   */
  getEpisodeProgress(episodeId) {
    return api.get(`/english-videos/episodes/${episodeId}/progress`)
  },

  /**
   * 更新剧集观看进度
   * @param {number} episodeId - 剧集ID
   * @param {Object} progressData - 进度数据
   * @returns {Promise} API响应
   */
  updateEpisodeProgress(episodeId, progressData) {
    return api.post(`/english-videos/episodes/${episodeId}/progress`, progressData)
  },

  /**
   * 收藏/取消收藏系列
   * @param {number} seriesId - 系列ID
   * @returns {Promise} API响应
   */
  toggleSeriesLike(seriesId) {
    return api.post(`/english-videos/series/${seriesId}/toggle-like`)
  },

  /**
   * 获取用户的视频观看统计
   * @returns {Promise} API响应
   */
  getUserVideoStats() {
    return api.get('/english-videos/user/stats')
  },

  /**
   * 获取视频统计信息
   * @returns {Promise} API响应
   */
  getVideoStats() {
    return api.get('/english-videos/stats')
  },

  /**
   * 获取用户进度
   * @param {Object} params - 查询参数
   * @returns {Promise} API响应
   */
  getUserProgress(params = {}) {
    const defaultParams = {
      page: 1,
      page_size: 20,
      sort_order: 'desc',
      ...params
    }
    return api.get('/english-videos/progress', { params: defaultParams })
  },

  /**
   * 搜索视频系列
   * @param {string} query - 搜索关键词
   * @param {Object} filters - 过滤条件
   * @returns {Promise} API响应
   */
  searchVideoSeries(query, filters = {}) {
    const defaultParams = {
      q: query,
      page: 1,
      page_size: 20,
      ...filters
    }
    return api.get('/english-videos/series/search', { params: defaultParams })
  },

  /**
   * 获取推荐的视频系列
   * @param {number} limit - 限制数量
   * @returns {Promise} API响应
   */
  getRecommendedSeries(limit = 10) {
    return api.get('/english-videos/series/recommended', { params: { limit } })
  },

  /**
   * 管理员接口 - 创建视频系列
   * @param {Object} seriesData - 系列数据
   * @returns {Promise} API响应
   */
  createVideoSeries(seriesData) {
    return api.post('/english-videos/series/admin', seriesData)
  },

  /**
   * 管理员接口 - 更新视频系列
   * @param {number} seriesId - 系列ID
   * @param {Object} seriesData - 系列数据
   * @returns {Promise} API响应
   */
  updateVideoSeries(seriesId, seriesData) {
    return api.put(`/english-videos/series/admin/${seriesId}`, seriesData)
  },

  /**
   * 管理员接口 - 删除视频系列
   * @param {number} seriesId - 系列ID
   * @returns {Promise} API响应
   */
  deleteVideoSeries(seriesId) {
    return api.delete(`/english-videos/series/admin/${seriesId}`)
  },

  /**
   * 管理员接口 - 创建剧集
   * @param {number} seriesId - 系列ID
   * @param {Object} episodeData - 剧集数据
   * @returns {Promise} API响应
   */
  createEpisode(seriesId, episodeData) {
    return api.post(`/english-videos/series/admin/${seriesId}/episodes`, episodeData)
  },

  /**
   * 管理员接口 - 更新剧集
   * @param {number} episodeId - 剧集ID
   * @param {Object} episodeData - 剧集数据
   * @returns {Promise} API响应
   */
  updateEpisode(episodeId, episodeData) {
    return api.put(`/english-videos/episodes/admin/${episodeId}`, episodeData)
  },

  /**
   * 管理员接口 - 删除剧集
   * @param {number} episodeId - 剧集ID
   * @returns {Promise} API响应
   */
  deleteEpisode(episodeId) {
    return api.delete(`/english-videos/episodes/admin/${episodeId}`)
  },

  /**
   * 管理员接口 - 批量导入剧集
   * @param {number} seriesId - 系列ID
   * @param {Array} episodesData - 剧集数据数组
   * @returns {Promise} API响应
   */
  batchImportEpisodes(seriesId, episodesData) {
    return api.post(`/english-videos/series/admin/${seriesId}/episodes/batch`, { episodes: episodesData })
  },

  /**
   * 管理员接口 - 创建未分类剧集
   * @param {Object} episodeData - 剧集数据
   * @returns {Promise} API响应
   */
  createUncategorizedEpisode(episodeData) {
    return api.post('/english-videos/episodes/admin/uncategorized', episodeData)
  }
}

export default englishVideosApi