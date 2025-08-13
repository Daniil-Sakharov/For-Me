package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// 🏛️ ДОМЕННАЯ СУЩНОСТЬ ПРОДУКТА
// Представляет товар в системе e-commerce

// Product представляет товар в каталоге
type Product struct {
	id           uuid.UUID
	name         string
	description  string
	sku          string    // артикул
	category     Category
	price        Money
	weight       Weight
	dimensions   Dimensions
	stockLevel   int
	minStockLevel int
	status       Status
	images       []string
	tags         []string
	createdAt    time.Time
	updatedAt    time.Time
}

// Category представляет категорию товара
type Category struct {
	id   uuid.UUID
	name string
	path string // иерархический путь: "Electronics > Computers > Laptops"
}

// Money представляет цену товара
type Money struct {
	amount   int64  // в копейках
	currency string
}

// Weight представляет вес товара
type Weight struct {
	value int64  // в граммах
	unit  string // g, kg
}

// Dimensions представляет размеры товара
type Dimensions struct {
	length int64 // в миллиметрах
	width  int64
	height int64
	unit   string // mm, cm, m
}

// Status представляет статус товара
type Status int

const (
	StatusActive Status = iota
	StatusInactive
	StatusOutOfStock
	StatusDiscontinued
)

// 🏗️ КОНСТРУКТОРЫ

// NewProduct создает новый товар
func NewProduct(name, description, sku string, category Category, price Money) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	
	if sku == "" {
		return nil, errors.New("product SKU cannot be empty")
	}
	
	if price.amount <= 0 {
		return nil, errors.New("product price must be positive")
	}
	
	return &Product{
		id:            uuid.New(),
		name:          name,
		description:   description,
		sku:           sku,
		category:      category,
		price:         price,
		status:        StatusActive,
		stockLevel:    0,
		minStockLevel: 1,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
		images:        make([]string, 0),
		tags:          make([]string, 0),
	}, nil
}

// 📦 БИЗНЕС-ЛОГИКА СКЛАДА

// IsInStock проверяет, есть ли товар в наличии
func (p *Product) IsInStock() bool {
	return p.stockLevel > 0 && p.status == StatusActive
}

// IsLowStock проверяет, мало ли товара на складе
func (p *Product) IsLowStock() bool {
	return p.stockLevel <= p.minStockLevel && p.stockLevel > 0
}

// CanReserve проверяет, можно ли зарезервировать количество
func (p *Product) CanReserve(quantity int) bool {
	return p.IsInStock() && p.stockLevel >= quantity
}

// Reserve резервирует товар
func (p *Product) Reserve(quantity int) error {
	if !p.CanReserve(quantity) {
		return errors.New("insufficient stock")
	}
	
	p.stockLevel -= quantity
	p.updatedAt = time.Now()
	
	if p.stockLevel == 0 {
		p.status = StatusOutOfStock
	}
	
	return nil
}

// RestoreStock восстанавливает товар (отмена резерва)
func (p *Product) RestoreStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	
	p.stockLevel += quantity
	p.updatedAt = time.Now()
	
	if p.status == StatusOutOfStock && p.stockLevel > 0 {
		p.status = StatusActive
	}
	
	return nil
}

// UpdateStock обновляет количество на складе
func (p *Product) UpdateStock(newStock int) error {
	if newStock < 0 {
		return errors.New("stock cannot be negative")
	}
	
	p.stockLevel = newStock
	p.updatedAt = time.Now()
	
	if newStock == 0 {
		p.status = StatusOutOfStock
	} else if p.status == StatusOutOfStock {
		p.status = StatusActive
	}
	
	return nil
}

// 💰 БИЗНЕС-ЛОГИКА ЦЕН

// UpdatePrice обновляет цену товара
func (p *Product) UpdatePrice(newPrice Money) error {
	if newPrice.amount <= 0 {
		return errors.New("price must be positive")
	}
	
	if newPrice.currency != p.price.currency {
		return errors.New("currency mismatch")
	}
	
	p.price = newPrice
	p.updatedAt = time.Now()
	return nil
}

// CalculateDiscountedPrice рассчитывает цену со скидкой
func (p *Product) CalculateDiscountedPrice(discountPercent float64) (Money, error) {
	if discountPercent < 0 || discountPercent > 100 {
		return Money{}, errors.New("discount must be between 0 and 100")
	}
	
	discountAmount := float64(p.price.amount) * discountPercent / 100
	finalAmount := p.price.amount - int64(discountAmount)
	
	return Money{
		amount:   finalAmount,
		currency: p.price.currency,
	}, nil
}

