// =====================================================================
// ⚙️ КОНФИГУРАЦИЯ ПРИЛОЖЕНИЯ - Управление настройками системы
// =====================================================================
//
// Этот файл содержит всю конфигурацию для Tender Automation MVP.
// Следует принципам 12-factor app:
// 1. Конфигурация хранится в переменных окружения
// 2. Строгое разделение между кодом и конфигурацией
// 3. Дефолтные значения для development окружения
// 4. Валидация критически важных настроек
//
// TODO: При расширении функциональности добавить:
// - Конфигурацию для Redis (кеширование, очереди)
// - Настройки для email отправки (SMTP/IMAP)
// - Конфигурацию для файлового хранилища
// - Настройки мониторинга и метрик

package configs

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// =====================================================================
// 📊 ОСНОВНАЯ СТРУКТУРА КОНФИГУРАЦИИ
// =====================================================================

// Config содержит всю конфигурацию приложения
//
// TODO: При добавлении новых компонентов расширить соответствующими секциями
// TODO: Добавить hot-reload конфигурации для некритичных настроек
type Config struct {
	// 🖥️ Настройки HTTP сервера
	Server ServerConfig `mapstructure:"server" validate:"required"`

	// 🗃️ Настройки базы данных
	Database DatabaseConfig `mapstructure:"database" validate:"required"`

	// 🤖 Настройки AI интеграции (Ollama/Llama)
	AI AIConfig `mapstructure:"ai" validate:"required"`

	// 🕷️ Настройки web scraping
	Scraping ScrapingConfig `mapstructure:"scraping" validate:"required"`

	// 📝 Настройки логирования
	Logging LoggingConfig `mapstructure:"logging" validate:"required"`

	// 🎯 Бизнес-логика настройки
	Business BusinessConfig `mapstructure:"business" validate:"required"`

	// 🔒 Настройки безопасности
	Security SecurityConfig `mapstructure:"security" validate:"required"`

	// 📊 Настройки мониторинга
	Monitoring MonitoringConfig `mapstructure:"monitoring" validate:"required"`
}

// =====================================================================
// 🖥️ КОНФИГУРАЦИЯ HTTP СЕРВЕРА
// =====================================================================

// ServerConfig содержит настройки HTTP сервера
//
// TODO: Добавить настройки TLS для production
// TODO: Добавить настройки rate limiting
type ServerConfig struct {
	// 🌐 Сетевые настройки
	Host string `mapstructure:"host" validate:"required" default:"0.0.0.0"`
	Port int    `mapstructure:"port" validate:"required,min=1,max=65535" default:"8080"`

	// ⏱️ Таймауты
	ReadTimeout     time.Duration `mapstructure:"read_timeout" default:"30s"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" default:"30s"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" default:"10s"`

	// 📏 Лимиты
	MaxRequestSize string `mapstructure:"max_request_size" default:"10MB"`

	// 🔧 Режим работы
	Mode string `mapstructure:"mode" validate:"oneof=debug release test" default:"debug"`
}

