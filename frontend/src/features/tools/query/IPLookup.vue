<template>
  <div class="ip-lookup-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-query' })">查询类</el-breadcrumb-item>
          <el-breadcrumb-item>IP地址查询</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>IP地址查询</h1>
      <p>查询IP地址的归属地信息、ISP运营商和地理位置</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Position /></el-icon>
              <span>IP地址输入</span>
            </div>
          </template>
          <div class="input-content">
            <div class="ip-input-group">
              <el-input
                v-model="ipAddress"
                placeholder="请输入IP地址，如：8.8.8.8"
                size="large"
                clearable
                @keyup.enter="lookupIP"
              >
                <template #prepend>
                  <el-button @click="getCurrentIP" :loading="gettingCurrentIP">
                    <el-icon><Location /></el-icon>
                    获取当前IP
                  </el-button>
                </template>
              </el-input>
            </div>
            <div class="action-buttons">
              <el-button
                @click="lookupIP"
                type="primary"
                size="large"
                :loading="loading"
                :disabled="!ipAddress.trim()"
              >
                <el-icon><Search /></el-icon>
                查询IP信息
              </el-button>
              <el-button @click="clearResults" size="large">
                <el-icon><RefreshLeft /></el-icon>
                清空结果
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="ipResult || error" class="result-section">
        <el-card class="result-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>查询结果</span>
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

          <div v-else-if="ipResult" class="ip-info-content">
            <div class="ip-info-grid">
              <div class="info-item">
                <div class="info-label">IP地址</div>
                <div class="info-value primary">{{ ipResult.ip }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">国家/地区</div>
                <div class="info-value">
                  {{ ipResult.country }}
                  <span v-if="ipResult.countryCode" class="country-code">{{ ipResult.countryCode }}</span>
                </div>
              </div>
              <div class="info-item">
                <div class="info-label">省份/州</div>
                <div class="info-value">{{ ipResult.region || '-' }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">城市</div>
                <div class="info-value">{{ ipResult.city || '-' }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">ISP运营商</div>
                <div class="info-value">{{ ipResult.isp || '-' }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">组织机构</div>
                <div class="info-value">{{ ipResult.org || '-' }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">时区</div>
                <div class="info-value">{{ ipResult.timezone || '-' }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">坐标位置</div>
                <div class="info-value">
                  {{ ipResult.latitude && ipResult.longitude 
                    ? `${ipResult.latitude}, ${ipResult.longitude}` 
                    : '-' }}
                </div>
              </div>
            </div>

            <div class="additional-info">
              <el-divider content-position="left">其他信息</el-divider>
              <div class="extra-info-grid">
                <div class="extra-info-item">
                  <el-icon><Timer /></el-icon>
                  <span>查询时间：{{ new Date(ipResult.queryTime).toLocaleString() }}</span>
                </div>
                <div class="extra-info-item">
                  <el-icon><Connection /></el-icon>
                  <span>IP类型：{{ getIPType(ipResult.ip) }}</span>
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
              <div class="history-ip">{{ item.ip }}</div>
              <div class="history-location">{{ item.country }} {{ item.city }}</div>
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
  Position,
  Location,
  Search,
  RefreshLeft,
  InfoFilled,
  Timer,
  Connection,
  DocumentCopy,
  StarFilled,
  Clock,
  Delete,
  Close
} from '@element-plus/icons-vue'

const router = useRouter()

const ipAddress = ref('')
const loading = ref(false)
const gettingCurrentIP = ref(false)
const ipResult = ref(null)
const error = ref('')
const history = ref([])

onMounted(() => {
  loadHistory()
})

const getCurrentIP = async () => {
  gettingCurrentIP.value = true
  try {
    const response = await fetch('https://api.ipify.org?format=json')
    const data = await response.json()
    ipAddress.value = data.ip
    ElMessage.success('已获取当前IP地址')
  } catch (err) {
    ElMessage.error('获取当前IP失败：' + err.message)
  } finally {
    gettingCurrentIP.value = false
  }
}

const lookupIP = async () => {
  if (!ipAddress.value.trim()) {
    ElMessage.warning('请输入IP地址')
    return
  }

  if (!isValidIP(ipAddress.value.trim())) {
    ElMessage.error('请输入有效的IP地址格式')
    return
  }

  loading.value = true
  error.value = ''
  ipResult.value = null

  try {
    const ip = ipAddress.value.trim()
    const response = await fetch(`http://ip-api.com/json/${ip}?fields=status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query`)
    
    const data = await response.json()
    
    if (data.status === 'success') {
      ipResult.value = {
        ip: data.query,
        country: data.country,
        countryCode: data.countryCode,
        region: data.regionName,
        city: data.city,
        isp: data.isp,
        org: data.org,
        latitude: data.lat,
        longitude: data.lon,
        timezone: data.timezone,
        queryTime: Date.now()
      }
      ElMessage.success('IP查询成功')
    } else {
      error.value = data.message || 'IP查询失败'
      ElMessage.error(error.value)
    }
  } catch (err) {
    error.value = '网络请求失败：' + err.message
    ElMessage.error(error.value)
  } finally {
    loading.value = false
  }
}

const isValidIP = (ip) => {
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  const ipv6Regex = /^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/
  return ipv4Regex.test(ip) || ipv6Regex.test(ip)
}

const getIPType = (ip) => {
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  return ipv4Regex.test(ip) ? 'IPv4' : 'IPv6'
}

const formatResult = () => {
  if (!ipResult.value) return ''
  
  return `IP地址查询结果：
IP地址：${ipResult.value.ip}
国家/地区：${ipResult.value.country} (${ipResult.value.countryCode})
省份/州：${ipResult.value.region || '-'}
城市：${ipResult.value.city || '-'}
ISP运营商：${ipResult.value.isp || '-'}
组织机构：${ipResult.value.org || '-'}
时区：${ipResult.value.timezone || '-'}
坐标位置：${ipResult.value.latitude && ipResult.value.longitude ? `${ipResult.value.latitude}, ${ipResult.value.longitude}` : '-'}
查询时间：${new Date(ipResult.value.queryTime).toLocaleString()}`
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
  if (!ipResult.value) return
  
  const existingIndex = history.value.findIndex(item => item.ip === ipResult.value.ip)
  if (existingIndex >= 0) {
    history.value.splice(existingIndex, 1)
  }
  
  history.value.unshift({ ...ipResult.value })
  
  if (history.value.length > 20) {
    history.value = history.value.slice(0, 20)
  }
  
  saveHistory()
  ElMessage.success('已保存到查询历史')
}

const loadFromHistory = (item) => {
  ipAddress.value = item.ip
  ipResult.value = { ...item }
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
  ipAddress.value = ''
  ipResult.value = null
  error.value = ''
}

const saveHistory = () => {
  localStorage.setItem('ip-lookup-history', JSON.stringify(history.value))
}

const loadHistory = () => {
  try {
    const saved = localStorage.getItem('ip-lookup-history')
    if (saved) {
      history.value = JSON.parse(saved)
    }
  } catch (err) {
    console.error('加载历史记录失败：', err)
  }
}
</script>

<style scoped>
.ip-lookup-container {
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

.ip-input-group {
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

.ip-info-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.ip-info-grid {
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

.country-code {
  font-size: 12px;
  color: #6b7280;
  margin-left: 8px;
  padding: 2px 6px;
  background: #e5e7eb;
  border-radius: 4px;
}

.extra-info-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.extra-info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b7280;
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

.history-ip {
  font-weight: 600;
  color: #3b82f6;
  min-width: 120px;
}

.history-location {
  flex: 1;
  color: #64748b;
}

.history-time {
  font-size: 12px;
  color: #9ca3af;
  min-width: 120px;
}

@media (max-width: 768px) {
  .ip-lookup-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .ip-info-grid {
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
  
  .history-ip,
  .history-location,
  .history-time {
    min-width: unset;
  }
}
</style>