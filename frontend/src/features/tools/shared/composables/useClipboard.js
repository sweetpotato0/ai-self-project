import { ref } from 'vue'
import { ElMessage } from 'element-plus'

export function useClipboard() {
  const isSupported = ref(navigator && 'clipboard' in navigator)
  const isLoading = ref(false)
  const error = ref(null)

  const copy = async (text) => {
    if (!isSupported.value) {
      error.value = '您的浏览器不支持剪贴板API'
      ElMessage.error(error.value)
      return false
    }

    isLoading.value = true
    error.value = null

    try {
      await navigator.clipboard.writeText(text)
      ElMessage.success('复制成功！')
      return true
    } catch (err) {
      error.value = '复制失败：' + err.message
      ElMessage.error(error.value)
      return false
    } finally {
      isLoading.value = false
    }
  }

  const read = async () => {
    if (!isSupported.value) {
      error.value = '您的浏览器不支持剪贴板API'
      return null
    }

    isLoading.value = true
    error.value = null

    try {
      const text = await navigator.clipboard.readText()
      return text
    } catch (err) {
      error.value = '读取剪贴板失败：' + err.message
      return null
    } finally {
      isLoading.value = false
    }
  }

  return {
    isSupported,
    isLoading,
    error,
    copy,
    read
  }
}