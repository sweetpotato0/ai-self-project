<template>
  <div class="regex-tester-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-others' })">其它</el-breadcrumb-item>
          <el-breadcrumb-item>正则表达式测试</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>正则表达式测试器</h1>
      <p>测试和验证正则表达式，支持实时匹配和详细解释</p>
    </div>

    <div class="tools-content">
      <!-- 正则输入区域 -->
      <div class="regex-section">
        <el-card class="regex-card">
          <template #header>
            <div class="card-header">
              <el-icon><Document /></el-icon>
              <span>正则表达式</span>
            </div>
          </template>
          
          <div class="regex-input">
            <div class="input-group">
              <label>正则表达式：</label>
              <div class="regex-input-wrapper">
                <span class="regex-delimiter">/</span>
                <el-input
                  v-model="regexPattern"
                  placeholder="输入正则表达式..."
                  @input="testRegex"
                  class="regex-field"
                />
                <span class="regex-delimiter">/</span>
                <el-input
                  v-model="regexFlags"
                  placeholder="flags"
                  @input="testRegex"
                  class="flags-field"
                  maxlength="5"
                />
              </div>
            </div>

            <div class="flags-options">
              <label>标志选项：</label>
              <div class="flag-checkboxes">
                <el-checkbox v-model="flags.global" @change="updateFlags">g (全局匹配)</el-checkbox>
                <el-checkbox v-model="flags.ignoreCase" @change="updateFlags">i (忽略大小写)</el-checkbox>
                <el-checkbox v-model="flags.multiline" @change="updateFlags">m (多行模式)</el-checkbox>
                <el-checkbox v-model="flags.dotAll" @change="updateFlags">s (dotAll模式)</el-checkbox>
                <el-checkbox v-model="flags.unicode" @change="updateFlags">u (Unicode模式)</el-checkbox>
              </div>
            </div>

            <div class="regex-info" v-if="regexError">
              <el-alert
                :title="regexError"
                type="error"
                show-icon
                :closable="false"
              />
            </div>

            <div class="regex-info" v-else-if="regexPattern">
              <div class="info-item">
                <span class="info-label">完整表达式：</span>
                <code class="regex-display">/{{ regexPattern }}/{{ regexFlags }}</code>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 测试文本区域 -->
      <div class="text-section">
        <el-card class="text-card">
          <template #header>
            <div class="card-header">
              <el-icon><Document /></el-icon>
              <span>测试文本</span>
              <div class="card-actions">
                <el-button @click="loadSampleText" size="small" text>
                  <el-icon><MagicStick /></el-icon>
                  样例文本
                </el-button>
              </div>
            </div>
          </template>
          
          <div class="text-input">
            <el-input
              v-model="testText"
              type="textarea"
              :rows="8"
              placeholder="输入要测试的文本..."
              @input="testRegex"
            />
            <div class="text-stats" v-if="testText">
              <span>字符数: {{ testText.length }}</span>
              <span>行数: {{ testText.split('\n').length }}</span>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 匹配结果区域 -->
      <div class="results-section">
        <el-card class="results-card">
          <template #header>
            <div class="card-header">
              <el-icon><CircleCheck /></el-icon>
              <span>匹配结果</span>
              <span v-if="matchResults.length > 0" class="match-count">
                找到 {{ matchResults.length }} 个匹配
              </span>
            </div>
          </template>
          
          <div class="results-content">
            <div v-if="!regexPattern || !testText" class="empty-state">
              <el-icon class="empty-icon"><Search /></el-icon>
              <p>请输入正则表达式和测试文本</p>
            </div>

            <div v-else-if="regexError" class="error-state">
              <el-icon class="error-icon"><WarningFilled /></el-icon>
              <p>正则表达式有误，请检查语法</p>
            </div>

            <div v-else-if="matchResults.length === 0" class="no-matches">
              <el-icon class="no-match-icon"><CircleCloseFilled /></el-icon>
              <p>没有找到匹配项</p>
            </div>

            <div v-else class="match-list">
              <div 
                v-for="(match, index) in matchResults" 
                :key="index"
                class="match-item"
              >
                <div class="match-header">
                  <span class="match-index">匹配 {{ index + 1 }}</span>
                  <span class="match-position">位置: {{ match.index }} - {{ match.index + match[0].length - 1 }}</span>
                </div>
                <div class="match-content">
                  <div class="match-text">
                    <strong>完整匹配:</strong>
                    <code class="match-value">{{ match[0] }}</code>
                  </div>
                  <div v-if="match.length > 1" class="match-groups">
                    <div 
                      v-for="(group, groupIndex) in match.slice(1)" 
                      :key="groupIndex"
                      class="group-item"
                    >
                      <strong>分组 {{ groupIndex + 1 }}:</strong>
                      <code class="group-value">{{ group || '(未匹配)' }}</code>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 高亮显示区域 -->
      <div class="highlight-section" v-if="regexPattern && testText && !regexError">
        <el-card class="highlight-card">
          <template #header>
            <div class="card-header">
              <el-icon><View /></el-icon>
              <span>高亮显示</span>
            </div>
          </template>
          
          <div class="highlighted-text" v-html="highlightedText"></div>
        </el-card>
      </div>
    </div>

    <!-- 替换功能 -->
    <div class="replace-section" v-if="regexPattern && testText && !regexError">
      <el-card class="replace-card">
        <template #header>
          <div class="card-header">
            <el-icon><Edit /></el-icon>
            <span>查找和替换</span>
          </div>
        </template>
        
        <div class="replace-content">
          <div class="replace-input">
            <div class="input-group">
              <label>替换为：</label>
              <el-input
                v-model="replaceText"
                placeholder="输入替换文本..."
                @input="performReplace"
              />
            </div>
          </div>
          
          <div class="replace-result" v-if="replacedText !== testText">
            <h4>替换结果：</h4>
            <div class="result-text">{{ replacedText }}</div>
            <div class="replace-actions">
              <el-button @click="copyReplacedText" size="small" type="primary">
                <el-icon><DocumentCopy /></el-icon>
                复制结果
              </el-button>
              <el-button @click="useReplacedText" size="small">
                <el-icon><RefreshLeft /></el-icon>
                应用到测试文本
              </el-button>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 常用正则表达式 -->
    <div class="patterns-section">
      <el-card class="patterns-card">
        <template #header>
          <div class="card-header">
            <el-icon><CollectionTag /></el-icon>
            <span>常用正则表达式</span>
          </div>
        </template>
        
        <div class="patterns-grid">
          <div 
            v-for="pattern in commonPatterns" 
            :key="pattern.name"
            class="pattern-item"
            @click="usePattern(pattern)"
          >
            <div class="pattern-header">
              <span class="pattern-name">{{ pattern.name }}</span>
            </div>
            <div class="pattern-regex">
              <code>{{ pattern.regex }}</code>
            </div>
            <div class="pattern-desc">{{ pattern.description }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 正则语法参考 -->
    <div class="reference-section">
      <el-card class="reference-card">
        <template #header>
          <div class="card-header">
            <el-icon><InfoFilled /></el-icon>
            <span>语法参考</span>
          </div>
        </template>
        
        <div class="reference-content">
          <div class="reference-grid">
            <div class="reference-group">
              <h4>基础语法</h4>
              <ul>
                <li><code>.</code> - 任意字符（除换行符）</li>
                <li><code>*</code> - 0次或多次</li>
                <li><code>+</code> - 1次或多次</li>
                <li><code>?</code> - 0次或1次</li>
                <li><code>{n}</code> - 恰好n次</li>
                <li><code>{n,m}</code> - n到m次</li>
              </ul>
            </div>
            <div class="reference-group">
              <h4>字符类</h4>
              <ul>
                <li><code>[abc]</code> - 字符集合</li>
                <li><code>[^abc]</code> - 非字符集合</li>
                <li><code>\d</code> - 数字 [0-9]</li>
                <li><code>\w</code> - 单词字符 [a-zA-Z0-9_]</li>
                <li><code>\s</code> - 空白字符</li>
                <li><code>\D \W \S</code> - 相应的非字符</li>
              </ul>
            </div>
            <div class="reference-group">
              <h4>位置锚点</h4>
              <ul>
                <li><code>^</code> - 行开始</li>
                <li><code>$</code> - 行结束</li>
                <li><code>\b</code> - 单词边界</li>
                <li><code>\B</code> - 非单词边界</li>
              </ul>
            </div>
            <div class="reference-group">
              <h4>分组和引用</h4>
              <ul>
                <li><code>()</code> - 捕获分组</li>
                <li><code>(?:)</code> - 非捕获分组</li>
                <li><code>(?=)</code> - 正向前瞻</li>
                <li><code>(?!)</code> - 负向前瞻</li>
                <li><code>\1</code> - 反向引用</li>
              </ul>
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
  Document,
  CircleCheck,
  Search,
  WarningFilled,
  CircleCloseFilled,
  View,
  Edit,
  DocumentCopy,
  RefreshLeft,
  CollectionTag,
  InfoFilled,
  MagicStick
} from '@element-plus/icons-vue'

