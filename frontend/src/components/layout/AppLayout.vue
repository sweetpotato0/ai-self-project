<template>
  <div class="app-layout" :class="layoutClasses">
    <!-- 侧边栏 -->
    <AppSidebar
      v-if="showSidebar"
      :collapsed="sidebarCollapsed"
      :title="appTitle"
      :active-menu="activeMenu"
      :menu-items="menuItems"
      :user="user"
      @toggle="toggleSidebar"
      @menu-click="handleMenuClick"
      @user-click="handleUserClick"
    />
    
    <!-- 移动端遮罩 -->
    <div
      v-if="showSidebar && mobileMenuOpen"
      class="mobile-overlay"
      @click="closeMobileMenu"
    />

    <!-- 主内容区域 -->
    <div class="app-main" :class="{ 'sidebar-collapsed': sidebarCollapsed, 'no-sidebar': !showSidebar }">
      <!-- 顶部导航 -->
      <AppHeader
        v-if="showHeader"
        :title="pageTitle"
        :breadcrumbs="breadcrumbs"
        :show-search="showSearch"
        :search-placeholder="searchPlaceholder"
        :show-mobile-menu="showSidebar && isMobile"
        :hide-name="hideUserName"
        :user="user"
        :notifications="notifications"
        :unread-count="unreadCount"
        :is-dark="isDarkMode"
        @search="handleSearch"
        @toggle-mobile-menu="toggleMobileMenu"
        @toggle-fullscreen="handleToggleFullscreen"
        @toggle-theme="handleToggleTheme"
        @user-command="handleUserCommand"
        @notification-command="handleNotificationCommand"
      >
        <template #actions>
          <slot name="header-actions" />
        </template>
      </AppHeader>

      <!-- 页面内容 -->
      <main class="app-content" :class="{ 'with-footer': showFooter }">
        <div v-if="showBreadcrumb && breadcrumbs.length > 0" class="content-breadcrumb">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item
              v-for="(item, index) in breadcrumbs"
              :key="index"
              :to="item.path"
            >
              <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <!-- 页面标题区域 -->
        <div v-if="showPageHeader" class="page-header">
          <div class="page-header-left">
            <h1 v-if="pageTitle" class="page-title">{{ pageTitle }}</h1>
            <p v-if="pageDescription" class="page-description">{{ pageDescription }}</p>
          </div>
          <div class="page-header-right">
            <slot name="page-actions" />
          </div>
        </div>

        <!-- 内容插槽 -->
        <div class="page-content">
          <slot />
        </div>
      </main>

      <!-- 底部 -->
      <AppFooter
        v-if="showFooter"
        :minimal="minimalFooter"
        :fixed="fixedFooter"
        :app-name="appTitle"
        :version="appVersion"
        :copyright="copyright"
        :links="footerLinks"
        :social-links="socialLinks"
        :show-language="showLanguageSwitch"
        :languages="languages"
        :current-language="currentLanguage"
        :show-status="showStatus"
        :status="systemStatus"
        :last-updated="lastUpdated"
        :server-info="serverInfo"
        :beian="beianInfo"
        @language-change="handleLanguageChange"
        @link-click="handleFooterLinkClick"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import AppFooter from './AppFooter.vue'

const props = defineProps({
  // 应用配置
  appTitle: {
    type: String,
    default: 'TaskMaster'
  },
  appVersion: {
    type: String,
    default: '1.0.0'
  },
  
  // 布局配置
  showSidebar: {
    type: Boolean,
    default: true
  },
  showHeader: {
    type: Boolean,
    default: true
  },
  showFooter: {
    type: Boolean,
    default: true
  },
  showBreadcrumb: {
    type: Boolean,
    default: false
  },
  showPageHeader: {
    type: Boolean,
    default: true
  },
  
  // 页面配置
  pageTitle: {
    type: String,
    default: ''
  },
  pageDescription: {
    type: String,
    default: ''
  },
  breadcrumbs: {
    type: Array,
    default: () => []
  },
  
  // 搜索配置
  showSearch: {
    type: Boolean,
    default: false
  },
  searchPlaceholder: {
    type: String,
    default: '搜索...'
  },
  
  // 菜单配置
  menuItems: {
    type: Array,
    default: () => []
  },
  activeMenu: {
    type: String,
    default: ''
  },
  
  // 用户配置
  user: {
    type: Object,
    default: () => null
  },
  hideUserName: {
    type: Boolean,
    default: false
  },
  
  // 通知配置
  notifications: {
    type: Array,
    default: () => []
  },
  unreadCount: {
    type: Number,
    default: 0
  },
  
  // 主题配置
  isDarkMode: {
    type: Boolean,
    default: false
  },
  
  // 底部配置
  minimalFooter: {
    type: Boolean,
    default: false
  },
  fixedFooter: {
    type: Boolean,
    default: false
  },
  copyright: {
    type: String,
    default: ''
  },
  footerLinks: {
    type: Array,
    default: () => []
  },
  socialLinks: {
    type: Array,
    default: () => []
  },
  
  // 语言配置
  showLanguageSwitch: {
    type: Boolean,
    default: false
  },
  languages: {
    type: Array,
    default: () => []
  },
  currentLanguage: {
    type: Object,
    default: () => ({ code: 'zh-CN', name: '简体中文' })
  },
  
  // 状态配置
  showStatus: {
    type: Boolean,
    default: false
  },
  systemStatus: {
    type: String,
    default: 'online'
  },
  lastUpdated: {
    type: [Date, String, Number],
    default: () => new Date()
  },
  serverInfo: {
    type: String,
    default: ''
  },
  
  // 备案信息
  beianInfo: {
    type: Object,
    default: null
  }
})

