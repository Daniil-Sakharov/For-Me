package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// 🏛️ ДОМЕННАЯ СУЩНОСТЬ ЗАКАЗА
// Содержит всю бизнес-логику, связанную с заказом
// НЕ ЗАВИСИТ от внешних технологий (БД, HTTP, etc.)

// Order представляет заказ в системе e-commerce
type Order struct {
	id          uuid.UUID
	customerID  uuid.UUID
	items       []OrderItem
	status      Status
	totalAmount Money
	currency    string
	createdAt   time.Time
	updatedAt   time.Time
	
	// Адресная информация
	shippingAddress Address
	billingAddress  Address
	
	// Дополнительная информация
	notes       string
	priority    Priority
	source      Source // web, mobile, api
}

// OrderItem представляет товар в заказе
type OrderItem struct {
	productID uuid.UUID
	quantity  int
	unitPrice Money
	discount  Money
}

// Address представляет адрес доставки/счета
type Address struct {
	street     string
	city       string
	postalCode string
	country    string
	phone      string
}

// Money представляет денежную сумму (Value Object)
type Money struct {
	amount   int64 // в копейках для точности
	currency string
}

// Status представляет статус заказа
type Status int

const (
	StatusPending Status = iota
	StatusValidated
	StatusPaymentProcessing
	StatusPaid
	StatusInventoryChecked
	StatusFulfillment
	StatusShipped
	StatusDelivered
	StatusCancelled
	StatusRefunded
)

// Priority представляет приоритет заказа
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityUrgent
)

// Source представляет источник заказа
type Source string

const (
	SourceWeb    Source = "web"
	SourceMobile Source = "mobile"
	SourceAPI    Source = "api"
)

// 🏗️ КОНСТРУКТОРЫ (Фабричные методы)

// NewOrder создает новый заказ
func NewOrder(customerID uuid.UUID, items []OrderItem, shippingAddr Address) (*Order, error) {
	if customerID == uuid.Nil {
		return nil, errors.New("customer ID cannot be empty")
	}
	
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}
	
	if err := shippingAddr.Validate(); err != nil {
		return nil, err
	}
	
	order := &Order{
		id:              uuid.New(),
		customerID:      customerID,
		items:           items,
		status:          StatusPending,
		shippingAddress: shippingAddr,
		billingAddress:  shippingAddr, // По умолчанию такой же как доставки
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
		priority:        PriorityNormal,
		source:          SourceWeb,
		currency:        "RUB",
	}
	
	// Рассчитываем общую стоимость
	order.calculateTotalAmount()
	
	return order, nil
}

// 💰 БИЗНЕС-ЛОГИКА РАСЧЕТОВ

// calculateTotalAmount рассчитывает общую стоимость заказа
func (o *Order) calculateTotalAmount() {
	var total int64
	for _, item := range o.items {
		itemTotal := item.unitPrice.amount * int64(item.quantity)
		itemTotal -= item.discount.amount
		total += itemTotal
	}
	
	o.totalAmount = Money{
		amount:   total,
		currency: o.currency,
	}
}

// AddItem добавляет товар в заказ
func (o *Order) AddItem(productID uuid.UUID, quantity int, unitPrice Money) error {
	if productID == uuid.Nil {
		return errors.New("product ID cannot be empty")
	}
	
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	
	if unitPrice.amount <= 0 {
		return errors.New("unit price must be positive")
	}
	
	// Проверяем, что заказ можно изменять
	if !o.CanModifyItems() {
		return errors.New("cannot modify items in current status")
	}
	
	item := OrderItem{
		productID: productID,
		quantity:  quantity,
		unitPrice: unitPrice,
		discount:  Money{amount: 0, currency: unitPrice.currency},
	}
	
	o.items = append(o.items, item)
	o.calculateTotalAmount()
	o.updatedAt = time.Now()
	
	return nil
}

// ApplyDiscount применяет скидку к товару
func (o *Order) ApplyDiscount(productID uuid.UUID, discount Money) error {
	if !o.CanModifyItems() {
		return errors.New("cannot apply discount in current status")
	}
	
	for i := range o.items {
		if o.items[i].productID == productID {
			o.items[i].discount = discount
			o.calculateTotalAmount()
			o.updatedAt = time.Now()
			return nil
		}
	}
	
	return errors.New("product not found in order")
}

// 📊 БИЗНЕС-ЛОГИКА СТАТУСОВ

// CanModifyItems проверяет, можно ли изменять товары в заказе
func (o *Order) CanModifyItems() bool {
	return o.status == StatusPending || o.status == StatusValidated
}

// CanCancel проверяет, можно ли отменить заказ
func (o *Order) CanCancel() bool {
	return o.status != StatusShipped && 
		   o.status != StatusDelivered && 
		   o.status != StatusCancelled &&
		   o.status != StatusRefunded
}

// Cancel отменяет заказ
func (o *Order) Cancel() error {
	if !o.CanCancel() {
		return errors.New("cannot cancel order in current status")
	}
	
	o.status = StatusCancelled
	o.updatedAt = time.Now()
	return nil
}

// MarkAsValidated помечает заказ как валидированный
func (o *Order) MarkAsValidated() error {
	if o.status != StatusPending {
		return errors.New("can only validate pending orders")
	}
	
	o.status = StatusValidated
	o.updatedAt = time.Now()
	return nil
}

// MarkAsPaymentProcessing помечает заказ как обрабатываемый платеж
func (o *Order) MarkAsPaymentProcessing() error {
	if o.status != StatusValidated {
		return errors.New("can only process payment for validated orders")
	}
	
	o.status = StatusPaymentProcessing
	o.updatedAt = time.Now()
	return nil
}

