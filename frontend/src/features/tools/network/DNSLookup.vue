<template>
  <div class="dns-lookup-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="$router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="$router.push({ name: 'tools-operations' })">运维类</el-breadcrumb-item>
          <el-breadcrumb-item>DNS查询</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>DNS查询工具</h1>
      <p>查询域名的DNS记录信息，支持A、AAAA、CNAME、MX、NS、TXT、PTR等记录类型</p>
    </div>

    <div class="tools-content">
        <!-- 查询表单 -->
        <el-card class="query-form-card">
          <el-form @submit.prevent="performLookup" class="lookup-form">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="域名" required>
                  <el-input 
                    v-model="lookupForm.domain"
                    placeholder="请输入域名，如 example.com"
                    size="large"
                    :disabled="loading"
                    @keyup.enter="performLookup"
                  >
                    <template #prefix>
                      <el-icon><Position /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="记录类型" required>
                  <el-select 
                    v-model="lookupForm.type"
                    placeholder="选择记录类型"
                    size="large"
                    :disabled="loading"
                  >
                    <el-option label="A记录 (IPv4)" value="A" />
                    <el-option label="AAAA记录 (IPv6)" value="AAAA" />
                    <el-option label="CNAME记录" value="CNAME" />
                    <el-option label="MX记录" value="MX" />
                    <el-option label="NS记录" value="NS" />
                    <el-option label="TXT记录" value="TXT" />
                    <el-option label="PTR记录" value="PTR" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label=" ">
                  <el-button 
                    type="primary"
                    size="large"
                    :loading="loading"
                    @click="performLookup"
                    style="width: 100%"
                  >
                    <el-icon v-if="!loading"><Search /></el-icon>
                    {{ loading ? '查询中...' : '查询' }}
                  </el-button>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-card>

        <!-- 快速域名选择 -->
        <el-card class="quick-domains-card" v-if="!result">
          <template #header>
            <div class="card-header">
              <el-icon><Star /></el-icon>
              <span>快速测试</span>
            </div>
          </template>
          <div class="quick-domains">
            <el-tag 
              v-for="domain in quickDomains"
              :key="domain"
              @click="setDomain(domain)"
              class="domain-tag"
              effect="plain"
            >
              {{ domain }}
            </el-tag>
          </div>
        </el-card>

        <!-- 查询结果 -->
        <el-card v-if="result" class="result-card">
          <template #header>
            <div class="result-header">
              <div class="result-title">
                <el-icon><DocumentCopy /></el-icon>
                <span>查询结果</span>
              </div>
              <div class="result-actions">
                <el-button size="small" @click="copyResult">
                  <el-icon><CopyDocument /></el-icon>
                  复制结果
                </el-button>
                <el-button size="small" @click="clearResult">清空</el-button>
              </div>
            </div>
          </template>

          <!-- 查询信息 -->
          <div class="query-info">
            <el-descriptions :column="3" size="default" border>
              <el-descriptions-item label="查询域名">{{ result.domain }}</el-descriptions-item>
              <el-descriptions-item label="记录类型">{{ result.type }}</el-descriptions-item>
              <el-descriptions-item label="查询状态">
                <el-tag :type="result.status === 'success' ? 'success' : 'danger'" size="small">
                  {{ result.status === 'success' ? '成功' : '失败' }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </div>

          <!-- 错误信息 -->
          <el-alert
            v-if="result.status === 'error'"
            :title="result.message"
            type="error"
            show-icon
            :closable="false"
            style="margin-top: 16px;"
          />

          <!-- DNS记录列表 -->
          <div v-if="result.status === 'success' && result.records.length > 0" class="dns-records">
            <h4>DNS记录 ({{ result.records.length }}条)</h4>
            <el-table :data="result.records" style="width: 100%" border>
              <el-table-column prop="type" label="类型" width="80" align="center">
                <template #default="{ row }">
                  <el-tag :type="getRecordTypeColor(row.type)" size="small">{{ row.type }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="value" label="记录值" min-width="300">
                <template #default="{ row }">
                  <div class="record-value">
                    <code>{{ row.value }}</code>
                    <el-button 
                      size="small" 
                      text 
                      @click="copyText(row.value)"
                      style="margin-left: 8px;"
                    >
                      <el-icon><CopyDocument /></el-icon>
                    </el-button>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="ttl" label="TTL" width="80" align="center">
                <template #default="{ row }">
                  <span v-if="row.ttl">{{ row.ttl }}s</span>
                  <span v-else>-</span>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <!-- 无记录提示 -->
          <el-empty 
            v-if="result.status === 'success' && result.records.length === 0"
            description="未找到相关DNS记录"
            :image-size="120"
          />
        </el-card>

        <!-- 使用说明 -->
        <el-card class="usage-card">
          <template #header>
            <div class="card-header">
              <el-icon><QuestionFilled /></el-icon>
              <span>DNS记录类型说明</span>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="A记录">将域名指向IPv4地址</el-descriptions-item>
            <el-descriptions-item label="AAAA记录">将域名指向IPv6地址</el-descriptions-item>
            <el-descriptions-item label="CNAME记录">将域名指向另一个域名</el-descriptions-item>
            <el-descriptions-item label="MX记录">指定邮件服务器</el-descriptions-item>
            <el-descriptions-item label="NS记录">指定权威域名服务器</el-descriptions-item>
            <el-descriptions-item label="TXT记录">存储文本信息，常用于验证</el-descriptions-item>
            <el-descriptions-item label="PTR记录">反向DNS查询，IP地址到域名</el-descriptions-item>
          </el-descriptions>
        </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Position, 
  Search, 
  DocumentCopy,
  CopyDocument,
  QuestionFilled,
  Star
} from '@element-plus/icons-vue'
import { networkApi } from '@/features/tools/network/api'

// 响应式数据
const loading = ref(false)
const result = ref(null)

const lookupForm = reactive({
  domain: '',
  type: 'A'
})

const quickDomains = [
  'google.com',
  'github.com', 
  'cloudflare.com',
  'baidu.com',
  'qq.com'
]

// DNS查询
const performLookup = async () => {
  if (!lookupForm.domain.trim()) {
    ElMessage.warning('请输入要查询的域名')
    return
  }

  loading.value = true
  try {
    const response = await networkApi.dnsLookup({
      domain: lookupForm.domain.trim(),
      type: lookupForm.type
    })
    
    result.value = response.data
    ElMessage.success('DNS查询完成')
  } catch (error) {
    console.error('DNS查询失败:', error)
    ElMessage.error(error.response?.data?.message || 'DNS查询失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 设置域名
const setDomain = (domain) => {
  lookupForm.domain = domain
}

// 清空结果
const clearResult = () => {
  result.value = null
}

// 复制结果
const copyResult = () => {
  if (!result.value) return

  const resultText = [
    `域名: ${result.value.domain}`,
    `记录类型: ${result.value.type}`,
    `查询状态: ${result.value.status}`,
    result.value.message ? `消息: ${result.value.message}` : '',
    '',
    'DNS记录:',
    ...result.value.records.map(record => `${record.type}: ${record.value}${record.ttl ? ` (TTL: ${record.ttl}s)` : ''}`)
  ].filter(Boolean).join('\n')

  copyText(resultText)
}

// 复制文本
const copyText = (text) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 获取记录类型颜色
const getRecordTypeColor = (type) => {
  const colors = {
    'A': 'success',
    'AAAA': 'warning', 
    'CNAME': 'info',
    'MX': 'danger',
    'NS': 'primary',
    'TXT': 'success',
    'PTR': 'warning'
  }
  return colors[type] || 'info'
}
</script>

<style scoped>
.dns-lookup-container {
  width: 100%;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.tools-header {
  margin-bottom: 24px;
  text-align: center;
  width: 100%;
  max-width: 1200px;
}

.breadcrumb {
  margin-bottom: 16px;
  text-align: left;
}

.breadcrumb :deep(.el-breadcrumb__item a) {
  color: var(--el-color-primary);
  text-decoration: none;
  cursor: pointer;
}

.breadcrumb :deep(.el-breadcrumb__item a:hover) {
  color: var(--el-color-primary-light-3);
}

.tools-header h1 {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
}

.tools-header p {
  font-size: 16px;
  color: var(--el-text-color-regular);
  margin: 0;
  line-height: 1.5;
}

.tools-content {
  width: 100%;
  max-width: 1200px;
}

.query-form-card {
  margin-bottom: 24px;
}

.lookup-form {
  margin: 0;
}

.quick-domains-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.quick-domains {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.domain-tag {
  cursor: pointer;
  transition: all 0.2s ease;
}

.domain-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.result-card {
  margin-bottom: 24px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.result-actions {
  display: flex;
  gap: 8px;
}

.query-info {
  margin-bottom: 16px;
}

.dns-records {
  margin-top: 16px;
}

.dns-records h4 {
  margin: 0 0 16px 0;
  color: var(--el-text-color-primary);
}

.record-value {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.record-value code {
  flex: 1;
  padding: 4px 8px;
  background: var(--el-bg-color-page);
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  word-break: break-all;
}

.usage-card {
  margin-bottom: 24px;
}

.usage-card :deep(.el-descriptions-item__content) {
  color: var(--el-text-color-regular);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .result-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .result-actions {
    justify-content: center;
  }
  
  .quick-domains {
    justify-content: center;
  }
}
</style>