// GetAddress возвращает полный адрес сервера
func (s ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// =====================================================================
// 🗃️ КОНФИГУРАЦИЯ БАЗЫ ДАННЫХ
// =====================================================================

// DatabaseConfig содержит настройки подключения к PostgreSQL
//
// TODO: Добавить настройки для connection pooling
// TODO: Добавить настройки для read replicas
type DatabaseConfig struct {
	// 🔗 Параметры подключения
	Host     string `mapstructure:"host" validate:"required" default:"localhost"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535" default:"5432"`
	User     string `mapstructure:"user" validate:"required" default:"postgres"`
	Password string `mapstructure:"password" validate:"required"`
	DBName   string `mapstructure:"dbname" validate:"required" default:"tender_automation"`
	SSLMode  string `mapstructure:"sslmode" validate:"oneof=disable require verify-ca verify-full" default:"disable"`

	// 🏊 Connection pool настройки
	MaxOpenConns    int           `mapstructure:"max_open_conns" validate:"min=1" default:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" validate:"min=1" default:"5"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" default:"1h"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" default:"30m"`

	// 🕐 Таймауты
	ConnectTimeout time.Duration `mapstructure:"connect_timeout" default:"10s"`
	QueryTimeout   time.Duration `mapstructure:"query_timeout" default:"30s"`
}

// GetDSN возвращает Data Source Name для подключения к PostgreSQL
func (d DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// =====================================================================
// 🤖 КОНФИГУРАЦИЯ AI ИНТЕГРАЦИИ
// =====================================================================

// AIConfig содержит настройки для работы с AI (Ollama/Llama)
//
// TODO: Добавить настройки для fallback моделей
// TODO: Добавить настройки для кеширования AI ответов
type AIConfig struct {
	// 🔗 Подключение к Ollama
	URL   string `mapstructure:"url" validate:"required,url" default:"http://localhost:11434"`
	Model string `mapstructure:"model" validate:"required" default:"llama2"`

	// ⏱️ Таймауты и retry
	Timeout    time.Duration `mapstructure:"timeout" default:"60s"`
	MaxRetries int           `mapstructure:"max_retries" validate:"min=0" default:"3"`
	RetryDelay time.Duration `mapstructure:"retry_delay" default:"5s"`

	// 🎯 Настройки анализа
	RelevanceThreshold float64 `mapstructure:"relevance_threshold" validate:"min=0,max=1" default:"0.7"`
	BatchSize          int     `mapstructure:"batch_size" validate:"min=1" default:"10"`

	// 🔧 Дополнительные параметры
	Temperature float64 `mapstructure:"temperature" validate:"min=0,max=2" default:"0.1"`
	MaxTokens   int     `mapstructure:"max_tokens" validate:"min=1" default:"1000"`
}

// =====================================================================
// 🕷️ КОНФИГУРАЦИЯ WEB SCRAPING
// =====================================================================

// ScrapingConfig содержит настройки для парсинга тендерных площадок
//
// TODO: Добавить настройки для proxy rotation
// TODO: Добавить настройки для captcha решения
type ScrapingConfig struct {
	// 🚦 Rate limiting для избежания блокировок
	Delay          time.Duration `mapstructure:"delay" default:"2s"`
	MaxConcurrent  int           `mapstructure:"max_concurrent" validate:"min=1" default:"3"`
	UserAgent      string        `mapstructure:"user_agent" default:"Mozilla/5.0 (compatible; TenderBot/1.0)"`

	// 🔄 Retry настройки
	MaxRetries   int           `mapstructure:"max_retries" validate:"min=0" default:"3"`
	RetryDelay   time.Duration `mapstructure:"retry_delay" default:"5s"`

	// ⏱️ Таймауты
	RequestTimeout time.Duration `mapstructure:"request_timeout" default:"30s"`
	PageTimeout    time.Duration `mapstructure:"page_timeout" default:"60s"`

	// 🎯 Платформы для парсинга
	EnabledPlatforms []string `mapstructure:"enabled_platforms" default:"zakupki"`

	// 📊 Batch настройки
	BatchSize     int           `mapstructure:"batch_size" validate:"min=1" default:"50"`
	ScanInterval  time.Duration `mapstructure:"scan_interval" default:"1h"`
}

// IsPlatformEnabled проверяет, включена ли платформа для парсинга
func (s ScrapingConfig) IsPlatformEnabled(platform string) bool {
	for _, enabled := range s.EnabledPlatforms {
		if enabled == platform {
			return true
		}
	}
	return false
}

// =====================================================================
// 📝 КОНФИГУРАЦИЯ ЛОГИРОВАНИЯ
// =====================================================================

// LoggingConfig содержит настройки системы логирования
//
// TODO: Добавить настройки для отправки логов в внешние системы
// TODO: Добавить настройки для ротации логов
type LoggingConfig struct {
	// 📊 Уровень и формат логирования
	Level  string `mapstructure:"level" validate:"oneof=debug info warn error fatal" default:"info"`
	Format string `mapstructure:"format" validate:"oneof=json console" default:"json"`
	Output string `mapstructure:"output" validate:"oneof=stdout stderr file" default:"stdout"`

	// 🗂️ Файловое логирование (если output=file)
	FileName   string `mapstructure:"file_name" default:"app.log"`
	MaxSize    int    `mapstructure:"max_size" default:"100"` // MB
	MaxBackups int    `mapstructure:"max_backups" default:"3"`
	MaxAge     int    `mapstructure:"max_age" default:"28"` // days

	// 🔧 Дополнительные настройки
	EnableCaller     bool `mapstructure:"enable_caller" default:"true"`
	EnableStacktrace bool `mapstructure:"enable_stacktrace" default:"false"`
}

// =====================================================================
// 🎯 КОНФИГУРАЦИЯ БИЗНЕС-ЛОГИКИ
// =====================================================================

// BusinessConfig содержит настройки бизнес-логики приложения
//
// TODO: Добавить настройки для различных алгоритмов анализа
// TODO: Добавить настройки для интеграции с внешними сервисами
type BusinessConfig struct {
	// 🤖 Настройки AI анализа
	MinAIScoreThreshold float64 `mapstructure:"min_ai_score_threshold" validate:"min=0,max=1" default:"0.7"`
	AnalysisBatchSize   int     `mapstructure:"analysis_batch_size" validate:"min=1" default:"10"`

	// 📊 Настройки пагинации
	DefaultPageSize int `mapstructure:"default_page_size" validate:"min=1" default:"20"`
	MaxPageSize     int `mapstructure:"max_page_size" validate:"min=1" default:"100"`

	// 🏷️ Статусы по умолчанию
	DefaultTenderStatus string `mapstructure:"default_tender_status" default:"active"`

	// ⏱️ Интервалы обработки
	TenderDiscoveryInterval time.Duration `mapstructure:"tender_discovery_interval" default:"1h"`
	AIAnalysisInterval      time.Duration `mapstructure:"ai_analysis_interval" default:"30m"`
	CleanupInterval         time.Duration `mapstructure:"cleanup_interval" default:"24h"`
}

// =====================================================================
// 🔒 КОНФИГУРАЦИЯ БЕЗОПАСНОСТИ
// =====================================================================

// SecurityConfig содержит настройки безопасности
//
// TODO: Добавить настройки JWT аутентификации
// TODO: Добавить настройки CORS для фронтенда
type SecurityConfig struct {
	// 🌐 CORS настройки
	CORSAllowedOrigins     []string `mapstructure:"cors_allowed_origins" default:"http://localhost:3000,http://localhost:8080"`
	CORSAllowedMethods     []string `mapstructure:"cors_allowed_methods" default:"GET,POST,PUT,DELETE,OPTIONS"`
	CORSAllowedHeaders     []string `mapstructure:"cors_allowed_headers" default:"*"`
	CORSAllowCredentials   bool     `mapstructure:"cors_allow_credentials" default:"true"`

	// 🚦 Rate limiting (если потребуется)
	RateLimitEnabled bool          `mapstructure:"rate_limit_enabled" default:"false"`
	RateLimitRPS     int           `mapstructure:"rate_limit_rps" default:"100"`
	RateLimitBurst   int           `mapstructure:"rate_limit_burst" default:"200"`

	// ⏱️ Таймауты запросов
	RequestTimeout time.Duration `mapstructure:"request_timeout" default:"30s"`
}

// =====================================================================
// 📊 КОНФИГУРАЦИЯ МОНИТОРИНГА
// =====================================================================

// MonitoringConfig содержит настройки мониторинга и метрик
//
// TODO: Добавить настройки для Prometheus и Grafana
// TODO: Добавить настройки для health checks
type MonitoringConfig struct {
	// 📈 Метрики
	MetricsEnabled bool   `mapstructure:"metrics_enabled" default:"true"`
	MetricsPath    string `mapstructure:"metrics_path" default:"/metrics"`

	// 🔍 Health check
	HealthCheckPath string `mapstructure:"health_check_path" default:"/health"`

	// 📊 Настройки профилирования (для разработки)
	ProfilerEnabled bool   `mapstructure:"profiler_enabled" default:"false"`
	ProfilerPath    string `mapstructure:"profiler_path" default:"/debug/pprof"`
}

// =====================================================================
// 🔧 ФУНКЦИИ ЗАГРУЗКИ И ВАЛИДАЦИИ КОНФИГУРАЦИИ
// =====================================================================

// Load загружает конфигурацию из переменных окружения и файлов
//
// TODO: Добавить поддержку конфигурационных файлов (config.yml)
// TODO: Добавить hot-reload для некритичных настроек
func Load() (*Config, error) {
	// Настраиваем viper для загрузки конфигурации
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Настраиваем преобразование имен переменных
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Загружаем .env файл (опционально)
	if err := viper.ReadInConfig(); err != nil {
		// Не критично если .env файл не найден
		// Переменные окружения все равно будут загружены
	}

	// Устанавливаем дефолтные значения
	setDefaults()

	// Создаем структуру конфигурации
	var config Config

	// Загружаем конфигурацию из viper
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Валидируем конфигурацию
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// setDefaults устанавливает дефолтные значения для всех настроек
//
// TODO: Вынести дефолтные значения в отдельный файл
func setDefaults() {
	// 🖥️ Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.shutdown_timeout", "10s")
	viper.SetDefault("server.max_request_size", "10MB")
	viper.SetDefault("server.mode", "debug")

	// 🗃️ Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.dbname", "tender_automation")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "1h")
	viper.SetDefault("database.conn_max_idle_time", "30m")
	viper.SetDefault("database.connect_timeout", "10s")
	viper.SetDefault("database.query_timeout", "30s")

	// 🤖 AI defaults
	viper.SetDefault("ai.url", "http://localhost:11434")
	viper.SetDefault("ai.model", "llama2")
	viper.SetDefault("ai.timeout", "60s")
	viper.SetDefault("ai.max_retries", 3)
	viper.SetDefault("ai.retry_delay", "5s")
	viper.SetDefault("ai.relevance_threshold", 0.7)
	viper.SetDefault("ai.batch_size", 10)
	viper.SetDefault("ai.temperature", 0.1)
	viper.SetDefault("ai.max_tokens", 1000)

	// 🕷️ Scraping defaults
	viper.SetDefault("scraping.delay", "2s")
	viper.SetDefault("scraping.max_concurrent", 3)
	viper.SetDefault("scraping.user_agent", "Mozilla/5.0 (compatible; TenderBot/1.0)")
	viper.SetDefault("scraping.max_retries", 3)
	viper.SetDefault("scraping.retry_delay", "5s")
	viper.SetDefault("scraping.request_timeout", "30s")
	viper.SetDefault("scraping.page_timeout", "60s")
	viper.SetDefault("scraping.enabled_platforms", []string{"zakupki"})
	viper.SetDefault("scraping.batch_size", 50)
	viper.SetDefault("scraping.scan_interval", "1h")

	// 📝 Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.file_name", "app.log")
	viper.SetDefault("logging.max_size", 100)
	viper.SetDefault("logging.max_backups", 3)
	viper.SetDefault("logging.max_age", 28)
	viper.SetDefault("logging.enable_caller", true)
	viper.SetDefault("logging.enable_stacktrace", false)

	// 🎯 Business defaults
	viper.SetDefault("business.min_ai_score_threshold", 0.7)
	viper.SetDefault("business.analysis_batch_size", 10)
	viper.SetDefault("business.default_page_size", 20)
	viper.SetDefault("business.max_page_size", 100)
	viper.SetDefault("business.default_tender_status", "active")
	viper.SetDefault("business.tender_discovery_interval", "1h")
	viper.SetDefault("business.ai_analysis_interval", "30m")
	viper.SetDefault("business.cleanup_interval", "24h")

	// 🔒 Security defaults
	viper.SetDefault("security.cors_allowed_origins", []string{"http://localhost:3000", "http://localhost:8080"})
	viper.SetDefault("security.cors_allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("security.cors_allowed_headers", []string{"*"})
	viper.SetDefault("security.cors_allow_credentials", true)
	viper.SetDefault("security.rate_limit_enabled", false)
	viper.SetDefault("security.rate_limit_rps", 100)
	viper.SetDefault("security.rate_limit_burst", 200)
	viper.SetDefault("security.request_timeout", "30s")

	// 📊 Monitoring defaults
	viper.SetDefault("monitoring.metrics_enabled", true)
	viper.SetDefault("monitoring.metrics_path", "/metrics")
	viper.SetDefault("monitoring.health_check_path", "/health")
	viper.SetDefault("monitoring.profiler_enabled", false)
	viper.SetDefault("monitoring.profiler_path", "/debug/pprof")
}

// validateConfig валидирует загруженную конфигурацию
//
// TODO: Добавить кастомные валидаторы для специфичных правил
// TODO: Добавить проверку доступности внешних сервисов
func validateConfig(config *Config) error {
	validate := validator.New()

	// Валидируем структуру конфигурации
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Дополнительные проверки
	if config.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}

	if config.Business.DefaultPageSize > config.Business.MaxPageSize {
		return fmt.Errorf("default page size cannot be greater than max page size")
	}

	if config.AI.RelevanceThreshold < 0 || config.AI.RelevanceThreshold > 1 {
		return fmt.Errorf("AI relevance threshold must be between 0 and 1")
	}

	return nil
}

// =====================================================================
// 🔧 ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// =====================================================================

// GetEnv возвращает переменную окружения или дефолтное значение
func GetEnv(key, defaultValue string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

// IsProduction проверяет, запущено ли приложение в production режиме
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Server.Mode) == "release"
}

// IsDevelopment проверяет, запущено ли приложение в development режиме
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Server.Mode) == "debug"
}

