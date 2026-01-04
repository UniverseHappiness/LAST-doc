<template>
  <div class="reset-password-container">
    <div class="reset-password-form">
      <h2>重置密码</h2>
      <p class="instruction-text">请输入您的新密码</p>
      <form @submit.prevent="handleResetPassword">
        <div class="form-group">
          <label for="token">重置令牌</label>
          <input
            type="text"
            id="token"
            v-model="form.token"
            required
            placeholder="请输入重置令牌"
          />
          <small class="form-text">重置令牌已从URL中自动填充</small>
        </div>
        <div class="form-group">
          <label for="password">新密码</label>
          <input
            type="password"
            id="password"
            v-model="form.password"
            required
            minlength="6"
            placeholder="请输入新密码（至少6位）"
          />
        </div>
        <div class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input
            type="password"
            id="confirmPassword"
            v-model="form.confirmPassword"
            required
            minlength="6"
            placeholder="请再次输入新密码"
          />
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="loading || !isFormValid">
            {{ loading ? '重置中...' : '重置密码' }}
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
      </div>
    </div>
  </div>
</template>

<script>
import authService from '@/utils/authService';

export default {
  name: 'ResetPasswordView',
  props: {
    token: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      form: {
        token: '',
        password: '',
        confirmPassword: ''
      },
      loading: false,
      error: '',
      success: ''
    };
  },
  computed: {
    isFormValid() {
      return this.form.token && 
             this.form.password && 
             this.form.password.length >= 6 &&
             this.form.password === this.form.confirmPassword;
    }
  },
  created() {
    // 从URL参数或props中获取token
    const urlToken = this.$route?.query?.token;
    if (urlToken) {
      this.form.token = urlToken;
    } else if (this.token) {
      this.form.token = this.token;
    }
  },
  methods: {
    async handleResetPassword() {
      try {
        this.loading = true;
        this.error = '';
        this.success = '';

        // 验证密码匹配
        if (this.form.password !== this.form.confirmPassword) {
          this.error = '两次输入的密码不一致';
          return;
        }

        // 验证密码长度
        if (this.form.password.length < 6) {
          this.error = '密码长度不能少于6位';
          return;
        }

        await authService.resetPassword({
          token: this.form.token,
          password: this.form.password
        });

        this.success = '密码重置成功，请使用新密码登录';

        // 3秒后跳转到登录页面
        setTimeout(() => {
          this.$emit('back-to-login');
        }, 3000);

      } catch (error) {
        this.error = error.error || '重置密码失败，请检查重置令牌是否有效';
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.reset-password-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.reset-password-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.reset-password-form h2 {
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

.form-text {
  display: block;
  margin-top: 0.5rem;
  font-size: 0.85rem;
  color: #6c757d;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
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