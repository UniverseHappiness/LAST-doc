<template>
  <div class="forgot-password-container">
    <div class="forgot-password-form">
      <h2>忘记密码</h2>
      <p class="instruction-text">请输入您的注册邮箱，我们将向您发送密码重置链接</p>
      <form @submit.prevent="handleRequestReset">
        <div class="form-group">
          <label for="email">邮箱地址</label>
          <input
            type="email"
            id="email"
            v-model="form.email"
            required
            placeholder="请输入您的注册邮箱"
          />
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '发送中...' : '发送重置链接' }}
          </button>
          <button type="button" class="btn btn-secondary" @click="$emit('back-to-login')" :disabled="loading">
            返回登录
          </button>
        </div>
      </form>
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
      <div v-if="success" class="success-message">
        {{ success }}
        <div v-if="showTokenInput" class="token-input-section">
          <h4>收到令牌？</h4>
          <div class="form-group">
            <label for="token">重置令牌</label>
            <input
              type="text"
              id="token"
              v-model="tokenInput"
              placeholder="请输入管理员提供的重置令牌"
            />
          </div>
          <button type="button" class="btn btn-primary" @click="goToResetPassword" :disabled="!tokenInput">
            使用令牌重置密码
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import authService from '@/utils/authService';

export default {
  name: 'ForgotPasswordView',
  data() {
    return {
      form: {
        email: ''
      },
      tokenInput: '',
      loading: false,
      error: '',
      success: '',
      showTokenInput: false
    };
  },
  methods: {
    async handleRequestReset() {
      try {
        this.loading = true;
        this.error = '';
        this.success = '';

        await authService.requestPasswordReset(this.form.email);
        
        this.success = '密码重置链接已发送到您的邮箱，请查收';
        this.showTokenInput = true;
        this.form.email = '';
      } catch (error) {
        this.error = error.error || '发送密码重置链接失败，请检查邮箱地址';
      } finally {
        this.loading = false;
      }
    },
    goToResetPassword() {
      if (this.tokenInput) {
        // 通过URL参数传递token
        this.$emit('use-token', this.tokenInput);
      }
    }
  }
};
</script>

<style scoped>
.forgot-password-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.forgot-password-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.forgot-password-form h2 {
  text-align: center;
  margin-bottom: 0.5rem;
  color: #333;
}

.instruction-text {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #666;
  font-size: 0.9rem;
}

.form-group {
  margin-bottom: 1.5rem;
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
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
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

.btn-secondary {
  background-color: #6c757d;
  color: white;
  width: 100%;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #5a6268;
}

.btn-secondary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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