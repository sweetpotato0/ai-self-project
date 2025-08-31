<template>
  <div class="articles-container">
    <div class="articles-header">
      <div class="header-left">
        <h2 class="page-title">个人文章</h2>
        <p class="page-subtitle">管理您的个人文章和博客</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          新建文章
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon total">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ stats.total || 0 }}</div>
            <div class="stat-label">总文章</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon published">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ stats.published || 0 }}</div>
            <div class="stat-label">已发布</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon draft">
            <el-icon><Edit /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ stats.draft || 0 }}</div>
            <div class="stat-label">草稿</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon views">
            <el-icon><View /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ stats.total_views || 0 }}</div>
            <div class="stat-label">总浏览量</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 筛选和搜索 -->
    <div class="filter-section">
      <div class="filter-row">
        <div class="filter-left">
          <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="handleFilterChange" style="width: 120px;">
            <el-option label="全部状态" value="" />
            <el-option label="已发布" value="published" />
            <el-option label="草稿" value="draft" />
            <el-option label="已归档" value="archived" />
          </el-select>
          
          <el-select v-model="filterTag" placeholder="标签筛选" clearable @change="handleFilterChange" style="width: 140px;">
            <el-option label="全部标签" value="" />
            <el-option
              v-for="tag in availableTags"
              :key="tag"
              :label="tag"
              :value="tag"
            />
          </el-select>
          
          <el-select v-model="sortBy" @change="handleSortChange" style="width: 140px;">
            <el-option label="最新创建" value="created_at:desc" />
            <el-option label="最早创建" value="created_at:asc" />
            <el-option label="最新更新" value="updated_at:desc" />
            <el-option label="最多浏览" value="view_count:desc" />
            <el-option label="最多点赞" value="like_count:desc" />
            <el-option label="标题排序" value="title:asc" />
          </el-select>
        </div>
        
        <div class="filter-right">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文章标题、内容..."
            prefix-icon="Search"
            clearable
            @keyup.enter="handleSearch"
            @clear="handleSearchClear"
            style="width: 280px; margin-right: 10px;"
          />
          <el-button type="primary" @click="handleSearch" :loading="searchLoading">
            搜索
          </el-button>
          <el-button @click="handleResetFilters">
            重置
          </el-button>
        </div>
      </div>
    </div>

    <!-- 文章列表 -->
    <div class="articles-list">
      <el-card v-if="articleStore.loading" class="loading-card">
        <div class="loading-content">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>加载中...</span>
        </div>
      </el-card>

      <div v-else-if="filteredArticles.length === 0" class="empty-state">
        <el-empty description="暂无文章" />
      </div>

      <div v-else class="articles-grid">
        <el-card
          v-for="article in filteredArticles"
          :key="article.id"
          class="article-card"
          shadow="hover"
          @click="viewArticle(article)"
        >
          <div class="article-cover" v-if="article.cover_image">
            <img :src="article.cover_image" :alt="article.title" />
          </div>
          <div class="article-content">
            <h3 class="article-title">{{ article.title }}</h3>
            <p class="article-summary">{{ getArticleSummary(article) }}</p>
            <div class="article-meta">
              <el-tag :type="getStatusType(article.status)" size="small">
                {{ getStatusText(article.status) }}
              </el-tag>
              <span class="article-date">{{ formatDate(article.created_at) }}</span>
            </div>
            <div class="article-stats">
              <span class="stat-item">
                <el-icon><View /></el-icon>
                {{ article.view_count }}
              </span>
              <span class="stat-item">
                <el-icon><Star /></el-icon>
                {{ article.like_count }}
              </span>
            </div>
          </div>
          <div class="article-actions">
            <el-button type="primary" size="small" @click.stop="editArticle(article)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click.stop="deleteArticle(article)">
              删除
            </el-button>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination-section" v-if="articleStore.total > 0">
      <el-pagination
        v-model:current-page="articleStore.page"
        v-model:page-size="articleStore.limit"
        :total="articleStore.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 创建/编辑文章对话框 -->
    <ArticleDialog
      v-model="showCreateDialog"
      :article="currentArticle"
      @success="handleArticleSuccess"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useArticleStore } from '@/stores/article'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Document, Check, Edit, View, Star, Search, Loading } from '@element-plus/icons-vue'
import ArticleDialog from '@/components/article/ArticleDialog.vue'

const router = useRouter()

const articleStore = useArticleStore()

// 响应式数据
const showCreateDialog = ref(false)
const currentArticle = ref(null)
const filterStatus = ref('')
const filterTag = ref('')
const sortBy = ref('created_at:desc')
const searchKeyword = ref('')
const searchLoading = ref(false)
const availableTags = ref([])

// 统计数据
const stats = computed(() => articleStore.stats)

// 过滤后的文章列表（现在主要由后端处理）
const filteredArticles = computed(() => {
  return articleStore.articles || []
})

