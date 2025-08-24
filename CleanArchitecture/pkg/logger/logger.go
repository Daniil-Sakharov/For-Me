package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// 📚 ИЗУЧИТЬ: Structured Logging в Go
// https://pkg.go.dev/log/slog
// https://betterstack.com/community/guides/logging/logging-in-go/
// https://go.dev/blog/slog
//
// ПРИНЦИПЫ ХОРОШЕГО ЛОГИРОВАНИЯ:
// 1. СТРУКТУРИРОВАННОСТЬ - используйте ключ-значение, не конкатенацию строк
// 2. КОНТЕКСТНОСТЬ - передавайте context для трейсинга
// 3. УРОВНИ - используйте правильные уровни логирования
// 4. ПРОИЗВОДИТЕЛЬНОСТЬ - логирование не должно замедлять приложение
// 5. БЕЗОПАСНОСТЬ - никогда не логируйте пароли и токены!

// Logger определяет интерфейс для логирования в Clean Architecture
// ВАЖНО: Интерфейс находится в pkg/, потому что он используется во всех слоях
type Logger interface {
	// Debug логирует отладочную информацию (только в dev окружении)
	Debug(msg string, args ...any)
	
	// Info логирует общую информацию о работе приложения
	Info(msg string, args ...any)
	
	// Warn логирует предупреждения (не критичные ошибки)
	Warn(msg string, args ...any)
	
	// Error логирует ошибки
	Error(msg string, args ...any)
	
	// With добавляет поля к логгеру (для создания контекстного логгера)
	With(args ...any) Logger
	
	// WithContext добавляет информацию из context (trace_id, user_id и т.д.)
	WithContext(ctx context.Context) Logger
}

// SlogLogger реализует интерфейс Logger используя стандартный slog
// 📚 ИЗУЧИТЬ: Adapter Pattern
// https://refactoring.guru/design-patterns/adapter
type SlogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger создает новый structured logger
func NewSlogLogger(env string) Logger {
	var handler slog.Handler
	
	// Настраиваем формат логов в зависимости от окружения
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// Добавляем информацию о файле и строке для отладки
		AddSource: true,
		// Кастомная функция для замены чувствительных данных
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// БЕЗОПАСНОСТЬ: Скрываем чувствительные поля
			if a.Key == "password" || a.Key == "token" || a.Key == "secret" {
				return slog.String(a.Key, "[REDACTED]")
			}
			
			// Форматируем время в ISO 8601
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, a.Value.Time().Format(time.RFC3339))
			}
			
			return a
		},
	}
	
	switch env {
	case "development", "dev":
		// В dev окружении используем текстовый формат для читаемости
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "production", "prod":
		// В prod окружении используем JSON для парсинга системами мониторинга
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		// По умолчанию JSON формат
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	
	logger := slog.New(handler)
	
	return &SlogLogger{
		logger: logger,
	}
}

// Debug реализует интерфейс Logger
func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info реализует интерфейс Logger
func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn реализует интерфейс Logger
func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error реализует интерфейс Logger
func (l *SlogLogger) Error(msg string, args ...any) {
	// Добавляем информацию о стеке вызовов для ошибок
	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = append(args, "caller_file", file, "caller_line", line)
	}
	
	l.logger.Error(msg, args...)
}

// With реализует интерфейс Logger
func (l *SlogLogger) With(args ...any) Logger {
	return &SlogLogger{
		logger: l.logger.With(args...),
	}
}

// WithContext реализует интерфейс Logger
func (l *SlogLogger) WithContext(ctx context.Context) Logger {
	// Извлекаем полезную информацию из context
	var contextArgs []any
	
	// Проверяем наличие trace_id (для distributed tracing)
	if traceID := getTraceIDFromContext(ctx); traceID != "" {
		contextArgs = append(contextArgs, "trace_id", traceID)
	}
	
	// Проверяем наличие user_id (для аудита)
	if userID := getUserIDFromContext(ctx); userID != 0 {
		contextArgs = append(contextArgs, "user_id", userID)
	}
	
	// Проверяем наличие request_id (для корреляции запросов)
	if requestID := getRequestIDFromContext(ctx); requestID != "" {
		contextArgs = append(contextArgs, "request_id", requestID)
	}
	
	if len(contextArgs) > 0 {
		return l.With(contextArgs...)
	}
	
	return l
}

// HELPER FUNCTIONS для извлечения данных из context
// В реальном приложении эти функции должны соответствовать вашему middleware

// ContextKey тип для ключей context (избегаем коллизий)
type ContextKey string

const (
	TraceIDKey   ContextKey = "trace_id"
	UserIDKey    ContextKey = "user_id"
	RequestIDKey ContextKey = "request_id"
)

// getTraceIDFromContext извлекает trace ID из context
func getTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// getUserIDFromContext извлекает user ID из context
func getUserIDFromContext(ctx context.Context) uint {
	if userID, ok := ctx.Value(UserIDKey).(uint); ok {
		return userID
	}
	return 0
}

// getRequestIDFromContext извлекает request ID из context
func getRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// UTILITY FUNCTIONS для создания context с метаданными

// WithTraceID добавляет trace ID в context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithUserID добавляет user ID в context
func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithRequestID добавляет request ID в context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// СПЕЦИАЛИЗИРОВАННЫЕ ЛОГГЕРЫ для разных слоев архитектуры

