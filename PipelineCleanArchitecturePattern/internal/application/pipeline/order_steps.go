package pipeline

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pipeline-clean-architecture/internal/domain/order"
	"pipeline-clean-architecture/internal/domain/payment"
	"pipeline-clean-architecture/pkg/pipeline"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 🔄 КОНКРЕТНЫЕ ШАГИ ПАЙПЛАЙНА ДЛЯ ОБРАБОТКИ ЗАКАЗОВ
//
// ============================================================================
// APPLICATION LAYER - КООРДИНАЦИЯ БИЗНЕС-ПРОЦЕССОВ
// ============================================================================
//
// Этот файл содержит реализации шагов пайплайна для обработки заказов.
// Каждый шаг реализует интерфейс StepHandler из pipeline engine.
//
// ПРИНЦИПЫ APPLICATION LAYER:
//
// 1. 🎯 КООРДИНАЦИЯ, НЕ БИЗНЕС-ЛОГИКА:
//    - Шаги НЕ содержат бизнес-правила
//    - Они координируют вызовы доменных сущностей
//    - Используют сервисы для внешних операций
//
// 2. 🔗 ОРКЕСТРАЦИЯ ЗАВИСИМОСТЕЙ:
//    - Получают данные от предыдущих шагов
//    - Вызывают нужные сервисы (Payment, Inventory, etc.)
//    - Обновляют доменные сущности
//    - Передают результат следующему шагу
//
// 3. 📊 ОБРАБОТКА ОШИБОК:
//    - Логируют подробности выполнения
//    - Обрабатывают технические ошибки
//    - Возвращают структурированные результаты
//
// 4. 🔄 STATELESS:
//    - Каждый шаг независим
//    - Не хранят состояние между вызовами
//    - Получают все данные через StepData
//
// ============================================================================
// ПОТОК ДАННЫХ МЕЖДУ ШАГАМИ:
// ============================================================================
//
// Step 1 (Validate) → { order_id, is_valid, order_status }
//        ↓
// Step 2 (Payment) → { payment_id, payment_status }  
//        ↓
// Step 3 (Inventory) → { reservation_id, all_available }
//        ↓
// Step 4 (Notifications) → { notifications_sent }
//
// ============================================================================

// ============================================================================
// ШАГ 1: ВАЛИДАЦИЯ ЗАКАЗА
// ============================================================================
//
// НАЗНАЧЕНИЕ:
// Проверяет корректность заказа перед его дальнейшей обработкой
//
// ЧТО ПРОВЕРЯЕТ:
// - Статус заказа (должен быть Pending)
// - Наличие товаров в каталоге
// - Актуальность цен товаров  
// - Валидность адреса доставки
//
// ЗАВИСИМОСТИ:
// - order.Repository: для получения и сохранения заказа
// - ProductService: для проверки товаров и цен
// - zap.Logger: для логирования процесса
//
// ВХОДНЫЕ ДАННЫЕ:
// - order_id: ID заказа для валидации
//
// ВЫХОДНЫЕ ДАННЫЕ:  
// - is_valid: true/false - прошел ли заказ валидацию
// - validation_errors: список ошибок валидации
// - order_status: обновленный статус заказа
// ============================================================================

type ValidateOrderStep struct {
	logger         *zap.Logger       // Логгер для отслеживания выполнения
	orderRepo      order.Repository  // Репозиторий для работы с заказами (DOMAIN интерфейс)
	productService ProductService    // Сервис для проверки товаров (APPLICATION интерфейс)
}

// ProductService интерфейс для работы с товарами
type ProductService interface {
	CheckProductAvailability(ctx context.Context, productID uuid.UUID, quantity int) (bool, error)
	GetProductPrice(ctx context.Context, productID uuid.UUID) (order.Money, error)
}

