package link

import (
	"context"
	"errors"
	"clean-url-shortener/internal/domain/link"
)

// CreateLinkRequest представляет запрос на создание короткой ссылки
type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	CustomCode  string `json:"custom_code,omitempty" validate:"omitempty,min=3,max=10"`
	UserID      uint   `json:"-"` // Передается из контекста аутентификации
}

// CreateLinkResponse представляет ответ при создании ссылки
type CreateLinkResponse struct {
	ID          uint   `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"` // Полный URL для использования
	CreatedAt   string `json:"created_at"`
}

// Ошибки для создания ссылок
var (
	ErrShortCodeAlreadyExists = errors.New("short code already exists")
	ErrUnsafeURL             = errors.New("URL is not safe")
	ErrInvalidCustomCode     = errors.New("invalid custom code")
)

// CreateLinkUseCase содержит бизнес-логику создания короткой ссылки
type CreateLinkUseCase struct {
	linkRepo          link.Repository    // Из domain слоя
	urlValidator      URLValidator       // Из usecase слоя
	shortCodeGen      ShortCodeGenerator // Из usecase слоя
	baseURL           string             // Базовый URL для формирования коротких ссылок
}

// NewCreateLinkUseCase создает новый Use Case для создания ссылок
func NewCreateLinkUseCase(
	linkRepo link.Repository,
	urlValidator URLValidator,
	shortCodeGen ShortCodeGenerator,
	baseURL string,
) *CreateLinkUseCase {
	return &CreateLinkUseCase{
		linkRepo:     linkRepo,
		urlValidator: urlValidator,
		shortCodeGen: shortCodeGen,
		baseURL:      baseURL,
	}
}

// Execute выполняет бизнес-логику создания короткой ссылки
func (uc *CreateLinkUseCase) Execute(ctx context.Context, req CreateLinkRequest) (*CreateLinkResponse, error) {
	// ШАГ 1: Валидируем URL на корректность и безопасность
	err := uc.urlValidator.Validate(req.OriginalURL)
	if err != nil {
		return nil, err
	}

	// ШАГ 2: Проверяем безопасность URL (не фишинг, не malware)
	isSafe, err := uc.urlValidator.IsSafe(req.OriginalURL)
	if err != nil {
		return nil, err
	}
	if !isSafe {
		return nil, ErrUnsafeURL
	}

	var newLink *link.Link

	// ШАГ 3: Создаем ссылку с пользовательским или автоматическим кодом
	if req.CustomCode != "" {
		// Пользователь предоставил свой код
		newLink, err = uc.createLinkWithCustomCode(ctx, req)
		if err != nil {
			return nil, err
		}
	} else {
		// Генерируем автоматический код
		newLink, err = uc.createLinkWithAutoCode(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	// ШАГ 4: Сохраняем ссылку в базе данных
	err = uc.linkRepo.Save(ctx, newLink)
	if err != nil {
		return nil, err
	}

	// ШАГ 5: Формируем ответ
	return &CreateLinkResponse{
		ID:          newLink.ID,
		OriginalURL: newLink.OriginalURL,
		ShortCode:   newLink.ShortCode,
		ShortURL:    uc.baseURL + "/" + newLink.ShortCode,
		CreatedAt:   newLink.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// createLinkWithCustomCode создает ссылку с пользовательским кодом
func (uc *CreateLinkUseCase) createLinkWithCustomCode(ctx context.Context, req CreateLinkRequest) (*link.Link, error) {
	// Проверяем уникальность пользовательского кода
	exists, err := uc.linkRepo.ExistsByShortCode(ctx, req.CustomCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrShortCodeAlreadyExists
	}

	// Валидируем пользовательский код через генератор
	validatedCode, err := uc.shortCodeGen.GenerateCustom(req.CustomCode)
	if err != nil {
		return nil, ErrInvalidCustomCode
	}

	// Создаем ссылку с пользовательским кодом
	return link.NewLinkWithCustomCode(req.OriginalURL, validatedCode, req.UserID)
}

// createLinkWithAutoCode создает ссылку с автоматически сгенерированным кодом
func (uc *CreateLinkUseCase) createLinkWithAutoCode(ctx context.Context, req CreateLinkRequest) (*link.Link, error) {
	// Генерируем уникальный код (попытки до 10 раз)
	const maxAttempts = 10
	
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Генерируем новый код
		shortCode, err := uc.shortCodeGen.Generate()
		if err != nil {
			return nil, err
		}

		// Проверяем уникальность
		exists, err := uc.linkRepo.ExistsByShortCode(ctx, shortCode)
		if err != nil {
			return nil, err
		}

		if !exists {
			// Найден уникальный код, создаем ссылку
			return link.NewLinkWithCustomCode(req.OriginalURL, shortCode, req.UserID)
		}
	}

	// Если за 10 попыток не удалось сгенерировать уникальный код
	return nil, errors.New("failed to generate unique short code")
}

// ПРИНЦИПЫ, КОТОРЫЕ МЫ СОБЛЮДАЕМ:
// 1. Валидация входных данных на уровне Use Case
// 2. Проверка бизнес-правил (уникальность короткого кода)
// 3. Координация между различными сервисами
// 4. Обработка ошибок и edge cases
// 5. Безопасность (проверка на фишинг/malware)