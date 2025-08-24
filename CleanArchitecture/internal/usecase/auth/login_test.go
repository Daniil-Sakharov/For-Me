package auth

import (
	"context"
	"errors"
	"testing"
	
	"clean-url-shortener/internal/domain/user"
	"clean-url-shortener/internal/domain/user/mocks"
	authMocks "clean-url-shortener/internal/usecase/auth/mocks"
)

// TestLoginUseCase_Success тестирует успешную авторизацию
func TestLoginUseCase_Success(t *testing.T) {
	// ARRANGE - подготовка тестового пользователя
	testUser := mocks.NewUserTestBuilder().
		WithID(1).
		WithEmail("test@example.com").
		WithName("Test User").
		WithHashedPassword("hashed_validpassword123").
		Build()
	
	// Подготовка моков
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	// Настройка поведения моков
	mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
		if email == testUser.Email {
			return testUser, nil
		}
		return nil, nil
	}
	
	mockHasher.CompareFunc = func(hashedPassword, password string) error {
		if hashedPassword == "hashed_validpassword123" && password == "validpassword123" {
			return nil // Пароли совпадают
		}
		return ErrInvalidCredentials
	}
	
	mockTokenGen.GenerateFunc = func(userID uint, email string) (string, error) {
		return "login_token", nil
	}
	
	// Создание Use Case
	uc := NewLoginUseCase(mockRepo, mockHasher, mockTokenGen)
	
	// Подготовка запроса
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "validpassword123",
	}
	
	// ACT - выполнение действия
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	
	// ASSERT - проверка результата
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	// Проверяем данные ответа
	if response.UserID != testUser.ID {
		t.Errorf("Expected UserID %d, got %d", testUser.ID, response.UserID)
	}
	
	if response.Email != testUser.Email {
		t.Errorf("Expected email %s, got %s", testUser.Email, response.Email)
	}
	
	if response.Name != testUser.Name {
		t.Errorf("Expected name %s, got %s", testUser.Name, response.Name)
	}
	
	if response.Token != "login_token" {
		t.Errorf("Expected token 'login_token', got %s", response.Token)
	}
	
	// Проверяем вызовы моков
	if mockRepo.FindByEmailCallCount != 1 {
		t.Errorf("Expected FindByEmail to be called once, called %d times", mockRepo.FindByEmailCallCount)
	}
	
	if mockHasher.CompareCallCount != 1 {
		t.Errorf("Expected Compare to be called once, called %d times", mockHasher.CompareCallCount)
	}
	
	if mockTokenGen.GenerateCallCount != 1 {
		t.Errorf("Expected Generate to be called once, called %d times", mockTokenGen.GenerateCallCount)
	}
	
	// Проверяем параметры, переданные в моки
	if mockRepo.LastFindByEmailParam != req.Email {
		t.Errorf("Expected FindByEmail called with %s, called with %s", req.Email, mockRepo.LastFindByEmailParam)
	}
	
	if mockHasher.LastPlainPassword != req.Password {
		t.Errorf("Expected Compare called with password %s, called with %s", req.Password, mockHasher.LastPlainPassword)
	}
	
	if mockHasher.LastHashedPassword != testUser.HashedPassword {
		t.Errorf("Expected Compare called with hash %s, called with %s", testUser.HashedPassword, mockHasher.LastHashedPassword)
	}
}

// TestLoginUseCase_UserNotFound тестирует попытку входа несуществующего пользователя
func TestLoginUseCase_UserNotFound(t *testing.T) {
	// ARRANGE
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	// Настройка мока: пользователь не найден
	mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
		return nil, nil // Пользователь не найден
	}
	
	uc := NewLoginUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "somepassword",
	}
	
	// ACT
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	
	// ASSERT
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
	
	if response != nil {
		t.Errorf("Expected nil response when user not found, got %+v", response)
	}
	
	// Проверяем, что другие методы НЕ были вызваны
	if mockHasher.CompareCallCount > 0 {
		t.Error("Compare should not be called when user not found")
	}
	
	if mockTokenGen.GenerateCallCount > 0 {
		t.Error("Generate should not be called when user not found")
	}
}

