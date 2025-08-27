import api from '@/api/index'

/**
 * 认证模块API接口
 * 包含用户注册、登录、个人资料管理等功能
 */
export const authApi = {
  /**
   * 用户注册
   * @param {Object} userData - 注册数据
   * @param {string} userData.username - 用户名
   * @param {string} userData.email - 邮箱
   * @param {string} userData.password - 密码
   * @param {string} userData.confirmPassword - 确认密码
   * @returns {Promise} 注册结果
   */
  register(userData) {
    return api.post('/users/register', {
      username: userData.username,
      email: userData.email,
      password: userData.password,
      confirm_password: userData.confirmPassword
    })
  },

  /**
   * 用户登录
   * @param {Object} credentials - 登录凭据
   * @param {string} credentials.email - 邮箱或用户名
   * @param {string} credentials.password - 密码
   * @param {boolean} credentials.rememberMe - 记住我
   * @returns {Promise} 登录结果包含token和用户信息
   */
  login(credentials) {
    return api.post('/users/login', {
      username: credentials.username || credentials.email,
      password: credentials.password,
      remember_me: credentials.rememberMe || false
    })
  },

  /**
   * 用户登出
   * @returns {Promise} 登出结果
   */
  logout() {
    return api.post('/auth/logout')
  },

  /**
   * 刷新token
   * @param {string} refreshToken - 刷新令牌
   * @returns {Promise} 新的token信息
   */
  refreshToken(refreshToken) {
    return api.post('/auth/refresh', {
      refresh_token: refreshToken
    })
  },

  /**
   * 获取当前用户信息
   * @returns {Promise} 用户详细信息
   */
  getCurrentUser() {
    return api.get('/users/profile')
  },

  /**
   * 更新用户个人资料
   * @param {Object} profileData - 个人资料数据
   * @param {string} profileData.username - 用户名
   * @param {string} profileData.email - 邮箱
   * @param {string} profileData.avatar - 头像URL
   * @param {string} profileData.bio - 个人简介
   * @returns {Promise} 更新结果
   */
  updateProfile(profileData) {
    return api.put('/users/profile', profileData)
  },

  /**
   * 修改密码
   * @param {Object} passwordData - 密码数据
   * @param {string} passwordData.currentPassword - 当前密码
   * @param {string} passwordData.newPassword - 新密码
   * @param {string} passwordData.confirmPassword - 确认新密码
   * @returns {Promise} 修改结果
   */
  changePassword(passwordData) {
    return api.put('/auth/password', {
      current_password: passwordData.currentPassword,
      new_password: passwordData.newPassword,
      confirm_password: passwordData.confirmPassword
    })
  },

  /**
   * 上传用户头像
   * @param {File} file - 头像文件
   * @returns {Promise} 上传结果包含新头像URL
   */
  uploadAvatar(file) {
    const formData = new FormData()
    formData.append('avatar', file)
    
    return api.post('/auth/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 发送密码重置邮件
   * @param {string} email - 邮箱地址
   * @returns {Promise} 发送结果
   */
  sendPasswordResetEmail(email) {
    return api.post('/auth/password/reset', {
      email: email
    })
  },

  /**
   * 重置密码
   * @param {Object} resetData - 重置数据
   * @param {string} resetData.token - 重置令牌
   * @param {string} resetData.email - 邮箱
   * @param {string} resetData.password - 新密码
   * @param {string} resetData.confirmPassword - 确认密码
   * @returns {Promise} 重置结果
   */
  resetPassword(resetData) {
    return api.post('/auth/password/reset/confirm', {
      token: resetData.token,
      email: resetData.email,
      password: resetData.password,
      confirm_password: resetData.confirmPassword
    })
  },

  /**
   * 发送邮箱验证邮件
   * @returns {Promise} 发送结果
   */
  sendEmailVerification() {
    return api.post('/auth/email/verify')
  },

  /**
   * 验证邮箱
   * @param {string} token - 验证令牌
   * @returns {Promise} 验证结果
   */
  verifyEmail(token) {
    return api.post('/auth/email/verify/confirm', {
      token: token
    })
  },

  /**
   * 检查用户名是否可用
   * @param {string} username - 用户名
   * @returns {Promise} 可用性检查结果
   */
  checkUsernameAvailability(username) {
    return api.get(`/auth/check/username/${username}`)
  },

  /**
   * 检查邮箱是否可用
   * @param {string} email - 邮箱地址
   * @returns {Promise} 可用性检查结果
   */
  checkEmailAvailability(email) {
    return api.get(`/auth/check/email/${email}`)
  },

  /**
   * 获取用户会话列表
   * @returns {Promise} 会话列表
   */
  getUserSessions() {
    return api.get('/auth/sessions')
  },

  /**
   * 终止指定会话
   * @param {string} sessionId - 会话ID
   * @returns {Promise} 终止结果
   */
  terminateSession(sessionId) {
    return api.delete(`/auth/sessions/${sessionId}`)
  },

  /**
   * 终止所有其他会话
   * @returns {Promise} 终止结果
   */
  terminateAllOtherSessions() {
    return api.delete('/auth/sessions/others')
  }
}

export default authApi