<template>
  <el-dialog
    v-model="visible"
    title="图表生成器"
    width="800px"
    :before-close="handleClose"
  >
    <div class="chart-generator">
      <div class="chart-config">
        <el-form :model="chartConfig" label-width="100px">
          <el-form-item label="图表类型">
            <el-select v-model="chartConfig.type" @change="updateChart">
              <el-option label="柱状图" value="bar" />
              <el-option label="折线图" value="line" />
              <el-option label="饼图" value="pie" />
              <el-option label="散点图" value="scatter" />
            </el-select>
          </el-form-item>

          <el-form-item label="图表标题">
            <el-input v-model="chartConfig.title" @input="updateChart" />
          </el-form-item>

          <el-form-item label="数据">
            <el-input
              v-model="chartConfig.data"
              type="textarea"
              :rows="6"
              placeholder="输入JSON格式数据，例如：[{name: '类别1', value: 100}, {name: '类别2', value: 200}]"
              @input="updateChart"
            />
          </el-form-item>
        </el-form>
      </div>

      <div class="chart-preview">
        <h4>图表预览</h4>
        <div class="preview-container">
          <v-chart
            v-if="chartOption"
            :option="chartOption"
            style="width: 100%; height: 300px;"
          />
          <div v-else class="empty-preview">
            配置数据后查看图表预览
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="insertChart">插入图表</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, LineChart, PieChart, ScatterChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'

use([
  CanvasRenderer,
  BarChart,
  LineChart,
  PieChart,
  ScatterChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

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

const chartConfig = ref({
  type: 'bar',
  title: '示例图表',
  data: JSON.stringify([
    { name: '类别1', value: 100 },
    { name: '类别2', value: 200 },
    { name: '类别3', value: 150 },
    { name: '类别4', value: 300 }
  ], null, 2)
})

const chartOption = computed(() => {
  try {
    const data = JSON.parse(chartConfig.value.data)

    if (chartConfig.value.type === 'pie') {
      return {
        title: {
          text: chartConfig.value.title,
          left: 'center'
        },
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        series: [
          {
            name: chartConfig.value.title,
            type: 'pie',
            radius: '50%',
            data: data.map(item => ({
              name: item.name,
              value: item.value
            }))
          }
        ]
      }
    } else {
      return {
        title: {
          text: chartConfig.value.title,
          left: 'center'
        },
        tooltip: {
          trigger: 'axis'
        },
        xAxis: {
          type: 'category',
          data: data.map(item => item.name)
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: chartConfig.value.title,
            type: chartConfig.value.type,
            data: data.map(item => item.value)
          }
        ]
      }
    }
  } catch (error) {
    return null
  }
})

const updateChart = () => {
  // 图表会自动更新
}

const insertChart = () => {
  if (!chartOption.value) {
    ElMessage.warning('请配置有效的图表数据')
    return
  }

  const chartHtml = `<div class="chart-container" style="margin: 20px 0; padding: 20px; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">
    <h3 style="text-align: center; margin-bottom: 15px;">${chartConfig.value.title}</h3>
    <div class="chart-placeholder" style="height: 300px; display: flex; align-items: center; justify-content: center; background: #f8f9fa; border-radius: 4px;">
      <p style="color: #666;">图表: ${chartConfig.value.type} - ${chartConfig.value.title}</p>
    </div>
  </div>`

  emit('insert', chartHtml)
  handleClose()
}

const handleClose = () => {
  chartConfig.value = {
    type: 'bar',
    title: '示例图表',
    data: JSON.stringify([
      { name: '类别1', value: 100 },
      { name: '类别2', value: 200 },
      { name: '类别3', value: 150 },
      { name: '类别4', value: 300 }
    ], null, 2)
  }
  visible.value = false
}
</script>

<style scoped>
.chart-generator {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  padding: 20px 0;
}

.chart-config {
  border-right: 1px solid #e9ecef;
  padding-right: 20px;
}

.chart-preview h4 {
  margin: 0 0 15px 0;
  color: #333;
}

.preview-container {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  overflow: hidden;
}

.empty-preview {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-style: italic;
  background: #f8f9fa;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
