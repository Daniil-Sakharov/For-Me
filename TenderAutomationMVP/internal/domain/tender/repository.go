// =====================================================================
// 🏛️ ИНТЕРФЕЙС REPOSITORY ДЛЯ TENDER - Контракт доступа к данным
// =====================================================================
//
// Этот файл определяет интерфейс TenderRepository в доменном слое.
// Следует принципам Clean Architecture:
// 1. Интерфейс определяется в доменном слое (внутренний)
// 2. Реализация будет в infrastructure слое (внешний)
// 3. Dependency Inversion - внешние слои зависят от внутренних
// 4. Интерфейс использует только доменные типы
//
// TODO: При расширении функциональности добавить:
// - Методы для работы с агрегатами (тендер + документы + товары)
// - Методы для bulk операций (создание множественных тендеров)
// - Методы для статистики и аналитики
// - Методы для full-text поиска

package tender

import (
	"context"
	"time"
)

// =====================================================================
// 🗃️ ОСНОВНОЙ ИНТЕРФЕЙС REPOSITORY
// =====================================================================

// TenderRepository определяет контракт для работы с тендерами в хранилище
//
// TODO: Реализовать следующие методы в infrastructure слое:
// 1. Все CRUD операции с proper error handling
// 2. Эффективные запросы с индексами
// 3. Транзакционность для сложных операций
// 4. Пагинация для больших объемов данных
// 5. Кеширование часто запрашиваемых данных
//
// ВАЖНО: Все методы должны:
// - Принимать context.Context для таймаутов и отмены
// - Возвращать доменные сущности, а не DTO или database модели
// - Использовать доменные ошибки (из errors.go)
// - Быть готовыми к concurrent access
type TenderRepository interface {
	// =====================================================================
	// 📝 ОСНОВНЫЕ CRUD ОПЕРАЦИИ
	// =====================================================================

	// Create создает новый тендер в хранилище
	//
	// TODO: Реализовать в PostgreSQL repository:
	// - INSERT запрос с RETURNING для получения ID
	// - Проверка уникальности external_id
	// - Валидация данных на уровне БД
	// - Proper error mapping (constraint violations → domain errors)
	//
	// Параметры:
	//   - ctx: контекст для таймаутов и отмены операции
	//   - tender: доменная сущность для создания
	//
	// Возвращает:
	//   - error: ErrDuplicateTender если external_id уже существует
	//
	// Пример использования:
	//   tender, _ := tender.NewTender("123", "Test tender", "zakupki", "http://...")
	//   err := repo.Create(ctx, tender)
	Create(ctx context.Context, tender *Tender) error

	// GetByID получает тендер по внутреннему ID
	//
	// TODO: Реализовать в PostgreSQL repository:
	// - SELECT запрос с оптимизированными индексами
	// - Маппинг database модели в доменную сущность
	// - Обработка NULL значений для опциональных полей
	//
	// Параметры:
	//   - ctx: контекст операции
	//   - id: внутренний ID тендера
	//
	// Возвращает:
	//   - *Tender: найденный тендер
	//   - error: ErrTenderNotFound если не найден
	GetByID(ctx context.Context, id uint) (*Tender, error)

	// GetByExternalID получает тендер по внешнему ID (с платформы)
	//
	// TODO: Реализовать эффективный поиск по unique индексу
	// TODO: Добавить кеширование для часто запрашиваемых тендеров
	//
	// Параметры:
	//   - ctx: контекст операции  
	//   - externalID: ID тендера на внешней платформе
	//
	// Возвращает:
	//   - *Tender: найденный тендер
	//   - error: ErrTenderNotFound если не найден
	GetByExternalID(ctx context.Context, externalID string) (*Tender, error)

	// Update обновляет существующий тендер
	//
	// TODO: Реализовать в PostgreSQL repository:
	// - UPDATE запрос с WHERE по ID
	// - Автоматическое обновление updated_at через триггер
	// - Оптимистичная блокировка для concurrent updates
	// - Валидация что тендер существует
	//
	// Параметры:
	//   - ctx: контекст операции
	//   - tender: обновленная доменная сущность
	//
	// Возвращает:
	//   - error: ErrTenderNotFound если тендер не существует
	Update(ctx context.Context, tender *Tender) error

	// Delete удаляет тендер (soft delete)
	//
	// TODO: Реализовать soft delete через deleted_at поле
	// TODO: Добавить каскадное удаление связанных сущностей
	//
	// Параметры:
	//   - ctx: контекст операции
	//   - id: ID тендера для удаления
	//
	// Возвращает:
	//   - error: ErrTenderNotFound если тендер не существует
	Delete(ctx context.Context, id uint) error

	// =====================================================================
	// 🔍 МЕТОДЫ ПОИСКА И ФИЛЬТРАЦИИ
	// =====================================================================

	// List возвращает список тендеров с фильтрацией и пагинацией
	//
	// TODO: Реализовать эффективную пагинацию и фильтрацию:
	// - WHERE условия на основе фильтров
	// - ORDER BY для сортировки
	// - LIMIT/OFFSET для пагинации
	// - COUNT(*) для общего количества (если нужно)
	//
	// Параметры:
	//   - ctx: контекст операции
	//   - filters: фильтры для поиска
	//
	// Возвращает:
	//   - []*Tender: список найденных тендеров
	//   - error: ошибка выполнения запроса
	List(ctx context.Context, filters TenderFilters) ([]*Tender, error)

	// GetPendingAnalysis возвращает тендеры, ожидающие AI анализа
	//
	// TODO: Реализовать запрос:
	// - WHERE ai_analyzed_at IS NULL
	// - AND status = 'active'  
	// - ORDER BY created_at ASC (старые сначала)
	// - LIMIT для batch processing
	//
	// Параметры:
	//   - ctx: контекст операции
	//   - limit: максимальное количество тендеров
	//
	// Возвращает:
	//   - []*Tender: тендеры для анализа
	//   - error: ошибка выполнения
	GetPendingAnalysis(ctx context.Context, limit int) ([]*Tender, error)

	// GetExpiredTenders возвращает тендеры с истекшим сроком подачи заявок
	//
	// TODO: Реализовать для автоматического обновления статусов:
	// - WHERE deadline_at < NOW()
	// - AND status = 'active'
	// - Batch update статусов на 'expired'
	//
	// Параметры:
	//   - ctx: контекст операции
	//
	// Возвращает:
	//   - []*Tender: истекшие тендеры
	//   - error: ошибка выполнения
	GetExpiredTenders(ctx context.Context) ([]*Tender, error)

	// =====================================================================
	// 📊 МЕТОДЫ ДЛЯ АНАЛИТИКИ И СТАТИСТИКИ
	// =====================================================================

	// GetStatistics возвращает общую статистику по тендерам
	//
	// TODO: Реализовать агрегирующие запросы:
	// - COUNT по статусам
	// - COUNT по платформам  
	// - AVG(ai_score) для проанализированных
	// - SUM(start_price) по валютам
	//
	// Параметры:
	//   - ctx: контекст операции
	//   - period: период для статистики
	//
	// Возвращает:
	//   - TenderStatistics: структура со статистикой
	//   - error: ошибка выполнения
	GetStatistics(ctx context.Context, period StatisticsPeriod) (*TenderStatistics, error)

	// CountByStatus возвращает количество тендеров по статусам
	//
	// TODO: Реализовать простой GROUP BY запрос
	//
	// Параметры:
	//   - ctx: контекст операции
	//
	// Возвращает:
	//   - map[TenderStatus]int: количество по каждому статусу
	//   - error: ошибка выполнения
	CountByStatus(ctx context.Context) (map[TenderStatus]int, error)

	// =====================================================================
	// 🔄 BATCH ОПЕРАЦИИ (для производительности)
	// =====================================================================

	// CreateBatch создает множественные тендеры за одну операцию
	//
	// TODO: Реализовать для эффективного импорта:
	// - INSERT с multiple VALUES
	// - Транзакция для atomicity
	// - Batch size limits для больших объемов
	// - Proper error handling для partial failures
	//
	// Параметры:
	//   - ctx: контекст операции
	//   - tenders: список тендеров для создания
	//
	// Возвращает:
	//   - error: ошибка создания (может содержать partial failures)
	CreateBatch(ctx context.Context, tenders []*Tender) error

	// UpdateStatusBatch обновляет статус множественных тендеров
	//
	// TODO: Реализовать для массовых обновлений:
	// - UPDATE WHERE id IN (...)
	// - Транзакция для consistency
	//
	// Параметры:
	//   - ctx: контекст операции
	//   - ids: список ID тендеров
	//   - status: новый статус
	//
	// Возвращает:
	//   - error: ошибка обновления
	UpdateStatusBatch(ctx context.Context, ids []uint, status TenderStatus) error
}

