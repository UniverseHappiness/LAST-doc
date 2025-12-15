# MCP本地使用指南

本指南将帮助您在本地环境中设置和使用MCP（Model Context Protocol）来访问AI技术文档库。

## 前提条件

1. 已安装并启动了AI技术文档库后端服务
2. 已有有效的用户账户
3. 已创建了API密钥（通过Web界面的API密钥管理页面）

## 获取API密钥

1. 登录AI技术文档库Web界面
2. 导航到"API密钥管理"页面
3. 点击"创建新密钥"按钮
4. 输入密钥名称并设置过期时间（可选）
5. 复制生成的API密钥并妥善保存

## MCP配置方法

### 方法一：使用Claude Desktop

1. 找到Claude Desktop的配置文件：
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`

2. 编辑配置文件，添加以下内容：

```json
{
  "mcpServers": {
    "ai-doc-library": {
      "command": "node",
      "args": ["path/to/mcp-server.js"],
      "env": {
        "API_KEY": "your-api-key-here",
        "SERVER_URL": "http://localhost:8080"
      }
    }
  }
}
```

### 方法二：使用自定义MCP客户端

如果您有自己的MCP客户端，可以使用以下HTTP配置：

```json
{
  "mcpServers": {
    "ai-doc-library": {
      "type": "streamable-http",
      "url": "http://localhost:8080/mcp",
      "headers": {
        "API_KEY": "your-api-key-here"
      }
    }
  }
}
```

## 可用的MCP工具

### 1. search_documents

搜索技术文档，支持关键词和语义搜索。

**参数:**
- `query` (必需): 搜索查询关键词
- `types` (可选): 文档类型过滤器，如 ["pdf", "docx", "markdown"]
- `version` (可选): 文档版本过滤器
- `limit` (可选): 返回结果数量限制，默认为10

**示例请求:**
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/call",
  "params": {
    "name": "search_documents",
    "arguments": {
      "query": "Vue组件开发",
      "types": ["markdown"],
      "limit": 5
    }
  }
}
```

### 2. get_document_content

获取指定文档的详细内容。

**参数:**
- `document_id` (必需): 文档ID
- `version` (可选): 文档版本，如果未指定则使用最新版本

**示例请求:**
```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "method": "tools/call",
  "params": {
    "name": "get_document_content",
    "arguments": {
      "document_id": "doc-123",
      "version": "v1.0"
    }
  }
}
```

## 测试MCP连接

您可以使用以下方法测试MCP连接是否正常：

### 使用curl测试

1. 测试MCP连接：
```bash
curl -X GET http://localhost:8080/api/v1/mcp/test \
  -H "API_KEY: your-api-key-here"
```

2. 发送MCP请求：
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "API_KEY: your-api-key-here" \
  -d '{
    "jsonrpc": "2.0",
    "id": "test",
    "method": "tools/list",
    "params": {}
  }'
```

### 使用Python脚本测试

创建一个简单的Python测试脚本：

```python
import requests
import json

# 配置
API_KEY = "your-api-key-here"
BASE_URL = "http://localhost:8080"

# 测试MCP连接
def test_mcp_connection():
    headers = {
        "API_KEY": API_KEY,
        "Content-Type": "application/json"
    }
    
    # 测试连接
    response = requests.get(f"{BASE_URL}/api/v1/mcp/test", headers=headers)
    print("连接测试:", response.json())
    
    # 获取工具列表
    tools_request = {
        "jsonrpc": "2.0",
        "id": "1",
        "method": "tools/list",
        "params": {}
    }
    
    response = requests.post(f"{BASE_URL}/mcp", headers=headers, json=tools_request)
    print("工具列表:", response.json())
    
    # 搜索文档
    search_request = {
        "jsonrpc": "2.0",
        "id": "2",
        "method": "tools/call",
        "params": {
            "name": "search_documents",
            "arguments": {
                "query": "测试",
                "limit": 3
            }
        }
    }
    
    response = requests.post(f"{BASE_URL}/mcp", headers=headers, json=search_request)
    print("搜索结果:", response.json())

if __name__ == "__main__":
    test_mcp_connection()
```

## 常见问题

### 1. API密钥无效

确保：
- API密钥已正确复制
- API密钥未过期
- 用户账户仍然活跃

### 2. 连接超时

检查：
- 后端服务是否正在运行
- 端口8080是否可访问
- 防火墙设置是否正确

### 3. 文档搜索无结果

确认：
- 数据库中是否有文档
- 文档已建立搜索索引
- 搜索关键词是否正确

## 高级用法

### 自定义MCP服务器

您可以创建自己的MCP服务器来与AI技术文档库交互：

```javascript
// mcp-server.js
const { spawn } = require('child_process');
const http = require('http');

const API_KEY = process.env.API_KEY;
const SERVER_URL = process.env.SERVER_URL || 'http://localhost:8080';

// MCP服务器逻辑
const server = http.createServer((req, res) => {
  // 处理MCP协议请求
  // ...
});

server.listen(3000, () => {
  console.log('MCP服务器运行在端口3000');
});
```

### 集成到其他应用

您可以将MCP协议集成到其他应用程序中，实现与AI技术文档库的无缝交互。

## 安全注意事项

1. **保护API密钥**: 不要在代码中硬编码API密钥，使用环境变量或安全的密钥管理
2. **权限控制**: 确保API密钥只有必要的权限
3. **定期轮换**: 定期更换API密钥以提高安全性
4. **监控使用**: 监控API密钥的使用情况，及时发现异常

## 更多资源

- [MCP协议官方文档](https://modelcontextprotocol.io/)
- [AI技术文档库API文档](./api_documentation.md)
- [Claude Desktop集成指南](https://docs.anthropic.com/claude/docs/mcp)

如有问题，请查看日志文件或联系技术支持。