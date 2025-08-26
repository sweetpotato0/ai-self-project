// 认证相关路由配置
export const authRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/features/auth/views/Login.vue'),
    meta: { requiresGuest: true }
  }
]