/**
 * 优化的WebSocket管理器
 * 解决硬编码URL、token安全传输、连接管理等问题
 */
class WebSocketManager {
  constructor(options = {}) {
    this.options = {
      baseURL: this.getWebSocketURL(),
      maxReconnectAttempts: options.maxReconnectAttempts || 10,
      reconnectInterval: options.reconnectInterval || 1000,
      maxReconnectInterval: options.maxReconnectInterval || 30000,
      heartbeatInterval: options.heartbeatInterval || 30000,
      debug: options.debug || false,
      ...options
    }

    this.ws = null
    this.reconnectAttempts = 0
    this.reconnectTimer = null
    this.heartbeatTimer = null
    this.connectionStatus = 'disconnected'
    this.eventHandlers = new Map()
    this.messageQueue = []
    this.isDestroyed = false

    this.setupEventHandlers()
  }

  getWebSocketURL() {
    // 从环境变量或当前域名动态生成WebSocket URL
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = process.env.VUE_APP_WS_HOST || window.location.host
    return `${protocol}//${host}/api/v1/ws`
  }

  setupEventHandlers() {
    // 绑定事件处理器，确保正确的this上下文
    this.onOpen = this.onOpen.bind(this)
    this.onClose = this.onClose.bind(this)
    this.onError = this.onError.bind(this)
    this.onMessage = this.onMessage.bind(this)
  }

  connect(token) {
    if (this.isDestroyed) {
      console.warn('WebSocketManager已被销毁，无法重新连接')
      return
    }

    if (!token) {
      this.log('未提供认证令牌，跳过WebSocket连接')
      return
    }

    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.log('WebSocket已连接')
      return
    }

    this.connectionStatus = 'connecting'
    this.emit('statusChange', { status: 'connecting' })

