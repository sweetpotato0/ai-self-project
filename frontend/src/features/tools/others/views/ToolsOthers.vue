<template>
  <ToolsPageLayout
    category-name="其它"
    title="其它工具"
    description="其他实用工具和小功能"
    :tools="otherTools"
    @tool-click="navigateToTool"
  />
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Grid,
  Brush,
  Key,
  Setting,
  MoreFilled
} from '@element-plus/icons-vue'
import ToolsPageLayout from '@/features/tools/shared/components/ToolsPageLayout.vue'

const router = useRouter()

const otherTools = reactive([
  {
    id: 'qr-code-generator',
    name: '二维码生成器',
    description: '生成各种类型的二维码，支持文本、网址、WiFi等格式',
    icon: Grid,
    tags: ['二维码', '生成器', 'WiFi', '分享'],
    route: 'tools-qr-code-generator'
  },
  {
    id: 'color-picker',
    name: '颜色选择器',
    description: '专业的颜色选择和转换工具，支持多种颜色格式',
    icon: Brush,
    tags: ['颜色', 'HEX', 'RGB', 'HSL'],
    route: 'tools-color-picker'
  },
  {
    id: 'password-strength-checker',
    name: '密码强度检测',
    description: '检测密码强度并提供安全建议，生成安全密码',
    icon: Key,
    tags: ['密码', '安全', '强度', '生成'],
    route: 'tools-password-strength-checker'
  },
  {
    id: 'regex-tester',
    name: '正则表达式测试',
    description: '测试和验证正则表达式，支持实时匹配和替换',
    icon: Setting,
    tags: ['正则', '匹配', '测试', '替换'],
    route: 'tools-regex-tester'
  },
  {
    id: 'hash-calculator',
    name: 'Hash计算器',
    description: '计算文本或文件的Hash值，支持MD5、SHA等算法',
    icon: MoreFilled,
    tags: ['Hash', 'MD5', 'SHA', '校验'],
    route: 'tools-hash-calculator'
  }
])

const navigateToTool = (tool) => {
  router.push({ name: tool.route })
}
</script>