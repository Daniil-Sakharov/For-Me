package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// 🔄 PIPELINE ENGINE - ЯДРО СИСТЕМЫ ОБРАБОТКИ ДАННЫХ
//
// ============================================================================
// ЧТО ТАКОЕ PIPELINE PATTERN:
// ============================================================================
//
// Pipeline Pattern - это архитектурный паттерн, где данные проходят через
// последовательность этапов обработки. Каждый этап выполняет определенную
// трансформацию данных и передает результат следующему этапу.
//
// АНАЛОГИЯ: Конвейер на заводе
// - Деталь поступает на конвейер
// - Каждый рабочий выполняет свою операцию
// - Деталь передается дальше по конвейеру
// - В конце получаем готовый продукт
//
// ============================================================================
// ПРЕИМУЩЕСТВА PIPELINE:
// ============================================================================
//
// 1. 🔧 МОДУЛЬНОСТЬ:
//    - Каждый шаг независим и решает одну задачу
//    - Легко добавлять/удалять/переставлять шаги
//    - Переиспользование шагов в разных пайплайнах
//
// 2. 🧪 ТЕСТИРУЕМОСТЬ:
//    - Каждый шаг тестируется изолированно
//    - Легко мокать зависимости
//    - Быстрые и надежные тесты
//
// 3. 🔄 ОТКАЗОУСТОЙЧИВОСТЬ:
//    - Retry логика для каждого шага
//    - Graceful обработка ошибок
//    - Возможность продолжить с определенного шага
//
// 4. 📊 НАБЛЮДАЕМОСТЬ:
//    - Логирование каждого шага
//    - Метрики времени выполнения
//    - Трассировка выполнения
//
// ============================================================================
// КОМПОНЕНТЫ PIPELINE ENGINE:
// ============================================================================
//
// - Engine: Основной движок, управляет выполнением
// - StepHandler: Интерфейс для реализации шагов
// - Middleware: Промежуточная обработка (логирование, метрики)
// - ErrorHandler: Обработка ошибок и решения о retry
// - MetricsCollector: Сбор метрик производительности
//
// ============================================================================

// Engine представляет движок пайплайнов
//
// ОТВЕТСТВЕННОСТИ:
// - Регистрация и управление шагами
// - Выполнение пайплайнов в правильной последовательности  
// - Обработка ошибок и retry логика
// - Сбор метрик и логирование
// - Управление таймаутами и отменой операций
type Engine struct {
	logger        *zap.Logger
	stepRegistry  map[string]StepHandler
	middleware    []Middleware
	errorHandler  ErrorHandler
	metrics       MetricsCollector
	config        Config
	mu            sync.RWMutex
}

// StepHandler интерфейс для обработчиков шагов пайплайна
type StepHandler interface {
	// Execute выполняет шаг пайплайна
	Execute(ctx context.Context, data *StepData) (*StepResult, error)
	
	// Name возвращает название шага
	Name() string
	
	// Validate проверяет корректность входных данных
	Validate(data *StepData) error
	
	// Timeout возвращает максимальное время выполнения
	Timeout() time.Duration
	
	// Dependencies возвращает зависимости шага
	Dependencies() []string
	
	// CanRetry определяет, можно ли повторить шаг при ошибке
	CanRetry(err error) bool
}

// Middleware интерфейс для промежуточного ПО
type Middleware interface {
	// Before выполняется перед шагом
	Before(ctx context.Context, step string, data *StepData) error
	
	// After выполняется после шага
	After(ctx context.Context, step string, result *StepResult) error
	
	// OnError выполняется при ошибке
	OnError(ctx context.Context, step string, err error) error
}

// ErrorHandler интерфейс для обработки ошибок
type ErrorHandler interface {
	// Handle обрабатывает ошибку
	Handle(ctx context.Context, step string, err error) error
	
	// ShouldRetry определяет, нужно ли повторить шаг
	ShouldRetry(err error) bool
	
	// ShouldContinue определяет, нужно ли продолжить пайплайн
	ShouldContinue(err error) bool
}

// MetricsCollector интерфейс для сбора метрик
type MetricsCollector interface {
	// RecordStepDuration записывает время выполнения шага
	RecordStepDuration(step string, duration time.Duration)
	
	// RecordStepSuccess записывает успешное выполнение шага
	RecordStepSuccess(step string)
	
	// RecordStepError записывает ошибку шага
	RecordStepError(step string, err error)
	
	// RecordPipelineDuration записывает время выполнения пайплайна
	RecordPipelineDuration(pipeline string, duration time.Duration)
}

