// 仪表盘相关路由配置
export const dashboardRoutes = [
  {
    path: '',
    name: 'dashboard',
    component: () => import('@/features/dashboard/views/DashboardHome.vue')
  },
  {
    path: 'todos',
    name: 'todos',
    component: () => import('@/features/todos/views/Todos.vue')
  },
  {
    path: 'calendar',
    name: 'calendar',
    component: () => import('@/features/calendar/views/Calendar.vue')
  },
  {
    path: 'analytics',
    name: 'analytics',
    component: () => import('@/features/analytics/views/Analytics.vue')
  },
  {
    path: 'settings',
    name: 'settings',
    component: () => import('@/features/settings/views/Settings.vue')
  },
  {
    path: 'profile',
    name: 'profile',
    component: () => import('@/features/profile/views/Profile.vue')
  },
  {
    path: 'articles',
    name: 'articles',
    component: () => import('@/features/articles/views/Articles.vue')
  },
  {
    path: 'articles/:id',
    name: 'article-detail',
    component: () => import('@/features/articles/views/ArticleDetail.vue')
  },
  {
    path: 'notifications',
    name: 'notifications',
    component: () => import('@/features/notifications/views/Notifications.vue')
  }
]