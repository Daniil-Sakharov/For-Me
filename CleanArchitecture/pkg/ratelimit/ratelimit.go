package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// 📚 ИЗУЧИТЬ: Rate Limiting алгоритмы и паттерны
// https://en.wikipedia.org/wiki/Rate_limiting
// https://blog.cloudflare.com/counting-things-a-lot-of-different-things/
// https://stripe.com/blog/rate-limiters
// https://konghq.com/blog/how-to-design-a-scalable-rate-limiting-algorithm
//
// ОСНОВНЫЕ АЛГОРИТМЫ RATE LIMITING:
// 1. TOKEN BUCKET - токены добавляются с фиксированной скоростью
// 2. LEAKY BUCKET - запросы обрабатываются с постоянной скоростью
// 3. FIXED WINDOW - фиксированные временные окна
// 4. SLIDING WINDOW - скользящие временные окна
// 5. SLIDING WINDOW LOG - точное скользящее окно с логом запросов
//
// ЗАЧЕМ НУЖЕН RATE LIMITING:
// 1. ЗАЩИТА ОТ DDoS атак
// 2. ПРЕДОТВРАЩЕНИЕ ЗЛОУПОТРЕБЛЕНИЙ API
// 3. ОБЕСПЕЧЕНИЕ СПРАВЕДЛИВОГО ИСПОЛЬЗОВАНИЯ РЕСУРСОВ
// 4. ЗАЩИТА ОТ BRUTE FORCE атак
// 5. КОНТРОЛЬ НАГРУЗКИ НА СИСТЕМУ

// RateLimiter интерфейс для rate limiting
type RateLimiter interface {
	// Allow проверяет, разрешен ли запрос для данного ключа
	Allow(ctx context.Context, key string) (*Result, error)
	
	// Reset сбрасывает лимит для ключа
	Reset(ctx context.Context, key string) error
	
	// GetInfo возвращает информацию о текущем состоянии лимита
	GetInfo(ctx context.Context, key string) (*Info, error)
}

// Result результат проверки rate limit
type Result struct {
	// Allowed указывает, разрешен ли запрос
	Allowed bool
	
	// Limit максимальное количество запросов
	Limit int64
	
	// Remaining количество оставшихся запросов
	Remaining int64
	
	// ResetTime время сброса счетчика
	ResetTime time.Time
	
	// RetryAfter через сколько можно повторить запрос (если Allowed = false)
	RetryAfter time.Duration
}

// Info информация о состоянии rate limiter
type Info struct {
	// Key ключ лимита
	Key string
	
	// Limit максимальное количество запросов
	Limit int64
	
	// Used количество использованных запросов
	Used int64
	
	// Remaining количество оставшихся запросов
	Remaining int64
	
	// ResetTime время сброса счетчика
	ResetTime time.Time
	
	// WindowSize размер временного окна
	WindowSize time.Duration
}

// Config конфигурация rate limiter
type Config struct {
	// Limit количество запросов
	Limit int64
	
	// Window временное окно
	Window time.Duration
	
	// KeyPrefix префикс для ключей в хранилище
	KeyPrefix string
}

// Стандартные ошибки
var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrInvalidConfig     = errors.New("invalid rate limiter configuration")
	ErrStorageUnavailable = errors.New("rate limiter storage unavailable")
)

// РЕАЛИЗАЦИЯ FIXED WINDOW RATE LIMITER
// 📚 ИЗУЧИТЬ: Fixed Window Counter
// https://en.wikipedia.org/wiki/Rate_limiting#Fixed_window_counter

// FixedWindowLimiter реализация Fixed Window алгоритма
type FixedWindowLimiter struct {
	storage Storage
	config  Config
}

