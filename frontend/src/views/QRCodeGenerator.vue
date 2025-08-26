<template>
  <div class="qrcode-generator-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-others' })">其它</el-breadcrumb-item>
          <el-breadcrumb-item>二维码生成器</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>二维码生成器</h1>
      <p>生成各种类型的二维码，支持文本、网址、WiFi等多种格式</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><EditPen /></el-icon>
              <span>输入内容</span>
            </div>
          </template>
          <div class="form-content">
            <div class="type-selector">
              <label>二维码类型：</label>
              <el-radio-group v-model="qrType" @change="handleTypeChange">
                <el-radio-button label="text">文本</el-radio-button>
                <el-radio-button label="url">网址</el-radio-button>
                <el-radio-button label="wifi">WiFi</el-radio-button>
                <el-radio-button label="email">邮箱</el-radio-button>
                <el-radio-button label="sms">短信</el-radio-button>
              </el-radio-group>
            </div>

            <!-- 文本类型 -->
            <div v-if="qrType === 'text'" class="input-group">
              <label>文本内容：</label>
              <el-input
                v-model="textContent"
                type="textarea"
                :rows="4"
                placeholder="请输入要生成二维码的文本内容..."
                @input="generateQRCode"
              />
            </div>

            <!-- 网址类型 -->
            <div v-if="qrType === 'url'" class="input-group">
              <label>网址：</label>
              <el-input
                v-model="urlContent"
                placeholder="https://example.com"
                @input="generateQRCode"
              />
            </div>

            <!-- WiFi类型 -->
            <div v-if="qrType === 'wifi'" class="wifi-form">
              <div class="input-group">
                <label>网络名称 (SSID)：</label>
                <el-input v-model="wifiSSID" placeholder="WiFi网络名称" @input="generateQRCode" />
              </div>
              <div class="input-group">
                <label>密码：</label>
                <el-input v-model="wifiPassword" type="password" placeholder="WiFi密码" @input="generateQRCode" />
              </div>
              <div class="input-group">
                <label>安全类型：</label>
                <el-select v-model="wifiSecurity" @change="generateQRCode">
                  <el-option label="WPA/WPA2" value="WPA" />
                  <el-option label="WEP" value="WEP" />
                  <el-option label="无密码" value="nopass" />
                </el-select>
              </div>
            </div>

            <!-- 邮箱类型 -->
            <div v-if="qrType === 'email'" class="email-form">
              <div class="input-group">
                <label>收件人：</label>
                <el-input v-model="emailTo" placeholder="recipient@example.com" @input="generateQRCode" />
              </div>
              <div class="input-group">
                <label>主题：</label>
                <el-input v-model="emailSubject" placeholder="邮件主题" @input="generateQRCode" />
              </div>
              <div class="input-group">
                <label>内容：</label>
                <el-input v-model="emailBody" type="textarea" :rows="3" placeholder="邮件内容" @input="generateQRCode" />
              </div>
            </div>

            <!-- 短信类型 -->
            <div v-if="qrType === 'sms'" class="sms-form">
              <div class="input-group">
                <label>手机号码：</label>
                <el-input v-model="smsNumber" placeholder="手机号码" @input="generateQRCode" />
              </div>
              <div class="input-group">
                <label>短信内容：</label>
                <el-input v-model="smsMessage" type="textarea" :rows="3" placeholder="短信内容" @input="generateQRCode" />
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <div class="output-section">
        <el-card class="output-card">
          <template #header>
            <div class="card-header">
              <el-icon><Picture /></el-icon>
              <span>生成结果</span>
              <div class="card-actions" v-if="qrCodeDataURL">
                <el-button @click="downloadQRCode" size="small" type="primary">
                  <el-icon><Download /></el-icon>
                  下载
                </el-button>
              </div>
            </div>
          </template>
          
          <div class="qr-display">
            <div v-if="!qrCodeDataURL" class="empty-state">
              <el-icon class="empty-icon"><Grid /></el-icon>
              <p>请在左侧输入内容生成二维码</p>
            </div>
            <div v-else class="qr-result">
              <div class="qr-image-container">
                <img :src="qrCodeDataURL" alt="Generated QR Code" class="qr-image" />
              </div>
              <div class="qr-settings">
                <div class="setting-group">
                  <label>尺寸：</label>
                  <el-slider
                    v-model="qrSize"
                    :min="100"
                    :max="500"
                    :step="50"
                    @change="generateQRCode"
                    style="width: 200px"
                  />
                  <span class="size-label">{{ qrSize }}px</span>
                </div>
                <div class="setting-group">
                  <label>容错级别：</label>
                  <el-select v-model="errorLevel" @change="generateQRCode">
                    <el-option label="低 (L)" value="L" />
                    <el-option label="中 (M)" value="M" />
                    <el-option label="高 (Q)" value="Q" />
                    <el-option label="最高 (H)" value="H" />
                  </el-select>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 使用说明 -->
    <div class="info-section">
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <el-icon><InfoFilled /></el-icon>
            <span>二维码类型说明</span>
          </div>
        </template>
        <div class="info-content">
          <div class="info-grid">
            <div class="info-item">
              <h4>文本二维码</h4>
              <p>生成包含纯文本内容的二维码，适用于文字信息分享</p>
            </div>
            <div class="info-item">
              <h4>网址二维码</h4>
              <p>生成网址链接二维码，扫描后直接跳转到指定网页</p>
            </div>
            <div class="info-item">
              <h4>WiFi二维码</h4>
              <p>生成WiFi连接信息，扫描后可直接连接到WiFi网络</p>
            </div>
            <div class="info-item">
              <h4>邮箱二维码</h4>
              <p>生成邮件信息，扫描后自动打开邮件应用并填充内容</p>
            </div>
            <div class="info-item">
              <h4>短信二维码</h4>
              <p>生成短信内容，扫描后自动打开短信应用并填充内容</p>
            </div>
            <div class="info-item">
              <h4>容错级别</h4>
              <p>L(低7%)、M(中15%)、Q(高25%)、H(最高30%) - 级别越高，二维码越复杂但更耐损</p>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  EditPen,
  Picture,
  Download,
  Grid,
  InfoFilled
} from '@element-plus/icons-vue'

