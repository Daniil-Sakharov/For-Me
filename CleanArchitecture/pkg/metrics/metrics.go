package metrics

import (
	"context"
	"sync"
	"time"
)

// 📚 ИЗУЧИТЬ: Application Metrics и Observability
// https://prometheus.io/docs/concepts/metric_types/
// https://opentelemetry.io/docs/concepts/observability-primer/
// https://sre.google/sre-book/monitoring-distributed-systems/
// https://grafana.com/blog/2018/08/02/the-red-method-how-to-instrument-your-services/
//
// ОСНОВНЫЕ ТИПЫ МЕТРИК:
// 1. COUNTER - монотонно возрастающий счетчик (количество запросов, ошибок)
// 2. GAUGE - текущее значение (количество активных пользователей, память)
// 3. HISTOGRAM - распределение значений (время отклика, размеры запросов)
// 4. SUMMARY - квантили и суммы (95-й перцентиль времени отклика)
//
// МЕТОДОЛОГИИ МОНИТОРИНГА:
// 1. RED Method (Rate, Errors, Duration) - для сервисов
// 2. USE Method (Utilization, Saturation, Errors) - для ресурсов
// 3. Four Golden Signals (Latency, Traffic, Errors, Saturation)

// MetricsCollector интерфейс для сбора метрик
type MetricsCollector interface {
	// Counter methods
	IncCounter(name string, tags map[string]string)
	IncCounterBy(name string, value int64, tags map[string]string)
	
	// Gauge methods
	SetGauge(name string, value float64, tags map[string]string)
	IncGauge(name string, tags map[string]string)
	DecGauge(name string, tags map[string]string)
	
	// Histogram methods
	RecordDuration(name string, duration time.Duration, tags map[string]string)
	RecordValue(name string, value float64, tags map[string]string)
	
	// Distribution methods
	RecordDistribution(name string, value float64, tags map[string]string)
	
	// Utility methods
	WithTags(tags map[string]string) MetricsCollector
	Close() error
}

// ПРОСТАЯ IN-MEMORY РЕАЛИЗАЦИЯ ДЛЯ РАЗРАБОТКИ

// InMemoryMetrics простая реализация метрик в памяти
type InMemoryMetrics struct {
	counters    map[string]*Counter
	gauges      map[string]*Gauge
	histograms  map[string]*Histogram
	mu          sync.RWMutex
	defaultTags map[string]string
}

// Counter счетчик
type Counter struct {
	Value int64
	Tags  map[string]string
	mu    sync.RWMutex
}

// Gauge измеритель
type Gauge struct {
	Value float64
	Tags  map[string]string
	mu    sync.RWMutex
}

// Histogram гистограмма
type Histogram struct {
	Count    int64
	Sum      float64
	Min      float64
	Max      float64
	Buckets  map[float64]int64 // bucket -> count
	Tags     map[string]string
	mu       sync.RWMutex
}

// NewInMemoryMetrics создает новый сборщик метрик в памяти
func NewInMemoryMetrics() *InMemoryMetrics {
	return &InMemoryMetrics{
		counters:   make(map[string]*Counter),
		gauges:     make(map[string]*Gauge),
		histograms: make(map[string]*Histogram),
	}
}

// IncCounter увеличивает счетчик на 1
func (m *InMemoryMetrics) IncCounter(name string, tags map[string]string) {
	m.IncCounterBy(name, 1, tags)
}

// IncCounterBy увеличивает счетчик на указанное значение
func (m *InMemoryMetrics) IncCounterBy(name string, value int64, tags map[string]string) {
	key := m.createKey(name, tags)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	counter, exists := m.counters[key]
	if !exists {
		counter = &Counter{
			Value: 0,
			Tags:  m.mergeTags(tags),
		}
		m.counters[key] = counter
	}
	
	counter.mu.Lock()
	counter.Value += value
	counter.mu.Unlock()
}

