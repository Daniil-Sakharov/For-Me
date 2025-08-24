package user

import (
	"context"
	"time"
)

// 📚 ИЗУЧИТЬ: Repository Pattern с Pagination и Filtering
// https://martinfowler.com/eaaCatalog/repository.html
// https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer-design
// https://www.postgresql.org/docs/current/queries-limit.html
//
// ПРИНЦИПЫ PAGINATION И FILTERING:
// 1. Pagination должна быть ЭФФЕКТИВНОЙ (используйте OFFSET/LIMIT или cursor-based)
// 2. Фильтры должны быть ТИПИЗИРОВАННЫМИ
// 3. Сортировка должна быть ПРЕДСКАЗУЕМОЙ
// 4. Всегда возвращайте ОБЩЕЕ КОЛИЧЕСТВО для UI
// 5. Валидируйте параметры pagination

// PaginationParams параметры пагинации
type PaginationParams struct {
	// Page номер страницы (начинается с 1)
	Page int `json:"page" validate:"min=1"`
	
	// Limit количество элементов на странице
	Limit int `json:"limit" validate:"min=1,max=100"`
	
	// Offset смещение (вычисляется автоматически)
	Offset int `json:"offset"`
}

// NewPaginationParams создает параметры пагинации с валидацией
func NewPaginationParams(page, limit int) PaginationParams {
	// Устанавливаем значения по умолчанию
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Максимум 100 элементов на странице
	}
	
	offset := (page - 1) * limit
	
	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// SortOrder направление сортировки
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// SortParams параметры сортировки
type SortParams struct {
	// Field поле для сортировки
	Field string `json:"field"`
	
	// Order направление сортировки
	Order SortOrder `json:"order"`
}

// NewSortParams создает параметры сортировки с валидацией
func NewSortParams(field string, order SortOrder) SortParams {
	// Валидация поля сортировки
	validFields := []string{"id", "email", "name", "created_at", "updated_at"}
	isValidField := false
	for _, validField := range validFields {
		if field == validField {
			isValidField = true
			break
		}
	}
	
	if !isValidField {
		field = "created_at" // Поле по умолчанию
	}
	
	// Валидация направления сортировки
	if order != SortOrderAsc && order != SortOrderDesc {
		order = SortOrderDesc // Направление по умолчанию
	}
	
	return SortParams{
		Field: field,
		Order: order,
	}
}

// FilterParams параметры фильтрации пользователей
type FilterParams struct {
	// Email фильтр по email (частичное совпадение)
	Email string `json:"email,omitempty"`
	
	// Name фильтр по имени (частичное совпадение)
	Name string `json:"name,omitempty"`
	
	// CreatedAfter пользователи, созданные после указанной даты
	CreatedAfter *time.Time `json:"created_after,omitempty"`
	
	// CreatedBefore пользователи, созданные до указанной даты
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	
	// UpdatedAfter пользователи, обновленные после указанной даты
	UpdatedAfter *time.Time `json:"updated_after,omitempty"`
	
	// UpdatedBefore пользователи, обновленные до указанной даты
	UpdatedBefore *time.Time `json:"updated_before,omitempty"`
	
	// IDs фильтр по списку ID
	IDs []uint `json:"ids,omitempty"`
	
	// ExcludeIDs исключить пользователей с указанными ID
	ExcludeIDs []uint `json:"exclude_ids,omitempty"`
}

// IsEmpty проверяет, пустые ли фильтры
func (f FilterParams) IsEmpty() bool {
	return f.Email == "" && 
		   f.Name == "" && 
		   f.CreatedAfter == nil && 
		   f.CreatedBefore == nil &&
		   f.UpdatedAfter == nil &&
		   f.UpdatedBefore == nil &&
		   len(f.IDs) == 0 &&
		   len(f.ExcludeIDs) == 0
}

// FindUsersQuery запрос для поиска пользователей с пагинацией и фильтрацией
type FindUsersQuery struct {
	// Pagination параметры пагинации
	Pagination PaginationParams `json:"pagination"`
	
	// Sort параметры сортировки
	Sort SortParams `json:"sort"`
	
	// Filter параметры фильтрации
	Filter FilterParams `json:"filter"`
}

// NewFindUsersQuery создает запрос с параметрами по умолчанию
func NewFindUsersQuery() FindUsersQuery {
	return FindUsersQuery{
		Pagination: NewPaginationParams(1, 10),
		Sort:       NewSortParams("created_at", SortOrderDesc),
		Filter:     FilterParams{},
	}
}

// PaginatedResult результат пагинированного запроса
type PaginatedResult struct {
	// Items найденные пользователи
	Items []*User `json:"items"`
	
	// Total общее количество пользователей (без учета пагинации)
	Total int64 `json:"total"`
	
	// Page текущая страница
	Page int `json:"page"`
	
	// Limit количество элементов на странице
	Limit int `json:"limit"`
	
	// TotalPages общее количество страниц
	TotalPages int `json:"total_pages"`
	
	// HasNext есть ли следующая страница
	HasNext bool `json:"has_next"`
	
	// HasPrev есть ли предыдущая страница
	HasPrev bool `json:"has_prev"`
}

// NewPaginatedResult создает результат пагинации
func NewPaginatedResult(items []*User, total int64, pagination PaginationParams) *PaginatedResult {
	totalPages := int((total + int64(pagination.Limit) - 1) / int64(pagination.Limit))
	if totalPages == 0 {
		totalPages = 1
	}
	
	return &PaginatedResult{
		Items:      items,
		Total:      total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: totalPages,
		HasNext:    pagination.Page < totalPages,
		HasPrev:    pagination.Page > 1,
	}
}

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
	
	// FindUsers находит пользователей с пагинацией, сортировкой и фильтрацией
	FindUsers(ctx context.Context, query FindUsersQuery) (*PaginatedResult, error)
	
	// CountUsers возвращает общее количество пользователей с учетом фильтров
	CountUsers(ctx context.Context, filter FilterParams) (int64, error)

	// Delete удаляет пользователя по ID
	Delete(ctx context.Context, id uint) error
	
	// DeleteByIDs удаляет пользователей по списку ID (batch delete)
	DeleteByIDs(ctx context.Context, ids []uint) error

	// Exists проверяет существование пользователя по email
	// Используется для проверки уникальности email
	Exists(ctx context.Context, email string) (bool, error)
	
	// ExistsByIDs проверяет существование пользователей по списку ID
	ExistsByIDs(ctx context.Context, ids []uint) (map[uint]bool, error)
}

// ВАЖНО: Интерфейс Repository находится в domain слое, потому что:
// 1. Доменный слой определяет ПРАВИЛА работы с данными
// 2. Use Case слой будет зависеть от этого интерфейса
// 3. Infrastructure слой будет РЕАЛИЗОВЫВАТЬ этот интерфейс
// 4. Это обеспечивает ИНВЕРСИЮ ЗАВИСИМОСТЕЙ - внутренние слои не зависят от внешних