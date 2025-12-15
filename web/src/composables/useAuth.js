import { ref, computed } from 'vue'
import authService from '../utils/authService'

// 创建响应式引用
const token = ref(localStorage.getItem('token'))
const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

// 监听localStorage变化
window.addEventListener('storage', (e) => {
  if (e.key === 'token') {
    token.value = e.newValue
  } else if (e.key === 'user') {
    user.value = e.newValue ? JSON.parse(e.newValue) : null
  }
})

export function useAuth() {
  const isAuthenticated = computed(() => {
    console.log('DEBUG: isAuthenticated计算属性被调用', { token: token.value, result: !!token.value })
    return !!token.value
  })
  const currentUser = computed(() => user.value)
  
  const login = async (credentials) => {
    try {
      const response = await authService.login(credentials)
      // 更新本地状态
      token.value = authService.token
      user.value = authService.user
      console.log('DEBUG: 登录成功，更新认证状态', {
        authServiceToken: authService.token,
        tokenValue: token.value,
        userValue: user.value,
        isAuthenticated: !!token.value
      })
      return response
    } catch (error) {
      throw error
    }
  }
  
  const logout = () => {
    authService.logout()
    // 更新本地状态
    token.value = null
    user.value = null
  }
  
  const register = async (userData) => {
    try {
      const response = await authService.register(userData)
      return response
    } catch (error) {
      throw error
    }
  }
  
  const checkAuthStatus = async () => {
    try {
      // 先检查localStorage中的token
      const localToken = localStorage.getItem('token')
      console.log('DEBUG: checkAuthStatus - localStorage中的token', localToken)
      
      if (!localToken) {
        console.log('DEBUG: 没有找到token，未认证')
        return false
      }
      
      // 更新本地的token值
      token.value = localToken
      console.log('DEBUG: 更新本地token值', { token: token.value })
      
      // 尝试获取当前用户信息
      const currentUser = await authService.getCurrentUser()
      // 更新本地状态
      user.value = currentUser
      console.log('DEBUG: 检查认证状态成功', {
        token: token.value,
        user: user.value,
        isAuthenticated: !!token.value
      })
      return true
    } catch (error) {
      // 如果获取用户信息失败，清除认证状态
      console.error('DEBUG: 检查认证状态失败', error)
      logout()
      return false
    }
  }
  
  const updateProfile = async (userData) => {
    try {
      const response = await authService.updateProfile(userData)
      // 更新本地状态
      user.value = authService.user
      return response
    } catch (error) {
      throw error
    }
  }
  
  const changePassword = async (passwordData) => {
    try {
      const response = await authService.changePassword(passwordData)
      return response
    } catch (error) {
      throw error
    }
  }
  
  const requestPasswordReset = async (email) => {
    try {
      const response = await authService.requestPasswordReset(email)
      return response
    } catch (error) {
      throw error
    }
  }
  
  const resetPassword = async (resetData) => {
    try {
      const response = await authService.resetPassword(resetData)
      return response
    } catch (error) {
      throw error
    }
  }
  
  const hasRole = (role) => {
    return authService.hasRole(role)
  }
  
  const isAdmin = () => {
    return authService.isAdmin()
  }
  
  return {
    isAuthenticated,
    currentUser,
    login,
    logout,
    register,
    checkAuthStatus,
    updateProfile,
    changePassword,
    requestPasswordReset,
    resetPassword,
    hasRole,
    isAdmin
  }
}