<template>
  <div class="video-series-management">
    <!-- 剧集管理视图 -->
    <EpisodeManagement 
      v-if="showEpisodeManagement && selectedSeries" 
      :seriesId="selectedSeries.id" 
      @back="handleBackFromEpisodes" 
    />
    
    <!-- 系列管理视图 -->
    <div v-else>
      <div class="management-header">
        <h2>视频系列管理</h2>
        <el-button type="primary" :icon="Plus" @click="showCreateDialog">
          新建系列
        </el-button>
      </div>

    <!-- 系列列表 -->
    <el-card class="series-list">
      <el-table :data="seriesList" style="width: 100%">
        <el-table-column label="ID" width="80">
          <template #default="{ row }">
            <span v-if="row.id === 0" class="virtual-series">虚拟</span>
            <span v-else>{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image
              :src="row.cover_image"
              fit="cover"
              style="width: 60px; height: 40px; border-radius: 4px;"
            />
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="title_cn" label="中文标题" />
        <el-table-column prop="difficulty" label="难度" width="80">
          <template #default="{ row }">
            <el-tag :type="getDifficultyType(row.difficulty)">
              {{ row.difficulty }}级
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="age_range" label="年龄段" width="100" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_published ? 'success' : 'info'">
              {{ row.is_published ? '已发布' : '未发布' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="episode_count" label="剧集数" width="80" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <template v-if="row.id === 0">
                <!-- 虚拟未分类系列只显示管理剧集按钮 -->
                <el-button size="small" type="info" @click="manageEpisodes(row)">
                  管理剧集
                </el-button>
              </template>
              <template v-else>
                <!-- 正常系列显示所有操作 -->
                <el-button size="small" @click="editSeries(row)">
                  编辑
                </el-button>
                <el-button size="small" type="info" @click="manageEpisodes(row)">
                  管理剧集
                </el-button>
                <el-button size="small" type="danger" @click="deleteSeries(row)">
                  删除
                </el-button>
              </template>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      :title="isEditing ? '编辑系列' : '新建系列'"
      v-model="showDialog"
      width="600px"
    >
      <el-form
        :model="formData"
        :rules="formRules"
        ref="formRef"
        label-width="100px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="英文标题" />
        </el-form-item>
        <el-form-item label="中文标题" prop="title_cn">
          <el-input v-model="formData.title_cn" placeholder="中文标题" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="系列描述"
          />
        </el-form-item>
        <el-form-item label="封面图" prop="cover_image">
          <el-input v-model="formData.cover_image" placeholder="封面图片URL" />
        </el-form-item>
        <el-form-item label="难度等级" prop="difficulty">
          <el-select v-model="formData.difficulty" placeholder="选择难度">
            <el-option label="1级 (入门)" :value="1" />
            <el-option label="2级 (初级)" :value="2" />
            <el-option label="3级 (中级)" :value="3" />
            <el-option label="4级 (高级)" :value="4" />
            <el-option label="5级 (专家)" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item label="年龄段" prop="age_range">
          <el-input v-model="formData.age_range" placeholder="如: 3-8" />
        </el-form-item>
        <el-form-item label="标签" prop="tags">
          <el-input
            v-model="tagsInput"
            placeholder="用逗号分隔的标签，如: 小猪佩奇,英语启蒙,儿歌"
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="1" :max="1000" />
        </el-form-item>
        <el-form-item label="发布状态">
          <el-switch v-model="formData.is_published" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEditing ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
    </div> <!-- 结束系列管理视图 -->
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useEnglishVideosStore } from '../../stores/englishVideosStore'
import EpisodeManagement from './EpisodeManagement.vue'

const englishVideosStore = useEnglishVideosStore()

// 响应式数据
const seriesList = ref([])
const showDialog = ref(false)
const showEpisodeManagement = ref(false)
const selectedSeries = ref(null)
const isEditing = ref(false)
const submitting = ref(false)
const formRef = ref()

// 表单数据
const formData = reactive({
  id: null,
  title: '',
  title_cn: '',
  description: '',
  cover_image: '',
  difficulty: 1,
  age_range: '',
  sort: 1,
  is_published: true
})

const tagsInput = ref('')

// 表单验证规则
const formRules = {
  title: [
    { required: true, message: '请输入英文标题', trigger: 'blur' }
  ],
  title_cn: [
    { required: true, message: '请输入中文标题', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入描述', trigger: 'blur' }
  ],
  cover_image: [
    { required: true, message: '请输入封面图片URL', trigger: 'blur' }
  ],
  age_range: [
    { required: true, message: '请输入适合年龄段', trigger: 'blur' }
  ]
}

// 方法
const loadSeriesList = async () => {
  try {
    const response = await englishVideosStore.fetchVideoSeries({
      page: 1,
      page_size: 50,
      sort_order: 'desc'
    })
    seriesList.value = response || []
  } catch (error) {
    ElMessage.error('获取系列列表失败: ' + error.message)
  }
}

const showCreateDialog = () => {
  isEditing.value = false
  resetForm()
  showDialog.value = true
}

const editSeries = (series) => {
  isEditing.value = true
  Object.assign(formData, series)
  
  // 解析标签
  try {
    const tags = JSON.parse(series.tags || '[]')
    tagsInput.value = tags.join(', ')
  } catch {
    tagsInput.value = ''
  }
  
  showDialog.value = true
}

const resetForm = () => {
  Object.assign(formData, {
    id: null,
    title: '',
    title_cn: '',
    description: '',
    cover_image: '',
    difficulty: 1,
    age_range: '',
    sort: 1,
    is_published: true
  })
  tagsInput.value = ''
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true
  
  try {
    // 处理标签
    const tags = tagsInput.value
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag)
    
    const payload = {
      ...formData,
      tags: JSON.stringify(tags)
    }

    if (isEditing.value) {
      await englishVideosStore.updateVideoSeries(formData.id, payload)
      ElMessage.success('更新系列成功')
    } else {
      await englishVideosStore.createVideoSeries(payload)
      ElMessage.success('创建系列成功')
    }
    
    showDialog.value = false
    loadSeriesList()
  } catch (error) {
    ElMessage.error(`${isEditing.value ? '更新' : '创建'}系列失败: ${error.message}`)
  } finally {
    submitting.value = false
  }
}

const deleteSeries = async (series) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除系列"${series.title_cn}"吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await englishVideosStore.deleteVideoSeries(series.id)
    ElMessage.success('删除系列成功')
    loadSeriesList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除系列失败: ' + error.message)
    }
  }
}

const manageEpisodes = (series) => {
  selectedSeries.value = series
  showEpisodeManagement.value = true
}

const handleBackFromEpisodes = () => {
  showEpisodeManagement.value = false
  selectedSeries.value = null
}

const getDifficultyType = (difficulty) => {
  const types = ['', 'success', 'info', 'warning', 'danger', 'danger']
  return types[difficulty] || 'info'
}

// 生命周期
onMounted(() => {
  loadSeriesList()
})
</script>

<style scoped>
.video-series-management {
  padding: 20px;
}

.management-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.series-list {
  margin-top: 20px;
}

.el-table {
  margin-top: 10px;
}

.el-image {
  display: block;
}

.virtual-series {
  color: #909399;
  font-style: italic;
  font-size: 12px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: flex-start;
  align-items: center;
  flex-wrap: wrap;
}

.action-buttons .el-button {
  margin: 0;
  min-width: 60px;
}
</style>