# 🏗️ Архитектура системы автоматизации тендеров

## 📋 Общее описание проекта

Интеллектуальная система автоматизации полного цикла работы с тендерами: от поиска и анализа до ценообразования и участия в торгах. Система использует AI для анализа документов, накопления знаний о товарах и оптимизации ценовых стратегий.

## 🎯 Основные функции системы

1. **Интеллектуальный парсинг тендерных площадок** (szvo.gov35.ru, zakupki.gov.ru и др.)
2. **AI-фильтрация тендеров** по медицинскому оборудованию (исключая услуги)
3. **Автоматическое извлечение всех документов** тендера (техзадания + дополнительные файлы)
4. **Глубокий анализ технических заданий** с извлечением товаров из таблиц
5. **Накопление и каталогизация товаров** с сайтов поставщиков в БД
6. **AI-коммуникация с менеджерами** для получения коммерческих предложений
7. **Анализ истории тендеров** и AI-рекомендации по ценообразованию
8. **Полный цикл участия в тендере** с прогнозированием результатов

---

## 📁 Структура проекта

```
tender-automation-system/
├── cmd/
│   ├── api/                    # REST API сервер
│   │   └── main.go
│   ├── parser/                 # Сервис парсинга
│   │   └── main.go
│   ├── analyzer/               # AI анализатор
│   │   └── main.go
│   └── emailer/                # Email сервис
│       └── main.go
├── internal/
│   ├── api/                    # REST API слой
│   │   ├── handlers/
│   │   │   ├── tender.go
│   │   │   ├── analysis.go
│   │   │   ├── supplier.go
│   │   │   ├── product.go
│   │   │   ├── pricing.go
│   │   │   └── auth.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logging.go
│   │   └── routes/
│   │       └── routes.go
│   ├── core/                   # Бизнес логика (доменный слой)
│   │   ├── entities/
│   │   │   ├── tender.go
│   │   │   ├── document.go
│   │   │   ├── supplier.go
│   │   │   ├── product.go       # Каталог товаров
│   │   │   ├── pricing.go       # Ценовая аналитика
│   │   │   ├── tender_result.go # Результаты тендеров
│   │   │   └── analysis.go
│   │   ├── services/
│   │   │   ├── tender_service.go
│   │   │   ├── analysis_service.go
│   │   │   ├── document_service.go
│   │   │   ├── product_catalog_service.go
│   │   │   ├── pricing_ai_service.go
│   │   │   ├── email_ai_service.go
│   │   │   └── tender_history_service.go
│   │   └── ports/              # Интерфейсы (ports & adapters)
│   │       ├── repositories.go
│   │       ├── parsers.go
│   │       ├── analyzers.go
│   │       ├── ai_models.go     # AI интерфейсы
│   │       └── notifiers.go
│   ├── adapters/               # Внешние адаптеры
│   │   ├── repositories/
│   │   │   ├── postgres/
│   │   │   │   ├── tender_repo.go
│   │   │   │   ├── document_repo.go
│   │   │   │   ├── supplier_repo.go
│   │   │   │   ├── product_repo.go
│   │   │   │   ├── pricing_repo.go
│   │   │   │   └── tender_result_repo.go
│   │   │   └── redis/
│   │   │       └── cache_repo.go
│   │   ├── parsers/
│   │   │   ├── zakupki_parser.go
│   │   │   ├── szvo_parser.go
│   │   │   ├── supplier_site_parser.go
│   │   │   └── base_parser.go
│   │   ├── ai/                 # AI модели и адаптеры
│   │   │   ├── llama_adapter.go
│   │   │   ├── openai_adapter.go
│   │   │   ├── document_processor.go
│   │   │   ├── pricing_analyzer.go
│   │   │   └── email_bot.go
│   │   ├── notifiers/
│   │   │   ├── email_sender.go
│   │   │   └── telegram_bot.go
│   │   └── external/
│   │       ├── file_downloader.go
│   │       ├── web_scraper.go
│   │       └── product_finder.go
│   ├── config/
│   │   ├── config.go
│   │   └── database.go
│   └── shared/                 # Общие утилиты
│       ├── errors/
│       │   └── errors.go
│       ├── logger/
│       │   └── logger.go
│       ├── utils/
│       │   ├── text_processor.go
│       │   ├── file_utils.go
│       │   └── validation.go
│       └── constants/
│           └── constants.go
├── pkg/                        # Публичные библиотеки
│   ├── database/
│   │   ├── postgres.go
│   │   └── migrations.go
│   ├── queue/
│   │   └── redis_queue.go
│   ├── cache/
│   │   └── redis_cache.go
│   └── http/
│       └── client.go
├── migrations/                 # SQL миграции
│   ├── 001_create_tenders.sql
│   ├── 002_create_documents.sql
│   ├── 003_create_suppliers.sql
│   └── 004_create_analysis.sql
├── configs/                    # Конфигурационные файлы
│   ├── config.yaml
│   ├── docker-compose.yml
│   └── .env.example
├── scripts/                    # Скрипты для развертывания
│   ├── setup.sh
│   └── deploy.sh
├── docs/                       # Документация
│   ├── api.md
│   ├── deployment.md
│   └── examples.md
├── tests/                      # Тесты
│   ├── integration/
│   ├── unit/
│   └── fixtures/
├── web/                        # Frontend (опционально)
│   ├── static/
│   └── templates/
├── .gitignore
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile
└── README.md
```

