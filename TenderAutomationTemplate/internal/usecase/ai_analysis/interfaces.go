// 🤖 ИНТЕРФЕЙСЫ ДЛЯ ИИ АНАЛИЗА ТЕНДЕРОВ
//
// Определяет интерфейсы для интеграции с Llama 4 и другими LLM
// для анализа тендеров, извлечения товаров и оптимизации цен

package ai_analysis

import (
	"context"
	"time"
	
	"tender-automation-system/internal/domain/tender"
	"tender-automation-system/internal/domain/product"
)

// 🧠 ОСНОВНОЙ ИНТЕРФЕЙС LLM ПРОВАЙДЕРА
type LLMProvider interface {
	// Анализ релевантности тендера
	AnalyzeTenderRelevance(ctx context.Context, req *TenderAnalysisRequest) (*TenderAnalysisResponse, error)
	
	// Извлечение товаров из технического задания
	ExtractProductsFromDocument(ctx context.Context, req *DocumentAnalysisRequest) (*ProductExtractionResponse, error)
	
	// Анализ цен и рекомендации
	AnalyzePricing(ctx context.Context, req *PricingAnalysisRequest) (*PricingAnalysisResponse, error)
	
	// Генерация email для поставщиков
	GenerateSupplierEmail(ctx context.Context, req *EmailGenerationRequest) (*EmailGenerationResponse, error)
	
	// Анализ ответов поставщиков
	AnalyzeSupplierResponse(ctx context.Context, req *SupplierResponseRequest) (*SupplierResponseAnalysis, error)
}

// 📄 ИНТЕРФЕЙС ДЛЯ ОБРАБОТКИ ДОКУМЕНТОВ
type DocumentProcessor interface {
	// Извлечение текста из DOC/DOCX
	ExtractTextFromDoc(ctx context.Context, filePath string) (string, error)
	
	// Извлечение текста из PDF
	ExtractTextFromPDF(ctx context.Context, filePath string) (string, error)
	
	// Извлечение таблиц из документа
	ExtractTablesFromDoc(ctx context.Context, filePath string) ([]DocumentTable, error)
	
	// Очистка и нормализация текста
	CleanText(text string) string
	
	// Определение типа документа
	DetectDocumentType(ctx context.Context, filePath string) (DocumentType, error)
}

// 📊 ИНТЕРФЕЙС ДЛЯ АНАЛИЗА ЦЕН
type PriceAnalyzer interface {
	// Анализ исторических данных тендеров
	AnalyzeHistoricalPrices(ctx context.Context, productIDs []uint, region string) (*PriceHistoryAnalysis, error)
	
	// Прогнозирование выигрышной цены
	PredictWinningPrice(ctx context.Context, req *WinningPricePredictionRequest) (*WinningPricePrediction, error)
	
	// Оптимизация цены для максимизации прибыли
	OptimizePrice(ctx context.Context, req *PriceOptimizationRequest) (*PriceOptimizationResult, error)
	
	// Анализ конкурентов
	AnalyzeCompetitors(ctx context.Context, tenderID uint) (*CompetitorAnalysis, error)
}

// 🔍 ИНТЕРФЕЙС ДЛЯ ПОИСКА ТОВАРОВ
type ProductMatcher interface {
	// Поиск похожих товаров в базе
	FindSimilarProducts(ctx context.Context, productName string, specifications map[string]interface{}) ([]product.Product, error)
	
	// Сопоставление товара с каталогом поставщиков
	MatchWithSupplierCatalog(ctx context.Context, productName string, supplierID uint) (*SupplierProductMatch, error)
	
	// Нечеткий поиск товаров
	FuzzySearchProducts(ctx context.Context, query string, threshold float64) ([]ProductMatch, error)
	
	// Извлечение ключевых характеристик товара
	ExtractProductFeatures(ctx context.Context, productDescription string) (*ProductFeatures, error)
}

// 📧 ИНТЕРФЕЙС ДЛЯ ГЕНЕРАЦИИ КОНТЕНТА
type ContentGenerator interface {
	// Генерация персонализированного email
	GeneratePersonalizedEmail(ctx context.Context, req *PersonalizedEmailRequest) (string, error)
	
	// Генерация коммерческого предложения
	GenerateCommercialOffer(ctx context.Context, req *CommercialOfferRequest) (string, error)
	
	// Генерация технического описания
	GenerateTechnicalDescription(ctx context.Context, product *product.Product) (string, error)
	
	// Перевод и локализация
	TranslateContent(ctx context.Context, text, targetLanguage string) (string, error)
}

