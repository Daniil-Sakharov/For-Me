// =====================================================================
// 🏛️ ДОМЕННАЯ СУЩНОСТЬ TENDER - Центральная бизнес-модель системы
// =====================================================================
//
// Пакет tender содержит доменную модель тендера - основную сущность системы.
// Эта сущность представляет тендер на закупку медицинского оборудования.
//
// ПРИНЦИПЫ CLEAN ARCHITECTURE:
// 1. Содержит ТОЛЬКО бизнес-логику, никаких зависимостей от внешних систем
// 2. Определяет правила и ограничения предметной области
// 3. Не знает о базах данных, API, UI или других технических деталях
// 4. Может быть протестирована изолированно
// 5. Изменяется только при изменении бизнес-требований
//
// TODO: При расширении функциональности добавить:
// - Доменные события (TenderCreated, TenderAnalyzed и т.д.)
// - Более сложную бизнес-логику валидации
// - Агрегаты для связанных сущностей (Products, Documents)
// - Value Objects для типизированных полей

package tender

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// =====================================================================
// 📊 ОСНОВНАЯ ДОМЕННАЯ СУЩНОСТЬ TENDER
// =====================================================================

// Tender представляет тендер на закупку товаров/услуг
//
// TODO: Реализовать следующую бизнес-логику:
// 1. Валидация всех полей при создании и обновлении
// 2. Бизнес-правила (например, нельзя анализировать неактивный тендер)
// 3. Вычисляемые поля (DaysUntilDeadline, IsExpired и т.д.)
// 4. Методы для изменения состояния (SetAIAnalysis, MarkAsAnalyzed)
// 5. Проверки бизнес-инвариантов
//
// ВАЖНО: Эта структура НЕ должна содержать теги для сериализации (json, db),
// так как это зона ответственности внешних слоев (interfaces, infrastructure)
type Tender struct {
	// 🔑 Уникальные идентификаторы
	ID         uint   // Внутренний ID в нашей системе
	ExternalID string // ID тендера на внешней платформе (zakupki.gov.ru и т.д.)

	// 📝 Основная информация
	Title       string // Название тендера
	Description string // Подробное описание

	// 🌐 Информация об источнике
	Platform string // Платформа-источник (zakupki, szvo, spb)
	URL      string // Прямая ссылка на тендер

	// 👤 Информация о заказчике
	Customer    string // Наименование заказчика
	CustomerINN string // ИНН заказчика

	// 💰 Финансовая информация
	StartPrice float64  // Начальная (максимальная) цена
	Currency   Currency // Валюта тендера

	// 📅 Временные рамки
	PublishedAt time.Time  // Дата публикации тендера
	DeadlineAt  *time.Time // Крайний срок подачи заявок

	// 🏷️ Классификация
	Status   TenderStatus // Текущий статус тендера
	Category string       // Категория товаров/услуг

	// 🤖 Результаты AI анализа
	AIScore            *float64         // Оценка релевантности (0.0-1.0)
	AIRecommendation   *AIRecommendation // Рекомендация AI
	AIAnalysisReason   string           // Обоснование решения AI
	AIAnalyzedAt       *time.Time       // Время проведения анализа

	// 📊 Служебные поля
	CreatedAt time.Time // Время создания записи
	UpdatedAt time.Time // Время последнего обновления
}

// =====================================================================
// 🏷️ ДОМЕННЫЕ ЭНУМЫ И VALUE OBJECTS
// =====================================================================

// TenderStatus представляет возможные статусы тендера
//
// TODO: Добавить методы для проверки возможных переходов между статусами
// TODO: Добавить бизнес-логику для автоматического изменения статусов
type TenderStatus string

const (
	StatusActive    TenderStatus = "active"    // Активный тендер (идет прием заявок)
	StatusCompleted TenderStatus = "completed" // Завершенный тендер
	StatusCancelled TenderStatus = "cancelled" // Отмененный тендер
	StatusExpired   TenderStatus = "expired"   // Истекший тендер
	StatusDraft     TenderStatus = "draft"     // Черновик (еще не опубликован)
)

// AIRecommendation представляет рекомендацию AI по участию в тендере
//
// TODO: Добавить более детальные рекомендации (например, "high_priority", "low_priority")
// TODO: Добавить уровень уверенности в рекомендации
type AIRecommendation string

