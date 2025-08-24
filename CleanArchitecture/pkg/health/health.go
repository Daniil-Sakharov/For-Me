package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 📚 ИЗУЧИТЬ: Health Check patterns
// https://microservices.io/patterns/observability/health-check-api.html
// https://tools.ietf.org/id/draft-inadarei-api-health-check-06.html
// https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
//
// ТИПЫ HEALTH CHECKS:
// 1. LIVENESS - приложение работает и может обрабатывать запросы
// 2. READINESS - приложение готово принимать трафик
// 3. STARTUP - приложение завершило инициализацию
//
// КОМПОНЕНТЫ ДЛЯ ПРОВЕРКИ:
// 1. База данных
// 2. Внешние сервисы (Redis, очереди)
// 3. Файловая система
// 4. Сетевые соединения
// 5. Бизнес-логические проверки

// Status статус health check
type Status string

const (
	StatusPass Status = "pass" // Компонент работает нормально
	StatusFail Status = "fail" // Компонент не работает
	StatusWarn Status = "warn" // Компонент работает, но есть проблемы
)

// CheckResult результат проверки компонента
type CheckResult struct {
	// Status статус проверки
	Status Status `json:"status"`
	
	// ComponentName имя компонента
	ComponentName string `json:"componentName,omitempty"`
	
	// ComponentType тип компонента (database, cache, external-service)
	ComponentType string `json:"componentType,omitempty"`
	
	// Time время проверки
	Time time.Time `json:"time"`
	
	// Duration время выполнения проверки
	Duration time.Duration `json:"duration,omitempty"`
	
	// Output дополнительная информация
	Output string `json:"output,omitempty"`
	
	// Error сообщение об ошибке
	Error string `json:"error,omitempty"`
	
	// Links полезные ссылки
	Links map[string]string `json:"links,omitempty"`
	
	// Metrics метрики компонента
	Metrics map[string]interface{} `json:"metrics,omitempty"`
}

// HealthResponse ответ health check API
type HealthResponse struct {
	// Status общий статус приложения
	Status Status `json:"status"`
	
	// Version версия приложения
	Version string `json:"version,omitempty"`
	
	// ReleaseID ID релиза
	ReleaseID string `json:"releaseId,omitempty"`
	
	// Notes дополнительные заметки
	Notes []string `json:"notes,omitempty"`
	
	// Output общий вывод
	Output string `json:"output,omitempty"`
	
	// ServiceID идентификатор сервиса
	ServiceID string `json:"serviceId,omitempty"`
	
	// Description описание сервиса
	Description string `json:"description,omitempty"`
	
	// Checks результаты проверок компонентов
	Checks map[string]CheckResult `json:"checks,omitempty"`
	
	// Links полезные ссылки
	Links map[string]string `json:"links,omitempty"`
	
	// Time время проверки
	Time time.Time `json:"time"`
	
	// Duration общее время проверки
	Duration time.Duration `json:"duration,omitempty"`
}

// Checker интерфейс для проверки компонента
type Checker interface {
	// Check выполняет проверку компонента
	Check(ctx context.Context) CheckResult
	
	// Name возвращает имя проверки
	Name() string
	
	// Type возвращает тип компонента
	Type() string
}

// CheckerFunc функция-адаптер для простых проверок
type CheckerFunc func(ctx context.Context) CheckResult

// Check реализует интерфейс Checker
func (f CheckerFunc) Check(ctx context.Context) CheckResult {
	return f(ctx)
}

// Name возвращает имя по умолчанию
func (f CheckerFunc) Name() string {
	return "anonymous"
}

// Type возвращает тип по умолчанию
func (f CheckerFunc) Type() string {
	return "component"
}

// HealthChecker основной health checker
type HealthChecker struct {
	checkers    map[string]Checker
	mu          sync.RWMutex
	config      Config
	version     string
	releaseID   string
	serviceID   string
	description string
}

// Config конфигурация health checker
type Config struct {
	// Timeout максимальное время выполнения всех проверок
	Timeout time.Duration `yaml:"timeout" env:"HEALTH_CHECK_TIMEOUT" default:"30s"`
	
	// CheckTimeout максимальное время выполнения одной проверки
	CheckTimeout time.Duration `yaml:"check_timeout" env:"HEALTH_CHECK_ITEM_TIMEOUT" default:"5s"`
	
	// FailFast прерывать проверку при первой ошибке
	FailFast bool `yaml:"fail_fast" env:"HEALTH_CHECK_FAIL_FAST" default:"false"`
	
	// Parallel выполнять проверки параллельно
	Parallel bool `yaml:"parallel" env:"HEALTH_CHECK_PARALLEL" default:"true"`
}

