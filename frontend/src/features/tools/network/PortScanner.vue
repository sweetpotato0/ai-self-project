<template>
  <div class="port-scanner-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-operations' })">运维类</el-breadcrumb-item>
          <el-breadcrumb-item>端口扫描</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>端口扫描工具</h1>
      <p>检测目标主机的开放端口，分析网络服务状态</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Monitor /></el-icon>
              <span>扫描配置</span>
            </div>
          </template>
          <div class="input-content">
            <div class="basic-config">
              <div class="input-group">
                <label>目标主机</label>
                <el-input
                  v-model="targetHost"
                  placeholder="请输入IP地址或域名，如：192.168.1.1 或 www.example.com"
                  size="large"
                  clearable
                  @keyup.enter="startScan"
                >
                  <template #prepend>
                    <el-select v-model="selectedExample" placeholder="示例" style="width: 120px" @change="loadExample">
                      <el-option label="本地主机" value="127.0.0.1" />
                      <el-option label="路由器" value="192.168.1.1" />
                      <el-option label="Google DNS" value="8.8.8.8" />
                      <el-option label="百度" value="www.baidu.com" />
                    </el-select>
                  </template>
                </el-input>
              </div>

              <div class="port-config">
                <div class="port-input-group">
                  <div class="input-group">
                    <label>扫描模式</label>
                    <el-radio-group v-model="scanMode" @change="onScanModeChange">
                      <el-radio-button label="common">常用端口</el-radio-button>
                      <el-radio-button label="range">端口范围</el-radio-button>
                      <el-radio-button label="custom">自定义</el-radio-button>
                    </el-radio-group>
                  </div>
                </div>

                <div v-if="scanMode === 'range'" class="range-input">
                  <div class="input-group">
                    <label>端口范围</label>
                    <div class="range-inputs">
                      <el-input-number 
                        v-model="portRange.start" 
                        :min="1" 
                        :max="65535" 
                        size="large"
                        placeholder="起始端口"
                      />
                      <span class="range-separator">-</span>
                      <el-input-number 
                        v-model="portRange.end" 
                        :min="portRange.start || 1" 
                        :max="65535" 
                        size="large"
                        placeholder="结束端口"
                      />
                    </div>
                  </div>
                </div>

                <div v-if="scanMode === 'custom'" class="custom-ports">
                  <div class="input-group">
                    <label>自定义端口</label>
                    <el-input
                      v-model="customPorts"
                      type="textarea"
                      :rows="3"
                      placeholder="请输入端口号，用逗号分隔，如：22,80,443,3389,8080"
                    />
                  </div>
                </div>

                <div v-if="scanMode === 'common'" class="common-ports-display">
                  <div class="common-ports-info">
                    <span class="info-label">将扫描常用端口：</span>
                    <div class="port-tags">
                      <el-tag
                        v-for="port in commonPorts.slice(0, 10)"
                        :key="port.port"
                        size="small"
                        :title="port.description"
                      >
                        {{ port.port }}
                      </el-tag>
                      <el-tag v-if="commonPorts.length > 10" size="small">
                        +{{ commonPorts.length - 10 }}个
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>

              <div class="advanced-options">
                <el-collapse>
                  <el-collapse-item title="高级选项" name="advanced">
                    <div class="advanced-grid">
                      <div class="option-item">
                        <label>超时时间</label>
                        <el-input-number 
                          v-model="scanOptions.timeout" 
                          :min="1000" 
                          :max="10000" 
                          :step="500"
                          size="small"
                        />
                        <span class="unit">ms</span>
                      </div>
                      <div class="option-item">
                        <label>并发数</label>
                        <el-input-number 
                          v-model="scanOptions.concurrent" 
                          :min="1" 
                          :max="100" 
                          size="small"
                        />
                        <span class="unit">个</span>
                      </div>
                    </div>
                  </el-collapse-item>
                </el-collapse>
              </div>
            </div>

            <div class="action-buttons">
              <el-button
                @click="startScan"
                type="primary"
                size="large"
                :loading="scanning"
                :disabled="!targetHost.trim()"
              >
                <el-icon><Search /></el-icon>
                {{ scanning ? '扫描中...' : '开始扫描' }}
              </el-button>
              <el-button @click="stopScan" size="large" :disabled="!scanning">
                <el-icon><Close /></el-icon>
                停止扫描
              </el-button>
              <el-button @click="clearResults" size="large">
                <el-icon><RefreshLeft /></el-icon>
                清空结果
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="scanResults.length > 0 || scanning || error" class="result-section">
        <el-card class="result-card">
          <template #header>
            <div class="card-header">
              <el-icon><DataLine /></el-icon>
              <span>扫描结果</span>
              <div class="result-stats" v-if="scanStats.total > 0">
                <el-tag type="success" size="small">开放: {{ scanStats.open }}</el-tag>
                <el-tag type="danger" size="small">关闭: {{ scanStats.closed }}</el-tag>
                <el-tag type="info" size="small">总计: {{ scanStats.total }}</el-tag>
              </div>
            </div>
          </template>

          <div v-if="error" class="error-content">
            <el-alert
              :title="error"
              type="error"
              show-icon
              :closable="true"
              @close="error = ''"
            >
              <template #default>
                <div class="error-details">
                  <p v-if="error.includes('目标地址')">
                    <strong>建议：</strong>
                    <br>• IP地址格式：192.168.1.1（四段完整格式）
                    <br>• 域名格式：www.example.com
                  </p>
                  <p v-else-if="error.includes('网络连接')">
                    <strong>建议：</strong>
                    <br>• 检查网络连接是否正常
                    <br>• 确认后端服务是否启动
                  </p>
                  <p v-else-if="error.includes('服务器内部')">
                    <strong>建议：</strong>
                    <br>• 稍后重试扫描
                    <br>• 检查后端服务日志
                  </p>
                </div>
              </template>
            </el-alert>
          </div>

          <div class="scan-content">
            <div v-if="scanning" class="scanning-status">
              <div class="scan-progress">
                <el-progress 
                  :percentage="scanProgress" 
                  :status="scanProgress === 100 ? 'success' : undefined"
                />
                <div class="progress-info">
                  <span>{{ currentScanInfo }}</span>
                </div>
              </div>
            </div>

            <div v-if="scanResults.length > 0" class="results-content">
              <div class="results-filters">
                <el-radio-group v-model="filterStatus" size="small">
                  <el-radio-button label="all">全部</el-radio-button>
                  <el-radio-button label="open">开放</el-radio-button>
                  <el-radio-button label="closed">关闭</el-radio-button>
                </el-radio-group>
              </div>

              <div class="results-table">
                <div class="table-header">
                  <div class="header-item">端口</div>
                  <div class="header-item">状态</div>
                  <div class="header-item">服务</div>
                  <div class="header-item">描述</div>
                  <div class="header-item">响应时间</div>
                </div>
                <div class="table-body">
                  <div
                    v-for="result in filteredResults"
                    :key="result.port"
                    class="result-row"
                    :class="{ 'open-port': result.status === 'open' }"
                  >
                    <div class="result-cell port-cell">{{ result.port }}</div>
                    <div class="result-cell status-cell">
                      <el-tag 
                        :type="result.status === 'open' ? 'success' : 'danger'" 
                        size="small"
                      >
                        {{ result.status === 'open' ? '开放' : '关闭' }}
                      </el-tag>
                    </div>
                    <div class="result-cell service-cell">{{ result.service || '-' }}</div>
                    <div class="result-cell description-cell">{{ result.description || '-' }}</div>
                    <div class="result-cell time-cell">{{ result.responseTime || '-' }}ms</div>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="scanResults.length > 0" class="action-section">
              <el-button @click="exportResults" type="success">
                <el-icon><Download /></el-icon>
                导出结果
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
              <span>扫描历史</span>
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
              <div class="history-target">{{ item.target }}</div>
              <div class="history-stats">
                <el-tag type="success" size="small">{{ item.openCount }}开放</el-tag>
                <el-tag type="info" size="small">{{ item.totalCount }}总计</el-tag>
              </div>
              <div class="history-time">{{ new Date(item.scanTime).toLocaleString() }}</div>
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { toolsApi } from '@/api/toolsApi'
import {
  Monitor,
  Search,
  Close,
  RefreshLeft,
  DataLine,
  Download,
  StarFilled,
  Clock,
  Delete
} from '@element-plus/icons-vue'

