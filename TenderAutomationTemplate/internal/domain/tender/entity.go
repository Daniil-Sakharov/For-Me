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

// 📋 ОСНОВНАЯ МОДЕЛЬ ТЕНДЕРА
type Tender struct {
	// ✅ Базовые поля
	ID        uint      `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
	// 🔍 Идентификация тендера
	ExternalID     string `json:"external_id" db:"external_id"`         // ID в системе закупок
	Number         string `json:"number" db:"number"`                   // Номер тендера
	Platform       string `json:"platform" db:"platform"`              // zakupki.gov.ru, szvo.gov35.ru
	OriginalURL    string `json:"original_url" db:"original_url"`       // Ссылка на тендер
	
	// 📄 Основная информация
	Title          string    `json:"title" db:"title"`                 // Название тендера
	Description    string    `json:"description" db:"description"`     // Описание
	CustomerName   string    `json:"customer_name" db:"customer_name"` // Заказчик
	CustomerINN    string    `json:"customer_inn" db:"customer_inn"`   // ИНН заказчика
	Region         string    `json:"region" db:"region"`               // Регион
	
	// 💰 Финансовая информация
	StartPrice     float64 `json:"start_price" db:"start_price"`         // Начальная цена
	Currency       string  `json:"currency" db:"currency"`               // Валюта (RUB)
	
	// 📅 Временные рамки
	PublishedAt    time.Time  `json:"published_at" db:"published_at"`       // Дата публикации
	SubmissionEnd  time.Time  `json:"submission_end" db:"submission_end"`   // Окончание подачи заявок
	TenderEnd      time.Time  `json:"tender_end" db:"tender_end"`           // Окончание торгов
	DeliveryPeriod int        `json:"delivery_period" db:"delivery_period"` // Срок поставки (дни)
	
	// 🏥 Специфика медицинского оборудования
	Category        TenderCategory `json:"category" db:"category"`               // Категория товаров
	IsMedical       bool          `json:"is_medical" db:"is_medical"`           // Медицинское оборудование
	IsEquipment     bool          `json:"is_equipment" db:"is_equipment"`       // Оборудование vs расходники
	IsGoods         bool          `json:"is_goods" db:"is_goods"`               // Поставка товаров vs услуг
	
	// 🤖 ИИ анализ
	AIRelevanceScore  float64        `json:"ai_relevance_score" db:"ai_relevance_score"`   // Оценка релевантности ИИ (0-1)
	AIAnalysis        string         `json:"ai_analysis" db:"ai_analysis"`                 // Анализ ИИ в JSON
	AIRecommendation  AIRecommendation `json:"ai_recommendation" db:"ai_recommendation"`   // Рекомендация ИИ
	
	// 📊 Статус обработки
	Status            TenderStatus `json:"status" db:"status"`                     // Статус тендера
	ProcessingStatus  ProcessingStatus `json:"processing_status" db:"processing_status"` // Статус обработки
	DocumentsDownloaded bool       `json:"documents_downloaded" db:"documents_downloaded"` // Документы скачаны
	ProductsExtracted   bool       `json:"products_extracted" db:"products_extracted"`     // Товары извлечены
	SuppliersContacted  bool       `json:"suppliers_contacted" db:"suppliers_contacted"`   // Поставщики найдены
	
	// 📄 Связанные документы
	DocumentsCount    int      `json:"documents_count" db:"documents_count"`       // Количество документов
	TechnicalTaskURL  string   `json:"technical_task_url" db:"technical_task_url"` // Ссылка на техзадание
	DocumentURLs      []string `json:"document_urls" db:"-"`                       // Все документы (JSON в БД)
	
	// 📈 Результаты тендера (заполняется после завершения)
	WinnerCompany     string  `json:"winner_company" db:"winner_company"`         // Компания-победитель
	WinnerPrice       float64 `json:"winner_price" db:"winner_price"`             // Цена победителя
	WinnerDiscount    float64 `json:"winner_discount" db:"winner_discount"`       // Скидка (%)
	TotalParticipants int     `json:"total_participants" db:"total_participants"` // Участников
	
	// 🔗 Связи с другими сущностями
	Products         []TenderProduct `json:"products" db:"-"`           // Товары в тендере
	EmailCampaigns   []EmailCampaign `json:"email_campaigns" db:"-"`    // Email кампании
	PriceAnalysis    *PriceAnalysis  `json:"price_analysis" db:"-"`     // Анализ цен
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