import api from './index'

/**
 * 工具相关的API接口
 */
export const toolsApi = {
  /**
   * 端口扫描
   * @param {Object} data - 扫描参数
   * @param {string} data.target - 目标主机
   * @param {number[]} data.ports - 端口列表
   * @param {number} data.timeout - 超时时间(ms)
   * @param {number} data.concurrent - 并发数
   * @returns {Promise<Object>} 扫描结果
   */
  portScan: (data) => {
    return api.post('/tools/network/port-scan', data)
  }
}