<template>
  <div class="profile-container">
    <div class="profile-header">
      <div class="profile-avatar-section">
        <div class="avatar-wrapper">
          <el-avatar :size="120" :src="userAvatar" class="profile-avatar">
            {{ authStore.user?.username?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <div class="avatar-overlay">
            <el-button type="primary" size="small" @click="changeAvatar">
              <el-icon><Camera /></el-icon>
              更换头像
            </el-button>
          </div>
        </div>
        <h2 class="profile-name">{{ authStore.user?.username }}</h2>
        <p class="profile-email">{{ authStore.user?.email }}</p>
        <div class="profile-stats">
          <div class="stat-item">
            <div class="stat-number">{{ stats.totalTasks }}</div>
            <div class="stat-label">总任务</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ stats.completedTasks }}</div>
            <div class="stat-label">已完成</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ stats.totalHours }}</div>
            <div class="stat-label">总工时</div>
          </div>
        </div>
      </div>
    </div>

    <div class="profile-content">
      <el-tabs v-model="activeTab" class="profile-tabs">
        <el-tab-pane label="基本信息" name="basic">
          <div class="tab-content">
            <el-form
              ref="profileFormRef"
              :model="profileForm"
              :rules="profileRules"
              label-width="100px"
              class="profile-form"
            >
              <el-form-item label="用户名" prop="username">
                <el-input v-model="profileForm.username" placeholder="请输入用户名" />
              </el-form-item>

              <el-form-item label="邮箱" prop="email">
                <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
              </el-form-item>

              <el-form-item label="昵称" prop="nickname">
                <el-input v-model="profileForm.nickname" placeholder="请输入昵称" />
              </el-form-item>

              <el-form-item label="个人简介" prop="bio">
                <el-input
                  v-model="profileForm.bio"
                  type="textarea"
                  :rows="4"
                  placeholder="请输入个人简介"
                />
              </el-form-item>

              <el-form-item>
                <el-button type="primary" @click="saveProfile" :loading="saving">
                  保存修改
                </el-button>
                <el-button @click="resetForm">重置</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="password">
          <div class="tab-content">
            <el-form
              ref="passwordFormRef"
              :model="passwordForm"
              :rules="passwordRules"
              label-width="100px"
              class="password-form"
            >
              <el-form-item label="当前密码" prop="currentPassword">
                <el-input
                  v-model="passwordForm.currentPassword"
                  type="password"
                  placeholder="请输入当前密码"
                  show-password
                />
              </el-form-item>

              <el-form-item label="新密码" prop="newPassword">
                <el-input
                  v-model="passwordForm.newPassword"
                  type="password"
                  placeholder="请输入新密码"
                  show-password
                />
              </el-form-item>

              <el-form-item label="确认密码" prop="confirmPassword">
                <el-input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  placeholder="请再次输入新密码"
                  show-password
                />
              </el-form-item>

              <el-form-item>
                <el-button type="primary" @click="changePassword" :loading="changingPassword">
                  修改密码
                </el-button>
                <el-button @click="resetPasswordForm">重置</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useTodoStore } from '@/stores/todo'
import { ElMessage } from 'element-plus'
import { Camera } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const todoStore = useTodoStore()

const activeTab = ref('basic')
const userAvatar = ref('')
const saving = ref(false)
const changingPassword = ref(false)

// 表单引用
const profileFormRef = ref()
const passwordFormRef = ref()

// 个人资料表单
const profileForm = reactive({
  username: '',
  email: '',
  nickname: '',
  bio: ''
})

// 密码表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 统计数据
const stats = reactive({
  totalTasks: 0,
  completedTasks: 0,
  totalHours: 0
})

// 表单验证规则
const profileRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
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
const changeAvatar = () => {
  ElMessage.info('头像更换功能开发中...')
}

const saveProfile = async () => {
  try {
    await profileFormRef.value.validate()
    saving.value = true

    // 这里调用API保存个人资料
    await new Promise(resolve => setTimeout(resolve, 1000))

    ElMessage.success('个人资料保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const resetForm = () => {
  profileFormRef.value.resetFields()
  loadProfileData()
}

const changePassword = async () => {
  try {
    await passwordFormRef.value.validate()
    changingPassword.value = true

    // 这里调用API修改密码
    await new Promise(resolve => setTimeout(resolve, 1000))

    ElMessage.success('密码修改成功')
    resetPasswordForm()
  } catch (error) {
    ElMessage.error('密码修改失败')
  } finally {
    changingPassword.value = false
  }
}

const resetPasswordForm = () => {
  passwordFormRef.value.resetFields()
}

const loadProfileData = () => {
  if (authStore.user) {
    profileForm.username = authStore.user.username || ''
    profileForm.email = authStore.user.email || ''
    profileForm.nickname = authStore.user.nickname || ''
    profileForm.bio = authStore.user.bio || ''
  }
}

const loadStats = () => {
  const todos = todoStore.todos || []
  stats.totalTasks = todos.length
  stats.completedTasks = todos.filter(todo => todo.status === 'completed').length
  stats.totalHours = todos.reduce((sum, todo) => sum + (todo.actual_hours || 0), 0)
}

// 生命周期
onMounted(() => {
  // 初始化用户头像
  if (authStore.user?.username) {
    userAvatar.value = `https://api.dicebear.com/7.x/avataaars/svg?seed=${authStore.user.username}`
  }

  loadProfileData()
  loadStats()
})
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.profile-header {
  text-align: center;
  margin-bottom: 40px;
  padding: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  color: white;
  position: relative;
  overflow: hidden;
}

.profile-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="white" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="white" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="white" opacity="0.1"/><circle cx="10" cy="60" r="0.5" fill="white" opacity="0.1"/><circle cx="90" cy="40" r="0.5" fill="white" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
  opacity: 0.3;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 20px;
}

.profile-avatar {
  border: 4px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.avatar-overlay {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  opacity: 0;
  transition: opacity 0.3s;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

.profile-name {
  font-size: 28px;
  font-weight: 600;
  margin: 10px 0 5px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.profile-email {
  font-size: 16px;
  opacity: 0.9;
  margin-bottom: 30px;
}

.profile-stats {
  display: flex;
  justify-content: center;
  gap: 40px;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.8;
}

.profile-content {
  background: white;
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.profile-tabs {
  padding: 20px;
}

.tab-content {
  padding: 20px 0;
}

.profile-form,
.password-form {
  max-width: 500px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .profile-container {
    padding: 10px;
  }

  .profile-header {
    padding: 20px;
  }

  .profile-stats {
    gap: 20px;
  }

    .stat-number {
    font-size: 24px;
  }
}
</style>