// DomainLogger создает логгер для доменного слоя
func (l *SlogLogger) DomainLogger() Logger {
	return l.With("layer", "domain")
}

// UseCaseLogger создает логгер для слоя Use Cases
func (l *SlogLogger) UseCaseLogger() Logger {
	return l.With("layer", "usecase")
}

// InfrastructureLogger создает логгер для инфраструктурного слоя
func (l *SlogLogger) InfrastructureLogger() Logger {
	return l.With("layer", "infrastructure")
}

// InterfaceLogger создает логгер для слоя интерфейсов
func (l *SlogLogger) InterfaceLogger() Logger {
	return l.With("layer", "interface")
}

// ПРОИЗВОДИТЕЛЬНОСТЬ: Lazy logging для дорогих операций
// 📚 ИЗУЧИТЬ: Lazy evaluation
// https://en.wikipedia.org/wiki/Lazy_evaluation

// LazyValue позволяет отложить вычисление значения до момента логирования
type LazyValue func() any

// LazyString создает ленивое строковое значение
func LazyString(fn func() string) LazyValue {
	return func() any { return fn() }
}

// LazyInt создает ленивое числовое значение
func LazyInt(fn func() int) LazyValue {
	return func() any { return fn() }
}

// ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ В КОММЕНТАРИЯХ:

/*
ПРАВИЛЬНОЕ использование логгера:

// В Use Case:
func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
    logger := uc.logger.WithContext(ctx).With(
        "operation", "user_login",
        "email", req.Email, // Email можно логировать
        // "password", req.Password, // НИКОГДА не логируйте пароли!
    )
    
    logger.Info("Login attempt started")
    
    user, err := uc.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        logger.Error("Failed to find user", "error", err)
        return nil, err
    }
    
    if user == nil {
        logger.Warn("Login attempt for non-existent user")
        return nil, ErrUserNotFound
    }
    
    // ... бизнес-логика ...
    
    logger.Info("Login successful", "user_id", user.ID)
    return response, nil
}

// В Repository:
func (r *UserRepository) Save(ctx context.Context, user *user.User) error {
    logger := r.logger.WithContext(ctx).With(
        "operation", "user_save",
        "user_id", user.ID,
    )
    
    start := time.Now()
    defer func() {
        logger.Debug("Save operation completed", 
            "duration_ms", time.Since(start).Milliseconds())
    }()
    
    logger.Debug("Saving user to database")
    
    // ... SQL операции ...
    
    return nil
}

// В Controller:
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
    // Создаем request ID для трейсинга
    requestID := generateRequestID()
    ctx := WithRequestID(r.Context(), requestID)
    
    logger := c.logger.WithContext(ctx).With(
        "handler", "auth_login",
        "method", r.Method,
        "path", r.URL.Path,
        "user_agent", r.UserAgent(),
        "ip", getClientIP(r),
    )
    
    logger.Info("HTTP request received")
    
    // ... обработка запроса ...
    
    logger.Info("HTTP response sent", "status_code", http.StatusOK)
}
*/

// БЕЗОПАСНОСТЬ: Функции для безопасного логирования

// SafeEmail маскирует часть email для логирования
func SafeEmail(email string) string {
	if len(email) < 3 {
		return "***"
	}
	
	atIndex := -1
	for i, r := range email {
		if r == '@' {
			atIndex = i
			break
		}
	}
	
	if atIndex <= 0 {
		return "***"
	}
	
	// Показываем первый символ + *** + домен
	return email[:1] + "***" + email[atIndex:]
}

// SafeUserID безопасно логирует user ID (в некоторых случаях может потребоваться маскировка)
func SafeUserID(userID uint) any {
	// В большинстве случаев user ID можно логировать как есть
	// Но в некоторых высокобезопасных системах может потребоваться хеширование
	return userID
}

// ErrorWithStack добавляет stack trace к ошибке для логирования
func ErrorWithStack(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "nil")
	}
	
	// Получаем stack trace
	stack := make([]byte, 4096)
	n := runtime.Stack(stack, false)
	
	return slog.Group("error",
		"message", err.Error(),
		"stack", string(stack[:n]),
	)
}

// МЕТРИКИ: Интеграция с системами мониторинга
// 📚 ИЗУЧИТЬ: Observability (Logs, Metrics, Traces)
// https://opentelemetry.io/docs/concepts/observability-primer/

// MetricLogger расширяет обычный логгер метриками
type MetricLogger struct {
	Logger
	metrics MetricsCollector
}

// MetricsCollector интерфейс для сбора метрик
type MetricsCollector interface {
	IncCounter(name string, tags map[string]string)
	RecordDuration(name string, duration time.Duration, tags map[string]string)
	SetGauge(name string, value float64, tags map[string]string)
}

// NewMetricLogger создает логгер с метриками
func NewMetricLogger(logger Logger, metrics MetricsCollector) Logger {
	return &MetricLogger{
		Logger:  logger,
		metrics: metrics,
	}
}

// Error переопределяет Error для сбора метрик ошибок
func (ml *MetricLogger) Error(msg string, args ...any) {
	// Увеличиваем счетчик ошибок
	ml.metrics.IncCounter("errors_total", map[string]string{
		"level": "error",
	})
	
	// Вызываем оригинальный Error
	ml.Logger.Error(msg, args...)
}

// ПРИМЕЧАНИЕ: В реальном приложении MetricsCollector должен быть реализован
// с использованием Prometheus, StatsD или другой системы метрик