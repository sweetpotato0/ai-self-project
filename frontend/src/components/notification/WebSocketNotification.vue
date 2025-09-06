<template>
  <div class="websocket-notification">
    <!-- 连接状态指示器 -->
    <div v-if="showConnectionStatus" class="connection-status" :class="connectionStatus">
      <el-icon v-if="connectionStatus === 'connected'"><CircleCheck /></el-icon>
      <el-icon v-else-if="connectionStatus === 'connecting'"><Loading /></el-icon>
      <el-icon v-else><CircleClose /></el-icon>
      <span>{{ connectionStatusText }}</span>
    </div>

    <!-- 实时通知弹窗 -->
    <el-dialog
      v-model="showNotificationDialog"
      title="🔔 新消息提醒"
      width="400px"
      :show-close="true"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
    >
      <div class="notification-content">
        <div class="notification-header">
          <el-icon :class="getNotificationIcon(currentNotification?.type)">
            <Bell v-if="currentNotification?.type === 'overdue'" />
            <Clock v-else-if="currentNotification?.type === 'due_soon'" />
            <CircleCheck v-else-if="currentNotification?.type === 'completed'" />
            <Bell v-else />
          </el-icon>
          <span class="notification-title">{{ currentNotification?.title }}</span>
        </div>

        <div class="notification-message">
          {{ currentNotification?.message }}
        </div>

        <div v-if="currentNotification?.data" class="notification-data">
          <div v-if="currentNotification.data.task_title" class="data-item">
            <strong>任务:</strong> {{ currentNotification.data.task_title }}
          </div>
          <div v-if="currentNotification.data.due_date" class="data-item">
                            <strong>截止时间:</strong> {{ formatDateTimeLocal(currentNotification.data.due_date) }}
          </div>
          <div v-if="currentNotification.data.hours_left" class="data-item">
            <strong>剩余时间:</strong> {{ currentNotification.data.hours_left }} 小时
          </div>
          <div v-if="currentNotification.data.days_overdue" class="data-item">
            <strong>逾期天数:</strong> {{ currentNotification.data.days_overdue }} 天
          </div>
        </div>
      </div>

      <template #footer>
        <div class="notification-actions">
          <el-button @click="markAsRead" type="primary" size="small">
            标记已读
          </el-button>
          <el-button @click="viewTask" type="success" size="small" v-if="currentNotification?.data?.task_id">
            查看任务
          </el-button>
          <el-button @click="showNotificationDialog = false" size="small">
            关闭
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { formatDateTime } from '@/utils/dateTime'
import { Bell, CircleCheck, CircleClose, Clock, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

// 响应式数据
const ws = ref(null)
const connectionStatus = ref('disconnected') // disconnected, connecting, connected
const showConnectionStatus = ref(false)
const showNotificationDialog = ref(false)
const currentNotification = ref(null)
const reconnectAttempts = ref(0)
const maxReconnectAttempts = 5
const reconnectInterval = ref(null)

// 计算属性
const connectionStatusText = computed(() => {
  switch (connectionStatus.value) {
    case 'connected':
      return '实时通知已连接'
    case 'connecting':
      return '正在连接...'
    case 'disconnected':
      return '实时通知已断开'
    default:
      return '连接状态未知'
  }
})

// 连接WebSocket
const connectWebSocket = () => {
  if (!authStore.token) {
    console.log('No auth token, skipping WebSocket connection')
    return
  }

  connectionStatus.value = 'connecting'
  showConnectionStatus.value = true

  const wsUrl = `ws://localhost:8080/api/v1/ws?token=${authStore.token}`

  try {
    ws.value = new WebSocket(wsUrl)

    ws.value.onopen = () => {
      // console.log('WebSocket connected')
      connectionStatus.value = 'connected'
      reconnectAttempts.value = 0

      // 3秒后隐藏连接状态
      setTimeout(() => {
        showConnectionStatus.value = false
      }, 3000)
    }

    ws.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        handleWebSocketMessage(data)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    ws.value.onclose = () => {
      // console.log('WebSocket disconnected')
      connectionStatus.value = 'disconnected'

      // 自动重连
      if (reconnectAttempts.value < maxReconnectAttempts) {
        reconnectAttempts.value++
        const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.value), 30000)

        // console.log(`Reconnecting in ${delay}ms (attempt ${reconnectAttempts.value}/${maxReconnectAttempts})`)

        reconnectInterval.value = setTimeout(() => {
          connectWebSocket()
        }, delay)
      } else {
        console.log('Max reconnection attempts reached')
        showConnectionStatus.value = true
      }
    }

    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
      connectionStatus.value = 'disconnected'
    }

  } catch (error) {
    console.error('Failed to create WebSocket connection:', error)
    connectionStatus.value = 'disconnected'
  }
}

