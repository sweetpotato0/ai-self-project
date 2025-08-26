<template>
  <div class="domain-info-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-query' })">查询类</el-breadcrumb-item>
          <el-breadcrumb-item>域名信息查询</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>域名信息查询</h1>
      <p>查询域名的SSL证书、网站信息、服务器配置等详细信息</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Connection /></el-icon>
              <span>域名输入</span>
            </div>
          </template>
          <div class="input-content">
            <div class="domain-input-group">
              <el-input
                v-model="domainName"
                placeholder="请输入域名，如：google.com"
                size="large"
                clearable
                @keyup.enter="queryDomainInfo"
              >
                <template #prepend>
                  <el-select v-model="selectedExample" placeholder="示例" style="width: 120px" @change="loadExample">
                    <el-option label="Google" value="google.com" />
                    <el-option label="GitHub" value="github.com" />
                    <el-option label="百度" value="baidu.com" />
                    <el-option label="淘宝" value="taobao.com" />
                  </el-select>
                </template>
                <template #append>
                  <el-checkbox v-model="includeSubdomains">包含www</el-checkbox>
                </template>
              </el-input>
            </div>
            <div class="action-buttons">
              <el-button
                @click="queryDomainInfo"
                type="primary"
                size="large"
                :loading="loading"
                :disabled="!domainName.trim()"
              >
                <el-icon><Search /></el-icon>
                查询域名信息
              </el-button>
              <el-button @click="clearResults" size="large">
                <el-icon><RefreshLeft /></el-icon>
                清空结果
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="domainResult || error" class="result-section">
        <el-card class="result-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>域名信息查询结果</span>
            </div>
          </template>
          
          <div v-if="error" class="error-content">
            <el-alert
              :title="error"
              type="error"
              show-icon
              :closable="false"
            />
          </div>

          <div v-else-if="domainResult" class="domain-info-content">
            <!-- 基本信息 -->
            <div class="info-section">
              <h3><el-icon><Monitor /></el-icon>基本信息</h3>
              <div class="info-grid">
                <div class="info-item">
                  <div class="info-label">域名</div>
                  <div class="info-value primary">{{ domainResult.domain }}</div>
                </div>
                <div class="info-item">
                  <div class="info-label">IP地址</div>
                  <div class="info-value">{{ domainResult.ipAddress || '-' }}</div>
                </div>
                <div class="info-item">
                  <div class="info-label">服务器类型</div>
                  <div class="info-value">{{ domainResult.serverType || '-' }}</div>
                </div>
                <div class="info-item">
                  <div class="info-label">响应时间</div>
                  <div class="info-value">
                    {{ domainResult.responseTime ? `${domainResult.responseTime}ms` : '-' }}
                  </div>
                </div>
                <div class="info-item">
                  <div class="info-label">HTTP状态</div>
                  <div class="info-value" :class="getStatusClass(domainResult.httpStatus)">
                    {{ domainResult.httpStatus || '-' }}
                  </div>
                </div>
                <div class="info-item">
                  <div class="info-label">域名年龄</div>
                  <div class="info-value">{{ domainResult.domainAge || '-' }}</div>
                </div>
              </div>
            </div>

            <!-- SSL证书信息 -->
            <div v-if="domainResult.sslInfo" class="info-section">
              <h3><el-icon><Lock /></el-icon>SSL证书信息</h3>
              <div class="ssl-info">
                <div class="ssl-status">
                  <el-tag 
                    :type="domainResult.sslInfo.valid ? 'success' : 'danger'" 
                    size="large"
                  >
                    <el-icon>
                      <component :is="domainResult.sslInfo.valid ? 'Check' : 'Close'" />
                    </el-icon>
                    {{ domainResult.sslInfo.valid ? 'SSL证书有效' : 'SSL证书无效' }}
                  </el-tag>
                </div>
                <div class="ssl-details">
                  <div class="ssl-grid">
                    <div class="ssl-item">
                      <div class="ssl-label">颁发给</div>
                      <div class="ssl-value">{{ domainResult.sslInfo.subject || '-' }}</div>
                    </div>
                    <div class="ssl-item">
                      <div class="ssl-label">颁发机构</div>
                      <div class="ssl-value">{{ domainResult.sslInfo.issuer || '-' }}</div>
                    </div>
                    <div class="ssl-item">
                      <div class="ssl-label">有效期从</div>
                      <div class="ssl-value">{{ formatDate(domainResult.sslInfo.validFrom) }}</div>
                    </div>
                    <div class="ssl-item">
                      <div class="ssl-label">有效期到</div>
                      <div class="ssl-value" :class="getSSLExpiryClass()">
                        {{ formatDate(domainResult.sslInfo.validTo) }}
                        <span v-if="getSSLDaysLeft() !== null" class="expiry-info">
                          ({{ getSSLDaysLeft() > 0 ? '还有' : '已过期' }}{{ Math.abs(getSSLDaysLeft()) }}天)
                        </span>
                      </div>
                    </div>
                    <div class="ssl-item">
                      <div class="ssl-label">证书类型</div>
                      <div class="ssl-value">{{ domainResult.sslInfo.type || '-' }}</div>
                    </div>
                    <div class="ssl-item">
                      <div class="ssl-label">加密强度</div>
                      <div class="ssl-value">{{ domainResult.sslInfo.keySize || '-' }}</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 网站信息 -->
            <div v-if="domainResult.websiteInfo" class="info-section">
              <h3><el-icon><Document /></el-icon>网站信息</h3>
              <div class="website-info">
                <div class="website-item">
                  <div class="website-label">网站标题</div>
                  <div class="website-value">{{ domainResult.websiteInfo.title || '-' }}</div>
                </div>
                <div class="website-item">
                  <div class="website-label">网站描述</div>
                  <div class="website-value">{{ domainResult.websiteInfo.description || '-' }}</div>
                </div>
                <div class="website-item">
                  <div class="website-label">关键词</div>
                  <div class="website-value">
                    <div v-if="domainResult.websiteInfo.keywords" class="keywords">
                      <el-tag
                        v-for="keyword in domainResult.websiteInfo.keywords.split(',').slice(0, 5)"
                        :key="keyword"
                        size="small"
                      >
                        {{ keyword.trim() }}
                      </el-tag>
                    </div>
                    <span v-else>-</span>
                  </div>
                </div>
                <div class="website-item">
                  <div class="website-label">语言</div>
                  <div class="website-value">{{ domainResult.websiteInfo.language || '-' }}</div>
                </div>
              </div>
            </div>

            <!-- 技术信息 -->
            <div v-if="domainResult.techInfo" class="info-section">
              <h3><el-icon><Tools /></el-icon>技术信息</h3>
              <div class="tech-info">
                <div class="tech-grid">
                  <div class="tech-item">
                    <div class="tech-label">Content-Type</div>
                    <div class="tech-value">{{ domainResult.techInfo.contentType || '-' }}</div>
                  </div>
                  <div class="tech-item">
                    <div class="tech-label">字符编码</div>
                    <div class="tech-value">{{ domainResult.techInfo.charset || '-' }}</div>
                  </div>
                  <div class="tech-item">
                    <div class="tech-label">压缩方式</div>
                    <div class="tech-value">{{ domainResult.techInfo.compression || '-' }}</div>
                  </div>
                  <div class="tech-item">
                    <div class="tech-label">页面大小</div>
                    <div class="tech-value">{{ domainResult.techInfo.contentLength || '-' }}</div>
                  </div>
                </div>
                <div v-if="domainResult.techInfo.technologies" class="technologies">
                  <div class="tech-label">检测到的技术：</div>
                  <div class="tech-tags">
                    <el-tag
                      v-for="tech in domainResult.techInfo.technologies"
                      :key="tech"
                      type="info"
                      size="small"
                    >
                      {{ tech }}
                    </el-tag>
                  </div>
                </div>
              </div>
            </div>

            <!-- DNS记录 -->
            <div v-if="domainResult.dnsRecords && domainResult.dnsRecords.length" class="info-section">
              <h3><el-icon><Connection /></el-icon>DNS记录</h3>
              <div class="dns-records">
                <div
                  v-for="record in domainResult.dnsRecords"
                  :key="record.type + record.value"
                  class="dns-record"
                >
                  <el-tag :type="getDNSRecordType(record.type)" size="small">{{ record.type }}</el-tag>
                  <span class="dns-value">{{ record.value }}</span>
                  <span v-if="record.ttl" class="dns-ttl">TTL: {{ record.ttl }}</span>
                </div>
              </div>
            </div>

            <div class="action-section">
              <el-button @click="copyToClipboard(formatResult())" type="success">
                <el-icon><DocumentCopy /></el-icon>
                复制结果
              </el-button>
              <el-button @click="addToHistory" type="primary">
                <el-icon><StarFilled /></el-icon>
                保存到历史
              </el-button>
              <el-button @click="openDomain" type="info">
                <el-icon><View /></el-icon>
                访问网站
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="history.length > 0" class="history-section">
        <el-card class="history-card">
          <template #header>
            <div class="card-header">
              <el-icon><Clock /></el-icon>
              <span>查询历史</span>
              <div class="card-actions">
                <el-button @click="clearHistory" size="small" text type="danger">
                  <el-icon><Delete /></el-icon>
                  清空历史
                </el-button>
              </div>
            </div>
          </template>
          <div class="history-list">
            <div 
              v-for="(item, index) in history" 
              :key="index"
              class="history-item"
              @click="loadFromHistory(item)"
            >
              <div class="history-domain">{{ item.domain }}</div>
              <div class="history-status">
                <el-tag :type="getStatusClass(item.httpStatus)" size="small">
                  {{ item.httpStatus || '未知' }}
                </el-tag>
              </div>
              <div class="history-ssl">
                <el-tag 
                  :type="item.sslInfo?.valid ? 'success' : 'danger'" 
                  size="small"
                >
                  {{ item.sslInfo?.valid ? 'SSL有效' : 'SSL无效' }}
                </el-tag>
              </div>
              <div class="history-time">{{ new Date(item.queryTime).toLocaleString() }}</div>
              <el-button 
                @click.stop="removeFromHistory(index)" 
                size="small" 
                text 
                type="danger"
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Connection,
  Search,
  RefreshLeft,
  InfoFilled,
  Monitor,
  Lock,
  Document,
  Tools,
  DocumentCopy,
  StarFilled,
  View,
  Clock,
  Delete,
  Close,
  Check
} from '@element-plus/icons-vue'

