<template>
  <div class="breadcrumb">
    <el-breadcrumb separator="/">
      <el-breadcrumb-item @click="navigateBack">工具箱</el-breadcrumb-item>
      <el-breadcrumb-item>{{ categoryName }}</el-breadcrumb-item>
    </el-breadcrumb>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

defineProps({
  categoryName: {
    type: String,
    required: true
  }
})

const router = useRouter()

const navigateBack = () => {
  router.push({ name: 'tools' })
}
</script>

<style scoped>
.breadcrumb {
  margin-bottom: 16px;
}

.breadcrumb .el-breadcrumb-item {
  cursor: pointer;
}
</style>