// NewHealthChecker создает новый health checker
func NewHealthChecker(config Config) *HealthChecker {
	return &HealthChecker{
		checkers: make(map[string]Checker),
		config:   config,
	}
}

// SetInfo устанавливает информацию о сервисе
func (hc *HealthChecker) SetInfo(version, releaseID, serviceID, description string) {
	hc.version = version
	hc.releaseID = releaseID
	hc.serviceID = serviceID
	hc.description = description
}

// Register регистрирует новую проверку
func (hc *HealthChecker) Register(name string, checker Checker) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checkers[name] = checker
}

// RegisterFunc регистрирует функцию как проверку
func (hc *HealthChecker) RegisterFunc(name, componentType string, checkFunc func(ctx context.Context) CheckResult) {
	checker := &namedChecker{
		name:         name,
		componentType: componentType,
		checkFunc:    checkFunc,
	}
	hc.Register(name, checker)
}

// Unregister удаляет проверку
func (hc *HealthChecker) Unregister(name string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	delete(hc.checkers, name)
}

// Check выполняет все проверки
func (hc *HealthChecker) Check(ctx context.Context) HealthResponse {
	start := time.Now()
	
	// Создаем контекст с таймаутом
	checkCtx, cancel := context.WithTimeout(ctx, hc.config.Timeout)
	defer cancel()
	
	hc.mu.RLock()
	checkers := make(map[string]Checker, len(hc.checkers))
	for name, checker := range hc.checkers {
		checkers[name] = checker
	}
	hc.mu.RUnlock()
	
	var checks map[string]CheckResult
	
	if hc.config.Parallel {
		checks = hc.runChecksParallel(checkCtx, checkers)
	} else {
		checks = hc.runChecksSequential(checkCtx, checkers)
	}
	
	// Определяем общий статус
	overallStatus := hc.calculateOverallStatus(checks)
	
	response := HealthResponse{
		Status:      overallStatus,
		Version:     hc.version,
		ReleaseID:   hc.releaseID,
		ServiceID:   hc.serviceID,
		Description: hc.description,
		Checks:      checks,
		Time:        start,
		Duration:    time.Since(start),
	}
	
	// Добавляем полезные ссылки
	response.Links = map[string]string{
		"self": "/health",
	}
	
	return response
}

// runChecksParallel выполняет проверки параллельно
func (hc *HealthChecker) runChecksParallel(ctx context.Context, checkers map[string]Checker) map[string]CheckResult {
	results := make(map[string]CheckResult)
	resultsMu := sync.Mutex{}
	
	var wg sync.WaitGroup
	
	for name, checker := range checkers {
		wg.Add(1)
		go func(name string, checker Checker) {
			defer wg.Done()
			
			// Контекст с таймаутом для одной проверки
			checkCtx, cancel := context.WithTimeout(ctx, hc.config.CheckTimeout)
			defer cancel()
			
			result := hc.runSingleCheck(checkCtx, checker)
			
			resultsMu.Lock()
			results[name] = result
			resultsMu.Unlock()
			
			// Если включен FailFast и проверка провалилась
			if hc.config.FailFast && result.Status == StatusFail {
				// В реальной реализации можно отменить остальные проверки
			}
		}(name, checker)
	}
	
	wg.Wait()
	return results
}

// runChecksSequential выполняет проверки последовательно
func (hc *HealthChecker) runChecksSequential(ctx context.Context, checkers map[string]Checker) map[string]CheckResult {
	results := make(map[string]CheckResult)
	
	for name, checker := range checkers {
		// Контекст с таймаутом для одной проверки
		checkCtx, cancel := context.WithTimeout(ctx, hc.config.CheckTimeout)
		
		result := hc.runSingleCheck(checkCtx, checker)
		results[name] = result
		
		cancel()
		
		// Если включен FailFast и проверка провалилась
		if hc.config.FailFast && result.Status == StatusFail {
			break
		}
	}
	
	return results
}

// runSingleCheck выполняет одну проверку с обработкой паники
func (hc *HealthChecker) runSingleCheck(ctx context.Context, checker Checker) CheckResult {
	start := time.Now()
	
	defer func() {
		if r := recover(); r != nil {
			// Обработка паники в проверке
		}
	}()
	
	result := checker.Check(ctx)
	result.Duration = time.Since(start)
	result.Time = start
	
	if result.ComponentName == "" {
		result.ComponentName = checker.Name()
	}
	if result.ComponentType == "" {
		result.ComponentType = checker.Type()
	}
	
	return result
}

