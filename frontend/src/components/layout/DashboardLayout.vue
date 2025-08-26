<template>
  <AppLayout
    :app-title="appConfig.title"
    :app-version="appConfig.version"
    :page-title="pageConfig.title"
    :page-description="pageConfig.description"
    :breadcrumbs="breadcrumbs"
    :show-search="pageConfig.showSearch"
    :search-placeholder="pageConfig.searchPlaceholder"
    :menu-items="filteredMenuItems"
    :active-menu="activeMenu"
    :user="user"
    :notifications="notifications"
    :unread-count="unreadCount"
    :is-dark-mode="isDarkMode"
    :show-language-switch="true"
    :languages="supportedLanguages"
    :current-language="currentLanguage"
    :show-status="true"
    :system-status="systemStatus"
    :footer-links="footerLinks"
    :copyright="copyright"
    @menu-click="handleMenuClick"
    @user-click="handleUserClick"
    @user-command="handleUserCommand"
    @search="handleSearch"
    @toggle-fullscreen="handleToggleFullscreen"
    @toggle-theme="handleToggleTheme"
    @notification-command="handleNotificationCommand"
    @language-change="handleLanguageChange"
    @footer-link-click="handleFooterLinkClick"
  >
    <!-- 页面操作按钮插槽 -->
    <template #page-actions>
      <slot name="page-actions" />
    </template>

    <!-- 头部操作按钮插槽 -->
    <template #header-actions>
      <slot name="header-actions" />
    </template>

    <!-- 主要内容插槽 -->
    <template #default>
      <slot />
    </template>
  </AppLayout>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from './AppLayout.vue'
import {
  defaultMenuItems,
  filterMenuByPermissions,
  generateBreadcrumbs,
  getActiveMenu,
  getPageConfig,
  supportedLanguages,
  defaultFooterLinks
} from './layoutConfig'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { useSettingsStore } from '@/stores/settings'

// Store
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const settingsStore = useSettingsStore()

const route = useRoute()
const router = useRouter()

// 应用配置
const appConfig = {
  title: 'TaskMaster',
  version: '1.0.0'
}

// 计算属性
const user = computed(() => authStore.user)
const notifications = computed(() => notificationStore.notifications.slice(0, 10))
const unreadCount = computed(() => notificationStore.unreadCount)
const isDarkMode = computed(() => settingsStore.settings?.theme === 'dark')
const currentLanguage = computed(() => 
  supportedLanguages.find(lang => lang.code === settingsStore.settings?.language) || 
  supportedLanguages[0]
)
const systemStatus = computed(() => {
  // 可以根据实际情况返回系统状态
  return 'online'
})

// 过滤后的菜单项
const filteredMenuItems = computed(() => {
  const permissions = user.value?.permissions || []
  return filterMenuByPermissions(defaultMenuItems, user.value, permissions)
})

// 当前路由的页面配置
const pageConfig = computed(() => getPageConfig(route.path))

// 面包屑导航
const breadcrumbs = computed(() => generateBreadcrumbs(route.path, filteredMenuItems.value))

// 当前激活菜单
const activeMenu = computed(() => getActiveMenu(route.path, filteredMenuItems.value))

// 底部配置
const footerLinks = computed(() => defaultFooterLinks)
const copyright = computed(() => `© ${new Date().getFullYear()} ${appConfig.title}. All rights reserved.`)

// 事件处理
const handleMenuClick = (item) => {
  if (item.route && item.route !== route.path) {
    router.push(item.route)
  }
}

const handleUserClick = () => {
  // 用户头像点击事件
  router.push('/dashboard/profile')
}

const handleUserCommand = async (command) => {
  switch (command) {
    case 'profile':
      router.push('/dashboard/profile')
      break
    case 'settings':
      router.push('/dashboard/settings')
      break
    case 'logout':
      try {
        await authStore.logout()
        router.push('/login')
      } catch (error) {
        console.error('Logout failed:', error)
      }
      break
  }
}

const handleSearch = (query) => {
  // 全局搜索功能
  const currentPath = route.path
  
  if (currentPath.includes('/todos')) {
    // 在待办事项页面搜索
    // 这里可以调用相应的搜索方法
    console.log('Searching todos:', query)
  } else if (currentPath.includes('/articles')) {
    // 在文章页面搜索
    console.log('Searching articles:', query)
  } else if (currentPath.includes('/tools')) {
    // 在工具页面搜索
    console.log('Searching tools:', query)
  } else {
    // 全局搜索
    router.push({ path: '/dashboard/search', query: { q: query } })
  }
}

const handleToggleFullscreen = (isFullscreen) => {
  // 全屏切换处理
  console.log('Fullscreen toggled:', isFullscreen)
}

const handleToggleTheme = () => {
  // 主题切换
  const newTheme = isDarkMode.value ? 'light' : 'dark'
  settingsStore.updateInterfaceSettings({ theme: newTheme })
}

const handleNotificationCommand = async (command) => {
  switch (command.action) {
    case 'view':
      // 查看特定通知
      await notificationStore.markAsRead(command.id)
      break
    case 'mark-all-read':
      // 标记所有通知为已读
      await notificationStore.markAllAsRead()
      break
    case 'view-all':
      // 查看所有通知
      router.push('/dashboard/notifications')
      break
  }
}

const handleLanguageChange = (language) => {
  // 语言切换
  settingsStore.updateInterfaceSettings({ language: language.code })
}

const handleFooterLinkClick = (link) => {
  // 底部链接点击处理
  if (link.external) {
    window.open(link.url, '_blank')
  } else {
    router.push(link.url)
  }
}

// 监听路由变化，更新页面标题
watch(
  () => pageConfig.value.title,
  (newTitle) => {
    if (newTitle) {
      document.title = `${newTitle} - ${appConfig.title}`
    } else {
      document.title = appConfig.title
    }
  },
  { immediate: true }
)
</script>