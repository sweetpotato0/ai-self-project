<template>
  <div class="article-detail-container">
    <!-- 返回按钮 -->
    <div class="back-section">
      <div class="back-button" @click="goBack">
        <div class="back-icon">
          <el-icon><ArrowLeft /></el-icon>
        </div>
        <div class="back-content">
          <span class="back-text">返回</span>
          <span class="back-subtext">文章列表</span>
        </div>
      </div>
    </div>

    <!-- 调试信息 -->
    <!-- <div v-if="isDev" class="debug-info">
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
    </div> -->

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

      <!-- 阅读设置 -->
      <div class="reading-controls">
        <el-button-group>
          <el-button
            :type="isWideMode ? '' : 'primary'"
            @click="isWideMode = false"
            size="small"
          >
            <el-icon><Monitor /></el-icon>
            标准模式
          </el-button>
          <el-button
            :type="isWideMode ? 'primary' : ''"
            @click="isWideMode = true"
            size="small"
          >
            <el-icon><FullScreen /></el-icon>
            宽屏模式
          </el-button>
        </el-button-group>
        <el-button @click="showToc = !showToc" size="small">
          <el-icon><Menu /></el-icon>
          {{ showToc ? '隐藏目录' : '显示目录' }}
        </el-button>
      </div>

      <!-- 文章正文 -->
      <div class="article-body" :class="{ 'wide-mode': isWideMode }">
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

      <!-- 悬浮式目录 -->
      <transition name="toc-fade">
        <div
          v-if="showToc && tocItems.length > 0"
          class="floating-toc"
          :class="{
            'toc-collapsed': tocCollapsed,
            'toc-wide': isWideMode
          }"
        >
          <!-- 目录标题栏 -->
          <div class="toc-header" @click="tocCollapsed = !tocCollapsed">
            <div class="toc-title">
              <el-icon><Menu /></el-icon>
              <span v-if="!tocCollapsed">目录</span>
            </div>
            <el-icon class="toc-toggle" :class="{ 'rotated': tocCollapsed }">
              <ArrowDown />
            </el-icon>
          </div>

          <!-- 目录内容 -->
          <div class="toc-content" v-show="!tocCollapsed">
            <div class="toc-progress">
              <div class="progress-bar" :style="{ height: readingProgress + '%' }"></div>
            </div>

            <ul class="toc-list">
              <li
                v-for="item in tocItems"
                :key="item.anchor"
                :class="`toc-level-${item.level}`"
                class="toc-item"
              >
                <a
                  :href="`#${item.anchor}`"
                  @click="scrollToHeading($event, item.anchor)"
                  :class="{ 'active': activeHeading === item.anchor }"
                  class="toc-link"
                >
                  <span class="toc-dot"></span>
                  <span class="toc-text">{{ item.text }}</span>
                </a>
              </li>
            </ul>
          </div>
        </div>
      </transition>

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
import { ArrowDown, ArrowLeft, Edit, FullScreen, Loading, Menu, Monitor, Share, Star, User, View } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import 'highlight.js/styles/github.css'
import MarkdownIt from 'markdown-it'
import highlightPlugin from 'markdown-it-highlightjs'
import { computed, onMounted, onUnmounted, ref } from 'vue'
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
const isWideMode = ref(false)
const showToc = ref(true)
const tocItems = ref([])
const activeHeading = ref('')
const tocCollapsed = ref(false)
const readingProgress = ref(0)

// 创建 markdown-it 实例
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

// 添加代码高亮插件
md.use(highlightPlugin, {
  auto: true,
  code: true
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

    // 提取目录
    extractTocFromContent(rendered)

    return addHeadingAnchors(rendered)
  } catch (error) {
    console.error('Markdown rendering error:', error)
    return '<p>内容渲染错误</p>'
  }
})

