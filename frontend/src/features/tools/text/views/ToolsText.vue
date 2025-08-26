<template>
  <div class="tools-text-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="navigateBack">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item>文本类</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>文本类工具</h1>
      <p>文本处理、格式转换、编码解码等工具</p>
    </div>

    <div class="tools-grid">
      <div 
        v-for="tool in textTools" 
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
  Document,
  Key,
  Connection,
  ArrowRight
} from '@element-plus/icons-vue'

const router = useRouter()

const textTools = reactive([
  {
    id: 'base64-encoder',
    name: 'Base64编码',
    description: 'Base64编码解码工具，支持文本和文件',
    icon: Key,
    tags: ['Base64', '编码', '解码', '加密'],
    route: 'tools-base64-encoder'
  },
  {
    id: 'url-encoder',
    name: 'URL编码',
    description: 'URL编码解码工具，处理特殊字符转换',
    icon: Connection,
    tags: ['URL', '编码', '解码', '转义'],
    route: 'tools-url-encoder'
  },
  {
    id: 'text-processor',
    name: '文本处理器',
    description: '文本格式化、大小写转换、去重等处理',
    icon: Document,
    tags: ['文本', '格式化', '转换', '处理'],
    route: 'tools-text-processor'
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
.tools-text-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.tools-header {
  margin-bottom: 30px;
  text-align: left;
}

.tools-header h1,
.tools-header p {
  text-align: center;
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
  width: 100%;
  justify-content: center;
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