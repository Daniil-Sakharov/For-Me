package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// TODO: Основная структура конфигурации приложения
// Содержит все настройки для:
// - Подключений к базам данных
// - Интеграции с внешними сервисами
// - Параметров работы пайплайнов
// - Настроек ИИ и email

type Config struct {
	// TODO: Основные настройки приложения
	App AppConfig `mapstructure:"app"`

	// TODO: Настройки баз данных
	Database DatabaseConfig `mapstructure:"database"`

	// TODO: Настройки ИИ (Llama 4 Maviric)
	AI AIConfig `mapstructure:"ai"`

	// TODO: Настройки web scraping
	Scraping ScrapingConfig `mapstructure:"scraping"`

	// TODO: Настройки email автоматизации
	Email EmailConfig `mapstructure:"email"`

	// TODO: Настройки обработки документов
	DocumentProcessing DocumentProcessingConfig `mapstructure:"document_processing"`

	// TODO: Настройки очередей и фоновых задач
	Queue QueueConfig `mapstructure:"queue"`

	// TODO: Настройки мониторинга
	Monitoring MonitoringConfig `mapstructure:"monitoring"`

	// TODO: Настройки поиска
	Search SearchConfig `mapstructure:"search"`

	// TODO: Настройки файлового хранилища
	Storage StorageConfig `mapstructure:"storage"`
}

type AppConfig struct {
	// TODO: Базовые настройки приложения
	Name        string `mapstructure:"name" default:"Tender Automation System"`
	Version     string `mapstructure:"version" default:"1.0.0"`
	Environment string `mapstructure:"environment" default:"development"` // development, production
	Port        int    `mapstructure:"port" default:"8080"`
	Host        string `mapstructure:"host" default:"0.0.0.0"`

	// TODO: Настройки безопасности
	JWTSecret      string        `mapstructure:"jwt_secret"`
	SessionTimeout time.Duration `mapstructure:"session_timeout" default:"24h"`
}

type DatabaseConfig struct {
	// TODO: PostgreSQL для основных данных
	Postgres PostgresConfig `mapstructure:"postgres"`

	// TODO: MongoDB для документов и неструктурированных данных
	MongoDB MongoDBConfig `mapstructure:"mongodb"`

	// TODO: Redis для кеширования и очередей
	Redis RedisConfig `mapstructure:"redis"`
}

type PostgresConfig struct {
	// TODO: Настройки подключения к PostgreSQL
	Host            string        `mapstructure:"host" default:"localhost"`
	Port            int           `mapstructure:"port" default:"5432"`
	User            string        `mapstructure:"user" default:"postgres"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname" default:"tender_automation"`
	SSLMode         string        `mapstructure:"sslmode" default:"disable"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" default:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" default:"5"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" default:"1h"`
}

type MongoDBConfig struct {
	// TODO: Настройки MongoDB для хранения документов
	URI      string `mapstructure:"uri" default:"mongodb://localhost:27017"`
	Database string `mapstructure:"database" default:"tender_documents"`
	Timeout  time.Duration `mapstructure:"timeout" default:"30s"`
}

type RedisConfig struct {
	// TODO: Настройки Redis для кеширования и очередей
	Host     string `mapstructure:"host" default:"localhost"`
	Port     int    `mapstructure:"port" default:"6379"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db" default:"0"`
	PoolSize int    `mapstructure:"pool_size" default:"10"`
}

type AIConfig struct {
	// TODO: Настройки ИИ провайдера (Llama 4 Maviric через Ollama)
	Provider string `mapstructure:"provider" default:"ollama"` // ollama, openai, anthropic

	// TODO: Ollama настройки
	Ollama OllamaConfig `mapstructure:"ollama"`

	// TODO: OpenAI настройки (альтернатива)
	OpenAI OpenAIConfig `mapstructure:"openai"`

	// TODO: Общие настройки
	MaxTokens         int           `mapstructure:"max_tokens" default:"4000"`
	Temperature       float32       `mapstructure:"temperature" default:"0.1"`
	RequestTimeout    time.Duration `mapstructure:"request_timeout" default:"5m"`
	MaxRetries        int           `mapstructure:"max_retries" default:"3"`
	ConcurrentRequests int          `mapstructure:"concurrent_requests" default:"5"`
}

type OllamaConfig struct {
	// TODO: Настройки для локального Ollama сервера
	Host  string `mapstructure:"host" default:"http://localhost:11434"`
	Model string `mapstructure:"model" default:"llama4-maviric:latest"`
}

type OpenAIConfig struct {
	// TODO: Настройки OpenAI API (если не используется локальная модель)
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model" default:"gpt-4"`
}

type ScrapingConfig struct {
	// TODO: Настройки web scraping для поиска тендеров
	UserAgent       string        `mapstructure:"user_agent" default:"TenderBot/1.0"`
	RequestDelay    time.Duration `mapstructure:"request_delay" default:"2s"`
	Timeout         time.Duration `mapstructure:"timeout" default:"30s"`
	MaxConcurrency  int           `mapstructure:"max_concurrency" default:"5"`
	RetryCount      int           `mapstructure:"retry_count" default:"3"`
	UseHeadless     bool          `mapstructure:"use_headless" default:"true"`
	ProxyRotation   bool          `mapstructure:"proxy_rotation" default:"false"`

	// TODO: Настройки для конкретных платформ
	Platforms map[string]PlatformConfig `mapstructure:"platforms"`
}

