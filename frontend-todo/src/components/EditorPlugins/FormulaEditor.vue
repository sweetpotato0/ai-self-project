<template>
  <el-dialog
    v-model="visible"
    title="数学公式编辑器"
    width="600px"
    :before-close="handleClose"
  >
    <div class="formula-editor">
      <div class="formula-templates">
        <h4>常用公式模板</h4>
        <div class="template-grid">
          <el-button
            v-for="template in formulaTemplates"
            :key="template.name"
            size="small"
            @click="insertTemplate(template.formula)"
          >
            {{ template.name }}
          </el-button>
        </div>
      </div>

      <div class="formula-input">
        <h4>自定义公式</h4>
        <el-input
          v-model="formulaText"
          type="textarea"
          :rows="4"
          placeholder="输入LaTeX公式，例如: \frac{a}{b} = \sqrt{c^2 + d^2}"
        />
        <div class="formula-preview">
          <h5>预览</h5>
          <div class="preview-container" v-html="renderedFormula"></div>
        </div>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="insertFormula">插入公式</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'insert'])

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const formulaText = ref('')

const formulaTemplates = [
  { name: '分数', formula: '\\frac{a}{b}' },
  { name: '平方根', formula: '\\sqrt{x}' },
  { name: '求和', formula: '\\sum_{i=1}^{n} x_i' },
  { name: '积分', formula: '\\int_{a}^{b} f(x) dx' },
  { name: '矩阵', formula: '\\begin{pmatrix} a & b \\\\ c & d \\end{pmatrix}' },
  { name: '极限', formula: '\\lim_{x \\to \\infty} f(x)' },
  { name: '导数', formula: '\\frac{d}{dx} f(x)' },
  { name: '指数', formula: 'e^{x}' }
]

const renderedFormula = computed(() => {
  if (!formulaText.value) return '<div class="empty-preview">输入公式查看预览</div>'

  try {
    // 这里可以集成 MathJax 或 KaTeX 来渲染公式
    // 暂时使用简单的显示
    return `<div class="formula-display">$${formulaText.value}$</div>`
  } catch (error) {
    return '<div class="error-preview">公式格式错误</div>'
  }
})

const insertTemplate = (formula) => {
  formulaText.value = formula
}

const insertFormula = () => {
  if (!formulaText.value.trim()) {
    ElMessage.warning('请输入公式')
    return
  }

  const formulaHtml = `<div class="math-formula" style="text-align: center; margin: 15px 0; padding: 10px; background: #f8f9fa; border-radius: 8px; font-family: 'Times New Roman', serif;">$${formulaText.value}$</div>`

  emit('insert', formulaHtml)
  handleClose()
}

const handleClose = () => {
  formulaText.value = ''
  visible.value = false
}
</script>

<style scoped>
.formula-editor {
  padding: 20px 0;
}

.formula-templates {
  margin-bottom: 20px;
}

.formula-templates h4 {
  margin: 0 0 15px 0;
  color: #333;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 10px;
}

.formula-input h4 {
  margin: 0 0 15px 0;
  color: #333;
}

.formula-preview {
  margin-top: 15px;
}

.formula-preview h5 {
  margin: 0 0 10px 0;
  color: #666;
}

.preview-container {
  min-height: 60px;
  padding: 15px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  text-align: center;
}

.empty-preview {
  color: #999;
  font-style: italic;
}

.error-preview {
  color: #f56c6c;
  font-style: italic;
}

.formula-display {
  font-size: 16px;
  line-height: 1.5;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
