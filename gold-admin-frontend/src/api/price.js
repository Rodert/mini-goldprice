import request from '@/utils/request'

/**
 * 获取价格列表
 */
export function getPriceList(params) {
  return request({
    url: '/prices',
    method: 'get',
    params
  })
}

/**
 * 获取价格详情
 */
export function getPrice(id) {
  return request({
    url: `/prices/${id}`,
    method: 'get'
  })
}

/**
 * 创建价格
 */
export function createPrice(data) {
  return request({
    url: '/prices',
    method: 'post',
    data
  })
}

/**
 * 更新价格
 */
export function updatePrice(id, data) {
  return request({
    url: `/prices/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除价格
 */
export function deletePrice(id) {
  return request({
    url: `/prices/${id}`,
    method: 'delete'
  })
}

/**
 * 同步基础价格
 */
export function syncBasePrice() {
  return request({
    url: '/prices/sync',
    method: 'post'
  })
}