const emit = defineEmits([
  'menu-click',
  'user-click',
  'user-command',
  'search',
  'toggle-fullscreen',
  'toggle-theme',
  'notification-command',
  'language-change',
  'footer-link-click'
])

const route = useRoute()
const router = useRouter()

// 响应式状态
const sidebarCollapsed = ref(false)
const mobileMenuOpen = ref(false)
const isMobile = ref(false)

// 计算属性
const layoutClasses = computed(() => ({
  'dark': props.isDarkMode,
  'mobile': isMobile.value,
  'sidebar-collapsed': sidebarCollapsed.value,
  'no-sidebar': !props.showSidebar
}))

// 方法
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('sidebar-collapsed', sidebarCollapsed.value.toString())
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const closeMobileMenu = () => {
  mobileMenuOpen.value = false
}

const handleMenuClick = (item) => {
  emit('menu-click', item)
  closeMobileMenu()
}

const handleUserClick = () => {
  emit('user-click')
}

const handleUserCommand = (command) => {
  emit('user-command', command)
}

const handleSearch = (query) => {
  emit('search', query)
}

const handleToggleFullscreen = (isFullscreen) => {
  emit('toggle-fullscreen', isFullscreen)
}

const handleToggleTheme = () => {
  emit('toggle-theme')
}

const handleNotificationCommand = (command) => {
  emit('notification-command', command)
}

const handleLanguageChange = (language) => {
  emit('language-change', language)
}

const handleFooterLinkClick = (link) => {
  emit('footer-link-click', link)
}

// 检测屏幕大小
const checkScreenSize = () => {
  isMobile.value = window.innerWidth <= 768
  if (isMobile.value) {
    mobileMenuOpen.value = false
  }
}

// 生命周期
onMounted(() => {
  // 恢复侧边栏状态
  const savedCollapsed = localStorage.getItem('sidebar-collapsed')
  if (savedCollapsed !== null) {
    sidebarCollapsed.value = savedCollapsed === 'true'
  }
  
  // 检测屏幕大小
  checkScreenSize()
  window.addEventListener('resize', checkScreenSize)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize)
})

// 导出一些方法供父组件使用
defineExpose({
  toggleSidebar,
  toggleMobileMenu,
  closeMobileMenu
})
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  background: #f8fafc;
}

.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: 250px;
  transition: margin-left 0.3s ease;
}

.app-main.sidebar-collapsed {
  margin-left: 64px;
}

.app-main.no-sidebar {
  margin-left: 0;
}

.mobile-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

.app-content {
  flex: 1;
  padding: 24px;
  overflow-x: auto;
}

.app-content.with-footer {
  padding-bottom: 100px;
}

.content-breadcrumb {
  margin-bottom: 16px;
  padding: 12px 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.page-header-left {
  flex: 1;
  min-width: 0;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.page-description {
  font-size: 16px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

.page-header-right {
  margin-left: 24px;
  flex-shrink: 0;
}

.page-content {
  flex: 1;
}

/* 深色模式 */
.app-layout.dark {
  background: #111827;
}

.dark .app-content {
  background: #111827;
}

.dark .page-header {
  border-bottom-color: #374151;
}

.dark .page-title {
  color: #f9fafb;
}

.dark .page-description {
  color: #9ca3af;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .app-main {
    margin-left: 0;
  }
  
  .app-main.sidebar-collapsed {
    margin-left: 0;
  }
  
  .app-content {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .page-header-right {
    margin-left: 0;
    align-self: stretch;
  }
  
  .page-title {
    font-size: 20px;
  }
}

@media (max-width: 480px) {
  .app-content {
    padding: 12px;
  }
  
  .page-title {
    font-size: 18px;
  }
  
  .page-description {
    font-size: 14px;
  }
}
</style>