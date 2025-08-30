<template>
  <div class="article-detail-container">
    <!-- 返回按钮 -->
    <div class="back-section">
      <el-button @click="goBack" type="text">
        <el-icon><ArrowLeft /></el-icon>
        返回文章列表
      </el-button>
    </div>

    <!-- 调试信息 -->
    <div v-if="isDev" class="debug-info">
      <h3>调试信息</h3>
      <p>路由参数: {{ route.params }}</p>
      <p>文章ID: {{ route.params.id }}</p>
      <p>加载状态: {{ loading }}</p>
      <p>认证状态: {{ authStore.isAuthenticated ? '已登录' : '未登录' }}</p>
      <p>用户信息: {{ authStore.user?.username || '无' }}</p>
      <p>Token存在: {{ hasToken ? '是' : '否' }}</p>
      <p>文章数据: {{ article ? '已加载' : '未加载' }}</p>
      <p>文章内容长度: {{ article?.content?.length || 0 }}</p>
      <p>当前路由名称: {{ route.name }}</p>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-card class="loading-card">
        <div class="loading-content">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>加载中...</span>
        </div>
      </el-card>
    </div>

    <!-- 文章内容 -->
    <div v-else-if="article" class="article-content">
      <!-- 文章头部 -->
      <div class="article-header">
        <div class="article-meta">
          <el-tag :type="getStatusType(article.status)" size="small">
            {{ getStatusText(article.status) }}
          </el-tag>
          <span class="publish-date">{{ formatDate(article.created_at) }}</span>
        </div>
        <h1 class="article-title">{{ article.title }}</h1>
        <p v-if="article.summary" class="article-summary">{{ article.summary }}</p>

        <!-- 封面图片 -->
        <div v-if="article.cover_image" class="article-cover">
          <img :src="article.cover_image" :alt="article.title" />
        </div>

        <!-- 文章统计 -->
        <div class="article-stats">
          <div class="stat-item">
            <el-icon><View /></el-icon>
            <span>{{ article.view_count || 0 }} 次浏览</span>
          </div>
          <div class="stat-item">
            <el-icon><Star /></el-icon>
            <span>{{ article.like_count || 0 }} 次点赞</span>
          </div>
          <div class="stat-item">
            <el-icon><User /></el-icon>
            <span>{{ article.user?.username || '未知作者' }}</span>
          </div>
        </div>

        <!-- 标签 -->
        <div v-if="article.tags" class="article-tags">
          <el-tag
            v-for="tag in parseTags(article.tags)"
            :key="tag"
            type="info"
            size="small"
            class="tag-item"
          >
            {{ tag }}
          </el-tag>
        </div>
      </div>

      <!-- 文章正文 -->
      <div class="article-body">
        <div class="content-wrapper">
          <!-- 调试信息 -->
          <div v-if="!article.content" class="debug-info">
            <p>调试信息：文章内容为空</p>
            <p>Article object: {{ JSON.stringify(article, null, 2) }}</p>
          </div>

          <!-- 使用 markdown-it 渲染 Markdown 内容 -->
          <div v-html="renderedContent" class="markdown-content"></div>

          <!-- 原始内容（调试用） -->
          <details class="debug-section">
            <summary>显示原始内容（调试用）</summary>
            <pre>{{ article.content }}</pre>
          </details>
        </div>
      </div>

      <!-- 文章操作 -->
      <div class="article-actions">
        <el-button type="primary" @click="editArticle">
          <el-icon><Edit /></el-icon>
          编辑文章
        </el-button>
        <el-button @click="likeArticle" :type="isLiked ? 'danger' : 'default'">
          <el-icon><Star /></el-icon>
          {{ isLiked ? '取消点赞' : '点赞' }}
        </el-button>
        <el-button @click="shareArticle">
          <el-icon><Share /></el-icon>
          分享
        </el-button>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else class="error-container">
      <el-card class="error-card">
        <el-empty description="文章不存在或已被删除" />
      </el-card>
    </div>

    <!-- 编辑文章对话框 -->
    <ArticleDialog
      v-model="showEditDialog"
      :article="article"
      @success="handleEditSuccess"
    />
  </div>
