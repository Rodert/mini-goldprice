import Vue from 'vue'
import VueRouter from 'vue-router'
import Layout from '@/layout'

Vue.use(VueRouter)

/**
 * 公共路由
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/display',
    component: () => import('@/views/display/index'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index'),
        meta: { title: '首页', icon: 'el-icon-s-home' }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    redirect: '/system/user',
    name: 'System',
    meta: { title: '系统管理', icon: 'el-icon-setting' },
    children: [
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/system/user/index'),
        meta: { title: '用户管理', icon: 'el-icon-user' }
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/system/role/index'),
        meta: { title: '角色管理', icon: 'el-icon-s-custom' }
      },
      {
        path: 'menu',
        name: 'Menu',
        component: () => import('@/views/system/menu/index'),
        meta: { title: '菜单管理', icon: 'el-icon-menu' }
      }
    ]
  },
  {
    path: '/business',
    component: Layout,
    redirect: '/business/price',
    name: 'Business',
    meta: { title: '业务管理', icon: 'el-icon-s-management' },
    children: [
      {
        path: 'price',
        name: 'Price',
        component: () => import('@/views/price/index'),
        meta: { title: '价格管理', icon: 'el-icon-coin' }
      },
      {
        path: 'appointment',
        name: 'Appointment',
        component: () => import('@/views/appointment/index'),
        meta: { title: '预约管理', icon: 'el-icon-date' }
      }
    ]
  },
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new VueRouter({
  mode: 'history',
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher
}

// 解决动态路由添加时的导航重定向警告
const originalPush = VueRouter.prototype.push
VueRouter.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => {
    if (err.name !== 'NavigationDuplicated' && !err.message.includes('Redirected')) {
      return Promise.reject(err)
    }
  })
}

const originalReplace = VueRouter.prototype.replace
VueRouter.prototype.replace = function replace(location) {
  return originalReplace.call(this, location).catch(err => {
    if (err.name !== 'NavigationDuplicated' && !err.message.includes('Redirected')) {
      return Promise.reject(err)
    }
  })
}

export default router



