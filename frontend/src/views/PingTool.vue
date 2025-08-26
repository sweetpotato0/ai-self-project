<template>
  <div class="ping-tool-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-operations' })">运维类</el-breadcrumb-item>
          <el-breadcrumb-item>Ping测试</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>Ping网络测试工具</h1>
      <p>学习Ping命令的使用方法和网络诊断技巧</p>
    </div>

    <div class="tools-content">
      <div class="simulator-section">
        <el-card class="simulator-card">
          <template #header>
            <div class="card-header">
              <el-icon><Monitor /></el-icon>
              <span>Ping命令模拟器</span>
            </div>
          </template>
          <div class="simulator-content">
            <div class="input-section">
              <div class="os-selector">
                <label>选择操作系统：</label>
                <el-radio-group v-model="selectedOS" @change="onOSChange">
                  <el-radio-button label="windows">Windows</el-radio-button>
                  <el-radio-button label="linux">Linux</el-radio-button>
                  <el-radio-button label="macos">macOS</el-radio-button>
                </el-radio-group>
              </div>
              
              <div class="input-group">
                <label>目标主机或IP地址：</label>
                <div class="input-wrapper">
                  <el-input 
                    v-model="targetHost" 
                    placeholder="例: www.baidu.com 或 8.8.8.8"
                    @keyup.enter="startPing"
                  />
                  <el-button @click="startPing" type="primary" :loading="pinging">
                    <el-icon><Connection /></el-icon>
                    {{ pinging ? 'Ping中...' : '开始Ping' }}
                  </el-button>
                </div>
              </div>
              
              <div class="options-section">
                <h4>Ping选项：</h4>
                <div class="options-grid">
                  <el-checkbox v-model="options.continuous">持续Ping (-t)</el-checkbox>
                  <el-checkbox v-model="options.ipv6">使用IPv6 (-6)</el-checkbox>
                  <el-checkbox v-model="options.dontFragment">不分段 (-f)</el-checkbox>
                  <el-checkbox v-model="options.recordRoute">记录路由 (-r)</el-checkbox>
                </div>
                <div class="params-grid">
                  <div class="param-item">
                    <label>包大小 (-l):</label>
                    <el-input-number v-model="options.packetSize" :min="1" :max="65000" size="small" />
                    <span class="unit">bytes</span>
                  </div>
                  <div class="param-item">
                    <label>超时时间 (-w):</label>
                    <el-input-number v-model="options.timeout" :min="1000" :max="10000" size="small" />
                    <span class="unit">ms</span>
                  </div>
                  <div class="param-item">
                    <label>Ping次数 (-n):</label>
                    <el-input-number v-model="options.count" :min="1" :max="100" size="small" />
                    <span class="unit">次</span>
                  </div>
                </div>
              </div>
            </div>
            
            <div class="output-section">
              <div class="terminal-header">
                <span>命令行输出</span>
                <el-button @click="clearOutput" size="small" text>清空</el-button>
              </div>
              <div class="terminal-output" ref="terminalRef">
                <div v-for="(line, index) in outputLines" :key="index" class="output-line">
                  {{ line }}
                </div>
                <div v-if="outputLines.length === 0" class="empty-output">
                  点击"开始Ping"查看模拟输出结果...
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <div class="info-section">
        <div class="info-grid">
          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <el-icon><Document /></el-icon>
                <span>Ping命令详解</span>
              </div>
            </template>
            <div class="info-content">
              <h4>基本语法</h4>
              <div class="command-example">
                ping [选项] 目标主机
              </div>
              
              <h4>常用参数</h4>
              <div class="os-params">
                <div class="param-section">
                  <h5>Windows</h5>
                  <div class="params-list">
                    <div class="param-row">
                      <code>-t</code>
                      <span>持续Ping直到手动停止</span>
                    </div>
                    <div class="param-row">
                      <code>-n 次数</code>
                      <span>指定发送Echo请求的次数</span>
                    </div>
                    <div class="param-row">
                      <code>-l 大小</code>
                      <span>指定数据包大小（字节）</span>
                    </div>
                    <div class="param-row">
                      <code>-w 超时</code>
                      <span>指定超时时间（毫秒）</span>
                    </div>
                  </div>
                </div>
                
                <div class="param-section">
                  <h5>Linux</h5>
                  <div class="params-list">
                    <div class="param-row">
                      <code>-c 次数</code>
                      <span>指定发送包的次数</span>
                    </div>
                    <div class="param-row">
                      <code>-s 大小</code>
                      <span>指定数据包大小（字节）</span>
                    </div>
                    <div class="param-row">
                      <code>-W 超时</code>
                      <span>指定超时时间（秒）</span>
                    </div>
                    <div class="param-row">
                      <code>-M do</code>
                      <span>设置不分段标志</span>
                    </div>
                  </div>
                </div>
                
                <div class="param-section">
                  <h5>macOS</h5>
                  <div class="params-list">
                    <div class="param-row">
                      <code>-c 次数</code>
                      <span>指定发送包的次数</span>
                    </div>
                    <div class="param-row">
                      <code>-s 大小</code>
                      <span>指定数据包大小（字节）</span>
                    </div>
                    <div class="param-row">
                      <code>-W 超时</code>
                      <span>指定超时时间（毫秒）</span>
                    </div>
                    <div class="param-row">
                      <code>-D</code>
                      <span>设置不分段标志</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </el-card>

          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <el-icon><DataAnalysis /></el-icon>
                <span>结果分析</span>
              </div>
            </template>
            <div class="info-content">
              <h4>返回信息含义</h4>
              <div class="result-explain">
                <div class="result-item">
                  <strong>Reply from IP:</strong> 成功收到回复
                </div>
                <div class="result-item">
                  <strong>bytes=32:</strong> 数据包大小
                </div>
                <div class="result-item">
                  <strong>time&lt;1ms:</strong> 往返时间（延迟）
                </div>
                <div class="result-item">
                  <strong>TTL=64:</strong> 生存时间
                </div>
              </div>
              
              <h4>常见错误</h4>
              <div class="error-list">
                <div class="error-item">
                  <code>Request timeout</code>
                  <span>请求超时，可能网络不通或防火墙阻拦</span>
                </div>
                <div class="error-item">
                  <code>Destination unreachable</code>
                  <span>目标不可达，检查网络路由</span>
                </div>
                <div class="error-item">
                  <code>Unknown host</code>
                  <span>未知主机，DNS解析失败</span>
                </div>
              </div>
            </div>
          </el-card>

          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <el-icon><Tools /></el-icon>
                <span>使用场景</span>
              </div>
            </template>
            <div class="info-content">
              <h4>网络诊断</h4>
              <ul>
                <li>测试网络连通性</li>
                <li>检查网络延迟</li>
                <li>排查网络故障</li>
                <li>测试DNS解析</li>
              </ul>
              
              <h4>性能监控</h4>
              <ul>
                <li>监控网络质量</li>
                <li>测量网络延迟</li>
                <li>检测丢包率</li>
                <li>分析网络稳定性</li>
              </ul>
              
              <h4>实际示例</h4>
              <div class="example-commands">
                <div class="example-os-section">
                  <h5>Windows</h5>
                  <div class="command-item">
                    <code>ping www.google.com</code>
                    <span>测试到Google的连通性</span>
                  </div>
                  <div class="command-item">
                    <code>ping -t 192.168.1.1</code>
                    <span>持续Ping路由器</span>
                  </div>
                  <div class="command-item">
                    <code>ping -l 1024 -n 10 baidu.com</code>
                    <span>用1KB包Ping百度10次</span>
                  </div>
                </div>
                
                <div class="example-os-section">
                  <h5>Linux/macOS</h5>
                  <div class="command-item">
                    <code>ping -c 4 www.google.com</code>
                    <span>Ping Google 4次</span>
                  </div>
                  <div class="command-item">
                    <code>ping 192.168.1.1</code>
                    <span>持续Ping路由器（Ctrl+C停止）</span>
                  </div>
                  <div class="command-item">
                    <code>ping -s 1024 -c 10 baidu.com</code>
                    <span>用1KB包Ping百度10次</span>
                  </div>
                  <div class="command-item">
                    <code>ping -W 2 -c 5 example.com</code>
                    <span>2秒超时，Ping 5次</span>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Monitor,
  Connection,
  Document,
  DataAnalysis,
  Tools
} from '@element-plus/icons-vue'