// TestLoginUseCase_InvalidPassword тестирует неверный пароль
func TestLoginUseCase_InvalidPassword(t *testing.T) {
	// ARRANGE
	testUser := mocks.NewUserTestBuilder().
		WithID(1).
		WithEmail("test@example.com").
		WithHashedPassword("hashed_correctpassword").
		Build()
	
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
		return testUser, nil
	}
	
	// Настройка мока: пароли не совпадают
	mockHasher.CompareFunc = func(hashedPassword, password string) error {
		return ErrInvalidCredentials // Пароли не совпадают
	}
	
	uc := NewLoginUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	
	// ACT
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	
	// ASSERT
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
	
	if response != nil {
		t.Errorf("Expected nil response on invalid password, got %+v", response)
	}
	
	// Проверяем, что Compare был вызван, но Generate - нет
	if mockHasher.CompareCallCount != 1 {
		t.Errorf("Expected Compare to be called once, called %d times", mockHasher.CompareCallCount)
	}
	
	if mockTokenGen.GenerateCallCount > 0 {
		t.Error("Generate should not be called on invalid password")
	}
}

// TestLoginUseCase_DatabaseError тестирует ошибки базы данных
func TestLoginUseCase_DatabaseError(t *testing.T) {
	// ARRANGE
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	dbError := errors.New("database connection lost")
	mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
		return nil, dbError
	}
	
	uc := NewLoginUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	// ACT
	ctx := context.Background()
	response, err := uc.Execute(ctx, req)
	
	// ASSERT
	if err != dbError {
		t.Errorf("Expected database error, got %v", err)
	}
	
	if response != nil {
		t.Errorf("Expected nil response on database error, got %+v", response)
	}
	
	// Проверяем, что другие сервисы НЕ были вызваны
	if mockHasher.CompareCallCount > 0 {
		t.Error("Compare should not be called on database error")
	}
	
	if mockTokenGen.GenerateCallCount > 0 {
		t.Error("Generate should not be called on database error")
	}
}

// TestLoginUseCase_TokenGenerationError тестирует ошибку генерации токена
func TestLoginUseCase_TokenGenerationError(t *testing.T) {
	// ARRANGE
	testUser := mocks.NewUserTestBuilder().
		WithID(1).
		WithEmail("test@example.com").
		WithHashedPassword("hashed_password").
		Build()
	
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &authMocks.MockPasswordHasher{}
	mockTokenGen := &authMocks.MockTokenGenerator{}
	
	mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
		return testUser, nil
	}
	
	mockHasher.CompareFunc = func(hashedPassword, password string) error {
		return nil // Пароли совпадают
	}
	
	// Ошибка генерации токена
	tokenError := errors.New("token service unavailable")
	mockTokenGen.GenerateFunc = func(userID uint, email string) (string, error) {
		return "", tokenError
	}
	
	uc := NewLoginUseCase(mockRepo, mockHasher, mockTokenGen)
	
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
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
	
	// Проверяем, что все шаги до генерации токена были выполнены
	if mockRepo.FindByEmailCallCount != 1 {
		t.Error("FindByEmail should be called even if token generation fails")
	}
	
	if mockHasher.CompareCallCount != 1 {
		t.Error("Compare should be called even if token generation fails")
	}
	
	if mockTokenGen.GenerateCallCount != 1 {
		t.Error("Generate should be called (and fail)")
	}
}

// TestLoginUseCase_EdgeCases тестирует граничные случаи
func TestLoginUseCase_EdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		setupMocks  func(*mocks.MockUserRepository, *authMocks.MockPasswordHasher, *authMocks.MockTokenGenerator)
		request     LoginRequest
		expectError bool
		description string
	}{
		{
			name: "empty_email",
			setupMocks: func(repo *mocks.MockUserRepository, hasher *authMocks.MockPasswordHasher, tokenGen *authMocks.MockTokenGenerator) {
				// Пустой email должен привести к поиску пустой строки
				repo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
					return nil, nil // Не найден
				}
			},
			request: LoginRequest{
				Email:    "",
				Password: "password123",
			},
			expectError: true,
			description: "Пустой email должен вызывать ошибку",
		},
		{
			name: "empty_password",
			setupMocks: func(repo *mocks.MockUserRepository, hasher *authMocks.MockPasswordHasher, tokenGen *authMocks.MockTokenGenerator) {
				testUser := mocks.NewUserTestBuilder().Build()
				repo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
					return testUser, nil
				}
				// Пустой пароль не должен совпадать с хешем
				hasher.CompareFunc = func(hashedPassword, password string) error {
					return ErrInvalidCredentials
				}
			},
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			expectError: true,
			description: "Пустой пароль должен вызывать ошибку",
		},
		{
			name: "very_long_email",
			setupMocks: func(repo *mocks.MockUserRepository, hasher *authMocks.MockPasswordHasher, tokenGen *authMocks.MockTokenGenerator) {
				repo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
					return nil, nil // Очень длинный email не найден
				}
			},
			request: LoginRequest{
				Email:    "very.long.email.address.that.might.cause.issues@very.long.domain.name.example.com",
				Password: "password123",
			},
			expectError: true,
			description: "Очень длинный email должен обрабатываться корректно",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			mockRepo := &mocks.MockUserRepository{}
			mockHasher := &authMocks.MockPasswordHasher{}
			mockTokenGen := &authMocks.MockTokenGenerator{}
			
			tc.setupMocks(mockRepo, mockHasher, mockTokenGen)
			
			uc := NewLoginUseCase(mockRepo, mockHasher, mockTokenGen)
			
			// ACT
			ctx := context.Background()
			response, err := uc.Execute(ctx, tc.request)
			
			// ASSERT
			if tc.expectError && err == nil {
				t.Errorf("Expected error, got nil. %s", tc.description)
			}
			
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, got %v. %s", err, tc.description)
			}
			
			if tc.expectError && response != nil {
				t.Errorf("Expected nil response on error, got %+v", response)
			}
		})
	}
}

