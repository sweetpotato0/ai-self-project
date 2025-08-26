/**
 * 通知系统模块API导出文件
 * 统一管理通知相关的所有API接口
 */

import notificationsApi from './notificationsApi'

// 导出通知API
export { notificationsApi }

// 默认导出所有API
export default {
  notificationsApi
}

/**
 * API使用示例：
 * 
 * 在Store中使用：
 * import { notificationsApi } from '@/features/notifications/api'
 * 
 * export const useNotificationsStore = defineStore('notifications', () => {
 *   const notifications = ref([])
 *   const unreadCount = ref(0)
 *   const loading = ref(false)
 *   
 *   const fetchNotifications = async (params = {}) => {
 *     loading.value = true
 *     try {
 *       const response = await notificationsApi.getNotifications(params)
 *       notifications.value = response.data.notifications
 *       return response
 *     } catch (error) {
 *       ElMessage.error('获取通知失败')
 *       throw error
 *     } finally {
 *       loading.value = false
 *     }
 *   }
 *   
 *   const fetchUnreadCount = async () => {
 *     try {
 *       const response = await notificationsApi.getUnreadCount()
 *       unreadCount.value = response.data.count
 *       return response
 *     } catch (error) {
 *       console.error('获取未读数量失败:', error)
 *     }
 *   }
 *   
 *   const markAsRead = async (notificationId) => {
 *     try {
 *       await notificationsApi.markAsRead(notificationId)
 *       // 更新本地状态
 *       const notification = notifications.value.find(n => n.id === notificationId)
 *       if (notification && !notification.read) {
 *         notification.read = true
 *         unreadCount.value = Math.max(0, unreadCount.value - 1)
 *       }
 *       return true
 *     } catch (error) {
 *       ElMessage.error('标记已读失败')
 *       throw error
 *     }
 *   }
 *   
 *   return {
 *     notifications,
 *     unreadCount,
 *     loading,
 *     fetchNotifications,
 *     fetchUnreadCount,
 *     markAsRead
 *   }
 * })
 * 
 * 在组件中使用：
 * import { notificationsApi } from '@/features/notifications/api'
 * import { ElMessage } from 'element-plus'
 * 
 * const handleBatchRead = async (notificationIds) => {
 *   try {
 *     await notificationsApi.batchMarkAsRead(notificationIds)
 *     ElMessage.success('批量标记成功')
 *     // 刷新通知列表
 *     await fetchNotifications()
 *   } catch (error) {
 *     ElMessage.error('批量标记失败')
 *   }
 * }
 * 
 * 实时通知WebSocket使用示例：
 * const ws = notificationsApi.connectWebSocket()
 * ws.onmessage = (event) => {
 *   const notification = JSON.parse(event.data)
 *   // 处理实时通知
 *   notifications.value.unshift(notification)
 *   unreadCount.value += 1
 * }
 */