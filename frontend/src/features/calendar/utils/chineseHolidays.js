/**
 * 中国节日数据和工具函数
 */

// 2025年中国法定节假日和重要节日数据
export const chineseHolidays2025 = [
  // 法定节假日
  { date: '2025-01-01', name: '元旦', type: 'legal', isWorkday: false },
  { date: '2025-01-28', name: '春节前夜', type: 'legal', isWorkday: false },
  { date: '2025-01-29', name: '春节', type: 'legal', isWorkday: false },
  { date: '2025-01-30', name: '春节', type: 'legal', isWorkday: false },
  { date: '2025-01-31', name: '春节', type: 'legal', isWorkday: false },
  { date: '2025-02-01', name: '春节', type: 'legal', isWorkday: false },
  { date: '2025-02-02', name: '春节', type: 'legal', isWorkday: false },
  { date: '2025-04-05', name: '清明节', type: 'legal', isWorkday: false },
  { date: '2025-04-06', name: '清明节', type: 'legal', isWorkday: false },
  { date: '2025-04-07', name: '清明节', type: 'legal', isWorkday: false },
  { date: '2025-05-01', name: '劳动节', type: 'legal', isWorkday: false },
  { date: '2025-05-02', name: '劳动节', type: 'legal', isWorkday: false },
  { date: '2025-05-03', name: '劳动节', type: 'legal', isWorkday: false },
  { date: '2025-05-04', name: '劳动节', type: 'legal', isWorkday: false },
  { date: '2025-05-05', name: '劳动节', type: 'legal', isWorkday: false },
  { date: '2025-05-31', name: '端午节', type: 'legal', isWorkday: false },
  { date: '2025-06-02', name: '端午节', type: 'legal', isWorkday: false },
  
  // 当前和即将到来的节日（便于测试）
  { date: '2025-08-29', name: '测试节日', type: 'legal', isWorkday: false },
  { date: '2025-08-30', name: '测试传统节', type: 'traditional', isWorkday: true },
  { date: '2025-09-01', name: '开学日', type: 'traditional', isWorkday: true },
  { date: '2025-09-07', name: '中秋节', type: 'legal', isWorkday: false },
  { date: '2025-09-10', name: '教师节', type: 'traditional', isWorkday: true },
  { date: '2025-10-01', name: '国庆节', type: 'legal', isWorkday: false },
  { date: '2025-10-02', name: '国庆节', type: 'legal', isWorkday: false },
  { date: '2025-10-03', name: '国庆节', type: 'legal', isWorkday: false },
  { date: '2025-10-04', name: '国庆节', type: 'legal', isWorkday: false },
  { date: '2025-10-05', name: '国庆节', type: 'legal', isWorkday: false },
  { date: '2025-10-06', name: '国庆节', type: 'legal', isWorkday: false },
  { date: '2025-10-07', name: '国庆节', type: 'legal', isWorkday: false },
  { date: '2025-10-08', name: '国庆节', type: 'legal', isWorkday: false },

  // 调休工作日
  { date: '2025-01-26', name: '春节调休', type: 'workday', isWorkday: true },
  { date: '2025-02-08', name: '春节调休', type: 'workday', isWorkday: true },
  { date: '2025-04-27', name: '劳动节调休', type: 'workday', isWorkday: true },
  { date: '2025-09-28', name: '国庆节调休', type: 'workday', isWorkday: true },
  { date: '2025-10-11', name: '国庆节调休', type: 'workday', isWorkday: true },

  // 传统节日
  { date: '2025-02-14', name: '情人节', type: 'traditional', isWorkday: true },
  { date: '2025-03-08', name: '妇女节', type: 'traditional', isWorkday: true },
  { date: '2025-03-12', name: '植树节', type: 'traditional', isWorkday: true },
  { date: '2025-04-01', name: '愚人节', type: 'traditional', isWorkday: true },
  { date: '2025-05-04', name: '青年节', type: 'traditional', isWorkday: true },
  { date: '2025-05-11', name: '母亲节', type: 'traditional', isWorkday: true },
  { date: '2025-06-01', name: '儿童节', type: 'traditional', isWorkday: true },
  { date: '2025-06-15', name: '父亲节', type: 'traditional', isWorkday: true },
  { date: '2025-08-07', name: '七夕节', type: 'traditional', isWorkday: true },
  { date: '2025-09-10', name: '教师节', type: 'traditional', isWorkday: true },
  { date: '2025-10-31', name: '万圣节', type: 'traditional', isWorkday: true },
  { date: '2025-12-24', name: '平安夜', type: 'traditional', isWorkday: true },
  { date: '2025-12-25', name: '圣诞节', type: 'traditional', isWorkday: true }
]

