package stat

import "context"

// Repository определяет интерфейс для работы со статистикой в базе данных
type Repository interface {
	// Save сохраняет запись статистики в базе данных
	Save(ctx context.Context, stat *Stat) error

	// FindByLinkID находит все записи статистики для указанной ссылки
	// с возможностью фильтрации по периоду времени
	FindByLinkID(ctx context.Context, linkID uint, period Period) ([]*Stat, error)

	// GetLinkStatistics возвращает агрегированную статистику для ссылки
	// за указанный период времени
	GetLinkStatistics(ctx context.Context, linkID uint, period Period) (*LinkStatistics, error)

	// GetUserLinkStatistics возвращает статистику для всех ссылок пользователя
	GetUserLinkStatistics(ctx context.Context, userID uint, period Period) ([]*LinkStatistics, error)

	// CountTotalClicks возвращает общее количество кликов по ссылке
	CountTotalClicks(ctx context.Context, linkID uint) (int, error)

	// CountUniqueClicks возвращает количество уникальных кликов по ссылке
	// (уникальность определяется по IP адресу)
	CountUniqueClicks(ctx context.Context, linkID uint, period Period) (int, error)

	// GetClicksByCountry возвращает статистику кликов по странам
	GetClicksByCountry(ctx context.Context, linkID uint, period Period) (map[string]int, error)

	// GetClicksByHour возвращает статистику кликов по часам
	// для построения временных графиков
	GetClicksByHour(ctx context.Context, linkID uint, period Period) ([]int, error)

	// GetTopReferers возвращает топ источников трафика
	GetTopReferers(ctx context.Context, linkID uint, period Period, limit int) ([]RefererStat, error)

	// DeleteByLinkID удаляет всю статистику для указанной ссылки
	// Используется при удалении ссылки
	DeleteByLinkID(ctx context.Context, linkID uint) error

	// CleanupOldStats удаляет старые записи статистики
	// Может использоваться для очистки данных старше определенного периода
	CleanupOldStats(ctx context.Context, olderThan time.Time) error
}

// ВАЖНО: Интерфейс Repository находится в domain слое для:
// 1. Определения правил работы с данными статистики
// 2. Использования в Use Case слое
// 3. Реализации в Infrastructure слое
// 4. Обеспечения инверсии зависимостей

import "time"