const router = useRouter()

const targetHost = ref('www.baidu.com')
const pinging = ref(false)
const outputLines = ref([])
const terminalRef = ref()
const selectedOS = ref('windows')

const options = ref({
  continuous: false,
  ipv6: false,
  dontFragment: false,
  recordRoute: false,
  packetSize: 32,
  timeout: 4000,
  count: 4
})

const osConfigs = {
  windows: {
    prompt: 'C:\\>',
    packetSizeParam: '-l',
    timeoutParam: '-w',
    countParam: '-n',
    continuousParam: '-t',
    ipv6Param: '-6',
    dontFragmentParam: '-f',
    recordRouteParam: '-r',
    defaultPacketSize: 32,
    timeoutUnit: 'ms'
  },
  linux: {
    prompt: '$',
    packetSizeParam: '-s',
    timeoutParam: '-W',
    countParam: '-c',
    continuousParam: '', // Linux默认就是持续的，需要用-c限制次数
    ipv6Param: '-6',
    dontFragmentParam: '-M do',
    recordRouteParam: '-R',
    defaultPacketSize: 56,
    timeoutUnit: 'seconds'
  },
  macos: {
    prompt: '$',
    packetSizeParam: '-s',
    timeoutParam: '-W',
    countParam: '-c',
    continuousParam: '',
    ipv6Param: '-6',
    dontFragmentParam: '-D',
    recordRouteParam: '-R',
    defaultPacketSize: 56,
    timeoutUnit: 'ms'
  }
}

