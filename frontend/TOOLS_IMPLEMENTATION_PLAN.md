# 工具箱系统实现方案

## 📊 当前工具实现状态分析

### ✅ 已完全实现的工具（24个）

#### 1. 开发类工具（4个）
- ✅ 时间戳转换工具
- ✅ JSON工具（格式化/压缩/转YAML）
- ✅ 字符串生成工具
- ✅ HTTP状态码查询工具

#### 2. 文本类工具（3个）
- ✅ Base64编码器
- ✅ URL编码器  
- ✅ 文本处理器

#### 3. 图像类工具（3个）
- ✅ 图片压缩器
- ✅ 图片格式转换器
- ✅ 图片尺寸调整器

#### 4. 其它类工具（5个）
- ✅ 二维码生成器
- ✅ 颜色选择器
- ✅ 密码强度检测器
- ✅ 正则表达式测试器
- ✅ Hash计算器

#### 5. 运维类工具（1个）
- ✅ Ping测试工具

#### 6. 学术类工具（2个）
- ✅ 引用生成器
- ✅ 数学计算器

#### 7. 查询类工具（3个）
- ✅ IP地址查询
- ✅ Whois查询
- ✅ 域名信息查询

#### 8. 文档类工具（0个）
- ❌ 所有工具待实现

---

## 🚀 待实现工具详细方案（9个）

**✅ 已完成查询类工具（3个）：**
- ✅ IP地址查询 - 支持IP归属地、ISP信息、地理位置查询
- ✅ Whois查询 - 域名注册信息、联系人信息、DNS服务器查询
- ✅ 域名信息查询 - SSL证书、网站信息、服务器配置查询

### 1. 运维类工具（2个待实现）

#### 1.1 端口扫描器（Port Scanner）
**文件名**: `PortScanner.vue`  
**路由**: `/tools/port-scanner`

**核心功能**:
- 单端口/端口范围扫描
- 常见端口预设（21,22,23,25,53,80,110,443,993,995等）
- TCP连接检测
- 扫描结果统计和导出

**技术实现**:
```javascript
// 主要逻辑 - 使用WebSocket检测端口连通性
const scanPort = async (host, port) => {
  return new Promise((resolve) => {
    const ws = new WebSocket(`ws://${host}:${port}`)
    const timeout = setTimeout(() => {
      ws.close()
      resolve({ port, status: 'closed', latency: null })
    }, 3000)
    
    const startTime = Date.now()
    ws.onopen = () => {
      clearTimeout(timeout)
      const latency = Date.now() - startTime
      ws.close()
      resolve({ port, status: 'open', latency })
    }
    
    ws.onerror = () => {
      clearTimeout(timeout)
      resolve({ port, status: 'closed', latency: null })
    }
  })
}

// 批量扫描
const scanPorts = async (host, ports) => {
  const results = []
  for (const port of ports) {
    const result = await scanPort(host, port)
    results.push(result)
  }
  return results
}
```

**UI组件设计**:
- 主机输入框（域名/IP地址）
- 端口范围选择器（起始端口-结束端口）
- 预设端口快选按钮组
- 实时扫描进度条和状态显示
- 结果表格（端口/状态/服务名/延迟）
- 导出功能（CSV/JSON格式）

**常见端口预设**:
```javascript
const commonPorts = {
  web: [80, 443, 8080, 8443, 3000, 5000],
  email: [25, 110, 143, 465, 587, 993, 995],
  ftp: [20, 21, 22],
  database: [1433, 1521, 3306, 5432, 6379, 27017],
  remote: [22, 23, 3389, 5900, 5901],
  dns: [53, 853]
}
```

#### 1.2 DNS查询工具（DNS Lookup）
**文件名**: `DNSLookup.vue`  
**路由**: `/tools/dns-lookup`

**核心功能**:
- A/AAAA/CNAME/MX/TXT/NS/SOA记录查询
- 反向DNS查询（IP到域名）
- DNS服务器选择（8.8.8.8, 1.1.1.1, 114.114.114.114等）
- 查询历史记录和收藏功能

**技术实现**:
```javascript
// 使用Google Public DNS API
const queryDNS = async (domain, type = 'A', resolver = '8.8.8.8') => {
  const url = `https://dns.google/resolve?name=${domain}&type=${type}&cd=0&do=0`
  try {
    const response = await fetch(url)
    const data = await response.json()
    return {
      success: true,
      answers: data.Answer || [],
      authority: data.Authority || [],
      additional: data.Additional || []
    }
  } catch (error) {
    return { success: false, error: error.message }
  }
}

