import api from './index'

// 获取用户设置
export function getSettings() {
  return api.get('/settings')
}

// 更新个人资料
export function updateProfile(data) {
  return api.put('/settings/profile', data)
}

// 修改密码
export function changePassword(data) {
  return api.put('/settings/password', data)
}

// 更新通知设置
export function updateNotificationSettings(data) {
  return api.put('/settings/notifications', data)
}

// 更新界面设置
export function updateInterfaceSettings(data) {
  return api.put('/settings/interface', data)
}

// 导出数据
export function exportData() {
  return api.get('/settings/export', {
    responseType: 'blob'
  })
}

// 清理已完成任务
export function clearCompletedTasks() {
  return api.delete('/settings/completed-tasks')
}