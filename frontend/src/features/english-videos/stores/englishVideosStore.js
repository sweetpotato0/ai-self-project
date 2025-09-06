import { ref } from 'vue'
import { defineStore } from 'pinia'
import { englishVideosApi } from '../api/englishVideosApi'

export const useEnglishVideosStore = defineStore('englishVideos', () => {
  // 状态
  const videoSeries = ref([])
  const currentSeries = ref(null)
  const episodes = ref([])
  const userProgress = ref({})
  const loading = ref(false)
  const error = ref(null)

  // 获取视频系列列表
  const fetchVideoSeries = async (params = {}) => {
    loading.value = true
    error.value = null
    try {
      const response = await englishVideosApi.getVideoSeries(params)
      videoSeries.value = response.data?.list || response.data || []
      return videoSeries.value
    } catch (err) {
      error.value = err.message || '获取视频系列失败'
      console.error('Fetch video series error:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 获取单个视频系列详情
  const getVideoSeries = async (seriesId) => {
    try {
      const response = await englishVideosApi.getVideoSeriesDetail(seriesId)
      currentSeries.value = response.data
      return response.data
    } catch (err) {
      error.value = err.message || '获取系列详情失败'
      console.error('Get video series error:', err)
      throw err
    }
  }

  // 获取系列的剧集列表
  const getEpisodes = async (seriesId) => {
    try {
      const response = await englishVideosApi.getEpisodes(seriesId)
      episodes.value = response.data?.list || response.data || []
      return episodes.value
    } catch (err) {
      error.value = err.message || '获取剧集列表失败'
      console.error('Get episodes error:', err)
      throw err
    }
  }

  // 获取单集详情
  const getEpisode = async (episodeId) => {
    try {
      const response = await englishVideosApi.getEpisodeDetail(episodeId)
      return response.data
    } catch (err) {
      error.value = err.message || '获取剧集详情失败'
      console.error('Get episode error:', err)
      throw err
    }
  }

  // 获取用户观看进度
  const getEpisodeProgress = async (episodeId) => {
    try {
      const response = await englishVideosApi.getEpisodeProgress(episodeId)
      return response.data
    } catch (err) {
      console.error('Get episode progress error:', err)
      return null
    }
  }

  // 更新观看进度
  const updateEpisodeProgress = async (episodeId, progressData) => {
    try {
      const response = await englishVideosApi.updateEpisodeProgress(episodeId, progressData)
      
      // 更新本地进度数据
      if (!userProgress.value) {
        userProgress.value = {}
      }
      userProgress.value[episodeId] = {
        ...userProgress.value[episodeId],
        ...progressData
      }

      return response.data
    } catch (err) {
      error.value = err.message || '更新进度失败'
      console.error('Update episode progress error:', err)
      throw err
    }
  }

  // 收藏/取消收藏系列
  const toggleSeriesLike = async (seriesId) => {
    try {
      const response = await englishVideosApi.toggleSeriesLike(seriesId)
      
      // 更新本地数据
      const series = videoSeries.value.find(s => s.id === seriesId)
      if (series) {
        series.is_liked = !series.is_liked
        series.like_count = response.data.like_count || series.like_count
      }

      return response.data
    } catch (err) {
      error.value = err.message || '操作失败'
      console.error('Toggle series like error:', err)
      throw err
    }
  }

  // 获取用户的观看统计
  const getUserVideoStats = async () => {
    try {
      const response = await englishVideosApi.getUserVideoStats()
      return response.data
    } catch (err) {
      console.error('Get user video stats error:', err)
      return {
        total_watch_time: 0,
        completed_episodes: 0,
        favorite_series: 0,
        learning_streak: 0
      }
    }
  }

  // 获取视频统计信息
  const getVideoStats = async () => {
    try {
      const response = await englishVideosApi.getVideoStats()
      return response.data
    } catch (err) {
      console.error('Get video stats error:', err)
      return {
        total_series: 0,
        total_episodes: 0,
        total_views: 0,
        watched_series: 0
      }
    }
  }

  // 获取用户进度
  const getUserProgress = async (params = {}) => {
    try {
      const response = await englishVideosApi.getUserProgress(params)
      return response.data
    } catch (err) {
      console.error('Get user progress error:', err)
      return []
    }
  }

  // 搜索视频系列
  const searchVideoSeries = async (query, filters = {}) => {
    loading.value = true
    error.value = null
    try {
      const response = await englishVideosApi.searchVideoSeries(query, filters)
      return response.data?.list || response.data || []
    } catch (err) {
      error.value = err.message || '搜索失败'
      console.error('Search video series error:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 获取推荐视频系列
  const getRecommendedSeries = async (limit = 10) => {
    try {
      const response = await englishVideosApi.getRecommendedSeries(limit)
      return response.data?.list || response.data || []
    } catch (err) {
      console.error('Get recommended series error:', err)
      return []
    }
  }

  // 管理员功能 - 创建视频系列
  const createVideoSeries = async (seriesData) => {
    try {
      const response = await englishVideosApi.createVideoSeries(seriesData)
      return response.data
    } catch (err) {
      error.value = err.message || '创建系列失败'
      console.error('Create video series error:', err)
      throw err
    }
  }

  // 管理员功能 - 更新视频系列
  const updateVideoSeries = async (seriesId, seriesData) => {
    try {
      const response = await englishVideosApi.updateVideoSeries(seriesId, seriesData)
      
      // 更新本地数据
      const index = videoSeries.value.findIndex(s => s.id === seriesId)
      if (index !== -1) {
        videoSeries.value[index] = { ...videoSeries.value[index], ...response.data }
      }
      
      return response.data
    } catch (err) {
      error.value = err.message || '更新系列失败'
      console.error('Update video series error:', err)
      throw err
    }
  }

  // 管理员功能 - 删除视频系列
  const deleteVideoSeries = async (seriesId) => {
    try {
      const response = await englishVideosApi.deleteVideoSeries(seriesId)
      
      // 从本地数据中移除
      videoSeries.value = videoSeries.value.filter(s => s.id !== seriesId)
      
      return response.data
    } catch (err) {
      error.value = err.message || '删除系列失败'
      console.error('Delete video series error:', err)
      throw err
    }
  }

  // 管理员功能 - 创建剧集
  const createEpisode = async (seriesId, episodeData) => {
    try {
      const response = await englishVideosApi.createEpisode(seriesId, episodeData)
      return response.data
    } catch (err) {
      error.value = err.message || '创建剧集失败'
      console.error('Create episode error:', err)
      throw err
    }
  }

  // 管理员功能 - 更新剧集
  const updateEpisode = async (episodeId, episodeData) => {
    try {
      const response = await englishVideosApi.updateEpisode(episodeId, episodeData)
      
      // 更新本地数据
      const index = episodes.value.findIndex(e => e.id === episodeId)
      if (index !== -1) {
        episodes.value[index] = { ...episodes.value[index], ...response.data }
      }
      
      return response.data
    } catch (err) {
      error.value = err.message || '更新剧集失败'
      console.error('Update episode error:', err)
      throw err
    }
  }

  // 管理员功能 - 删除剧集
  const deleteEpisode = async (episodeId) => {
    try {
      const response = await englishVideosApi.deleteEpisode(episodeId)
      
      // 从本地数据中移除
      episodes.value = episodes.value.filter(e => e.id !== episodeId)
      
      return response.data
    } catch (err) {
      error.value = err.message || '删除剧集失败'
      console.error('Delete episode error:', err)
      throw err
    }
  }

  // 管理员功能 - 批量导入剧集
  const batchImportEpisodes = async (seriesId, episodesData) => {
    try {
      const response = await englishVideosApi.batchImportEpisodes(seriesId, episodesData)
      return response.data
    } catch (err) {
      error.value = err.message || '批量导入剧集失败'
      console.error('Batch import episodes error:', err)
      throw err
    }
  }

  // 管理员功能 - 创建未分类剧集
  const createUncategorizedEpisode = async (episodeData) => {
    try {
      const response = await englishVideosApi.createUncategorizedEpisode(episodeData)
      return response.data
    } catch (err) {
      error.value = err.message || '创建未分类剧集失败'
      console.error('Create uncategorized episode error:', err)
      throw err
    }
  }

  // 重置状态
  const resetState = () => {
    videoSeries.value = []
    currentSeries.value = null
    episodes.value = []
    userProgress.value = {}
    error.value = null
  }

  return {
    // 状态
    videoSeries,
    currentSeries,
    episodes,
    userProgress,
    loading,
    error,

    // 方法
    fetchVideoSeries,
    getVideoSeries,
    getEpisodes,
    getEpisode,
    getEpisodeProgress,
    updateEpisodeProgress,
    toggleSeriesLike,
    getUserVideoStats,
    getVideoStats,
    getUserProgress,
    searchVideoSeries,
    getRecommendedSeries,
    resetState,

    // 管理员方法
    createVideoSeries,
    updateVideoSeries,
    deleteVideoSeries,
    createEpisode,
    updateEpisode,
    deleteEpisode,
    batchImportEpisodes,
    createUncategorizedEpisode
  }
})