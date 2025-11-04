import request from '@/utils/request'

/**
 * 获取预约列表
 */
export function getAppointmentList(params) {
  return request({
    url: '/appointments',
    method: 'get',
    params
  })
}

/**
 * 获取预约详情
 */
export function getAppointment(id) {
  return request({
    url: `/appointments/${id}`,
    method: 'get'
  })
}

/**
 * 创建预约
 */
export function createAppointment(data) {
  return request({
    url: '/appointments',
    method: 'post',
    data
  })
}

/**
 * 更新预约
 */
export function updateAppointment(id, data) {
  return request({
    url: `/appointments/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除预约
 */
export function deleteAppointment(id) {
  return request({
    url: `/appointments/${id}`,
    method: 'delete'
  })
}

/**
 * 获取预约统计
 */
export function getAppointmentStats() {
  return request({
    url: '/appointments/stats',
    method: 'get'
  })
}

