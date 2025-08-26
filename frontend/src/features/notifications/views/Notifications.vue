<template>
  <div class="notifications-page">
    <div class="page-header">
      <h1>消息通知</h1>
      <div class="header-actions">
        <el-button @click="markAllAsRead" type="primary" :loading="loading">
          全部标记为已读
        </el-button>
        <el-button @click="refreshNotifications" :loading="loading">
          刷新
        </el-button>
      </div>
    </div>

    <div class="notifications-container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>

      <div v-else-if="notifications.length === 0" class="empty-state">
        <el-icon class="empty-icon"><Bell /></el-icon>
        <h3>暂无消息通知</h3>
        <p>当有新的任务提醒或系统通知时，会在这里显示</p>
      </div>

      <div v-else class="notifications-list">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          :class="[
            'notification-item',
            { 'unread': !notification.is_read }
          ]"
          @click="handleNotificationClick(notification)"
        >
          <div class="notification-icon">
            <el-icon v-if="notification.type === 'due_soon'"><Clock /></el-icon>
            <el-icon v-else-if="notification.type === 'overdue'"><Warning /></el-icon>
            <el-icon v-else-if="notification.type === 'completed'"><Check /></el-icon>
            <el-icon v-else-if="notification.type === 'article_published'"><Document /></el-icon>
            <el-icon v-else-if="notification.type === 'system'"><InfoFilled /></el-icon>
            <el-icon v-else-if="notification.type === 'welcome'"><User /></el-icon>
            <el-icon v-else-if="notification.type === 'daily_summary'"><DataBoard /></el-icon>
            <el-icon v-else-if="notification.type === 'weekly_report'"><TrendCharts /></el-icon>
            <el-icon v-else><Bell /></el-icon>
          </div>

          <div class="notification-content">
            <div class="notification-header">
              <h4 class="notification-title">{{ notification.title }}</h4>
              <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
            </div>
            <p class="notification-message">{{ notification.message }}</p>

            <div v-if="notification.data" class="notification-data">
              <div v-if="getNotificationData(notification, 'task_title')" class="data-item">
                <strong>任务:</strong> {{ getNotificationData(notification, 'task_title') }}
              </div>
              <div v-if="getNotificationData(notification, 'due_date')" class="data-item">
                <strong>截止时间:</strong> {{ formatDateTime(getNotificationData(notification, 'due_date')) }}
              </div>
              <div v-if="getNotificationData(notification, 'hours_left')" class="data-item">
                <strong>剩余时间:</strong> {{ getNotificationData(notification, 'hours_left') }} 小时
              </div>
              <div v-if="getNotificationData(notification, 'days_overdue')" class="data-item">
                <strong>逾期天数:</strong> {{ getNotificationData(notification, 'days_overdue') }} 天
              </div>
              <div v-if="getNotificationData(notification, 'article_title')" class="data-item">
                <strong>文章:</strong> {{ getNotificationData(notification, 'article_title') }}
              </div>
            </div>
          </div>

          <div class="notification-actions">
            <el-button
              v-if="!notification.is_read"
              type="text"
              size="small"
              @click.stop="markAsRead(notification)"
            >
              标记已读
            </el-button>
            <el-button
              type="text"
              size="small"
              @click.stop="deleteNotification(notification)"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="notifications.length > 0" class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Bell, Clock, Warning, Check, Document, InfoFilled, User, DataBoard, TrendCharts, Delete } from '@element-plus/icons-vue'
import { notificationApi } from '@/api/notification'
import { formatRelativeTime, formatDateTime } from '@/utils/dateTime'

const router = useRouter()

// 响应式数据
const notifications = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 方法
const formatTime = formatRelativeTime

const getNotificationData = (notification, key) => {
  try {
    const data = JSON.parse(notification.data)
    return data[key]
  } catch (error) {
    return null
  }
}

const fetchNotifications = async () => {
  loading.value = true
  try {
    const response = await notificationApi.getNotifications({
      page: currentPage.value,
      limit: pageSize.value
    })
    notifications.value = response.data.notifications || []
    total.value = response.data.total || notifications.value.length
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
    ElMessage.error('获取通知失败')
  } finally {
    loading.value = false
  }
}

const markAsRead = async (notification) => {
  try {
    await notificationApi.markAsRead(notification.id)
    notification.is_read = true
    ElMessage.success('已标记为已读')
  } catch (error) {
    ElMessage.error('标记已读失败')
  }
}

const markAllAsRead = async () => {
  try {
    await notificationApi.markAllAsRead()
    notifications.value.forEach(n => n.is_read = true)
    ElMessage.success('已全部标记为已读')
  } catch (error) {
    ElMessage.error('标记全部已读失败')
  }
}

const deleteNotification = async (notification) => {
  try {
    await ElMessageBox.confirm('确定要删除这条通知吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await notificationApi.deleteNotification(notification.id)
    const index = notifications.value.findIndex(n => n.id === notification.id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
    ElMessage.success('通知已删除')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除通知失败')
    }
  }
}

const handleNotificationClick = async (notification) => {
  // 标记为已读
  if (!notification.is_read) {
    await markAsRead(notification)
  }

  // 根据通知类型执行不同操作
  switch (notification.type) {
    case 'due_soon':
    case 'overdue':
    case 'completed':
      // 跳转到TODO详情页面
      router.push({ name: 'todos' })
      break
    case 'article_published':
      // 跳转到文章详情页面
      const articleId = getNotificationData(notification, 'article_id')
      if (articleId) {
        router.push({ name: 'article-detail', params: { id: articleId } })
      }
      break
    default:
      break
  }
}

const refreshNotifications = () => {
  fetchNotifications()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  fetchNotifications()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchNotifications()
}

// 生命周期
onMounted(() => {
  fetchNotifications()
})
</script>

<style scoped>
.notifications-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.notifications-container {
  min-height: 400px;
}

.loading-container {
  padding: 40px;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #909399;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  color: #c0c4cc;
}

.empty-state h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #606266;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
}

.notification-item:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.notification-item.unread {
  border-left: 4px solid #409eff;
  background: #f0f9ff;
}

.notification-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #e6f7ff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
  flex-shrink: 0;
  font-size: 18px;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.notification-title {
  margin: 0;
  font-size: 16px;
  color: #303133;
  font-weight: 600;
}

.notification-time {
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
  margin-left: 12px;
}

.notification-message {
  margin: 0 0 12px 0;
  color: #606266;
  line-height: 1.5;
}

.notification-data {
  background: #f8f9fa;
  padding: 12px;
  border-radius: 6px;
  font-size: 14px;
}

.data-item {
  margin-bottom: 6px;
}

.data-item:last-child {
  margin-bottom: 0;
}

.data-item strong {
  color: #303133;
  margin-right: 8px;
}

.notification-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}

.notification-item:hover .notification-actions {
  opacity: 1;
}

.pagination-container {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}
</style>
