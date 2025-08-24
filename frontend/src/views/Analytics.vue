<template>
  <div class="analytics-page">
    <div class="page-header">
      <h1>数据分析</h1>
      <p>查看您的任务和文章统计信息</p>
    </div>

    <div class="analytics-grid">
      <!-- 任务统计卡片 -->
      <el-card class="stat-card">
        <template #header>
          <div class="card-header">
            <span>任务统计</span>
            <el-button type="primary" size="small" @click="refreshTodoStats">刷新</el-button>
          </div>
        </template>
        <div class="stat-content">
          <div class="stat-item">
            <div class="stat-number">{{ todoStats.total }}</div>
            <div class="stat-label">总任务数</div>
          </div>
          <div class="stat-item">
            <div class="stat-number success">{{ todoStats.completed }}</div>
            <div class="stat-label">已完成</div>
          </div>
          <div class="stat-item">
            <div class="stat-number warning">{{ todoStats.inProgress }}</div>
            <div class="stat-label">进行中</div>
          </div>
          <div class="stat-item">
            <div class="stat-number danger">{{ todoStats.pending }}</div>
            <div class="stat-label">待处理</div>
          </div>
        </div>
      </el-card>

      <!-- 文章统计卡片 -->
      <el-card class="stat-card">
        <template #header>
          <div class="card-header">
            <span>文章统计</span>
            <el-button type="primary" size="small" @click="refreshArticleStats">刷新</el-button>
          </div>
        </template>
        <div class="stat-content">
          <div class="stat-item">
            <div class="stat-number">{{ articleStats.total }}</div>
            <div class="stat-label">总文章数</div>
          </div>
          <div class="stat-item">
            <div class="stat-number success">{{ articleStats.published }}</div>
            <div class="stat-label">已发布</div>
          </div>
          <div class="stat-item">
            <div class="stat-number warning">{{ articleStats.draft }}</div>
            <div class="stat-label">草稿</div>
          </div>
          <div class="stat-item">
            <div class="stat-number info">{{ articleStats.totalViews }}</div>
            <div class="stat-label">总浏览量</div>
          </div>
        </div>
      </el-card>

      <!-- 任务完成趋势 -->
      <el-card class="chart-card">
        <template #header>
          <div class="card-header">
            <span>任务完成趋势</span>
            <el-select v-model="todoTimeRange" size="small" style="width: 120px">
              <el-option label="最近7天" value="7" />
              <el-option label="最近30天" value="30" />
              <el-option label="最近90天" value="90" />
            </el-select>
          </div>
        </template>
        <div class="chart-container">
          <div v-if="todoTrendData.length === 0" class="chart-placeholder">
            <el-icon size="48" color="#909399"><TrendCharts /></el-icon>
            <p>暂无数据</p>
            <p class="chart-desc">暂无任务完成趋势数据</p>
          </div>
          <div v-else class="trend-chart">

            <div class="chart-scroll-container" ref="todoChartContainer">
              <div class="chart-bars">
                <div
                  v-for="(item, index) in todoTrendData"
                  :key="index"
                  class="chart-bar"
                  :class="{ 'zero-bar': (Number(item.count) || 0) === 0 }"
                  :style="{ height: getBarHeight(item.count, todoTrendData) }"
                  :title="`${item.date}: ${item.count} 个任务`"
                  @click="handleBarClick('todo', item, index)"
                >
                  <div v-if="(Number(item.count) || 0) > 0" class="bar-value">{{ item.count }}</div>
                </div>
              </div>
              <div class="chart-labels">
                <span
                  v-for="(item, index) in todoTrendData"
                  :key="index"
                  class="chart-label"
                >
                  {{ formatChartDate(item.date, index, todoTrendData.length) }}
                </span>
              </div>
            </div>
            <!-- 移除滚动提示，避免遮挡坐标 -->
          </div>
        </div>
      </el-card>

      <!-- 文章浏览量趋势 -->
      <el-card class="chart-card">
        <template #header>
          <div class="card-header">
            <span>文章浏览量趋势</span>
            <el-select v-model="articleTimeRange" size="small" style="width: 120px">
              <el-option label="最近7天" value="7" />
              <el-option label="最近30天" value="30" />
              <el-option label="最近90天" value="90" />
            </el-select>
          </div>
        </template>
        <div class="chart-container">
          <div v-if="articleTrendData.length === 0" class="chart-placeholder">
            <el-icon size="48" color="#909399"><TrendCharts /></el-icon>
            <p>暂无数据</p>
            <p class="chart-desc">暂无文章浏览量趋势数据</p>
          </div>
          <div v-else class="trend-chart">

            <div class="chart-scroll-container" ref="articleChartContainer">
              <div class="chart-bars">
                <div
                  v-for="(item, index) in articleTrendData"
                  :key="index"
                  class="chart-bar"
                  :class="{ 'zero-bar': (Number(item.count) || 0) === 0 }"
                  :style="{ height: getBarHeight(item.count, articleTrendData) }"
                  :title="`${item.date}: ${item.count || 0} 次浏览`"
                  @click="handleBarClick('article', item, index)"
                >
                  <div v-if="(Number(item.count) || 0) > 0" class="bar-value">{{ item.count || 0 }}</div>
                </div>
              </div>
              <div class="chart-labels">
                <span
                  v-for="(item, index) in articleTrendData"
                  :key="index"
                  class="chart-label"
                >
                  {{ formatChartDate(item.date, index, articleTrendData.length) }}
                </span>
              </div>
            </div>
            <!-- 移除滚动提示，避免遮挡坐标 -->
          </div>
        </div>
      </el-card>

      <!-- 热门文章 -->
      <el-card class="list-card">
        <template #header>
          <div class="card-header">
            <span>热门文章</span>
            <el-button type="primary" size="small" @click="viewAllArticles">查看全部</el-button>
          </div>
        </template>
        <div class="article-list">
          <div v-for="article in popularArticles" :key="article.id" class="article-item">
            <div class="article-info">
              <h4>{{ article.title }}</h4>
              <p>{{ article.summary }}</p>
              <div class="article-meta">
                <span>浏览量: {{ article.view_count || 0 }}</span>
                <span>点赞: {{ article.like_count || 0 }}</span>
                <span>{{ formatDate(article.created_at) }}</span>
              </div>
            </div>
            <el-button type="primary" size="small" @click="viewArticle(article.id)">查看</el-button>
          </div>
        </div>
      </el-card>

      <!-- 最近任务 -->
      <el-card class="list-card">
        <template #header>
          <div class="card-header">
            <span>最近任务</span>
            <el-button type="primary" size="small" @click="viewAllTodos">查看全部</el-button>
          </div>
        </template>
        <div class="todo-list">
          <div v-for="todo in recentTodos" :key="todo.id" class="todo-item">
            <div class="todo-info">
              <h4>{{ todo.title }}</h4>
              <p>{{ todo.description }}</p>
              <div class="todo-meta">
                <el-tag :type="getStatusType(todo.status)" size="small">
                  {{ getStatusText(todo.status) }}
                </el-tag>
                <span>{{ formatDate(todo.created_at) }}</span>
              </div>
            </div>
            <el-button type="primary" size="small" @click="viewTodo(todo.id)">查看</el-button>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useStatisticsStore } from '@/stores/statistics'