// 📊 СТРУКТУРЫ ДАННЫХ

// StepData данные для выполнения шага
type StepData struct {
	// ID уникальный идентификатор выполнения
	ID string `json:"id"`
	
	// Step название шага
	Step string `json:"step"`
	
	// Input входные данные
	Input map[string]interface{} `json:"input"`
	
	// Context контекстные данные
	Context map[string]interface{} `json:"context"`
	
	// Metadata метаданные
	Metadata map[string]string `json:"metadata"`
	
	// StartedAt время начала
	StartedAt time.Time `json:"started_at"`
}

// StepResult результат выполнения шага
type StepResult struct {
	// Step название шага
	Step string `json:"step"`
	
	// Success успешно ли выполнен шаг
	Success bool `json:"success"`
	
	// Output выходные данные
	Output map[string]interface{} `json:"output"`
	
	// Duration время выполнения
	Duration time.Duration `json:"duration"`
	
	// StartedAt время начала
	StartedAt time.Time `json:"started_at"`
	
	// CompletedAt время завершения
	CompletedAt time.Time `json:"completed_at"`
	
	// Error ошибка (если есть)
	Error error `json:"error,omitempty"`
	
	// Retryable можно ли повторить при ошибке
	Retryable bool `json:"retryable"`
	
	// Metadata метаданные результата
	Metadata map[string]string `json:"metadata"`
}

// PipelineDefinition определение пайплайна
type PipelineDefinition struct {
	// Name название пайплайна
	Name string `json:"name"`
	
	// Description описание
	Description string `json:"description"`
	
	// Steps шаги пайплайна в порядке выполнения
	Steps []string `json:"steps"`
	
	// Parallel параллельные шаги (группы)
	Parallel [][]string `json:"parallel,omitempty"`
	
	// Conditional условные шаги
	Conditional map[string]Condition `json:"conditional,omitempty"`
	
	// Timeout общий таймаут пайплайна
	Timeout time.Duration `json:"timeout"`
	
	// Retry настройки повторов
	Retry RetryConfig `json:"retry"`
	
	// Metadata метаданные пайплайна
	Metadata map[string]string `json:"metadata"`
}

// Condition условие для выполнения шага
type Condition struct {
	// Field поле для проверки
	Field string `json:"field"`
	
	// Operator оператор сравнения
	Operator string `json:"operator"` // eq, ne, gt, lt, contains, exists
	
	// Value значение для сравнения
	Value interface{} `json:"value"`
	
	// Logic логический оператор для объединения условий
	Logic string `json:"logic,omitempty"` // and, or
	
	// Nested вложенные условия
	Nested []Condition `json:"nested,omitempty"`
}

// RetryConfig настройки повторов
type RetryConfig struct {
	// MaxAttempts максимальное количество попыток
	MaxAttempts int `json:"max_attempts"`
	
	// BaseDelay базовая задержка между попытками
	BaseDelay time.Duration `json:"base_delay"`
	
	// MaxDelay максимальная задержка
	MaxDelay time.Duration `json:"max_delay"`
	
	// Multiplier множитель для экспоненциального backoff
	Multiplier float64 `json:"multiplier"`
	
	// Jitter добавлять ли случайное отклонение
	Jitter bool `json:"jitter"`
}

// ExecutionContext контекст выполнения пайплайна
type ExecutionContext struct {
	// ID идентификатор выполнения
	ID string `json:"id"`
	
	// Pipeline определение пайплайна
	Pipeline PipelineDefinition `json:"pipeline"`
	
	// StartedAt время начала
	StartedAt time.Time `json:"started_at"`
	
	// CompletedAt время завершения
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	// CurrentStep текущий шаг
	CurrentStep string `json:"current_step"`
	
	// CompletedSteps завершенные шаги
	CompletedSteps []string `json:"completed_steps"`
	
	// Results результаты шагов
	Results map[string]*StepResult `json:"results"`
	
	// GlobalData глобальные данные пайплайна
	GlobalData map[string]interface{} `json:"global_data"`
	
	// Status статус выполнения
	Status ExecutionStatus `json:"status"`
	
	// Error общая ошибка пайплайна
	Error error `json:"error,omitempty"`
}

