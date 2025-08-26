import { createRouter, createWebHistory } from 'vue-router'
import { setupRouterGuards } from './guards'
import { authRoutes } from './modules/authRoutes'
import { dashboardRoutes } from './modules/dashboardRoutes'
import { toolsRoutes } from './modules/toolsRoutes'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  // 认证路由
  ...authRoutes,
  // 仪表盘主路由
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/features/dashboard/views/Dashboard.vue'),
    meta: { requiresAuth: true },
    children: [
      // 仪表盘子路由
      ...dashboardRoutes,
      // 工具路由
      ...toolsRoutes
    ]
  },
  // 独立的文章详情路由
  {
    path: '/articles/:id',
    name: 'article-detail-standalone',
    component: () => import('@/features/articles/views/ArticleDetail.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 设置路由守卫
setupRouterGuards(router)

export default router