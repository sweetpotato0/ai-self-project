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
  },
  {
    path: 'audit-logs',
    name: 'audit-logs',
    component: () => import('@/features/audit/views/AuditLogs.vue')
  },
  {
    path: 'english-learning',
    name: 'english-learning',
    component: () => import('@/features/english-learning/views/EnglishLearning.vue')
  },
  {
    path: 'english-learning/play/:id',
    name: 'song-player',
    component: () => import('@/features/english-learning/views/SongPlayer.vue')
  },
  {
    path: 'english-videos',
    name: 'EnglishVideos',
    component: () => import('@/features/english-videos/views/EnglishVideos.vue')
  },
  {
    path: 'english-videos/series/:seriesId',
    name: 'VideoSeries',
    component: () => import('@/features/english-videos/views/VideoPlayer.vue')
  },
  {
    path: 'english-videos/play/:seriesId/:episodeId',
    name: 'VideoPlayer',
    component: () => import('@/features/english-videos/views/VideoPlayer.vue')
  },
  {
    path: 'admin/english-videos',
    name: 'admin-english-videos',
    component: () => import('@/features/english-videos/views/admin/AdminDashboard.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: 'admin/english-videos/series',
    name: 'admin-video-series',
    component: () => import('@/features/english-videos/views/admin/VideoSeriesAdmin.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: 'admin/english-videos/episodes',
    name: 'admin-video-episodes',
    component: () => import('@/features/english-videos/views/admin/EpisodeAdmin.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  }
]