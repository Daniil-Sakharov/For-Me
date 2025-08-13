package auth

// ИНТЕРФЕЙСЫ ДЛЯ ЗАВИСИМОСТЕЙ USE CASES
// Эти интерфейсы определяют контракты для внешних сервисов,
// которые будут реализованы в infrastructure слое

// PasswordHasher определяет интерфейс для хеширования паролей
// Скрывает детали реализации (bcrypt, argon2 и т.д.)
type PasswordHasher interface {
	// Hash хеширует пароль
	Hash(password string) (string, error)
	
	// Compare сравнивает пароль с хешем
	Compare(hashedPassword, password string) error
}

// TokenGenerator определяет интерфейс для генерации JWT токенов
// Скрывает детали реализации JWT
type TokenGenerator interface {
	// Generate создает токен для пользователя
	Generate(userID uint, email string) (string, error)
	
	// Validate проверяет валидность токена и возвращает данные
	Validate(token string) (*TokenData, error)
}

// TokenData содержит данные из токена
type TokenData struct {
	UserID uint
	Email  string
}

// ВАЖНО: Эти интерфейсы находятся в usecase слое, потому что:
// 1. Use Cases определяют ЧТО им нужно от внешних сервисов
// 2. Infrastructure слой будет РЕАЛИЗОВЫВАТЬ эти интерфейсы
// 3. Это обеспечивает инверсию зависимостей
// 4. Use Cases не зависят от конкретных реализаций (bcrypt, JWT библиотеки)