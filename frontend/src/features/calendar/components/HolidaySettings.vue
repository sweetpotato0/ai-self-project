<template>
  <div class="holiday-settings">
    <el-tabs v-model="activeTab" class="settings-tabs">
      <!-- 订阅源管理 -->
      <el-tab-pane label="订阅源管理" name="sources">
        <div class="tab-content">
          <!-- 添加自定义订阅源 -->
          <div class="add-source-section">
            <div class="section-header">
              <h4>节日数据源订阅</h4>
              <el-button type="primary" size="small" @click="showAddSourceDialog = true">
                <el-icon><Plus /></el-icon>
                添加订阅
              </el-button>
            </div>
            
            <!-- 预定义数据源 -->
            <div class="sources-list">
              <h5>预定义数据源</h5>
              <div class="source-items">
                <div
                  v-for="source in predefinedSources"
                  :key="source.id"
                  class="source-item"
                  :class="{ active: source.enabled }"
                >
                  <div class="source-info">
                    <div class="source-header">
                      <h6>{{ source.name }}</h6>
                      <el-tag :type="getSourceTypeTag(source.type)" size="small">{{ source.type }}</el-tag>
                    </div>
                    <p>{{ source.description }}</p>
                    <div class="source-meta">
                      <span class="source-url">{{ source.url }}</span>
                      <span v-if="source.lastSync" class="sync-time">
                        上次同步: {{ formatSyncTime(source.lastSync) }}
                      </span>
                    </div>
                  </div>
                  <div class="source-actions">
                    <el-switch
                      v-model="source.enabled"
                      @change="toggleSource(source.id, $event)"
                    />
                    <el-button
                      size="small"
                      :loading="source.syncing"
                      @click="syncSource(source)"
                    >
                      同步
                    </el-button>
                  </div>
                </div>
              </div>

              <!-- 自定义数据源 -->
              <h5 v-if="customSources.length > 0">自定义数据源</h5>
              <div class="source-items" v-if="customSources.length > 0">
                <div
                  v-for="source in customSources"
                  :key="source.id"
                  class="source-item"
                  :class="{ active: source.enabled }"
                >
                  <div class="source-info">
                    <div class="source-header">
                      <h6>{{ source.name }}</h6>
                      <el-tag type="info" size="small">自定义</el-tag>
                    </div>
                    <p>{{ source.description || '无描述' }}</p>
                    <div class="source-meta">
                      <span class="source-url">{{ source.url }}</span>
                      <span v-if="source.lastSync" class="sync-time">
                        上次同步: {{ formatSyncTime(source.lastSync) }}
                      </span>
                    </div>
                  </div>
                  <div class="source-actions">
                    <el-switch
                      v-model="source.enabled"
                      @change="toggleSource(source.id, $event)"
                    />
                    <el-button
                      size="small"
                      :loading="source.syncing"
                      @click="syncSource(source)"
                    >
                      同步
                    </el-button>
                    <el-button
                      size="small"
                      type="danger"
                      @click="removeSource(source.id)"
                    >
                      删除
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- 节日预览 -->
      <el-tab-pane label="节日预览" name="preview">
        <div class="tab-content">
          <div class="preview-actions">
            <el-button @click="refreshHolidays" :loading="refreshing">
              <el-icon><Refresh /></el-icon>
              刷新数据
            </el-button>
            <el-button @click="clearCache">
              <el-icon><Delete /></el-icon>
              清除缓存
            </el-button>
          </div>

          <!-- 即将到来的节日 -->
          <div class="upcoming-section">
            <h4>即将到来的节日 (30天内)</h4>
            <div class="upcoming-holidays">
              <div 
                v-for="holiday in upcomingHolidays"
                :key="holiday.date"
                class="holiday-item"
                :class="`holiday-type-${holiday.type}`"
              >
                <div class="holiday-date">
                  <span class="month-day">{{ formatMonthDay(holiday.date) }}</span>
                  <span class="weekday">{{ formatWeekday(holiday.date) }}</span>
                </div>
                <div class="holiday-details">
                  <span class="holiday-name">{{ holiday.name }}</span>
                  <el-tag :type="getSourceTypeTag(holiday.type)" size="small">
                    {{ holiday.type }}
                  </el-tag>
                </div>
                <div class="days-until">
                  {{ getDaysUntil(holiday.date) }}天后
                </div>
              </div>
              <div v-if="upcomingHolidays.length === 0" class="no-holidays">
                <el-icon><InfoFilled /></el-icon>
                <span>接下来30天内没有节日数据</span>
              </div>
            </div>
          </div>

          <!-- 节日统计 -->
          <div class="stats-section">
            <h4>节日统计</h4>
            <div class="stats-grid">
              <div class="stat-item">
                <div class="stat-number">{{ holidayStats.total }}</div>
                <div class="stat-label">总节日数</div>
              </div>
              <div class="stat-item">
                <div class="stat-number">{{ holidayStats.legal }}</div>
                <div class="stat-label">法定节假日</div>
              </div>
              <div class="stat-item">
                <div class="stat-number">{{ holidayStats.traditional }}</div>
                <div class="stat-label">传统节日</div>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 添加数据源对话框 -->
    <el-dialog
      v-model="showAddSourceDialog"
      title="添加节日数据源"
      width="500px"
    >
      <el-form :model="newSourceForm" label-width="100px" ref="newSourceFormRef">
        <el-form-item label="数据源名称" required>
          <el-input 
            v-model="newSourceForm.name" 
            placeholder="请输入数据源名称"
            maxlength="50"
          />
        </el-form-item>
        
        <el-form-item label="数据源URL" required>
          <el-input 
            v-model="newSourceForm.url" 
            placeholder="https://api.example.com/holidays"
            type="url"
          />
        </el-form-item>
        
        <el-form-item label="数据源类型">
          <el-select v-model="newSourceForm.type" placeholder="请选择类型">
            <el-option label="法定节假日" value="legal" />
            <el-option label="传统节日" value="traditional" />
            <el-option label="国际节日" value="international" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述">
          <el-input 
            v-model="newSourceForm.description" 
            type="textarea" 
            placeholder="请输入数据源描述（可选）"
            :rows="3"
            maxlength="200"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showAddSourceDialog = false">取消</el-button>
          <el-button type="primary" @click="addCustomSource" :loading="adding">
            添加
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, Refresh, Delete, InfoFilled 
} from '@element-plus/icons-vue'
import { holidayManager } from '../utils/holidaySubscription.js'

