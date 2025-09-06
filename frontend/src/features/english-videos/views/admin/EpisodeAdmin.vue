<template>
  <div class="episode-admin">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="$router.go(-1)" :icon="ArrowLeft">返回</el-button>
        <h1>剧集管理</h1>
      </div>
    </div>

    <!-- 剧集管理组件 -->
    <EpisodeManagement />
  </div>
</template>

<script setup>
import { ArrowLeft } from '@element-plus/icons-vue'
import EpisodeManagement from '../../components/admin/EpisodeManagement.vue'
</script>

<style scoped>
.episode-admin {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.header-left h1 {
  margin: 0;
  color: #303133;
  font-size: 24px;
  font-weight: 600;
}
</style>