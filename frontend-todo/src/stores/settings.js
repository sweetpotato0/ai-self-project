import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus'
import {
  getSettings,
  updateProfile,
  changePassword,
  updateNotificationSettings,
  updateInterfaceSettings,
  exportData,
  clearCompletedTasks
} from '@/api/settings'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: {
      id: null,
      user_id: null,
      due_reminder: true,
      completion_notification: true,
      new_task_notification: true,
      email_notification: false,
      theme: 'light',
      language: 'zh-CN',
      timezone: 'Asia/Shanghai'
    },
    loading: false
  }),

  actions: {
    // 获取设置
    async fetchSettings() {
      this.loading = true
      try {
        const response = await getSettings()
        if (response.data) {
          this.settings = response.data
        }
      } catch (error) {
        console.error('获取设置失败:', error)
        ElMessage.error('获取设置失败')
      } finally {
        this.loading = false
      }
    },

    // 更新个人资料
    async updateUserProfile(profileData) {
      this.loading = true
      try {
        await updateProfile(profileData)
        ElMessage.success('个人资料更新成功')
        return true
      } catch (error) {
        console.error('更新个人资料失败:', error)
        ElMessage.error(error.response?.data?.error || '更新失败')
        return false
      } finally {
        this.loading = false
      }
    },

    // 修改密码
    async changeUserPassword(passwordData) {
      this.loading = true
      try {
        await changePassword(passwordData)
        ElMessage.success('密码修改成功')
        return true
      } catch (error) {
        console.error('修改密码失败:', error)
        ElMessage.error(error.response?.data?.error || '密码修改失败')
        return false
      } finally {
        this.loading = false
      }
    },

    // 保存通知设置
    async saveNotificationSettings(notificationData) {
      this.loading = true
      try {
        await updateNotificationSettings(notificationData)
        this.settings = { ...this.settings, ...notificationData }
        ElMessage.success('通知设置已保存')
        return true
      } catch (error) {
        console.error('保存通知设置失败:', error)
        ElMessage.error('保存通知设置失败')
        return false
      } finally {
        this.loading = false
      }
    },

    // 保存界面设置
    async saveInterfaceSettings(interfaceData) {
      this.loading = true
      try {
        await updateInterfaceSettings(interfaceData)
        this.settings = { ...this.settings, ...interfaceData }
        
        // 应用主题设置
        this.applyTheme(interfaceData.theme)
        
        ElMessage.success('界面设置已保存')
        return true
      } catch (error) {
        console.error('保存界面设置失败:', error)
        ElMessage.error('保存界面设置失败')
        return false
      } finally {
        this.loading = false
      }
    },

    // 导出数据
    async exportUserData() {
      this.loading = true
      try {
        const response = await exportData()
        
        // 创建下载链接
        const blob = new Blob([response.data], { type: 'application/json' })
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `taskmaster-export-${new Date().toISOString().split('T')[0]}.json`
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        window.URL.revokeObjectURL(url)
        
        ElMessage.success('数据导出成功')
        return true
      } catch (error) {
        console.error('导出数据失败:', error)
        ElMessage.error('导出数据失败')
        return false
      } finally {
        this.loading = false
      }
    },

    // 清理已完成任务
    async clearCompleted() {
      this.loading = true
      try {
        await clearCompletedTasks()
        ElMessage.success('已完成任务清理成功')
        return true
      } catch (error) {
        console.error('清理任务失败:', error)
        ElMessage.error('清理任务失败')
        return false
      } finally {
        this.loading = false
      }
    },

    // 应用主题
    applyTheme(theme) {
      const root = document.documentElement
      
      if (theme === 'dark') {
        root.classList.add('dark')
        localStorage.setItem('theme', 'dark')
      } else if (theme === 'light') {
        root.classList.remove('dark')
        localStorage.setItem('theme', 'light')
      } else if (theme === 'auto') {
        const isDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches
        if (isDarkMode) {
          root.classList.add('dark')
        } else {
          root.classList.remove('dark')
        }
        localStorage.setItem('theme', 'auto')
      }
    },

    // 初始化设置（应用程序启动时调用）
    async initializeSettings() {
      try {
        await this.fetchSettings()
        
        // 应用保存的主题设置
        const savedTheme = localStorage.getItem('theme') || this.settings.theme
        this.applyTheme(savedTheme)
        
        // 监听系统主题变化（当设置为auto时）
        if (savedTheme === 'auto') {
          const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
          mediaQuery.addEventListener('change', (e) => {
            if (this.settings.theme === 'auto') {
              this.applyTheme('auto')
            }
          })
        }
      } catch (error) {
        console.error('初始化设置失败:', error)
      }
    }
  },

  getters: {
    // 获取通知设置
    notificationSettings: (state) => ({
      due_reminder: state.settings.due_reminder,
      completion_notification: state.settings.completion_notification,
      new_task_notification: state.settings.new_task_notification,
      email_notification: state.settings.email_notification
    }),

    // 获取界面设置
    interfaceSettings: (state) => ({
      theme: state.settings.theme,
      language: state.settings.language,
      timezone: state.settings.timezone
    }),

    // 是否为暗色主题
    isDarkTheme: (state) => {
      if (state.settings.theme === 'dark') return true
      if (state.settings.theme === 'light') return false
      // auto模式下检查系统设置
      return window.matchMedia('(prefers-color-scheme: dark)').matches
    }
  }
})