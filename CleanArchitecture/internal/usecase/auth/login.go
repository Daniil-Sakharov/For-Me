package auth

import (
	"context"
	"errors"
	"time"
	"clean-url-shortener/internal/domain/user"
	"clean-url-shortener/pkg/logger"
)

// LoginRequest представляет запрос на авторизацию
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse представляет ответ при авторизации
type LoginResponse struct {
	Token  string `json:"token"`
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// Ошибки для логина
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound      = errors.New("user not found")
)

// LoginUseCase содержит бизнес-логику авторизации пользователя
// 📚 ИЗУЧИТЬ: Logging в Use Cases
// https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
// 
// ПРИНЦИПЫ ЛОГИРОВАНИЯ В USE CASES:
// 1. Логируем БИЗНЕС-СОБЫТИЯ, а не технические детали
// 2. Используем структурированное логирование
// 3. НЕ логируем чувствительные данные (пароли, токены)
// 4. Логируем начало и конец операций
// 5. Логируем ошибки с контекстом
type LoginUseCase struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
	logger         logger.Logger  // ← ДОБАВИЛИ логгер
}

// NewLoginUseCase создает новый Use Case для авторизации
func NewLoginUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
	logger logger.Logger,  // ← ДОБАВИЛИ параметр логгера
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
		logger:         logger.UseCaseLogger(), // ← Создаем логгер для Use Case слоя
	}
}

// Execute выполняет бизнес-логику авторизации пользователя
func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// ЛОГИРОВАНИЕ: Создаем контекстный логгер для этой операции
	// 📚 ИЗУЧИТЬ: Structured logging best practices
	// https://betterstack.com/community/guides/logging/logging-in-go/
	operationLogger := uc.logger.WithContext(ctx).With(
		"operation", "user_login",
		"email", logger.SafeEmail(req.Email), // Безопасно логируем email
		// ВАЖНО: НЕ логируем пароль! Это критическая ошибка безопасности
	)
	
	// Засекаем время выполнения операции
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		operationLogger.Debug("Login operation completed", 
			"duration_ms", duration.Milliseconds())
	}()
	
	operationLogger.Info("Login attempt started")
	
	// ШАГ 1: Находим пользователя по email
	operationLogger.Debug("Looking up user by email")
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		operationLogger.Error("Failed to find user by email", 
			"error", err,
			"step", "user_lookup")
		return nil, err
	}
	
	if existingUser == nil {
		// БЕЗОПАСНОСТЬ: Логируем попытку входа несуществующего пользователя
		// Это может указывать на попытки подбора аккаунтов
		operationLogger.Warn("Login attempt for non-existent user",
			"step", "user_lookup")
		return nil, ErrUserNotFound
	}
	
	operationLogger.Debug("User found", "user_id", existingUser.ID)

	// ШАГ 2: Проверяем пароль
	operationLogger.Debug("Verifying password")
	err = uc.passwordHasher.Compare(existingUser.HashedPassword, req.Password)
	if err != nil {
		// БЕЗОПАСНОСТЬ: Логируем неудачную попытку входа
		// Это может указывать на попытки подбора паролей
		operationLogger.Warn("Invalid password provided",
			"user_id", existingUser.ID,
			"step", "password_verification")
		
		// Не раскрываем информацию о том, что именно неверно (email или пароль)
		// Это важно для безопасности
		return nil, ErrInvalidCredentials
	}
	
	operationLogger.Debug("Password verified successfully")

	// ШАГ 3: Генерируем JWT токен
	operationLogger.Debug("Generating authentication token")
	token, err := uc.tokenGenerator.Generate(existingUser.ID, existingUser.Email)
	if err != nil {
		operationLogger.Error("Failed to generate authentication token",
			"error", err,
			"user_id", existingUser.ID,
			"step", "token_generation")
		return nil, err
	}
	
	// БЕЗОПАСНОСТЬ: НЕ логируем сам токен!
	operationLogger.Debug("Authentication token generated successfully")

	// ШАГ 4: Возвращаем успешный ответ
	operationLogger.Info("Login successful",
		"user_id", existingUser.ID,
		"user_name", existingUser.Name)
	
	return &LoginResponse{
		Token:  token,
		UserID: existingUser.ID,
		Email:  existingUser.Email,
		Name:   existingUser.Name,
	}, nil
}

// ПРИНЦИПЫ БЕЗОПАСНОСТИ:
// 1. Не раскрываем, существует ли пользователь с таким email
// 2. Используем константное время для проверки пароля
// 3. Генерируем безопасные JWT токены
// 4. Не логируем чувствительные данные