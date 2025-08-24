package mocks

import (
	"context"
	"clean-url-shortener/internal/domain/user"
)

// MockUserRepository - мок репозитория для тестирования
// 📚 ИЗУЧИТЬ: Mocking patterns in Go
// https://blog.codecentric.de/en/2017/08/gomock-tutorial/
// https://github.com/golang/mock
//
// ЗАЧЕМ НУЖНЫ МОКИ:
// 1. Тестируем бизнес-логику БЕЗ реальной БД
// 2. Контролируем поведение зависимостей
// 3. Тесты работают быстро и независимо
// 4. Можем симулировать ошибки и edge cases
type MockUserRepository struct {
	// Поля для контроля поведения мока
	SaveFunc       func(ctx context.Context, user *user.User) error
	FindByIDFunc   func(ctx context.Context, id uint) (*user.User, error)
	FindByEmailFunc func(ctx context.Context, email string) (*user.User, error)
	DeleteFunc     func(ctx context.Context, id uint) error
	ExistsFunc     func(ctx context.Context, email string) (bool, error)
	
	// Счетчики вызовов для проверки в тестах
	SaveCallCount       int
	FindByIDCallCount   int
	FindByEmailCallCount int
	DeleteCallCount     int
	ExistsCallCount     int
	
	// Последние переданные параметры для проверки
	LastSavedUser    *user.User
	LastFindByIDParam uint
	LastFindByEmailParam string
	LastDeleteParam  uint
	LastExistsParam  string
}

// Save реализует интерфейс user.Repository
func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
	m.SaveCallCount++
	m.LastSavedUser = u
	
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, u)
	}
	
	// Поведение по умолчанию: успешное сохранение
	if u.ID == 0 {
		u.ID = 1 // Симулируем автоинкремент
	}
	return nil
}

// FindByID реализует интерфейс user.Repository
func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	m.FindByIDCallCount++
	m.LastFindByIDParam = id
	
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	
	// Поведение по умолчанию: пользователь не найден
	return nil, nil
}

// FindByEmail реализует интерфейс user.Repository
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	m.FindByEmailCallCount++
	m.LastFindByEmailParam = email
	
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	
	// Поведение по умолчанию: пользователь не найден
	return nil, nil
}

// Delete реализует интерфейс user.Repository
func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	m.DeleteCallCount++
	m.LastDeleteParam = id
	
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	
	// Поведение по умолчанию: успешное удаление
	return nil
}

// Exists реализует интерфейс user.Repository
func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	m.ExistsCallCount++
	m.LastExistsParam = email
	
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, email)
	}
	
	// Поведение по умолчанию: пользователь не существует
	return false, nil
}

// Reset сбрасывает состояние мока (полезно в тестах)
func (m *MockUserRepository) Reset() {
	m.SaveCallCount = 0
	m.FindByIDCallCount = 0
	m.FindByEmailCallCount = 0
	m.DeleteCallCount = 0
	m.ExistsCallCount = 0
	
	m.LastSavedUser = nil
	m.LastFindByIDParam = 0
	m.LastFindByEmailParam = ""
	m.LastDeleteParam = 0
	m.LastExistsParam = ""
	
	m.SaveFunc = nil
	m.FindByIDFunc = nil
	m.FindByEmailFunc = nil
	m.DeleteFunc = nil
	m.ExistsFunc = nil
}

// ПАТТЕРН: Test Builder для создания тестовых пользователей
// 📚 ИЗУЧИТЬ: Test Data Builders
// https://martinfowler.com/bliki/ObjectMother.html

// UserTestBuilder помогает создавать пользователей для тестов
type UserTestBuilder struct {
	user *user.User
}

// NewUserTestBuilder создает билдер с базовыми значениями
func NewUserTestBuilder() *UserTestBuilder {
	u, _ := user.NewUser("test@example.com", "Test User")
	u.ID = 1
	
	return &UserTestBuilder{user: u}
}

// WithID устанавливает ID пользователя
func (b *UserTestBuilder) WithID(id uint) *UserTestBuilder {
	b.user.ID = id
	return b
}

// WithEmail устанавливает email пользователя
func (b *UserTestBuilder) WithEmail(email string) *UserTestBuilder {
	b.user.Email = email
	return b
}

// WithName устанавливает имя пользователя
func (b *UserTestBuilder) WithName(name string) *UserTestBuilder {
	b.user.Name = name
	return b
}

// WithHashedPassword устанавливает хешированный пароль
func (b *UserTestBuilder) WithHashedPassword(hash string) *UserTestBuilder {
	b.user.HashedPassword = hash
	return b
}

// Build возвращает готового пользователя
func (b *UserTestBuilder) Build() *user.User {
	return b.user
}