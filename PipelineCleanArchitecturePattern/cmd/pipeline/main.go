package main

import (
	"context"
	"log"
	"time"

	"pipeline-clean-architecture/internal/application/pipeline"
	"pipeline-clean-architecture/internal/domain/order"
	"pipeline-clean-architecture/internal/domain/payment"
	pipelineEngine "pipeline-clean-architecture/pkg/pipeline"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"errors"
)

// 🚀 ДЕМОНСТРАЦИЯ PIPELINE CLEAN ARCHITECTURE
// Показывает работу пайплайна обработки заказов

func main() {
	// Инициализируем логгер
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("🚀 Starting Pipeline Clean Architecture Demo")

	// Создаем моки сервисов для демонстрации
	orderRepo := &MockOrderRepository{}
	productService := &MockProductService{}
	paymentService := &MockPaymentService{}
	inventoryService := &MockInventoryService{}
	notificationService := &MockNotificationService{}

	// Создаем конфигурацию pipeline engine
	config := pipelineEngine.Config{
		MaxConcurrentPipelines: 10,
		DefaultTimeout:         5 * time.Minute,
		EnableMetrics:          true,
		EnableTracing:          true,
		RetryConfig: pipelineEngine.RetryConfig{
			MaxAttempts: 3,
			BaseDelay:   time.Second,
			MaxDelay:    30 * time.Second,
			Multiplier:  2.0,
			Jitter:      true,
		},
		BufferSize: 100,
	}

	// Создаем pipeline engine
	engine := pipelineEngine.NewEngine(logger, config)

	// Регистрируем шаги пайплайна
	validateStep := pipeline.NewValidateOrderStep(logger, orderRepo, productService)
	paymentStep := pipeline.NewProcessPaymentStep(logger, orderRepo, paymentService)
	inventoryStep := pipeline.NewCheckInventoryStep(logger, orderRepo, inventoryService)
	notificationStep := pipeline.NewSendNotificationsStep(logger, notificationService)

	if err := engine.RegisterStep("validate_order", validateStep); err != nil {
		logger.Fatal("Failed to register validate step", zap.Error(err))
	}

	if err := engine.RegisterStep("process_payment", paymentStep); err != nil {
		logger.Fatal("Failed to register payment step", zap.Error(err))
	}

	if err := engine.RegisterStep("check_inventory", inventoryStep); err != nil {
		logger.Fatal("Failed to register inventory step", zap.Error(err))
	}

	if err := engine.RegisterStep("send_notifications", notificationStep); err != nil {
		logger.Fatal("Failed to register notification step", zap.Error(err))
	}

	// Создаем тестовый заказ для демонстрации
	testOrderID := uuid.New()
	testOrder := createTestOrder(testOrderID)
	orderRepo.orders[testOrderID] = testOrder

	logger.Info("📋 Created test order", 
		zap.String("order_id", testOrderID.String()),
		zap.String("customer_id", testOrder.CustomerID().String()),
		zap.Int("items_count", len(testOrder.Items())),
		zap.String("total_amount", testOrder.TotalAmount().String()))

	// Определяем пайплайн обработки заказов
	pipelineDefinition := pipelineEngine.PipelineDefinition{
		Name:        "order_processing",
		Description: "Полный пайплайн обработки e-commerce заказа",
		Steps: []string{
			"validate_order",
			"process_payment", 
			"check_inventory",
			"send_notifications",
		},
		Timeout: 10 * time.Minute,
		Retry: pipelineEngine.RetryConfig{
			MaxAttempts: 3,
			BaseDelay:   2 * time.Second,
			MaxDelay:    30 * time.Second,
			Multiplier:  2.0,
			Jitter:      true,
		},
		Metadata: map[string]string{
			"version":     "1.0",
			"environment": "demo",
			"owner":       "order_processing_team",
		},
	}

	// Подготавливаем входные данные для пайплайна
	initialData := map[string]interface{}{
		"order_id": testOrderID.String(),
	}

	logger.Info("🔄 Starting pipeline execution",
		zap.String("pipeline", pipelineDefinition.Name),
		zap.Strings("steps", pipelineDefinition.Steps))

	// Выполняем пайплайн
	ctx := context.Background()
	executionContext, err := engine.Execute(ctx, pipelineDefinition, initialData)

	if err != nil {
		logger.Error("❌ Pipeline execution failed", 
			zap.Error(err),
			zap.String("status", string(executionContext.Status)))
		
		// Показываем результаты даже при ошибке
		displayResults(logger, executionContext)
		return
	}

	logger.Info("✅ Pipeline execution completed successfully!",
		zap.String("status", string(executionContext.Status)),
		zap.Duration("total_duration", time.Since(executionContext.StartedAt)),
		zap.Int("completed_steps", len(executionContext.CompletedSteps)))

	// Показываем подробные результаты
	displayResults(logger, executionContext)

	// Показываем финальное состояние заказа
	finalOrder, _ := orderRepo.GetByID(ctx, testOrderID)
	logger.Info("📦 Final order status",
		zap.String("order_id", testOrderID.String()),
		zap.String("status", finalOrder.Status().String()),
		zap.String("total_amount", finalOrder.TotalAmount().String()))

	logger.Info("🎉 Demo completed successfully!")
}

