import api from './index'

export const notificationApi = {
  // 获取通知列表
  getNotifications(params = {}) {
    return api.get('/notifications', { params })
  },

  // 获取未读通知数量
  getUnreadCount() {
    return api.get('/notifications/unread-count')
  },

  // 标记通知为已读
  markAsRead(id) {
    return api.put(`/notifications/${id}/read`)
  },

  // 标记所有通知为已读
  markAllAsRead() {
    return api.put('/notifications/mark-all-read')
  },

  // 删除通知
  deleteNotification(id) {
    return api.delete(`/notifications/${id}`)
  }
}
