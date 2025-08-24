package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// 📚 ИЗУЧИТЬ: JWT Refresh Tokens
// https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/
// https://datatracker.ietf.org/doc/html/rfc6749#section-1.5
// https://tools.ietf.org/html/draft-ietf-oauth-security-topics-13
//
// ЗАЧЕМ НУЖНЫ REFRESH TOKENS:
// 1. БЕЗОПАСНОСТЬ - Access tokens имеют короткий срок жизни
// 2. ОТЗЫВ - Можно отозвать refresh token без ожидания истечения access token
// 3. АУДИТ - Можно отслеживать активность пользователей
// 4. РОТАЦИЯ - Можно обновлять токены без повторного ввода пароля
//
// ПРИНЦИПЫ БЕЗОПАСНОСТИ:
// 1. Refresh tokens должны быть длинными и случайными
// 2. Refresh tokens должны храниться в безопасном месте (БД, не в localStorage!)
// 3. Refresh tokens должны иметь срок жизни
// 4. Refresh tokens должны ротироваться при использовании
// 5. Refresh tokens должны быть привязаны к конкретному пользователю и устройству

// RefreshToken представляет доменную модель refresh token
// ЭТО ЧИСТАЯ ДОМЕННАЯ СУЩНОСТЬ - содержит только бизнес-логику
type RefreshToken struct {
	// ID уникальный идентификатор токена
	ID uint
	
	// Token сам токен (длинная случайная строка)
	Token string
	
	// UserID ID пользователя, которому принадлежит токен
	UserID uint
	
	// DeviceInfo информация об устройстве для безопасности
	DeviceInfo string
	
	// ExpiresAt время истечения токена
	ExpiresAt time.Time
	
	// CreatedAt время создания токена
	CreatedAt time.Time
	
	// LastUsedAt время последнего использования
	LastUsedAt *time.Time
	
	// IsRevoked флаг отзыва токена
	IsRevoked bool
	
	// RevokedAt время отзыва токена
	RevokedAt *time.Time
	
	// RevokedReason причина отзыва токена
	RevokedReason string
}

// Доменные ошибки для refresh tokens
var (
	ErrRefreshTokenExpired    = errors.New("refresh token has expired")
	ErrRefreshTokenRevoked    = errors.New("refresh token has been revoked")
	ErrRefreshTokenInvalid    = errors.New("refresh token is invalid")
	ErrRefreshTokenNotFound   = errors.New("refresh token not found")
	ErrInvalidDeviceInfo      = errors.New("invalid device information")
	ErrTokenGenerationFailed  = errors.New("failed to generate secure token")
)

// NewRefreshToken создает новый refresh token
// 📚 ИЗУЧИТЬ: Cryptographically secure random generation
// https://pkg.go.dev/crypto/rand
func NewRefreshToken(userID uint, deviceInfo string, expiryDuration time.Duration) (*RefreshToken, error) {
	// Валидация входных параметров
	if userID == 0 {
		return nil, errors.New("user ID cannot be zero")
	}
	
	if deviceInfo == "" {
		return nil, ErrInvalidDeviceInfo
	}
	
	if expiryDuration <= 0 {
		return nil, errors.New("expiry duration must be positive")
	}
	
	// Генерируем криптографически стойкий токен
	token, err := generateSecureToken(32) // 32 байта = 256 бит
	if err != nil {
		return nil, ErrTokenGenerationFailed
	}
	
	now := time.Now()
	
	return &RefreshToken{
		Token:      token,
		UserID:     userID,
		DeviceInfo: deviceInfo,
		ExpiresAt:  now.Add(expiryDuration),
		CreatedAt:  now,
		IsRevoked:  false,
	}, nil
}

