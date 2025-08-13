package stat

import (
	"context"
	"clean-url-shortener/internal/domain/stat"
)

// TrackClickUseCase содержит бизнес-логику отслеживания кликов по ссылкам
// Это Use Case обрабатывает события кликов и сохраняет статистику
type TrackClickUseCase struct {
	statRepo        stat.Repository    // Из domain слоя
	geoService      GeoLocationService // Из usecase слоя
	eventSubscriber EventSubscriber    // Из usecase слоя
}

// NewTrackClickUseCase создает новый Use Case для отслеживания кликов
func NewTrackClickUseCase(
	statRepo stat.Repository,
	geoService GeoLocationService,
	eventSubscriber EventSubscriber,
) *TrackClickUseCase {
	return &TrackClickUseCase{
		statRepo:        statRepo,
		geoService:      geoService,
		eventSubscriber: eventSubscriber,
	}
}

// StartTracking запускает отслеживание кликов по ссылкам
// Это долго работающий процесс, который слушает события и обрабатывает их
func (uc *TrackClickUseCase) StartTracking(ctx context.Context) error {
	// Подписываемся на события кликов
	eventChan, err := uc.eventSubscriber.SubscribeLinkClicked()
	if err != nil {
		return err
	}

	// Запускаем обработку событий в горутине
	go func() {
		defer uc.eventSubscriber.Close()

		for {
			select {
			case <-ctx.Done():
				// Контекст отменен, завершаем работу
				return

			case event, ok := <-eventChan:
				if !ok {
					// Канал закрыт, завершаем работу
					return
				}

				// Обрабатываем событие клика
				err := uc.processClickEvent(ctx, event)
				if err != nil {
					// TODO: логирование ошибки
					// В продакшене здесь должно быть полноценное логирование
					continue
				}
			}
		}
	}()

	return nil
}

// processClickEvent обрабатывает одно событие клика
func (uc *TrackClickUseCase) processClickEvent(ctx context.Context, event LinkClickedEvent) error {
	// ШАГ 1: Создаем доменную модель статистики
	clickStat := stat.NewStat(
		event.LinkID,
		event.UserAgent,
		event.IPAddress,
		event.Referer,
	)

	// ШАГ 2: Обогащаем данные геолокацией
	if uc.geoService.IsValidIP(event.IPAddress) {
		country, city, err := uc.geoService.GetLocation(event.IPAddress)
		if err == nil {
			clickStat.SetLocation(country, city)
		}
		// Если геолокация не работает, продолжаем без нее
	}

	// ШАГ 3: Сохраняем статистику в базе данных
	err := uc.statRepo.Save(ctx, clickStat)
	if err != nil {
		return err
	}

	return nil
}

// TrackSingleClick обрабатывает единичный клик синхронно
// Используется для прямого вызова без событийной системы
func (uc *TrackClickUseCase) TrackSingleClick(
	ctx context.Context,
	linkID uint,
	userAgent, ipAddress, referer string,
) error {
	// Создаем событие
	event := LinkClickedEvent{
		LinkID:    linkID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		Referer:   referer,
	}

	// Обрабатываем его
	return uc.processClickEvent(ctx, event)
}

// АРХИТЕКТУРНЫЕ ПРИНЦИПЫ:
// 1. АСИНХРОННАЯ ОБРАБОТКА - события обрабатываются в фоне
// 2. ОТКАЗОУСТОЙЧИВОСТЬ - ошибки в статистике не влияют на основную функциональность
// 3. РАЗДЕЛЕНИЕ ОТВЕТСТВЕННОСТИ - Use Case координирует, domain модели содержат логику
// 4. ТЕСТИРУЕМОСТЬ - все зависимости инжектируются через интерфейсы