// ExecutionStatus статус выполнения
type ExecutionStatus string

const (
	StatusPending    ExecutionStatus = "pending"
	StatusRunning    ExecutionStatus = "running"
	StatusCompleted  ExecutionStatus = "completed"
	StatusFailed     ExecutionStatus = "failed"
	StatusCancelled  ExecutionStatus = "cancelled"
)

// Config конфигурация движка
type Config struct {
	// MaxConcurrentPipelines максимальное количество одновременных пайплайнов
	MaxConcurrentPipelines int `json:"max_concurrent_pipelines"`
	
	// DefaultTimeout таймаут по умолчанию
	DefaultTimeout time.Duration `json:"default_timeout"`
	
	// EnableMetrics включить сбор метрик
	EnableMetrics bool `json:"enable_metrics"`
	
	// EnableTracing включить трассировку
	EnableTracing bool `json:"enable_tracing"`
	
	// RetryConfig настройки повторов по умолчанию
	RetryConfig RetryConfig `json:"retry_config"`
	
	// BufferSize размер буфера для каналов
	BufferSize int `json:"buffer_size"`
}

// 🏗️ КОНСТРУКТОР

// NewEngine создает новый движок пайплайнов
func NewEngine(logger *zap.Logger, config Config) *Engine {
	return &Engine{
		logger:       logger,
		stepRegistry: make(map[string]StepHandler),
		middleware:   make([]Middleware, 0),
		config:       config,
	}
}

// 📝 РЕГИСТРАЦИЯ КОМПОНЕНТОВ

// RegisterStep регистрирует обработчик шага
func (e *Engine) RegisterStep(name string, handler StepHandler) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if _, exists := e.stepRegistry[name]; exists {
		return fmt.Errorf("step %s already registered", name)
	}
	
	e.stepRegistry[name] = handler
	e.logger.Info("Registered pipeline step", zap.String("step", name))
	return nil
}

// RegisterMiddleware регистрирует промежуточное ПО
func (e *Engine) RegisterMiddleware(middleware Middleware) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.middleware = append(e.middleware, middleware)
	e.logger.Info("Registered pipeline middleware")
}

// SetErrorHandler устанавливает обработчик ошибок
func (e *Engine) SetErrorHandler(handler ErrorHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.errorHandler = handler
	e.logger.Info("Set error handler")
}

// SetMetricsCollector устанавливает сборщик метрик
func (e *Engine) SetMetricsCollector(collector MetricsCollector) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	e.metrics = collector
	e.logger.Info("Set metrics collector")
}

// 🚀 ВЫПОЛНЕНИЕ ПАЙПЛАЙНОВ

// Execute выполняет пайплайн
func (e *Engine) Execute(ctx context.Context, definition PipelineDefinition, initialData map[string]interface{}) (*ExecutionContext, error) {
	e.logger.Info("Starting pipeline execution", 
		zap.String("pipeline", definition.Name))
	
	// Создаем контекст выполнения
	execCtx := &ExecutionContext{
		ID:             generateID(),
		Pipeline:       definition,
		StartedAt:      time.Now(),
		CurrentStep:    "",
		CompletedSteps: make([]string, 0),
		Results:        make(map[string]*StepResult),
		GlobalData:     initialData,
		Status:         StatusRunning,
	}
	
	// Создаем контекст с таймаутом
	timeout := definition.Timeout
	if timeout == 0 {
		timeout = e.config.DefaultTimeout
	}
	
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	
	// Выполняем пайплайн
	startTime := time.Now()
	err := e.executePipeline(ctxWithTimeout, execCtx)
	duration := time.Since(startTime)
	
	// Завершаем контекст
	now := time.Now()
	execCtx.CompletedAt = &now
	
	if err != nil {
		execCtx.Status = StatusFailed
		execCtx.Error = err
		e.logger.Error("Pipeline execution failed",
			zap.String("pipeline", definition.Name),
			zap.Error(err),
			zap.Duration("duration", duration))
	} else {
		execCtx.Status = StatusCompleted
		e.logger.Info("Pipeline execution completed",
			zap.String("pipeline", definition.Name),
			zap.Duration("duration", duration))
	}
	
	// Записываем метрики
	if e.metrics != nil {
		e.metrics.RecordPipelineDuration(definition.Name, duration)
	}
	
	return execCtx, err
}

