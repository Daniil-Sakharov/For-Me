package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// 📚 ИЗУЧИТЬ: Caching patterns и стратегии
// https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Strategies.html
// https://redis.io/docs/manual/patterns/
// https://martinfowler.com/bliki/TwoHardThings.html
//
// ОСНОВНЫЕ ПАТТЕРНЫ КЭШИРОВАНИЯ:
// 1. CACHE-ASIDE (Lazy Loading) - приложение управляет кэшем
// 2. WRITE-THROUGH - запись в кэш и БД одновременно  
// 3. WRITE-BEHIND (Write Back) - отложенная запись в БД
// 4. REFRESH-AHEAD - проактивное обновление кэша
//
// ПРИНЦИПЫ КЭШИРОВАНИЯ В CLEAN ARCHITECTURE:
// 1. Кэш - это ДЕТАЛЬ РЕАЛИЗАЦИИ (infrastructure)
// 2. Бизнес-логика НЕ ДОЛЖНА зависеть от кэша
// 3. Кэш должен быть ПРОЗРАЧНЫМ для Use Cases
// 4. Используем интерфейсы для тестируемости

// Cache определяет интерфейс для кэширования в Clean Architecture
// ВАЖНО: Интерфейс находится в pkg/, потому что используется разными слоями
type Cache interface {
	// Get получает значение по ключу
	Get(ctx context.Context, key string) ([]byte, error)
	
	// Set сохраняет значение с TTL (Time To Live)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	
	// Delete удаляет значение по ключу
	Delete(ctx context.Context, key string) error
	
	// Exists проверяет существование ключа
	Exists(ctx context.Context, key string) (bool, error)
	
	// SetNX устанавливает значение только если ключ не существует (для блокировок)
	SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error)
	
	// Expire устанавливает TTL для существующего ключа
	Expire(ctx context.Context, key string, ttl time.Duration) error
	
	// Keys возвращает ключи по паттерну (осторожно в production!)
	Keys(ctx context.Context, pattern string) ([]string, error)
	
	// FlushAll очищает весь кэш (только для тестов!)
	FlushAll(ctx context.Context) error
}

// Стандартные ошибки кэша
var (
	ErrCacheKeyNotFound = errors.New("cache key not found")
	ErrCacheUnavailable = errors.New("cache service unavailable")
	ErrInvalidTTL       = errors.New("invalid TTL value")
)

// ТИПИЗИРОВАННЫЙ КЭШИРУЮЩИЙ СЛОЙ
// 📚 ИЗУЧИТЬ: Generic programming в Go
// https://go.dev/doc/tutorial/generics

// TypedCache предоставляет типизированный интерфейс для кэширования
type TypedCache[T any] struct {
	cache  Cache
	prefix string // Префикс для ключей (namespace)
}

// NewTypedCache создает типизированный кэш
func NewTypedCache[T any](cache Cache, prefix string) *TypedCache[T] {
	return &TypedCache[T]{
		cache:  cache,
		prefix: prefix,
	}
}

// Get получает типизированное значение
func (tc *TypedCache[T]) Get(ctx context.Context, key string) (*T, error) {
	fullKey := tc.prefix + ":" + key
	
	data, err := tc.cache.Get(ctx, fullKey)
	if err != nil {
		if errors.Is(err, ErrCacheKeyNotFound) {
			return nil, nil // Возвращаем nil, если ключ не найден
		}
		return nil, err
	}
	
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	
	return &value, nil
}

// Set сохраняет типизированное значение
func (tc *TypedCache[T]) Set(ctx context.Context, key string, value *T, ttl time.Duration) error {
	fullKey := tc.prefix + ":" + key
	
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return tc.cache.Set(ctx, fullKey, data, ttl)
}

// Delete удаляет значение
func (tc *TypedCache[T]) Delete(ctx context.Context, key string) error {
	fullKey := tc.prefix + ":" + key
	return tc.cache.Delete(ctx, fullKey)
}

// СПЕЦИАЛИЗИРОВАННЫЕ КЭШИ ДЛЯ РАЗНЫХ ТИПОВ ДАННЫХ

// UserCache кэш для пользователей
type UserCache struct {
	*TypedCache[CachedUser]
}