// SetGauge устанавливает значение измерителя
func (m *InMemoryMetrics) SetGauge(name string, value float64, tags map[string]string) {
	key := m.createKey(name, tags)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	gauge, exists := m.gauges[key]
	if !exists {
		gauge = &Gauge{
			Tags: m.mergeTags(tags),
		}
		m.gauges[key] = gauge
	}
	
	gauge.mu.Lock()
	gauge.Value = value
	gauge.mu.Unlock()
}

// IncGauge увеличивает измеритель на 1
func (m *InMemoryMetrics) IncGauge(name string, tags map[string]string) {
	m.adjustGauge(name, 1, tags)
}

// DecGauge уменьшает измеритель на 1
func (m *InMemoryMetrics) DecGauge(name string, tags map[string]string) {
	m.adjustGauge(name, -1, tags)
}

// adjustGauge изменяет измеритель на указанное значение
func (m *InMemoryMetrics) adjustGauge(name string, delta float64, tags map[string]string) {
	key := m.createKey(name, tags)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	gauge, exists := m.gauges[key]
	if !exists {
		gauge = &Gauge{
			Value: 0,
			Tags:  m.mergeTags(tags),
		}
		m.gauges[key] = gauge
	}
	
	gauge.mu.Lock()
	gauge.Value += delta
	gauge.mu.Unlock()
}

// RecordDuration записывает длительность в гистограмму
func (m *InMemoryMetrics) RecordDuration(name string, duration time.Duration, tags map[string]string) {
	m.RecordValue(name, float64(duration.Nanoseconds())/1e6, tags) // в миллисекундах
}

// RecordValue записывает значение в гистограмму
func (m *InMemoryMetrics) RecordValue(name string, value float64, tags map[string]string) {
	key := m.createKey(name, tags)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	histogram, exists := m.histograms[key]
	if !exists {
		histogram = &Histogram{
			Min:     value,
			Max:     value,
			Buckets: make(map[float64]int64),
			Tags:    m.mergeTags(tags),
		}
		m.histograms[key] = histogram
	}
	
	histogram.mu.Lock()
	histogram.Count++
	histogram.Sum += value
	
	if value < histogram.Min {
		histogram.Min = value
	}
	if value > histogram.Max {
		histogram.Max = value
	}
	
	// Добавляем в buckets (простая реализация)
	bucket := m.findBucket(value)
	histogram.Buckets[bucket]++
	histogram.mu.Unlock()
}

// RecordDistribution записывает значение в распределение (алиас для RecordValue)
func (m *InMemoryMetrics) RecordDistribution(name string, value float64, tags map[string]string) {
	m.RecordValue(name, value, tags)
}

// WithTags создает новый сборщик с дополнительными тегами
func (m *InMemoryMetrics) WithTags(tags map[string]string) MetricsCollector {
	return &InMemoryMetrics{
		counters:    m.counters,
		gauges:      m.gauges,
		histograms:  m.histograms,
		defaultTags: m.mergeTags(tags),
	}
}

// Close закрывает сборщик метрик
func (m *InMemoryMetrics) Close() error {
	// Для in-memory реализации ничего не делаем
	return nil
}

// UTILITY METHODS

// createKey создает уникальный ключ для метрики с тегами
func (m *InMemoryMetrics) createKey(name string, tags map[string]string) string {
	key := name
	allTags := m.mergeTags(tags)
	
	for k, v := range allTags {
		key += "," + k + "=" + v
	}
	
	return key
}

// mergeTags объединяет теги с дефолтными
func (m *InMemoryMetrics) mergeTags(tags map[string]string) map[string]string {
	result := make(map[string]string)
	
	// Сначала дефолтные теги
	for k, v := range m.defaultTags {
		result[k] = v
	}
	
	// Затем переданные теги (перезаписывают дефолтные)
	for k, v := range tags {
		result[k] = v
	}
	
	return result
}