const generatePingOutput = () => {
  const host = targetHost.value || 'example.com'
  const isIP = /^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/.test(host)
  const targetIP = isIP ? host : `${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`
  
  const config = osConfigs[selectedOS.value]
  const lines = []
  
  // 构建命令行
  let command = `ping ${host}`
  
  if (selectedOS.value === 'windows') {
    if (options.value.continuous) command += ` ${config.continuousParam}`
    if (options.value.ipv6) command += ` ${config.ipv6Param}`
    if (options.value.dontFragment) command += ` ${config.dontFragmentParam}`
    if (options.value.recordRoute) command += ` ${config.recordRouteParam}`
    if (options.value.packetSize !== config.defaultPacketSize) command += ` ${config.packetSizeParam} ${options.value.packetSize}`
    if (options.value.timeout !== 4000) command += ` ${config.timeoutParam} ${options.value.timeout}`
    if (!options.value.continuous) command += ` ${config.countParam} ${options.value.count}`
  } else {
    // Linux/macOS
    if (!options.value.continuous) command += ` ${config.countParam} ${options.value.count}`
    if (options.value.ipv6) command += ` ${config.ipv6Param}`
    if (options.value.dontFragment) command += ` ${config.dontFragmentParam}`
    if (options.value.recordRoute) command += ` ${config.recordRouteParam}`
    if (options.value.packetSize !== config.defaultPacketSize) command += ` ${config.packetSizeParam} ${options.value.packetSize}`
    if (selectedOS.value === 'linux') {
      if (options.value.timeout !== 4000) command += ` ${config.timeoutParam} ${Math.floor(options.value.timeout / 1000)}`
    } else {
      if (options.value.timeout !== 4000) command += ` ${config.timeoutParam} ${options.value.timeout}`
    }
  }
  
  lines.push(`${config.prompt} ${command}`)
  lines.push('')
  
  // 不同系统的输出格式
  if (selectedOS.value === 'windows') {
    if (!isIP) {
      lines.push(`Pinging ${host} [${targetIP}] with ${options.value.packetSize} bytes of data:`)
    } else {
      lines.push(`Pinging ${targetIP} with ${options.value.packetSize} bytes of data:`)
    }
  } else {
    // Linux/macOS
    const packetSize = options.value.packetSize + 8 // 加上ICMP头
    if (!isIP) {
      lines.push(`PING ${host} (${targetIP}) ${options.value.packetSize}(${packetSize}) bytes of data.`)
    } else {
      lines.push(`PING ${targetIP} (${targetIP}) ${options.value.packetSize}(${packetSize}) bytes of data.`)
    }
  }
  lines.push('')
  
  return { lines, targetIP, config }
}

const startPing = async () => {
  if (!targetHost.value.trim()) {
    ElMessage.warning('请输入目标主机或IP地址')
    return
  }
  
  pinging.value = true
  outputLines.value = []
  
  const { lines, targetIP, config } = generatePingOutput()
  
  // 逐行显示输出
  for (let i = 0; i < lines.length; i++) {
    outputLines.value.push(lines[i])
    await new Promise(resolve => setTimeout(resolve, 100))
    scrollToBottom()
  }
  
  // 模拟Ping结果
  const count = options.value.continuous ? 10 : options.value.count
  
  for (let i = 0; i < count; i++) {
    if (Math.random() > 0.1) { // 90% 成功率
      const delay = Math.floor(Math.random() * 50) + 1
      const ttl = Math.floor(Math.random() * 10) + 55
      
      if (selectedOS.value === 'windows') {
        outputLines.value.push(`Reply from ${targetIP}: bytes=${options.value.packetSize} time=${delay}ms TTL=${ttl}`)
      } else {
        // Linux/macOS格式
        const packetSize = options.value.packetSize + 8
        const icmpSeq = i + 1
        outputLines.value.push(`${packetSize} bytes from ${targetIP}: icmp_seq=${icmpSeq} ttl=${ttl} time=${delay}.${Math.floor(Math.random() * 1000).toString().padStart(3, '0')} ms`)
      }
    } else {
      if (selectedOS.value === 'windows') {
        outputLines.value.push('Request timed out.')
      } else {
        outputLines.value.push(`ping: sendto: Host is down`)
      }
    }
    
    await new Promise(resolve => setTimeout(resolve, 800))
    scrollToBottom()
    
    if (options.value.continuous && i === 9) {
      generateStatistics(targetIP, i + 1, i, config)
      break
    }
  }
  
  if (!options.value.continuous) {
    generateStatistics(targetIP, count, Math.floor(count * 0.9), config)
  }
  
  scrollToBottom()
  pinging.value = false
  ElMessage.success('Ping测试完成')
}