// 📊 УПРАВЛЕНИЕ МЕТАДАННЫМИ

// AddTag добавляет тег к товару
func (p *Product) AddTag(tag string) {
	if tag == "" {
		return
	}
	
	// Проверяем, что тег еще не добавлен
	for _, existingTag := range p.tags {
		if existingTag == tag {
			return
		}
	}
	
	p.tags = append(p.tags, tag)
	p.updatedAt = time.Now()
}

// RemoveTag удаляет тег
func (p *Product) RemoveTag(tag string) {
	for i, existingTag := range p.tags {
		if existingTag == tag {
			p.tags = append(p.tags[:i], p.tags[i+1:]...)
			p.updatedAt = time.Now()
			return
		}
	}
}

// AddImage добавляет изображение
func (p *Product) AddImage(imageURL string) error {
	if imageURL == "" {
		return errors.New("image URL cannot be empty")
	}
	
	p.images = append(p.images, imageURL)
	p.updatedAt = time.Now()
	return nil
}

// SetCategory устанавливает категорию
func (p *Product) SetCategory(category Category) {
	p.category = category
	p.updatedAt = time.Now()
}

// SetDimensions устанавливает размеры
func (p *Product) SetDimensions(dims Dimensions) error {
	if dims.length <= 0 || dims.width <= 0 || dims.height <= 0 {
		return errors.New("dimensions must be positive")
	}
	
	p.dimensions = dims
	p.updatedAt = time.Now()
	return nil
}

// SetWeight устанавливает вес
func (p *Product) SetWeight(weight Weight) error {
	if weight.value <= 0 {
		return errors.New("weight must be positive")
	}
	
	p.weight = weight
	p.updatedAt = time.Now()
	return nil
}

// 📊 СТАТУС И ЖИЗНЕННЫЙ ЦИКЛ

// Activate активирует товар
func (p *Product) Activate() error {
	if p.status == StatusDiscontinued {
		return errors.New("cannot activate discontinued product")
	}
	
	if p.stockLevel > 0 {
		p.status = StatusActive
	} else {
		p.status = StatusOutOfStock
	}
	
	p.updatedAt = time.Now()
	return nil
}

// Deactivate деактивирует товар
func (p *Product) Deactivate() {
	p.status = StatusInactive
	p.updatedAt = time.Now()
}

// Discontinue снимает товар с продажи
func (p *Product) Discontinue() {
	p.status = StatusDiscontinued
	p.updatedAt = time.Now()
}

// 🔍 ГЕТТЕРЫ

func (p *Product) ID() uuid.UUID       { return p.id }
func (p *Product) Name() string        { return p.name }
func (p *Product) Description() string { return p.description }
func (p *Product) SKU() string         { return p.sku }
func (p *Product) Category() Category  { return p.category }
func (p *Product) Price() Money        { return p.price }
func (p *Product) Weight() Weight      { return p.weight }
func (p *Product) Dimensions() Dimensions { return p.dimensions }
func (p *Product) StockLevel() int     { return p.stockLevel }
func (p *Product) MinStockLevel() int  { return p.minStockLevel }
func (p *Product) Status() Status      { return p.status }
func (p *Product) Images() []string    { return p.images }
func (p *Product) Tags() []string      { return p.tags }
func (p *Product) CreatedAt() time.Time { return p.createdAt }
func (p *Product) UpdatedAt() time.Time { return p.updatedAt }

// 🔍 МЕТОДЫ ДЛЯ ВСПОМОГАТЕЛЬНЫХ ТИПОВ

// Category methods
func (c Category) ID() uuid.UUID { return c.id }
func (c Category) Name() string  { return c.name }
func (c Category) Path() string  { return c.path }

// Money methods
func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }
func (m Money) ToFloat() float64 { return float64(m.amount) / 100.0 }

// Weight methods
func (w Weight) Value() int64   { return w.value }
func (w Weight) Unit() string   { return w.unit }
func (w Weight) ToKilograms() float64 { return float64(w.value) / 1000.0 }

// Dimensions methods
func (d Dimensions) Length() int64 { return d.length }
func (d Dimensions) Width() int64  { return d.width }
func (d Dimensions) Height() int64 { return d.height }
func (d Dimensions) Unit() string  { return d.unit }

// Volume возвращает объем в кубических миллиметрах
func (d Dimensions) Volume() int64 {
	return d.length * d.width * d.height
}

// Status methods
func (s Status) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	case StatusOutOfStock:
		return "out_of_stock"
	case StatusDiscontinued:
		return "discontinued"
	default:
		return "unknown"
	}
}