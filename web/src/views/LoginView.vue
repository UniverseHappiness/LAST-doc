<template>
  <div class="login-container">
    <div class="login-form">
      <h2>用户登录</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">用户名/邮箱</label>
          <input
            type="text"
            id="username"
            v-model="loginForm.username"
            required
            placeholder="请输入用户名或邮箱"
          />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input
            type="password"
            id="password"
            v-model="loginForm.password"
            required
            placeholder="请输入密码"
          />
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </div>
      </form>
      <div class="form-links">
        <a href="#" @click.prevent="$emit('show-register')">还没有账号？立即注册</a>
        <a href="#" @click.prevent="showForgotPassword">忘记密码？</a>
      </div>
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script>
import authService from '@/utils/authService';

export default {
  name: 'LoginView',
  data() {
    return {
      loginForm: {
        username: '',
        password: ''
      },
      loading: false,
      error: ''
    };
  },
  methods: {
    async handleLogin() {
      try {
        this.loading = true;
        this.error = '';
        
        const response = await authService.login(this.loginForm);
        
        // 显示成功消息
        this.$emit('login-success', response.user);
      } catch (error) {
        this.error = error.error || '登录失败，请检查用户名和密码';
      } finally {
        this.loading = false;
      }
    },
    showForgotPassword() {
      // 显示忘记密码提示
      alert('忘记密码功能暂未实现，请联系管理员重置密码');
    }
  },
  created() {
    // 如果已经登录，触发登录成功事件
    if (authService.isAuthenticated()) {
      this.$emit('login-success', authService.getUser());
    }
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.login-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.login-form h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #555;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-group input:focus {
  border-color: #007bff;
  outline: none;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.form-actions {
  margin-top: 1.5rem;
}

.btn {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 500;
  text-align: center;
  text-decoration: none;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background-color: #007bff;
  color: white;
  width: 100%;
}

.btn-primary:hover:not(:disabled) {
  background-color: #0069d9;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-links {
  margin-top: 1.5rem;
  text-align: center;
}

.form-links a {
  display: block;
  margin: 0.5rem 0;
  color: #007bff;
  text-decoration: none;
}

.form-links a:hover {
  text-decoration: underline;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
}
</style>