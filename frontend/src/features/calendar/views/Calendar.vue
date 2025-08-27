<template>
  <div class="calendar-page">
    <!-- 炫酷的头部 -->
    <div class="calendar-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="header-title">📅 智能日程管理</h1>
          <p class="header-subtitle">可视化您的每一天，让时间更有价值</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" size="large" @click="showQuickAddDialog = true">
            <el-icon><Plus /></el-icon>
            快速添加
          </el-button>
          <el-button size="large" @click="goToTodos">
            <el-icon><List /></el-icon>
            任务列表
          </el-button>
        </div>
      </div>
    </div>

    <!-- 大日历主体 -->
    <div class="calendar-main">
      <!-- 日历导航栏 -->
      <div class="calendar-nav">
        <div class="nav-left">
          <el-button-group>
            <el-button @click="prevMonth" :icon="ArrowLeft" />
            <el-button @click="nextMonth" :icon="ArrowRight" />
          </el-button-group>
          <el-button @click="goToToday" type="primary" plain>今天</el-button>
        </div>

        <div class="nav-center">
          <h2 class="current-month">{{ currentMonthYear }}</h2>
        </div>

                 <div class="nav-right">
           <el-select v-model="viewMode" size="large" style="width: 120px">
             <el-option label="月视图" value="month" />
             <el-option label="周视图" value="week" />
             <el-option label="日视图" value="day" />
           </el-select>
         </div>
      </div>

             <!-- 月视图 -->
       <div v-if="viewMode === 'month'" class="calendar-grid">
         <!-- 星期标题 -->
         <div class="calendar-weekdays">
           <div
             v-for="day in weekdays"
             :key="day"
             class="weekday-header"
             :class="{ 'weekend': day === '周六' || day === '周日' }"
           >
             {{ day }}
           </div>
         </div>

         <!-- 日期网格 -->
         <div class="calendar-days">
           <div
             v-for="day in calendarDays"
             :key="day.date"
             class="calendar-day"
             :class="[
               { 'other-month': !day.isCurrentMonth },
               { 'today': day.isToday },
               { 'selected': day.isSelected },
               { 'weekend': day.isWeekend },
               { 'has-events': day.events.length > 0 }
             ]"
             @click="handleDateClick(day)"
             @dblclick="handleDateDoubleClick(day)"
           >
             <!-- 日期数字 -->
             <div class="day-number">
               <span class="date-number">{{ day.dayNumber }}</span>
               <span v-if="day.lunarDate" class="lunar-date">{{ day.lunarDate }}</span>
             </div>

             <!-- 事件列表 -->
             <div class="day-events">
               <div
                 v-for="(event, index) in day.events.slice(0, 3)"
                 :key="event.id"
                 class="day-event"
                 :class="getEventClass(event)"
                 @click.stop="viewEvent(event)"
               >
                 <div class="event-dot"></div>
                 <span class="event-title">{{ event.title }}</span>
                 <!-- 跨天任务连接指示器 -->
                 <div v-if="event.isMultiDay && !event.isLastDay" class="event-continue-indicator">
                   <el-icon><ArrowRight /></el-icon>
                 </div>
               </div>

               <!-- 更多事件指示器 -->
               <div
                 v-if="day.events.length > 3"
                 class="more-events"
                 @click.stop="viewAllEvents(day)"
               >
                 +{{ day.events.length - 3 }} 更多
               </div>
             </div>

             <!-- 添加按钮 -->
             <div class="day-add-btn" @click.stop="handleAddButtonClick(day)">
               <el-icon><Plus /></el-icon>
             </div>
           </div>
         </div>
       </div>

       <!-- 周视图 -->
       <div v-if="viewMode === 'week'" class="week-view">
         <div class="week-header">
           <div class="week-day-header">时间</div>
           <div
             v-for="day in weekDays"
             :key="day.date"
             class="week-day-header"
             :class="{ 'weekend': day.isWeekend, 'today': day.isToday, 'selected': day.isSelected }"
             @click="selectDate(day)"
           >
             <div class="week-day-name">{{ day.dayName }}</div>
             <div class="week-day-number">{{ day.dayNumber }}</div>
           </div>
         </div>

         <div class="week-content">
           <div class="time-column">
             <div v-for="hour in 24" :key="hour" class="time-slot">
               {{ (hour - 1).toString().padStart(2, '0') }}:00
             </div>
           </div>

           <div
             v-for="day in weekDays"
             :key="day.date"
             class="week-day-column"
             :class="{ 'weekend': day.isWeekend, 'today': day.isToday, 'selected': day.isSelected }"
           >
             <div
               v-for="hour in 24"
               :key="hour"
               class="time-slot"
               @click="handleTimeSlotClick(day, hour - 1)"
             >
               <div
                 v-for="event in getEventsForTimeSlot(day, hour - 1)"
                 :key="event.id"
                 class="week-event"
                 :class="getEventClass(event)"
                 @click.stop="viewEvent(event)"
               >
                 {{ event.title }}
                 <!-- 跨天任务连接指示器 -->
                 <div v-if="event.isMultiDay && !event.isLastDay" class="event-continue-indicator">
                   <el-icon><ArrowRight /></el-icon>
                 </div>
               </div>
             </div>
           </div>
         </div>
       </div>

       <!-- 日视图 -->
       <div v-if="viewMode === 'day'" class="day-view">
         <div class="day-header">
           <h3>{{ selectedDate.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }) }}</h3>
         </div>

         <div class="day-timeline">
           <div
             v-for="timeSlot in dayViewEvents"
             :key="timeSlot.hour"
             class="time-slot"
             :class="{ 'current-hour': timeSlot.hour === new Date().getHours() }"
           >
             <div class="time-label">{{ timeSlot.time }}</div>
             <div class="time-content">
               <div
                 v-for="event in timeSlot.events"
                 :key="event.id"
                 class="day-event"
                 :class="getEventClass(event)"
                 @click="viewEvent(event)"
               >
                 <div class="event-time">{{ formatEventTime(event) }}</div>
                 <div class="event-title">{{ event.title }}</div>
                 <div class="event-description">{{ event.description }}</div>
               </div>

               <!-- 添加按钮 -->
               <div
                 class="add-event-btn"
                 @click="handleTimeSlotClick({ date: selectedDate.toISOString().split('T')[0] }, timeSlot.hour)"
               >
                 <el-icon><Plus /></el-icon>
               </div>
             </div>
           </div>
         </div>
       </div>
    </div>

    <!-- 快速添加对话框 -->
    <el-dialog
      v-model="showQuickAddDialog"
      title="快速添加任务"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="quickAddForm" label-width="100px">
        <el-form-item label="任务标题" required>
          <el-input v-model="quickAddForm.title" placeholder="输入任务标题" />
        </el-form-item>

                <el-form-item label="全天任务">
          <el-switch v-model="quickAddForm.isAllDay" />
        </el-form-item>

        <!-- 全天任务：只需要一个日期 -->
        <el-form-item label="任务日期" v-if="quickAddForm.isAllDay" required>
          <el-date-picker
            v-model="quickAddForm.startDate"
            type="date"
            placeholder="选择任务日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>

        <!-- 非全天任务：开始日期+时间，结束日期+时间 -->
        <template v-if="!quickAddForm.isAllDay">
          <el-form-item label="开始时间" required>
            <div class="datetime-group">
              <el-date-picker
                v-model="quickAddForm.startDate"
                type="date"
                placeholder="开始日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 60%"
              />
              <el-time-picker
                v-model="quickAddForm.startTime"
                placeholder="开始时间"
                format="HH:mm"
                value-format="HH:mm"
                :disabled-hours="() => []"
                :disabled-minutes="() => []"
                style="width: 35%"
              />
            </div>
          </el-form-item>

          <el-form-item label="结束时间">
            <div class="datetime-group">
              <el-date-picker
                v-model="quickAddForm.endDate"
                type="date"
                placeholder="结束日期（可选，默认同一天）"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 60%"
              />
              <el-time-picker
                v-model="quickAddForm.endTime"
                placeholder="结束时间"
                format="HH:mm"
                value-format="HH:mm"
                :disabled-hours="() => []"
                :disabled-minutes="() => []"
                style="width: 35%"
              />
            </div>
          </el-form-item>

          <el-form-item label="时间预设">
            <div class="time-presets">
              <el-button
                v-for="preset in timePresets"
                :key="preset.label"
                size="small"
                @click="setTimePreset(preset)"
                :type="isTimePresetActive(preset) ? 'primary' : 'default'"
              >
                {{ preset.label }}
              </el-button>
            </div>
          </el-form-item>
        </template>

        <el-form-item label="优先级">
          <el-select v-model="quickAddForm.priority" placeholder="选择优先级">
            <el-option label="低" value="low" />
            <el-option label="中" value="medium" />
            <el-option label="高" value="high" />
            <el-option label="紧急" value="urgent" />
          </el-select>
        </el-form-item>

        <el-form-item label="分类">
          <hover-category-selector
            v-model="quickAddForm.category_id"
            :categories="todoStore.categories"
            placeholder="选择分类"
          />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="quickAddForm.description"
            type="textarea"
            :rows="3"
            placeholder="任务描述（可选）"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showQuickAddDialog = false">取消</el-button>
          <el-button type="primary" @click="addQuickTask">添加任务</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 事件详情对话框 -->
    <el-dialog
      v-model="showEventDialog"
      title="任务详情"
      width="700px"
    >
      <div v-if="selectedEvent" class="event-detail">
        <!-- 任务标题 -->
        <div class="event-header">
          <h3>{{ selectedEvent.title }}</h3>
          <div class="event-actions">
            <el-button size="small" type="primary" @click="editEvent(selectedEvent)">
              编辑
            </el-button>
            <el-button size="small" type="danger" @click="deleteEvent(selectedEvent)">
              删除
            </el-button>
          </div>
        </div>

        <!-- 任务描述 -->
        <div v-if="selectedEvent.description" class="event-section">
          <h4>任务描述</h4>
          <p class="event-description">{{ selectedEvent.description }}</p>
        </div>

        <!-- 任务状态和优先级 -->
        <div class="event-section">
          <h4>任务信息</h4>
          <div class="event-meta">
            <div class="meta-item">
              <span class="meta-label">状态:</span>
              <el-tag :type="getStatusType(selectedEvent.status)" size="small">
                {{ getStatusText(selectedEvent.status) }}
              </el-tag>
            </div>
            <div class="meta-item">
              <span class="meta-label">优先级:</span>
              <el-tag :type="getPriorityType(selectedEvent.priority)" size="small">
                {{ getPriorityText(selectedEvent.priority) }}
              </el-tag>
            </div>
          </div>
        </div>

        <!-- 时间信息 -->
        <div class="event-section">
          <h4>时间信息</h4>
          <div class="time-info">
            <div class="time-item">
              <span class="time-label">开始时间:</span>
              <span class="time-value">{{ formatEventDateTime(selectedEvent.start_date || selectedEvent.startDate) }}</span>
            </div>
            <div class="time-item">
              <span class="time-label">截止时间:</span>
              <span class="time-value" :class="{ 'overdue': isOverdue(selectedEvent) }">
                {{ formatEventDateTime(selectedEvent.due_date || selectedEvent.dueDate) }}
              </span>
            </div>
            <div v-if="selectedEvent.completed_at || selectedEvent.completedAt" class="time-item">
              <span class="time-label">完成时间:</span>
              <span class="time-value">{{ formatEventDateTime(selectedEvent.completed_at || selectedEvent.completedAt) }}</span>
            </div>
            <div v-if="selectedEvent.estimated_hours || selectedEvent.estimatedHours" class="time-item">
              <span class="time-label">预估工时:</span>
              <span class="time-value">{{ selectedEvent.estimated_hours || selectedEvent.estimatedHours }} 小时</span>
            </div>
            <div v-if="selectedEvent.actual_hours || selectedEvent.actualHours" class="time-item">
              <span class="time-label">实际工时:</span>
              <span class="time-value">{{ selectedEvent.actual_hours || selectedEvent.actualHours }} 小时</span>
            </div>
          </div>
        </div>

        <!-- 跨天信息 -->
        <div v-if="isMultiDayEvent(selectedEvent)" class="event-section">
          <h4>跨天信息</h4>
          <div class="multiday-info">
            <el-icon><Calendar /></el-icon>
            <span>这是一个跨天任务，持续 {{ getEventDuration(selectedEvent) }} 天</span>
          </div>
        </div>

        <!-- 创建信息 -->
        <div class="event-section">
          <h4>创建信息</h4>
          <div class="create-info">
            <div class="info-item">
              <span class="info-label">创建时间:</span>
              <span class="info-value">{{ formatEventDateTime(selectedEvent.created_at || selectedEvent.createdAt) }}</span>
            </div>
            <div v-if="selectedEvent.updated_at || selectedEvent.updatedAt" class="info-item">
              <span class="info-label">更新时间:</span>
              <span class="info-value">{{ formatEventDateTime(selectedEvent.updated_at || selectedEvent.updatedAt) }}</span>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="event-actions-bottom">
          <el-button
            v-if="selectedEvent.status !== 'completed'"
            type="success"
            @click="completeEvent(selectedEvent)"
          >
            标记完成
          </el-button>
          <el-button
            v-if="selectedEvent.status === 'pending'"
            type="warning"
            @click="startEvent(selectedEvent)"
          >
            开始任务
          </el-button>
          <el-button
            v-if="selectedEvent.status === 'in_progress'"
            type="info"
            @click="pauseEvent(selectedEvent)"
          >
            暂停任务
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTodoStore } from '@/stores/todo'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, ArrowLeft, ArrowRight, Calendar, List } from '@element-plus/icons-vue'
import HoverCategorySelector from '@/components/common/HoverCategorySelector.vue'

