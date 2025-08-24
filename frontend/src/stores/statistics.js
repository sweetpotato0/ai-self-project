import { statisticsApi } from '@/api/statistics'
import { ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useStatisticsStore = defineStore('statistics', () => {
  // 状态
  const todoStats = ref({
    total: 0,
    pending: 0,
    inProgress: 0,
    completed: 0,
    cancelled: 0,
    overdue: 0
  })

  const articleStats = ref({
    total: 0,
    published: 0,
    draft: 0,
    archived: 0,
    totalViews: 0,
    totalLikes: 0
  })

  const todoTrends = ref([])
  const articleTrends = ref([])
  const loading = ref(false)

  // 计算属性
  const allStats = computed(() => ({
    todo: todoStats.value,
    article: articleStats.value
  }))

  // 方法
  const fetchStatistics = async (types = ['todo', 'article']) => {
    loading.value = true
    try {
      const promises = types.map(type =>
        statisticsApi.getStatistics({ type })
      )

      const results = await Promise.all(promises)

      results.forEach((result, index) => {
        const type = types[index]
        if (result.code === 200) {
          if (type === 'todo') {
            todoStats.value = result.data
          } else if (type === 'article') {
            articleStats.value = result.data
          }
        }
      })
    } catch (error) {
      console.error('Failed to fetch statistics:', error)
      ElMessage.error('获取统计数据失败')
    } finally {
      loading.value = false
    }
  }

  const fetchTodoStats = async () => {
    try {
      const response = await statisticsApi.getTodoStats()
      if (response.code === 200) {
        todoStats.value = response.data
      }
    } catch (error) {
      console.error('Failed to fetch todo stats:', error)
      ElMessage.error('获取任务统计失败')
    }
  }

  const fetchArticleStats = async () => {
    try {
      const response = await statisticsApi.getArticleStats()
      if (response.code === 200) {
        articleStats.value = response.data
      }
    } catch (error) {
      console.error('Failed to fetch article stats:', error)
      ElMessage.error('获取文章统计失败')
    }
  }

  const fetchTrends = async (type, params = {}) => {
    try {
      const response = await statisticsApi.getTrends({ type, ...params })
      if (response.code === 200) {
        if (type === 'todo') {
          todoTrends.value = response.data
        } else if (type === 'article') {
          articleTrends.value = response.data
        }
      }
    } catch (error) {
      console.error(`Failed to fetch ${type} trends:`, error)
      ElMessage.error(`获取${type === 'todo' ? '任务' : '文章'}趋势失败`)
    }
  }

  const fetchTodoTrends = async (params = {}) => {
    try {
      const response = await statisticsApi.getTodoTrends(params)
      if (response.code === 200) {
        todoTrends.value = response.data
      }
    } catch (error) {
      console.error('Failed to fetch todo trends:', error)
      ElMessage.error('获取任务趋势失败')
    }
  }

  const fetchArticleTrends = async (params = {}) => {
    try {
      const response = await statisticsApi.getArticleTrends(params)
      if (response.code === 200) {
        articleTrends.value = response.data
      }
    } catch (error) {
      console.error('Failed to fetch article trends:', error)
      ElMessage.error('获取文章趋势失败')
    }
  }

  const refreshAll = async () => {
    await Promise.all([
      fetchStatistics(),
      fetchTrends('todo'),
      fetchTrends('article')
    ])
  }

  return {
    // 状态
    todoStats,
    articleStats,
    todoTrends,
    articleTrends,
    loading,

    // 计算属性
    allStats,

    // 方法
    fetchStatistics,
    fetchTodoStats,
    fetchArticleStats,
    fetchTrends,
    fetchTodoTrends,
    fetchArticleTrends,
    refreshAll
  }
})
