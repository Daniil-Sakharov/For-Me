// 🏛️ ДОМЕННАЯ МОДЕЛЬ ТЕНДЕРА для системы автоматизации
//
// Представляет тендер на медицинское оборудование из государственных
// закупочных систем (zakupki.gov.ru, szvo.gov35.ru и др.)

package tender

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// TODO: Основная доменная сущность тендера
// Представляет чистую доменную модель без зависимостей от внешних систем
// Содержит всю бизнес-логику, связанную с тендером
type Tender struct {
	// TODO: Базовые поля тендера
	ID          uint      `json:"id"`
	TenderID    string    `json:"tender_id"`    // Номер тендера на площадке
	Title       string    `json:"title"`        // Название тендера
	Description string    `json:"description"`  // Описание закупки
	Platform    string    `json:"platform"`     // Платформа (zakupki, szvo, spb)
	URL         string    `json:"url"`          // Ссылка на тендер

	// TODO: Заказчик и контактная информация
	Customer        string `json:"customer"`         // Наименование заказчика
	CustomerINN     string `json:"customer_inn"`     // ИНН заказчика
	CustomerAddress string `json:"customer_address"` // Адрес заказчика
	ContactPerson   string `json:"contact_person"`   // Контактное лицо
	ContactPhone    string `json:"contact_phone"`    // Телефон
	ContactEmail    string `json:"contact_email"`    // Email

	// TODO: Финансовая информация
	StartPrice    float64 `json:"start_price"`    // Начальная (максимальная) цена
	Currency      string  `json:"currency"`       // Валюта (обычно RUB)
	PriceType     string  `json:"price_type"`     // Тип цены (за весь объем/за единицу)
	VAT           bool    `json:"vat"`            // НДС включен/исключен
	EstimatedCost float64 `json:"estimated_cost"` // Ориентировочная стоимость

	// TODO: Временные рамки
	PublishedAt        time.Time  `json:"published_at"`        // Дата публикации
	ApplicationDeadline time.Time `json:"application_deadline"` // Дедлайн подачи заявок
	ContractStartDate  *time.Time `json:"contract_start_date"` // Дата начала поставки
	ContractEndDate    *time.Time `json:"contract_end_date"`   // Дата окончания поставки
	AuctionDate        *time.Time `json:"auction_date"`        // Дата аукциона (если есть)

	// TODO: Категории и теги
	Category       string   `json:"category"`        // Основная категория
	Subcategory    string   `json:"subcategory"`     // Подкатегория
	Keywords       []string `json:"keywords"`        // Ключевые слова
	Tags           []string `json:"tags"`            // Теги для фильтрации
	ProductTypes   []string `json:"product_types"`   // Типы продукции
	MedicalCategory string  `json:"medical_category"` // Медицинская категория

	// TODO: Статус тендера
	Status           TenderStatus `json:"status"`            // Текущий статус
	ProcessingStage  string       `json:"processing_stage"`  // Этап обработки в нашей системе
	ParticipantCount int          `json:"participant_count"` // Количество участников
	IsActive         bool         `json:"is_active"`         // Активен ли тендер

	// TODO: AI анализ
	AIRelevanceScore    float64           `json:"ai_relevance_score"`    // Оценка релевантности (0-1)
	AIRecommendation    AIRecommendation  `json:"ai_recommendation"`     // Рекомендация ИИ
	AIConfidenceLevel   float64           `json:"ai_confidence_level"`   // Уровень уверенности ИИ
	AIAnalysisReason    string            `json:"ai_analysis_reason"`    // Обоснование анализа
	CompetitionLevel    CompetitionLevel  `json:"competition_level"`     // Уровень конкуренции
	RiskFactors         []string          `json:"risk_factors"`          // Выявленные риски
	Opportunities       []string          `json:"opportunities"`         // Возможности

	// TODO: Документы
	DocumentsURL       string    `json:"documents_url"`       // Ссылка на документы
	DocumentsCount     int       `json:"documents_count"`     // Количество документов
	TechnicalTaskFound bool      `json:"technical_task_found"` // Найдено ли техзадание
	DocumentsProcessed bool      `json:"documents_processed"` // Обработаны ли документы
	LastDocumentCheck  time.Time `json:"last_document_check"` // Последняя проверка документов

	// TODO: Извлеченные продукты
	ProductsExtracted     bool `json:"products_extracted"`      // Извлечены ли товары
	ProductsCount         int  `json:"products_count"`          // Количество извлеченных товаров
	ProductExtractionDate *time.Time `json:"product_extraction_date"` // Дата извлечения

	// TODO: Поставщики и цены
	SuppliersFound       bool       `json:"suppliers_found"`       // Найдены ли поставщики
	SuppliersCount       int        `json:"suppliers_count"`       // Количество найденных поставщиков
	EmailCampaignSent    bool       `json:"email_campaign_sent"`   // Отправлена ли email кампания
	EmailResponsesCount  int        `json:"email_responses_count"` // Количество ответов
	LastEmailSent        *time.Time `json:"last_email_sent"`       // Дата последнего письма
	OptimalPriceCalculated bool     `json:"optimal_price_calculated"` // Рассчитана ли оптимальная цена
	RecommendedPrice     *float64   `json:"recommended_price"`     // Рекомендуемая цена для участия

	// TODO: Результаты участия
	ParticipationDecision string     `json:"participation_decision"` // Решение об участии (yes/no/pending)
	ApplicationSubmitted  bool       `json:"application_submitted"`  // Подана ли заявка
	WinProbability       float64    `json:"win_probability"`        // Вероятность выигрыша (0-1)
	ActualResult         *string    `json:"actual_result"`          // Фактический результат (won/lost/cancelled)
	WinnerPrice          *float64   `json:"winner_price"`           // Цена победителя
	ResultAnalyzed       bool       `json:"result_analyzed"`        // Проанализирован ли результат

	// TODO: Метаданные
	CreatedAt time.Time  `json:"created_at"` // Дата создания записи в системе
	UpdatedAt time.Time  `json:"updated_at"` // Дата последнего обновления
	DeletedAt *time.Time `json:"deleted_at"` // Дата удаления (soft delete)
	
	// TODO: Связи с другими сущностями (будут определены в репозитории)
	// Products []Product - товары из техзадания
	// Suppliers []Supplier - найденные поставщики  
	// EmailCampaigns []EmailCampaign - email кампании
	// Documents []Document - загруженные документы
}

