package auth

import (
	"context"
	"errors"
	"time"
	"clean-url-shortener/internal/domain/user"
	"clean-url-shortener/pkg/logger"
)

// СТРУКТУРЫ ДЛЯ ЗАПРОСОВ И ОТВЕТОВ

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

// RegisterResponse представляет ответ при регистрации
type RegisterResponse struct {
	Token  string `json:"token"`
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// ДОМЕННЫЕ ОШИБКИ ДЛЯ USE CASE

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrWeakPassword      = errors.New("password is too weak")
)

// RegisterUseCase содержит бизнес-логику регистрации пользователя
// ЭТО ОСНОВНОЙ СЦЕНАРИЙ ИСПОЛЬЗОВАНИЯ (USE CASE)
// 📚 ИЗУЧИТЬ: Logging в Clean Architecture
// https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
type RegisterUseCase struct {
	// Зависимости инжектируются через интерфейсы
	userRepo       user.Repository  // Из domain слоя
	passwordHasher PasswordHasher   // Из usecase слоя
	tokenGenerator TokenGenerator   // Из usecase слоя
	logger         logger.Logger    // ← ДОБАВИЛИ логгер
}

// NewRegisterUseCase создает новый Use Case для регистрации
// Все зависимости передаются через интерфейсы (Dependency Injection)
func NewRegisterUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
	logger logger.Logger,  // ← ДОБАВИЛИ параметр логгера
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
		logger:         logger.UseCaseLogger(), // ← Создаем логгер для Use Case слоя
	}
}

// Execute выполняет бизнес-логику регистрации пользователя
// Это ОРКЕСТРАТОР - координирует работу между доменными объектами и сервисами
func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// ЛОГИРОВАНИЕ: Создаем контекстный логгер для операции регистрации
	operationLogger := uc.logger.WithContext(ctx).With(
		"operation", "user_registration",
		"email", logger.SafeEmail(req.Email),
		"name", req.Name,
		// ВАЖНО: НЕ логируем пароль!
	)
	
	// Засекаем время выполнения
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		operationLogger.Debug("Registration operation completed", 
			"duration_ms", duration.Milliseconds())
	}()
	
	operationLogger.Info("User registration started")

	// ШАГ 1: Проверяем, что пользователь с таким email не существует
	operationLogger.Debug("Checking if user already exists")
	exists, err := uc.userRepo.Exists(ctx, req.Email)
	if err != nil {
		operationLogger.Error("Failed to check user existence", 
			"error", err,
			"step", "user_existence_check")
		return nil, err
	}
	
	if exists {
		operationLogger.Warn("Registration attempt for existing user",
			"step", "user_existence_check")
		return nil, ErrUserAlreadyExists
	}
	
	operationLogger.Debug("User does not exist, proceeding with registration")

	// ШАГ 2: Валидируем силу пароля (доменное правило)
	operationLogger.Debug("Validating password strength")
	if len(req.Password) < 8 {
		operationLogger.Warn("Weak password provided",
			"password_length", len(req.Password),
			"step", "password_validation")
		return nil, ErrWeakPassword
	}

	// ШАГ 3: Создаем доменную модель пользователя
	operationLogger.Debug("Creating domain user entity")
	newUser, err := user.NewUser(req.Email, req.Name)
	if err != nil {
		operationLogger.Error("Failed to create domain user entity", 
			"error", err,
			"step", "domain_entity_creation")
		return nil, err
	}
	
	operationLogger.Debug("Domain user entity created successfully")

	// ШАГ 4: Хешируем пароль через внешний сервис
	operationLogger.Debug("Hashing user password")
	hashedPassword, err := uc.passwordHasher.Hash(req.Password)
	if err != nil {
		operationLogger.Error("Failed to hash password", 
			"error", err,
			"step", "password_hashing")
		return nil, err
	}
	
	operationLogger.Debug("Password hashed successfully")

	// ШАГ 5: Устанавливаем захешированный пароль
	operationLogger.Debug("Setting hashed password to user entity")
	err = newUser.SetPassword(hashedPassword)
	if err != nil {
		operationLogger.Error("Failed to set hashed password", 
			"error", err,
			"step", "password_setting")
		return nil, err
	}

	// ШАГ 6: Сохраняем пользователя в базе данных
	operationLogger.Debug("Saving user to database")
	err = uc.userRepo.Save(ctx, newUser)
	if err != nil {
		operationLogger.Error("Failed to save user to database", 
			"error", err,
			"step", "user_persistence")
		return nil, err
	}
	
	operationLogger.Info("User saved to database successfully", 
		"user_id", newUser.ID)

	// ШАГ 7: Генерируем JWT токен для автоматического входа
	operationLogger.Debug("Generating authentication token")
	token, err := uc.tokenGenerator.Generate(newUser.ID, newUser.Email)
	if err != nil {
		operationLogger.Error("Failed to generate authentication token", 
			"error", err,
			"user_id", newUser.ID,
			"step", "token_generation")
		return nil, err
	}
	
	// БЕЗОПАСНОСТЬ: НЕ логируем сам токен!
	operationLogger.Debug("Authentication token generated successfully")

	// ШАГ 8: Возвращаем ответ
	operationLogger.Info("User registration completed successfully",
		"user_id", newUser.ID,
		"user_name", newUser.Name)
	
	return &RegisterResponse{
		Token:  token,
		UserID: newUser.ID,
		Email:  newUser.Email,
		Name:   newUser.Name,
	}, nil
}

// ПРИНЦИПЫ, КОТОРЫЕ МЫ СОБЛЮДАЕМ:
// 1. SINGLE RESPONSIBILITY - Use Case отвечает только за один сценарий
// 2. DEPENDENCY INVERSION - зависим от интерфейсов, а не от реализаций
// 3. OPEN/CLOSED - можем расширять функциональность, не изменяя код
// 4. SEPARATION OF CONCERNS - разделение ответственности между слоями