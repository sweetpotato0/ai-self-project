<template>
  <div class="color-picker-container">
    <div class="tools-header">
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item @click="router.push({ name: 'tools' })">工具箱</el-breadcrumb-item>
          <el-breadcrumb-item @click="router.push({ name: 'tools-others' })">其它</el-breadcrumb-item>
          <el-breadcrumb-item>颜色选择器</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <h1>颜色选择器</h1>
      <p>专业的颜色选择和转换工具，支持多种颜色格式</p>
    </div>

    <div class="tools-content">
      <!-- 颜色选择区域 -->
      <div class="picker-section">
        <el-card class="picker-card">
          <template #header>
            <div class="card-header">
              <el-icon><Brush /></el-icon>
              <span>颜色选择</span>
            </div>
          </template>
          
          <div class="picker-content">
            <!-- 主颜色展示 -->
            <div class="color-display">
              <div 
                class="color-preview" 
                :style="{ backgroundColor: currentColor }"
                @click="openColorPicker"
              >
                <div class="color-overlay">
                  <span class="color-text">点击选择颜色</span>
                </div>
              </div>
              <input 
                ref="colorInput"
                type="color" 
                v-model="currentColor" 
                @change="updateColor"
                class="hidden-color-input"
              />
            </div>

            <!-- 预设颜色 -->
            <div class="preset-colors">
              <h4>预设颜色</h4>
              <div class="preset-grid">
                <div 
                  v-for="color in presetColors" 
                  :key="color"
                  class="preset-color"
                  :style="{ backgroundColor: color }"
                  @click="selectPresetColor(color)"
                  :class="{ active: color === currentColor }"
                >
                </div>
              </div>
            </div>

            <!-- 最近使用的颜色 -->
            <div class="recent-colors" v-if="recentColors.length > 0">
              <h4>最近使用</h4>
              <div class="recent-grid">
                <div 
                  v-for="color in recentColors" 
                  :key="color"
                  class="recent-color"
                  :style="{ backgroundColor: color }"
                  @click="selectPresetColor(color)"
                  :class="{ active: color === currentColor }"
                >
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 颜色信息区域 -->
      <div class="info-section">
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <el-icon><DataAnalysis /></el-icon>
              <span>颜色信息</span>
              <el-button @click="copyAllFormats" size="small" type="primary">
                <el-icon><DocumentCopy /></el-icon>
                复制全部
              </el-button>
            </div>
          </template>
          
          <div class="color-formats">
            <!-- HEX格式 -->
            <div class="format-item">
              <label>HEX:</label>
              <div class="format-value">
                <el-input v-model="hexValue" @blur="updateFromHex" />
                <el-button @click="copyFormat(hexValue)" size="small" text>
                  <el-icon><DocumentCopy /></el-icon>
                </el-button>
              </div>
            </div>

            <!-- RGB格式 -->
            <div class="format-item">
              <label>RGB:</label>
              <div class="format-value">
                <el-input v-model="rgbValue" @blur="updateFromRGB" />
                <el-button @click="copyFormat(rgbValue)" size="small" text>
                  <el-icon><DocumentCopy /></el-icon>
                </el-button>
              </div>
            </div>

            <!-- HSL格式 -->
            <div class="format-item">
              <label>HSL:</label>
              <div class="format-value">
                <el-input v-model="hslValue" @blur="updateFromHSL" />
                <el-button @click="copyFormat(hslValue)" size="small" text>
                  <el-icon><DocumentCopy /></el-icon>
                </el-button>
              </div>
            </div>

            <!-- HSV格式 -->
            <div class="format-item">
              <label>HSV:</label>
              <div class="format-value">
                <el-input v-model="hsvValue" readonly />
                <el-button @click="copyFormat(hsvValue)" size="small" text>
                  <el-icon><DocumentCopy /></el-icon>
                </el-button>
              </div>
            </div>

            <!-- CMYK格式 -->
            <div class="format-item">
              <label>CMYK:</label>
              <div class="format-value">
                <el-input v-model="cmykValue" readonly />
                <el-button @click="copyFormat(cmykValue)" size="small" text>
                  <el-icon><DocumentCopy /></el-icon>
                </el-button>
              </div>
            </div>

            <!-- CSS变量格式 -->
            <div class="format-item">
              <label>CSS变量:</label>
              <div class="format-value">
                <el-input v-model="cssVariableValue" readonly />
                <el-button @click="copyFormat(cssVariableValue)" size="small" text>
                  <el-icon><DocumentCopy /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 调色板区域 -->
      <div class="palette-section">
        <el-card class="palette-card">
          <template #header>
            <div class="card-header">
              <el-icon><Grid /></el-icon>
              <span>调色板</span>
            </div>
          </template>
          
          <div class="palette-content">
            <!-- 色调变化 -->
            <div class="palette-group">
              <h4>色调变化</h4>
              <div class="palette-colors">
                <div 
                  v-for="color in hueVariations" 
                  :key="color"
                  class="palette-color"
                  :style="{ backgroundColor: color }"
                  @click="selectPresetColor(color)"
                  :title="color"
                >
                </div>
              </div>
            </div>

            <!-- 明度变化 -->
            <div class="palette-group">
              <h4>明度变化</h4>
              <div class="palette-colors">
                <div 
                  v-for="color in lightnessVariations" 
                  :key="color"
                  class="palette-color"
                  :style="{ backgroundColor: color }"
                  @click="selectPresetColor(color)"
                  :title="color"
                >
                </div>
              </div>
            </div>

            <!-- 饱和度变化 -->
            <div class="palette-group">
              <h4>饱和度变化</h4>
              <div class="palette-colors">
                <div 
                  v-for="color in saturationVariations" 
                  :key="color"
                  class="palette-color"
                  :style="{ backgroundColor: color }"
                  @click="selectPresetColor(color)"
                  :title="color"
                >
                </div>
              </div>
            </div>

            <!-- 互补色 -->
            <div class="palette-group">
              <h4>配色方案</h4>
              <div class="palette-colors">
                <div 
                  v-for="color in complementaryColors" 
                  :key="color.name"
                  class="palette-color"
                  :style="{ backgroundColor: color.color }"
                  @click="selectPresetColor(color.color)"
                  :title="`${color.name}: ${color.color}`"
                >
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Brush,
  DataAnalysis,
  DocumentCopy,
  Grid
} from '@element-plus/icons-vue'

