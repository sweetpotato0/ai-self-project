/**
 * 英语学习模块API导出文件
 * 统一管理英语学习相关的所有API接口
 */

import englishLearningApi from './englishLearningApi'

// 导出英语学习API
export { englishLearningApi }

// 默认导出所有API
export default {
  englishLearningApi
}

/**
 * API使用示例：
 * 
 * 在Store中使用：
 * import { englishLearningApi } from '@/features/english-learning/api'
 * 
 * export const useEnglishLearningStore = defineStore('englishLearning', () => {
 *   const categories = ref([])
 *   const songs = ref([])
 *   const userProgress = ref([])
 *   const loading = ref(false)
 *   
 *   // 获取分类列表
 *   const fetchCategories = async (params = {}) => {
 *     loading.value = true
 *     try {
 *       const response = await englishLearningApi.getCategories(params)
 *       categories.value = response.categories
 *       return response
 *     } catch (error) {
 *       ElMessage.error('获取分类失败')
 *       throw error
 *     } finally {
 *       loading.value = false
 *     }
 *   }
 *   
 *   // 获取歌曲列表
 *   const fetchSongs = async (params = {}) => {
 *     loading.value = true
 *     try {
 *       const response = await englishLearningApi.getSongs(params)
 *       songs.value = response.songs
 *       return response
 *     } catch (error) {
 *       ElMessage.error('获取歌曲失败')
 *       throw error
 *     } finally {
 *       loading.value = false
 *     }
 *   }
 *   
 *   // 点赞歌曲
 *   const likeSong = async (songId) => {
 *     try {
 *       await englishLearningApi.likeSong(songId)
 *       ElMessage.success('点赞成功')
 *       // 更新本地数据
 *       const song = songs.value.find(s => s.id === songId)
 *       if (song) {
 *         song.like_count++
 *         song.is_liked = true
 *       }
 *     } catch (error) {
 *       ElMessage.error('点赞失败')
 *       throw error
 *     }
 *   }
 *   
 *   // 更新学习进度
 *   const updateProgress = async (songId, progressData) => {
 *     try {
 *       await englishLearningApi.updateProgress(songId, progressData)
 *       ElMessage.success('进度更新成功')
 *       // 刷新进度数据
 *       await fetchUserProgress()
 *     } catch (error) {
 *       ElMessage.error('进度更新失败')
 *       throw error
 *     }
 *   }
 *   
 *   // 获取用户进度
 *   const fetchUserProgress = async (params = {}) => {
 *     try {
 *       const response = await englishLearningApi.getUserProgress(params)
 *       userProgress.value = response
 *       return response
 *     } catch (error) {
 *       ElMessage.error('获取学习进度失败')
 *       throw error
 *     }
 *   }
 *   
 *   return {
 *     categories,
 *     songs,
 *     userProgress,
 *     loading,
 *     fetchCategories,
 *     fetchSongs,
 *     likeSong,
 *     updateProgress,
 *     fetchUserProgress
 *   }
 * })
 * 
 * 在组件中使用：
 * import { englishLearningApi } from '@/features/english-learning/api'
 * import { ElMessage } from 'element-plus'
 * 
 * // 获取推荐歌曲
 * const getRecommendations = async () => {
 *   try {
 *     const response = await englishLearningApi.getRecommendations({ limit: 10 })
 *     console.log('推荐歌曲:', response.songs)
 *   } catch (error) {
 *     ElMessage.error('获取推荐失败')
 *   }
 * }
 * 
 * // 搜索歌曲
 * const searchSongs = async (searchTerm) => {
 *   try {
 *     const response = await englishLearningApi.searchSongs({
 *       query: searchTerm,
 *       only_published: true,
 *       sort_by: 'view_count',
 *       sort_order: 'desc'
 *     })
 *     console.log('搜索结果:', response.songs)
 *   } catch (error) {
 *     ElMessage.error('搜索失败')
 *   }
 * }
 * 
 * // 创建学习计划
 * const createPlan = async () => {
 *   try {
 *     await englishLearningApi.createLearningPlan({
 *       title: '每日英语学习',
 *       description: '每天学习30分钟英语歌曲',
 *       goal_type: 'daily',
 *       target_minutes: 30,
 *       start_date: new Date(),
 *       end_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000) // 30天后
 *     })
 *     ElMessage.success('学习计划创建成功')
 *   } catch (error) {
 *     ElMessage.error('创建学习计划失败')
 *   }
 * }
 */