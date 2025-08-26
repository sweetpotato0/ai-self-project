<template>
  <div class="dashboard-home">
    <div class="welcome-section">
      <h1>欢迎回来，{{ authStore.user?.username }}！</h1>
      <p>今天是 {{ currentDate }}</p>
    </div>

    <div class="stats-grid">
      <el-card
        v-for="stat in stats"
        :key="stat.title"
        class="stat-card"
        @click="viewStatDetails(stat)"
      >
        <div class="stat-content">
          <div class="stat-icon" :style="{ background: stat.color }">
            <el-icon><component :is="stat.icon" /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ stat.value }}</div>
            <div class="stat-title">{{ stat.title }}</div>
          </div>
          <div class="stat-arrow">
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
      </el-card>
    </div>

    <div class="quick-actions">
      <h3>快速操作</h3>
      <div class="actions-grid">
        <el-card
          v-for="action in quickActions"
          :key="action.title"
          @click="handleQuickAction(action)"
          class="action-card"
        >
          <div class="action-content">
            <el-icon class="action-icon" :style="{ color: action.color }">
              <component :is="action.icon" />
            </el-icon>
            <div class="action-text">
              <div class="action-title">{{ action.title }}</div>
              <div class="action-desc">{{ action.description }}</div>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useStatisticsStore } from '@/stores/statistics'
import { useTodoStore } from '@/stores/todo'
import { useArticleStore } from '@/stores/article'
import { List, Check, Clock, Document, Plus, Calendar, TrendCharts, Setting, ArrowRight, Edit, View } from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()
const statisticsStore = useStatisticsStore()
const todoStore = useTodoStore()
const articleStore = useArticleStore()

const currentDate = ref('')

const stats = computed(() => {
  const todoStats = statisticsStore.todoStats
  const articleStats = statisticsStore.articleStats
  return [
    { title: '总任务数', value: todoStats.total, icon: 'List', color: '#409EFF', type: 'todo' },
    { title: '已完成', value: todoStats.completed, icon: 'Check', color: '#67C23A', type: 'todo' },
    { title: '进行中', value: todoStats.inProgress, icon: 'Clock', color: '#E6A23C', type: 'todo' },
    { title: '待处理', value: todoStats.pending, icon: 'Document', color: '#F56C6C', type: 'todo' },
    { title: '文章总数', value: articleStats.total, icon: 'Edit', color: '#9C27B0', type: 'article' },
    { title: '已发布', value: articleStats.published, icon: 'View', color: '#4CAF50', type: 'article' },
    { title: '草稿', value: articleStats.draft, icon: 'Document', color: '#FF9800', type: 'article' },
    { title: '总浏览量', value: articleStats.totalViews, icon: 'View', color: '#2196F3', type: 'article' }
  ]
})

const quickActions = ref([
  { title: '新建任务', description: '快速创建新的TODO任务', icon: 'Plus', color: '#409EFF', action: 'create-todo' },
  { title: '新建文章', description: '创建新的技术文章', icon: 'Edit', color: '#9C27B0', action: 'create-article' },
  { title: '查看日程', description: '查看今日和本周安排', icon: 'Calendar', color: '#67C23A', action: 'view-calendar' },
  { title: '个人文章', description: '管理您的文章库', icon: 'Document', color: '#4CAF50', action: 'view-articles' },
  { title: '数据分析', description: '查看任务完成统计', icon: 'TrendCharts', color: '#E6A23C', action: 'view-analytics' },
  { title: '系统设置', description: '个性化您的使用体验', icon: 'Setting', color: '#909399', action: 'view-settings' }
])

const handleQuickAction = (action) => {
  switch (action.action) {
    case 'create-todo': router.push({ name: 'todos' }); break
    case 'create-article': router.push({ name: 'articles' }); break
    case 'view-calendar': router.push({ name: 'calendar' }); break
    case 'view-articles': router.push({ name: 'articles' }); break
    case 'view-analytics': router.push({ name: 'analytics' }); break
    case 'view-settings': router.push({ name: 'settings' }); break
  }
}

const viewStatDetails = (stat) => {
  // 根据统计类型跳转到对应的详细页面
  if (stat.type === 'todo') {
    router.push({ name: 'todos' })
  } else if (stat.type === 'article') {
    router.push({ name: 'articles' })
  } else {
    router.push({ name: 'todos' })
  }
}

onMounted(() => {
  const now = new Date()
  currentDate.value = now.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' })
  statisticsStore.fetchStatistics()
  todoStore.fetchTodos()
  articleStore.fetchArticles()
})
</script>

<style scoped>
.dashboard-home {
  padding: 20px;
}

.welcome-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 40px;
  border-radius: 16px;
  margin-bottom: 30px;
  text-align: center;
}

.welcome-section h1 {
  margin: 0 0 10px 0;
  font-size: 28px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
  position: relative;
}

.stat-arrow {
  position: absolute;
  right: 0;
  color: #c0c4cc;
  transition: color 0.3s ease;
}

.stat-card:hover .stat-arrow {
  color: #409eff;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: #2c3e50;
}

.stat-title {
  font-size: 14px;
  color: #606266;
  margin-top: 5px;
}

.quick-actions h3 {
  margin: 0 0 20px 0;
  font-size: 20px;
  color: #2c3e50;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.action-card {
  cursor: pointer;
  transition: transform 0.3s ease;
}

.action-card:hover {
  transform: translateY(-3px);
}

.action-content {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 20px;
}

.action-icon {
  font-size: 32px;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(64, 158, 255, 0.1);
  border-radius: 12px;
}

.action-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 5px;
}

.action-desc {
  font-size: 14px;
  color: #606266;
}
</style>
