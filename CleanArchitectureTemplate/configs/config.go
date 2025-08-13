// ⚙️ УНИВЕРСАЛЬНЫЙ ШАБЛОН КОНФИГУРАЦИИ для Clean Architecture
//
// ✅ ВСЕГДА одинаково:
// - Структура Config с вложенными секциями
// - Функция LoadConfig() для загрузки из ENV
// - Валидация конфигурации
//
// 🔀 МЕНЯЕТСЯ в зависимости от проекта:
// - Конкретные поля конфигурации
// - Дополнительные секции (Redis, Kafka, AWS, etc.)
// - Значения по умолчанию

package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// 🏗️ ОСНОВНАЯ СТРУКТУРА КОНФИГУРАЦИИ
// ✅ ВСЕГДА имеет такую структуру, но поля могут изменяться
type Config struct {
	// 🔀 ОСНОВНЫЕ СЕКЦИИ (адаптируйте под ваш проект):
	
	App      AppConfig      `json:"app"`      // Общие настройки приложения
	Server   ServerConfig   `json:"server"`   // HTTP/Web сервер
	Database DatabaseConfig `json:"database"` // База данных
	Auth     AuthConfig     `json:"auth"`     // Аутентификация
	
	// 📝 ДОПОЛНИТЕЛЬНЫЕ СЕКЦИИ (раскомментируйте при необходимости):
	
	// Redis    RedisConfig    `json:"redis"`    // Кеширование
	// Email    EmailConfig    `json:"email"`    // Отправка email
	// AWS      AWSConfig      `json:"aws"`      // Amazon Web Services
	// Kafka    KafkaConfig    `json:"kafka"`    // Message broker
	// GRPC     GRPCConfig     `json:"grpc"`     // gRPC сервер
	// GraphQL  GraphQLConfig  `json:"graphql"`  // GraphQL сервер
	// Metrics  MetricsConfig  `json:"metrics"`  // Метрики и мониторинг
	// Logging  LoggingConfig  `json:"logging"`  // Логирование
	// Tracing  TracingConfig  `json:"tracing"`  // Трейсинг
}

// 📱 СЕКЦИЯ: Общие настройки приложения
// ✅ ВСЕГДА нужна базовая информация о приложении
type AppConfig struct {
	Name        string `json:"name"`         // 🔀 ИЗМЕНИТЕ: название приложения
	Version     string `json:"version"`      // 🔀 ИЗМЕНИТЕ: версия приложения
	Environment string `json:"environment"`  // dev, staging, production
	Debug       bool   `json:"debug"`        // Режим отладки
}

// 🌐 СЕКЦИЯ: HTTP/Web сервер
// ✅ ВСЕГДА нужна для веб-приложений
type ServerConfig struct {
	Port            string        `json:"port"`             // Порт сервера
	Host            string        `json:"host"`             // Хост сервера
	ReadTimeout     time.Duration `json:"read_timeout"`     // Таймаут чтения
	WriteTimeout    time.Duration `json:"write_timeout"`    // Таймаут записи
	ShutdownTimeout time.Duration `json:"shutdown_timeout"` // Таймаут graceful shutdown
	
	// 📝 ДОПОЛНИТЕЛЬНЫЕ НАСТРОЙКИ HTTP:
	MaxHeaderBytes int  `json:"max_header_bytes"` // Макс размер заголовков
	EnableCORS     bool `json:"enable_cors"`      // Включить CORS
	EnableHTTPS    bool `json:"enable_https"`     // Включить HTTPS
	
	// 📝 TLS НАСТРОЙКИ (если нужен HTTPS):
	CertFile string `json:"cert_file"` // Путь к сертификату
	KeyFile  string `json:"key_file"`  // Путь к приватному ключу
}