const router = useRouter()

const targetHost = ref('')
const selectedExample = ref('')
const scanMode = ref('common')
const portRange = ref({ start: 1, end: 1000 })
const customPorts = ref('')
const scanning = ref(false)
const scanResults = ref([])
const error = ref('')
const filterStatus = ref('all')
const history = ref([])

const scanOptions = ref({
  timeout: 3000,
  concurrent: 20
})

const scanProgress = ref(0)
const currentScanInfo = ref('')

// 常用端口配置
const commonPorts = ref([
  { port: 21, service: 'FTP', description: 'File Transfer Protocol' },
  { port: 22, service: 'SSH', description: 'Secure Shell' },
  { port: 23, service: 'Telnet', description: 'Telnet Protocol' },
  { port: 25, service: 'SMTP', description: 'Simple Mail Transfer Protocol' },
  { port: 53, service: 'DNS', description: 'Domain Name System' },
  { port: 80, service: 'HTTP', description: 'HyperText Transfer Protocol' },
  { port: 110, service: 'POP3', description: 'Post Office Protocol v3' },
  { port: 143, service: 'IMAP', description: 'Internet Message Access Protocol' },
  { port: 443, service: 'HTTPS', description: 'HTTP over TLS/SSL' },
  { port: 993, service: 'IMAPS', description: 'IMAP over TLS/SSL' },
  { port: 995, service: 'POP3S', description: 'POP3 over TLS/SSL' },
  { port: 1433, service: 'MSSQL', description: 'Microsoft SQL Server' },
  { port: 3306, service: 'MySQL', description: 'MySQL Database' },
  { port: 3389, service: 'RDP', description: 'Remote Desktop Protocol' },
  { port: 5432, service: 'PostgreSQL', description: 'PostgreSQL Database' },
  { port: 5900, service: 'VNC', description: 'Virtual Network Computing' },
  { port: 8080, service: 'HTTP-Alt', description: 'HTTP Alternative Port' },
  { port: 8443, service: 'HTTPS-Alt', description: 'HTTPS Alternative Port' }
])