// PaymentService интерфейс для работы с платежами
type PaymentService interface {
	ProcessPayment(ctx context.Context, orderID uuid.UUID, method payment.Method) (*payment.Payment, error)
	ValidatePaymentMethod(method payment.Method) error
}

// InventoryService интерфейс для работы со складом
type InventoryService interface {
	CheckAvailability(ctx context.Context, productID uuid.UUID, quantity int) (bool, error)
	ReserveItems(ctx context.Context, items []ReservationItem) (*Reservation, error)
	ReleaseReservation(ctx context.Context, reservationID uuid.UUID) error
}

// NotificationService интерфейс для уведомлений
type NotificationService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendSMS(ctx context.Context, phone, message string) error
	SendPushNotification(ctx context.Context, userID uuid.UUID, title, body string) error
}

// ReservationItem элемент резервирования
type ReservationItem struct {
	ProductID uuid.UUID
	Quantity  int
}

// Reservation резервирование товаров
type Reservation struct {
	ID        uuid.UUID
	Items     []ReservationItem
	ExpiresAt time.Time
}

// 🔄 STEP 1: VALIDATE ORDER

// NewValidateOrderStep создает новый шаг валидации заказа
func NewValidateOrderStep(logger *zap.Logger, orderRepo order.Repository, productService ProductService) *ValidateOrderStep {
	return &ValidateOrderStep{
		logger:         logger,
		orderRepo:      orderRepo,
		productService: productService,
	}
}