import { useTodoStore } from '@/stores/todo'
import { useArticleStore } from '@/stores/article'
import { ElMessage } from 'element-plus'
import { TrendCharts, ArrowRight } from '@element-plus/icons-vue'

const router = useRouter()
const statisticsStore = useStatisticsStore()
const todoStore = useTodoStore()
const articleStore = useArticleStore()

// 响应式数据
const todoTimeRange = ref('7')
const articleTimeRange = ref('7')

// 计算属性
const todoStats = computed(() => statisticsStore.todoStats)
const articleStats = computed(() => statisticsStore.articleStats)

const popularArticles = computed(() => {
  if (!articleStore.articles || !Array.isArray(articleStore.articles)) {
    return []
  }
  return articleStore.articles
    .sort((a, b) => (b.view_count || 0) - (a.view_count || 0))
    .slice(0, 5)
})

const recentTodos = computed(() => {
  if (!todoStore.todos || !Array.isArray(todoStore.todos)) {
    return []
  }
  return todoStore.todos
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    .slice(0, 5)
})

// 趋势数据
const todoTrendData = computed(() => statisticsStore.todoTrends)
const articleTrendData = computed(() => statisticsStore.articleTrends)

// 方法
const refreshTodoStats = () => {
  statisticsStore.fetchTodoStats()
}

