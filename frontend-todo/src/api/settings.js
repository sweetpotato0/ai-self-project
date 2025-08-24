import { request } from './index'

// 获取用户设置
export function getSettings() {
  return request.get('/settings')
}

// 更新个人资料
export function updateProfile(data) {
  return request.put('/settings/profile', data)
}

// 修改密码
export function changePassword(data) {
  return request.put('/settings/password', data)
}

// 更新通知设置
export function updateNotificationSettings(data) {
  return request.put('/settings/notifications', data)
}

// 更新界面设置
export function updateInterfaceSettings(data) {
  return request.put('/settings/interface', data)
}

// 导出数据
export function exportData() {
  return request.get('/settings/export', {
    responseType: 'blob'
  })
}

// 清理已完成任务
export function clearCompletedTasks() {
  return request.delete('/settings/completed-tasks')
}