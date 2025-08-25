<template>
  <div class="tools-development-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="navigateBack">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item>开发类</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>开发类工具</h1>
      <p>开发过程中的实用工具集合</p>
    </div>

    <div class="tools-grid">
      <div 
        v-for="tool in developmentTools" 
        :key="tool.id"
        class="tool-card"
        @click="navigateToTool(tool)"
      >
        <div class="tool-icon">
          <el-icon :size="32">
            <component :is="tool.icon" />
          </el-icon>
        </div>
        <div class="tool-content">
          <h3>{{ tool.name }}</h3>
          <p>{{ tool.description }}</p>
          <div class="tool-tags">
            <span 
              v-for="tag in tool.tags" 
              :key="tag" 
              class="tool-tag"
            >
              {{ tag }}
            </span>
          </div>
        </div>
        <div class="tool-arrow">
          <el-icon>
            <ArrowRight />
          </el-icon>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Timer,
  Document,
  Key,
  Connection,
  ArrowRight
} from '@element-plus/icons-vue'

const router = useRouter()

const developmentTools = reactive([
  {
    id: 'timestamp-converter',
    name: '时间戳转换',
    description: '时间戳与日期时间相互转换，支持多种格式',
    icon: Timer,
    tags: ['时间', '转换', '格式化'],
    route: 'tools-timestamp-converter'
  },
  {
    id: 'json-tools',
    name: 'JSON工具',
    description: 'JSON格式化、压缩、转YAML，支持语法高亮',
    icon: Document,
    tags: ['JSON', 'YAML', '格式化', '转换'],
    route: 'tools-json-tools'
  },
  {
    id: 'string-generator',
    name: '字符串生成',
    description: '生成随机字符串，支持自定义长度和字符集',
    icon: Key,
    tags: ['密码', '随机', '字符串', 'API密钥'],
    route: 'tools-string-generator'
  },
  {
    id: 'http-status-codes',
    name: 'HTTP状态码',
    description: 'HTTP状态码查询工具，包含详细说明和使用场景',
    icon: Connection,
    tags: ['HTTP', '状态码', 'API', '网络'],
    route: 'tools-http-status-codes'
  }
])

const navigateBack = () => {
  router.push({ name: 'tools' })
}

const navigateToTool = (tool) => {
  router.push({ name: tool.route })
}
</script>

<style scoped>
.tools-development-container {
  padding: 20px 0;
}

.tools-header {
  margin-bottom: 30px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.breadcrumb .el-breadcrumb-item {
  cursor: pointer;
}

.tools-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.tools-header p {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
  max-width: 1200px;
}

.tool-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 20px;
  border: 1px solid #e5e7eb;
}

.tool-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.1);
  border-color: #409eff;
}

.tool-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  flex-shrink: 0;
}

.tool-content {
  flex: 1;
}

.tool-content h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.tool-content p {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.tool-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tool-tag {
  font-size: 12px;
  color: #409eff;
  background: rgba(64, 158, 255, 0.1);
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid rgba(64, 158, 255, 0.2);
}

.tool-arrow {
  color: #9ca3af;
  transition: all 0.3s ease;
}

.tool-card:hover .tool-arrow {
  color: #409eff;
  transform: translateX(4px);
}

@media (max-width: 768px) {
  .tools-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .tool-card {
    padding: 20px;
    gap: 16px;
  }
  
  .tool-icon {
    width: 56px;
    height: 56px;
  }
  
  .tool-icon .el-icon {
    font-size: 24px !important;
  }
  
  .tool-content h3 {
    font-size: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
}
</style>