---

## 🏛️ Архитектурные паттерны

### 1. **Clean Architecture (Hexagonal Architecture)**

```
External → Adapters → Core (Business Logic) → Entities
```

**Преимущества:**
- Независимость от внешних сервисов
- Легкость тестирования
- Возможность замены компонентов

### 2. **Repository Pattern**
```go
type TenderRepository interface {
    Create(tender *entities.Tender) error
    FindByID(id string) (*entities.Tender, error)
    FindByFilters(filters TenderFilters) ([]*entities.Tender, error)
    Update(tender *entities.Tender) error
    Delete(id string) error
}
```

### 3. **Service Layer Pattern**
```go
type TenderService struct {
    tenderRepo     ports.TenderRepository
    documentRepo   ports.DocumentRepository
    analyzer       ports.AIAnalyzer
    logger         logger.Logger
}
```

### 4. **Factory Pattern** для парсеров
```go
type ParserFactory interface {
    CreateParser(site string) (Parser, error)
}
```

### 5. **Strategy Pattern** для AI анализаторов
```go
type AnalysisStrategy interface {
    Analyze(content string) (*AnalysisResult, error)
}
```

### 6. **Observer Pattern** для уведомлений
```go
type EventBus interface {
    Subscribe(event string, handler EventHandler)
    Publish(event string, data interface{})
}
```

---

## 🔄 Стратегия разработки

### **Фаза 1: Основа (2 недели)**

#### Неделя 1: Инфраструктура
- [x] Настройка проекта и структуры
- [x] База данных (PostgreSQL) + миграции
- [x] Базовые модели (Tender, Document, Supplier)
- [x] REST API основа с аутентификацией
- [x] Логирование и конфигурация

#### Неделя 2: Базовый CRUD
- [x] CRUD операции для тендеров
- [x] Система пользователей и ролей
- [x] Базовые эндпоинты API
- [x] Простые тесты

### **Фаза 2: Парсинг и автоматизация (2 недели)**

#### Неделя 3: Web Scraping
- [x] Парсер для zakupki.gov.ru
- [x] Система очередей (Redis)
- [x] Фоновые задачи для парсинга
- [x] Обработка ошибок и ретраи

#### Неделя 4: Обработка документов
- [x] Загрузка файлов (DOC/DOCX)
- [x] Извлечение текста из документов
- [x] Парсинг таблиц технических заданий
- [x] Сохранение структурированных данных

### **Фаза 3: AI и умная обработка (1.5 недели)**

#### Неделя 5-6: AI Интеграция
- [x] OpenAI API интеграция
- [x] Промпты для анализа тендеров
- [x] Классификация медицинского оборудования
- [x] Извлечение ключевых характеристик

### **Фаза 4: Email автоматизация (0.5 недели)**

#### Неделя 6-7: Коммуникации
- [x] Email отправка
- [x] Шаблоны писем
- [x] Интеграция с поставщиками
- [x] Отслеживание ответов