type PlatformConfig struct {
	// TODO: Конфигурация для конкретной закупочной платформы
	Enabled     bool          `mapstructure:"enabled" default:"true"`
	BaseURL     string        `mapstructure:"base_url"`
	SearchURL   string        `mapstructure:"search_url"`
	RateLimit   time.Duration `mapstructure:"rate_limit" default:"5s"`
	MaxPages    int           `mapstructure:"max_pages" default:"10"`
	Keywords    []string      `mapstructure:"keywords"`
}

type EmailConfig struct {
	// TODO: Настройки email автоматизации
	SMTP SMTPConfig `mapstructure:"smtp"`
	IMAP IMAPConfig `mapstructure:"imap"`

	// TODO: Общие настройки email
	FromEmail       string        `mapstructure:"from_email"`
	FromName        string        `mapstructure:"from_name" default:"Tender Automation System"`
	ReplyTo         string        `mapstructure:"reply_to"`
	MaxEmailsPerDay int           `mapstructure:"max_emails_per_day" default:"100"`
	SendDelay       time.Duration `mapstructure:"send_delay" default:"1m"`
	
	// TODO: Настройки ИИ для генерации писем
	AITemplates bool `mapstructure:"ai_templates" default:"true"`
}

type SMTPConfig struct {
	// TODO: SMTP настройки для отправки писем
	Host     string `mapstructure:"host" default:"smtp.gmail.com"`
	Port     int    `mapstructure:"port" default:"587"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	UseTLS   bool   `mapstructure:"use_tls" default:"true"`
}

type IMAPConfig struct {
	// TODO: IMAP настройки для чтения ответов
	Host     string `mapstructure:"host" default:"imap.gmail.com"`
	Port     int    `mapstructure:"port" default:"993"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	UseTLS   bool   `mapstructure:"use_tls" default:"true"`
}

type DocumentProcessingConfig struct {
	// TODO: Настройки обработки документов
	TempDir           string        `mapstructure:"temp_dir" default:"/tmp/tender_docs"`
	MaxFileSize       int64         `mapstructure:"max_file_size" default:"52428800"` // 50MB
	ProcessingTimeout time.Duration `mapstructure:"processing_timeout" default:"10m"`
	PreserveOriginal  bool          `mapstructure:"preserve_original" default:"true"`
	
	// TODO: Настройки извлечения текста
	OCREnabled       bool `mapstructure:"ocr_enabled" default:"true"`
	ExtractTables    bool `mapstructure:"extract_tables" default:"true"`
	ExtractImages    bool `mapstructure:"extract_images" default:"false"`
	
	// TODO: Поддерживаемые форматы
	SupportedFormats []string `mapstructure:"supported_formats" default:"pdf,doc,docx,xlsx,xls,rtf,txt"`
}

type QueueConfig struct {
	// TODO: Настройки очередей задач (Asynq + Redis)
	RedisAddr     string        `mapstructure:"redis_addr" default:"localhost:6379"`
	Concurrency   int           `mapstructure:"concurrency" default:"10"`
	RetryDelay    time.Duration `mapstructure:"retry_delay" default:"1m"`
	MaxRetries    int           `mapstructure:"max_retries" default:"5"`
	
	// TODO: Настройки воркеров
	Workers WorkersConfig `mapstructure:"workers"`
}

type WorkersConfig struct {
	// TODO: Количество воркеров для разных типов задач
	ScrapingWorkers           int `mapstructure:"scraping_workers" default:"3"`
	DocumentProcessingWorkers int `mapstructure:"document_processing_workers" default:"5"`
	AIAnalysisWorkers        int `mapstructure:"ai_analysis_workers" default:"2"`
	EmailWorkers             int `mapstructure:"email_workers" default:"2"`
	PriceAnalysisWorkers     int `mapstructure:"price_analysis_workers" default:"1"`
}

type MonitoringConfig struct {
	// TODO: Настройки мониторинга и метрик
	Enabled     bool   `mapstructure:"enabled" default:"true"`
	MetricsPort int    `mapstructure:"metrics_port" default:"9090"`
	MetricsPath string `mapstructure:"metrics_path" default:"/metrics"`
	
	// TODO: Alerting
	SlackWebhook string `mapstructure:"slack_webhook"`
	EmailAlerts  bool   `mapstructure:"email_alerts" default:"true"`
	
	// TODO: Логирование
	LogLevel  string `mapstructure:"log_level" default:"info"`
	LogFormat string `mapstructure:"log_format" default:"json"` // json, text
}

type SearchConfig struct {
	// TODO: Настройки поискового движка
	Provider string `mapstructure:"provider" default:"bleve"` // elasticsearch, bleve
	
	// TODO: Elasticsearch настройки
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	
	// TODO: Bleve настройки (встроенный поиск)
	Bleve BleveConfig `mapstructure:"bleve"`
}