// Execute выполняет валидацию заказа
//
// АЛГОРИТМ ВЫПОЛНЕНИЯ:
// 1. Извлекаем order_id из входных данных
// 2. Получаем заказ из репозитория  
// 3. Проверяем статус заказа
// 4. Валидируем каждый товар в заказе
// 5. Проверяем адрес доставки
// 6. Обновляем статус заказа при успешной валидации
// 7. Возвращаем результат для следующего шага
func (s *ValidateOrderStep) Execute(ctx context.Context, data *pipeline.StepData) (*pipeline.StepResult, error) {
	s.logger.Info("Validating order", zap.String("execution_id", data.ID))
	
	startTime := time.Now()
	
	// ============================================================================
	// ЭТАП 1: ИЗВЛЕЧЕНИЕ И ВАЛИДАЦИЯ ВХОДНЫХ ДАННЫХ
	// ============================================================================
	
	// Получаем order_id из данных предыдущего шага (или начальных данных)
	// Это демонстрирует как данные передаются между шагами в пайплайне
	orderIDStr, ok := data.Input["order_id"].(string)
	if !ok {
		return nil, errors.New("order_id is required")
	}
	
	// Преобразуем string в UUID с проверкой корректности формата
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid order_id: %w", err)
	}
	
	// ============================================================================
	// ЭТАП 2: ПОЛУЧЕНИЕ ЗАКАЗА ИЗ ДОМЕНА
	// ============================================================================
	
	// Используем репозиторий для получения доменной сущности
	// Репозиторий определен в DOMAIN слое, реализован в INFRASTRUCTURE слое
	// Это пример Dependency Inversion Principle
	ord, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	// ============================================================================
	// ЭТАП 3: ВЫПОЛНЕНИЕ БИЗНЕС-ВАЛИДАЦИИ
	// ============================================================================
	
	// Собираем все ошибки валидации в список
	validationErrors := make([]string, 0)
	
	// 1. Проверяем статус заказа (бизнес-правило из DOMAIN слоя)
	if ord.Status() != order.StatusPending {
		validationErrors = append(validationErrors, 
			fmt.Sprintf("order status is %s, expected pending", ord.Status().String()))
	}
	
	// 2. Проверяем товары
	for _, item := range ord.Items() {
		// Проверяем доступность товара
		available, err := s.productService.CheckProductAvailability(ctx, item.ProductID(), item.Quantity())
		if err != nil {
			s.logger.Error("Failed to check product availability", 
				zap.String("product_id", item.ProductID().String()),
				zap.Error(err))
			validationErrors = append(validationErrors, 
				fmt.Sprintf("failed to check availability for product %s", item.ProductID().String()))
			continue
		}
		
		if !available {
			validationErrors = append(validationErrors, 
				fmt.Sprintf("product %s is not available in quantity %d", 
					item.ProductID().String(), item.Quantity()))
		}
		
		// Проверяем актуальность цены
		currentPrice, err := s.productService.GetProductPrice(ctx, item.ProductID())
		if err != nil {
			s.logger.Error("Failed to get current product price", 
				zap.String("product_id", item.ProductID().String()),
				zap.Error(err))
			validationErrors = append(validationErrors, 
				fmt.Sprintf("failed to get price for product %s", item.ProductID().String()))
			continue
		}
		
		if currentPrice.Amount() != item.UnitPrice().Amount() {
			validationErrors = append(validationErrors, 
				fmt.Sprintf("price for product %s has changed from %d to %d", 
					item.ProductID().String(), item.UnitPrice().Amount(), currentPrice.Amount()))
		}
	}
	
	// 3. Проверяем адрес доставки
	if err := ord.ShippingAddress().Validate(); err != nil {
		validationErrors = append(validationErrors, 
			fmt.Sprintf("invalid shipping address: %s", err.Error()))
	}
	
	// Определяем успешность валидации
	isValid := len(validationErrors) == 0
	
	// Обновляем статус заказа если валидация прошла успешно
	if isValid {
		if err := ord.MarkAsValidated(); err != nil {
			return nil, fmt.Errorf("failed to mark order as validated: %w", err)
		}
		
		// Сохраняем обновленный заказ
		if err := s.orderRepo.Save(ctx, ord); err != nil {
			return nil, fmt.Errorf("failed to save validated order: %w", err)
		}
	}
	
	// Создаем результат
	result := &pipeline.StepResult{
		Step:        "validate_order",
		Success:     isValid,
		StartedAt:   startTime,
		CompletedAt: time.Now(),
		Duration:    time.Since(startTime),
		Output: map[string]interface{}{
			"order_id":          orderID.String(),
			"is_valid":          isValid,
			"validation_errors": validationErrors,
			"order_status":      ord.Status().String(),
			"total_amount":      ord.TotalAmount().Amount(),
		},
		Metadata: map[string]string{
			"validator": "order_validation_service",
			"version":   "1.0",
		},
	}
	
	if !isValid {
		result.Error = errors.New("order validation failed")
		result.Retryable = false // Не повторяем валидацию
	}
	
	s.logger.Info("Order validation completed", 
		zap.String("order_id", orderID.String()),
		zap.Bool("is_valid", isValid),
		zap.Int("error_count", len(validationErrors)),
		zap.Duration("duration", result.Duration))
	
	return result, nil
}

// Name возвращает название шага
func (s *ValidateOrderStep) Name() string {
	return "validate_order"
}

// Validate проверяет корректность входных данных
func (s *ValidateOrderStep) Validate(data *pipeline.StepData) error {
	if _, ok := data.Input["order_id"]; !ok {
		return errors.New("order_id is required")
	}
	return nil
}

// Timeout возвращает максимальное время выполнения
func (s *ValidateOrderStep) Timeout() time.Duration {
	return 30 * time.Second
}

// Dependencies возвращает зависимости шага
func (s *ValidateOrderStep) Dependencies() []string {
	return []string{} // Первый шаг, нет зависимостей
}

// CanRetry определяет, можно ли повторить шаг при ошибке
func (s *ValidateOrderStep) CanRetry(err error) bool {
	// Повторяем только технические ошибки, не ошибки валидации
	return !errors.Is(err, errors.New("order validation failed"))
}

// 🔄 STEP 2: PROCESS PAYMENT

// ProcessPaymentStep шаг обработки платежа
type ProcessPaymentStep struct {
	logger         *zap.Logger
	orderRepo      order.Repository
	paymentService PaymentService
}

