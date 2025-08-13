package auth

import (
	"context"
	"errors"
	"clean-url-shortener/internal/domain/user"
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
type RegisterUseCase struct {
	// Зависимости инжектируются через интерфейсы
	userRepo       user.Repository  // Из domain слоя
	passwordHasher PasswordHasher   // Из usecase слоя
	tokenGenerator TokenGenerator   // Из usecase слоя
}

// NewRegisterUseCase создает новый Use Case для регистрации
// Все зависимости передаются через интерфейсы (Dependency Injection)
func NewRegisterUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

// Execute выполняет бизнес-логику регистрации пользователя
// Это ОРКЕСТРАТОР - координирует работу между доменными объектами и сервисами
func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// ШАГ 1: Проверяем, что пользователь с таким email не существует
	exists, err := uc.userRepo.Exists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// ШАГ 2: Валидируем силу пароля (доменное правило)
	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	// ШАГ 3: Создаем доменную модель пользователя
	// Используем конструктор из domain слоя для валидации доменных правил
	newUser, err := user.NewUser(req.Email, req.Name)
	if err != nil {
		return nil, err
	}

	// ШАГ 4: Хешируем пароль через внешний сервис
	hashedPassword, err := uc.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// ШАГ 5: Устанавливаем захешированный пароль
	err = newUser.SetPassword(hashedPassword)
	if err != nil {
		return nil, err
	}

	// ШАГ 6: Сохраняем пользователя в базе данных
	err = uc.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// ШАГ 7: Генерируем JWT токен для автоматического входа
	token, err := uc.tokenGenerator.Generate(newUser.ID, newUser.Email)
	if err != nil {
		return nil, err
	}

	// ШАГ 8: Возвращаем ответ
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