// 响应式数据
const activeTab = ref('sources')
const showAddSourceDialog = ref(false)
const refreshing = ref(false)
const adding = ref(false)

const predefinedSources = ref([])
const customSources = ref([])
const allHolidays = ref([])

const newSourceForm = ref({
  name: '',
  url: '',
  type: 'custom',
  description: ''
})

// 计算属性
const upcomingHolidays = computed(() => {
  const now = new Date()
  const futureDate = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000) // 30天后
  
  return allHolidays.value.filter(holiday => {
    const holidayDate = new Date(holiday.date)
    return holidayDate >= now && holidayDate <= futureDate
  }).sort((a, b) => new Date(a.date) - new Date(b.date))
})

const holidayStats = computed(() => {
  const stats = {
    total: allHolidays.value.length,
    legal: 0,
    traditional: 0,
    international: 0,
    custom: 0
  }
  
  allHolidays.value.forEach(holiday => {
    if (stats[holiday.type] !== undefined) {
      stats[holiday.type]++
    }
  })
  
  return stats
})

// 方法
const loadSources = () => {
  predefinedSources.value = holidayManager.predefinedSources
  customSources.value = holidayManager.customSources
}

const toggleSource = (sourceId, enabled) => {
  holidayManager.toggleSource(sourceId, enabled)
  loadSources()
  if (enabled) {
    // 如果启用了数据源，自动同步一次
    const source = [...predefinedSources.value, ...customSources.value].find(s => s.id === sourceId)
    if (source) {
      syncSource(source)
    }
  }
}

const syncSource = async (source) => {
  source.syncing = true
  try {
    await holidayManager.fetchHolidayData(source)
    ElMessage.success(`${source.name} 同步成功`)
    loadSources()
    await loadHolidays()
  } catch (error) {
    ElMessage.error(`${source.name} 同步失败: ${error.message}`)
  } finally {
    source.syncing = false
  }
}

const removeSource = async (sourceId) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个数据源吗？删除后相关的节日数据也会被清除。',
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    holidayManager.removeCustomSource(sourceId)
    loadSources()
    await loadHolidays()
    ElMessage.success('数据源删除成功')
  } catch (error) {
    // 用户取消删除
  }
}