const generateStatistics = (targetIP, sent, received, config) => {
  const lost = sent - received
  const lossPercent = Math.floor((lost / sent) * 100)
  
  outputLines.value.push('')
  
  if (selectedOS.value === 'windows') {
    outputLines.value.push('Ping statistics for ' + targetIP + ':')
    outputLines.value.push(`    Packets: Sent = ${sent}, Received = ${received}, Lost = ${lost} (${lossPercent}% loss),`)
    
    if (received > 0) {
      outputLines.value.push('Approximate round trip times in milli-seconds:')
      outputLines.value.push('    Minimum = 1ms, Maximum = 45ms, Average = 15ms')
    }
  } else {
    // Linux/macOS格式
    outputLines.value.push(`--- ${targetIP} ping statistics ---`)
    outputLines.value.push(`${sent} packets transmitted, ${received} received, ${lossPercent}% packet loss, time ${sent * 1000}ms`)
    
    if (received > 0) {
      outputLines.value.push(`rtt min/avg/max/mdev = 1.234/15.678/45.123/12.456 ms`)
    }
  }
}

const onOSChange = () => {
  const config = osConfigs[selectedOS.value]
  options.value.packetSize = config.defaultPacketSize
  clearOutput()
}

const clearOutput = () => {
  outputLines.value = []
}

const scrollToBottom = async () => {
  await nextTick()
  if (terminalRef.value) {
    terminalRef.value.scrollTop = terminalRef.value.scrollHeight
  }
}
</script>

<style scoped>
.ping-tool-container {
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

.simulator-card {
  border: 1px solid #e5e7eb;
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #2c3e50;
}

.input-section {
  margin-bottom: 24px;
}

.os-selector {
  margin-bottom: 20px;
}

.os-selector label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #374151;
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #374151;
}

.input-wrapper {
  display: flex;
  gap: 12px;
}

.input-wrapper .el-input {
  flex: 1;
}

.options-section h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 16px 0;
}

.options-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.params-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.param-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.param-item label {
  font-size: 14px;
  color: #374151;
  white-space: nowrap;
}

.unit {
  font-size: 12px;
  color: #6b7280;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #374151;
  color: white;
  font-weight: 500;
  border-radius: 8px 8px 0 0;
}

.terminal-output {
  background: #1f2937;
  color: #e5e7eb;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  padding: 16px;
  min-height: 300px;
  max-height: 400px;
  overflow-y: auto;
  border-radius: 0 0 8px 8px;
  border: 1px solid #374151;
  border-top: none;
}

.output-line {
  margin: 2px 0;
  white-space: pre-wrap;
}

.empty-output {
  color: #9ca3af;
  font-style: italic;
  text-align: center;
  margin-top: 60px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.info-card {
  border: 1px solid #e5e7eb;
}

.info-content {
  line-height: 1.6;
}

.info-content h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 12px 0;
}

.command-example {
  background: #f3f4f6;
  padding: 12px;
  border-radius: 6px;
  font-family: monospace;
  font-size: 14px;
  margin-bottom: 16px;
}

.os-params {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.param-section h5 {
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 12px 0;
  padding: 8px 12px;
  background: #f3f4f6;
  border-radius: 6px;
  border-left: 3px solid #409eff;
}

.params-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-left: 16px;
}

.param-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.param-row code {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
  min-width: 80px;
  display: inline-block;
}

.result-explain,
.error-list,
.example-commands {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.example-os-section h5 {
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 12px 0;
  padding: 8px 12px;
  background: #f3f4f6;
  border-radius: 6px;
  border-left: 3px solid #10b981;
}

.example-os-section .command-item {
  margin-left: 16px;
}

.result-item,
.error-item,
.command-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  background: #f8fafc;
  border-radius: 6px;
}

.result-item strong,
.error-item code,
.command-item code {
  font-family: monospace;
  font-size: 13px;
}

.error-item code {
  background: #fee2e2;
  color: #dc2626;
  padding: 2px 6px;
  border-radius: 4px;
}

.command-item code {
  background: #e0f2fe;
  color: #0891b2;
  padding: 2px 6px;
  border-radius: 4px;
}

.info-content ul {
  margin: 0;
  padding-left: 20px;
}

.info-content li {
  margin: 6px 0;
  color: #374151;
}

@media (max-width: 768px) {
  .ping-tool-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .input-wrapper {
    flex-direction: column;
  }
  
  .params-grid {
    grid-template-columns: 1fr;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>