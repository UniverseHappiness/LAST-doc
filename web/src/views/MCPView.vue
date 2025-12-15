<template>
  <div class="mcp-view">
    <div class="page-header">
      <h1>MCP协议管理</h1>
      <p>管理MCP协议配置和API密钥</p>
    </div>

    <div class="content-container">
      <!-- MCP配置卡片 -->
      <div class="card">
        <div class="card-header">
          <h2>MCP配置</h2>
        </div>
        <div class="card-body">
          <div class="config-info">
            <h3>CoStrict插件配置</h3>
            <p>将以下配置添加到CoStrict的MCP服务器配置中：</p>
            
            <div class="config-box">
              <pre><code>{{ mcpConfig }}</code></pre>
            </div>
            
            <div class="config-actions">
              <button @click="copyConfig" class="btn btn-primary">复制配置</button>
              <button @click="testConnection" class="btn btn-secondary" :disabled="testing">
                {{ testing ? '测试中...' : '测试连接' }}
              </button>
            </div>
          </div>
        </div>
      </div>


      <!-- 使用说明卡片 -->
      <div class="card">
        <div class="card-header">
          <h2>使用说明</h2>
        </div>
        <div class="card-body">
          <div class="instructions">
            <h3>如何在CoStrict中使用MCP协议</h3>
            <ol>
              <li>打开CoStrict设置，找到MCP服务器配置</li>
              <li>点击"添加MCP服务器"</li>
              <li>将上方提供的配置粘贴到配置框中</li>
              <li>保存配置并重启CoStrict</li>
              <li>在对话中即可通过MCP协议访问技术文档库</li>
            </ol>
            
            <h3>可用的MCP工具</h3>
            <ul>
              <li><strong>search_documents</strong>: 搜索技术文档，支持关键词和语义搜索</li>
              <li><strong>get_document_content</strong>: 获取指定文档的详细内容</li>
            </ul>
          </div>
        </div>
      </div>
    </div>


    <!-- 消息提示 -->
    <div v-if="message" class="message" :class="messageType">
      {{ message }}
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';

export default {
  name: 'MCPView',
  setup() {
    const testing = ref(false);
    const message = ref('');
    const messageType = ref('success');
    
    // 获取MCP配置（用于配置显示）
    const mcpConfig = computed(() => {
      // 使用占位符，用户需要从API密钥管理页面获取实际密钥
      const apiKey = 'your-api-key-here';
      const serverUrl = `${window.location.origin}/mcp`;
      
      return `{
  "mcpServers": {
    "ai-doc-library": {
      "type": "streamable-http",
      "url": "${serverUrl}",
      "headers": {
        "API_KEY": "${apiKey}"
      }
    }
  }
}`;
    });
    
    // 测试MCP连接
    const testConnection = async () => {
      testing.value = true;
      try {
        showMessage('请先从API密钥管理页面获取有效的API密钥', 'info');
      } catch (error) {
        showMessage('MCP连接测试失败: ' + error.message, 'error');
      } finally {
        testing.value = false;
      }
    };

    // 复制配置到剪贴板
    const copyConfig = async () => {
      try {
        await navigator.clipboard.writeText(mcpConfig.value);
        showMessage('配置已复制到剪贴板', 'success');
      } catch (error) {
        showMessage('复制失败，请手动复制', 'error');
      }
    };
    
    // 显示消息
    const showMessage = (text, type = 'success') => {
      message.value = text;
      messageType.value = type;
      setTimeout(() => {
        message.value = '';
      }, 3000);
    };
    
    return {
      testing,
      message,
      messageType,
      mcpConfig,
      testConnection,
      copyConfig,
      showMessage
    };
  }
};
</script>

<style scoped>
.mcp-view {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h1 {
  margin: 0 0 10px 0;
  color: #333;
}

.page-header p {
  margin: 0;
  color: #666;
}

.content-container {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
}

.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
}

.card-header h2 {
  margin: 0;
  color: #333;
}

.card-body {
  padding: 20px;
}

.config-box {
  background: #f5f5f5;
  border-radius: 4px;
  padding: 15px;
  margin: 15px 0;
  overflow-x: auto;
}

.config-box pre {
  margin: 0;
  white-space: pre-wrap;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.4;
}

.config-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
}

.api-keys-table th,
.api-keys-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.api-keys-table th {
  background-color: #f8f9fa;
  font-weight: 600;
}

.key-display {
  font-family: 'Courier New', monospace;
  font-size: 14px;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  margin-left: 8px;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.btn-icon:hover {
  background-color: #f0f0f0;
}

.btn-danger:hover {
  background-color: #ffdddd;
}

.instructions {
  line-height: 1.6;
}

.instructions h3 {
  margin-top: 20px;
  margin-bottom: 10px;
}

.instructions ol,
.instructions ul {
  padding-left: 20px;
}

.instructions li {
  margin-bottom: 8px;
}

.modal-overlay {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex !important;
  align-items: center;
  justify-content: center;
  z-index: 9999 !important;
}

.modal {
  background: white !important;
  border-radius: 8px !important;
  width: 90% !important;
  max-width: 500px !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15) !important;
  display: block !important;
  position: relative !important;
  min-height: 300px !important;
  z-index: 10000 !important;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
}

.btn-close {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 0;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
}

.form-control {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 15px 20px;
  border-top: 1px solid #eee;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-primary:hover {
  background-color: #0069d9;
}

.btn-primary:disabled {
  background-color: #a0c3ff;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.message {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: 4px;
  color: white;
  font-weight: 500;
  z-index: 1001;
  max-width: 300px;
}

.message.success {
  background-color: #28a745;
}

.message.error {
  background-color: #dc3545;
}

@media (max-width: 768px) {
  .mcp-view {
    padding: 10px;
  }
  
  .content-container {
    grid-template-columns: 1fr;
  }
  
  .config-actions {
    flex-direction: column;
  }
  
  .api-keys-table {
    font-size: 14px;
  }
  
  .api-keys-table th,
  .api-keys-table td {
    padding: 8px 10px;
  }
}
</style>