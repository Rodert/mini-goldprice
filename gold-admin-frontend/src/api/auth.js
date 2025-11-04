import request from '@/utils/request'

/**
 * 登录
 */
export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

/**
 * 获取用户信息
 */
export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

/**
 * 登出
 */
export function logout() {
  return request({
    url: '/logout',
    method: 'post'
  })
}

/**
 * 获取用户菜单
 */
export function getUserMenus() {
  return request({
    url: '/user/menus',
    method: 'get'
  })
}



