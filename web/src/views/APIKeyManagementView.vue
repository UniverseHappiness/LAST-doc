<template>
  <div class="api-key-management">
    <div class="management-header">
      <h2>API密钥管理</h2>
      <div class="header-actions">
        <button @click="handleCreateKeyClick" class="btn btn-primary">
          创建新密钥
        </button>
      </div>
    </div>
    
    <div class="api-keys-table">
      <table>
        <thead>
          <tr>
            <th>名称</th>
            <th>密钥</th>
            <th>用户</th>
            <th>过期时间</th>
            <th>最后使用</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="7" class="loading-row">加载中...</td>
          </tr>
          <tr v-else-if="apiKeys.length === 0">
            <td colspan="7" class="empty-row">暂无数据</td>
          </tr>
          <tr v-else v-for="apiKey in apiKeys" :key="apiKey.id">
            <td>{{ apiKey.name }}</td>
            <td class="key-cell">
              <code>{{ maskedKey(apiKey.key) }}</code>
              <button 
                @click="copyToClipboard(apiKey.key)" 
                class="btn-copy"
                title="复制完整密钥"
              >
                📋
              </button>
            </td>
            <td>{{ apiKey.user_id }}</td>
            <td>{{ formatDate(apiKey.expires_at) }}</td>
            <td>{{ formatDate(apiKey.last_used) }}</td>
            <td>{{ formatDate(apiKey.created_at) }}</td>
            <td class="actions-cell">
              <button
                @click="deleteApiKey(apiKey)"
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
    
    <!-- 创建API密钥模态框 -->
    <div v-if="showCreateModal" class="modal-overlay" @click="closeCreateModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>创建API密钥</h3>
          <button @click="closeCreateModal" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label for="keyName">密钥名称</label>
            <input
              type="text"
              id="keyName"
              v-model="createForm.name"
              placeholder="请输入密钥名称"
              required
            />
          </div>
          <div class="form-group">
            <label for="expiresAt">过期时间（可选）</label>
            <input
              type="datetime-local"
              id="expiresAt"
              v-model="createForm.expires_at"
            />
            <small>留空表示永不过期</small>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeCreateModal" class="btn btn-secondary">取消</button>
          <button @click="createApiKey" class="btn btn-primary" :disabled="creating">
            {{ creating ? '创建中...' : '创建' }}
          </button>
        </div>
      </div>
    </div>
    
    <!-- 创建成功模态框 -->
    <div v-if="showSuccessModal" class="modal-overlay" @click="closeSuccessModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>API密钥创建成功</h3>
          <button @click="closeSuccessModal" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>请复制并妥善保存您的API密钥</label>
            <div class="key-display">
              <code>{{ newApiKey.key }}</code>
              <button @click="copyToClipboard(newApiKey.key)" class="btn-copy">
                📋 复制
              </button>
            </div>
            <p class="warning">此密钥只显示一次，请立即保存！</p>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeSuccessModal" class="btn btn-primary">我已保存</button>
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
  name: 'APIKeyManagementView',
  data() {
    return {
      apiKeys: [],
      loading: false,
      error: '',
      success: '',
      
      // 创建模态框
      showCreateModal: false,
      createForm: {
        name: '',
        expires_at: ''
      },
      creating: false,
      
      // 成功模态框
      showSuccessModal: false,
      newApiKey: null
    };
  },
  methods: {
    async loadApiKeys() {
      try {
        this.loading = true;
        this.error = '';
        
        const params = {};
        // 如果不是管理员，只能查看自己的API密钥
        if (!authService.isAdmin()) {
          params.user_id = authService.getUserId();
        }
        
        const response = await axios.get('/api/v1/mcp/keys', { params });
        this.apiKeys = response.data.keys;
      } catch (error) {
        this.error = error.response?.data?.error || '加载API密钥列表失败';
      } finally {
        this.loading = false;
      }
    },
    
    async createApiKey() {
      try {
        console.log('[DEBUG] 开始创建API密钥', this.createForm);
        this.creating = true;
        this.error = '';
        
        // 验证必填字段
        if (!this.createForm.name || this.createForm.name.trim() === '') {
          this.error = '请输入API密钥名称';
          return;
        }
        
        // 准备请求数据，处理过期时间格式
        const requestData = {
          name: this.createForm.name.trim()
        };
        
        // 只有当过期时间不为空时才添加到请求中
        if (this.createForm.expires_at && this.createForm.expires_at.trim() !== '') {
          // 将本地时间字符串转换为RFC3339格式
          const date = new Date(this.createForm.expires_at);
          if (isNaN(date.getTime())) {
            this.error = '无效的过期时间格式';
            return;
          }
          requestData.expires_at = date.toISOString();
          console.log('[DEBUG] 转换后的过期时间', requestData.expires_at);
        }
        
        console.log('[DEBUG] 发送的请求数据', requestData);
        const response = await axios.post('/api/v1/mcp/keys', requestData);
        console.log('[DEBUG] API密钥创建成功', response.data);
        this.newApiKey = response.data;
        
        console.log('[DEBUG] 关闭创建模态框，显示成功模态框');
        this.showCreateModal = false;
        this.showSuccessModal = true;
        
        // 重置表单
        this.createForm = {
          name: '',
          expires_at: ''
        };
        
        // 刷新列表
        this.loadApiKeys();
      } catch (error) {
        console.error('[DEBUG] 创建API密钥失败', error);
        this.error = error.response?.data?.error || '创建API密钥失败';
      } finally {
        this.creating = false;
      }
    },
    
    async deleteApiKey(apiKey) {
      if (confirm(`确定要删除API密钥 "${apiKey.name}" 吗？`)) {
        try {
          console.log('[DEBUG] 尝试删除API密钥', { id: apiKey.id, name: apiKey.name });
          await axios.delete(`/api/v1/mcp/keys/${apiKey.id}`);
          this.success = 'API密钥删除成功';
          this.loadApiKeys();
        } catch (error) {
          console.error('[DEBUG] 删除API密钥失败', error);
          this.error = error.response?.data?.error || '删除API密钥失败';
        }
      }
    },
    
    closeCreateModal() {
      console.log('[DEBUG] 关闭创建模态框');
      this.showCreateModal = false;
      this.createForm = {
        name: '',
        expires_at: ''
      };
    },
    
    closeSuccessModal() {
      this.showSuccessModal = false;
      this.newApiKey = null;
    },
    
    maskedKey(key) {
      if (!key) return '';
      return key.substring(0, 8) + '...' + key.substring(key.length - 4);
    },
    
    formatDate(dateString) {
      if (!dateString) return '-';
      return new Date(dateString).toLocaleString();
    },
    
    async copyToClipboard(text) {
      try {
        await navigator.clipboard.writeText(text);
        this.success = '已复制到剪贴板';
        setTimeout(() => {
          this.success = '';
        }, 3000);
      } catch (error) {
        this.error = '复制失败，请手动复制';
      }
    },
    
    handleCreateKeyClick() {
      console.log('[DEBUG] 点击创建密钥按钮');
      try {
        this.showCreateModal = true;
        console.log('[DEBUG] 创建模态框显示状态:', this.showCreateModal);
      } catch (error) {
        console.error('[DEBUG] 显示创建模态框时出错', error);
      }
    }
  },
  created() {
    // 检查是否已登录
    if (!authService.isAuthenticated()) {
      this.$router.push('/login');
      return;
    }
    
    this.loadApiKeys();
  }
};
</script>

