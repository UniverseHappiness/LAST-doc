package service

import (
	"testing"
	"time"
)

// TestCacheItem_IsExpired 测试缓存项过期检查
func TestCacheItem_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		item     CacheItem
		expected bool
	}{
		{
			name: "已过期项",
			item: CacheItem{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "未过期项",
			item: CacheItem{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "刚好过期",
			item: CacheItem{
				ExpiresAt: time.Now().Add(-time.Second),
			},
			expected: true,
		},
		{
			name: "即将过期",
			item: CacheItem{
				ExpiresAt: time.Now().Add(time.Second),
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsExpired()
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestNewMemoryCache 测试创建内存缓存实例
func TestNewMemoryCache(t *testing.T) {
	cache := NewMemoryCache()
	if cache == nil {
		t.Fatal("NewMemoryCache() 返回 nil")
	}

	// 验证返回的类型是否正确
	if _, ok := cache.(*memoryCache); !ok {
		t.Error("NewMemoryCache() 返回的类型不正确")
	}
}

// TestMemoryCache_SetAndGet 测试设置和获取缓存
func TestMemoryCache_SetAndGet(t *testing.T) {
	cache := NewMemoryCache()

	tests := []struct {
		name  string
		key   string
		value interface{}
		ttl   time.Duration
	}{
		{
			name:  "字符串值",
			key:   "key1",
			value: "test value",
			ttl:   time.Minute,
		},
		{
			name:  "整数值",
			key:   "key2",
			value: 12345,
			ttl:   time.Minute * 30,
		},
		{
			name:  "nil值",
			key:   "key5",
			value: nil,
			ttl:   time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置缓存
			err := cache.Set(tt.key, tt.value, tt.ttl)
			if err != nil {
				t.Errorf("Set() error = %v", err)
				return
			}

			// 获取缓存
			value, found := cache.Get(tt.key)
			if !found {
				t.Error("Get() 未找到已设置的缓存")
				return
			}

			// 验证值
			if value != tt.value {
				t.Errorf("Get() value = %v, expected %v", value, tt.value)
			}
		})
	}
}

// TestMemoryCache_NotFound 测试获取不存在的缓存
func TestMemoryCache_NotFound(t *testing.T) {
	cache := NewMemoryCache()

	value, found := cache.Get("non_existent_key")
	if found {
		t.Error("Get() 找到了不应该存在的缓存")
	}
	if value != nil {
		t.Errorf("Get() value = %v, expected nil", value)
	}
}

// TestMemoryCache_Expired 测试过期缓存
func TestMemoryCache_Expired(t *testing.T) {
	cache := NewMemoryCache()

	key := "expired_key"
	value := "test"

	// 设置一个已经过期的缓存
	err := cache.Set(key, value, -time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 获取已过期的缓存
	retrievedValue, found := cache.Get(key)
	if found {
		t.Error("Get() 找到了已过期的缓存")
	}
	if retrievedValue != nil {
		t.Errorf("Get() value = %v, expected nil", retrievedValue)
	}
}

// TestMemoryCache_Delete 测试删除缓存
func TestMemoryCache_Delete(t *testing.T) {
	cache := NewMemoryCache()

	key := "delete_key"
	value := "test"

	// 设置缓存
	err := cache.Set(key, value, time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 删除缓存
	err = cache.Delete(key)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// 验证缓存已被删除
	_, found := cache.Get(key)
	if found {
		t.Error("Delete() 后缓存仍然存在")
	}
}

// TestMemoryCache_DeleteNonExistent 测试删除不存在的缓存
func TestMemoryCache_DeleteNonExistent(t *testing.T) {
	cache := NewMemoryCache()

	err := cache.Delete("non_existent_key")
	if err != nil {
		t.Errorf("Delete() error = %v, expected nil", err)
	}
}

// TestMemoryCache_Clear 测试清空缓存
func TestMemoryCache_Clear(t *testing.T) {
	cache := NewMemoryCache()

	// 设置多个缓存项
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	cache.Set("key3", "value3", time.Minute)

	// 清空缓存
	err := cache.Clear()
	if err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	// 验证所有缓存项都被清除
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")
	_, found3 := cache.Get("key3")

	if found1 || found2 || found3 {
		t.Error("Clear() 后仍有缓存项存在")
	}
}

// TestMemoryCache_ClearEmpty 测试清空已空的缓存
func TestMemoryCache_ClearEmpty(t *testing.T) {
	cache := NewMemoryCache()

	err := cache.Clear()
	if err != nil {
		t.Errorf("Clear() error = %v, expected nil", err)
	}
}

// TestMemoryCache_Overwrite 测试覆盖已存在的缓存
func TestMemoryCache_Overwrite(t *testing.T) {
	cache := NewMemoryCache()

	key := "overwrite_key"

	// 设置初始值
	err := cache.Set(key, "value1", time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 覆盖为新值
	err = cache.Set(key, "value2", time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 验证新值
	value, found := cache.Get(key)
	if !found {
		t.Fatal("Get() 未找到缓存")
	}

	if value != "value2" {
		t.Errorf("Get() value = %v, expected value2", value)
	}
}

// TestMemoryCache_SameKeyDifferentTTL 测试相同键不同TTL
func TestMemoryCache_SameKeyDifferentTTL(t *testing.T) {
	cache := NewMemoryCache()

	key := "ttl_test_key"

	// 设置短TTL
	cache.Set(key, "value1", time.Millisecond*100)

	// 立即覆盖为长TTL
	cache.Set(key, "value2", time.Hour)

	// 短暂等待
	time.Sleep(time.Millisecond * 50)

	// 验证缓存仍然存在（因为使用长TTL）
	value, found := cache.Get(key)
	if !found {
		t.Fatal("Get() 未找到缓存，可能被错误地过期")
	}

	if value != "value2" {
		t.Errorf("Get() value = %v, expected value2", value)
	}
}

// TestSearchCacheKey 测试生成搜索缓存键
func TestSearchCacheKey(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		searchType    string
		filters       map[string]interface{}
		page          int
		size          int
		expectInvalid bool
	}{
		{
			name:       "标准搜索",
			query:      "test query",
			searchType: "keyword",
			page:       1,
			size:       10,
		},
		{
			name:       "语义搜索",
			query:      "semantic search",
			searchType: "semantic",
			page:       2,
			size:       20,
		},
		{
			name:       "混合搜索",
			query:      "hybrid search",
			searchType: "hybrid",
			page:       1,
			size:       50,
		},
		{
			name:       "空查询",
			query:      "",
			searchType: "keyword",
			page:       1,
			size:       10,
		},
		{
			name:       "第一页",
			query:      "first",
			searchType: "keyword",
			page:       1,
			size:       100,
		},
		{
			name:       "大页码",
			query:      "test",
			searchType: "keyword",
			page:       100,
			size:       10,
		},
		{
			name:       "特殊字符查询",
			query:      "test:query",
			searchType: "keyword",
			page:       1,
			size:       10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := searchCacheKey(tt.query, tt.searchType, tt.filters, tt.page, tt.size)

			// 验证key包含查询
			if !containsSubstrings(key, tt.query) {
				t.Errorf("key应该包含query: query=%s, key=%s", tt.query, key)
			}

			// 验证key包含搜索类型
			if !containsSubstrings(key, tt.searchType) {
				t.Errorf("key应该包含searchType: type=%s, key=%s", tt.searchType, key)
			}

			// 验证key不为空
			if key == "" {
				t.Error("生成的cache key为空")
			}
		})
	}
}

// containsSubstrings 辅助函数：检查字符串是否包含子串
func containsSubstrings(s, substr string) bool {
	if substr == "" {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestSearchCacheKey_Uniqueness 测试不同参数生成不同key
func TestSearchCacheKey_Uniqueness(t *testing.T) {
	tests := []struct {
		name   string
		query1 string
		query2 string
		type1  string
		type2  string
		page1  int
		page2  int
		size1  int
		size2  int
	}{
		{
			name:   "不同查询",
			query1: "query1",
			query2: "query2",
			type1:  "keyword",
			type2:  "keyword",
			page1:  1,
			page2:  1,
			size1:  10,
			size2:  10,
		},
		{
			name:   "不同类型",
			query1: "same query",
			query2: "same query",
			type1:  "keyword",
			type2:  "semantic",
			page1:  1,
			page2:  1,
			size1:  10,
			size2:  10,
		},
		{
			name:   "不同页码",
			query1: "query",
			query2: "query",
			type1:  "keyword",
			type2:  "keyword",
			page1:  1,
			page2:  2,
			size1:  10,
			size2:  10,
		},
		{
			name:   "不同大小",
			query1: "query",
			query2: "query",
			type1:  "keyword",
			type2:  "keyword",
			page1:  1,
			page2:  1,
			size1:  10,
			size2:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key1 := searchCacheKey(tt.query1, tt.type1, nil, tt.page1, tt.size1)
			key2 := searchCacheKey(tt.query2, tt.type2, nil, tt.page2, tt.size2)

			if key1 == key2 {
				t.Error("不同参数应该生成不同的key")
			}
		})
	}
}

// TestMemoryCache_Concurrent 测试并发访问
func TestMemoryCache_Concurrent(t *testing.T) {
	cache := NewMemoryCache()
	done := make(chan bool)

	// 并发写入
	for i := 0; i < 100; i++ {
		go func(n int) {
			key := "concurrent_key_" + string(rune(n))
			cache.Set(key, n, time.Minute)
			done <- true
		}(i)
	}

	// 等待所有写入完成
	for i := 0; i < 100; i++ {
		<-done
	}
}

// TestMemoryCache_EmptyStringKey 测试空字符串键
func TestMemoryCache_EmptyStringKey(t *testing.T) {
	cache := NewMemoryCache()

	value := "test"
	err := cache.Set("", value, time.Minute)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	retrievedValue, found := cache.Get("")
	if !found {
		t.Fatal("Get() 未找到空键的缓存")
	}

	if retrievedValue != value {
		t.Errorf("Get() value = %v, expected %v", retrievedValue, value)
	}
}
