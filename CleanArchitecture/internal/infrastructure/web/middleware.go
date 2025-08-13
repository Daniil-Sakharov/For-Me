package web

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
	
	"clean-url-shortener/internal/usecase/auth"
)

// ContextKey тип для ключей контекста
type ContextKey string

const (
	// UserIDKey ключ для ID пользователя в контексте
	UserIDKey ContextKey = "user_id"
	// UserEmailKey ключ для email пользователя в контексте
	UserEmailKey ContextKey = "user_email"
)

// AuthMiddleware middleware для аутентификации
type AuthMiddleware struct {
	tokenGenerator auth.TokenGenerator
}

// NewAuthMiddleware создает новый middleware для аутентификации
func NewAuthMiddleware(tokenGenerator auth.TokenGenerator) *AuthMiddleware {
	return &AuthMiddleware{
		tokenGenerator: tokenGenerator,
	}
}

// RequireAuth middleware который требует аутентификации
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Валидируем токен
		tokenData, err := m.tokenGenerator.Validate(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем данные пользователя в контекст
		ctx := context.WithValue(r.Context(), UserIDKey, tokenData.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, tokenData.Email)
		r = r.WithContext(ctx)

		// Передаем управление следующему handler
		next(w, r)
	}
}

// OptionalAuth middleware который не требует аутентификации, но извлекает данные если токен есть
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				if tokenData, err := m.tokenGenerator.Validate(token); err == nil {
					ctx := context.WithValue(r.Context(), UserIDKey, tokenData.UserID)
					ctx = context.WithValue(ctx, UserEmailKey, tokenData.Email)
					r = r.WithContext(ctx)
				}
			}
		}
		next(w, r)
	}
}

// CORSMiddleware middleware для CORS
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Обрабатываем preflight запросы
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

// LoggingMiddleware middleware для логирования запросов
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем wrapper для ResponseWriter чтобы захватить статус код
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Выполняем запрос
		next(wrapped, r)

		// Логируем информацию о запросе
		duration := time.Since(start)
		log.Printf(
			"%s %s %d %v %s",
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			duration,
			r.RemoteAddr,
		)
	}
}

// responseWriter обертка для захвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RecoveryMiddleware middleware для восстановления от паник
func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// ChainMiddleware объединяет несколько middleware в цепочку
func ChainMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	// Применяем middleware в обратном порядке
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// GetUserIDFromContext извлекает ID пользователя из контекста
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

// GetUserEmailFromContext извлекает email пользователя из контекста
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// ПРИНЦИПЫ MIDDLEWARE:
// 1. Каждый middleware имеет одну ответственность
// 2. Middleware можно комбинировать в любом порядке
// 3. Контекст используется для передачи данных между middleware
// 4. Graceful обработка ошибок без падения сервера