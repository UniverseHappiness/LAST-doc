import axios from 'axios';

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 可以在这里添加认证头等
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    // 统一错误处理
    const message = error.response?.data?.error || error.message || '请求失败';
    return Promise.reject(new Error(message));
  }
);

// MCP服务对象
const mcpService = {
  // 获取MCP配置
  async getConfig(apiKey) {
    try {
      const response = await api.get('/mcp/config', {
        headers: {
          'API_KEY': apiKey
        }
      });
      return response;
    } catch (error) {
      throw new Error(`获取MCP配置失败: ${error.message}`);
    }
  },

  // 测试MCP连接
  async testConnection(apiKey) {
    try {
      const response = await api.get('/mcp/test', {
        headers: {
          'API_KEY': apiKey
        }
      });
      return response;
    } catch (error) {
      throw new Error(`测试MCP连接失败: ${error.message}`);
    }
  },

  // 发送MCP消息
  async sendMessage(message, apiKey) {
    try {
      const response = await api.post('/mcp', message, {
        headers: {
          'Content-Type': 'application/json',
          'API_KEY': apiKey
        }
      });
      return response;
    } catch (error) {
      throw new Error(`发送MCP消息失败: ${error.message}`);
    }
  },

  // 创建API密钥
  async createApiKey(keyData) {
    try {
      const response = await api.post('/mcp/keys', keyData);
      return response;
    } catch (error) {
      throw new Error(`创建API密钥失败: ${error.message}`);
    }
  },

  // 获取API密钥列表
  async getApiKeys(userId) {
    try {
      const response = await api.get('/mcp/keys', {
        params: {
          user_id: userId
        }
      });
      return response;
    } catch (error) {
      throw new Error(`获取API密钥列表失败: ${error.message}`);
    }
  },

  // 删除API密钥
  async deleteApiKey(keyId) {
    try {
      const response = await api.delete(`/mcp/keys/${keyId}`);
      return response;
    } catch (error) {
      throw new Error(`删除API密钥失败: ${error.message}`);
    }
  },

  // 搜索文档（通过MCP协议）
  async searchDocuments(query, options = {}) {
    const { types, version, limit = 10 } = options;
    
    // 这里需要一个有效的API密钥，实际使用时应该从用户配置中获取
    const apiKey = localStorage.getItem('mcp_api_key') || 'demo-key';
    
    try {
      const message = {
        jsonrpc: '2.0',
        id: Date.now().toString(),
        method: 'tools/call',
        params: {
          name: 'search_documents',
          arguments: {
            query,
            types,
            version,
            limit
          }
        }
      };

      const response = await this.sendMessage(message, apiKey);
      return response.result;
    } catch (error) {
      throw new Error(`搜索文档失败: ${error.message}`);
    }
  },

  // 获取文档内容（通过MCP协议）
  async getDocumentContent(documentId, version) {
    // 这里需要一个有效的API密钥，实际使用时应该从用户配置中获取
    const apiKey = localStorage.getItem('mcp_api_key') || 'demo-key';
    
    try {
      const message = {
        jsonrpc: '2.0',
        id: Date.now().toString(),
        method: 'tools/call',
        params: {
          name: 'get_document_content',
          arguments: {
            document_id: documentId,
            version
          }
        }
      };

      const response = await this.sendMessage(message, apiKey);
      return response.result;
    } catch (error) {
      throw new Error(`获取文档内容失败: ${error.message}`);
    }
  },

  // 获取可用工具列表（通过MCP协议）
  async getTools() {
    // 这里需要一个有效的API密钥，实际使用时应该从用户配置中获取
    const apiKey = localStorage.getItem('mcp_api_key') || 'demo-key';
    
    try {
      const message = {
        jsonrpc: '2.0',
        id: Date.now().toString(),
        method: 'tools/list',
        params: {}
      };

      const response = await this.sendMessage(message, apiKey);
      return response.result;
    } catch (error) {
      throw new Error(`获取工具列表失败: ${error.message}`);
    }
  },

  // 初始化MCP连接
  async initialize(clientInfo) {
    // 这里需要一个有效的API密钥，实际使用时应该从用户配置中获取
    const apiKey = localStorage.getItem('mcp_api_key') || 'demo-key';
    
    try {
      const message = {
        jsonrpc: '2.0',
        id: Date.now().toString(),
        method: 'initialize',
        params: {
          protocolVersion: '2024-11-05',
          capabilities: {},
          clientInfo: clientInfo || {
            name: 'AI技术文档库Web客户端',
            version: '1.0.0'
          }
        }
      };

      const response = await this.sendMessage(message, apiKey);
      return response.result;
    } catch (error) {
      throw new Error(`初始化MCP连接失败: ${error.message}`);
    }
  }
};

export { mcpService };