// TODO: Перечисление статусов тендера
type TenderStatus string

const (
	// TODO: Статусы на стороне закупочной площадки
	TenderStatusPublished    TenderStatus = "published"     // Опубликован
	TenderStatusApplications TenderStatus = "applications"  // Прием заявок
	TenderStatusAuction      TenderStatus = "auction"       // Аукцион
	TenderStatusEvaluation   TenderStatus = "evaluation"    // Рассмотрение заявок
	TenderStatusCompleted    TenderStatus = "completed"     // Завершен
	TenderStatusCancelled    TenderStatus = "cancelled"     // Отменен

	// TODO: Статусы в нашей системе
	TenderStatusNew              TenderStatus = "new"                // Новый тендер
	TenderStatusAIAnalysis       TenderStatus = "ai_analysis"        // AI анализ
	TenderStatusDocumentsPending TenderStatus = "documents_pending"  // Ожидание документов
	TenderStatusProcessing       TenderStatus = "processing"         // Обработка документов
	TenderStatusSupplierSearch   TenderStatus = "supplier_search"    // Поиск поставщиков
	TenderStatusEmailCampaign    TenderStatus = "email_campaign"     // Email кампания
	TenderStatusPriceAnalysis    TenderStatus = "price_analysis"     // Анализ цен
	TenderStatusReady            TenderStatus = "ready"              // Готов к участию
	TenderStatusParticipating    TenderStatus = "participating"      // Участвуем
	TenderStatusResultsPending   TenderStatus = "results_pending"    // Ожидание результатов
	TenderStatusResultsAnalyzed  TenderStatus = "results_analyzed"   // Результаты проанализированы
	TenderStatusArchived         TenderStatus = "archived"           // Архивирован
)