// CachedUser структура пользователя для кэширования
type CachedUser struct {
	ID             uint      `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// NewUserCache создает кэш для пользователей
func NewUserCache(cache Cache) *UserCache {
	return &UserCache{
		TypedCache: NewTypedCache[CachedUser](cache, "user"),
	}
}

// GetByEmail получает пользователя по email
func (uc *UserCache) GetByEmail(ctx context.Context, email string) (*CachedUser, error) {
	return uc.Get(ctx, "email:"+email)
}

// SetByEmail сохраняет пользователя по email
func (uc *UserCache) SetByEmail(ctx context.Context, email string, user *CachedUser, ttl time.Duration) error {
	return uc.Set(ctx, "email:"+email, user, ttl)
}

// GetByID получает пользователя по ID
func (uc *UserCache) GetByID(ctx context.Context, userID uint) (*CachedUser, error) {
	return uc.Get(ctx, "id:"+string(rune(userID)))
}

// SetByID сохраняет пользователя по ID
func (uc *UserCache) SetByID(ctx context.Context, userID uint, user *CachedUser, ttl time.Duration) error {
	return uc.Set(ctx, "id:"+string(rune(userID)), user, ttl)
}

// SessionCache кэш для сессий
type SessionCache struct {
	*TypedCache[CachedSession]
}

// CachedSession структура сессии для кэширования
type CachedSession struct {
	UserID    uint      `json:"user_id"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// NewSessionCache создает кэш для сессий
func NewSessionCache(cache Cache) *SessionCache {
	return &SessionCache{
		TypedCache: NewTypedCache[CachedSession](cache, "session"),
	}
}

// GetByToken получает сессию по токену
func (sc *SessionCache) GetByToken(ctx context.Context, token string) (*CachedSession, error) {
	// БЕЗОПАСНОСТЬ: Хешируем токен для использования в качестве ключа
	// Это предотвращает утечку токенов через логи кэша
	hashedToken := hashToken(token)
	return sc.Get(ctx, hashedToken)
}

// SetByToken сохраняет сессию по токену
func (sc *SessionCache) SetByToken(ctx context.Context, token string, session *CachedSession, ttl time.Duration) error {
	hashedToken := hashToken(token)
	return sc.Set(ctx, hashedToken, session, ttl)
}

// DeleteByToken удаляет сессию по токену
func (sc *SessionCache) DeleteByToken(ctx context.Context, token string) error {
	hashedToken := hashToken(token)
	return sc.Delete(ctx, hashedToken)
}

// UTILITY FUNCTIONS

// hashToken хеширует токен для безопасного использования в качестве ключа кэша
func hashToken(token string) string {
	// В реальном приложении используйте криптографический хеш (SHA256)
	// Здесь упрощенная реализация для демонстрации
	hash := 0
	for _, char := range token {
		hash = hash*31 + int(char)
	}
	return string(rune(hash))
}

// ПАТТЕРНЫ КЭШИРОВАНИЯ ДЛЯ USE CASES

// CachingUserRepository декоратор для репозитория с кэшированием
// 📚 ИЗУЧИТЬ: Decorator Pattern
// https://refactoring.guru/design-patterns/decorator
type CachingUserRepository struct {
	repository UserRepository // Оригинальный репозиторий
	cache      *UserCache     // Кэш пользователей
	ttl        time.Duration  // TTL для кэшированных данных
}

// UserRepository интерфейс репозитория пользователей (упрощенный)
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*CachedUser, error)
	FindByID(ctx context.Context, id uint) (*CachedUser, error)
	Save(ctx context.Context, user *CachedUser) error
}

// NewCachingUserRepository создает репозиторий с кэшированием
func NewCachingUserRepository(repo UserRepository, cache *UserCache, ttl time.Duration) *CachingUserRepository {
	return &CachingUserRepository{
		repository: repo,
		cache:      cache,
		ttl:        ttl,
	}
}

// FindByEmail реализует паттерн Cache-Aside
func (r *CachingUserRepository) FindByEmail(ctx context.Context, email string) (*CachedUser, error) {
	// ШАГ 1: Пытаемся найти в кэше
	cachedUser, err := r.cache.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrCacheKeyNotFound) {
		// Если ошибка кэша не критична, продолжаем с БД
		// В production можно добавить логирование
	}
	
	if cachedUser != nil {
		// Cache HIT - возвращаем из кэша
		return cachedUser, nil
	}
	
	// Cache MISS - ищем в БД
	user, err := r.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	
	if user != nil {
		// Сохраняем в кэш для следующих запросов
		if cacheErr := r.cache.SetByEmail(ctx, email, user, r.ttl); cacheErr != nil {
			// Ошибка кэширования не критична, продолжаем
			// В production добавить логирование
		}
	}
	
	return user, nil
}

// FindByID реализует паттерн Cache-Aside
func (r *CachingUserRepository) FindByID(ctx context.Context, id uint) (*CachedUser, error) {
	// Аналогично FindByEmail
	cachedUser, err := r.cache.GetByID(ctx, id)
	if err != nil && !errors.Is(err, ErrCacheKeyNotFound) {
		// Логируем ошибку кэша
	}
	
	if cachedUser != nil {
		return cachedUser, nil
	}
	
	user, err := r.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if user != nil {
		if cacheErr := r.cache.SetByID(ctx, id, user, r.ttl); cacheErr != nil {
			// Логируем ошибку кэширования
		}
	}
	
	return user, nil
}