// displayResults показывает подробные результаты выполнения пайплайна
func displayResults(logger *zap.Logger, execCtx *pipelineEngine.ExecutionContext) {
	logger.Info("📊 Pipeline Execution Results",
		zap.String("execution_id", execCtx.ID),
		zap.String("pipeline", execCtx.Pipeline.Name),
		zap.String("status", string(execCtx.Status)),
		zap.Time("started_at", execCtx.StartedAt))

	if execCtx.CompletedAt != nil {
		logger.Info("⏱️ Timing Information",
			zap.Duration("total_duration", execCtx.CompletedAt.Sub(execCtx.StartedAt)),
			zap.Time("completed_at", *execCtx.CompletedAt))
	}

	logger.Info("🔄 Completed Steps",
		zap.Strings("steps", execCtx.CompletedSteps),
		zap.Int("count", len(execCtx.CompletedSteps)))

	// Показываем результаты каждого шага
	for stepName, result := range execCtx.Results {
		stepLogger := logger.With(zap.String("step", stepName))
		
		if result.Success {
			stepLogger.Info("✅ Step completed successfully",
				zap.Duration("duration", result.Duration),
				zap.Any("output", result.Output))
		} else {
			stepLogger.Error("❌ Step failed",
				zap.Duration("duration", result.Duration),
				zap.Error(result.Error),
				zap.Bool("retryable", result.Retryable))
		}

		// Показываем метаданные шага
		if len(result.Metadata) > 0 {
			stepLogger.Info("📋 Step metadata", zap.Any("metadata", result.Metadata))
		}
	}

	// Показываем глобальные данные
	if len(execCtx.GlobalData) > 0 {
		logger.Info("🌐 Global Data", zap.Any("data", execCtx.GlobalData))
	}

	// Показываем ошибку пайплайна если есть
	if execCtx.Error != nil {
		logger.Error("💥 Pipeline Error", zap.Error(execCtx.Error))
	}
}

// createTestOrder создает тестовый заказ для демонстрации
func createTestOrder(orderID uuid.UUID) *order.Order {
	customerID := uuid.New()
	
	// Создаем товары заказа
	items := []order.OrderItem{
		createOrderItem(uuid.New(), 2, 150000), // Товар 1: 2 шт по 1500 руб
		createOrderItem(uuid.New(), 1, 75000),  // Товар 2: 1 шт по 750 руб
	}

	// Создаем адрес доставки
	shippingAddress := order.Address{
		Street:     "ул. Примерная, д. 123, кв. 45",
		City:       "Москва", 
		PostalCode: "123456",
		Country:    "Россия",
		Phone:      "+7 (900) 123-45-67",
	}

	// Создаем заказ
	testOrder, err := order.NewOrder(customerID, items, shippingAddress)
	if err != nil {
		log.Fatalf("Failed to create test order: %v", err)
	}

	// Устанавливаем дополнительные свойства
	testOrder.SetPriority(order.PriorityNormal)
	testOrder.SetNotes("Тестовый заказ для демонстрации пайплайна")

	return testOrder
}

// createOrderItem создает товар заказа (хелпер функция)
func createOrderItem(productID uuid.UUID, quantity int, priceKopecks int64) order.OrderItem {
	// В реальном приложении это было бы через фабричный метод
	// Здесь используем reflection для доступа к приватным полям (только для демо!)
	money := order.Money{Amount: priceKopecks, Currency: "RUB"}
	
	return order.OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: money,
		Discount:  order.Money{Amount: 0, Currency: "RUB"},
	}
}

// 🔧 MOCK РЕАЛИЗАЦИИ СЕРВИСОВ ДЛЯ ДЕМОНСТРАЦИИ

// MockOrderRepository мок репозитория заказов
type MockOrderRepository struct {
	orders map[uuid.UUID]*order.Order
}

func (r *MockOrderRepository) Save(ctx context.Context, ord *order.Order) error {
	if r.orders == nil {
		r.orders = make(map[uuid.UUID]*order.Order)
	}
	r.orders[ord.ID()] = ord
	return nil
}

