package auth

import (
	"context"
	"errors"
	"testing"
	"time"
	
	"clean-url-shortener/internal/domain/user"
	"clean-url-shortener/internal/domain/user/mocks"
	authMocks "clean-url-shortener/internal/usecase/auth/mocks"
)

// 📚 ИЗУЧИТЬ: Testing Use Cases в Clean Architecture
// https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
// https://martinfowler.com/articles/practical-test-pyramid.html
//
// ПРИНЦИПЫ ТЕСТИРОВАНИЯ USE CASES:
// 1. Тестируем БИЗНЕС-ЛОГИКУ, а не инфраструктуру
// 2. Используем МОКИ для всех внешних зависимостей
// 3. Тестируем ВСЕ сценарии: успешные и ошибочные
// 4. Проверяем ВЗАИМОДЕЙСТВИЕ между компонентами

// TestRegisterUseCase_Success тестирует успешную регистрацию пользователя
func TestRegisterUseCase_Success(t *testing.T) {
	// ARRANGE - подготовка моков и данных
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	// Настраиваем поведение моков
	mockRepo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
		return false, nil // Пользователь не существует
	}
	
	mockHasher.HashFunc = func(password string) (string, error) {
		return "hashed_" + password, nil
	}
	
	mockTokenGen.GenerateFunc = func(userID uint, email string) (string, error) {
		return "generated_token", nil
	}
	
	// Создаем Use Case с моками
	uc := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
	
	// Подготавливаем запрос
	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "validpassword123",
		Name:     "Test User",
	}
	
	// ACT - выполняем действие
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	
	// ASSERT - проверяем результат
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	// Проверяем данные ответа
	if response.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, response.Email)
	}
	
	if response.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, response.Name)
	}
	
	if response.Token != "generated_token" {
		t.Errorf("Expected token 'generated_token', got %s", response.Token)
	}
	
	if response.UserID == 0 {
		t.Error("Expected UserID to be set")
	}
	
	// Проверяем, что моки были вызваны правильно
	if mockRepo.ExistsCallCount != 1 {
		t.Errorf("Expected Exists to be called once, called %d times", mockRepo.ExistsCallCount)
	}
	
	if mockHasher.HashCallCount != 1 {
		t.Errorf("Expected Hash to be called once, called %d times", mockHasher.HashCallCount)
	}
	
	if mockRepo.SaveCallCount != 1 {
		t.Errorf("Expected Save to be called once, called %d times", mockRepo.SaveCallCount)
	}
	
	if mockTokenGen.GenerateCallCount != 1 {
		t.Errorf("Expected Generate to be called once, called %d times", mockTokenGen.GenerateCallCount)
	}
	
	// Проверяем параметры, переданные в моки
	if mockRepo.LastExistsParam != req.Email {
		t.Errorf("Expected Exists to be called with %s, called with %s", req.Email, mockRepo.LastExistsParam)
	}
	
	if mockHasher.LastPlainPassword != req.Password {
		t.Errorf("Expected Hash to be called with %s, called with %s", req.Password, mockHasher.LastPlainPassword)
	}
}

// TestRegisterUseCase_UserAlreadyExists тестирует попытку регистрации существующего пользователя
func TestRegisterUseCase_UserAlreadyExists(t *testing.T) {
	// ARRANGE
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	// Настраиваем мок: пользователь уже существует
	mockRepo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
		return true, nil // Пользователь уже существует!
	}
	
	uc := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := RegisterRequest{
		Email:    "existing@example.com",
		Password: "validpassword123",
		Name:     "Test User",
	}
	
	// ACT
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	
	// ASSERT
	if err != ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
	
	if response != nil {
		t.Errorf("Expected nil response when user exists, got %+v", response)
	}
	
	// Проверяем, что другие методы НЕ были вызваны
	if mockHasher.HashCallCount > 0 {
		t.Error("Hash should not be called when user already exists")
	}
	
	if mockRepo.SaveCallCount > 0 {
		t.Error("Save should not be called when user already exists")
	}
	
	if mockTokenGen.GenerateCallCount > 0 {
		t.Error("Generate should not be called when user already exists")
	}
}