// 节日类型配置
export const holidayTypes = {
  legal: {
    name: '法定节假日',
    color: '#f56c6c',
    backgroundColor: '#fef0f0',
    icon: '🎉',
    priority: 1
  },
  workday: {
    name: '调休工作日',
    color: '#e6a23c',
    backgroundColor: '#fdf6ec',
    icon: '💼',
    priority: 2
  },
  traditional: {
    name: '传统节日',
    color: '#909399',
    backgroundColor: '#f4f4f5',
    icon: '🏮',
    priority: 3
  }
}

/**
 * 获取指定日期的节日信息
 * @param {Date|string} date 日期
 * @returns {Object|null} 节日信息或null
 */
export function getHolidayInfo(date) {
  const dateString = formatDateString(date)
  return chineseHolidays2025.find(holiday => holiday.date === dateString) || null
}

/**
 * 获取指定月份的所有节日
 * @param {number} year 年份
 * @param {number} month 月份 (1-12)
 * @returns {Array} 节日列表
 */
export function getHolidaysForMonth(year, month) {
  const monthString = `${year}-${month.toString().padStart(2, '0')}`
  return chineseHolidays2025.filter(holiday => holiday.date.startsWith(monthString))
}

/**
 * 获取指定年份的所有节日
 * @param {number} year 年份
 * @returns {Array} 节日列表
 */
export function getHolidaysForYear(year) {
  const yearString = year.toString()
  return chineseHolidays2025.filter(holiday => holiday.date.startsWith(yearString))
}

/**
 * 检查指定日期是否为节假日
 * @param {Date|string} date 日期
 * @returns {boolean} 是否为节假日
 */
export function isHoliday(date) {
  const holiday = getHolidayInfo(date)
  return holiday && holiday.type === 'legal' && !holiday.isWorkday
}

/**
 * 检查指定日期是否为工作日
 * @param {Date|string} date 日期
 * @returns {boolean} 是否为工作日
 */
export function isWorkday(date) {
  const holiday = getHolidayInfo(date)
  if (holiday) {
    return holiday.isWorkday
  }
  
  // 如果没有特殊节日，按正常工作日判断（周一到周五）
  const dateObj = new Date(date)
  const dayOfWeek = dateObj.getDay()
  return dayOfWeek >= 1 && dayOfWeek <= 5
}

/**
 * 格式化日期为字符串
 * @param {Date|string} date 日期
 * @returns {string} YYYY-MM-DD格式字符串
 */
function formatDateString(date) {
  const dateObj = typeof date === 'string' ? new Date(date) : date
  return dateObj.toISOString().split('T')[0]
}

/**
 * 获取节日显示样式
 * @param {string} type 节日类型
 * @returns {Object} 样式对象
 */
export function getHolidayStyle(type) {
  return holidayTypes[type] || holidayTypes.traditional
}

/**
 * 获取即将到来的节日
 * @param {number} days 未来天数，默认30天
 * @returns {Array} 即将到来的节日列表
 */
export function getUpcomingHolidays(days = 30) {
  const today = new Date()
  const futureDate = new Date(today.getTime() + days * 24 * 60 * 60 * 1000)
  
  return chineseHolidays2025.filter(holiday => {
    const holidayDate = new Date(holiday.date)
    return holidayDate >= today && holidayDate <= futureDate
  }).sort((a, b) => new Date(a.date) - new Date(b.date))
}