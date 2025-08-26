<template>
  <div class="tools-page-container">
    <div class="tools-header">
      <ToolsBreadcrumb :category-name="categoryName" />
      <h1>{{ title }}</h1>
      <p>{{ description }}</p>
    </div>

    <div class="tools-grid">
      <ToolCard 
        v-for="tool in tools" 
        :key="tool.id"
        :tool="tool"
        @click="handleToolClick"
      />
    </div>
  </div>
</template>

<script setup>
import ToolsBreadcrumb from './ToolsBreadcrumb.vue'
import ToolCard from './ToolCard.vue'

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
  tools: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['tool-click'])

const handleToolClick = (tool) => {
  emit('tool-click', tool)
}
</script>

<style scoped>
.tools-page-container {
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

@media (max-width: 768px) {
  .tools-page-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .tools-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>