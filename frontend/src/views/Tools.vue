<template>
  <div class="tools-container">
    <div class="tools-header">
      <h1>工具箱</h1>
      <p>实用工具集合，提升工作效率</p>
    </div>

    <div class="tools-grid">
      <div 
        v-for="category in toolCategories" 
        :key="category.id"
        class="category-card"
        @click="navigateToCategory(category)"
      >
        <div class="category-icon">
          <el-icon :size="40">
            <component :is="category.icon" />
          </el-icon>
        </div>
        <div class="category-content">
          <h3>{{ category.name }}</h3>
          <p>{{ category.description }}</p>
          <div class="category-count">
            {{ category.toolCount }} 个工具
          </div>
        </div>
        <div class="category-arrow">
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
  Picture, 
  Monitor, 
  Tools as ToolsIcon, 
  Connection, 
  Reading, 
  Search, 
  MoreFilled,
  ArrowRight
} from '@element-plus/icons-vue'

const router = useRouter()

const toolCategories = reactive([
  {
    id: 'text',
    name: '文本类',
    description: '文本处理、格式转换、编码解码等工具',
    icon: Document,
    toolCount: 3,
    route: 'tools-text'
  },
  {
    id: 'document',
    name: '文档类',
    description: '文档转换、PDF处理、格式化等工具',
    icon: Reading,
    toolCount: 0,
    route: 'tools-document'
  },
  {
    id: 'image',
    name: '图像类',
    description: '图片压缩、格式转换、编辑处理等工具',
    icon: Picture,
    toolCount: 3,
    route: 'tools-image'
  },
  {
    id: 'development',
    name: '开发类',
    description: '开发工具、代码格式化、API测试等工具',
    icon: Monitor,
    toolCount: 4,
    route: 'tools-development'
  },
  {
    id: 'operations',
    name: '运维类',
    description: '服务器监控、网络工具、系统管理等工具',
    icon: Connection,
    toolCount: 3,
    route: 'tools-operations'
  },
  {
    id: 'academic',
    name: '学术类',
    description: '论文工具、数学计算、学术格式等工具',
    icon: Reading,
    toolCount: 3,
    route: 'tools-academic'
  },
  {
    id: 'query',
    name: '查询类',
    description: 'IP查询、域名查询、whois等查询工具',
    icon: Search,
    toolCount: 3,
    route: 'tools-query'
  },
  {
    id: 'others',
    name: '其它',
    description: '其他实用工具和小功能',
    icon: MoreFilled,
    toolCount: 5,
    route: 'tools-others'
  }
])

const navigateToCategory = (category) => {
  router.push({ name: category.route })
}
</script>

<style scoped>
.tools-container {
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.tools-header {
  text-align: center;
  margin-bottom: 40px;
}

.tools-header h1 {
  font-size: 32px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 10px 0;
}

.tools-header p {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 24px;
  max-width: 1200px;
  width: 100%;
  justify-content: center;
}

.category-card {
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

.category-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.1);
  border-color: #409eff;
}

.category-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  flex-shrink: 0;
}

.category-content {
  flex: 1;
}

.category-content h3 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.category-content p {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.category-count {
  font-size: 12px;
  color: #409eff;
  font-weight: 500;
  background: rgba(64, 158, 255, 0.1);
  padding: 4px 8px;
  border-radius: 6px;
  display: inline-block;
}

.category-arrow {
  color: #9ca3af;
  transition: all 0.3s ease;
}

.category-card:hover .category-arrow {
  color: #409eff;
  transform: translateX(4px);
}

@media (max-width: 768px) {
  .tools-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .category-card {
    padding: 20px;
    gap: 16px;
  }
  
  .category-icon {
    width: 60px;
    height: 60px;
  }
  
  .category-icon .el-icon {
    font-size: 24px !important;
  }
  
  .category-content h3 {
    font-size: 18px;
  }
  
  .tools-header h1 {
    font-size: 28px;
  }
}
</style>