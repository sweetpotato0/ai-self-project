<template>
  <div class="sidebar" :class="{ collapsed: collapsed }">
    <div class="sidebar-header">
      <div class="logo">
        <el-icon class="logo-icon"><Grid /></el-icon>
        <span v-if="!collapsed" class="logo-text">{{ title || 'TaskMaster' }}</span>
      </div>
      <el-button
        v-if="!collapsed"
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
      :collapse="collapsed"
      background-color="transparent"
    >
      <template v-if="collapsed">
        <!-- 折叠状态的菜单项 -->
        <template v-for="item in menuItems" :key="item.index">
          <el-tooltip :content="item.title" placement="right">
            <el-menu-item :index="item.index" @click="handleMenuClick(item)">
              <el-icon><component :is="item.icon" /></el-icon>
            </el-menu-item>
          </el-tooltip>
        </template>
      </template>

      <template v-else>
        <!-- 展开状态的菜单项 -->
        <template v-for="item in menuItems" :key="item.index">
          <el-sub-menu v-if="item.children" :index="item.index">
            <template #title>
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item 
              v-for="child in item.children" 
              :key="child.index"
              :index="child.index"
              @click="handleMenuClick(child)"
            >
              <el-icon><component :is="child.icon" /></el-icon>
              <span>{{ child.title }}</span>
            </el-menu-item>
          </el-sub-menu>
          
          <el-menu-item v-else :index="item.index" @click="handleMenuClick(item)">
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.title }}</span>
          </el-menu-item>
        </template>
      </template>
    </el-menu>

    <!-- 底部用户信息 -->
    <div class="sidebar-footer">
      <div v-if="!collapsed" class="user-info" @click="handleUserClick">
        <el-avatar :size="32" :src="user?.avatar">
          <el-icon><User /></el-icon>
        </el-avatar>
        <div class="user-details">
          <div class="username">{{ user?.username || '用户' }}</div>
          <div class="user-status">{{ user?.email || '未登录' }}</div>
        </div>
        <el-icon class="user-arrow"><ArrowRight /></el-icon>
      </div>
      <el-tooltip v-else content="用户设置" placement="right">
        <el-avatar :size="32" :src="user?.avatar" @click="handleUserClick" class="user-avatar-collapsed">
          <el-icon><User /></el-icon>
        </el-avatar>
      </el-tooltip>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  Grid,
  Fold,
  Expand,
  User,
  ArrowRight
} from '@element-plus/icons-vue'

const props = defineProps({
  // 是否折叠
  collapsed: {
    type: Boolean,
    default: false
  },
  // 应用标题
  title: {
    type: String,
    default: 'TaskMaster'
  },
  // 当前激活菜单
  activeMenu: {
    type: String,
    default: ''
  },
  // 菜单项配置
  menuItems: {
    type: Array,
    default: () => []
  },
  // 用户信息
  user: {
    type: Object,
    default: () => null
  }
})

const emit = defineEmits(['toggle', 'menu-click', 'user-click'])
const router = useRouter()

const toggleSidebar = () => {
  emit('toggle')
}

const handleMenuClick = (item) => {
  emit('menu-click', item)
  
  // 如果有路由，则进行导航
  if (item.route) {
    router.push(item.route)
  }
}

const handleUserClick = () => {
  emit('user-click')
}
</script>

<style scoped>
.sidebar {
  width: 250px;
  height: 100vh;
  background: linear-gradient(180deg, #1e293b 0%, #334155 100%);
  border-right: 1px solid #334155;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  position: fixed;
  left: 0;
  top: 0;
  z-index: 1000;
}

.sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  color: white;
}

.logo-icon {
  font-size: 24px;
  color: #3b82f6;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
}

.collapse-btn {
  color: #94a3b8;
  padding: 4px;
}

.collapse-btn:hover {
  color: white;
  background: rgba(255, 255, 255, 0.1);
}

.collapsed-toggle {
  margin-left: auto;
}

.sidebar-menu {
  flex: 1;
  border: none;
  padding: 16px 0;
  overflow-y: auto;
}

.sidebar-menu :deep(.el-menu-item) {
  color: #94a3b8;
  height: 48px;
  line-height: 48px;
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.2s;
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
}

.sidebar-menu :deep(.el-sub-menu__title) {
  color: #94a3b8;
  height: 48px;
  line-height: 48px;
  margin: 4px 12px;
  border-radius: 8px;
}

.sidebar-menu :deep(.el-sub-menu__title:hover) {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  color: #94a3b8;
}

.user-info:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.username {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-status {
  font-size: 12px;
  opacity: 0.7;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-arrow {
  font-size: 16px;
  opacity: 0.5;
}

.user-avatar-collapsed {
  cursor: pointer;
  transition: all 0.2s;
}

.user-avatar-collapsed:hover {
  transform: scale(1.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 100%;
    transform: translateX(-100%);
  }
  
  .sidebar.collapsed {
    width: 100%;
    transform: translateX(-100%);
  }
  
  .sidebar.mobile-open {
    transform: translateX(0);
  }
}
</style>