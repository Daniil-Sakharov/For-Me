package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Auth     AuthConfig     `json:"auth"`
	App      AppConfig      `json:"app"`
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode"`
}

// ServerConfig конфигурация HTTP сервера
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// AuthConfig конфигурация аутентификации
type AuthConfig struct {
	JWTSecret    string        `json:"jwt_secret"`
	JWTIssuer    string        `json:"jwt_issuer"`
	JWTExpiry    time.Duration `json:"jwt_expiry"`
	BcryptCost   int           `json:"bcrypt_cost"`
}

// AppConfig общие настройки приложения
type AppConfig struct {
	BaseURL     string `json:"base_url"`
	Environment string `json:"environment"` // development, staging, production
	LogLevel    string `json:"log_level"`   // debug, info, warn, error
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "url_shortener"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Auth: AuthConfig{
			JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			JWTIssuer:  getEnv("JWT_ISSUER", "url-shortener"),
			JWTExpiry:  getEnvDuration("JWT_EXPIRY", 24*time.Hour),
			BcryptCost: getEnvInt("BCRYPT_COST", 12),
		},
		App: AppConfig{
			BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
	}

	// Валидируем конфигурацию
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Проверяем обязательные поля
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Auth.JWTSecret == "" || c.Auth.JWTSecret == "your-secret-key-change-in-production" {
		return fmt.Errorf("JWT secret must be set and changed from default")
	}
	if c.App.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}

	// Проверяем допустимые значения
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Auth.BcryptCost < 4 || c.Auth.BcryptCost > 31 {
		return fmt.Errorf("invalid bcrypt cost: %d (must be between 4 and 31)", c.Auth.BcryptCost)
	}

	return nil
}

// IsDevelopment проверяет, запущено ли приложение в режиме разработки
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction проверяет, запущено ли приложение в продакшене
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// УТИЛИТЫ ДЛЯ РАБОТЫ С ПЕРЕМЕННЫМИ ОКРУЖЕНИЯ

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt возвращает целочисленное значение переменной окружения
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDuration возвращает duration из переменной окружения
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// ПРИНЦИПЫ КОНФИГУРАЦИИ:
// 1. Конфигурация через переменные окружения (12-factor app)
// 2. Значения по умолчанию для удобства разработки
// 3. Валидация конфигурации при запуске
// 4. Разделение конфигурации по доменам (DB, Server, Auth)
// 5. Типобезопасность (int для портов, duration для таймаутов)