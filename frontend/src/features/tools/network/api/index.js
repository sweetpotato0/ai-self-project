import api from '@/api/index'

/**
 * 网络工具相关的API接口
 */
export const networkApi = {
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
  },

  /**
   * DNS查询
   * @param {Object} data - 查询参数
   * @param {string} data.domain - 要查询的域名
   * @param {string} data.type - 记录类型 (A, AAAA, CNAME, MX, NS, TXT, PTR)
   * @returns {Promise<Object>} DNS查询结果
   */
  dnsLookup: (data) => {
    return api.post('/tools/network/dns-lookup', data)
  }
}