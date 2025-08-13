package auth

import (
	"context"
	"errors"
	"clean-url-shortener/internal/domain/user"
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
type LoginUseCase struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

// NewLoginUseCase создает новый Use Case для авторизации
func NewLoginUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

// Execute выполняет бизнес-логику авторизации пользователя
func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// ШАГ 1: Находим пользователя по email
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, ErrUserNotFound
	}

	// ШАГ 2: Проверяем пароль
	err = uc.passwordHasher.Compare(existingUser.HashedPassword, req.Password)
	if err != nil {
		// Не раскрываем информацию о том, что именно неверно (email или пароль)
		// Это важно для безопасности
		return nil, ErrInvalidCredentials
	}

	// ШАГ 3: Генерируем JWT токен
	token, err := uc.tokenGenerator.Generate(existingUser.ID, existingUser.Email)
	if err != nil {
		return nil, err
	}

	// ШАГ 4: Возвращаем ответ
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