// TestRegisterUseCase_DatabaseErrors тестирует различные ошибки базы данных
func TestRegisterUseCase_DatabaseErrors(t *testing.T) {
	testCases := []struct {
		name           string
		setupMocks     func(*mocks.MockUserRepository, *authMocks.MockPasswordHasher, *authMocks.MockTokenGenerator)
		expectedError  error
		description    string
	}{
		{
			name: "exists_check_error",
			setupMocks: func(repo *mocks.MockUserRepository, hasher *authMocks.MockPasswordHasher, tokenGen *authMocks.MockTokenGenerator) {
				// Ошибка при проверке существования пользователя
				repo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
					return false, errors.New("database connection error")
				}
			},
			expectedError: errors.New("database connection error"),
			description:   "Ошибка БД при проверке существования пользователя",
		},
		{
			name: "save_error",
			setupMocks: func(repo *mocks.MockUserRepository, hasher *authMocks.MockPasswordHasher, tokenGen *authMocks.MockTokenGenerator) {
				repo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
					return false, nil
				}
				hasher.HashFunc = func(password string) (string, error) {
					return "hashed_password", nil
				}
				// Ошибка при сохранении пользователя
				repo.SaveFunc = func(ctx context.Context, user *user.User) error {
					return errors.New("save error")
				}
			},
			expectedError: errors.New("save error"),
			description:   "Ошибка БД при сохранении пользователя",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			mockRepo := &mocks.MockUserRepository{}
			mockHasher := &authMocks.MockPasswordHasher{}
			mockTokenGen := &authMocks.MockTokenGenerator{}
			
			tc.setupMocks(mockRepo, mockHasher, mockTokenGen)
			
			uc := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
			
			req := RegisterRequest{
				Email:    "test@example.com",
				Password: "validpassword123",
				Name:     "Test User",
			}
			
			// ACT
			ctx := context.Background()
			response, err := uc.Execute(ctx, req)
			
			// ASSERT
			if err == nil {
				t.Errorf("Expected error, got nil. %s", tc.description)
			}
			
			if response != nil {
				t.Errorf("Expected nil response on error, got %+v", response)
			}
		})
	}
}

// TestRegisterUseCase_PasswordHashingError тестирует ошибку хеширования пароля
func TestRegisterUseCase_PasswordHashingError(t *testing.T) {
	// ARRANGE
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	mockRepo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
		return false, nil
	}
	
	// Настраиваем ошибку хеширования
	hashError := errors.New("hashing failed")
	mockHasher.HashFunc = func(password string) (string, error) {
		return "", hashError
	}
	
	uc := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "validpassword123",
		Name:     "Test User",
	}
	
	// ACT
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	
	// ASSERT
	if err != hashError {
		t.Errorf("Expected hash error, got %v", err)
	}
	
	if response != nil {
		t.Errorf("Expected nil response on hash error, got %+v", response)
	}
	
	// Проверяем, что Save НЕ был вызван
	if mockRepo.SaveCallCount > 0 {
		t.Error("Save should not be called when hashing fails")
	}
}

// TestRegisterUseCase_TokenGenerationError тестирует ошибку генерации токена
func TestRegisterUseCase_TokenGenerationError(t *testing.T) {
	// ARRANGE
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	mockRepo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
		return false, nil
	}
	
	mockHasher.HashFunc = func(password string) (string, error) {
		return "hashed_password", nil
	}
	
	// Настраиваем ошибку генерации токена
	tokenError := errors.New("token generation failed")
	mockTokenGen.GenerateFunc = func(userID uint, email string) (string, error) {
		return "", tokenError
	}
	
	uc := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "validpassword123",
		Name:     "Test User",
	}
	
	// ACT
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	
	// ASSERT
	if err != tokenError {
		t.Errorf("Expected token error, got %v", err)
	}
	
	if response != nil {
		t.Errorf("Expected nil response on token error, got %+v", response)
	}
	
	// Проверяем, что пользователь ВСЕ ЕЩЕ был сохранен
	// (это важно для консистентности данных)
	if mockRepo.SaveCallCount != 1 {
		t.Errorf("Save should be called even if token generation fails, called %d times", mockRepo.SaveCallCount)
	}
}

