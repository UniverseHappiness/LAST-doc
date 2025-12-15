package middleware

import (
	"net/http"
	"strings"

	"github.com/UniverseHappiness/LAST-doc/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	userService service.UserService
}

// NewAuthMiddleware 创建认证中间件实例
func NewAuthMiddleware(userService service.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少授权令牌",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的授权格式",
			})
			c.Abort()
			return
		}

		// 提取令牌
		tokenString := authHeader[7:]

		// 验证令牌
		claims, err := m.userService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的授权令牌",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireRole 需要特定角色的中间件
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先进行认证
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// 获取用户角色
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "无法获取用户角色",
			})
			c.Abort()
			return
		}

		// 检查角色权限
		userRoleStr := userRole.(string)
		for _, role := range roles {
			if userRoleStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
		})
		c.Abort()
	}
}

// RequireAdmin 需要管理员权限的中间件
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

// OptionalAuth 可选认证中间件
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		// 提取令牌
		tokenString := authHeader[7:]

		// 验证令牌
		claims, err := m.userService.ValidateJWT(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// APIKeyAuth API密钥认证中间件
func (m *AuthMiddleware) APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
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
				"error": "API密钥是必需的",
			})
			c.Abort()
			return
		}

		// 将API密钥存储到上下文中
		c.Set("api_key", apiKey)

		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, API_KEY")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