const router = useRouter()
const todoStore = useTodoStore()

// 响应式数据
const currentDate = ref(new Date())
const selectedDate = ref(new Date())
const viewMode = ref('month') // 'month', 'week', 'day'
const showQuickAddDialog = ref(false)
const showEventDialog = ref(false)
const selectedEvent = ref(null)

// 快速添加表单
// 时间预设
const timePresets = [
  { label: '上午', start: '09:00', end: '12:00' },
  { label: '下午', start: '14:00', end: '18:00' },
  { label: '晚上', start: '19:00', end: '22:00' },
  { label: '全天', start: '00:00', end: '23:59' },
  { label: '1小时', start: '09:00', end: '10:00' },
  { label: '2小时', start: '09:00', end: '11:00' },
  { label: '4小时', start: '09:00', end: '13:00' }
]

const quickAddForm = ref({
  title: '',
  date: '',
  startTime: '09:00',
  endTime: '10:00',
  isAllDay: false,
  startDate: '',
  endDate: '',
  priority: 'medium',
  description: '',
  category_id: null
})

// 计算属性
const currentMonthYear = computed(() => {
  if (viewMode.value === 'month') {
    return currentDate.value.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long'
    })
  } else if (viewMode.value === 'week') {
    const startOfWeek = new Date(selectedDate.value)
    startOfWeek.setDate(selectedDate.value.getDate() - selectedDate.value.getDay())
    const endOfWeek = new Date(startOfWeek)
    endOfWeek.setDate(startOfWeek.getDate() + 6)

    return `${startOfWeek.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' })} - ${endOfWeek.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', year: 'numeric' })}`
  } else if (viewMode.value === 'day') {
    return selectedDate.value.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      weekday: 'long'
    })
  }
  return ''
})

