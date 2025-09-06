<template>
  <div class="admin-dashboard">
    <div class="page-header">
      <h1>英语视频管理后台</h1>
      <p class="page-description">管理视频系列、剧集和相关内容</p>
    </div>

    <!-- 统计概览 -->
    <div class="stats-grid">
      <el-card class="stats-card">
        <div class="stats-content">
          <div class="stats-icon">
            <el-icon size="24" color="#409EFF"><VideoPlay /></el-icon>
          </div>
          <div class="stats-info">
            <div class="stats-value">{{ stats.totalSeries || 0 }}</div>
            <div class="stats-label">视频系列</div>
          </div>
        </div>
      </el-card>

      <el-card class="stats-card">
        <div class="stats-content">
          <div class="stats-icon">
            <el-icon size="24" color="#67C23A"><List /></el-icon>
          </div>
          <div class="stats-info">
            <div class="stats-value">{{ stats.totalEpisodes || 0 }}</div>
            <div class="stats-label">总剧集数</div>
          </div>
        </div>
      </el-card>

      <el-card class="stats-card">
        <div class="stats-content">
          <div class="stats-icon">
            <el-icon size="24" color="#E6A23C"><View /></el-icon>
          </div>
          <div class="stats-info">
            <div class="stats-value">{{ formatNumber(stats.totalViews) || 0 }}</div>
            <div class="stats-label">总观看次数</div>
          </div>
        </div>
      </el-card>

      <el-card class="stats-card">
        <div class="stats-content">
          <div class="stats-icon">
            <el-icon size="24" color="#F56C6C"><User /></el-icon>
          </div>
          <div class="stats-info">
            <div class="stats-value">{{ stats.activeUsers || 0 }}</div>
            <div class="stats-label">活跃用户</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 快速操作 -->
    <el-card class="quick-actions">
      <template #header>
        <h3>快速操作</h3>
      </template>
      <div class="actions-grid">
        <el-button 
          type="primary" 
          size="large" 
          @click="$router.push('/dashboard/admin/english-videos/series')"
        >
          <el-icon><VideoPlay /></el-icon>
          管理视频系列
        </el-button>
        
        <el-button 
          type="success" 
          size="large"
          @click="$router.push('/dashboard/admin/english-videos/episodes')"
        >
          <el-icon><List /></el-icon>
          管理剧集
        </el-button>
        
        <el-button 
          type="warning" 
          size="large"
          @click="showBatchImport = true"
        >
          <el-icon><Upload /></el-icon>
          批量导入
        </el-button>
        
        <el-button 
          type="info" 
          size="large"
          @click="refreshStats"
        >
          <el-icon><Refresh /></el-icon>
          刷新统计
        </el-button>
      </div>
    </el-card>

    <!-- 最近活动 -->
    <el-card class="recent-activity">
      <template #header>
        <h3>最近更新的系列</h3>
      </template>
      <el-table :data="recentSeries" style="width: 100%">
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image
              :src="row.cover_image"
              fit="cover"
              style="width: 60px; height: 40px; border-radius: 4px;"
            />
          </template>
        </el-table-column>
        <el-table-column prop="title_cn" label="系列名称" />
        <el-table-column prop="episode_count" label="剧集数" width="80" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_published ? 'success' : 'info'">
              {{ row.is_published ? '已发布' : '未发布' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="更新时间" width="150">
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" @click="editSeries(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 批量导入对话框 -->
    <el-dialog v-model="showBatchImport" title="批量导入向导" width="600px">
      <div class="batch-import-guide">
        <el-steps :active="1" align-center>
          <el-step title="选择导入类型" />
          <el-step title="准备数据" />
          <el-step title="开始导入" />
        </el-steps>
        
        <div class="import-options">
          <el-button 
            type="primary" 
            size="large"
            @click="goToBatchImport('series')"
          >
            导入视频系列
          </el-button>
          <el-button 
            type="success" 
            size="large"
            @click="goToBatchImport('episodes')"
          >
            导入剧集
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  VideoPlay, 
  List, 
  View, 
  User, 
  Upload, 
  Refresh 
} from '@element-plus/icons-vue'
import { useEnglishVideosStore } from '../../stores/englishVideosStore'

const router = useRouter()
const englishVideosStore = useEnglishVideosStore()

// 响应式数据
const stats = ref({
  totalSeries: 0,
  totalEpisodes: 0,
  totalViews: 0,
  activeUsers: 0
})

const recentSeries = ref([])
const showBatchImport = ref(false)
const loading = ref(false)

// 方法
const loadStats = async () => {
  try {
    loading.value = true
    const [seriesResponse, statsResponse] = await Promise.all([
      englishVideosStore.fetchVideoSeries(),
      englishVideosStore.getUserVideoStats()
    ])
    
    const series = seriesResponse || []
    stats.value = {
      totalSeries: series.length,
      totalEpisodes: series.reduce((sum, s) => sum + (s.episode_count || 0), 0),
      totalViews: statsResponse.total_watch_time || 0,
      activeUsers: statsResponse.learning_streak || 0
    }
    
    // 获取最近更新的系列（取前5个）
    recentSeries.value = series
      .sort((a, b) => new Date(b.updated_at || b.created_at) - new Date(a.updated_at || a.created_at))
      .slice(0, 5)
      
  } catch (error) {
    ElMessage.error('加载统计数据失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const refreshStats = () => {
  ElMessage.info('正在刷新统计数据...')
  loadStats()
}

const editSeries = (series) => {
  router.push({
    name: 'admin-video-series',
    query: { edit: series.id }
  })
}

const goToBatchImport = (type) => {
  showBatchImport.value = false
  if (type === 'series') {
    router.push('/dashboard/admin/english-videos/series')
  } else {
    router.push('/dashboard/admin/english-videos/episodes')
  }
}

const formatNumber = (num) => {
  if (!num) return '0'
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

const formatDate = (dateString) => {
  if (!dateString) return '--'
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 生命周期
onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.admin-dashboard {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 30px;
  text-align: center;
}

.page-header h1 {
  margin: 0 0 10px 0;
  color: #303133;
  font-size: 28px;
  font-weight: 600;
}

.page-description {
  margin: 0;
  color: #909399;
  font-size: 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stats-card {
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.stats-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stats-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  border-radius: 12px;
  background: #f5f7fa;
}

.stats-info {
  flex: 1;
}

.stats-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.stats-label {
  font-size: 14px;
  color: #909399;
}

.quick-actions {
  margin-bottom: 30px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 15px;
}

.actions-grid .el-button {
  height: 60px;
  font-size: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.recent-activity {
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.recent-activity .el-table {
  margin-top: 15px;
}

.batch-import-guide {
  padding: 20px;
  text-align: center;
}

.el-steps {
  margin-bottom: 40px;
}

.import-options {
  display: flex;
  justify-content: center;
  gap: 20px;
}

.import-options .el-button {
  height: 50px;
  min-width: 150px;
  font-size: 16px;
}
</style>