// NewProcessPaymentStep создает новый шаг обработки платежа
func NewProcessPaymentStep(logger *zap.Logger, orderRepo order.Repository, paymentService PaymentService) *ProcessPaymentStep {
	return &ProcessPaymentStep{
		logger:         logger,
		orderRepo:      orderRepo,
		paymentService: paymentService,
	}
}

// Execute выполняет обработку платежа
func (s *ProcessPaymentStep) Execute(ctx context.Context, data *pipeline.StepData) (*pipeline.StepResult, error) {
	s.logger.Info("Processing payment", zap.String("execution_id", data.ID))
	
	startTime := time.Now()
	
	// Получаем данные из предыдущего шага
	orderIDStr, ok := data.Input["order_id"].(string)
	if !ok {
		return nil, errors.New("order_id is required")
	}
	
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid order_id: %w", err)
	}
	
	// Проверяем, что предыдущий шаг прошел успешно
	isValid, ok := data.Input["is_valid"].(bool)
	if !ok || !isValid {
		return &pipeline.StepResult{
			Step:        "process_payment",
			Success:     false,
			StartedAt:   startTime,
			CompletedAt: time.Now(),
			Duration:    time.Since(startTime),
			Error:       errors.New("cannot process payment for invalid order"),
			Retryable:   false,
		}, nil
	}
	
	// Получаем заказ
	ord, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	// Получаем способ оплаты из метаданных или используем карту по умолчанию
	paymentMethodStr := data.Metadata["payment_method"]
	if paymentMethodStr == "" {
		paymentMethodStr = "card"
	}
	
	paymentMethod := payment.Method(paymentMethodStr)
	
	// Валидируем способ оплаты
	if err := s.paymentService.ValidatePaymentMethod(paymentMethod); err != nil {
		return nil, fmt.Errorf("invalid payment method: %w", err)
	}
	
	// Обновляем статус заказа
	if err := ord.MarkAsPaymentProcessing(); err != nil {
		return nil, fmt.Errorf("failed to mark order as payment processing: %w", err)
	}
	
	// Сохраняем заказ
	if err := s.orderRepo.Save(ctx, ord); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	
	// Обрабатываем платеж
	pmt, err := s.paymentService.ProcessPayment(ctx, orderID, paymentMethod)
	if err != nil {
		// Возвращаем заказ в предыдущий статус при ошибке
		ord.MarkAsValidated()
		s.orderRepo.Save(ctx, ord)
		
		return &pipeline.StepResult{
			Step:        "process_payment",
			Success:     false,
			StartedAt:   startTime,
			CompletedAt: time.Now(),
			Duration:    time.Since(startTime),
			Output: map[string]interface{}{
				"order_id":       orderID.String(),
				"payment_failed": true,
				"error_message":  err.Error(),
			},
			Error:     fmt.Errorf("payment processing failed: %w", err),
			Retryable: true, // Можно повторить обработку платежа
		}, nil
	}
	
	// Обновляем статус заказа на оплаченный
	if err := ord.MarkAsPaid(); err != nil {
		return nil, fmt.Errorf("failed to mark order as paid: %w", err)
	}
	
	// Сохраняем заказ
	if err := s.orderRepo.Save(ctx, ord); err != nil {
		return nil, fmt.Errorf("failed to save paid order: %w", err)
	}
	
	// Создаем результат
	result := &pipeline.StepResult{
		Step:        "process_payment",
		Success:     true,
		StartedAt:   startTime,
		CompletedAt: time.Now(),
		Duration:    time.Since(startTime),
		Output: map[string]interface{}{
			"order_id":          orderID.String(),
			"payment_id":        pmt.ID().String(),
			"payment_status":    pmt.Status().String(),
			"payment_method":    pmt.Method().String(),
			"payment_amount":    pmt.Amount().Amount(),
			"order_status":      ord.Status().String(),
		},
		Metadata: map[string]string{
			"payment_processor": "payment_service",
			"payment_provider":  pmt.Provider().String(),
		},
	}
	
	s.logger.Info("Payment processing completed", 
		zap.String("order_id", orderID.String()),
		zap.String("payment_id", pmt.ID().String()),
		zap.String("payment_status", pmt.Status().String()),
		zap.Duration("duration", result.Duration))
	
	return result, nil
}

