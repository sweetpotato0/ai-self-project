import { notificationsApi } from '@/features/notifications/api'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  // 状态
  const notifications = ref([])
  const loading = ref(false)
  const unreadCount = ref(0)

  // 计算属性
  const unreadNotifications = computed(() => {
    return notifications.value.filter(n => !n.is_read)
  })

  const readNotifications = computed(() => {
    return notifications.value.filter(n => n.is_read)
  })

  // 方法
  const fetchNotifications = async (params = {}) => {
    loading.value = true
    try {
      const response = await notificationsApi.getNotifications(params)
      notifications.value = response.data.notifications || []
      unreadCount.value = unreadNotifications.value.length
    } catch (error) {
      console.error('Failed to fetch notifications:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const addNotification = (notification) => {
    // 检查是否已存在相同ID的通知
    const existingIndex = notifications.value.findIndex(n => n.id === notification.id)
    if (existingIndex >= 0) {
      // 更新现有通知
      notifications.value[existingIndex] = notification
    } else {
      // 添加新通知到开头
      notifications.value.unshift(notification)
    }
    unreadCount.value = unreadNotifications.value.length
  }

  const setNotifications = (newNotifications) => {
    notifications.value = newNotifications
    unreadCount.value = unreadNotifications.value.length
  }

  const markAsRead = async (notificationId) => {
    try {
      await notificationsApi.markAsRead(notificationId)
      const notification = notifications.value.find(n => n.id === notificationId)
      if (notification) {
        notification.is_read = true
        unreadCount.value = unreadNotifications.value.length
      }
    } catch (error) {
      console.error('Failed to mark notification as read:', error)
      throw error
    }
  }

  const markAllAsRead = async () => {
    try {
      await notificationsApi.markAllAsRead()
      notifications.value.forEach(n => n.is_read = true)
      unreadCount.value = 0
    } catch (error) {
      console.error('Failed to mark all notifications as read:', error)
      throw error
    }
  }

  const deleteNotification = async (notificationId) => {
    try {
      await notificationsApi.deleteNotification(notificationId)
      const index = notifications.value.findIndex(n => n.id === notificationId)
      if (index > -1) {
        const notification = notifications.value[index]
        notifications.value.splice(index, 1)
        if (!notification.is_read) {
          unreadCount.value = unreadNotifications.value.length
        }
      }
    } catch (error) {
      console.error('Failed to delete notification:', error)
      throw error
    }
  }

  const clearNotifications = () => {
    notifications.value = []
    unreadCount.value = 0
  }

  const updateUnreadCount = (count) => {
    unreadCount.value = count
  }

  return {
    // 状态
    notifications,
    loading,
    unreadCount,

    // 计算属性
    unreadNotifications,
    readNotifications,

    // 方法
    fetchNotifications,
    addNotification,
    setNotifications,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    clearNotifications,
    updateUnreadCount
  }
})