const router = useRouter()

const domainName = ref('')
const selectedExample = ref('')
const includeSubdomains = ref(false)
const loading = ref(false)
const domainResult = ref(null)
const error = ref('')
const history = ref([])

onMounted(() => {
  loadHistory()
})

const loadExample = () => {
  domainName.value = selectedExample.value
}

const queryDomainInfo = async () => {
  if (!domainName.value.trim()) {
    ElMessage.warning('请输入域名')
    return
  }

  let domain = domainName.value.trim().toLowerCase()
  if (!isValidDomain(domain)) {
    ElMessage.error('请输入有效的域名格式')
    return
  }

  if (includeSubdomains.value && !domain.startsWith('www.')) {
    domain = 'www.' + domain
  }

  loading.value = true
  error.value = ''
  domainResult.value = null

  try {
    const startTime = Date.now()
    const result = await simulateDomainInfo(domain)
    const endTime = Date.now()
    
    domainResult.value = {
      ...result,
      responseTime: endTime - startTime,
      queryTime: Date.now()
    }
    ElMessage.success('域名信息查询成功')
  } catch (err) {
    error.value = '域名信息查询失败：' + err.message
    ElMessage.error(error.value)
  } finally {
    loading.value = false
  }
}

// 模拟域名信息查询（实际项目中应使用真实API）
const simulateDomainInfo = async (domain) => {
  // 模拟网络延迟
  await new Promise(resolve => setTimeout(resolve, 2000))
  
  const mockData = {
    'google.com': {
      domain: 'google.com',
      ipAddress: '172.217.160.14',
      serverType: 'gws',
      httpStatus: '200 OK',
      domainAge: '26年',
      sslInfo: {
        valid: true,
        subject: '*.google.com',
        issuer: 'GTS CA 1C3',
        validFrom: '2023-10-09T08:21:04Z',
        validTo: '2024-01-01T08:21:03Z',
        type: 'Domain Validation (DV)',
        keySize: '2048-bit RSA'
      },
      websiteInfo: {
        title: 'Google',
        description: 'Search the world\'s information',
        keywords: 'search,google,engine',
        language: 'en'
      },
      techInfo: {
        contentType: 'text/html',
        charset: 'UTF-8',
        compression: 'gzip',
        contentLength: '47KB',
        technologies: ['HTTP/2', 'Google Analytics', 'Google Tag Manager']
      },
      dnsRecords: [
        { type: 'A', value: '172.217.160.14', ttl: 300 },
        { type: 'AAAA', value: '2404:6800:4008:c06::8e', ttl: 300 },
        { type: 'MX', value: '10 smtp.google.com', ttl: 3600 }
      ]
    }
  }

  if (mockData[domain]) {
    return mockData[domain]
  }

  // 为其他域名生成模拟数据
  const httpStatuses = ['200 OK', '301 Moved Permanently', '404 Not Found', '503 Service Unavailable']
  const serverTypes = ['nginx', 'Apache', 'cloudflare', 'Microsoft-IIS']
  const randomStatus = httpStatuses[Math.floor(Math.random() * httpStatuses.length)]
  const randomServer = serverTypes[Math.floor(Math.random() * serverTypes.length)]

  return {
    domain,
    ipAddress: `${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`,
    serverType: randomServer,
    httpStatus: randomStatus,
    domainAge: `${Math.floor(Math.random() * 20) + 1}年`,
    sslInfo: {
      valid: Math.random() > 0.3,
      subject: domain,
      issuer: 'Let\'s Encrypt Authority X3',
      validFrom: new Date(Date.now() - Math.random() * 365 * 24 * 60 * 60 * 1000).toISOString(),
      validTo: new Date(Date.now() + Math.random() * 365 * 24 * 60 * 60 * 1000).toISOString(),
      type: 'Domain Validation (DV)',
      keySize: '2048-bit RSA'
    },
    websiteInfo: {
      title: domain.charAt(0).toUpperCase() + domain.slice(1),
      description: `Welcome to ${domain}`,
      keywords: `${domain},website,online`,
      language: 'zh-CN'
    },
    techInfo: {
      contentType: 'text/html',
      charset: 'UTF-8',
      compression: 'gzip',
      contentLength: `${Math.floor(Math.random() * 500) + 50}KB`,
      technologies: ['HTML5', 'CSS3', 'JavaScript']
    },
    dnsRecords: [
      { type: 'A', value: `${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`, ttl: 300 }
    ]
  }
}

