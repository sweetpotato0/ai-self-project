/**
 * 节日订阅管理系统
 */

// 预定义的节日数据源
export const predefinedSources = [
  {
    id: 'china-legal',
    name: '中国法定节假日',
    url: 'https://api.example.com/holidays/china/legal',
    description: '包含所有中国法定节假日和调休安排',
    type: 'legal',
    enabled: true,
    lastSync: null
  },
  {
    id: 'china-traditional',
    name: '中国传统节日',
    url: 'https://api.example.com/holidays/china/traditional',
    description: '包含春节、端午、中秋等传统节日',
    type: 'traditional',
    enabled: false,
    lastSync: null
  },
  {
    id: 'international',
    name: '国际节日',
    url: 'https://api.example.com/holidays/international',
    description: '包含元旦、情人节、圣诞节等国际节日',
    type: 'international',
    enabled: false,
    lastSync: null
  }
]

// 本地存储键名
const STORAGE_KEYS = {
  CUSTOM_SOURCES: 'holiday_custom_sources',
  PREDEFINED_SOURCES: 'holiday_predefined_sources',
  HOLIDAY_CACHE: 'holiday_data_cache'
}

/**
 * 节日订阅管理类
 */
export class HolidaySubscriptionManager {
  constructor() {
    this.customSources = this.loadCustomSources()
    this.predefinedSources = this.loadPredefinedSources()
    this.holidayCache = this.loadHolidayCache()
  }

  // 加载自定义数据源
  loadCustomSources() {
    const saved = localStorage.getItem(STORAGE_KEYS.CUSTOM_SOURCES)
    return saved ? JSON.parse(saved) : []
  }

  // 加载预定义数据源状态
  loadPredefinedSources() {
    const saved = localStorage.getItem(STORAGE_KEYS.PREDEFINED_SOURCES)
    if (saved) {
      const savedSources = JSON.parse(saved)
      // 合并预定义源和保存的状态
      return predefinedSources.map(source => {
        const savedSource = savedSources.find(s => s.id === source.id)
        return savedSource ? { ...source, ...savedSource } : source
      })
    }
    return [...predefinedSources]
  }

  // 加载节日缓存数据
  loadHolidayCache() {
    const saved = localStorage.getItem(STORAGE_KEYS.HOLIDAY_CACHE)
    return saved ? JSON.parse(saved) : {}
  }

  // 保存自定义数据源
  saveCustomSources() {
    localStorage.setItem(STORAGE_KEYS.CUSTOM_SOURCES, JSON.stringify(this.customSources))
  }

  // 保存预定义数据源状态
  savePredefinedSources() {
    localStorage.setItem(STORAGE_KEYS.PREDEFINED_SOURCES, JSON.stringify(this.predefinedSources))
  }

  // 保存节日缓存
  saveHolidayCache() {
    localStorage.setItem(STORAGE_KEYS.HOLIDAY_CACHE, JSON.stringify(this.holidayCache))
  }

  // 添加自定义订阅源
  addCustomSource(sourceData) {
    const newSource = {
      id: `custom_${Date.now()}`,
      name: sourceData.name,
      url: sourceData.url,
      description: sourceData.description || '',
      type: sourceData.type || 'custom',
      enabled: true,
      lastSync: null,
      isCustom: true
    }

    this.customSources.push(newSource)
    this.saveCustomSources()
    return newSource
  }

  // 删除自定义订阅源
  removeCustomSource(sourceId) {
    this.customSources = this.customSources.filter(source => source.id !== sourceId)
    this.saveCustomSources()
    // 同时删除缓存数据
    delete this.holidayCache[sourceId]
    this.saveHolidayCache()
  }

  // 启用/禁用数据源
  toggleSource(sourceId, enabled) {
    // 在自定义源中查找
    const customSource = this.customSources.find(s => s.id === sourceId)
    if (customSource) {
      customSource.enabled = enabled
      this.saveCustomSources()
      return
    }

    // 在预定义源中查找
    const predefinedSource = this.predefinedSources.find(s => s.id === sourceId)
    if (predefinedSource) {
      predefinedSource.enabled = enabled
      this.savePredefinedSources()
    }
  }

  // 获取所有数据源
  getAllSources() {
    return [...this.predefinedSources, ...this.customSources]
  }

  // 获取已启用的数据源
  getEnabledSources() {
    return this.getAllSources().filter(source => source.enabled)
  }

