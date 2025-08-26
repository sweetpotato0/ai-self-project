<template>
  <ToolsPageLayout
    category-name="查询类"
    title="查询类工具"
    description="IP查询、域名查询、whois等查询工具"
    :tools="queryTools"
    @tool-click="navigateToTool"
  />
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Search,
  Position,
  Connection
} from '@element-plus/icons-vue'
import ToolsPageLayout from '@/features/tools/shared/components/ToolsPageLayout.vue'

const router = useRouter()

const queryTools = reactive([
  {
    id: 'ip-lookup',
    name: 'IP地址查询',
    description: 'IP地址归属地查询，获取地理位置和ISP信息',
    icon: Position,
    tags: ['IP', '地址', '归属地', 'ISP'],
    route: 'tools-ip-lookup'
  },
  {
    id: 'whois-lookup',
    name: 'Whois查询',
    description: '域名whois信息查询，获取注册信息和DNS记录',
    icon: Search,
    tags: ['Whois', '域名', '注册', 'DNS'],
    route: 'tools-whois-lookup'
  },
  {
    id: 'domain-info',
    name: '域名信息',
    description: '域名详细信息查询，包括SSL证书和网站信息',
    icon: Connection,
    tags: ['域名', 'SSL', '证书', '网站'],
    route: 'tools-domain-info'
  }
])

const navigateToTool = (tool) => {
  router.push({ name: tool.route })
}
</script>