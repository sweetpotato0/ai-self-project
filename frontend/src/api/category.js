import request from '@/utils/request'

// 创建分类
export function createCategory(data) {
  return request({
    url: '/categories',
    method: 'post',
    data
  })
}

// 获取分类树
export function getCategoryTree() {
  return request({
    url: '/categories/tree',
    method: 'get'
  })
}

// 获取所有分类（扁平结构）
export function getAllCategories() {
  return request({
    url: '/categories',
    method: 'get'
  })
}

// 根据ID获取分类
export function getCategoryById(id) {
  return request({
    url: `/categories/${id}`,
    method: 'get'
  })
}

// 更新分类
export function updateCategory(id, data) {
  return request({
    url: `/categories/${id}`,
    method: 'put',
    data
  })
}

// 删除分类
export function deleteCategory(id) {
  return request({
    url: `/categories/${id}`,
    method: 'delete'
  })
}

// 获取分类统计
export function getCategoryStats() {
  return request({
    url: '/categories/stats',
    method: 'get'
  })
}