// 处理WebSocket消息
const handleWebSocketMessage = (data) => {
  switch (data.type) {
    case 'connection':
      // console.log('WebSocket connection established:', data.data)
      break

    case 'notification':
      console.log('Received notification:', data.data)
      showNotificationPopup(data.data)
      // 更新通知store
      notificationStore.addNotification(data.data)
      break

    case 'notifications':
      console.log('Received notifications list:', data.data)
      notificationStore.setNotifications(data.data)
      break

    case 'pong':
      console.log('Received pong:', data.data)
      break

    default:
      console.log('Unknown message type:', data.type)
  }
}

// 显示通知弹窗
const showNotificationPopup = (notification) => {
  currentNotification.value = notification
  showNotificationDialog.value = true

  // 播放提示音（如果浏览器支持）
  try {
    const audio = new Audio('/notification.mp3')
    audio.play().catch(() => {
      // 如果无法播放音频，忽略错误
    })
  } catch (error) {
    // 忽略音频播放错误
  }

  // 显示浏览器通知（如果用户允许）
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification(notification.title, {
      body: notification.message,
      icon: '/favicon.ico',
      tag: 'task-notification'
    })
  }
}

// 获取通知图标
const getNotificationIcon = (type) => {
  switch (type) {
    case 'overdue':
      return 'notification-overdue'
    case 'due_soon':
      return 'notification-due-soon'
    case 'completed':
      return 'notification-completed'
    default:
      return 'notification-default'
  }
}

// 格式化日期时间
const formatDateTimeLocal = formatDateTime

// 标记通知为已读
const markAsRead = async () => {
  if (currentNotification.value?.id) {
    try {
      await notificationStore.markAsRead(currentNotification.value.id)
      showNotificationDialog.value = false
      ElMessage.success('已标记为已读')
    } catch (error) {
      ElMessage.error('标记已读失败')
    }
  }
}

// 查看任务
const viewTask = () => {
  if (currentNotification.value?.data?.task_id) {
    showNotificationDialog.value = false
    router.push(`/dashboard/todos?task=${currentNotification.value.data.task_id}`)
  }
}

// 发送ping消息保持连接
const sendPing = () => {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({ type: 'ping' }))
  }
}

// 请求通知权限
const requestNotificationPermission = () => {
  if ('Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission()
  }
}

// 生命周期
onMounted(() => {
  // 请求通知权限
  requestNotificationPermission()

  // 连接WebSocket
  connectWebSocket()

  // 定期发送ping保持连接
  const pingInterval = setInterval(sendPing, 30000) // 每30秒ping一次

  onUnmounted(() => {
    clearInterval(pingInterval)
    if (reconnectInterval.value) {
      clearTimeout(reconnectInterval.value)
    }
    if (ws.value) {
      ws.value.close()
    }
  })
})
</script>

<style scoped>
.websocket-notification {
  position: relative;
}

.connection-status {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 9999;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.3s ease;
}

.connection-status.connected {
  background: rgba(103, 194, 58, 0.9);
  color: white;
}

.connection-status.connecting {
  background: rgba(230, 162, 60, 0.9);
  color: white;
}

.connection-status.disconnected {
  background: rgba(245, 108, 108, 0.9);
  color: white;
}

.notification-content {
  padding: 10px 0;
}

.notification-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.notification-title {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

.notification-message {
  margin-bottom: 15px;
  line-height: 1.5;
  color: #666;
}

.notification-data {
  background: #f8f9fa;
  padding: 12px;
  border-radius: 6px;
  margin-bottom: 15px;
}

.data-item {
  margin-bottom: 8px;
  font-size: 14px;
}

.data-item:last-child {
  margin-bottom: 0;
}

.data-item strong {
  color: #333;
  margin-right: 8px;
}

.notification-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.notification-overdue {
  color: #f56c6c;
}

.notification-due-soon {
  color: #e6a23c;
}

.notification-completed {
  color: #67c23a;
}

.notification-default {
  color: #409eff;
}
</style>
