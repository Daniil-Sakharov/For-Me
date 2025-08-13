// 🏛️ УНИВЕРСАЛЬНЫЙ ШАБЛОН DOMAIN ENTITY для Clean Architecture
//
// ✅ ВСЕГДА одинаково:
// - Структура файла с доменной моделью
// - Конструктор с валидацией (New{Entity})
// - Доменные ошибки (Err{Something})
// - Методы для изменения состояния
// - НЕТ внешних зависимостей!
//
// 🔀 МЕНЯЕТСЯ в зависимости от проекта:
// - Название сущности (User, Product, Order, etc.)
// - Поля структуры
// - Доменные правила валидации
// - Конкретные методы

package example_entity // 🔀 ИЗМЕНИТЕ: название сущности (user, product, order, etc.)

import (
	"errors"
	"regexp"
	"time"
	
	// ✅ ВСЕГДА: только стандартная библиотека Go!
	// ❌ НЕ ИМПОРТИРУЙТЕ: gorm, database/sql, http, json, и другие внешние библиотеки
)

// 🏗️ ДОМЕННАЯ МОДЕЛЬ
// ✅ ВСЕГДА: основная структура сущности
type ExampleEntity struct { // 🔀 ИЗМЕНИТЕ: название сущности (User, Product, Order, etc.)
	
	// ✅ ВСЕГДА: базовые поля (почти для всех сущностей)
	ID        uint      `json:"id"`         // Уникальный идентификатор
	CreatedAt time.Time `json:"created_at"` // Время создания
	UpdatedAt time.Time `json:"updated_at"` // Время обновления
	
	// 🔀 ДОБАВЬТЕ: специфичные поля для вашей сущности
	// Примеры для разных сущностей:
	
	// ДЛЯ USER:
	// Email          string `json:"email"`
	// HashedPassword string `json:"-"`        // НЕ сериализуется в JSON
	// Name           string `json:"name"`
	// IsActive       bool   `json:"is_active"`
	// Role           string `json:"role"`     // admin, user, moderator
	
	// ДЛЯ PRODUCT:
	// Name        string  `json:"name"`
	// Description string  `json:"description"`
	// Price       float64 `json:"price"`
	// SKU         string  `json:"sku"`        // Артикул
	// CategoryID  uint    `json:"category_id"`
	// IsAvailable bool    `json:"is_available"`
	// Stock       int     `json:"stock"`      // Количество на складе
	
	// ДЛЯ ORDER:
	// UserID      uint    `json:"user_id"`
	// TotalAmount float64 `json:"total_amount"`
	// Status      string  `json:"status"`     // pending, confirmed, shipped, delivered
	// ShippingAddress string `json:"shipping_address"`
	
	// 🔀 ВАШИ ПОЛЯ:
	SampleField string `json:"sample_field"` // 🔀 ЗАМЕНИТЕ на реальные поля
}

// 🚨 ДОМЕННЫЕ ОШИБКИ
// ✅ ВСЕГДА: определите ошибки для вашей сущности
var (
	// 🔀 ИЗМЕНИТЕ: ошибки под вашу сущность
	ErrInvalidExampleEntity     = errors.New("invalid example entity")          // 🔀 ИЗМЕНИТЕ
	ErrEmptySampleField        = errors.New("sample field cannot be empty")     // 🔀 ИЗМЕНИТЕ
	ErrExampleEntityNotFound   = errors.New("example entity not found")         // 🔀 ИЗМЕНИТЕ
	ErrExampleEntityExists     = errors.New("example entity already exists")    // 🔀 ИЗМЕНИТЕ
	
	// ПРИМЕРЫ ошибок для разных сущностей:
	
	// ДЛЯ USER:
	// ErrInvalidEmail     = errors.New("invalid email format")
	// ErrWeakPassword     = errors.New("password too weak")
	// ErrUserNotFound     = errors.New("user not found")
	// ErrUserExists       = errors.New("user already exists")
	// ErrEmptyName        = errors.New("name cannot be empty")
	// ErrInvalidRole      = errors.New("invalid user role")
	
	// ДЛЯ PRODUCT:
	// ErrInvalidPrice     = errors.New("price must be positive")
	// ErrEmptyName        = errors.New("product name cannot be empty")
	// ErrInvalidSKU       = errors.New("invalid SKU format")
	// ErrOutOfStock       = errors.New("product out of stock")
	// ErrProductNotFound  = errors.New("product not found")
	
	// ДЛЯ ORDER:
	// ErrInvalidAmount    = errors.New("invalid order amount")
	// ErrInvalidStatus    = errors.New("invalid order status")
	// ErrOrderNotFound    = errors.New("order not found")
	// ErrEmptyAddress     = errors.New("shipping address cannot be empty")
)