// 🗄️ СЕКЦИЯ: База данных
// ✅ ВСЕГДА нужна для приложений с БД
type DatabaseConfig struct {
	// 🔀 ОСНОВНЫЕ НАСТРОЙКИ (адаптируйте под вашу БД):
	Driver   string `json:"driver"`   // postgres, mysql, mongodb, sqlite
	Host     string `json:"host"`     // Хост БД
	Port     int    `json:"port"`     // Порт БД
	User     string `json:"user"`     // Пользователь БД
	Password string `json:"password"` // Пароль БД
	Name     string `json:"name"`     // Название БД
	SSLMode  string `json:"ssl_mode"` // SSL режим
	
	// 📝 НАСТРОЙКИ ПУЛА СОЕДИНЕНИЙ:
	MaxOpenConns    int           `json:"max_open_conns"`    // Макс открытых соединений
	MaxIdleConns    int           `json:"max_idle_conns"`    // Макс простаивающих соединений
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"` // Время жизни соединения
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"` // Время простоя соединения
	
	// 📝 ДОПОЛНИТЕЛЬНЫЕ НАСТРОЙКИ:
	MigrationsPath string `json:"migrations_path"` // Путь к миграциям
	AutoMigrate    bool   `json:"auto_migrate"`    // Автоматические миграции
}

// 🔐 СЕКЦИЯ: Аутентификация
// ✅ ВСЕГДА нужна для приложений с авторизацией
type AuthConfig struct {
	// JWT НАСТРОЙКИ:
	JWTSecret     string        `json:"jwt_secret"`      // Секретный ключ JWT
	JWTExpiry     time.Duration `json:"jwt_expiry"`      // Время жизни токена
	JWTIssuer     string        `json:"jwt_issuer"`      // Издатель токена
	RefreshExpiry time.Duration `json:"refresh_expiry"`  // Время жизни refresh токена
	
	// 📝 НАСТРОЙКИ ПАРОЛЕЙ:
	BcryptCost       int `json:"bcrypt_cost"`        // Стоимость bcrypt хеширования
	MinPasswordLen   int `json:"min_password_len"`   // Минимальная длина пароля
	RequireUppercase bool `json:"require_uppercase"` // Требовать заглавные буквы
	RequireNumbers   bool `json:"require_numbers"`   // Требовать цифры
	RequireSymbols   bool `json:"require_symbols"`   // Требовать символы
	
	// 📝 OAUTH НАСТРОЙКИ (если используется):
	GoogleClientID     string `json:"google_client_id"`     // Google OAuth
	GoogleClientSecret string `json:"google_client_secret"`
	FacebookAppID      string `json:"facebook_app_id"`      // Facebook OAuth
	FacebookAppSecret  string `json:"facebook_app_secret"`
}

// 📝 ДОПОЛНИТЕЛЬНЫЕ СЕКЦИИ КОНФИГУРАЦИИ
// Раскомментируйте и адаптируйте под ваш проект:

// // 🗄️ REDIS КОНФИГУРАЦИЯ
// type RedisConfig struct {
//     Host     string `json:"host"`     // Хост Redis
//     Port     int    `json:"port"`     // Порт Redis
//     Password string `json:"password"` // Пароль Redis
//     DB       int    `json:"db"`       // Номер БД Redis
//     
//     // Настройки пула:
//     PoolSize     int           `json:"pool_size"`     // Размер пула
//     MinIdleConns int           `json:"min_idle_conns"` // Мин простаивающих соединений
//     DialTimeout  time.Duration `json:"dial_timeout"`  // Таймаут подключения
//     ReadTimeout  time.Duration `json:"read_timeout"`  // Таймаут чтения
//     WriteTimeout time.Duration `json:"write_timeout"` // Таймаут записи
// }

// // 📧 EMAIL КОНФИГУРАЦИЯ
// type EmailConfig struct {
//     Provider string `json:"provider"` // smtp, sendgrid, ses, mailgun
//     
//     // SMTP настройки:
//     SMTPHost     string `json:"smtp_host"`     // SMTP хост
//     SMTPPort     int    `json:"smtp_port"`     // SMTP порт
//     SMTPUser     string `json:"smtp_user"`     // SMTP пользователь
//     SMTPPassword string `json:"smtp_password"` // SMTP пароль
//     SMTPFrom     string `json:"smtp_from"`     // Email отправителя
//     
//     // API ключи для сервисов:
//     SendGridAPIKey string `json:"sendgrid_api_key"` // SendGrid
//     MailgunAPIKey  string `json:"mailgun_api_key"`  // Mailgun
//     MailgunDomain  string `json:"mailgun_domain"`   // Mailgun домен
// }

