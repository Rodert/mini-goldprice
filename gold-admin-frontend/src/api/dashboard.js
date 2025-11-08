import request from '@/utils/request'

/**
 * 获取首页统计数据
 */
export function getDashboardStats() {
  return request({
    url: '/dashboard/stats',
    method: 'get'
  })
}

/**
 * 获取最近动态
 */
export function getRecentActivities() {
  return request({
    url: '/dashboard/activities',
    method: 'get'
  })
}

/**
 * 获取预约趋势（最近7天）
 */
export function getAppointmentTrend() {
  return request({
    url: '/dashboard/trend',
    method: 'get'
  })
}