// executePipeline выполняет шаги пайплайна
func (e *Engine) executePipeline(ctx context.Context, execCtx *ExecutionContext) error {
	// Выполняем шаги последовательно
	for _, step := range execCtx.Pipeline.Steps {
		// Проверяем контекст
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		// Проверяем условие выполнения шага
		if !e.shouldExecuteStep(step, execCtx) {
			e.logger.Debug("Skipping step due to condition",
				zap.String("step", step))
			continue
		}
		
		// Выполняем шаг
		execCtx.CurrentStep = step
		result, err := e.executeStep(ctx, step, execCtx)
		
		// Сохраняем результат
		execCtx.Results[step] = result
		
		if err != nil {
			// Обрабатываем ошибку
			if e.errorHandler != nil {
				if handleErr := e.errorHandler.Handle(ctx, step, err); handleErr != nil {
					e.logger.Error("Error handler failed", zap.Error(handleErr))
				}
				
				// Проверяем, нужно ли продолжить
				if !e.errorHandler.ShouldContinue(err) {
					return fmt.Errorf("pipeline stopped at step %s: %w", step, err)
				}
			} else {
				return fmt.Errorf("step %s failed: %w", step, err)
			}
		}
		
		// Добавляем в завершенные шаги
		execCtx.CompletedSteps = append(execCtx.CompletedSteps, step)
		
		// Обновляем глобальные данные из результата
		e.updateGlobalData(execCtx, result)
	}
	
	return nil
}

// executeStep выполняет отдельный шаг
func (e *Engine) executeStep(ctx context.Context, stepName string, execCtx *ExecutionContext) (*StepResult, error) {
	e.logger.Debug("Executing step", zap.String("step", stepName))
	
	// Получаем обработчик шага
	e.mu.RLock()
	handler, exists := e.stepRegistry[stepName]
	e.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("step handler not found: %s", stepName)
	}
	
	// Подготавливаем данные для шага
	stepData := &StepData{
		ID:        generateID(),
		Step:      stepName,
		Input:     e.prepareStepInput(stepName, execCtx),
		Context:   execCtx.GlobalData,
		Metadata:  execCtx.Pipeline.Metadata,
		StartedAt: time.Now(),
	}
	
	// Валидируем входные данные
	if err := handler.Validate(stepData); err != nil {
		return &StepResult{
			Step:        stepName,
			Success:     false,
			StartedAt:   stepData.StartedAt,
			CompletedAt: time.Now(),
			Error:       fmt.Errorf("validation failed: %w", err),
		}, err
	}
	
	// Выполняем middleware Before
	for _, mw := range e.middleware {
		if err := mw.Before(ctx, stepName, stepData); err != nil {
			return nil, fmt.Errorf("middleware before failed: %w", err)
		}
	}
	
	// Создаем контекст с таймаутом для шага
	stepTimeout := handler.Timeout()
	if stepTimeout == 0 {
		stepTimeout = e.config.DefaultTimeout
	}
	
	ctxWithTimeout, cancel := context.WithTimeout(ctx, stepTimeout)
	defer cancel()
	
	// Выполняем шаг с повторами
	var result *StepResult
	var err error
	
	retryConfig := execCtx.Pipeline.Retry
	if retryConfig.MaxAttempts == 0 {
		retryConfig = e.config.RetryConfig
	}
	
	for attempt := 1; attempt <= retryConfig.MaxAttempts; attempt++ {
		startTime := time.Now()
		result, err = handler.Execute(ctxWithTimeout, stepData)
		duration := time.Since(startTime)
		
		if result == nil {
			result = &StepResult{
				Step:        stepName,
				Success:     err == nil,
				StartedAt:   startTime,
				CompletedAt: time.Now(),
				Duration:    duration,
				Error:       err,
			}
		}
		
		// Записываем метрики
		if e.metrics != nil {
			e.metrics.RecordStepDuration(stepName, duration)
			if err == nil {
				e.metrics.RecordStepSuccess(stepName)
			} else {
				e.metrics.RecordStepError(stepName, err)
			}
		}
		
		// Если успешно или нельзя повторить - выходим
		if err == nil || !handler.CanRetry(err) || attempt == retryConfig.MaxAttempts {
			break
		}
		
		// Ждем перед повтором
		delay := e.calculateRetryDelay(attempt, retryConfig)
		e.logger.Warn("Retrying step",
			zap.String("step", stepName),
			zap.Int("attempt", attempt),
			zap.Duration("delay", delay),
			zap.Error(err))
		
		timer := time.NewTimer(delay)
		select {
		case <-ctxWithTimeout.Done():
			timer.Stop()
			return result, ctxWithTimeout.Err()
		case <-timer.C:
		}
	}
	
	// Выполняем middleware After
	for _, mw := range e.middleware {
		if mwErr := mw.After(ctx, stepName, result); mwErr != nil {
			e.logger.Error("Middleware after failed", zap.Error(mwErr))
		}
	}
	
	// Выполняем middleware OnError если была ошибка
	if err != nil {
		for _, mw := range e.middleware {
			if mwErr := mw.OnError(ctx, stepName, err); mwErr != nil {
				e.logger.Error("Middleware onError failed", zap.Error(mwErr))
			}
		}
	}
	
	return result, err
}

