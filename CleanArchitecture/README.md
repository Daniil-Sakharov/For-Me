# 🏗️ Clean Architecture URL Shortener

Это **полноценная реализация Clean Architecture** на Go для сервиса сокращения URL. Проект демонстрирует все принципы чистой архитектуры с детальными комментариями в коде.

## 📋 Содержание

- [🎯 Что такое Clean Architecture](#-что-такое-clean-architecture)
- [🏛️ Архитектура проекта](#️-архитектура-проекта)
- [📁 Структура папок](#-структура-папок)
- [🔄 Поток данных](#-поток-данных)
- [🚀 Запуск проекта](#-запуск-проекта)
- [📡 API Endpoints](#-api-endpoints)
- [🧪 Тестирование](#-тестирование)
- [🔧 Конфигурация](#-конфигурация)
- [📈 Преимущества архитектуры](#-преимущества-архитектуры)

## 🎯 Что такое Clean Architecture

**Clean Architecture** - это архитектурный подход, предложенный Робертом Мартином (Uncle Bob), который обеспечивает:

### ✅ Основные принципы:

1. **Независимость от фреймворков** - бизнес-логика не зависит от внешних библиотек
2. **Тестируемость** - можно тестировать без UI, БД и веб-сервера
3. **Независимость от UI** - можно легко заменить веб-интерфейс на консольный
4. **Независимость от БД** - можно заменить PostgreSQL на MongoDB без изменения бизнес-логики
5. **Независимость от внешних сервисов** - бизнес-правила не знают о внешнем мире

### 🏗️ Слои архитектуры:

```
┌─────────────────────────────────────────────────────────────┐
│                    🌐 FRAMEWORKS & DRIVERS                  │
│              (Web, DB, External APIs, Devices)             │
├─────────────────────────────────────────────────────────────┤
│                   🔌 INTERFACE ADAPTERS                     │
│            (Controllers, Gateways, Presenters)             │
├─────────────────────────────────────────────────────────────┤
│                     💼 USE CASES                            │
│              (Application Business Rules)                  │
├─────────────────────────────────────────────────────────────┤
│                     🏛️ ENTITIES                             │
│               (Enterprise Business Rules)                  │
└─────────────────────────────────────────────────────────────┘
```

### 📏 Правило зависимостей:

Зависимости **ВСЕГДА** направлены внутрь. Внутренние слои НЕ ЗНАЮТ о внешних слоях.

## 🏛️ Архитектура проекта

### 🔵 Domain Layer (Entities)
**Расположение:** `internal/domain/`

**Что содержит:**
- 🏛️ **Entities** - доменные модели с бизнес-логикой
- 🔗 **Repository Interfaces** - интерфейсы для работы с данными
- ❌ **Domain Errors** - доменные ошибки

**Принципы:**
- ✅ НЕТ зависимостей от внешних библиотек (кроме стандартной библиотеки Go)
- ✅ Содержит чистую бизнес-логику
- ✅ Определяет интерфейсы для внешних сервисов

**Пример:**
```go
// internal/domain/user/entity.go
type User struct {
    ID             uint
    Email          string
    HashedPassword string
    Name           string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func NewUser(email, name string) (*User, error) {
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    // Доменные правила валидации
}
```

### 🔵 Use Cases Layer (Application Business Rules)
**Расположение:** `internal/usecase/`

**Что содержит:**
- 🎯 **Use Cases** - сценарии использования приложения
- 🔌 **Interfaces** - интерфейсы для внешних зависимостей
- 🔄 **Request/Response** - структуры для входных и выходных данных

**Принципы:**
- ✅ Оркестрирует взаимодействие между Entities
- ✅ Содержит application-специфичную бизнес-логику
- ✅ Зависит только от Domain слоя

**Пример:**
```go
// internal/usecase/auth/register.go
func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
    // 1. Проверяем существование пользователя
    exists, err := uc.userRepo.Exists(ctx, req.Email)
    
    // 2. Создаем доменную модель
    newUser, err := user.NewUser(req.Email, req.Name)
    
    // 3. Хешируем пароль через внешний сервис
    hashedPassword, err := uc.passwordHasher.Hash(req.Password)
    
    // 4. Сохраняем через репозиторий
    err = uc.userRepo.Save(ctx, newUser)
    
    return &RegisterResponse{Token: token}, nil
}
```

### 🔵 Interface Adapters Layer
**Расположение:** `internal/interfaces/`

**Что содержит:**
- 🎮 **Controllers** - HTTP handlers
- 📊 **Presenters** - форматирование ответов
- 🔗 **Gateways** - адаптеры для внешних сервисов

**Принципы:**
- ✅ Преобразует данные между внешним миром и Use Cases
- ✅ Обрабатывает HTTP протокол
- ✅ НЕ содержит бизнес-логики

**Пример:**
```go
// internal/interfaces/controllers/auth_controller.go
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Извлекаем данные из HTTP запроса
    var req auth.RegisterRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 2. Вызываем Use Case
    response, err := c.registerUC.Execute(r.Context(), req)
    
    // 3. Возвращаем HTTP ответ
    c.writeJSONResponse(w, response, http.StatusCreated)
}
```

### 🔵 Infrastructure Layer (Frameworks & Drivers)
**Расположение:** `internal/infrastructure/`

**Что содержит:**
- 🗄️ **Database** - работа с PostgreSQL
- 🌐 **Web** - HTTP сервер и middleware
- 🔐 **External** - внешние сервисы (JWT, bcrypt, etc.)

**Принципы:**
- ✅ Реализует интерфейсы из внутренних слоев
- ✅ Содержит технические детали
- ✅ Может быть заменен без изменения бизнес-логики

**Пример:**
```go
// internal/infrastructure/database/user_repository.go
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
    query := `INSERT INTO users (email, hashed_password, name) VALUES ($1, $2, $3)`
    _, err := r.db.ExecContext(ctx, query, u.Email, u.HashedPassword, u.Name)
    return err
}
```

## 📁 Структура папок

```
CleanArchitecture/
├── 📁 cmd/                          # Точки входа приложения
│   └── main.go                      # Главная функция
├── 📁 configs/                      # Конфигурация
│   └── config.go                    # Загрузка настроек
├── 📁 internal/                     # Внутренний код приложения
│   ├── 📁 domain/                   # 🏛️ DOMAIN LAYER
│   │   ├── 📁 user/
│   │   │   ├── entity.go            # Доменная модель User
│   │   │   └── repository.go        # Интерфейс UserRepository
│   │   ├── 📁 link/
│   │   │   ├── entity.go            # Доменная модель Link
│   │   │   └── repository.go        # Интерфейс LinkRepository
│   │   └── 📁 stat/
│   │       ├── entity.go            # Доменная модель Stat
│   │       └── repository.go        # Интерфейс StatRepository
│   ├── 📁 usecase/                  # 💼 USE CASES LAYER
│   │   ├── 📁 auth/
│   │   │   ├── interfaces.go        # Интерфейсы зависимостей
│   │   │   ├── register.go          # Use Case регистрации
│   │   │   └── login.go             # Use Case входа
│   │   ├── 📁 link/
│   │   │   ├── interfaces.go        # Интерфейсы зависимостей
│   │   │   ├── create.go            # Use Case создания ссылки
│   │   │   └── redirect.go          # Use Case редиректа
│   │   └── 📁 stat/
│   │       ├── interfaces.go        # Интерфейсы зависимостей
│   │       └── track_click.go       # Use Case отслеживания кликов
│   ├── 📁 infrastructure/           # 🌐 INFRASTRUCTURE LAYER
│   │   ├── 📁 database/
│   │   │   ├── connection.go        # Подключение к PostgreSQL
│   │   │   ├── user_repository.go   # Реализация UserRepository
│   │   │   └── link_repository.go   # Реализация LinkRepository
│   │   ├── 📁 web/
│   │   │   ├── server.go            # HTTP сервер
│   │   │   └── middleware.go        # HTTP middleware
│   │   └── 📁 external/
│   │       ├── password_hasher.go   # Реализация bcrypt
│   │       └── jwt_generator.go     # Реализация JWT
│   └── 📁 interfaces/               # 🔌 INTERFACE ADAPTERS LAYER
│       └── 📁 controllers/
│           ├── auth_controller.go   # HTTP контроллер авторизации
│           └── link_controller.go   # HTTP контроллер ссылок
├── 📁 pkg/                          # Публичные пакеты
│   └── 📁 di/
│       └── container.go             # Dependency Injection контейнер
├── 📁 migrations/                   # SQL миграции (пусто - в коде)
├── go.mod                           # Go модуль
├── go.sum                           # Зависимости
└── README.md                        # Эта документация
```

## 🔄 Поток данных

### 📤 Входящий запрос (например, создание ссылки):

```
1. 🌐 HTTP Request
   ↓
2. 🎮 LinkController.CreateLink()
   ├── Извлекает данные из HTTP
   ├── Валидирует JSON
   └── Извлекает User ID из JWT токена
   ↓
3. 💼 CreateLinkUseCase.Execute()
   ├── Валидирует URL (через URLValidator)
   ├── Проверяет безопасность URL
   ├── Генерирует короткий код (через ShortCodeGenerator)
   ├── Создает доменную модель Link
   └── Сохраняет через LinkRepository
   ↓
4. 🗄️ LinkRepository.Save()
   ├── Преобразует доменную модель в SQL
   ├── Выполняет INSERT запрос в PostgreSQL
   └── Возвращает результат
   ↓
5. 📤 HTTP Response (JSON с короткой ссылкой)
```

### 🔍 Редирект по короткой ссылке:

```
1. 🌐 GET /{shortCode}
   ↓
2. 🎮 LinkController.Redirect()
   ├── Извлекает shortCode из URL
   ├── Получает User-Agent, IP, Referer
   └── Вызывает RedirectUseCase
   ↓
3. 💼 RedirectUseCase.Execute()
   ├── Находит ссылку по коду (через LinkRepository)
   ├── Увеличивает счетчик кликов
   └── Публикует событие (через EventPublisher)
   ↓
4. 🌐 HTTP Redirect (302) на оригинальный URL
```

## 🚀 Запуск проекта

### 🐘 Настройка PostgreSQL

```bash
# Запуск PostgreSQL в Docker
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=url_shortener \
  -p 5432:5432 \
  postgres:15

# Или используйте docker-compose.yml
docker-compose up -d
```

### ⚙️ Настройка переменных окружения

```bash
# Создайте .env файл
cp .env.example .env

# Основные настройки
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=url_shortener
export DB_SSL_MODE=disable

export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080

export JWT_SECRET=your-super-secret-key-change-in-production
export JWT_ISSUER=url-shortener
export JWT_EXPIRY=24h

export BASE_URL=http://localhost:8080
export ENVIRONMENT=development
```

### 🏃‍♂️ Запуск приложения

```bash
# Установка зависимостей
go mod download

# Запуск приложения
go run cmd/main.go

# Или сборка и запуск
go build -o bin/url-shortener cmd/main.go
./bin/url-shortener
```

### 📱 Вывод при запуске:

```
🚀 Starting URL Shortener Service...
📋 Configuration loaded for environment: development
🔧 Initializing dependencies...
✅ Dependencies initialized successfully
🌐 HTTP server starting on 0.0.0.0:8080
```

## 📡 API Endpoints

### 🔐 Аутентификация

#### Регистрация пользователя
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "strongpassword",
  "name": "John Doe"
}
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1,
  "email": "user@example.com",
  "name": "John Doe"
}
```

#### Вход пользователя
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "strongpassword"
}
```

### 🔗 Работа со ссылками

#### Создание короткой ссылки
```http
POST /links
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "original_url": "https://example.com/very/long/url",
  "custom_code": "my-link"  // опционально
}
```

**Ответ:**
```json
{
  "id": 1,
  "original_url": "https://example.com/very/long/url",
  "short_code": "my-link",
  "short_url": "http://localhost:8080/my-link",
  "created_at": "2024-01-01T12:00:00Z"
}
```

#### Получение ссылок пользователя
```http
GET /links?limit=20&offset=0
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Редирект по короткой ссылке
```http
GET /{shortCode}
```
Возвращает `HTTP 302 Redirect` на оригинальный URL.

### 🏥 Health Check
```http
GET /health
```

**Ответ:**
```json
{
  "status": "ok",
  "service": "url-shortener"
}
```

## 🧪 Тестирование

### 🔧 Юнит-тесты

```bash
# Запуск всех тестов
go test ./...

# Тесты с покрытием
go test -cover ./...

# Verbose режим
go test -v ./...
```

### 📊 Тесты по слоям

```bash
# Тестирование доменного слоя
go test ./internal/domain/...

# Тестирование Use Cases
go test ./internal/usecase/...

# Тестирование контроллеров
go test ./internal/interfaces/...
```

### 🐳 Интеграционные тесты

```bash
# Запуск с тестовой БД
export ENVIRONMENT=test
export DB_NAME=url_shortener_test
go test -tags=integration ./...
```

## 🔧 Конфигурация

### 📋 Переменные окружения

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `DB_HOST` | `localhost` | Хост базы данных |
| `DB_PORT` | `5432` | Порт базы данных |
| `DB_USER` | `postgres` | Пользователь БД |
| `DB_PASSWORD` | `postgres` | Пароль БД |
| `DB_NAME` | `url_shortener` | Имя базы данных |
| `DB_SSL_MODE` | `disable` | Режим SSL для БД |
| `SERVER_HOST` | `0.0.0.0` | Хост HTTP сервера |
| `SERVER_PORT` | `8080` | Порт HTTP сервера |
| `JWT_SECRET` | ⚠️ **Обязательно** | Секретный ключ для JWT |
| `JWT_ISSUER` | `url-shortener` | Издатель JWT токенов |
| `JWT_EXPIRY` | `24h` | Время жизни JWT |
| `BASE_URL` | `http://localhost:8080` | Базовый URL для коротких ссылок |
| `ENVIRONMENT` | `development` | Окружение (dev/staging/prod) |

### 🏭 Продакшен настройки

```bash
export ENVIRONMENT=production
export JWT_SECRET=very-strong-secret-key-256-bits
export DB_SSL_MODE=require
export SERVER_READ_TIMEOUT=10s
export SERVER_WRITE_TIMEOUT=10s
export BCRYPT_COST=14
```

## 📈 Преимущества архитектуры

### ✅ Тестируемость

```go
// Легко мокать зависимости для тестов
func TestRegisterUseCase(t *testing.T) {
    mockUserRepo := &MockUserRepository{}
    mockPasswordHasher := &MockPasswordHasher{}
    mockTokenGenerator := &MockTokenGenerator{}
    
    useCase := auth.NewRegisterUseCase(
        mockUserRepo,
        mockPasswordHasher, 
        mockTokenGenerator,
    )
    
    // Тест без реальной БД, bcrypt и JWT
}
```

### 🔄 Заменяемость компонентов

```go
// Легко заменить PostgreSQL на MongoDB
userRepo := mongodb.NewUserRepository(mongoClient)

// Или заменить bcrypt на argon2
passwordHasher := external.NewArgon2PasswordHasher()

// Или добавить Redis кеш
userRepo = cache.NewCachedUserRepository(userRepo, redisClient)
```

### 🎯 Независимость от фреймворков

- **Можно заменить** `net/http` на Gin/Echo без изменения Use Cases
- **Можно заменить** PostgreSQL на MySQL без изменения бизнес-логики  
- **Можно добавить** gRPC API рядом с HTTP API
- **Можно добавить** CLI интерфейс

### 🛡️ Безопасность зависимостей

- Доменный слой **не знает** о HTTP, SQL, JWT
- Use Cases **не знают** о конкретных БД и библиотеках
- Бизнес-логика **защищена** от изменений во внешних зависимостях

### 📊 Масштабируемость

- **Микросервисы**: каждый Use Case может стать отдельным сервисом
- **CQRS**: легко разделить команды и запросы
- **Event Sourcing**: легко добавить событийную модель
- **DDD**: естественно расширяется до Domain-Driven Design

---

## 👨‍💻 Для разработчиков

### 🔍 Как читать код

1. **Начните с** `cmd/main.go` - точка входа
2. **Изучите** `pkg/di/container.go` - вирирование зависимостей  
3. **Посмотрите** `internal/domain/` - доменные модели
4. **Изучите** `internal/usecase/` - бизнес-логику
5. **Посмотрите** `internal/interfaces/controllers/` - HTTP API

### 🛠️ Как добавить новую функциональность

1. **Создайте** доменную модель в `internal/domain/`
2. **Определите** интерфейс репозитория
3. **Создайте** Use Case в `internal/usecase/`
4. **Реализуйте** репозиторий в `internal/infrastructure/database/`
5. **Создайте** контроллер в `internal/interfaces/controllers/`
6. **Зарегистрируйте** в DI контейнере

### 📝 Стиль кода

- ✅ **Подробные комментарии** объясняют принципы архитектуры
- ✅ **Константы ошибок** определены в доменном слое
- ✅ **Валидация данных** на каждом уровне
- ✅ **Graceful shutdown** для продакшена
- ✅ **Structured logging** (можно добавить)

---

**🎉 Поздравляем! Вы изучили полную реализацию Clean Architecture на Go!**

Этот проект демонстрирует все ключевые принципы и может служить основой для реальных продакшен приложений.