// =====================================================================
// 📋 СТРУКТУРЫ ДЛЯ ФИЛЬТРАЦИИ И ПАГИНАЦИИ
// =====================================================================

// TenderFilters определяет параметры для фильтрации списка тендеров
//
// TODO: Добавить дополнительные фильтры по мере необходимости:
// - PriceRange (min/max price)
// - DeadlineRange (from/to dates)
// - CustomerFilter (по заказчику)
// - FullTextSearch (поиск по тексту)
type TenderFilters struct {
	// 🔍 Основные фильтры
	Status    *TenderStatus // Фильтр по статусу (optional)
	Platform  *string       // Фильтр по платформе (optional)
	Category  *string       // Фильтр по категории (optional)

	// 🤖 AI фильтры  
	MinAIScore       *float64         // Минимальная AI оценка
	AIRecommendation *AIRecommendation // Фильтр по AI рекомендации

	// 📅 Временные фильтры
	CreatedAfter  *time.Time // Созданы после даты
	CreatedBefore *time.Time // Созданы до даты
	DeadlineAfter *time.Time // Дедлайн после даты

	// 📊 Пагинация
	Limit  int // Максимальное количество результатов (default: 20)
	Offset int // Смещение для пагинации (default: 0)

	// 🔄 Сортировка
	SortBy    string // Поле для сортировки (default: "created_at")
	SortOrder string // Порядок: "asc" или "desc" (default: "desc")
}