// 📊 СТРУКТУРЫ ЗАПРОСОВ И ОТВЕТОВ

// 🔍 АНАЛИЗ ТЕНДЕРА
type TenderAnalysisRequest struct {
	Tender      *tender.Tender `json:"tender"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	CustomerInfo string        `json:"customer_info"`
	Budget      float64        `json:"budget"`
}

type TenderAnalysisResponse struct {
	IsRelevant       bool                      `json:"is_relevant"`
	RelevanceScore   float64                   `json:"relevance_score"`    // 0-1
	Recommendation   tender.AIRecommendation   `json:"recommendation"`
	Reasoning        string                    `json:"reasoning"`
	Category         tender.TenderCategory     `json:"category"`
	IsMedical        bool                      `json:"is_medical"`
	IsEquipment      bool                      `json:"is_equipment"`
	IsGoods          bool                      `json:"is_goods"`
	RiskFactors      []string                  `json:"risk_factors"`
	Opportunities    []string                  `json:"opportunities"`
	CompetitionLevel CompetitionLevel          `json:"competition_level"`
	Confidence       float64                   `json:"confidence"`
}

// 📄 АНАЛИЗ ДОКУМЕНТОВ
type DocumentAnalysisRequest struct {
	DocumentPath string            `json:"document_path"`
	DocumentText string            `json:"document_text"`
	DocumentType DocumentType      `json:"document_type"`
	TenderID     uint             `json:"tender_id"`
	Context      map[string]interface{} `json:"context"`
}

type ProductExtractionResponse struct {
	Products    []ExtractedProduct `json:"products"`
	Tables      []DocumentTable    `json:"tables"`
	TotalItems  int               `json:"total_items"`
	Confidence  float64           `json:"confidence"`
	Summary     string            `json:"summary"`
	Issues      []string          `json:"issues"`
}

type ExtractedProduct struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Quantity       int                    `json:"quantity"`
	Unit           string                 `json:"unit"`
	MaxPrice       float64                `json:"max_price"`
	Specifications map[string]interface{} `json:"specifications"`
	Requirements   string                 `json:"requirements"`
	Category       product.ProductCategory `json:"category"`
	GOST          string                 `json:"gost"`
	Manufacturer   string                 `json:"manufacturer"`
	OriginalText   string                 `json:"original_text"`
	Confidence     float64                `json:"confidence"`
	TableIndex     int                    `json:"table_index"`
	RowIndex       int                    `json:"row_index"`
}

// 💰 АНАЛИЗ ЦЕН
type PricingAnalysisRequest struct {
	TenderID         uint                   `json:"tender_id"`
	Products         []ExtractedProduct     `json:"products"`
	HistoricalData   []HistoricalTender     `json:"historical_data"`
	MarketConditions *MarketConditions      `json:"market_conditions"`
	CompanyProfile   *CompanyProfile        `json:"company_profile"`
}

type PricingAnalysisResponse struct {
	RecommendedDiscount float64                `json:"recommended_discount"`  // Рекомендуемая скидка (%)
	TotalPrice         float64                `json:"total_price"`           // Итоговая цена
	WinProbability     float64                `json:"win_probability"`       // Вероятность выигрыша
	ProfitMargin       float64                `json:"profit_margin"`         // Маржа прибыли
	RiskLevel          RiskLevel              `json:"risk_level"`            // Уровень риска
	ProductPrices      []ProductPriceAdvice   `json:"product_prices"`        // Цены по товарам
	Strategy           PricingStrategy        `json:"strategy"`              // Стратегия ценообразования
	Justification      string                 `json:"justification"`         // Обоснование цены
	Alternatives       []PricingAlternative   `json:"alternatives"`          // Альтернативные варианты
}

// 📧 ГЕНЕРАЦИЯ EMAIL
type EmailGenerationRequest struct {
	SupplierName    string               `json:"supplier_name"`
	SupplierType    string               `json:"supplier_type"`
	Products        []ExtractedProduct   `json:"products"`
	TenderInfo      *tender.Tender       `json:"tender_info"`
	CompanyInfo     *CompanyProfile      `json:"company_info"`
	EmailType       EmailType            `json:"email_type"`
	PersonalizedData map[string]interface{} `json:"personalized_data"`
}

type EmailGenerationResponse struct {
	Subject      string   `json:"subject"`
	Body         string   `json:"body"`
	Attachments  []string `json:"attachments"`
	Priority     EmailPriority `json:"priority"`
	FollowUpDays int      `json:"follow_up_days"`
	ExpectedResponse string `json:"expected_response"`
}

// 📊 ДОПОЛНИТЕЛЬНЫЕ ТИПЫ

type DocumentType string

const (
	DocumentTypeTechnicalTask   DocumentType = "technical_task"
	DocumentTypeContract        DocumentType = "contract"
	DocumentTypeSpecifications  DocumentType = "specifications"
	DocumentTypePriceList       DocumentType = "price_list"
	DocumentTypeOther          DocumentType = "other"
)

type CompetitionLevel string

const (
	CompetitionLow    CompetitionLevel = "low"
	CompetitionMedium CompetitionLevel = "medium"
	CompetitionHigh   CompetitionLevel = "high"
)

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type EmailType string

const (
	EmailTypeInitialInquiry    EmailType = "initial_inquiry"
	EmailTypeFollowUp         EmailType = "follow_up"
	EmailTypeUrgentRequest    EmailType = "urgent_request"
	EmailTypePriceNegotiation EmailType = "price_negotiation"
)

type EmailPriority string

const (
	PriorityLow    EmailPriority = "low"
	PriorityNormal EmailPriority = "normal"
	PriorityHigh   EmailPriority = "high"
	PriorityUrgent EmailPriority = "urgent"
)

type PricingStrategy string

const (
	StrategyAggressive     PricingStrategy = "aggressive"      // Агрессивная - максимальная скидка
	StrategyCompetitive    PricingStrategy = "competitive"     // Конкурентная - рыночная цена
	StrategyPremium        PricingStrategy = "premium"         // Премиум - высокая цена, качество
	StrategyValue          PricingStrategy = "value"           // Оптимальное соотношение цена/качество
)

// 📊 ВСПОМОГАТЕЛЬНЫЕ СТРУКТУРЫ

type DocumentTable struct {
	Headers []string              `json:"headers"`
	Rows    [][]string           `json:"rows"`
	Title   string               `json:"title"`
	PageNum int                  `json:"page_num"`
	Metadata map[string]interface{} `json:"metadata"`
}

type HistoricalTender struct {
	ID               uint      `json:"id"`
	Title            string    `json:"title"`
	StartPrice       float64   `json:"start_price"`
	WinningPrice     float64   `json:"winning_price"`
	Discount         float64   `json:"discount"`
	ParticipantCount int       `json:"participant_count"`
	CompletedAt      time.Time `json:"completed_at"`
	Region           string    `json:"region"`
	Category         string    `json:"category"`
}

type MarketConditions struct {
	EconomicIndex    float64           `json:"economic_index"`
	SeasonalFactor   float64           `json:"seasonal_factor"`
	CompetitorCount  int               `json:"competitor_count"`
	DemandLevel      string            `json:"demand_level"`
	SupplyLevel      string            `json:"supply_level"`
	RegionalFactors  map[string]float64 `json:"regional_factors"`
}

type CompanyProfile struct {
	Name             string    `json:"name"`
	Experience       int       `json:"experience"`        // Лет опыта
	WinRate          float64   `json:"win_rate"`          // Процент выигранных тендеров
	AverageDiscount  float64   `json:"average_discount"`  // Средняя скидка
	Specialization   []string  `json:"specialization"`    // Специализация
	Region           []string  `json:"region"`            // Регионы работы
	CertificateLevel string    `json:"certificate_level"` // Уровень сертификации
}

type ProductPriceAdvice struct {
	ProductName      string  `json:"product_name"`
	OriginalPrice    float64 `json:"original_price"`
	RecommendedPrice float64 `json:"recommended_price"`
	Discount         float64 `json:"discount"`
	Confidence       float64 `json:"confidence"`
	Reasoning        string  `json:"reasoning"`
}

type PricingAlternative struct {
	Strategy        PricingStrategy `json:"strategy"`
	TotalPrice      float64         `json:"total_price"`
	WinProbability  float64         `json:"win_probability"`
	ProfitMargin    float64         `json:"profit_margin"`
	Description     string          `json:"description"`
	Pros            []string        `json:"pros"`
	Cons            []string        `json:"cons"`
}

type WinningPricePredictionRequest struct {
	TenderID          uint                   `json:"tender_id"`
	StartPrice        float64                `json:"start_price"`
	ProductCategories []product.ProductCategory `json:"product_categories"`
	Region           string                 `json:"region"`
	CompetitorCount  int                    `json:"competitor_count"`
	MarketConditions *MarketConditions      `json:"market_conditions"`
}

type WinningPricePrediction struct {
	PredictedPrice    float64   `json:"predicted_price"`
	PredictedDiscount float64   `json:"predicted_discount"`
	Confidence        float64   `json:"confidence"`
	PriceRange        PriceRange `json:"price_range"`
	Factors           []PriceFactor `json:"factors"`
}

type PriceRange struct {
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Median float64 `json:"median"`
}

type PriceFactor struct {
	Name   string  `json:"name"`
	Impact float64 `json:"impact"`  // Влияние на цену (-1 до 1)
	Weight float64 `json:"weight"`  // Вес фактора (0 до 1)
}

type PriceOptimizationRequest struct {
	TenderID         uint    `json:"tender_id"`
	TargetMargin     float64 `json:"target_margin"`     // Целевая маржа
	MaxRisk          float64 `json:"max_risk"`          // Максимальный риск
	MinWinProbability float64 `json:"min_win_probability"` // Минимальная вероятность выигрыша
}

type PriceOptimizationResult struct {
	OptimalPrice      float64              `json:"optimal_price"`
	OptimalDiscount   float64              `json:"optimal_discount"`
	ExpectedMargin    float64              `json:"expected_margin"`
	WinProbability    float64              `json:"win_probability"`
	RiskScore         float64              `json:"risk_score"`
	OptimizationSteps []OptimizationStep   `json:"optimization_steps"`
}

type OptimizationStep struct {
	Step        string  `json:"step"`
	Price       float64 `json:"price"`
	Probability float64 `json:"probability"`
	Margin      float64 `json:"margin"`
	Reasoning   string  `json:"reasoning"`
}

type CompetitorAnalysis struct {
	TotalCompetitors     int                    `json:"total_competitors"`
	KnownCompetitors     []CompetitorProfile    `json:"known_competitors"`
	AverageWinRate       float64                `json:"average_win_rate"`
	AverageDiscount      float64                `json:"average_discount"`
	CompetitionIntensity CompetitionLevel       `json:"competition_intensity"`
	ThreatLevel          string                 `json:"threat_level"`
	Recommendations      []string               `json:"recommendations"`
}

type CompetitorProfile struct {
	Name         string  `json:"name"`
	WinRate      float64 `json:"win_rate"`
	AverageDiscount float64 `json:"average_discount"`
	Specialization []string `json:"specialization"`
	LastActivity time.Time `json:"last_activity"`
	ThreatLevel  string   `json:"threat_level"`
}

type SupplierProductMatch struct {
	ProductID     uint    `json:"product_id"`
	SupplierSKU   string  `json:"supplier_sku"`
	MatchScore    float64 `json:"match_score"`
	Price         float64 `json:"price"`
	Available     bool    `json:"available"`
	DeliveryDays  int     `json:"delivery_days"`
	ProductURL    string  `json:"product_url"`
	Confidence    float64 `json:"confidence"`
}

type ProductMatch struct {
	Product    *product.Product `json:"product"`
	MatchScore float64          `json:"match_score"`
	MatchType  string           `json:"match_type"`  // exact, fuzzy, semantic
	Reasoning  string           `json:"reasoning"`
}

type ProductFeatures struct {
	MainFeatures    []Feature `json:"main_features"`
	TechnicalSpecs  []Feature `json:"technical_specs"`
	Keywords        []string  `json:"keywords"`
	Category        product.ProductCategory `json:"category"`
	Confidence      float64   `json:"confidence"`
}

type Feature struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Unit  string      `json:"unit"`
	Type  string      `json:"type"`  // numeric, text, boolean, enum
}

type PriceHistoryAnalysis struct {
	AveragePrice      float64         `json:"average_price"`
	PriceRange        PriceRange      `json:"price_range"`
	PriceTrend        string          `json:"price_trend"`    // increasing, decreasing, stable
	SeasonalFactors   []SeasonalFactor `json:"seasonal_factors"`
	RegionalPrices    map[string]float64 `json:"regional_prices"`
	Confidence        float64         `json:"confidence"`
}

type SeasonalFactor struct {
	Month      int     `json:"month"`
	Factor     float64 `json:"factor"`  // Коэффициент относительно среднего
	Confidence float64 `json:"confidence"`
}

type SupplierResponseRequest struct {
	EmailContent    string            `json:"email_content"`
	SupplierID      uint             `json:"supplier_id"`
	OriginalRequest *EmailGenerationRequest `json:"original_request"`
	AttachmentPaths []string         `json:"attachment_paths"`
}

type SupplierResponseAnalysis struct {
	HasOffer           bool              `json:"has_offer"`
	OfferDetails       *OfferDetails     `json:"offer_details"`
	ResponseType       string            `json:"response_type"`    // positive, negative, request_info
	Sentiment          string            `json:"sentiment"`        // positive, neutral, negative
	NextAction         string            `json:"next_action"`      // call, email, wait, negotiate
	ExtractedData      map[string]interface{} `json:"extracted_data"`
	Confidence         float64           `json:"confidence"`
	NeedsHumanReview   bool              `json:"needs_human_review"`
}

type OfferDetails struct {
	Products      []OfferProduct `json:"products"`
	TotalPrice    float64        `json:"total_price"`
	Currency      string         `json:"currency"`
	DeliveryTerms string         `json:"delivery_terms"`
	PaymentTerms  string         `json:"payment_terms"`
	ValidUntil    *time.Time     `json:"valid_until"`
	Notes         string         `json:"notes"`
}

type OfferProduct struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TotalPrice   float64 `json:"total_price"`
	DeliveryDays int     `json:"delivery_days"`
	Availability string  `json:"availability"`
	SKU          string  `json:"sku"`
}

type PersonalizedEmailRequest struct {
	SupplierName     string                 `json:"supplier_name"`
	ContactPerson    string                 `json:"contact_person"`
	SupplierProfile  *SupplierProfile       `json:"supplier_profile"`
	Products         []ExtractedProduct     `json:"products"`
	TenderDeadline   time.Time              `json:"tender_deadline"`
	EmailTemplate    string                 `json:"email_template"`
	PersonalizationData map[string]interface{} `json:"personalization_data"`
}

type SupplierProfile struct {
	CompanySize     string   `json:"company_size"`
	Specialization  []string `json:"specialization"`
	Region          string   `json:"region"`
	PreferredContact string  `json:"preferred_contact"`
	BusinessHours   string   `json:"business_hours"`
	Language        string   `json:"language"`
	PreviousDeals   int      `json:"previous_deals"`
}

type CommercialOfferRequest struct {
	TenderID        uint                   `json:"tender_id"`
	Products        []OfferProduct         `json:"products"`
	CompanyInfo     *CompanyProfile        `json:"company_info"`
	DeliveryTerms   string                 `json:"delivery_terms"`
	PaymentTerms    string                 `json:"payment_terms"`
	ValidityPeriod  int                    `json:"validity_period"`
	SpecialOffers   []string               `json:"special_offers"`
	Template        string                 `json:"template"`
}

// 📝 КОММЕНТАРИИ ПО ИНТЕГРАЦИИ С LLAMA 4 MAVIRIC:
//
// Llama 4 Maviric - отличный выбор для вашей системы:
//
// ✅ ПРЕИМУЩЕСТВА:
// - Высокое качество анализа текста на русском языке
// - Хорошо понимает контекст бизнес-документов
// - Способен к сложному reasoning и анализу
// - Может работать локально (конфиденциальность)
// - Поддерживает большие контексты для документов
//
// 🔧 РЕКОМЕНДАЦИИ ПО ИНТЕГРАЦИИ:
// 1. Используйте Ollama для локального развертывания
// 2. Настройте специализированные промпты для каждой задачи
// 3. Реализуйте кеширование результатов анализа
// 4. Добавьте fallback на другие модели при недоступности
// 5. Используйте RAG (Retrieval-Augmented Generation) для работы с базой знаний
//
// 📊 АЛЬТЕРНАТИВЫ ДЛЯ СРАВНЕНИЯ:
// - OpenAI GPT-4 (облачное решение)
// - Anthropic Claude (хорош для анализа документов)
// - Yandex GPT (локальная русская модель)
// - GigaChat (российское решение)