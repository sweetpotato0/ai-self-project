/**
 * 系统设置模块API导出文件
 * 统一管理系统设置相关的所有API接口
 */

import settingsApi from './settingsApi'

// 导出设置API
export { settingsApi }

// 默认导出所有API
export default {
  settingsApi
}

/**
 * API使用示例：
 * 
 * 在Store中使用：
 * import { settingsApi } from '@/features/settings/api'
 * 
 * export const useSettingsStore = defineStore('settings', () => {
 *   const settings = ref({})
 *   const loading = ref(false)
 *   
 *   const fetchSettings = async (section = null) => {
 *     loading.value = true
 *     try {
 *       const response = await settingsApi.getSettings(section)
 *       settings.value = response.data
 *       return response
 *     } catch (error) {
 *       ElMessage.error('获取设置失败')
 *       throw error
 *     } finally {
 *       loading.value = false
 *     }
 *   }
 *   
 *   const updateSettings = async (settingsData) => {
 *     try {
 *       const response = await settingsApi.updateSettings(settingsData)
 *       settings.value = { ...settings.value, ...settingsData }
 *       ElMessage.success('设置更新成功')
 *       return response
 *     } catch (error) {
 *       ElMessage.error('设置更新失败')
 *       throw error
 *     }
 *   }
 *   
 *   return {
 *     settings,
 *     loading,
 *     fetchSettings,
 *     updateSettings
 *   }
 * })
 * 
 * 在组件中使用：
 * import { settingsApi } from '@/features/settings/api'
 * import { ElMessage } from 'element-plus'
 * 
 * const handlePasswordChange = async (passwordData) => {
 *   try {
 *     await settingsApi.changePassword(passwordData)
 *     ElMessage.success('密码修改成功')
 *   } catch (error) {
 *     ElMessage.error('密码修改失败')
 *   }
 * }
 */