<template>
  <div class="app-header">
    <!-- 左侧区域 -->
    <div class="header-left">
      <!-- 移动端菜单按钮 -->
      <el-button
        v-if="showMobileMenu"
        type="text"
        @click="toggleMobileMenu"
        class="mobile-menu-btn"
      >
        <el-icon><Menu /></el-icon>
      </el-button>
      
      <!-- 面包屑导航 -->
      <el-breadcrumb v-if="breadcrumbs.length > 0" separator="/" class="breadcrumb">
        <el-breadcrumb-item
          v-for="(item, index) in breadcrumbs"
          :key="index"
          :to="item.path"
        >
          <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
          {{ item.title }}
        </el-breadcrumb-item>
      </el-breadcrumb>
      
      <!-- 页面标题 -->
      <h1 v-if="title" class="page-title">{{ title }}</h1>
    </div>

    <!-- 中间区域 - 搜索 -->
    <div v-if="showSearch" class="header-center">
      <div class="search-container">
        <el-input
          v-model="searchQuery"
          :placeholder="searchPlaceholder"
          @keyup.enter="handleSearch"
          class="search-input"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <!-- 右侧区域 -->
    <div class="header-right">
      <!-- 自定义操作插槽 -->
      <slot name="actions" />
      
      <!-- 快捷操作按钮 -->
      <div class="quick-actions">
        <!-- 全屏切换 -->
        <el-tooltip :content="isFullscreen ? '退出全屏' : '全屏'" placement="bottom">
          <el-button type="text" @click="toggleFullscreen" class="action-btn">
            <el-icon>
              <FullScreen v-if="!isFullscreen" />
              <Aim v-else />
            </el-icon>
          </el-button>
        </el-tooltip>

        <!-- 主题切换 -->
        <el-tooltip :content="isDark ? '浅色模式' : '深色模式'" placement="bottom">
          <el-button type="text" @click="toggleTheme" class="action-btn">
            <el-icon>
              <Sunny v-if="isDark" />
              <Moon v-else />
            </el-icon>
          </el-button>
        </el-tooltip>

        <!-- 通知中心 -->
        <el-dropdown @command="handleNotificationCommand" placement="bottom-end">
          <el-button type="text" class="action-btn notification-btn">
            <el-badge :value="unreadCount" :hidden="unreadCount === 0">
              <el-icon><Bell /></el-icon>
            </el-badge>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu class="notification-dropdown">
              <div class="notification-header">
                <span>通知</span>
                <el-button
                  type="text"
                  size="small"
                  @click="handleNotificationCommand('mark-all-read')"
                >
                  全部已读
                </el-button>
              </div>
              <el-dropdown-item
                v-for="notification in notifications.slice(0, 5)"
                :key="notification.id"
                :command="{ action: 'view', id: notification.id }"
                class="notification-item"
              >
                <div class="notification-content">
                  <div class="notification-title">{{ notification.title }}</div>
                  <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
                </div>
                <el-badge v-if="!notification.read" is-dot />
              </el-dropdown-item>
              <el-dropdown-item
                v-if="notifications.length === 0"
                disabled
                class="empty-notification"
              >
                暂无通知
              </el-dropdown-item>
              <div class="notification-footer">
                <el-button
                  type="text"
                  size="small"
                  @click="handleNotificationCommand('view-all')"
                >
                  查看全部
                </el-button>
              </div>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>

      <!-- 用户菜单 -->
      <el-dropdown @command="handleUserCommand" placement="bottom-end">
        <div class="user-menu">
          <el-avatar :size="32" :src="user?.avatar" class="user-avatar">
            <el-icon><User /></el-icon>
          </el-avatar>
          <span v-if="!hideName" class="username">{{ user?.username || '用户' }}</span>
          <el-icon class="dropdown-arrow"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人资料
            </el-dropdown-item>
            <el-dropdown-item command="settings">
              <el-icon><Setting /></el-icon>
              设置
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  Menu,
  Search,
  FullScreen,
  Aim,
  Sunny,
  Moon,
  Bell,
  User,
  ArrowDown,
  Setting,
  SwitchButton
} from '@element-plus/icons-vue'