// calculateOverallStatus вычисляет общий статус
func (hc *HealthChecker) calculateOverallStatus(checks map[string]CheckResult) Status {
	if len(checks) == 0 {
		return StatusPass
	}
	
	hasWarnings := false
	
	for _, check := range checks {
		switch check.Status {
		case StatusFail:
			return StatusFail // Любая неудачная проверка = общий fail
		case StatusWarn:
			hasWarnings = true
		}
	}
	
	if hasWarnings {
		return StatusWarn
	}
	
	return StatusPass
}

// ПРЕДОПРЕДЕЛЕННЫЕ ПРОВЕРКИ

// DatabaseChecker проверка подключения к базе данных
type DatabaseChecker struct {
	db   *sql.DB
	name string
}

// NewDatabaseChecker создает проверку БД
func NewDatabaseChecker(db *sql.DB, name string) *DatabaseChecker {
	return &DatabaseChecker{
		db:   db,
		name: name,
	}
}

// Check выполняет проверку БД
func (dc *DatabaseChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()
	
	// Проверяем подключение
	err := dc.db.PingContext(ctx)
	if err != nil {
		return CheckResult{
			Status:        StatusFail,
			ComponentName: dc.name,
			ComponentType: "database",
			Time:          start,
			Duration:      time.Since(start),
			Error:         err.Error(),
			Output:        "Database connection failed",
		}
	}
	
	// Получаем статистику подключений
	stats := dc.db.Stats()
	metrics := map[string]interface{}{
		"open_connections":     stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration_ms":    stats.WaitDuration.Milliseconds(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
	
	status := StatusPass
	output := "Database connection successful"
	
	// Проверяем на предупреждения
	if stats.OpenConnections > 80 { // Если больше 80% от максимума
		status = StatusWarn
		output = "High number of database connections"
	}
	
	return CheckResult{
		Status:        status,
		ComponentName: dc.name,
		ComponentType: "database",
		Time:          start,
		Duration:      time.Since(start),
		Output:        output,
		Metrics:       metrics,
	}
}

// Name возвращает имя проверки
func (dc *DatabaseChecker) Name() string {
	return dc.name
}

// Type возвращает тип компонента
func (dc *DatabaseChecker) Type() string {
	return "database"
}

// HTTPChecker проверка внешнего HTTP сервиса
type HTTPChecker struct {
	url        string
	name       string
	client     *http.Client
	expectedStatus int
}

// NewHTTPChecker создает проверку HTTP сервиса
func NewHTTPChecker(url, name string, timeout time.Duration) *HTTPChecker {
	return &HTTPChecker{
		url:  url,
		name: name,
		client: &http.Client{
			Timeout: timeout,
		},
		expectedStatus: http.StatusOK,
	}
}

// WithExpectedStatus устанавливает ожидаемый статус код
func (hc *HTTPChecker) WithExpectedStatus(status int) *HTTPChecker {
	hc.expectedStatus = status
	return hc
}

// Check выполняет HTTP проверку
func (hc *HTTPChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(ctx, "GET", hc.url, nil)
	if err != nil {
		return CheckResult{
			Status:        StatusFail,
			ComponentName: hc.name,
			ComponentType: "http",
			Time:          start,
			Duration:      time.Since(start),
			Error:         err.Error(),
			Output:        "Failed to create HTTP request",
		}
	}
	
	resp, err := hc.client.Do(req)
	if err != nil {
		return CheckResult{
			Status:        StatusFail,
			ComponentName: hc.name,
			ComponentType: "http",
			Time:          start,
			Duration:      time.Since(start),
			Error:         err.Error(),
			Output:        "HTTP request failed",
		}
	}
	defer resp.Body.Close()
	
	status := StatusPass
	output := fmt.Sprintf("HTTP check successful (status: %d)", resp.StatusCode)
	
	if resp.StatusCode != hc.expectedStatus {
		status = StatusFail
		output = fmt.Sprintf("Unexpected status code: %d (expected: %d)", resp.StatusCode, hc.expectedStatus)
	}
	
	metrics := map[string]interface{}{
		"status_code":   resp.StatusCode,
		"response_time": time.Since(start).Milliseconds(),
	}
	
	return CheckResult{
		Status:        status,
		ComponentName: hc.name,
		ComponentType: "http",
		Time:          start,
		Duration:      time.Since(start),
		Output:        output,
		Metrics:       metrics,
		Links: map[string]string{
			"endpoint": hc.url,
		},
	}
}

// Name возвращает имя проверки
func (hc *HTTPChecker) Name() string {
	return hc.name
}

// Type возвращает тип компонента
func (hc *HTTPChecker) Type() string {
	return "http"
}

// ВСПОМОГАТЕЛЬНЫЕ ТИПЫ

// namedChecker обертка для функций с именем и типом
type namedChecker struct {
	name          string
	componentType string
	checkFunc     func(ctx context.Context) CheckResult
}

func (nc *namedChecker) Check(ctx context.Context) CheckResult {
	return nc.checkFunc(ctx)
}

func (nc *namedChecker) Name() string {
	return nc.name
}

func (nc *namedChecker) Type() string {
	return nc.componentType
}

// HTTP HANDLER для health checks

// Handler HTTP обработчик для health checks
func (hc *HealthChecker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		
		// Выполняем проверки
		result := hc.Check(r.Context())
		
		// Устанавливаем HTTP статус на основе результата
		switch result.Status {
		case StatusPass:
			w.WriteHeader(http.StatusOK)
		case StatusWarn:
			w.WriteHeader(http.StatusOK) // 200 для warnings
		case StatusFail:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		
		// Отправляем JSON ответ
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to encode health check response", http.StatusInternalServerError)
		}
	}
}

