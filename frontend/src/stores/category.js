import {
    createCategory,
    deleteCategory,
    getAllCategories,
    getCategoryById,
    getCategoryStats,
    getCategoryTree,
    updateCategory
} from '@/api/category'
import { ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useCategoryStore = defineStore('category', () => {
  // 状态
  const categories = ref([])
  const categoryTree = ref([])
  const categoryStats = ref({})
  const loading = ref(false)

  // 计算属性
  const flatCategories = computed(() => {
    const flatten = (cats, level = 0) => {
      let result = []
      cats.forEach(cat => {
        result.push({
          ...cat,
          level,
          indent: '　'.repeat(level) // 使用全角空格缩进
        })
        if (cat.children && cat.children.length > 0) {
          result = result.concat(flatten(cat.children, level + 1))
        }
      })
      return result
    }
    return flatten(categoryTree.value)
  })

  // 获取分类树
  const fetchCategoryTree = async () => {
    try {
      loading.value = true
      const response = await getCategoryTree()
      categoryTree.value = response.data || []
      return response.data
    } catch (error) {
      ElMessage.error('获取分类树失败')
      console.error('获取分类树失败:', error)
      return []
    } finally {
      loading.value = false
    }
  }

  // 获取所有分类（扁平结构）
  const fetchAllCategories = async () => {
    try {
      loading.value = true
      const response = await getAllCategories()
      categories.value = response.data || []
      return response.data
    } catch (error) {
      ElMessage.error('获取分类列表失败')
      console.error('获取分类列表失败:', error)
      return []
    } finally {
      loading.value = false
    }
  }

  // 创建分类
  const addCategory = async (categoryData) => {
    try {
      loading.value = true
      const response = await createCategory(categoryData)
      ElMessage.success('创建分类成功')
      // 重新获取分类树
      await fetchCategoryTree()
      return response.data
    } catch (error) {
      ElMessage.error('创建分类失败')
      console.error('创建分类失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新分类
  const editCategory = async (id, categoryData) => {
    try {
      loading.value = true
      const response = await updateCategory(id, categoryData)
      ElMessage.success('更新分类成功')
      // 重新获取分类树
      await fetchCategoryTree()
      return response.data
    } catch (error) {
      ElMessage.error('更新分类失败')
      console.error('更新分类失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除分类
  const removeCategory = async (id) => {
    try {
      loading.value = true
      await deleteCategory(id)
      ElMessage.success('删除分类成功')
      // 重新获取分类树
      await fetchCategoryTree()
    } catch (error) {
      ElMessage.error('删除分类失败')
      console.error('删除分类失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取分类统计
  const fetchCategoryStats = async () => {
    try {
      const response = await getCategoryStats()
      categoryStats.value = response.data || {}
      return response.data
    } catch (error) {
      console.error('获取分类统计失败:', error)
      return {}
    }
  }

  // 根据ID获取分类
  const getCategory = async (id) => {
    try {
      const response = await getCategoryById(id)
      return response.data
    } catch (error) {
      ElMessage.error('获取分类详情失败')
      console.error('获取分类详情失败:', error)
      return null
    }
  }

  // 刷新所有数据
  const refreshAll = async () => {
    await Promise.all([
      fetchCategoryTree(),
      fetchAllCategories(),
      fetchCategoryStats()
    ])
  }

  // 根据ID查找分类（在树中）
  const findCategoryById = (id, cats = categoryTree.value) => {
    for (const cat of cats) {
      if (cat.id === id) {
        return cat
      }
      if (cat.children && cat.children.length > 0) {
        const found = findCategoryById(id, cat.children)
        if (found) return found
      }
    }
    return null
  }

  // 获取分类路径（从根到当前分类）
  const getCategoryPath = (id, cats = categoryTree.value, path = []) => {
    for (const cat of cats) {
      const currentPath = [...path, cat]
      if (cat.id === id) {
        return currentPath
      }
      if (cat.children && cat.children.length > 0) {
        const found = getCategoryPath(id, cat.children, currentPath)
        if (found) return found
      }
    }
    return null
  }

  return {
    // 状态
    categories,
    categoryTree,
    categoryStats,
    loading,

    // 计算属性
    flatCategories,

    // 方法
    fetchCategoryTree,
    fetchAllCategories,
    addCategory,
    editCategory,
    removeCategory,
    fetchCategoryStats,
    getCategory,
    refreshAll,
    findCategoryById,
    getCategoryPath
  }
})
