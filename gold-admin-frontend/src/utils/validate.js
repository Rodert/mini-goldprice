/**
 * 验证用户名
 * @param {string} str
 * @returns {Boolean}
 */
export function validUsername(str) {
  return str.length >= 3
}

/**
 * 验证密码
 * @param {string} str
 * @returns {Boolean}
 */
export function validPassword(str) {
  return str.length >= 6
}

/**
 * 验证手机号
 * @param {string} str
 * @returns {Boolean}
 */
export function validPhone(str) {
  const reg = /^1[3-9]\d{9}$/
  return reg.test(str)
}

/**
 * 验证邮箱
 * @param {string} str
 * @returns {Boolean}
 */
export function validEmail(str) {
  const reg = /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/
  return reg.test(str)
}

/**
 * 验证是否为外部链接
 * @param {string} path
 * @returns {Boolean}
 */
export function isExternal(path) {
  return /^(https?:|mailto:|tel:)/.test(path)
}

