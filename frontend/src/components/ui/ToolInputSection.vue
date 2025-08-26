<template>
  <div class="input-section">
    <div class="input-header" v-if="title">
      <h3>{{ title }}</h3>
      <p v-if="description">{{ description }}</p>
    </div>
    
    <div class="input-content">
      <slot />
    </div>
    
    <div class="input-actions" v-if="$slots.actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.input-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.input-header {
  margin-bottom: 20px;
  text-align: center;
}

.input-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.input-header p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

.input-content {
  margin-bottom: 20px;
}

.input-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}

@media (max-width: 768px) {
  .input-section {
    padding: 16px;
  }
  
  .input-actions {
    flex-direction: column;
  }
  
  .input-actions :deep(.el-button) {
    width: 100%;
  }
}
</style>