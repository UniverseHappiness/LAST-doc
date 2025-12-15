<template>
  <div class="user-management">
    <div class="management-header">
      <h2>用户管理</h2>
      <div class="header-actions">
        <div class="search-box">
          <input
            type="text"
            v-model="searchQuery"
            placeholder="搜索用户..."
            @input="handleSearch"
          />
        </div>
        <div class="filter-box">
          <select v-model="roleFilter" @change="loadUsers">
            <option value="">所有角色</option>
            <option value="admin">管理员</option>
            <option value="editor">编辑员</option>
            <option value="user">普通用户</option>
          </select>
          <select v-model="statusFilter" @change="loadUsers">
            <option value="">所有状态</option>
            <option value="true">激活</option>
            <option value="false">禁用</option>
          </select>
        </div>
      </div>
    </div>
    
    <div class="users-table">
      <table>
        <thead>
          <tr>
            <th>用户名</th>
            <th>邮箱</th>
            <th>姓名</th>
            <th>角色</th>
            <th>状态</th>
            <th>最后登录</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="8" class="loading-row">加载中...</td>
          </tr>
          <tr v-else-if="users.length === 0">
            <td colspan="8" class="empty-row">暂无数据</td>
          </tr>
          <tr v-else v-for="user in users" :key="user.id">
            <td>{{ user.username }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.first_name }} {{ user.last_name }}</td>
            <td>
              <span :class="['role-badge', `role-${user.role}`]">
                {{ getRoleLabel(user.role) }}
              </span>
            </td>
            <td>
              <span :class="['status-badge', user.is_active ? 'status-active' : 'status-inactive']">
                {{ user.is_active ? '激活' : '禁用' }}
              </span>
            </td>
            <td>{{ formatDate(user.last_login) }}</td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td class="actions-cell">
              <button
                @click="editUser(user)"
                class="btn-icon"
                title="编辑"
              >
                ✏️
              </button>
              <button
                @click="toggleUserStatus(user)"
                :class="['btn-icon', user.is_active ? 'btn-disable' : 'btn-enable']"
                :title="user.is_active ? '禁用' : '启用'"
              >
                {{ user.is_active ? '🚫' : '✅' }}
              </button>
              <button
                @click="deleteUser(user)"
                class="btn-icon btn-delete"
                title="删除"
              >
                🗑️
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <div class="pagination">
      <button
        @click="goToPage(currentPage - 1)"
        :disabled="currentPage <= 1"
        class="btn btn-secondary"
      >
        上一页
      </button>
      <span class="page-info">
        第 {{ currentPage }} 页，共 {{ totalPages }} 页
      </span>
      <button
        @click="goToPage(currentPage + 1)"
        :disabled="currentPage >= totalPages"
        class="btn btn-secondary"
      >
        下一页
      </button>
    </div>
    
    <!-- 编辑用户模态框 -->
    <div v-if="showEditModal" class="modal-overlay" @click="closeEditModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>编辑用户</h3>
          <button @click="closeEditModal" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label for="editEmail">邮箱</label>
            <input
              type="email"
              id="editEmail"
              v-model="editForm.email"
              required
            />
          </div>
          <div class="form-group">
            <label for="editFirstName">姓名</label>
            <input
              type="text"
              id="editFirstName"
              v-model="editForm.first_name"
              required
            />
          </div>
          <div class="form-group">
            <label for="editRole">角色</label>
            <select id="editRole" v-model="editForm.role">
              <option value="user">普通用户</option>
              <option value="editor">编辑员</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          <div class="form-group">
            <label>
              <input
                type="checkbox"
                v-model="editForm.is_active"
              />
              账户激活
            </label>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeEditModal" class="btn btn-secondary">取消</button>
          <button @click="saveUser" class="btn btn-primary" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
    
    <!-- 删除确认模态框 -->
    <div v-if="showDeleteModal" class="modal-overlay" @click="closeDeleteModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>确认删除</h3>
          <button @click="closeDeleteModal" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <p>确定要删除用户 <strong>{{ deleteTarget.username }}</strong> 吗？</p>
          <p class="warning">此操作不可撤销！</p>
        </div>
        <div class="modal-footer">
          <button @click="closeDeleteModal" class="btn btn-secondary">取消</button>
          <button @click="confirmDelete" class="btn btn-danger" :disabled="deleting">
            {{ deleting ? '删除中...' : '确认删除' }}
          </button>
        </div>
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
import axios from 'axios';
import authService from '@/utils/authService';

