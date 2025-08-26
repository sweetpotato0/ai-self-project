import {
  DataBoard,
  List,
  Calendar,
  Document,
  Setting,
  DataAnalysis,
  Bell,
  User,
  Tools,
  Connection,
  Search,
  MagicStick,
  Cpu,
  Picture,
  EditPen,
  More,
  School
} from '@element-plus/icons-vue'

/**
 * 主要菜单配置
 */
export const defaultMenuItems = [
  {
    index: 'dashboard',
    title: '仪表盘',
    icon: DataBoard,
    route: '/dashboard'
  },
  {
    index: 'todos',
    title: 'TODO清单',
    icon: List,
    route: '/dashboard/todos'
  },
  {
    index: 'calendar',
    title: '日程安排',
    icon: Calendar,
    route: '/dashboard/calendar'
  },
  {
    index: 'articles',
    title: '文章管理',
    icon: Document,
    route: '/dashboard/articles'
  },
  {
    index: 'analytics',
    title: '数据分析',
    icon: DataAnalysis,
    route: '/dashboard/analytics'
  },
  {
    index: 'tools',
    title: '工具箱',
    icon: Tools,
    route: '/dashboard/tools',
    children: [
      {
        index: 'tools-development',
        title: '开发工具',
        icon: Cpu,
        route: '/dashboard/tools/development'
      },
      {
        index: 'tools-text',
        title: '文本处理',
        icon: EditPen,
        route: '/dashboard/tools/text'
      },
      {
        index: 'tools-image',
        title: '图像处理',
        icon: Picture,
        route: '/dashboard/tools/image'
      },
      {
        index: 'tools-network',
        title: '网络工具',
        icon: Connection,
        route: '/dashboard/tools/network'
      },
      {
        index: 'tools-query',
        title: '查询工具',
        icon: Search,
        route: '/dashboard/tools/query'
      },
      {
        index: 'tools-academic',
        title: '学术工具',
        icon: School,
        route: '/dashboard/tools/academic'
      },
      {
        index: 'tools-others',
        title: '其他工具',
        icon: More,
        route: '/dashboard/tools/others'
      }
    ]
  },
  {
    index: 'notifications',
    title: '通知中心',
    icon: Bell,
    route: '/dashboard/notifications'
  },
  {
    index: 'profile',
    title: '个人资料',
    icon: User,
    route: '/dashboard/profile'
  },
  {
    index: 'settings',
    title: '系统设置',
    icon: Setting,
    route: '/dashboard/settings'
  }
]

/**
 * 根据用户权限过滤菜单项
 * @param {Array} menuItems - 菜单项
 * @param {Object} user - 用户信息
 * @param {Array} permissions - 用户权限
 * @returns {Array} 过滤后的菜单项
 */
export function filterMenuByPermissions(menuItems, user, permissions = []) {
  if (!user) return []
  
  return menuItems.filter(item => {
    // 如果菜单项有权限要求
    if (item.permission && !permissions.includes(item.permission)) {
      return false
    }
    
    // 递归过滤子菜单
    if (item.children) {
      item.children = filterMenuByPermissions(item.children, user, permissions)
    }
    
    return true
  })
}

/**
 * 获取面包屑导航
 * @param {string} currentRoute - 当前路由
 * @param {Array} menuItems - 菜单项
 * @returns {Array} 面包屑导航数组
 */
export function generateBreadcrumbs(currentRoute, menuItems = defaultMenuItems) {
  const breadcrumbs = []
  
  function findPath(items, targetRoute, path = []) {
    for (const item of items) {
      const currentPath = [...path, item]
      
      if (item.route === targetRoute) {
        return currentPath
      }
      
      if (item.children) {
        const found = findPath(item.children, targetRoute, currentPath)
        if (found) return found
      }
    }
    return null
  }
  
  const path = findPath(menuItems, currentRoute)
  if (path) {
    return path.map(item => ({
      title: item.title,
      path: item.route,
      icon: item.icon
    }))
  }
  
  return breadcrumbs
}

/**
 * 获取当前激活菜单
 * @param {string} currentRoute - 当前路由
 * @param {Array} menuItems - 菜单项
 * @returns {string} 激活菜单索引
 */
export function getActiveMenu(currentRoute, menuItems = defaultMenuItems) {
  function findActive(items) {
    for (const item of items) {
      if (currentRoute.startsWith(item.route)) {
        if (item.children) {
          // 如果有子菜单，递归查找更精确的匹配
          const childActive = findActive(item.children)
          if (childActive) return childActive
        }
        return item.index
      }
    }
    return ''
  }
  
  return findActive(menuItems)
}

/**
 * 页面配置映射
 */
export const pageConfigs = {
  '/dashboard': {
    title: '仪表盘',
    description: '查看系统概览和重要数据',
    showSearch: false
  },
  '/dashboard/todos': {
    title: 'TODO清单',
    description: '管理您的待办事项和任务',
    showSearch: true,
    searchPlaceholder: '搜索任务...'
  },
  '/dashboard/calendar': {
    title: '日程安排',
    description: '查看和管理您的日程安排',
    showSearch: false
  },
  '/dashboard/articles': {
    title: '文章管理',
    description: '创建、编辑和管理您的文章',
    showSearch: true,
    searchPlaceholder: '搜索文章...'
  },
  '/dashboard/analytics': {
    title: '数据分析',
    description: '查看系统数据和统计分析',
    showSearch: false
  },
  '/dashboard/tools': {
    title: '工具箱',
    description: '实用工具集合',
    showSearch: true,
    searchPlaceholder: '搜索工具...'
  },
  '/dashboard/notifications': {
    title: '通知中心',
    description: '查看和管理系统通知',
    showSearch: true,
    searchPlaceholder: '搜索通知...'
  },
  '/dashboard/profile': {
    title: '个人资料',
    description: '管理您的个人信息和偏好设置',
    showSearch: false
  },
  '/dashboard/settings': {
    title: '系统设置',
    description: '配置系统参数和个人偏好',
    showSearch: false
  }
}

/**
 * 获取页面配置
 * @param {string} route - 当前路由
 * @returns {Object} 页面配置对象
 */
export function getPageConfig(route) {
  // 精确匹配
  if (pageConfigs[route]) {
    return pageConfigs[route]
  }
  
  // 前缀匹配
  for (const [path, config] of Object.entries(pageConfigs)) {
    if (route.startsWith(path) && path !== '/dashboard') {
      return config
    }
  }
  
  // 默认配置
  return {
    title: '',
    description: '',
    showSearch: false,
    searchPlaceholder: '搜索...'
  }
}

/**
 * 默认底部链接配置
 */
export const defaultFooterLinks = [
  { name: '关于我们', url: '/about' },
  { name: '隐私政策', url: '/privacy' },
  { name: '服务条款', url: '/terms' },
  { name: '帮助中心', url: '/help', external: true }
]

/**
 * 默认社交媒体链接
 */
export const defaultSocialLinks = [
  // 可以根据需要添加社交媒体链接
  // { name: 'GitHub', url: 'https://github.com', icon: 'GitHubIcon' }
]

/**
 * 支持的语言配置
 */
export const supportedLanguages = [
  { code: 'zh-CN', name: '简体中文' },
  { code: 'en-US', name: 'English' },
  { code: 'zh-TW', name: '繁體中文' },
  { code: 'ja-JP', name: '日本語' }
]

export default {
  defaultMenuItems,
  filterMenuByPermissions,
  generateBreadcrumbs,
  getActiveMenu,
  getPageConfig,
  pageConfigs,
  defaultFooterLinks,
  defaultSocialLinks,
  supportedLanguages
}