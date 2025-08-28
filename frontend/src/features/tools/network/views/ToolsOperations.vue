<template>
  <ToolsPageLayout
    category-name="运维类"
    title="运维类工具"
    description="服务器监控、网络工具、系统管理等工具"
    :tools="operationsTools"
    @tool-click="navigateToTool"
  />
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Monitor,
  Connection,
  Position
} from '@element-plus/icons-vue'
import ToolsPageLayout from '@/features/tools/shared/components/ToolsPageLayout.vue'

const router = useRouter()

const operationsTools = reactive([
  {
    id: 'ping-tool',
    name: 'Ping测试',
    description: '网络连通性测试工具，检测主机是否可达',
    icon: Connection,
    tags: ['网络', 'Ping', '连通性', '延迟'],
    route: 'tools-ping-tool'
  },
  {
    id: 'port-scanner',
    name: '端口扫描',
    description: '检测主机开放端口，分析网络服务状态',
    icon: Monitor,
    tags: ['端口', '扫描', '网络', '服务'],
    route: 'tools-port-scanner'
  },
  {
    id: 'dns-lookup',
    name: 'DNS查询',
    description: 'DNS解析查询工具，查看域名解析信息',
    icon: Position,
    tags: ['DNS', '域名', '解析', '查询'],
    route: 'tools-dns-lookup'
  }
])

const navigateToTool = (tool) => {
  if (tool.id === 'ping-tool' || tool.id === 'port-scanner') {
    router.push({ name: tool.route })
  } else {
    ElMessage.info(`${tool.name} 即将上线，敬请期待！`)
  }
}
</script>