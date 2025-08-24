import { ElMessage, ElNotification } from 'element-plus'

class NotificationManager {
  constructor() {
    this.settings = {
      due_reminder: true,
      completion_notification: true,
      new_task_notification: true,
      email_notification: false
    }
  }

  // 更新通知设置
  updateSettings(newSettings) {
    this.settings = { ...this.settings, ...newSettings }
  }

  // 任务到期提醒
  showDueReminder(task) {
    if (!this.settings.due_reminder) return
    
    ElNotification({
      title: '任务即将到期',
      message: `任务「${task.title}」将在今天到期，请及时处理`,
      type: 'warning',
      duration: 5000,
      position: 'top-right'
    })
  }

  // 任务完成通知
  showCompletionNotification(task) {
    if (!this.settings.completion_notification) return
    
    ElNotification({
      title: '任务已完成',
      message: `恭喜！您已完成任务「${task.title}」`,
      type: 'success',
      duration: 3000,
      position: 'top-right'
    })
  }

  // 新任务创建通知
  showNewTaskNotification(task) {
    if (!this.settings.new_task_notification) return
    
    ElMessage({
      message: `新任务「${task.title}」已创建`,
      type: 'success',
      duration: 2000
    })
  }

  // 邮件通知（这里只是占位，实际需要后端支持）
  sendEmailNotification(type, task) {
    if (!this.settings.email_notification) return
    
    console.log(`发送邮件通知: ${type}`, task)
    // 实际实现需要调用后端邮件服务
  }

  // 检查任务到期情况
  checkTasksDue(tasks) {
    if (!this.settings.due_reminder) return
    
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    
    const tomorrow = new Date(today)
    tomorrow.setDate(today.getDate() + 1)
    
    tasks.forEach(task => {
      if (task.completed) return
      
      const dueDate = new Date(task.due_date)
      dueDate.setHours(0, 0, 0, 0)
      
      // 检查今天到期的任务
      if (dueDate.getTime() === today.getTime()) {
        this.showDueReminder(task)
      }
    })
  }

  // 显示桌面通知（需要用户授权）
  async showDesktopNotification(title, message, options = {}) {
    if (!this.settings.due_reminder && !this.settings.completion_notification) return
    
    // 检查浏览器是否支持通知
    if (!('Notification' in window)) {
      console.log('此浏览器不支持桌面通知')
      return
    }
    
    // 请求通知权限
    if (Notification.permission === 'default') {
      const permission = await Notification.requestPermission()
      if (permission !== 'granted') {
        console.log('用户拒绝了通知权限')
        return
      }
    }
    
    if (Notification.permission === 'granted') {
      new Notification(title, {
        body: message,
        icon: '/favicon.ico',
        badge: '/favicon.ico',
        ...options
      })
    }
  }

  // 批量处理通知设置更新
  applyNotificationSettings(todos) {
    if (!todos || todos.length === 0) return
    
    // 检查到期任务
    this.checkTasksDue(todos)
    
    // 可以添加更多批量处理逻辑
  }
}

// 创建全局实例
export const notificationManager = new NotificationManager()

// 导出类以便在需要时创建新实例
export default NotificationManager