const addCustomSource = async () => {
  if (!newSourceForm.value.name || !newSourceForm.value.url) {
    ElMessage.error('请填写数据源名称和URL')
    return
  }
  
  adding.value = true
  try {
    const newSource = holidayManager.addCustomSource(newSourceForm.value)
    ElMessage.success('数据源添加成功')
    showAddSourceDialog.value = false
    
    // 重置表单
    newSourceForm.value = {
      name: '',
      url: '',
      type: 'custom',
      description: ''
    }
    
    loadSources()
    
    // 自动同步新添加的数据源
    await syncSource(newSource)
  } catch (error) {
    ElMessage.error(`添加数据源失败: ${error.message}`)
  } finally {
    adding.value = false
  }
}

const refreshHolidays = async () => {
  refreshing.value = true
  try {
    await holidayManager.getAllHolidays(true) // 强制刷新
    await loadHolidays()
    ElMessage.success('节日数据刷新成功')
  } catch (error) {
    ElMessage.error(`刷新失败: ${error.message}`)
  } finally {
    refreshing.value = false
  }
}

const clearCache = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清除所有缓存的节日数据吗？清除后需要重新同步数据源。',
      '确认清除缓存',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    holidayManager.clearCache()
    allHolidays.value = []
    ElMessage.success('缓存清除成功')
  } catch (error) {
    // 用户取消
  }
}

const loadHolidays = async () => {
  try {
    allHolidays.value = await holidayManager.getAllHolidays()
  } catch (error) {
    console.error('加载节日数据失败:', error)
  }
}

const getSourceTypeTag = (type) => {
  const typeMap = {
    'legal': 'danger',
    'traditional': 'warning',
    'international': 'primary',
    'custom': 'info'
  }
  return typeMap[type] || 'info'
}

const formatSyncTime = (timeString) => {
  if (!timeString) return '从未同步'
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN')
}

const formatMonthDay = (dateString) => {
  const date = new Date(dateString)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

const formatWeekday = (dateString) => {
  const date = new Date(dateString)
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return weekdays[date.getDay()]
}

const getDaysUntil = (dateString) => {
  const today = new Date()
  const targetDate = new Date(dateString)
  const diffTime = targetDate.getTime() - today.getTime()
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

// 生命周期
onMounted(async () => {
  loadSources()
  await loadHolidays()
})

// 暴露给父组件的方法
defineExpose({
  getAllHolidays: () => allHolidays.value,
  refreshHolidays
})
</script>

<style scoped>
.holiday-settings {
  max-width: 900px;
  margin: 0 auto;
}

.settings-tabs {
  margin-top: 16px;
}

.tab-content {
  padding: 16px 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.sources-list h5 {
  margin: 24px 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #606266;
}

.source-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.source-item {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  transition: all 0.2s ease;
}

.source-item:hover {
  background: #f1f3f4;
}

.source-item.active {
  border-color: #409eff;
  background: linear-gradient(135deg, #f0f9ff, #fff);
}

.source-info {
  flex: 1;
  margin-right: 16px;
}

.source-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.source-header h6 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.source-info p {
  margin: 0 0 8px 0;
  font-size: 13px;
  color: #606266;
  line-height: 1.4;
}

.source-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.source-url {
  font-size: 12px;
  color: #909399;
  font-family: monospace;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  word-break: break-all;
}

.sync-time {
  font-size: 11px;
  color: #c0c4cc;
}

.source-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-end;
}

.preview-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.upcoming-section h4,
.stats-section h4 {
  margin: 24px 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.upcoming-holidays {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.holiday-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #e9ecef;
}

.holiday-type-legal {
  border-left-color: #f56c6c;
}

.holiday-type-traditional {
  border-left-color: #e6a23c;
}

.holiday-type-international {
  border-left-color: #409eff;
}

.holiday-type-custom {
  border-left-color: #909399;
}

.holiday-date {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-right: 16px;
  min-width: 60px;
}

.month-day {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.weekday {
  font-size: 12px;
  color: #909399;
}

.holiday-details {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
}

.holiday-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.days-until {
  font-size: 12px;
  color: #606266;
  font-weight: 500;
}

.no-holidays {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 32px;
  color: #909399;
  font-size: 14px;
}

.stats-section {
  border-top: 1px solid #f0f0f0;
  padding-top: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  border-radius: 8px;
  background: linear-gradient(135deg, #f8f9ff, #fff);
  border: 1px solid #e8eaff;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #409eff;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .source-item {
    flex-direction: column;
    align-items: stretch;
  }
  
  .source-actions {
    flex-direction: row;
    justify-content: space-between;
    margin-top: 12px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>