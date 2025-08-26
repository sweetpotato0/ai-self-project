/**
 * 认证模块API导出文件
 * 统一管理认证相关的所有API接口
 */

import authApi from './authApi'

// 导出认证API
export { authApi }

// 默认导出所有API
export default {
  authApi
}

/**
 * API使用示例：
 * 
 * 方式一：按需导入
 * import { authApi } from '@/features/auth/api'
 * const result = await authApi.login({ email, password })
 * 
 * 方式二：默认导入
 * import authApis from '@/features/auth/api'
 * const result = await authApis.authApi.login({ email, password })
 * 
 * 方式三：在组合式API中使用
 * import { authApi } from '@/features/auth/api'
 * 
 * export function useAuth() {
 *   const login = async (credentials) => {
 *     try {
 *       const response = await authApi.login(credentials)
 *       // 处理登录成功逻辑
 *       return response
 *     } catch (error) {
 *       // 处理登录失败逻辑
 *       throw error
 *     }
 *   }
 *   
 *   return { login }
 * }
 */