const router = useRouter()

// 响应式数据
const regexPattern = ref('')
const regexFlags = ref('')
const testText = ref('')
const replaceText = ref('')
const regexError = ref('')

const flags = ref({
  global: false,
  ignoreCase: false,
  multiline: false,
  dotAll: false,
  unicode: false
})

// 常用正则表达式
const commonPatterns = ref([
  {
    name: '邮箱地址',
    regex: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$',
    description: '匹配常见的邮箱地址格式'
  },
  {
    name: '手机号码',
    regex: '^1[3-9]\\d{9}$',
    description: '匹配中国大陆手机号码'
  },
  {
    name: 'URL地址',
    regex: 'https?:\\/\\/[^\\s]+',
    description: '匹配HTTP或HTTPS网址'
  },
  {
    name: 'IP地址',
    regex: '^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$',
    description: '匹配IPv4地址'
  },
  {
    name: '身份证号',
    regex: '^\\d{15}|\\d{17}[\\dXx]$',
    description: '匹配15位或18位身份证号'
  },
  {
    name: '中文字符',
    regex: '[\\u4e00-\\u9fa5]',
    description: '匹配中文汉字'
  },
  {
    name: '数字',
    regex: '^-?\\d+(\\.\\d+)?$',
    description: '匹配整数或小数'
  },
  {
    name: '日期格式',
    regex: '^\\d{4}-\\d{2}-\\d{2}$',
    description: '匹配YYYY-MM-DD格式日期'
  }
])

