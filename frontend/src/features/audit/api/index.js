import api from '@/api/index'

/**
 * 审计日志相关的API接口
 */
export const auditApi = {
  /**
   * 获取审计日志列表
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.limit - 每页数量
   * @param {string} params.user_id - 用户ID
   * @param {string} params.username - 用户名
   * @param {string} params.action - 操作类型
   * @param {string} params.resource - 资源类型
   * @param {string} params.method - HTTP方法
   * @param {string} params.ip_address - IP地址
   * @param {number} params.status_code - 状态码
   * @param {string} params.start_time - 开始时间
   * @param {string} params.end_time - 结束时间
   * @param {string} params.order_by - 排序字段
   * @param {string} params.order - 排序方向
   * @returns {Promise} 审计日志列表
   */
  getAuditLogs(params = {}) {
    return api.get('/audit-logs', { params })
  },

  /**
   * 获取单个审计日志详情
   * @param {number} id - 审计日志ID
   * @returns {Promise} 审计日志详情
   */
  getAuditLogById(id) {
    return api.get(`/audit-logs/${id}`)
  },

  /**
   * 获取审计日志统计
   * @param {Object} params - 查询参数
   * @param {string} params.start_time - 开始时间
   * @param {string} params.end_time - 结束时间
   * @returns {Promise} 审计日志统计
   */
  getAuditLogStats(params = {}) {
    return api.get('/audit-logs/stats', { params })
  },

  /**
   * 批量删除审计日志
   * @param {Object} data - 删除参数
   * @param {string} data.before_time - 删除此时间之前的日志
   * @returns {Promise} 删除结果
   */
  deleteAuditLogs(data) {
    return api.delete('/audit-logs', { data })
  }
}