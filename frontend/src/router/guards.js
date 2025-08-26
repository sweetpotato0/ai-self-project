import { useAuthStore } from '@/stores/auth'

// 路由守卫
export function setupRouterGuards(router) {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()
    const isAuthenticated = authStore.isAuthenticated

    // 检查认证要求
    if (to.meta.requiresAuth && !isAuthenticated) {
      next('/login')
    } else if (to.meta.requiresGuest && isAuthenticated) {
      next('/dashboard')
    } else {
      next()
    }
  })

  // 页面标题守卫
  router.afterEach((to) => {
    const title = to.meta.title || '工具箱管理系统'
    document.title = title
  })
}