// Name возвращает название шага
func (s *ProcessPaymentStep) Name() string {
	return "process_payment"
}

// Validate проверяет корректность входных данных
func (s *ProcessPaymentStep) Validate(data *pipeline.StepData) error {
	if _, ok := data.Input["order_id"]; !ok {
		return errors.New("order_id is required")
	}
	if isValid, ok := data.Input["is_valid"].(bool); !ok || !isValid {
		return errors.New("order must be valid before payment processing")
	}
	return nil
}

// Timeout возвращает максимальное время выполнения
func (s *ProcessPaymentStep) Timeout() time.Duration {
	return 60 * time.Second // Больше времени для платежа
}

// Dependencies возвращает зависимости шага
func (s *ProcessPaymentStep) Dependencies() []string {
	return []string{"validate_order"}
}

// CanRetry определяет, можно ли повторить шаг при ошибке
func (s *ProcessPaymentStep) CanRetry(err error) bool {
	// Повторяем большинство ошибок платежей
	return true
}

// 🔄 STEP 3: CHECK INVENTORY

// CheckInventoryStep шаг проверки склада
type CheckInventoryStep struct {
	logger           *zap.Logger
	orderRepo        order.Repository
	inventoryService InventoryService
}

// NewCheckInventoryStep создает новый шаг проверки склада
func NewCheckInventoryStep(logger *zap.Logger, orderRepo order.Repository, inventoryService InventoryService) *CheckInventoryStep {
	return &CheckInventoryStep{
		logger:           logger,
		orderRepo:        orderRepo,
		inventoryService: inventoryService,
	}
}

