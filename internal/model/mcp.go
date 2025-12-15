package model

import (
	"time"
)

// MCPRequest MCP协议请求结构
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"` // 固定为"2.0"
	ID      interface{} `json:"id"`      // 请求ID，可以是字符串或数字
	Method  string      `json:"method"`  // 方法名
	Params  interface{} `json:"params"`  // 参数，根据方法不同而变化
}

// MCPResponse MCP协议响应结构
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"` // 固定为"2.0"
	ID      interface{} `json:"id"`      // 对应请求的ID
	Result  interface{} `json:"result"`  // 成功时的结果
	Error   *MCPError   `json:"error"`   // 错误时的错误信息
}

// MCPError MCP协议错误结构
type MCPError struct {
	Code    int         `json:"code"`    // 错误代码
	Message string      `json:"message"` // 错误消息
	Data    interface{} `json:"data"`    // 错误详细数据（可选）
}

// MCPInitializeParams MCP初始化参数
type MCPInitializeParams struct {
	ProtocolVersion       string                 `json:"protocolVersion"`
	Capabilities          map[string]interface{} `json:"capabilities"`
	ClientInfo            MCPClientInfo          `json:"clientInfo"`
	Trace                 string                 `json:"trace"`
	Locale                string                 `json:"locale"`
	Root                  *MCPRoot               `json:"root"`
	RootPath              string                 `json:"rootPath"`
	WorkspaceFolders      []MCPWorkspaceFolder   `json:"workspaceFolders"`
	InitializationOptions map[string]interface{} `json:"initializationOptions"`
}

// MCPClientInfo MCP客户端信息
type MCPClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// MCPRoot MCP根目录信息
type MCPRoot struct {
	URI string `json:"uri"`
}

// MCPWorkspaceFolder MCP工作区文件夹
type MCPWorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// MCPCapabilities MCP服务器能力
type MCPCapabilities struct {
	Tools     *MCPToolsCapability     `json:"tools,omitempty"`
	Resources *MCPResourcesCapability `json:"resources,omitempty"`
	Prompts   *MCPPromptsCapability   `json:"prompts,omitempty"`
	Logging   *MCPLoggingCapability   `json:"logging,omitempty"`
	Sampling  *MCPSamplingCapability  `json:"sampling,omitempty"`
}

// MCPToolsCapability MCP工具能力
type MCPToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// MCPResourcesCapability MCP资源能力
type MCPResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// MCPPromptsCapability MCP提示能力
type MCPPromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// MCPLoggingCapability MCP日志能力
type MCPLoggingCapability struct{}

// MCPSamplingCapability MCP采样能力
type MCPSamplingCapability struct{}

// MCPInitializeResult MCP初始化结果
type MCPInitializeResult struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    MCPCapabilities `json:"capabilities"`
	ServerInfo      MCPServerInfo   `json:"serverInfo"`
}

// MCPServerInfo MCP服务器信息
type MCPServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// MCPTool MCP工具定义
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPToolListParams MCP工具列表参数
type MCPToolListParams struct {
	Cursor *string `json:"cursor,omitempty"`
}

// MCPToolListResult MCP工具列表结果
type MCPToolListResult struct {
	Tools      []MCPTool `json:"tools"`
	NextCursor *string   `json:"nextCursor,omitempty"`
}

// MCPToolCallParams MCP工具调用参数
type MCPToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// MCPToolResult MCP工具结果
type MCPToolResult struct {
	Content []interface{} `json:"content"`
	IsError bool          `json:"isError,omitempty"`
	_Meta   interface{}   `json:"_meta,omitempty"`
}

// MCPTextContent MCP文本内容
type MCPTextContent struct {
	Type string `json:"type"` // 固定为"text"
	Text string `json:"text"`
}

// MCPSearchToolParams MCP搜索工具参数
type MCPSearchToolParams struct {
	Query   string   `json:"query"`
	Types   []string `json:"types,omitempty"`
	Version string   `json:"version,omitempty"`
	Limit   int      `json:"limit,omitempty"`
}

// MCPSearchToolResult MCP搜索工具结果
type MCPSearchToolResult struct {
	Documents []MCPSearchDocument `json:"documents"`
	Total     int                 `json:"total"`
	Query     string              `json:"query"`
}

// MCPSearchDocument MCP搜索文档
type MCPSearchDocument struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Version     string  `json:"version"`
	Score       float64 `json:"score"`
	Content     string  `json:"content"`
	ContentType string  `json:"content_type"`
	Section     string  `json:"section"`
}

// MCPAPIKey API密钥模型
type MCPAPIKey struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string     `json:"name" gorm:"not null"`
	Key       string     `json:"key" gorm:"uniqueIndex;not null"`
	UserID    string     `json:"user_id" gorm:"not null"`
	Enabled   bool       `json:"enabled" gorm:"default:true"`
	ExpiresAt *time.Time `json:"expires_at"`
	LastUsed  *time.Time `json:"last_used"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName 返回API密钥表名
func (MCPAPIKey) TableName() string {
	return "mcp_api_keys"
}

// MCPConfig MCP配置模型
type MCPConfig struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Endpoint    string    `json:"endpoint" gorm:"not null"`
	APIKey      string    `json:"api_key" gorm:"not null"`
	Enabled     bool      `json:"enabled" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 返回MCP配置表名
func (MCPConfig) TableName() string {
	return "mcp_configs"
}