// findBucket находит подходящий bucket для значения
func (m *InMemoryMetrics) findBucket(value float64) float64 {
	// Простые buckets: 1, 5, 10, 25, 50, 100, 250, 500, 1000, +Inf
	buckets := []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, float64(^uint(0) >> 1)} // +Inf
	
	for _, bucket := range buckets {
		if value <= bucket {
			return bucket
		}
	}
	
	return buckets[len(buckets)-1]
}

// СТАТИСТИКА И ЭКСПОРТ МЕТРИК

// GetStats возвращает статистику всех метрик
func (m *InMemoryMetrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stats := make(map[string]interface{})
	
	// Счетчики
	counters := make(map[string]interface{})
	for key, counter := range m.counters {
		counter.mu.RLock()
		counters[key] = map[string]interface{}{
			"value": counter.Value,
			"tags":  counter.Tags,
		}
		counter.mu.RUnlock()
	}
	stats["counters"] = counters
	
	// Измерители
	gauges := make(map[string]interface{})
	for key, gauge := range m.gauges {
		gauge.mu.RLock()
		gauges[key] = map[string]interface{}{
			"value": gauge.Value,
			"tags":  gauge.Tags,
		}
		gauge.mu.RUnlock()
	}
	stats["gauges"] = gauges
	
	// Гистограммы
	histograms := make(map[string]interface{})
	for key, histogram := range m.histograms {
		histogram.mu.RLock()
		avg := float64(0)
		if histogram.Count > 0 {
			avg = histogram.Sum / float64(histogram.Count)
		}
		
		histograms[key] = map[string]interface{}{
			"count":   histogram.Count,
			"sum":     histogram.Sum,
			"min":     histogram.Min,
			"max":     histogram.Max,
			"avg":     avg,
			"buckets": histogram.Buckets,
			"tags":    histogram.Tags,
		}
		histogram.mu.RUnlock()
	}
	stats["histograms"] = histograms
	
	return stats
}

// СПЕЦИАЛИЗИРОВАННЫЕ МЕТРИКИ ДЛЯ CLEAN ARCHITECTURE

// ApplicationMetrics метрики для приложения
type ApplicationMetrics struct {
	collector MetricsCollector
}

// NewApplicationMetrics создает метрики приложения
func NewApplicationMetrics(collector MetricsCollector) *ApplicationMetrics {
	return &ApplicationMetrics{
		collector: collector.WithTags(map[string]string{
			"service": "url-shortener",
		}),
	}
}

// HTTP МЕТРИКИ

// RecordHTTPRequest записывает HTTP запрос
func (m *ApplicationMetrics) RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	tags := map[string]string{
		"method": method,
		"path":   path,
		"status": getStatusClass(statusCode),
	}
	
	// Счетчик запросов
	m.collector.IncCounter("http_requests_total", tags)
	
	// Длительность запроса
	m.collector.RecordDuration("http_request_duration_ms", duration, tags)
	
	// Счетчик ошибок
	if statusCode >= 400 {
		m.collector.IncCounter("http_errors_total", tags)
	}
}

// USE CASE МЕТРИКИ

// RecordUseCaseExecution записывает выполнение Use Case
func (m *ApplicationMetrics) RecordUseCaseExecution(useCaseName string, success bool, duration time.Duration) {
	tags := map[string]string{
		"use_case": useCaseName,
		"status":   getSuccessStatus(success),
	}
	
	// Счетчик выполнений
	m.collector.IncCounter("usecase_executions_total", tags)
	
	// Длительность выполнения
	m.collector.RecordDuration("usecase_duration_ms", duration, tags)
	
	// Счетчик ошибок
	if !success {
		m.collector.IncCounter("usecase_errors_total", tags)
	}
}

// DATABASE МЕТРИКИ