const selectedDateFormatted = computed(() => {
  return selectedDate.value.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
})

const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

// 辅助函数：获取本地日期字符串（不考虑时区）
const getLocalDateString = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  const firstDay = new Date(year, month, 1)
  const startDate = new Date(firstDay)
  startDate.setDate(startDate.getDate() - firstDay.getDay())

  const days = []
  const today = new Date()

  for (let i = 0; i < 42; i++) {
    const date = new Date(startDate)
    date.setDate(startDate.getDate() + i)

    const dayEvents = getEventsForDate(date)
    const isWeekend = date.getDay() === 0 || date.getDay() === 6

    days.push({
      date: getLocalDateString(date),
      dayNumber: date.getDate(),
      lunarDate: getLunarDate(date),
      isCurrentMonth: date.getMonth() === month,
      isToday: getLocalDateString(date) === getLocalDateString(today),
      isSelected: getLocalDateString(date) === getLocalDateString(selectedDate.value),
      isWeekend: isWeekend,
      events: dayEvents
    })
  }

  return days
})

const selectedDateEvents = computed(() => {
  return getEventsForDate(selectedDate.value)
})

// 周视图数据
const weekDays = computed(() => {
  if (viewMode.value !== 'week') return []

  const startOfWeek = new Date(selectedDate.value)
  startOfWeek.setDate(selectedDate.value.getDate() - selectedDate.value.getDay())

  const days = []
  for (let i = 0; i < 7; i++) {
    const date = new Date(startOfWeek)
    date.setDate(startOfWeek.getDate() + i)

    const dayEvents = getEventsForDate(date)
    const isWeekend = date.getDay() === 0 || date.getDay() === 6

    days.push({
      date: getLocalDateString(date),
      dayNumber: date.getDate(),
      dayName: date.toLocaleDateString('zh-CN', { weekday: 'short' }),
      lunarDate: getLunarDate(date),
      isToday: getLocalDateString(date) === getLocalDateString(new Date()),
      isSelected: getLocalDateString(date) === getLocalDateString(selectedDate.value),
      isWeekend: isWeekend,
      events: dayEvents
    })
  }

  return days
})

