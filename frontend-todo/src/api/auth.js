import api from './index'

export const authApi = {
  // 用户注册
  register(data) {
    return api.post('/users/register', data)
  },

  // 用户登录
  login(data) {
    return api.post('/users/login', data)
  },

  // 获取用户信息
  getProfile() {
    return api.get('/users/profile')
  },

  // 更新用户信息
  updateProfile(data) {
    return api.put('/users/profile', data)
  }
}
