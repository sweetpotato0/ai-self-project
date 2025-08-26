<template>
  <ToolsPageLayout
    category-name="学术类"
    title="学术类工具"
    description="论文工具、数学计算、学术格式等工具"
    :tools="academicTools"
    @tool-click="navigateToTool"
  />
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Reading,
  Operation,
  DataAnalysis
} from '@element-plus/icons-vue'
import ToolsPageLayout from '@/features/tools/shared/components/ToolsPageLayout.vue'

const router = useRouter()

const academicTools = reactive([
  {
    id: 'citation-generator',
    name: '引用生成器',
    description: '生成各种学术引用格式，支持APA、MLA、Chicago等',
    icon: Reading,
    tags: ['引用', 'APA', 'MLA', '论文'],
    route: 'tools-citation-generator'
  },
  {
    id: 'math-calculator',
    name: '数学计算器',
    description: '高级数学计算工具，支持复杂函数和公式计算',
    icon: Operation,
    tags: ['数学', '计算', '函数', '公式'],
    route: 'tools-math-calculator'
  },
  {
    id: 'data-analyzer',
    name: '数据分析',
    description: '简单的数据分析工具，统计计算和图表生成',
    icon: DataAnalysis,
    tags: ['数据', '统计', '分析', '图表'],
    route: 'tools-data-analyzer'
  }
])

const navigateToTool = (tool) => {
  if (tool.id === 'citation-generator' || tool.id === 'math-calculator') {
    router.push({ name: tool.route })
  } else {
    ElMessage.info(`${tool.name} 即将上线，敬请期待！`)
  }
}
</script>