// 方法
const loadData = async () => {
  try {
    const params = {
      status: filterStatus.value,
      tag: filterTag.value,
      sort: sortBy.value,
      search: searchKeyword.value
    }
    
    await Promise.all([
      articleStore.fetchArticles(params),
      articleStore.fetchStats(),
      loadAvailableTags()
    ])
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
}

const loadAvailableTags = async () => {
  try {
    // 从文章中提取所有可用标签
    const allTags = new Set()
    articleStore.articles?.forEach(article => {
      if (article.tags) {
        try {
          const tags = JSON.parse(article.tags)
          if (Array.isArray(tags)) {
            tags.forEach(tag => allTags.add(tag))
          }
        } catch (e) {
          // 忽略解析错误
        }
      }
    })
    availableTags.value = Array.from(allTags).sort()
  } catch (error) {
    console.error('加载标签失败:', error)
  }
}

const handleFilterChange = () => {
  // 筛选改变时重新加载数据
  articleStore.page = 1 // 重置到第一页
  loadData()
}

const handleSortChange = () => {
  // 排序改变时重新加载数据
  articleStore.page = 1 // 重置到第一页
  loadData()
}

const handleSearch = async () => {
  if (searchLoading.value) return
  
  try {
    searchLoading.value = true
    articleStore.page = 1 // 重置到第一页
    await loadData()
  } catch (error) {
    ElMessage.error('搜索失败')
  } finally {
    searchLoading.value = false
  }
}

const handleSearchClear = () => {
  searchKeyword.value = ''
  loadData()
}

const handleResetFilters = () => {
  filterStatus.value = ''
  filterTag.value = ''
  sortBy.value = 'created_at:desc'
  searchKeyword.value = ''
  articleStore.page = 1
  loadData()
}

const handleSizeChange = (size) => {
  articleStore.limit = size
  loadData()
}

const handleCurrentChange = (page) => {
  articleStore.page = page
  loadData()
}

const viewArticle = (article) => {
  // 跳转到文章详情页面
  router.push(`/articles/${article.id}`)
}

const editArticle = (article) => {
  currentArticle.value = article
  showCreateDialog.value = true
}

const deleteArticle = async (article) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文章"${article.title}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await articleStore.deleteArticle(article.id)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleArticleSuccess = () => {
  showCreateDialog.value = false
  currentArticle.value = null
  loadData()
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
  return date.toLocaleDateString('zh-CN')
}

const getArticleSummary = (article) => {
  if (article.summary && article.summary.trim()) {
    return article.summary
  }
  
  if (article.content && article.content.trim()) {
    // 去除 Markdown 格式标记
    let plainText = article.content
      // 去除标题标记
      .replace(/#{1,6}\s+/g, '')
      // 去除粗体和斜体
      .replace(/\*\*([^*]+)\*\*/g, '$1')
      .replace(/\*([^*]+)\*/g, '$1')
      // 去除代码块
      .replace(/```[\s\S]*?```/g, '')
      // 去除行内代码
      .replace(/`([^`]+)`/g, '$1')
      // 去除链接
      .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
      // 去除图片
      .replace(/!\[([^\]]*)\]\([^)]+\)/g, '')
      // 去除换行符和多余空格
      .replace(/\n+/g, ' ')
      .replace(/\s+/g, ' ')
      .trim()
    
    // 截取前80个字符
    if (plainText.length > 80) {
      return plainText.substring(0, 80) + '...'
    }
    
    return plainText || '暂无摘要'
  }
  
  return '暂无摘要'
}

// 生命周期
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.articles-container {
  padding: 20px;
}

.articles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  color: #666;
  margin: 0;
}

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  border: none;
  border-radius: 12px;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.published {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.draft {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.stat-icon.views {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

/* 筛选和搜索 */
.filter-section {
  margin-bottom: 20px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
}

.filter-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 15px;
}

.filter-left {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.filter-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* 文章列表 */
.articles-list {
  margin-bottom: 30px;
}

.loading-card {
  text-align: center;
  padding: 40px;
}

.loading-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #666;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.articles-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: flex-start;
}

.articles-grid .article-card {
  flex: 0 0 calc(20% - 13px);
  max-width: calc(20% - 13px);
  min-width: 240px;
}

/* 当空间不足以放置5列时，自动调整为4列 */
@media (max-width: 1600px) {
  .articles-grid .article-card {
    flex: 0 0 calc(25% - 12px);
    max-width: calc(25% - 12px);
  }
}

/* 当空间不足以放置4列时，自动调整为3列 */
@media (max-width: 1200px) {
  .articles-grid .article-card {
    flex: 0 0 calc(33.333% - 11px);
    max-width: calc(33.333% - 11px);
  }
}

/* 当空间不足以放置3列时，自动调整为2列 */
@media (max-width: 900px) {
  .articles-grid .article-card {
    flex: 0 0 calc(50% - 8px);
    max-width: calc(50% - 8px);
  }
}

/* 当空间不足以放置2列时，自动调整为1列 */
@media (max-width: 800px) {
  .articles-grid .article-card {
    flex: 0 0 100%;
    max-width: 100%;
    min-width: unset;
  }
}

.article-card {
  border: none;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.article-cover {
  height: 120px;
  overflow: hidden;
  border-radius: 8px 8px 0 0;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-content {
  padding: 15px;
}

.article-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: #333;
  line-height: 1.3;
}

.article-summary {
  color: #666;
  font-size: 13px;
  line-height: 1.5;
  margin: 0 0 10px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.article-date {
  font-size: 12px;
  color: #999;
}

.article-stats {
  display: flex;
  gap: 15px;
  margin-bottom: 10px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #666;
}

.article-actions {
  display: flex;
  gap: 8px;
  padding: 0 15px 15px;
}

/* 分页 */
.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .articles-container {
    padding: 10px;
  }

  .articles-header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .filter-section {
    padding: 15px;
  }

  .filter-row {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }

  .filter-left {
    justify-content: center;
  }

  .filter-right {
    justify-content: center;
  }

  .articles-grid {
    grid-template-columns: 1fr;
  }
}
</style>
