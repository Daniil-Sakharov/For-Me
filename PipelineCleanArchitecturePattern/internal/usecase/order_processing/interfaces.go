package order_processing

import (
	"context"
	"time"

	"pipeline-clean-architecture/internal/domain/order"
	"pipeline-clean-architecture/internal/domain/payment"
	"pipeline-clean-architecture/internal/domain/product"

	"github.com/google/uuid"
)

// 💼 ИНТЕРФЕЙСЫ USE CASE СЛОЯ
// Определяют контракты для бизнес-логики обработки заказов
// Эти интерфейсы реализуются в Application Layer

// OrderProcessor определяет основные операции с заказами
type OrderProcessor interface {
	// 📝 СОЗДАНИЕ И ВАЛИДАЦИЯ ЗАКАЗОВ
	
	// CreateOrder создает новый заказ и валидирует его
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error)
	
	// ValidateOrder проверяет корректность заказа
	ValidateOrder(ctx context.Context, orderID uuid.UUID) (*ValidateOrderResponse, error)
	
	// 💰 ОБРАБОТКА ПЛАТЕЖЕЙ
	
	// ProcessPayment обрабатывает платеж для заказа
	ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*ProcessPaymentResponse, error)
	
	// 📦 ПРОВЕРКА СКЛАДА
	
	// CheckInventory проверяет наличие товаров на складе
	CheckInventory(ctx context.Context, orderID uuid.UUID) (*CheckInventoryResponse, error)
	
	// ReserveInventory резервирует товары для заказа
	ReserveInventory(ctx context.Context, orderID uuid.UUID) (*ReserveInventoryResponse, error)
	
	// 🚚 ОРГАНИЗАЦИЯ ДОСТАВКИ
	
	// ArrangeDelivery организует доставку заказа
	ArrangeDelivery(ctx context.Context, orderID uuid.UUID) (*ArrangeDeliveryResponse, error)
	
	// 📧 УВЕДОМЛЕНИЯ
	
	// SendNotifications отправляет уведомления о статусе заказа
	SendNotifications(ctx context.Context, orderID uuid.UUID, event NotificationEvent) error
	
	// 📊 СТАТИСТИКА И АНАЛИТИКА
	
	// GetOrderStats возвращает статистику по заказам
	GetOrderStats(ctx context.Context, req GetOrderStatsRequest) (*GetOrderStatsResponse, error)
}

// PipelineExecutor определяет выполнение пайплайна обработки заказов
type PipelineExecutor interface {
	// 🔄 ПАЙПЛАЙН ОБРАБОТКИ
	
	// ExecuteOrderPipeline выполняет полный пайплайн обработки заказа
	ExecuteOrderPipeline(ctx context.Context, orderID uuid.UUID) (*PipelineResult, error)
	
	// ExecuteStep выполняет отдельный шаг пайплайна
	ExecuteStep(ctx context.Context, step PipelineStep, orderID uuid.UUID) (*StepResult, error)
	
	// GetPipelineStatus возвращает статус выполнения пайплайна
	GetPipelineStatus(ctx context.Context, orderID uuid.UUID) (*PipelineStatus, error)
}

// 📊 REQUEST/RESPONSE СТРУКТУРЫ

// CreateOrderRequest запрос на создание заказа
type CreateOrderRequest struct {
	CustomerID      uuid.UUID             `json:"customer_id" validate:"required"`
	Items           []CreateOrderItem     `json:"items" validate:"required,dive"`
	ShippingAddress order.Address         `json:"shipping_address" validate:"required"`
	BillingAddress  *order.Address        `json:"billing_address,omitempty"`
	PaymentMethod   payment.Method        `json:"payment_method" validate:"required"`
	Notes           string                `json:"notes,omitempty"`
	Priority        order.Priority        `json:"priority"`
	Source          order.Source          `json:"source"`
}

// CreateOrderItem товар в запросе создания заказа
type CreateOrderItem struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

// CreateOrderResponse ответ создания заказа
type CreateOrderResponse struct {
	Order     *order.Order `json:"order"`
	Warnings  []string     `json:"warnings,omitempty"`
	NextSteps []string     `json:"next_steps"`
}