  // 从URL获取节日数据
  async fetchHolidayData(source) {
    try {
      console.log(`正在从 ${source.name} 获取节日数据...`)
      
      // 模拟API调用（实际项目中这里会是真实的fetch请求）
      const response = await this.simulateApiCall(source.url)
      
      if (response.success) {
        const holidays = this.normalizeHolidayData(response.data, source.type)
        
        // 缓存数据
        this.holidayCache[source.id] = {
          data: holidays,
          lastSync: new Date().toISOString(),
          expireAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString() // 24小时过期
        }
        this.saveHolidayCache()

        // 更新数据源的最后同步时间
        const sourceToUpdate = this.getAllSources().find(s => s.id === source.id)
        if (sourceToUpdate) {
          sourceToUpdate.lastSync = new Date().toISOString()
          if (sourceToUpdate.isCustom) {
            this.saveCustomSources()
          } else {
            this.savePredefinedSources()
          }
        }

        console.log(`从 ${source.name} 获取到 ${holidays.length} 个节日`)
        return holidays
      } else {
        throw new Error(response.error || '获取数据失败')
      }
    } catch (error) {
      console.error(`从 ${source.name} 获取节日数据失败:`, error)
      throw error
    }
  }

  // 模拟API调用（实际项目中替换为真实API）
  async simulateApiCall(url) {
    // 模拟网络延迟
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 根据URL返回不同的模拟数据
    if (url.includes('china/legal')) {
      return {
        success: true,
        data: [
          { date: '2025-09-07', name: '中秋节', type: 'legal' },
          { date: '2025-10-01', name: '国庆节', type: 'legal' },
          { date: '2025-10-02', name: '国庆节', type: 'legal' },
          { date: '2025-10-03', name: '国庆节', type: 'legal' }
        ]
      }
    } else if (url.includes('china/traditional')) {
      return {
        success: true,
        data: [
          { date: '2025-09-10', name: '教师节', type: 'traditional' },
          { date: '2025-12-22', name: '冬至', type: 'traditional' },
          { date: '2025-12-24', name: '平安夜', type: 'traditional' }
        ]
      }
    } else if (url.includes('international')) {
      return {
        success: true,
        data: [
          { date: '2025-12-25', name: '圣诞节', type: 'international' },
          { date: '2025-09-21', name: '国际和平日', type: 'international' }
        ]
      }
    } else {
      // 自定义URL - 返回空数据或错误
      return {
        success: false,
        error: '无法获取数据，请检查URL是否正确'
      }
    }
  }

  // 标准化节日数据格式
  normalizeHolidayData(rawData, sourceType) {
    return rawData.map(item => ({
      date: item.date,
      name: item.name,
      type: item.type || sourceType,
      isWorkday: item.isWorkday !== undefined ? item.isWorkday : true,
      source: sourceType
    }))
  }

  // 获取缓存的节日数据
  getCachedHolidayData(sourceId) {
    const cached = this.holidayCache[sourceId]
    if (!cached) return null

    // 检查是否过期
    const now = new Date()
    const expireAt = new Date(cached.expireAt)
    if (now > expireAt) {
      delete this.holidayCache[sourceId]
      this.saveHolidayCache()
      return null
    }

    return cached.data
  }

  // 获取所有节日数据（从缓存或网络）
  async getAllHolidays(forceRefresh = false) {
    const enabledSources = this.getEnabledSources()
    const allHolidays = []

    for (const source of enabledSources) {
      try {
        let holidays = null

        if (!forceRefresh) {
          holidays = this.getCachedHolidayData(source.id)
        }

        if (!holidays) {
          holidays = await this.fetchHolidayData(source)
        }

        if (holidays) {
          allHolidays.push(...holidays)
        }
      } catch (error) {
        console.error(`获取 ${source.name} 的节日数据失败:`, error)
      }
    }

    return allHolidays
  }

  // 根据日期查找节日
  getHolidayInfo(date, holidays = null) {
    if (!holidays) {
      // 如果没有提供节日数据，尝试从缓存获取
      const enabledSources = this.getEnabledSources()
      const allCachedHolidays = []
      
      for (const source of enabledSources) {
        const cached = this.getCachedHolidayData(source.id)
        if (cached) {
          allCachedHolidays.push(...cached)
        }
      }
      holidays = allCachedHolidays
    }

    const dateString = typeof date === 'string' ? date : date.toISOString().split('T')[0]
    return holidays.find(holiday => holiday.date === dateString) || null
  }

  // 清除所有缓存
  clearCache() {
    this.holidayCache = {}
    this.saveHolidayCache()
  }
}

// 导出单例实例
export const holidayManager = new HolidaySubscriptionManager()