package user

import (
	"errors"
	"regexp"
	"time"
)

// User представляет доменную модель пользователя
// ЭТО ЧИСТАЯ ДОМЕННАЯ СУЩНОСТЬ - НЕТ ЗАВИСИМОСТЕЙ ОТ ВНЕШНИХ БИБЛИОТЕК!
// Содержит только бизнес-логику и правила домена
type User struct {
	// ID - уникальный идентификатор пользователя
	ID uint

	// Email - электронная почта (должна быть уникальной)
	Email string

	// HashedPassword - захешированный пароль (не храним пароль в открытом виде)
	HashedPassword string

	// Name - имя пользователя
	Name string

	// CreatedAt - время создания аккаунта
	CreatedAt time.Time

	// UpdatedAt - время последнего обновления
	UpdatedAt time.Time
}

// Доменные ошибки - определяем их в доменном слое
var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrEmptyName       = errors.New("name cannot be empty")
	ErrEmptyPassword   = errors.New("password cannot be empty")
	ErrPasswordTooWeak = errors.New("password must be at least 8 characters long")
)

// NewUser - конструктор для создания нового пользователя
// Содержит валидацию доменных правил
func NewUser(email, name string) (*User, error) {
	// Валидация email согласно доменным правилам
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	// Валидация имени согласно доменным правилам
	if name == "" {
		return nil, ErrEmptyName
	}

	now := time.Now()

	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetPassword устанавливает захешированный пароль
// Принимает уже захешированный пароль (хеширование - это техническая деталь,
// которая должна быть в infrastructure слое)
func (u *User) SetPassword(hashedPassword string) error {
	if hashedPassword == "" {
		return ErrEmptyPassword
	}

	u.HashedPassword = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateName обновляет имя пользователя
func (u *User) UpdateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}

	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail обновляет email пользователя
func (u *User) UpdateEmail(email string) error {
	if !isValidEmail(email) {
		return ErrInvalidEmail
	}

	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// isValidEmail проверяет валидность email согласно доменным правилам
func isValidEmail(email string) bool {
	// Простая регулярка для проверки email
	// В реальном проекте можно использовать более сложную валидацию
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}