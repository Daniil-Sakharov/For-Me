package link

import (
	"context"
	"errors"
	"clean-url-shortener/internal/domain/link"
)

// RedirectRequest представляет запрос на редирект по короткой ссылке
type RedirectRequest struct {
	ShortCode  string `json:"short_code" validate:"required"`
	UserAgent  string `json:"user_agent"`
	IPAddress  string `json:"ip_address"`
	Referer    string `json:"referer"`
}

// RedirectResponse представляет ответ с URL для редиректа
type RedirectResponse struct {
	OriginalURL string `json:"original_url"`
	LinkID      uint   `json:"link_id"`
}

// Ошибки для редиректа
var (
	ErrLinkNotFound = errors.New("link not found")
	ErrLinkExpired  = errors.New("link has expired")
)

// RedirectUseCase содержит бизнес-логику редиректа по короткой ссылке
type RedirectUseCase struct {
	linkRepo       link.Repository // Из domain слоя
	eventPublisher EventPublisher  // Из usecase слоя
}

// NewRedirectUseCase создает новый Use Case для редиректа
func NewRedirectUseCase(
	linkRepo link.Repository,
	eventPublisher EventPublisher,
) *RedirectUseCase {
	return &RedirectUseCase{
		linkRepo:       linkRepo,
		eventPublisher: eventPublisher,
	}
}

// Execute выполняет бизнес-логику редиректа по короткой ссылке
func (uc *RedirectUseCase) Execute(ctx context.Context, req RedirectRequest) (*RedirectResponse, error) {
	// ШАГ 1: Находим ссылку по короткому коду
	foundLink, err := uc.linkRepo.FindByShortCode(ctx, req.ShortCode)
	if err != nil {
		return nil, err
	}
	if foundLink == nil {
		return nil, ErrLinkNotFound
	}

	// ШАГ 2: Увеличиваем счетчик кликов в доменной модели
	// Это доменная операция - каждый переход увеличивает счетчик
	foundLink.IncrementClicks()

	// ШАГ 3: Обновляем счетчик в базе данных
	// Можно сделать это асинхронно для улучшения производительности
	err = uc.linkRepo.IncrementClicks(ctx, foundLink.ID)
	if err != nil {
		// Логируем ошибку, но не прерываем редирект
		// Пользователь должен быть перенаправлен даже если счетчик не обновился
		// TODO: логирование ошибки
	}

	// ШАГ 4: Публикуем событие о клике для сбора статистики
	// Это делается асинхронно через event publisher
	err = uc.eventPublisher.PublishLinkClicked(
		foundLink.ID,
		req.UserAgent,
		req.IPAddress,
		req.Referer,
	)
	if err != nil {
		// Логируем ошибку, но не прерываем редирект
		// Статистика важна, но не критична для основной функциональности
		// TODO: логирование ошибки
	}

	// ШАГ 5: Возвращаем URL для редиректа
	return &RedirectResponse{
		OriginalURL: foundLink.OriginalURL,
		LinkID:      foundLink.ID,
	}, nil
}

// ПРИНЦИПЫ ДИЗАЙНА:
// 1. ПРОИЗВОДИТЕЛЬНОСТЬ - редирект должен быть максимально быстрым
// 2. НАДЕЖНОСТЬ - даже если статистика не работает, редирект должен работать
// 3. АСИНХРОННОСТЬ - счетчики и события обрабатываются асинхронно
// 4. ОТКАЗОУСТОЙЧИВОСТЬ - ошибки в дополнительной функциональности не ломают основную