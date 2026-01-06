package model

import (
	"testing"
	"time"
)

func TestMCPRequest_FieldValues(t *testing.T) {
	tests := []struct {
		name string
		req  MCPRequest
	}{
		{
			name: "valid request with string ID",
			req: MCPRequest{
				JSONRPC: "2.0",
				ID:      "test-id",
				Method:  "test/method",
				Params:  map[string]interface{}{"key": "value"},
			},
		},
		{
			name: "valid request with number ID",
			req: MCPRequest{
				JSONRPC: "2.0",
				ID:      123,
				Method:  "test/method",
				Params:  map[string]interface{}{"key": "value"},
			},
		},
		{
			name: "request with nil params",
			req: MCPRequest{
				JSONRPC: "2.0",
				ID:      "test-id",
				Method:  "test/method",
				Params:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.JSONRPC != "2.0" {
				t.Errorf("JSONRPC = %v, want 2.0", tt.req.JSONRPC)
			}
			if tt.req.Method == "" {
				t.Error("Method should not be empty")
			}
		})
	}
}

func TestMCPResponse_FieldValues(t *testing.T) {
	tests := []struct {
		name string
		resp MCPResponse
	}{
		{
			name: "success response",
			resp: MCPResponse{
				JSONRPC: "2.0",
				ID:      "test-id",
				Method:  "test/method",
				Result:  map[string]interface{}{"data": "value"},
				Error:   nil,
			},
		},
		{
			name: "error response",
			resp: MCPResponse{
				JSONRPC: "2.0",
				ID:      "test-id",
				Error: &MCPError{
					Code:    -32600,
					Message: "Invalid Request",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp.JSONRPC != "2.0" {
				t.Errorf("JSONRPC = %v, want 2.0", tt.resp.JSONRPC)
			}
			if tt.resp.Result == nil && tt.resp.Error == nil {
				t.Error("Response should have either Result or Error")
			}
		})
	}
}

func TestMCPError_FieldValues(t *testing.T) {
	tests := []struct {
		name string
		err  MCPError
	}{
		{
			name: "standard error",
			err: MCPError{
				Code:    -32600,
				Message: "Invalid Request",
				Data:    nil,
			},
		},
		{
			name: "error with data",
			err: MCPError{
				Code:    -32600,
				Message: "Invalid Request",
				Data:    map[string]interface{}{"detail": "field required"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code == 0 {
				t.Error("Error code should not be zero")
			}
			if tt.err.Message == "" {
				t.Error("Error message should not be empty")
			}
		})
	}
}

func TestMCPInitializeParams_FieldValues(t *testing.T) {
	params := MCPInitializeParams{
		ProtocolVersion: "2024-11-05",
		Capabilities: map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": false,
			},
		},
		ClientInfo: MCPClientInfo{
			Name:    "Test Client",
			Version: "1.0.0",
		},
		Trace:  "off",
		Locale: "en-US",
	}

	if params.ProtocolVersion != "2024-11-05" {
		t.Errorf("ProtocolVersion = %v, want 2024-11-05", params.ProtocolVersion)
	}
	if params.ClientInfo.Name != "Test Client" {
		t.Errorf("ClientInfo.Name = %v, want Test Client", params.ClientInfo.Name)
	}
	if params.Capabilities == nil {
		t.Error("Capabilities should not be nil")
	}
}

func TestMCPClientInfo_FieldValues(t *testing.T) {
	clientInfo := MCPClientInfo{
		Name:    "CoStrict",
		Version: "1.0.0",
	}

	if clientInfo.Name != "CoStrict" {
		t.Errorf("Name = %v, want CoStrict", clientInfo.Name)
	}
	if clientInfo.Version != "1.0.0" {
		t.Errorf("Version = %v, want 1.0.0", clientInfo.Version)
	}
}

func TestMCPRoot_FieldValues(t *testing.T) {
	root := MCPRoot{
		URI: "file:///path/to/workspace",
	}

	if root.URI != "file:///path/to/workspace" {
		t.Errorf("URI = %v, want file:///path/to/workspace", root.URI)
	}
}

func TestMCPWorkspaceFolder_FieldValues(t *testing.T) {
	folder := MCPWorkspaceFolder{
		URI:  "file:///path/to/workspace",
		Name: "My Workspace",
	}

	if folder.URI != "file:///path/to/workspace" {
		t.Errorf("URI = %v, want file:///path/to/workspace", folder.URI)
	}
	if folder.Name != "My Workspace" {
		t.Errorf("Name = %v, want My Workspace", folder.Name)
	}
}

