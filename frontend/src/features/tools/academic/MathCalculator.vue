<template>
  <div class="math-calculator-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-academic' })">学术类</el-breadcrumb-item>
          <el-breadcrumb-item>数学计算器</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>高级数学计算器</h1>
      <p>支持基础运算、三角函数、对数、幂次等数学计算</p>
    </div>

    <div class="tools-content">
      <div class="calculator-section">
        <div class="calculator-main">
          <div class="display-section">
            <div class="expression-display">{{ displayExpression }}</div>
            <div class="result-display">{{ displayResult }}</div>
          </div>
          
          <div class="buttons-section">
            <div class="function-buttons">
              <el-button @click="insertFunction('sin(')" class="func-btn">sin</el-button>
              <el-button @click="insertFunction('cos(')" class="func-btn">cos</el-button>
              <el-button @click="insertFunction('tan(')" class="func-btn">tan</el-button>
              <el-button @click="insertFunction('log(')" class="func-btn">log</el-button>
              <el-button @click="insertFunction('ln(')" class="func-btn">ln</el-button>
              <el-button @click="insertFunction('sqrt(')" class="func-btn">√</el-button>
              <el-button @click="insertOperator('^')" class="func-btn">x²</el-button>
              <el-button @click="insertConstant('π')" class="func-btn">π</el-button>
              <el-button @click="insertConstant('e')" class="func-btn">e</el-button>
            </div>
            
            <div class="main-buttons">
              <el-button @click="clearAll" class="clear-btn">C</el-button>
              <el-button @click="clearEntry" class="clear-btn">CE</el-button>
              <el-button @click="deleteLast" class="operator-btn">⌫</el-button>
              <el-button @click="insertOperator('/')" class="operator-btn">÷</el-button>
              
              <el-button @click="insertNumber('7')" class="number-btn">7</el-button>
              <el-button @click="insertNumber('8')" class="number-btn">8</el-button>
              <el-button @click="insertNumber('9')" class="number-btn">9</el-button>
              <el-button @click="insertOperator('*')" class="operator-btn">×</el-button>
              
              <el-button @click="insertNumber('4')" class="number-btn">4</el-button>
              <el-button @click="insertNumber('5')" class="number-btn">5</el-button>
              <el-button @click="insertNumber('6')" class="number-btn">6</el-button>
              <el-button @click="insertOperator('-')" class="operator-btn">-</el-button>
              
              <el-button @click="insertNumber('1')" class="number-btn">1</el-button>
              <el-button @click="insertNumber('2')" class="number-btn">2</el-button>
              <el-button @click="insertNumber('3')" class="number-btn">3</el-button>
              <el-button @click="insertOperator('+')" class="operator-btn">+</el-button>
              
              <el-button @click="insertNumber('0')" class="number-btn zero-btn">0</el-button>
              <el-button @click="insertOperator('.')" class="number-btn">.</el-button>
              <el-button @click="calculate" class="equals-btn">=</el-button>
            </div>
          </div>
        </div>

        <div class="history-section">
          <el-card class="history-card">
            <template #header>
              <div class="card-header">
                <el-icon><Clock /></el-icon>
                <span>计算历史</span>
                <el-button @click="clearHistory" size="small" text>清空</el-button>
              </div>
            </template>
            <div class="history-list">
              <div 
                v-for="(item, index) in history" 
                :key="index"
                class="history-item"
                @click="useHistoryItem(item)"
              >
                <div class="history-expression">{{ item.expression }}</div>
                <div class="history-result">= {{ item.result }}</div>
              </div>
              <div v-if="history.length === 0" class="no-history">
                暂无计算历史
              </div>
            </div>
          </el-card>
        </div>
      </div>

      <div class="info-section">
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>支持的函数</span>
            </div>
          </template>
          <div class="functions-grid">
            <div class="function-group">
              <h4>基础运算</h4>
              <ul>
                <li>+ - × ÷ (加减乘除)</li>
                <li>^ (幂次运算)</li>
                <li>() (括号)</li>
              </ul>
            </div>
            <div class="function-group">
              <h4>三角函数</h4>
              <ul>
                <li>sin() cos() tan()</li>
                <li>支持弧度制</li>
                <li>例: sin(π/2) = 1</li>
              </ul>
            </div>
            <div class="function-group">
              <h4>对数函数</h4>
              <ul>
                <li>log() (以10为底)</li>
                <li>ln() (自然对数)</li>
                <li>sqrt() (平方根)</li>
              </ul>
            </div>
            <div class="function-group">
              <h4>常数</h4>
              <ul>
                <li>π ≈ 3.14159</li>
                <li>e ≈ 2.71828</li>
              </ul>
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
  Clock,
  InfoFilled
} from '@element-plus/icons-vue'