// LivenessHandler простой liveness check (всегда возвращает 200)
func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status": "pass",
			"time":   time.Now(),
		}
		json.NewEncoder(w).Encode(response)
	}
}

// ReadinessHandler readiness check с минимальными проверками
func (hc *HealthChecker) ReadinessHandler(criticalChecks []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Выполняем только критичные проверки
		result := hc.checkCriticalComponents(r.Context(), criticalChecks)
		
		switch result.Status {
		case StatusPass:
			w.WriteHeader(http.StatusOK)
		case StatusWarn:
			w.WriteHeader(http.StatusOK)
		case StatusFail:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		
		json.NewEncoder(w).Encode(result)
	}
}

// checkCriticalComponents проверяет только критичные компоненты
func (hc *HealthChecker) checkCriticalComponents(ctx context.Context, criticalChecks []string) HealthResponse {
	start := time.Now()
	
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	checks := make(map[string]CheckResult)
	
	for _, name := range criticalChecks {
		if checker, exists := hc.checkers[name]; exists {
			checkCtx, cancel := context.WithTimeout(ctx, hc.config.CheckTimeout)
			result := hc.runSingleCheck(checkCtx, checker)
			checks[name] = result
			cancel()
		}
	}
	
	return HealthResponse{
		Status:   hc.calculateOverallStatus(checks),
		Checks:   checks,
		Time:     start,
		Duration: time.Since(start),
	}
}

// ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ:

/*
// В main.go:
func main() {
    // Создаем health checker
    healthChecker := health.NewHealthChecker(health.Config{
        Timeout:      30 * time.Second,
        CheckTimeout: 5 * time.Second,
        Parallel:     true,
    })
    
    // Устанавливаем информацию о сервисе
    healthChecker.SetInfo("1.0.0", "abc123", "url-shortener", "URL Shortening Service")
    
    // Регистрируем проверки
    healthChecker.Register("database", health.NewDatabaseChecker(db, "postgres"))
    healthChecker.Register("redis", health.NewHTTPChecker("http://redis:6379/ping", "redis", 2*time.Second))
    
    // Кастомная проверка
    healthChecker.RegisterFunc("disk_space", "filesystem", func(ctx context.Context) health.CheckResult {
        // Проверяем свободное место на диске
        return health.CheckResult{
            Status: health.StatusPass,
            Output: "Disk space OK",
        }
    })
    
    // Настраиваем HTTP endpoints
    http.HandleFunc("/health", healthChecker.Handler())
    http.HandleFunc("/health/live", health.LivenessHandler())
    http.HandleFunc("/health/ready", healthChecker.ReadinessHandler([]string{"database"}))
    
    // Запускаем сервер
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// В Kubernetes:
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    image: url-shortener:latest
    livenessProbe:
      httpGet:
        path: /health/live
        port: 8080
      initialDelaySeconds: 30
      periodSeconds: 10
    readinessProbe:
      httpGet:
        path: /health/ready
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5

// В Use Case (опционально):
type LoginUseCase struct {
    healthChecker *health.HealthChecker
    // ... другие зависимости
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
    // Можно проверить критичные компоненты перед выполнением
    if uc.healthChecker != nil {
        result := uc.healthChecker.Check(ctx)
        if result.Status == health.StatusFail {
            return nil, errors.New("system is unhealthy")
        }
    }
    
    // Бизнес-логика...
}
*/

// ВАЖНЫЕ ПРИНЦИПЫ:
// 1. Health checks должны быть быстрыми (< 5 секунд)
// 2. Различайте liveness, readiness и startup probes
// 3. Не включайте в readiness проверки внешних сервисов, если приложение может работать без них
// 4. Логируйте результаты health checks для диагностики
// 5. Используйте таймауты для всех проверок
// 6. Кэшируйте результаты проверок, если они дорогие
// 7. Предоставляйте детальную информацию о проблемах