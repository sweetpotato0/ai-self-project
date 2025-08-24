<template>
  <div class="settings-page">
    <h2>系统设置</h2>

    <el-tabs v-model="activeTab" class="settings-tabs">
      <!-- 个人设置 -->
      <el-tab-pane label="个人设置" name="profile">
        <div class="settings-section">
          <h3>个人信息</h3>
          <el-form :model="profileForm" :rules="profileRules" ref="profileFormRef" label-width="100px">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="profileForm.username" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="profileForm.email" />
            </el-form-item>
            <el-form-item label="昵称">
              <el-input v-model="profileForm.nickname" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="updateProfile">保存更改</el-button>
            </el-form-item>
          </el-form>
        </div>

        <div class="settings-section">
          <h3>密码设置</h3>
          <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
            <el-form-item label="当前密码" prop="currentPassword">
              <el-input v-model="passwordForm.currentPassword" type="password" />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="passwordForm.newPassword" type="password" />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="passwordForm.confirmPassword" type="password" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="updatePassword">修改密码</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <!-- 通知设置 -->
      <el-tab-pane label="通知设置" name="notifications">
        <div class="settings-section">
          <h3>通知偏好</h3>
          <el-form label-width="200px">
            <el-form-item label="任务到期提醒">
              <el-switch v-model="notificationSettings.dueReminder" />
            </el-form-item>
            <el-form-item label="任务完成通知">
              <el-switch v-model="notificationSettings.completionNotification" />
            </el-form-item>
            <el-form-item label="新任务通知">
              <el-switch v-model="notificationSettings.newTaskNotification" />
            </el-form-item>
            <el-form-item label="邮件通知">
              <el-switch v-model="notificationSettings.emailNotification" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveNotificationSettings">保存设置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <!-- 界面设置 -->
      <el-tab-pane label="界面设置" name="interface">
        <div class="settings-section">
          <h3>主题设置</h3>
          <el-form label-width="100px">
            <el-form-item label="主题模式">
              <el-radio-group v-model="interfaceSettings.theme">
                <el-radio label="light">浅色模式</el-radio>
                <el-radio label="dark">深色模式</el-radio>
                <el-radio label="auto">跟随系统</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="语言">
              <el-select v-model="interfaceSettings.language">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en-US" />
              </el-select>
            </el-form-item>
            <el-form-item label="时区">
              <el-select v-model="interfaceSettings.timezone">
                <el-option label="北京时间 (UTC+8)" value="Asia/Shanghai" />
                <el-option label="UTC" value="UTC" />
                <el-option label="美国东部时间" value="America/New_York" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveInterfaceSettings">保存设置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <!-- 数据管理 -->
      <el-tab-pane label="数据管理" name="data">
        <div class="settings-section">
          <h3>数据导出</h3>
          <p>导出您的所有任务数据为JSON格式</p>
          <el-button type="primary" @click="exportData">导出数据</el-button>
        </div>

        <div class="settings-section">
          <h3>数据导入</h3>
          <p>从JSON文件导入任务数据</p>
          <el-upload
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            accept=".json"
          >
            <el-button type="primary">选择文件</el-button>
          </el-upload>
        </div>

        <div class="settings-section">
          <h3>数据清理</h3>
          <p>删除所有已完成的任务（此操作不可恢复）</p>
          <el-button type="danger" @click="clearCompletedTasks">清理已完成任务</el-button>
        </div>
      </el-tab-pane>

      <!-- 关于 -->
      <el-tab-pane label="关于" name="about">
        <div class="settings-section">
          <h3>系统信息</h3>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="应用名称">TaskMaster</el-descriptions-item>
            <el-descriptions-item label="版本">1.0.0</el-descriptions-item>
            <el-descriptions-item label="构建时间">2024-08-23</el-descriptions-item>
            <el-descriptions-item label="技术栈">Vue 3 + Element Plus + Gin</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useTodoStore } from '@/stores/todo'
import { useSettingsStore } from '@/stores/settings'
import { ElMessage, ElMessageBox } from 'element-plus'

const authStore = useAuthStore()
const todoStore = useTodoStore()
const settingsStore = useSettingsStore()