// TODO: Рекомендации AI по участию в тендере
type AIRecommendation string

const (
	AIRecommendationParticipate AIRecommendation = "participate" // Рекомендуется участвовать
	AIRecommendationSkip        AIRecommendation = "skip"        // Пропустить
	AIRecommendationAnalyze     AIRecommendation = "analyze"     // Требует дополнительного анализа
	AIRecommendationWatchlist   AIRecommendation = "watchlist"   // Добавить в список наблюдения
)

// TODO: Уровень конкуренции в тендере
type CompetitionLevel string

const (
	CompetitionLevelLow    CompetitionLevel = "low"    // Низкая конкуренция (1-3 участника)
	CompetitionLevelMedium CompetitionLevel = "medium" // Средняя конкуренция (4-7 участников)
	CompetitionLevelHigh   CompetitionLevel = "high"   // Высокая конкуренция (8+ участников)
)

// TODO: Бизнес-методы доменной сущности

// IsActive проверяет, активен ли тендер (можно ли еще подать заявку)
func (t *Tender) IsActiveForParticipation() bool {
	now := time.Now()
	return t.IsActive && 
		   t.ApplicationDeadline.After(now) && 
		   (t.Status == TenderStatusPublished || t.Status == TenderStatusApplications)
}

// DaysUntilDeadline возвращает количество дней до дедлайна подачи заявок
func (t *Tender) DaysUntilDeadline() int {
	if t.ApplicationDeadline.IsZero() {
		return -1
	}
	duration := t.ApplicationDeadline.Sub(time.Now())
	return int(duration.Hours() / 24)
}

// IsHighPriority определяет, является ли тендер высокоприоритетным
func (t *Tender) IsHighPriority() bool {
	return t.AIRelevanceScore >= 0.8 && 
		   t.StartPrice >= 1000000 && // >= 1 млн рублей
		   t.DaysUntilDeadline() >= 7 && // >= 7 дней до дедлайна
		   t.CompetitionLevel != CompetitionLevelHigh
}

// CanProcessDocuments проверяет, можно ли обрабатывать документы
func (t *Tender) CanProcessDocuments() bool {
	return t.AIRelevanceScore >= 0.7 && 
		   t.DocumentsURL != "" && 
		   !t.DocumentsProcessed
}

// ShouldStartEmailCampaign проверяет, стоит ли запускать email кампанию
func (t *Tender) ShouldStartEmailCampaign() bool {
	return t.ProductsExtracted && 
		   t.SuppliersFound && 
		   !t.EmailCampaignSent &&
		   t.DaysUntilDeadline() >= 3 // Минимум 3 дня до дедлайна
}

// CalculateWinProbability рассчитывает вероятность выигрыша на основе различных факторов
func (t *Tender) CalculateWinProbability() float64 {
	// TODO: Сложная логика расчета вероятности выигрыша
	// Учитывает:
	// - AI оценку релевантности
	// - Уровень конкуренции
	// - Соответствие цены
	// - Историческую статистику
	// - Качество технического предложения
	
	probability := t.AIRelevanceScore * 0.3 // 30% - релевантность
	
	// Учитываем уровень конкуренции
	switch t.CompetitionLevel {
	case CompetitionLevelLow:
		probability += 0.4 // +40% за низкую конкуренцию
	case CompetitionLevelMedium:
		probability += 0.2 // +20% за среднюю конкуренцию
	case CompetitionLevelHigh:
		probability += 0.0 // Высокая конкуренция не добавляет
	}
	
	// Учитываем наличие оптимальной цены
	if t.OptimalPriceCalculated {
		probability += 0.3 // +30% за оптимизированную цену
	}
	
	// Ограничиваем значение в диапазоне [0, 1]
	if probability > 1.0 {
		probability = 1.0
	}
	if probability < 0.0 {
		probability = 0.0
	}
	
	return probability
}

