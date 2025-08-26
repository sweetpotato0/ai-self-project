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
      },
      {
        path: 'tools/text',
        name: 'tools-text',
        component: () => import('@/views/ToolsText.vue')
      },
      {
        path: 'tools/base64-encoder',
        name: 'tools-base64-encoder',
        component: () => import('@/views/Base64Encoder.vue')
      },
      {
        path: 'tools/url-encoder',
        name: 'tools-url-encoder',
        component: () => import('@/views/UrlEncoder.vue')
      },
      {
        path: 'tools/text-processor',
        name: 'tools-text-processor',
        component: () => import('@/views/TextProcessor.vue')
      },
      {
        path: 'tools/image',
        name: 'tools-image',
        component: () => import('@/views/ToolsImage.vue')
      },
      {
        path: 'tools/image-compressor',
        name: 'tools-image-compressor',
        component: () => import('@/views/ImageCompressor.vue')
      },
      {
        path: 'tools/image-converter',
        name: 'tools-image-converter',
        component: () => import('@/views/ImageConverter.vue')
      },
      {
        path: 'tools/image-resizer',
        name: 'tools-image-resizer',
        component: () => import('@/views/ImageResizer.vue')
      },
      {
        path: 'tools/operations',
        name: 'tools-operations',
        component: () => import('@/views/ToolsOperations.vue')
      },
      {
        path: 'tools/academic',
        name: 'tools-academic',
        component: () => import('@/views/ToolsAcademic.vue')
      },
      {
        path: 'tools/query',
        name: 'tools-query',
        component: () => import('@/views/ToolsQuery.vue')
      },
      {
        path: 'tools/document',
        name: 'tools-document',
        component: () => import('@/views/ToolsDocument.vue')
      },
      {
        path: 'tools/others',
        name: 'tools-others',
        component: () => import('@/views/ToolsOthers.vue')
      },
      {
        path: 'tools/citation-generator',
        name: 'tools-citation-generator',
        component: () => import('@/views/CitationGenerator.vue')
      },
      {
        path: 'tools/math-calculator',
        name: 'tools-math-calculator',
        component: () => import('@/views/MathCalculator.vue')
      },
      {
        path: 'tools/ping-tool',
        name: 'tools-ping-tool',
        component: () => import('@/views/PingTool.vue')
      },
      // 其他工具
      {
        path: 'tools/qr-code-generator',
        name: 'tools-qr-code-generator',
        component: () => import('@/views/QRCodeGenerator.vue')
      },
      {
        path: 'tools/color-picker',
        name: 'tools-color-picker',
        component: () => import('@/views/ColorPicker.vue')
      },
      {
        path: 'tools/password-strength-checker',
        name: 'tools-password-strength-checker',
        component: () => import('@/views/PasswordStrengthChecker.vue')
      },
      {
        path: 'tools/regex-tester',
        name: 'tools-regex-tester',
        component: () => import('@/views/RegexTester.vue')
      },
      {
        path: 'tools/hash-calculator',
        name: 'tools-hash-calculator',
        component: () => import('@/views/HashCalculator.vue')
      },
      // 查询类工具
      {
        path: 'tools/ip-lookup',
        name: 'tools-ip-lookup',
        component: () => import('@/views/IPLookup.vue')
      },
      {
        path: 'tools/whois-lookup',
        name: 'tools-whois-lookup',
        component: () => import('@/views/WhoisLookup.vue')
      },
      {
        path: 'tools/domain-info',
        name: 'tools-domain-info',
        component: () => import('@/views/DomainInfo.vue')
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
