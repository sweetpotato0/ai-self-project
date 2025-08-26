/**
 * 验证邮箱格式
 * @param {string} email - 邮箱地址
 * @returns {boolean} 是否为有效邮箱
 */
export function isValidEmail(email) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * 验证URL格式
 * @param {string} url - URL地址
 * @returns {boolean} 是否为有效URL
 */
export function isValidURL(url) {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * 验证IP地址格式
 * @param {string} ip - IP地址
 * @returns {boolean} 是否为有效IP地址
 */
export function isValidIP(ip) {
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  const ipv6Regex = /^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/
  return ipv4Regex.test(ip) || ipv6Regex.test(ip)
}

/**
 * 验证域名格式
 * @param {string} domain - 域名
 * @returns {boolean} 是否为有效域名
 */
export function isValidDomain(domain) {
  const domainRegex = /^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/i
  return domainRegex.test(domain)
}

/**
 * 验证手机号码格式（中国）
 * @param {string} phone - 手机号码
 * @returns {boolean} 是否为有效手机号码
 */
export function isValidPhone(phone) {
  const phoneRegex = /^1[3-9]\d{9}$/
  return phoneRegex.test(phone)
}

/**
 * 验证身份证号码格式（中国）
 * @param {string} idCard - 身份证号码
 * @returns {boolean} 是否为有效身份证号码
 */
export function isValidIDCard(idCard) {
  const idCardRegex = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/
  return idCardRegex.test(idCard)
}

/**
 * 验证密码强度
 * @param {string} password - 密码
 * @returns {object} 验证结果
 */
export function validatePasswordStrength(password) {
  const result = {
    score: 0,
    level: 'weak',
    suggestions: []
  }

  if (password.length < 8) {
    result.suggestions.push('密码长度至少8位')
  } else if (password.length >= 8) {
    result.score += 1
  }

  if (password.length >= 12) {
    result.score += 1
  }

  if (/[a-z]/.test(password)) {
    result.score += 1
  } else {
    result.suggestions.push('包含小写字母')
  }

  if (/[A-Z]/.test(password)) {
    result.score += 1
  } else {
    result.suggestions.push('包含大写字母')
  }

  if (/[0-9]/.test(password)) {
    result.score += 1
  } else {
    result.suggestions.push('包含数字')
  }

  if (/[^a-zA-Z0-9]/.test(password)) {
    result.score += 1
  } else {
    result.suggestions.push('包含特殊字符')
  }

  if (result.score <= 2) {
    result.level = 'weak'
  } else if (result.score <= 4) {
    result.level = 'medium'
  } else {
    result.level = 'strong'
  }

  return result
}

/**
 * 验证JSON格式
 * @param {string} jsonString - JSON字符串
 * @returns {boolean} 是否为有效JSON
 */
export function isValidJSON(jsonString) {
  try {
    JSON.parse(jsonString)
    return true
  } catch {
    return false
  }
}

/**
 * 验证Base64格式
 * @param {string} base64String - Base64字符串
 * @returns {boolean} 是否为有效Base64
 */
export function isValidBase64(base64String) {
  const base64Regex = /^[A-Za-z0-9+/]*={0,2}$/
  return base64Regex.test(base64String) && base64String.length % 4 === 0
}

/**
 * 验证颜色值格式
 * @param {string} color - 颜色值
 * @returns {boolean} 是否为有效颜色值
 */
export function isValidColor(color) {
  const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/
  const rgbRegex = /^rgb\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)$/
  const rgbaRegex = /^rgba\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(0|1|0?\.\d+)\s*\)$/
  
  return hexRegex.test(color) || rgbRegex.test(color) || rgbaRegex.test(color)
}