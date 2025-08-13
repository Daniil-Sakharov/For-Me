package payment

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// 🏛️ ДОМЕННАЯ СУЩНОСТЬ ПЛАТЕЖА

// Payment представляет платеж в системе
type Payment struct {
	id              uuid.UUID
	orderID         uuid.UUID
	customerID      uuid.UUID
	amount          Money
	method          Method
	status          Status
	provider        Provider
	externalID      string    // ID во внешней платежной системе
	cardLast4       string    // последние 4 цифры карты
	cardBrand       string    // Visa, MasterCard, etc.
	fee             Money     // комиссия
	description     string
	metadata        map[string]string
	processedAt     *time.Time
	createdAt       time.Time
	updatedAt       time.Time
}

// Money представляет денежную сумму
type Money struct {
	amount   int64  // в копейках
	currency string
}

// Method представляет способ оплаты
type Method string

const (
	MethodCard           Method = "card"
	MethodBankTransfer   Method = "bank_transfer"
	MethodDigitalWallet  Method = "digital_wallet"
	MethodCryptocurrency Method = "cryptocurrency"
	MethodCash           Method = "cash"
)

// Status представляет статус платежа
type Status int

const (
	StatusPending Status = iota
	StatusProcessing
	StatusSucceeded
	StatusFailed
	StatusCancelled
	StatusRefunded
	StatusPartiallyRefunded
)

// Provider представляет платежный провайдер
type Provider string

const (
	ProviderStripe    Provider = "stripe"
	ProviderPayPal    Provider = "paypal"
	ProviderSberbank  Provider = "sberbank"
	ProviderYandex    Provider = "yandex"
	ProviderInternal  Provider = "internal"
)

// 🏗️ КОНСТРУКТОРЫ

// NewPayment создает новый платеж
func NewPayment(orderID, customerID uuid.UUID, amount Money, method Method, provider Provider) (*Payment, error) {
	if orderID == uuid.Nil {
		return nil, errors.New("order ID cannot be empty")
	}
	
	if customerID == uuid.Nil {
		return nil, errors.New("customer ID cannot be empty")
	}
	
	if amount.amount <= 0 {
		return nil, errors.New("payment amount must be positive")
	}
	
	if method == "" {
		return nil, errors.New("payment method cannot be empty")
	}
	
	return &Payment{
		id:         uuid.New(),
		orderID:    orderID,
		customerID: customerID,
		amount:     amount,
		method:     method,
		provider:   provider,
		status:     StatusPending,
		metadata:   make(map[string]string),
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}, nil
}

// 💳 БИЗНЕС-ЛОГИКА ОБРАБОТКИ ПЛАТЕЖЕЙ

// CanProcess проверяет, можно ли обработать платеж
func (p *Payment) CanProcess() bool {
	return p.status == StatusPending
}

// MarkAsProcessing помечает платеж как обрабатываемый
func (p *Payment) MarkAsProcessing() error {
	if !p.CanProcess() {
		return errors.New("cannot process payment in current status")
	}
	
	p.status = StatusProcessing
	p.updatedAt = time.Now()
	return nil
}

// MarkAsSucceeded помечает платеж как успешный
func (p *Payment) MarkAsSucceeded(externalID string) error {
	if p.status != StatusProcessing {
		return errors.New("can only succeed processing payments")
	}
	
	p.status = StatusSucceeded
	p.externalID = externalID
	now := time.Now()
	p.processedAt = &now
	p.updatedAt = now
	return nil
}

// MarkAsFailed помечает платеж как неудачный
func (p *Payment) MarkAsFailed(reason string) error {
	if p.status != StatusProcessing && p.status != StatusPending {
		return errors.New("can only fail pending or processing payments")
	}
	
	p.status = StatusFailed
	p.SetMetadata("failure_reason", reason)
	p.updatedAt = time.Now()
	return nil
}

// Cancel отменяет платеж
func (p *Payment) Cancel() error {
	if p.status != StatusPending {
		return errors.New("can only cancel pending payments")
	}
	
	p.status = StatusCancelled
	p.updatedAt = time.Now()
	return nil
}

// 💰 БИЗНЕС-ЛОГИКА ВОЗВРАТОВ

// CanRefund проверяет, можно ли сделать возврат
func (p *Payment) CanRefund() bool {
	return p.status == StatusSucceeded || p.status == StatusPartiallyRefunded
}

// Refund выполняет полный возврат
func (p *Payment) Refund(reason string) error {
	if !p.CanRefund() {
		return errors.New("cannot refund payment in current status")
	}
	
	p.status = StatusRefunded
	p.SetMetadata("refund_reason", reason)
	p.updatedAt = time.Now()
	return nil
}

// PartialRefund выполняет частичный возврат
func (p *Payment) PartialRefund(amount Money, reason string) error {
	if !p.CanRefund() {
		return errors.New("cannot refund payment in current status")
	}
	
	if amount.amount >= p.amount.amount {
		return errors.New("refund amount cannot exceed payment amount")
	}
	
	if amount.currency != p.amount.currency {
		return errors.New("refund currency must match payment currency")
	}
	
	p.status = StatusPartiallyRefunded
	p.SetMetadata("partial_refund_amount", fmt.Sprintf("%d", amount.amount))
	p.SetMetadata("refund_reason", reason)
	p.updatedAt = time.Now()
	return nil
}

