package mocks

import (
	"clean-url-shortener/internal/usecase/auth"
)

// MockPasswordHasher - мок для хеширования паролей
// 📚 ИЗУЧИТЬ: Password hashing best practices
// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
// https://pkg.go.dev/golang.org/x/crypto/bcrypt
type MockPasswordHasher struct {
	HashFunc    func(password string) (string, error)
	CompareFunc func(hashedPassword, password string) error
	
	// Счетчики и последние параметры для проверки в тестах
	HashCallCount    int
	CompareCallCount int
	LastHashedPassword string
	LastPlainPassword  string
}

// Hash реализует интерфейс auth.PasswordHasher
func (m *MockPasswordHasher) Hash(password string) (string, error) {
	m.HashCallCount++
	m.LastPlainPassword = password
	
	if m.HashFunc != nil {
		return m.HashFunc(password)
	}
	
	// Поведение по умолчанию: простое "хеширование" для тестов
	return "hashed_" + password, nil
}

// Compare реализует интерфейс auth.PasswordHasher
func (m *MockPasswordHasher) Compare(hashedPassword, password string) error {
	m.CompareCallCount++
	m.LastHashedPassword = hashedPassword
	m.LastPlainPassword = password
	
	if m.CompareFunc != nil {
		return m.CompareFunc(hashedPassword, password)
	}
	
	// Поведение по умолчанию: пароли совпадают, если хеш начинается с "hashed_"
	expectedHash := "hashed_" + password
	if hashedPassword != expectedHash {
		return auth.ErrInvalidCredentials
	}
	
	return nil
}

// Reset сбрасывает состояние мока
func (m *MockPasswordHasher) Reset() {
	m.HashCallCount = 0
	m.CompareCallCount = 0
	m.LastHashedPassword = ""
	m.LastPlainPassword = ""
	m.HashFunc = nil
	m.CompareFunc = nil
}

// MockTokenGenerator - мок для генерации JWT токенов
// 📚 ИЗУЧИТЬ: JWT токены и их безопасность
// https://jwt.io/introduction/
// https://auth0.com/blog/a-look-at-the-latest-draft-for-jwt-bcp/
// https://tools.ietf.org/html/rfc7519
type MockTokenGenerator struct {
	GenerateFunc func(userID uint, email string) (string, error)
	ValidateFunc func(token string) (*auth.TokenData, error)
	
	// Счетчики и последние параметры для проверки в тестах
	GenerateCallCount int
	ValidateCallCount int
	LastGeneratedUserID uint
	LastGeneratedEmail  string
	LastValidatedToken  string
}

// Generate реализует интерфейс auth.TokenGenerator
func (m *MockTokenGenerator) Generate(userID uint, email string) (string, error) {
	m.GenerateCallCount++
	m.LastGeneratedUserID = userID
	m.LastGeneratedEmail = email
	
	if m.GenerateFunc != nil {
		return m.GenerateFunc(userID, email)
	}
	
	// Поведение по умолчанию: простой токен для тестов
	return "mock_token_" + email, nil
}

// Validate реализует интерфейс auth.TokenGenerator
func (m *MockTokenGenerator) Validate(token string) (*auth.TokenData, error) {
	m.ValidateCallCount++
	m.LastValidatedToken = token
	
	if m.ValidateFunc != nil {
		return m.ValidateFunc(token)
	}
	
	// Поведение по умолчанию: валидный токен для тестов
	if token == "invalid_token" {
		return nil, auth.ErrInvalidCredentials
	}
	
	return &auth.TokenData{
		UserID: 1,
		Email:  "test@example.com",
	}, nil
}

// Reset сбрасывает состояние мока
func (m *MockTokenGenerator) Reset() {
	m.GenerateCallCount = 0
	m.ValidateCallCount = 0
	m.LastGeneratedUserID = 0
	m.LastGeneratedEmail = ""
	m.LastValidatedToken = ""
	m.GenerateFunc = nil
	m.ValidateFunc = nil
}

// HELPER FUNCTIONS для тестов

// NewMockPasswordHasherWithError создает мок, который всегда возвращает ошибку
func NewMockPasswordHasherWithError(err error) *MockPasswordHasher {
	return &MockPasswordHasher{
		HashFunc: func(password string) (string, error) {
			return "", err
		},
		CompareFunc: func(hashedPassword, password string) error {
			return err
		},
	}
}

// NewMockTokenGeneratorWithError создает мок, который всегда возвращает ошибку
func NewMockTokenGeneratorWithError(err error) *MockTokenGenerator {
	return &MockTokenGenerator{
		GenerateFunc: func(userID uint, email string) (string, error) {
			return "", err
		},
		ValidateFunc: func(token string) (*auth.TokenData, error) {
			return nil, err
		},
	}
}