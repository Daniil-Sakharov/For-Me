package main

import (
	"context"
	"log"
	"time"

	// ИМПОРТЫ ПО СЛОЯМ CLEAN ARCHITECTURE:
	
	// APPLICATION LAYER - шаги пайплайна, которые координируют бизнес-процессы
	"pipeline-clean-architecture/internal/application/pipeline"
	
	// DOMAIN LAYER - доменные сущности с бизнес-логикой (НЕ зависят от внешних систем)
	"pipeline-clean-architecture/internal/domain/order"
	"pipeline-clean-architecture/internal/domain/payment"
	
	// PKG LAYER - переиспользуемые компоненты (Pipeline Engine)
	pipelineEngine "pipeline-clean-architecture/pkg/pipeline"

	// ВНЕШНИЕ ЗАВИСИМОСТИ
	"github.com/google/uuid"
	"go.uber.org/zap"
	"errors"
)

// 🚀 ДЕМОНСТРАЦИЯ PIPELINE CLEAN ARCHITECTURE
// 
// Этот файл демонстрирует полный цикл работы пайплайна обработки заказов
// в архитектуре Clean Architecture + Pipeline Pattern
//
// ПОТОК ВЫПОЛНЕНИЯ:
// 1. Инициализация всех зависимостей (логгер, сервисы-моки)
// 2. Создание и настройка Pipeline Engine
// 3. Регистрация шагов пайплайна 
// 4. Создание тестового заказа
// 5. Запуск пайплайна обработки
// 6. Демонстрация результатов
//
// ПРИНЦИПЫ ДЕМОНСТРАЦИИ:
// - Показать как данные проходят через все слои архитектуры
// - Продемонстрировать инверсию зависимостей (DI)
// - Показать как каждый шаг пайплайна использует доменные сущности
// - Демонстрировать обработку ошибок и логирование

func main() {
	// ================================================
	// ЭТАП 1: ИНИЦИАЛИЗАЦИЯ ИНФРАСТРУКТУРЫ
	// ================================================
	
	// Инициализируем структурированный логгер (zap)
	// В production здесь была бы конфигурация из файла/env переменных
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("🚀 Starting Pipeline Clean Architecture Demo")

	// ================================================
	// ЭТАП 2: DEPENDENCY INJECTION - СОЗДАНИЕ ЗАВИСИМОСТЕЙ
	// ================================================
	//
	// В реальной системе здесь был бы DI контейнер, который создавал бы:
	// - Подключения к PostgreSQL/MongoDB
	// - HTTP клиенты для внешних API
	// - Redis для кеширования
	// - Kafka/RabbitMQ для очередей
	//
	// Для демонстрации создаем моки, которые имитируют поведение реальных сервисов
	
	// INFRASTRUCTURE LAYER - моки репозиториев и внешних сервисов
	// В реальности это были бы PostgreSQLOrderRepository, StripePaymentService и т.д.
	orderRepo := &MockOrderRepository{}                // Имитация PostgreSQL репозитория
	productService := &MockProductService{}            // Имитация сервиса каталога товаров
	paymentService := &MockPaymentService{}            // Имитация Stripe/PayPal API
	inventoryService := &MockInventoryService{}        // Имитация системы управления складом
	notificationService := &MockNotificationService{}  // Имитация email/SMS сервисов

	// ================================================
	// ЭТАП 3: КОНФИГУРАЦИЯ PIPELINE ENGINE
	// ================================================
	//
	// Pipeline Engine - это ядро нашей системы, которое управляет выполнением шагов
	// Конфигурация определяет поведение движка:
	
	config := pipelineEngine.Config{
		MaxConcurrentPipelines: 10,               // Сколько пайплайнов можно выполнять параллельно
		DefaultTimeout:         5 * time.Minute,  // Максимальное время выполнения пайплайна
		EnableMetrics:          true,             // Включить сбор метрик (Prometheus)
		EnableTracing:          true,             // Включить трассировку для отладки
		RetryConfig: pipelineEngine.RetryConfig{  // Настройки повторов при ошибках
			MaxAttempts: 3,                       // Максимум 3 попытки
			BaseDelay:   time.Second,             // Начальная задержка 1 сек
			MaxDelay:    30 * time.Second,        // Максимальная задержка 30 сек
			Multiplier:  2.0,                     // Экспоненциальный backoff (1s, 2s, 4s...)
			Jitter:      true,                    // Добавить случайное отклонение для избежания thundering herd
		},
		BufferSize: 100,                          // Размер буферов каналов
	}

	// Создаем экземпляр Pipeline Engine
	// Это центральный компонент, который будет координировать выполнение всех шагов
	engine := pipelineEngine.NewEngine(logger, config)

	// ================================================
	// ЭТАП 4: СОЗДАНИЕ И РЕГИСТРАЦИЯ ШАГОВ ПАЙПЛАЙНА
	// ================================================
	//
	// Каждый шаг пайплайна инкапсулирует определенную бизнес-логику.
	// Шаги реализуют интерфейс StepHandler и могут выполняться независимо.
	//
	// ПРИНЦИП DEPENDENCY INJECTION:
	// Каждый шаг получает только те зависимости, которые ему нужны.
	// Это делает шаги тестируемыми и переиспользуемыми.
	
	// Шаг 1: Валидация заказа
	// Зависимости: логгер + репозиторий заказов + сервис каталога товаров
	// Что делает: проверяет корректность заказа, наличие товаров, актуальность цен
	validateStep := pipeline.NewValidateOrderStep(logger, orderRepo, productService)
	
	// Шаг 2: Обработка платежа  
	// Зависимости: логгер + репозиторий заказов + платежный сервис
	// Что делает: обрабатывает платеж через внешний сервис (Stripe/PayPal)
	paymentStep := pipeline.NewProcessPaymentStep(logger, orderRepo, paymentService)
	
	// Шаг 3: Проверка склада
	// Зависимости: логгер + репозиторий заказов + сервис управления складом
	// Что делает: проверяет наличие товаров на складе и резервирует их
	inventoryStep := pipeline.NewCheckInventoryStep(logger, orderRepo, inventoryService)
	
	// Шаг 4: Отправка уведомлений
	// Зависимости: логгер + сервис уведомлений
	// Что делает: отправляет email/SMS клиенту о статусе заказа
	notificationStep := pipeline.NewSendNotificationsStep(logger, notificationService)

	// Регистрируем все шаги в Pipeline Engine
	// Движок будет использовать эти имена для поиска и выполнения шагов
	
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