// 天视图数据
const dayViewEvents = computed(() => {
  if (viewMode.value !== 'day') return []

  const events = getEventsForDate(selectedDate.value)
  const timeSlots = []

  // 生成24小时时间段
  for (let hour = 0; hour < 24; hour++) {
    const timeSlot = {
      hour: hour,
      time: `${hour.toString().padStart(2, '0')}:00`,
      events: []
    }

    // 为每个时间段分配事件
    events.forEach(event => {
      const startTime = event.startDate || event.start_date ? new Date(event.startDate || event.start_date).getHours() : 9
      const endTime = event.dueDate || event.due_date ? new Date(event.dueDate || event.due_date).getHours() : 10

      if (hour >= startTime && hour <= endTime) {
        timeSlot.events.push(event)
      }
    })

    timeSlots.push(timeSlot)
  }

  return timeSlots
})

// 方法
const prevMonth = () => {
  if (viewMode.value === 'month') {
    currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() - 1, 1)
  } else if (viewMode.value === 'week') {
    // 切换到上一周
    const newDate = new Date(selectedDate.value)
    newDate.setDate(selectedDate.value.getDate() - 7)
    selectedDate.value = newDate
  } else if (viewMode.value === 'day') {
    // 切换到上一天
    const newDate = new Date(selectedDate.value)
    newDate.setDate(selectedDate.value.getDate() - 1)
    selectedDate.value = newDate
  }
}

const nextMonth = () => {
  if (viewMode.value === 'month') {
    currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 1)
  } else if (viewMode.value === 'week') {
    // 切换到下一周
    const newDate = new Date(selectedDate.value)
    newDate.setDate(selectedDate.value.getDate() + 7)
    selectedDate.value = newDate
  } else if (viewMode.value === 'day') {
    // 切换到下一天
    const newDate = new Date(selectedDate.value)
    newDate.setDate(selectedDate.value.getDate() + 1)
    selectedDate.value = newDate
  }
}

// 点击事件处理
let clickTimer = null

const handleDateClick = (day) => {
  // 清除之前的定时器
  if (clickTimer) {
    clearTimeout(clickTimer)
  }

  // 设置新的定时器，延迟处理单击事件
  clickTimer = setTimeout(() => {
    selectDate(day)
  }, 200) // 200ms延迟，避免与双击冲突
}

const handleDateDoubleClick = (day) => {
  // 清除单击定时器
  if (clickTimer) {
    clearTimeout(clickTimer)
    clickTimer = null
  }

  // 处理双击事件
  handleAddButtonClick(day)
}

const handleAddButtonClick = (day) => {
  showQuickAddDialog.value = true
  selectedDate.value = new Date(day.date)
  quickAddForm.value.startDate = day.date
  quickAddForm.value.endDate = day.date
  quickAddForm.value.isAllDay = true
  console.log('添加按钮点击，日期:', day.date, '设置到表单:', quickAddForm.value.startDate)
}

const selectDate = (day) => {
  // 修复日期选择偏移问题，直接使用传入的日期对象
  selectedDate.value = new Date(day.date)
  console.log('选择日期:', day.date, '实际选择:', selectedDate.value.toISOString().split('T')[0])
}

const goToToday = () => {
  const today = new Date()
  currentDate.value = new Date()
  selectedDate.value = new Date()

  // 根据视图模式调整显示
  if (viewMode.value === 'month') {
    // 月视图：跳转到当前月份
    currentDate.value = new Date(today.getFullYear(), today.getMonth(), 1)
  } else if (viewMode.value === 'week') {
    // 周视图：跳转到当前周
    selectedDate.value = new Date()
  } else if (viewMode.value === 'day') {
    // 日视图：跳转到今天
    selectedDate.value = new Date()
  }
}

const goToTodos = () => {
  router.push('/dashboard/todos')
}

const viewEvent = (event) => {
  selectedEvent.value = event
  showEventDialog.value = true
}

const viewAllEvents = (day) => {
  selectedDate.value = new Date(day.date)
  // 可以跳转到详细视图或显示更多事件
}

const handleTimeSlotClick = (day, hour) => {
  showQuickAddDialog.value = true
  selectedDate.value = new Date(day.date)
  quickAddForm.value.startDate = day.date
  quickAddForm.value.endDate = day.date
  quickAddForm.value.isAllDay = false
  quickAddForm.value.startTime = `${hour.toString().padStart(2, '0')}:00`
  quickAddForm.value.endTime = `${(hour + 1).toString().padStart(2, '0')}:00`
  console.log('时间段点击，日期:', day.date, '小时:', hour, '设置到表单:', quickAddForm.value.startDate, '时间:', quickAddForm.value.startTime)
}

// 设置时间预设
const setTimePreset = (preset) => {
  quickAddForm.value.startTime = preset.start
  quickAddForm.value.endTime = preset.end
}

// 检查时间预设是否激活
const isTimePresetActive = (preset) => {
  return quickAddForm.value.startTime === preset.start && quickAddForm.value.endTime === preset.end
}

