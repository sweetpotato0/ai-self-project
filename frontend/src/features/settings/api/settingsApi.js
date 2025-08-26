import api from '@/api/index'

/**
 * 系统设置模块API接口
 * 包含用户设置、界面设置、通知设置、隐私设置等功能
 */
export const settingsApi = {
  /**
   * 获取用户所有设置
   * @param {string} section - 设置分组 (profile|interface|notifications|privacy|security)
   * @returns {Promise} 设置数据
   */
  getSettings(section = null) {
    const params = section ? { section } : {}
    return api.get('/settings', { params })
  },

  /**
   * 批量更新设置
   * @param {Object} settingsData - 设置数据
   * @returns {Promise} 更新结果
   */
  updateSettings(settingsData) {
    return api.put('/settings', settingsData)
  },

  /**
   * 重置设置为默认值
   * @param {string} section - 要重置的设置分组
   * @returns {Promise} 重置结果
   */
  resetSettings(section = null) {
    return api.post('/settings/reset', { section })
  },

  // ========== 个人资料设置 ==========
  
  /**
   * 获取个人资料设置
   * @returns {Promise} 个人资料数据
   */
  getProfileSettings() {
    return api.get('/settings/profile')
  },

  /**
   * 更新个人资料
   * @param {Object} profileData - 个人资料数据
   * @param {string} profileData.username - 用户名
   * @param {string} profileData.email - 邮箱
   * @param {string} profileData.firstName - 名
   * @param {string} profileData.lastName - 姓
   * @param {string} profileData.bio - 个人简介
   * @param {string} profileData.location - 所在地
   * @param {string} profileData.website - 个人网站
   * @param {string} profileData.timezone - 时区
   * @returns {Promise} 更新结果
   */
  updateProfile(profileData) {
    return api.put('/settings/profile', profileData)
  },

  /**
   * 上传头像
   * @param {File} avatarFile - 头像文件
   * @returns {Promise} 上传结果包含头像URL
   */
  uploadAvatar(avatarFile) {
    const formData = new FormData()
    formData.append('avatar', avatarFile)
    
    return api.post('/settings/profile/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 删除头像
   * @returns {Promise} 删除结果
   */
  deleteAvatar() {
    return api.delete('/settings/profile/avatar')
  },

  // ========== 安全设置 ==========
  
  /**
   * 修改密码
   * @param {Object} passwordData - 密码数据
   * @param {string} passwordData.currentPassword - 当前密码
   * @param {string} passwordData.newPassword - 新密码
   * @param {string} passwordData.confirmPassword - 确认新密码
   * @returns {Promise} 修改结果
   */
  changePassword(passwordData) {
    return api.put('/settings/security/password', {
      current_password: passwordData.currentPassword,
      new_password: passwordData.newPassword,
      confirm_password: passwordData.confirmPassword
    })
  },

  /**
   * 启用/禁用两步验证
   * @param {boolean} enabled - 是否启用
   * @returns {Promise} 设置结果
   */
  updateTwoFactorAuth(enabled) {
    return api.put('/settings/security/2fa', { enabled })
  },

  /**
   * 获取两步验证设置
   * @returns {Promise} 两步验证设置信息
   */
  getTwoFactorAuth() {
    return api.get('/settings/security/2fa')
  },

  /**
   * 生成两步验证密钥
   * @returns {Promise} 密钥和二维码
   */
  generateTwoFactorSecret() {
    return api.post('/settings/security/2fa/generate')
  },

  /**
   * 验证两步验证码
   * @param {string} code - 验证码
   * @returns {Promise} 验证结果
   */
  verifyTwoFactorCode(code) {
    return api.post('/settings/security/2fa/verify', { code })
  },

  /**
   * 获取活跃会话列表
   * @returns {Promise} 会话列表
   */
  getActiveSessions() {
    return api.get('/settings/security/sessions')
  },

  /**
   * 终止会话
   * @param {string} sessionId - 会话ID
   * @returns {Promise} 终止结果
   */
  terminateSession(sessionId) {
    return api.delete(`/settings/security/sessions/${sessionId}`)
  },

  /**
   * 终止所有其他会话
   * @returns {Promise} 终止结果
   */
  terminateAllOtherSessions() {
    return api.delete('/settings/security/sessions/others')
  },

  // ========== 界面设置 ==========
  
  /**
   * 获取界面设置
   * @returns {Promise} 界面设置数据
   */
  getInterfaceSettings() {
    return api.get('/settings/interface')
  },

  /**
   * 更新界面设置
   * @param {Object} interfaceData - 界面设置数据
   * @param {string} interfaceData.theme - 主题 (light|dark|auto)
   * @param {string} interfaceData.language - 语言
   * @param {string} interfaceData.timezone - 时区
   * @param {string} interfaceData.dateFormat - 日期格式
   * @param {string} interfaceData.timeFormat - 时间格式
   * @param {boolean} interfaceData.compactMode - 紧凑模式
   * @param {boolean} interfaceData.showSidebar - 显示侧边栏
   * @param {number} interfaceData.itemsPerPage - 每页显示条数
   * @returns {Promise} 更新结果
   */
  updateInterfaceSettings(interfaceData) {
    return api.put('/settings/interface', interfaceData)
  },

  // ========== 通知设置 ==========
  
  /**
   * 获取通知设置
   * @returns {Promise} 通知设置数据
   */
  getNotificationSettings() {
    return api.get('/settings/notifications')
  },

  /**
   * 更新通知设置
   * @param {Object} notificationData - 通知设置数据
   * @param {boolean} notificationData.emailEnabled - 邮件通知
   * @param {boolean} notificationData.pushEnabled - 推送通知
   * @param {boolean} notificationData.desktopEnabled - 桌面通知
   * @param {Object} notificationData.types - 通知类型设置
   * @returns {Promise} 更新结果
   */
  updateNotificationSettings(notificationData) {
    return api.put('/settings/notifications', notificationData)
  },

  /**
   * 测试通知
   * @param {string} type - 通知类型 (email|push|desktop)
   * @returns {Promise} 测试结果
   */
  testNotification(type) {
    return api.post('/settings/notifications/test', { type })
  },

  // ========== 隐私设置 ==========
  
  /**
   * 获取隐私设置
   * @returns {Promise} 隐私设置数据
   */
  getPrivacySettings() {
    return api.get('/settings/privacy')
  },

  /**
   * 更新隐私设置
   * @param {Object} privacyData - 隐私设置数据
   * @param {boolean} privacyData.profilePublic - 公开个人资料
   * @param {boolean} privacyData.showEmail - 显示邮箱
   * @param {boolean} privacyData.showActivity - 显示活动状态
   * @param {boolean} privacyData.allowMessages - 允许私信
   * @param {string} privacyData.dataRetention - 数据保留期限
   * @returns {Promise} 更新结果
   */
  updatePrivacySettings(privacyData) {
    return api.put('/settings/privacy', privacyData)
  },

  // ========== 数据管理 ==========
  
  /**
   * 导出用户数据
   * @param {Object} exportOptions - 导出选项
   * @param {Array} exportOptions.types - 要导出的数据类型
   * @param {string} exportOptions.format - 导出格式 (json|csv|xml)
   * @returns {Promise} 导出文件
   */
  exportData(exportOptions = {}) {
    return api.post('/settings/data/export', exportOptions, {
      responseType: 'blob'
    })
  },

  /**
   * 获取数据使用统计
   * @returns {Promise} 数据统计信息
   */
  getDataUsage() {
    return api.get('/settings/data/usage')
  },

  /**
   * 清理数据
   * @param {Object} cleanOptions - 清理选项
   * @param {boolean} cleanOptions.completedTasks - 清理已完成任务
   * @param {boolean} cleanOptions.oldNotifications - 清理旧通知
   * @param {boolean} cleanOptions.trashItems - 清理回收站
   * @param {number} cleanOptions.daysOld - 清理多少天前的数据
   * @returns {Promise} 清理结果
   */
  cleanupData(cleanOptions) {
    return api.post('/settings/data/cleanup', cleanOptions)
  },

  /**
   * 删除账户
   * @param {Object} deleteData - 删除确认数据
   * @param {string} deleteData.password - 当前密码
   * @param {string} deleteData.reason - 删除原因
   * @returns {Promise} 删除结果
   */
  deleteAccount(deleteData) {
    return api.delete('/settings/account', {
      data: deleteData
    })
  },

  // ========== 集成设置 ==========
  
  /**
   * 获取第三方集成设置
   * @returns {Promise} 集成设置列表
   */
  getIntegrations() {
    return api.get('/settings/integrations')
  },

  /**
   * 连接第三方服务
   * @param {string} service - 服务名称
   * @param {Object} credentials - 凭据信息
   * @returns {Promise} 连接结果
   */
  connectIntegration(service, credentials) {
    return api.post(`/settings/integrations/${service}/connect`, credentials)
  },

  /**
   * 断开第三方服务
   * @param {string} service - 服务名称
   * @returns {Promise} 断开结果
   */
  disconnectIntegration(service) {
    return api.delete(`/settings/integrations/${service}`)
  },

  /**
   * 测试第三方服务连接
   * @param {string} service - 服务名称
   * @returns {Promise} 测试结果
   */
  testIntegration(service) {
    return api.post(`/settings/integrations/${service}/test`)
  }
}

export default settingsApi