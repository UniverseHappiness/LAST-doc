package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/service"
	"github.com/gin-gonic/gin"
)

// MCPHandler MCP处理器
type MCPHandler struct {
	mcpService service.MCPService
}

// NewMCPHandler 创建MCP处理器实例
func NewMCPHandler(mcpService service.MCPService) *MCPHandler {
	return &MCPHandler{
		mcpService: mcpService,
	}
}

// HandleMCPRequest 处理MCP请求
func (h *MCPHandler) HandleMCPRequest(c *gin.Context) {
	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// 解析MCP请求
	var req model.MCPRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON-RPC request",
		})
		return
	}

	// 验证JSON-RPC版本
	if req.JSONRPC != "2.0" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported JSON-RPC version",
		})
		return
	}

	// 获取API密钥
	apiKey := c.GetHeader("API_KEY")
	if apiKey == "" {
		// 尝试从Authorization头获取
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "API key is required",
		})
		return
	}

	// 处理MCP请求
	ctx := context.Background()
	resp, err := h.mcpService.HandleRequest(ctx, &req, apiKey)
	if err != nil {
		log.Printf("Error handling MCP request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// 返回MCP响应
	c.JSON(http.StatusOK, resp)
}

// GetMCPConfig 获取MCP配置
func (h *MCPHandler) GetMCPConfig(c *gin.Context) {
	// 获取API密钥
	apiKey := c.GetHeader("API_KEY")
	if apiKey == "" {
		// 尝试从Authorization头获取
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "API key is required",
		})
		return
	}

	// 验证API密钥
	_, err := h.mcpService.ValidateAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid API key",
		})
		return
	}

	// 获取服务器地址
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	serverURL := fmt.Sprintf("%s://%s", scheme, host)

	// 返回MCP配置
	config := gin.H{
		"mcpServers": gin.H{
			"ai-doc-library": gin.H{
				"type": "streamable-http",
				"url":  serverURL + "/mcp",
				"headers": gin.H{
					"API_KEY": apiKey,
				},
			},
		},
	}

	c.JSON(http.StatusOK, config)
}

// CreateAPIKey 创建API密钥
func (h *MCPHandler) CreateAPIKey(c *gin.Context) {
	// 检查用户认证
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权访问",
		})
		return
	}

	// 解析请求体
	var req struct {
		Name      string  `json:"name" binding:"required"`
		ExpiresAt *string `json:"expires_at,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	// 解析过期时间
	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		// 尝试解析ISO 8601格式的时间字符串
		parsedTime, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid expires_at format, expected RFC3339 format",
			})
			return
		}
		expiresAt = &parsedTime
	}

	// 创建API密钥
	ctx := context.Background()
	apiKey, err := h.mcpService.CreateAPIKey(ctx, req.Name, userID.(string), expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create API key: %v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, apiKey)
}

// GetAPIKeys 获取API密钥列表
func (h *MCPHandler) GetAPIKeys(c *gin.Context) {
	// 检查用户认证
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权访问",
		})
		return
	}

	// 管理员可以查看所有用户的API密钥
	var targetUserID string
	if role, exists := c.Get("role"); exists && role == "admin" {
		if queryUserID := c.Query("user_id"); queryUserID != "" {
			targetUserID = queryUserID
		} else {
			targetUserID = userID.(string)
		}
	} else {
		targetUserID = userID.(string)
	}

	// 获取API密钥列表
	ctx := context.Background()
	keys, err := h.mcpService.GetAPIKeys(ctx, targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get API keys: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keys": keys,
	})
}

// DeleteAPIKey 删除API密钥
func (h *MCPHandler) DeleteAPIKey(c *gin.Context) {
	// 检查用户认证
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权访问",
		})
		return
	}

	// 获取密钥ID
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "key_id is required",
		})
		return
	}

	// 获取API密钥信息以验证权限
	ctx := context.Background()
	apiKey, err := h.mcpService.ValidateAPIKeyByKeyID(keyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "API密钥不存在",
		})
		return
	}

	// 检查权限：只有管理员或密钥所有者可以删除
	userRole, _ := c.Get("role")
	if userRole != "admin" && apiKey.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
		})
		return
	}

	// 删除API密钥
	err = h.mcpService.DeleteAPIKey(ctx, keyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to delete API key: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key deleted successfully",
	})
}

// TestMCPConnection 测试MCP连接
func (h *MCPHandler) TestMCPConnection(c *gin.Context) {
	// 获取API密钥
	apiKey := c.GetHeader("API_KEY")
	if apiKey == "" {
		// 尝试从Authorization头获取
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "API key is required",
		})
		return
	}

	// 验证API密钥
	key, err := h.mcpService.ValidateAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid API key",
		})
		return
	}

	// 更新API密钥最后使用时间
	ctx := context.Background()
	if err := h.mcpService.UpdateAPIKeyLastUsed(ctx, key.ID); err != nil {
		log.Printf("Failed to update API key last used time: %v", err)
	}

	// 构建测试请求
	testReq := &model.MCPRequest{
		JSONRPC: "2.0",
		ID:      "test",
		Method:  "tools/list",
		Params:  nil,
	}

	// 处理测试请求
	resp, err := h.mcpService.HandleRequest(ctx, testReq, apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("MCP connection test failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "MCP connection test successful",
		"response": resp,
	})
}

// SendMessage 发送消息到MCP服务器
func (h *MCPHandler) SendMessage(c *gin.Context) {
	// 获取API密钥
	apiKey := c.GetHeader("API_KEY")
	if apiKey == "" {
		// 尝试从Authorization头获取
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "API key is required",
		})
		return
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// 验证是否是JSON-RPC请求
	if !bytes.HasPrefix(body, []byte(`{"jsonrpc":"2.0"`)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON-RPC request format",
		})
		return
	}

	// 解析MCP请求
	var req model.MCPRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON-RPC request",
		})
		return
	}

	// 验证JSON-RPC版本
	if req.JSONRPC != "2.0" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported JSON-RPC version",
		})
		return
	}

	// 处理MCP请求
	ctx := context.Background()
	resp, err := h.mcpService.HandleRequest(ctx, &req, apiKey)
	if err != nil {
		log.Printf("Error handling MCP request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// 返回MCP响应
	c.JSON(http.StatusOK, resp)
}

// getRequestBody 获取请求体内容（用于调试）
func getRequestBody(c *gin.Context) string {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Sprintf("读取请求体失败: %v", err)
	}
	// 重新设置请求体，以便后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body)
}