// IsValid проверяет валидность refresh token
func (rt *RefreshToken) IsValid() error {
	// Проверяем, не отозван ли токен
	if rt.IsRevoked {
		return ErrRefreshTokenRevoked
	}
	
	// Проверяем, не истек ли токен
	if time.Now().After(rt.ExpiresAt) {
		return ErrRefreshTokenExpired
	}
	
	// Проверяем базовые поля
	if rt.Token == "" || rt.UserID == 0 {
		return ErrRefreshTokenInvalid
	}
	
	return nil
}

// Use отмечает использование токена
func (rt *RefreshToken) Use() error {
	// Проверяем валидность перед использованием
	if err := rt.IsValid(); err != nil {
		return err
	}
	
	now := time.Now()
	rt.LastUsedAt = &now
	
	return nil
}

// Revoke отзывает токен с указанием причины
func (rt *RefreshToken) Revoke(reason string) {
	if rt.IsRevoked {
		return // Уже отозван
	}
	
	now := time.Now()
	rt.IsRevoked = true
	rt.RevokedAt = &now
	rt.RevokedReason = reason
}

// IsExpiringSoon проверяет, истекает ли токен в ближайшее время
func (rt *RefreshToken) IsExpiringSoon(threshold time.Duration) bool {
	return time.Until(rt.ExpiresAt) < threshold
}

// DaysUntilExpiry возвращает количество дней до истечения
func (rt *RefreshToken) DaysUntilExpiry() int {
	duration := time.Until(rt.ExpiresAt)
	return int(duration.Hours() / 24)
}

// REPOSITORY INTERFACE для refresh tokens

// RefreshTokenRepository определяет интерфейс для работы с refresh tokens
// ИНТЕРФЕЙС В ДОМЕННОМ СЛОЕ - определяет ЧТО нужно делать, но НЕ КАК
type RefreshTokenRepository interface {
	// Save сохраняет refresh token
	Save(token *RefreshToken) error
	
	// FindByToken находит токен по значению
	FindByToken(tokenValue string) (*RefreshToken, error)
	
	// FindByUserID находит все токены пользователя
	FindByUserID(userID uint) ([]*RefreshToken, error)
	
	// Delete удаляет токен
	Delete(tokenID uint) error
	
	// DeleteByUserID удаляет все токены пользователя
	DeleteByUserID(userID uint) error
	
	// DeleteExpired удаляет истекшие токены
	DeleteExpired() (int, error)
	
	// RevokeAllByUserID отзывает все токены пользователя
	RevokeAllByUserID(userID uint, reason string) error
	
	// RevokeByDevice отзывает токены конкретного устройства
	RevokeByDevice(userID uint, deviceInfo string, reason string) error
}

// UTILITY FUNCTIONS

// generateSecureToken генерирует криптографически стойкий токен
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	
	// Используем crypto/rand для генерации случайных байтов
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	
	// Кодируем в hex для удобства хранения и передачи
	return hex.EncodeToString(bytes), nil
}

// ДОМЕННЫЕ СЕРВИСЫ для работы с refresh tokens

// RefreshTokenService доменный сервис для работы с refresh tokens
// 📚 ИЗУЧИТЬ: Domain Services в DDD
// https://martinfowler.com/bliki/EvansClassification.html
type RefreshTokenService struct {
	repository RefreshTokenRepository
}

// NewRefreshTokenService создает новый доменный сервис
func NewRefreshTokenService(repository RefreshTokenRepository) *RefreshTokenService {
	return &RefreshTokenService{
		repository: repository,
	}
}

// CreateToken создает новый refresh token для пользователя
func (s *RefreshTokenService) CreateToken(userID uint, deviceInfo string, expiryDuration time.Duration) (*RefreshToken, error) {
	// Создаем новый токен
	token, err := NewRefreshToken(userID, deviceInfo, expiryDuration)
	if err != nil {
		return nil, err
	}
	
	// Сохраняем в репозитории
	if err := s.repository.Save(token); err != nil {
		return nil, err
	}
	
	return token, nil
}