</template>

<script setup>
import ArticleDialog from '@/components/article/ArticleDialog.vue'
import { useArticleStore } from '@/stores/article'
import { useAuthStore } from '@/stores/auth'
import { ArrowLeft, Edit, Loading, Share, Star, User, View } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import MarkdownIt from 'markdown-it'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const articleStore = useArticleStore()
const authStore = useAuthStore()

// 响应式数据
const loading = ref(true)
const article = ref(null)
const showEditDialog = ref(false)
const isLiked = ref(false)

// 创建 markdown-it 实例
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

// 计算属性
const isDev = computed(() => import.meta.env.DEV)
const hasToken = computed(() => !!localStorage.getItem('token'))

const renderedContent = computed(() => {
  console.log('Article content:', article.value?.content)
  if (!article.value?.content) {
    console.log('No article content found')
    return '<p>暂无内容</p>'
  }
  try {
    const rendered = md.render(article.value.content)
    console.log('Rendered content:', rendered)
    return rendered
  } catch (error) {
    console.error('Markdown rendering error:', error)
    return '<p>内容渲染错误</p>'
  }
})

// 方法
const loadArticle = async () => {
  try {
    loading.value = true
    const articleId = parseInt(route.params.id)
    console.log('Loading article with ID:', articleId)
    console.log('Route params:', route.params)
    console.log('Auth status:', authStore.isAuthenticated)
    console.log('User:', authStore.user)

    // 检查用户是否已登录
    if (!authStore.isAuthenticated) {
      console.error('User not authenticated')
      ElMessage.error('请先登录')
      router.push('/login')
      return
    }

    // 检查token是否存在
    const token = localStorage.getItem('token')
    if (!token) {
      console.error('No token found')
      ElMessage.error('登录已过期，请重新登录')
      router.push('/login')
      return
    }

    console.log('Token exists:', token.substring(0, 50) + '...')

    // 直接调用API进行测试
    console.log('Calling articleStore.getArticleById...')
    const result = await articleStore.getArticleById(articleId)
    console.log('Article result:', result)

    if (result) {
      article.value = result
      // 设置用户点赞状态
      isLiked.value = result.is_liked_by_user || false
      console.log('Article loaded successfully:', article.value)
      console.log('User like status:', isLiked.value)
    } else {
      console.error('No article data returned')
      ElMessage.error('文章数据为空')
    }

    // 增加浏览量
    await articleStore.incrementViewCount(articleId)
  } catch (error) {
    console.error('Load article error:', error)
    console.error('Error details:', error.response || error.message)

    if (error.response?.status === 401) {
      ElMessage.error('登录已过期，请重新登录')
      router.push('/login')
    } else {
      ElMessage.error('加载文章失败: ' + (error.message || '未知错误'))
    }
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  // 检查当前路由，如果是独立路由则返回文章列表，否则返回仪表盘
  if (route.name === 'article-detail-standalone') {
    router.push('/dashboard/articles')
  } else {
    router.push('/articles')
  }
}

const editArticle = () => {
  showEditDialog.value = true
}

const handleEditSuccess = () => {
  showEditDialog.value = false
  loadArticle() // 重新加载文章
  ElMessage.success('文章更新成功')
}

const likeArticle = async () => {
  try {
    if (isLiked.value) {
      await articleStore.unlikeArticle(article.value.id)
      isLiked.value = false
      ElMessage.success('取消点赞成功')
    } else {
      await articleStore.likeArticle(article.value.id)
      isLiked.value = true
      ElMessage.success('点赞成功')
    }
  } catch (error) {
    console.error('Like article error:', error)
    if (error.response?.data?.message === 'Already liked this article') {
      ElMessage.warning('您已经点赞过这篇文章')
      // 如果后端说已经点赞过了，更新前端状态
      isLiked.value = true
    } else if (error.response?.data?.message === 'Not liked yet') {
      ElMessage.warning('您还没有点赞过这篇文章')
      // 如果后端说还没点赞过，更新前端状态
      isLiked.value = false
    } else {
      ElMessage.error('操作失败: ' + (error.response?.data?.error || error.message))
    }
  }
}

const shareArticle = () => {
  const url = window.location.href
  if (navigator.share) {
    navigator.share({
      title: article.value.title,
      text: article.value.summary,
      url: url
    })
  } else {
    // 复制链接到剪贴板
    navigator.clipboard.writeText(url).then(() => {
      ElMessage.success('链接已复制到剪贴板')
    })
  }
}

const getStatusType = (status) => {
  const types = {
    published: 'success',
    draft: 'warning',
    archived: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    published: '已发布',
    draft: '草稿',
    archived: '已归档'
  }
  return texts[status] || status
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const parseTags = (tagsString) => {
  if (!tagsString) return []
  try {
    return JSON.parse(tagsString)
  } catch {
    return tagsString.split(',').map(tag => tag.trim())
  }
}

// 生命周期
onMounted(() => {
  console.log('ArticleDetail component mounted')
  console.log('Route params:', route.params)
  console.log('Route name:', route.name)
  loadArticle()
})
</script>

<style scoped>
.article-detail-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.back-section {
  margin-bottom: 20px;
}

.loading-container,
.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.loading-card,
.error-card {
  text-align: center;
}

.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.article-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.article-header {
  padding: 60px 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  color: white;
  position: relative;
  overflow: hidden;
}

.article-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="white" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="white" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="white" opacity="0.1"/><circle cx="10" cy="60" r="0.5" fill="white" opacity="0.1"/><circle cx="90" cy="40" r="0.5" fill="white" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
  opacity: 0.3;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
}

.publish-date {
  font-size: 14px;
  opacity: 0.8;
}

.article-title {
  font-size: 3rem;
  font-weight: 800;
  margin: 0 0 30px 0;
  line-height: 1.2;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  position: relative;
  z-index: 1;
}

.article-title::after {
  content: '';
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 100px;
  height: 4px;
  background: linear-gradient(90deg, #fff, transparent);
  border-radius: 2px;
}

.article-summary {
  font-size: 1.1rem;
  line-height: 1.6;
  opacity: 0.9;
  margin-bottom: 30px;
}

.article-cover {
  margin: 30px 0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.article-cover img {
  width: 100%;
  height: auto;
  display: block;
}

.article-stats {
  display: flex;
  gap: 40px;
  margin-bottom: 30px;
  position: relative;
  z-index: 1;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  opacity: 0.9;
  background: rgba(255, 255, 255, 0.1);
  padding: 8px 16px;
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

.stat-item:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.article-tags {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  position: relative;
  z-index: 1;
}

.tag-item {
  margin: 0;
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.tag-item:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.article-body {
  padding: 60px 40px;
  background: white;
  position: relative;
}

.article-body::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, #e2e8f0, transparent);
}

.content-wrapper {
  max-width: 900px;
  margin: 0 auto;
  position: relative;
}

.markdown-content {
  line-height: 1.8;
  font-size: 18px;
  color: #2c3e50;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  margin-top: 2.5em;
  margin-bottom: 1.2em;
  font-weight: 700;
  line-height: 1.3;
  color: #1a202c;
}

.markdown-content :deep(h1) {
  font-size: 2.5rem;
  border-bottom: 3px solid #667eea;
  padding-bottom: 15px;
  margin-top: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.markdown-content :deep(h2) {
  font-size: 2rem;
  color: #2d3748;
  border-left: 4px solid #667eea;
  padding-left: 20px;
  margin-left: -20px;
}

.markdown-content :deep(h3) {
  font-size: 1.5rem;
  color: #4a5568;
  margin-top: 2em;
}

.markdown-content :deep(h4) {
  font-size: 1.25rem;
  color: #718096;
}

.markdown-content :deep(p) {
  margin-bottom: 1.5em;
  text-align: justify;
  line-height: 1.8;
  color: #4a5568;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin-bottom: 1.5em;
  padding-left: 2.5em;
}

.markdown-content :deep(li) {
  margin-bottom: 0.8em;
  line-height: 1.7;
  color: #4a5568;
}

.markdown-content :deep(ul li) {
  position: relative;
}

.markdown-content :deep(ul li::before) {
  content: '•';
  color: #667eea;
  font-weight: bold;
  position: absolute;
  left: -1.5em;
}

.markdown-content :deep(ol li) {
  counter-increment: list-counter;
}

.markdown-content :deep(ol li::before) {
  content: counter(list-counter) '.';
  color: #667eea;
  font-weight: bold;
  position: absolute;
  left: -2em;
}

.markdown-content :deep(blockquote) {
  border-left: 5px solid #667eea;
  padding: 1.5em 2em;
  margin: 2em 0;
  background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
  border-radius: 8px;
  font-style: italic;
  position: relative;
}

.markdown-content :deep(blockquote::before) {
  content: '"';
  font-size: 4rem;
  color: #667eea;
  position: absolute;
  top: -10px;
  left: 10px;
  opacity: 0.3;
}

.markdown-content :deep(code) {
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
  padding: 4px 8px;
  border-radius: 6px;
  font-family: 'JetBrains Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9em;
  color: #e53e3e;
  border: 1px solid #e2e8f0;
}

.markdown-content :deep(pre) {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
  color: #e2e8f0;
  padding: 2em;
  border-radius: 12px;
  overflow-x: auto;
  margin: 2em 0;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  position: relative;
}

.markdown-content :deep(pre::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #667eea, #764ba2, #f093fb);
  border-radius: 12px 12px 0 0;
}

.markdown-content :deep(pre code) {
  background: none;
  padding: 0;
  color: inherit;
  font-size: 0.95em;
  line-height: 1.6;
}

.markdown-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 2em 0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid #e2e8f0;
  padding: 15px;
  text-align: left;
}

.markdown-content :deep(th) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.9em;
  letter-spacing: 0.5px;
}

.markdown-content :deep(td) {
  background: white;
  color: #4a5568;
}

.markdown-content :deep(tr:nth-child(even) td) {
  background: #f7fafc;
}

.markdown-content :deep(strong) {
  color: #2d3748;
  font-weight: 700;
}

.markdown-content :deep(em) {
  color: #667eea;
  font-style: italic;
}

.markdown-content :deep(a) {
  color: #667eea;
  text-decoration: none;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
}

.markdown-content :deep(a:hover) {
  border-bottom-color: #667eea;
  color: #764ba2;
}

.markdown-content :deep(hr) {
  border: none;
  height: 2px;
  background: linear-gradient(90deg, transparent, #667eea, transparent);
  margin: 3em 0;
}

.article-actions {
  padding: 40px;
  background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
  border-top: 1px solid #e2e8f0;
  display: flex;
  gap: 20px;
  justify-content: center;
  position: relative;
}

.article-actions::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, #667eea, transparent);
}

.debug-info {
  background: #fff3cd;
  border: 1px solid #ffeaa7;
  border-radius: 4px;
  padding: 15px;
  margin-bottom: 20px;
  color: #856404;
}

.debug-section {
  margin-top: 30px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 4px;
  border: 1px solid #dee2e6;
}

.debug-section summary {
  cursor: pointer;
  font-weight: 600;
  color: #495057;
  margin-bottom: 10px;
}

.debug-section pre {
  background: #e9ecef;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-word;
}

@media (max-width: 768px) {
  .article-detail-container {
    padding: 10px;
  }

  .article-header {
    padding: 20px;
  }

  .article-title {
    font-size: 1.8rem;
  }

  .article-body {
    padding: 20px;
  }

  .article-stats {
    flex-direction: column;
    gap: 15px;
  }

  .article-actions {
    padding: 20px;
    flex-direction: column;
  }
}
</style>
