<template>
  <div class="citation-generator-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-academic' })">学术类</el-breadcrumb-item>
          <el-breadcrumb-item>引用生成器</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>学术引用生成器</h1>
      <p>生成APA、MLA、Chicago等标准学术引用格式</p>
    </div>

    <div class="tools-content">
      <div class="source-type-section">
        <el-card class="source-card">
          <template #header>
            <div class="card-header">
              <el-icon><Document /></el-icon>
              <span>选择资源类型</span>
            </div>
          </template>
          <div class="source-types">
            <div 
              v-for="type in sourceTypes" 
              :key="type.id"
              class="source-type"
              :class="{ active: selectedSourceType === type.id }"
              @click="selectSourceType(type.id)"
            >
              <div class="type-icon">{{ type.icon }}</div>
              <div class="type-name">{{ type.name }}</div>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="selectedSourceType" class="input-section">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><EditPen /></el-icon>
              <span>填写资源信息</span>
            </div>
          </template>
          <div class="form-content">
            <div class="input-grid">
              <div 
                v-for="field in currentFields" 
                :key="field.key"
                class="input-item"
              >
                <label>{{ field.label }} <span v-if="field.required" class="required">*</span></label>
                <el-input 
                  v-model="formData[field.key]" 
                  :placeholder="field.placeholder"
                  clearable
                />
              </div>
            </div>
            <div class="generate-section">
              <el-button @click="generateCitations" type="primary" size="large">
                <el-icon><MagicStick /></el-icon>
                生成引用
              </el-button>
              <el-button @click="clearForm" type="default" size="large">
                <el-icon><RefreshLeft /></el-icon>
                清空
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="citations.length > 0" class="output-section">
        <el-card class="output-card">
          <template #header>
            <div class="card-header">
              <el-icon><Collection /></el-icon>
              <span>生成的引用</span>
            </div>
          </template>
          <div class="citations-content">
            <div 
              v-for="citation in citations" 
              :key="citation.style"
              class="citation-item"
            >
              <div class="citation-header">
                <h3>{{ citation.style }}</h3>
                <el-button @click="copyCitation(citation.text)" size="small" type="primary">
                  <el-icon><DocumentCopy /></el-icon>
                  复制
                </el-button>
              </div>
              <div class="citation-text">{{ citation.text }}</div>
            </div>
          </div>
        </el-card>
      </div>

      <div class="info-section">
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>引用格式说明</span>
            </div>
          </template>
          <div class="info-content">
            <div class="format-item">
              <h4>APA格式</h4>
              <p>美国心理学会格式，常用于心理学、教育学、社会科学等领域</p>
            </div>
            <div class="format-item">
              <h4>MLA格式</h4>
              <p>现代语言协会格式，常用于文学、语言学、人文学科等领域</p>
            </div>
            <div class="format-item">
              <h4>Chicago格式</h4>
              <p>芝加哥格式，常用于历史学、文学、艺术等领域</p>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Document,
  EditPen,
  MagicStick,
  RefreshLeft,
  Collection,
  DocumentCopy,
  InfoFilled
} from '@element-plus/icons-vue'

const router = useRouter()

const selectedSourceType = ref('')
const formData = ref({})
const citations = ref([])

const sourceTypes = ref([
  { id: 'book', name: '图书', icon: '📚' },
  { id: 'journal', name: '期刊文章', icon: '📰' },
  { id: 'website', name: '网站', icon: '🌐' },
  { id: 'thesis', name: '学位论文', icon: '🎓' }
])

const fieldTemplates = {
  book: [
    { key: 'author', label: '作者', placeholder: '姓, 名', required: true },
    { key: 'title', label: '书名', placeholder: '书籍标题', required: true },
    { key: 'publisher', label: '出版社', placeholder: '出版社名称', required: true },
    { key: 'year', label: '出版年份', placeholder: '2023', required: true },
    { key: 'location', label: '出版地', placeholder: '城市', required: false }
  ],
  journal: [
    { key: 'author', label: '作者', placeholder: '姓, 名', required: true },
    { key: 'title', label: '文章标题', placeholder: '文章标题', required: true },
    { key: 'journal', label: '期刊名称', placeholder: '期刊名称', required: true },
    { key: 'volume', label: '卷号', placeholder: '卷号', required: false },
    { key: 'issue', label: '期号', placeholder: '期号', required: false },
    { key: 'pages', label: '页码', placeholder: '1-10', required: false },
    { key: 'year', label: '发表年份', placeholder: '2023', required: true }
  ],
  website: [
    { key: 'author', label: '作者', placeholder: '姓, 名（可选）', required: false },
    { key: 'title', label: '页面标题', placeholder: '网页标题', required: true },
    { key: 'website', label: '网站名称', placeholder: '网站名称', required: true },
    { key: 'url', label: '网址', placeholder: 'https://example.com', required: true },
    { key: 'accessDate', label: '访问日期', placeholder: '2023-12-01', required: true }
  ],
  thesis: [
    { key: 'author', label: '作者', placeholder: '姓, 名', required: true },
    { key: 'title', label: '论文标题', placeholder: '论文标题', required: true },
    { key: 'degree', label: '学位类型', placeholder: '硕士/博士', required: true },
    { key: 'institution', label: '学校', placeholder: '学校名称', required: true },
    { key: 'year', label: '年份', placeholder: '2023', required: true }
  ]
}

