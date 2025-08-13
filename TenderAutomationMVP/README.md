# 🚀 Tender Automation MVP - Руководство по реализации

**Минимальная рабочая версия системы автоматизации тендеров с возможностью расширения**

## 📋 Содержание

1. [Обзор MVP](#обзор-mvp)
2. [Архитектура и принципы](#архитектура-и-принципы)
3. [Структура проекта](#структура-проекта)
4. [Детальное руководство по реализации](#детальное-руководство-по-реализации)
5. [Этапы разработки](#этапы-разработки)
6. [Расширение функциональности](#расширение-функциональности)
7. [Deployment и DevOps](#deployment-и-devops)

## 🎯 Обзор MVP

### Что делает MVP:
1. **🔍 Поиск тендеров** - парсинг zakupki.gov.ru
2. **🤖 AI анализ** - базовая классификация с Llama
3. **💾 Хранение данных** - PostgreSQL для тендеров
4. **🌐 REST API** - CRUD операции и управление
5. **📊 Базовая аналитика** - статистика и метрики

### Что НЕ входит в MVP (для расширения):
- Обработка документов
- Email автоматизация
- Множественные площадки
- Продвинутая аналитика
- Web UI

## 🏗️ Архитектура и принципы

### Clean Architecture Layers:

```
┌─────────────────────────────────────┐
│            🔌 INTERFACES            │  ← HTTP API, CLI
├─────────────────────────────────────┤
│            💼 USE CASES             │  ← Business Logic
├─────────────────────────────────────┤
│            🏛️ DOMAIN               │  ← Entities & Rules
├─────────────────────────────────────┤
│          🌐 INFRASTRUCTURE          │  ← DB, AI, Scraping
└─────────────────────────────────────┘
```

### Ключевые принципы:

1. **🔄 Dependency Inversion** - зависимости направлены внутрь
2. **🔗 Interface Segregation** - маленькие, специфичные интерфейсы
3. **🧩 Single Responsibility** - каждый модуль одну задачу
4. **🚀 Extensibility** - легкое добавление новой функциональности
5. **🧪 Testability** - все компоненты тестируемы

## 📁 Структура проекта

```
TenderAutomationMVP/
├── 📚 README.md                     # Это файл - полная документация
├── 🔧 go.mod                        # Go модуль с зависимостями
├── ⚙️ .env.example                  # Пример конфигурации
├── 🗃️ migrations/                   # Миграции базы данных
│   ├── 001_create_tenders.up.sql
│   └── 001_create_tenders.down.sql
├── 🚀 cmd/                          # Точки входа приложения
│   ├── api/                         # HTTP API сервер
│   │   └── main.go
│   ├── scraper/                     # Standalone scraper
│   │   └── main.go
│   └── cli/                         # CLI утилиты
│       └── main.go
├── ⚙️ configs/                      # Конфигурация
│   └── config.go
├── 🏛️ internal/                     # Основная логика приложения
│   ├── domain/                      # 🏛️ ДОМЕННЫЙ СЛОЙ
│   │   ├── tender/                  # Доменная модель тендера
│   │   │   ├── entity.go            # Сущность тендера
│   │   │   ├── repository.go        # Интерфейс репозитория
│   │   │   └── errors.go            # Доменные ошибки
│   │   └── analysis/                # Доменная модель анализа
│   │       ├── entity.go
│   │       └── repository.go
│   ├── usecase/                     # 💼 СЛОЙ USE CASES
│   │   ├── tender/                  # Use cases для тендеров
│   │   │   ├── interfaces.go        # Интерфейсы для dependencies
│   │   │   ├── create_tender.go     # Создание тендера
│   │   │   ├── get_tenders.go       # Получение списка тендеров
│   │   │   └── update_tender.go     # Обновление тендера
│   │   ├── discovery/               # Use cases для поиска
│   │   │   ├── interfaces.go
│   │   │   └── discover_tenders.go  # Поиск новых тендеров
│   │   └── analysis/                # Use cases для AI анализа
│   │       ├── interfaces.go
│   │       └── analyze_tender.go    # Анализ релевантности
│   ├── infrastructure/              # 🌐 ИНФРАСТРУКТУРНЫЙ СЛОЙ
│   │   ├── database/                # Работа с БД
│   │   │   ├── postgres.go          # Подключение к PostgreSQL
│   │   │   └── tender_repository.go # Реализация tender repository
│   │   ├── scraping/                # Web scraping
│   │   │   └── zakupki_scraper.go   # Scraper для zakupki.gov.ru
│   │   └── ai/                      # AI интеграция
│   │       └── llama_client.go      # Клиент для Llama через Ollama
│   └── interfaces/                  # 🔌 СЛОЙ ИНТЕРФЕЙСОВ
│       ├── http/                    # HTTP API
│       │   ├── server.go            # HTTP сервер
│       │   ├── handlers/            # HTTP handlers
│       │   │   ├── tender_handler.go
│       │   │   └── health_handler.go
│       │   ├── middleware/          # HTTP middleware
│       │   │   └── logging.go
│       │   └── routes.go            # Маршруты
│       └── cli/                     # CLI интерфейс
│           └── commands.go
├── 🧰 pkg/                          # Переиспользуемые утилиты
│   ├── logger/                      # Structured logging
│   │   └── logger.go
│   ├── validator/                   # Валидация данных
│   │   └── validator.go
│   └── container/                   # DI контейнер
│       └── container.go
├── 🧪 test/                         # Интеграционные тесты
│   └── integration/
│       └── api_test.go
├── 🐳 deployments/                  # Docker и развертывание
│   ├── Dockerfile
│   └── docker-compose.yml
└── 📝 docs/                         # Дополнительная документация
    ├── api.md                       # API документация
    └── deployment.md                # Инструкции по развертыванию
```

## 🛠️ Детальное руководство по реализации

### 1️⃣ Настройка окружения

#### Шаг 1: Инициализация Go модуля
```bash
cd TenderAutomationMVP
go mod init tender-automation-mvp
```

#### Шаг 2: Установка зависимостей
- **Web framework**: `github.com/gin-gonic/gin`
- **Database**: `github.com/lib/pq` (PostgreSQL)
- **Configuration**: `github.com/spf13/viper`
- **Logging**: `go.uber.org/zap`
- **Web scraping**: `github.com/gocolly/colly/v2`
- **AI client**: `github.com/ollama/ollama` (для Llama)
- **Validation**: `github.com/go-playground/validator/v10`
- **Testing**: `github.com/stretchr/testify`

### 2️⃣ Реализация Domain Layer

#### `internal/domain/tender/entity.go`
```go
// TODO: Реализовать доменную сущность Tender
// 
// Основные требования:
// 1. Содержит ТОЛЬКО бизнес-логику, без зависимостей от внешних систем
// 2. Методы для валидации данных
// 3. Бизнес-правила (например, IsActive(), CanBeAnalyzed())
// 4. Value Objects для типизированных полей (Status, Platform)
// 5. Доменные события (если нужны)
//
// Основные поля:
// - ID, ExternalID, Title, Description
// - Platform, URL, Status
// - StartPrice, Currency
// - PublishedAt, DeadlineAt
// - AIScore, AIRecommendation
// - CreatedAt, UpdatedAt
//
// Методы:
// - NewTender() - конструктор с валидацией
// - Validate() - валидация всех полей
// - IsActive() - проверка активности
// - CanBeAnalyzed() - можно ли анализировать
// - UpdateAIAnalysis() - обновление AI анализа
```

#### `internal/domain/tender/repository.go`
```go
// TODO: Определить интерфейс TenderRepository
//
// Основные требования:
// 1. Интерфейс определяется в domain слое
// 2. Реализация будет в infrastructure слое
// 3. Методы должны возвращать доменные сущности
// 4. Использовать context.Context для таймаутов
//
// Основные методы:
// - Create(ctx, tender) error
// - GetByID(ctx, id) (*Tender, error)
// - GetByExternalID(ctx, externalID) (*Tender, error)
// - List(ctx, filters) ([]*Tender, error)
// - Update(ctx, tender) error
// - Delete(ctx, id) error
// - GetPendingAnalysis(ctx) ([]*Tender, error)
```

### 3️⃣ Реализация Use Cases Layer

#### `internal/usecase/tender/interfaces.go`
```go
// TODO: Определить интерфейсы для dependencies use cases
//
// Основные интерфейсы:
// 1. TenderRepository (из domain слоя)
// 2. Logger - для логирования
// 3. Validator - для валидации входных данных
//
// Принципы:
// - Use cases НЕ знают о деталях реализации
// - Зависимости инжектируются через интерфейсы
// - Каждый use case решает одну конкретную задачу
```

#### `internal/usecase/tender/create_tender.go`
```go
// TODO: Реализовать use case для создания тендера
//
// Алгоритм:
// 1. Валидировать входные данные
// 2. Проверить, что тендер не существует (по ExternalID)
// 3. Создать доменную сущность Tender
// 4. Сохранить в репозитории
// 5. Логировать операцию
// 6. Вернуть созданный тендер
//
// Обработка ошибок:
// - Валидация данных
// - Дублирование тендеров
// - Ошибки БД
//
// Структура:
// type CreateTenderUseCase struct {
//     repo TenderRepository
//     logger Logger
//     validator Validator
// }
```

#### `internal/usecase/discovery/discover_tenders.go`
```go
// TODO: Реализовать use case для поиска новых тендеров
//
// Алгоритм:
// 1. Запустить web scraper
// 2. Получить список новых тендеров
// 3. Для каждого тендера:
//    - Проверить, существует ли уже в БД
//    - Если нет - создать новый через CreateTenderUseCase
//    - Логировать результат
// 4. Вернуть статистику (найдено, создано, ошибки)
//
// Dependencies:
// - Scraper interface
// - TenderRepository
// - CreateTenderUseCase
// - Logger
```

#### `internal/usecase/analysis/analyze_tender.go`
```go
// TODO: Реализовать use case для AI анализа тендера
//
// Алгоритм:
// 1. Получить тендер из репозитория
// 2. Проверить, можно ли его анализировать
// 3. Подготовить prompt для AI
// 4. Отправить запрос к AI модели
// 5. Парсить ответ AI
// 6. Обновить тендер результатами анализа
// 7. Сохранить в репозитории
//
// Dependencies:
// - AIClient interface
// - TenderRepository
// - Logger
//
// Промпт для AI:
// "Проанализируй тендер: {title}. 
//  Это медицинское оборудование? 
//  Оценка релевантности (0-1):
//  Рекомендация (participate/skip/analyze):"
```

### 4️⃣ Реализация Infrastructure Layer

#### `internal/infrastructure/database/postgres.go`
```go
// TODO: Реализовать подключение к PostgreSQL
//
// Основные задачи:
// 1. Создание connection pool
// 2. Конфигурация таймаутов и лимитов
// 3. Health check для БД
// 4. Graceful shutdown
// 5. Миграции (опционально)
//
// Настройки:
// - MaxOpenConns: 25
// - MaxIdleConns: 5
// - ConnMaxLifetime: 1h
// - ConnMaxIdleTime: 30m
//
// Функции:
// - NewPostgresDB(config) (*sql.DB, error)
// - HealthCheck(db) error
// - Close(db) error
```

#### `internal/infrastructure/database/tender_repository.go`
```go
// TODO: Реализовать TenderRepository для PostgreSQL
//
// Основные задачи:
// 1. Реализовать все методы интерфейса TenderRepository
// 2. SQL запросы с использованием подготовленных statements
// 3. Маппинг между БД моделями и доменными сущностями
// 4. Обработка ошибок БД
// 5. Оптимизация запросов (индексы, LIMIT, OFFSET)
//
// SQL таблица tenders:
// - id SERIAL PRIMARY KEY
// - external_id VARCHAR UNIQUE
// - title TEXT NOT NULL
// - platform VARCHAR NOT NULL
// - url TEXT
// - start_price DECIMAL
// - status VARCHAR
// - ai_score DECIMAL
// - ai_recommendation VARCHAR
// - published_at TIMESTAMP
// - deadline_at TIMESTAMP
// - created_at TIMESTAMP DEFAULT NOW()
// - updated_at TIMESTAMP DEFAULT NOW()
//
// Примеры запросов:
// - INSERT INTO tenders (...) VALUES (...) RETURNING *
// - SELECT * FROM tenders WHERE id = $1
// - UPDATE tenders SET ... WHERE id = $1
```

#### `internal/infrastructure/scraping/zakupki_scraper.go`
```go
// TODO: Реализовать scraper для zakupki.gov.ru
//
// Основные задачи:
// 1. Использовать colly для парсинга
// 2. Поиск по ключевым словам "медицинское оборудование"
// 3. Извлечение данных о тендерах
// 4. Rate limiting для избежания блокировок
// 5. User-Agent rotation
// 6. Обработка ошибок и retry логика
//
// Алгоритм:
// 1. Настроить colly collector
// 2. Установить поисковые параметры
// 3. Парсить страницы результатов
// 4. Извлекать данные каждого тендера:
//    - Номер, название, цена
//    - Ссылка, срок подачи заявок
//    - Заказчик, статус
// 5. Вернуть список TenderData структур
//
// Структура TenderData:
// - ExternalID, Title, Description
// - URL, Platform, StartPrice
// - PublishedAt, DeadlineAt
// - Customer, Status
//
// Rate limiting: 1 запрос в 2 секунды
// User-Agent: "Mozilla/5.0 (compatible; TenderBot/1.0)"
```

#### `internal/infrastructure/ai/llama_client.go`
```go
// TODO: Реализовать клиент для Llama через Ollama
//
// Основные задачи:
// 1. HTTP клиент для Ollama API
// 2. Подготовка промптов для анализа
// 3. Парсинг ответов AI
// 4. Обработка ошибок и таймаутов
// 5. Retry логика
//
// Алгоритм:
// 1. Подготовить промпт с данными тендера
// 2. Отправить POST запрос к Ollama API
// 3. Дождаться ответа (может быть долго)
// 4. Парсить JSON ответ
// 5. Извлечь релевантность и рекомендацию
//
// Промпт шаблон:
// "Проанализируй тендер на закупку:
//  Название: {title}
//  Описание: {description}
//  Заказчик: {customer}
//  
//  Вопросы:
//  1. Это медицинское оборудование? (да/нет)
//  2. Оценка релевантности от 0 до 1:
//  3. Рекомендация (participate/skip/analyze):
//  
//  Ответ в формате JSON:
//  {\"is_medical\": true, \"score\": 0.85, \"recommendation\": \"participate\"}"
//
// Ollama API endpoint: http://localhost:11434/api/generate
// Model: llama4-maviric или llama2
// Timeout: 60 секунд
```

### 5️⃣ Реализация Interfaces Layer

#### `internal/interfaces/http/server.go`
```go
// TODO: Реализовать HTTP сервер
//
// Основные задачи:
// 1. Настройка Gin router
// 2. Middleware (logging, recovery, CORS)
// 3. Регистрация маршрутов
// 4. Graceful shutdown
// 5. Health check endpoint
//
// Middleware:
// - Request/Response logging
// - Panic recovery
// - CORS для фронтенда
// - Request timeout
//
// Основные маршруты:
// GET /health - health check
// GET /api/v1/tenders - список тендеров
// GET /api/v1/tenders/:id - конкретный тендер
// POST /api/v1/tenders - создать тендер
// PUT /api/v1/tenders/:id - обновить тендер
// POST /api/v1/discover - запустить поиск тендеров
// POST /api/v1/analyze/:id - запустить анализ тендера
//
// Graceful shutdown:
// - Обработка SIGTERM/SIGINT
// - Завершение текущих запросов
// - Закрытие соединений
```

#### `internal/interfaces/http/handlers/tender_handler.go`
```go
// TODO: Реализовать HTTP handlers для тендеров
//
// Основные задачи:
// 1. Валидация входных данных
// 2. Вызов соответствующих use cases
// 3. Обработка ошибок
// 4. Форматирование ответов
// 5. HTTP статус коды
//
// Handlers:
// - GetTenders() - GET /api/v1/tenders
//   - Query параметры: page, limit, status, platform
//   - Response: {data: [...], total: 100, page: 1, limit: 20}
//
// - GetTenderByID() - GET /api/v1/tenders/:id
//   - Path параметр: id
//   - Response: {data: {...}}
//   - 404 если не найден
//
// - CreateTender() - POST /api/v1/tenders
//   - Request body: {external_id, title, platform, url, ...}
//   - Validation обязательных полей
//   - Response: 201 + created tender
//
// - UpdateTender() - PUT /api/v1/tenders/:id
//   - Path параметр: id
//   - Request body: partial tender data
//   - Response: updated tender
//
// - DiscoverTenders() - POST /api/v1/discover
//   - Query параметр: platform (опционально)
//   - Response: {found: 15, created: 12, errors: 0}
//
// - AnalyzeTender() - POST /api/v1/analyze/:id
//   - Path параметр: id
//   - Response: analysis result
//
// Обработка ошибок:
// - 400 для validation errors
// - 404 для not found
// - 409 для conflicts
// - 500 для internal errors
```

### 6️⃣ Конфигурация и DI

#### `configs/config.go`
```go
// TODO: Реализовать конфигурацию приложения
//
// Основные задачи:
// 1. Структуры для всех настроек
// 2. Загрузка из .env файлов
// 3. Переменные окружения
// 4. Валидация конфигурации
// 5. Дефолтные значения
//
// Разделы конфигурации:
// - Server (host, port, timeouts)
// - Database (postgres connection)
// - AI (ollama URL, model name)
// - Scraping (intervals, rate limits)
// - Logging (level, format)
//
// Источники конфигурации (приоритет):
// 1. Переменные окружения
// 2. .env файл
// 3. Дефолтные значения
//
// Использовать viper для загрузки
// Использовать validator для валидации
```

#### `pkg/container/container.go`
```go
// TODO: Реализовать DI контейнер
//
// Основные задачи:
// 1. Инициализация всех dependencies
// 2. Правильный порядок создания объектов
// 3. Передача dependencies через интерфейсы
// 4. Graceful shutdown всех ресурсов
//
// Порядок инициализации:
// 1. Config, Logger
// 2. Database connection
// 3. Repositories (зависят от БД)
// 4. Infrastructure services (AI, Scraper)
// 5. Use cases (зависят от repositories и services)
// 6. HTTP handlers (зависят от use cases)
// 7. HTTP server
//
// Методы:
// - NewContainer(config) (*Container, error)
// - Run(ctx) error - запуск всех сервисов
// - Shutdown() error - graceful shutdown
//
// Принципы:
// - Все зависимости инжектируются через интерфейсы
// - Никаких глобальных переменных
// - Легко тестировать с mock'ами
```

### 7️⃣ База данных и миграции

#### `migrations/001_create_tenders.up.sql`
```sql
-- TODO: Создать таблицу tenders
--
-- Основные требования:
-- 1. Эффективные индексы для поиска
-- 2. Правильные типы данных
-- 3. Ограничения целостности
-- 4. Значения по умолчанию
--
-- Поля таблицы:
-- - id: SERIAL PRIMARY KEY
-- - external_id: VARCHAR(100) UNIQUE NOT NULL (индекс)
-- - title: TEXT NOT NULL
-- - description: TEXT
-- - platform: VARCHAR(50) NOT NULL (индекс)
-- - url: TEXT
-- - customer: TEXT
-- - start_price: DECIMAL(15,2)
-- - currency: VARCHAR(3) DEFAULT 'RUB'
-- - status: VARCHAR(20) NOT NULL (индекс)
-- - ai_score: DECIMAL(3,2) CHECK (ai_score >= 0 AND ai_score <= 1)
-- - ai_recommendation: VARCHAR(20)
-- - ai_analyzed_at: TIMESTAMP
-- - published_at: TIMESTAMP
-- - deadline_at: TIMESTAMP (индекс)
-- - created_at: TIMESTAMP DEFAULT NOW()
-- - updated_at: TIMESTAMP DEFAULT NOW()
--
-- Индексы:
-- - unique_external_id на (external_id)
-- - idx_platform на (platform)
-- - idx_status на (status)
-- - idx_deadline на (deadline_at)
-- - idx_ai_score на (ai_score) WHERE ai_score IS NOT NULL
-- - idx_created_at на (created_at)
--
-- Constraints:
-- - CHECK status IN ('active', 'completed', 'cancelled', 'expired')
-- - CHECK ai_recommendation IN ('participate', 'skip', 'analyze')
-- - CHECK published_at <= deadline_at
```

## 📋 Этапы разработки

### Этап 1: Основа (3-5 дней) 🏗️
```bash
✅ Настройка проекта (go.mod, структура папок)
✅ Конфигурация и логирование
✅ Подключение к PostgreSQL
✅ Миграции БД
✅ Базовый HTTP сервер с health check
✅ DI контейнер
```

### Этап 2: Domain & Repository (2-3 дня) 🏛️
```bash
✅ Доменная модель Tender
✅ Интерфейс TenderRepository
✅ Реализация PostgreSQL repository
✅ Базовые unit тесты
```

### Этап 3: Use Cases (2-3 дня) 💼
```bash
✅ CreateTenderUseCase
✅ GetTendersUseCase
✅ UpdateTenderUseCase
✅ Тесты use cases с mock'ами
```

### Этап 4: Web Scraping (3-4 дня) 🕷️
```bash
✅ Zakupki scraper
✅ DiscoverTendersUseCase
✅ Интеграция с БД
✅ Rate limiting и error handling
```

### Этап 5: AI Integration (2-3 дня) 🤖
```bash
✅ Ollama клиент
✅ AnalyzeTenderUseCase
✅ Промпты для анализа
✅ Обработка AI ответов
```

### Этап 6: HTTP API (2-3 дня) 🌐
```bash
✅ HTTP handlers
✅ Middleware
✅ API documentation
✅ Интеграционные тесты
```

### Этап 7: Deployment (1-2 дня) 🚀
```bash
✅ Dockerfile
✅ docker-compose.yml
✅ Environment configuration
✅ Production готовность
```

**Общее время: 15-23 дня** (3-4 недели)

## 🔄 Расширение функциональности

После завершения MVP, вы легко можете добавить:

### 📄 Обработка документов
```bash
# Новые компоненты:
internal/domain/document/          # Доменная модель документа
internal/usecase/document/         # Use cases для обработки
internal/infrastructure/document/  # PDF/DOC парсеры
```

### 📧 Email автоматизация
```bash
# Новые компоненты:
internal/domain/supplier/          # Доменная модель поставщика
internal/domain/email_campaign/    # Доменная модель кампании
internal/usecase/supplier/         # Use cases для поставщиков
internal/infrastructure/email/     # SMTP/IMAP клиенты
```

### 🌐 Web UI
```bash
# Новые компоненты:
web/                               # Frontend приложение
internal/interfaces/web/           # Web handlers
templates/                         # HTML шаблоны
static/                            # CSS/JS файлы
```

### 📊 Продвинутая аналитика
```bash
# Новые компоненты:
internal/domain/analytics/         # Доменная модель аналитики
internal/usecase/analytics/        # Use cases для аналитики
internal/infrastructure/search/    # Elasticsearch/ClickHouse
```

## 🐳 Deployment и DevOps

### Локальная разработка
```bash
# Запуск БД
docker-compose up -d postgres

# Миграции
migrate -path migrations -database "postgres://..." up

# Запуск приложения
go run cmd/api/main.go
```

### Production deployment
```bash
# Сборка образа
docker build -t tender-automation-mvp .

# Запуск в production
docker-compose -f docker-compose.prod.yml up -d
```

### Environment variables
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=tender_automation

# AI
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama2

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🧪 Тестирование

### Unit тесты
```bash
# Тестирование доменных сущностей
go test ./internal/domain/...

# Тестирование use cases
go test ./internal/usecase/...
```

### Интеграционные тесты
```bash
# Тестирование API
go test ./test/integration/...

# Тестирование БД
go test ./internal/infrastructure/database/...
```

### Покрытие
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📚 Дополнительные ресурсы

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Gin Web Framework](https://gin-gonic.com/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Ollama Documentation](https://ollama.ai/docs)

---

## 🚀 Быстрый старт

1. **Склонируйте структуру проекта**
2. **Настройте окружение** (PostgreSQL, Ollama)
3. **Заполните конфигурацию** (.env файл)
4. **Реализуйте компоненты** поэтапно
5. **Тестируйте каждый этап**
6. **Разворачивайте в production**

**Удачи в разработке! 🎉**