// 响应式数据
const activeTab = ref('profile')
const profileFormRef = ref()
const passwordFormRef = ref()

const profileForm = ref({
  username: '',
  email: '',
  nickname: ''
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const notificationSettings = ref({
  dueReminder: true,
  completionNotification: true,
  newTaskNotification: true,
  emailNotification: false
})

const interfaceSettings = ref({
  theme: 'light',
  language: 'zh-CN',
  timezone: 'Asia/Shanghai'
})

// 表单验证规则
const profileRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 方法
const updateProfile = async () => {
  try {
    await profileFormRef.value.validate()
    const success = await settingsStore.updateUserProfile(profileForm.value)
    if (success) {
      // 更新本地用户信息
      authStore.user.username = profileForm.value.username
      authStore.user.email = profileForm.value.email
      authStore.user.nickname = profileForm.value.nickname
    }
  } catch (error) {
    console.error('更新失败:', error)
  }
}

const updatePassword = async () => {
  try {
    await passwordFormRef.value.validate()
    const success = await settingsStore.changeUserPassword({
      current_password: passwordForm.value.currentPassword,
      new_password: passwordForm.value.newPassword
    })
    if (success) {
      passwordForm.value = {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      }
    }
  } catch (error) {
    console.error('密码修改失败:', error)
  }
}

const saveNotificationSettings = async () => {
  await settingsStore.saveNotificationSettings({
    due_reminder: notificationSettings.value.dueReminder,
    completion_notification: notificationSettings.value.completionNotification,
    new_task_notification: notificationSettings.value.newTaskNotification,
    email_notification: notificationSettings.value.emailNotification
  })
}

const saveInterfaceSettings = async () => {
  await settingsStore.saveInterfaceSettings({
    theme: interfaceSettings.value.theme,
    language: interfaceSettings.value.language,
    timezone: interfaceSettings.value.timezone
  })
}

const exportData = async () => {
  await settingsStore.exportUserData()
}

const handleFileChange = async (file) => {
  const reader = new FileReader()
  reader.onload = async (e) => {
    try {
      const data = JSON.parse(e.target.result)
      
      // 验证数据格式
      if (!data.todos || !Array.isArray(data.todos)) {
        throw new Error('无效的数据格式')
      }
      
      // 这里可以添加更多的数据导入逻辑
      // 例如：批量创建待办事项
      ElMessage.success('数据导入成功')
    } catch (error) {
      console.error('数据导入失败:', error)
      ElMessage.error('文件格式错误或数据无效')
    }
  }
  reader.readAsText(file.raw)
}

const clearCompletedTasks = async () => {
  try {
    await ElMessageBox.confirm(
      '此操作将删除所有已完成的任务，且不可恢复。确定继续吗？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const success = await settingsStore.clearCompleted()
    if (success) {
      // 刷新待办事项列表
      await todoStore.fetchTodos()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清理失败:', error)
    }
  }
}

// 生命周期
onMounted(async () => {
  // 初始化表单数据
  if (authStore.user) {
    profileForm.value.username = authStore.user.username || ''
    profileForm.value.email = authStore.user.email || ''
    profileForm.value.nickname = authStore.user.nickname || ''
  }
  
  // 加载设置
  await settingsStore.fetchSettings()
  
  // 同步设置到表单
  const settings = settingsStore.settings
  notificationSettings.value = {
    dueReminder: settings.due_reminder,
    completionNotification: settings.completion_notification,
    newTaskNotification: settings.new_task_notification,
    emailNotification: settings.email_notification
  }
  
  interfaceSettings.value = {
    theme: settings.theme,
    language: settings.language,
    timezone: settings.timezone
  }
})
</script>

<style scoped>
.settings-page {
  padding: 20px;
}

.settings-tabs {
  margin-top: 20px;
}

.settings-section {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fff;
}

.settings-section h3 {
  margin: 0 0 20px 0;
  color: #303133;
  font-size: 16px;
  font-weight: 600;
}

.settings-section p {
  color: #606266;
  margin-bottom: 15px;
}
</style>
