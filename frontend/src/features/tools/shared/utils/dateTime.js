/**
 * 时间格式化工具函数
 */

/**
 * 格式化相对时间（如：5分钟前、2小时前、3天前）
 * @param {string|Date} date - 日期字符串或Date对象
 * @returns {string} 格式化后的相对时间
 */
export function formatRelativeTime(date) {
  // 确保date是有效的Date对象
  let dateObj
  if (typeof date === 'string') {
    dateObj = new Date(date)
  } else if (date instanceof Date) {
    dateObj = date
  } else {
    console.warn('Invalid date format:', date)
    return '时间未知'
  }

  // 检查日期是否有效
  if (isNaN(dateObj.getTime())) {
    console.warn('Invalid date:', date)
    return '时间未知'
  }

  const now = new Date()
  const diff = now.getTime() - dateObj.getTime()

  // 处理负数时间差（未来时间）
  if (diff < 0) {
    return '刚刚'
  }

  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 1) {
    return '刚刚'
  } else if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else if (days < 30) {
    return `${days}天前`
  } else {
    // 超过30天显示具体日期
    return dateObj.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }
}

/**
 * 格式化日期时间（如：2024-01-15 14:30）
 * @param {string|Date} dateTime - 日期时间字符串或Date对象
 * @returns {string} 格式化后的日期时间
 */
export function formatDateTime(dateTime) {
  if (!dateTime) return ''

  // 确保dateTime是有效的Date对象
  let dateObj
  if (typeof dateTime === 'string') {
    dateObj = new Date(dateTime)
  } else if (dateTime instanceof Date) {
    dateObj = dateTime
  } else {
    console.warn('Invalid dateTime format:', dateTime)
    return '时间未知'
  }

  // 检查日期是否有效
  if (isNaN(dateObj.getTime())) {
    console.warn('Invalid dateTime:', dateTime)
    return '时间未知'
  }

  return dateObj.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * 格式化日期（如：2024-01-15）
 * @param {string|Date} date - 日期字符串或Date对象
 * @returns {string} 格式化后的日期
 */
export function formatDate(date) {
  if (!date) return ''

  // 确保date是有效的Date对象
  let dateObj
  if (typeof date === 'string') {
    dateObj = new Date(date)
  } else if (date instanceof Date) {
    dateObj = date
  } else {
    console.warn('Invalid date format:', date)
    return '日期未知'
  }

  // 检查日期是否有效
  if (isNaN(dateObj.getTime())) {
    console.warn('Invalid date:', date)
    return '日期未知'
  }

  return dateObj.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

/**
 * 格式化时间（如：14:30）
 * @param {string|Date} time - 时间字符串或Date对象
 * @returns {string} 格式化后的时间
 */
export function formatTime(time) {
  if (!time) return ''

  // 确保time是有效的Date对象
  let timeObj
  if (typeof time === 'string') {
    timeObj = new Date(time)
  } else if (time instanceof Date) {
    timeObj = time
  } else {
    console.warn('Invalid time format:', time)
    return '时间未知'
  }

  // 检查时间是否有效
  if (isNaN(timeObj.getTime())) {
    console.warn('Invalid time:', time)
    return '时间未知'
  }

  return timeObj.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * 检查日期是否有效
 * @param {string|Date} date - 日期字符串或Date对象
 * @returns {boolean} 是否有效
 */
export function isValidDate(date) {
  if (!date) return false

  let dateObj
  if (typeof date === 'string') {
    dateObj = new Date(date)
  } else if (date instanceof Date) {
    dateObj = date
  } else {
    return false
  }

  return !isNaN(dateObj.getTime())
}

/**
 * 获取两个日期之间的天数差
 * @param {string|Date} date1 - 第一个日期
 * @param {string|Date} date2 - 第二个日期
 * @returns {number} 天数差
 */
export function getDaysDiff(date1, date2) {
  if (!isValidDate(date1) || !isValidDate(date2)) {
    return 0
  }

  const d1 = new Date(date1)
  const d2 = new Date(date2)
  const diffTime = Math.abs(d2.getTime() - d1.getTime())
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

/**
 * 检查日期是否过期
 * @param {string|Date} dueDate - 截止日期
 * @returns {boolean} 是否过期
 */
export function isOverdue(dueDate) {
  if (!isValidDate(dueDate)) {
    return false
  }
  return new Date(dueDate) < new Date()
}