// Storage интерфейс для хранения счетчиков
type Storage interface {
	// Increment увеличивает счетчик и возвращает новое значение
	Increment(ctx context.Context, key string, expiry time.Duration) (int64, error)
	
	// Get получает текущее значение счетчика
	Get(ctx context.Context, key string) (int64, error)
	
	// Set устанавливает значение счетчика
	Set(ctx context.Context, key string, value int64, expiry time.Duration) error
	
	// Delete удаляет ключ
	Delete(ctx context.Context, key string) error
	
	// TTL возвращает время жизни ключа
	TTL(ctx context.Context, key string) (time.Duration, error)
}

// NewFixedWindowLimiter создает новый Fixed Window rate limiter
func NewFixedWindowLimiter(storage Storage, config Config) (*FixedWindowLimiter, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	
	return &FixedWindowLimiter{
		storage: storage,
		config:  config,
	}, nil
}

// Allow проверяет, разрешен ли запрос
func (rl *FixedWindowLimiter) Allow(ctx context.Context, key string) (*Result, error) {
	// Создаем ключ с временным окном
	windowKey := rl.createWindowKey(key, time.Now())
	
	// Увеличиваем счетчик
	count, err := rl.storage.Increment(ctx, windowKey, rl.config.Window)
	if err != nil {
		return nil, fmt.Errorf("failed to increment counter: %w", err)
	}
	
	// Вычисляем время сброса (начало следующего окна)
	resetTime := rl.getNextWindowStart(time.Now())
	
	result := &Result{
		Allowed:   count <= rl.config.Limit,
		Limit:     rl.config.Limit,
		Remaining: max(0, rl.config.Limit-count),
		ResetTime: resetTime,
	}
	
	if !result.Allowed {
		result.RetryAfter = time.Until(resetTime)
	}
	
	return result, nil
}

// Reset сбрасывает счетчик для ключа
func (rl *FixedWindowLimiter) Reset(ctx context.Context, key string) error {
	windowKey := rl.createWindowKey(key, time.Now())
	return rl.storage.Delete(ctx, windowKey)
}

// GetInfo возвращает информацию о лимите
func (rl *FixedWindowLimiter) GetInfo(ctx context.Context, key string) (*Info, error) {
	windowKey := rl.createWindowKey(key, time.Now())
	
	used, err := rl.storage.Get(ctx, windowKey)
	if err != nil {
		return nil, err
	}
	
	return &Info{
		Key:        key,
		Limit:      rl.config.Limit,
		Used:       used,
		Remaining:  max(0, rl.config.Limit-used),
		ResetTime:  rl.getNextWindowStart(time.Now()),
		WindowSize: rl.config.Window,
	}, nil
}

// createWindowKey создает ключ для текущего временного окна
func (rl *FixedWindowLimiter) createWindowKey(key string, now time.Time) string {
	windowStart := rl.getWindowStart(now)
	return fmt.Sprintf("%s:%s:%d", rl.config.KeyPrefix, key, windowStart.Unix())
}

// getWindowStart возвращает начало текущего временного окна
func (rl *FixedWindowLimiter) getWindowStart(now time.Time) time.Time {
	windowSeconds := int64(rl.config.Window.Seconds())
	return time.Unix((now.Unix()/windowSeconds)*windowSeconds, 0)
}

// getNextWindowStart возвращает начало следующего временного окна
func (rl *FixedWindowLimiter) getNextWindowStart(now time.Time) time.Time {
	return rl.getWindowStart(now).Add(rl.config.Window)
}

// РЕАЛИЗАЦИЯ SLIDING WINDOW RATE LIMITER
// 📚 ИЗУЧИТЬ: Sliding Window Log
// https://konghq.com/blog/how-to-design-a-scalable-rate-limiting-algorithm

// SlidingWindowLimiter реализация Sliding Window алгоритма
type SlidingWindowLimiter struct {
	storage Storage
	config  Config
}

// NewSlidingWindowLimiter создает новый Sliding Window rate limiter
func NewSlidingWindowLimiter(storage Storage, config Config) (*SlidingWindowLimiter, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	
	return &SlidingWindowLimiter{
		storage: storage,
		config:  config,
	}, nil
}

