<template>
  <div class="user-profile">
    <div class="profile-header">
      <h2>用户资料</h2>
    </div>
    
    <div class="profile-content">
      <div class="profile-section">
        <h3>基本信息</h3>
        <form @submit.prevent="updateProfile" class="profile-form">
          <div class="form-row">
            <div class="form-group">
              <label for="username">用户名</label>
              <input
                type="text"
                id="username"
                v-model="userForm.username"
                disabled
                class="disabled-input"
              />
              <small>用户名不可修改</small>
            </div>
            <div class="form-group">
              <label for="email">邮箱</label>
              <input
                type="email"
                id="email"
                v-model="userForm.email"
                required
              />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label for="firstName">姓名</label>
              <input
                type="text"
                id="firstName"
                v-model="userForm.first_name"
                required
              />
            </div>
            <div class="form-group">
              <label for="role">角色</label>
              <input
                type="text"
                id="role"
                v-model="userForm.role"
                disabled
                class="disabled-input"
              />
              <small>角色由管理员分配</small>
            </div>
          </div>
          
          <div class="form-actions">
            <button type="submit" class="btn btn-primary" :disabled="loading">
              {{ loading ? '更新中...' : '更新资料' }}
            </button>
          </div>
        </form>
      </div>
      
      <div class="profile-section">
        <h3>修改密码</h3>
        <form @submit.prevent="changePassword" class="password-form">
          <div class="form-group">
            <label for="oldPassword">当前密码</label>
            <input
              type="password"
              id="oldPassword"
              v-model="passwordForm.old_password"
              required
            />
          </div>
          
          <div class="form-group">
            <label for="newPassword">新密码</label>
            <input
              type="password"
              id="newPassword"
              v-model="passwordForm.new_password"
              required
              minlength="6"
            />
          </div>
          
          <div class="form-group">
            <label for="confirmPassword">确认新密码</label>
            <input
              type="password"
              id="confirmPassword"
              v-model="confirmPassword"
              required
              minlength="6"
            />
          </div>
          
          <div class="form-actions">
            <button type="submit" class="btn btn-secondary" :disabled="passwordLoading">
              {{ passwordLoading ? '修改中...' : '修改密码' }}
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <div v-if="error" class="error-message">
      {{ error }}
    </div>
    <div v-if="success" class="success-message">
      {{ success }}
    </div>
  </div>
</template>

<script>
import authService from '@/utils/authService';

export default {
  name: 'UserProfileView',
  data() {
    return {
      userForm: {
        username: '',
        email: '',
        first_name: '',
        role: ''
      },
      passwordForm: {
        old_password: '',
        new_password: ''
      },
      confirmPassword: '',
      loading: false,
      passwordLoading: false,
      error: '',
      success: ''
    };
  },
  methods: {
    async loadUserProfile() {
      try {
        const user = await authService.getCurrentUser();
        this.userForm = {
          username: user.username,
          email: user.email,
          first_name: user.first_name,
          role: user.role
        };
      } catch (error) {
        this.error = '加载用户资料失败';
      }
    },
    
    async updateProfile() {
      try {
        this.loading = true;
        this.error = '';
        this.success = '';
        
        const updatedUser = await authService.updateProfile(this.userForm);
        this.success = '资料更新成功';
        
        // 触发个人资料更新事件
        this.$emit('profile-updated', updatedUser);
      } catch (error) {
        this.error = error.error || '更新资料失败';
      } finally {
        this.loading = false;
      }
    },
    
    async changePassword() {
      try {
        // 验证密码确认
        if (this.passwordForm.new_password !== this.confirmPassword) {
          this.error = '两次输入的密码不一致';
          return;
        }
        
        this.passwordLoading = true;
        this.error = '';
        this.success = '';
        
        await authService.changePassword(this.passwordForm);
        this.success = '密码修改成功';
        
        // 清空密码表单
        this.passwordForm = {
          old_password: '',
          new_password: ''
        };
        this.confirmPassword = '';
      } catch (error) {
        this.error = error.error || '修改密码失败';
      } finally {
        this.passwordLoading = false;
      }
    }
  },
  created() {
    // 检查是否已登录
    if (!authService.isAuthenticated()) {
      this.$router.push('/login');
      return;
    }
    
    this.loadUserProfile();
  }
};
</script>

<style scoped>
.user-profile {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.profile-header {
  margin-bottom: 2rem;
}

.profile-header h2 {
  color: #333;
  margin-bottom: 0.5rem;
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.profile-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
}

.profile-section h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #333;
  border-bottom: 1px solid #eee;
  padding-bottom: 0.5rem;
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

.form-group input.disabled-input {
  background-color: #f8f9fa;
  color: #6c757d;
}

.form-group small {
  display: block;
  margin-top: 0.25rem;
  color: #6c757d;
  font-size: 0.875rem;
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
}

.btn-primary:hover:not(:disabled) {
  background-color: #0069d9;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #5a6268;
}

.btn:disabled {
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

@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
    gap: 0;
  }
  
  .user-profile {
    padding: 1rem;
  }
}
</style>