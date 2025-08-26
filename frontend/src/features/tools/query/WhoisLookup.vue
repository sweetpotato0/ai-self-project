<template>
  <div class="whois-lookup-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-query' })">查询类</el-breadcrumb-item>
          <el-breadcrumb-item>Whois查询</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>Whois查询</h1>
      <p>查询域名的注册信息、所有者信息和DNS配置</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Search /></el-icon>
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
                @keyup.enter="lookupWhois"
              >
                <template #prepend>
                  <el-select v-model="selectedExample" placeholder="示例" style="width: 120px" @change="loadExample">
                    <el-option label="Google" value="google.com" />
                    <el-option label="GitHub" value="github.com" />
                    <el-option label="百度" value="baidu.com" />
                    <el-option label="阿里巴巴" value="alibaba.com" />
                  </el-select>
                </template>
              </el-input>
            </div>
            <div class="action-buttons">
              <el-button
                @click="lookupWhois"
                type="primary"
                size="large"
                :loading="loading"
                :disabled="!domainName.trim()"
              >
                <el-icon><Search /></el-icon>
                查询Whois信息
              </el-button>
              <el-button @click="clearResults" size="large">
                <el-icon><RefreshLeft /></el-icon>
                清空结果
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="whoisResult || error" class="result-section">
        <el-card class="result-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>Whois查询结果</span>
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

          <div v-else-if="whoisResult" class="whois-info-content">
            <!-- 域名基本信息 -->
            <div class="info-section">
              <h3><el-icon><Document /></el-icon>域名信息</h3>
              <div class="info-grid">
                <div class="info-item">
                  <div class="info-label">域名</div>
                  <div class="info-value primary">{{ whoisResult.domain }}</div>
                </div>
                <div class="info-item">
                  <div class="info-label">注册商</div>
                  <div class="info-value">{{ whoisResult.registrar || '-' }}</div>
                </div>
                <div class="info-item">
                  <div class="info-label">注册时间</div>
                  <div class="info-value">{{ formatDate(whoisResult.creationDate) }}</div>
                </div>
                <div class="info-item">
                  <div class="info-label">过期时间</div>
                  <div class="info-value" :class="getExpiryClass()">
                    {{ formatDate(whoisResult.expirationDate) }}
                    <span v-if="getDaysUntilExpiry() !== null" class="expiry-info">
                      ({{ getDaysUntilExpiry() > 0 ? '还有' : '已过期' }}{{ Math.abs(getDaysUntilExpiry()) }}天)
                    </span>
                  </div>
                </div>
                <div class="info-item">
                  <div class="info-label">更新时间</div>
                  <div class="info-value">{{ formatDate(whoisResult.updatedDate) }}</div>
                </div>
                <div class="info-item">
                  <div class="info-label">域名状态</div>
                  <div class="info-value">
                    <div v-if="whoisResult.status && whoisResult.status.length" class="status-list">
                      <el-tag
                        v-for="status in whoisResult.status.slice(0, 3)"
                        :key="status"
                        size="small"
                        :type="getStatusType(status)"
                      >
                        {{ status }}
                      </el-tag>
                      <el-tag v-if="whoisResult.status.length > 3" size="small">
                        +{{ whoisResult.status.length - 3 }}
                      </el-tag>
                    </div>
                    <span v-else>-</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- DNS信息 -->
            <div v-if="whoisResult.nameServers && whoisResult.nameServers.length" class="info-section">
              <h3><el-icon><Connection /></el-icon>DNS服务器</h3>
              <div class="dns-list">
                <div
                  v-for="ns in whoisResult.nameServers"
                  :key="ns"
                  class="dns-item"
                >
                  <el-icon><Monitor /></el-icon>
                  <span>{{ ns }}</span>
                </div>
              </div>
            </div>

            <!-- 联系信息 -->
            <div v-if="whoisResult.contacts && Object.keys(whoisResult.contacts).length" class="info-section">
              <h3><el-icon><UserFilled /></el-icon>联系信息</h3>
              <el-tabs>
                <el-tab-pane
                  v-for="(contact, type) in whoisResult.contacts"
                  :key="type"
                  :label="getContactTypeLabel(type)"
                  :name="type"
                >
                  <div class="contact-info">
                    <div v-if="contact.name" class="contact-item">
                      <span class="contact-label">姓名:</span>
                      <span class="contact-value">{{ contact.name }}</span>
                    </div>
                    <div v-if="contact.organization" class="contact-item">
                      <span class="contact-label">组织:</span>
                      <span class="contact-value">{{ contact.organization }}</span>
                    </div>
                    <div v-if="contact.email" class="contact-item">
                      <span class="contact-label">邮箱:</span>
                      <span class="contact-value">{{ contact.email }}</span>
                    </div>
                    <div v-if="contact.phone" class="contact-item">
                      <span class="contact-label">电话:</span>
                      <span class="contact-value">{{ contact.phone }}</span>
                    </div>
                    <div v-if="contact.country" class="contact-item">
                      <span class="contact-label">国家:</span>
                      <span class="contact-value">{{ contact.country }}</span>
                    </div>
                  </div>
                </el-tab-pane>
              </el-tabs>
            </div>

            <!-- 原始数据 -->
            <div class="info-section">
              <h3><el-icon><Document /></el-icon>原始Whois数据</h3>
              <el-input
                v-model="whoisResult.rawData"
                type="textarea"
                :rows="10"
                readonly
                class="raw-data-textarea"
              />
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
              <div class="history-registrar">{{ item.registrar || '-' }}</div>
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
  Search,
  RefreshLeft,
  InfoFilled,
  Document,
  Connection,
  Monitor,
  UserFilled,
  DocumentCopy,
  StarFilled,
  Clock,
  Delete,
  Close
} from '@element-plus/icons-vue'