// 🔧 ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ

// shouldExecuteStep проверяет, нужно ли выполнить шаг
func (e *Engine) shouldExecuteStep(step string, execCtx *ExecutionContext) bool {
	condition, exists := execCtx.Pipeline.Conditional[step]
	if !exists {
		return true
	}
	
	return e.evaluateCondition(condition, execCtx)
}

// evaluateCondition оценивает условие
func (e *Engine) evaluateCondition(condition Condition, execCtx *ExecutionContext) bool {
	// Получаем значение поля
	fieldValue := e.getFieldValue(condition.Field, execCtx)
	
	// Сравниваем с условием
	result := e.compareValues(fieldValue, condition.Operator, condition.Value)
	
	// Обрабатываем вложенные условия
	if len(condition.Nested) > 0 {
		nestedResults := make([]bool, len(condition.Nested))
		for i, nested := range condition.Nested {
			nestedResults[i] = e.evaluateCondition(nested, execCtx)
		}
		
		// Применяем логический оператор
		switch condition.Logic {
		case "and":
			for _, nr := range nestedResults {
				result = result && nr
			}
		case "or":
			for _, nr := range nestedResults {
				result = result || nr
			}
		}
	}
	
	return result
}

// getFieldValue получает значение поля из контекста
func (e *Engine) getFieldValue(field string, execCtx *ExecutionContext) interface{} {
	// Поиск в глобальных данных
	if val, exists := execCtx.GlobalData[field]; exists {
		return val
	}
	
	// Поиск в результатах шагов
	for _, result := range execCtx.Results {
		if val, exists := result.Output[field]; exists {
			return val
		}
	}
	
	return nil
}

// compareValues сравнивает значения
func (e *Engine) compareValues(fieldValue interface{}, operator string, conditionValue interface{}) bool {
	switch operator {
	case "eq":
		return fieldValue == conditionValue
	case "ne":
		return fieldValue != conditionValue
	case "exists":
		return fieldValue != nil
	// Дополнительные операторы можно добавить здесь
	default:
		return false
	}
}

// prepareStepInput подготавливает входные данные для шага
func (e *Engine) prepareStepInput(stepName string, execCtx *ExecutionContext) map[string]interface{} {
	input := make(map[string]interface{})
	
	// Копируем глобальные данные
	for k, v := range execCtx.GlobalData {
		input[k] = v
	}
	
	// Добавляем результаты предыдущих шагов
	for stepName, result := range execCtx.Results {
		for k, v := range result.Output {
			input[fmt.Sprintf("%s.%s", stepName, k)] = v
		}
	}
	
	return input
}

// updateGlobalData обновляет глобальные данные из результата шага
func (e *Engine) updateGlobalData(execCtx *ExecutionContext, result *StepResult) {
	if result != nil && result.Output != nil {
		for k, v := range result.Output {
			execCtx.GlobalData[k] = v
		}
	}
}

// calculateRetryDelay вычисляет задержку для повтора
func (e *Engine) calculateRetryDelay(attempt int, config RetryConfig) time.Duration {
	delay := config.BaseDelay
	
	// Экспоненциальный backoff
	for i := 1; i < attempt; i++ {
		delay = time.Duration(float64(delay) * config.Multiplier)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
			break
		}
	}
	
	// Добавляем jitter если нужно
	if config.Jitter {
		jitter := time.Duration(float64(delay) * 0.1) // 10% jitter
		delay += time.Duration(int64(jitter) * (int64(time.Now().UnixNano()) % 2 - 1))
	}
	
	return delay
}

// generateID генерирует уникальный ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}