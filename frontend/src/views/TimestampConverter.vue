<template>
  <div class="timestamp-converter-container">
    <div class="tool-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-development' })">开发类</el-breadcrumb-item>
          <el-breadcrumb-item>时间戳转换</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>时间戳转换工具</h1>
      <p>时间戳与日期时间的相互转换工具</p>
    </div>

    <!-- 当前时间戳实时显示 -->
    <div class="current-timestamp-section">
      <div class="current-timestamp-card">
        <div class="timestamp-display">
          <div class="timestamp-label">当前时间戳</div>
          <div class="timestamp-value">{{ currentTimestamp }}</div>
          <div class="timestamp-human">{{ currentHumanTime }}</div>
        </div>
        <div class="timestamp-formats">
          <div class="format-item">
            <span class="format-label">秒级:</span>
            <span class="format-value">{{ Math.floor(currentTimestamp / 1000) }}</span>
          </div>
          <div class="format-item">
            <span class="format-label">毫秒级:</span>
            <span class="format-value">{{ currentTimestamp }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="converter-sections">
      <!-- 时间戳转时间 -->
      <div class="converter-section">
        <h3>时间戳转时间</h3>
        <div class="converter-card">
          <div class="input-group">
            <el-input
              v-model="timestampInput"
              placeholder="请输入时间戳（支持10位或13位）"
              clearable
              @input="convertTimestampToDate"
            >
              <template #prepend>时间戳</template>
            </el-input>
          </div>
          
          <div v-if="timestampResult" class="result-group">
            <div class="result-item">
              <span class="result-label">标准格式:</span>
              <span class="result-value">{{ timestampResult.standard }}</span>
              <el-button 
                type="text" 
                @click="copyToClipboard(timestampResult.standard)"
                class="copy-btn"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
            <div class="result-item">
              <span class="result-label">ISO格式:</span>
              <span class="result-value">{{ timestampResult.iso }}</span>
              <el-button 
                type="text" 
                @click="copyToClipboard(timestampResult.iso)"
                class="copy-btn"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
            <div class="result-item">
              <span class="result-label">本地时间:</span>
              <span class="result-value">{{ timestampResult.locale }}</span>
              <el-button 
                type="text" 
                @click="copyToClipboard(timestampResult.locale)"
                class="copy-btn"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
            <div class="result-item">
              <span class="result-label">相对时间:</span>
              <span class="result-value">{{ timestampResult.relative }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 时间转时间戳 -->
      <div class="converter-section">
        <h3>时间转时间戳</h3>
        <div class="converter-card">
          <div class="input-group">
            <el-date-picker
              v-model="dateInput"
              type="datetime"
              placeholder="选择日期时间"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              @change="convertDateToTimestamp"
              style="width: 100%"
            />
          </div>
          
          <div v-if="dateResult" class="result-group">
            <div class="result-item">
              <span class="result-label">秒级时间戳:</span>
              <span class="result-value">{{ dateResult.seconds }}</span>
              <el-button 
                type="text" 
                @click="copyToClipboard(dateResult.seconds)"
                class="copy-btn"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
            <div class="result-item">
              <span class="result-label">毫秒级时间戳:</span>
              <span class="result-value">{{ dateResult.milliseconds }}</span>
              <el-button 
                type="text" 
                @click="copyToClipboard(dateResult.milliseconds)"
                class="copy-btn"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 常用时间戳 -->
    <div class="common-timestamps">
      <h3>常用时间戳</h3>
      <div class="timestamp-grid">
        <div 
          v-for="common in commonTimestamps" 
          :key="common.label"
          class="timestamp-item"
          @click="useCommonTimestamp(common.timestamp)"
        >
          <div class="common-label">{{ common.label }}</div>
          <div class="common-timestamp">{{ common.timestamp }}</div>
          <div class="common-date">{{ formatTimestamp(common.timestamp) }}</div>
        </div>
      </div>
    </div>

    <!-- 编程语言代码示例 -->
    <div class="code-examples-section">
      <h3>编程语言时间处理示例</h3>
      <p class="section-description">各种编程语言中获取和转换时间戳的代码示例</p>
      
      <div class="code-examples-grid">
        <div 
          v-for="example in codeExamples" 
          :key="example.language"
          class="code-example-card"
        >
          <div class="code-header">
            <div class="language-info">
              <div class="language-icon" :style="{ backgroundColor: example.color }">
                <span class="language-abbr">{{ example.abbr }}</span>
              </div>
              <div class="language-name">{{ example.language }}</div>
            </div>
            <el-button 
              type="text" 
              size="small" 
              @click="copyCode(example.code)"
              class="copy-code-btn"
            >
              <el-icon><CopyDocument /></el-icon>
              复制
            </el-button>
          </div>
          
          <div class="code-content">
            <div class="code-section">
              <div class="code-label">获取当前时间戳:</div>
              <pre class="code-block">{{ example.getCurrentTimestamp }}</pre>
            </div>
            
            <div class="code-section">
              <div class="code-label">时间戳转时间:</div>
              <pre class="code-block">{{ example.timestampToDate }}</pre>
            </div>
            
            <div class="code-section">
              <div class="code-label">时间转时间戳:</div>
              <pre class="code-block">{{ example.dateToTimestamp }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()

// 响应式数据
const currentTimestamp = ref(Date.now())
const currentHumanTime = ref('')
const timestampInput = ref('')
const timestampResult = ref(null)
const dateInput = ref('')
const dateResult = ref(null)

// 定时器
let timer = null

// 常用时间戳
const commonTimestamps = ref([
  { label: '当前时间', timestamp: 0 },
  { label: '今天开始', timestamp: 0 },
  { label: '今天结束', timestamp: 0 },
  { label: '昨天开始', timestamp: 0 },
  { label: '一周前', timestamp: 0 },
  { label: '一月前', timestamp: 0 }
])

// 编程语言代码示例
const codeExamples = ref([
  {
    language: 'JavaScript',
    abbr: 'JS',
    color: '#f7df1e',
    getCurrentTimestamp: `// 毫秒级时间戳\nDate.now()\n// 或\nnew Date().getTime()\n\n// 秒级时间戳\nMath.floor(Date.now() / 1000)`,
    timestampToDate: `// 毫秒级时间戳转换\nconst date = new Date(1692960000000)\nconsole.log(date.toLocaleString())\n\n// 秒级时间戳转换\nconst date = new Date(1692960000 * 1000)\nconsole.log(date.toISOString())`,
    dateToTimestamp: `// 当前时间转时间戳\nconst timestamp = new Date().getTime()\n\n// 指定时间转时间戳\nconst timestamp = new Date('2023-08-25 12:00:00').getTime()`
  },
  {
    language: 'Python',
    abbr: 'PY',
    color: '#3776ab',
    getCurrentTimestamp: `import time\n\n# 秒级时间戳\ntime.time()\n\n# 毫秒级时间戳\nint(time.time() * 1000)`,
    timestampToDate: `import datetime\n\n# 秒级时间戳转换\ntimestamp = 1692960000\ndt = datetime.datetime.fromtimestamp(timestamp)\nprint(dt.strftime('%Y-%m-%d %H:%M:%S'))`,
    dateToTimestamp: `import datetime\n\n# 当前时间转时间戳\ntimestamp = int(datetime.datetime.now().timestamp())\n\n# 指定时间转时间戳\ndt = datetime.datetime(2023, 8, 25, 12, 0, 0)\ntimestamp = int(dt.timestamp())`
  },
  {
    language: 'Java',
    abbr: 'JAVA',
    color: '#ed8b00',
    getCurrentTimestamp: `// 毫秒级时间戳\nlong timestamp = System.currentTimeMillis();\n\n// 秒级时间戳\nlong timestamp = System.currentTimeMillis() / 1000;`,
    timestampToDate: `import java.time.*;\n\n// 秒级时间戳转换\nlong timestamp = 1692960000;\nLocalDateTime dateTime = LocalDateTime.ofEpochSecond(\n    timestamp, 0, ZoneOffset.UTC);\nSystem.out.println(dateTime);`,
    dateToTimestamp: `import java.time.*;\n\n// 当前时间转时间戳\nlong timestamp = Instant.now().getEpochSecond();\n\n// 指定时间转时间戳\nLocalDateTime dateTime = LocalDateTime.of(2023, 8, 25, 12, 0);\nlong timestamp = dateTime.toEpochSecond(ZoneOffset.UTC);`
  },
  {
    language: 'Go',
    abbr: 'GO',
    color: '#00add8',
    getCurrentTimestamp: `package main\n\nimport "time"\n\n// 秒级时间戳\ntimestamp := time.Now().Unix()\n\n// 毫秒级时间戳\ntimestamp := time.Now().UnixMilli()`,
    timestampToDate: `package main\n\nimport (\n    "fmt"\n    "time"\n)\n\n// 秒级时间戳转换\ntimestamp := int64(1692960000)\nt := time.Unix(timestamp, 0)\nfmt.Println(t.Format("2006-01-02 15:04:05"))`,
    dateToTimestamp: `package main\n\nimport "time"\n\n// 当前时间转时间戳\ntimestamp := time.Now().Unix()\n\n// 指定时间转时间戳\nt, _ := time.Parse("2006-01-02 15:04:05", "2023-08-25 12:00:00")\ntimestamp := t.Unix()`
  },
  {
    language: 'Rust',
    abbr: 'RS',
    color: '#dea584',
    getCurrentTimestamp: `use std::time::{SystemTime, UNIX_EPOCH};\n\n// 秒级时间戳\nlet timestamp = SystemTime::now()\n    .duration_since(UNIX_EPOCH)\n    .unwrap()\n    .as_secs();\n\n// 毫秒级时间戳\nlet timestamp = SystemTime::now()\n    .duration_since(UNIX_EPOCH)\n    .unwrap()\n    .as_millis();`,
    timestampToDate: `use chrono::{DateTime, Utc, NaiveDateTime};\n\n// 时间戳转日期\nlet timestamp = 1692960000;\nlet naive = NaiveDateTime::from_timestamp(timestamp, 0).unwrap();\nlet datetime: DateTime<Utc> = DateTime::from_utc(naive, Utc);\nprintln!("{}", datetime.format("%Y-%m-%d %H:%M:%S"));`,
    dateToTimestamp: `use chrono::{Utc, TimeZone};\n\n// 当前时间转时间戳\nlet timestamp = Utc::now().timestamp();\n\n// 指定时间转时间戳\nlet dt = Utc.ymd(2023, 8, 25).and_hms(12, 0, 0);\nlet timestamp = dt.timestamp();`
  }
])

// 更新当前时间戳
const updateCurrentTimestamp = () => {
  currentTimestamp.value = Date.now()
  currentHumanTime.value = dayjs(currentTimestamp.value).format('YYYY-MM-DD HH:mm:ss')
  
  // 更新常用时间戳
  const now = dayjs()
  commonTimestamps.value = [
    { label: '当前时间', timestamp: now.valueOf() },
    { label: '今天开始', timestamp: now.startOf('day').valueOf() },
    { label: '今天结束', timestamp: now.endOf('day').valueOf() },
    { label: '昨天开始', timestamp: now.subtract(1, 'day').startOf('day').valueOf() },
    { label: '一周前', timestamp: now.subtract(1, 'week').valueOf() },
    { label: '一月前', timestamp: now.subtract(1, 'month').valueOf() }
  ]
}

// 时间戳转时间
const convertTimestampToDate = () => {
  if (!timestampInput.value) {
    timestampResult.value = null
    return
  }
  
  try {
    let timestamp = parseInt(timestampInput.value)
    
    // 处理10位时间戳（秒级）
    if (timestamp.toString().length === 10) {
      timestamp = timestamp * 1000
    }
    
    if (timestamp.toString().length !== 13) {
      timestampResult.value = null
      return
    }
    
    const date = dayjs(timestamp)
    
    if (!date.isValid()) {
      timestampResult.value = null
      return
    }
    
    timestampResult.value = {
      standard: date.format('YYYY-MM-DD HH:mm:ss'),
      iso: date.toISOString(),
      locale: date.format('YYYY年MM月DD日 HH:mm:ss'),
      relative: date.fromNow()
    }
  } catch (error) {
    timestampResult.value = null
  }
}

// 时间转时间戳
const convertDateToTimestamp = () => {
  if (!dateInput.value) {
    dateResult.value = null
    return
  }
  
  try {
    const date = dayjs(dateInput.value)
    
    if (!date.isValid()) {
      dateResult.value = null
      return
    }
    
    const timestamp = date.valueOf()
    
    dateResult.value = {
      seconds: Math.floor(timestamp / 1000),
      milliseconds: timestamp
    }
  } catch (error) {
    dateResult.value = null
  }
}

// 格式化时间戳显示
const formatTimestamp = (timestamp) => {
  return dayjs(timestamp).format('MM-DD HH:mm')
}

// 使用常用时间戳
const useCommonTimestamp = (timestamp) => {
  timestampInput.value = timestamp.toString()
  convertTimestampToDate()
}

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text.toString())
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 复制代码到剪贴板
const copyCode = async (code) => {
  const fullCode = `${code.getCurrentTimestamp}\n\n${code.timestampToDate}\n\n${code.dateToTimestamp}`
  try {
    await navigator.clipboard.writeText(fullCode)
    ElMessage.success('代码已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 生命周期
onMounted(() => {
  updateCurrentTimestamp()
  timer = setInterval(updateCurrentTimestamp, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.timestamp-converter-container {
  padding: 20px 0;
  max-width: 1200px;
  margin: 0 auto;
}

.tool-header {
  margin-bottom: 30px;
  text-align: left;
}

.tool-header h1,
.tool-header p {
  text-align: center;
}

.breadcrumb {
  margin-bottom: 16px;
}

.breadcrumb .el-breadcrumb-item {
  cursor: pointer;
}

.tool-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.tool-header p {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

.current-timestamp-section {
  margin-bottom: 40px;
}

.current-timestamp-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 16px;
  padding: 32px;
  text-align: center;
}

.timestamp-display {
  margin-bottom: 24px;
}

.timestamp-label {
  font-size: 16px;
  opacity: 0.9;
  margin-bottom: 8px;
}

.timestamp-value {
  font-size: 36px;
  font-weight: 600;
  font-family: 'Monaco', 'Menlo', monospace;
  margin-bottom: 8px;
  letter-spacing: 1px;
}

.timestamp-human {
  font-size: 18px;
  opacity: 0.9;
}

.timestamp-formats {
  display: flex;
  justify-content: center;
  gap: 40px;
}

.format-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.format-label {
  font-size: 14px;
  opacity: 0.8;
}

.format-value {
  font-size: 16px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-weight: 600;
}

.converter-sections {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: 32px;
  margin-bottom: 40px;
}

.converter-section h3 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 16px;
}

.converter-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.input-group {
  margin-bottom: 24px;
}

.result-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.result-item {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.result-label {
  font-weight: 500;
  color: #374151;
  min-width: 100px;
  flex-shrink: 0;
}

.result-value {
  flex: 1;
  font-family: 'Monaco', 'Menlo', monospace;
  color: #1f2937;
  margin-left: 12px;
}

.copy-btn {
  margin-left: 8px;
  padding: 4px;
}

.common-timestamps h3 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 16px;
}

.timestamp-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.timestamp-item {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: center;
}

.timestamp-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-color: #409eff;
}

.common-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.common-timestamp {
  font-size: 12px;
  font-family: 'Monaco', 'Menlo', monospace;
  color: #6b7280;
  margin-bottom: 4px;
}

.common-date {
  font-size: 12px;
  color: #9ca3af;
}

/* 编程语言代码示例样式 */
.code-examples-section {
  margin-top: 40px;
}

.code-examples-section h3 {
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.section-description {
  font-size: 16px;
  color: #6b7280;
  margin-bottom: 24px;
}

.code-examples-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(450px, 1fr));
  gap: 24px;
}

.code-example-card {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.code-example-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 25px -5px rgba(0, 0, 0, 0.15);
  border-color: #409eff;
}

.code-header {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  padding: 20px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e2e8f0;
}

.language-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.language-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.language-abbr {
  color: white;
  font-weight: 700;
  font-size: 14px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.language-name {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.copy-code-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  background: rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.2);
  color: #409eff;
  font-weight: 500;
  transition: all 0.2s ease;
}

.copy-code-btn:hover {
  background: rgba(64, 158, 255, 0.15);
  transform: translateY(-1px);
}

.code-content {
  padding: 0;
}

.code-section {
  padding: 24px;
  border-bottom: 1px solid #f1f5f9;
}

.code-section:last-child {
  border-bottom: none;
}

.code-label {
  font-size: 14px;
  font-weight: 600;
  color: #4b5563;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.code-label::before {
  content: "▶";
  color: #10b981;
  font-size: 12px;
}

.code-block {
  background: #0f172a;
  color: #e2e8f0;
  padding: 16px 20px;
  border-radius: 12px;
  font-family: 'Fira Code', 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  line-height: 1.6;
  margin: 0;
  overflow-x: auto;
  border: 1px solid #334155;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
}

.code-block::-webkit-scrollbar {
  height: 6px;
}

.code-block::-webkit-scrollbar-track {
  background: #1e293b;
  border-radius: 3px;
}

.code-block::-webkit-scrollbar-thumb {
  background: #475569;
  border-radius: 3px;
}

.code-block::-webkit-scrollbar-thumb:hover {
  background: #64748b;
}

@media (max-width: 768px) {
  .converter-sections {
    grid-template-columns: 1fr;
    gap: 24px;
  }
  
  .timestamp-formats {
    flex-direction: column;
    gap: 20px;
  }
  
  .timestamp-value {
    font-size: 28px;
  }
  
  .timestamp-grid {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 12px;
  }
  
  .current-timestamp-card {
    padding: 24px;
  }
  
  .code-examples-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .code-example-card {
    margin: 0 -8px;
  }
  
  .code-header {
    padding: 16px 20px;
  }
  
  .language-icon {
    width: 40px;
    height: 40px;
  }
  
  .language-name {
    font-size: 16px;
  }
  
  .code-section {
    padding: 20px;
  }
  
  .code-block {
    padding: 12px 16px;
    font-size: 12px;
  }
}
</style>