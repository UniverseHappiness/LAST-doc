package service

import (
	"sync"
	"time"
)

// CacheItem 缓存项
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired 检查缓存项是否过期
func (c *CacheItem) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// CacheService 缓存服务接口
type CacheService interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
	Clear() error
}

// memoryCache 内存缓存实现
type memoryCache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() CacheService {
	return &memoryCache{
		items: make(map[string]CacheItem),
	}
}

// Set 设置缓存
func (m *memoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.items[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	return nil
}

// Get 获取缓存
func (m *memoryCache) Get(key string) (interface{}, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, found := m.items[key]
	if !found {
		return nil, false
	}

	// 检查是否过期
	if item.IsExpired() {
		return nil, false
	}

	return item.Value, true
}

// Delete 删除缓存
func (m *memoryCache) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.items, key)
	return nil
}

// Clear 清空缓存
func (m *memoryCache) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.items = make(map[string]CacheItem)
	return nil
}

// searchCacheKey 生成搜索缓存键
func searchCacheKey(query, searchType string, filters map[string]interface{}, page, size int) string {
	// 简单的键生成，实际项目中可以使用更复杂的哈希算法
	return query + "|" + searchType + "|" + string(page) + "|" + string(size)
}
