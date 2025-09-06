<template>
  <div class="video-series-admin">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="$router.go(-1)" :icon="ArrowLeft">返回</el-button>
        <h1>视频系列管理</h1>
      </div>
      <div class="header-right">
        <el-button type="primary" :icon="Plus" @click="showCreateDialog">
          新建系列
        </el-button>
      </div>
    </div>

    <!-- 系列管理组件 -->
    <VideoSeriesManagement />
  </div>
</template>

<script setup>
import { ArrowLeft, Plus } from '@element-plus/icons-vue'
import VideoSeriesManagement from '../../components/admin/VideoSeriesManagement.vue'
</script>

<style scoped>
.video-series-admin {
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