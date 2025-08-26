<template>
  <div class="password-strength-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-others' })">其它</el-breadcrumb-item>
          <el-breadcrumb-item>密码强度检测</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>密码强度检测器</h1>
      <p>检测密码强度并提供安全建议，帮助您创建更安全的密码</p>
    </div>

    <div class="tools-content">
      <!-- 密码输入区域 -->
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Key /></el-icon>
              <span>密码检测</span>
            </div>
          </template>
          
          <div class="password-input">
            <div class="input-group">
              <label>输入密码进行检测：</label>
              <el-input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="请输入您的密码..."
                @input="checkPassword"
                class="password-field"
              >
                <template #suffix>
                  <el-button
                    @click="showPassword = !showPassword"
                    :icon="showPassword ? Hide : View"
                    text
                    class="toggle-visibility"
                  />
                </template>
              </el-input>
            </div>

            <!-- 强度指示器 -->
            <div class="strength-indicator" v-if="password">
              <div class="strength-bar">
                <div 
                  class="strength-fill"
                  :class="strengthClass"
                  :style="{ width: `${strengthScore * 20}%` }"
                ></div>
              </div>
              <div class="strength-label" :class="strengthClass">
                {{ strengthText }}
              </div>
            </div>

            <!-- 密码信息 -->
            <div class="password-info" v-if="password">
              <div class="info-grid">
                <div class="info-item">
                  <span class="info-label">长度:</span>
                  <span class="info-value">{{ passwordStats.length }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">字符类型:</span>
                  <span class="info-value">{{ passwordStats.characterTypes }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">估算破解时间:</span>
                  <span class="info-value">{{ estimatedCrackTime }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">熵值:</span>
                  <span class="info-value">{{ passwordStats.entropy.toFixed(1) }} bits</span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 检测结果区域 -->
      <div class="result-section">
        <el-card class="result-card">
          <template #header>
            <div class="card-header">
              <el-icon><CircleCheck /></el-icon>
              <span>检测结果</span>
            </div>
          </template>
          
          <div v-if="!password" class="empty-state">
            <el-icon class="empty-icon"><Lock /></el-icon>
            <p>请输入密码开始检测</p>
          </div>

          <div v-else class="check-results">
            <!-- 检查项目 -->
            <div class="check-items">
              <div 
                v-for="check in passwordChecks" 
                :key="check.key"
                class="check-item"
                :class="check.status"
              >
                <el-icon class="check-icon">
                  <SuccessFilled v-if="check.status === 'pass'" />
                  <WarningFilled v-else-if="check.status === 'warning'" />
                  <CircleCloseFilled v-else />
                </el-icon>
                <span class="check-text">{{ check.text }}</span>
              </div>
            </div>

            <!-- 建议 -->
            <div class="suggestions" v-if="suggestions.length > 0">
              <h4>安全建议</h4>
              <ul class="suggestion-list">
                <li v-for="suggestion in suggestions" :key="suggestion">
                  {{ suggestion }}
                </li>
              </ul>
            </div>

            <!-- 常见弱密码警告 -->
            <div class="weak-password-warning" v-if="isCommonPassword">
              <el-alert
                title="警告：这是一个常见弱密码"
                type="error"
                :description="commonPasswordMessage"
                show-icon
                :closable="false"
              />
            </div>
          </div>
        </el-card>
      </div>

      <!-- 密码生成器 -->
      <div class="generator-section">
        <el-card class="generator-card">
          <template #header>
            <div class="card-header">
              <el-icon><MagicStick /></el-icon>
              <span>安全密码生成</span>
            </div>
          </template>
          
          <div class="generator-content">
            <div class="generator-options">
              <div class="option-group">
                <label>长度：</label>
                <el-slider
                  v-model="generateOptions.length"
                  :min="8"
                  :max="32"
                  show-input
                  style="width: 200px"
                />
              </div>
              
              <div class="option-checkboxes">
                <el-checkbox v-model="generateOptions.uppercase">大写字母 (A-Z)</el-checkbox>
                <el-checkbox v-model="generateOptions.lowercase">小写字母 (a-z)</el-checkbox>
                <el-checkbox v-model="generateOptions.numbers">数字 (0-9)</el-checkbox>
                <el-checkbox v-model="generateOptions.symbols">特殊符号 (!@#$%)</el-checkbox>
                <el-checkbox v-model="generateOptions.excludeSimilar">排除易混淆字符 (0O1lI)</el-checkbox>
              </div>
            </div>

            <div class="generator-actions">
              <el-button @click="generateSecurePassword" type="primary" size="large">
                <el-icon><Refresh /></el-icon>
                生成安全密码
              </el-button>
            </div>

            <div class="generated-passwords" v-if="generatedPasswords.length > 0">
              <h4>生成的密码</h4>
              <div 
                v-for="(genPassword, index) in generatedPasswords" 
                :key="index"
                class="generated-item"
              >
                <el-input 
                  :value="genPassword.password" 
                  readonly
                  class="generated-password"
                >
                  <template #suffix>
                    <div class="password-actions">
                      <el-button 
                        @click="copyPassword(genPassword.password)" 
                        :icon="DocumentCopy" 
                        text 
                        size="small"
                      />
                      <el-button 
                        @click="usePassword(genPassword.password)" 
                        :icon="Select" 
                        text 
                        size="small"
                        title="使用此密码"
                      />
                    </div>
                  </template>
                </el-input>
                <div class="generated-strength">
                  <span :class="`strength-${genPassword.strength.toLowerCase()}`">
                    {{ genPassword.strength }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 密码安全知识 -->
    <div class="info-section">
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <el-icon><InfoFilled /></el-icon>
            <span>密码安全知识</span>
          </div>
        </template>
        <div class="security-tips">
          <div class="tips-grid">
            <div class="tip-item">
              <h4>密码长度</h4>
              <p>至少使用12个字符，越长越安全</p>
            </div>
            <div class="tip-item">
              <h4>字符组合</h4>
              <p>包含大小写字母、数字和特殊符号</p>
            </div>
            <div class="tip-item">
              <h4>避免常见词汇</h4>
              <p>不使用字典词汇、个人信息或键盘序列</p>
            </div>
            <div class="tip-item">
              <h4>定期更换</h4>
              <p>重要账户密码定期更换</p>
            </div>
            <div class="tip-item">
              <h4>独一无二</h4>
              <p>每个账户使用不同的密码</p>
            </div>
            <div class="tip-item">
              <h4>使用密码管理器</h4>
              <p>使用专业工具管理复杂密码</p>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Key,
  CircleCheck,
  Lock,
  MagicStick,
  Refresh,
  DocumentCopy,
  Select,
  InfoFilled,
  View,
  Hide,
  SuccessFilled,
  WarningFilled,
  CircleCloseFilled
} from '@element-plus/icons-vue'

const router = useRouter()

// 响应式数据
const password = ref('')
const showPassword = ref(false)
const generatedPasswords = ref([])
const generateOptions = ref({
  length: 16,
  uppercase: true,
  lowercase: true,
  numbers: true,
  symbols: true,
  excludeSimilar: true
})

// 常见弱密码列表（部分）
const commonPasswords = [
  '123456', 'password', '123456789', '12345678', '12345', '1234567',
  'qwerty', 'abc123', 'password123', 'admin', '123123', 'welcome',
  'login', '111111', '000000', '1234', '123', 'qwerty123'
]

// 密码统计信息
const passwordStats = computed(() => {
  if (!password.value) {
    return {
      length: 0,
      characterTypes: 0,
      entropy: 0
    }
  }

  const pwd = password.value
  let charTypes = 0
  let charsetSize = 0

  // 检查字符类型
  if (/[a-z]/.test(pwd)) { charTypes++; charsetSize += 26 }
  if (/[A-Z]/.test(pwd)) { charTypes++; charsetSize += 26 }
  if (/[0-9]/.test(pwd)) { charTypes++; charsetSize += 10 }
  if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?`~]/.test(pwd)) { charTypes++; charsetSize += 32 }

  // 计算熵值
  const entropy = pwd.length * Math.log2(charsetSize || 1)

  return {
    length: pwd.length,
    characterTypes: charTypes,
    entropy: entropy
  }
})

// 强度得分 (0-5)
const strengthScore = computed(() => {
  if (!password.value) return 0

  let score = 0
  const pwd = password.value

  // 长度检查
  if (pwd.length >= 8) score++
  if (pwd.length >= 12) score++
  if (pwd.length >= 16) score++

  // 字符类型检查
  if (/[a-z]/.test(pwd) && /[A-Z]/.test(pwd)) score++
  if (/[0-9]/.test(pwd)) score++
  if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?`~]/.test(pwd)) score++

  // 减分项
  if (isCommonPassword.value) score -= 2
  if (/(.)\1{2,}/.test(pwd)) score-- // 重复字符
  if (/123|abc|qwe/i.test(pwd)) score-- // 连续字符

  return Math.max(0, Math.min(5, score))
})

const strengthClass = computed(() => {
  const score = strengthScore.value
  if (score <= 1) return 'strength-weak'
  if (score <= 2) return 'strength-fair'
  if (score <= 3) return 'strength-good'
  if (score <= 4) return 'strength-strong'
  return 'strength-very-strong'
})

const strengthText = computed(() => {
  const score = strengthScore.value
  if (score <= 1) return '很弱'
  if (score <= 2) return '较弱'
  if (score <= 3) return '一般'
  if (score <= 4) return '强'
  return '很强'
})

// 估算破解时间
const estimatedCrackTime = computed(() => {
  if (!password.value) return ''
  
  const entropy = passwordStats.value.entropy
  const attemptsPerSecond = 1e9 // 假设每秒10亿次尝试
  const combinations = Math.pow(2, entropy)
  const seconds = combinations / (2 * attemptsPerSecond) // 平均需要尝试一半

  if (seconds < 1) return '瞬间'
  if (seconds < 60) return `${Math.round(seconds)}秒`
  if (seconds < 3600) return `${Math.round(seconds / 60)}分钟`
  if (seconds < 86400) return `${Math.round(seconds / 3600)}小时`
  if (seconds < 31536000) return `${Math.round(seconds / 86400)}天`
  if (seconds < 31536000000) return `${Math.round(seconds / 31536000)}年`
  return '数万亿年'
})

// 是否为常见密码
const isCommonPassword = computed(() => {
  return commonPasswords.includes(password.value.toLowerCase())
})

const commonPasswordMessage = computed(() => {
  if (isCommonPassword.value) {
    return '此密码在常见弱密码列表中，极易被破解。请立即更换为更强的密码。'
  }
  return ''
})

// 密码检查项目
const passwordChecks = computed(() => {
  if (!password.value) return []

  const pwd = password.value
  const checks = [
    {
      key: 'length',
      text: '长度至少8个字符',
      status: pwd.length >= 8 ? 'pass' : 'fail'
    },
    {
      key: 'length-good',
      text: '长度至少12个字符（推荐）',
      status: pwd.length >= 12 ? 'pass' : pwd.length >= 8 ? 'warning' : 'fail'
    },
    {
      key: 'lowercase',
      text: '包含小写字母',
      status: /[a-z]/.test(pwd) ? 'pass' : 'fail'
    },
    {
      key: 'uppercase',
      text: '包含大写字母',
      status: /[A-Z]/.test(pwd) ? 'pass' : 'warning'
    },
    {
      key: 'numbers',
      text: '包含数字',
      status: /[0-9]/.test(pwd) ? 'pass' : 'warning'
    },
    {
      key: 'symbols',
      text: '包含特殊符号',
      status: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?`~]/.test(pwd) ? 'pass' : 'warning'
    },
    {
      key: 'no-common',
      text: '非常见弱密码',
      status: !isCommonPassword.value ? 'pass' : 'fail'
    },
    {
      key: 'no-repeat',
      text: '无连续重复字符',
      status: !/(.)\1{2,}/.test(pwd) ? 'pass' : 'warning'
    },
    {
      key: 'no-sequence',
      text: '无明显序列',
      status: !/123|abc|qwe/i.test(pwd) ? 'pass' : 'warning'
    }
  ]

  return checks
})

// 建议
const suggestions = computed(() => {
  if (!password.value) return []

  const suggestions = []
  const pwd = password.value

  if (pwd.length < 12) {
    suggestions.push('增加密码长度到至少12个字符')
  }
  if (!/[A-Z]/.test(pwd)) {
    suggestions.push('添加大写字母')
  }
  if (!/[0-9]/.test(pwd)) {
    suggestions.push('添加数字')
  }
  if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?`~]/.test(pwd)) {
    suggestions.push('添加特殊符号')
  }
  if (isCommonPassword.value) {
    suggestions.push('避免使用常见密码')
  }
  if (/(.)\1{2,}/.test(pwd)) {
    suggestions.push('避免连续重复字符')
  }
  if (/123|abc|qwe/i.test(pwd)) {
    suggestions.push('避免使用键盘序列或字母序列')
  }

  return suggestions
})

// 检查密码
const checkPassword = () => {
  // 密码检查逻辑在计算属性中处理
}

// 生成安全密码
const generateSecurePassword = () => {
  const options = generateOptions.value
  let charset = ''
  
  if (options.lowercase) charset += 'abcdefghijklmnopqrstuvwxyz'
  if (options.uppercase) charset += 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  if (options.numbers) charset += '0123456789'
  if (options.symbols) charset += '!@#$%^&*()_+-=[]{}|;:,.<>?'
  
  if (options.excludeSimilar) {
    charset = charset.replace(/[0O1lI]/g, '')
  }
  
  if (!charset) {
    ElMessage.warning('请至少选择一种字符类型')
    return
  }
  
  const passwords = []
  for (let i = 0; i < 3; i++) {
    let newPassword = ''
    for (let j = 0; j < options.length; j++) {
      newPassword += charset.charAt(Math.floor(Math.random() * charset.length))
    }
    
    // 评估生成的密码强度
    const tempPassword = password.value
    password.value = newPassword
    const strength = strengthText.value
    password.value = tempPassword
    
    passwords.push({
      password: newPassword,
      strength: strength
    })
  }
  
  generatedPasswords.value = passwords
  ElMessage.success('已生成3个安全密码')
}

// 复制密码
const copyPassword = async (pwd) => {
  try {
    await navigator.clipboard.writeText(pwd)
    ElMessage.success('密码已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 使用密码
const usePassword = (pwd) => {
  password.value = pwd
  checkPassword()
  ElMessage.success('已应用此密码进行检测')
}
</script>

<style scoped>
.password-strength-container {
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
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: auto auto;
  gap: 24px;
  margin-bottom: 32px;
}

.generator-section {
  grid-column: span 2;
}

.input-card,
.result-card,
.generator-card,
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

.password-input {
  padding: 8px 0;
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 12px;
  font-weight: 500;
  color: #374151;
}

.password-field {
  font-size: 16px;
}

.toggle-visibility {
  color: #6b7280;
}

.strength-indicator {
  margin-top: 16px;
}

.strength-bar {
  height: 8px;
  background: #f1f5f9;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
}

.strength-fill {
  height: 100%;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.strength-weak { background: #ef4444; color: #ef4444; }
.strength-fair { background: #f97316; color: #f97316; }
.strength-good { background: #eab308; color: #eab308; }
.strength-strong { background: #22c55e; color: #22c55e; }
.strength-very-strong { background: #16a34a; color: #16a34a; }

.strength-label {
  font-weight: 600;
  font-size: 14px;
  text-align: right;
}

.password-info {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e5e7eb;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f8fafc;
  border-radius: 6px;
}

.info-label {
  color: #6b7280;
  font-size: 14px;
}

.info-value {
  font-weight: 600;
  color: #374151;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.check-results {
  padding: 8px 0;
}

.check-items {
  margin-bottom: 24px;
}

.check-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  font-size: 14px;
}

.check-item.pass .check-icon { color: #22c55e; }
.check-item.warning .check-icon { color: #f97316; }
.check-item.fail .check-icon { color: #ef4444; }

.check-item.pass .check-text { color: #374151; }
.check-item.warning .check-text { color: #92400e; }
.check-item.fail .check-text { color: #991b1b; }

.suggestions h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 12px 0;
}

.suggestion-list {
  margin: 0;
  padding-left: 20px;
  color: #6b7280;
}

.suggestion-list li {
  margin: 6px 0;
  font-size: 14px;
}

.weak-password-warning {
  margin-top: 20px;
}

.generator-content {
  padding: 8px 0;
}

.generator-options {
  margin-bottom: 24px;
}

.option-group {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.option-group label {
  min-width: 60px;
  font-weight: 500;
  color: #374151;
}

.option-checkboxes {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.generator-actions {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.generated-passwords h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 16px 0;
}

.generated-item {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.generated-password {
  flex: 1;
  font-family: 'Monaco', 'Menlo', monospace;
}

.password-actions {
  display: flex;
  gap: 4px;
}

.generated-strength {
  min-width: 60px;
  font-size: 12px;
  font-weight: 600;
  text-align: center;
}

.security-tips {
  padding: 16px 0;
}

.tips-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
}

.tip-item h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.tip-item p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

@media (max-width: 1024px) {
  .tools-content {
    grid-template-columns: 1fr;
  }
  
  .generator-section {
    grid-column: auto;
  }
}

@media (max-width: 768px) {
  .password-strength-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .option-checkboxes {
    grid-template-columns: 1fr;
  }
  
  .generated-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .generated-password {
    width: 100%;
  }
  
  .tips-grid {
    grid-template-columns: 1fr;
  }
}
</style>