const router = useRouter()

// 响应式数据
const currentColor = ref('#3498db')
const colorInput = ref()
const recentColors = ref([])

// 预设颜色
const presetColors = ref([
  '#FF0000', '#FF8C00', '#FFD700', '#ADFF2F', '#00FF00', '#00CED1',
  '#0000FF', '#8A2BE2', '#FF1493', '#FF69B4', '#000000', '#808080',
  '#FFFFFF', '#F0F8FF', '#FAEBD7', '#00FFFF', '#7FFFD4', '#F0FFFF',
  '#F5F5DC', '#FFE4C4', '#FFEBCD', '#0000FF', '#8A2BE2', '#A52A2A'
])

// 颜色转换函数
const hexToRgb = (hex) => {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  return result ? {
    r: parseInt(result[1], 16),
    g: parseInt(result[2], 16),
    b: parseInt(result[3], 16)
  } : null
}

const rgbToHsl = (r, g, b) => {
  r /= 255; g /= 255; b /= 255
  const max = Math.max(r, g, b), min = Math.min(r, g, b)
  let h, s, l = (max + min) / 2

  if (max === min) {
    h = s = 0
  } else {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break
      case g: h = (b - r) / d + 2; break
      case b: h = (r - g) / d + 4; break
    }
    h /= 6
  }

  return {
    h: Math.round(h * 360),
    s: Math.round(s * 100),
    l: Math.round(l * 100)
  }
}

const rgbToHsv = (r, g, b) => {
  r /= 255; g /= 255; b /= 255
  const max = Math.max(r, g, b), min = Math.min(r, g, b)
  let h, s, v = max

  const d = max - min
  s = max === 0 ? 0 : d / max

  if (max === min) {
    h = 0
  } else {
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break
      case g: h = (b - r) / d + 2; break
      case b: h = (r - g) / d + 4; break
    }
    h /= 6
  }

  return {
    h: Math.round(h * 360),
    s: Math.round(s * 100),
    v: Math.round(v * 100)
  }
}

const rgbToCmyk = (r, g, b) => {
  r /= 255; g /= 255; b /= 255
  const k = 1 - Math.max(r, Math.max(g, b))
  const c = (1 - r - k) / (1 - k) || 0
  const m = (1 - g - k) / (1 - k) || 0
  const y = (1 - b - k) / (1 - k) || 0

  return {
    c: Math.round(c * 100),
    m: Math.round(m * 100),
    y: Math.round(y * 100),
    k: Math.round(k * 100)
  }
}