const refreshArticleStats = () => {
  statisticsStore.fetchArticleStats()
}

const refreshAllStats = () => {
  statisticsStore.fetchStatistics()
}

const viewAllArticles = () => {
  router.push('/dashboard/articles')
}

const viewAllTodos = () => {
  router.push('/dashboard/todos')
}

const viewArticle = (id) => {
  router.push(`/articles/${id}`)
}

const viewTodo = (id) => {
  router.push('/dashboard/todos')
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

const getStatusType = (status) => {
  switch (status) {
    case 'completed': return 'success'
    case 'in_progress': return 'warning'
    case 'pending': return 'info'
    default: return 'info'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'completed': return '已完成'
    case 'in_progress': return '进行中'
    case 'pending': return '待处理'
    default: return '未知'
  }
}

// 图表相关方法
const getBarHeight = (value, data) => {
  if (!data || data.length === 0) return '0px'

  // 统一使用count字段，通过type字段判断数据类型
  const maxValue = Math.max(...data.map(item => Number(item.count || 0)))

  // 确保value是数字
  const numValue = Number(value) || 0

  // 如果所有值都是0，则0值显示4px，其他情况显示0px
  if (maxValue === 0) {
    return numValue === 0 ? '4px' : '0px'
  }

  // 0值显示4px的灰色小柱子
  if (numValue === 0) {
    return '4px'
  }

  // 计算非0值的高度
  const percentage = (numValue / maxValue) * 100
  const height = Math.max(percentage * 2, 20) // 最小高度20px，最大200px

  return `${height}px`
}

const formatChartDate = (dateString, index, totalLength) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const currentDays = parseInt(todoTimeRange.value) || parseInt(articleTimeRange.value) || 7

  // 根据时间范围调整日期格式
  if (currentDays <= 7) {
    // 7天显示：月/日
    return `${date.getMonth() + 1}/${date.getDate()}`
  } else if (currentDays <= 30) {
    // 30天显示：每隔几天显示一次，确保不会太密集
    if (index % 3 === 0 || index === totalLength - 1) {
      return `${date.getMonth() + 1}/${date.getDate()}`
    }
    return ''
  } else {
    // 90天显示：每隔几天显示一次
    if (index % 7 === 0 || index === totalLength - 1) {
      return `${date.getMonth() + 1}/${date.getDate()}`
    }
    return ''
  }
}

// 处理柱状图点击事件
const handleBarClick = (type, item, index) => {
  const date = new Date(item.date)
  const formattedDate = date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })

  if (type === 'todo') {
    ElMessage.info(`${formattedDate}: 完成了 ${item.count || 0} 个任务`)
    // 可以在这里添加跳转到任务列表的逻辑
    // router.push({ name: 'todos', query: { date: item.date } })
  } else if (type === 'article') {
    ElMessage.info(`${formattedDate}: 文章浏览量 ${item.count || 0} 次`)
    // 可以在这里添加跳转到文章列表的逻辑
    // router.push({ name: 'articles', query: { date: item.date } })
  }
}

// 监听时间范围变化
watch(todoTimeRange, (newValue) => {
  statisticsStore.fetchTodoTrends({ days: newValue })
  updateChartGap(newValue)
})

watch(articleTimeRange, (newValue) => {
  statisticsStore.fetchArticleTrends({ days: newValue })
  updateChartGap(newValue)
})

