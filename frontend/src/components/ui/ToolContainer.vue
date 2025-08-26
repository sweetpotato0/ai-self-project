<template>
  <div class="tool-container">
    <div class="tool-header">
      <ToolsBreadcrumb :category-name="categoryName" />
      <div class="header-content">
        <div class="header-icon">
          <el-icon :size="48">
            <component :is="icon" />
          </el-icon>
        </div>
        <div class="header-text">
          <h1>{{ title }}</h1>
          <p>{{ description }}</p>
        </div>
      </div>
    </div>

    <div class="tool-content">
      <slot />
    </div>
  </div>
</template>

<script setup>
import ToolsBreadcrumb from './ToolsBreadcrumb.vue'

defineProps({
  categoryName: {
    type: String,
    required: true
  },
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  icon: {
    type: [String, Object],
    required: true
  }
})
</script>

<style scoped>
.tool-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.tool-header {
  margin-bottom: 30px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  color: white;
  text-align: center;
  margin-top: 20px;
}

.header-icon {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 20px;
  backdrop-filter: blur(10px);
}

.header-text h1 {
  font-size: 32px;
  font-weight: 600;
  margin: 0 0 8px 0;
}

.header-text p {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
}

.tool-content {
  background: white;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

@media (max-width: 768px) {
  .tool-container {
    padding: 16px;
  }
  
  .header-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-text h1 {
    font-size: 24px;
  }
  
  .tool-content {
    padding: 24px 16px;
  }
}
</style>