const router = useRouter()

const expression = ref('')
const result = ref('')
const history = ref([])

const displayExpression = computed(() => {
  return expression.value || '0'
})

const displayResult = computed(() => {
  return result.value || ''
})

const insertNumber = (num) => {
  if (result.value && !isOperator(expression.value.slice(-1))) {
    // 如果有结果且最后不是操作符，开始新的计算
    expression.value = num
    result.value = ''
  } else {
    expression.value += num
  }
}

const insertOperator = (op) => {
  if (expression.value === '' && op === '-') {
    expression.value += op
    return
  }
  
  if (expression.value === '') return
  
  const lastChar = expression.value.slice(-1)
  if (isOperator(lastChar) && op !== '.') {
    // 替换最后的操作符
    expression.value = expression.value.slice(0, -1) + op
  } else {
    expression.value += op
  }
  result.value = ''
}

const insertFunction = (func) => {
  if (result.value && !isOperator(expression.value.slice(-1))) {
    expression.value = func
    result.value = ''
  } else {
    expression.value += func
  }
}

const insertConstant = (constant) => {
  const value = constant === 'π' ? 'π' : 'e'
  if (result.value && !isOperator(expression.value.slice(-1))) {
    expression.value = value
    result.value = ''
  } else {
    expression.value += value
  }
}

const isOperator = (char) => {
  return ['+', '-', '*', '/', '^', '.'].includes(char)
}

const clearAll = () => {
  expression.value = ''
  result.value = ''
}

const clearEntry = () => {
  result.value = ''
}

const deleteLast = () => {
  expression.value = expression.value.slice(0, -1)
  result.value = ''
}

