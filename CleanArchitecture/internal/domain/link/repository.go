package link

import "context"

// Repository определяет интерфейс для работы со ссылками в базе данных
type Repository interface {
	// Save сохраняет ссылку в базе данных
	// Если ID = 0, создает новую ссылку и устанавливает ID
	// Если ID != 0, обновляет существующую ссылку
	Save(ctx context.Context, link *Link) error

	// FindByID находит ссылку по ID
	FindByID(ctx context.Context, id uint) (*Link, error)

	// FindByShortCode находит ссылку по короткому коду
	// Это основной метод для редиректа по короткой ссылке
	FindByShortCode(ctx context.Context, shortCode string) (*Link, error)

	// FindByUserID находит все ссылки пользователя с пагинацией
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]*Link, error)

	// Delete удаляет ссылку по ID
	Delete(ctx context.Context, id uint) error

	// ExistsByShortCode проверяет существование ссылки с указанным коротким кодом
	// Используется для проверки уникальности короткого кода
	ExistsByShortCode(ctx context.Context, shortCode string) (bool, error)

	// CountByUserID возвращает общее количество ссылок пользователя
	CountByUserID(ctx context.Context, userID uint) (int, error)

	// IncrementClicks увеличивает счетчик переходов по ссылке
	// Это может быть отдельной операцией для оптимизации производительности
	IncrementClicks(ctx context.Context, linkID uint) error
}

// ПРИМЕЧАНИЕ: Этот интерфейс находится в domain слое, потому что:
// 1. Определяет правила работы с данными для доменной модели Link
// 2. Use Case слой будет использовать этот интерфейс
// 3. Infrastructure слой будет реализовывать этот интерфейс
// 4. Обеспечивает инверсию зависимостей