const isValidDomain = (domain) => {
  const domainRegex = /^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/i
  return domainRegex.test(domain)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return dateString
  }
}

const getSSLDaysLeft = () => {
  if (!domainResult.value?.sslInfo?.validTo) return null
  const expiry = new Date(domainResult.value.sslInfo.validTo)
  const now = new Date()
  return Math.ceil((expiry - now) / (1000 * 60 * 60 * 24))
}

const getSSLExpiryClass = () => {
  const days = getSSLDaysLeft()
  if (days === null) return ''
  if (days < 0) return 'expired'
  if (days < 7) return 'warning'
  if (days < 30) return 'caution'
  return ''
}

const getStatusClass = (status) => {
  if (!status) return 'info'
  if (status.startsWith('2')) return 'success'
  if (status.startsWith('3')) return 'warning'
  if (status.startsWith('4') || status.startsWith('5')) return 'danger'
  return 'info'
}

const getDNSRecordType = (type) => {
  const typeMap = {
    'A': 'primary',
    'AAAA': 'success',
    'CNAME': 'info',
    'MX': 'warning',
    'TXT': 'default',
    'NS': 'danger'
  }
  return typeMap[type] || 'default'
}

const formatResult = () => {
  if (!domainResult.value) return ''
  
  let result = `域名信息查询结果：\n`
  result += `域名：${domainResult.value.domain}\n`
  result += `IP地址：${domainResult.value.ipAddress || '-'}\n`
  result += `HTTP状态：${domainResult.value.httpStatus || '-'}\n`
  result += `服务器类型：${domainResult.value.serverType || '-'}\n`
  result += `响应时间：${domainResult.value.responseTime || '-'}ms\n`
  result += `域名年龄：${domainResult.value.domainAge || '-'}\n`
  
  if (domainResult.value.sslInfo) {
    result += `\nSSL证书信息：\n`
    result += `证书状态：${domainResult.value.sslInfo.valid ? '有效' : '无效'}\n`
    result += `颁发给：${domainResult.value.sslInfo.subject || '-'}\n`
    result += `颁发机构：${domainResult.value.sslInfo.issuer || '-'}\n`
    result += `有效期：${formatDate(domainResult.value.sslInfo.validFrom)} - ${formatDate(domainResult.value.sslInfo.validTo)}\n`
  }
  
  result += `\n查询时间：${new Date(domainResult.value.queryTime).toLocaleString()}`
  
  return result
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('结果已复制到剪贴板')
  } catch (err) {
    ElMessage.error('复制失败：' + err.message)
  }
}