// ValidateOrderResponse ответ валидации заказа
type ValidateOrderResponse struct {
	IsValid      bool                    `json:"is_valid"`
	Errors       []ValidationError       `json:"errors,omitempty"`
	Warnings     []ValidationWarning     `json:"warnings,omitempty"`
	ValidatedAt  time.Time              `json:"validated_at"`
	ValidatedBy  string                 `json:"validated_by"`
}

// ValidationError ошибка валидации
type ValidationError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ValidationWarning предупреждение валидации
type ValidationWarning struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ProcessPaymentRequest запрос обработки платежа
type ProcessPaymentRequest struct {
	OrderID       uuid.UUID         `json:"order_id" validate:"required"`
	PaymentMethod payment.Method    `json:"payment_method" validate:"required"`
	Provider      payment.Provider  `json:"provider" validate:"required"`
	CardToken     string           `json:"card_token,omitempty"`
	ReturnURL     string           `json:"return_url,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// ProcessPaymentResponse ответ обработки платежа
type ProcessPaymentResponse struct {
	Payment       *payment.Payment `json:"payment"`
	Status        payment.Status   `json:"status"`
	RedirectURL   string          `json:"redirect_url,omitempty"`
	QRCodeData    string          `json:"qr_code_data,omitempty"`
	ExpiresAt     *time.Time      `json:"expires_at,omitempty"`
}

// CheckInventoryResponse ответ проверки склада
type CheckInventoryResponse struct {
	AvailableItems   []InventoryItem `json:"available_items"`
	UnavailableItems []InventoryItem `json:"unavailable_items"`
	IsFullyAvailable bool           `json:"is_fully_available"`
	EstimatedDelay   *time.Duration `json:"estimated_delay,omitempty"`
	Alternatives     []Alternative  `json:"alternatives,omitempty"`
}

// InventoryItem информация о товаре на складе
type InventoryItem struct {
	ProductID     uuid.UUID `json:"product_id"`
	RequestedQty  int       `json:"requested_qty"`
	AvailableQty  int       `json:"available_qty"`
	ReservedQty   int       `json:"reserved_qty"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	EstimatedDate *time.Time `json:"estimated_date,omitempty"`
}

// Alternative альтернативный товар
type Alternative struct {
	ProductID   uuid.UUID      `json:"product_id"`
	Name        string         `json:"name"`
	Price       product.Money  `json:"price"`
	Similarity  float64        `json:"similarity"` // 0-1
	Description string         `json:"description"`
}

// ReserveInventoryResponse ответ резервирования склада
type ReserveInventoryResponse struct {
	ReservedItems []ReservedItem `json:"reserved_items"`
	ReservationID uuid.UUID     `json:"reservation_id"`
	ExpiresAt     time.Time     `json:"expires_at"`
	TotalReserved int           `json:"total_reserved"`
}

// ReservedItem зарезервированный товар
type ReservedItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	ReservedAt  time.Time `json:"reserved_at"`
}

// ArrangeDeliveryResponse ответ организации доставки
type ArrangeDeliveryResponse struct {
	DeliveryID      uuid.UUID               `json:"delivery_id"`
	CarrierID       uuid.UUID              `json:"carrier_id"`
	CarrierName     string                 `json:"carrier_name"`
	TrackingNumber  string                 `json:"tracking_number"`
	EstimatedDate   time.Time              `json:"estimated_date"`
	DeliveryOptions []DeliveryOption       `json:"delivery_options"`
	Cost            product.Money          `json:"cost"`
}

// DeliveryOption опция доставки
type DeliveryOption struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Cost        product.Money `json:"cost"`
	Duration    time.Duration `json:"duration"`
	IsSelected  bool         `json:"is_selected"`
}

// NotificationEvent событие для уведомления
type NotificationEvent string

const (
	NotificationOrderCreated    NotificationEvent = "order_created"
	NotificationOrderValidated  NotificationEvent = "order_validated"
	NotificationPaymentReceived NotificationEvent = "payment_received"
	NotificationPaymentFailed   NotificationEvent = "payment_failed"
	NotificationInventoryReserved NotificationEvent = "inventory_reserved"
	NotificationShipped         NotificationEvent = "shipped"
	NotificationDelivered       NotificationEvent = "delivered"
	NotificationCancelled       NotificationEvent = "cancelled"
)

