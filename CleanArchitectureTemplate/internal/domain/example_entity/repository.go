// 🗄️ УНИВЕРСАЛЬНЫЙ ШАБЛОН REPOSITORY INTERFACE для Clean Architecture
//
// ✅ ВСЕГДА одинаково:
// - Интерфейс Repository в domain слое
// - Базовые CRUD методы (Save, FindByID, Delete)
// - Контекст во всех методах
// - НЕТ импортов конкретных технологий БД!
//
// 🔀 МЕНЯЕТСЯ в зависимости от проекта:
// - Название сущности (User, Product, Order, etc.)
// - Специфичные методы поиска (FindByEmail, FindByCategory, etc.)
// - Методы для бизнес-логики (ExistsByEmail, CountByStatus, etc.)

package example_entity // 🔀 ИЗМЕНИТЕ: название сущности (user, product, order, etc.)

import (
	"context"
	
	// ✅ ВСЕГДА: только стандартная библиотека!
	// ❌ НЕ ИМПОРТИРУЙТЕ: database/sql, gorm, mongo-driver и другие БД библиотеки
)

// 🏛️ ИНТЕРФЕЙС РЕПОЗИТОРИЯ
// ✅ ВСЕГДА: определен в domain слое, реализуется в infrastructure
type Repository interface {
	
	// ✅ БАЗОВЫЕ CRUD ОПЕРАЦИИ (всегда нужны):
	
	// Сохранить сущность (создать новую или обновить существующую)
	Save(ctx context.Context, entity *ExampleEntity) error // 🔀 ИЗМЕНИТЕ: тип сущности
	
	// Найти сущность по ID
	FindByID(ctx context.Context, id uint) (*ExampleEntity, error) // 🔀 ИЗМЕНИТЕ: тип сущности
	
	// Удалить сущность по ID
	Delete(ctx context.Context, id uint) error
	
	// 🔀 СПЕЦИФИЧНЫЕ МЕТОДЫ ПОИСКА (добавьте под вашу сущность):
	
	// Примеры методов для ExampleEntity (замените на реальные):
	FindBySampleField(ctx context.Context, field string) (*ExampleEntity, error) // 🔀 ИЗМЕНИТЕ
	ExistsBySampleField(ctx context.Context, field string) (bool, error)         // 🔀 ИЗМЕНИТЕ
	
	// 📝 ПРИМЕРЫ методов для разных сущностей:
	
	// ДЛЯ USER:
	// FindByEmail(ctx context.Context, email string) (*User, error)
	// FindByEmailAndPassword(ctx context.Context, email, hashedPassword string) (*User, error)
	// ExistsByEmail(ctx context.Context, email string) (bool, error)
	// FindActiveUsers(ctx context.Context) ([]*User, error)
	// FindUsersByRole(ctx context.Context, role string) ([]*User, error)
	// CountUsers(ctx context.Context) (int64, error)
	// UpdateLastLogin(ctx context.Context, userID uint) error
	
	// ДЛЯ PRODUCT:
	// FindBySKU(ctx context.Context, sku string) (*Product, error)
	// FindByCategory(ctx context.Context, categoryID uint) ([]*Product, error)
	// FindAvailableProducts(ctx context.Context) ([]*Product, error)
	// FindProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*Product, error)
	// ExistsBySKU(ctx context.Context, sku string) (bool, error)
	// UpdateStock(ctx context.Context, productID uint, newStock int) error
	// DecrementStock(ctx context.Context, productID uint, quantity int) error
	// CountProductsByCategory(ctx context.Context, categoryID uint) (int64, error)
	
	// ДЛЯ ORDER:
	// FindByUserID(ctx context.Context, userID uint) ([]*Order, error)
	// FindByStatus(ctx context.Context, status string) ([]*Order, error)
	// FindOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*Order, error)
	// UpdateStatus(ctx context.Context, orderID uint, status string) error
	// CountOrdersByStatus(ctx context.Context, status string) (int64, error)
	// CalculateTotalRevenue(ctx context.Context) (float64, error)
	// FindPendingOrders(ctx context.Context) ([]*Order, error)
	
	// 🔀 ДОПОЛНИТЕЛЬНЫЕ МЕТОДЫ (добавьте под ваши потребности):
	
	// Пагинация:
	// FindAll(ctx context.Context, limit, offset int) ([]*ExampleEntity, error)
	// FindWithPagination(ctx context.Context, page, size int) ([]*ExampleEntity, int64, error)
	
	// Поиск и фильтрация:
	// Search(ctx context.Context, query string) ([]*ExampleEntity, error)
	// FindByFilters(ctx context.Context, filters map[string]interface{}) ([]*ExampleEntity, error)
	
	// Счетчики и агрегация:
	// Count(ctx context.Context) (int64, error)
	// CountByCondition(ctx context.Context, condition string) (int64, error)
	
	// Batch операции:
	// SaveBatch(ctx context.Context, entities []*ExampleEntity) error
	// DeleteBatch(ctx context.Context, ids []uint) error
	
	// Soft delete (если нужен):
	// SoftDelete(ctx context.Context, id uint) error
	// Restore(ctx context.Context, id uint) error
	// FindDeleted(ctx context.Context) ([]*ExampleEntity, error)
}

