<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="isEdit ? '编辑歌曲' : '新增歌曲'"
    width="800px"
    :before-close="handleClose"
    class="song-edit-dialog"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
      @submit.prevent
    >
      <el-tabs v-model="activeTab" class="form-tabs">
        <!-- 基础信息 -->
        <el-tab-pane label="基础信息" name="basic">
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="歌曲标题" prop="title" required>
                <el-input
                  v-model="formData.title"
                  placeholder="请输入英文歌曲标题"
                  clearable
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="中文标题" prop="title_cn">
                <el-input
                  v-model="formData.title_cn"
                  placeholder="请输入中文标题（可选）"
                  clearable
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="分类" prop="category_id" required>
                <el-select
                  v-model.number="formData.category_id"
                  placeholder="选择或输入歌曲分类"
                  filterable
                  allow-create
                  default-first-option
                  :reserve-keyword="false"
                  style="width: 100%"
                >
                  <el-option
                    v-for="category in categories"
                    :key="category.id"
                    :label="category.name_cn"
                    :value="category.id"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="难度等级" prop="difficulty" required>
                <el-select
                  v-model.number="formData.difficulty"
                  placeholder="选择难度等级"
                  style="width: 100%"
                >
                  <el-option label="入门 (1星)" :value="1" />
                  <el-option label="初级 (2星)" :value="2" />
                  <el-option label="中级 (3星)" :value="3" />
                  <el-option label="高级 (4星)" :value="4" />
                  <el-option label="专家 (5星)" :value="5" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="适合年龄" prop="age_range">
                <el-select
                  v-model="formData.age_range"
                  placeholder="选择适合年龄范围"
                  style="width: 100%"
                >
                  <el-option label="3-6岁" value="3-6" />
                  <el-option label="7-12岁" value="7-12" />
                  <el-option label="13-18岁" value="13-18" />
                  <el-option label="18岁以上" value="18+" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="时长" prop="duration">
                <el-input
                  v-model.number="formData.duration"
                  placeholder="歌曲时长（秒）"
                  type="number"
                  :min="1"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="描述" prop="description">
            <el-input
              v-model="formData.description"
              type="textarea"
              :rows="3"
              placeholder="请输入歌曲描述"
            />
          </el-form-item>

          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="发布状态">
                <el-switch
                  v-model="formData.is_published"
                  active-text="已发布"
                  inactive-text="草稿"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="排序">
                <el-input
                  v-model.number="formData.sort"
                  placeholder="排序值（数字越小越靠前）"
                  type="number"
                  :min="0"
                />
              </el-form-item>
            </el-col>
          </el-row>
        </el-tab-pane>

        <!-- 媒体文件 -->
        <el-tab-pane label="媒体文件" name="media">
          <el-form-item label="封面图片" prop="cover_image">
            <div class="media-upload">
              <div v-if="formData.cover_image" class="preview-container">
                <img :src="formData.cover_image" alt="封面预览" class="cover-preview" />
                <el-button
                  type="danger"
                  size="small"
                  :icon="Delete"
                  circle
                  class="remove-btn"
                  @click="formData.cover_image = ''"
                />
              </div>
              <el-input
                v-model="formData.cover_image"
                placeholder="请输入封面图片URL"
                clearable
              >
                <template #append>
                  <el-button :icon="Upload" @click="handleUploadCover">上传</el-button>
                </template>
              </el-input>
            </div>
          </el-form-item>

          <el-form-item label="音频文件" prop="audio_url">
            <div class="media-upload">
              <div v-if="formData.audio_url" class="audio-preview">
                <audio :src="formData.audio_url" controls class="audio-player" />
                <el-button
                  type="danger"
                  size="small"
                  :icon="Delete"
                  circle
                  class="remove-btn"
                  @click="formData.audio_url = ''"
                />
              </div>
              <el-input
                v-model="formData.audio_url"
                placeholder="请输入音频文件URL（可选）"
                clearable
              >
                <template #append>
                  <el-button :icon="Upload" @click="handleUploadAudio">上传</el-button>
                </template>
              </el-input>
            </div>
          </el-form-item>

          <el-form-item label="视频文件" prop="video_url">
            <div class="media-upload">
              <div v-if="formData.video_url" class="video-preview">
                <video :src="formData.video_url" controls class="video-player" />
                <el-button
                  type="danger"
                  size="small"
                  :icon="Delete"
                  circle
                  class="remove-btn"
                  @click="formData.video_url = ''"
                />
              </div>
              <el-input
                v-model="formData.video_url"
                placeholder="请输入视频文件URL（可选）"
                clearable
              >
                <template #append>
                  <el-button :icon="Upload" @click="handleUploadVideo">上传</el-button>
                </template>
              </el-input>
            </div>
          </el-form-item>
        </el-tab-pane>

        <!-- 歌词内容 -->
        <el-tab-pane label="歌词内容" name="lyrics">
          <el-form-item label="英文歌词" prop="lyrics" required>
            <el-input
              v-model="formData.lyrics"
              type="textarea"
              :rows="8"
              placeholder="请输入英文歌词&#10;可使用换行符分隔段落"
            />
          </el-form-item>

          <el-form-item label="中文翻译" prop="lyrics_cn">
            <el-input
              v-model="formData.lyrics_cn"
              type="textarea"
              :rows="8"
              placeholder="请输入中文歌词翻译（可选）&#10;建议与英文歌词行数对应"
            />
          </el-form-item>
        </el-tab-pane>

        <!-- 标签设置 -->
        <el-tab-pane label="标签设置" name="tags">
          <el-form-item label="歌曲标签" prop="tags">
            <div class="tags-input">
              <el-tag
                v-for="(tag, index) in formData.tags"
                :key="index"
                closable
                @close="removeTag(index)"
                class="tag-item"
              >
                {{ tag }}
              </el-tag>
              <el-input
                v-if="inputVisible"
                ref="tagInputRef"
                v-model="inputValue"
                size="small"
                class="tag-input"
                @keyup.enter="handleInputConfirm"
                @blur="handleInputConfirm"
              />
              <el-button
                v-else
                size="small"
                :icon="Plus"
                @click="showInput"
                class="add-tag-btn"
              >
                添加标签
              </el-button>
            </div>
            <div class="tags-help">
              <el-text type="info" size="small">
                建议标签：儿歌、动画、英语启蒙、语法、词汇、口语、听力等
              </el-text>
            </div>
          </el-form-item>
        </el-tab-pane>
      </el-tabs>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Upload, Delete } from '@element-plus/icons-vue'