// Allow проверяет, разрешен ли запрос (упрощенная реализация)
func (rl *SlidingWindowLimiter) Allow(ctx context.Context, key string) (*Result, error) {
	now := time.Now()
	
	// Для упрощения используем два окна: текущее и предыдущее
	currentWindowKey := rl.createWindowKey(key, now)
	previousWindowKey := rl.createWindowKey(key, now.Add(-rl.config.Window))
	
	// Получаем счетчики для обоих окон
	currentCount, err := rl.storage.Get(ctx, currentWindowKey)
	if err != nil {
		return nil, err
	}
	
	previousCount, err := rl.storage.Get(ctx, previousWindowKey)
	if err != nil {
		return nil, err
	}
	
	// Вычисляем вес предыдущего окна
	windowStart := rl.getWindowStart(now)
	elapsedInWindow := now.Sub(windowStart)
	weightOfPreviousWindow := float64(rl.config.Window-elapsedInWindow) / float64(rl.config.Window)
	
	// Приблизительное количество запросов в скользящем окне
	estimatedCount := int64(float64(previousCount)*weightOfPreviousWindow) + currentCount
	
	allowed := estimatedCount < rl.config.Limit
	
	if allowed {
		// Увеличиваем счетчик текущего окна
		currentCount, err = rl.storage.Increment(ctx, currentWindowKey, rl.config.Window)
		if err != nil {
			return nil, err
		}
		
		// Пересчитываем после инкремента
		estimatedCount = int64(float64(previousCount)*weightOfPreviousWindow) + currentCount
	}
	
	resetTime := rl.getNextWindowStart(now)
	
	result := &Result{
		Allowed:   allowed,
		Limit:     rl.config.Limit,
		Remaining: max(0, rl.config.Limit-estimatedCount),
		ResetTime: resetTime,
	}
	
	if !result.Allowed {
		result.RetryAfter = time.Until(resetTime)
	}
	
	return result, nil
}

// Reset и GetInfo аналогичны FixedWindowLimiter
func (rl *SlidingWindowLimiter) Reset(ctx context.Context, key string) error {
	currentWindowKey := rl.createWindowKey(key, time.Now())
	previousWindowKey := rl.createWindowKey(key, time.Now().Add(-rl.config.Window))
	
	_ = rl.storage.Delete(ctx, currentWindowKey)
	_ = rl.storage.Delete(ctx, previousWindowKey)
	
	return nil
}

func (rl *SlidingWindowLimiter) GetInfo(ctx context.Context, key string) (*Info, error) {
	// Аналогично Allow, но без инкремента
	now := time.Now()
	currentWindowKey := rl.createWindowKey(key, now)
	previousWindowKey := rl.createWindowKey(key, now.Add(-rl.config.Window))
	
	currentCount, _ := rl.storage.Get(ctx, currentWindowKey)
	previousCount, _ := rl.storage.Get(ctx, previousWindowKey)
	
	windowStart := rl.getWindowStart(now)
	elapsedInWindow := now.Sub(windowStart)
	weightOfPreviousWindow := float64(rl.config.Window-elapsedInWindow) / float64(rl.config.Window)
	
	estimatedCount := int64(float64(previousCount)*weightOfPreviousWindow) + currentCount
	
	return &Info{
		Key:        key,
		Limit:      rl.config.Limit,
		Used:       estimatedCount,
		Remaining:  max(0, rl.config.Limit-estimatedCount),
		ResetTime:  rl.getNextWindowStart(now),
		WindowSize: rl.config.Window,
	}, nil
}

func (rl *SlidingWindowLimiter) createWindowKey(key string, now time.Time) string {
	windowStart := rl.getWindowStart(now)
	return fmt.Sprintf("%s:%s:%d", rl.config.KeyPrefix, key, windowStart.Unix())
}