---

## 🔧 Технический стек

### **Backend**
- **Язык:** Go 1.21+
- **Фреймворк:** Gin или Chi (легковесные)
- **БД:** PostgreSQL 15+
- **Кэш:** Redis 7+
- **Очереди:** Redis + asynq или RabbitMQ

### **AI & ML**
- **Primary AI:** Llama 4 (через Ollama или API) для основного анализа
- **Backup AI:** OpenAI GPT-4 как резервный вариант
- **Document Processing:** UniDoc для DOC/DOCX, pdfcpu для PDF
- **Text Processing:** GoNLP, regexp для предобработки
- **Price Analytics:** Собственные алгоритмы + AI для ценового анализа

### **Инфраструктура**
- **Контейнеризация:** Docker + Docker Compose
- **Логирование:** Zap или Logrus
- **Мониторинг:** Prometheus + Grafana
- **Деплой:** GitHub Actions

### **Внешние сервисы**
- **Email:** SendGrid или SMTP
- **Уведомления:** Telegram Bot API
- **Файловое хранилище:** MinIO или AWS S3

---

## 🧠 Детальный разбор компонентов

### 1. **🕷️ Интеллектуальный парсинг тендерных площадок**

#### Технические детали:
```go
type TenderParser interface {
    ParseTenders(filters ParseFilters) ([]*RawTender, error)
    ParseTenderDetails(tenderID string) (*TenderDetails, error)
    DownloadAllDocuments(tender *TenderDetails) ([]*Document, error)
    NavigateToDocumentsSection(tenderURL string) error
}

type ParseFilters struct {
    Keywords       []string
    Categories     []string
    Region         string
    DateFrom       time.Time
    DateTo         time.Time
    MaxPrice       decimal.Decimal
    ExcludeServices bool `default:"true"` // Исключать услуги
    MedicalOnly     bool `default:"true"` // Только медицинское оборудование
}

type DocumentType struct {
    Name           string
    Type           string // "technical_spec", "additional", "contract"
    IsTechnicalSpec bool
    Priority       int
}
```

#### Расширенная стратегия парсинга:
1. **Автоматическая навигация** к разделу "Документы" в каждом тендере
2. **Распознавание типов документов** (техзадание vs дополнительные)
3. **Скачивание ВСЕХ документов** для полного анализа
4. **User-Agent ротация** для избежания блокировок
5. **Задержки между запросами** (1-3 сек)
6. **Proxy ротация** при необходимости
7. **Распределенная нагрузка** через очереди
8. **Инкрементальный парсинг** (только новые тендеры)

#### Обработка ошибок:
- Exponential backoff при ошибках
- Логирование всех неудачных попыток
- Уведомления при критических ошибках
- Fallback на ручной режим
- Автоматический retry при сбоях скачивания

### 2. **🤖 Мультимодальный AI Анализ**

#### Архитектура AI модуля:
```go
type AIAnalyzer interface {
    AnalyzeTender(tender *entities.Tender) (*AnalysisResult, error)
    ClassifyCategory(description string) (Category, float64, error)
    ExtractRequirements(techSpec string) (*Requirements, error)
    ExtractProductsFromTable(documentText string) ([]*Product, error)
    AnalyzeAllDocuments(documents []*Document) (*CompleteAnalysis, error)
    GenerateSearchTerms(requirements *Requirements) ([]string, error)
}

type AnalysisResult struct {
    IsRelevant         bool               `json:"is_relevant"`
    Confidence         float64            `json:"confidence"`
    Category           string             `json:"category"`
    Requirements       []*Requirement     `json:"requirements"`
    ExtractedProducts  []*Product         `json:"extracted_products"`
    DeliveryTimeframe  string             `json:"delivery_timeframe"`
    SpecialConditions  []string           `json:"special_conditions"`
    KeyTerms           []string           `json:"key_terms"`
    EstimatedValue     decimal.Decimal    `json:"estimated_value"`
    Reasoning          string             `json:"reasoning"`
}

type Product struct {
    Name           string          `json:"name"`
    Description    string          `json:"description"`
    Quantity       int             `json:"quantity"`
    Unit           string          `json:"unit"`
    Specifications map[string]string `json:"specifications"`
    Category       string          `json:"category"`
    SearchTerms    []string        `json:"search_terms"`
}
```

