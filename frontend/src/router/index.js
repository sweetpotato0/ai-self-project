import { useAuthStore } from '@/stores/auth'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresGuest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('@/views/DashboardHome.vue')
      },
      {
        path: 'todos',
        name: 'todos',
        component: () => import('@/views/Todos.vue')
      },
      {
        path: 'calendar',
        name: 'calendar',
        component: () => import('@/views/Calendar.vue')
      },
      {
        path: 'analytics',
        name: 'analytics',
        component: () => import('@/views/Analytics.vue')
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/views/Settings.vue')
      },
      {
        path: 'profile',
        name: 'profile',
        component: () => import('@/views/Profile.vue')
      },
      {
        path: 'articles',
        name: 'articles',
        component: () => import('@/views/Articles.vue')
      },
      {
        path: 'articles/:id',
        name: 'article-detail',
        component: () => import('@/views/ArticleDetail.vue')
      },
      {
        path: 'notifications',
        name: 'notifications',
        component: () => import('@/views/Notifications.vue')
      },
      {
        path: 'tools',
        name: 'tools',
        component: () => import('@/views/Tools.vue')
      },
      {
        path: 'tools/development',
        name: 'tools-development',
        component: () => import('@/views/ToolsDevelopment.vue')
      },
      {
        path: 'tools/timestamp-converter',
        name: 'tools-timestamp-converter',
        component: () => import('@/views/TimestampConverter.vue')
      },
      {
        path: 'tools/json-tools',
        name: 'tools-json-tools',
        component: () => import('@/views/JsonToolsSimple.vue')
      },
      {
        path: 'tools/string-generator',
        name: 'tools-string-generator',
        component: () => import('@/views/StringGenerator.vue')
      },
      {
        path: 'tools/http-status-codes',
        name: 'tools-http-status-codes',
        component: () => import('@/views/HttpStatusCodes.vue')
      }
    ]
  },
  // 独立的路由，可以直接访问 /articles/:id
  {
    path: '/articles/:id',
    name: 'article-detail-standalone',
    component: () => import('@/views/ArticleDetail.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const isAuthenticated = authStore.isAuthenticated

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresGuest && isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