// 提取目录
const extractTocFromContent = (htmlContent) => {
  const parser = new DOMParser()
  const doc = parser.parseFromString(htmlContent, 'text/html')
  const headings = doc.querySelectorAll('h1, h2, h3, h4, h5, h6')

  tocItems.value = Array.from(headings).map(heading => {
    const level = parseInt(heading.tagName.charAt(1))
    let text = heading.textContent.trim()

    // 多步骤清理数字编号
    // 1. 移除开头的数字编号格式：1. 、2. 、1.1 、1.2.3 等
    text = text.replace(/^\s*\d+(\.\d+)*\.?\s*/, '')

    // 2. 如果还有重复，再次清理（防止嵌套编号）
    text = text.replace(/^\s*\d+(\.\d+)*\.?\s*/, '')

    // 3. 移除可能的其他编号格式
    text = text.replace(/^[（(]\d+[)）]\s*/, '') // 移除 (1) 格式
    text = text.replace(/^[一二三四五六七八九十]+[、.]\s*/, '') // 移除中文编号


    const anchor = text.toLowerCase()
      .replace(/[^\w\u4e00-\u9fa5\s-]/g, '')
      .replace(/\s+/g, '-')

    return { level, text, anchor }
  })
}

// 为标题添加锚点
const addHeadingAnchors = (htmlContent) => {
  return htmlContent.replace(
    /<h([1-6])([^>]*)>([^<]+)<\/h[1-6]>/g,
    (match, level, attrs, text) => {
      // 多步骤清理数字编号来生成锚点
      let cleanText = text.trim()
      cleanText = cleanText.replace(/^\s*\d+(\.\d+)*\.?\s*/, '')
      cleanText = cleanText.replace(/^\s*\d+(\.\d+)*\.?\s*/, '') // 再次清理
      cleanText = cleanText.replace(/^[（(]\d+[)）]\s*/, '')
      cleanText = cleanText.replace(/^[一二三四五六七八九十]+[、.]\s*/, '')

      const anchor = cleanText.toLowerCase()
        .replace(/[^\w\u4e00-\u9fa5\s-]/g, '')
        .replace(/\s+/g, '-')
      return `<h${level}${attrs} id="${anchor}">${text}</h${level}>`
    }
  )
}

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

const scrollToHeading = (event, anchor) => {
  event.preventDefault()
  const element = document.getElementById(anchor)
  if (element) {
    element.scrollIntoView({
      behavior: 'smooth',
      block: 'start',
      inline: 'nearest'
    })
    activeHeading.value = anchor
  }
}

// 滚动监听
const handleScroll = () => {
  // 计算阅读进度
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop
  const scrollHeight = document.documentElement.scrollHeight - window.innerHeight
  readingProgress.value = Math.min((scrollTop / scrollHeight) * 100, 100)

  // 更新当前活跃的标题
  const headings = tocItems.value
  let currentHeading = ''

  for (let i = headings.length - 1; i >= 0; i--) {
    const element = document.getElementById(headings[i].anchor)
    if (element) {
      const rect = element.getBoundingClientRect()
      if (rect.top <= 100) { // 标题距离顶部100px内认为是当前标题
        currentHeading = headings[i].anchor
        break
      }
    }
  }

  if (currentHeading) {
    activeHeading.value = currentHeading
  }
}

// 生命周期
onMounted(() => {
  console.log('ArticleDetail component mounted')
  console.log('Route params:', route.params)
  console.log('Route name:', route.name)
  loadArticle()

  // 添加滚动监听
  window.addEventListener('scroll', handleScroll)
  handleScroll() // 初始化
})

onUnmounted(() => {
  // 移除滚动监听
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.article-detail-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px 40px;
}

.back-section {
  margin-bottom: 30px;
  position: relative;
}

.back-button {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9) 0%, rgba(248, 250, 252, 0.9) 100%);
  border: 1px solid rgba(226, 232, 240, 0.8);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  position: relative;
  overflow: hidden;
}

.back-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.1), transparent);
  transition: left 0.5s ease;
}

.back-button:hover::before {
  left: 100%;
}

.back-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.15);
  border-color: rgba(102, 126, 234, 0.3);
}

.back-button:active {
  transform: translateY(0);
  transition: all 0.1s ease;
}