// GetOrderStatsRequest запрос статистики заказов
type GetOrderStatsRequest struct {
	CustomerID *uuid.UUID         `json:"customer_id,omitempty"`
	DateFrom   *time.Time        `json:"date_from,omitempty"`
	DateTo     *time.Time        `json:"date_to,omitempty"`
	Status     *order.Status     `json:"status,omitempty"`
	Source     *order.Source     `json:"source,omitempty"`
	GroupBy    StatsGroupBy      `json:"group_by"`
}

// StatsGroupBy группировка статистики
type StatsGroupBy string

const (
	StatsGroupByDay    StatsGroupBy = "day"
	StatsGroupByWeek   StatsGroupBy = "week"
	StatsGroupByMonth  StatsGroupBy = "month"
	StatsGroupByStatus StatsGroupBy = "status"
	StatsGroupBySource StatsGroupBy = "source"
)

// GetOrderStatsResponse ответ статистики заказов
type GetOrderStatsResponse struct {
	TotalOrders    int                    `json:"total_orders"`
	TotalRevenue   product.Money          `json:"total_revenue"`
	AverageOrder   product.Money          `json:"average_order"`
	StatusBreakdown map[string]int        `json:"status_breakdown"`
	TrendData      []StatsTrendPoint      `json:"trend_data"`
	TopProducts    []StatsProduct         `json:"top_products"`
}

// StatsTrendPoint точка тренда
type StatsTrendPoint struct {
	Period time.Time     `json:"period"`
	Count  int           `json:"count"`
	Revenue product.Money `json:"revenue"`
}

// StatsProduct статистика по товару
type StatsProduct struct {
	ProductID uuid.UUID     `json:"product_id"`
	Name      string        `json:"name"`
	Quantity  int           `json:"quantity"`
	Revenue   product.Money `json:"revenue"`
}

// 🔄 ПАЙПЛАЙН СТРУКТУРЫ

// PipelineStep шаг пайплайна
type PipelineStep string

const (
	StepValidation         PipelineStep = "validation"
	StepPaymentProcessing  PipelineStep = "payment_processing"
	StepInventoryCheck     PipelineStep = "inventory_check"
	StepInventoryReserve   PipelineStep = "inventory_reserve"
	StepDeliveryArrange    PipelineStep = "delivery_arrange"
	StepNotifications      PipelineStep = "notifications"
)

// PipelineResult результат выполнения пайплайна
type PipelineResult struct {
	OrderID     uuid.UUID               `json:"order_id"`
	Success     bool                   `json:"success"`
	Steps       []StepResult           `json:"steps"`
	Duration    time.Duration          `json:"duration"`
	CompletedAt time.Time             `json:"completed_at"`
	FinalStatus order.Status          `json:"final_status"`
	Errors      []PipelineError        `json:"errors,omitempty"`
}

// StepResult результат выполнения шага
type StepResult struct {
	Step        PipelineStep   `json:"step"`
	Success     bool          `json:"success"`
	Duration    time.Duration `json:"duration"`
	StartedAt   time.Time     `json:"started_at"`
	CompletedAt time.Time     `json:"completed_at"`
	Data        interface{}   `json:"data,omitempty"`
	Error       *PipelineError `json:"error,omitempty"`
}

// PipelineError ошибка пайплайна
type PipelineError struct {
	Step    PipelineStep `json:"step"`
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Retry   bool         `json:"retry"`
}

// PipelineStatus статус пайплайна
type PipelineStatus struct {
	OrderID        uuid.UUID     `json:"order_id"`
	CurrentStep    PipelineStep  `json:"current_step"`
	CompletedSteps []PipelineStep `json:"completed_steps"`
	Progress       float64       `json:"progress"` // 0-1
	StartedAt      time.Time     `json:"started_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	EstimatedEnd   *time.Time    `json:"estimated_end,omitempty"`
}