// GetRiskLevel возвращает общий уровень риска участия
func (t *Tender) GetRiskLevel() string {
	riskCount := len(t.RiskFactors)
	
	switch {
	case riskCount == 0:
		return "low"
	case riskCount <= 2:
		return "medium"
	default:
		return "high"
	}
}

// ShouldParticipate принимает решение об участии на основе всех факторов
func (t *Tender) ShouldParticipate() bool {
	// TODO: Комплексная логика принятия решения
	// Должна учитывать:
	// - AI рекомендацию
	// - Вероятность выигрыша
	// - Уровень риска
	// - Финансовую привлекательность
	// - Стратегические цели компании
	
	if t.AIRecommendation == AIRecommendationSkip {
		return false
	}
	
	if t.AIRecommendation == AIRecommendationParticipate && 
	   t.CalculateWinProbability() >= 0.6 &&
	   t.GetRiskLevel() != "high" {
		return true
	}
	
	return false
}

// UpdateStatus обновляет статус тендера с валидацией переходов
func (t *Tender) UpdateStatus(newStatus TenderStatus) error {
	// TODO: Валидация допустимых переходов между статусами
	// Например: new -> ai_analysis -> documents_pending -> processing и т.д.
	
	validTransitions := map[TenderStatus][]TenderStatus{
		TenderStatusNew: {TenderStatusAIAnalysis, TenderStatusCancelled},
		TenderStatusAIAnalysis: {TenderStatusDocumentsPending, TenderStatusArchived},
		TenderStatusDocumentsPending: {TenderStatusProcessing, TenderStatusArchived},
		// ... остальные переходы
	}
	
	if allowedStatuses, exists := validTransitions[t.Status]; exists {
		for _, allowedStatus := range allowedStatuses {
			if allowedStatus == newStatus {
				t.Status = newStatus
				t.UpdatedAt = time.Now()
				return nil
			}
		}
	}
	
	return fmt.Errorf("invalid status transition from %s to %s", t.Status, newStatus)
}

// TODO: Дополнительные методы для бизнес-логики
// - Валидация данных тендера
// - Расчет метрик эффективности
// - Генерация отчетов
// - Интеграция с внешними системами
// - Уведомления и алерты

// Validate проверяет корректность данных тендера
func (t *Tender) Validate() error {
	// TODO: Валидация всех полей тендера
	if t.TenderID == "" {
		return fmt.Errorf("tender ID is required")
	}
	
	if t.Title == "" {
		return fmt.Errorf("tender title is required")
	}
	
	if t.Platform == "" {
		return fmt.Errorf("platform is required")
	}
	
	if t.StartPrice < 0 {
		return fmt.Errorf("start price cannot be negative")
	}
	
	if t.ApplicationDeadline.Before(time.Now()) {
		return fmt.Errorf("application deadline cannot be in the past")
	}
	
	return nil
}

// 📊 ЕНУМЫ И КОНСТАНТЫ

type TenderCategory string

const (
	CategoryMedicalEquipment  TenderCategory = "medical_equipment"   // Медицинское оборудование
	CategoryMedicalSupplies   TenderCategory = "medical_supplies"    // Медицинские расходники
	CategoryLaboratoryEquip   TenderCategory = "laboratory_equipment" // Лабораторное оборудование
	CategoryDiagnosticEquip   TenderCategory = "diagnostic_equipment" // Диагностическое оборудование
	CategorySurgicalEquip     TenderCategory = "surgical_equipment"   // Хирургическое оборудование
	CategoryRehabilitationEquip TenderCategory = "rehabilitation_equipment" // Реабилитационное оборудование
	CategoryOther             TenderCategory = "other"               // Прочее
)

