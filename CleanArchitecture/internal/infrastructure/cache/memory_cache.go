package cache

import (
	"context"
	"strings"
	"sync"
	"time"
	
	"clean-url-shortener/pkg/cache"
	"clean-url-shortener/pkg/logger"
)

// 📚 ИЗУЧИТЬ: In-Memory кэширование
// https://en.wikipedia.org/wiki/Cache_(computing)
// https://github.com/patrickmn/go-cache
// https://github.com/allegro/bigcache
//
// КОГДА ИСПОЛЬЗОВАТЬ IN-MEMORY CACHE:
// 1. Разработка и тестирование
// 2. Небольшие приложения (single instance)
// 3. Кэширование вычислений (не данных между запросами)
// 4. Временное решение перед Redis
//
// ОГРАНИЧЕНИЯ IN-MEMORY CACHE:
// 1. Не работает в кластере (каждый экземпляр имеет свой кэш)
// 2. Данные теряются при перезапуске
// 3. Потребляет память приложения
// 4. Нет персистентности

// cacheItem элемент кэша с метаданными
type cacheItem struct {
	data      []byte    // Закэшированные данные
	expiresAt time.Time // Время истечения
	createdAt time.Time // Время создания
}

// isExpired проверяет истек ли элемент кэша
func (item *cacheItem) isExpired() bool {
	return time.Now().After(item.expiresAt)
}

// MemoryCache реализация кэша в памяти
// 📚 ИЗУЧИТЬ: Thread-safe programming в Go
// https://golang.org/doc/effective_go.html#sharing
// https://pkg.go.dev/sync
type MemoryCache struct {
	// items хранилище элементов кэша
	items map[string]*cacheItem
	
	// mutex для thread-safe операций
	// ВАЖНО: В Go нужно защищать доступ к shared state
	mutex sync.RWMutex
	
	// logger для логирования операций кэша
	logger logger.Logger
	
	// cleanupInterval интервал очистки истекших элементов
	cleanupInterval time.Duration
	
	// stopCleanup канал для остановки фоновой очистки
	stopCleanup chan struct{}
	
	// metrics метрики кэша
	metrics *cacheMetrics
}

// cacheMetrics внутренние метрики кэша
type cacheMetrics struct {
	hits   int64
	misses int64
	errors int64
	mutex  sync.RWMutex
}

// NewMemoryCache создает новый in-memory кэш
func NewMemoryCache(logger logger.Logger, cleanupInterval time.Duration) cache.Cache {
	mc := &MemoryCache{
		items:           make(map[string]*cacheItem),
		logger:          logger.InfrastructureLogger().With("component", "memory_cache"),
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
		metrics:         &cacheMetrics{},
	}
	
	// Запускаем фоновую очистку истекших элементов
	go mc.startCleanup()
	
	mc.logger.Info("Memory cache initialized", 
		"cleanup_interval", cleanupInterval)
	
	return mc
}

// Get получает значение по ключу
func (mc *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	item, exists := mc.items[key]
	if !exists {
		mc.recordMiss()
		mc.logger.Debug("Cache miss", "key", key)
		return nil, cache.ErrCacheKeyNotFound
	}
	
	// Проверяем истечение TTL
	if item.isExpired() {
		mc.recordMiss()
		mc.logger.Debug("Cache expired", "key", key, "expired_at", item.expiresAt)
		
		// Удаляем истекший элемент (требует write lock)
		mc.mutex.RUnlock()
		mc.mutex.Lock()
		delete(mc.items, key)
		mc.mutex.Unlock()
		mc.mutex.RLock()
		
		return nil, cache.ErrCacheKeyNotFound
	}
	
	mc.recordHit()
	mc.logger.Debug("Cache hit", "key", key, "data_size", len(item.data))
	
	// Возвращаем копию данных для безопасности
	result := make([]byte, len(item.data))
	copy(result, item.data)
	
	return result, nil
}

// Set сохраняет значение с TTL
func (mc *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		return cache.ErrInvalidTTL
	}
	
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	// Создаем копию данных для безопасности
	data := make([]byte, len(value))
	copy(data, value)
	
	item := &cacheItem{
		data:      data,
		expiresAt: time.Now().Add(ttl),
		createdAt: time.Now(),
	}
	
	mc.items[key] = item
	
	mc.logger.Debug("Cache set", 
		"key", key, 
		"data_size", len(value), 
		"ttl", ttl,
		"expires_at", item.expiresAt)
	
	return nil
}