const hslToRgb = (h, s, l) => {
  h /= 360; s /= 100; l /= 100
  const hue2rgb = (p, q, t) => {
    if (t < 0) t += 1
    if (t > 1) t -= 1
    if (t < 1/6) return p + (q - p) * 6 * t
    if (t < 1/2) return q
    if (t < 2/3) return p + (q - p) * (2/3 - t) * 6
    return p
  }

  if (s === 0) {
    return { r: l, g: l, b: l }
  } else {
    const q = l < 0.5 ? l * (1 + s) : l + s - l * s
    const p = 2 * l - q
    return {
      r: Math.round(hue2rgb(p, q, h + 1/3) * 255),
      g: Math.round(hue2rgb(p, q, h) * 255),
      b: Math.round(hue2rgb(p, q, h - 1/3) * 255)
    }
  }
}

const rgbToHex = (r, g, b) => {
  return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)
}

// 计算属性
const hexValue = computed({
  get: () => currentColor.value.toUpperCase(),
  set: (val) => {
    if (/^#[0-9A-F]{6}$/i.test(val)) {
      currentColor.value = val
    }
  }
})

const rgbValue = computed(() => {
  const rgb = hexToRgb(currentColor.value)
  return rgb ? `rgb(${rgb.r}, ${rgb.g}, ${rgb.b})` : ''
})

const hslValue = computed(() => {
  const rgb = hexToRgb(currentColor.value)
  if (!rgb) return ''
  const hsl = rgbToHsl(rgb.r, rgb.g, rgb.b)
  return `hsl(${hsl.h}, ${hsl.s}%, ${hsl.l}%)`
})

const hsvValue = computed(() => {
  const rgb = hexToRgb(currentColor.value)
  if (!rgb) return ''
  const hsv = rgbToHsv(rgb.r, rgb.g, rgb.b)
  return `hsv(${hsv.h}, ${hsv.s}%, ${hsv.v}%)`
})

const cmykValue = computed(() => {
  const rgb = hexToRgb(currentColor.value)
  if (!rgb) return ''
  const cmyk = rgbToCmyk(rgb.r, rgb.g, rgb.b)
  return `cmyk(${cmyk.c}%, ${cmyk.m}%, ${cmyk.y}%, ${cmyk.k}%)`
})

const cssVariableValue = computed(() => {
  return `--primary-color: ${currentColor.value};`
})

// 调色板变化
const hueVariations = computed(() => {
  const rgb = hexToRgb(currentColor.value)
  if (!rgb) return []
  const hsl = rgbToHsl(rgb.r, rgb.g, rgb.b)
  const variations = []
  for (let i = 0; i < 12; i++) {
    const newHue = (hsl.h + i * 30) % 360
    const newRgb = hslToRgb(newHue, hsl.s, hsl.l)
    variations.push(rgbToHex(newRgb.r, newRgb.g, newRgb.b))
  }
  return variations
})

const lightnessVariations = computed(() => {
  const rgb = hexToRgb(currentColor.value)
  if (!rgb) return []
  const hsl = rgbToHsl(rgb.r, rgb.g, rgb.b)
  const variations = []
  for (let i = 0; i < 10; i++) {
    const newLightness = Math.max(5, Math.min(95, 10 + i * 10))
    const newRgb = hslToRgb(hsl.h, hsl.s, newLightness)
    variations.push(rgbToHex(newRgb.r, newRgb.g, newRgb.b))
  }
  return variations
})

const saturationVariations = computed(() => {
  const rgb = hexToRgb(currentColor.value)
  if (!rgb) return []
  const hsl = rgbToHsl(rgb.r, rgb.g, rgb.b)
  const variations = []
  for (let i = 0; i < 10; i++) {
    const newSaturation = i * 11
    const newRgb = hslToRgb(hsl.h, newSaturation, hsl.l)
    variations.push(rgbToHex(newRgb.r, newRgb.g, newRgb.b))
  }
  return variations
})

const complementaryColors = computed(() => {
  const rgb = hexToRgb(currentColor.value)
  if (!rgb) return []
  const hsl = rgbToHsl(rgb.r, rgb.g, rgb.b)
  
  const schemes = [
    { name: '原色', color: currentColor.value },
    { name: '互补色', color: rgbToHex(...Object.values(hslToRgb((hsl.h + 180) % 360, hsl.s, hsl.l))) },
    { name: '三元色1', color: rgbToHex(...Object.values(hslToRgb((hsl.h + 120) % 360, hsl.s, hsl.l))) },
    { name: '三元色2', color: rgbToHex(...Object.values(hslToRgb((hsl.h + 240) % 360, hsl.s, hsl.l))) },
    { name: '类似色1', color: rgbToHex(...Object.values(hslToRgb((hsl.h + 30) % 360, hsl.s, hsl.l))) },
    { name: '类似色2', color: rgbToHex(...Object.values(hslToRgb((hsl.h - 30 + 360) % 360, hsl.s, hsl.l))) }
  ]
  
  return schemes
})

