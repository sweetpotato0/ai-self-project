<template>
  <div class="word-counter-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-text' })">文本处理</el-breadcrumb-item>
          <el-breadcrumb-item>字符统计</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>字符统计工具</h1>
      <p>统计文本的字符数、单词数、行数等信息</p>
    </div>

    <div class="tools-content">
      <div class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Edit /></el-icon>
              <span>文本输入</span>
            </div>
          </template>
          <div class="input-content">
            <el-input
              v-model="inputText"
              type="textarea"
              :rows="12"
              placeholder="请输入要统计的文本..."
              @input="updateStats"
            />
            <div class="input-actions">
              <el-button @click="clearText">
                <el-icon><Delete /></el-icon>
                清空文本
              </el-button>
              <el-button @click="pasteFromClipboard">
                <el-icon><DocumentCopy /></el-icon>
                粘贴文本
              </el-button>
              <el-button @click="loadSample">
                <el-icon><Document /></el-icon>
                加载示例
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div class="stats-section">
        <el-card class="stats-card">
          <template #header>
            <div class="card-header">
              <el-icon><DataAnalysis /></el-icon>
              <span>统计结果</span>
            </div>
          </template>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-label">总字符数</div>
              <div class="stat-value">{{ stats.totalChars.toLocaleString() }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">字符数(不含空格)</div>
              <div class="stat-value">{{ stats.charsNoSpaces.toLocaleString() }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">中文字符数</div>
              <div class="stat-value">{{ stats.chineseChars.toLocaleString() }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">英文字符数</div>
              <div class="stat-value">{{ stats.englishChars.toLocaleString() }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">数字字符数</div>
              <div class="stat-value">{{ stats.numberChars.toLocaleString() }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">标点符号数</div>
              <div class="stat-value">{{ stats.punctuationChars.toLocaleString() }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">单词数</div>
              <div class="stat-value">{{ stats.words.toLocaleString() }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">行数</div>
              <div class="stat-value">{{ stats.lines.toLocaleString() }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">段落数</div>
              <div class="stat-value">{{ stats.paragraphs.toLocaleString() }}</div>
            </div>
          </div>
        </el-card>
      </div>

      <div class="details-section">
        <div class="details-grid">
          <el-card class="detail-card">
            <template #header>
              <div class="card-header">
                <el-icon><List /></el-icon>
                <span>词频统计</span>
              </div>
            </template>
            <div class="word-frequency">
              <div v-if="wordFrequency.length === 0" class="empty-state">
                暂无词频数据
              </div>
              <div v-else class="frequency-list">
                <div
                  v-for="(item, index) in wordFrequency.slice(0, 20)"
                  :key="index"
                  class="frequency-item"
                >
                  <span class="word">{{ item.word }}</span>
                  <span class="count">{{ item.count }}次</span>
                  <div class="frequency-bar">
                    <div 
                      class="frequency-fill"
                      :style="{ width: (item.count / wordFrequency[0].count * 100) + '%' }"
                    />
                  </div>
                </div>
              </div>
            </div>
          </el-card>

          <el-card class="detail-card">
            <template #header>
              <div class="card-header">
                <el-icon><Timer /></el-icon>
                <span>阅读信息</span>
              </div>
            </template>
            <div class="reading-stats">
              <div class="reading-item">
                <span class="reading-label">预计阅读时间:</span>
                <span class="reading-value">{{ readingTime }}</span>
              </div>
              <div class="reading-item">
                <span class="reading-label">预计说话时间:</span>
                <span class="reading-value">{{ speakingTime }}</span>
              </div>
              <div class="reading-item">
                <span class="reading-label">平均字符/行:</span>
                <span class="reading-value">{{ averageCharsPerLine }}</span>
              </div>
              <div class="reading-item">
                <span class="reading-label">平均单词/句:</span>
                <span class="reading-value">{{ averageWordsPerSentence }}</span>
              </div>
            </div>
          </el-card>
        </div>
      </div>

      <div class="action-section">
        <el-button @click="copyStats" type="success">
          <el-icon><DocumentCopy /></el-icon>
          复制统计结果
        </el-button>
        <el-button @click="exportStats" type="primary">
          <el-icon><Download /></el-icon>
          导出详细报告
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@/composables/useClipboard'
import {
  Edit,
  Delete,
  DocumentCopy,
  Document,
  DataAnalysis,
  List,
  Timer,
  Download
} from '@element-plus/icons-vue'

const router = useRouter()
const { copy } = useClipboard()

const inputText = ref('')
const stats = ref({
  totalChars: 0,
  charsNoSpaces: 0,
  chineseChars: 0,
  englishChars: 0,
  numberChars: 0,
  punctuationChars: 0,
  words: 0,
  lines: 0,
  paragraphs: 0,
  sentences: 0
})

const wordFrequency = computed(() => {
  if (!inputText.value.trim()) return []
  
  const words = inputText.value
    .toLowerCase()
    .replace(/[^\w\u4e00-\u9fa5\s]/g, ' ')
    .split(/\s+/)
    .filter(word => word.length > 1)
  
  const frequency = {}
  words.forEach(word => {
    frequency[word] = (frequency[word] || 0) + 1
  })
  
  return Object.entries(frequency)
    .map(([word, count]) => ({ word, count }))
    .sort((a, b) => b.count - a.count)
})

const readingTime = computed(() => {
  const wordsPerMinute = 200 // 平均阅读速度
  const minutes = Math.ceil(stats.value.words / wordsPerMinute)
  if (minutes < 1) return '不到1分钟'
  if (minutes === 1) return '1分钟'
  if (minutes < 60) return `${minutes}分钟`
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  return `${hours}小时${remainingMinutes}分钟`
})

const speakingTime = computed(() => {
  const wordsPerMinute = 150 // 平均说话速度
  const minutes = Math.ceil(stats.value.words / wordsPerMinute)
  if (minutes < 1) return '不到1分钟'
  if (minutes === 1) return '1分钟'
  if (minutes < 60) return `${minutes}分钟`
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  return `${hours}小时${remainingMinutes}分钟`
})

const averageCharsPerLine = computed(() => {
  if (stats.value.lines === 0) return '0'
  return Math.round(stats.value.totalChars / stats.value.lines)
})

const averageWordsPerSentence = computed(() => {
  if (stats.value.sentences === 0) return '0'
  return Math.round(stats.value.words / stats.value.sentences)
})

const updateStats = () => {
  const text = inputText.value
  
  stats.value = {
    totalChars: text.length,
    charsNoSpaces: text.replace(/\s/g, '').length,
    chineseChars: (text.match(/[\u4e00-\u9fa5]/g) || []).length,
    englishChars: (text.match(/[a-zA-Z]/g) || []).length,
    numberChars: (text.match(/[0-9]/g) || []).length,
    punctuationChars: (text.match(/[^\w\s\u4e00-\u9fa5]/g) || []).length,
    words: text.trim() ? text.trim().split(/\s+/).length : 0,
    lines: text ? text.split('\n').length : 0,
    paragraphs: text.trim() ? text.trim().split(/\n\s*\n/).length : 0,
    sentences: (text.match(/[.!?。！？]/g) || []).length
  }
}

const clearText = () => {
  inputText.value = ''
  updateStats()
}

const pasteFromClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    inputText.value = text
    updateStats()
    ElMessage.success('文本已粘贴')
  } catch (err) {
    ElMessage.error('粘贴失败：' + err.message)
  }
}

const loadSample = () => {
  inputText.value = `这是一个示例文本，用于演示字符统计工具的功能。

This is a sample text to demonstrate the character counting tool features. It contains both Chinese and English characters, numbers like 123, and various punctuation marks!

本工具可以统计：
1. 总字符数和不含空格的字符数
2. 中文字符、英文字符、数字和标点符号的数量
3. 单词数、行数和段落数
4. 词频统计和阅读时间预估

希望这个工具对您有所帮助！`
  updateStats()
  ElMessage.success('示例文本已加载')
}

const copyStats = async () => {
  const statsText = `文本统计结果：
总字符数：${stats.value.totalChars.toLocaleString()}
字符数(不含空格)：${stats.value.charsNoSpaces.toLocaleString()}
中文字符数：${stats.value.chineseChars.toLocaleString()}
英文字符数：${stats.value.englishChars.toLocaleString()}
数字字符数：${stats.value.numberChars.toLocaleString()}
标点符号数：${stats.value.punctuationChars.toLocaleString()}
单词数：${stats.value.words.toLocaleString()}
行数：${stats.value.lines.toLocaleString()}
段落数：${stats.value.paragraphs.toLocaleString()}
预计阅读时间：${readingTime.value}
预计说话时间：${speakingTime.value}`

  await copy(statsText)
}

const exportStats = () => {
  const reportData = {
    timestamp: new Date().toLocaleString(),
    text: inputText.value,
    statistics: stats.value,
    wordFrequency: wordFrequency.value.slice(0, 50),
    readingTime: readingTime.value,
    speakingTime: speakingTime.value
  }
  
  const blob = new Blob([JSON.stringify(reportData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `文本统计报告_${new Date().toISOString().slice(0, 10)}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  
  ElMessage.success('统计报告已导出')
}

onMounted(() => {
  updateStats()
})
</script>

<style scoped>
.word-counter-container {
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

.input-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px;
  text-align: center;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
}

.details-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.empty-state {
  text-align: center;
  color: #9ca3af;
  padding: 40px 20px;
}

.frequency-list {
  max-height: 400px;
  overflow-y: auto;
}

.frequency-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
}

.word {
  font-weight: 500;
  min-width: 80px;
}

.count {
  font-size: 12px;
  color: #6b7280;
  min-width: 40px;
}

.frequency-bar {
  flex: 1;
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
}

.frequency-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6, #06b6d4);
  border-radius: 3px;
}

.reading-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.reading-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f1f5f9;
}

.reading-label {
  color: #6b7280;
  font-weight: 500;
}

.reading-value {
  color: #1f2937;
  font-weight: 600;
}

.action-section {
  display: flex;
  gap: 12px;
  justify-content: center;
  padding-top: 20px;
  border-top: 1px solid #e5e7eb;
}

@media (max-width: 768px) {
  .word-counter-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .details-grid {
    grid-template-columns: 1fr;
  }
  
  .input-actions {
    flex-direction: column;
  }
  
  .action-section {
    flex-direction: column;
  }
  
  .stat-value {
    font-size: 24px;
  }
}
</style>