type TenderStatus string

const (
	StatusActive      TenderStatus = "active"       // Активный тендер
	StatusCompleted   TenderStatus = "completed"    // Завершенный тендер
	StatusCancelled   TenderStatus = "cancelled"    // Отмененный тендер
	StatusPostponed   TenderStatus = "postponed"    // Отложенный тендер
)

type ProcessingStatus string

const (
	ProcessingNew           ProcessingStatus = "new"              // Новый, не обработан
	ProcessingDiscovered    ProcessingStatus = "discovered"      // Найден скрапером
	ProcessingAnalyzing     ProcessingStatus = "analyzing"       // Анализируется ИИ
	ProcessingDocuments     ProcessingStatus = "downloading_docs" // Скачиваются документы
	ProcessingProducts      ProcessingStatus = "extracting_products" // Извлекаются товары
	ProcessingSuppliers     ProcessingStatus = "finding_suppliers"   // Ищутся поставщики
	ProcessingContacting    ProcessingStatus = "contacting_suppliers" // Связь с поставщиками
	ProcessingCompleted     ProcessingStatus = "completed"       // Обработка завершена
	ProcessingFailed        ProcessingStatus = "failed"          // Ошибка обработки
	ProcessingSkipped       ProcessingStatus = "skipped"         // Пропущен (не релевантен)
)

type AIRecommendation string

const (
	RecommendationParticipate AIRecommendation = "participate"    // Участвовать в тендере
	RecommendationSkip        AIRecommendation = "skip"          // Пропустить тендер
	RecommendationAnalyze     AIRecommendation = "need_analysis" // Нужен дополнительный анализ
	RecommendationWatchList   AIRecommendation = "watch_list"    // Добавить в список наблюдения
)

// 🏥 СПЕЦИАЛИЗИРОВАННЫЕ СТРУКТУРЫ

// Товар в тендере (связан с Product entity)
type TenderProduct struct {
	TenderID    uint    `json:"tender_id" db:"tender_id"`
	ProductID   uint    `json:"product_id" db:"product_id"`
	Quantity    int     `json:"quantity" db:"quantity"`
	Unit        string  `json:"unit" db:"unit"`         // шт, кг, м, л
	MaxPrice    float64 `json:"max_price" db:"max_price"` // Максимальная цена за единицу
	
	// Требования из техзадания
	Specifications string `json:"specifications" db:"specifications"` // Технические требования
	Requirements   string `json:"requirements" db:"requirements"`     // Дополнительные требования
}

// Email кампания для тендера
type EmailCampaign struct {
	TenderID     uint      `json:"tender_id" db:"tender_id"`
	SupplierID   uint      `json:"supplier_id" db:"supplier_id"`
	SentAt       time.Time `json:"sent_at" db:"sent_at"`
	ResponseAt   *time.Time `json:"response_at" db:"response_at"`
	Status       string    `json:"status" db:"status"`    // sent, delivered, opened, replied
	OfferReceived bool     `json:"offer_received" db:"offer_received"`
}

// Анализ цен для тендера
type PriceAnalysis struct {
	TenderID          uint      `json:"tender_id" db:"tender_id"`
	RecommendedPrice  float64   `json:"recommended_price" db:"recommended_price"`
	MinViablePrice    float64   `json:"min_viable_price" db:"min_viable_price"`
	MaxCompetitivePrice float64 `json:"max_competitive_price" db:"max_competitive_price"`
	WinProbability    float64   `json:"win_probability" db:"win_probability"`
	AnalysisDate      time.Time `json:"analysis_date" db:"analysis_date"`
	Confidence        float64   `json:"confidence" db:"confidence"` // Уверенность прогноза
}