// 方法
const openColorPicker = () => {
  colorInput.value.click()
}

const updateColor = () => {
  addToRecentColors(currentColor.value)
}

const selectPresetColor = (color) => {
  currentColor.value = color
  addToRecentColors(color)
}

const addToRecentColors = (color) => {
  const index = recentColors.value.indexOf(color)
  if (index > -1) {
    recentColors.value.splice(index, 1)
  }
  recentColors.value.unshift(color)
  if (recentColors.value.length > 12) {
    recentColors.value.pop()
  }
}

const updateFromHex = () => {
  if (/^#[0-9A-F]{6}$/i.test(hexValue.value)) {
    currentColor.value = hexValue.value
    addToRecentColors(hexValue.value)
  }
}

const updateFromRGB = () => {
  const match = rgbValue.value.match(/rgb\((\d+),\s*(\d+),\s*(\d+)\)/)
  if (match) {
    const r = parseInt(match[1])
    const g = parseInt(match[2])
    const b = parseInt(match[3])
    if (r <= 255 && g <= 255 && b <= 255) {
      currentColor.value = rgbToHex(r, g, b)
      addToRecentColors(currentColor.value)
    }
  }
}

const updateFromHSL = () => {
  const match = hslValue.value.match(/hsl\((\d+),\s*(\d+)%,\s*(\d+)%\)/)
  if (match) {
    const h = parseInt(match[1])
    const s = parseInt(match[2])
    const l = parseInt(match[3])
    if (h <= 360 && s <= 100 && l <= 100) {
      const rgb = hslToRgb(h, s, l)
      currentColor.value = rgbToHex(rgb.r, rgb.g, rgb.b)
      addToRecentColors(currentColor.value)
    }
  }
}

const copyFormat = async (value) => {
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const copyAllFormats = async () => {
  const allFormats = [
    `HEX: ${hexValue.value}`,
    `RGB: ${rgbValue.value}`,
    `HSL: ${hslValue.value}`,
    `HSV: ${hsvValue.value}`,
    `CMYK: ${cmykValue.value}`,
    `CSS: ${cssVariableValue.value}`
  ].join('\n')
  
  try {
    await navigator.clipboard.writeText(allFormats)
    ElMessage.success('已复制所有格式到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}
</script>

<style scoped>
.color-picker-container {
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
  grid-template-columns: 300px 1fr;
  grid-template-rows: auto auto;
  gap: 24px;
}

.picker-section {
  grid-row: span 2;
}

.picker-card,
.info-card,
.palette-card {
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

.card-header .el-icon + span {
  margin-left: 8px;
}

.picker-content {
  padding: 8px 0;
}

.color-display {
  position: relative;
  margin-bottom: 24px;
}

.color-preview {
  width: 100%;
  height: 120px;
  border-radius: 12px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  border: 2px solid #e5e7eb;
  transition: all 0.3s ease;
}

.color-preview:hover {
  border-color: #409eff;
  transform: scale(1.02);
}

.color-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.color-preview:hover .color-overlay {
  opacity: 1;
}

.color-text {
  color: white;
  font-weight: 500;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

.hidden-color-input {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

.preset-colors h4,
.recent-colors h4 {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 12px 0;
}

.preset-grid,
.recent-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
}

.preset-color,
.recent-color {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s ease;
}

.preset-color:hover,
.recent-color:hover {
  border-color: #409eff;
  transform: scale(1.1);
}

.preset-color.active,
.recent-color.active {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.recent-colors {
  margin-top: 24px;
}

.color-formats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.format-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.format-item label {
  min-width: 80px;
  font-weight: 500;
  color: #374151;
}

.format-value {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.palette-content {
  padding: 8px 0;
}

.palette-group {
  margin-bottom: 24px;
}

.palette-group:last-child {
  margin-bottom: 0;
}

.palette-group h4 {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 12px 0;
}

.palette-colors {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.palette-color {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  cursor: pointer;
  border: 1px solid rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
}

.palette-color:hover {
  transform: scale(1.2);
  border-color: #409eff;
  z-index: 1;
}

@media (max-width: 1024px) {
  .tools-content {
    grid-template-columns: 1fr;
    grid-template-rows: auto;
  }
  
  .picker-section {
    grid-row: auto;
  }
}

@media (max-width: 768px) {
  .color-picker-container {
    padding: 16px;
  }
  
  .tools-header h1 {
    font-size: 24px;
  }
  
  .format-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .format-item label {
    min-width: auto;
  }
  
  .format-value {
    width: 100%;
  }
}
</style>