// 反向DNS查询
const reverseDNS = async (ip) => {
  const reverseIP = ip.split('.').reverse().join('.') + '.in-addr.arpa'
  return await queryDNS(reverseIP, 'PTR')
}
```

**UI组件设计**:
- 域名/IP输入框（支持批量查询）
- DNS记录类型选择器（下拉多选）
- DNS服务器选择（预设+自定义）
- 查询结果卡片展示
- 历史记录和收藏管理
- 导出查询结果

---

### 2. ✅ 查询类工具（已完成3个）

所有查询类工具已完成实现：

#### 2.1 ✅ IP地址查询（IP Lookup） - 已实现
- ✅ IP归属地查询（国家/省份/城市）
- ✅ ISP信息显示（运营商/组织）
- ✅ IPv4/IPv6支持
- ✅ 获取当前IP功能
- ✅ 查询历史记录
- ✅ 响应式设计

#### 2.2 ✅ Whois查询（Whois Lookup） - 已实现
- ✅ 域名注册信息查询
- ✅ 注册商/注册人信息展示
- ✅ 域名过期时间和状态
- ✅ DNS服务器信息
- ✅ 联系人信息管理
- ✅ 原始Whois数据展示

#### 2.3 ✅ 域名信息查询（Domain Info） - 已实现
- ✅ SSL证书信息查询
- ✅ 网站基本信息（标题/描述/关键词）
- ✅ 服务器信息（IP地址/服务器类型）
- ✅ 域名年龄计算
- ✅ 技术信息分析
- ✅ DNS记录查询

---

### 3. 学术类工具（1个待实现）

#### 3.1 数据分析工具（Data Analyzer）
**文件名**: `DataAnalyzer.vue`  
**路由**: `/tools/data-analyzer`

**核心功能**:
- CSV/Excel文件导入解析
- 基础统计计算（均值/中位数/标准差/方差）
- 数据可视化（柱状图/折线图/散点图/饼图）
- 数据清洗和预处理
- 分析结果导出

**技术实现**:
```javascript
// 使用Papa Parse处理CSV文件
import Papa from 'papaparse'

const parseCSV = (file) => {
  return new Promise((resolve, reject) => {
    Papa.parse(file, {
      header: true,
      skipEmptyLines: true,
      complete: (results) => {
        resolve(results.data)
      },
      error: (error) => {
        reject(error)
      }
    })
  })
}

// 基础统计计算
const calculateStatistics = (data, column) => {
  const values = data.map(row => parseFloat(row[column])).filter(val => !isNaN(val))
  
  if (values.length === 0) return null
  
  const sorted = values.sort((a, b) => a - b)
  const sum = values.reduce((acc, val) => acc + val, 0)
  const mean = sum / values.length
  
  const variance = values.reduce((acc, val) => acc + Math.pow(val - mean, 2), 0) / values.length
  const stdDev = Math.sqrt(variance)
  
  return {
    count: values.length,
    sum,
    mean,
    median: sorted[Math.floor(sorted.length / 2)],
    min: Math.min(...values),
    max: Math.max(...values),
    variance,
    standardDeviation: stdDev
  }
}

// 使用ECharts生成图表
import * as echarts from 'echarts'

const generateChart = (container, data, type) => {
  const chart = echarts.init(container)
  
  const option = {
    title: { text: '数据可视化' },
    tooltip: {},
    xAxis: { data: data.labels },
    yAxis: {},
    series: [{
      name: '数值',
      type: type, // 'bar', 'line', 'scatter', 'pie'
      data: data.values
    }]
  }
  
  chart.setOption(option)
  return chart
}
```

---

### 4. 文档类工具（3个全新工具）

#### 4.1 Markdown编辑器（Markdown Editor）
**文件名**: `MarkdownEditor.vue`  
**路由**: `/tools/markdown-editor`

**核心功能**:
- 实时预览（分屏显示）
- 语法高亮编辑
- 导出HTML/PDF功能
- 表格编辑器
- 图片上传和粘贴支持
- 常用语法快捷插入

**技术实现**:
```javascript
// 使用marked.js解析Markdown
import { marked } from 'marked'
import hljs from 'highlight.js'

// 配置marked
marked.setOptions({
  highlight: function(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value
    }
    return hljs.highlightAuto(code).value
  },
  breaks: true,
  gfm: true
})

const renderMarkdown = (markdown) => {
  return marked(markdown)
}

// PDF导出功能
import jsPDF from 'jspdf'

const exportToPDF = (html) => {
  const pdf = new jsPDF()
  pdf.html(html, {
    callback: function (doc) {
      doc.save('markdown-export.pdf')
    }
  })
}
```

#### 4.2 PDF工具集（PDF Tools）
**文件名**: `PDFTools.vue`  
**路由**: `/tools/pdf-tools`

**核心功能**:
- PDF文件合并
- PDF页面拆分和提取
- PDF转图片（PNG/JPG）
- PDF信息查看（页数/大小/创建时间）
- PDF页面旋转
- PDF加密和解密

**技术实现**:
```javascript
// 使用PDF-lib处理PDF
import { PDFDocument } from 'pdf-lib'

// 合并PDF文件
const mergePDFs = async (pdfFiles) => {
  const mergedPdf = await PDFDocument.create()
  
  for (const pdfFile of pdfFiles) {
    const pdf = await PDFDocument.load(pdfFile)
    const copiedPages = await mergedPdf.copyPages(pdf, pdf.getPageIndices())
    copiedPages.forEach((page) => mergedPdf.addPage(page))
  }
  
  return await mergedPdf.save()
}