const currentFields = computed(() => {
  return fieldTemplates[selectedSourceType.value] || []
})

const selectSourceType = (typeId) => {
  selectedSourceType.value = typeId
  formData.value = {}
  citations.value = []
}

const generateCitations = () => {
  if (!validateForm()) return
  
  citations.value = [
    { style: 'APA格式', text: generateAPA() },
    { style: 'MLA格式', text: generateMLA() },
    { style: 'Chicago格式', text: generateChicago() }
  ]
  
  ElMessage.success('引用格式生成完成')
}

const validateForm = () => {
  const requiredFields = currentFields.value.filter(field => field.required)
  for (const field of requiredFields) {
    if (!formData.value[field.key] || !formData.value[field.key].trim()) {
      ElMessage.warning(`请填写${field.label}`)
      return false
    }
  }
  return true
}

const generateAPA = () => {
  const data = formData.value
  const type = selectedSourceType.value
  
  switch (type) {
    case 'book':
      return `${data.author} (${data.year}). ${data.title}. ${data.publisher}.`
    case 'journal':
      const volume = data.volume ? `, ${data.volume}` : ''
      const issue = data.issue ? `(${data.issue})` : ''
      const pages = data.pages ? `, ${data.pages}` : ''
      return `${data.author} (${data.year}). ${data.title}. ${data.journal}${volume}${issue}${pages}.`
    case 'website':
      const author = data.author ? `${data.author}. ` : ''
      return `${author}${data.title}. ${data.website}. Retrieved ${data.accessDate}, from ${data.url}`
    case 'thesis':
      return `${data.author} (${data.year}). ${data.title} (${data.degree} thesis). ${data.institution}.`
    default:
      return ''
  }
}

const generateMLA = () => {
  const data = formData.value
  const type = selectedSourceType.value
  
  switch (type) {
    case 'book':
      return `${data.author}. ${data.title}. ${data.publisher}, ${data.year}.`
    case 'journal':
      const volume = data.volume ? `, vol. ${data.volume}` : ''
      const issue = data.issue ? `, no. ${data.issue}` : ''
      const pages = data.pages ? `, pp. ${data.pages}` : ''
      return `${data.author}. "${data.title}." ${data.journal}${volume}${issue}, ${data.year}${pages}.`
    case 'website':
      const author = data.author ? `${data.author}. ` : ''
      return `${author}"${data.title}." ${data.website}, ${data.accessDate}, ${data.url}.`
    case 'thesis':
      return `${data.author}. ${data.title}. ${data.year}. ${data.institution}, ${data.degree} thesis.`
    default:
      return ''
  }
}

const generateChicago = () => {
  const data = formData.value
  const type = selectedSourceType.value
  
  switch (type) {
    case 'book':
      return `${data.author}. ${data.title}. ${data.location || ''}: ${data.publisher}, ${data.year}.`
    case 'journal':
      const volume = data.volume ? ` ${data.volume}` : ''
      const issue = data.issue ? `, no. ${data.issue}` : ''
      const pages = data.pages ? ` (${data.year}): ${data.pages}` : ` (${data.year})`
      return `${data.author}. "${data.title}." ${data.journal}${volume}${issue}${pages}.`
    case 'website':
      const author = data.author ? `${data.author}. ` : ''
      return `${author}"${data.title}." ${data.website}. Accessed ${data.accessDate}. ${data.url}.`
    case 'thesis':
      return `${data.author}. "${data.title}." ${data.degree} thesis, ${data.institution}, ${data.year}.`
    default:
      return ''
  }
}

const copyCitation = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('引用已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const clearForm = () => {
  formData.value = {}
  citations.value = []
  ElMessage.info('表单已清空')
}
</script>

<style scoped>
.citation-generator-container {
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

.source-card,
.input-card,
.output-card,
.info-card {
  border: 1px solid #e5e7eb;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #2c3e50;
}

.source-types {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
}

.source-type {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: center;
}

.source-type:hover {
  border-color: #409eff;
  background-color: #f8fafc;
}

.source-type.active {
  border-color: #409eff;
  background-color: rgba(64, 158, 255, 0.05);
}

.type-icon {
  font-size: 32px;
}

.type-name {
  font-weight: 500;
  color: #2c3e50;
}

.input-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.input-item label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #374151;
}

.required {
  color: #ef4444;
}

.generate-section {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.citations-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.citation-item {
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.citation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.citation-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.citation-text {
  font-family: serif;
  line-height: 1.6;
  color: #374151;
  background: #fff;
  padding: 16px;
  border-radius: 6px;
  border-left: 4px solid #409eff;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.format-item h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.format-item p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .citation-generator-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .source-types {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .input-grid {
    grid-template-columns: 1fr;
  }
  
  .generate-section {
    flex-direction: column;
  }
  
  .citation-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>