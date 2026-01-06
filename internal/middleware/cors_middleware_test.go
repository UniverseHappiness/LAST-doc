package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORS_Headers(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由并应用CORS中间件
	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建请求
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证CORS头
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, API_KEY",
		"Access-Control-Max-Age":       "86400",
	}

	for key, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(key)
		if actualValue != expectedValue {
			t.Errorf("Header %s = %s, want %s", key, actualValue, expectedValue)
		}
	}
}

func TestCORS_OptionsRequest(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	handlerCalled := false

	// 创建路由并应用CORS中间件
	router := gin.New()
	router.Use(CORS())
	router.OPTIONS("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "should not be called"})
	})

	// 创建OPTIONS预检请求
	req, err := http.NewRequest("OPTIONS", "/test", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// OPTIONS请求应该在中间件中返回204，不应该调用handler
	if handlerCalled {
		t.Error("OPTIONS请求应该在中间件中终止，不应该调用handler")
	}

	// 验证状态码为204
	if w.Code != http.StatusNoContent {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestCORS_DifferentMethods(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			// 创建路由并应用CORS中间件
			router := gin.New()
			router.Use(CORS())
			router.Any("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			// 创建请求
			req, err := http.NewRequest(method, "/test", nil)
			if err != nil {
				t.Fatalf("创建请求失败: %v", err)
			}

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证状态码为200
			if w.Code != http.StatusOK {
				t.Errorf("状态码 = %d, want %d", w.Code, http.StatusOK)
			}

			// 验证CORS头存在
			origin := w.Header().Get("Access-Control-Allow-Origin")
			if origin != "*" {
				t.Errorf("Access-Control-Allow-Origin = %s, want *", origin)
			}

			allowMethods := w.Header().Get("Access-Control-Allow-Methods")
			expectedMethods := "GET, POST, PUT, DELETE, OPTIONS"
			if allowMethods != expectedMethods {
				t.Errorf("Access-Control-Allow-Methods = %s, want %s", allowMethods, expectedMethods)
			}
		})
	}
}

func TestCORS_AllowsRequestToPassThrough(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由并应用CORS中间件
	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// 创建请求
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证调用handler并返回200
	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusOK)
	}

	// 验证响应体
	body := w.Body.String()
	if body != `{"status":"success"}` && body != `{"status":"success"}`+"\n" {
		t.Errorf("响应体 = %s, want success", body)
	}
}