const (
	RecommendationParticipate AIRecommendation = "participate" // Рекомендуется участвовать
	RecommendationSkip        AIRecommendation = "skip"        // Пропустить тендер
	RecommendationAnalyze     AIRecommendation = "analyze"     // Требует дополнительного анализа
)

// Currency представляет валюту тендера
//
// TODO: Добавить методы для конвертации валют
// TODO: Добавить валидацию поддерживаемых валют
type Currency string

const (
	CurrencyRUB Currency = "RUB" // Российский рубль
	CurrencyUSD Currency = "USD" // Доллар США
	CurrencyEUR Currency = "EUR" // Евро
)

// Platform представляет платформу-источник тендера
//
// TODO: Вынести в отдельный value object с валидацией
// TODO: Добавить специфичные для платформы настройки
type Platform string

const (
	PlatformZakupki Platform = "zakupki" // zakupki.gov.ru
	PlatformSZVO    Platform = "szvo"    // szvo.gov35.ru
	PlatformSPB     Platform = "spb"     // gz-spb.ru
)

// =====================================================================
// 🏗️ КОНСТРУКТОР И МЕТОДЫ СОЗДАНИЯ
// =====================================================================

// NewTender создает новый тендер с валидацией обязательных полей
//
// TODO: Добавить дополнительные параметры для полной инициализации
// TODO: Реализовать паттерн Builder для сложных случаев создания
// TODO: Добавить создание доменных событий (TenderCreated)
//
// Параметры:
//   - externalID: уникальный ID на внешней платформе
//   - title: название тендера
//   - platform: платформа-источник
//   - url: ссылка на тендер
//
// Возвращает:
//   - *Tender: созданный тендер
//   - error: ошибка валидации
func NewTender(externalID, title, platform, url string) (*Tender, error) {
	// TODO: Реализовать полную валидацию входных параметров
	if err := validateExternalID(externalID); err != nil {
		return nil, fmt.Errorf("invalid external ID: %w", err)
	}

	if err := validateTitle(title); err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}

	if err := validatePlatform(platform); err != nil {
		return nil, fmt.Errorf("invalid platform: %w", err)
	}

	if err := validateURL(url); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	now := time.Now()

	tender := &Tender{
		ExternalID:  externalID,
		Title:       strings.TrimSpace(title),
		Platform:    platform,
		URL:         url,
		Status:      StatusActive,
		Currency:    CurrencyRUB, // Значение по умолчанию
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: now, // TODO: Получать из внешнего источника
	}

	// TODO: Добавить создание доменного события TenderCreated
	// TODO: Добавить генерацию внутреннего UUID если нужно

	return tender, nil
}

// =====================================================================
// 🔍 МЕТОДЫ ЗАПРОСОВ (QUERY METHODS)
// =====================================================================

// IsActive проверяет, активен ли тендер для подачи заявок
//
// TODO: Добавить проверку дополнительных условий (например, минимальная цена)
// TODO: Учитывать рабочие дни и часы для более точной проверки
func (t *Tender) IsActive() bool {
	if t.Status != StatusActive {
		return false
	}

	// Если дедлайн не установлен, считаем что тендер активен
	if t.DeadlineAt == nil {
		return true
	}

	// Проверяем что дедлайн еще не истек
	return time.Now().Before(*t.DeadlineAt)
}

// DaysUntilDeadline возвращает количество дней до истечения срока подачи заявок
//
// TODO: Добавить расчет рабочих дней (исключая выходные и праздники)
// TODO: Добавить учет часовых поясов
func (t *Tender) DaysUntilDeadline() int {
	if t.DeadlineAt == nil {
		return -1 // Дедлайн не установлен
	}

	duration := t.DeadlineAt.Sub(time.Now())
	if duration < 0 {
		return 0 // Дедлайн уже истек
	}

	return int(duration.Hours() / 24)
}

// CanBeAnalyzed проверяет, можно ли проводить AI анализ тендера
//
// TODO: Добавить дополнительные бизнес-правила для анализа
// TODO: Проверить наличие необходимых данных для анализа
func (t *Tender) CanBeAnalyzed() bool {
	// Нельзя анализировать уже проанализированные тендеры
	if t.AIAnalyzedAt != nil {
		return false
	}

	// Нельзя анализировать неактивные тендеры
	if !t.IsActive() {
		return false
	}

	// Должны быть заполнены минимально необходимые поля
	return t.Title != "" && t.Platform != ""
}