// PDF转图片
const pdfToImages = async (pdfFile) => {
  const pdf = await PDFDocument.load(pdfFile)
  const pageCount = pdf.getPageCount()
  const images = []
  
  for (let i = 0; i < pageCount; i++) {
    // 使用Canvas API渲染PDF页面为图片
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    // 渲染逻辑...
    images.push(canvas.toDataURL('image/png'))
  }
  
  return images
}
```

#### 4.3 词云生成器（Word Cloud Generator）
**文件名**: `WordCloudGenerator.vue`  
**路由**: `/tools/word-cloud`

**核心功能**:
- 文本词频统计
- 词云可视化生成
- 字体样式自定义
- 颜色主题选择
- 形状模板（圆形/方形/自定义）
- 高清图片导出

**技术实现**:
```javascript
// 使用wordcloud2.js生成词云
import WordCloud from 'wordcloud'

// 文本预处理和词频统计
const processText = (text) => {
  const words = text.toLowerCase()
    .replace(/[^\w\s]/g, '') // 移除标点
    .split(/\s+/) // 按空格分割
    .filter(word => word.length > 2) // 过滤短词
  
  const frequency = {}
  words.forEach(word => {
    frequency[word] = (frequency[word] || 0) + 1
  })
  
  return Object.entries(frequency)
    .sort(([,a], [,b]) => b - a)
    .slice(0, 100) // 取前100个高频词
}

// 生成词云
const generateWordCloud = (canvas, wordList, options = {}) => {
  WordCloud(canvas, {
    list: wordList,
    gridSize: options.gridSize || 8,
    weightFactor: options.weightFactor || 16,
    fontFamily: options.fontFamily || 'Arial',
    color: options.color || 'random-dark',
    backgroundColor: options.backgroundColor || '#ffffff',
    rotateRatio: options.rotateRatio || 0.5,
    shape: options.shape || 'circle'
  })
}
```

---

## 📋 实现优先级建议

### 🥇 第一批（✅ 已完成）
1. **✅ IP地址查询** - 使用免费API，实现简单，用户需求高
2. **✅ Whois查询** - 域名管理和分析需求  
3. **✅ 域名信息查询** - SSL证书查询很实用

### 🥈 第二批（推荐优先实现）
4. **DNS查询工具** - 网络管理员和开发者常用工具
5. **端口扫描器** - 网络安全测试的基础需求
6. **数据分析工具** - 学术和商业价值高

### 🥉 第三批（功能完整性，实现较复杂）
7. **Markdown编辑器** - 技术文档编辑需求
8. **词云生成器** - 数据可视化应用
9. **PDF工具集** - 文件处理需求广泛

---

## 🛠️ 技术实现注意事项

### 1. API选择原则
- **免费优先**: 优先使用稳定的免费公共API
- **备用方案**: 为每个API准备2-3个备选方案
- **错误处理**: 完善的网络请求超时和错误处理机制
- **限频处理**: 实现请求限频和重试机制

### 2. 性能优化
- **大文件处理**: 使用Web Workers处理大文件解析
- **内存管理**: 及时释放大对象引用，避免内存泄漏
- **懒加载**: 图表库等重型依赖采用动态导入
- **缓存策略**: 合理使用本地存储缓存查询结果

### 3. 安全考虑
- **客户端处理**: 所有敏感操作在客户端完成
- **数据验证**: 严格的输入验证和格式检查
- **CORS处理**: 正确处理跨域请求限制
- **XSS防护**: 对用户输入进行适当的转义处理

### 4. 用户体验
- **加载状态**: 统一的Loading组件和骨架屏
- **错误提示**: 友好的错误信息和解决建议
- **操作反馈**: 及时的操作成功/失败反馈
- **快捷操作**: 常用功能的快捷键支持

### 5. 代码规范
- **组件复用**: 抽象公共组件，避免代码重复
- **错误边界**: React/Vue错误边界处理
- **类型安全**: 使用TypeScript增强类型安全
- **测试覆盖**: 为核心功能编写单元测试

---

## 📁 文件结构建议

```
src/views/
├── tools/
│   ├── network/
│   │   ├── PortScanner.vue
│   │   ├── DNSLookup.vue
│   │   └── PingTool.vue (已实现)
│   ├── query/
│   │   ├── IPLookup.vue
│   │   ├── WhoisLookup.vue
│   │   └── DomainInfo.vue
│   ├── academic/
│   │   ├── DataAnalyzer.vue
│   │   ├── CitationGenerator.vue (已实现)
│   │   └── MathCalculator.vue (已实现)
│   └── document/
│       ├── MarkdownEditor.vue
│       ├── PDFTools.vue
│       └── WordCloudGenerator.vue
└── components/
    ├── common/
    │   ├── LoadingSpinner.vue
    │   ├── ErrorAlert.vue
    │   └── ExportButton.vue
    └── charts/
        ├── BarChart.vue
        ├── LineChart.vue
        └── PieChart.vue
```

这个实现方案涵盖了所有未实现的工具，提供了详细的技术实现思路、UI设计建议和开发优先级，可以作为后续开发的详细指导文档。