<template>
  <div class="dashboard-container">
    <!-- 侧边栏 -->
    <div class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <div class="logo">
          <el-icon class="logo-icon"><Grid /></el-icon>
          <span v-if="!sidebarCollapsed" class="logo-text">TaskMaster</span>
        </div>
        <el-button
          v-if="!sidebarCollapsed"
          type="text"
          @click="toggleSidebar"
          class="collapse-btn"
        >
          <el-icon><Fold /></el-icon>
        </el-button>
        <el-button
          v-else
          type="text"
          @click="toggleSidebar"
          class="collapse-btn collapsed-toggle"
        >
          <el-icon><Expand /></el-icon>
        </el-button>
      </div>

      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        :collapse="sidebarCollapsed"
        background-color="transparent"
      >
        <template v-if="sidebarCollapsed">
          <el-tooltip content="仪表盘" placement="right">
            <el-menu-item index="dashboard" @click="navigateTo('dashboard')">
              <el-icon><DataBoard /></el-icon>
            </el-menu-item>
          </el-tooltip>

          <el-tooltip content="TODO清单" placement="right">
            <el-menu-item index="todos" @click="navigateTo('todos')">
              <el-icon><List /></el-icon>
            </el-menu-item>
          </el-tooltip>

          <el-tooltip content="日程安排" placement="right">
            <el-menu-item index="calendar" @click="navigateTo('calendar')">
              <el-icon><Calendar /></el-icon>
            </el-menu-item>
          </el-tooltip>

          <el-tooltip content="数据分析" placement="right">
            <el-menu-item index="analytics" @click="navigateTo('analytics')">
              <el-icon><TrendCharts /></el-icon>
            </el-menu-item>
          </el-tooltip>

          <el-tooltip content="个人文章" placement="right">
            <el-menu-item index="articles" @click="navigateTo('articles')">
              <el-icon><Document /></el-icon>
            </el-menu-item>
          </el-tooltip>

          <el-tooltip content="系统设置" placement="right">
            <el-menu-item index="settings" @click="navigateTo('settings')">
              <el-icon><Setting /></el-icon>
            </el-menu-item>
          </el-tooltip>
        </template>

        <template v-else>
          <el-menu-item index="dashboard" @click="navigateTo('dashboard')">
            <el-icon><DataBoard /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>

          <el-menu-item index="todos" @click="navigateTo('todos')">
            <el-icon><List /></el-icon>
            <span>TODO清单</span>
          </el-menu-item>

          <el-menu-item index="calendar" @click="navigateTo('calendar')">
            <el-icon><Calendar /></el-icon>
            <span>日程安排</span>
          </el-menu-item>

          <el-menu-item index="analytics" @click="navigateTo('analytics')">
            <el-icon><TrendCharts /></el-icon>
            <span>数据分析</span>
          </el-menu-item>

          <el-menu-item index="articles" @click="navigateTo('articles')">
            <el-icon><Document /></el-icon>
            <span>个人文章</span>
          </el-menu-item>

          <el-menu-item index="settings" @click="navigateTo('settings')">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </template>
      </el-menu>

      <div class="sidebar-footer">
        <div class="user-info">
          <el-tooltip 
            v-if="sidebarCollapsed" 
            :content="authStore.user?.username || '用户'" 
            placement="right"
          >
            <el-avatar :size="32" :src="userAvatar">
              {{ authStore.user?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
          </el-tooltip>
          <el-avatar 
            v-else 
            :size="32" 
            :src="userAvatar"
          >
            {{ authStore.user?.username?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <div v-if="!sidebarCollapsed" class="user-details">
            <div class="username">{{ authStore.user?.username }}</div>
            <div class="user-role">管理员</div>
          </div>
        </div>
        <el-tooltip 
          v-if="sidebarCollapsed" 
          content="退出登录" 
          placement="right"
        >
          <el-button type="text" @click="logout" class="logout-btn">
            <el-icon><SwitchButton /></el-icon>
          </el-button>
        </el-tooltip>
        <el-button 
          v-else 
          type="text" 
          @click="logout" 
          class="logout-btn"
        >
          <el-icon><SwitchButton /></el-icon>
        </el-button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 顶部导航栏 -->
      <div class="top-navbar">
        <div class="navbar-left">
          <h2 class="page-title">{{ pageTitle }}</h2>
        </div>
        <div class="navbar-right">
          <NotificationPanel />
          <WebSocketNotification />
          <el-dropdown>
            <el-avatar :size="40" :src="userAvatar">
              {{ authStore.user?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="navigateTo('profile')">
                  <el-icon><User /></el-icon>个人资料
                </el-dropdown-item>
                <el-dropdown-item @click="navigateTo('settings')">
                  <el-icon><Setting /></el-icon>设置
                </el-dropdown-item>
                <el-dropdown-item divided @click="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 页面内容 -->
      <div class="page-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import {
  Grid, Fold, Expand, DataBoard, List, Calendar,
  TrendCharts, Setting, User, SwitchButton, Document
} from '@element-plus/icons-vue'
import NotificationPanel from '@/components/NotificationPanel.vue'
import WebSocketNotification from '@/components/WebSocketNotification.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const sidebarCollapsed = ref(false)
const userAvatar = ref('')

// 计算属性
const activeMenu = computed(() => route.name)
const pageTitle = computed(() => {
  const titleMap = {
    dashboard: '仪表盘',
    todos: 'TODO清单',
    calendar: '日程安排',
    analytics: '数据分析',
    articles: '个人文章',
    settings: '系统设置',
    profile: '个人资料'
  }
  return titleMap[route.name] || '仪表盘'
})

// 方法
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const navigateTo = (routeName) => {
  router.push({ name: routeName })
}

const logout = async () => {
  try {
    authStore.logout()
    router.push('/login')
    ElMessage.success('已退出登录')
  } catch (error) {
    ElMessage.error('退出登录失败')
  }
}

// 生命周期
onMounted(() => {
  // 初始化用户头像
  if (authStore.user?.username) {
    userAvatar.value = `https://api.dicebear.com/7.x/avataaars/svg?seed=${authStore.user.username}`
  }
})
</script>

<style scoped>
.dashboard-container {
  display: flex;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.sidebar {
  width: 260px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  font-size: 24px;
  color: #409EFF;
}

.logo-text {
  font-size: 18px;
  font-weight: bold;
}

.collapse-btn {
  font-size: 16px;
}

.sidebar-menu {
  flex: 1;
  border: none;
  margin-top: 20px;
}

.sidebar-menu .el-menu-item {
  margin: 4px 12px;
  border-radius: 8px;
  height: 48px;
  line-height: 48px;
}

.sidebar-menu .el-menu-item:hover {
  background: rgba(64, 158, 255, 0.1);
}

.sidebar-menu .el-menu-item.is-active {
  background: rgba(64, 158, 255, 0.2);
  color: #409EFF;
}

.sidebar-footer {
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-details {
  display: flex;
  flex-direction: column;
}

.username {
  font-size: 14px;
  font-weight: 500;
}

.user-role {
  font-size: 12px;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px 0 0 20px;
  margin: 20px 20px 20px 0;
  overflow: hidden;
}

.top-navbar {
  height: 80px;
  padding: 0 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.9);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.page-content {
  flex: 1;
  padding: 30px;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    z-index: 1000;
    transform: translateX(-100%);
  }

  .sidebar.collapsed {
    transform: translateX(0);
  }

  .main-content {
    margin: 0;
    border-radius: 0;
  }
}
</style>