// Execute выполняет проверку склада
func (s *CheckInventoryStep) Execute(ctx context.Context, data *pipeline.StepData) (*pipeline.StepResult, error) {
	s.logger.Info("Checking inventory", zap.String("execution_id", data.ID))
	
	startTime := time.Now()
	
	// Получаем данные
	orderIDStr, ok := data.Input["order_id"].(string)
	if !ok {
		return nil, errors.New("order_id is required")
	}
	
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid order_id: %w", err)
	}
	
	// Проверяем, что платеж прошел успешно
	paymentStatus, ok := data.Input["payment_status"].(string)
	if !ok || paymentStatus != "succeeded" {
		return &pipeline.StepResult{
			Step:        "check_inventory",
			Success:     false,
			StartedAt:   startTime,
			CompletedAt: time.Now(),
			Duration:    time.Since(startTime),
			Error:       errors.New("cannot check inventory for unpaid order"),
			Retryable:   false,
		}, nil
	}
	
	// Получаем заказ
	ord, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	// Проверяем наличие всех товаров
	allAvailable := true
	unavailableItems := make([]string, 0)
	reservationItems := make([]ReservationItem, 0)
	
	for _, item := range ord.Items() {
		available, err := s.inventoryService.CheckAvailability(ctx, item.ProductID(), item.Quantity())
		if err != nil {
			s.logger.Error("Failed to check inventory", 
				zap.String("product_id", item.ProductID().String()),
				zap.Error(err))
			return nil, fmt.Errorf("failed to check inventory for product %s: %w", 
				item.ProductID().String(), err)
		}
		
		if !available {
			allAvailable = false
			unavailableItems = append(unavailableItems, item.ProductID().String())
		} else {
			reservationItems = append(reservationItems, ReservationItem{
				ProductID: item.ProductID(),
				Quantity:  item.Quantity(),
			})
		}
	}
	
	// Создаем результат
	result := &pipeline.StepResult{
		Step:        "check_inventory",
		Success:     allAvailable,
		StartedAt:   startTime,
		CompletedAt: time.Now(),
		Duration:    time.Since(startTime),
		Output: map[string]interface{}{
			"order_id":           orderID.String(),
			"all_available":      allAvailable,
			"unavailable_items":  unavailableItems,
			"available_items":    len(reservationItems),
		},
		Metadata: map[string]string{
			"inventory_service": "warehouse_management",
		},
	}
	
	if allAvailable {
		// Резервируем товары
		reservation, err := s.inventoryService.ReserveItems(ctx, reservationItems)
		if err != nil {
			return nil, fmt.Errorf("failed to reserve items: %w", err)
		}
		
		// Обновляем статус заказа
		if err := ord.MarkAsInventoryChecked(); err != nil {
			// Освобождаем резервирование при ошибке
			s.inventoryService.ReleaseReservation(ctx, reservation.ID)
			return nil, fmt.Errorf("failed to mark order as inventory checked: %w", err)
		}
		
		// Сохраняем заказ
		if err := s.orderRepo.Save(ctx, ord); err != nil {
			// Освобождаем резервирование при ошибке
			s.inventoryService.ReleaseReservation(ctx, reservation.ID)
			return nil, fmt.Errorf("failed to save order: %w", err)
		}
		
		result.Output["reservation_id"] = reservation.ID.String()
		result.Output["reservation_expires_at"] = reservation.ExpiresAt.Format(time.RFC3339)
		result.Output["order_status"] = ord.Status().String()
		
		s.logger.Info("Inventory check completed successfully", 
			zap.String("order_id", orderID.String()),
			zap.String("reservation_id", reservation.ID.String()),
			zap.Duration("duration", result.Duration))
	} else {
		result.Error = errors.New("insufficient inventory")
		result.Retryable = false // Не повторяем если товаров нет
		
		s.logger.Warn("Inventory check failed - insufficient stock", 
			zap.String("order_id", orderID.String()),
			zap.Strings("unavailable_items", unavailableItems),
			zap.Duration("duration", result.Duration))
	}
	
	return result, nil
}

// Name возвращает название шага
func (s *CheckInventoryStep) Name() string {
	return "check_inventory"
}

// Validate проверяет корректность входных данных
func (s *CheckInventoryStep) Validate(data *pipeline.StepData) error {
	if _, ok := data.Input["order_id"]; !ok {
		return errors.New("order_id is required")
	}
	if paymentStatus, ok := data.Input["payment_status"].(string); !ok || paymentStatus != "succeeded" {
		return errors.New("payment must be successful before inventory check")
	}
	return nil
}

// Timeout возвращает максимальное время выполнения
func (s *CheckInventoryStep) Timeout() time.Duration {
	return 45 * time.Second
}

// Dependencies возвращает зависимости шага
func (s *CheckInventoryStep) Dependencies() []string {
	return []string{"process_payment"}
}

// CanRetry определяет, можно ли повторить шаг при ошибке
func (s *CheckInventoryStep) CanRetry(err error) bool {
	// Повторяем только технические ошибки
	return !errors.Is(err, errors.New("insufficient inventory"))
}

// 🔄 STEP 4: SEND NOTIFICATIONS

// SendNotificationsStep шаг отправки уведомлений
type SendNotificationsStep struct {
	logger              *zap.Logger
	notificationService NotificationService
}

// NewSendNotificationsStep создает новый шаг отправки уведомлений
func NewSendNotificationsStep(logger *zap.Logger, notificationService NotificationService) *SendNotificationsStep {
	return &SendNotificationsStep{
		logger:              logger,
		notificationService: notificationService,
	}
}

