import axios from 'axios'
import { MessageBox, Message } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'

// 创建 axios 实例
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API || '/api',
  timeout: 5000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 在请求发送之前做一些处理
    if (store.getters.token) {
      // 让每个请求携带 token
      config.headers['Authorization'] = 'Bearer ' + getToken()
    }
    return config
  },
  error => {
    // 请求错误处理
    console.log(error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data

    // 如果返回的状态码不是 200，则判断为错误
    if (res.code !== 200) {
      Message({
        message: res.message || 'Error',
        type: 'error',
        duration: 5 * 1000
      })

      // 401: 未授权
      if (res.code === 401) {
        MessageBox.confirm('登录已过期，请重新登录', '提示', {
          confirmButtonText: '重新登录',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          store.dispatch('user/resetToken').then(() => {
            location.reload()
          })
        })
      }
      // 403: 禁止访问
      if (res.code === 403) {
        MessageBox.confirm(res.message || '无权限访问', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          // 可以跳转到首页或者不处理
        })
      }
      return Promise.reject(new Error(res.message || 'Error'))
    } else {
      return res
    }
  },
  error => {
    console.log('err' + error)
    // 处理 HTTP 状态码错误（如网络错误、超时等）
    let message = '请求失败'
    if (error.response) {
      // 服务器返回了错误响应
      const status = error.response.status
      if (status === 403) {
        message = '无权限访问'
      } else if (status === 401) {
        message = '未授权，请重新登录'
        MessageBox.confirm('登录已过期，请重新登录', '提示', {
          confirmButtonText: '重新登录',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          store.dispatch('user/resetToken').then(() => {
            location.reload()
          })
        })
      } else if (status === 404) {
        message = '请求的资源不存在'
      } else if (status >= 500) {
        message = '服务器错误'
      }
      // 尝试从响应数据中获取错误信息
      if (error.response.data && error.response.data.message) {
        message = error.response.data.message
      }
    } else if (error.request) {
      message = '网络错误，请检查网络连接'
    } else {
      message = error.message || '请求失败'
    }
    
    Message({
      message: message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service

