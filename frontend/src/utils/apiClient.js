import axios from 'axios'
import { ElMessage } from 'element-plus'

/**
 * 优化的API客户端类
 * 解决全局状态污染和内存泄漏问题
 */
class ApiClient {
  constructor(config = {}) {
    this.instance = axios.create({
      baseURL: config.baseURL || '/api/v1',
      timeout: config.timeout || 30000,
      headers: {
        'Content-Type': 'application/json',
        ...config.headers
      }
    })

    // 私有状态，避免全局污染
    this.tokenRefreshState = {
      isRefreshing: false,
      failedQueue: [],
      maxQueueSize: 50 // 防止内存泄漏
    }

    this.setupInterceptors()
  }

  setupInterceptors() {
    // 请求拦截器
    this.instance.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(this.handleError(error))
      }
    )

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response) => response.data,
      async (error) => {
        const originalRequest = error.config

        // 处理401错误（token过期）
        if (error.response?.status === 401 && !originalRequest._retry) {
          return this.handleUnauthorized(error, originalRequest)
        }

        return Promise.reject(this.handleError(error))
      }
    )
  }

  async handleUnauthorized(error, originalRequest) {
    const { tokenRefreshState } = this

    // 如果正在刷新token，将请求加入队列
    if (tokenRefreshState.isRefreshing) {
      return new Promise((resolve, reject) => {
        if (tokenRefreshState.failedQueue.length >= tokenRefreshState.maxQueueSize) {
          // 队列满了，直接拒绝以防止内存泄漏
          reject(new Error('请求队列已满，请稍后重试'))
          return
        }

        tokenRefreshState.failedQueue.push({ resolve, reject, originalRequest })
      })
    }

    originalRequest._retry = true
    tokenRefreshState.isRefreshing = true

    const token = localStorage.getItem('token')
    
    if (!token) {
      this.processQueue(new Error('未找到认证令牌'))
      this.redirectToLogin()
      return Promise.reject(error)
    }

    try {
      const response = await axios.post('/api/v1/auth/refresh', { token })
      const { token: newToken, user } = response.data.data

      // 更新存储
      localStorage.setItem('token', newToken)
      localStorage.setItem('user', JSON.stringify(user))

      // 更新axios默认header
      this.instance.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      
      // 处理队列中的请求
      this.processQueue(null, newToken)

      // 重试原始请求
      originalRequest.headers.Authorization = `Bearer ${newToken}`
      return this.instance(originalRequest)

    } catch (refreshError) {
      this.processQueue(refreshError)
      this.showErrorMessage('登录已过期，请重新登录')
      this.redirectToLogin()
      return Promise.reject(refreshError)
    }
  }

  processQueue(error, token = null) {
    const { tokenRefreshState } = this

    tokenRefreshState.failedQueue.forEach(({ resolve, reject, originalRequest }) => {
      if (error) {
        reject(error)
      } else {
        originalRequest.headers.Authorization = `Bearer ${token}`
        resolve(this.instance(originalRequest))
      }
    })
    
    // 清理队列和状态
    tokenRefreshState.failedQueue = []
    tokenRefreshState.isRefreshing = false
  }

  handleError(error) {
    const { response } = error

    if (!response) {
      // 网络错误
      this.showErrorMessage('网络连接失败，请检查网络状态')
      return new Error('NETWORK_ERROR')
    }

    const { status, data } = response

    switch (status) {
      case 400:
        this.showErrorMessage(data.message || '请求参数错误')
        break
      case 403:
        this.showErrorMessage('权限不足')
        break
      case 404:
        this.showErrorMessage('请求的资源不存在')
        break
      case 429:
        this.showErrorMessage('请求过于频繁，请稍后重试')
        break
      case 500:
        this.showErrorMessage('服务器内部错误')
        break
      case 502:
      case 503:
      case 504:
        this.showErrorMessage('服务暂时不可用，请稍后重试')
        break
      default:
        this.showErrorMessage(data.message || '请求失败')
    }

    return error
  }

  showErrorMessage(message) {
    ElMessage.error(message)
  }

  redirectToLogin() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    // 避免在非浏览器环境中调用window
    if (typeof window !== 'undefined') {
      window.location.href = '/login'
    }
  }

  // 清理方法，防止内存泄漏
  destroy() {
    this.tokenRefreshState.failedQueue = []
    this.tokenRefreshState.isRefreshing = false
    
    // 移除拦截器
    this.instance.interceptors.request.clear()
    this.instance.interceptors.response.clear()
  }

  // 代理方法
  get(url, config) {
    return this.instance.get(url, config)
  }

  post(url, data, config) {
    return this.instance.post(url, data, config)
  }

  put(url, data, config) {
    return this.instance.put(url, data, config)
  }

  delete(url, config) {
    return this.instance.delete(url, config)
  }

  patch(url, data, config) {
    return this.instance.patch(url, data, config)
  }

  // 获取实例用于特殊需求
  getInstance() {
    return this.instance
  }
}

// 创建默认实例
const apiClient = new ApiClient()

// 导出实例和类
export { ApiClient }
export default apiClient