// IsRelevant проверяет, является ли тендер релевантным на основе AI анализа
//
// TODO: Добавить настраиваемый порог релевантности
// TODO: Учитывать дополнительные факторы (цена, заказчик и т.д.)
func (t *Tender) IsRelevant() bool {
	if t.AIScore == nil {
		return false
	}

	// Порог релевантности - TODO: вынести в конфигурацию
	const relevanceThreshold = 0.7
	return *t.AIScore >= relevanceThreshold
}

// ShouldParticipate возвращает рекомендацию об участии в тендере
//
// TODO: Добавить комплексную логику принятия решения
// TODO: Учитывать исторические данные и статистику
func (t *Tender) ShouldParticipate() bool {
	if t.AIRecommendation == nil {
		return false
	}

	return *t.AIRecommendation == RecommendationParticipate
}

// =====================================================================
// 🔄 МЕТОДЫ ИЗМЕНЕНИЯ СОСТОЯНИЯ (COMMAND METHODS)
// =====================================================================

// SetAIAnalysis устанавливает результаты AI анализа
//
// TODO: Добавить валидацию входных параметров
// TODO: Добавить создание доменного события AIAnalysisCompleted
// TODO: Добавить проверку прав на изменение
//
// Параметры:
//   - score: оценка релевантности от 0.0 до 1.0
//   - recommendation: рекомендация AI
//   - reason: обоснование решения
func (t *Tender) SetAIAnalysis(score float64, recommendation AIRecommendation, reason string) error {
	// TODO: Реализовать полную валидацию
	if score < 0.0 || score > 1.0 {
		return ErrInvalidAIScore
	}

	if !isValidRecommendation(recommendation) {
		return ErrInvalidAIRecommendation
	}

	// Устанавливаем результаты анализа
	t.AIScore = &score
	t.AIRecommendation = &recommendation
	t.AIAnalysisReason = reason
	now := time.Now()
	t.AIAnalyzedAt = &now
	t.UpdatedAt = now

	// TODO: Добавить создание доменного события
	// TODO: Добавить логирование изменения состояния

	return nil
}

// UpdateStatus изменяет статус тендера с проверкой допустимых переходов
//
// TODO: Реализовать state machine для контроля переходов
// TODO: Добавить создание доменных событий для изменения статуса
func (t *Tender) UpdateStatus(newStatus TenderStatus) error {
	// TODO: Реализовать логику проверки допустимых переходов
	if !isValidStatusTransition(t.Status, newStatus) {
		return ErrInvalidStatusTransition
	}

	oldStatus := t.Status
	t.Status = newStatus
	t.UpdatedAt = time.Now()

	// TODO: Добавить создание доменного события StatusChanged
	// TODO: Добавить логирование перехода статуса

	_ = oldStatus // Используется для логирования
	return nil
}

// =====================================================================
// 🛡️ МЕТОДЫ ВАЛИДАЦИИ
// =====================================================================

// Validate проверяет корректность всех полей тендера
//
// TODO: Добавить валидацию всех полей
// TODO: Добавить кастомные правила валидации
// TODO: Интегрировать с библиотекой валидации
func (t *Tender) Validate() error {
	var errors []string

	// Проверяем обязательные поля
	if t.ExternalID == "" {
		errors = append(errors, "external ID is required")
	}

	if t.Title == "" {
		errors = append(errors, "title is required")
	}

	if t.Platform == "" {
		errors = append(errors, "platform is required")
	}

	// Проверяем валидность URL
	if t.URL != "" {
		if err := validateURL(t.URL); err != nil {
			errors = append(errors, fmt.Sprintf("invalid URL: %v", err))
		}
	}

	// Проверяем финансовые данные
	if t.StartPrice < 0 {
		errors = append(errors, "start price cannot be negative")
	}

	// Проверяем даты
	if t.DeadlineAt != nil && t.DeadlineAt.Before(t.PublishedAt) {
		errors = append(errors, "deadline cannot be before published date")
	}

	// Проверяем AI анализ
	if t.AIScore != nil && (*t.AIScore < 0 || *t.AIScore > 1) {
		errors = append(errors, "AI score must be between 0 and 1")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}

// =====================================================================
// 🔧 ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ВАЛИДАЦИИ
// =====================================================================

// validateExternalID проверяет корректность внешнего ID
func validateExternalID(externalID string) error {
	externalID = strings.TrimSpace(externalID)
	if externalID == "" {
		return ErrEmptyExternalID
	}

	// TODO: Добавить проверку формата ID для разных платформ
	// TODO: Добавить проверку уникальности (через repository)

	return nil
}

// validateTitle проверяет корректность названия тендера
func validateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrEmptyTitle
	}

	// TODO: Добавить проверку максимальной длины
	// TODO: Добавить проверку запрещенных символов

	return nil
}

