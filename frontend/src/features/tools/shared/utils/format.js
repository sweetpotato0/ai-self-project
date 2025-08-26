/**
 * 格式化文件大小
 * @param {number} bytes - 字节数
 * @param {number} decimals - 小数位数
 * @returns {string} 格式化后的文件大小
 */
export function formatFileSize(bytes, decimals = 2) {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

/**
 * 格式化数字，添加千位分隔符
 * @param {number} num - 要格式化的数字
 * @returns {string} 格式化后的数字字符串
 */
export function formatNumber(num) {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

/**
 * 格式化百分比
 * @param {number} value - 数值 (0-1)
 * @param {number} decimals - 小数位数
 * @returns {string} 百分比字符串
 */
export function formatPercentage(value, decimals = 1) {
  return (value * 100).toFixed(decimals) + '%'
}

/**
 * 格式化时间戳
 * @param {number|Date} timestamp - 时间戳或Date对象
 * @param {string} format - 格式化字符串
 * @returns {string} 格式化后的时间字符串
 */
export function formatTimestamp(timestamp, format = 'YYYY-MM-DD HH:mm:ss') {
  const date = new Date(timestamp)
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化持续时间
 * @param {number} milliseconds - 毫秒数
 * @returns {string} 格式化后的持续时间
 */
export function formatDuration(milliseconds) {
  const seconds = Math.floor(milliseconds / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  
  if (days > 0) {
    return `${days}天 ${hours % 24}时 ${minutes % 60}分`
  } else if (hours > 0) {
    return `${hours}时 ${minutes % 60}分 ${seconds % 60}秒`
  } else if (minutes > 0) {
    return `${minutes}分 ${seconds % 60}秒`
  } else {
    return `${seconds}秒`
  }
}

/**
 * 截断文本
 * @param {string} text - 要截断的文本
 * @param {number} length - 最大长度
 * @param {string} suffix - 后缀
 * @returns {string} 截断后的文本
 */
export function truncateText(text, length, suffix = '...') {
  if (text.length <= length) return text
  return text.slice(0, length) + suffix
}

/**
 * 格式化JSON字符串
 * @param {any} obj - 要格式化的对象
 * @param {number} spaces - 缩进空格数
 * @returns {string} 格式化后的JSON字符串
 */
export function formatJSON(obj, spaces = 2) {
  try {
    return JSON.stringify(obj, null, spaces)
  } catch (error) {
    return String(obj)
  }
}