export default {
  name: 'UserManagementView',
  data() {
    return {
      users: [],
      loading: false,
      error: '',
      success: '',
      currentPage: 1,
      totalPages: 1,
      pageSize: 10,
      totalUsers: 0,
      searchQuery: '',
      roleFilter: '',
      statusFilter: '',
      searchTimeout: null,
      
      // 编辑模态框
      showEditModal: false,
      editForm: {
        id: '',
        email: '',
        first_name: '',
        role: 'user',
        is_active: true
      },
      saving: false,
      
      // 删除模态框
      showDeleteModal: false,
      deleteTarget: null,
      deleting: false
    };
  },
  methods: {
    async loadUsers() {
      try {
        this.loading = true;
        this.error = '';
        
        const params = {
          page: this.currentPage,
          limit: this.pageSize
        };
        
        if (this.roleFilter) {
          params.role = this.roleFilter;
        }
        
        if (this.statusFilter !== '') {
          params.is_active = this.statusFilter === 'true';
        }
        
        if (this.searchQuery) {
          params.search = this.searchQuery;
        }
        
        const response = await axios.get('/api/v1/users', { params });
        
        this.users = response.data.users;
        this.totalUsers = response.data.total;
        this.totalPages = Math.ceil(this.totalUsers / this.pageSize);
      } catch (error) {
        this.error = error.response?.data?.error || '加载用户列表失败';
      } finally {
        this.loading = false;
      }
    },
    
    handleSearch() {
      // 防抖搜索
      clearTimeout(this.searchTimeout);
      this.searchTimeout = setTimeout(() => {
        this.currentPage = 1;
        this.loadUsers();
      }, 500);
    },
    
    goToPage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page;
        this.loadUsers();
      }
    },
    
    getRoleLabel(role) {
      const roleMap = {
        'admin': '管理员',
        'editor': '编辑员',
        'user': '普通用户'
      };
      return roleMap[role] || role;
    },
    
    formatDate(dateString) {
      if (!dateString) return '-';
      return new Date(dateString).toLocaleString();
    },
    
    editUser(user) {
      this.editForm = {
        id: user.id,
        email: user.email,
        first_name: user.first_name,
        role: user.role,
        is_active: user.is_active
      };
      this.showEditModal = true;
    },
    
    closeEditModal() {
      this.showEditModal = false;
      this.editForm = {
        id: '',
        email: '',
        first_name: '',
        role: 'user',
        is_active: true
      };
    },
    
    async saveUser() {
      try {
        this.saving = true;
        this.error = '';
        this.success = '';
        
        await axios.put(`/api/v1/users/${this.editForm.id}`, this.editForm);
        
        this.success = '用户信息更新成功';
        this.closeEditModal();
        this.loadUsers();
      } catch (error) {
        this.error = error.response?.data?.error || '更新用户信息失败';
      } finally {
        this.saving = false;
      }
    },
    
    toggleUserStatus(user) {
      const action = user.is_active ? '禁用' : '启用';
      if (confirm(`确定要${action}用户 ${user.username} 吗？`)) {
        this.updateUserStatus(user.id, !user.is_active);
      }
    },
    
    async updateUserStatus(userId, isActive) {
      try {
        await axios.put(`/api/v1/users/${userId}`, {
          is_active: isActive
        });
        
        this.success = `用户状态更新成功`;
        this.loadUsers();
      } catch (error) {
        this.error = error.response?.data?.error || '更新用户状态失败';
      }
    },
    
    deleteUser(user) {
      this.deleteTarget = user;
      this.showDeleteModal = true;
    },
    
    closeDeleteModal() {
      this.showDeleteModal = false;
      this.deleteTarget = null;
    },
    
    async confirmDelete() {
      try {
        this.deleting = true;
        this.error = '';
        this.success = '';
        
        await axios.delete(`/api/v1/users/${this.deleteTarget.id}`);
        
        this.success = '用户删除成功';
        this.closeDeleteModal();
        this.loadUsers();
      } catch (error) {
        this.error = error.response?.data?.error || '删除用户失败';
      } finally {
        this.deleting = false;
      }
    }
  },
  created() {
    // 检查是否是管理员
    if (!authService.isAdmin()) {
      this.$router.push('/');
      return;
    }
    
    this.loadUsers();
  }
};
</script>

<style scoped>
.user-management {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.management-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.management-header h2 {
  margin: 0;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
}

.search-box input {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 200px;
}

.filter-box select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
}

.users-table {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.users-table table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th,
.users-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.users-table th {
  background-color: #f8f9fa;
  font-weight: 600;
  color: #495057;
}

.loading-row,
.empty-row {
  text-align: center;
  color: #6c757d;
  font-style: italic;
}

.role-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
}

.role-admin {
  background-color: #dc3545;
  color: white;
}

.role-editor {
  background-color: #fd7e14;
  color: white;
}

.role-user {
  background-color: #6c757d;
  color: white;
}

.status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
}

.status-active {
  background-color: #28a745;
  color: white;
}

.status-inactive {
  background-color: #6c757d;
  color: white;
}

.actions-cell {
  display: flex;
  gap: 0.5rem;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
  font-size: 1rem;
}

.btn-icon:hover {
  background-color: #f0f0f0;
}

.btn-disable:hover {
  background-color: #fff3cd;
}

.btn-enable:hover {
  background-color: #d4edda;
}

.btn-delete:hover {
  background-color: #f8d7da;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
}

.page-info {
  color: #6c757d;
}

.btn {
  display: inline-block;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
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

.btn-danger {
  background-color: #dc3545;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background-color: #c82333;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6c757d;
}

.modal-body {
  padding: 1rem;
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

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 1rem;
  border-top: 1px solid #eee;
}

.warning {
  color: #dc3545;
  font-weight: 500;
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
  .management-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-actions {
    width: 100%;
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-box input {
    width: 100%;
  }
  
  .users-table {
    overflow-x: auto;
  }
  
  .users-table table {
    min-width: 800px;
  }
}
</style>