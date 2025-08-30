import { articlesApi } from '@/features/articles/api'
import { defineStore } from 'pinia'

export const useArticleStore = defineStore('article', {
  state: () => ({
    articles: [],
    currentArticle: null,
    stats: {},
    loading: false,
    total: 0,
    page: 1,
    limit: 10
  }),

  getters: {
    publishedArticles: (state) => state.articles.filter(article => article.status === 'published'),
    draftArticles: (state) => state.articles.filter(article => article.status === 'draft'),
    archivedArticles: (state) => state.articles.filter(article => article.status === 'archived'),

    // 获取文章统计信息
    getStats: (state) => {
      const total = state.articles.length
      const published = state.articles.filter(article => article.status === 'published').length
      const draft = state.articles.filter(article => article.status === 'draft').length
      const archived = state.articles.filter(article => article.status === 'archived').length
      const totalViews = state.articles.reduce((sum, article) => sum + (article.view_count || 0), 0)
      const totalLikes = state.articles.reduce((sum, article) => sum + (article.like_count || 0), 0)

      return {
        total,
        published,
        draft,
        archived,
        totalViews,
        totalLikes
      }
    }
  },

  actions: {
    // 获取文章列表
    async fetchArticles(params = {}) {
      this.loading = true
      try {
        const response = await articlesApi.getArticles({
          page: this.page,
          limit: this.limit,
          ...params
        })

        if (response.code === 200) {
          this.articles = response.data.articles
          this.total = response.data.total
          this.page = response.data.page
          this.limit = response.data.limit
        }
      } catch (error) {
        console.error('Failed to fetch articles:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // 获取文章统计信息
    async fetchStats() {
      try {
        const response = await articlesApi.getArticleStats()
        if (response.code === 200) {
          this.stats = response.data
        }
      } catch (error) {
        console.error('Failed to fetch article stats:', error)
        throw error
      }
    },

        // 获取单个文章
    async getArticle(id) {
      try {
        console.log('Store: Fetching article with ID:', id)
        console.log('Store: Calling articlesApi.getArticle...')

        const response = await articlesApi.getArticle(id)
        console.log('Store: API response received:', response)
        console.log('Store: Response type:', typeof response)
        console.log('Store: Response keys:', Object.keys(response || {}))

        if (response && response.code === 200) {
          this.currentArticle = response.data
          console.log('Store: Article data set successfully:', response.data)
          return response.data
        } else {
          console.error('Store: API returned error code:', response?.code)
          console.error('Store: API error message:', response?.message)
          throw new Error(response?.message || '获取文章失败')
        }
      } catch (error) {
        console.error('Store: Failed to fetch article:', error)
        console.error('Store: Error details:', error.response || error.message)
        throw error
      }
    },

    // 增加文章浏览量
    async incrementViewCount(id) {
      try {
        // 这里可以调用后端API来增加浏览量
        // 目前先在前端更新
        const article = this.articles.find(a => a.id === id)
        if (article) {
          article.view_count = (article.view_count || 0) + 1
        }
        if (this.currentArticle && this.currentArticle.id === id) {
          this.currentArticle.view_count = (this.currentArticle.view_count || 0) + 1
        }
      } catch (error) {
        console.error('Failed to increment view count:', error)
      }
    },

    // 点赞文章
    async likeArticle(id) {
      try {
        console.log('Store: Liking article with ID:', id)
        const response = await articlesApi.likeArticle(id)
        console.log('Store: Like API response:', response)
        
        if (response && response.code === 200) {
          // 更新本地数据
          const article = this.articles.find(a => a.id === id)
          if (article) {
            article.like_count = (article.like_count || 0) + 1
          }
          if (this.currentArticle && this.currentArticle.id === id) {
            this.currentArticle.like_count = (this.currentArticle.like_count || 0) + 1
          }
          console.log('Store: Article like count updated locally')
          return true
        } else {
          throw new Error(response?.message || '点赞失败')
        }
      } catch (error) {
        console.error('Failed to like article:', error)
        throw error
      }
    },

    // 取消点赞文章
    async unlikeArticle(id) {
      try {
        console.log('Store: Unliking article with ID:', id)
        const response = await articlesApi.unlikeArticle(id)
        console.log('Store: Unlike API response:', response)
        
        if (response && response.code === 200) {
          // 更新本地数据
          const article = this.articles.find(a => a.id === id)
          if (article) {
            article.like_count = Math.max(0, (article.like_count || 0) - 1)
          }
          if (this.currentArticle && this.currentArticle.id === id) {
            this.currentArticle.like_count = Math.max(0, (this.currentArticle.like_count || 0) - 1)
          }
          console.log('Store: Article like count updated locally')
          return true
        } else {
          throw new Error(response?.message || '取消点赞失败')
        }
      } catch (error) {
        console.error('Failed to unlike article:', error)
        throw error
      }
    },

    // 创建文章
    async createArticle(articleData) {
      try {
        const response = await articlesApi.createArticle(articleData)
        if (response.code === 200) {
          // 重新获取文章列表
          await this.fetchArticles()
          return response.data
        }
      } catch (error) {
        console.error('Failed to create article:', error)
        throw error
      }
    },

    // 更新文章
    async updateArticle(id, articleData) {
      try {
        const response = await articlesApi.updateArticle(id, articleData)
        if (response.code === 200) {
          // 更新当前文章
          if (this.currentArticle && this.currentArticle.id === id) {
            this.currentArticle = response.data
          }
          // 更新列表中的文章
          const index = this.articles.findIndex(article => article.id === id)
          if (index !== -1) {
            this.articles[index] = response.data
          }
          return response.data
        }
      } catch (error) {
        console.error('Failed to update article:', error)
        throw error
      }
    },

    // 删除文章
    async deleteArticle(id) {
      try {
        const response = await articlesApi.deleteArticle(id)
        if (response.code === 200) {
          // 从列表中移除
          this.articles = this.articles.filter(article => article.id !== id)
          // 如果删除的是当前文章，清空当前文章
          if (this.currentArticle && this.currentArticle.id === id) {
            this.currentArticle = null
          }
          return true
        }
      } catch (error) {
        console.error('Failed to delete article:', error)
        throw error
      }
    },

    // 设置当前文章
    setCurrentArticle(article) {
      this.currentArticle = article
    },

    // 清空当前文章
    clearCurrentArticle() {
      this.currentArticle = null
    },

    // 获取单个文章的别名方法
    async getArticleById(id) {
      return await this.getArticle(id)
    },

    // 重置状态
    reset() {
      this.articles = []
      this.currentArticle = null
      this.stats = {}
      this.loading = false
      this.total = 0
      this.page = 1
      this.limit = 10
    }
  }
})
