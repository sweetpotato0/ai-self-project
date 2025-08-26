<template>
  <div class="loading-container" :class="{ 'full-screen': fullScreen }">
    <div class="loading-content">
      <el-icon class="loading-icon" :size="size">
        <Loading />
      </el-icon>
      <p v-if="text" class="loading-text">{{ text }}</p>
    </div>
  </div>
</template>

<script setup>
import { Loading } from '@element-plus/icons-vue'

defineProps({
  text: {
    type: String,
    default: '加载中...'
  },
  size: {
    type: Number,
    default: 40
  },
  fullScreen: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.loading-container.full-screen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  z-index: 9999;
}

.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.loading-icon {
  color: #409eff;
  animation: spin 1s linear infinite;
}

.loading-text {
  color: #606266;
  font-size: 14px;
  margin: 0;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>