// 格式化日期
export const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 获取状态文本
export const getStatusText = (status) => {
  const statusMap = {
    pending: '待处理',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

// 获取优先级文本
export const getPriorityText = (level) => {
  const priorityMap = {
    1: '低',
    2: '中',
    3: '高',
    4: '紧急',
    5: '立即'
  }
  return priorityMap[level] || level
}