// RecordDatabaseQuery записывает выполнение запроса к БД
func (m *ApplicationMetrics) RecordDatabaseQuery(operation string, table string, success bool, duration time.Duration) {
	tags := map[string]string{
		"operation": operation,
		"table":     table,
		"status":    getSuccessStatus(success),
	}
	
	// Счетчик запросов
	m.collector.IncCounter("db_queries_total", tags)
	
	// Длительность запроса
	m.collector.RecordDuration("db_query_duration_ms", duration, tags)
	
	// Счетчик ошибок
	if !success {
		m.collector.IncCounter("db_errors_total", tags)
	}
}

// CACHE МЕТРИКИ

// RecordCacheOperation записывает операцию с кэшем
func (m *ApplicationMetrics) RecordCacheOperation(operation string, hit bool) {
	tags := map[string]string{
		"operation": operation,
		"result":    getCacheResult(hit),
	}
	
	// Счетчик операций
	m.collector.IncCounter("cache_operations_total", tags)
	
	// Отдельные счетчики для hits/misses
	if operation == "get" {
		if hit {
			m.collector.IncCounter("cache_hits_total", nil)
		} else {
			m.collector.IncCounter("cache_misses_total", nil)
		}
	}
}

// BUSINESS МЕТРИКИ

// RecordUserRegistration записывает регистрацию пользователя
func (m *ApplicationMetrics) RecordUserRegistration(success bool) {
	tags := map[string]string{
		"status": getSuccessStatus(success),
	}
	
	m.collector.IncCounter("user_registrations_total", tags)
}

// RecordUserLogin записывает попытку входа
func (m *ApplicationMetrics) RecordUserLogin(success bool) {
	tags := map[string]string{
		"status": getSuccessStatus(success),
	}
	
	m.collector.IncCounter("user_logins_total", tags)
}

// RecordLinkCreation записывает создание ссылки
func (m *ApplicationMetrics) RecordLinkCreation(success bool) {
	tags := map[string]string{
		"status": getSuccessStatus(success),
	}
	
	m.collector.IncCounter("links_created_total", tags)
}

// RecordLinkRedirect записывает переход по ссылке
func (m *ApplicationMetrics) RecordLinkRedirect(found bool) {
	tags := map[string]string{
		"found": getBooleanString(found),
	}
	
	m.collector.IncCounter("link_redirects_total", tags)
}

// SYSTEM МЕТРИКИ

// SetActiveConnections устанавливает количество активных соединений
func (m *ApplicationMetrics) SetActiveConnections(count int) {
	m.collector.SetGauge("active_connections", float64(count), nil)
}

// SetActiveUsers устанавливает количество активных пользователей
func (m *ApplicationMetrics) SetActiveUsers(count int) {
	m.collector.SetGauge("active_users", float64(count), nil)
}

// RecordMemoryUsage записывает использование памяти
func (m *ApplicationMetrics) RecordMemoryUsage(heapMB, stackMB float64) {
	m.collector.SetGauge("memory_heap_mb", heapMB, nil)
	m.collector.SetGauge("memory_stack_mb", stackMB, nil)
}

// RATE LIMITING МЕТРИКИ

// RecordRateLimitCheck записывает проверку rate limit
func (m *ApplicationMetrics) RecordRateLimitCheck(limiterName string, allowed bool) {
	tags := map[string]string{
		"limiter": limiterName,
		"result":  getRateLimitResult(allowed),
	}
	
	m.collector.IncCounter("rate_limit_checks_total", tags)
	
	if !allowed {
		m.collector.IncCounter("rate_limit_exceeded_total", tags)
	}
}

// UTILITY FUNCTIONS

func getStatusClass(statusCode int) string {
	switch {
	case statusCode < 300:
		return "2xx"
	case statusCode < 400:
		return "3xx"
	case statusCode < 500:
		return "4xx"
	default:
		return "5xx"
	}
}

func getSuccessStatus(success bool) string {
	if success {
		return "success"
	}
	return "error"
}

func getCacheResult(hit bool) string {
	if hit {
		return "hit"
	}
	return "miss"
}

func getBooleanString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func getRateLimitResult(allowed bool) string {
	if allowed {
		return "allowed"
	}
	return "blocked"
}