.back-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
  font-size: 16px;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.back-button:hover .back-icon {
  transform: translateX(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.back-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.back-text {
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  line-height: 1.2;
}

.back-subtext {
  font-size: 12px;
  color: #718096;
  line-height: 1;
  opacity: 0.8;
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

.reading-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 40px;
  background: white;
  border-bottom: 1px solid #e2e8f0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.article-body {
  padding: 60px 40px;
  background: white;
  position: relative;
}

.article-body.wide-mode {
  max-width: none;
  padding: 60px 60px;
}

.article-body.wide-mode .content-wrapper {
  max-width: 1400px;
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

/* 悬浮目录样式 */
.floating-toc {
  position: fixed;
  right: 30px;
  top: 50%;
  transform: translateY(-50%);
  width: 300px;
  max-height: 70vh;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.floating-toc.toc-collapsed {
  width: 60px;
  height: 60px;
}

.floating-toc.toc-wide {
  right: 50px;
  width: 320px;
}

.toc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 16px 16px 0 0;
}

.toc-collapsed .toc-header {
  border-radius: 16px;
  justify-content: center;
}

.toc-header:hover {
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
}

.toc-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 14px;
}

.toc-toggle {
  transition: transform 0.3s ease;
  opacity: 0.8;
}

.toc-toggle:hover {
  opacity: 1;
}

.toc-toggle.rotated {
  transform: rotate(180deg);
}

.toc-content {
  display: flex;
  max-height: calc(70vh - 60px);
  overflow: hidden;
}

.toc-progress {
  width: 4px;
  background: rgba(0, 0, 0, 0.05);
  position: relative;
  flex-shrink: 0;
}

.progress-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  border-radius: 0 2px 2px 0;
  transition: height 0.3s ease;
}

.toc-list {
  flex: 1;
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  padding: 20px;
}

.toc-list::-webkit-scrollbar {
  width: 6px;
}

.toc-list::-webkit-scrollbar-track {
  background: transparent;
}

.toc-list::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

.toc-list::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}

.toc-item {
  margin-bottom: 2px;
}

.toc-link {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #6c757d;
  text-decoration: none;
  padding: 8px 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-size: 13px;
  line-height: 1.4;
  position: relative;
  word-break: break-word;
}

.toc-link:hover {
  color: #667eea;
  background: rgba(102, 126, 234, 0.08);
  transform: translateX(4px);
}

.toc-link.active {
  color: #667eea;
  background: rgba(102, 126, 234, 0.12);
  font-weight: 600;
  transform: translateX(4px);
}

.toc-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.4;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.toc-link:hover .toc-dot,
.toc-link.active .toc-dot {
  opacity: 1;
  transform: scale(1.5);
}

.toc-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 不同层级的样式 */
.toc-level-1 .toc-link {
  font-weight: 600;
  font-size: 14px;
}

.toc-level-1 .toc-dot {
  width: 8px;
  height: 8px;
}

.toc-level-2 .toc-link {
  margin-left: 12px;
  font-size: 13px;
}

.toc-level-3 .toc-link {
  margin-left: 24px;
  font-size: 12px;
  opacity: 0.8;
}

.toc-level-4 .toc-link {
  margin-left: 36px;
  font-size: 12px;
  opacity: 0.7;
}

.toc-level-5 .toc-link,
.toc-level-6 .toc-link {
  margin-left: 48px;
  font-size: 11px;
  opacity: 0.6;
}

/* 动画效果 */
.toc-fade-enter-active,
.toc-fade-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.toc-fade-enter-from,
.toc-fade-leave-to {
  opacity: 0;
  transform: translateY(-50%) translateX(100px) scale(0.8);
}

.content-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
}

.markdown-content {
  line-height: 1.8;
  font-size: 15px;
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
}

.markdown-content :deep(pre) {
  background: #f6f8fa;
  border: none;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 16px 0;
  font-size: 13px;
  line-height: 1.6;
  font-family: 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  position: relative;
}

.markdown-content :deep(pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
  border-radius: 0;
  font-size: inherit;
  display: block;
  overflow-x: auto;
}

/* 代码高亮样式 - 与编辑器保持一致 */
.markdown-content :deep(.hljs) {
  display: block;
  overflow-x: auto;
  padding: 0;
  color: #24292e;
  background: transparent;
}

.markdown-content :deep(.hljs-comment),
.markdown-content :deep(.hljs-quote) {
  color: #6a737d;
  font-style: italic;
}

