# 🏗️ Архитектура системы автоматизации тендеров

## 📋 Общее описание проекта

Система автоматизации поиска и анализа медицинских тендеров с интеграцией AI для анализа технических заданий и автоматической отправки коммерческих предложений.

## 🎯 Основные функции системы

1. **Парсинг тендерных площадок** (szvo.gov35.ru, zakupki.gov.ru)
2. **AI-анализ тендеров** на соответствие критериям (медицинское оборудование)
3. **Извлечение и анализ технических заданий** из документов
4. **Поиск поставщиков** медицинского оборудования
5. **Автоматическая отправка запросов** на коммерческие предложения
6. **Управление процессом** через REST API и веб-интерфейс

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
│   │   │   └── analysis.go
│   │   ├── services/
│   │   │   ├── tender_service.go
│   │   │   ├── analysis_service.go
│   │   │   ├── document_service.go
│   │   │   └── email_service.go
│   │   └── ports/              # Интерфейсы (ports & adapters)
│   │       ├── repositories.go
│   │       ├── parsers.go
│   │       ├── analyzers.go
│   │       └── notifiers.go
│   ├── adapters/               # Внешние адаптеры
│   │   ├── repositories/
│   │   │   ├── postgres/
│   │   │   │   ├── tender_repo.go
│   │   │   │   ├── document_repo.go
│   │   │   │   └── supplier_repo.go
│   │   │   └── redis/
│   │   │       └── cache_repo.go
│   │   ├── parsers/
│   │   │   ├── zakupki_parser.go
│   │   │   ├── szvo_parser.go
│   │   │   └── base_parser.go
│   │   ├── analyzers/
│   │   │   ├── openai_analyzer.go
│   │   │   ├── claude_analyzer.go
│   │   │   └── document_processor.go
│   │   ├── notifiers/
│   │   │   ├── email_sender.go
│   │   │   └── telegram_bot.go
│   │   └── external/
│   │       ├── file_downloader.go
│   │       └── web_scraper.go
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
- **OpenAI API:** GPT-4 для анализа текстов
- **Document Processing:** UniDoc или Aspose
- **Text Processing:** GoNLP, regexp

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

### 1. **🕷️ Парсинг тендерных площадок**

#### Технические детали:
```go
type TenderParser interface {
    ParseTenders(filters ParseFilters) ([]*RawTender, error)
    ParseTenderDetails(tenderID string) (*TenderDetails, error)
    DownloadDocuments(tender *TenderDetails) ([]*Document, error)
}

type ParseFilters struct {
    Keywords    []string
    Categories  []string
    Region      string
    DateFrom    time.Time
    DateTo      time.Time
    MaxPrice    decimal.Decimal
}
```

#### Стратегия парсинга:
1. **User-Agent ротация** для избежания блокировок
2. **Задержки между запросами** (1-3 сек)
3. **Proxy ротация** при необходимости
4. **Распределенная нагрузка** через очереди
5. **Инкрементальный парсинг** (только новые тендеры)

#### Обработка ошибок:
- Exponential backoff при ошибках
- Логирование всех неудачных попыток
- Уведомления при критических ошибках
- Fallback на ручной режим

### 2. **🤖 AI Анализ тендеров**

#### Архитектура AI модуля:
```go
type AIAnalyzer interface {
    AnalyzeTender(tender *entities.Tender) (*AnalysisResult, error)
    ClassifyCategory(description string) (Category, float64, error)
    ExtractRequirements(techSpec string) (*Requirements, error)
    GenerateSearchTerms(requirements *Requirements) ([]string, error)
}

type AnalysisResult struct {
    IsRelevant      bool              `json:"is_relevant"`
    Confidence      float64           `json:"confidence"`
    Category        string            `json:"category"`
    Requirements    []*Requirement    `json:"requirements"`
    KeyTerms        []string          `json:"key_terms"`
    EstimatedValue  decimal.Decimal   `json:"estimated_value"`
    Reasoning       string            `json:"reasoning"`
}
```

#### Промпты для OpenAI:
```json
{
  "system_prompt": "Ты эксперт по медицинским тендерам. Анализируй тендеры и определяй, относятся ли они к поставке медицинского оборудования (не услуги!).",
  "user_prompt": "Проанализируй этот тендер: {tender_text}. Ответь в JSON формате с полями: is_relevant, confidence, category, requirements, reasoning."
}
```

#### Обработка различных типов документов:
- **DOC/DOCX:** UniDoc для извлечения текста и таблиц
- **PDF:** pdfcpu или Apache Tika
- **XLS/XLSX:** excelize для табличных данных
- **RTF:** специализированные парсеры

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

### 4. **🔍 Поиск поставщиков**

#### Стратегия поиска:
```go
type SupplierFinder interface {
    FindSuppliers(equipment *Equipment) ([]*Supplier, error)
    VerifySupplier(supplier *Supplier) (*SupplierVerification, error)
    GetContactInfo(supplier *Supplier) (*ContactInfo, error)
}
```

#### Источники поставщиков:
1. **Поиск в Google/Yandex** по ключевым словам
2. **Отраслевые каталоги** (medprom.ru, medrussia.ru)
3. **Социальные сети** и профессиональные платформы
4. **Собственная база данных** проверенных поставщиков

### 5. **📊 Система аналитики**

#### Метрики для отслеживания:
- Количество найденных тендеров
- Процент релевантных тендеров
- Успешность отправки запросов
- Количество полученных КП
- ROI по каждому тендеру

#### Дашборд компоненты:
- Графики динамики тендеров
- Тепловая карта регионов
- Топ категорий оборудования
- Эффективность поставщиков

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

-- Отправленные запросы
CREATE TABLE commercial_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tender_id UUID REFERENCES tenders(id),
    supplier_id UUID REFERENCES suppliers(id),
    email_sent_at TIMESTAMP,
    reply_received_at TIMESTAMP,
    status VARCHAR(50) DEFAULT 'sent',
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

## 🔍 Рекомендации по собеседованиям

Этот проект покрывает все ключевые темы для Middle Go Developer:

1. **Архитектура:** Clean Architecture, микросервисы
2. **Базы данных:** PostgreSQL, миграции, оптимизация
3. **Concurrency:** Goroutines, channels, worker pools
4. **HTTP:** REST API, middleware, аутентификация  
5. **Testing:** Unit, integration, mocking
6. **DevOps:** Docker, CI/CD, мониторинг
7. **Внешние API:** HTTP клиенты, rate limiting
8. **Обработка данных:** Парсинг, валидация, трансформация

**Совет:** Подготовь примеры кода и архитектурных решений из этого проекта для демонстрации на собеседованиях.

---

*Удачи в разработке! Этот проект станет отличным портфолио для твоей карьеры Middle Go Developer.*