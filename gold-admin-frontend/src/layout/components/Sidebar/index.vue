<template>
  <div>
    <div class="logo-container">
      <h1 class="sidebar-title">沪金汇管理系统</h1>
    </div>
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
        :default-active="activeMenu"
        :background-color="variables.menuBg"
        :text-color="variables.menuText"
        :active-text-color="variables.menuActiveText"
        :unique-opened="false"
        :collapse-transition="false"
        mode="vertical"
      >
        <sidebar-item v-for="route in routes" :key="route.path" :item="route" :base-path="route.path" />
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import SidebarItem from './SidebarItem'
import variables from '@/styles/variables.scss'
import Layout from '@/layout'

export default {
  components: { SidebarItem },
  computed: {
    ...mapGetters(['sidebar', 'menus']),
    routes() {
      // 如果有后端返回的用户菜单，使用用户菜单；否则使用所有路由
      if (this.menus && this.menus.length > 0) {
        return this.convertMenusToRoutes(this.menus)
      }
      return this.$router.options.routes
    },
    activeMenu() {
      const route = this.$route
      const { meta, path } = route
      if (meta.activeMenu) {
        return meta.activeMenu
      }
      return path
    },
    variables() {
      return variables
    }
  },
  methods: {
    // 将后端菜单树转换为前端路由格式
    convertMenusToRoutes(menus) {
      const routes = []
      
      // 首页路由（始终显示）
      routes.push({
        path: '/',
        component: Layout,
        redirect: '/dashboard',
        children: [
          {
            path: 'dashboard',
            name: 'Dashboard',
            meta: { title: '首页', icon: 'el-icon-s-home' }
          }
        ]
      })

      // 转换后端菜单为路由
      menus.forEach(menu => {
        const route = {
          path: menu.path || '/',
          component: Layout,
          redirect: menu.redirect || '',
          name: menu.name || menu.title,
          meta: {
            title: menu.title,
            icon: menu.icon
          }
        }

        // 处理子菜单
        if (menu.children && menu.children.length > 0) {
          route.children = menu.children.map(child => ({
            path: child.path,
            name: child.name || child.title,
            meta: {
              title: child.title,
              icon: child.icon
            }
          }))
        }

        routes.push(route)
      })

      return routes
    }
  }
}
</script>

<style lang="scss" scoped>
.logo-container {
  height: 50px;
  line-height: 50px;
  background: #2b2f3a;
  text-align: center;
  overflow: hidden;

  .sidebar-title {
    font-size: 18px;
    font-weight: 600;
    color: #fff;
    margin: 0;
  }
}
</style>



