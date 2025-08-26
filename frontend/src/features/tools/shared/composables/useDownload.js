import { ref } from 'vue'
import { ElMessage } from 'element-plus'

export function useDownload() {
  const isLoading = ref(false)
  
  const downloadText = (text, filename = 'download.txt') => {
    try {
      const blob = new Blob([text], { type: 'text/plain;charset=utf-8' })
      downloadBlob(blob, filename)
      ElMessage.success('下载成功！')
    } catch (error) {
      ElMessage.error('下载失败：' + error.message)
    }
  }

  const downloadJson = (data, filename = 'data.json') => {
    try {
      const jsonString = JSON.stringify(data, null, 2)
      const blob = new Blob([jsonString], { type: 'application/json;charset=utf-8' })
      downloadBlob(blob, filename)
      ElMessage.success('下载成功！')
    } catch (error) {
      ElMessage.error('下载失败：' + error.message)
    }
  }

  const downloadCsv = (data, filename = 'data.csv') => {
    try {
      let csvContent = ''
      if (Array.isArray(data) && data.length > 0) {
        // 添加表头
        const headers = Object.keys(data[0])
        csvContent += headers.join(',') + '\n'
        
        // 添加数据行
        data.forEach(row => {
          const values = headers.map(header => {
            let value = row[header] || ''
            // 转义逗号和引号
            if (typeof value === 'string' && (value.includes(',') || value.includes('"'))) {
              value = `"${value.replace(/"/g, '""')}"`
            }
            return value
          })
          csvContent += values.join(',') + '\n'
        })
      }
      
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8' })
      downloadBlob(blob, filename)
      ElMessage.success('下载成功！')
    } catch (error) {
      ElMessage.error('下载失败：' + error.message)
    }
  }

  const downloadImage = (canvas, filename = 'image.png', quality = 0.9) => {
    try {
      canvas.toBlob(blob => {
        downloadBlob(blob, filename)
        ElMessage.success('下载成功！')
      }, 'image/png', quality)
    } catch (error) {
      ElMessage.error('下载失败：' + error.message)
    }
  }

  const downloadBlob = (blob, filename) => {
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  return {
    isLoading,
    downloadText,
    downloadJson,
    downloadCsv,
    downloadImage,
    downloadBlob
  }
}