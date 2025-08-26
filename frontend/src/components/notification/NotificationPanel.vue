<template>
  <el-popover
    placement="bottom-end"
    :width="350"
    trigger="click"
    popper-class="notification-popover"
  >
    <template #reference>
      <el-button type="text" class="notification-btn" @click="handleClick">
        <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
          <el-icon><Bell /></el-icon>
        </el-badge>
      </el-button>
    </template>

    <div class="notification-panel">
      <div class="notification-header">
        <h3>消息通知 ({{ unreadCount }} 未读)</h3>
        <el-button type="text" size="small" @click="markAllAsRead">
          全部已读
        </el-button>
      </div>

      <div class="notification-list">
        <div v-if="notifications.length === 0" class="no-notifications">
          <el-icon><Bell /></el-icon>
          <p>暂无消息通知</p>
        </div>

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
            <el-icon v-else><InfoFilled /></el-icon>
          </div>

          <div class="notification-content">
            <div class="notification-title">{{ notification.message }}</div>
            <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
          </div>

          <div class="notification-actions">
            <el-button
              type="text"
              size="small"
              @click.stop="markAsRead(notification)"
              v-if="!notification.is_read"
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

      <div class="notification-footer">
        <el-button type="text" size="small" @click="viewAllNotifications">
          查看全部
        </el-button>
      </div>
    </div>
  </el-popover>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Bell, Clock, Warning, Check, InfoFilled, Delete } from '@element-plus/icons-vue'
import { notificationsApi } from '@/features/notifications/api'
import { formatRelativeTime } from '@/utils/dateTime'

const router = useRouter()

// 通知数据
const notifications = ref([])
const loading = ref(false)

// 计算属性
const unreadCount = computed(() => {
  return notifications.value.filter(n => !n.is_read).length
})

// 方法
const formatTime = formatRelativeTime

const markAsRead = async (notification) => {
  try {
    await notificationsApi.markAsRead(notification.id)
    notification.is_read = true
    ElMessage.success('已标记为已读')
  } catch (error) {
    ElMessage.error('标记已读失败')
  }
}

const markAllAsRead = async () => {
  try {
    await notificationsApi.markAllAsRead()
    notifications.value.forEach(n => n.is_read = true)
    ElMessage.success('已全部标记为已读')
  } catch (error) {
    ElMessage.error('标记全部已读失败')
  }
}

const deleteNotification = async (notification) => {
  try {
    await notificationsApi.deleteNotification(notification.id)
    const index = notifications.value.findIndex(n => n.id === notification.id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
    ElMessage.success('通知已删除')
  } catch (error) {
    ElMessage.error('删除通知失败')
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
    default:
      break
  }
}

const viewAllNotifications = () => {
  // 跳转到通知管理页面
  router.push({ name: 'notifications' })
}

const handleClick = () => {
  console.log('Notification button clicked!')
  ElMessage.info('通知按钮被点击了！')
}

// 获取通知列表
const fetchNotifications = async () => {
  loading.value = true
  try {
    const response = await notificationsApi.getNotifications({ limit: 20 })
    notifications.value = response.data.notifications || []
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
    ElMessage.error('获取通知失败')
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  console.log('NotificationPanel component mounted')
  fetchNotifications()
})
</script>

<style scoped>
.notification-btn {
  font-size: 20px;
  color: #606266;
  background: none;
  border: none;
  padding: 8px;
  border-radius: 4px;
  transition: all 0.2s;
  position: relative;
  z-index: 1000;
}

.notification-btn:hover {
  background: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.notification-btn:active {
  background: rgba(64, 158, 255, 0.2);
}

.notification-panel {
  max-height: 500px;
  display: flex;
  flex-direction: column;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.notification-header h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.notification-list {
  flex: 1;
  overflow-y: auto;
  max-height: 350px;
}

.no-notifications {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

.no-notifications .el-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.no-notifications p {
  margin: 0;
  font-size: 14px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.2s;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-item.unread {
  background-color: #f0f9ff;
}

.notification-item.unread:hover {
  background-color: #e6f7ff;
}

.notification-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #e6f7ff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
  line-height: 1.4;
}

.notification-time {
  font-size: 12px;
  color: #909399;
}

.notification-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.notification-item:hover .notification-actions {
  opacity: 1;
}

.notification-footer {
  padding: 12px 16px;
  border-top: 1px solid #e4e7ed;
  text-align: center;
}

/* 自定义弹出框样式 */
:deep(.notification-popover) {
  padding: 0;
}
</style>