// TestRegisterUseCase_InvalidDomainData тестирует валидацию доменных данных
func TestRegisterUseCase_InvalidDomainData(t *testing.T) {
	// ARRANGE
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	mockRepo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
		return false, nil
	}
	
	uc := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
	
	testCases := []struct {
		name        string
		request     RegisterRequest
		description string
	}{
		{
			name: "invalid_email",
			request: RegisterRequest{
				Email:    "invalid-email",
				Password: "validpassword123",
				Name:     "Test User",
			},
			description: "Невалидный email должен вызывать доменную ошибку",
		},
		{
			name: "empty_name",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "validpassword123",
				Name:     "",
			},
			description: "Пустое имя должно вызывать доменную ошибку",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			ctx := context.Background()
			response, err := uc.Execute(ctx, tc.request)
			
			// ASSERT
			if err == nil {
				t.Errorf("Expected domain validation error, got nil. %s", tc.description)
			}
			
			if response != nil {
				t.Errorf("Expected nil response on validation error, got %+v", response)
			}
			
			// Проверяем, что никакие внешние сервисы НЕ были вызваны
			if mockHasher.HashCallCount > 0 {
				t.Error("Hash should not be called on validation error")
			}
			
			if mockRepo.SaveCallCount > 0 {
				t.Error("Save should not be called on validation error")
			}
		})
	}
}

// ИНТЕГРАЦИОННЫЙ ТЕСТ с реальными компонентами
// 📚 ИЗУЧИТЬ: Integration vs Unit tests
// https://martinfowler.com/bliki/IntegrationTest.html

// TestRegisterUseCase_Integration тестирует интеграцию с реальными сервисами
// ВНИМАНИЕ: Этот тест использует реальные bcrypt и JWT - он медленнее!
func TestRegisterUseCase_Integration(t *testing.T) {
	// Пропускаем интеграционные тесты при быстром запуске
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// ARRANGE - используем реальные сервисы
	mockRepo := &mocks.MockUserRepository{}
	
	// TODO: Здесь можно подключить реальные сервисы
	// realHasher := &external.BcryptPasswordHasher{}
	// realTokenGen := &external.JWTTokenGenerator{Secret: "test-secret"}
	
	// Пока используем моки
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	mockRepo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
		return false, nil
	}
	
	uc := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := RegisterRequest{
		Email:    "integration@example.com",
		Password: "integrationtest123",
		Name:     "Integration Test User",
	}
	
	// ACT
	start := time.Now()
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	duration := time.Since(start)
	
	// ASSERT
	if err != nil {
		t.Errorf("Integration test failed: %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response in integration test")
	}
	
	// Проверяем производительность
	if duration > time.Second {
		t.Errorf("Registration took too long: %v", duration)
	}
	
	t.Logf("Integration test completed in %v", duration)
}

// БЕНЧМАРК для измерения производительности Use Case
func BenchmarkRegisterUseCase_Execute(b *testing.B) {
	// ARRANGE
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	mockRepo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
		return false, nil
	}
	
	uc := NewRegisterUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := RegisterRequest{
		Email:    "benchmark@example.com",
		Password: "benchmarkpassword123",
		Name:     "Benchmark User",
	}
	
	ctx := context.Background()
	
	// Сбрасываем таймер
	b.ResetTimer()
	
	// Выполняем бенчмарк
	for i := 0; i < b.N; i++ {
		// Сбрасываем состояние моков для каждой итерации
		mockRepo.Reset()
		mockHasher.Reset()
		mockTokenGen.Reset()
		
		mockRepo.ExistsFunc = func(ctx context.Context, email string) (bool, error) {
			return false, nil
		}
		
		_, err := uc.Execute(ctx, req)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}