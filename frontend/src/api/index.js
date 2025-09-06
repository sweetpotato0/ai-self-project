import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000, // 增加到30秒，适配DNS查询的较长响应时间
  headers: {
    'Content-Type': 'application/json'
  }
})

// 正在刷新token的标志
let isRefreshing = false
// 存储等待的请求
let failedQueue = []

// 处理队列中的请求
const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })
  
  failedQueue = []
}

// 解析JWT token获取过期时间
const parseJWT = (token) => {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    }).join(''))
    return JSON.parse(jsonPayload)
  } catch (error) {
    return null
  }
}

// 检查token是否即将过期（30分钟内）
const isTokenExpiringSoon = (token) => {
  const payload = parseJWT(token)
  if (!payload || !payload.exp) return false
  
  const now = Date.now() / 1000
  const timeToExpiry = payload.exp - now
  
  // 如果30分钟内过期，返回true
  return timeToExpiry < 30 * 60
}

// 主动刷新token
const proactiveRefreshToken = async () => {
  const token = localStorage.getItem('token')
  if (!token || isRefreshing) return
  
  if (isTokenExpiringSoon(token)) {
    try {
      isRefreshing = true
      const response = await axios.post('/api/v1/auth/refresh', { token })
      const newToken = response.data.data.token
      const user = response.data.data.user
      
      localStorage.setItem('token', newToken)
      localStorage.setItem('user', JSON.stringify(user))
      api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      
      console.log('Token主动刷新成功')
    } catch (error) {
      console.error('Token主动刷新失败:', error)
    } finally {
      isRefreshing = false
    }
  }
}

// 每5分钟检查一次token是否需要刷新
setInterval(proactiveRefreshToken, 5 * 60 * 1000)

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  async (error) => {
    const originalRequest = error.config

    if (error.response) {
      const { status, data } = error.response

      if (status === 401 && !originalRequest._retry) {
        if (isRefreshing) {
          // 如果正在刷新token，将请求加入队列
          return new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject })
          }).then(token => {
            originalRequest.headers.Authorization = `Bearer ${token}`
            return api(originalRequest)
          }).catch(err => {
            return Promise.reject(err)
          })
        }

        originalRequest._retry = true
        isRefreshing = true

        const token = localStorage.getItem('token')
        
        if (!token) {
          // 没有token，直接跳转登录
          processQueue(error, null)
          ElMessage.error('登录已过期，请重新登录')
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          return Promise.reject(error)
        }

        try {
          // 尝试刷新token
          const response = await axios.post('/api/v1/auth/refresh', { token })
          const newToken = response.data.data.token
          const user = response.data.data.user

          // 更新localStorage
          localStorage.setItem('token', newToken)
          localStorage.setItem('user', JSON.stringify(user))

          // 更新axios默认header
          api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
          
          // 处理队列中的请求
          processQueue(null, newToken)

          // 重试原始请求
          originalRequest.headers.Authorization = `Bearer ${newToken}`
          return api(originalRequest)

        } catch (refreshError) {
          // 刷新失败，清除数据并跳转登录
          processQueue(refreshError, null)
          ElMessage.error('登录已过期，请重新登录')
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          return Promise.reject(refreshError)
        } finally {
          isRefreshing = false
        }
      }

      // 处理其他状态码
      switch (status) {
        case 403:
          ElMessage.error('权限不足')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          if (status !== 401) { // 401已经在上面处理了
            ElMessage.error(data.message || '请求失败')
          }
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

export default api