// // ☁️ AWS КОНФИГУРАЦИЯ
// type AWSConfig struct {
//     Region          string `json:"region"`           // AWS регион
//     AccessKeyID     string `json:"access_key_id"`    // AWS Access Key
//     SecretAccessKey string `json:"secret_access_key"` // AWS Secret Key
//     
//     // S3 настройки:
//     S3Bucket string `json:"s3_bucket"` // S3 бакет
//     S3Region string `json:"s3_region"` // S3 регион
//     
//     // SQS настройки:
//     SQSQueueURL string `json:"sqs_queue_url"` // SQS очередь
// }

// // 📊 KAFKA КОНФИГУРАЦИЯ
// type KafkaConfig struct {
//     Brokers []string `json:"brokers"` // Список брокеров
//     
//     // Producer настройки:
//     ProducerRetries int           `json:"producer_retries"` // Повторы отправки
//     ProducerTimeout time.Duration `json:"producer_timeout"` // Таймаут отправки
//     
//     // Consumer настройки:
//     ConsumerGroup   string        `json:"consumer_group"`   // Группа потребителей
//     ConsumerTimeout time.Duration `json:"consumer_timeout"` // Таймаут потребления
//     
//     // Топики:
//     Topics map[string]string `json:"topics"` // Карта топиков
// }

// // 🔗 gRPC КОНФИГУРАЦИЯ
// type GRPCConfig struct {
//     Port    string        `json:"port"`    // Порт gRPC сервера
//     Timeout time.Duration `json:"timeout"` // Таймаут запросов
//     
//     // TLS настройки:
//     EnableTLS bool   `json:"enable_tls"` // Включить TLS
//     CertFile  string `json:"cert_file"`  // Сертификат
//     KeyFile   string `json:"key_file"`   // Приватный ключ
// }

// // 📈 МЕТРИКИ КОНФИГУРАЦИЯ
// type MetricsConfig struct {
//     Enabled bool   `json:"enabled"` // Включить метрики
//     Port    string `json:"port"`    // Порт metrics сервера
//     Path    string `json:"path"`    // Путь к метрикам
//     
//     // Prometheus настройки:
//     PrometheusEnabled bool `json:"prometheus_enabled"` // Включить Prometheus
// }

// // 📝 ЛОГИРОВАНИЕ КОНФИГУРАЦИЯ
// type LoggingConfig struct {
//     Level  string `json:"level"`  // debug, info, warn, error
//     Format string `json:"format"` // json, text
//     Output string `json:"output"` // stdout, file
//     
//     // Файловое логирование:
//     FilePath   string `json:"file_path"`   // Путь к файлу логов
//     MaxSize    int    `json:"max_size"`    // Макс размер файла (MB)
//     MaxBackups int    `json:"max_backups"` // Количество backup файлов
//     MaxAge     int    `json:"max_age"`     // Время хранения (дни)
// }

// // 🔍 ТРЕЙСИНГ КОНФИГУРАЦИЯ
// type TracingConfig struct {
//     Enabled     bool    `json:"enabled"`      // Включить трейсинг
//     ServiceName string  `json:"service_name"` // Название сервиса
//     JaegerURL   string  `json:"jaeger_url"`   // URL Jaeger
//     SampleRate  float64 `json:"sample_rate"`  // Частота сэмплирования
// }

