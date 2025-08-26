/**
 * 布局组件导出文件
 * 统一管理所有布局相关的组件和工具
 */

// 布局组件
export { default as AppLayout } from './AppLayout.vue'
export { default as AppSidebar } from './AppSidebar.vue'
export { default as AppHeader } from './AppHeader.vue'
export { default as AppFooter } from './AppFooter.vue'
export { default as DashboardLayout } from './DashboardLayout.vue'

// 布局配置和工具
export {
  defaultMenuItems,
  filterMenuByPermissions,
  generateBreadcrumbs,
  getActiveMenu,
  getPageConfig,
  pageConfigs,
  defaultFooterLinks,
  defaultSocialLinks,
  supportedLanguages
} from './layoutConfig'

/**
 * 使用示例：
 * 
 * 1. 基础布局使用：
 * import { AppLayout } from '@/components/layout'
 * 
 * <AppLayout
 *   :page-title="'我的页面'"
 *   :menu-items="menuItems"
 *   :user="user"
 * >
 *   <div>页面内容</div>
 * </AppLayout>
 * 
 * 2. Dashboard布局使用：
 * import { DashboardLayout } from '@/components/layout'
 * 
 * <DashboardLayout>
 *   <template #page-actions>
 *     <el-button type="primary">新建</el-button>
 *   </template>
 *   <div>Dashboard内容</div>
 * </DashboardLayout>
 * 
 * 3. 配置工具使用：
 * import { generateBreadcrumbs, getActiveMenu } from '@/components/layout'
 * 
 * const breadcrumbs = generateBreadcrumbs('/dashboard/todos')
 * const activeMenu = getActiveMenu('/dashboard/todos')
 */