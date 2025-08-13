package link

import (
	"crypto/rand"
	"errors"
	"net/url"
	"time"
)

// Link представляет доменную модель короткой ссылки
// ЧИСТАЯ ДОМЕННАЯ СУЩНОСТЬ без зависимостей от внешних библиотек
type Link struct {
	// ID - уникальный идентификатор ссылки
	ID uint

	// OriginalURL - исходный URL, который нужно сократить
	OriginalURL string

	// ShortCode - короткий код для ссылки (например, "abc123")
	ShortCode string

	// UserID - ID пользователя, создавшего ссылку
	UserID uint

	// CreatedAt - время создания ссылки
	CreatedAt time.Time

	// UpdatedAt - время последнего обновления
	UpdatedAt time.Time

	// ClicksCount - количество переходов по ссылке
	ClicksCount uint
}

// Доменные ошибки
var (
	ErrInvalidURL       = errors.New("invalid URL format")
	ErrEmptyURL         = errors.New("URL cannot be empty")
	ErrInvalidShortCode = errors.New("invalid short code")
	ErrEmptyShortCode   = errors.New("short code cannot be empty")
)

// NewLink создает новую ссылку с автоматически сгенерированным коротким кодом
func NewLink(originalURL string, userID uint) (*Link, error) {
	// Валидация URL согласно доменным правилам
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	// Генерируем уникальный короткий код
	shortCode, err := generateShortCode()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Link{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
		ClicksCount: 0,
	}, nil
}

// NewLinkWithCustomCode создает новую ссылку с пользовательским коротким кодом
func NewLinkWithCustomCode(originalURL, shortCode string, userID uint) (*Link, error) {
	// Валидация URL
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	// Валидация пользовательского короткого кода
	if err := validateShortCode(shortCode); err != nil {
		return nil, err
	}

	now := time.Now()

	return &Link{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
		ClicksCount: 0,
	}, nil
}

// UpdateURL обновляет исходный URL
func (l *Link) UpdateURL(newURL string) error {
	if err := validateURL(newURL); err != nil {
		return err
	}

	l.OriginalURL = newURL
	l.UpdatedAt = time.Now()
	return nil
}

// UpdateShortCode обновляет короткий код
func (l *Link) UpdateShortCode(newShortCode string) error {
	if err := validateShortCode(newShortCode); err != nil {
		return err
	}

	l.ShortCode = newShortCode
	l.UpdatedAt = time.Now()
	return nil
}

// IncrementClicks увеличивает счетчик переходов на 1
// Это доменная операция - каждый переход по ссылке увеличивает счетчик
func (l *Link) IncrementClicks() {
	l.ClicksCount++
	l.UpdatedAt = time.Now()
}

// IsOwner проверяет, принадлежит ли ссылка указанному пользователю
func (l *Link) IsOwner(userID uint) bool {
	return l.UserID == userID
}

// validateURL проверяет валидность URL согласно доменным правилам
func validateURL(rawURL string) error {
	if rawURL == "" {
		return ErrEmptyURL
	}

	// Парсим URL для проверки его корректности
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	// Проверяем, что URL имеет схему (http или https)
	if parsedURL.Scheme == "" {
		return ErrInvalidURL
	}

	// Проверяем, что URL имеет хост
	if parsedURL.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

// validateShortCode проверяет валидность короткого кода
func validateShortCode(shortCode string) error {
	if shortCode == "" {
		return ErrEmptyShortCode
	}

	// Короткий код должен содержать только буквы и цифры
	for _, r := range shortCode {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return ErrInvalidShortCode
		}
	}

	// Длина кода должна быть от 3 до 10 символов
	if len(shortCode) < 3 || len(shortCode) > 10 {
		return ErrInvalidShortCode
	}

	return nil
}

// generateShortCode генерирует случайный короткий код
func generateShortCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	bytes := make([]byte, codeLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}

// RegenerateShortCode генерирует новый короткий код для существующей ссылки
func (l *Link) RegenerateShortCode() error {
	newCode, err := generateShortCode()
	if err != nil {
		return err
	}

	l.ShortCode = newCode
	l.UpdatedAt = time.Now()
	return nil
}