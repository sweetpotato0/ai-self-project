import { authApi } from '@/features/auth/api'
import { ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)

  // 登录
  const login = async (credentials) => {
    try {
      const response = await authApi.login(credentials)
      const { token: newToken, user: userData } = response.data

      token.value = newToken
      user.value = userData

      localStorage.setItem('token', newToken)
      localStorage.setItem('user', JSON.stringify(userData))

      ElMessage.success('登录成功')
      
      // 登录成功后初始化设置
      const { useSettingsStore } = await import('./settings')
      const settingsStore = useSettingsStore()
      settingsStore.initializeSettings()
      
      return true
    } catch (error) {
      ElMessage.error('登录失败')
      return false
    }
  }

  // 注册
  const register = async (userData) => {
    try {
      await authApi.register(userData)
      ElMessage.success('注册成功，请登录')
      return true
    } catch (error) {
      ElMessage.error('注册失败')
      return false
    }
  }

  // 登出
  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    ElMessage.success('已退出登录')
  }

  // 获取用户信息
  const fetchProfile = async () => {
    try {
      const response = await authApi.getProfile()
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(response.data))
    } catch (error) {
      logout()
    }
  }

  // 初始化
  const init = async () => {
    const savedUser = localStorage.getItem('user')
    if (savedUser && token.value) {
      user.value = JSON.parse(savedUser)
      
      // 如果用户已经登录，初始化设置
      const { useSettingsStore } = await import('./settings')
      const settingsStore = useSettingsStore()
      settingsStore.initializeSettings()
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    register,
    logout,
    fetchProfile,
    init
  }
})