// 🚨 ДОМЕННЫЕ ОШИБКИ
var (
	ErrInvalidTender        = errors.New("invalid tender data")
	ErrInvalidURL          = errors.New("invalid tender URL")
	ErrInvalidPrice        = errors.New("invalid price value")
	ErrInvalidDate         = errors.New("invalid date")
	ErrTenderNotFound      = errors.New("tender not found")
	ErrTenderExpired       = errors.New("tender submission period expired")
	ErrInvalidPlatform     = errors.New("unsupported tender platform")
	ErrEmptyTitle          = errors.New("tender title cannot be empty")
	ErrInvalidCategory     = errors.New("invalid tender category")
	ErrDocumentsNotReady   = errors.New("tender documents not downloaded yet")
)

// 🏭 КОНСТРУКТОР ТЕНДЕРА
func NewTender(externalID, number, platform, originalURL, title string) (*Tender, error) {
	// Валидация обязательных полей
	if strings.TrimSpace(title) == "" {
		return nil, ErrEmptyTitle
	}
	
	if strings.TrimSpace(externalID) == "" {
		return nil, fmt.Errorf("external ID is required")
	}
	
	// Валидация URL
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}
	
	// Валидация платформы
	if err := validatePlatform(platform); err != nil {
		return nil, err
	}
	
	now := time.Now()
	
	return &Tender{
		ExternalID:       externalID,
		Number:          number,
		Platform:        platform,
		OriginalURL:     originalURL,
		Title:           strings.TrimSpace(title),
		Currency:        "RUB",
		Status:          StatusActive,
		ProcessingStatus: ProcessingNew,
		AIRecommendation: RecommendationAnalyze,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// 🔧 МЕТОДЫ ДЛЯ ИЗМЕНЕНИЯ СОСТОЯНИЯ

// Установить результат анализа ИИ
func (t *Tender) SetAIAnalysis(score float64, analysis string, recommendation AIRecommendation) error {
	if score < 0 || score > 1 {
		return fmt.Errorf("AI relevance score must be between 0 and 1")
	}
	
	t.AIRelevanceScore = score
	t.AIAnalysis = analysis
	t.AIRecommendation = recommendation
	t.UpdatedAt = time.Now()
	
	// Автоматически обновляем статус обработки
	if recommendation == RecommendationSkip {
		t.ProcessingStatus = ProcessingSkipped
	}
	
	return nil
}

// Отметить, что документы скачаны
func (t *Tender) MarkDocumentsDownloaded(documentURLs []string, technicalTaskURL string) {
	t.DocumentsDownloaded = true
	t.DocumentURLs = documentURLs
	t.DocumentsCount = len(documentURLs)
	t.TechnicalTaskURL = technicalTaskURL
	t.UpdatedAt = time.Now()
	
	if t.ProcessingStatus == ProcessingDocuments {
		t.ProcessingStatus = ProcessingProducts
	}
}

// Отметить, что товары извлечены
func (t *Tender) MarkProductsExtracted() {
	t.ProductsExtracted = true
	t.UpdatedAt = time.Now()
	
	if t.ProcessingStatus == ProcessingProducts {
		t.ProcessingStatus = ProcessingSuppliers
	}
}

// Отметить, что поставщики найдены и связь установлена
func (t *Tender) MarkSuppliersContacted() {
	t.SuppliersContacted = true
	t.UpdatedAt = time.Now()
	
	if t.ProcessingStatus == ProcessingContacting {
		t.ProcessingStatus = ProcessingCompleted
	}
}

// Установить результаты тендера
func (t *Tender) SetResults(winnerCompany string, winnerPrice float64, totalParticipants int) error {
	if winnerPrice < 0 {
		return ErrInvalidPrice
	}
	
	t.WinnerCompany = winnerCompany
	t.WinnerPrice = winnerPrice
	t.TotalParticipants = totalParticipants
	
	// Рассчитать скидку
	if t.StartPrice > 0 {
		t.WinnerDiscount = ((t.StartPrice - winnerPrice) / t.StartPrice) * 100
	}
	
	t.Status = StatusCompleted
	t.UpdatedAt = time.Now()
	
	return nil
}

// Обновить статус обработки
func (t *Tender) UpdateProcessingStatus(status ProcessingStatus) {
	t.ProcessingStatus = status
	t.UpdatedAt = time.Now()
}

// 🔍 МЕТОДЫ ЗАПРОСОВ

// Проверить, активен ли тендер
func (t *Tender) IsActive() bool {
	return t.Status == StatusActive && time.Now().Before(t.SubmissionEnd)
}

// Проверить, завершился ли прием заявок
func (t *Tender) IsSubmissionClosed() bool {
	return time.Now().After(t.SubmissionEnd)
}

// Проверить, релевантен ли тендер для участия
func (t *Tender) IsRelevant() bool {
	return t.IsMedical && 
		   t.IsGoods && 
		   t.AIRelevanceScore >= 0.7 &&
		   t.AIRecommendation == RecommendationParticipate
}

// Проверить, готов ли тендер для участия
func (t *Tender) IsReadyForParticipation() bool {
	return t.DocumentsDownloaded &&
		   t.ProductsExtracted &&
		   t.SuppliersContacted &&
		   t.ProcessingStatus == ProcessingCompleted
}

// Получить время до окончания подачи заявок
func (t *Tender) TimeUntilSubmissionEnd() time.Duration {
	return time.Until(t.SubmissionEnd)
}

// Проверить, есть ли достаточно времени для участия
func (t *Tender) HasSufficientTime(minDays int) bool {
	return t.TimeUntilSubmissionEnd() > time.Duration(minDays)*24*time.Hour
}

// 🔧 ПРИВАТНЫЕ ФУНКЦИИ ВАЛИДАЦИИ

func validateURL(rawURL string) error {
	if rawURL == "" {
		return ErrInvalidURL
	}
	
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}
	
	if u.Scheme == "" || u.Host == "" {
		return ErrInvalidURL
	}
	
	return nil
}

