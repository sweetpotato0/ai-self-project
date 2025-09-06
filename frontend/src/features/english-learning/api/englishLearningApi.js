import api from '@/api/index'

/**
 * 英语学习模块API接口
 * 包含分类管理、歌曲管理、学习进度、推荐系统等功能
 */
export const englishLearningApi = {
  // ====== 分类管理 ======
  
  /**
   * 获取学习分类列表
   * @param {Object} params - 查询参数
   * @param {boolean} params.is_active - 是否激活
   * @param {number} params.page - 页码
   * @param {number} params.limit - 每页数量
   * @returns {Promise} 分类列表和分页信息
   */
  getCategories(params = {}) {
    return api.get('/learning/categories', { params })
  },

  /**
   * 创建学习分类
   * @param {Object} categoryData - 分类数据
   * @param {string} categoryData.name - 分类名称
   * @param {string} categoryData.name_cn - 中文名称
   * @param {string} categoryData.description - 描述
   * @param {string} categoryData.icon - 图标
   * @param {string} categoryData.color - 颜色
   * @param {boolean} categoryData.is_active - 是否激活
   * @param {number} categoryData.sort - 排序
   * @returns {Promise} 创建的分类信息
   */
  createCategory(categoryData) {
    return api.post('/learning/categories', categoryData)
  },

  /**
   * 更新学习分类
   * @param {number|string} categoryId - 分类ID
   * @param {Object} categoryData - 更新数据
   * @returns {Promise} 更新结果
   */
  updateCategory(categoryId, categoryData) {
    return api.put(`/learning/categories/${categoryId}`, categoryData)
  },

  /**
   * 删除学习分类
   * @param {number|string} categoryId - 分类ID
   * @returns {Promise} 删除结果
   */
  deleteCategory(categoryId) {
    return api.delete(`/learning/categories/${categoryId}`)
  },

  // ====== 歌曲/学习材料管理 ======
  
  /**
   * 获取歌曲列表
   * @param {Object} params - 查询参数
   * @param {number} params.category_id - 分类ID
   * @param {number} params.difficulty - 难度等级 (1-5)
   * @param {string} params.age_range - 年龄范围 (3-6|7-12|13-18|18+)
   * @param {boolean} params.is_published - 是否已发布
   * @param {string} params.search - 搜索关键词
   * @param {string} params.tags - 标签筛选
   * @param {string} params.sort_by - 排序字段 (created_at|title|view_count|like_count|difficulty)
   * @param {string} params.sort_order - 排序方向 (asc|desc)
   * @param {string} params.sort - 排序参数 (格式：field:direction)
   * @param {number} params.page - 页码
   * @param {number} params.limit - 每页数量
   * @returns {Promise} 歌曲列表和分页信息
   */
  getSongs(params = {}) {
    return api.get('/learning/songs', { params })
  },

  /**
   * 获取单个歌曲详情
   * @param {number|string} songId - 歌曲ID
   * @returns {Promise} 歌曲详细信息
   */
  getSong(songId) {
    return api.get(`/learning/songs/${songId}`)
  },

  /**
   * 创建新歌曲
   * @param {Object} songData - 歌曲数据
   * @param {string} songData.title - 歌曲标题
   * @param {string} songData.title_cn - 中文标题
   * @param {string} songData.description - 描述
   * @param {string} songData.lyrics - 歌词
   * @param {string} songData.lyrics_cn - 中文歌词翻译
   * @param {string} songData.audio_url - 音频URL
   * @param {string} songData.video_url - 视频URL
   * @param {string} songData.cover_image - 封面图片
   * @param {number} songData.category_id - 分类ID
   * @param {number} songData.difficulty - 难度等级
   * @param {string} songData.age_range - 适合年龄范围
   * @param {Array} songData.tags - 标签数组
   * @param {boolean} songData.is_published - 是否发布
   * @param {number} songData.duration_seconds - 时长（秒）
   * @param {number} songData.sort - 排序
   * @returns {Promise} 创建的歌曲信息
   */
  createSong(songData) {
    return api.post('/learning/songs', songData)
  },

  /**
   * 更新歌曲
   * @param {number|string} songId - 歌曲ID
   * @param {Object} songData - 更新数据
   * @returns {Promise} 更新结果
   */
  updateSong(songId, songData) {
    return api.put(`/learning/songs/${songId}`, songData)
  },

  /**
   * 删除歌曲
   * @param {number|string} songId - 歌曲ID
   * @returns {Promise} 删除结果
   */
  deleteSong(songId) {
    return api.delete(`/learning/songs/${songId}`)
  },

  /**
   * 点赞歌曲
   * @param {number|string} songId - 歌曲ID
   * @returns {Promise} 点赞结果
   */
  likeSong(songId) {
    return api.post(`/learning/songs/${songId}/like`, {})
  },

  /**
   * 取消点赞歌曲
   * @param {number|string} songId - 歌曲ID
   * @returns {Promise} 取消点赞结果
   */
  unlikeSong(songId) {
    return api.delete(`/learning/songs/${songId}/like`)
  },

  // ====== 用户学习进度 ======
  
  /**
   * 更新学习进度
   * @param {number|string} songId - 歌曲ID
   * @param {Object} progressData - 进度数据
   * @param {number} progressData.progress - 进度百分比 (0-100)
   * @param {boolean} progressData.is_completed - 是否完成
   * @param {number} progressData.play_count - 播放次数
   * @param {number} progressData.study_time_minutes - 学习时长（分钟）
   * @param {string} progressData.notes - 学习笔记
   * @returns {Promise} 更新结果
   */
  updateProgress(songId, progressData) {
    return api.put(`/learning/songs/${songId}/progress`, progressData)
  },

  /**
   * 获取用户学习进度
   * @param {Object} params - 查询参数
   * @param {number} params.song_id - 指定歌曲ID（可选）
   * @returns {Promise} 用户学习进度列表
   */
  getUserProgress(params = {}) {
    return api.get('/learning/user/progress', { params })
  },

  /**
   * 获取用户学习统计
   * @returns {Promise} 用户学习统计数据
   */
  getUserStats() {
    return api.get('/learning/user/stats')
  },

  /**
   * 获取推荐歌曲
   * @param {Object} params - 查询参数
   * @param {number} params.limit - 推荐数量限制，默认10
   * @returns {Promise} 推荐歌曲列表
   */
  getRecommendations(params = {}) {
    return api.get('/learning/user/recommendations', { params })
  },

  // ====== 词汇管理 ======

  /**
   * 获取歌曲相关词汇
   * @param {number|string} songId - 歌曲ID
   * @returns {Promise} 词汇列表
   */
  getSongVocabularies(songId) {
    return api.get(`/learning/songs/${songId}/vocabularies`)
  },

  /**
   * 添加词汇到歌曲
   * @param {number|string} songId - 歌曲ID
   * @param {Array} vocabularyIds - 词汇ID数组
   * @returns {Promise} 添加结果
   */
  addVocabulariesToSong(songId, vocabularyIds) {
    return api.post(`/learning/songs/${songId}/vocabularies`, {
      vocabulary_ids: vocabularyIds
    })
  },

  /**
   * 从歌曲移除词汇
   * @param {number|string} songId - 歌曲ID
   * @param {number|string} vocabularyId - 词汇ID
   * @returns {Promise} 移除结果
   */
  removeVocabularyFromSong(songId, vocabularyId) {
    return api.delete(`/learning/songs/${songId}/vocabularies/${vocabularyId}`)
  },

  /**
   * 获取用户词汇掌握情况
   * @param {Object} params - 查询参数
   * @param {string} params.mastery_level - 掌握程度 (not_learned|learning|mastered)
   * @param {number} params.page - 页码
   * @param {number} params.limit - 每页数量
   * @returns {Promise} 用户词汇掌握情况
   */
  getUserVocabularies(params = {}) {
    return api.get('/learning/user/vocabularies', { params })
  },

  /**
   * 更新用户词汇掌握情况
   * @param {number|string} vocabularyId - 词汇ID
   * @param {Object} masteryData - 掌握情况数据
   * @param {string} masteryData.mastery_level - 掌握程度
   * @param {number} masteryData.practice_count - 练习次数
   * @param {string} masteryData.notes - 学习笔记
   * @returns {Promise} 更新结果
   */
  updateVocabularyMastery(vocabularyId, masteryData) {
    return api.put(`/learning/vocabularies/${vocabularyId}/mastery`, masteryData)
  },

  // ====== 学习计划 ======

  /**
   * 获取用户学习计划
   * @param {Object} params - 查询参数
   * @param {string} params.status - 计划状态 (active|completed|paused)
   * @returns {Promise} 学习计划列表
   */
  getLearningPlans(params = {}) {
    return api.get('/learning/plans', { params })
  },

  /**
   * 创建学习计划
   * @param {Object} planData - 计划数据
   * @param {string} planData.title - 计划标题
   * @param {string} planData.description - 计划描述
   * @param {string} planData.goal_type - 目标类型 (daily|weekly|monthly)
   * @param {number} planData.target_minutes - 目标学习分钟数
   * @param {Date} planData.start_date - 开始日期
   * @param {Date} planData.end_date - 结束日期
   * @returns {Promise} 创建的计划信息
   */
  createLearningPlan(planData) {
    return api.post('/learning/plans', planData)
  },

  /**
   * 更新学习计划
   * @param {number|string} planId - 计划ID
   * @param {Object} planData - 更新数据
   * @returns {Promise} 更新结果
   */
  updateLearningPlan(planId, planData) {
    return api.put(`/learning/plans/${planId}`, planData)
  },

  /**
   * 删除学习计划
   * @param {number|string} planId - 计划ID
   * @returns {Promise} 删除结果
   */
  deleteLearningPlan(planId) {
    return api.delete(`/learning/plans/${planId}`)
  },

  // ====== 学习会话 ======

  /**
   * 开始学习会话
   * @param {Object} sessionData - 会话数据
   * @param {number} sessionData.song_id - 歌曲ID
   * @param {number} sessionData.plan_id - 计划ID（可选）
   * @returns {Promise} 会话信息
   */
  startStudySession(sessionData) {
    return api.post('/learning/sessions', sessionData)
  },

  /**
   * 结束学习会话
   * @param {number|string} sessionId - 会话ID
   * @param {Object} sessionData - 会话结束数据
   * @param {number} sessionData.duration_minutes - 学习时长
   * @param {number} sessionData.progress_gained - 进度增长
   * @param {Array} sessionData.vocabularies_learned - 学会的词汇ID数组
   * @returns {Promise} 结束结果
   */
  endStudySession(sessionId, sessionData) {
    return api.put(`/learning/sessions/${sessionId}/end`, sessionData)
  },

  /**
   * 获取学习会话历史
   * @param {Object} params - 查询参数
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @param {number} params.page - 页码
   * @param {number} params.limit - 每页数量
   * @returns {Promise} 学习会话历史
   */
  getStudySessions(params = {}) {
    return api.get('/learning/sessions', { params })
  },

  // ====== 搜索和筛选 ======

  /**
   * 高级搜索歌曲
   * @param {Object} searchParams - 搜索参数
   * @param {string} searchParams.query - 搜索关键词
   * @param {Array} searchParams.categories - 分类筛选
   * @param {Array} searchParams.difficulty_range - 难度范围 [min, max]
   * @param {Array} searchParams.age_ranges - 年龄范围筛选
   * @param {Array} searchParams.tags - 标签筛选
   * @param {boolean} searchParams.only_published - 只显示已发布
   * @param {string} searchParams.sort_by - 排序字段
   * @param {string} searchParams.sort_order - 排序方向
   * @returns {Promise} 搜索结果
   */
  searchSongs(searchParams) {
    return api.get('/learning/songs/search', { params: searchParams })
  },

  /**
   * 获取热门歌曲
   * @param {Object} params - 查询参数
   * @param {number} params.limit - 数量限制
   * @param {string} params.period - 时间范围 (day|week|month|all)
   * @param {number} params.category_id - 分类筛选
   * @returns {Promise} 热门歌曲列表
   */
  getPopularSongs(params = {}) {
    return api.get('/learning/songs/popular', { params })
  },

  /**
   * 获取最新发布的歌曲
   * @param {Object} params - 查询参数
   * @param {number} params.limit - 数量限制
   * @param {number} params.category_id - 分类筛选
   * @returns {Promise} 最新歌曲列表
   */
  getLatestSongs(params = {}) {
    return api.get('/learning/songs/latest', { params })
  }
}

export default englishLearningApi