// RotateToken создает новый токен и отзывает старый (Token Rotation)
// 📚 ИЗУЧИТЬ: Refresh Token Rotation
// https://auth0.com/docs/secure/tokens/refresh-tokens/refresh-token-rotation
func (s *RefreshTokenService) RotateToken(oldTokenValue string, deviceInfo string, expiryDuration time.Duration) (*RefreshToken, error) {
	// Находим старый токен
	oldToken, err := s.repository.FindByToken(oldTokenValue)
	if err != nil {
		return nil, err
	}
	
	if oldToken == nil {
		return nil, ErrRefreshTokenNotFound
	}
	
	// Проверяем валидность старого токена
	if err := oldToken.IsValid(); err != nil {
		return nil, err
	}
	
	// Отмечаем использование старого токена
	if err := oldToken.Use(); err != nil {
		return nil, err
	}
	
	// Создаем новый токен
	newToken, err := NewRefreshToken(oldToken.UserID, deviceInfo, expiryDuration)
	if err != nil {
		return nil, err
	}
	
	// Сохраняем новый токен
	if err := s.repository.Save(newToken); err != nil {
		return nil, err
	}
	
	// Отзываем старый токен
	oldToken.Revoke("rotated")
	if err := s.repository.Save(oldToken); err != nil {
		// Логируем ошибку, но не прерываем операцию
		// Новый токен уже создан
	}
	
	return newToken, nil
}

// ValidateToken проверяет валидность токена
func (s *RefreshTokenService) ValidateToken(tokenValue string) (*RefreshToken, error) {
	token, err := s.repository.FindByToken(tokenValue)
	if err != nil {
		return nil, err
	}
	
	if token == nil {
		return nil, ErrRefreshTokenNotFound
	}
	
	if err := token.IsValid(); err != nil {
		return nil, err
	}
	
	return token, nil
}

// RevokeUserTokens отзывает все токены пользователя
func (s *RefreshTokenService) RevokeUserTokens(userID uint, reason string) error {
	return s.repository.RevokeAllByUserID(userID, reason)
}

// RevokeDeviceTokens отзывает токены конкретного устройства
func (s *RefreshTokenService) RevokeDeviceTokens(userID uint, deviceInfo string, reason string) error {
	return s.repository.RevokeByDevice(userID, deviceInfo, reason)
}

// CleanupExpiredTokens удаляет истекшие токены
func (s *RefreshTokenService) CleanupExpiredTokens() (int, error) {
	return s.repository.DeleteExpired()
}

// GetUserTokens возвращает все активные токены пользователя
func (s *RefreshTokenService) GetUserTokens(userID uint) ([]*RefreshToken, error) {
	allTokens, err := s.repository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	// Фильтруем только активные токены
	var activeTokens []*RefreshToken
	for _, token := range allTokens {
		if token.IsValid() == nil {
			activeTokens = append(activeTokens, token)
		}
	}
	
	return activeTokens, nil
}

// КОНФИГУРАЦИЯ для refresh tokens

// RefreshTokenConfig конфигурация refresh tokens
type RefreshTokenConfig struct {
	// Срок жизни refresh token (обычно 7-30 дней)
	ExpiryDuration time.Duration `yaml:"expiry_duration" env:"REFRESH_TOKEN_EXPIRY" default:"168h"` // 7 дней
	
	// Максимальное количество активных токенов на пользователя
	MaxTokensPerUser int `yaml:"max_tokens_per_user" env:"REFRESH_TOKEN_MAX_PER_USER" default:"5"`
	
	// Включить ротацию токенов
	EnableRotation bool `yaml:"enable_rotation" env:"REFRESH_TOKEN_ENABLE_ROTATION" default:"true"`
	
	// Интервал очистки истекших токенов
	CleanupInterval time.Duration `yaml:"cleanup_interval" env:"REFRESH_TOKEN_CLEANUP_INTERVAL" default:"24h"`
	
	// Предупреждать об истечении за N дней
	ExpiryWarningDays int `yaml:"expiry_warning_days" env:"REFRESH_TOKEN_EXPIRY_WARNING_DAYS" default:"3"`
}