<style scoped>
.api-key-management {
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
}

.api-keys-table {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.api-keys-table table {
  width: 100%;
  border-collapse: collapse;
}

.api-keys-table th,
.api-keys-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.api-keys-table th {
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

.key-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.key-cell code {
  background-color: #f8f9fa;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.875rem;
  flex: 1;
  word-break: break-all;
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

.btn-delete:hover {
  background-color: #f8d7da;
}

.btn-copy {
  background: none;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-copy:hover {
  background-color: #f0f0f0;
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

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 确保模态框在所有情况下都能正确显示 */
.modal-overlay {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  background-color: rgba(0, 0, 0, 0.5) !important;
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
  z-index: 9999 !important; /* 提高z-index确保在最上层 */
  padding: 20px;
  box-sizing: border-box;
}

.modal {
  background: white !important;
  border-radius: 8px !important;
  width: 90% !important;
  max-width: 500px !important;
  max-height: 90vh !important;
  overflow-y: auto !important;
  position: relative !important;
  z-index: 10000 !important; /* 确保模态框内容在遮罩层之上 */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15) !important;
  transform: translateZ(0) !important; /* 确保在所有浏览器中都能正确显示 */
  display: flex !important;
  flex-direction: column !important;
}

/* 确保模态框内容正确显示 */
.modal-header, .modal-body, .modal-footer {
  position: relative !important;
  z-index: inherit !important;
}

/* 调试样式 - 确保模态框正确显示 */
.modal-overlay::before {
  content: "模态框已显示 (z-index: 9999)" !important;
  position: absolute !important;
  top: 10px !important;
  left: 10px !important;
  background: rgba(255, 255, 255, 0.8) !important;
  padding: 5px 10px !important;
  border-radius: 4px !important;
  font-size: 12px !important;
  color: #333 !important;
  z-index: 10001 !important;
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

.form-group input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}

.form-group small {
  display: block;
  margin-top: 0.25rem;
  color: #6c757d;
  font-size: 0.875rem;
}

.key-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.key-display code {
  flex: 1;
  background-color: #f8f9fa;
  padding: 0.5rem;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.875rem;
  word-break: break-all;
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
  margin-top: 1rem;
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
  
  .api-keys-table {
    overflow-x: auto;
  }
  
  .api-keys-table table {
    min-width: 800px;
  }
}
</style>