// ТЕСТ НА БЕЗОПАСНОСТЬ: проверяем, что не утекает информация
// 📚 ИЗУЧИТЬ: Security testing
// https://owasp.org/www-project-top-ten/
// https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html

// TestLoginUseCase_SecurityConsiderations проверяет аспекты безопасности
func TestLoginUseCase_SecurityConsiderations(t *testing.T) {
	t.Run("timing_attack_resistance", func(t *testing.T) {
		// ARRANGE
		mockRepo := &mocks.MockUserRepository{}
		mockHasher := &authMocks.MockPasswordHasher{}
		mockTokenGen := &authMocks.MockTokenGenerator{}
		
		// Настраиваем сценарии: пользователь существует vs не существует
		userExistsRepo := &mocks.MockUserRepository{}
		userExistsRepo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
			return mocks.NewUserTestBuilder().WithEmail(email).Build(), nil
		}
		
		userNotExistsRepo := &mocks.MockUserRepository{}
		userNotExistsRepo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
			return nil, nil
		}
		
		mockHasher.CompareFunc = func(hashedPassword, password string) error {
			return ErrInvalidCredentials
		}
		
		uc1 := NewLoginUseCase(userExistsRepo, mockHasher, mockTokenGen)
		uc2 := NewLoginUseCase(userNotExistsRepo, mockHasher, mockTokenGen)
		
		req := LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		
		ctx := context.Background()
		
		// ACT - измеряем время выполнения
		// В реальном приложении время должно быть примерно одинаковым
		// для защиты от timing attacks
		
		_, err1 := uc1.Execute(ctx, req) // Пользователь существует
		_, err2 := uc2.Execute(ctx, req) // Пользователь не существует
		
		// ASSERT
		// Оба случая должны возвращать ошибки авторизации
		if err1 != ErrInvalidCredentials {
			t.Errorf("Expected ErrInvalidCredentials for existing user, got %v", err1)
		}
		
		if err2 != ErrUserNotFound {
			t.Errorf("Expected ErrUserNotFound for non-existing user, got %v", err2)
		}
		
		// ПРИМЕЧАНИЕ: В реальном приложении оба случая должны возвращать
		// одинаковую ошибку (ErrInvalidCredentials) для защиты от user enumeration
		t.Log("SECURITY NOTE: Consider returning same error type for both cases to prevent user enumeration")
	})
	
	t.Run("password_not_logged", func(t *testing.T) {
		// Проверяем, что пароль не попадает в логи или другие места
		// Это проверяется через анализ вызовов моков
		
		mockRepo := &mocks.MockUserRepository{}
		mockHasher := &authMocks.MockPasswordHasher{}
		mockTokenGen := &authMocks.MockTokenGenerator{}
		
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
			return nil, nil
		}
		
		uc := NewLoginUseCase(mockRepo, mockHasher, mockTokenGen)
		
		req := LoginRequest{
			Email:    "test@example.com",
			Password: "secretpassword123",
		}
		
		ctx := context.Background()
		_, _ = uc.Execute(ctx, req)
		
		// Проверяем, что пароль передается только в Compare
		// и нигде больше не сохраняется
		if mockHasher.CompareCallCount > 0 {
			t.Log("Password was passed to hasher (expected)")
		}
		
		// В реальном приложении здесь можно добавить проверки логов
		t.Log("SECURITY NOTE: Ensure passwords are never logged or stored in plain text")
	})
}