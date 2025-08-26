import { ElMessage } from 'element-plus'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import en from 'element-plus/dist/locale/en.mjs'

class LocaleManager {
  constructor() {
    this.currentLanguage = 'zh-CN'
    this.currentTimezone = 'Asia/Shanghai'
    this.locales = {
      'zh-CN': zhCn,
      'en-US': en
    }
  }

  // 更新语言设置
  updateLanguage(language) {
    this.currentLanguage = language
    
    // 更新 Element Plus 语言
    this.applyElementLocale(language)
    
    // 更新 HTML lang 属性
    document.documentElement.lang = language === 'zh-CN' ? 'zh' : 'en'
    
    // 触发语言切换事件
    window.dispatchEvent(new CustomEvent('languageChanged', { 
      detail: { language } 
    }))
  }

  // 应用 Element Plus 语言设置
  applyElementLocale(language) {
    const locale = this.locales[language] || this.locales['zh-CN']
    
    // 这里需要在应用级别重新配置 Element Plus
    // 由于动态切换语言比较复杂，我们暂时通过事件通知其他组件
    window.elementLocale = locale
  }

  // 更新时区设置
  updateTimezone(timezone) {
    this.currentTimezone = timezone
    
    // 触发时区切换事件
    window.dispatchEvent(new CustomEvent('timezoneChanged', { 
      detail: { timezone } 
    }))
  }

  // 格式化日期时间（根据当前时区）
  formatDateTime(dateString, options = {}) {
    if (!dateString) return ''
    
    const date = new Date(dateString)
    
    const defaultOptions = {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: this.currentTimezone,
      ...options
    }

    return new Intl.DateTimeFormat(this.currentLanguage, defaultOptions).format(date)
  }

  // 格式化日期
  formatDate(dateString, options = {}) {
    if (!dateString) return ''
    
    const date = new Date(dateString)
    
    const defaultOptions = {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      timeZone: this.currentTimezone,
      ...options
    }

    return new Intl.DateTimeFormat(this.currentLanguage, defaultOptions).format(date)
  }

  // 格式化相对时间
  formatRelativeTime(dateString) {
    if (!dateString) return ''
    
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now - date
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
    
    if (diffDays === 0) {
      return this.currentLanguage === 'zh-CN' ? '今天' : 'Today'
    } else if (diffDays === 1) {
      return this.currentLanguage === 'zh-CN' ? '昨天' : 'Yesterday'
    } else if (diffDays === -1) {
      return this.currentLanguage === 'zh-CN' ? '明天' : 'Tomorrow'
    } else if (diffDays > 1) {
      return this.currentLanguage === 'zh-CN' ? `${diffDays}天前` : `${diffDays} days ago`
    } else if (diffDays < -1) {
      return this.currentLanguage === 'zh-CN' ? `${Math.abs(diffDays)}天后` : `in ${Math.abs(diffDays)} days`
    }
    
    return this.formatDate(dateString)
  }

  // 获取当前时区的时间
  getCurrentTime() {
    return new Date().toLocaleString(this.currentLanguage, {
      timeZone: this.currentTimezone,
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }

  // 转换时间到当前时区
  convertToCurrentTimezone(dateString) {
    if (!dateString) return ''
    
    const date = new Date(dateString)
    return date.toLocaleString(this.currentLanguage, {
      timeZone: this.currentTimezone
    })
  }

  // 获取时区偏移量
  getTimezoneOffset() {
    const date = new Date()
    const utcDate = new Date(date.toLocaleString('en-US', { timeZone: 'UTC' }))
    const targetDate = new Date(date.toLocaleString('en-US', { timeZone: this.currentTimezone }))
    
    return (targetDate.getTime() - utcDate.getTime()) / (1000 * 60 * 60) // 小时偏移
  }

  // 获取本地化文本
  getText(key, fallback = '') {
    const texts = {
      'zh-CN': {
        'task.created': '任务已创建',
        'task.updated': '任务已更新',
        'task.completed': '任务已完成',
        'task.deleted': '任务已删除',
        'notification.due_soon': '任务即将到期',
        'notification.overdue': '任务已过期',
        'date.today': '今天',
        'date.yesterday': '昨天',
        'date.tomorrow': '明天',
        'time.morning': '上午',
        'time.afternoon': '下午',
        'time.evening': '晚上'
      },
      'en-US': {
        'task.created': 'Task created',
        'task.updated': 'Task updated', 
        'task.completed': 'Task completed',
        'task.deleted': 'Task deleted',
        'notification.due_soon': 'Task due soon',
        'notification.overdue': 'Task overdue',
        'date.today': 'Today',
        'date.yesterday': 'Yesterday',
        'date.tomorrow': 'Tomorrow',
        'time.morning': 'AM',
        'time.afternoon': 'PM',
        'time.evening': 'Evening'
      }
    }

    return texts[this.currentLanguage]?.[key] || fallback || key
  }
}

// 创建全局实例
export const localeManager = new LocaleManager()

export default LocaleManager