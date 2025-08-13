// 🤖 МОДУЛЬ ДЛЯ СИСТЕМЫ АВТОМАТИЗАЦИИ ТЕНДЕРОВ С ИИ
//
// Специализированная система для:
// - Поиска и анализа тендеров на медицинское оборудование  
// - Автоматической обработки документов с ИИ
// - Интеграции с поставщиками через email
// - Оптимизации цен на основе исторических данных

module tender-automation-system

go 1.21

require (
	// ✅ БАЗОВЫЕ ЗАВИСИМОСТИ:
	github.com/go-playground/validator/v10 v10.16.0 // Валидация данных
	golang.org/x/crypto v0.17.0                     // Безопасность
	
	// 🗄️ БАЗЫ ДАННЫХ:
	github.com/lib/pq v1.10.9                       // PostgreSQL для основных данных
	go.mongodb.org/mongo-driver v1.13.1             // MongoDB для документов и неструктурированных данных
	github.com/redis/go-redis/v9 v9.3.0             // Redis для кеширования и очередей
	
	// 🕷️ WEB SCRAPING И ПАРСИНГ:
	github.com/gocolly/colly/v2 v2.1.0              // Мощный фреймворк для скрапинга
	github.com/PuerkitoBio/goquery v1.8.1           // jQuery-like селекторы для HTML
	github.com/chromedp/chromedp v0.9.3             // Headless Chrome для динамических сайтов
	github.com/playwright-community/playwright-go v0.4000.0 // Альтернатива для сложных сайтов
	
	// 📄 ОБРАБОТКА ДОКУМЕНТОВ:
	github.com/unidoc/unioffice v1.28.0             // Работа с DOC/DOCX файлами
	github.com/ledongthuc/pdf v0.0.0-20220302134840-0c2507a12d80 // PDF обработка
	github.com/tealeg/xlsx/v3 v3.3.2                // Excel файлы
	
	// 🤖 ИИ И МАШИННОЕ ОБУЧЕНИЕ:
	github.com/sashabaranov/go-openai v1.17.9       // OpenAI API (альтернатива для Llama)
	github.com/ggerganov/whisper.cpp/bindings/go v0.0.0-20231117122806-ed6bf6543fd9 // Local AI
	github.com/ollama/ollama v0.1.17                // Локальные LLM модели (для Llama)
	
	// 📧 EMAIL ОБРАБОТКА:
	github.com/go-mail/mail/v2 v2.3.0               // Отправка писем
	github.com/emersion/go-imap v1.2.1              // Чтение входящих писем
	github.com/jhillyerd/enmime v0.10.1             // Парсинг email сообщений
	
	// 🌐 HTTP И API:
	github.com/gin-gonic/gin v1.9.1                 // Web framework для API
	github.com/gorilla/websocket v1.5.1             // WebSocket для real-time уведомлений
	
	// 📊 ОЧЕРЕДИ И ФОНОВЫЕ ЗАДАЧИ:
	github.com/hibiken/asynq v0.24.1                // Redis-based task queue
	github.com/robfig/cron/v3 v3.0.1                // Cron jobs для расписания
	
	// 🔍 ПОИСК И ИНДЕКСАЦИЯ:
	github.com/elastic/go-elasticsearch/v8 v8.11.1  // Elasticsearch для поиска по документам
	github.com/blevesearch/bleve/v2 v2.3.10         // Встроенный full-text search
	
	// 📈 АНАЛИТИКА И МЕТРИКИ:
	github.com/prometheus/client_golang v1.17.0     // Метрики
	go.uber.org/zap v1.26.0                         // Логирование
	
	// 🔧 УТИЛИТЫ:
	github.com/spf13/cobra v1.8.0                   // CLI интерфейс
	github.com/spf13/viper v1.18.2                  // Конфигурация
	github.com/joho/godotenv v1.4.0                 // .env файлы
	
	// 🧪 ТЕСТИРОВАНИЕ:
	github.com/stretchr/testify v1.8.4              // Testing framework
	github.com/golang/mock v1.6.0                   // Mocking
	
	// 🛡️ БЕЗОПАСНОСТЬ:
	github.com/golang-jwt/jwt/v5 v5.2.0             // JWT токены
	golang.org/x/time v0.5.0                       // Rate limiting
	
	// 📂 ФАЙЛОВАЯ СИСТЕМА:
	github.com/aws/aws-sdk-go-v2 v1.24.0           // AWS S3 для хранения файлов
	github.com/minio/minio-go/v7 v7.0.66           // MinIO для локального хранения
	
	// 🔄 АВТОМАТИЗАЦИЯ:
	github.com/fsnotify/fsnotify v1.7.0             // Мониторинг файлов
	github.com/kardianos/service v1.2.2             // Системный сервис
)

// 📝 СПЕЦИФИЧНЫЕ ЗАВИСИМОСТИ ДЛЯ ТЕНДЕРНОЙ СИСТЕМЫ:
//
// 🕷️ Scraping Stack:
// - colly: основной scraper для статических сайтов
// - chromedp: для сайтов с JavaScript (zakupki.gov.ru)
// - playwright: бэкап для сложных сайтов
//
// 📄 Document Processing:
// - unioffice: для чтения технических заданий (DOC/DOCX)
// - pdf: для PDF документов
// - xlsx: для таблиц с ценами
//
// 🤖 AI Integration:
// - ollama: для запуска Llama модели локально
// - go-openai: совместимый API интерфейс
//
// 📧 Email Automation:
// - go-mail: отправка коммерческих предложений
// - go-imap: чтение ответов от поставщиков
//
// 🔍 Search & Analytics:
// - elasticsearch: индексация товаров и документов
// - bleve: встроенный поиск как fallback
//
// 📊 Background Processing:
// - asynq: очереди для обработки тендеров
// - cron: расписание автоматических проверок
//
// 🔧 ИНСТРУКЦИЯ ПО НАСТРОЙКЕ:
//
// 1. Установите зависимости:
//    go mod tidy
//
// 2. Настройте внешние сервисы:
//    - PostgreSQL для основных данных
//    - MongoDB для документов  
//    - Redis для очередей
//    - Elasticsearch для поиска (опционально)
//
// 3. Настройте AI модель:
//    - Скачайте Llama 4 через Ollama
//    - Или настройте API ключи для облачных сервисов
//
// 4. Настройте email:
//    - SMTP сервер для отправки
//    - IMAP для чтения ответов
//
// 5. Настройте хранилище файлов:
//    - AWS S3 или MinIO для документов тендеров