// 匹配结果
const matchResults = computed(() => {
  if (!regexPattern.value || !testText.value || regexError.value) {
    return []
  }

  try {
    const regex = new RegExp(regexPattern.value, regexFlags.value)
    const results = []
    
    if (flags.value.global) {
      let match
      while ((match = regex.exec(testText.value)) !== null) {
        results.push(match)
        if (!flags.value.global) break
      }
    } else {
      const match = regex.exec(testText.value)
      if (match) results.push(match)
    }
    
    return results
  } catch (e) {
    return []
  }
})

// 高亮显示的文本
const highlightedText = computed(() => {
  if (!regexPattern.value || !testText.value || regexError.value) {
    return testText.value
  }

  try {
    const regex = new RegExp(regexPattern.value, regexFlags.value)
    let result = testText.value
    let offset = 0
    
    const matches = matchResults.value
    matches.forEach((match) => {
      const start = match.index + offset
      const end = start + match[0].length
      const highlighted = `<mark class="regex-match">${match[0]}</mark>`
      result = result.slice(0, start) + highlighted + result.slice(end)
      offset += highlighted.length - match[0].length
    })
    
    return result.replace(/\n/g, '<br>')
  } catch (e) {
    return testText.value
  }
})

// 替换后的文本
const replacedText = computed(() => {
  if (!regexPattern.value || !testText.value || regexError.value || replaceText.value === '') {
    return testText.value
  }

  try {
    const regex = new RegExp(regexPattern.value, regexFlags.value)
    return testText.value.replace(regex, replaceText.value)
  } catch (e) {
    return testText.value
  }
})

// 更新标志
const updateFlags = () => {
  let flagStr = ''
  if (flags.value.global) flagStr += 'g'
  if (flags.value.ignoreCase) flagStr += 'i'
  if (flags.value.multiline) flagStr += 'm'
  if (flags.value.dotAll) flagStr += 's'
  if (flags.value.unicode) flagStr += 'u'
  
  regexFlags.value = flagStr
  testRegex()
}

// 测试正则表达式
const testRegex = () => {
  regexError.value = ''
  
  if (!regexPattern.value) return
  
  try {
    new RegExp(regexPattern.value, regexFlags.value)
  } catch (e) {
    regexError.value = `语法错误: ${e.message}`
  }
}

// 使用模式
const usePattern = (pattern) => {
  regexPattern.value = pattern.regex
  testRegex()
  ElMessage.success(`已应用模式: ${pattern.name}`)
}

// 加载样例文本
const loadSampleText = () => {
  testText.value = `示例文本：
邮箱地址：user@example.com, test@gmail.com
手机号码：13812345678, 18966778899
网址：https://www.example.com, http://test.org
日期：2023-12-25, 2024-01-01
中文内容：这是一段中文测试文本
数字：123, 45.67, -89.12
特殊字符：@#$%^&*()
`
  ElMessage.success('已加载样例文本')
}

// 执行替换
const performReplace = () => {
  // 替换逻辑在计算属性中处理
}