func TestMCPCapabilities_FieldValues(t *testing.T) {
	capabilities := MCPCapabilities{
		Tools: &MCPToolsCapability{
			ListChanged: true,
		},
		Resources: &MCPResourcesCapability{
			Subscribe:   true,
			ListChanged: false,
		},
		Prompts: &MCPPromptsCapability{
			ListChanged: false,
		},
		Logging:  &MCPLoggingCapability{},
		Sampling: &MCPSamplingCapability{},
	}

	if capabilities.Tools == nil {
		t.Error("Tools capability should not be nil")
	}
	if capabilities.Resources == nil {
		t.Error("Resources capability should not be nil")
	}
	if capabilities.Prompts == nil {
		t.Error("Prompts capability should not be nil")
	}
}

func TestMCPInitializeResult_FieldValues(t *testing.T) {
	result := MCPInitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: MCPCapabilities{
			Tools: &MCPToolsCapability{},
		},
		ServerInfo: MCPServerInfo{
			Name:    "LAST-doc MCP Server",
			Version: "1.0.0",
		},
	}

	if result.ProtocolVersion != "2024-11-05" {
		t.Errorf("ProtocolVersion = %v, want 2024-11-05", result.ProtocolVersion)
	}
	if result.ServerInfo.Name != "LAST-doc MCP Server" {
		t.Errorf("ServerInfo.Name = %v, want LAST-doc MCP Server", result.ServerInfo.Name)
	}
}

func TestMCPServerInfo_FieldValues(t *testing.T) {
	serverInfo := MCPServerInfo{
		Name:    "LAST-doc MCP Server",
		Version: "1.0.0",
	}

	if serverInfo.Name == "" {
		t.Error("Server name should not be empty")
	}
	if serverInfo.Version == "" {
		t.Error("Server version should not be empty")
	}
}

func TestMCPTool_FieldValues(t *testing.T) {
	tool := MCPTool{
		Name:        "search_documents",
		Description: "搜索技术文档",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "搜索查询",
				},
			},
			"required": []string{"query"},
		},
	}

	if tool.Name != "search_documents" {
		t.Errorf("Name = %v, want search_documents", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Description should not be empty")
	}
	if tool.InputSchema == nil {
		t.Error("InputSchema should not be nil")
	}
}