// 🏭 КОНСТРУКТОР
// ✅ ВСЕГДА: конструктор с валидацией доменных правил
func NewExampleEntity( /* 🔀 ПАРАМЕТРЫ под вашу сущность */ sampleField string) (*ExampleEntity, error) {
	
	// ✅ ВСЕГДА: валидация входных параметров
	if err := validateSampleField(sampleField); err != nil { // 🔀 ИЗМЕНИТЕ: валидация
		return nil, err
	}
	
	// 🔀 ДОБАВЬТЕ: дополнительную валидацию под вашу сущность
	// Примеры:
	// if err := validateEmail(email); err != nil {
	//     return nil, err
	// }
	// if err := validatePrice(price); err != nil {
	//     return nil, err
	// }
	
	// ✅ ВСЕГДА: установка временных меток
	now := time.Now()
	
	return &ExampleEntity{
		// ✅ ВСЕГДА: базовые поля
		CreatedAt: now,
		UpdatedAt: now,
		
		// 🔀 ВАШИ ПОЛЯ:
		SampleField: sampleField, // 🔀 ЗАМЕНИТЕ на реальные поля
		
		// Примеры для разных сущностей:
		// Email:          email,
		// HashedPassword: hashedPassword,
		// Name:           name,
		// IsActive:       true,
		// Role:           "user",
		
		// Name:        name,
		// Description: description,
		// Price:       price,
		// SKU:         sku,
		// IsAvailable: true,
		// Stock:       stock,
	}, nil
}

// 🔧 МЕТОДЫ ДЛЯ ИЗМЕНЕНИЯ СОСТОЯНИЯ
// ✅ ВСЕГДА: методы для безопасного изменения полей

// 🔀 ИЗМЕНИТЕ: методы под вашу сущность
func (e *ExampleEntity) UpdateSampleField(newValue string) error {
	if err := validateSampleField(newValue); err != nil {
		return err
	}
	
	e.SampleField = newValue
	e.UpdatedAt = time.Now() // ✅ ВСЕГДА: обновляем время
	return nil
}

// 📝 ПРИМЕРЫ методов для разных сущностей:

// ДЛЯ USER:
// func (u *User) UpdateEmail(newEmail string) error {
//     if err := validateEmail(newEmail); err != nil {
//         return err
//     }
//     u.Email = newEmail
//     u.UpdatedAt = time.Now()
//     return nil
// }
//
// func (u *User) UpdateName(newName string) error {
//     if newName == "" {
//         return ErrEmptyName
//     }
//     u.Name = newName
//     u.UpdatedAt = time.Now()
//     return nil
// }
//
// func (u *User) SetPassword(hashedPassword string) error {
//     if hashedPassword == "" {
//         return ErrWeakPassword
//     }
//     u.HashedPassword = hashedPassword
//     u.UpdatedAt = time.Now()
//     return nil
// }
//
// func (u *User) Activate() {
//     u.IsActive = true
//     u.UpdatedAt = time.Now()
// }
//
// func (u *User) Deactivate() {
//     u.IsActive = false
//     u.UpdatedAt = time.Now()
// }

// ДЛЯ PRODUCT:
// func (p *Product) UpdatePrice(newPrice float64) error {
//     if newPrice <= 0 {
//         return ErrInvalidPrice
//     }
//     p.Price = newPrice
//     p.UpdatedAt = time.Now()
//     return nil
// }
//
// func (p *Product) UpdateStock(newStock int) error {
//     if newStock < 0 {
//         return ErrInvalidStock
//     }
//     p.Stock = newStock
//     p.UpdatedAt = time.Now()
//     return nil
// }
//
// func (p *Product) MakeAvailable() {
//     p.IsAvailable = true
//     p.UpdatedAt = time.Now()
// }
//
// func (p *Product) MakeUnavailable() {
//     p.IsAvailable = false
//     p.UpdatedAt = time.Now()
// }

// ДЛЯ ORDER:
// func (o *Order) ConfirmOrder() error {
//     if o.Status != "pending" {
//         return ErrInvalidStatus
//     }
//     o.Status = "confirmed"
//     o.UpdatedAt = time.Now()
//     return nil
// }
//
// func (o *Order) ShipOrder() error {
//     if o.Status != "confirmed" {
//         return ErrInvalidStatus
//     }
//     o.Status = "shipped"
//     o.UpdatedAt = time.Now()
//     return nil
// }