// 🔧 ФУНКЦИЯ ЗАГРУЗКИ КОНФИГУРАЦИИ
// ✅ ВСЕГДА одинаковая структура, но ENV переменные могут отличаться
func LoadConfig() (*Config, error) {
	config := &Config{
		// 🔀 ЗНАЧЕНИЯ ПО УМОЛЧАНИЮ (адаптируйте под ваш проект):
		
		App: AppConfig{
			Name:        getEnv("APP_NAME", "Your Application"),        // 🔀 ИЗМЕНИТЕ
			Version:     getEnv("APP_VERSION", "1.0.0"),              // 🔀 ИЗМЕНИТЕ
			Environment: getEnv("APP_ENV", "development"),
			Debug:       getEnvBool("APP_DEBUG", true),
		},
		
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Host:            getEnv("SERVER_HOST", "localhost"),
			ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
			MaxHeaderBytes:  getEnvInt("SERVER_MAX_HEADER_BYTES", 1048576), // 1MB
			EnableCORS:      getEnvBool("SERVER_ENABLE_CORS", true),
			EnableHTTPS:     getEnvBool("SERVER_ENABLE_HTTPS", false),
			CertFile:        getEnv("SERVER_CERT_FILE", ""),
			KeyFile:         getEnv("SERVER_KEY_FILE", ""),
		},
		
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "postgres"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Name:            getEnv("DB_NAME", "your_app_db"),            // 🔀 ИЗМЕНИТЕ
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
			MigrationsPath:  getEnv("DB_MIGRATIONS_PATH", "migrations"),
			AutoMigrate:     getEnvBool("DB_AUTO_MIGRATE", false),
		},
		
		Auth: AuthConfig{
			JWTSecret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key"),  // 🔀 ИЗМЕНИТЕ в production!
			JWTExpiry:        getEnvDuration("JWT_EXPIRY", 24*time.Hour),
			JWTIssuer:        getEnv("JWT_ISSUER", "your-app"),                  // 🔀 ИЗМЕНИТЕ
			RefreshExpiry:    getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			BcryptCost:       getEnvInt("BCRYPT_COST", 12),
			MinPasswordLen:   getEnvInt("MIN_PASSWORD_LEN", 8),
			RequireUppercase: getEnvBool("REQUIRE_UPPERCASE", true),
			RequireNumbers:   getEnvBool("REQUIRE_NUMBERS", true),
			RequireSymbols:   getEnvBool("REQUIRE_SYMBOLS", false),
			
			// OAuth (если используется):
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			FacebookAppID:      getEnv("FACEBOOK_APP_ID", ""),
			FacebookAppSecret:  getEnv("FACEBOOK_APP_SECRET", ""),
		},
		
		// 📝 ДОПОЛНИТЕЛЬНЫЕ СЕКЦИИ (раскомментируйте при необходимости):
		
		// Redis: RedisConfig{
		//     Host:         getEnv("REDIS_HOST", "localhost"),
		//     Port:         getEnvInt("REDIS_PORT", 6379),
		//     Password:     getEnv("REDIS_PASSWORD", ""),
		//     DB:           getEnvInt("REDIS_DB", 0),
		//     PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),
		//     MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNS", 2),
		//     DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		//     ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		//     WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		// },
	}
	
	return config, nil
}

// ✅ ФУНКЦИЯ ВАЛИДАЦИИ КОНФИГУРАЦИИ
// Добавьте проверки критичных параметров
func (c *Config) Validate() error {
	// 🔀 ДОБАВЬТЕ валидацию под ваш проект:
	
	if c.App.Name == "" {
		return fmt.Errorf("app name is required")
	}
	
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	
	if c.Auth.JWTSecret == "" || c.Auth.JWTSecret == "your-super-secret-jwt-key" {
		return fmt.Errorf("JWT secret must be set and not use default value")
	}
	
	// 📝 ДОПОЛНИТЕЛЬНЫЕ ВАЛИДАЦИИ:
	// if c.Redis.Host == "" {
	//     return fmt.Errorf("redis host is required")
	// }
	
	return nil
}

// 🔧 ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ✅ ВСЕГДА одинаковые, используйте как есть

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// 📚 ИНСТРУКЦИЯ ПО АДАПТАЦИИ:
//
// 1. 🔀 ИЗМЕНИТЕ значения по умолчанию:
//    - APP_NAME, APP_VERSION
//    - DB_NAME, JWT_ISSUER
//    - JWT_SECRET (обязательно в production!)
//
// 2. 📝 РАСКОММЕНТИРУЙТЕ нужные секции:
//    - Redis, Email, AWS, Kafka, etc.
//    - Добавьте соответствующие ENV переменные
//
// 3. ✅ ДОБАВЬТЕ валидацию:
//    - Обязательные поля в Validate()
//    - Проверки форматов и диапазонов
//
// 4. 🔀 АДАПТИРУЙТЕ ENV переменные:
//    - Добавьте специфичные для вашего проекта
//    - Измените префиксы если нужно
//
// 5. 📝 СОЗДАЙТЕ .env файл:
//    - Скопируйте все ENV переменные
//    - Установите нужные значения