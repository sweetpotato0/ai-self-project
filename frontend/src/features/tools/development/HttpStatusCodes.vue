<template>
  <div class="http-status-container">
    <div class="tool-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-development' })">开发类</el-breadcrumb-item>
          <el-breadcrumb-item>HTTP状态码</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>HTTP 状态码查询</h1>
      <p>完整的HTTP状态码参考指南，按类别分组显示</p>
    </div>

    <!-- 搜索和过滤 -->
    <div class="search-section">
      <div class="search-controls">
        <el-input
          v-model="searchQuery"
          placeholder="搜索状态码或描述..."
          class="search-input"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="selectedCategory"
          placeholder="选择类别"
          class="category-select"
          clearable
        >
          <el-option label="全部类别" value="" />
          <el-option 
            v-for="category in statusCategories" 
            :key="category.code"
            :label="category.name"
            :value="category.code"
          />
        </el-select>
      </div>
    </div>

    <!-- 状态码统计 -->
    <div class="stats-section">
      <div class="stat-card" v-for="category in statusCategories" :key="category.code">
        <div class="stat-icon" :style="{ background: category.color }">
          <span>{{ category.code }}xx</span>
        </div>
        <div class="stat-content">
          <div class="stat-name">{{ category.name }}</div>
          <div class="stat-count">{{ getStatusCountByCategory(category.code) }} 个</div>
        </div>
      </div>
    </div>

    <!-- 状态码列表 -->
    <div class="status-sections">
      <div 
        v-for="category in statusCategories" 
        :key="category.code"
        v-show="shouldShowCategory(category.code)"
        class="category-section"
      >
        <div class="category-header" :style="{ borderLeftColor: category.color }">
          <div class="category-title">
            <span class="category-code">{{ category.code }}xx</span>
            <span class="category-name">{{ category.name }}</span>
          </div>
          <div class="category-description">{{ category.description }}</div>
        </div>

        <div class="status-grid">
          <div 
            v-for="status in getFilteredStatusesByCategory(category.code)"
            :key="status.code"
            class="status-card"
            @click="showStatusDetail(status)"
          >
            <div class="status-code" :style="{ color: category.color }">
              {{ status.code }}
            </div>
            <div class="status-phrase">{{ status.phrase }}</div>
            <div class="status-meaning">{{ status.meaning }}</div>
            <div class="status-category-tag" :style="{ background: category.color + '20', color: category.color }">
              {{ category.name }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="showDetailDialog"
      :title="`HTTP ${selectedStatus?.code} - ${selectedStatus?.phrase}`"
      width="600px"
      class="status-detail-dialog"
    >
      <div v-if="selectedStatus" class="status-detail">
        <div class="detail-header">
          <div class="detail-code" :style="{ color: getCategoryColor(selectedStatus.code) }">
            {{ selectedStatus.code }}
          </div>
          <div class="detail-info">
            <div class="detail-phrase">{{ selectedStatus.phrase }}</div>
            <div class="detail-category">{{ getCategoryName(selectedStatus.code) }}</div>
          </div>
        </div>
        
        <div class="detail-content">
          <div class="detail-section">
            <h4>含义说明</h4>
            <p>{{ selectedStatus.meaning }}</p>
          </div>
          
          <div class="detail-section">
            <h4>详细描述</h4>
            <p>{{ selectedStatus.description }}</p>
          </div>
          
          <div class="detail-section" v-if="selectedStatus.examples">
            <h4>使用场景</h4>
            <ul>
              <li v-for="example in selectedStatus.examples" :key="example">{{ example }}</li>
            </ul>
          </div>
          
          <div class="detail-section" v-if="selectedStatus.headers">
            <h4>相关响应头</h4>
            <div class="headers-list">
              <el-tag v-for="header in selectedStatus.headers" :key="header" type="info" class="header-tag">
                {{ header }}
              </el-tag>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'

const router = useRouter()

// 响应式数据
const searchQuery = ref('')
const selectedCategory = ref('')
const showDetailDialog = ref(false)
const selectedStatus = ref(null)

// 状态码分类
const statusCategories = ref([
  {
    code: '1',
    name: '信息响应',
    description: '请求已接收，继续处理',
    color: '#3b82f6'
  },
  {
    code: '2',
    name: '成功响应',
    description: '请求已成功接收、理解、接受',
    color: '#10b981'
  },
  {
    code: '3',
    name: '重定向',
    description: '需要进一步操作以完成请求',
    color: '#f59e0b'
  },
  {
    code: '4',
    name: '客户端错误',
    description: '请求包含错误语法或无法完成',
    color: '#ef4444'
  },
  {
    code: '5',
    name: '服务器错误',
    description: '服务器无法完成明显有效的请求',
    color: '#8b5cf6'
  }
])

// 完整的状态码列表
const httpStatusCodes = ref([
  // 1xx 信息响应
  {
    code: 100,
    phrase: 'Continue',
    meaning: '继续',
    description: '服务器已接收到请求头，客户端应继续发送请求体。通常在客户端需要发送大量数据时使用。',
    examples: ['POST请求发送大文件', 'PUT请求上传数据'],
    headers: ['Expect: 100-continue']
  },
  {
    code: 101,
    phrase: 'Switching Protocols',
    meaning: '切换协议',
    description: '服务器已理解并准备切换到客户端请求的协议。常用于WebSocket升级。',
    examples: ['HTTP升级到WebSocket', 'HTTP/1.1升级到HTTP/2'],
    headers: ['Upgrade', 'Connection']
  },
  {
    code: 102,
    phrase: 'Processing',
    meaning: '处理中',
    description: '服务器已接收并正在处理请求，但还没有可用的响应。防止客户端超时。',
    examples: ['WebDAV请求处理', '长时间运行的操作'],
    headers: []
  },
  {
    code: 103,
    phrase: 'Early Hints',
    meaning: '早期提示',
    description: '在最终HTTP消息之前发送，用于预加载资源。',
    examples: ['预加载CSS和JS文件', '提前DNS解析'],
    headers: ['Link']
  },

  // 2xx 成功响应
  {
    code: 200,
    phrase: 'OK',
    meaning: '成功',
    description: '请求已成功。响应的信息取决于所使用的方法。',
    examples: ['GET请求获取资源成功', 'POST请求处理成功', '表单提交成功'],
    headers: ['Content-Type', 'Content-Length']
  },
  {
    code: 201,
    phrase: 'Created',
    meaning: '已创建',
    description: '请求已成功，并且服务器创建了新的资源。',
    examples: ['POST创建新用户', 'PUT创建新文件', 'API创建新记录'],
    headers: ['Location', 'Content-Type']
  },
  {
    code: 202,
    phrase: 'Accepted',
    meaning: '已接受',
    description: '请求已接受处理，但处理还未完成。通常用于异步处理。',
    examples: ['异步任务提交', '批量处理请求', '邮件发送请求'],
    headers: ['Content-Type']
  },
  {
    code: 203,
    phrase: 'Non-Authoritative Information',
    meaning: '非授权信息',
    description: '请求成功，但返回的信息可能来自另一来源。',
    examples: ['代理服务器返回缓存信息', '镜像站点数据'],
    headers: ['Content-Type']
  },
  {
    code: 204,
    phrase: 'No Content',
    meaning: '无内容',
    description: '请求成功处理，但没有返回任何内容。',
    examples: ['DELETE操作成功', 'PUT更新成功', '表单提交无需返回数据'],
    headers: []
  },
  {
    code: 205,
    phrase: 'Reset Content',
    meaning: '重置内容',
    description: '请求成功处理，客户端应重置文档视图。',
    examples: ['表单提交后重置', '编辑器清空内容'],
    headers: []
  },
  {
    code: 206,
    phrase: 'Partial Content',
    meaning: '部分内容',
    description: '服务器成功处理了部分GET请求，返回部分内容。',
    examples: ['视频流播放', '文件断点续传', '分页数据获取'],
    headers: ['Content-Range', 'Content-Type']
  },
  {
    code: 207,
    phrase: 'Multi-Status',
    meaning: '多状态',
    description: '多个资源操作的状态信息。WebDAV扩展。',
    examples: ['批量文件操作', 'WebDAV PROPFIND'],
    headers: ['Content-Type: application/xml']
  },
  {
    code: 208,
    phrase: 'Already Reported',
    meaning: '已报告',
    description: 'DAV绑定的成员已在此请求的早期部分被枚举，不再重复包含。',
    examples: ['WebDAV资源枚举'],
    headers: []
  },
  {
    code: 226,
    phrase: 'IM Used',
    meaning: '使用实例操作',
    description: '服务器已完成GET请求，响应是一个或多个实例操作的结果。',
    examples: ['Delta编码响应', 'HTTP实例操作'],
    headers: ['IM', 'Delta-Base']
  },

  // 3xx 重定向
  {
    code: 300,
    phrase: 'Multiple Choices',
    meaning: '多种选择',
    description: '请求有多种可能的响应，用户或浏览器应选择其中一种。',
    examples: ['多语言版本选择', '多格式内容选择'],
    headers: ['Location', 'Content-Type']
  },
  {
    code: 301,
    phrase: 'Moved Permanently',
    meaning: '永久移动',
    description: '请求的资源已永久移动到新位置，搜索引擎会更新索引。',
    examples: ['网站域名更改', 'URL结构调整', 'HTTPS强制重定向'],
    headers: ['Location']
  },
  {
    code: 302,
    phrase: 'Found',
    meaning: '找到',
    description: '请求的资源现在临时从不同的URI获得。',
    examples: ['临时页面跳转', '登录后跳转', '临时维护页面'],
    headers: ['Location']
  },
  {
    code: 303,
    phrase: 'See Other',
    meaning: '查看其他位置',
    description: '对应请求的响应可以在另一个URI上找到，应使用GET方法获取。',
    examples: ['POST-Redirect-GET模式', '表单提交后跳转'],
    headers: ['Location']
  },
  {
    code: 304,
    phrase: 'Not Modified',
    meaning: '未修改',
    description: '资源未修改，可以使用缓存版本。',
    examples: ['浏览器缓存验证', 'CDN缓存检查', 'API缓存验证'],
    headers: ['Cache-Control', 'ETag', 'Last-Modified']
  },
  {
    code: 305,
    phrase: 'Use Proxy',
    meaning: '使用代理',
    description: '请求必须通过指定的代理进行访问。已废弃。',
    examples: ['代理服务器访问'],
    headers: ['Location']
  },
  {
    code: 307,
    phrase: 'Temporary Redirect',
    meaning: '临时重定向',
    description: '请求应使用相同方法重复到另一个URI。',
    examples: ['服务器维护跳转', '负载均衡重定向'],
    headers: ['Location']
  },
  {
    code: 308,
    phrase: 'Permanent Redirect',
    meaning: '永久重定向',
    description: '请求和所有将来的请求应使用相同方法重复到另一个URI。',
    examples: ['API版本升级', '永久URL更改'],
    headers: ['Location']
  },

  // 4xx 客户端错误
  {
    code: 400,
    phrase: 'Bad Request',
    meaning: '错误请求',
    description: '服务器无法理解请求的语法。',
    examples: ['JSON格式错误', '必填参数缺失', '请求体格式错误'],
    headers: ['Content-Type']
  },
  {
    code: 401,
    phrase: 'Unauthorized',
    meaning: '未授权',
    description: '请求需要用户验证。',
    examples: ['未提供认证信息', 'API密钥错误', 'Token过期'],
    headers: ['WWW-Authenticate']
  },
  {
    code: 402,
    phrase: 'Payment Required',
    meaning: '需要付费',
    description: '为将来使用而保留的状态码。',
    examples: ['付费内容访问', '超出免费额度'],
    headers: ['Content-Type']
  },
  {
    code: 403,
    phrase: 'Forbidden',
    meaning: '禁止',
    description: '服务器理解请求但拒绝执行。',
    examples: ['权限不足', 'IP被封禁', '资源访问受限'],
    headers: ['Content-Type']
  },
  {
    code: 404,
    phrase: 'Not Found',
    meaning: '未找到',
    description: '服务器找不到请求的资源。',
    examples: ['页面不存在', 'API接口不存在', '文件被删除'],
    headers: ['Content-Type']
  },
  {
    code: 405,
    phrase: 'Method Not Allowed',
    meaning: '方法不允许',
    description: '请求方法对于请求的资源不被允许。',
    examples: ['对只读资源使用POST', 'GET接口收到PUT请求'],
    headers: ['Allow']
  },
  {
    code: 406,
    phrase: 'Not Acceptable',
    meaning: '不可接受',
    description: '请求的资源无法满足Accept头中的条件。',
    examples: ['不支持的媒体类型', '不支持的语言'],
    headers: ['Content-Type']
  },
  {
    code: 407,
    phrase: 'Proxy Authentication Required',
    meaning: '需要代理验证',
    description: '客户端必须先使用代理进行验证。',
    examples: ['企业代理认证', 'HTTP代理登录'],
    headers: ['Proxy-Authenticate']
  },
  {
    code: 408,
    phrase: 'Request Time-out',
    meaning: '请求超时',
    description: '服务器等待请求时间过长超时。',
    examples: ['网络连接慢', '大文件上传超时', '数据库查询超时'],
    headers: ['Content-Type']
  },
  {
    code: 409,
    phrase: 'Conflict',
    meaning: '冲突',
    description: '请求与资源的当前状态冲突。',
    examples: ['数据版本冲突', '重复创建资源', '并发修改冲突'],
    headers: ['Content-Type']
  },
  {
    code: 410,
    phrase: 'Gone',
    meaning: '已删除',
    description: '请求的资源已被永久删除。',
    examples: ['文章被删除', '用户账号注销', '产品下架'],
    headers: ['Content-Type']
  },
  {
    code: 411,
    phrase: 'Length Required',
    meaning: '需要有效长度',
    description: '服务器需要Content-Length头字段。',
    examples: ['POST请求缺少Content-Length', '文件上传未指定大小'],
    headers: ['Content-Type']
  },
  {
    code: 412,
    phrase: 'Precondition Failed',
    meaning: '先决条件失败',
    description: '服务器无法满足请求头中的一个或多个先决条件。',
    examples: ['If-Match头验证失败', '条件请求不满足'],
    headers: ['Content-Type']
  },
  {
    code: 413,
    phrase: 'Request Entity Too Large',
    meaning: '请求实体过大',
    description: '请求实体超过服务器愿意或能够处理的大小。',
    examples: ['文件上传过大', '请求体超出限制', '图片文件太大'],
    headers: ['Retry-After']
  },
  {
    code: 414,
    phrase: 'Request-URI Too Large',
    meaning: '请求URI过长',
    description: '请求的URI长度超过服务器能够解释的长度。',
    examples: ['GET参数过多', 'URL过长', '查询字符串太大'],
    headers: ['Content-Type']
  },
  {
    code: 415,
    phrase: 'Unsupported Media Type',
    meaning: '不支持的媒体类型',
    description: '请求实体的媒体类型服务器不支持。',
    examples: ['上传不支持的文件格式', 'Content-Type不正确'],
    headers: ['Content-Type']
  },
  {
    code: 416,
    phrase: 'Requested Range Not Satisfiable',
    meaning: '请求范围不符合要求',
    description: '客户端请求的范围无效。',
    examples: ['Range头超出文件大小', '无效的字节范围'],
    headers: ['Content-Range']
  },
  {
    code: 417,
    phrase: 'Expectation Failed',
    meaning: '期望失败',
    description: '服务器无法满足Expect请求头的期望。',
    examples: ['Expect: 100-continue失败'],
    headers: ['Content-Type']
  },
  {
    code: 418,
    phrase: "I'm a teapot",
    meaning: '我是茶壶',
    description: '愚人节玩笑，服务器拒绝冲泡咖啡因为它是个茶壶。',
    examples: ['HTTP愚人节彩蛋', 'RFC 2324玩笑协议'],
    headers: ['Content-Type']
  },
  {
    code: 421,
    phrase: 'Misdirected Request',
    meaning: '错误定向请求',
    description: '请求被定向到无法产生响应的服务器。',
    examples: ['HTTP/2连接复用错误', 'SNI配置问题'],
    headers: ['Content-Type']
  },
  {
    code: 422,
    phrase: 'Unprocessable Entity',
    meaning: '无法处理的实体',
    description: '请求格式正确，但语义有误。',
    examples: ['表单验证失败', '业务逻辑错误', 'JSON结构正确但数据无效'],
    headers: ['Content-Type']
  },
  {
    code: 423,
    phrase: 'Locked',
    meaning: '已锁定',
    description: '当前资源被锁定。WebDAV扩展。',
    examples: ['文件被其他用户编辑', 'WebDAV资源锁定'],
    headers: ['Content-Type']
  },
  {
    code: 424,
    phrase: 'Failed Dependency',
    meaning: '失败的依赖',
    description: '由于之前的某个请求发生的错误，导致当前请求失败。',
    examples: ['WebDAV批量操作部分失败'],
    headers: ['Content-Type']
  },
  {
    code: 425,
    phrase: 'Too Early',
    meaning: '过早',
    description: '服务器不愿意冒险处理可能被重放的请求。',
    examples: ['TLS 1.3 Early Data安全限制'],
    headers: ['Content-Type']
  },
  {
    code: 426,
    phrase: 'Upgrade Required',
    meaning: '需要升级',
    description: '客户端应切换到不同的协议。',
    examples: ['强制HTTPS升级', 'WebSocket升级要求'],
    headers: ['Upgrade', 'Connection']
  },
  {
    code: 428,
    phrase: 'Precondition Required',
    meaning: '需要先决条件',
    description: '源服务器要求请求是有条件的。',
    examples: ['需要If-Match头', '防止并发修改'],
    headers: ['Content-Type']
  },
  {
    code: 429,
    phrase: 'Too Many Requests',
    meaning: '请求过多',
    description: '用户在给定时间内发送了太多请求（限流）。',
    examples: ['API调用频率限制', 'DDoS防护', '爬虫限制'],
    headers: ['Retry-After', 'X-RateLimit-Limit']
  },
  {
    code: 431,
    phrase: 'Request Header Fields Too Large',
    meaning: '请求头字段过大',
    description: '服务器不愿处理请求，因为头字段过大。',
    examples: ['Cookie过多', '请求头总大小超限'],
    headers: ['Content-Type']
  },
  {
    code: 451,
    phrase: 'Unavailable For Legal Reasons',
    meaning: '因法律原因不可用',
    description: '服务器因法律要求而拒绝提供资源访问。',
    examples: ['内容被政府审查', '版权限制', '地区屏蔽'],
    headers: ['Content-Type']
  },

  // 5xx 服务器错误
  {
    code: 500,
    phrase: 'Internal Server Error',
    meaning: '内部服务器错误',
    description: '服务器遇到意外情况，无法完成请求。',
    examples: ['程序代码错误', '数据库连接失败', '配置文件错误'],
    headers: ['Content-Type']
  },
  {
    code: 501,
    phrase: 'Not Implemented',
    meaning: '未实现',
    description: '服务器不支持实现请求所需要的功能。',
    examples: ['HTTP方法未实现', 'API功能未开发', '协议版本不支持'],
    headers: ['Content-Type']
  },
  {
    code: 502,
    phrase: 'Bad Gateway',
    meaning: '错误网关',
    description: '服务器作为网关或代理，从上游服务器收到无效响应。',
    examples: ['反向代理错误', '上游服务器宕机', '负载均衡器故障'],
    headers: ['Content-Type']
  },
  {
    code: 503,
    phrase: 'Service Unavailable',
    meaning: '服务不可用',
    description: '服务器目前无法使用，通常是临时状态。',
    examples: ['服务器维护', '流量过载', '临时关闭服务'],
    headers: ['Retry-After']
  },
  {
    code: 504,
    phrase: 'Gateway Time-out',
    meaning: '网关超时',
    description: '服务器作为网关或代理，没有及时从上游服务器收到请求。',
    examples: ['反向代理超时', '数据库查询超时', '微服务调用超时'],
    headers: ['Content-Type']
  },
  {
    code: 505,
    phrase: 'HTTP Version Not Supported',
    meaning: 'HTTP版本不受支持',
    description: '服务器不支持或拒绝支持请求消息中使用的HTTP版本。',
    examples: ['使用过旧HTTP版本', 'HTTP/3协议不支持'],
    headers: ['Content-Type']
  },
  {
    code: 506,
    phrase: 'Variant Also Negotiates',
    meaning: '变体也在协商',
    description: '服务器内部配置错误：所请求的协商变体资源被配置为在透明内容协商中使用自己。',
    examples: ['内容协商配置错误'],
    headers: ['Content-Type']
  },
  {
    code: 507,
    phrase: 'Insufficient Storage',
    meaning: '存储空间不足',
    description: '服务器无法存储完成请求所必须的内容。WebDAV扩展。',
    examples: ['磁盘空间不足', '上传文件超出存储限制'],
    headers: ['Content-Type']
  },
  {
    code: 508,
    phrase: 'Loop Detected',
    meaning: '检测到循环',
    description: '服务器在处理请求时陷入无限循环。WebDAV扩展。',
    examples: ['符号链接循环', 'WebDAV无限递归'],
    headers: ['Content-Type']
  },
  {
    code: 510,
    phrase: 'Not Extended',
    meaning: '未扩展',
    description: '获取资源所需要的策略并没有被满足。',
    examples: ['HTTP扩展框架策略不满足'],
    headers: ['Content-Type']
  },
  {
    code: 511,
    phrase: 'Network Authentication Required',
    meaning: '需要网络认证',
    description: '客户端需要进行身份验证才能获得网络访问权限。',
    examples: ['WiFi登录页面', '校园网认证', '企业网络认证'],
    headers: ['Content-Type']
  }
])

// 计算属性
const getStatusCountByCategory = (categoryCode) => {
  return httpStatusCodes.value.filter(status => 
    Math.floor(status.code / 100) === parseInt(categoryCode)
  ).length
}

const shouldShowCategory = (categoryCode) => {
  if (selectedCategory.value && selectedCategory.value !== categoryCode) {
    return false
  }
  return getFilteredStatusesByCategory(categoryCode).length > 0
}

const getFilteredStatusesByCategory = (categoryCode) => {
  return httpStatusCodes.value.filter(status => {
    const matchesCategory = Math.floor(status.code / 100) === parseInt(categoryCode)
    const matchesSearch = !searchQuery.value || 
      status.code.toString().includes(searchQuery.value) ||
      status.phrase.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      status.meaning.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      status.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    return matchesCategory && matchesSearch
  })
}

const getCategoryColor = (code) => {
  const categoryCode = Math.floor(code / 100).toString()
  const category = statusCategories.value.find(cat => cat.code === categoryCode)
  return category ? category.color : '#6b7280'
}

const getCategoryName = (code) => {
  const categoryCode = Math.floor(code / 100).toString()
  const category = statusCategories.value.find(cat => cat.code === categoryCode)
  return category ? category.name : '未知类别'
}

// 方法
const showStatusDetail = (status) => {
  selectedStatus.value = status
  showDetailDialog.value = true
}
</script>

<style scoped>
.http-status-container {
  padding: 20px 0;
  max-width: 1400px;
  margin: 0 auto;
}

.tool-header {
  margin-bottom: 30px;
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

.search-section {
  margin-bottom: 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  border-radius: 12px;
}

.search-controls {
  display: flex;
  gap: 16px;
  align-items: center;
}

.search-input {
  flex: 1;
  max-width: 400px;
}

.search-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.category-select {
  width: 200px;
}

.category-select :deep(.el-select__wrapper) {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.stats-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 30px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 14px;
}

.stat-content {
  flex: 1;
}

.stat-name {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.stat-count {
  font-size: 14px;
  color: #6b7280;
}

.status-sections {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.category-section {
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.category-header {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  padding: 24px;
  border-bottom: 1px solid #e2e8f0;
  border-left: 4px solid;
}

.category-title {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.category-code {
  font-size: 20px;
  font-weight: 700;
  color: #1f2937;
  font-family: 'Monaco', 'Menlo', monospace;
}

.category-name {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
}

.category-description {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  padding: 24px;
}

.status-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.status-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.15);
  border-color: #409eff;
}

.status-code {
  font-size: 24px;
  font-weight: 700;
  font-family: 'Monaco', 'Menlo', monospace;
  margin-bottom: 8px;
}

.status-phrase {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
  font-family: 'Monaco', 'Menlo', monospace;
}

.status-meaning {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
  margin-bottom: 12px;
}

.status-category-tag {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 6px;
  font-weight: 500;
  position: absolute;
  top: 16px;
  right: 16px;
}

.status-detail-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  margin: 0;
  padding: 20px 24px;
}

.status-detail-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.status-detail-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
}

.status-detail {
  padding: 0;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #e2e8f0;
}

.detail-code {
  font-size: 48px;
  font-weight: 700;
  font-family: 'Monaco', 'Menlo', monospace;
  line-height: 1;
}

.detail-info {
  flex: 1;
}

.detail-phrase {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
}

.detail-category {
  font-size: 14px;
  color: #6b7280;
}

.detail-content {
  padding: 24px;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.detail-section h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-section p {
  font-size: 14px;
  color: #4b5563;
  line-height: 1.6;
  margin: 0;
}

.detail-section ul {
  margin: 0;
  padding-left: 20px;
}

.detail-section li {
  font-size: 14px;
  color: #4b5563;
  line-height: 1.6;
  margin-bottom: 4px;
}

.headers-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.header-tag {
  font-family: 'Monaco', 'Menlo', monospace;
}

@media (max-width: 1024px) {
  .search-controls {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-input {
    max-width: none;
  }
  
  .category-select {
    width: 100%;
  }
  
  .stats-section {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }
  
  .status-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .detail-header {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
  
  .detail-code {
    font-size: 36px;
  }
}

@media (max-width: 768px) {
  .http-status-container {
    padding: 16px 0;
  }
  
  .category-header {
    padding: 20px;
  }
  
  .status-grid {
    padding: 20px;
  }
  
  .status-card {
    padding: 16px;
  }
}
</style>