    try {
      // 安全传输token - 使用消息而不是URL参数
      this.ws = new WebSocket(this.options.baseURL)
      this.ws.onopen = this.onOpen
      this.ws.onclose = this.onClose
      this.ws.onerror = this.onError
      this.ws.onmessage = this.onMessage

      // 存储token用于认证
      this.pendingAuthToken = token

    } catch (error) {
      this.log('WebSocket连接失败:', error)
      this.handleConnectionError(error)
    }
  }

  onOpen() {
    this.log('WebSocket连接已建立')
    this.connectionStatus = 'connected'
    this.reconnectAttempts = 0
    
    // 发送认证消息
    if (this.pendingAuthToken) {
      this.sendMessage({
        type: 'auth',
        token: this.pendingAuthToken
      })
      this.pendingAuthToken = null
    }

    // 发送队列中的消息
    this.flushMessageQueue()
    
    // 启动心跳
    this.startHeartbeat()
    
    this.emit('open')
    this.emit('statusChange', { status: 'connected' })
  }

  onClose(event) {
    this.log('WebSocket连接关闭:', event.code, event.reason)
    this.connectionStatus = 'disconnected'
    this.stopHeartbeat()
    
    this.emit('close', event)
    this.emit('statusChange', { status: 'disconnected' })

    // 自动重连（如果不是主动关闭）
    if (!this.isDestroyed && event.code !== 1000) {
      this.scheduleReconnect()
    }
  }

  onError(error) {
    this.log('WebSocket错误:', error)
    this.emit('error', error)
    this.handleConnectionError(error)
  }

  onMessage(event) {
    try {
      const data = JSON.parse(event.data)
      this.log('收到消息:', data)
      
      // 处理心跳响应
      if (data.type === 'pong') {
        return
      }
      
      // 处理认证响应
      if (data.type === 'auth_response') {
        if (data.success) {
          this.log('认证成功')
          this.emit('authenticated')
        } else {
          this.log('认证失败:', data.message)
          this.emit('authFailed', data.message)
          this.disconnect(1008, '认证失败')
        }
        return
      }

      this.emit('message', data)
      
      // 分发特定类型的消息
      if (data.type) {
        this.emit(data.type, data)
      }
      
    } catch (error) {
      this.log('消息解析失败:', error)
      this.emit('parseError', error)
    }
  }

  sendMessage(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      try {
        this.ws.send(JSON.stringify(message))
        this.log('发送消息:', message)
      } catch (error) {
        this.log('消息发送失败:', error)
        this.emit('sendError', error)
      }
    } else {
      // 连接未建立，加入队列
      if (this.messageQueue.length < 100) { // 防止队列过大
        this.messageQueue.push(message)
        this.log('消息已加入队列:', message)
      } else {
        this.log('消息队列已满，丢弃消息:', message)
      }
    }
  }

  flushMessageQueue() {
    while (this.messageQueue.length > 0) {
      const message = this.messageQueue.shift()
      this.sendMessage(message)
    }
  }

  scheduleReconnect() {
    if (this.reconnectAttempts >= this.options.maxReconnectAttempts) {
      this.log('达到最大重连次数，停止重连')
      this.emit('maxReconnectAttemptsReached')
      return
    }

    const delay = Math.min(
      this.options.reconnectInterval * Math.pow(2, this.reconnectAttempts),
      this.options.maxReconnectInterval
    )

    this.reconnectAttempts++
    this.log(`准备重连 (${this.reconnectAttempts}/${this.options.maxReconnectAttempts})，${delay}ms后`)

    this.reconnectTimer = setTimeout(() => {
      const token = localStorage.getItem('token')
      if (token) {
        this.connect(token)
      } else {
        this.log('未找到认证令牌，停止重连')
      }
    }, delay)
  }

  startHeartbeat() {
    this.stopHeartbeat()
    
    this.heartbeatTimer = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.sendMessage({ type: 'ping', timestamp: Date.now() })
      }
    }, this.options.heartbeatInterval)
  }

  stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  handleConnectionError(error) {
    this.connectionStatus = 'error'
    this.emit('statusChange', { status: 'error', error })
  }

  disconnect(code = 1000, reason = 'Normal closure') {
    this.log('主动断开WebSocket连接')
    
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    
    this.stopHeartbeat()
    
    if (this.ws) {
      this.ws.close(code, reason)
      this.ws = null
    }
    
    this.connectionStatus = 'disconnected'
    this.emit('statusChange', { status: 'disconnected' })
  }

  // 事件系统
  on(event, handler) {
    if (!this.eventHandlers.has(event)) {
      this.eventHandlers.set(event, [])
    }
    this.eventHandlers.get(event).push(handler)
  }

  off(event, handler) {
    if (this.eventHandlers.has(event)) {
      const handlers = this.eventHandlers.get(event)
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
      }
    }
  }

  emit(event, ...args) {
    if (this.eventHandlers.has(event)) {
      this.eventHandlers.get(event).forEach(handler => {
        try {
          handler(...args)
        } catch (error) {
          console.error('WebSocket事件处理器错误:', error)
        }
      })
    }
  }

  // 获取状态
  getStatus() {
    return {
      connectionStatus: this.connectionStatus,
      reconnectAttempts: this.reconnectAttempts,
      readyState: this.ws ? this.ws.readyState : WebSocket.CLOSED,
      queueLength: this.messageQueue.length
    }
  }

  // 清理资源
  destroy() {
    this.log('销毁WebSocketManager')
    this.isDestroyed = true
    
    this.disconnect()
    
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
    }
    
    // 清理事件处理器
    this.eventHandlers.clear()
    this.messageQueue = []
  }

  log(...args) {
    if (this.options.debug) {
      console.log('[WebSocketManager]', ...args)
    }
  }
}

// 创建单例实例
let globalWebSocketManager = null

export function createWebSocketManager(options = {}) {
  if (globalWebSocketManager) {
    globalWebSocketManager.destroy()
  }
  
  globalWebSocketManager = new WebSocketManager(options)
  return globalWebSocketManager
}

export function getWebSocketManager() {
  if (!globalWebSocketManager) {
    globalWebSocketManager = new WebSocketManager()
  }
  return globalWebSocketManager
}

export { WebSocketManager }
export default WebSocketManager