// Delete удаляет значение по ключу
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	_, exists := mc.items[key]
	if exists {
		delete(mc.items, key)
		mc.logger.Debug("Cache delete", "key", key)
	} else {
		mc.logger.Debug("Cache delete - key not found", "key", key)
	}
	
	return nil
}

// Exists проверяет существование ключа
func (mc *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	item, exists := mc.items[key]
	if !exists {
		return false, nil
	}
	
	// Проверяем истечение TTL
	if item.isExpired() {
		// Удаляем истекший элемент
		mc.mutex.RUnlock()
		mc.mutex.Lock()
		delete(mc.items, key)
		mc.mutex.Unlock()
		mc.mutex.RLock()
		
		return false, nil
	}
	
	return true, nil
}

// SetNX устанавливает значение только если ключ не существует
func (mc *MemoryCache) SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error) {
	if ttl <= 0 {
		return false, cache.ErrInvalidTTL
	}
	
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	// Проверяем существование ключа
	item, exists := mc.items[key]
	if exists && !item.isExpired() {
		mc.logger.Debug("Cache SetNX - key already exists", "key", key)
		return false, nil // Ключ уже существует
	}
	
	// Создаем новый элемент
	data := make([]byte, len(value))
	copy(data, value)
	
	newItem := &cacheItem{
		data:      data,
		expiresAt: time.Now().Add(ttl),
		createdAt: time.Now(),
	}
	
	mc.items[key] = newItem
	
	mc.logger.Debug("Cache SetNX success", 
		"key", key, 
		"data_size", len(value), 
		"ttl", ttl)
	
	return true, nil
}

// Expire устанавливает TTL для существующего ключа
func (mc *MemoryCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if ttl <= 0 {
		return cache.ErrInvalidTTL
	}
	
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	item, exists := mc.items[key]
	if !exists {
		mc.logger.Debug("Cache expire - key not found", "key", key)
		return cache.ErrCacheKeyNotFound
	}
	
	if item.isExpired() {
		delete(mc.items, key)
		return cache.ErrCacheKeyNotFound
	}
	
	// Обновляем время истечения
	item.expiresAt = time.Now().Add(ttl)
	
	mc.logger.Debug("Cache expire updated", 
		"key", key, 
		"new_ttl", ttl,
		"expires_at", item.expiresAt)
	
	return nil
}

// Keys возвращает ключи по паттерну
// ВНИМАНИЕ: Эта операция может быть медленной для больших кэшей!
func (mc *MemoryCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	var matchingKeys []string
	
	for key, item := range mc.items {
		// Пропускаем истекшие элементы
		if item.isExpired() {
			continue
		}
		
		// Простая проверка паттерна (поддерживает только * в конце)
		if matchesPattern(key, pattern) {
			matchingKeys = append(matchingKeys, key)
		}
	}
	
	mc.logger.Debug("Cache keys search", 
		"pattern", pattern, 
		"found_count", len(matchingKeys))
	
	return matchingKeys, nil
}

// FlushAll очищает весь кэш
// ВНИМАНИЕ: Используйте только для тестов!
func (mc *MemoryCache) FlushAll(ctx context.Context) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	itemCount := len(mc.items)
	mc.items = make(map[string]*cacheItem)
	
	// Сбрасываем метрики
	mc.metrics.mutex.Lock()
	mc.metrics.hits = 0
	mc.metrics.misses = 0
	mc.metrics.errors = 0
	mc.metrics.mutex.Unlock()
	
	mc.logger.Warn("Cache flushed", "deleted_items", itemCount)
	
	return nil
}

// ФОНОВАЯ ОЧИСТКА ИСТЕКШИХ ЭЛЕМЕНТОВ

// startCleanup запускает фоновую очистку истекших элементов
func (mc *MemoryCache) startCleanup() {
	ticker := time.NewTicker(mc.cleanupInterval)
	defer ticker.Stop()
	
	mc.logger.Info("Cache cleanup started", "interval", mc.cleanupInterval)
	
	for {
		select {
		case <-ticker.C:
			mc.cleanupExpired()
		case <-mc.stopCleanup:
			mc.logger.Info("Cache cleanup stopped")
			return
		}
	}
}