#### Промпты для Llama 4:
```json
{
  "system_prompt": "Ты эксперт-аналитик медицинских тендеров с глубокими знаниями медицинского оборудования. Твоя задача: анализировать документы тендеров и извлекать структурированную информацию о товарах. ВАЖНО: работаем ТОЛЬКО с поставкой товаров, НЕ с услугами!",
  
  "document_analysis_prompt": "Проанализируй все документы тендера. Извлеки:\n1. Список всех товаров из таблиц технического задания\n2. Сроки поставки\n3. Особые условия\n4. Оцени релевантность (медицинское оборудование)\nДокументы: {documents_text}",
  
  "product_extraction_prompt": "Из этой таблицы извлеки все товары с характеристиками:\n{table_text}\nВерни в JSON формате массив товаров с полями: name, description, quantity, unit, specifications, category."
}
```

#### Обработка различных типов документов:
- **DOC/DOCX:** UniDoc для извлечения текста и таблиц
- **PDF:** pdfcpu для текста + pdf2go для сложных макетов
- **XLS/XLSX:** excelize для табличных данных
- **RTF:** специализированные парсеры Go
- **Изображения таблиц:** OCR через Tesseract + AI для анализа

### 3. **📧 Email автоматизация**

#### Архитектура email системы:
```go
type EmailService interface {
    SendCommercialRequest(supplier *Supplier, requirements *Requirements) error
    SendNotification(user *User, notification *Notification) error
    ProcessReplies() ([]*EmailReply, error)
}

type EmailTemplate struct {
    Subject     string
    BodyHTML    string
    BodyText    string
    Attachments []Attachment
}
```

#### Шаблоны писем:
1. **Запрос коммерческого предложения**
2. **Уведомление о новых тендерах**
3. **Статистические отчеты**
4. **Напоминания о дедлайнах**

#### Интеграция с CRM:
- Отслеживание открытий писем
- Статистика ответов
- Автоматическая категоризация ответов

### 4. **🔍 Интеллектуальный поиск и каталогизация товаров**

#### Стратегия поиска и накопления:
```go
type ProductCatalogService interface {
    SearchProductOnline(product *Product) ([]*SupplierProduct, error)
    CatalogizeSupplierSite(siteURL string) error
    FindProductInCatalog(product *Product) ([]*CatalogedProduct, error)
    UpdateProductCatalog() error
}

type SupplierFinder interface {
    FindSuppliers(equipment *Equipment) ([]*Supplier, error)
    VerifySupplier(supplier *Supplier) (*SupplierVerification, error)
    GetContactInfo(supplier *Supplier) (*ContactInfo, error)
    ParseSupplierProducts(supplierURL string) ([]*Product, error)
}

type CatalogedProduct struct {
    ProductID       string          `json:"product_id"`
    Name           string          `json:"name"`
    Description    string          `json:"description"`
    SupplierID     string          `json:"supplier_id"`
    ProductURL     string          `json:"product_url"`
    Specifications map[string]string `json:"specifications"`
    PriceRange     string          `json:"price_range"`
    Availability   string          `json:"availability"`
    LastUpdated    time.Time       `json:"last_updated"`
}
```

#### Источники поставщиков и стратегия накопления:
1. **Поиск в Google/Yandex** по ключевым словам товаров
2. **Отраслевые каталоги** (medprom.ru, medrussia.ru, zdrav.ru)
3. **Социальные сети** и профессиональные платформы
4. **Собственная база данных** проверенных поставщиков
5. **Автоматическая каталогизация** - парсинг сайтов поставщиков для накопления товаров
6. **Сохранение URL товаров** для быстрого доступа в будущем
7. **Периодическое обновление** каталога товаров
8. **Поиск по накопленному каталогу** перед поиском в интернете

### 5. **🧠 AI-система ценообразования и прогнозирования**

