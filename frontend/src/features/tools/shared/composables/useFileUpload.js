import { ref } from 'vue'
import { ElMessage } from 'element-plus'

export function useFileUpload() {
  const isLoading = ref(false)
  const uploadedFiles = ref([])

  const readFileAsText = (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result)
      reader.onerror = () => reject(new Error('文件读取失败'))
      reader.readAsText(file)
    })
  }

  const readFileAsDataURL = (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result)
      reader.onerror = () => reject(new Error('文件读取失败'))
      reader.readAsDataURL(file)
    })
  }

  const readFileAsArrayBuffer = (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result)
      reader.onerror = () => reject(new Error('文件读取失败'))
      reader.readAsArrayBuffer(file)
    })
  }

  const validateFile = (file, options = {}) => {
    const {
      maxSize = 10 * 1024 * 1024, // 默认10MB
      allowedTypes = [],
      allowedExtensions = []
    } = options

    // 检查文件大小
    if (file.size > maxSize) {
      const maxSizeMB = (maxSize / 1024 / 1024).toFixed(1)
      throw new Error(`文件大小不能超过 ${maxSizeMB}MB`)
    }

    // 检查文件类型
    if (allowedTypes.length > 0 && !allowedTypes.includes(file.type)) {
      throw new Error(`不支持的文件类型: ${file.type}`)
    }

    // 检查文件扩展名
    if (allowedExtensions.length > 0) {
      const extension = file.name.split('.').pop().toLowerCase()
      if (!allowedExtensions.includes(extension)) {
        throw new Error(`不支持的文件扩展名: .${extension}`)
      }
    }

    return true
  }

  const uploadFile = async (file, options = {}) => {
    isLoading.value = true
    
    try {
      validateFile(file, options)
      
      const fileData = {
        file,
        name: file.name,
        size: file.size,
        type: file.type,
        lastModified: file.lastModified
      }

      // 根据选项读取文件内容
      if (options.readAs === 'text') {
        fileData.content = await readFileAsText(file)
      } else if (options.readAs === 'dataURL') {
        fileData.content = await readFileAsDataURL(file)
      } else if (options.readAs === 'arrayBuffer') {
        fileData.content = await readFileAsArrayBuffer(file)
      }

      uploadedFiles.value.push(fileData)
      ElMessage.success('文件上传成功！')
      
      return fileData
    } catch (error) {
      ElMessage.error(error.message)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const uploadMultipleFiles = async (files, options = {}) => {
    const results = []
    for (const file of files) {
      try {
        const result = await uploadFile(file, options)
        results.push(result)
      } catch (error) {
        results.push({ error: error.message, file })
      }
    }
    return results
  }

  const clearFiles = () => {
    uploadedFiles.value = []
  }

  const removeFile = (index) => {
    uploadedFiles.value.splice(index, 1)
  }

  return {
    isLoading,
    uploadedFiles,
    uploadFile,
    uploadMultipleFiles,
    readFileAsText,
    readFileAsDataURL,
    readFileAsArrayBuffer,
    validateFile,
    clearFiles,
    removeFile
  }
}