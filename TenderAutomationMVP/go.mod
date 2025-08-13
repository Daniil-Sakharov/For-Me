module tender-automation-mvp

go 1.21

require (
	// 🌐 Web Framework
	github.com/gin-gonic/gin v1.9.1

	// 🗃️ Database
	github.com/lib/pq v1.10.9 // PostgreSQL driver

	// ⚙️ Configuration
	github.com/spf13/viper v1.18.2   // Configuration management
	github.com/joho/godotenv v1.5.1  // .env file support

	// 📝 Logging
	go.uber.org/zap v1.26.0 // Structured logging

	// ✅ Validation
	github.com/go-playground/validator/v10 v10.16.0

	// 🕷️ Web Scraping
	github.com/gocolly/colly/v2 v2.1.0

	// 🤖 AI Integration (Ollama client)
	// TODO: Добавить Ollama Go client когда будет доступен
	// github.com/ollama/ollama v0.1.17

	// 🧪 Testing
	github.com/stretchr/testify v1.8.4

	// 🔧 Utilities
	github.com/google/uuid v1.5.0 // UUID generation
)

// TODO: Добавить дополнительные зависимости по мере необходимости:
//
// Для расширения функциональности:
// - github.com/robfig/cron/v3 v3.0.1               // Cron scheduler
// - github.com/hibiken/asynq v0.24.1               // Background jobs
// - github.com/prometheus/client_golang v1.18.0    // Metrics
// - github.com/unidoc/unioffice v1.28.0            // Document processing
// - gopkg.in/mail.v2 v2.3.1                       // Email
// - github.com/redis/go-redis/v9 v9.3.0            // Redis client
//
// Для тестирования:
// - github.com/golang-migrate/migrate/v4           // Database migrations
// - github.com/testcontainers/testcontainers-go    // Integration testing
//
// Примечания:
// 1. Используем минимальный набор зависимостей для MVP
// 2. Все зависимости должны быть совместимы с Go 1.21+
// 3. Предпочитаем стабильные, хорошо поддерживаемые библиотеки
// 4. Избегаем vendor lock-in - выбираем библиотеки с интерфейсами