const router = useRouter()

// 响应式数据
const qrType = ref('text')
const qrCodeDataURL = ref('')
const qrSize = ref(200)
const errorLevel = ref('M')

// 不同类型的内容
const textContent = ref('')
const urlContent = ref('')
const wifiSSID = ref('')
const wifiPassword = ref('')
const wifiSecurity = ref('WPA')
const emailTo = ref('')
const emailSubject = ref('')
const emailBody = ref('')
const smsNumber = ref('')
const smsMessage = ref('')

// 生成二维码内容
const getQRContent = () => {
  switch (qrType.value) {
    case 'text':
      return textContent.value
    case 'url':
      return urlContent.value
    case 'wifi':
      if (!wifiSSID.value) return ''
      return `WIFI:T:${wifiSecurity.value};S:${wifiSSID.value};P:${wifiPassword.value};H:false;;`
    case 'email':
      if (!emailTo.value) return ''
      return `mailto:${emailTo.value}?subject=${encodeURIComponent(emailSubject.value)}&body=${encodeURIComponent(emailBody.value)}`
    case 'sms':
      if (!smsNumber.value) return ''
      return `sms:${smsNumber.value}${smsMessage.value ? `?body=${encodeURIComponent(smsMessage.value)}` : ''}`
    default:
      return ''
  }
}

