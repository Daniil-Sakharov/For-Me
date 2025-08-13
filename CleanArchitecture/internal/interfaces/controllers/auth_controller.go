package controllers

import (
	"encoding/json"
	"net/http"
	
	"clean-url-shortener/internal/usecase/auth"
	"clean-url-shortener/internal/infrastructure/web"
	"github.com/go-playground/validator/v10"
)

// AuthController обрабатывает HTTP запросы для аутентификации
// ЭТО АДАПТЕР между HTTP и Use Cases
type AuthController struct {
	registerUC *auth.RegisterUseCase // Use Case для регистрации
	loginUC    *auth.LoginUseCase    // Use Case для входа
	validator  *validator.Validate   // Валидатор для входных данных
}

// NewAuthController создает новый контроллер аутентификации
func NewAuthController(
	registerUC *auth.RegisterUseCase,
	loginUC *auth.LoginUseCase,
) *AuthController {
	return &AuthController{
		registerUC: registerUC,
		loginUC:    loginUC,
		validator:  validator.New(),
	}
}

// Register обрабатывает HTTP запрос на регистрацию пользователя
// POST /auth/register
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	// ШАГ 1: Декодируем JSON из тела запроса
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// ШАГ 2: Валидируем входные данные
	if err := c.validator.Struct(req); err != nil {
		c.writeErrorResponse(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// ШАГ 3: Вызываем Use Case для регистрации
	response, err := c.registerUC.Execute(r.Context(), req)
	if err != nil {
		c.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ШАГ 4: Возвращаем успешный ответ
	c.writeJSONResponse(w, response, http.StatusCreated)
}

// Login обрабатывает HTTP запрос на вход пользователя
// POST /auth/login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	// ШАГ 1: Декодируем JSON из тела запроса
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// ШАГ 2: Валидируем входные данные
	if err := c.validator.Struct(req); err != nil {
		c.writeErrorResponse(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// ШАГ 3: Вызываем Use Case для входа
	response, err := c.loginUC.Execute(r.Context(), req)
	if err != nil {
		c.writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// ШАГ 4: Возвращаем успешный ответ
	c.writeJSONResponse(w, response, http.StatusOK)
}

// RegisterRoutes регистрирует маршруты контроллера
func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	// Применяем middleware для всех маршрутов аутентификации
	mux.HandleFunc("POST /auth/register", web.ChainMiddleware(
		c.Register,
		web.CORSMiddleware,
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
	))

	mux.HandleFunc("POST /auth/login", web.ChainMiddleware(
		c.Login,
		web.CORSMiddleware,
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
	))
}

// ErrorResponse структура для ошибок
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// writeErrorResponse записывает ошибку в HTTP ответ
func (c *AuthController) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}
	
	json.NewEncoder(w).Encode(errorResp)
}

// writeJSONResponse записывает успешный JSON ответ
func (c *AuthController) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Если не удалось закодировать ответ, отправляем ошибку сервера
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// ПРИНЦИПЫ КОНТРОЛЛЕРА:
// 1. Тонкий слой - только HTTP протокол и валидация
// 2. Делегирует всю бизнес-логику Use Cases
// 3. Не содержит бизнес-правил
// 4. Обрабатывает только HTTP специфичные задачи
// 5. Стандартизированные ответы об ошибках