// СОБЫТИЯ ДОМЕНА для refresh tokens
// 📚 ИЗУЧИТЬ: Domain Events
// https://martinfowler.com/eaaDev/DomainEvent.html

// RefreshTokenCreated событие создания токена
type RefreshTokenCreated struct {
	TokenID    uint      `json:"token_id"`
	UserID     uint      `json:"user_id"`
	DeviceInfo string    `json:"device_info"`
	CreatedAt  time.Time `json:"created_at"`
}

// RefreshTokenUsed событие использования токена
type RefreshTokenUsed struct {
	TokenID  uint      `json:"token_id"`
	UserID   uint      `json:"user_id"`
	UsedAt   time.Time `json:"used_at"`
}

// RefreshTokenRevoked событие отзыва токена
type RefreshTokenRevoked struct {
	TokenID       uint      `json:"token_id"`
	UserID        uint      `json:"user_id"`
	RevokedAt     time.Time `json:"revoked_at"`
	RevokedReason string    `json:"revoked_reason"`
}

// RefreshTokenRotated событие ротации токена
type RefreshTokenRotated struct {
	OldTokenID uint      `json:"old_token_id"`
	NewTokenID uint      `json:"new_token_id"`
	UserID     uint      `json:"user_id"`
	RotatedAt  time.Time `json:"rotated_at"`
}

// ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ:

/*
// В Use Case:
type RefreshTokenUseCase struct {
    tokenService *auth.RefreshTokenService
    jwtGenerator JWTGenerator
}

func (uc *RefreshTokenUseCase) RefreshAccessToken(refreshTokenValue string, deviceInfo string) (*TokenPair, error) {
    // Валидируем refresh token
    refreshToken, err := uc.tokenService.ValidateToken(refreshTokenValue)
    if err != nil {
        return nil, err
    }
    
    // Генерируем новый access token
    accessToken, err := uc.jwtGenerator.Generate(refreshToken.UserID)
    if err != nil {
        return nil, err
    }
    
    // Ротируем refresh token (опционально)
    newRefreshToken, err := uc.tokenService.RotateToken(refreshTokenValue, deviceInfo, 7*24*time.Hour)
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken.Token,
    }, nil
}

// В Controller:
func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
    var req RefreshTokenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    deviceInfo := extractDeviceInfo(r) // User-Agent, IP, etc.
    
    tokens, err := c.refreshTokenUC.RefreshAccessToken(req.RefreshToken, deviceInfo)
    if err != nil {
        if errors.Is(err, auth.ErrRefreshTokenExpired) {
            http.Error(w, "Token expired", http.StatusUnauthorized)
            return
        }
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }
    
    json.NewEncoder(w).Encode(tokens)
}

// Периодическая очистка:
func (s *RefreshTokenService) StartCleanupWorker(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            deleted, err := s.CleanupExpiredTokens()
            if err != nil {
                log.Printf("Failed to cleanup expired tokens: %v", err)
            } else if deleted > 0 {
                log.Printf("Cleaned up %d expired refresh tokens", deleted)
            }
        case <-ctx.Done():
            return
        }
    }
}
*/

// ВАЖНЫЕ ПРИНЦИПЫ БЕЗОПАСНОСТИ:
// 1. Refresh tokens должны храниться в БД, НЕ в localStorage браузера
// 2. Используйте HTTPS для всех операций с токенами
// 3. Реализуйте ротацию токенов для повышения безопасности
// 4. Ограничьте количество активных токенов на пользователя
// 5. Логируйте все операции с токенами для аудита
// 6. Предоставьте пользователю возможность отозвать токены
// 7. Очищайте истекшие токены регулярно