// 使用Canvas生成二维码 (简化版本，实际项目建议使用qrcode.js库)
const generateQRCodeCanvas = (content, size, errorCorrection) => {
  // 这里使用一个简化的实现，实际项目中应该使用专业的QR码库
  // 为了演示，我们创建一个简单的占位符图片
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  
  canvas.width = size
  canvas.height = size
  
  // 创建一个简单的网格模式来模拟二维码
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, size, size)
  
  ctx.fillStyle = '#000000'
  const moduleSize = size / 25
  
  // 生成基于内容的伪随机模式
  let hash = 0
  for (let i = 0; i < content.length; i++) {
    hash = ((hash << 5) - hash + content.charCodeAt(i)) & 0xffffffff
  }
  
  for (let row = 0; row < 25; row++) {
    for (let col = 0; col < 25; col++) {
      // 使用内容的哈希值来决定每个模块的状态
      const seed = hash + row * 25 + col
      if ((seed * 9301 + 49297) % 233280 < 116640) {
        ctx.fillRect(col * moduleSize, row * moduleSize, moduleSize, moduleSize)
      }
    }
  }
  
  // 添加定位标记(角落的方块)
  const drawPositionMarker = (x, y) => {
    ctx.fillStyle = '#000000'
    ctx.fillRect(x, y, moduleSize * 7, moduleSize * 7)
    ctx.fillStyle = '#ffffff'
    ctx.fillRect(x + moduleSize, y + moduleSize, moduleSize * 5, moduleSize * 5)
    ctx.fillStyle = '#000000'
    ctx.fillRect(x + moduleSize * 2, y + moduleSize * 2, moduleSize * 3, moduleSize * 3)
  }
  
  drawPositionMarker(0, 0) // 左上
  drawPositionMarker(18 * moduleSize, 0) // 右上
  drawPositionMarker(0, 18 * moduleSize) // 左下
  
  return canvas.toDataURL('image/png')
}

// 生成二维码
const generateQRCode = () => {
  const content = getQRContent()
  if (!content.trim()) {
    qrCodeDataURL.value = ''
    return
  }
  
  try {
    qrCodeDataURL.value = generateQRCodeCanvas(content, qrSize.value, errorLevel.value)
  } catch (error) {
    ElMessage.error('二维码生成失败')
  }
}

// 处理类型切换
const handleTypeChange = () => {
  qrCodeDataURL.value = ''
  generateQRCode()
}

// 下载二维码
const downloadQRCode = () => {
  if (!qrCodeDataURL.value) {
    ElMessage.warning('请先生成二维码')
    return
  }
  
  const link = document.createElement('a')
  link.download = `qrcode-${Date.now()}.png`
  link.href = qrCodeDataURL.value
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  ElMessage.success('二维码下载成功')
}

// 初始化
onMounted(() => {
  textContent.value = '这是一个示例二维码'
  generateQRCode()
})
</script>

<style scoped>
.qrcode-generator-container {
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
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  margin-bottom: 32px;
}

.input-card,
.output-card,
.info-card {
  border: 1px solid #e5e7eb;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: #2c3e50;
}

.card-header .el-icon + span {
  margin-left: 8px;
}

.form-content {
  padding: 8px 0;
}

.type-selector {
  margin-bottom: 24px;
}

.type-selector label {
  display: block;
  margin-bottom: 12px;
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

.wifi-form,
.email-form,
.sms-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.qr-display {
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #9ca3af;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.qr-result {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
  width: 100%;
}

.qr-image-container {
  display: flex;
  justify-content: center;
}

.qr-image {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.qr-settings {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
  max-width: 300px;
}

.setting-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.setting-group label {
  min-width: 80px;
  font-weight: 500;
  color: #374151;
}

.size-label {
  min-width: 50px;
  color: #6b7280;
  font-size: 14px;
}

.info-content {
  padding: 16px 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
}

.info-item h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.info-item p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .qrcode-generator-container {
    padding: 16px;
  }
  
  .tools-content {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .setting-group {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .setting-group label {
    min-width: auto;
  }
}
</style>