#### Архитектура модуля ценообразования:
```go
type PricingAI interface {
    AnalyzeTenderHistory(tenderID string) (*PricingAnalysis, error)
    RecommendOptimalPrice(tender *Tender, products []*Product) (*PriceRecommendation, error)
    PredictWinningChance(price decimal.Decimal, tender *Tender) (float64, error)
    LearnFromTenderResult(result *TenderResult) error
}

type PricingAnalysis struct {
    HistoricalPrices    []decimal.Decimal  `json:"historical_prices"`
    AverageWinningPrice decimal.Decimal    `json:"average_winning_price"`
    MarketTrends        []*MarketTrend     `json:"market_trends"`
    CompetitorAnalysis  []*Competitor      `json:"competitor_analysis"`
    Recommendations    *PriceRecommendation `json:"recommendations"`
}

type TenderResult struct {
    TenderID      string          `json:"tender_id"`
    WinnerCompany string          `json:"winner_company"`
    WinningPrice  decimal.Decimal `json:"winning_price"`
    OurBid        decimal.Decimal `json:"our_bid"`
    OurPosition   int             `json:"our_position"`
    ResultDate    time.Time       `json:"result_date"`
    Products      []*Product      `json:"products"`
}
```

#### AI Email Bot для коммуникации с менеджерами:
```go
type EmailAI interface {
    GenerateInitialRequest(product *Product, supplier *Supplier) (*EmailTemplate, error)
    ParseManagerReply(emailContent string) (*ReplyAnalysis, error)
    GenerateFollowUp(previousConversation []*Email) (*EmailTemplate, error)
    NegotiatePrice(targetPrice decimal.Decimal, currentOffer decimal.Decimal) (*EmailTemplate, error)
}
```

### 6. **📊 Расширенная система аналитики**

#### Метрики для отслеживания:
- Количество найденных тендеров
- Процент релевантных тендеров (с AI фильтрацией)
- Успешность отправки запросов
- Количество полученных КП
- ROI по каждому тендеру
- **Точность ценовых прогнозов AI**
- **Процент выигранных тендеров**
- **Эффективность накопленного каталога товаров**

#### Дашборд компоненты:
- Графики динамики тендеров
- Тепловая карта регионов
- Топ категорий оборудования
- Эффективность поставщиков
- **AI Pricing Dashboard** - анализ точности прогнозов
- **Каталог товаров** - статистика накопления и использования
- **Email Bot Performance** - эффективность AI-коммуникации

---

## 🗄️ Схема базы данных

### Основные таблицы:

```sql
-- Тендеры
CREATE TABLE tenders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) UNIQUE NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    customer_name VARCHAR(500),
    start_price DECIMAL(15,2),
    site_source VARCHAR(100) NOT NULL,
    publication_date TIMESTAMP,
    deadline TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Документы тендеров
CREATE TABLE tender_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tender_id UUID REFERENCES tenders(id),
    filename VARCHAR(500),
    file_type VARCHAR(50),
    file_size BIGINT,
    download_url TEXT,
    local_path TEXT,
    content_text TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AI анализ
CREATE TABLE tender_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tender_id UUID REFERENCES tenders(id),
    is_relevant BOOLEAN,
    confidence_score DECIMAL(3,2),
    category VARCHAR(200),
    extracted_requirements JSONB,
    reasoning TEXT,
    analyzed_at TIMESTAMP DEFAULT NOW()
);

-- Поставщики
CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(500) NOT NULL,
    website VARCHAR(500),
    email VARCHAR(255),
    phone VARCHAR(100),
    specializations TEXT[],
    reliability_score DECIMAL(3,2),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Каталог товаров от поставщиков
CREATE TABLE cataloged_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_name VARCHAR(500) NOT NULL,
    description TEXT,
    supplier_id UUID REFERENCES suppliers(id),
    product_url TEXT,
    specifications JSONB,
    price_range VARCHAR(100),
    availability VARCHAR(100),
    category VARCHAR(200),
    search_terms TEXT[],
    last_updated TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Результаты тендеров для обучения AI
CREATE TABLE tender_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tender_id UUID REFERENCES tenders(id),
    winner_company VARCHAR(500),
    winning_price DECIMAL(15,2),
    our_bid DECIMAL(15,2),
    our_position INTEGER,
    total_participants INTEGER,
    result_date TIMESTAMP,
    products_data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Ценовая аналитика AI
CREATE TABLE pricing_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tender_id UUID REFERENCES tenders(id),
    recommended_price DECIMAL(15,2),
    confidence_score DECIMAL(3,2),
    winning_probability DECIMAL(3,2),
    market_analysis JSONB,
    competitor_analysis JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Отправленные запросы и email переписка
CREATE TABLE commercial_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tender_id UUID REFERENCES tenders(id),
    supplier_id UUID REFERENCES suppliers(id),
    product_id UUID REFERENCES cataloged_products(id),
    email_thread JSONB, -- История переписки
    initial_request_sent_at TIMESTAMP,
    last_reply_at TIMESTAMP,
    commercial_offer_received BOOLEAN DEFAULT FALSE,
    final_price DECIMAL(15,2),
    status VARCHAR(50) DEFAULT 'sent',
    ai_conversation_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 🚀 План запуска и развертывания

### **Локальная разработка:**
```bash
# 1. Клонирование и настройка
git clone <repo>
cd tender-automation-system
cp configs/.env.example .env