const addToHistory = () => {
  if (!domainResult.value) return
  
  const existingIndex = history.value.findIndex(item => item.domain === domainResult.value.domain)
  if (existingIndex >= 0) {
    history.value.splice(existingIndex, 1)
  }
  
  history.value.unshift({ ...domainResult.value })
  
  if (history.value.length > 20) {
    history.value = history.value.slice(0, 20)
  }
  
  saveHistory()
  ElMessage.success('已保存到查询历史')
}

const loadFromHistory = (item) => {
  domainName.value = item.domain
  domainResult.value = { ...item }
  error.value = ''
}

const removeFromHistory = (index) => {
  history.value.splice(index, 1)
  saveHistory()
  ElMessage.success('已从历史记录中移除')
}

const clearHistory = () => {
  history.value = []
  saveHistory()
  ElMessage.success('已清空查询历史')
}

const clearResults = () => {
  domainName.value = ''
  selectedExample.value = ''
  domainResult.value = null
  error.value = ''
}

const openDomain = () => {
  if (domainResult.value?.domain) {
    window.open(`https://${domainResult.value.domain}`, '_blank')
  }
}

const saveHistory = () => {
  localStorage.setItem('domain-info-history', JSON.stringify(history.value))
}

const loadHistory = () => {
  try {
    const saved = localStorage.getItem('domain-info-history')
    if (saved) {
      history.value = JSON.parse(saved)
    }
  } catch (err) {
    console.error('加载历史记录失败：', err)
  }
}
</script>