// Execute выполняет отправку уведомлений
func (s *SendNotificationsStep) Execute(ctx context.Context, data *pipeline.StepData) (*pipeline.StepResult, error) {
	s.logger.Info("Sending notifications", zap.String("execution_id", data.ID))
	
	startTime := time.Now()
	
	// Получаем данные
	orderIDStr, ok := data.Input["order_id"].(string)
	if !ok {
		return nil, errors.New("order_id is required")
	}
	
	// Проверяем, что все предыдущие шаги прошли успешно
	allAvailable, ok := data.Input["all_available"].(bool)
	if !ok || !allAvailable {
		return &pipeline.StepResult{
			Step:        "send_notifications",
			Success:     false,
			StartedAt:   startTime,
			CompletedAt: time.Now(),
			Duration:    time.Since(startTime),
			Error:       errors.New("cannot send success notifications for failed order processing"),
			Retryable:   false,
		}, nil
	}
	
	// Отправляем уведомления (в реальной системе здесь был бы доступ к данным клиента)
	notifications := []string{
		"email_confirmation",
		"sms_confirmation",
	}
	
	successfulNotifications := make([]string, 0)
	failedNotifications := make([]string, 0)
	
	// Email уведомление
	err := s.notificationService.SendEmail(ctx, 
		"customer@example.com", 
		"Заказ подтвержден", 
		fmt.Sprintf("Ваш заказ %s успешно обработан и будет доставлен в ближайшее время.", orderIDStr))
	
	if err != nil {
		s.logger.Error("Failed to send email notification", zap.Error(err))
		failedNotifications = append(failedNotifications, "email")
	} else {
		successfulNotifications = append(successfulNotifications, "email")
	}
	
	// SMS уведомление
	err = s.notificationService.SendSMS(ctx, 
		"+7900123456789", 
		fmt.Sprintf("Заказ %s подтвержден. Ожидайте доставку.", orderIDStr))
	
	if err != nil {
		s.logger.Error("Failed to send SMS notification", zap.Error(err))
		failedNotifications = append(failedNotifications, "sms")
	} else {
		successfulNotifications = append(successfulNotifications, "sms")
	}
	
	// Определяем успешность шага
	success := len(successfulNotifications) > 0
	
	result := &pipeline.StepResult{
		Step:        "send_notifications",
		Success:     success,
		StartedAt:   startTime,
		CompletedAt: time.Now(),
		Duration:    time.Since(startTime),
		Output: map[string]interface{}{
			"order_id":                   orderIDStr,
			"successful_notifications":   successfulNotifications,
			"failed_notifications":       failedNotifications,
			"notification_count":         len(notifications),
			"successful_count":           len(successfulNotifications),
		},
		Metadata: map[string]string{
			"notification_service": "communication_service",
		},
	}
	
	if !success {
		result.Error = errors.New("all notifications failed")
		result.Retryable = true // Можно повторить отправку уведомлений
	}
	
	s.logger.Info("Notifications sending completed", 
		zap.String("order_id", orderIDStr),
		zap.Strings("successful", successfulNotifications),
		zap.Strings("failed", failedNotifications),
		zap.Duration("duration", result.Duration))
	
	return result, nil
}

// Name возвращает название шага
func (s *SendNotificationsStep) Name() string {
	return "send_notifications"
}

// Validate проверяет корректность входных данных
func (s *SendNotificationsStep) Validate(data *pipeline.StepData) error {
	if _, ok := data.Input["order_id"]; !ok {
		return errors.New("order_id is required")
	}
	return nil
}

// Timeout возвращает максимальное время выполнения
func (s *SendNotificationsStep) Timeout() time.Duration {
	return 30 * time.Second
}

// Dependencies возвращает зависимости шага
func (s *SendNotificationsStep) Dependencies() []string {
	return []string{"check_inventory"}
}

// CanRetry определяет, можно ли повторить шаг при ошибке
func (s *SendNotificationsStep) CanRetry(err error) bool {
	// Всегда можно повторить отправку уведомлений
	return true
}