const router = useRouter()

const domainName = ref('')
const selectedExample = ref('')
const loading = ref(false)
const whoisResult = ref(null)
const error = ref('')
const history = ref([])

onMounted(() => {
  loadHistory()
})

const loadExample = () => {
  domainName.value = selectedExample.value
}

const lookupWhois = async () => {
  if (!domainName.value.trim()) {
    ElMessage.warning('请输入域名')
    return
  }

  const domain = domainName.value.trim().toLowerCase()
  if (!isValidDomain(domain)) {
    ElMessage.error('请输入有效的域名格式')
    return
  }

  loading.value = true
  error.value = ''
  whoisResult.value = null

  try {
    // 模拟Whois查询结果（由于CORS限制，实际项目中需要后端代理）
    const mockResult = await simulateWhoisLookup(domain)
    whoisResult.value = {
      ...mockResult,
      queryTime: Date.now()
    }
    ElMessage.success('Whois查询成功')
  } catch (err) {
    error.value = 'Whois查询失败：' + err.message
    ElMessage.error(error.value)
  } finally {
    loading.value = false
  }
}

// 模拟Whois查询（实际项目中应使用真实API）
const simulateWhoisLookup = async (domain) => {
  // 模拟网络延迟
  await new Promise(resolve => setTimeout(resolve, 1500))
  
  const mockData = {
    'google.com': {
      domain: 'google.com',
      registrar: 'MarkMonitor Inc.',
      creationDate: '1997-09-15T04:00:00Z',
      expirationDate: '2028-09-14T04:00:00Z',
      updatedDate: '2019-09-09T15:39:04Z',
      status: ['clientDeleteProhibited', 'clientTransferProhibited', 'clientUpdateProhibited', 'serverDeleteProhibited'],
      nameServers: ['ns1.google.com', 'ns2.google.com', 'ns3.google.com', 'ns4.google.com'],
      contacts: {
        registrant: {
          name: 'Google Inc.',
          organization: 'Google Inc.',
          email: 'dns-admin@google.com',
          country: 'US'
        }
      },
      rawData: `Domain Name: GOOGLE.COM
Registry Domain ID: 2138514_DOMAIN_COM-VRSN
Registrar WHOIS Server: whois.markmonitor.com
Registrar URL: http://www.markmonitor.com
Updated Date: 2019-09-09T15:39:04Z
Creation Date: 1997-09-15T04:00:00Z
Registry Expiry Date: 2028-09-14T04:00:00Z
Registrar: MarkMonitor Inc.
Registrar IANA ID: 292
Registrar Abuse Contact Email: abusecomplaints@markmonitor.com
Registrar Abuse Contact Phone: +1.2086851750
Domain Status: clientDeleteProhibited
Domain Status: clientTransferProhibited
Domain Status: clientUpdateProhibited
Domain Status: serverDeleteProhibited
Name Server: NS1.GOOGLE.COM
Name Server: NS2.GOOGLE.COM
Name Server: NS3.GOOGLE.COM
Name Server: NS4.GOOGLE.COM`
    }
  }

  if (mockData[domain]) {
    return mockData[domain]
  }

  // 为其他域名生成模拟数据
  const now = new Date()
  const creationDate = new Date(now.getTime() - Math.random() * 10 * 365 * 24 * 60 * 60 * 1000)
  const expirationDate = new Date(now.getTime() + Math.random() * 5 * 365 * 24 * 60 * 60 * 1000)

  return {
    domain,
    registrar: ['GoDaddy.com LLC', 'Namecheap Inc.', 'Google Domains LLC', 'Amazon Registrar Inc.'][Math.floor(Math.random() * 4)],
    creationDate: creationDate.toISOString(),
    expirationDate: expirationDate.toISOString(),
    updatedDate: new Date(now.getTime() - Math.random() * 365 * 24 * 60 * 60 * 1000).toISOString(),
    status: ['clientTransferProhibited', 'clientUpdateProhibited'],
    nameServers: [`ns1.${domain}`, `ns2.${domain}`],
    contacts: {
      registrant: {
        name: 'Domain Owner',
        organization: 'Private Registration',
        email: 'admin@' + domain,
        country: 'US'
      }
    },
    rawData: `Domain Name: ${domain.toUpperCase()}
Registrar: Mock Registrar Inc.
Creation Date: ${creationDate.toISOString()}
Registry Expiry Date: ${expirationDate.toISOString()}
Domain Status: clientTransferProhibited
Name Server: ns1.${domain}
Name Server: ns2.${domain}`
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

const getDaysUntilExpiry = () => {
  if (!whoisResult.value?.expirationDate) return null
  const expiry = new Date(whoisResult.value.expirationDate)
  const now = new Date()
  return Math.ceil((expiry - now) / (1000 * 60 * 60 * 24))
}

const getExpiryClass = () => {
  const days = getDaysUntilExpiry()
  if (days === null) return ''
  if (days < 0) return 'expired'
  if (days < 30) return 'warning'
  if (days < 90) return 'caution'
  return ''
}

const getStatusType = (status) => {
  if (status.includes('Prohibited') || status.includes('prohibited')) return 'danger'
  if (status.includes('Hold') || status.includes('hold')) return 'warning'
  return 'info'
}

const getContactTypeLabel = (type) => {
  const labels = {
    registrant: '注册人',
    admin: '管理员',
    tech: '技术',
    billing: '账务'
  }
  return labels[type] || type
}

const formatResult = () => {
  if (!whoisResult.value) return ''
  
  let result = `域名Whois查询结果：\n`
  result += `域名：${whoisResult.value.domain}\n`
  result += `注册商：${whoisResult.value.registrar || '-'}\n`
  result += `注册时间：${formatDate(whoisResult.value.creationDate)}\n`
  result += `过期时间：${formatDate(whoisResult.value.expirationDate)}\n`
  result += `更新时间：${formatDate(whoisResult.value.updatedDate)}\n`
  
  if (whoisResult.value.nameServers && whoisResult.value.nameServers.length) {
    result += `DNS服务器：${whoisResult.value.nameServers.join(', ')}\n`
  }
  
  result += `查询时间：${new Date(whoisResult.value.queryTime).toLocaleString()}`
  
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
  if (!whoisResult.value) return
  
  const existingIndex = history.value.findIndex(item => item.domain === whoisResult.value.domain)
  if (existingIndex >= 0) {
    history.value.splice(existingIndex, 1)
  }
  
  history.value.unshift({ ...whoisResult.value })
  
  if (history.value.length > 20) {
    history.value = history.value.slice(0, 20)
  }
  
  saveHistory()
  ElMessage.success('已保存到查询历史')
}

const loadFromHistory = (item) => {
  domainName.value = item.domain
  whoisResult.value = { ...item }
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
  whoisResult.value = null
  error.value = ''
}

const saveHistory = () => {
  localStorage.setItem('whois-lookup-history', JSON.stringify(history.value))
}

const loadHistory = () => {
  try {
    const saved = localStorage.getItem('whois-lookup-history')
    if (saved) {
      history.value = JSON.parse(saved)
    }
  } catch (err) {
    console.error('加载历史记录失败：', err)
  }
}
</script>

<style scoped>
.whois-lookup-container {
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

.whois-info-content {
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

.info-value.expired {
  color: #dc2626;
}

.info-value.warning {
  color: #ea580c;
}

.info-value.caution {
  color: #ca8a04;
}

.expiry-info {
  font-size: 12px;
  font-weight: 400;
  margin-left: 8px;
}

.status-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.dns-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 12px;
}

.dns-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #f1f5f9;
  border-radius: 6px;
  font-family: monospace;
  font-size: 14px;
}

.contact-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.contact-label {
  font-weight: 500;
  color: #64748b;
  min-width: 60px;
}

.contact-value {
  color: #1e293b;
}

.raw-data-textarea {
  font-family: 'Courier New', monospace;
  font-size: 12px;
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

.history-registrar {
  flex: 1;
  color: #64748b;
}

.history-time {
  font-size: 12px;
  color: #9ca3af;
  min-width: 120px;
}

@media (max-width: 768px) {
  .whois-lookup-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .dns-list {
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
  .history-registrar,
  .history-time {
    min-width: unset;
  }
}
</style>