<style scoped>
.domain-info-container {
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

.tools-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.card-actions {
  margin-left: auto;
}

.input-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.domain-input-group {
  width: 100%;
}

.action-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.error-content {
  padding: 20px 0;
}

.domain-info-content {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.info-section h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 8px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.info-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #1e293b;
  font-weight: 600;
  word-break: break-all;
}

.info-value.primary {
  color: #3b82f6;
  font-size: 16px;
}

.info-value.success {
  color: #059669;
}

.info-value.warning {
  color: #ea580c;
}

.info-value.danger {
  color: #dc2626;
}

.info-value.expired {
  color: #dc2626;
}

.info-value.caution {
  color: #ca8a04;
}

.expiry-info {
  font-size: 12px;
  font-weight: 400;
  margin-left: 8px;
}

.ssl-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.ssl-status {
  text-align: center;
}

.ssl-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.ssl-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ssl-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

.ssl-value {
  font-size: 14px;
  color: #1e293b;
  font-weight: 600;
}

.website-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.website-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background: #f1f5f9;
  border-radius: 6px;
}

.website-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

.website-value {
  font-size: 14px;
  color: #1e293b;
}

.keywords {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tech-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.tech-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.tech-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tech-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

.tech-value {
  font-size: 14px;
  color: #1e293b;
  font-weight: 600;
}

.technologies {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tech-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.dns-records {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dns-record {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 6px;
  font-family: monospace;
}

.dns-value {
  flex: 1;
  font-size: 14px;
  color: #1e293b;
}

.dns-ttl {
  font-size: 12px;
  color: #64748b;
}

.action-section {
  display: flex;
  gap: 12px;
  justify-content: center;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 400px;
  overflow-y: auto;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.history-item:hover {
  background: #e2e8f0;
}

.history-domain {
  font-weight: 600;
  color: #3b82f6;
  min-width: 150px;
}

.history-status,
.history-ssl {
  min-width: 80px;
}

.history-time {
  font-size: 12px;
  color: #9ca3af;
  min-width: 120px;
  margin-left: auto;
}

@media (max-width: 768px) {
  .domain-info-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .ssl-grid,
  .tech-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .action-section {
    flex-direction: column;
  }
  
  .history-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .history-domain,
  .history-status,
  .history-ssl,
  .history-time {
    min-width: unset;
  }
  
  .history-time {
    margin-left: 0;
  }
}
</style>