func (rl *SlidingWindowLimiter) getWindowStart(now time.Time) time.Time {
	windowSeconds := int64(rl.config.Window.Seconds())
	return time.Unix((now.Unix()/windowSeconds)*windowSeconds, 0)
}

func (rl *SlidingWindowLimiter) getNextWindowStart(now time.Time) time.Time {
	return rl.getWindowStart(now).Add(rl.config.Window)
}

// МНОГОУРОВНЕВЫЙ RATE LIMITER
// 📚 ИЗУЧИТЬ: Hierarchical Rate Limiting
// https://stripe.com/blog/rate-limiters

// MultiLevelLimiter многоуровневый rate limiter
type MultiLevelLimiter struct {
	limiters []RateLimiter
}

// NewMultiLevelLimiter создает многоуровневый rate limiter
func NewMultiLevelLimiter(limiters ...RateLimiter) *MultiLevelLimiter {
	return &MultiLevelLimiter{
		limiters: limiters,
	}
}

// Allow проверяет все уровни лимитов
func (ml *MultiLevelLimiter) Allow(ctx context.Context, key string) (*Result, error) {
	var strictestResult *Result
	
	// Проверяем все лимиты
	for _, limiter := range ml.limiters {
		result, err := limiter.Allow(ctx, key)
		if err != nil {
			return nil, err
		}
		
		// Если любой лимит превышен, запрос запрещен
		if !result.Allowed {
			return result, nil
		}
		
		// Запоминаем самый строгий лимит
		if strictestResult == nil || result.Remaining < strictestResult.Remaining {
			strictestResult = result
		}
	}
	
	return strictestResult, nil
}

