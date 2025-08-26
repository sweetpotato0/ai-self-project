// 工具相关路由配置
export const toolsRoutes = [
  {
    path: 'tools',
    name: 'tools',
    component: () => import('@/features/tools/views/Tools.vue')
  },
  // 开发工具
  {
    path: 'tools/development',
    name: 'tools-development',
    component: () => import('@/features/tools/development/views/ToolsDevelopment.vue')
  },
  {
    path: 'tools/timestamp-converter',
    name: 'tools-timestamp-converter',
    component: () => import('@/features/tools/development/TimestampConverter.vue')
  },
  {
    path: 'tools/json-tools',
    name: 'tools-json-tools',
    component: () => import('@/features/tools/development/JsonToolsSimple.vue')
  },
  {
    path: 'tools/string-generator',
    name: 'tools-string-generator',
    component: () => import('@/features/tools/development/StringGenerator.vue')
  },
  {
    path: 'tools/http-status-codes',
    name: 'tools-http-status-codes',
    component: () => import('@/features/tools/development/HttpStatusCodes.vue')
  },
  // 文本工具
  {
    path: 'tools/text',
    name: 'tools-text',
    component: () => import('@/features/tools/text/views/ToolsText.vue')
  },
  {
    path: 'tools/base64-encoder',
    name: 'tools-base64-encoder',
    component: () => import('@/features/tools/text/Base64Encoder.vue')
  },
  {
    path: 'tools/url-encoder',
    name: 'tools-url-encoder',
    component: () => import('@/features/tools/text/UrlEncoder.vue')
  },
  {
    path: 'tools/text-processor',
    name: 'tools-text-processor',
    component: () => import('@/features/tools/text/TextProcessor.vue')
  },
  // 图片工具
  {
    path: 'tools/image',
    name: 'tools-image',
    component: () => import('@/features/tools/image/views/ToolsImage.vue')
  },
  {
    path: 'tools/image-compressor',
    name: 'tools-image-compressor',
    component: () => import('@/features/tools/image/ImageCompressor.vue')
  },
  {
    path: 'tools/image-converter',
    name: 'tools-image-converter',
    component: () => import('@/features/tools/image/ImageConverter.vue')
  },
  {
    path: 'tools/image-resizer',
    name: 'tools-image-resizer',
    component: () => import('@/features/tools/image/ImageResizer.vue')
  },
  // 网络工具
  {
    path: 'tools/operations',
    name: 'tools-operations',
    component: () => import('@/features/tools/network/views/ToolsOperations.vue')
  },
  {
    path: 'tools/ping-tool',
    name: 'tools-ping-tool',
    component: () => import('@/features/tools/network/PingTool.vue')
  },
  // 查询工具
  {
    path: 'tools/query',
    name: 'tools-query',
    component: () => import('@/features/tools/query/views/ToolsQuery.vue')
  },
  {
    path: 'tools/ip-lookup',
    name: 'tools-ip-lookup',
    component: () => import('@/features/tools/query/IPLookup.vue')
  },
  {
    path: 'tools/whois-lookup',
    name: 'tools-whois-lookup',
    component: () => import('@/features/tools/query/WhoisLookup.vue')
  },
  {
    path: 'tools/domain-info',
    name: 'tools-domain-info',
    component: () => import('@/features/tools/query/DomainInfo.vue')
  },
  // 学术工具
  {
    path: 'tools/academic',
    name: 'tools-academic',
    component: () => import('@/features/tools/academic/views/ToolsAcademic.vue')
  },
  {
    path: 'tools/citation-generator',
    name: 'tools-citation-generator',
    component: () => import('@/features/tools/academic/CitationGenerator.vue')
  },
  {
    path: 'tools/math-calculator',
    name: 'tools-math-calculator',
    component: () => import('@/features/tools/academic/MathCalculator.vue')
  },
  // 文档工具
  {
    path: 'tools/document',
    name: 'tools-document',
    component: () => import('@/features/tools/document/views/ToolsDocument.vue')
  },
  // 其他工具
  {
    path: 'tools/others',
    name: 'tools-others',
    component: () => import('@/features/tools/others/views/ToolsOthers.vue')
  },
  {
    path: 'tools/qr-code-generator',
    name: 'tools-qr-code-generator',
    component: () => import('@/features/tools/others/QRCodeGenerator.vue')
  },
  {
    path: 'tools/color-picker',
    name: 'tools-color-picker',
    component: () => import('@/features/tools/others/ColorPicker.vue')
  },
  {
    path: 'tools/password-strength-checker',
    name: 'tools-password-strength-checker',
    component: () => import('@/features/tools/others/PasswordStrengthChecker.vue')
  },
  {
    path: 'tools/regex-tester',
    name: 'tools-regex-tester',
    component: () => import('@/features/tools/others/RegexTester.vue')
  },
  {
    path: 'tools/hash-calculator',
    name: 'tools-hash-calculator',
    component: () => import('@/features/tools/others/HashCalculator.vue')
  }
]