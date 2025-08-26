/**
 * 文章模块API导出文件
 * 统一管理文章相关的所有API接口
 */

import articlesApi from './articlesApi'

// 导出文章API
export { articlesApi }

// 默认导出所有API
export default {
  articlesApi
}

/**
 * API使用示例：
 * 
 * 在Store中使用：
 * import { articlesApi } from '@/features/articles/api'
 * 
 * export const useArticlesStore = defineStore('articles', () => {
 *   const articles = ref([])
 *   const loading = ref(false)
 *   
 *   const fetchArticles = async (params = {}) => {
 *     loading.value = true
 *     try {
 *       const response = await articlesApi.getArticles(params)
 *       articles.value = response.data.articles
 *       return response
 *     } catch (error) {
 *       ElMessage.error('获取文章失败')
 *       throw error
 *     } finally {
 *       loading.value = false
 *     }
 *   }
 *   
 *   const createArticle = async (articleData) => {
 *     try {
 *       const response = await articlesApi.createArticle(articleData)
 *       articles.value.unshift(response.data.article)
 *       ElMessage.success('文章创建成功')
 *       return response
 *     } catch (error) {
 *       ElMessage.error('文章创建失败')
 *       throw error
 *     }
 *   }
 *   
 *   return {
 *     articles,
 *     loading,
 *     fetchArticles,
 *     createArticle
 *   }
 * })
 * 
 * 在组件中使用：
 * import { articlesApi } from '@/features/articles/api'
 * import { ElMessage } from 'element-plus'
 * 
 * const handleLikeArticle = async (articleId) => {
 *   try {
 *     await articlesApi.likeArticle(articleId)
 *     ElMessage.success('点赞成功')
 *     // 刷新数据
 *     await fetchArticles()
 *   } catch (error) {
 *     ElMessage.error('点赞失败')
 *   }
 * }
 */