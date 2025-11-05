import router from './router'
import store from './store'
import { Message } from 'element-ui'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { getToken } from '@/utils/auth'
import Layout from '@/layout'

NProgress.configure({ showSpinner: false })

const whiteList = ['/login', '/display']

// 组件映射表 - 添加新功能时在这里添加组件映射
const componentMap = {
  'dashboard/index': () => import('@/views/dashboard/index'),
  'price/index': () => import('@/views/price/index'),
  'appointment/index': () => import('@/views/appointment/index'),
  'shop/index': () => import('@/views/shop/index'),
  'system/user/index': () => import('@/views/system/user/index'),
  'system/role/index': () => import('@/views/system/role/index'),
  'system/menu/index': () => import('@/views/system/menu/index')
}

// 将后端菜单转换为Vue路由
function generateRoutes(menus) {
  const routes = []
  
  menus.forEach(menu => {
    // 只处理目录和菜单类型
    if (menu.type === 1) {
      // 目录 - 包含子菜单
      const route = {
        path: menu.path,
        component: Layout,
        name: menu.name,
        meta: {
          title: menu.title,
          icon: menu.icon
        },
        children: []
      }
      
      if (menu.children && menu.children.length > 0) {
        menu.children.forEach(child => {
          if (child.type === 2) {
            // 子菜单
            const component = componentMap[child.component]
            if (component) {
              route.children.push({
                path: child.path,
                name: child.name,
                component: component,
                meta: {
                  title: child.title,
                  icon: child.icon
                }
              })
            }
          }
        })
      }
      
      routes.push(route)
    } else if (menu.type === 2 && menu.parent_id === 0) {
      // 独立菜单（没有父级的菜单）
      const component = componentMap[menu.component]
      if (component) {
        routes.push({
          path: menu.path,
          component: Layout,
          children: [{
            path: '',
            name: menu.name,
            component: component,
            meta: {
              title: menu.title,
              icon: menu.icon
            }
          }]
        })
      }
    }
  })
  
  return routes
}

router.beforeEach(async(to, from, next) => {
  NProgress.start()

  const hasToken = getToken()

  if (hasToken) {
    if (to.path === '/login') {
      next({ path: '/' })
      NProgress.done()
    } else {
      const hasGetUserInfo = store.getters.roles && store.getters.roles.length > 0
      
      if (hasGetUserInfo) {
        next()
      } else {
        try {
          const { menus } = await store.dispatch('user/getInfo')
          
          const accessRoutes = menus && menus.length > 0 ? generateRoutes(menus) : []
          
          console.log('动态添加路由:', accessRoutes)
          
          accessRoutes.forEach(route => {
            router.addRoute(route)
          })
          
          next({ ...to, replace: true })
        } catch (error) {
          console.error('获取用户信息失败:', error)
          await store.dispatch('user/resetToken')
          Message.error(error.message || '验证失败，请重新登录')
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    if (whiteList.indexOf(to.path) !== -1) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})
