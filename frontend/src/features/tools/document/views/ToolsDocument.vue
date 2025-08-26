<template>
  <div class="tools-document-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="navigateBack">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item>文档类</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>文档类工具</h1>
      <p>文档转换、PDF处理、格式化等工具</p>
    </div>

    <div class="coming-soon-container">
      <div class="coming-soon-content">
        <div class="icon-container">
          <el-icon :size="120">
            <Reading />
          </el-icon>
        </div>
        <h2>即将上线</h2>
        <p>文档类工具正在开发中，敬请期待！</p>
        <div class="planned-tools">
          <h3>计划功能：</h3>
          <ul>
            <li>PDF合并与拆分</li>
            <li>Word文档转换</li>
            <li>Markdown编辑器</li>
            <li>文档格式转换</li>
          </ul>
        </div>
        <el-button @click="navigateBack" type="primary" size="large">
          返回工具箱
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { Reading } from '@element-plus/icons-vue'

const router = useRouter()

const navigateBack = () => {
  router.push({ name: 'tools' })
}
</script>

<style scoped>
.tools-document-container {
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

.coming-soon-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  width: 100%;
}

.coming-soon-content {
  text-align: center;
  max-width: 500px;
  padding: 40px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.icon-container {
  margin-bottom: 24px;
  color: #667eea;
}

.coming-soon-content h2 {
  font-size: 32px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 16px 0;
}

.coming-soon-content > p {
  font-size: 18px;
  color: #6b7280;
  margin: 0 0 32px 0;
}

.planned-tools {
  text-align: left;
  margin: 32px 0;
  padding: 24px;
  background: #f8fafc;
  border-radius: 8px;
}

.planned-tools h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 16px 0;
}

.planned-tools ul {
  margin: 0;
  padding-left: 20px;
  list-style-type: disc;
}

.planned-tools li {
  font-size: 14px;
  color: #6b7280;
  margin: 8px 0;
}

@media (max-width: 768px) {
  .tools-document-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .coming-soon-content {
    padding: 24px;
    margin: 0 16px;
  }
  
  .coming-soon-content h2 {
    font-size: 28px;
  }
}
</style>