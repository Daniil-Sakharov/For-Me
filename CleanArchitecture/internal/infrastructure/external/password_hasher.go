package external

import (
	"clean-url-shortener/internal/usecase/auth"
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher реализует интерфейс PasswordHasher используя bcrypt
// ЭТО РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА ИЗ USECASE СЛОЯ В INFRASTRUCTURE СЛОЕ
type BcryptPasswordHasher struct {
	cost int // Стоимость хеширования (чем больше, тем медленнее но безопаснее)
}

// NewBcryptPasswordHasher создает новый hasher с заданной стоимостью
func NewBcryptPasswordHasher(cost int) auth.PasswordHasher {
	// Если стоимость не задана, используем стандартную
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	
	return &BcryptPasswordHasher{
		cost: cost,
	}
}

// Hash хеширует пароль используя bcrypt
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	
	return string(hashedBytes), nil
}

// Compare сравнивает пароль с хешем
func (h *BcryptPasswordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ПРИНЦИПЫ РЕАЛИЗАЦИИ:
// 1. Реализует интерфейс из usecase слоя
// 2. Использует конкретную библиотеку (bcrypt)
// 3. Скрывает технические детали от верхних слоев
// 4. Может быть легко заменена на другую реализацию (например, argon2)