// 🔍 МЕТОДЫ ЗАПРОСОВ (БЕЗ ИЗМЕНЕНИЯ СОСТОЯНИЯ)
// ✅ ВСЕГДА: методы для проверки состояния

// 🔀 ИЗМЕНИТЕ: методы под вашу сущность
func (e *ExampleEntity) IsValid() bool {
	return e.SampleField != "" // 🔀 ИЗМЕНИТЕ: логика валидации
}

// 📝 ПРИМЕРЫ методов для разных сущностей:

// ДЛЯ USER:
// func (u *User) IsActive() bool {
//     return u.IsActive
// }
//
// func (u *User) IsAdmin() bool {
//     return u.Role == "admin"
// }
//
// func (u *User) HasValidEmail() bool {
//     return validateEmail(u.Email) == nil
// }

// ДЛЯ PRODUCT:
// func (p *Product) IsAvailable() bool {
//     return p.IsAvailable && p.Stock > 0
// }
//
// func (p *Product) IsInStock() bool {
//     return p.Stock > 0
// }
//
// func (p *Product) CanOrder(quantity int) bool {
//     return p.IsAvailable && p.Stock >= quantity
// }

// ДЛЯ ORDER:
// func (o *Order) IsPending() bool {
//     return o.Status == "pending"
// }
//
// func (o *Order) IsCompleted() bool {
//     return o.Status == "delivered"
// }
//
// func (o *Order) CanBeCancelled() bool {
//     return o.Status == "pending" || o.Status == "confirmed"
// }

// 🔧 ПРИВАТНЫЕ ФУНКЦИИ ВАЛИДАЦИИ
// ✅ ВСЕГДА: вспомогательные функции для валидации

// 🔀 ИЗМЕНИТЕ: валидацию под вашу сущность
func validateSampleField(field string) error {
	if field == "" {
		return ErrEmptySampleField
	}
	// 🔀 ДОБАВЬТЕ: дополнительную валидацию
	return nil
}

// 📝 ПРИМЕРЫ валидации для разных сущностей:

// ДЛЯ USER:
// var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
//
// func validateEmail(email string) error {
//     if email == "" {
//         return ErrEmptyEmail
//     }
//     if !emailRegex.MatchString(email) {
//         return ErrInvalidEmail
//     }
//     return nil
// }
//
// func validatePassword(password string) error {
//     if len(password) < 8 {
//         return ErrWeakPassword
//     }
//     // Дополнительные проверки пароля
//     return nil
// }

// ДЛЯ PRODUCT:
// var skuRegex = regexp.MustCompile(`^[A-Z0-9\-]+$`)
//
// func validateSKU(sku string) error {
//     if sku == "" {
//         return ErrEmptySKU
//     }
//     if !skuRegex.MatchString(sku) {
//         return ErrInvalidSKU
//     }
//     return nil
// }
//
// func validatePrice(price float64) error {
//     if price <= 0 {
//         return ErrInvalidPrice
//     }
//     return nil
// }

// 📚 ИНСТРУКЦИЯ ПО АДАПТАЦИИ:
//
// 1. 🔀 ПЕРЕИМЕНУЙТЕ файл и пакет:
//    example_entity → user (или product, order, etc.)
//
// 2. 🔀 ЗАМЕНИТЕ ExampleEntity:
//    На User, Product, Order или вашу сущность
//
// 3. 🔀 ОПРЕДЕЛИТЕ поля структуры:
//    Добавьте специфичные поля вместо SampleField
//
// 4. 🔀 СОЗДАЙТЕ доменные ошибки:
//    Определите ошибки специфичные для вашей сущности
//
// 5. 🔀 АДАПТИРУЙТЕ конструктор:
//    Измените параметры и валидацию под вашу логику
//
// 6. 🔀 ДОБАВЬТЕ методы изменения:
//    Создайте методы для изменения полей вашей сущности
//
// 7. 🔀 ДОБАВЬТЕ методы запросов:
//    Создайте методы для проверки состояния
//
// 8. 🔀 РЕАЛИЗУЙТЕ валидацию:
//    Создайте функции валидации для ваших полей
//
// 9. ✅ ВСЕГДА сохраняйте:
//    - НЕТ внешних зависимостей
//    - Валидация в конструкторе и методах
//    - Обновление UpdatedAt при изменениях
//    - Доменные ошибки для бизнес-правил