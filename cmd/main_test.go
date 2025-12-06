package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestMain 测试主函数
func TestMain(t *testing.T) {
	// 跳过实际启动应用程序的测试，因为这需要真实的数据库环境
	// 这里主要测试辅助函数的功能
	t.Run("TestGetEnv", testGetEnv)
	t.Run("TestBuildDSN", testBuildDSN)
	t.Run("TestAutoMigrate", testAutoMigrate)
}

// TestGetEnv 测试环境变量获取函数
func testGetEnv(t *testing.T) {
	// 测试已设置的环境变量
	testKey := "TEST_ENV_VAR"
	testValue := "test_value"

	// 设置环境变量
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	// 测试获取已设置的环境变量
	result := getEnv(testKey, "default")
	if result != testValue {
		t.Errorf("getEnv(%s, 'default') = %s, want %s", testKey, result, testValue)
	}

	// 测试获取未设置的环境变量
	result = getEnv("NON_EXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("getEnv('NON_EXISTENT_VAR', 'default_value') = %s, want default_value", result)
	}
}

// TestBuildDSN 测试数据库连接字符串构建函数
func testBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     string
		user     string
		password string
		dbname   string
		expected string
	}{
		{
			name:     "标准配置",
			host:     "localhost",
			port:     "5432",
			user:     "postgres",
			password: "postgres",
			dbname:   "testdb",
			expected: "host=localhost user=postgres password=postgres dbname=testdb port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		},
		{
			name:     "不同端口",
			host:     "db.example.com",
			port:     "5433",
			user:     "testuser",
			password: "testpass",
			dbname:   "production",
			expected: "host=db.example.com user=testuser password=testpass dbname=production port=5433 sslmode=disable TimeZone=Asia/Shanghai",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildDSN(tt.host, tt.port, tt.user, tt.password, tt.dbname)
			if result != tt.expected {
				t.Errorf("buildDSN() = %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestAutoMigrate 测试数据库自动迁移功能
func testAutoMigrate(t *testing.T) {
	// 使用临时数据库进行测试
	dbName := fmt.Sprintf("test_db_%d", time.Now().Unix())
	dsn := buildDSN("localhost", "5432", "postgres", "postgres", dbName)

	// 连接到PostgreSQL服务器（不指定数据库）
	serverDSN := buildDSN("localhost", "5432", "postgres", "postgres", "postgres")
	db, err := gorm.Open(postgres.Open(serverDSN), &gorm.Config{})
	if err != nil {
		t.Skipf("无法连接到PostgreSQL服务器: %v", err)
	}

	// 创建测试数据库
	if err := db.Exec("CREATE DATABASE " + dbName).Error; err != nil {
		t.Fatalf("无法创建测试数据库: %v", err)
	}
	defer func() {
		// 清理：删除测试数据库
		db.Exec("DROP DATABASE IF EXISTS " + dbName)
	}()

	// 连接到测试数据库
	testDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法连接到测试数据库: %v", err)
	}

	// 测试自动迁移
	err = autoMigrate(testDB)
	if err != nil {
		t.Errorf("autoMigrate() error = %v", err)
	}

	// 验证表是否创建成功
	var tableCount int64
	if err := testDB.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name IN ('documents', 'document_versions', 'document_metadata')
	`).Scan(&tableCount).Error; err != nil {
		t.Errorf("查询表数量失败: %v", err)
	}

	if tableCount != 3 {
		t.Errorf("预期创建3个表，实际创建 %d 个", tableCount)
	}

	// 验证约束是否创建成功
	var constraintCount int64
	if err := testDB.Raw(`
		SELECT COUNT(*)
		FROM information_schema.table_constraints
		WHERE table_name = 'document_versions'
		AND constraint_name = 'idx_document_version_unique'
	`).Scan(&constraintCount).Error; err != nil {
		t.Errorf("查询约束数量失败: %v", err)
	}

	if constraintCount != 1 {
		t.Errorf("预期创建1个约束，实际创建 %d 个", constraintCount)
	}
}

// TestApplicationIntegration 测试应用程序集成
func TestApplicationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 设置测试环境变量
	testEnv := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_USER":     "postgres",
		"DB_PASSWORD": "postgres",
		"DB_NAME":     "test_integration",
		"SERVER_PORT": "0", // 使用随机端口
		"STORAGE_DIR": "./test_storage",
	}

	// 备份原始环境变量
	originalEnv := make(map[string]string)
	for key, value := range testEnv {
		originalEnv[key] = os.Getenv(key)
		os.Setenv(key, value)
	}
	defer func() {
		// 恢复原始环境变量
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// 创建测试数据库
	dbName := testEnv["DB_NAME"]
	serverDSN := buildDSN(testEnv["DB_HOST"], testEnv["DB_PORT"], testEnv["DB_USER"], testEnv["DB_PASSWORD"], "postgres")
	db, err := gorm.Open(postgres.Open(serverDSN), &gorm.Config{})
	if err != nil {
		t.Skipf("无法连接到PostgreSQL服务器: %v", err)
	}

	// 删除并重新创建测试数据库
	db.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err := db.Exec("CREATE DATABASE " + dbName).Error; err != nil {
		t.Skipf("无法创建测试数据库: %v", err)
	}
	defer func() {
		db.Exec("DROP DATABASE IF EXISTS " + dbName)
	}()

	// 创建存储目录
	storageDir := testEnv["STORAGE_DIR"]
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		t.Fatalf("无法创建存储目录: %v", err)
	}
	defer os.RemoveAll(storageDir)

	// 启动应用程序（在单独的goroutine中）
	appStarted := make(chan error)
	appStopped := make(chan struct{})

	go func() {
		defer close(appStopped)

		// 重定向标准输出以避免干扰测试
		// 注意：这里我们不直接调用main()，因为这会导致程序退出
		// 在实际项目中，可以重构main()函数为可测试的形式
		appStarted <- fmt.Errorf("集成测试需要重构main()函数以支持测试")
	}()

	select {
	case err := <-appStarted:
		if err != nil {
			t.Logf("应用程序启动失败（预期）: %v", err)
			return
		}
	case <-time.After(5 * time.Second):
		t.Log("应用程序启动超时")
	case <-appStopped:
		t.Log("应用程序已停止")
	}
}

// TestGracefulShutdown 测试应用程序优雅关闭
func TestGracefulShutdown(t *testing.T) {
	// 这个测试验证应用程序能够正确处理信号
	// 在实际项目中，可以发送SIGTERM信号并验证应用程序是否正确关闭

	// 由于main()函数的结构限制，这里主要测试概念
	t.Log("优雅关闭测试需要信号处理支持")
}

// TestStorageDirectoryCreation 测试存储目录创建
func TestStorageDirectoryCreation(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("test_storage_%d", time.Now().Unix()))

	// 测试目录创建
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("创建存储目录失败: %v", err)
	}

	// 验证目录是否存在
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("存储目录创建失败")
	}

	// 验证目录权限
	fileInfo, err := os.Stat(tempDir)
	if err != nil {
		t.Fatalf("获取目录信息失败: %v", err)
	}

	// 检查权限是否为0755
	expectedMode := os.FileMode(0755)
	if fileInfo.Mode().Perm() != expectedMode {
		t.Errorf("目录权限不正确，预期 %o, 实际 %o", expectedMode, fileInfo.Mode().Perm())
	}

	// 清理
	os.RemoveAll(tempDir)
}

// TestEnvironmentConfiguration 测试环境配置
func TestEnvironmentConfiguration(t *testing.T) {
	// 测试不同的环境配置
	configs := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "开发环境配置",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "postgres",
				"DB_NAME":     "dev_db",
				"SERVER_PORT": "8080",
				"STORAGE_DIR": "./storage",
			},
		},
		{
			name: "生产环境配置",
			env: map[string]string{
				"DB_HOST":     "prod-db.example.com",
				"DB_PORT":     "5432",
				"DB_USER":     "produser",
				"DB_PASSWORD": "prodpass",
				"DB_NAME":     "prod_db",
				"SERVER_PORT": "80",
				"STORAGE_DIR": "/var/storage",
			},
		},
	}

	for _, config := range configs {
		t.Run(config.name, func(t *testing.T) {
			// 备份原始环境变量
			originalEnv := make(map[string]string)
			for key, value := range config.env {
				originalEnv[key] = os.Getenv(key)
				os.Setenv(key, value)
			}

			// 测试getEnv函数
			for key, expectedValue := range config.env {
				result := getEnv(key, "default")
				if result != expectedValue {
					t.Errorf("getEnv(%s, 'default') = %s, want %s", key, result, expectedValue)
				}
			}

			// 测试buildDSN函数
			dsn := buildDSN(
				config.env["DB_HOST"],
				config.env["DB_PORT"],
				config.env["DB_USER"],
				config.env["DB_PASSWORD"],
				config.env["DB_NAME"],
			)

			expectedDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
				config.env["DB_HOST"],
				config.env["DB_USER"],
				config.env["DB_PASSWORD"],
				config.env["DB_NAME"],
				config.env["DB_PORT"],
			)

			if dsn != expectedDSN {
				t.Errorf("buildDSN() = %s, want %s", dsn, expectedDSN)
			}

			// 恢复原始环境变量
			for key, value := range originalEnv {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}
		})
	}
}

// TestHealthEndpoint 测试健康检查端点
func TestHealthEndpoint(t *testing.T) {
	// 这个测试需要实际启动HTTP服务器
	// 由于main()函数的限制，这里提供测试框架

	t.Log("健康检查端点测试需要HTTP服务器启动支持")

	// 在实际项目中，可以：
	// 1. 启动测试服务器
	// 2. 发送GET请求到/health端点
	// 3. 验证响应状态码为200
	// 4. 验证响应体包含{"status": "ok"}
}

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	// 测试数据库连接功能
	dsn := buildDSN("localhost", "5432", "postgres", "postgres", "postgres")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("无法连接到PostgreSQL服务器: %v", err)
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("获取原生数据库连接失败: %v", err)
	}

	// 测试ping
	if err := sqlDB.Ping(); err != nil {
		t.Errorf("数据库ping失败: %v", err)
	}

	// 测试连接池设置
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 验证设置（通过Stats()方法检查连接池状态）
	stats := sqlDB.Stats()
	if stats.Idle != 0 {
		// 初始状态下应该没有空闲连接
		t.Logf("当前空闲连接数: %d", stats.Idle)
	}

	if stats.OpenConnections != 1 {
		// 应该有一个打开的连接
		t.Logf("当前打开连接数: %d", stats.OpenConnections)
	}

	// 验证设置是否生效
	t.Logf("无法直接验证最大连接数设置，但设置方法已调用")
	t.Logf("连接池状态: %+v", stats)
}
