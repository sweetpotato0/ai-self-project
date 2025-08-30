<template>
  <div class="audit-logs-container">
    <div class="page-header">
      <h1>系统审计日志</h1>
      <p>查看系统中所有用户的操作记录和审计信息</p>
    </div>

    <!-- 搜索过滤器 -->
    <el-card class="filter-card">
      <el-form :model="filters" @submit.prevent="loadAuditLogs" class="filter-form">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-form-item label="用户名">
              <el-input 
                v-model="filters.username" 
                placeholder="输入用户名" 
                clearable
                @keyup.enter="loadAuditLogs"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="操作类型">
              <el-select v-model="filters.action" placeholder="选择操作类型" clearable>
                <el-option label="查询" value="查询" />
                <el-option label="创建" value="创建" />
                <el-option label="更新" value="更新" />
                <el-option label="删除" value="删除" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="资源类型">
              <el-select v-model="filters.resource" placeholder="选择资源类型" clearable>
                <el-option label="审计日志" value="audit-logs" />
                <el-option label="文章" value="articles" />
                <el-option label="通知" value="notifications" />
                <el-option label="设置" value="settings" />
                <el-option label="统计" value="statistics" />
                <el-option label="任务" value="todos" />
                <el-option label="分类" value="categories" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="状态码">
              <el-select v-model="filters.status_code" placeholder="选择状态码" clearable>
                <el-option label="200 成功" :value="200" />
                <el-option label="400 客户端错误" :value="400" />
                <el-option label="401 未授权" :value="401" />
                <el-option label="403 禁止访问" :value="403" />
                <el-option label="404 未找到" :value="404" />
                <el-option label="500 服务器错误" :value="500" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="6">
            <el-form-item label="IP地址">
              <el-input 
                v-model="filters.ip_address" 
                placeholder="输入IP地址" 
                clearable
                @keyup.enter="loadAuditLogs"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="时间范围">
              <el-date-picker
                v-model="dateRange"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DDTHH:mm:ssZ"
                @change="handleDateRangeChange"
              />
            </el-form-item>
          </el-col>
          <el-col :span="4">
            <el-form-item label=" ">
              <el-button type="primary" @click="loadAuditLogs" :loading="loading">
                <el-icon><Search /></el-icon>
                搜索
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label=" ">
              <el-button @click="resetFilters">重置</el-button>
              <el-button @click="showStatsDialog = true">
                <el-icon><DataAnalysis /></el-icon>
                统计信息
              </el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 审计日志表格 -->
    <el-card class="table-card">
      <template #header>
        <div class="table-header">
          <span>审计日志 (共 {{ pagination.total }} 条)</span>
          <div>
            <el-button 
              size="small" 
              @click="loadAuditLogs"
              :loading="loading"
            >
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table 
        :data="auditLogs" 
        :loading="loading" 
        stripe
        @row-click="handleRowClick"
        style="cursor: pointer;"
      >
        <el-table-column prop="timestamp" label="时间" width="160" sortable>
          <template #default="{ row }">
            <div class="time-cell">
              <el-icon><Clock /></el-icon>
              {{ formatTime(row.timestamp) }}
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="username" label="用户" width="120">
          <template #default="{ row }">
            <div class="user-cell">
              <el-icon><User /></el-icon>
              {{ row.username || '未知用户' }}
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="action" label="操作" width="120">
          <template #default="{ row }">
            <el-tag :type="getActionTypeColor(row.action)" size="small">
              {{ row.action }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="resource" label="资源" width="100">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.resource }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="method" label="方法" width="80">
          <template #default="{ row }">
            <el-tag 
              :type="getMethodTypeColor(row.method)" 
              size="small"
            >
              {{ row.method }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status_code" label="状态" width="80">
          <template #default="{ row }">
            <el-tag 
              :type="getStatusTypeColor(row.status_code)" 
              size="small"
            >
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="duration" label="耗时" width="80">
          <template #default="{ row }">
            {{ formatDuration(row.duration) }}
          </template>
        </el-table-column>

        <el-table-column prop="ip_address" label="IP地址" width="140">
          <template #default="{ row }">
            <div class="ip-cell">
              <el-icon><LocationInformation /></el-icon>
              {{ row.ip_address }}
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="path" label="路径" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <code class="path-code">{{ row.path }}</code>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :page-sizes="[20, 50, 100, 200]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog 
      v-model="showDetailDialog" 
      title="审计日志详情" 
      width="80%"
      class="detail-dialog"
    >
      <div v-if="selectedLog" class="detail-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="操作时间">
            {{ formatTime(selectedLog.timestamp) }}
          </el-descriptions-item>
          <el-descriptions-item label="用户">
            {{ selectedLog.username || '未知用户' }} (ID: {{ selectedLog.user_id || 'N/A' }})
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <el-tag :type="getActionTypeColor(selectedLog.action)">
              {{ selectedLog.action }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="资源类型">
            <el-tag type="info">{{ selectedLog.resource }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="HTTP方法">
            <el-tag :type="getMethodTypeColor(selectedLog.method)">
              {{ selectedLog.method }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态码">
            <el-tag :type="getStatusTypeColor(selectedLog.status_code)">
              {{ selectedLog.status_code }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="请求耗时">
            {{ formatDuration(selectedLog.duration) }}
          </el-descriptions-item>
          <el-descriptions-item label="IP地址">
            {{ selectedLog.ip_address }}
          </el-descriptions-item>
          <el-descriptions-item label="请求路径" span="2">
            <code>{{ selectedLog.path }}</code>
          </el-descriptions-item>
          <el-descriptions-item label="用户代理" span="2">
            {{ selectedLog.user_agent || 'N/A' }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="selectedLog.request_body" class="detail-section">
          <h4>请求内容</h4>
          <el-input
            type="textarea"
            :model-value="formatJsonString(selectedLog.request_body)"
            :rows="6"
            readonly
          />
        </div>

        <div v-if="selectedLog.response_body" class="detail-section">
          <h4>响应内容</h4>
          <el-input
            type="textarea"
            :model-value="formatJsonString(selectedLog.response_body)"
            :rows="6"
            readonly
          />
        </div>
      </div>
    </el-dialog>

    <!-- 统计信息对话框 -->
    <el-dialog 
      v-model="showStatsDialog" 
      title="审计日志统计" 
      width="70%"
      @open="loadStats"
    >
      <div v-if="stats" class="stats-content">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-statistic title="总计日志" :value="stats.total_count" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="今日日志" :value="stats.today_count" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="成功请求" :value="stats.success_count" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="错误请求" :value="stats.error_count" />
          </el-col>
        </el-row>

        <el-row :gutter="20" style="margin-top: 20px;">
          <el-col :span="12">
            <h4>活跃用户 Top 10</h4>
            <el-table :data="stats.top_users" size="small">
              <el-table-column prop="username" label="用户名" />
              <el-table-column prop="count" label="操作次数" />
            </el-table>
          </el-col>
          <el-col :span="12">
            <h4>热门操作 Top 10</h4>
            <el-table :data="stats.top_actions" size="small">
              <el-table-column prop="action" label="操作类型" />
              <el-table-column prop="count" label="次数" />
            </el-table>
          </el-col>
        </el-row>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Search, 
  Refresh, 
  Clock, 
  User, 
  LocationInformation,
  DataAnalysis
} from '@element-plus/icons-vue'
import { auditApi } from '../api'

// 响应式数据
const loading = ref(false)
const auditLogs = ref([])
const selectedLog = ref(null)
const showDetailDialog = ref(false)
const showStatsDialog = ref(false)
const stats = ref(null)
const dateRange = ref([])

// 过滤器
const filters = reactive({
  username: '',
  action: '',
  resource: '',
  method: '',
  ip_address: '',
  status_code: null,
  start_time: '',
  end_time: '',
  order_by: 'timestamp',
  order: 'desc'
})

// 分页
const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

// 加载审计日志
const loadAuditLogs = async () => {
  loading.value = true
  try {
    const params = {
      ...filters,
      page: pagination.page,
      limit: pagination.limit
    }
    
    // 清理空值参数
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null) {
        delete params[key]
      }
    })

    const response = await auditApi.getAuditLogs(params)
    auditLogs.value = response.data.items || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('Failed to load audit logs:', error)
    ElMessage.error('加载审计日志失败')
  } finally {
    loading.value = false
  }
}

// 加载统计信息
const loadStats = async () => {
  try {
    const params = {}
    if (filters.start_time) params.start_time = filters.start_time
    if (filters.end_time) params.end_time = filters.end_time
    
    const response = await auditApi.getAuditLogStats(params)
    stats.value = response.data
  } catch (error) {
    console.error('Failed to load audit stats:', error)
    ElMessage.error('加载统计信息失败')
  }
}

// 处理日期范围变化
const handleDateRangeChange = (range) => {
  if (range && range.length === 2) {
    filters.start_time = range[0]
    filters.end_time = range[1]
  } else {
    filters.start_time = ''
    filters.end_time = ''
  }
}

// 重置过滤器
const resetFilters = () => {
  Object.keys(filters).forEach(key => {
    if (key === 'order_by') {
      filters[key] = 'timestamp'
    } else if (key === 'order') {
      filters[key] = 'desc'
    } else {
      filters[key] = ''
    }
  })
  dateRange.value = []
  pagination.page = 1
  loadAuditLogs()
}

// 处理行点击
const handleRowClick = (row) => {
  selectedLog.value = row
  showDetailDialog.value = true
}

// 处理页码变化
const handlePageChange = (page) => {
  pagination.page = page
  loadAuditLogs()
}

// 处理页大小变化
const handleSizeChange = (size) => {
  pagination.limit = size
  pagination.page = 1
  loadAuditLogs()
}

// 格式化时间
const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

// 格式化耗时
const formatDuration = (duration) => {
  if (duration > 1000000000) {
    return `${(duration / 1000000000).toFixed(2)}s`
  } else if (duration > 1000000) {
    return `${(duration / 1000000).toFixed(0)}ms`
  } else if (duration > 1000) {
    return `${(duration / 1000).toFixed(0)}μs`
  } else {
    return `${duration}ns`
  }
}

// 格式化JSON字符串
const formatJsonString = (str) => {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

// 获取操作类型颜色
const getActionTypeColor = (action) => {
  const colorMap = {
    create: 'success',
    view: 'info',
    list: 'info',
    update: 'warning',
    delete: 'danger',
    login: 'primary',
    logout: 'primary'
  }
  return colorMap[action] || 'info'
}

// 获取HTTP方法颜色
const getMethodTypeColor = (method) => {
  const colorMap = {
    GET: 'primary',
    POST: 'success',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'warning'
  }
  return colorMap[method] || 'info'
}

// 获取状态码颜色
const getStatusTypeColor = (statusCode) => {
  if (statusCode >= 200 && statusCode < 300) return 'success'
  if (statusCode >= 300 && statusCode < 400) return 'warning'
  if (statusCode >= 400 && statusCode < 500) return 'danger'
  if (statusCode >= 500) return 'danger'
  return 'info'
}

// 组件挂载
onMounted(() => {
  loadAuditLogs()
})
</script>

<style scoped>
.audit-logs-container {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
  text-align: center;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.page-header p {
  font-size: 16px;
  color: var(--el-text-color-regular);
  margin: 0;
}

.filter-card, .table-card {
  margin-bottom: 20px;
}

.filter-form {
  margin: 0;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.time-cell, .user-cell, .ip-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.path-code {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  background: var(--el-bg-color-page);
  padding: 2px 4px;
  border-radius: 2px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.detail-dialog :deep(.el-dialog__body) {
  max-height: 70vh;
  overflow-y: auto;
}

.detail-content .detail-section {
  margin-top: 20px;
}

.detail-section h4 {
  margin: 0 0 12px 0;
  color: var(--el-text-color-primary);
}

.stats-content h4 {
  margin: 16px 0 12px 0;
  color: var(--el-text-color-primary);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .audit-logs-container {
    padding: 12px;
  }
  
  .page-header h1 {
    font-size: 24px;
  }
  
  .table-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
}</style>