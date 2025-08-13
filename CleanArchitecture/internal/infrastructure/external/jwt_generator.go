package external

import (
	"fmt"
	"time"
	
	"clean-url-shortener/internal/usecase/auth"
	"github.com/golang-jwt/jwt/v5"
)

// JWTTokenGenerator реализует интерфейс TokenGenerator используя JWT
type JWTTokenGenerator struct {
	secretKey []byte
	issuer    string
	expiry    time.Duration
}

// NewJWTTokenGenerator создает новый JWT генератор
func NewJWTTokenGenerator(secretKey, issuer string, expiry time.Duration) auth.TokenGenerator {
	if expiry == 0 {
		expiry = 24 * time.Hour // По умолчанию 24 часа
	}
	
	return &JWTTokenGenerator{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		expiry:    expiry,
	}
}

// Claims структура для JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Generate создает JWT токен для пользователя
func (g *JWTTokenGenerator) Generate(userID uint, email string) (string, error) {
	now := time.Now()
	
	// Создаем claims
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.expiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	
	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Подписываем токен
	tokenString, err := token.SignedString(g.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	
	return tokenString, nil
}

// Validate проверяет валидность токена и возвращает данные
func (g *JWTTokenGenerator) Validate(tokenString string) (*auth.TokenData, error) {
	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return g.secretKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	// Проверяем валидность токена
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	// Извлекаем claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	
	// Возвращаем данные из токена
	return &auth.TokenData{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}

// ПРИНЦИПЫ:
// 1. Инкапсулирует все детали работы с JWT
// 2. Обеспечивает безопасность (проверка алгоритма подписи)
// 3. Настраиваемое время жизни токенов
// 4. Стандартные JWT claims для совместимости