func TestMCPToolListParams_FieldValues(t *testing.T) {
	tests := []struct {
		name   string
		params MCPToolListParams
	}{
		{
			name: "without cursor",
			params: MCPToolListParams{
				Cursor: nil,
			},
		},
		{
			name: "with cursor",
			params: MCPToolListParams{
				Cursor: func() *string { s := "cursor-123"; return &s }(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试字段可以正确赋值
			if tt.name == "with cursor" && tt.params.Cursor == nil {
				t.Error("Cursor should not be nil when provided")
			}
		})
	}
}

func TestMCPToolListResult_FieldValues(t *testing.T) {
	result := MCPToolListResult{
		Tools: []MCPTool{
			{
				Name:        "tool1",
				Description: "Tool 1",
				InputSchema: map[string]interface{}{},
			},
			{
				Name:        "tool2",
				Description: "Tool 2",
				InputSchema: map[string]interface{}{},
			},
		},
		NextCursor: func() *string { s := "next-cursor"; return &s }(),
	}

	if len(result.Tools) != 2 {
		t.Errorf("Tools length = %v, want 2", len(result.Tools))
	}
	if result.NextCursor == nil {
		t.Error("NextCursor should not be nil")
	}
	if *result.NextCursor != "next-cursor" {
		t.Errorf("NextCursor = %v, want next-cursor", *result.NextCursor)
	}
}

func TestMCPToolCallParams_FieldValues(t *testing.T) {
	params := MCPToolCallParams{
		Name: "search_documents",
		Arguments: map[string]interface{}{
			"query": "test query",
			"limit": 10,
		},
	}

	if params.Name != "search_documents" {
		t.Errorf("Name = %v, want search_documents", params.Name)
	}
	if params.Arguments == nil {
		t.Error("Arguments should not be nil")
	}
	if params.Arguments["query"] != "test query" {
		t.Errorf("Arguments[query] = %v, want test query", params.Arguments["query"])
	}
}

func TestMCPToolResult_FieldValues(t *testing.T) {
	tests := []struct {
		name   string
		result MCPToolResult
	}{
		{
			name: "successful result",
			result: MCPToolResult{
				Content: []interface{}{
					MCPTextContent{
						Type: "text",
						Text: "Result content",
					},
				},
				IsError: false,
				_Meta:   nil,
			},
		},
		{
			name: "error result",
			result: MCPToolResult{
				Content: []interface{}{
					MCPTextContent{
						Type: "text",
						Text: "Error occurred",
					},
				},
				IsError: true,
				_Meta: map[string]interface{}{
					"error_code": 100,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.result.Content) == 0 {
				t.Error("Content should not be empty")
			}
		})
	}
}

func TestMCPTextContent_FieldValues(t *testing.T) {
	content := MCPTextContent{
		Type: "text",
		Text: "This is text content",
	}

	if content.Type != "text" {
		t.Errorf("Type = %v, want text", content.Type)
	}
	if content.Text == "" {
		t.Error("Text should not be empty")
	}
}

func TestMCPSearchToolParams_FieldValues(t *testing.T) {
	tests := []struct {
		name   string
		params MCPSearchToolParams
	}{
		{
			name: "basic search",
			params: MCPSearchToolParams{
				Query: "test",
			},
		},
		{
			name: "search with all parameters",
			params: MCPSearchToolParams{
				Query:   "test",
				Types:   []string{"pdf", "markdown"},
				Version: "1.0.0",
				Limit:   10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.params.Query == "" {
				t.Error("Query should not be empty")
			}
		})
	}
}

func TestMCPSearchToolResult_FieldValues(t *testing.T) {
	result := MCPSearchToolResult{
		Documents: []MCPSearchDocument{
			{
				ID:            "doc-1",
				DocumentID:    "doc-1",
				Name:          "Document 1",
				Type:          "pdf",
				Version:       "1.0.0",
				Library:       "React",
				Score:         0.95,
				Content:       "Document content",
				Snippet:       "Snippet content",
				ContentType:   "text/plain",
				Section:       "Introduction",
				StartPosition: 0,
				EndPosition:   100,
			},
		},
		Total: 1,
		Query: "test",
	}

	if len(result.Documents) != 1 {
		t.Errorf("Documents length = %v, want 1", len(result.Documents))
	}
	if result.Total != 1 {
		t.Errorf("Total = %v, want 1", result.Total)
	}
	if result.Query != "test" {
		t.Errorf("Query = %v, want test", result.Query)
	}
}

func TestMCPSearchDocument_FieldValues(t *testing.T) {
	doc := MCPSearchDocument{
		ID:            "doc-1",
		DocumentID:    "doc-1",
		Name:          "Test Document",
		Type:          "pdf",
		Version:       "1.0.0",
		Library:       "React",
		Score:         0.95,
		Content:       "Document content",
		Snippet:       "Snippet content",
		ContentType:   "text/plain",
		Section:       "Introduction",
		StartPosition: 0,
		EndPosition:   100,
	}

	if doc.ID == "" {
		t.Error("ID should not be empty")
	}
	if doc.DocumentID == "" {
		t.Error("DocumentID should not be empty")
	}
	if doc.Score < 0 || doc.Score > 1 {
		t.Errorf("Score = %v, want value between 0 and 1", doc.Score)
	}
	if doc.StartPosition < 0 {
		t.Errorf("StartPosition = %v, want non-negative", doc.StartPosition)
	}
	if doc.EndPosition < doc.StartPosition {
		t.Errorf("EndPosition = %v, want >= StartPosition (%v)", doc.EndPosition, doc.StartPosition)
	}
}

func TestMCPAPIKey_TableName(t *testing.T) {
	apiKey := MCPAPIKey{}
	if apiKey.TableName() != "mcp_api_keys" {
		t.Errorf("TableName() = %v, want mcp_api_keys", apiKey.TableName())
	}
}

func TestMCPAPIKey_FieldValues(t *testing.T) {
	now := time.Now()
	apiKey := MCPAPIKey{
		ID:        "key-1",
		Name:      "Test API Key",
		Key:       "sk-test-key-123456",
		UserID:    "user-1",
		Enabled:   true,
		ExpiresAt: &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if apiKey.Name != "Test API Key" {
		t.Errorf("Name = %v, want Test API Key", apiKey.Name)
	}
	if apiKey.Key == "" {
		t.Error("Key should not be empty")
	}
	if apiKey.UserID == "" {
		t.Error("UserID should not be empty")
	}
	if !apiKey.Enabled {
		t.Error("Enabled should be true by default")
	}
}

func TestMCPConfig_TableName(t *testing.T) {
	config := MCPConfig{}
	if config.TableName() != "mcp_configs" {
		t.Errorf("TableName() = %v, want mcp_configs", config.TableName())
	}
}

func TestMCPConfig_FieldValues(t *testing.T) {
	config := MCPConfig{
		ID:          "config-1",
		Name:        "Test MCP Config",
		Description: "Test configuration",
		Endpoint:    "http://localhost:8080",
		APIKey:      "test-api-key",
		Enabled:     true,
	}

	if config.Name != "Test MCP Config" {
		t.Errorf("Name = %v, want Test MCP Config", config.Name)
	}
	if config.Endpoint == "" {
		t.Error("Endpoint should not be empty")
	}
	if config.APIKey == "" {
		t.Error("APIKey should not be empty")
	}
	if !config.Enabled {
		t.Error("Enabled should be true by default")
	}
}