const scanStats = computed(() => {
  const open = scanResults.value.filter(r => r.status === 'open').length
  const closed = scanResults.value.filter(r => r.status === 'closed').length
  return {
    open,
    closed,
    total: scanResults.value.length
  }
})

const filteredResults = computed(() => {
  if (filterStatus.value === 'all') {
    return scanResults.value
  }
  return scanResults.value.filter(result => result.status === filterStatus.value)
})

onMounted(() => {
  loadHistory()
})

const loadExample = () => {
  targetHost.value = selectedExample.value
}

const onScanModeChange = () => {
  // 清空结果当切换模式时
  scanResults.value = []
  error.value = ''
}

const startScan = async () => {
  if (!targetHost.value.trim()) {
    ElMessage.warning('请输入目标主机')
    return
  }

  scanning.value = true
  scanResults.value = []
  error.value = ''
  scanProgress.value = 0
  currentScanInfo.value = '准备扫描...'

  try {
    const ports = getPortsToScan()
    if (ports.length === 0) {
      throw new Error('没有有效的端口可扫描')
    }

    currentScanInfo.value = `开始扫描 ${targetHost.value} 的 ${ports.length} 个端口...`
    
    // 调用后端API
    const response = await toolsApi.portScan({
      target: targetHost.value,
      ports: ports,
      timeout: scanOptions.value.timeout,
      concurrent: scanOptions.value.concurrent
    })
    
    // 处理API响应
    if (response.data && response.data.results) {
      scanResults.value = response.data.results
      scanProgress.value = 100
      currentScanInfo.value = '扫描完成'
    } else {
      throw new Error('无效的API响应')
    }
    
    ElMessage.success('端口扫描完成')
  } catch (err) {
    // 提取更友好的错误信息
    let errorMessage = '扫描失败'
    if (err.response?.data?.message) {
      // 如果后端返回了错误消息，使用后端的消息
      errorMessage = err.response.data.message
    } else if (err.response?.status === 400) {
      errorMessage = '请检查目标地址格式是否正确'
    } else if (err.response?.status === 500) {
      errorMessage = '服务器内部错误，请稍后重试'
    } else if (err.message?.includes('Network Error')) {
      errorMessage = '网络连接错误，请检查网络连接'
    } else if (err.message?.includes('timeout')) {
      errorMessage = '请求超时，请稍后重试'
    }
    
    error.value = errorMessage
    // 只显示一次错误消息，不重复弹窗
  } finally {
    scanning.value = false
    scanProgress.value = 100
    currentScanInfo.value = '扫描完成'
  }
}

