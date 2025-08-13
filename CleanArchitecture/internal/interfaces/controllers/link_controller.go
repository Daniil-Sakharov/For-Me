package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	
	"clean-url-shortener/internal/usecase/link"
	"clean-url-shortener/internal/infrastructure/web"
	"github.com/go-playground/validator/v10"
)

// LinkController обрабатывает HTTP запросы для работы со ссылками
type LinkController struct {
	createLinkUC   *link.CreateLinkUseCase // Use Case для создания ссылок
	redirectUC     *link.RedirectUseCase   // Use Case для редиректа
	validator      *validator.Validate     // Валидатор входных данных
	authMiddleware *web.AuthMiddleware     // Middleware для аутентификации
}

// NewLinkController создает новый контроллер ссылок
func NewLinkController(
	createLinkUC *link.CreateLinkUseCase,
	redirectUC *link.RedirectUseCase,
	authMiddleware *web.AuthMiddleware,
) *LinkController {
	return &LinkController{
		createLinkUC:   createLinkUC,
		redirectUC:     redirectUC,
		validator:      validator.New(),
		authMiddleware: authMiddleware,
	}
}

// CreateLink обрабатывает HTTP запрос на создание короткой ссылки
// POST /links
func (c *LinkController) CreateLink(w http.ResponseWriter, r *http.Request) {
	// ШАГ 1: Извлекаем ID пользователя из контекста (установлен middleware)
	userID, ok := web.GetUserIDFromContext(r.Context())
	if !ok {
		c.writeErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// ШАГ 2: Декодируем JSON из тела запроса
	var req link.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// ШАГ 3: Устанавливаем ID пользователя в запрос
	req.UserID = userID

	// ШАГ 4: Валидируем входные данные
	if err := c.validator.Struct(req); err != nil {
		c.writeErrorResponse(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// ШАГ 5: Вызываем Use Case для создания ссылки
	response, err := c.createLinkUC.Execute(r.Context(), req)
	if err != nil {
		c.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ШАГ 6: Возвращаем успешный ответ
	c.writeJSONResponse(w, response, http.StatusCreated)
}

// Redirect обрабатывает HTTP запрос на редирект по короткой ссылке
// GET /{shortCode}
func (c *LinkController) Redirect(w http.ResponseWriter, r *http.Request) {
	// ШАГ 1: Извлекаем короткий код из URL пути
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		c.writeErrorResponse(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// ШАГ 2: Создаем запрос для Use Case
	req := link.RedirectRequest{
		ShortCode: path,
		UserAgent: r.Header.Get("User-Agent"),
		IPAddress: c.getClientIP(r),
		Referer:   r.Header.Get("Referer"),
	}

	// ШАГ 3: Вызываем Use Case для редиректа
	response, err := c.redirectUC.Execute(r.Context(), req)
	if err != nil {
		c.writeErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	// ШАГ 4: Выполняем HTTP редирект
	http.Redirect(w, r, response.OriginalURL, http.StatusFound)
}

// GetUserLinks возвращает все ссылки пользователя (заглушка для примера)
// GET /links
func (c *LinkController) GetUserLinks(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID пользователя из контекста
	userID, ok := web.GetUserIDFromContext(r.Context())
	if !ok {
		c.writeErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Извлекаем параметры пагинации из query параметров
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20 // По умолчанию
	}

	offset, _ := strconv.Atoi(offsetStr)
	if offset < 0 {
		offset = 0
	}

	// Для упрощения примера возвращаем заглушку
	// В реальном проекте здесь был бы Use Case для получения ссылок пользователя
	response := map[string]interface{}{
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
		"links":   []interface{}{}, // Пустой список для примера
		"total":   0,
	}

	c.writeJSONResponse(w, response, http.StatusOK)
}

// RegisterRoutes регистрирует маршруты контроллера
func (c *LinkController) RegisterRoutes(mux *http.ServeMux) {
	// Защищенные маршруты (требуют аутентификации)
	mux.HandleFunc("POST /links", web.ChainMiddleware(
		c.CreateLink,
		c.authMiddleware.RequireAuth, // Требуем аутентификации
		web.CORSMiddleware,
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
	))

	mux.HandleFunc("GET /links", web.ChainMiddleware(
		c.GetUserLinks,
		c.authMiddleware.RequireAuth, // Требуем аутентификации
		web.CORSMiddleware,
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
	))

	// Публичные маршруты (не требуют аутентификации)
	mux.HandleFunc("GET /", web.ChainMiddleware(
		c.Redirect,
		web.CORSMiddleware,
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
	))
}

// getClientIP извлекает IP адрес клиента из HTTP запроса
func (c *LinkController) getClientIP(r *http.Request) string {
	// Проверяем заголовки прокси
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// Берем первый IP из списка
		if parts := strings.Split(ip, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}

	// Возвращаем RemoteAddr как fallback
	if ip := strings.Split(r.RemoteAddr, ":"); len(ip) > 0 {
		return ip[0]
	}

	return r.RemoteAddr
}

// writeErrorResponse записывает ошибку в HTTP ответ
func (c *LinkController) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
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
func (c *LinkController) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// ВАЖНЫЕ ПРИНЦИПЫ:
// 1. Контроллер НЕ содержит бизнес-логики
// 2. Все данные извлекаются из HTTP запроса и передаются Use Case
// 3. Аутентификация обрабатывается через middleware
// 4. Валидация происходит на уровне HTTP, но бизнес-правила в Use Case
// 5. Обработка ошибок стандартизирована