// =====================================================================
// 📊 СТРУКТУРЫ ДЛЯ СТАТИСТИКИ
// =====================================================================

// StatisticsPeriod определяет период для статистики
type StatisticsPeriod string

const (
	PeriodToday     StatisticsPeriod = "today"      // За сегодня
	PeriodWeek      StatisticsPeriod = "week"       // За неделю
	PeriodMonth     StatisticsPeriod = "month"      // За месяц
	PeriodQuarter   StatisticsPeriod = "quarter"    // За квартал
	PeriodYear      StatisticsPeriod = "year"       // За год
	PeriodAllTime   StatisticsPeriod = "all_time"   // За все время
)

// TenderStatistics содержит агрегированную статистику по тендерам
//
// TODO: Расширить дополнительными метриками по мере необходимости
type TenderStatistics struct {
	// 📊 Общие метрики
	TotalCount    int                        // Общее количество тендеров
	StatusCounts  map[TenderStatus]int       // Количество по статусам
	PlatformCounts map[string]int            // Количество по платформам

	// 🤖 AI метрики
	AnalyzedCount     int     // Количество проанализированных
	AverageAIScore    float64 // Средняя AI оценка
	RelevantCount     int     // Количество релевантных (score >= 0.7)
	RecommendedCount  int     // Количество с рекомендацией "participate"

	// 💰 Финансовые метрики
	TotalValue        float64            // Общая стоимость тендеров
	AverageValue      float64            // Средняя стоимость
	CurrencyBreakdown map[Currency]float64 // Разбивка по валютам

	// 📅 Временные метрики  
	Period            StatisticsPeriod // Период статистики
	OldestTender      *time.Time       // Дата самого старого тендера
	NewestTender      *time.Time       // Дата самого нового тендера
}

// =====================================================================
// 🔧 ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ ДЛЯ ФИЛЬТРОВ
// =====================================================================

// NewTenderFilters создает новые фильтры с дефолтными значениями
//
// TODO: Добавить валидацию параметров фильтрации
func NewTenderFilters() TenderFilters {
	return TenderFilters{
		Limit:     20,    // Дефолтный размер страницы
		Offset:    0,     // Начинаем с первой страницы
		SortBy:    "created_at", // Сортируем по дате создания
		SortOrder: "desc", // Новые сначала
	}
}

// WithStatus добавляет фильтр по статусу
func (f TenderFilters) WithStatus(status TenderStatus) TenderFilters {
	f.Status = &status
	return f
}

// WithPlatform добавляет фильтр по платформе
func (f TenderFilters) WithPlatform(platform string) TenderFilters {
	f.Platform = &platform
	return f
}

// WithPagination устанавливает параметры пагинации
func (f TenderFilters) WithPagination(limit, offset int) TenderFilters {
	f.Limit = limit
	f.Offset = offset
	return f
}

// WithAIScoreFilter добавляет фильтр по AI оценке
func (f TenderFilters) WithAIScoreFilter(minScore float64) TenderFilters {
	f.MinAIScore = &minScore
	return f
}

// =====================================================================
// 📝 ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ REPOSITORY
// =====================================================================

// TODO: Примеры использования в use cases:
//
// 1. Создание тендера:
//    tender, _ := tender.NewTender("123", "Test", "zakupki", "http://...")
//    err := repo.Create(ctx, tender)
//
// 2. Поиск тендеров с фильтрацией:
//    filters := NewTenderFilters().
//        WithStatus(StatusActive).
//        WithPlatform("zakupki").
//        WithPagination(20, 0)
//    tenders, err := repo.List(ctx, filters)
//
// 3. Получение тендеров для анализа:
//    tenders, err := repo.GetPendingAnalysis(ctx, 10)
//
// 4. Статистика:
//    stats, err := repo.GetStatistics(ctx, PeriodMonth)
//
// 5. Обновление статуса:
//    tender.UpdateStatus(StatusExpired)
//    err := repo.Update(ctx, tender)

// =====================================================================
// ✅ ЗАКЛЮЧЕНИЕ
// =====================================================================

// Этот интерфейс определяет полный контракт для работы с тендерами.
// При реализации в infrastructure слое:
// 1. Используйте эффективные SQL запросы с правильными индексами
// 2. Реализуйте proper error handling и mapping к доменным ошибкам
// 3. Добавьте транзакционность для сложных операций
// 4. Рассмотрите кеширование для часто запрашиваемых данных
// 5. Обеспечьте thread safety для concurrent access
//
// Помните: интерфейс принадлежит доменному слою и не должен изменяться
// из-за технических ограничений внешних систем.