// 📊 УПРАВЛЕНИЕ МЕТАДАННЫМИ

// SetMetadata устанавливает метаданные
func (p *Payment) SetMetadata(key, value string) {
	if p.metadata == nil {
		p.metadata = make(map[string]string)
	}
	
	p.metadata[key] = value
	p.updatedAt = time.Now()
}

// GetMetadata получает значение метаданных
func (p *Payment) GetMetadata(key string) (string, bool) {
	if p.metadata == nil {
		return "", false
	}
	
	value, exists := p.metadata[key]
	return value, exists
}

// SetCardInfo устанавливает информацию о карте
func (p *Payment) SetCardInfo(last4, brand string) error {
	if p.method != MethodCard {
		return errors.New("card info can only be set for card payments")
	}
	
	if len(last4) != 4 {
		return errors.New("last4 must be exactly 4 digits")
	}
	
	p.cardLast4 = last4
	p.cardBrand = brand
	p.updatedAt = time.Now()
	return nil
}

// SetFee устанавливает комиссию
func (p *Payment) SetFee(fee Money) error {
	if fee.currency != p.amount.currency {
		return errors.New("fee currency must match payment currency")
	}
	
	if fee.amount < 0 {
		return errors.New("fee cannot be negative")
	}
	
	p.fee = fee
	p.updatedAt = time.Now()
	return nil
}

// 💰 РАСЧЕТЫ

// NetAmount возвращает сумму к получению (без комиссии)
func (p *Payment) NetAmount() Money {
	return Money{
		amount:   p.amount.amount - p.fee.amount,
		currency: p.amount.currency,
	}
}

// IsExpired проверяет, истек ли платеж
func (p *Payment) IsExpired(timeout time.Duration) bool {
	if p.status != StatusPending {
		return false
	}
	
	return time.Since(p.createdAt) > timeout
}

// 🔍 ПРОВЕРКИ СТАТУСА

// IsSuccessful проверяет, успешен ли платеж
func (p *Payment) IsSuccessful() bool {
	return p.status == StatusSucceeded
}

// IsFailed проверяет, неудачен ли платеж
func (p *Payment) IsFailed() bool {
	return p.status == StatusFailed
}

// IsRefunded проверяет, возвращен ли платеж
func (p *Payment) IsRefunded() bool {
	return p.status == StatusRefunded || p.status == StatusPartiallyRefunded
}

// IsPending проверяет, ожидает ли платеж обработки
func (p *Payment) IsPending() bool {
	return p.status == StatusPending
}

// IsProcessing проверяет, обрабатывается ли платеж
func (p *Payment) IsProcessing() bool {
	return p.status == StatusProcessing
}

// 🔍 ГЕТТЕРЫ

func (p *Payment) ID() uuid.UUID           { return p.id }
func (p *Payment) OrderID() uuid.UUID      { return p.orderID }
func (p *Payment) CustomerID() uuid.UUID   { return p.customerID }
func (p *Payment) Amount() Money           { return p.amount }
func (p *Payment) Method() Method          { return p.method }
func (p *Payment) Status() Status          { return p.status }
func (p *Payment) Provider() Provider      { return p.provider }
func (p *Payment) ExternalID() string      { return p.externalID }
func (p *Payment) CardLast4() string       { return p.cardLast4 }
func (p *Payment) CardBrand() string       { return p.cardBrand }
func (p *Payment) Fee() Money              { return p.fee }
func (p *Payment) Description() string     { return p.description }
func (p *Payment) Metadata() map[string]string { return p.metadata }
func (p *Payment) ProcessedAt() *time.Time { return p.processedAt }
func (p *Payment) CreatedAt() time.Time    { return p.createdAt }
func (p *Payment) UpdatedAt() time.Time    { return p.updatedAt }

// 🔍 МЕТОДЫ ДЛЯ ВСПОМОГАТЕЛЬНЫХ ТИПОВ

// Money methods
func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }
func (m Money) ToFloat() float64 { return float64(m.amount) / 100.0 }

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.ToFloat(), m.currency)
}

// Method methods
func (m Method) String() string { return string(m) }

func (m Method) IsCard() bool {
	return m == MethodCard
}

func (m Method) IsDigital() bool {
	return m == MethodDigitalWallet || m == MethodCryptocurrency
}

// Status methods
func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusProcessing:
		return "processing"
	case StatusSucceeded:
		return "succeeded"
	case StatusFailed:
		return "failed"
	case StatusCancelled:
		return "cancelled"
	case StatusRefunded:
		return "refunded"
	case StatusPartiallyRefunded:
		return "partially_refunded"
	default:
		return "unknown"
	}
}

// Provider methods
func (p Provider) String() string { return string(p) }

func (p Provider) IsExternal() bool {
	return p != ProviderInternal
}