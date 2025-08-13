package user

import "context"

// Repository определяет интерфейс для работы с пользователями в базе данных
// ЭТО ИНТЕРФЕЙС В ДОМЕННОМ СЛОЕ - он определяет ЧТО нужно делать, но НЕ КАК
// Конкретная реализация будет в infrastructure слое
type Repository interface {
	// Save сохраняет пользователя в базе данных
	// Если ID = 0, создает нового пользователя и устанавливает ID
	// Если ID != 0, обновляет существующего пользователя
	Save(ctx context.Context, user *User) error

	// FindByID находит пользователя по ID
	// Возвращает nil, если пользователь не найден
	FindByID(ctx context.Context, id uint) (*User, error)

	// FindByEmail находит пользователя по email
	// Возвращает nil, если пользователь не найден
	FindByEmail(ctx context.Context, email string) (*User, error)

	// Delete удаляет пользователя по ID
	Delete(ctx context.Context, id uint) error

	// Exists проверяет существование пользователя по email
	// Используется для проверки уникальности email
	Exists(ctx context.Context, email string) (bool, error)
}

// ВАЖНО: Интерфейс Repository находится в domain слое, потому что:
// 1. Доменный слой определяет ПРАВИЛА работы с данными
// 2. Use Case слой будет зависеть от этого интерфейса
// 3. Infrastructure слой будет РЕАЛИЗОВЫВАТЬ этот интерфейс
// 4. Это обеспечивает ИНВЕРСИЮ ЗАВИСИМОСТЕЙ - внутренние слои не зависят от внешних