import { useEnglishLearningStore } from '../../stores/englishLearningStore'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  song: {
    type: Object,
    default: null
  },
  categories: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const store = useEnglishLearningStore()

// 响应式数据
const formRef = ref()
const tagInputRef = ref()
const loading = ref(false)
const activeTab = ref('basic')
const inputVisible = ref(false)
const inputValue = ref('')

// 表单数据
const defaultFormData = {
  title: '',
  title_cn: '',
  description: '',
  lyrics: '',
  lyrics_cn: '',
  audio_url: '',
  video_url: '',
  cover_image: '',
  category_id: null,
  difficulty: 1,
  age_range: '',
  tags: [],
  is_published: false,
  duration: null,
  sort: 0
}

const formData = ref({ ...defaultFormData })

// 计算属性
const isEdit = computed(() => !!props.song)

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入歌曲标题', trigger: 'blur' },
    { min: 1, max: 200, message: '标题长度在 1 到 200 个字符', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  difficulty: [
    { required: true, message: '请选择难度等级', trigger: 'change' }
  ],
  audio_url: [
    { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' }
  ],
  lyrics: [
    { required: true, message: '请输入歌词内容', trigger: 'blur' }
  ]
}

// 监听器
watch(() => props.modelValue, (visible) => {
  if (visible) {
    initForm()
  }
})

watch(() => props.song, (song) => {
  if (song) {
    initForm()
  }
})

// 方法
const initForm = () => {
  if (props.song) {
    // 处理tags字段，确保它是数组格式
    let tags = []
    if (props.song.tags) {
      if (Array.isArray(props.song.tags)) {
        tags = props.song.tags
      } else if (typeof props.song.tags === 'string') {
        try {
          // 尝试解析JSON字符串
          const parsed = JSON.parse(props.song.tags)
          tags = Array.isArray(parsed) ? parsed : []
        } catch {
          // 如果不是JSON，按逗号分割
          tags = props.song.tags.split(',').map(tag => tag.trim()).filter(tag => tag)
        }
      }
    }
    
    // 只复制可编辑的字段，不包括嵌套的关联对象
    formData.value = {
      title: props.song.title || '',
      title_cn: props.song.title_cn || '',
      category_id: props.song.category_id || null,
      difficulty: props.song.difficulty || 1,
      age_range: props.song.age_range || '',
      duration: props.song.duration || 0,
      description: props.song.description || '',
      lyrics: props.song.lyrics || '',
      lyrics_cn: props.song.lyrics_cn || '',
      tags: tags,
      sort: props.song.sort || 0,
      is_published: props.song.is_published || false,
      cover_image: props.song.cover_image || '',
      audio_url: props.song.audio_url || '',
      video_url: props.song.video_url || ''
    }
  } else {
    formData.value = { ...defaultFormData }
  }
  activeTab.value = 'basic'
}

const handleClose = async () => {
  if (hasFormChanged()) {
    try {
      await ElMessageBox.confirm(
        '表单内容已修改，确定要关闭吗？',
        '确认关闭',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
    } catch {
      return
    }
  }
  emit('update:modelValue', false)
}

const hasFormChanged = () => {
  if (!props.song) {
    return Object.keys(formData.value).some(key => {
      if (key === 'tags') return Array.isArray(formData.value.tags) && formData.value.tags.length > 0
      return formData.value[key] !== defaultFormData[key]
    })
  }
  
  return Object.keys(formData.value).some(key => {
    if (key === 'tags') {
      // 特殊处理tags字段的比较
      const currentTags = Array.isArray(formData.value.tags) ? formData.value.tags : []
      const originalTags = props.song.tags ? 
        (Array.isArray(props.song.tags) ? props.song.tags : props.song.tags.split(',').map(tag => tag.trim()).filter(tag => tag)) : []
      return JSON.stringify(currentTags.sort()) !== JSON.stringify(originalTags.sort())
    }
    return JSON.stringify(formData.value[key]) !== JSON.stringify(props.song[key])
  })
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true

    const submitData = {
      ...formData.value,
      tags: Array.isArray(formData.value.tags) ? formData.value.tags.join(',') : formData.value.tags || ''
    }

    // 如果分类是字符串（用户输入的新分类），需要特殊处理
    if (typeof formData.value.category_id === 'string' && !props.categories.find(c => c.id === formData.value.category_id)) {
      // 这是一个新分类，发送分类名称而不是ID
      submitData.category_name = formData.value.category_id
      submitData.category_id = null
    }

    if (isEdit.value) {
      await store.updateSong(props.song.id, submitData)
      ElMessage.success('歌曲更新成功')
    } else {
      await store.createSong(submitData)
      ElMessage.success('歌曲创建成功')
    }

    emit('success')
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    loading.value = false
  }
}

// 标签管理
const removeTag = (index) => {
  formData.value.tags.splice(index, 1)
}

const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    tagInputRef.value?.focus()
  })
}