func validatePlatform(platform string) error {
	validPlatforms := map[string]bool{
		"zakupki.gov.ru":    true,
		"szvo.gov35.ru":     true,
		"gz-spb.ru":         true,
		"gz-nn.ru":          true,
		"gz.mos.ru":         true,
		"aggregator-kkr.ru": true,
		"roseltorg.ru":      true,
		"b2b-center.ru":     true,
		"sberbank-ast.ru":   true,
		"fabrikant.ru":      true,
	}
	
	if !validPlatforms[platform] {
		return ErrInvalidPlatform
	}
	
	return nil
}

// 📊 ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ

// Получить цвет статуса для UI
func (s TenderStatus) Color() string {
	switch s {
	case StatusActive:
		return "green"
	case StatusCompleted:
		return "blue"
	case StatusCancelled:
		return "red"
	case StatusPostponed:
		return "yellow"
	default:
		return "gray"
	}
}

// Получить человекочитаемое название категории
func (c TenderCategory) DisplayName() string {
	switch c {
	case CategoryMedicalEquipment:
		return "Медицинское оборудование"
	case CategoryMedicalSupplies:
		return "Медицинские расходники"
	case CategoryLaboratoryEquip:
		return "Лабораторное оборудование"
	case CategoryDiagnosticEquip:
		return "Диагностическое оборудование"
	case CategorySurgicalEquip:
		return "Хирургическое оборудование"
	case CategoryRehabilitationEquip:
		return "Реабилитационное оборудование"
	default:
		return "Прочее"
	}
}

// Получить прогресс обработки в процентах
func (p ProcessingStatus) Progress() int {
	switch p {
	case ProcessingNew:
		return 0
	case ProcessingDiscovered:
		return 10
	case ProcessingAnalyzing:
		return 20
	case ProcessingDocuments:
		return 40
	case ProcessingProducts:
		return 60
	case ProcessingSuppliers:
		return 80
	case ProcessingContacting:
		return 90
	case ProcessingCompleted:
		return 100
	case ProcessingFailed, ProcessingSkipped:
		return 0
	default:
		return 0
	}
}