const props = defineProps({
  // 页面标题
  title: {
    type: String,
    default: ''
  },
  // 面包屑导航
  breadcrumbs: {
    type: Array,
    default: () => []
  },
  // 是否显示搜索
  showSearch: {
    type: Boolean,
    default: false
  },
  // 搜索占位符
  searchPlaceholder: {
    type: String,
    default: '搜索...'
  },
  // 是否显示移动端菜单按钮
  showMobileMenu: {
    type: Boolean,
    default: false
  },
  // 是否隐藏用户名
  hideName: {
    type: Boolean,
    default: false
  },
  // 用户信息
  user: {
    type: Object,
    default: () => null
  },
  // 通知列表
  notifications: {
    type: Array,
    default: () => []
  },
  // 未读通知数量
  unreadCount: {
    type: Number,
    default: 0
  },
  // 是否深色模式
  isDark: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'search',
  'toggle-mobile-menu',
  'toggle-fullscreen',
  'toggle-theme',
  'user-command',
  'notification-command'
])

const router = useRouter()
const searchQuery = ref('')
const isFullscreen = ref(false)

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    emit('search', searchQuery.value)
  }
}

const toggleMobileMenu = () => {
  emit('toggle-mobile-menu')
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
  emit('toggle-fullscreen', isFullscreen.value)
}

const toggleTheme = () => {
  emit('toggle-theme')
}

const handleUserCommand = (command) => {
  emit('user-command', command)
  
  // 内置路由处理
  switch (command) {
    case 'profile':
      router.push('/dashboard/profile')
      break
    case 'settings':
      router.push('/dashboard/settings')
      break
    case 'logout':
      // 由父组件处理登出逻辑
      break
  }
}

const handleNotificationCommand = (command) => {
  if (typeof command === 'string') {
    emit('notification-command', { action: command })
  } else {
    emit('notification-command', command)
  }
  
  // 内置路由处理
  if (command === 'view-all' || command.action === 'view-all') {
    router.push('/dashboard/notifications')
  }
}

const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  return `${Math.floor(diff / 86400000)}天前`
}
</script>

<style scoped>
.app-header {
  height: 64px;
  background: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 999;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  min-width: 0;
}

.mobile-menu-btn {
  color: #6b7280;
  padding: 8px;
}

.breadcrumb {
  color: #6b7280;
}

.breadcrumb :deep(.el-breadcrumb__item) {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.header-center {
  flex: 1;
  max-width: 400px;
  margin: 0 24px;
}

.search-container {
  width: 100%;
}

.search-input {
  width: 100%;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.quick-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-right: 16px;
}

.action-btn {
  color: #6b7280;
  padding: 8px;
  transition: all 0.2s;
}

.action-btn:hover {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
}

.notification-btn {
  position: relative;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s;
}

.user-menu:hover {
  background: #f3f4f6;
}

.user-avatar {
  cursor: pointer;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.dropdown-arrow {
  color: #9ca3af;
  font-size: 12px;
}

/* 通知下拉菜单样式 */
.notification-dropdown {
  width: 320px;
  max-height: 400px;
  overflow-y: auto;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #e5e7eb;
  font-weight: 600;
}

.notification-item {
  padding: 12px 16px !important;
  border-bottom: 1px solid #f3f4f6;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notification-time {
  font-size: 12px;
  color: #9ca3af;
}

.empty-notification {
  text-align: center;
  color: #9ca3af;
  padding: 24px 16px !important;
}

.notification-footer {
  padding: 8px 16px;
  border-top: 1px solid #e5e7eb;
  text-align: center;
}

/* 深色模式 */
.dark .app-header {
  background: #1f2937;
  border-bottom-color: #374151;
}

.dark .page-title {
  color: #f9fafb;
}

.dark .username {
  color: #e5e7eb;
}

.dark .user-menu:hover {
  background: #374151;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-header {
    padding: 0 16px;
  }
  
  .header-center {
    display: none;
  }
  
  .username {
    display: none;
  }
  
  .quick-actions {
    margin-right: 8px;
  }
}

@media (max-width: 480px) {
  .breadcrumb {
    display: none;
  }
  
  .page-title {
    font-size: 16px;
  }
}
</style>