const getPortsToScan = () => {
  switch (scanMode.value) {
    case 'common':
      return commonPorts.value.map(p => p.port)
    case 'range':
      if (!portRange.value.start || !portRange.value.end) {
        throw new Error('请输入有效的端口范围')
      }
      const ports = []
      for (let i = portRange.value.start; i <= portRange.value.end; i++) {
        ports.push(i)
      }
      return ports
    case 'custom':
      if (!customPorts.value.trim()) {
        throw new Error('请输入自定义端口')
      }
      return customPorts.value.split(',')
        .map(p => parseInt(p.trim()))
        .filter(p => p > 0 && p <= 65535)
    default:
      return []
  }
}

// 模拟端口扫描（将来替换为真实API调用）
const simulatePortScan = async (ports) => {
  const total = ports.length
  
  for (let i = 0; i < ports.length; i++) {
    const port = ports[i]
    const progress = Math.floor((i / total) * 100)
    scanProgress.value = progress
    currentScanInfo.value = `正在扫描端口 ${port}... (${i + 1}/${total})`
    
    // 模拟扫描延迟
    await new Promise(resolve => setTimeout(resolve, 50))
    
    // 模拟扫描结果
    const isOpen = Math.random() > 0.8 // 20% 概率开放
    const responseTime = isOpen ? Math.floor(Math.random() * 100) + 1 : null
    
    const portInfo = commonPorts.value.find(p => p.port === port)
    
    scanResults.value.push({
      port,
      status: isOpen ? 'open' : 'closed',
      service: portInfo?.service || '',
      description: portInfo?.description || '',
      responseTime: responseTime
    })
  }
  
  scanProgress.value = 100
  currentScanInfo.value = '扫描完成'
}

const stopScan = () => {
  scanning.value = false
  currentScanInfo.value = '扫描已停止'
  ElMessage.info('扫描已停止')
}

const clearResults = () => {
  scanResults.value = []
  error.value = ''
  scanProgress.value = 0
  currentScanInfo.value = ''
}

