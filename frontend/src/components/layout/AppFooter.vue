<template>
  <footer class="app-footer" :class="{ minimal: minimal, fixed: fixed }">
    <div class="footer-content">
      <!-- 左侧版权信息 -->
      <div class="footer-left">
        <div class="copyright">
          {{ copyright || `© ${currentYear} ${appName || 'TaskMaster'}. All rights reserved.` }}
        </div>
        <div v-if="version" class="version">
          v{{ version }}
        </div>
      </div>

      <!-- 中间链接 -->
      <div v-if="!minimal && links.length > 0" class="footer-center">
        <div class="footer-links">
          <a
            v-for="link in links"
            :key="link.name"
            :href="link.url"
            :target="link.external ? '_blank' : '_self'"
            class="footer-link"
            @click="handleLinkClick(link)"
          >
            {{ link.name }}
            <el-icon v-if="link.external" class="external-icon">
              <TopRight />
            </el-icon>
          </a>
        </div>
      </div>

      <!-- 右侧信息 -->
      <div class="footer-right">
        <!-- 语言切换 -->
        <el-dropdown v-if="showLanguage" @command="handleLanguageChange" placement="top-end">
          <el-button type="text" size="small" class="language-btn">
            <el-icon><Globe /></el-icon>
            {{ currentLanguage.name }}
            <el-icon><ArrowUp /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item
                v-for="lang in languages"
                :key="lang.code"
                :command="lang.code"
                :class="{ active: lang.code === currentLanguage.code }"
              >
                {{ lang.name }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <!-- 社交媒体链接 -->
        <div v-if="!minimal && socialLinks.length > 0" class="social-links">
          <a
            v-for="social in socialLinks"
            :key="social.name"
            :href="social.url"
            target="_blank"
            class="social-link"
            :title="social.name"
          >
            <el-icon><component :is="social.icon" /></el-icon>
          </a>
        </div>

        <!-- 状态指示器 -->
        <div v-if="showStatus" class="status-indicator">
          <el-popover placement="top" :width="200" trigger="hover">
            <template #reference>
              <div class="status-dot" :class="statusClass">
                <span class="status-text">{{ statusText }}</span>
              </div>
            </template>
            <div class="status-details">
              <div class="status-item">
                <span>服务状态:</span>
                <span :class="statusClass">{{ statusText }}</span>
              </div>
              <div class="status-item">
                <span>最后更新:</span>
                <span>{{ formatTime(lastUpdated) }}</span>
              </div>
              <div v-if="serverInfo" class="status-item">
                <span>服务器:</span>
                <span>{{ serverInfo }}</span>
              </div>
            </div>
          </el-popover>
        </div>
      </div>
    </div>

    <!-- 备案信息 -->
    <div v-if="!minimal && beian" class="beian-info">
      <a :href="beian.url" target="_blank" class="beian-link">
        {{ beian.text }}
      </a>
    </div>
  </footer>
</template>

<script setup>
import { computed, ref } from 'vue'
import {
  Globe,
  ArrowUp,
  TopRight
} from '@element-plus/icons-vue'

const props = defineProps({
  // 是否为简化模式
  minimal: {
    type: Boolean,
    default: false
  },
  // 是否固定在底部
  fixed: {
    type: Boolean,
    default: false
  },
  // 应用名称
  appName: {
    type: String,
    default: 'TaskMaster'
  },
  // 版本号
  version: {
    type: String,
    default: ''
  },
  // 自定义版权信息
  copyright: {
    type: String,
    default: ''
  },
  // 链接列表
  links: {
    type: Array,
    default: () => [
      { name: '关于我们', url: '/about' },
      { name: '隐私政策', url: '/privacy' },
      { name: '服务条款', url: '/terms' },
      { name: '帮助中心', url: '/help', external: true }
    ]
  },
  // 社交媒体链接
  socialLinks: {
    type: Array,
    default: () => []
  },
  // 是否显示语言切换
  showLanguage: {
    type: Boolean,
    default: false
  },
  // 可用语言
  languages: {
    type: Array,
    default: () => [
      { code: 'zh-CN', name: '简体中文' },
      { code: 'en-US', name: 'English' }
    ]
  },
  // 当前语言
  currentLanguage: {
    type: Object,
    default: () => ({ code: 'zh-CN', name: '简体中文' })
  },
  // 是否显示状态指示器
  showStatus: {
    type: Boolean,
    default: false
  },
  // 系统状态
  status: {
    type: String,
    default: 'online', // online, offline, maintenance
    validator: (value) => ['online', 'offline', 'maintenance'].includes(value)
  },
  // 最后更新时间
  lastUpdated: {
    type: [Date, String, Number],
    default: () => new Date()
  },
  // 服务器信息
  serverInfo: {
    type: String,
    default: ''
  },
  // 备案信息
  beian: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['language-change', 'link-click'])

const currentYear = computed(() => new Date().getFullYear())

const statusClass = computed(() => {
  switch (props.status) {
    case 'online':
      return 'status-online'
    case 'offline':
      return 'status-offline'
    case 'maintenance':
      return 'status-maintenance'
    default:
      return 'status-online'
  }
})

const statusText = computed(() => {
  switch (props.status) {
    case 'online':
      return '正常运行'
    case 'offline':
      return '服务离线'
    case 'maintenance':
      return '维护中'
    default:
      return '正常运行'
  }
})

const handleLanguageChange = (langCode) => {
  const language = props.languages.find(lang => lang.code === langCode)
  if (language) {
    emit('language-change', language)
  }
}

const handleLinkClick = (link) => {
  emit('link-click', link)
}

const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleString()
}
</script>

<style scoped>
.app-footer {
  background: #f8fafc;
  border-top: 1px solid #e2e8f0;
  padding: 24px;
  margin-top: auto;
}

.app-footer.fixed {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.app-footer.minimal {
  padding: 12px 24px;
  background: transparent;
  border-top: none;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.footer-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.copyright {
  color: #6b7280;
  font-size: 14px;
}

.version {
  color: #9ca3af;
  font-size: 12px;
  background: #e5e7eb;
  padding: 2px 8px;
  border-radius: 12px;
}

.footer-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.footer-links {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.footer-link {
  color: #6b7280;
  text-decoration: none;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: color 0.2s;
}

.footer-link:hover {
  color: #3b82f6;
}

.external-icon {
  font-size: 12px;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.language-btn {
  color: #6b7280;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.social-links {
  display: flex;
  gap: 8px;
}

.social-link {
  color: #6b7280;
  font-size: 18px;
  padding: 6px;
  border-radius: 6px;
  transition: all 0.2s;
  text-decoration: none;
}

.social-link:hover {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
}

.status-indicator {
  display: flex;
  align-items: center;
}

.status-dot {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}

.status-dot::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: block;
}

.status-online::before {
  background: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.status-offline::before {
  background: #ef4444;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
}

.status-maintenance::before {
  background: #f59e0b;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
}

.status-text {
  font-size: 12px;
  color: #6b7280;
}

.status-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  font-size: 12px;
}

.status-item span:first-child {
  color: #6b7280;
}

.status-item span:last-child {
  font-weight: 500;
}

.beian-info {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
  text-align: center;
}

.beian-link {
  color: #9ca3af;
  text-decoration: none;
  font-size: 12px;
}

.beian-link:hover {
  color: #6b7280;
}

/* 深色模式 */
.dark .app-footer {
  background: #1f2937;
  border-top-color: #374151;
}

.dark .copyright,
.dark .footer-link,
.dark .language-btn,
.dark .social-link,
.dark .status-text {
  color: #9ca3af;
}

.dark .footer-link:hover,
.dark .social-link:hover {
  color: #60a5fa;
}

.dark .version {
  background: #374151;
  color: #d1d5db;
}

.dark .beian-info {
  border-top-color: #374151;
}

.dark .beian-link {
  color: #6b7280;
}

.dark .beian-link:hover {
  color: #9ca3af;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .footer-content {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
  
  .footer-left {
    flex-direction: column;
    gap: 8px;
  }
  
  .footer-links {
    justify-content: center;
    gap: 16px;
  }
  
  .footer-right {
    flex-wrap: wrap;
    justify-content: center;
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .app-footer {
    padding: 16px;
  }
  
  .footer-links {
    flex-direction: column;
    gap: 8px;
  }
}
</style>