func (r *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	if r.orders == nil {
		r.orders = make(map[uuid.UUID]*order.Order)
	}
	
	ord, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return ord, nil
}

// Остальные методы интерфейса (заглушки для демо)
func (r *MockOrderRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*order.Order, error) {
	return nil, errors.New("not implemented")
}

func (r *MockOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented") 
}

func (r *MockOrderRepository) GetByStatus(ctx context.Context, status order.Status) ([]*order.Order, error) {
	return nil, errors.New("not implemented")
}

func (r *MockOrderRepository) GetByStatusAndPriority(ctx context.Context, status order.Status, priority order.Priority) ([]*order.Order, error) {
	return nil, errors.New("not implemented")
}

func (r *MockOrderRepository) GetPendingOrders(ctx context.Context) ([]*order.Order, error) {
	return nil, errors.New("not implemented")
}

func (r *MockOrderRepository) GetActiveOrders(ctx context.Context) ([]*order.Order, error) {
	return nil, errors.New("not implemented")
}

func (r *MockOrderRepository) GetOrdersInDateRange(ctx context.Context, from, to time.Time) ([]*order.Order, error) {
	return nil, errors.New("not implemented")
}

func (r *MockOrderRepository) CountByStatus(ctx context.Context) (map[order.Status]int, error) {
	return nil, errors.New("not implemented")
}

func (r *MockOrderRepository) GetTotalRevenueByPeriod(ctx context.Context, from, to time.Time) (order.Money, error) {
	return order.Money{}, errors.New("not implemented")
}

func (r *MockOrderRepository) SaveBatch(ctx context.Context, orders []*order.Order) error {
	return errors.New("not implemented")
}

func (r *MockOrderRepository) UpdateStatusBatch(ctx context.Context, orderIDs []uuid.UUID, newStatus order.Status) error {
	return errors.New("not implemented")
}

// MockProductService мок сервиса товаров
type MockProductService struct{}

func (s *MockProductService) CheckProductAvailability(ctx context.Context, productID uuid.UUID, quantity int) (bool, error) {
	// Для демо всегда возвращаем true
	return true, nil
}

func (s *MockProductService) GetProductPrice(ctx context.Context, productID uuid.UUID) (order.Money, error) {
	// Возвращаем фиксированную цену для демо
	return order.Money{Amount: 150000, Currency: "RUB"}, nil
}

// MockPaymentService мок сервиса платежей
type MockPaymentService struct{}

func (s *MockPaymentService) ProcessPayment(ctx context.Context, orderID uuid.UUID, method payment.Method) (*payment.Payment, error) {
	// Создаем успешный платеж для демо
	pmt, err := payment.NewPayment(orderID, uuid.New(), 
		payment.Money{Amount: 375000, Currency: "RUB"}, 
		method, 
		payment.ProviderSberbank)
	
	if err != nil {
		return nil, err
	}

	// Имитируем успешную обработку
	pmt.MarkAsProcessing()
	pmt.MarkAsSucceeded("payment_ext_123456")
	
	return pmt, nil
}

func (s *MockPaymentService) ValidatePaymentMethod(method payment.Method) error {
	// Для демо принимаем все методы
	return nil
}

// MockInventoryService мок сервиса склада
type MockInventoryService struct{}

func (s *MockInventoryService) CheckAvailability(ctx context.Context, productID uuid.UUID, quantity int) (bool, error) {
	// Для демо всегда есть в наличии
	return true, nil
}

func (s *MockInventoryService) ReserveItems(ctx context.Context, items []pipeline.ReservationItem) (*pipeline.Reservation, error) {
	// Создаем успешное резервирование
	return &pipeline.Reservation{
		ID:        uuid.New(),
		Items:     items,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *MockInventoryService) ReleaseReservation(ctx context.Context, reservationID uuid.UUID) error {
	// Для демо всегда успешно
	return nil
}

// MockNotificationService мок сервиса уведомлений
type MockNotificationService struct{}

func (s *MockNotificationService) SendEmail(ctx context.Context, to, subject, body string) error {
	// Имитируем отправку email
	log.Printf("📧 EMAIL SENT to %s: %s", to, subject)
	return nil
}

func (s *MockNotificationService) SendSMS(ctx context.Context, phone, message string) error {
	// Имитируем отправку SMS
	log.Printf("📱 SMS SENT to %s: %s", phone, message)
	return nil
}

func (s *MockNotificationService) SendPushNotification(ctx context.Context, userID uuid.UUID, title, body string) error {
	// Имитируем push уведомление
	log.Printf("🔔 PUSH SENT to %s: %s", userID.String(), title)
	return nil
}