import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { englishLearningApi } from '../api'

export const useEnglishLearningStore = defineStore('englishLearning', () => {
  // 状态
  const categories = ref([])
  const songs = ref([])
  const userProgress = ref([])
  const userStats = ref(null)
  const recommendations = ref([])
  const currentSong = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // 分页信息
  const songsTotal = ref(0)
  const songsPage = ref(1)
  const songsLimit = ref(20)
  const songsTotalPages = ref(0)

  const categoriesTotal = ref(0)
  const categoriesPage = ref(1)
  const categoriesLimit = ref(20)
  const categoriesTotalPages = ref(0)

  // 计算属性
  const completedSongs = computed(() => {
    return userProgress.value.filter(p => p.is_completed).length
  })

  const totalStudyTime = computed(() => {
    return userProgress.value.reduce((total, p) => total + (p.study_time_minutes || 0), 0)
  })

  const favoriteCategory = computed(() => {
    if (!userStats.value) return null
    return userStats.value.favorite_category
  })

  const currentLevel = computed(() => {
    return userStats.value?.level || 1
  })

  // Actions - 分类管理
  const fetchCategories = async (params = {}) => {
    loading.value = true
    error.value = null
    try {
      const response = await englishLearningApi.getCategories(params)
      const data = response.data || response
      categories.value = data.categories || []
      categoriesTotal.value = data.total || 0
      categoriesPage.value = data.page || 1
      categoriesLimit.value = data.limit || 20
      categoriesTotalPages.value = data.total_pages || 0
      return response
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const createCategory = async (categoryData) => {
    try {
      const response = await englishLearningApi.createCategory(categoryData)
      categories.value.unshift(response)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const updateCategory = async (categoryId, categoryData) => {
    try {
      const response = await englishLearningApi.updateCategory(categoryId, categoryData)
      const index = categories.value.findIndex(c => c.id === categoryId)
      if (index !== -1) {
        categories.value[index] = { ...categories.value[index], ...categoryData }
      }
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const deleteCategory = async (categoryId) => {
    try {
      await englishLearningApi.deleteCategory(categoryId)
      categories.value = categories.value.filter(c => c.id !== categoryId)
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  // Actions - 歌曲管理
  const fetchSongs = async (params = {}) => {
    loading.value = true
    error.value = null
    try {
      const response = await englishLearningApi.getSongs(params)
      const data = response.data || response
      songs.value = data.songs || []
      songsTotal.value = data.total || 0
      songsPage.value = data.page || 1
      songsLimit.value = data.limit || 20
      songsTotalPages.value = data.total_pages || 0
      return response
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchSong = async (songId) => {
    try {
      const response = await englishLearningApi.getSong(songId)
      currentSong.value = response
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const createSong = async (songData) => {
    try {
      const response = await englishLearningApi.createSong(songData)
      songs.value.unshift(response)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const updateSong = async (songId, songData) => {
    try {
      const response = await englishLearningApi.updateSong(songId, songData)
      const index = songs.value.findIndex(s => s.id === songId)
      if (index !== -1) {
        songs.value[index] = { ...songs.value[index], ...songData }
      }
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const deleteSong = async (songId) => {
    try {
      await englishLearningApi.deleteSong(songId)
      songs.value = songs.value.filter(s => s.id !== songId)
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const likeSong = async (songId) => {
    try {
      await englishLearningApi.likeSong(songId)
      const song = songs.value.find(s => s.id === songId)
      if (song) {
        song.like_count = (song.like_count || 0) + 1
        song.is_liked = true
      }
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const unlikeSong = async (songId) => {
    try {
      await englishLearningApi.unlikeSong(songId)
      const song = songs.value.find(s => s.id === songId)
      if (song) {
        song.like_count = Math.max((song.like_count || 1) - 1, 0)
        song.is_liked = false
      }
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  // Actions - 学习进度管理
  const updateProgress = async (songId, progressData) => {
    try {
      await englishLearningApi.updateProgress(songId, progressData)
      
      // 更新本地进度数据
      const existingProgress = userProgress.value.find(p => p.song_id === songId)
      if (existingProgress) {
        Object.assign(existingProgress, progressData)
      } else {
        userProgress.value.push({
          song_id: songId,
          ...progressData
        })
      }

      // 如果完成了歌曲，刷新统计数据
      if (progressData.is_completed) {
        await fetchUserStats()
      }
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const fetchUserProgress = async (params = {}) => {
    try {
      const response = await englishLearningApi.getUserProgress(params)
      const data = response.data || response
      userProgress.value = data || []
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const fetchUserStats = async () => {
    try {
      const response = await englishLearningApi.getUserStats()
      const data = response.data || response
      userStats.value = data || {}
      return response
    } catch (err) {
      error.value = err.message
      // 不抛出错误，因为统计数据是可选的
      console.warn('Failed to fetch user stats:', err)
    }
  }

  // Actions - 推荐系统
  const fetchRecommendations = async (params = {}) => {
    try {
      const response = await englishLearningApi.getRecommendations(params)
      const data = response.data || response
      recommendations.value = data.songs || []
      return response
    } catch (err) {
      error.value = err.message
      // 不抛出错误，因为推荐是可选的
      console.warn('Failed to fetch recommendations:', err)
    }
  }

  // Actions - 搜索
  const searchSongs = async (searchParams) => {
    loading.value = true
    try {
      const response = await englishLearningApi.searchSongs(searchParams)
      const data = response.data || response
      songs.value = data.songs || []
      songsTotal.value = data.total || 0
      return response
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const getPopularSongs = async (params = {}) => {
    try {
      const response = await englishLearningApi.getPopularSongs(params)
      const data = response.data || response
      return data.songs || []
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const getLatestSongs = async (params = {}) => {
    try {
      const response = await englishLearningApi.getLatestSongs(params)
      const data = response.data || response
      return data.songs || []
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  // Actions - 词汇管理
  const getUserVocabularies = async (params = {}) => {
    try {
      const response = await englishLearningApi.getUserVocabularies(params)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const updateVocabularyMastery = async (vocabularyId, masteryData) => {
    try {
      const response = await englishLearningApi.updateVocabularyMastery(vocabularyId, masteryData)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  // Actions - 学习计划
  const getLearningPlans = async (params = {}) => {
    try {
      const response = await englishLearningApi.getLearningPlans(params)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const createLearningPlan = async (planData) => {
    try {
      const response = await englishLearningApi.createLearningPlan(planData)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  // Actions - 学习会话
  const startStudySession = async (sessionData) => {
    try {
      const response = await englishLearningApi.startStudySession(sessionData)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const endStudySession = async (sessionId, sessionData) => {
    try {
      const response = await englishLearningApi.endStudySession(sessionId, sessionData)
      // 刷新用户进度和统计
      await Promise.all([
        fetchUserProgress(),
        fetchUserStats()
      ])
      return response
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  // 重置状态
  const reset = () => {
    categories.value = []
    songs.value = []
    userProgress.value = []
    userStats.value = null
    recommendations.value = []
    currentSong.value = null
    loading.value = false
    error.value = null
  }

  // 清除错误
  const clearError = () => {
    error.value = null
  }

  return {
    // 状态
    categories,
    songs,
    userProgress,
    userStats,
    recommendations,
    currentSong,
    loading,
    error,
    
    // 分页信息
    songsTotal,
    songsPage,
    songsLimit,
    songsTotalPages,
    categoriesTotal,
    categoriesPage,
    categoriesLimit,
    categoriesTotalPages,

    // 计算属性
    completedSongs,
    totalStudyTime,
    favoriteCategory,
    currentLevel,

    // Actions
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    fetchSongs,
    fetchSong,
    createSong,
    updateSong,
    deleteSong,
    likeSong,
    unlikeSong,
    updateProgress,
    fetchUserProgress,
    fetchUserStats,
    fetchRecommendations,
    searchSongs,
    getPopularSongs,
    getLatestSongs,
    getUserVocabularies,
    updateVocabularyMastery,
    getLearningPlans,
    createLearningPlan,
    startStudySession,
    endStudySession,
    reset,
    clearError
  }
})