module tender-automation

go 1.21

require (
	// TODO: Web Framework (выберите один)
	github.com/gin-gonic/gin v1.9.1          // Популярный Go веб-фреймворк
	// github.com/labstack/echo/v4 v4.11.3    // Альтернатива Gin
	// github.com/gofiber/fiber/v2 v2.52.0    // Быстрый веб-фреймворк

	// TODO: База данных
	github.com/lib/pq v1.10.9                      // PostgreSQL драйвер
	go.mongodb.org/mongo-driver v1.13.1            // MongoDB драйвер
	github.com/go-redis/redis/v8 v8.11.5           // Redis client

	// TODO: ORM (опционально, предпочтительно plain SQL)
	// github.com/jmoiron/sqlx v1.3.5              // SQL расширения
	// gorm.io/gorm v1.25.5                        // ORM (не рекомендуется для Clean Architecture)

	// TODO: Конфигурация
	github.com/spf13/viper v1.18.2                 // Загрузка конфигурации
	github.com/joho/godotenv v1.5.1                // Загрузка .env файлов

	// TODO: Логирование
	go.uber.org/zap v1.26.0                        // Structured logging
	github.com/sirupsen/logrus v1.9.3              // Альтернатива Zap

	// TODO: Валидация
	github.com/go-playground/validator/v10 v10.16.0

	// TODO: Web Scraping
	github.com/gocolly/colly/v2 v2.1.0             // Web scraping framework
	github.com/chromedp/chromedp v0.9.3            // Headless Chrome automation
	github.com/playwright-community/playwright-go v0.4001.0 // Альтернатива ChromeDP

	// TODO: Document Processing
	github.com/unidoc/unioffice v1.28.0            // DOC/DOCX обработка
	github.com/ledongthuc/pdf v0.0.0-20220302134840-0c2394612674 // PDF обработка
	github.com/xuri/excelize/v2 v2.8.0             // Excel обработка
	github.com/PuerkitoBio/goquery v1.8.1          // HTML parsing

	// TODO: AI/ML Integration  
	// Ollama client для Llama 4 Maviric
	github.com/ollama/ollama v0.1.17               // Ollama Go client
	// github.com/sashabaranov/go-openai v1.17.9   // OpenAI API (альтернатива)

	// TODO: Email
	gopkg.in/mail.v2 v2.3.1                       // SMTP отправка
	github.com/emersion/go-imap v1.2.1             // IMAP для чтения писем
	github.com/jhillyerd/enmime v0.10.1            // Email parsing

	// TODO: Background Jobs & Scheduling
	github.com/hibiken/asynq v0.24.1               // Redis-based task queue
	github.com/robfig/cron/v3 v3.0.1               // Cron scheduler

	// TODO: Поиск и индексация
	github.com/olivere/elastic/v7 v7.0.32          // Elasticsearch client
	github.com/blevesearch/bleve/v2 v2.3.10        // Локальный поисковый движок

	// TODO: CLI
	github.com/spf13/cobra v1.8.0                  // CLI framework

	// TODO: Мониторинг
	github.com/prometheus/client_golang v1.18.0    // Prometheus metrics

	// TODO: HTTP клиент
	github.com/go-resty/resty/v2 v2.11.0           // HTTP client с retry и middleware

	// TODO: Утилиты
	github.com/google/uuid v1.5.0                  // UUID генерация
	github.com/stretchr/testify v1.8.4             // Testing utilities
	golang.org/x/crypto v0.17.0                    // Криптографические утилиты
	golang.org/x/time v0.5.0                       // Rate limiting

	// TODO: File Storage (опционально)
	// github.com/aws/aws-sdk-go v1.49.0           // AWS S3
	// github.com/minio/minio-go/v7 v7.0.66        // MinIO client

	// TODO: Message Queue (альтернативы Redis)
	// github.com/nats-io/nats.go v1.31.0          // NATS
	// github.com/rabbitmq/amqp091-go v1.9.0       // RabbitMQ
)

// TODO: Replace directives (если используете локальные модули)
// replace (
//     tender-automation/internal => ./internal
//     tender-automation/pkg => ./pkg
// )