// Save реализует паттерн Write-Through
func (r *CachingUserRepository) Save(ctx context.Context, user *CachedUser) error {
	// ШАГ 1: Сохраняем в БД
	err := r.repository.Save(ctx, user)
	if err != nil {
		return err
	}
	
	// ШАГ 2: Обновляем кэш
	// Сохраняем по ID и email одновременно
	if cacheErr := r.cache.SetByID(ctx, user.ID, user, r.ttl); cacheErr != nil {
		// Логируем ошибку, но не прерываем операцию
	}
	
	if cacheErr := r.cache.SetByEmail(ctx, user.Email, user, r.ttl); cacheErr != nil {
		// Логируем ошибку, но не прерываем операцию
	}
	
	return nil
}

// КОНФИГУРАЦИЯ КЭША

// CacheConfig конфигурация кэша
type CacheConfig struct {
	// Основные настройки
	Enabled bool          `yaml:"enabled" env:"CACHE_ENABLED" default:"true"`
	TTL     time.Duration `yaml:"ttl" env:"CACHE_TTL" default:"5m"`
	
	// Настройки для разных типов данных
	UserTTL    time.Duration `yaml:"user_ttl" env:"CACHE_USER_TTL" default:"10m"`
	SessionTTL time.Duration `yaml:"session_ttl" env:"CACHE_SESSION_TTL" default:"1h"`
	
	// Настройки Redis (если используется)
	RedisURL      string `yaml:"redis_url" env:"REDIS_URL"`
	RedisPassword string `yaml:"redis_password" env:"REDIS_PASSWORD"`
	RedisDB       int    `yaml:"redis_db" env:"REDIS_DB" default:"0"`
	
	// Настройки производительности
	MaxRetries    int           `yaml:"max_retries" env:"CACHE_MAX_RETRIES" default:"3"`
	RetryDelay    time.Duration `yaml:"retry_delay" env:"CACHE_RETRY_DELAY" default:"100ms"`
	DialTimeout   time.Duration `yaml:"dial_timeout" env:"CACHE_DIAL_TIMEOUT" default:"5s"`
	ReadTimeout   time.Duration `yaml:"read_timeout" env:"CACHE_READ_TIMEOUT" default:"3s"`
	WriteTimeout  time.Duration `yaml:"write_timeout" env:"CACHE_WRITE_TIMEOUT" default:"3s"`
}

// МОНИТОРИНГ КЭША

// CacheMetrics метрики кэша для мониторинга
type CacheMetrics struct {
	Hits   int64 `json:"hits"`   // Количество попаданий в кэш
	Misses int64 `json:"misses"` // Количество промахов
	Errors int64 `json:"errors"` // Количество ошибок
}

// HitRatio возвращает коэффициент попаданий в кэш
func (m *CacheMetrics) HitRatio() float64 {
	total := m.Hits + m.Misses
	if total == 0 {
		return 0
	}
	return float64(m.Hits) / float64(total)
}

// ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ В КОММЕНТАРИЯХ:

/*
ПРАВИЛЬНОЕ использование кэша в Use Cases:

// В Use Case НЕ ДОЛЖНО быть прямой зависимости от кэша!
// Кэш инкапсулирован в Repository

type LoginUseCase struct {
    userRepo user.Repository // Это может быть CachingUserRepository
    // ... другие зависимости
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
    // Use Case не знает о кэше!
    // Кэширование происходит прозрачно в Repository
    user, err := uc.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    
    // ... остальная бизнес-логика
}

// В DI Container:
func (c *Container) initRepositories() error {
    // Создаем обычный репозиторий
    baseUserRepo := database.NewUserRepository(c.DB)
    
    // Оборачиваем в кэширующий декоратор
    userCache := cache.NewUserCache(c.Cache)
    c.UserRepo = cache.NewCachingUserRepository(baseUserRepo, userCache, 10*time.Minute)
    
    return nil
}

// Для сессий:
type AuthMiddleware struct {
    sessionCache *cache.SessionCache
    tokenValidator TokenValidator
}

func (m *AuthMiddleware) ValidateToken(ctx context.Context, token string) (*Session, error) {
    // Проверяем в кэше сессий
    session, err := m.sessionCache.GetByToken(ctx, token)
    if err == nil && session != nil {
        return session, nil
    }
    
    // Если не в кэше, валидируем токен
    validatedSession, err := m.tokenValidator.Validate(token)
    if err != nil {
        return nil, err
    }
    
    // Сохраняем в кэш
    m.sessionCache.SetByToken(ctx, token, validatedSession, time.Hour)
    
    return validatedSession, nil
}
*/

// ВАЖНЫЕ ПРИНЦИПЫ:
// 1. Кэш НЕ ДОЛЖЕН влиять на бизнес-логику
// 2. Ошибки кэша НЕ ДОЛЖНЫ прерывать операции
// 3. Кэш должен быть прозрачным для Use Cases
// 4. Всегда имейте fallback на основной источник данных
// 5. Мониторьте hit ratio и производительность кэша