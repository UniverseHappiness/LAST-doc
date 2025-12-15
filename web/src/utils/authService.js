// 认证服务
import axios from 'axios';

const API_BASE_URL = '/api/v1';

class AuthService {
  constructor() {
    this.token = localStorage.getItem('token');
    this.user = JSON.parse(localStorage.getItem('user') || 'null');
    
    // 设置axios默认headers
    if (this.token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`;
    }
  }

  // 用户注册
  async register(userData) {
    try {
      const response = await axios.post(`${API_BASE_URL}/auth/register`, userData);
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '注册失败' };
    }
  }

  // 用户登录
  async login(credentials) {
    try {
      const response = await axios.post(`${API_BASE_URL}/auth/login`, credentials);
      const { user, access_token, refresh_token } = response.data;
      
      console.log('DEBUG: authService.login - 获取到token', { access_token, user });
      
      // 保存token和用户信息
      this.token = access_token;
      this.user = user;
      
      localStorage.setItem('token', access_token);
      localStorage.setItem('refreshToken', refresh_token);
      localStorage.setItem('user', JSON.stringify(user));
      
      // 设置axios默认headers
      axios.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
      
      // 确认token已保存
      console.log('DEBUG: authService.login - token已保存', {
        thisToken: this.token,
        localStorageToken: localStorage.getItem('token')
      });
      
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '登录失败' };
    }
  }

  // 登出
  logout() {
    this.token = null;
    this.user = null;
    
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('user');
    
    delete axios.defaults.headers.common['Authorization'];
  }

  // 刷新token
  async refreshToken() {
    try {
      const refreshToken = localStorage.getItem('refreshToken');
      if (!refreshToken) {
        throw new Error('没有刷新令牌');
      }
      
      const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
        refresh_token: refreshToken
      });
      
      const { access_token } = response.data;
      this.token = access_token;
      
      localStorage.setItem('token', access_token);
      axios.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
      
      return access_token;
    } catch (error) {
      // 刷新失败，清除登录信息
      this.logout();
      throw error.response?.data || { error: '令牌刷新失败' };
    }
  }

  // 获取当前用户信息
  async getCurrentUser() {
    try {
      const response = await axios.get(`${API_BASE_URL}/users/profile`);
      this.user = response.data.user;
      localStorage.setItem('user', JSON.stringify(this.user));
      return this.user;
    } catch (error) {
      throw error.response?.data || { error: '获取用户信息失败' };
    }
  }

  // 更新用户信息
  async updateProfile(userData) {
    try {
      const response = await axios.put(`${API_BASE_URL}/users/profile`, userData);
      this.user = response.data.user;
      localStorage.setItem('user', JSON.stringify(this.user));
      return this.user;
    } catch (error) {
      throw error.response?.data || { error: '更新用户信息失败' };
    }
  }

  // 修改密码
  async changePassword(passwordData) {
    try {
      const response = await axios.post(`${API_BASE_URL}/users/change-password`, passwordData);
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '修改密码失败' };
    }
  }

  // 请求密码重置
  async requestPasswordReset(email) {
    try {
      const response = await axios.post(`${API_BASE_URL}/auth/password/reset-request`, { email });
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '请求密码重置失败' };
    }
  }

  // 重置密码
  async resetPassword(resetData) {
    try {
      const response = await axios.post(`${API_BASE_URL}/auth/password/reset`, resetData);
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '重置密码失败' };
    }
  }

  // 检查是否已登录
  isAuthenticated() {
    return !!this.token;
  }

  // 检查是否有特定角色
  hasRole(role) {
    return this.user && this.user.role === role;
  }

  // 检查是否是管理员
  isAdmin() {
    return this.hasRole('admin');
  }

  // 获取当前用户
  getUser() {
    return this.user;
  }

  // 获取当前用户ID
  getUserId() {
    return this.user ? this.user.id : null;
  }

  // 获取当前用户角色
  getUserRole() {
    return this.user ? this.user.role : null;
  }
}

export default new AuthService();