// 复制替换结果
const copyReplacedText = async () => {
  try {
    await navigator.clipboard.writeText(replacedText.value)
    ElMessage.success('替换结果已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 应用替换结果
const useReplacedText = () => {
  testText.value = replacedText.value
  replaceText.value = ''
  ElMessage.success('已应用替换结果到测试文本')
}

// 监听flags输入的变化
const watchFlagsInput = () => {
  const flagStr = regexFlags.value
  flags.value = {
    global: flagStr.includes('g'),
    ignoreCase: flagStr.includes('i'),
    multiline: flagStr.includes('m'),
    dotAll: flagStr.includes('s'),
    unicode: flagStr.includes('u')
  }
}

// 监听regexFlags的变化
const watchRegexFlags = () => {
  watchFlagsInput()
  testRegex()
}
</script>

<style scoped>
.regex-tester-container {
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
  gap: 24px;
  margin-bottom: 32px;
}

.highlight-section {
  grid-column: span 2;
}

.regex-card,
.text-card,
.results-card,
.highlight-card,
.replace-card,
.patterns-card,
.reference-card {
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

.match-count {
  font-size: 14px;
  color: #22c55e;
  font-weight: normal;
}

.regex-input {
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

.regex-input-wrapper {
  display: flex;
  align-items: center;
  gap: 0;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.regex-delimiter {
  padding: 0 8px;
  background: #f5f7fa;
  color: #909399;
  font-family: 'Monaco', 'Menlo', monospace;
  border-right: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  height: 32px;
}

.regex-field {
  flex: 1;
  border: none;
}

.regex-field :deep(.el-input__wrapper) {
  border: none;
  box-shadow: none;
}

.flags-field {
  width: 60px;
  border: none;
}

.flags-field :deep(.el-input__wrapper) {
  border: none;
  box-shadow: none;
}

.flags-options {
  margin-bottom: 20px;
}

.flags-options label {
  display: block;
  margin-bottom: 12px;
  font-weight: 500;
  color: #374151;
}

.flag-checkboxes {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.regex-info {
  margin-top: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.info-label {
  color: #6b7280;
}

.regex-display {
  background: #f1f5f9;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
  color: #1f2937;
}

.text-input {
  padding: 8px 0;
}

.text-stats {
  margin-top: 8px;
  font-size: 12px;
  color: #6b7280;
  text-align: right;
  display: flex;
  justify-content: space-between;
}

.results-content {
  padding: 8px 0;
  min-height: 200px;
}

.empty-state,
.error-state,
.no-matches {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 150px;
  color: #9ca3af;
}

.empty-icon,
.error-icon,
.no-match-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.error-state {
  color: #ef4444;
}

.match-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.match-item {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  background: #f8fafc;
}

.match-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 14px;
}

.match-index {
  font-weight: 600;
  color: #2c3e50;
}

.match-position {
  color: #6b7280;
  font-family: 'Monaco', 'Menlo', monospace;
}

.match-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.match-text,
.group-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.match-value,
.group-value {
  background: #fff;
  border: 1px solid #d1d5db;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
}

.match-groups {
  margin-left: 16px;
  padding-left: 16px;
  border-left: 2px solid #d1d5db;
}

.highlighted-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 14px;
  line-height: 1.6;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  max-height: 300px;
  overflow-y: auto;
}

.highlighted-text :deep(.regex-match) {
  background: #fbbf24;
  padding: 2px 4px;
  border-radius: 2px;
  font-weight: 600;
}

.replace-section {
  margin-bottom: 32px;
}

.replace-content {
  padding: 8px 0;
}

.replace-result h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 16px 0 12px 0;
}

.result-text {
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 14px;
  line-height: 1.6;
  max-height: 200px;
  overflow-y: auto;
  margin-bottom: 16px;
}

.replace-actions {
  display: flex;
  gap: 12px;
}

.patterns-section,
.reference-section {
  margin-bottom: 32px;
}

.patterns-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.pattern-item {
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: #fff;
}

.pattern-item:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.pattern-header {
  margin-bottom: 8px;
}

.pattern-name {
  font-weight: 600;
  color: #2c3e50;
}

.pattern-regex {
  margin-bottom: 8px;
}

.pattern-regex code {
  background: #f1f5f9;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  word-break: break-all;
}

.pattern-desc {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.4;
}

.reference-content {
  padding: 16px 0;
}

.reference-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
}

.reference-group h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 12px 0;
}

.reference-group ul {
  margin: 0;
  padding-left: 20px;
}

.reference-group li {
  margin: 6px 0;
  font-size: 14px;
  color: #374151;
  line-height: 1.4;
}

.reference-group code {
  background: #f1f5f9;
  padding: 2px 4px;
  border-radius: 2px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
  color: #1f2937;
}

@media (max-width: 1024px) {
  .tools-content {
    grid-template-columns: 1fr;
  }
  
  .highlight-section {
    grid-column: auto;
  }
}

@media (max-width: 768px) {
  .regex-tester-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .flag-checkboxes {
    grid-template-columns: 1fr;
  }
  
  .patterns-grid {
    grid-template-columns: 1fr;
  }
  
  .reference-grid {
    grid-template-columns: 1fr;
  }
  
  .replace-actions {
    flex-direction: column;
  }
}
</style>