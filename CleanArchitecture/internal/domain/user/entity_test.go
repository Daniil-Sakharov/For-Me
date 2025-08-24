package user

import (
	"testing"
	"time"
)

// 📚 ИЗУЧИТЬ: Table-driven tests в Go
// https://github.com/golang/go/wiki/TableDrivenTests
// https://dave.cheney.net/2019/05/07/prefer-table-driven-tests

// TestNewUser_Success тестирует успешное создание пользователя
func TestNewUser_Success(t *testing.T) {
	// ARRANGE (подготовка данных)
	email := "test@example.com"
	name := "Test User"
	
	// ACT (выполнение действия)
	user, err := NewUser(email, name)
	
	// ASSERT (проверка результата)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if user == nil {
		t.Fatal("Expected user to be created, got nil")
	}
	
	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}
	
	if user.Name != name {
		t.Errorf("Expected name %s, got %s", name, user.Name)
	}
	
	// Проверяем, что время создания установлено
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	
	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
	
	// ID должен быть 0 для нового пользователя
	if user.ID != 0 {
		t.Errorf("Expected ID to be 0 for new user, got %d", user.ID)
	}
}

// TestNewUser_ValidationErrors тестирует валидацию при создании пользователя
// ПАТТЕРН: Table-driven tests для множественных сценариев
func TestNewUser_ValidationErrors(t *testing.T) {
	// Таблица тестовых случаев
	testCases := []struct {
		name          string  // Название теста
		email         string  // Входной email
		userName      string  // Входное имя
		expectedError error   // Ожидаемая ошибка
		description   string  // Описание теста
	}{
		{
			name:          "invalid_email_empty",
			email:         "",
			userName:      "Valid Name",
			expectedError: ErrInvalidEmail,
			description:   "Пустой email должен вызывать ошибку",
		},
		{
			name:          "invalid_email_format",
			email:         "not-an-email",
			userName:      "Valid Name",
			expectedError: ErrInvalidEmail,
			description:   "Неверный формат email должен вызывать ошибку",
		},
		{
			name:          "invalid_email_no_domain",
			email:         "test@",
			userName:      "Valid Name",
			expectedError: ErrInvalidEmail,
			description:   "Email без домена должен вызывать ошибку",
		},
		{
			name:          "empty_name",
			email:         "test@example.com",
			userName:      "",
			expectedError: ErrEmptyName,
			description:   "Пустое имя должно вызывать ошибку",
		},
		{
			name:          "whitespace_only_name",
			email:         "test@example.com",
			userName:      "   ",
			expectedError: ErrEmptyName,
			description:   "Имя из пробелов должно вызывать ошибку",
		},
	}
	
	// Выполняем каждый тестовый случай
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			user, err := NewUser(tc.email, tc.userName)
			
			// ASSERT
			if err != tc.expectedError {
				t.Errorf("Expected error %v, got %v. %s", tc.expectedError, err, tc.description)
			}
			
			if user != nil {
				t.Errorf("Expected user to be nil when error occurs, got %+v", user)
			}
		})
	}
}

// TestUser_SetPassword тестирует установку пароля
func TestUser_SetPassword(t *testing.T) {
	// ARRANGE
	user, err := NewUser("test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	password := "validpassword123"
	
	// ACT
	err = user.SetPassword(password)
	
	// ASSERT
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if user.HashedPassword == "" {
		t.Error("Expected hashed password to be set")
	}
	
	if user.HashedPassword == password {
		t.Error("Password should be hashed, not stored in plain text")
	}
}

// TestUser_SetPassword_ValidationErrors тестирует валидацию пароля
func TestUser_SetPassword_ValidationErrors(t *testing.T) {
	// ARRANGE
	user, err := NewUser("test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	testCases := []struct {
		name          string
		password      string
		expectedError error
		description   string
	}{
		{
			name:          "empty_password",
			password:      "",
			expectedError: ErrEmptyPassword,
			description:   "Пустой пароль должен вызывать ошибку",
		},
		{
			name:          "short_password",
			password:      "123",
			expectedError: ErrPasswordTooWeak,
			description:   "Короткий пароль должен вызывать ошибку",
		},
		{
			name:          "seven_char_password",
			password:      "1234567",
			expectedError: ErrPasswordTooWeak,
			description:   "Пароль из 7 символов должен вызывать ошибку",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			err := user.SetPassword(tc.password)
			
			// ASSERT
			if err != tc.expectedError {
				t.Errorf("Expected error %v, got %v. %s", tc.expectedError, err, tc.description)
			}
		})
	}
}

// TestUser_CheckPassword тестирует проверку пароля
func TestUser_CheckPassword(t *testing.T) {
	// ARRANGE
	user, err := NewUser("test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	password := "validpassword123"
	err = user.SetPassword(password)
	if err != nil {
		t.Fatalf("Failed to set password: %v", err)
	}
	
	// ACT & ASSERT - правильный пароль
	err = user.CheckPassword(password)
	if err != nil {
		t.Errorf("Expected no error for correct password, got %v", err)
	}
	
	// ACT & ASSERT - неправильный пароль
	err = user.CheckPassword("wrongpassword")
	if err == nil {
		t.Error("Expected error for incorrect password")
	}
}

// TestUser_UpdatedAt тестирует обновление времени изменения
func TestUser_UpdatedAt(t *testing.T) {
	// ARRANGE
	user, err := NewUser("test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	originalUpdatedAt := user.UpdatedAt
	
	// Небольшая задержка, чтобы время изменилось
	time.Sleep(time.Millisecond)
	
	// ACT
	user.Touch()
	
	// ASSERT
	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated after Touch()")
	}
}

// БЕНЧМАРКИ для измерения производительности
// 📚 ИЗУЧИТЬ: Benchmarking в Go
// https://golang.org/pkg/testing/#hdr-Benchmarks
// https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go

// BenchmarkNewUser измеряет производительность создания пользователя
func BenchmarkNewUser(b *testing.B) {
	email := "test@example.com"
	name := "Test User"
	
	// b.N автоматически подбирается для получения точных измерений
	for i := 0; i < b.N; i++ {
		_, err := NewUser(email, name)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// BenchmarkSetPassword измеряет производительность хеширования пароля
func BenchmarkSetPassword(b *testing.B) {
	user, err := NewUser("test@example.com", "Test User")
	if err != nil {
		b.Fatalf("Failed to create user: %v", err)
	}
	
	password := "validpassword123"
	
	// Сбрасываем таймер, чтобы не учитывать время подготовки
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := user.SetPassword(password)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// ПРИМЕРЫ для документации
// 📚 ИЗУЧИТЬ: Example tests в Go
// https://golang.org/pkg/testing/#hdr-Examples

// ExampleNewUser демонстрирует создание нового пользователя
func ExampleNewUser() {
	user, err := NewUser("john@example.com", "John Doe")
	if err != nil {
		panic(err)
	}
	
	println("User created:", user.Email, user.Name)
	// Output: User created: john@example.com John Doe
}

// ExampleUser_SetPassword демонстрирует установку пароля
func ExampleUser_SetPassword() {
	user, _ := NewUser("john@example.com", "John Doe")
	
	err := user.SetPassword("securepassword123")
	if err != nil {
		panic(err)
	}
	
	println("Password set successfully")
	// Output: Password set successfully
}