// 📚 ИНСТРУКЦИЯ ПО АДАПТАЦИИ:
//
// 1. 🔀 ПЕРЕИМЕНУЙТЕ файл и пакет:
//    example_entity → user (или product, order, etc.)
//
// 2. 🔀 ЗАМЕНИТЕ ExampleEntity:
//    На User, Product, Order или вашу сущность
//
// 3. ✅ СОХРАНИТЕ базовые CRUD методы:
//    Save, FindByID, Delete - нужны почти всегда
//
// 4. 🔀 ДОБАВЬТЕ специфичные методы:
//    Выберите из примеров те, что подходят вашей сущности
//
// 5. 🔀 СОЗДАЙТЕ свои методы:
//    Добавьте методы специфичные для вашей бизнес-логики
//
// 6. ✅ ВСЕГДА используйте context.Context:
//    Первый параметр в каждом методе
//
// 7. ✅ НЕ импортируйте БД библиотеки:
//    Интерфейс должен быть чистым от технических деталей
//
// 8. 📝 ДОКУМЕНТИРУЙТЕ методы:
//    Добавьте комментарии для сложных методов
//
// 📋 ПРИНЦИПЫ REPOSITORY PATTERN:
//
// ✅ ДЕЛАЙТЕ:
// - Определяйте ЧТО нужно делать, не КАК
// - Группируйте методы логически
// - Используйте понятные названия методов
// - Возвращайте доменные типы (не DTO)
// - Обрабатывайте ошибки на уровне интерфейса
//
// ❌ НЕ ДЕЛАЙТЕ:
// - Не привязывайтесь к конкретной БД
// - Не добавляйте SQL в интерфейс
// - Не смешивайте разные сущности в одном репозитории
// - Не делайте методы слишком специфичными
// - Не забывайте про контекст
//
// 🔄 СВЯЗЬ С ДРУГИМИ СЛОЯМИ:
//
// Domain → Repository Interface (этот файл)
//    ↑
// Infrastructure → Repository Implementation
//    ↑
// Use Cases → Используют Repository Interface
//    ↑
// Controllers → Вызывают Use Cases
//
// 📝 ПРИМЕРЫ РЕАЛИЗАЦИИ (в infrastructure/database):
//
// type UserRepository struct {
//     db *sql.DB  // или любая другая БД
// }
//
// func NewUserRepository(db *sql.DB) user.Repository {
//     return &UserRepository{db: db}
// }
//
// func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
//     // SQL или MongoDB запрос
// }
//
// func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
//     // SQL или MongoDB запрос
// }