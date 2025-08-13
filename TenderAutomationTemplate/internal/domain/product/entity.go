// 🏥 ДОМЕННАЯ МОДЕЛЬ ТОВАРА МЕДИЦИНСКОГО ОБОРУДОВАНИЯ
//
// Представляет товар/оборудование из техзаданий тендеров
// Связывается с поставщиками и их каталогами

package product

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// 🏥 ОСНОВНАЯ МОДЕЛЬ ТОВАРА
type Product struct {
	// ✅ Базовые поля
	ID        uint      `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
	// 📄 Основная информация о товаре
	Name            string `json:"name" db:"name"`                       // Название товара
	Description     string `json:"description" db:"description"`         // Описание
	Category        ProductCategory `json:"category" db:"category"`       // Категория
	Subcategory     string `json:"subcategory" db:"subcategory"`         // Подкатегория
	
	// 🔍 Идентификация и поиск
	Keywords        []string `json:"keywords" db:"-"`                    // Ключевые слова для поиска (JSON в БД)
	Synonyms        []string `json:"synonyms" db:"-"`                    // Синонимы названия
	GOST           string   `json:"gost" db:"gost"`                      // ГОСТ или стандарт
	ArticleNumber   string   `json:"article_number" db:"article_number"` // Артикул производителя
	Barcode         string   `json:"barcode" db:"barcode"`               // Штрих-код
	
	// 🏭 Производство и сертификация
	Manufacturer    string   `json:"manufacturer" db:"manufacturer"`      // Производитель
	CountryOfOrigin string   `json:"country_of_origin" db:"country_of_origin"` // Страна производства
	CertificateType string   `json:"certificate_type" db:"certificate_type"`   // Тип сертификата (МЗ РФ, ТР ТС)
	CertificateNumber string `json:"certificate_number" db:"certificate_number"` // Номер сертификата
	CertificateValidUntil *time.Time `json:"certificate_valid_until" db:"certificate_valid_until"` // Срок действия
	
	// 📐 Технические характеристики
	Specifications  map[string]interface{} `json:"specifications" db:"-"` // Техспецификации (JSON в БД)
	Dimensions      *ProductDimensions     `json:"dimensions" db:"-"`     // Размеры
	Weight          *float64              `json:"weight" db:"weight"`     // Вес (кг)
	PowerConsumption *float64             `json:"power_consumption" db:"power_consumption"` // Потребляемая мощность (Вт)
	
	// 💰 Ценовая информация
	AveragePrice    *float64 `json:"average_price" db:"average_price"`     // Средняя рыночная цена
	MinPrice        *float64 `json:"min_price" db:"min_price"`             // Минимальная найденная цена
	MaxPrice        *float64 `json:"max_price" db:"max_price"`             // Максимальная найденная цена
	LastPriceUpdate *time.Time `json:"last_price_update" db:"last_price_update"` // Последнее обновление цен
	
	// 🔗 Связи с поставщиками
	SupplierProducts []SupplierProduct `json:"supplier_products" db:"-"` // Товар у разных поставщиков
	
	// 🤖 ИИ анализ
	AIExtracted     bool     `json:"ai_extracted" db:"ai_extracted"`       // Извлечен ИИ из документа
	AIConfidence    float64  `json:"ai_confidence" db:"ai_confidence"`     // Уверенность ИИ в корректности
	OriginalText    string   `json:"original_text" db:"original_text"`     // Оригинальный текст из документа
	
	// 📊 Статистика использования
	TimesFound      int       `json:"times_found" db:"times_found"`         // Сколько раз встречался в тендерах
	LastFoundAt     *time.Time `json:"last_found_at" db:"last_found_at"`    // Последний раз найден
	PopularityScore float64   `json:"popularity_score" db:"popularity_score"` // Индекс популярности
}

// 📊 ЕНУМЫ И КОНСТАНТЫ

type ProductCategory string

const (
	CategoryDiagnostic     ProductCategory = "diagnostic"     // Диагностическое оборудование
	CategorySurgical       ProductCategory = "surgical"       // Хирургическое оборудование
	CategoryLaboratory     ProductCategory = "laboratory"     // Лабораторное оборудование
	CategoryReanimation    ProductCategory = "reanimation"    // Реанимационное оборудование
	CategoryRadiology      ProductCategory = "radiology"      // Рентгенологическое оборудование
	CategoryPhysical       ProductCategory = "physical"       // Физиотерапевтическое оборудование
	CategoryDental         ProductCategory = "dental"         // Стоматологическое оборудование
	CategoryConsumables    ProductCategory = "consumables"    // Расходные материалы
	CategoryFurniture      ProductCategory = "furniture"      // Медицинская мебель
	CategoryOther          ProductCategory = "other"          // Прочее
)

// 📐 РАЗМЕРЫ ТОВАРА
type ProductDimensions struct {
	Length float64 `json:"length" db:"length"` // Длина (см)
	Width  float64 `json:"width" db:"width"`   // Ширина (см)
	Height float64 `json:"height" db:"height"` // Высота (см)
}

// 🏪 ТОВАР У ПОСТАВЩИКА
type SupplierProduct struct {
	ProductID      uint      `json:"product_id" db:"product_id"`
	SupplierID     uint      `json:"supplier_id" db:"supplier_id"`
	SupplierSKU    string    `json:"supplier_sku" db:"supplier_sku"`       // SKU у поставщика
	ProductURL     string    `json:"product_url" db:"product_url"`         // Ссылка на товар у поставщика
	Price          float64   `json:"price" db:"price"`                     // Цена у поставщика
	Currency       string    `json:"currency" db:"currency"`               // Валюта
	InStock        bool      `json:"in_stock" db:"in_stock"`               // В наличии
	DeliveryDays   int       `json:"delivery_days" db:"delivery_days"`     // Срок поставки (дни)
	MinOrderQty    int       `json:"min_order_qty" db:"min_order_qty"`     // Минимальная партия
	LastChecked    time.Time `json:"last_checked" db:"last_checked"`       // Последняя проверка
	FoundByParser  bool      `json:"found_by_parser" db:"found_by_parser"` // Найден парсером
}

// 🚨 ДОМЕННЫЕ ОШИБКИ
var (
	ErrInvalidProduct       = errors.New("invalid product data")
	ErrEmptyProductName     = errors.New("product name cannot be empty")
	ErrInvalidCategory      = errors.New("invalid product category")
	ErrInvalidPrice         = errors.New("invalid price value")
	ErrProductNotFound      = errors.New("product not found")
	ErrDuplicateProduct     = errors.New("product already exists")
	ErrInvalidDimensions    = errors.New("invalid product dimensions")
	ErrInvalidWeight        = errors.New("invalid weight value")
	ErrMissingSpecifications = errors.New("product specifications are required")
)

// 🏭 КОНСТРУКТОР ТОВАРА
func NewProduct(name, description string, category ProductCategory) (*Product, error) {
	// Валидация обязательных полей
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyProductName
	}
	
	if err := validateCategory(category); err != nil {
		return nil, err
	}
	
	now := time.Now()
	
	return &Product{
		Name:            strings.TrimSpace(name),
		Description:     strings.TrimSpace(description),
		Category:        category,
		Keywords:        extractKeywords(name, description),
		Specifications:  make(map[string]interface{}),
		CreatedAt:       now,
		UpdatedAt:       now,
		TimesFound:      1,
		LastFoundAt:     &now,
		PopularityScore: 1.0,
	}, nil
}

// 🤖 СОЗДАНИЕ ТОВАРА ИЗ АНАЛИЗА ИИ
func NewProductFromAI(name, description, originalText string, category ProductCategory, confidence float64) (*Product, error) {
	product, err := NewProduct(name, description, category)
	if err != nil {
		return nil, err
	}
	
	product.AIExtracted = true
	product.AIConfidence = confidence
	product.OriginalText = originalText
	product.Synonyms = generateSynonyms(name)
	
	return product, nil
}

// 🔧 МЕТОДЫ ДЛЯ ИЗМЕНЕНИЯ СОСТОЯНИЯ

// Обновить технические характеристики
func (p *Product) UpdateSpecifications(specs map[string]interface{}) error {
	if len(specs) == 0 {
		return ErrMissingSpecifications
	}
	
	if p.Specifications == nil {
		p.Specifications = make(map[string]interface{})
	}
	
	for key, value := range specs {
		p.Specifications[key] = value
	}
	
	p.UpdatedAt = time.Now()
	return nil
}

// Обновить размеры товара
func (p *Product) UpdateDimensions(length, width, height float64) error {
	if length <= 0 || width <= 0 || height <= 0 {
		return ErrInvalidDimensions
	}
	
	p.Dimensions = &ProductDimensions{
		Length: length,
		Width:  width,
		Height: height,
	}
	
	p.UpdatedAt = time.Now()
	return nil
}

// Обновить вес товара
func (p *Product) UpdateWeight(weight float64) error {
	if weight <= 0 {
		return ErrInvalidWeight
	}
	
	p.Weight = &weight
	p.UpdatedAt = time.Now()
	return nil
}

// Добавить поставщика для товара
func (p *Product) AddSupplier(supplierID uint, sku, productURL string, price float64, inStock bool, deliveryDays int) error {
	if price < 0 {
		return ErrInvalidPrice
	}
	
	supplierProduct := SupplierProduct{
		ProductID:     p.ID,
		SupplierID:    supplierID,
		SupplierSKU:   sku,
		ProductURL:    productURL,
		Price:         price,
		Currency:      "RUB",
		InStock:       inStock,
		DeliveryDays:  deliveryDays,
		LastChecked:   time.Now(),
		FoundByParser: true,
	}
	
	p.SupplierProducts = append(p.SupplierProducts, supplierProduct)
	p.updatePriceStatistics()
	p.UpdatedAt = time.Now()
	
	return nil
}

// Отметить, что товар найден в тендере
func (p *Product) MarkFoundInTender() {
	p.TimesFound++
	now := time.Now()
	p.LastFoundAt = &now
	p.updatePopularityScore()
	p.UpdatedAt = now
}

// Добавить ключевые слова
func (p *Product) AddKeywords(keywords ...string) {
	for _, keyword := range keywords {
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		if keyword != "" && !contains(p.Keywords, keyword) {
			p.Keywords = append(p.Keywords, keyword)
		}
	}
	p.UpdatedAt = time.Now()
}

// 🔍 МЕТОДЫ ЗАПРОСОВ

// Проверить, доступен ли товар для поставки
func (p *Product) IsAvailable() bool {
	for _, sp := range p.SupplierProducts {
		if sp.InStock {
			return true
		}
	}
	return false
}

// Получить минимальную цену среди поставщиков
func (p *Product) GetMinPrice() float64 {
	if len(p.SupplierProducts) == 0 {
		if p.MinPrice != nil {
			return *p.MinPrice
		}
		return 0
	}
	
	minPrice := p.SupplierProducts[0].Price
	for _, sp := range p.SupplierProducts[1:] {
		if sp.Price < minPrice && sp.InStock {
			minPrice = sp.Price
		}
	}
	
	return minPrice
}

// Получить максимальную цену среди поставщиков
func (p *Product) GetMaxPrice() float64 {
	if len(p.SupplierProducts) == 0 {
		if p.MaxPrice != nil {
			return *p.MaxPrice
		}
		return 0
	}
	
	maxPrice := p.SupplierProducts[0].Price
	for _, sp := range p.SupplierProducts[1:] {
		if sp.Price > maxPrice && sp.InStock {
			maxPrice = sp.Price
		}
	}
	
	return maxPrice
}

// Получить минимальный срок поставки
func (p *Product) GetMinDeliveryDays() int {
	if len(p.SupplierProducts) == 0 {
		return 0
	}
	
	minDays := p.SupplierProducts[0].DeliveryDays
	for _, sp := range p.SupplierProducts[1:] {
		if sp.DeliveryDays < minDays && sp.InStock {
			minDays = sp.DeliveryDays
		}
	}
	
	return minDays
}

// Проверить, нужно ли обновить цены
func (p *Product) NeedsPriceUpdate() bool {
	if p.LastPriceUpdate == nil {
		return true
	}
	
	// Обновляем цены раз в неделю
	return time.Since(*p.LastPriceUpdate) > 7*24*time.Hour
}

// Проверить релевантность для поискового запроса
func (p *Product) IsRelevantFor(query string) bool {
	query = strings.ToLower(query)
	
	// Проверка в названии
	if strings.Contains(strings.ToLower(p.Name), query) {
		return true
	}
	
	// Проверка в ключевых словах
	for _, keyword := range p.Keywords {
		if strings.Contains(keyword, query) {
			return true
		}
	}
	
	// Проверка в синонимах
	for _, synonym := range p.Synonyms {
		if strings.Contains(strings.ToLower(synonym), query) {
			return true
		}
	}
	
	return false
}

// 🔧 ПРИВАТНЫЕ МЕТОДЫ

// Обновить статистику цен
func (p *Product) updatePriceStatistics() {
	if len(p.SupplierProducts) == 0 {
		return
	}
	
	var prices []float64
	for _, sp := range p.SupplierProducts {
		if sp.InStock {
			prices = append(prices, sp.Price)
		}
	}
	
	if len(prices) == 0 {
		return
	}
	
	// Минимальная и максимальная цена
	minPrice := prices[0]
	maxPrice := prices[0]
	sum := 0.0
	
	for _, price := range prices {
		if price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}
		sum += price
	}
	
	// Средняя цена
	avgPrice := sum / float64(len(prices))
	
	p.MinPrice = &minPrice
	p.MaxPrice = &maxPrice
	p.AveragePrice = &avgPrice
	
	now := time.Now()
	p.LastPriceUpdate = &now
}

// Обновить индекс популярности
func (p *Product) updatePopularityScore() {
	// Базовый скор на основе количества упоминаний
	baseScore := float64(p.TimesFound)
	
	// Бонус за недавнее обнаружение
	if p.LastFoundAt != nil {
		daysSinceLastFound := time.Since(*p.LastFoundAt).Hours() / 24
		if daysSinceLastFound < 30 {
			baseScore *= (1 + (30-daysSinceLastFound)/30*0.5) // До 50% бонуса
		}
	}
	
	// Бонус за количество поставщиков
	supplierBonus := float64(len(p.SupplierProducts)) * 0.1
	
	p.PopularityScore = baseScore + supplierBonus
}

// 🔧 ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ

func validateCategory(category ProductCategory) error {
	validCategories := map[ProductCategory]bool{
		CategoryDiagnostic:  true,
		CategorySurgical:    true,
		CategoryLaboratory:  true,
		CategoryReanimation: true,
		CategoryRadiology:   true,
		CategoryPhysical:    true,
		CategoryDental:      true,
		CategoryConsumables: true,
		CategoryFurniture:   true,
		CategoryOther:       true,
	}
	
	if !validCategories[category] {
		return ErrInvalidCategory
	}
	
	return nil
}

func extractKeywords(name, description string) []string {
	text := strings.ToLower(name + " " + description)
	
	// Простое извлечение ключевых слов (можно улучшить с помощью NLP)
	words := strings.Fields(text)
	
	var keywords []string
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:")
		if len(word) > 3 { // Игнорируем короткие слова
			keywords = append(keywords, word)
		}
	}
	
	return removeDuplicates(keywords)
}

func generateSynonyms(name string) []string {
	// Простой генератор синонимов (можно расширить)
	synonymMap := map[string][]string{
		"аппарат":      {"устройство", "прибор", "система"},
		"монитор":      {"дисплей", "экран"},
		"насос":        {"помпа"},
		"стол":         {"стойка", "тумба"},
		"шкаф":         {"ящик", "бокс"},
		"кровать":      {"койка", "ложе"},
		"кресло":       {"стул"},
		"анализатор":   {"анализатор крови", "гематологический анализатор"},
	}
	
	nameLower := strings.ToLower(name)
	var synonyms []string
	
	for key, values := range synonymMap {
		if strings.Contains(nameLower, key) {
			synonyms = append(synonyms, values...)
		}
	}
	
	return synonyms
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// 📊 МЕТОДЫ ДЛЯ UI

// Получить человекочитаемое название категории
func (c ProductCategory) DisplayName() string {
	switch c {
	case CategoryDiagnostic:
		return "Диагностическое оборудование"
	case CategorySurgical:
		return "Хирургическое оборудование"
	case CategoryLaboratory:
		return "Лабораторное оборудование"
	case CategoryReanimation:
		return "Реанимационное оборудование"
	case CategoryRadiology:
		return "Рентгенологическое оборудование"
	case CategoryPhysical:
		return "Физиотерапевтическое оборудование"
	case CategoryDental:
		return "Стоматологическое оборудование"
	case CategoryConsumables:
		return "Расходные материалы"
	case CategoryFurniture:
		return "Медицинская мебель"
	default:
		return "Прочее"
	}
}

// Получить краткое описание для поиска
func (p *Product) GetSearchSummary() string {
	summary := p.Name
	if p.Manufacturer != "" {
		summary += fmt.Sprintf(" (%s)", p.Manufacturer)
	}
	if p.AveragePrice != nil {
		summary += fmt.Sprintf(" - от %.0f руб.", *p.AveragePrice)
	}
	return summary
}