// 格式化事件日期时间
const formatEventDateTime = (dateTime) => {
  if (!dateTime) return '未设置'
  const date = new Date(dateTime)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 检查是否超时
const isOverdue = (event) => {
  const dueDate = event.due_date || event.dueDate
  if (!dueDate) return false
  return new Date(dueDate) < new Date() && event.status !== 'completed'
}

// 检查是否为跨天事件
const isMultiDayEvent = (event) => {
  const startDate = event.start_date || event.startDate
  const dueDate = event.due_date || event.dueDate
  if (!startDate || !dueDate) return false

  const start = new Date(startDate)
  const end = new Date(dueDate)
  const startStr = getLocalDateString(start)
  const endStr = getLocalDateString(end)

  return startStr !== endStr
}

// 获取事件持续时间
const getEventDuration = (event) => {
  const startDate = event.start_date || event.startDate
  const dueDate = event.due_date || event.dueDate
  if (!startDate || !dueDate) return 1

  const start = new Date(startDate)
  const end = new Date(dueDate)
  const diffTime = Math.abs(end - start)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  return diffDays
}

// 事件操作函数
const editEvent = (event) => {
  // TODO: 实现编辑功能
  ElMessage.info('编辑功能开发中...')
}

const deleteEvent = async (event) => {
  try {
    await ElMessageBox.confirm('确定要删除这个任务吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const success = await todoStore.deleteTodo(event.id)
    if (success) {
      showEventDialog.value = false
      await todoStore.fetchTodos()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除任务失败')
    }
  }
}

const completeEvent = async (event) => {
  try {
    const success = await todoStore.updateTodo(event.id, { status: 'completed' })
    if (success) {
      showEventDialog.value = false
      await todoStore.fetchTodos()
      ElMessage.success('任务已完成')
    }
  } catch (error) {
    ElMessage.error('更新任务状态失败')
  }
}

const startEvent = async (event) => {
  try {
    const success = await todoStore.updateTodo(event.id, { status: 'in_progress' })
    if (success) {
      showEventDialog.value = false
      await todoStore.fetchTodos()
      ElMessage.success('任务已开始')
    }
  } catch (error) {
    ElMessage.error('更新任务状态失败')
  }
}

const pauseEvent = async (event) => {
  try {
    const success = await todoStore.updateTodo(event.id, { status: 'pending' })
    if (success) {
      showEventDialog.value = false
      await todoStore.fetchTodos()
      ElMessage.success('任务已暂停')
    }
  } catch (error) {
    ElMessage.error('更新任务状态失败')
  }
}

const getEventsForTimeSlot = (day, hour) => {
  const events = getEventsForDate(new Date(day.date))
  return events.filter(event => {
    const startTime = event.startDate || event.start_date ? new Date(event.startDate || event.start_date).getHours() : 9
    const endTime = event.dueDate || event.due_date ? new Date(event.dueDate || event.due_date).getHours() : 10
    return hour >= startTime && hour <= endTime
  })
}

const addQuickTask = async () => {
  if (!quickAddForm.value.title.trim()) {
    ElMessage.warning('请输入任务标题')
    return
  }

  if (!quickAddForm.value.startDate) {
    ElMessage.warning('请选择开始日期')
    return
  }

  try {
    const priorityMap = {
      'low': 1,
      'medium': 2,
      'high': 3,
      'urgent': 4
    }

    // 处理日期和时间格式
    let startDate = null
    let dueDate = null

        if (quickAddForm.value.isAllDay) {
      // 全天任务 - 只需要一个日期
      const [startYear, startMonth, startDay] = quickAddForm.value.startDate.split('-').map(Number)
      const startDateTime = new Date(startYear, startMonth - 1, startDay, 0, 0, 0)
      const endDateTime = new Date(startYear, startMonth - 1, startDay, 23, 59, 59)
      startDate = startDateTime.toISOString()
      dueDate = endDateTime.toISOString()
        } else {
      // 时间段任务
      const [startYear, startMonth, startDay] = quickAddForm.value.startDate.split('-').map(Number)
      const [startHour, startMinute] = quickAddForm.value.startTime.split(':').map(Number)

      const startDateTime = new Date(startYear, startMonth - 1, startDay, startHour, startMinute, 0)
      startDate = startDateTime.toISOString()

      // 检查是否有结束日期
      if (quickAddForm.value.endDate && quickAddForm.value.endDate !== quickAddForm.value.startDate) {
        // 跨天时间段任务
        const [endYear, endMonth, endDay] = quickAddForm.value.endDate.split('-').map(Number)
        const [endHour, endMinute] = quickAddForm.value.endTime.split(':').map(Number)

        const endDateTime = new Date(endYear, endMonth - 1, endDay, endHour, endMinute, 0)
        dueDate = endDateTime.toISOString()
      } else {
        // 同一天时间段任务
        const [endHour, endMinute] = quickAddForm.value.endTime.split(':').map(Number)
        const endDateTime = new Date(startYear, startMonth - 1, startDay, endHour, endMinute, 0)
        dueDate = endDateTime.toISOString()
      }
    }

    console.log('添加任务，开始时间:', startDate, '结束时间:', dueDate)

    const taskData = {
      title: quickAddForm.value.title,
      description: quickAddForm.value.description,
      priority_id: priorityMap[quickAddForm.value.priority] || 2,
      start_date: startDate,
      due_date: dueDate
    }

    await todoStore.createTodo(taskData)
    ElMessage.success('任务添加成功')
    showQuickAddDialog.value = false

    // 重新获取任务列表以更新日历显示
    await todoStore.fetchTodos()

    // 重置表单
    quickAddForm.value = {
      title: '',
      startDate: '',
      endDate: '',
      startTime: '09:00',
      endTime: '10:00',
      isAllDay: false,
      priority: 'medium',
      description: '',
      category_id: null
    }
  } catch (error) {
    console.error('添加任务失败:', error)
    ElMessage.error('添加任务失败: ' + (error.response?.data?.message || error.message))
  }
}

const getLunarDate = (date) => {
  // 简单的农历转换，实际项目中可以使用专门的农历库
  const lunarDates = ['初一', '初二', '初三', '初四', '初五', '初六', '初七', '初八', '初九', '初十',
                     '十一', '十二', '十三', '十四', '十五', '十六', '十七', '十八', '十九', '二十',
                     '廿一', '廿二', '廿三', '廿四', '廿五', '廿六', '廿七', '廿八', '廿九', '三十']
  return lunarDates[date.getDate() % 30]
}

const getEventClass = (event) => {
  // 根据优先级ID或名称确定样式类
  const priorityId = event.priority_id || event.priorityId
  const priorityName = event.priority

  let priorityClass = 'event-medium'

  if (priorityId) {
    // 根据优先级ID
    if (priorityId === 1) priorityClass = 'event-low'
    else if (priorityId === 2) priorityClass = 'event-medium'
    else if (priorityId === 3) priorityClass = 'event-high'
    else if (priorityId === 4 || priorityId === 5) priorityClass = 'event-urgent'
  } else if (priorityName) {
    // 根据优先级名称
    const priorityClasses = {
      'low': 'event-low',
      'medium': 'event-medium',
      'high': 'event-high',
      'urgent': 'event-urgent'
    }
    priorityClass = priorityClasses[priorityName] || 'event-medium'
  }

  // 添加跨天样式类
  if (event.isMultiDay) {
    priorityClass += ' multi-day'
    if (event.isFirstDay) {
      priorityClass += ' multi-day-start'
    } else if (event.isLastDay) {
      priorityClass += ' multi-day-end'
    } else if (event.isMiddleDay) {
      priorityClass += ' multi-day-middle'
    }
  }

  return priorityClass
}

const getEventsForDate = (date) => {
  const dateStr = getLocalDateString(date)
  const todos = todoStore.todos || []
  return todos.filter(todo => {
    // 兼容不同的字段名称
    const startDate = todo.startDate || todo.start_date ? new Date(todo.startDate || todo.start_date) : null
    const dueDate = todo.dueDate || todo.due_date ? new Date(todo.dueDate || todo.due_date) : null

    if (startDate && dueDate) {
      const startDateStr = getLocalDateString(startDate)
      const dueDateStr = getLocalDateString(dueDate)
      return dateStr >= startDateStr && dateStr <= dueDateStr
    } else if (startDate) {
      return dateStr === getLocalDateString(startDate)
    } else if (dueDate) {
      return dateStr === getLocalDateString(dueDate)
    }
    return false
  }).map(todo => {
    // 为每个任务添加跨天信息
    const startDate = todo.startDate || todo.start_date ? new Date(todo.startDate || todo.start_date) : null
    const dueDate = todo.dueDate || todo.due_date ? new Date(todo.dueDate || todo.due_date) : null

    if (startDate && dueDate) {
      const startDateStr = getLocalDateString(startDate)
      const dueDateStr = getLocalDateString(dueDate)
      const currentDateStr = getLocalDateString(date)

      return {
        ...todo,
        isMultiDay: startDateStr !== dueDateStr,
        isFirstDay: currentDateStr === startDateStr,
        isLastDay: currentDateStr === dueDateStr,
        isMiddleDay: currentDateStr > startDateStr && currentDateStr < dueDateStr
      }
    }

    return {
      ...todo,
      isMultiDay: false,
      isFirstDay: true,
      isLastDay: true,
      isMiddleDay: false
    }
  })
}

const formatEventTime = (event) => {
  const start = event.startDate || event.start_date ? new Date(event.startDate || event.start_date) : null
  const end = event.dueDate || event.due_date ? new Date(event.dueDate || event.due_date) : null

  if (start && end) {
    return `${start.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })} - ${end.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}`
  } else if (start) {
    return start.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  return '全天'
}

const getPriorityType = (priority) => {
  const typeMap = { 1: 'success', 2: 'warning', 3: 'danger', 4: 'danger', 5: 'danger' }
  return typeMap[priority] || 'info'
}

const getPriorityText = (priority) => {
  const textMap = { 1: '低', 2: '中', 3: '高', 4: '紧急', 5: '立即' }
  return textMap[priority] || '未知'
}

const getStatusType = (status) => {
  const typeMap = {
    'pending': 'info',
    'in_progress': 'warning',
    'completed': 'success',
    'cancelled': 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    'pending': '待处理',
    'in_progress': '进行中',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return textMap[status] || '未知'
}



// 生命周期
onMounted(async () => {
  await todoStore.fetchTodos()
  await todoStore.fetchCategories()
})
</script>

<style scoped>
.calendar-page {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

/* 炫酷头部 */
.calendar-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 30px;
  margin-bottom: 30px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  flex: 1;
}

.header-title {
  font-size: 2.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 10px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-subtitle {
  font-size: 1.1rem;
  color: #666;
  margin: 0;
  font-weight: 400;
}

.header-actions {
  display: flex;
  gap: 15px;
}

/* 大日历主体 */
.calendar-main {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

/* 日历导航栏 */
.calendar-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid rgba(102, 126, 234, 0.1);
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.nav-center {
  flex: 1;
  text-align: center;
}

.current-month {
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

/* 大日历表格 */
.calendar-grid {
  background: white;
  border-radius: 15px;
  overflow: visible;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  padding: 10px;
}

/* 星期标题 */
.calendar-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.weekday-header {
  padding: 20px 10px;
  text-align: center;
  font-weight: 600;
  font-size: 1.1rem;
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
}

.weekday-header.weekend {
  background: rgba(255, 255, 255, 0.1);
}

.weekday-header:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

/* 日期网格 */
.calendar-days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  background: white;
}

.calendar-day {
  min-height: 120px;
  padding: 15px 10px 15px 10px;
  border-right: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  background: white;
  overflow: visible;
}

.calendar-day:hover {
  background: linear-gradient(135deg, #f8f9ff 0%, #e8f2ff 100%);
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.2);
  z-index: 10;
}

.calendar-day.other-month {
  background: #fafafa;
  color: #ccc;
}

.calendar-day.today {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
  font-weight: 700;
  box-shadow: 0 8px 25px rgba(255, 107, 107, 0.4);
  border: 3px solid #fff;
  transform: scale(1.05);
  z-index: 20;
  position: relative;
}

.calendar-day.today::before {
  content: '今天';
  position: absolute;
  top: -8px;
  right: 2px;
  background: #ff6b6b;
  color: white;
  font-size: 0.7rem;
  padding: 2px 6px;
  border-radius: 10px;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  z-index: 30;
}

.calendar-day.selected {
  background: linear-gradient(135deg, #409eff 0%, #36a3f7 100%);
  color: white;
  font-weight: 700;
  box-shadow: 0 5px 15px rgba(64, 158, 255, 0.3);
}

.calendar-day.weekend {
  background: linear-gradient(135deg, #fff5f5 0%, #ffe8e8 100%);
}

.calendar-day.weekend:hover {
  background: linear-gradient(135deg, #ffe8e8 0%, #ffd6d6 100%);
}

.calendar-day.has-events {
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
}

.calendar-day.has-events:hover {
  background: linear-gradient(135deg, #e0f2fe 0%, #bae6fd 100%);
}

/* 日期数字 */
.day-number {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 10px;
}

.date-number {
  font-size: 1.2rem;
  font-weight: 600;
  margin-bottom: 2px;
}

.lunar-date {
  font-size: 0.7rem;
  color: #999;
  opacity: 0.8;
}

.calendar-day.today .lunar-date,
.calendar-day.selected .lunar-date {
  color: rgba(255, 255, 255, 0.8);
}

/* 事件列表 */
.day-events {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.day-event {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.7rem;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.day-event:hover {
  transform: scale(1.05);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.event-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.event-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 事件优先级样式 */
.event-urgent {
  background: rgba(245, 108, 108, 0.2);
  color: #f56c6c;
  border: 1px solid rgba(245, 108, 108, 0.3);
}

.event-urgent .event-dot {
  background: #f56c6c;
}

.event-high {
  background: rgba(230, 162, 60, 0.2);
  color: #e6a23c;
  border: 1px solid rgba(230, 162, 60, 0.3);
}

.event-high .event-dot {
  background: #e6a23c;
}

.event-medium {
  background: rgba(64, 158, 255, 0.2);
  color: #409eff;
  border: 1px solid rgba(64, 158, 255, 0.3);
}

.event-medium .event-dot {
  background: #409eff;
}

.event-low {
  background: rgba(103, 194, 58, 0.2);
  color: #67c23a;
  border: 1px solid rgba(103, 194, 58, 0.3);
}

.event-low .event-dot {
  background: #67c23a;
}

/* 跨天任务连接样式 */
.multi-day {
  position: relative;
  margin: 0 -2px;
  z-index: 2;
}

.multi-day-start {
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  margin-right: 0;
  position: relative;
}

.multi-day-start::after {
  content: '';
  position: absolute;
  top: 0;
  right: -2px;
  width: 4px;
  height: 100%;
  background: inherit;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
}

.multi-day-middle {
  border-radius: 0;
  margin: 0 -2px;
  position: relative;
}

.multi-day-middle::before {
  content: '';
  position: absolute;
  top: 0;
  left: -2px;
  width: 4px;
  height: 100%;
  background: inherit;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
}

.multi-day-middle::after {
  content: '';
  position: absolute;
  top: 0;
  right: -2px;
  width: 4px;
  height: 100%;
  background: inherit;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
}

.multi-day-end {
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
  margin-left: 0;
  position: relative;
}

.multi-day-end::before {
  content: '';
  position: absolute;
  top: 0;
  left: -2px;
  width: 4px;
  height: 100%;
  background: inherit;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
}

/* 时间预设样式 */
.time-range {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.time-separator {
  color: #666;
  font-weight: 500;
}

.time-presets {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 10px;
}

.time-presets .el-button {
  font-size: 0.8rem;
  padding: 4px 8px;
  border-radius: 6px;
}

/* 日期时间组合样式 */
.datetime-group {
  display: flex;
  gap: 10px;
  align-items: center;
}

.datetime-group .el-date-picker {
  flex: 1;
}

.datetime-group .el-time-picker {
  flex: 1;
}

/* 跨天任务连接指示器 */
.event-continue-indicator {
  position: absolute;
  right: -8px;
  top: 50%;
  transform: translateY(-50%);
  width: 16px;
  height: 16px;
  background: inherit;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: inherit;
  font-size: 8px;
  z-index: 3;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.event-continue-indicator::before {
  content: '';
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  width: 8px;
  height: 2px;
  background: inherit;
  border-radius: 1px;
}

/* 增强跨天任务的视觉效果 */
.multi-day {
  position: relative;
  margin: 0 -1px;
  z-index: 2;
  border-radius: 6px;
}

.multi-day-start {
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
  margin-right: 0;
  position: relative;
}

.multi-day-start::after {
  content: '';
  position: absolute;
  top: 0;
  right: -1px;
  width: 2px;
  height: 100%;
  background: inherit;
  border-top-right-radius: 1px;
  border-bottom-right-radius: 1px;
}

.multi-day-middle {
  border-radius: 0;
  margin: 0 -1px;
  position: relative;
}

.multi-day-middle::before {
  content: '';
  position: absolute;
  top: 0;
  left: -1px;
  width: 2px;
  height: 100%;
  background: inherit;
  border-top-left-radius: 1px;
  border-bottom-left-radius: 1px;
}

.multi-day-middle::after {
  content: '';
  position: absolute;
  top: 0;
  right: -1px;
  width: 2px;
  height: 100%;
  background: inherit;
  border-top-right-radius: 1px;
  border-bottom-right-radius: 1px;
}

.multi-day-end {
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
  margin-left: 0;
  position: relative;
}

.multi-day-end::before {
  content: '';
  position: absolute;
  top: 0;
  left: -1px;
  width: 2px;
  height: 100%;
  background: inherit;
  border-top-left-radius: 1px;
  border-bottom-left-radius: 1px;
}

/* 更多事件指示器 */
.more-events {
  font-size: 0.6rem;
  color: #999;
  text-align: center;
  padding: 2px;
  cursor: pointer;
  border-radius: 3px;
  transition: all 0.2s ease;
}

.more-events:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

/* 添加按钮 */
.day-add-btn {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 20px;
  height: 20px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0;
  transition: all 0.3s ease;
  color: #667eea;
}

.calendar-day:hover .day-add-btn {
  opacity: 1;
  transform: scale(1.1);
}

.day-add-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: scale(1.2);
}

/* 对话框样式 */
.el-dialog {
  border-radius: 15px;
  overflow: hidden;
}

.el-dialog__header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px 30px;
}

.el-dialog__title {
  color: white;
  font-weight: 600;
}

.el-dialog__body {
  padding: 30px;
}

.el-dialog__footer {
  padding: 20px 30px;
  border-top: 1px solid #f0f0f0;
}

/* 事件详情样式 */
.event-detail {
  padding: 20px;
}

.event-detail h3 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 1.5rem;
}

.event-description {
  color: #666;
  line-height: 1.6;
  margin-bottom: 20px;
}

.event-meta {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .calendar-page {
    padding: 10px;
  }

  .calendar-header {
    padding: 20px;
  }

  .header-title {
    font-size: 1.8rem;
  }

  .calendar-main {
    padding: 20px;
  }

  .current-month {
    font-size: 1.5rem;
  }

  .calendar-day {
    min-height: 80px;
    padding: 10px 5px;
  }

  .date-number {
    font-size: 1rem;
  }

  .day-event {
    font-size: 0.6rem;
    padding: 1px 4px;
  }
}

/* 动画效果 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.calendar-main {
  animation: fadeInUp 0.6s ease-out;
}

.calendar-day {
  animation: fadeInUp 0.3s ease-out;
}

/* 滚动条样式 */
.calendar-days::-webkit-scrollbar {
  width: 6px;
}

.calendar-days::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.calendar-days::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 3px;
}

.calendar-days::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%);
}

/* 周视图样式 */
.week-view {
  background: white;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.week-header {
  display: grid;
  grid-template-columns: 80px repeat(7, 1fr);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.week-day-header {
  padding: 15px 10px;
  text-align: center;
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
}

.week-day-header:hover {
  background: rgba(255, 255, 255, 0.1);
}

.week-day-header.weekend {
  background: rgba(255, 255, 255, 0.1);
}

.week-day-header.today {
  background: rgba(255, 107, 107, 0.3);
  font-weight: 700;
}

.week-day-header.selected {
  background: rgba(64, 158, 255, 0.3);
  font-weight: 700;
}

.week-day-name {
  font-size: 0.9rem;
  margin-bottom: 5px;
}

.week-day-number {
  font-size: 1.2rem;
  font-weight: 600;
}

.week-content {
  display: grid;
  grid-template-columns: 80px repeat(7, 1fr);
  height: 600px;
  overflow-y: auto;
}

.time-column {
  background: #f8f9fa;
  border-right: 1px solid #e9ecef;
}

.time-slot {
  height: 60px;
  padding: 5px;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  color: #666;
  position: relative;
}

.week-day-column {
  border-right: 1px solid #e9ecef;
  position: relative;
}

.week-day-column.weekend {
  background: rgba(255, 245, 245, 0.3);
}

.week-day-column.today {
  background: rgba(255, 107, 107, 0.1);
}

.week-day-column.selected {
  background: rgba(64, 158, 255, 0.1);
}

.week-day-column .time-slot {
  cursor: pointer;
  transition: all 0.2s ease;
}

.week-day-column .time-slot:hover {
  background: rgba(102, 126, 234, 0.1);
}

.week-event {
  background: rgba(64, 158, 255, 0.2);
  color: #409eff;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.7rem;
  margin: 1px 0;
  cursor: pointer;
  border: 1px solid rgba(64, 158, 255, 0.3);
  transition: all 0.2s ease;
}

.week-event:hover {
  background: rgba(64, 158, 255, 0.3);
  transform: scale(1.02);
}

/* 日视图样式 */
.day-view {
  background: white;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.day-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
  text-align: center;
}

.day-header h3 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.day-timeline {
  height: 600px;
  overflow-y: auto;
}

.day-timeline .time-slot {
  display: flex;
  border-bottom: 1px solid #e9ecef;
  min-height: 80px;
  transition: all 0.2s ease;
}

.day-timeline .time-slot:hover {
  background: rgba(102, 126, 234, 0.05);
}

.day-timeline .time-slot.current-hour {
  background: rgba(255, 107, 107, 0.1);
  border-left: 4px solid #ff6b6b;
}

.time-label {
  width: 80px;
  padding: 10px;
  background: #f8f9fa;
  border-right: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: #666;
  font-size: 0.9rem;
}

.time-content {
  flex: 1;
  padding: 10px;
  position: relative;
}

.day-event {
  background: rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.3);
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.day-event:hover {
  background: rgba(64, 158, 255, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.day-event .event-time {
  font-size: 0.8rem;
  color: #666;
  margin-bottom: 5px;
}

.day-event .event-title {
  font-weight: 600;
  color: #333;
  margin-bottom: 5px;
}

.day-event .event-description {
  font-size: 0.8rem;
  color: #666;
  line-height: 1.4;
}

.add-event-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 30px;
  height: 30px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0;
  transition: all 0.3s ease;
  color: #667eea;
}

.time-content:hover .add-event-btn {
  opacity: 1;
}

.add-event-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: scale(1.1);
}

/* 周期TODO的深色显示 */
.event-urgent {
  background: rgba(245, 108, 108, 0.3) !important;
  color: #f56c6c !important;
  border: 1px solid rgba(245, 108, 108, 0.5) !important;
}

.event-high {
  background: rgba(230, 162, 60, 0.3) !important;
  color: #e6a23c !important;
  border: 1px solid rgba(230, 162, 60, 0.5) !important;
}

.event-medium {
  background: rgba(64, 158, 255, 0.3) !important;
  color: #409eff !important;
  border: 1px solid rgba(64, 158, 255, 0.5) !important;
}

.event-low {
  background: rgba(103, 194, 58, 0.3) !important;
  color: #67c23a !important;
  border: 1px solid rgba(103, 194, 58, 0.5) !important;
}
</style>
