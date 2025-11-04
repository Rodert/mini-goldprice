import { login, logout, getUserInfo, getUserMenus } from '@/api/auth'
import { getToken, setToken, removeToken } from '@/utils/auth'

const state = {
  token: getToken(),
  name: '',
  avatar: '',
  roles: [],
  menus: []
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  },
  SET_MENUS: (state, menus) => {
    state.menus = menus
  }
}

const actions = {
  // 用户登录
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    return new Promise((resolve, reject) => {
      login({ username: username.trim(), password: password }).then(response => {
        const { data } = response
        commit('SET_TOKEN', data.token)
        setToken(data.token)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // 获取用户信息
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getUserInfo(state.token).then(response => {
        const { data } = response

        if (!data) {
          reject(new Error('验证失败，请重新登录'))
        }

        const { real_name, avatar, roles, menus } = data

        commit('SET_NAME', real_name)
        commit('SET_AVATAR', avatar)
        commit('SET_ROLES', roles)
        commit('SET_MENUS', menus || [])
        resolve(data)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // 获取用户菜单
  getMenus({ commit }) {
    return new Promise((resolve, reject) => {
      getUserMenus().then(response => {
        const { data } = response
        commit('SET_MENUS', data || [])
        resolve(data)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // 用户登出
  logout({ commit, state }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        commit('SET_TOKEN', '')
        commit('SET_ROLES', [])
        removeToken()
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // 清除 Token
  resetToken({ commit }) {
    return new Promise(resolve => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      commit('SET_MENUS', [])
      removeToken()
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}