const handleInputConfirm = () => {
  const value = inputValue.value.trim()
  if (value && !formData.value.tags.includes(value)) {
    formData.value.tags.push(value)
  }
  inputVisible.value = false
  inputValue.value = ''
}

// 文件上传（这里可以根据实际需求实现文件上传逻辑）
const handleUploadCover = () => {
  ElMessage.info('文件上传功能待实现')
}

const handleUploadAudio = () => {
  ElMessage.info('文件上传功能待实现')
}

const handleUploadVideo = () => {
  ElMessage.info('文件上传功能待实现')
}
</script>

<style scoped>
.song-edit-dialog :deep(.el-dialog__body) {
  padding: 20px 24px;
}

.form-tabs {
  margin-top: -10px;
}

.form-tabs :deep(.el-tabs__content) {
  padding-top: 20px;
}

.media-upload {
  width: 100%;
}

.preview-container,
.audio-preview,
.video-preview {
  position: relative;
  margin-bottom: 12px;
  display: inline-block;
}

.cover-preview {
  width: 120px;
  height: 120px;
  object-fit: cover;
  border-radius: 8px;
  border: 1px solid #dcdfe6;
}

.audio-player {
  width: 100%;
  max-width: 400px;
  height: 40px;
}

.video-player {
  width: 100%;
  max-width: 400px;
  height: 200px;
}

.remove-btn {
  position: absolute;
  top: -8px;
  right: -8px;
  z-index: 1;
}

.tags-input {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.tag-item {
  margin: 0;
}

.tag-input {
  width: 120px;
}

.add-tag-btn {
  height: 24px;
  line-height: 1;
  border: 1px dashed #d9d9d9;
  background: transparent;
}

.tags-help {
  margin-top: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .song-edit-dialog :deep(.el-dialog) {
    width: 95vw !important;
    margin: 20px auto;
  }
  
  .form-tabs :deep(.el-tabs__nav-scroll) {
    padding: 0 10px;
  }
}
</style>