const calculate = () => {
  if (!expression.value) return
  
  try {
    let expr = expression.value
    
    // 替换数学常数
    expr = expr.replace(/π/g, Math.PI.toString())
    expr = expr.replace(/e/g, Math.E.toString())
    
    // 替换数学函数
    expr = expr.replace(/sin\(/g, 'Math.sin(')
    expr = expr.replace(/cos\(/g, 'Math.cos(')
    expr = expr.replace(/tan\(/g, 'Math.tan(')
    expr = expr.replace(/log\(/g, 'Math.log10(')
    expr = expr.replace(/ln\(/g, 'Math.log(')
    expr = expr.replace(/sqrt\(/g, 'Math.sqrt(')
    
    // 替换运算符
    expr = expr.replace(/×/g, '*')
    expr = expr.replace(/÷/g, '/')
    expr = expr.replace(/\^/g, '**')
    
    // 计算结果
    const calcResult = Function('"use strict"; return (' + expr + ')')()
    
    if (isNaN(calcResult) || !isFinite(calcResult)) {
      throw new Error('计算错误')
    }
    
    const formattedResult = Number(calcResult.toPrecision(10)).toString()
    result.value = formattedResult
    
    // 添加到历史记录
    history.value.unshift({
      expression: expression.value,
      result: formattedResult,
      timestamp: new Date().toLocaleTimeString()
    })
    
    // 限制历史记录数量
    if (history.value.length > 20) {
      history.value.pop()
    }
    
  } catch (error) {
    ElMessage.error('计算表达式有误')
    result.value = 'Error'
  }
}

const useHistoryItem = (item) => {
  expression.value = item.expression
  result.value = item.result
}

const clearHistory = () => {
  history.value = []
  ElMessage.success('历史记录已清空')
}

// 键盘支持
const handleKeyPress = (event) => {
  const key = event.key
  
  if (key >= '0' && key <= '9') {
    insertNumber(key)
  } else if (['+', '-', '*', '/', '.'].includes(key)) {
    insertOperator(key)
  } else if (key === 'Enter' || key === '=') {
    calculate()
  } else if (key === 'Escape') {
    clearAll()
  } else if (key === 'Backspace') {
    deleteLast()
  }
}

// 添加键盘事件监听
if (typeof window !== 'undefined') {
  window.addEventListener('keydown', handleKeyPress)
}
</script>

<style scoped>
.math-calculator-container {
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

.calculator-section {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 24px;
  margin-bottom: 24px;
}

.calculator-main {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.display-section {
  background: #1f2937;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  min-height: 100px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.expression-display {
  color: #9ca3af;
  font-size: 16px;
  font-family: 'Courier New', monospace;
  word-break: break-all;
  margin-bottom: 8px;
}

.result-display {
  color: #fff;
  font-size: 32px;
  font-weight: 600;
  font-family: 'Courier New', monospace;
  text-align: right;
  word-break: break-all;
}

.function-buttons {
  display: grid;
  grid-template-columns: repeat(9, 1fr);
  gap: 8px;
  margin-bottom: 12px;
}

.main-buttons {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.func-btn {
  height: 45px;
  font-size: 14px;
  font-weight: 500;
}

.number-btn {
  height: 50px;
  font-size: 18px;
  font-weight: 500;
  background: #f8fafc;
}

.number-btn:hover {
  background: #f1f5f9;
}

.zero-btn {
  grid-column: span 2;
}

.operator-btn {
  height: 50px;
  font-size: 18px;
  font-weight: 500;
  background: #3b82f6;
  color: white;
}

.operator-btn:hover {
  background: #2563eb;
}

.clear-btn {
  height: 50px;
  font-size: 16px;
  font-weight: 500;
  background: #ef4444;
  color: white;
}

.clear-btn:hover {
  background: #dc2626;
}

.equals-btn {
  height: 50px;
  font-size: 20px;
  font-weight: 600;
  background: #10b981;
  color: white;
}

.equals-btn:hover {
  background: #059669;
}

.history-card {
  border: 1px solid #e5e7eb;
  height: fit-content;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: #2c3e50;
}

.history-list {
  max-height: 400px;
  overflow-y: auto;
}

.history-item {
  padding: 12px;
  border-bottom: 1px solid #f3f4f6;
  cursor: pointer;
  transition: background-color 0.2s;
}

.history-item:hover {
  background-color: #f8fafc;
}

.history-expression {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.history-result {
  font-size: 16px;
  font-weight: 500;
  color: #2c3e50;
  font-family: 'Courier New', monospace;
}

.no-history {
  text-align: center;
  color: #9ca3af;
  padding: 20px;
  font-style: italic;
}

.info-card {
  border: 1px solid #e5e7eb;
}

.functions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.function-group h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 12px 0;
}

.function-group ul {
  margin: 0;
  padding-left: 16px;
  list-style-type: disc;
}

.function-group li {
  font-size: 14px;
  color: #6b7280;
  margin: 6px 0;
}

@media (max-width: 768px) {
  .math-calculator-container {
    padding: 16px;
  }
  
  .calculator-section {
    grid-template-columns: 1fr;
  }
  
  .function-buttons {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .functions-grid {
    grid-template-columns: 1fr;
  }
  
  .result-display {
    font-size: 24px;
  }
}
</style>