# 2. Запуск инфраструктуры
docker-compose up -d postgres redis

# 3. Миграции
make migrate

# 4. Запуск сервисов
make run-api
make run-parser
make run-analyzer
```

### **Production деплой:**
- **Kubernetes** кластер
- **Helm charts** для управления
- **CI/CD** через GitHub Actions
- **Мониторинг** с Prometheus/Grafana

---

## 🎯 Ключевые особенности реализации

### **1. Масштабируемость:**
- Микросервисная архитектура
- Горизонтальное масштабирование парсеров
- Очереди для асинхронной обработки

### **2. Надежность:**
- Graceful shutdown
- Health checks
- Circuit breaker pattern
- Distributed tracing

### **3. Безопасность:**
- JWT аутентификация
- Rate limiting
- Input validation
- SQL injection protection

### **4. Производительность:**
- Redis кэширование
- Database connection pooling
- Оптимизированные SQL запросы
- Batch processing

---

## 📈 Этапы внедрения

### **MVP (4 недели):**
- Парсинг одной площадки
- Базовый AI анализ
- Простая отправка email
- Web интерфейс для управления

### **v1.0 (6-8 недель):**
- Несколько тендерных площадок
- Продвинутый AI анализ
- Система уведомлений
- Аналитика и отчеты

### **v2.0 (3-4 месяца):**
- Машинное обучение для улучшения точности
- Интеграция с CRM системами
- Mobile приложение
- API для интеграций

---

## 💡 О модели Llama 4

**Важное замечание:** На момент написания (январь 2025) модели "Llama 4 Maverick" не существует. Возможные варианты:

1. **Llama 3.1 или 3.2** - актуальные версии от Meta
2. **Llama 3 70B/405B** - для сложных задач анализа
3. **Запуск через Ollama** локально для экономии
4. **API провайдеры** - Together AI, Replicate для облачного доступа

**Рекомендация:** Начни с **Llama 3.1 8B** через Ollama для разработки, затем переходи на **70B** для продакшена.

## 🔍 Рекомендации по собеседованиям

Этот проект покрывает все ключевые темы для Middle Go Developer:

1. **Архитектура:** Clean Architecture, микросервисы
2. **Базы данных:** PostgreSQL, миграции, оптимизация, JSONB
3. **Concurrency:** Goroutines, channels, worker pools
4. **HTTP:** REST API, middleware, аутентификация  
5. **AI Integration:** Работа с LLM API, prompt engineering
6. **Web Scraping:** Сложный парсинг, обход защиты
7. **Document Processing:** Работа с различными форматами файлов
8. **Testing:** Unit, integration, mocking
9. **DevOps:** Docker, CI/CD, мониторинг
10. **Внешние API:** HTTP клиенты, rate limiting
11. **Обработка данных:** Парсинг, валидация, трансформация
12. **Машинное обучение:** Интеграция AI в бизнес-логику

**Совет:** Подготовь примеры кода и архитектурных решений из этого проекта для демонстрации на собеседованиях. Особенно впечатлят AI-компоненты и сложная архитектура.

---

*Удачи в разработке! Этот проект станет отличным портфолио для твоей карьеры Middle Go Developer.*