// Reset сбрасывает все лимиты
func (ml *MultiLevelLimiter) Reset(ctx context.Context, key string) error {
	for _, limiter := range ml.limiters {
		if err := limiter.Reset(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// GetInfo возвращает информацию о самом строгом лимите
func (ml *MultiLevelLimiter) GetInfo(ctx context.Context, key string) (*Info, error) {
	var strictestInfo *Info
	
	for _, limiter := range ml.limiters {
		info, err := limiter.GetInfo(ctx, key)
		if err != nil {
			return nil, err
		}
		
		if strictestInfo == nil || info.Remaining < strictestInfo.Remaining {
			strictestInfo = info
		}
	}
	
	return strictestInfo, nil
}

// СПЕЦИАЛИЗИРОВАННЫЕ RATE LIMITERS

// LoginRateLimiter специальный лимитер для попыток входа
type LoginRateLimiter struct {
	ipLimiter    RateLimiter // Лимит по IP
	userLimiter  RateLimiter // Лимит по пользователю
	globalLimiter RateLimiter // Глобальный лимит
}

// NewLoginRateLimiter создает лимитер для логина
func NewLoginRateLimiter(storage Storage) *LoginRateLimiter {
	return &LoginRateLimiter{
		ipLimiter: &FixedWindowLimiter{
			storage: storage,
			config: Config{
				Limit:     10,           // 10 попыток
				Window:    time.Minute,  // за минуту
				KeyPrefix: "login_ip",
			},
		},
		userLimiter: &FixedWindowLimiter{
			storage: storage,
			config: Config{
				Limit:     5,            // 5 попыток
				Window:    time.Minute,  // за минуту
				KeyPrefix: "login_user",
			},
		},
		globalLimiter: &FixedWindowLimiter{
			storage: storage,
			config: Config{
				Limit:     1000,         // 1000 попыток
				Window:    time.Minute,  // за минуту
				KeyPrefix: "login_global",
			},
		},
	}
}

// CheckLogin проверяет лимиты для попытки входа
func (lr *LoginRateLimiter) CheckLogin(ctx context.Context, ip, userID string) (*Result, error) {
	// Проверяем лимит по IP
	ipResult, err := lr.ipLimiter.Allow(ctx, ip)
	if err != nil {
		return nil, err
	}
	if !ipResult.Allowed {
		return ipResult, nil
	}
	
	// Проверяем лимит по пользователю
	if userID != "" {
		userResult, err := lr.userLimiter.Allow(ctx, userID)
		if err != nil {
			return nil, err
		}
		if !userResult.Allowed {
			return userResult, nil
		}
	}
	
	// Проверяем глобальный лимит
	globalResult, err := lr.globalLimiter.Allow(ctx, "global")
	if err != nil {
		return nil, err
	}
	
	return globalResult, nil
}

// UTILITY FUNCTIONS

// validateConfig проверяет конфигурацию
func validateConfig(config Config) error {
	if config.Limit <= 0 {
		return ErrInvalidConfig
	}
	if config.Window <= 0 {
		return ErrInvalidConfig
	}
	return nil
}

// max возвращает максимальное значение
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// КОНФИГУРАЦИИ ДЛЯ РАЗНЫХ СЛУЧАЕВ ИСПОЛЬЗОВАНИЯ

// GetAPIRateLimitConfig конфигурация для API endpoints
func GetAPIRateLimitConfig() Config {
	return Config{
		Limit:     100,          // 100 запросов
		Window:    time.Minute,  // за минуту
		KeyPrefix: "api",
	}
}

// GetAuthRateLimitConfig конфигурация для аутентификации
func GetAuthRateLimitConfig() Config {
	return Config{
		Limit:     5,            // 5 попыток
		Window:    time.Minute,  // за минуту
		KeyPrefix: "auth",
	}
}

// GetRegistrationRateLimitConfig конфигурация для регистрации
func GetRegistrationRateLimitConfig() Config {
	return Config{
		Limit:     3,             // 3 регистрации
		Window:    time.Hour,     // за час
		KeyPrefix: "registration",
	}
}

// ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ:

/*
// В middleware:
func RateLimitMiddleware(limiter ratelimit.RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := getClientKey(r) // IP или user ID
            
            result, err := limiter.Allow(r.Context(), key)
            if err != nil {
                http.Error(w, "Rate limiter error", http.StatusInternalServerError)
                return
            }
            
            // Добавляем заголовки с информацией о лимитах
            w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
            w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
            w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime.Unix(), 10))
            
            if !result.Allowed {
                w.Header().Set("Retry-After", strconv.FormatInt(int64(result.RetryAfter.Seconds()), 10))
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// В Use Case:
type LoginUseCase struct {
    rateLimiter *ratelimit.LoginRateLimiter
    // ... другие зависимости
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest, clientIP string) (*LoginResponse, error) {
    // Проверяем rate limit перед попыткой входа
    result, err := uc.rateLimiter.CheckLogin(ctx, clientIP, req.Email)
    if err != nil {
        return nil, err
    }
    
    if !result.Allowed {
        return nil, ratelimit.ErrRateLimitExceeded
    }
    
    // Продолжаем с логикой входа...
}

// В DI Container:
func (c *Container) initRateLimiters() error {
    storage := redis.NewRedisStorage(c.Redis) // или memory storage
    
    // API rate limiter
    c.APIRateLimiter = ratelimit.NewFixedWindowLimiter(storage, ratelimit.GetAPIRateLimitConfig())
    
    // Login rate limiter
    c.LoginRateLimiter = ratelimit.NewLoginRateLimiter(storage)
    
    return nil
}
*/

// ВАЖНЫЕ ПРИНЦИПЫ:
// 1. Выбирайте подходящий алгоритм для вашего случая использования
// 2. Используйте разные лимиты для разных типов операций
// 3. Предоставляйте клиентам информацию о лимитах через HTTP заголовки
// 4. Логируйте превышения лимитов для мониторинга
// 5. Рассмотрите whitelist для доверенных клиентов
// 6. Тестируйте rate limiting в нагрузочных тестах
// 7. Мониторьте производительность rate limiter'а