<template>
  <el-alert
    v-if="visible"
    :title="title"
    :type="type"
    :description="description"
    :closable="closable"
    :center="center"
    show-icon
    @close="handleClose"
  >
    <template v-if="$slots.default">
      <slot />
    </template>
  </el-alert>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: '错误'
  },
  description: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'error',
    validator: (value) => ['success', 'warning', 'info', 'error'].includes(value)
  },
  closable: {
    type: Boolean,
    default: true
  },
  center: {
    type: Boolean,
    default: false
  },
  modelValue: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'close'])

const visible = ref(props.modelValue)

const handleClose = () => {
  visible.value = false
  emit('update:modelValue', false)
  emit('close')
}
</script>