type ElasticsearchConfig struct {
	// TODO: Настройки Elasticsearch
	URLs     []string `mapstructure:"urls" default:"http://localhost:9200"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	Index    string   `mapstructure:"index" default:"tender_data"`
}

type BleveConfig struct {
	// TODO: Настройки встроенного поиска Bleve
	IndexPath string `mapstructure:"index_path" default:"./data/search_index"`
}

type StorageConfig struct {
	// TODO: Настройки файлового хранилища
	Provider string `mapstructure:"provider" default:"local"` // local, s3, minio
	
	// TODO: Локальное хранилище
	Local LocalStorageConfig `mapstructure:"local"`
	
	// TODO: AWS S3
	S3 S3Config `mapstructure:"s3"`
	
	// TODO: MinIO
	MinIO MinIOConfig `mapstructure:"minio"`
}

type LocalStorageConfig struct {
	// TODO: Локальное файловое хранилище
	BasePath string `mapstructure:"base_path" default:"./data/files"`
}

type S3Config struct {
	// TODO: AWS S3 настройки
	Region    string `mapstructure:"region" default:"us-east-1"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type MinIOConfig struct {
	// TODO: MinIO настройки
	Endpoint  string `mapstructure:"endpoint" default:"localhost:9000"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket" default:"tender-documents"`
	UseSSL    bool   `mapstructure:"use_ssl" default:"false"`
}

// TODO: Функция загрузки конфигурации
// Использует Viper для чтения из:
// - .env файлов
// - переменных окружения
// - YAML/JSON конфигов
// - значений по умолчанию
func Load() (*Config, error) {
	// TODO: Настройка Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/tender-automation")

	// TODO: Автоматическое чтение переменных окружения
	viper.AutomaticEnv()
	viper.SetEnvPrefix("TENDER")

	// TODO: Загрузка .env файла если существует
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read .env file: %w", err)
		}
	}

	// TODO: Попытка прочитать основной конфиг файл
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Конфиг файл не найден - продолжаем с переменными окружения
	}

	// TODO: Маппинг в структуру
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// TODO: Валидация конфигурации
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// TODO: Инициализация platform конфигураций по умолчанию
	initDefaultPlatforms(&cfg)

	return &cfg, nil
}

// TODO: Валидация конфигурации
func validateConfig(cfg *Config) error {
	// TODO: Проверка обязательных полей
	if cfg.Database.Postgres.Password == "" {
		return fmt.Errorf("postgres password is required")
	}

	if cfg.AI.Provider == "ollama" && cfg.AI.Ollama.Host == "" {
		return fmt.Errorf("ollama host is required when using ollama provider")
	}

	if cfg.AI.Provider == "openai" && cfg.AI.OpenAI.APIKey == "" {
		return fmt.Errorf("openai api key is required when using openai provider")
	}

	if cfg.Email.SMTP.Username == "" || cfg.Email.SMTP.Password == "" {
		return fmt.Errorf("smtp credentials are required")
	}

	// TODO: Проверка диапазонов значений
	if cfg.App.Port < 1 || cfg.App.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", cfg.App.Port)
	}

	return nil
}

// TODO: Инициализация конфигураций платформ по умолчанию
func initDefaultPlatforms(cfg *Config) {
	if cfg.Scraping.Platforms == nil {
		cfg.Scraping.Platforms = make(map[string]PlatformConfig)
	}

	// TODO: Настройки для zakupki.gov.ru
	if _, exists := cfg.Scraping.Platforms["zakupki"]; !exists {
		cfg.Scraping.Platforms["zakupki"] = PlatformConfig{
			Enabled:   true,
			BaseURL:   "https://zakupki.gov.ru",
			SearchURL: "https://zakupki.gov.ru/epz/order/extendedsearch/results.html",
			RateLimit: 5 * time.Second,
			MaxPages:  10,
			Keywords: []string{
				"медицинское оборудование",
				"диагностическое оборудование",
				"лабораторное оборудование",
			},
		}
	}

	// TODO: Настройки для szvo.gov35.ru
	if _, exists := cfg.Scraping.Platforms["szvo"]; !exists {
		cfg.Scraping.Platforms["szvo"] = PlatformConfig{
			Enabled:   true,
			BaseURL:   "https://szvo.gov35.ru",
			SearchURL: "https://szvo.gov35.ru/purchases",
			RateLimit: 3 * time.Second,
			MaxPages:  5,
			Keywords: []string{
				"медицинское оборудование",
				"диагностическое оборудование",
			},
		}
	}

	// TODO: Настройки для gz-spb.ru
	if _, exists := cfg.Scraping.Platforms["spb"]; !exists {
		cfg.Scraping.Platforms["spb"] = PlatformConfig{
			Enabled:   true,
			BaseURL:   "https://gz-spb.ru",
			SearchURL: "https://gz-spb.ru/search",
			RateLimit: 7 * time.Second,
			MaxPages:  3,
			Keywords: []string{
				"медицинское оборудование",
			},
		}
	}
}

// TODO: Получение строки подключения к PostgreSQL
func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// TODO: Проверка development окружения
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// TODO: Проверка production окружения
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}