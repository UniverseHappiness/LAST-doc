<template>
  <div class="register-container">
    <div class="register-form">
      <h2>用户注册</h2>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            type="text"
            id="username"
            v-model="registerForm.username"
            required
            placeholder="请输入用户名"
            minlength="3"
            maxlength="50"
          />
        </div>
        <div class="form-group">
          <label for="email">邮箱</label>
          <input
            type="email"
            id="email"
            v-model="registerForm.email"
            required
            placeholder="请输入邮箱地址"
          />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label for="firstName">姓名</label>
            <input
              type="text"
              id="firstName"
              v-model="registerForm.first_name"
              required
              placeholder="请输入姓名"
            />
          </div>
          <div class="form-group">
            <label for="lastName">姓氏</label>
            <input
              type="text"
              id="lastName"
              v-model="registerForm.last_name"
              required
              placeholder="请输入姓氏"
            />
          </div>
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input
            type="password"
            id="password"
            v-model="registerForm.password"
            required
            placeholder="请输入密码"
            minlength="6"
          />
        </div>
        <div class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input
            type="password"
            id="confirmPassword"
            v-model="confirmPassword"
            required
            placeholder="请再次输入密码"
            minlength="6"
          />
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '注册中...' : '注册' }}
          </button>
        </div>
      </form>
      <div class="form-links">
        <a href="#" @click.prevent="$emit('back-to-login')">已有账号？立即登录</a>
      </div>
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
      <div v-if="success" class="success-message">
        {{ success }}
      </div>
    </div>
  </div>
</template>

<script>
import authService from '@/utils/authService';

export default {
  name: 'RegisterView',
  data() {
    return {
      registerForm: {
        username: '',
        email: '',
        first_name: '',
        last_name: '',
        password: ''
      },
      confirmPassword: '',
      loading: false,
      error: '',
      success: ''
    };
  },
  methods: {
    async handleRegister() {
      try {
        // 验证密码确认
        if (this.registerForm.password !== this.confirmPassword) {
          this.error = '两次输入的密码不一致';
          return;
        }

        this.loading = true;
        this.error = '';
        this.success = '';
        
        const response = await authService.register(this.registerForm);
        
        this.success = '注册成功！请登录您的账号';
        
        // 触发注册成功事件
        this.$emit('register-success', response.user);
        
        // 3秒后触发返回登录事件
        setTimeout(() => {
          this.$emit('back-to-login');
        }, 3000);
      } catch (error) {
        this.error = error.error || '注册失败，请稍后重试';
      } finally {
        this.loading = false;
      }
    }
  },
  created() {
    // 如果已经登录，直接跳转到首页
    if (authService.isAuthenticated()) {
      this.$router.push('/');
    }
  }
};
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 1rem;
}

.register-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.register-form h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-row {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group {
  flex: 1;
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

.success-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
  border-radius: 4px;
}
</style>