.markdown-content :deep(.hljs-keyword),
.markdown-content :deep(.hljs-selector-tag),
.markdown-content :deep(.hljs-subst) {
  color: #d73a49;
  font-weight: bold;
}

.markdown-content :deep(.hljs-number),
.markdown-content :deep(.hljs-literal),
.markdown-content :deep(.hljs-variable),
.markdown-content :deep(.hljs-template-variable),
.markdown-content :deep(.hljs-tag .hljs-attr) {
  color: #005cc5;
}

.markdown-content :deep(.hljs-string),
.markdown-content :deep(.hljs-doctag) {
  color: #032f62;
}

.markdown-content :deep(.hljs-title),
.markdown-content :deep(.hljs-section),
.markdown-content :deep(.hljs-selector-id) {
  color: #6f42c1;
  font-weight: bold;
}

.markdown-content :deep(.hljs-class .hljs-title) {
  color: #6f42c1;
}

.markdown-content :deep(.hljs-tag),
.markdown-content :deep(.hljs-name),
.markdown-content :deep(.hljs-attribute) {
  color: #22863a;
  font-weight: normal;
}

.markdown-content :deep(.hljs-regexp),
.markdown-content :deep(.hljs-link) {
  color: #032f62;
}

.markdown-content :deep(.hljs-symbol),
.markdown-content :deep(.hljs-bullet) {
  color: #e36209;
}

.markdown-content :deep(.hljs-built_in),
.markdown-content :deep(.hljs-builtin-name) {
  color: #005cc5;
}

.markdown-content :deep(.hljs-meta) {
  color: #6a737d;
}

.markdown-content :deep(.hljs-deletion) {
  background: #ffeef0;
}

.markdown-content :deep(.hljs-addition) {
  background: #f0fff4;
}

.markdown-content :deep(.hljs-emphasis) {
  font-style: italic;
}

.markdown-content :deep(.hljs-strong) {
  font-weight: bold;
}

.markdown-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 2em 0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: none;
  border-right: 1px solid #e2e8f0;
  border-bottom: 1px solid #e2e8f0;
  padding: 15px;
  text-align: left;
}

.markdown-content :deep(th:last-child),
.markdown-content :deep(td:last-child) {
  border-right: none;
}

.markdown-content :deep(tr:last-child td) {
  border-bottom: none;
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

@media (max-width: 1200px) {
  .article-body.wide-mode {
    padding: 40px;
  }

  .floating-toc {
    right: 20px;
    width: 280px;
  }

  .floating-toc.toc-wide {
    right: 30px;
    width: 300px;
  }
}

@media (max-width: 768px) {
  .article-detail-container {
    padding: 10px;
  }

  .back-section {
    margin-bottom: 20px;
  }

  .back-button {
    padding: 10px 16px;
    gap: 10px;
  }

  .back-icon {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }

  .back-text {
    font-size: 14px;
  }

  .back-subtext {
    font-size: 11px;
  }

  .article-header {
    padding: 20px;
  }

  .article-title {
    font-size: 1.8rem;
  }

  .reading-controls {
    padding: 15px 20px;
    flex-direction: column;
    gap: 15px;
  }

  .article-body {
    padding: 20px;
  }

  .article-body.wide-mode {
    padding: 20px;
  }

  .content-wrapper {
    max-width: 100%;
  }

  .article-stats {
    flex-direction: column;
    gap: 15px;
  }

  .article-actions {
    padding: 20px;
    flex-direction: column;
  }

  /* 移动端悬浮目录优化 */
  .floating-toc {
    right: 10px;
    width: 280px;
    max-height: 60vh;
  }

  .floating-toc.toc-collapsed {
    width: 50px;
    height: 50px;
  }

  .floating-toc.toc-wide {
    right: 10px;
    width: 300px;
  }

  .toc-header {
    padding: 12px 16px;
  }

  .toc-list {
    padding: 15px;
  }

  .toc-link {
    font-size: 12px;
    padding: 6px 8px;
  }

  .toc-level-1 .toc-link {
    font-size: 13px;
  }
}
</style>