// МОНИТОРИНГ ПРОИЗВОДИТЕЛЬНОСТИ USE CASES

// UseCaseTimer помощник для измерения времени выполнения Use Cases
type UseCaseTimer struct {
	metrics     *ApplicationMetrics
	useCaseName string
	startTime   time.Time
}

// NewUseCaseTimer создает новый таймер для Use Case
func (m *ApplicationMetrics) NewUseCaseTimer(useCaseName string) *UseCaseTimer {
	return &UseCaseTimer{
		metrics:     m,
		useCaseName: useCaseName,
		startTime:   time.Now(),
	}
}

// Finish завершает измерение и записывает метрики
func (t *UseCaseTimer) Finish(err error) {
	duration := time.Since(t.startTime)
	success := err == nil
	t.metrics.RecordUseCaseExecution(t.useCaseName, success, duration)
}

// КОНФИГУРАЦИЯ МЕТРИК

// MetricsConfig конфигурация системы метрик
type MetricsConfig struct {
	Enabled    bool   `yaml:"enabled" env:"METRICS_ENABLED" default:"true"`
	Provider   string `yaml:"provider" env:"METRICS_PROVIDER" default:"memory"` // memory, prometheus, datadog
	Address    string `yaml:"address" env:"METRICS_ADDRESS" default:":9090"`
	Path       string `yaml:"path" env:"METRICS_PATH" default:"/metrics"`
	Namespace  string `yaml:"namespace" env:"METRICS_NAMESPACE" default:"urlshortener"`
	
	// Sampling (для высоконагруженных систем)
	SampleRate float64 `yaml:"sample_rate" env:"METRICS_SAMPLE_RATE" default:"1.0"`
}

// ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ:

/*
// В Use Case:
type LoginUseCase struct {
    metrics *metrics.ApplicationMetrics
    // ... другие зависимости
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
    timer := uc.metrics.NewUseCaseTimer("user_login")
    defer timer.Finish(err)
    
    // Бизнес-логика...
    
    uc.metrics.RecordUserLogin(true)
    return response, nil
}

// В HTTP middleware:
func MetricsMiddleware(metrics *metrics.ApplicationMetrics) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Обработка запроса
            recorder := &statusRecorder{ResponseWriter: w, statusCode: 200}
            next.ServeHTTP(recorder, r)
            
            // Записываем метрики
            duration := time.Since(start)
            metrics.RecordHTTPRequest(r.Method, r.URL.Path, recorder.statusCode, duration)
        })
    }
}

// В Repository:
type UserRepository struct {
    metrics *metrics.ApplicationMetrics
    // ... другие зависимости
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        success := err == nil
        r.metrics.RecordDatabaseQuery("select", "users", success, duration)
    }()
    
    // SQL запрос...
}

// В main.go:
func main() {
    metricsCollector := metrics.NewInMemoryMetrics()
    appMetrics := metrics.NewApplicationMetrics(metricsCollector)
    
    // Передаем метрики во все слои через DI
    container := di.NewContainer(config, appMetrics)
    
    // Запускаем HTTP endpoint для метрик
    go func() {
        http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
            stats := metricsCollector.GetStats()
            json.NewEncoder(w).Encode(stats)
        })
        log.Fatal(http.ListenAndServe(":9090", nil))
    }()
    
    // ... остальная логика приложения
}
*/

// ВАЖНЫЕ ПРИНЦИПЫ:
// 1. Метрики НЕ ДОЛЖНЫ влиять на бизнес-логику
// 2. Используйте осмысленные теги для группировки метрик
// 3. Не создавайте слишком много уникальных комбинаций тегов (cardinality explosion)
// 4. Мониторьте производительность самой системы метрик
// 5. Используйте sampling для высоконагруженных систем
// 6. Настройте алерты на критичные метрики
// 7. Регулярно пересматривайте и очищайте неиспользуемые метрики