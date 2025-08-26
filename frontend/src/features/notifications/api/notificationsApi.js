import api from '@/api/index'

/**
 * 通知系统模块API接口
 * 包含通知的CRUD操作、推送设置、订阅管理等功能
 */
export const notificationsApi = {
  /**
   * 获取通知列表
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.limit - 每页数量
   * @param {string} params.type - 通知类型筛选 (system|task|article|user)
   * @param {string} params.status - 状态筛选 (read|unread|all)
   * @param {string} params.priority - 优先级筛选 (low|medium|high|urgent)
   * @param {string} params.startDate - 开始日期
   * @param {string} params.endDate - 结束日期
   * @returns {Promise} 通知列表和分页信息
   */
  getNotifications(params = {}) {
    return api.get('/notifications', { params })
  },

  /**
   * 获取单个通知详情
   * @param {number|string} notificationId - 通知ID
   * @returns {Promise} 通知详细信息
   */
  getNotification(notificationId) {
    return api.get(`/notifications/${notificationId}`)
  },

  /**
   * 创建通知
   * @param {Object} notificationData - 通知数据
   * @param {string} notificationData.title - 标题
   * @param {string} notificationData.content - 内容
   * @param {string} notificationData.type - 类型
   * @param {string} notificationData.priority - 优先级
   * @param {Array} notificationData.recipients - 接收者列表
   * @param {string} notificationData.scheduledAt - 定时发送时间
   * @returns {Promise} 创建的通知信息
   */
  createNotification(notificationData) {
    return api.post('/notifications', {
      title: notificationData.title,
      content: notificationData.content,
      type: notificationData.type || 'system',
      priority: notificationData.priority || 'medium',
      recipients: notificationData.recipients || [],
      scheduled_at: notificationData.scheduledAt,
      action_url: notificationData.actionUrl,
      action_text: notificationData.actionText,
      icon: notificationData.icon,
      image: notificationData.image
    })
  },

  /**
   * 更新通知
   * @param {number|string} notificationId - 通知ID
   * @param {Object} notificationData - 更新数据
   * @returns {Promise} 更新后的通知信息
   */
  updateNotification(notificationId, notificationData) {
    return api.put(`/notifications/${notificationId}`, notificationData)
  },

  /**
   * 删除通知
   * @param {number|string} notificationId - 通知ID
   * @returns {Promise} 删除结果
   */
  deleteNotification(notificationId) {
    return api.delete(`/notifications/${notificationId}`)
  },

  /**
   * 批量删除通知
   * @param {Array} notificationIds - 通知ID数组
   * @returns {Promise} 批量删除结果
   */
  batchDeleteNotifications(notificationIds) {
    return api.delete('/notifications/batch', {
      data: { ids: notificationIds }
    })
  },

  /**
   * 标记通知为已读
   * @param {number|string} notificationId - 通知ID
   * @returns {Promise} 标记结果
   */
  markAsRead(notificationId) {
    return api.patch(`/notifications/${notificationId}/read`)
  },

  /**
   * 标记通知为未读
   * @param {number|string} notificationId - 通知ID
   * @returns {Promise} 标记结果
   */
  markAsUnread(notificationId) {
    return api.patch(`/notifications/${notificationId}/unread`)
  },

  /**
   * 批量标记通知为已读
   * @param {Array} notificationIds - 通知ID数组
   * @returns {Promise} 批量标记结果
   */
  batchMarkAsRead(notificationIds) {
    return api.patch('/notifications/batch/read', {
      ids: notificationIds
    })
  },

  /**
   * 标记所有通知为已读
   * @param {Object} filters - 筛选条件
   * @returns {Promise} 标记结果
   */
  markAllAsRead(filters = {}) {
    return api.patch('/notifications/mark-all-read', filters)
  },

  /**
   * 获取未读通知数量
   * @param {Object} filters - 筛选条件
   * @returns {Promise} 未读数量统计
   */
  getUnreadCount(filters = {}) {
    return api.get('/notifications/unread-count', { params: filters })
  },

  /**
   * 获取通知统计信息
   * @param {Object} params - 统计参数
   * @param {string} params.period - 统计周期 (day|week|month|year)
   * @param {string} params.startDate - 开始日期
   * @param {string} params.endDate - 结束日期
   * @returns {Promise} 统计数据
   */
  getNotificationStats(params = {}) {
    return api.get('/notifications/stats', { params })
  },

  /**
   * 归档通知
   * @param {number|string} notificationId - 通知ID
   * @returns {Promise} 归档结果
   */
  archiveNotification(notificationId) {
    return api.patch(`/notifications/${notificationId}/archive`)
  },

  /**
   * 批量归档通知
   * @param {Array} notificationIds - 通知ID数组
   * @returns {Promise} 批量归档结果
   */
  batchArchiveNotifications(notificationIds) {
    return api.patch('/notifications/batch/archive', {
      ids: notificationIds
    })
  },

  /**
   * 取消归档通知
   * @param {number|string} notificationId - 通知ID
   * @returns {Promise} 取消归档结果
   */
  unarchiveNotification(notificationId) {
    return api.patch(`/notifications/${notificationId}/unarchive`)
  },

  // ========== 推送设置 ==========

  /**
   * 获取推送设置
   * @returns {Promise} 推送设置信息
   */
  getPushSettings() {
    return api.get('/notifications/push/settings')
  },

  /**
   * 更新推送设置
   * @param {Object} settingsData - 推送设置数据
   * @param {boolean} settingsData.enabled - 是否启用推送
   * @param {Object} settingsData.types - 推送类型设置
   * @param {Array} settingsData.schedule - 推送时间安排
   * @returns {Promise} 更新结果
   */
  updatePushSettings(settingsData) {
    return api.put('/notifications/push/settings', settingsData)
  },

  /**
   * 订阅推送服务
   * @param {Object} subscription - 推送订阅信息
   * @returns {Promise} 订阅结果
   */
  subscribePush(subscription) {
    return api.post('/notifications/push/subscribe', subscription)
  },

  /**
   * 取消推送订阅
   * @param {string} endpoint - 订阅端点
   * @returns {Promise} 取消订阅结果
   */
  unsubscribePush(endpoint) {
    return api.post('/notifications/push/unsubscribe', { endpoint })
  },

  /**
   * 测试推送通知
   * @param {Object} testData - 测试数据
   * @returns {Promise} 测试结果
   */
  testPushNotification(testData) {
    return api.post('/notifications/push/test', testData)
  },

  // ========== 邮件通知 ==========

  /**
   * 获取邮件通知设置
   * @returns {Promise} 邮件通知设置
   */
  getEmailSettings() {
    return api.get('/notifications/email/settings')
  },

  /**
   * 更新邮件通知设置
   * @param {Object} emailSettings - 邮件设置数据
   * @param {boolean} emailSettings.enabled - 是否启用邮件通知
   * @param {Object} emailSettings.types - 邮件类型设置
   * @param {string} emailSettings.frequency - 发送频率 (instant|daily|weekly)
   * @param {Array} emailSettings.digestTime - 摘要发送时间
   * @returns {Promise} 更新结果
   */
  updateEmailSettings(emailSettings) {
    return api.put('/notifications/email/settings', emailSettings)
  },

  /**
   * 发送测试邮件
   * @param {string} recipient - 接收者邮箱
   * @returns {Promise} 发送结果
   */
  sendTestEmail(recipient) {
    return api.post('/notifications/email/test', { recipient })
  },

  // ========== 通知模板 ==========

  /**
   * 获取通知模板列表
   * @param {Object} params - 查询参数
   * @returns {Promise} 模板列表
   */
  getNotificationTemplates(params = {}) {
    return api.get('/notifications/templates', { params })
  },

  /**
   * 获取通知模板详情
   * @param {number|string} templateId - 模板ID
   * @returns {Promise} 模板详情
   */
  getNotificationTemplate(templateId) {
    return api.get(`/notifications/templates/${templateId}`)
  },

  /**
   * 创建通知模板
   * @param {Object} templateData - 模板数据
   * @param {string} templateData.name - 模板名称
   * @param {string} templateData.type - 模板类型
   * @param {string} templateData.subject - 主题模板
   * @param {string} templateData.content - 内容模板
   * @returns {Promise} 创建的模板信息
   */
  createNotificationTemplate(templateData) {
    return api.post('/notifications/templates', templateData)
  },

  /**
   * 更新通知模板
   * @param {number|string} templateId - 模板ID
   * @param {Object} templateData - 模板数据
   * @returns {Promise} 更新结果
   */
  updateNotificationTemplate(templateId, templateData) {
    return api.put(`/notifications/templates/${templateId}`, templateData)
  },

  /**
   * 删除通知模板
   * @param {number|string} templateId - 模板ID
   * @returns {Promise} 删除结果
   */
  deleteNotificationTemplate(templateId) {
    return api.delete(`/notifications/templates/${templateId}`)
  },

  /**
   * 预览通知模板
   * @param {number|string} templateId - 模板ID
   * @param {Object} variables - 模板变量
   * @returns {Promise} 预览结果
   */
  previewNotificationTemplate(templateId, variables) {
    return api.post(`/notifications/templates/${templateId}/preview`, { variables })
  },

  // ========== 实时通知 ==========

  /**
   * 连接WebSocket通知
   * @returns {WebSocket} WebSocket连接实例
   */
  connectWebSocket() {
    const wsUrl = `${api.defaults.baseURL.replace('http', 'ws')}/notifications/ws`
    return new WebSocket(wsUrl)
  },

  /**
   * 获取实时通知设置
   * @returns {Promise} 实时通知设置
   */
  getRealtimeSettings() {
    return api.get('/notifications/realtime/settings')
  },

  /**
   * 更新实时通知设置
   * @param {Object} settings - 实时通知设置
   * @returns {Promise} 更新结果
   */
  updateRealtimeSettings(settings) {
    return api.put('/notifications/realtime/settings', settings)
  },

  // ========== 通知日志 ==========

  /**
   * 获取通知发送日志
   * @param {Object} params - 查询参数
   * @returns {Promise} 发送日志列表
   */
  getNotificationLogs(params = {}) {
    return api.get('/notifications/logs', { params })
  },

  /**
   * 重试失败的通知
   * @param {number|string} logId - 日志ID
   * @returns {Promise} 重试结果
   */
  retryNotification(logId) {
    return api.post(`/notifications/logs/${logId}/retry`)
  },

  /**
   * 清理通知日志
   * @param {Object} cleanupOptions - 清理选项
   * @param {number} cleanupOptions.daysOld - 清理多少天前的日志
   * @param {string} cleanupOptions.status - 状态筛选
   * @returns {Promise} 清理结果
   */
  cleanupNotificationLogs(cleanupOptions) {
    return api.post('/notifications/logs/cleanup', cleanupOptions)
  }
}

export default notificationsApi