// =====================================================================
// 📝 ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ КОНФИГУРАЦИИ
// =====================================================================

// TODO: Примеры использования в коде:
//
// 1. Загрузка конфигурации:
//    config, err := configs.Load()
//    if err != nil {
//        log.Fatal("Failed to load config:", err)
//    }
//
// 2. Использование в HTTP сервере:
//    server := &http.Server{
//        Addr:         config.Server.GetAddress(),
//        ReadTimeout:  config.Server.ReadTimeout,
//        WriteTimeout: config.Server.WriteTimeout,
//    }
//
// 3. Подключение к базе данных:
//    db, err := sql.Open("postgres", config.Database.GetDSN())
//
// 4. Настройка AI клиента:
//    aiClient := ai.NewLlamaClient(config.AI.URL, config.AI.Model)
//
// 5. Проверка режима работы:
//    if config.IsProduction() {
//        // Production specific logic
//    }

// =====================================================================
// ✅ ЗАКЛЮЧЕНИЕ
// =====================================================================

// Этот файл предоставляет полную конфигурацию для MVP системы.
// При расширении функциональности:
// 1. Добавляйте новые секции конфигурации
// 2. Обновляйте setDefaults() с новыми дефолтными значениями
// 3. Добавляйте валидацию для новых настроек
// 4. Документируйте все новые параметры
//
// Помните: конфигурация должна быть thread-safe и immutable после загрузки.