// cleanupExpired удаляет истекшие элементы
func (mc *MemoryCache) cleanupExpired() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	now := time.Now()
	deletedCount := 0
	
	for key, item := range mc.items {
		if now.After(item.expiresAt) {
			delete(mc.items, key)
			deletedCount++
		}
	}
	
	if deletedCount > 0 {
		mc.logger.Debug("Cache cleanup completed", 
			"deleted_items", deletedCount,
			"remaining_items", len(mc.items))
	}
}

// Stop останавливает фоновую очистку
func (mc *MemoryCache) Stop() {
	close(mc.stopCleanup)
}

// МЕТРИКИ КЭША

// recordHit записывает попадание в кэш
func (mc *MemoryCache) recordHit() {
	mc.metrics.mutex.Lock()
	mc.metrics.hits++
	mc.metrics.mutex.Unlock()
}

// recordMiss записывает промах кэша
func (mc *MemoryCache) recordMiss() {
	mc.metrics.mutex.Lock()
	mc.metrics.misses++
	mc.metrics.mutex.Unlock()
}

// recordError записывает ошибку кэша
func (mc *MemoryCache) recordError() {
	mc.metrics.mutex.Lock()
	mc.metrics.errors++
	mc.metrics.mutex.Unlock()
}

// GetMetrics возвращает метрики кэша
func (mc *MemoryCache) GetMetrics() cache.CacheMetrics {
	mc.metrics.mutex.RLock()
	defer mc.metrics.mutex.RUnlock()
	
	return cache.CacheMetrics{
		Hits:   mc.metrics.hits,
		Misses: mc.metrics.misses,
		Errors: mc.metrics.errors,
	}
}

// GetStats возвращает расширенную статистику кэша
func (mc *MemoryCache) GetStats() map[string]interface{} {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	mc.metrics.mutex.RLock()
	defer mc.metrics.mutex.RUnlock()
	
	totalMemory := 0
	expiredItems := 0
	now := time.Now()
	
	for _, item := range mc.items {
		totalMemory += len(item.data)
		if now.After(item.expiresAt) {
			expiredItems++
		}
	}
	
	metrics := mc.GetMetrics()
	
	return map[string]interface{}{
		"total_items":    len(mc.items),
		"expired_items":  expiredItems,
		"memory_bytes":   totalMemory,
		"hits":           metrics.Hits,
		"misses":         metrics.Misses,
		"errors":         metrics.Errors,
		"hit_ratio":      metrics.HitRatio(),
	}
}

// UTILITY FUNCTIONS

// matchesPattern проверяет соответствие ключа паттерну
func matchesPattern(key, pattern string) bool {
	// Простая реализация: поддерживаем только * в конце
	if pattern == "*" {
		return true
	}
	
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(key, prefix)
	}
	
	return key == pattern
}

// КОНФИГУРАЦИЯ ДЛЯ PRODUCTION

// MemoryCacheConfig конфигурация memory cache
type MemoryCacheConfig struct {
	CleanupInterval time.Duration `yaml:"cleanup_interval" env:"MEMORY_CACHE_CLEANUP_INTERVAL" default:"5m"`
	MaxItems        int           `yaml:"max_items" env:"MEMORY_CACHE_MAX_ITEMS" default:"10000"`
	MaxMemoryMB     int           `yaml:"max_memory_mb" env:"MEMORY_CACHE_MAX_MEMORY_MB" default:"100"`
}

// ПРИМЕЧАНИЯ ПО PRODUCTION:
// 1. Для production используйте Redis вместо memory cache
// 2. Memory cache подходит только для single-instance приложений
// 3. Мониторьте потребление памяти
// 4. Настройте подходящие TTL для ваших данных
// 5. Рассмотрите использование LRU eviction policy для ограничения памяти

// ПРИМЕР ИСПОЛЬЗОВАНИЯ:

/*
// В main.go или DI container:
func initCache(logger logger.Logger) cache.Cache {
    // Для development
    if env == "development" {
        return cache.NewMemoryCache(logger, 5*time.Minute)
    }
    
    // Для production
    return redis.NewRedisCache(redisConfig, logger)
}

// В Use Case:
type LoginUseCase struct {
    userRepo user.Repository // CachingUserRepository
    // кэш инкапсулирован в репозитории
}

// В Repository:
userCache := cache.NewUserCache(cacheInstance)
cachingRepo := cache.NewCachingUserRepository(baseRepo, userCache, 10*time.Minute)
*/