// 根据时间范围调整图表空隙
const updateChartGap = (timeRange) => {
  const root = document.documentElement
  if (timeRange <= 7) {
    root.style.setProperty('--chart-gap', '4px')
  } else if (timeRange <= 30) {
    root.style.setProperty('--chart-gap', '8px')
  } else {
    root.style.setProperty('--chart-gap', '12px')
  }
}

onMounted(async () => {
  // 加载所有数据
  await Promise.all([
    refreshAllStats(),
    todoStore.fetchTodos(),
    articleStore.fetchArticles()
  ])

  // 获取趋势数据
  statisticsStore.fetchTodoTrends({ days: todoTimeRange.value })
  statisticsStore.fetchArticleTrends({ days: articleTimeRange.value })

  // 初始化图表空隙设置
  updateChartGap(todoTimeRange.value)
})
</script>

<style scoped>
.analytics-page {
  padding: 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 40px;
  border-radius: 16px;
}

.page-header h1 {
  margin: 0 0 10px 0;
  font-size: 32px;
}

.page-header p {
  margin: 0;
  font-size: 16px;
  opacity: 0.9;
}

.analytics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
}

.stat-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.chart-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.list-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.stat-content {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  transition: transform 0.3s ease;
}

.stat-item:hover {
  transform: translateY(-2px);
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 8px;
}

.stat-number.success {
  color: #67c23a;
}

.stat-number.warning {
  color: #e6a23c;
}

.stat-number.danger {
  color: #f56c6c;
}

.stat-number.info {
  color: #409eff;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

.chart-container {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  text-align: center;
  color: #909399;
}

.chart-placeholder p {
  margin: 10px 0 5px 0;
  font-size: 16px;
}

.chart-desc {
  font-size: 14px;
  opacity: 0.7;
}

.trend-chart {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.chart-scroll-container {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: thin;
  scrollbar-color: #c1c1c1 #f1f1f1;
}

.chart-scroll-container::-webkit-scrollbar {
  height: 6px;
}

.chart-scroll-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.chart-scroll-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.chart-scroll-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.chart-bars {
  flex: 1;
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: var(--chart-gap, 4px);
  padding: 40px 10px 10px 10px;
  min-height: 200px;
  min-width: max-content;
  position: relative;
}

.chart-bar {
  background: #409eff;
  border-radius: 4px 4px 0 0;
  min-width: 40px;
  width: 40px;
  position: relative;
  transition: all 0.3s ease;
  cursor: pointer;
  border: 1px solid #e1e8ed;
}

.chart-bar:hover {
  background: #337ecc;
  transform: scaleY(1.02);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.zero-bar {
  background: #e4e7ed !important;
  border: 1px solid #dcdfe6;
}

.zero-bar:hover {
  background: #d3d7e0 !important;
  transform: scaleY(1.02);
  box-shadow: 0 2px 8px rgba(228, 231, 237, 0.3);
}

.bar-value {
  position: absolute;
  top: -25px;
  left: 50%;
  transform: translateX(-50%);
  color: #2c3e50;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
  opacity: 1;
  transition: all 0.3s ease;
  z-index: 10;
}

.chart-bar:hover .bar-value {
  color: #1a252f;
  transform: translateX(-50%) scale(1.05);
}

.chart-labels {
  display: flex;
  justify-content: space-between;
  gap: var(--chart-gap, 4px);
  padding: 10px;
  border-top: 1px solid #f0f0f0;
  min-width: max-content;
}

.chart-label {
  min-width: 40px;
  text-align: center;
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
}



.article-list,
.todo-list {
  max-height: 400px;
  overflow-y: auto;
}

.article-item,
.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.article-item:last-child,
.todo-item:last-child {
  border-bottom: none;
}

.article-info,
.todo-info {
  flex: 1;
  margin-right: 15px;
}

.article-info h4,
.todo-info h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #2c3e50;
}

.article-info p,
.todo-info p {
  margin: 0 0 8px 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.4;
}

.article-meta,
.todo-meta {
  display: flex;
  gap: 15px;
  font-size: 12px;
  color: #909399;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .analytics-grid {
    grid-template-columns: 1fr;
  }

  .stat-content {
    grid-template-columns: 1fr;
  }

  .article-item,
  .todo-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