const exportResults = () => {
  const csvContent = [
    'Port,Status,Service,Description,ResponseTime',
    ...scanResults.value.map(r => 
      `${r.port},${r.status},${r.service},${r.description},${r.responseTime || ''}`
    )
  ].join('\n')
  
  const blob = new Blob([csvContent], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `port-scan-${targetHost.value}-${new Date().toISOString().split('T')[0]}.csv`
  a.click()
  URL.revokeObjectURL(url)
  
  ElMessage.success('扫描结果已导出')
}

const addToHistory = () => {
  if (scanResults.value.length === 0) return
  
  const existingIndex = history.value.findIndex(item => item.target === targetHost.value)
  if (existingIndex >= 0) {
    history.value.splice(existingIndex, 1)
  }
  
  history.value.unshift({
    target: targetHost.value,
    results: [...scanResults.value],
    openCount: scanStats.value.open,
    totalCount: scanStats.value.total,
    scanTime: Date.now()
  })
  
  if (history.value.length > 20) {
    history.value = history.value.slice(0, 20)
  }
  
  saveHistory()
  ElMessage.success('已保存到扫描历史')
}

const loadFromHistory = (item) => {
  targetHost.value = item.target
  scanResults.value = [...item.results]
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
  ElMessage.success('已清空扫描历史')
}

const saveHistory = () => {
  localStorage.setItem('port-scan-history', JSON.stringify(history.value))
}

const loadHistory = () => {
  try {
    const saved = localStorage.getItem('port-scan-history')
    if (saved) {
      history.value = JSON.parse(saved)
    }
  } catch (err) {
    console.error('加载历史记录失败：', err)
  }
}
</script>

<style scoped>
.port-scanner-container {
  padding: 20px;
  max-width: 1400px;
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

.result-stats {
  margin-left: auto;
  display: flex;
  gap: 8px;
}

.card-actions {
  margin-left: auto;
}

.input-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-group label {
  font-weight: 500;
  color: #374151;
}

.port-config {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.range-inputs {
  display: flex;
  align-items: center;
  gap: 12px;
}

.range-separator {
  font-size: 18px;
  color: #6b7280;
  font-weight: 500;
}

.common-ports-info {
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.info-label {
  font-size: 14px;
  color: #374151;
  font-weight: 500;
  margin-bottom: 8px;
  display: block;
}

.port-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.advanced-options {
  margin-top: 16px;
}

.advanced-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  padding: 16px 0;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.option-item label {
  font-size: 14px;
  color: #374151;
  white-space: nowrap;
}

.unit {
  font-size: 12px;
  color: #6b7280;
}

.action-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.error-content {
  padding: 20px 0;
}

.error-details {
  margin-top: 8px;
}

.error-details p {
  margin: 12px 0;
  font-size: 14px;
  line-height: 1.5;
  color: #6b7280;
}

.error-details strong {
  color: #ef4444;
  display: block;
  margin-bottom: 8px;
}

.scan-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.scanning-status {
  padding: 20px;
  text-align: center;
}

.scan-progress {
  max-width: 600px;
  margin: 0 auto;
}

.progress-info {
  margin-top: 12px;
  font-size: 14px;
  color: #6b7280;
}

.results-filters {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.results-table {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.table-header {
  display: grid;
  grid-template-columns: 80px 80px 120px 1fr 100px;
  background: #f8fafc;
  border-bottom: 1px solid #e5e7eb;
}

.header-item {
  padding: 12px;
  font-weight: 600;
  color: #374151;
  font-size: 14px;
  border-right: 1px solid #e5e7eb;
}

.header-item:last-child {
  border-right: none;
}

.table-body {
  max-height: 400px;
  overflow-y: auto;
}

.result-row {
  display: grid;
  grid-template-columns: 80px 80px 120px 1fr 100px;
  border-bottom: 1px solid #f3f4f6;
  transition: background-color 0.2s;
}

.result-row:hover {
  background: #f8fafc;
}

.result-row.open-port {
  background: rgba(34, 197, 94, 0.05);
}

.result-cell {
  padding: 12px;
  font-size: 14px;
  color: #374151;
  border-right: 1px solid #f3f4f6;
  display: flex;
  align-items: center;
}

.result-cell:last-child {
  border-right: none;
}

.port-cell {
  font-weight: 600;
  font-family: monospace;
}

.service-cell {
  font-weight: 500;
  color: #059669;
}

.time-cell {
  font-family: monospace;
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

.history-target {
  font-weight: 600;
  color: #3b82f6;
  min-width: 150px;
}

.history-stats {
  display: flex;
  gap: 6px;
  flex: 1;
}

.history-time {
  font-size: 12px;
  color: #9ca3af;
  min-width: 120px;
}

@media (max-width: 768px) {
  .port-scanner-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .range-inputs {
    flex-direction: column;
    gap: 8px;
  }
  
  .range-separator {
    display: none;
  }
  
  .advanced-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .action-section {
    flex-direction: column;
  }
  
  .table-header,
  .result-row {
    grid-template-columns: 60px 60px 80px 1fr 80px;
  }
  
  .result-cell {
    padding: 8px;
    font-size: 12px;
  }
  
  .history-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .history-target,
  .history-time {
    min-width: unset;
  }
}
</style>