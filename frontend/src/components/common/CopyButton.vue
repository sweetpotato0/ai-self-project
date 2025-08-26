<template>
  <el-button
    :type="buttonType"
    :size="size"
    :icon="Copy"
    :loading="isLoading"
    @click="handleCopy"
  >
    {{ buttonText }}
  </el-button>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Copy } from '@element-plus/icons-vue'
import { useClipboard } from '@/composables/useClipboard'

const props = defineProps({
  text: {
    type: String,
    required: true
  },
  buttonType: {
    type: String,
    default: 'primary'
  },
  size: {
    type: String,
    default: 'default'
  },
  successMessage: {
    type: String,
    default: '复制成功！'
  },
  buttonText: {
    type: String,
    default: '复制'
  }
})

const emit = defineEmits(['copy', 'error'])

const { copy, isLoading } = useClipboard()

const handleCopy = async () => {
  try {
    const success = await copy(props.text)
    if (success) {
      ElMessage.success(props.successMessage)
      emit('copy', props.text)
    }
  } catch (error) {
    ElMessage.error('复制失败：' + error.message)
    emit('error', error)
  }
}
</script>