// MarkAsPaid помечает заказ как оплаченный
func (o *Order) MarkAsPaid() error {
	if o.status != StatusPaymentProcessing {
		return errors.New("can only mark as paid orders in payment processing")
	}
	
	o.status = StatusPaid
	o.updatedAt = time.Now()
	return nil
}

// MarkAsInventoryChecked помечает заказ как проверенный по складу
func (o *Order) MarkAsInventoryChecked() error {
	if o.status != StatusPaid {
		return errors.New("can only check inventory for paid orders")
	}
	
	o.status = StatusInventoryChecked
	o.updatedAt = time.Time{}
	return nil
}

// MarkAsShipped помечает заказ как отправленный
func (o *Order) MarkAsShipped() error {
	if o.status != StatusFulfillment {
		return errors.New("can only ship orders in fulfillment")
	}
	
	o.status = StatusShipped
	o.updatedAt = time.Now()
	return nil
}

// MarkAsDelivered помечает заказ как доставленный
func (o *Order) MarkAsDelivered() error {
	if o.status != StatusShipped {
		return errors.New("can only deliver shipped orders")
	}
	
	o.status = StatusDelivered
	o.updatedAt = time.Now()
	return nil
}

// 🔍 ГЕТТЕРЫ (для доступа к приватным полям)

func (o *Order) ID() uuid.UUID          { return o.id }
func (o *Order) CustomerID() uuid.UUID  { return o.customerID }
func (o *Order) Items() []OrderItem     { return o.items }
func (o *Order) Status() Status         { return o.status }
func (o *Order) TotalAmount() Money     { return o.totalAmount }
func (o *Order) Currency() string       { return o.currency }
func (o *Order) CreatedAt() time.Time   { return o.createdAt }
func (o *Order) UpdatedAt() time.Time   { return o.updatedAt }
func (o *Order) ShippingAddress() Address { return o.shippingAddress }
func (o *Order) BillingAddress() Address  { return o.billingAddress }
func (o *Order) Notes() string          { return o.notes }
func (o *Order) Priority() Priority     { return o.priority }
func (o *Order) Source() Source         { return o.source }

// 📊 ДОПОЛНИТЕЛЬНЫЕ МЕТОДЫ

// SetPriority устанавливает приоритет заказа
func (o *Order) SetPriority(priority Priority) {
	o.priority = priority
	o.updatedAt = time.Now()
}

// SetNotes устанавливает заметки к заказу
func (o *Order) SetNotes(notes string) {
	o.notes = notes
	o.updatedAt = time.Now()
}

// SetBillingAddress устанавливает адрес для счета
func (o *Order) SetBillingAddress(address Address) error {
	if err := address.Validate(); err != nil {
		return err
	}
	
	o.billingAddress = address
	o.updatedAt = time.Time{}
	return nil
}

// ItemsCount возвращает общее количество товаров в заказе
func (o *Order) ItemsCount() int {
	total := 0
	for _, item := range o.items {
		total += item.quantity
	}
	return total
}

// 🔍 МЕТОДЫ ДЛЯ OrderItem

func (oi *OrderItem) ProductID() uuid.UUID { return oi.productID }
func (oi *OrderItem) Quantity() int        { return oi.quantity }
func (oi *OrderItem) UnitPrice() Money     { return oi.unitPrice }
func (oi *OrderItem) Discount() Money      { return oi.discount }

// Total возвращает общую стоимость позиции (количество * цена - скидка)
func (oi *OrderItem) Total() Money {
	total := oi.unitPrice.amount * int64(oi.quantity)
	total -= oi.discount.amount
	
	return Money{
		amount:   total,
		currency: oi.unitPrice.currency,
	}
}

// 🔍 МЕТОДЫ ДЛЯ Address

func (a *Address) Street() string     { return a.street }
func (a *Address) City() string       { return a.city }
func (a *Address) PostalCode() string { return a.postalCode }
func (a *Address) Country() string    { return a.country }
func (a *Address) Phone() string      { return a.phone }

// Validate валидирует адрес
func (a *Address) Validate() error {
	if a.street == "" {
		return errors.New("street is required")
	}
	if a.city == "" {
		return errors.New("city is required")
	}
	if a.postalCode == "" {
		return errors.New("postal code is required")
	}
	if a.country == "" {
		return errors.New("country is required")
	}
	return nil
}

// 🔍 МЕТОДЫ ДЛЯ Money

func (m *Money) Amount() int64   { return m.amount }
func (m *Money) Currency() string { return m.currency }

// ToFloat возвращает сумму в рублях (деля на 100 копеек)
func (m *Money) ToFloat() float64 {
	return float64(m.amount) / 100.0
}

// String возвращает строковое представление денег
func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.ToFloat(), m.currency)
}

// 🔍 МЕТОДЫ ДЛЯ Status

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusValidated:
		return "validated"
	case StatusPaymentProcessing:
		return "payment_processing"
	case StatusPaid:
		return "paid"
	case StatusInventoryChecked:
		return "inventory_checked"
	case StatusFulfillment:
		return "fulfillment"
	case StatusShipped:
		return "shipped"
	case StatusDelivered:
		return "delivered"
	case StatusCancelled:
		return "cancelled"
	case StatusRefunded:
		return "refunded"
	default:
		return "unknown"
	}
}

// 🔍 МЕТОДЫ ДЛЯ Priority

func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityUrgent:
		return "urgent"
	default:
		return "unknown"
	}
}