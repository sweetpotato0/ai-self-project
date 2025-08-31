<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑文章' : '新建文章'"
    width="80%"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      class="article-form"
    >
      <el-row :gutter="20">
        <el-col :span="16">
          <el-form-item label="文章标题" prop="title">
            <el-input v-model="form.title" placeholder="请输入文章标题" />
          </el-form-item>

          <el-form-item label="文章内容" prop="content">
            <MarkdownEditor
              v-model="form.content"
              placeholder="请使用 Markdown 语法编写你的专业文章..."
              height="600px"
            />
          </el-form-item>
        </el-col>

        <el-col :span="8">
          <el-form-item label="文章摘要" prop="summary">
            <el-input
              v-model="form.summary"
              type="textarea"
              :rows="4"
              placeholder="请输入文章摘要"
            />
          </el-form-item>

          <el-form-item label="封面图片" prop="cover_image">
            <el-input v-model="form.cover_image" placeholder="请输入封面图片URL" />
          </el-form-item>

          <el-form-item label="文章状态" prop="status">
            <el-select v-model="form.status" placeholder="请选择文章状态">
              <el-option label="草稿" value="draft" />
              <el-option label="已发布" value="published" />
              <el-option label="已归档" value="archived" />
            </el-select>
          </el-form-item>

          <el-form-item label="文章标签" prop="tags">
            <el-select
              v-model="form.tags"
              multiple
              filterable
              allow-create
              default-first-option
              placeholder="请选择或输入标签"
            >
              <el-option
                v-for="tag in commonTags"
                :key="tag"
                :label="tag"
                :value="tag"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useArticleStore } from '@/stores/article'
import { ElMessage } from 'element-plus'
import MarkdownEditor from '@/components/editor/MarkdownEditor.vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  article: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const articleStore = useArticleStore()

// 响应式数据
const formRef = ref()
const submitting = ref(false)

// 表单数据
const form = reactive({
  title: '',
  content: '',
  summary: '',
  cover_image: '',
  status: 'draft',
  tags: []
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择文章状态', trigger: 'change' }
  ]
}

// 常用标签
const commonTags = [
  '技术', '编程', '前端', '后端', 'Vue', 'React', 'JavaScript', 'Python',
  'Go', '数据库', '架构', '设计', '生活', '随笔', '教程', '分享'
]

// 计算属性
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const isEdit = computed(() => !!props.article)

// 方法
const resetForm = () => {
  // 重置表单数据
  Object.assign(form, {
    title: '',
    content: '',
    summary: '',
    cover_image: '',
    status: 'draft',
    tags: []
  })
  
  // 重置表单验证状态
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

// 监听文章变化，填充表单
watch(() => props.article, (newArticle) => {
  if (newArticle) {
    form.title = newArticle.title || ''
    form.content = newArticle.content || ''
    form.summary = newArticle.summary || ''
    form.cover_image = newArticle.cover_image || ''
    form.status = newArticle.status || 'draft'

    // 解析标签
    if (newArticle.tags) {
      try {
        form.tags = JSON.parse(newArticle.tags)
      } catch (e) {
        form.tags = []
      }
    } else {
      form.tags = []
    }
  } else {
    // 新建文章时重置表单
    resetForm()
  }
}, { immediate: true })

// 监听对话框打开状态，确保新建时表单被重置
watch(() => props.modelValue, (newValue) => {
  if (newValue && !props.article) {
    // 对话框打开且没有编辑文章时，重置表单
    resetForm()
  }
})

const handleClose = () => {
  resetForm()
  dialogVisible.value = false
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    const articleData = {
      title: form.title,
      content: form.content,
      summary: form.summary,
      cover_image: form.cover_image,
      status: form.status,
      tags: form.tags
    }

    if (isEdit.value) {
      await articleStore.updateArticle(props.article.id, articleData)
      ElMessage.success('文章更新成功')
    } else {
      await articleStore.createArticle(articleData)
      ElMessage.success('文章创建成功')
    }

    emit('success')
    resetForm()
    dialogVisible.value = false
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.article-form {
  max-height: 70vh;
  overflow-y: auto;
}

.dialog-footer {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .el-dialog {
    width: 95% !important;
  }

  .el-col {
    width: 100% !important;
  }
}
</style>

<style scoped>
.article-form {
  max-height: 70vh;
  overflow-y: auto;
}

.cover-preview {
  margin-top: 10px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e4e7ed;
}

.cover-preview img {
  width: 100%;
  height: 120px;
  object-fit: cover;
}

.dialog-footer {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .el-dialog {
    width: 95% !important;
  }

  .el-col {
    width: 100% !important;
  }
}
</style>
