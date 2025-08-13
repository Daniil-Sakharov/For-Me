package order

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// 🏛️ ИНТЕРФЕЙС РЕПОЗИТОРИЯ ЗАКАЗОВ
// Определяет контракт для работы с хранилищем заказов
// Находится в ДОМЕННОМ слое, но РЕАЛИЗУЕТСЯ в инфраструктурном

// Repository определяет методы для работы с заказами в хранилище
type Repository interface {
	// 💾 ОСНОВНЫЕ CRUD ОПЕРАЦИИ
	
	// Save сохраняет новый заказ или обновляет существующий
	Save(ctx context.Context, order *Order) error
	
	// GetByID получает заказ по ID
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	
	// GetByCustomerID получает все заказы клиента
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*Order, error)
	
	// Delete удаляет заказ (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error
	
	// 🔍 ПОИСК И ФИЛЬТРАЦИЯ
	
	// GetByStatus получает заказы по статусу
	GetByStatus(ctx context.Context, status Status) ([]*Order, error)
	
	// GetByStatusAndPriority получает заказы по статусу и приоритету
	GetByStatusAndPriority(ctx context.Context, status Status, priority Priority) ([]*Order, error)
	
	// GetPendingOrders получает все ожидающие обработки заказы
	GetPendingOrders(ctx context.Context) ([]*Order, error)
	
	// GetActiveOrders получает все активные заказы (не отмененные и не доставленные)
	GetActiveOrders(ctx context.Context) ([]*Order, error)
	
	// 📊 АНАЛИТИЧЕСКИЕ ЗАПРОСЫ
	
	// GetOrdersInDateRange получает заказы в диапазоне дат
	GetOrdersInDateRange(ctx context.Context, from, to time.Time) ([]*Order, error)
	
	// CountByStatus возвращает количество заказов по статусам
	CountByStatus(ctx context.Context) (map[Status]int, error)
	
	// GetTotalRevenueByPeriod возвращает общую выручку за период
	GetTotalRevenueByPeriod(ctx context.Context, from, to time.Time) (Money, error)
	
	// 🔄 ПАКЕТНЫЕ ОПЕРАЦИИ
	
	// SaveBatch сохраняет несколько заказов в одной транзакции
	SaveBatch(ctx context.Context, orders []*Order) error
	
	// UpdateStatusBatch обновляет статус у нескольких заказов
	UpdateStatusBatch(ctx context.Context, orderIDs []uuid.UUID, newStatus Status) error
}

// 🎯 СПЕЦИАЛИЗИРОВАННЫЕ РЕПОЗИТОРИИ

// QueryRepository определяет сложные запросы для аналитики
type QueryRepository interface {
	// GetTopCustomersByRevenue получает топ клиентов по выручке
	GetTopCustomersByRevenue(ctx context.Context, limit int, from, to time.Time) ([]CustomerRevenue, error)
	
	// GetPopularProducts получает популярные товары
	GetPopularProducts(ctx context.Context, limit int, from, to time.Time) ([]ProductSales, error)
	
	// GetOrderTrends получает тренды заказов по дням/месяцам
	GetOrderTrends(ctx context.Context, from, to time.Time, groupBy TrendGroupBy) ([]OrderTrend, error)
	
	// GetAverageOrderValue вычисляет средний чек
	GetAverageOrderValue(ctx context.Context, from, to time.Time) (Money, error)
}

// 📊 ВСПОМОГАТЕЛЬНЫЕ ТИПЫ ДЛЯ АНАЛИТИКИ

// CustomerRevenue представляет выручку от клиента
type CustomerRevenue struct {
	CustomerID    uuid.UUID
	TotalRevenue  Money
	OrdersCount   int
	AverageOrder  Money
}

// ProductSales представляет продажи товара
type ProductSales struct {
	ProductID     uuid.UUID
	TotalQuantity int
	TotalRevenue  Money
	OrdersCount   int
}

// OrderTrend представляет тренд заказов
type OrderTrend struct {
	Period       time.Time
	OrdersCount  int
	TotalRevenue Money
	AverageOrder Money
}

// TrendGroupBy определяет группировку для трендов
type TrendGroupBy string

const (
	TrendGroupByDay   TrendGroupBy = "day"
	TrendGroupByWeek  TrendGroupBy = "week"
	TrendGroupByMonth TrendGroupBy = "month"
	TrendGroupByYear  TrendGroupBy = "year"
)

// 🔍 ИНТЕРФЕЙСЫ ДЛЯ ПОИСКА

// SearchCriteria определяет критерии поиска заказов
type SearchCriteria struct {
	CustomerID    *uuid.UUID
	Status        *Status
	Priority      *Priority
	Source        *Source
	MinAmount     *Money
	MaxAmount     *Money
	DateFrom      *time.Time
	DateTo        *time.Time
	ProductID     *uuid.UUID
	ShippingCity  *string
	Limit         int
	Offset        int
	SortBy        SortField
	SortDirection SortDirection
}

// SortField определяет поле для сортировки
type SortField string

const (
	SortByCreatedAt   SortField = "created_at"
	SortByUpdatedAt   SortField = "updated_at"
	SortByTotalAmount SortField = "total_amount"
	SortByStatus      SortField = "status"
	SortByPriority    SortField = "priority"
)

// SortDirection определяет направление сортировки
type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// SearchRepository определяет расширенный поиск
type SearchRepository interface {
	// Search выполняет поиск по критериям
	Search(ctx context.Context, criteria SearchCriteria) ([]*Order, error)
	
	// Count возвращает количество заказов по критериям
	Count(ctx context.Context, criteria SearchCriteria) (int, error)
}