// validatePlatform проверяет корректность платформы
func validatePlatform(platform string) error {
	validPlatforms := map[string]bool{
		string(PlatformZakupki): true,
		string(PlatformSZVO):    true,
		string(PlatformSPB):     true,
	}

	if !validPlatforms[platform] {
		return ErrInvalidPlatform
	}

	return nil
}

// validateURL проверяет корректность URL
func validateURL(rawURL string) error {
	if rawURL == "" {
		return ErrEmptyURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

// isValidRecommendation проверяет корректность AI рекомендации
func isValidRecommendation(recommendation AIRecommendation) bool {
	validRecommendations := map[AIRecommendation]bool{
		RecommendationParticipate: true,
		RecommendationSkip:        true,
		RecommendationAnalyze:     true,
	}

	return validRecommendations[recommendation]
}

// isValidStatusTransition проверяет допустимость перехода между статусами
//
// TODO: Реализовать полную логику state machine
func isValidStatusTransition(from, to TenderStatus) bool {
	// Базовые правила переходов (TODO: расширить)
	validTransitions := map[TenderStatus][]TenderStatus{
		StatusDraft:     {StatusActive, StatusCancelled},
		StatusActive:    {StatusCompleted, StatusCancelled, StatusExpired},
		StatusCompleted: {}, // Финальное состояние
		StatusCancelled: {}, // Финальное состояние
		StatusExpired:   {}, // Финальное состояние
	}

	allowedStatuses, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == to {
			return true
		}
	}

	return false
}

// =====================================================================
// 🔧 ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ ДЛЯ УДОБСТВА
// =====================================================================

// String возвращает строковое представление тендера
//
// TODO: Добавить более информативное представление
func (t *Tender) String() string {
	return fmt.Sprintf("Tender{ID: %d, ExternalID: %s, Title: %s, Platform: %s, Status: %s}",
		t.ID, t.ExternalID, t.Title, t.Platform, t.Status)
}

// Clone создает глубокую копию тендера
//
// TODO: Реализовать полное клонирование всех полей
func (t *Tender) Clone() *Tender {
	clone := *t // Поверхностное копирование

	// Копируем указатели
	if t.DeadlineAt != nil {
		deadline := *t.DeadlineAt
		clone.DeadlineAt = &deadline
	}

	if t.AIScore != nil {
		score := *t.AIScore
		clone.AIScore = &score
	}

	if t.AIRecommendation != nil {
		recommendation := *t.AIRecommendation
		clone.AIRecommendation = &recommendation
	}

	if t.AIAnalyzedAt != nil {
		analyzedAt := *t.AIAnalyzedAt
		clone.AIAnalyzedAt = &analyzedAt
	}

	return &clone
}

// =====================================================================
// 🎯 МЕТОДЫ ДЛЯ БИЗНЕС-ЛОГИКИ (можно расширить)
// =====================================================================

// TODO: Добавить методы для расширенной бизнес-логики:
//
// - CalculatePriority() int - расчет приоритета тендера
// - GetRiskLevel() RiskLevel - оценка рисков участия
// - EstimateParticipationCost() float64 - оценка стоимости участия
// - GetCompetitionLevel() CompetitionLevel - уровень конкуренции
// - ShouldSendNotification() bool - нужно ли отправлять уведомление
// - GetRecommendedActions() []Action - рекомендуемые действия
// - CalculateTimeToParticipate() time.Duration - время на подготовку
//
// Примеры методов для событий:
// - OnCreated() - обработка создания тендера
// - OnAnalyzed() - обработка завершения анализа
// - OnStatusChanged() - обработка изменения статуса
//
// Примеры методов для интеграций:
// - ToDTO() TenderDTO - конвертация в DTO для API
// - FromExternalData() - создание из внешних данных
// - ToAnalysisRequest() - подготовка данных для AI анализа